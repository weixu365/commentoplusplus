package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

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

	comments, commenters, err := repository.Repo.CommentRepository.ListComments(commenterHex, domainName, path, isModerator)
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

	comments, commenters, err := repository.Repo.CommentRepository.ListUnapprovedComments(domain)
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

	comments, commenters, err := repository.Repo.CommentRepository.ListApprovedComments(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	_commenters := map[string]*model.Commenter{}
	for commenterHex, commenter := range commenters {
		commenter.Email = ""
		_commenters[commenterHex] = commenter
	}

	bodyMarshal(w, response{
		"success":    true,
		"domain":     domain,
		"comments":   comments,
		"commenters": _commenters,
	})

}
