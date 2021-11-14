package handler

import (
	"database/sql"
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func commentList(commenterHex string, domain string, path string, includeUnapproved bool) ([]model.Comment, map[string]model.Commenter, error) {
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

	commenters := make(map[string]model.Commenter)
	commenters["anonymous"] = model.Commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []model.Comment{}
	for rows.Next() {
		c := model.Comment{}
		if err = rows.Scan(
			&c.CommentHex,
			&c.CommenterHex,
			&c.Markdown,
			&c.Html,
			&c.ParentHex,
			&c.Score,
			&c.State,
			&c.Deleted,
			&c.CreationDate); err != nil {
			return nil, nil, app.ErrorInternal
		}

		if commenterHex != "anonymous" {
			statement = `
				SELECT direction
				FROM votes
				WHERE commentHex=$1 AND commenterHex=$2;
			`
			row := repository.Db.QueryRow(statement, c.CommentHex, commenterHex)

			if err = row.Scan(&c.Direction); err != nil {
				// TODO: is the only error here that there is no such entry?
				c.Direction = 0
			}
		}

		if commenterHex != c.CommenterHex {
			c.Markdown = ""
		}

		if !includeUnapproved {
			c.State = ""
		}

		comments = append(comments, c)

		if _, ok := commenters[c.CommenterHex]; !ok {
			commenters[c.CommenterHex], err = commenterGetByHex(c.CommenterHex)
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

	domain := util.DomainStrip(*x.Domain)
	path := *x.Path

	d, err := domainGet(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	p, err := repository.Repo.PageRepository.GetPageByPath(domain, path)
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

		for _, mod := range *d.Moderators {
			modList[mod.Email] = true
			if mod.Email == commenter.Email {
				isModerator = true
			}
		}
	} else {
		for _, mod := range *d.Moderators {
			modList[mod.Email] = true
		}
	}

	domainViewRecord(domain, commenterHex)

	comments, commenters, err := commentList(commenterHex, domain, path, isModerator)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	_commenters := map[string]model.Commenter{}
	for commenterHex, cr := range commenters {
		if _, ok := modList[cr.Email]; ok {
			cr.IsModerator = true
		}
		cr.Email = ""
		_commenters[commenterHex] = cr
	}

	bodyMarshal(w, response{
		"success":               true,
		"domain":                domain,
		"comments":              comments,
		"commenters":            _commenters,
		"requireModeration":     d.RequireModeration,
		"requireIdentification": d.RequireIdentification,
		"isFrozen":              d.State == "frozen",
		"isModerator":           isModerator,
		"defaultSortPolicy":     d.DefaultSortPolicy,
		"attributes":            p,
		"configuredOauths": map[string]bool{
			"commento": d.CommentoProvider,
			"google":   googleConfigured && d.GoogleProvider,
			"twitter":  twitterConfigured && d.TwitterProvider,
			"github":   githubConfigured && d.GithubProvider,
			"gitlab":   gitlabConfigured && d.GitlabProvider,
			"sso":      d.SsoProvider,
		},
	})
}

func commentListApprovals(domain string) ([]model.Comment, map[string]model.Commenter, error) {
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

	commenters := make(map[string]model.Commenter)
	commenters["anonymous"] = model.Commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []model.Comment{}
	for rows.Next() {
		c := model.Comment{}
		if err = rows.Scan(
			&c.Path,
			&c.CommentHex,
			&c.CommenterHex,
			&c.Markdown,
			&c.Html,
			&c.ParentHex,
			&c.Score,
			&c.State,
			&c.Deleted,
			&c.CreationDate); err != nil {
			return nil, nil, app.ErrorInternal
		}

		comments = append(comments, c)

		if _, ok := commenters[c.CommenterHex]; !ok {
			commenters[c.CommenterHex], err = commenterGetByHex(c.CommenterHex)
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

func commentListAll(domain string) ([]model.Comment, map[string]model.Commenter, error) {
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

	commenters := make(map[string]model.Commenter)
	commenters["anonymous"] = model.Commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []model.Comment{}
	for rows.Next() {
		c := model.Comment{}
		if err = rows.Scan(
			&c.Path,
			&c.CommentHex,
			&c.CommenterHex,
			&c.Markdown,
			&c.Html,
			&c.ParentHex,
			&c.Score,
			&c.State,
			&c.Deleted,
			&c.CreationDate); err != nil {
			return nil, nil, app.ErrorInternal
		}

		comments = append(comments, c)

		if _, ok := commenters[c.CommenterHex]; !ok {
			commenters[c.CommenterHex], err = commenterGetByHex(c.CommenterHex)
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
