package handler

import (
	"database/sql"
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func commentList(commenterHex string, domain string, path string, includeUnapproved bool) ([]*model.Comment, map[string]*model.Commenter, error) {
	// path can be empty
	if commenterHex == "" || domain == "" {
		return nil, nil, app.ErrorMissingField
	}

	statement := `
		SELECT
			commentHex,
			commenterHex,
			markdown,
			html,
			parentHex,
			score,
			state,
			deleted,
			creationDate
		FROM comments
		WHERE
			canon($1) LIKE canon(comments.domain) AND
			comments.path = $2
	`

	if !includeUnapproved {
		if commenterHex == "anonymous" {
			statement += `AND state = 'approved'`
		} else {
			statement += `AND (state = 'approved' OR commenterHex = $3)`
		}
	}

	statement += `;`

	var rows *sql.Rows
	var err error

	if !includeUnapproved && commenterHex != "anonymous" {
		rows, err = repository.Db.Query(statement, domain, path, commenterHex)
	} else {
		rows, err = repository.Db.Query(statement, domain, path)
	}

	if err != nil {
		util.GetLogger().Errorf("cannot get comments: %v", err)
		return nil, nil, app.ErrorInternal
	}
	defer rows.Close()

	commenters := make(map[string]*model.Commenter)
	commenters["anonymous"] = &model.Commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []*model.Comment{}
	for rows.Next() {
		comment := model.Comment{}
		if err = rows.Scan(
			&comment.CommentHex,
			&comment.CommenterHex,
			&comment.Markdown,
			&comment.Html,
			&comment.ParentHex,
			&comment.Score,
			&comment.State,
			&comment.Deleted,
			&comment.CreationDate); err != nil {
			return nil, nil, app.ErrorInternal
		}

		if commenterHex != "anonymous" {
			statement = `
				SELECT direction
				FROM votes
				WHERE commentHex=$1 AND commenterHex=$2;
			`
			row := repository.Db.QueryRow(statement, comment.CommentHex, commenterHex)

			if err = row.Scan(&comment.Direction); err != nil {
				// TODO: is the only error here that there is no such entry?
				comment.Direction = 0
			}
		}

		if commenterHex != comment.CommenterHex {
			comment.Markdown = ""
		}

		if !includeUnapproved {
			comment.State = ""
		}

		comments = append(comments, &comment)

		if _, ok := commenters[comment.CommenterHex]; !ok {
			commenters[comment.CommenterHex], err = repository.Repo.CommenterRepository.GetCommenterByHex(comment.CommenterHex)
			if err != nil {
				util.GetLogger().Errorf("cannot retrieve commenter: %v", err)
				return nil, nil, app.ErrorInternal
			}
		}
	}

	return comments, commenters, nil
}

func CommentListHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"CommenterToken"`
		Domain         *string `json:"domain"`
		Path           *string `json:"path"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domainName := util.DomainStrip(*x.Domain)
	path := *x.Path

	domain, err := domainGet(domainName)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	page, err := repository.Repo.PageRepository.GetPageByPath(domainName, path)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	commenterHex := "anonymous"
	isModerator := false
	modList := map[string]bool{}

	if *x.CommenterToken != "anonymous" {
		commenter, err := repository.Repo.CommenterRepository.GetCommenterByToken(*x.CommenterToken)
		if err != nil {
			if err == app.ErrorNoSuchToken {
				commenterHex = "anonymous"
			} else {
				bodyMarshal(w, response{"success": false, "message": err.Error()})
				return
			}
		} else {
			commenterHex = commenter.CommenterHex
		}

		for _, mod := range *domain.Moderators {
			modList[mod.Email] = true
			if mod.Email == commenter.Email {
				isModerator = true
			}
		}
	} else {
		for _, mod := range *domain.Moderators {
			modList[mod.Email] = true
		}
	}

	repository.Repo.LogRepository.LogDomainViewRecord(domainName, commenterHex)

	comments, commenters, err := commentList(commenterHex, domainName, path, isModerator)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	_commenters := map[string]*model.Commenter{}
	for commenterHex, cr := range commenters {
		if _, ok := modList[cr.Email]; ok {
			cr.IsModerator = true
		}
		cr.Email = ""
		_commenters[commenterHex] = cr
	}

	bodyMarshal(w, response{
		"success":               true,
		"domain":                domainName,
		"comments":              comments,
		"commenters":            _commenters,
		"requireModeration":     domain.RequireModeration,
		"requireIdentification": domain.RequireIdentification,
		"isFrozen":              domain.State == "frozen",
		"isModerator":           isModerator,
		"defaultSortPolicy":     domain.DefaultSortPolicy,
		"attributes":            page,
		"configuredOauths": map[string]bool{
			"commento": domain.CommentoProvider,
			"google":   googleConfigured && domain.GoogleProvider,
			"twitter":  twitterConfigured && domain.TwitterProvider,
			"github":   githubConfigured && domain.GithubProvider,
			"gitlab":   gitlabConfigured && domain.GitlabProvider,
			"sso":      domain.SsoProvider,
		},
	})
}

func commentListApprovals(domain string) ([]*model.Comment, map[string]*model.Commenter, error) {
	if domain == "" {
		return nil, nil, app.ErrorMissingField
	}

	statement := `
		SELECT
			path,
			commentHex,
			commenterHex,
			markdown,
			html,
			parentHex,
			score,
			state,
			deleted,
			creationDate
		FROM comments
		WHERE
		canon(comments.domain) LIKE canon($1) AND deleted = false AND
			( state = 'unapproved' OR state = 'flagged' );
	`

	var rows *sql.Rows
	var err error

	rows, err = repository.Db.Query(statement, domain)

	if err != nil {
		util.GetLogger().Errorf("cannot get comments: %v", err)
		return nil, nil, app.ErrorInternal
	}
	defer rows.Close()

	commenters := make(map[string]*model.Commenter)
	commenters["anonymous"] = &model.Commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []*model.Comment{}
	for rows.Next() {
		comment := model.Comment{}
		if err = rows.Scan(
			&comment.Path,
			&comment.CommentHex,
			&comment.CommenterHex,
			&comment.Markdown,
			&comment.Html,
			&comment.ParentHex,
			&comment.Score,
			&comment.State,
			&comment.Deleted,
			&comment.CreationDate); err != nil {
			return nil, nil, app.ErrorInternal
		}

		comments = append(comments, &comment)

		if _, ok := commenters[comment.CommenterHex]; !ok {
			commenters[comment.CommenterHex], err = repository.Repo.CommenterRepository.GetCommenterByHex(comment.CommenterHex)
			if err != nil {
				util.GetLogger().Errorf("cannot retrieve commenter: %v", err)
				return nil, nil, app.ErrorInternal
			}
		}
	}

	return comments, commenters, nil
}

func CommentListApprovalsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := util.DomainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotAuthorised.Error()})
		return
	}

	comments, commenters, err := commentListApprovals(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	_commenters := map[string]*model.Commenter{}
	for commenterHex, cr := range commenters {
		cr.Email = ""
		_commenters[commenterHex] = cr
	}

	bodyMarshal(w, response{
		"success":    true,
		"domain":     domain,
		"comments":   comments,
		"commenters": _commenters,
	})

}

func commentListAll(domain string) ([]*model.Comment, map[string]*model.Commenter, error) {
	if domain == "" {
		return nil, nil, app.ErrorMissingField
	}

	statement := `
		SELECT
			path,
			commentHex,
			commenterHex,
			markdown,
			html,
			parentHex,
			score,
			state,
			deleted,
			creationDate
		FROM comments
		WHERE
		canon(comments.domain) LIKE canon($1) AND deleted = false AND 
			( state = 'approved'  );
	`

	var rows *sql.Rows
	var err error

	rows, err = repository.Db.Query(statement, domain)

	if err != nil {
		util.GetLogger().Errorf("cannot get comments: %v", err)
		return nil, nil, app.ErrorInternal
	}
	defer rows.Close()

	commenters := make(map[string]*model.Commenter)
	commenters["anonymous"] = &model.Commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []*model.Comment{}
	for rows.Next() {
		comment := model.Comment{}
		if err = rows.Scan(
			&comment.Path,
			&comment.CommentHex,
			&comment.CommenterHex,
			&comment.Markdown,
			&comment.Html,
			&comment.ParentHex,
			&comment.Score,
			&comment.State,
			&comment.Deleted,
			&comment.CreationDate); err != nil {
			return nil, nil, app.ErrorInternal
		}

		comments = append(comments, &comment)

		if _, ok := commenters[comment.CommenterHex]; !ok {
			commenters[comment.CommenterHex], err = repository.Repo.CommenterRepository.GetCommenterByHex(comment.CommenterHex)
			if err != nil {
				util.GetLogger().Errorf("cannot retrieve commenter: %v", err)
				return nil, nil, app.ErrorInternal
			}
		}
	}

	return comments, commenters, nil
}

func CommentListAllHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := util.DomainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotAuthorised.Error()})
		return
	}

	comments, commenters, err := commentListAll(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	_commenters := map[string]model.Commenter{}
	for commenterHex, cr := range commenters {
		cr.Email = ""
		_commenters[commenterHex] = cr
	}

	bodyMarshal(w, response{
		"success":    true,
		"domain":     domain,
		"comments":   comments,
		"commenters": _commenters,
	})

}
