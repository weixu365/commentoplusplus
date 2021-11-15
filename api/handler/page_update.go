package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func PageUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string     `json:"commenterToken"`
		Domain         *string     `json:"domain"`
		Path           *string     `json:"path"`
		Attributes     *model.Page `json:"attributes"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := repository.Repo.CommenterRepository.GetCommenterByToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := util.DomainStrip(*x.Domain)

	isModerator, err := isDomainModerator(domain, c.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isModerator {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotModerator.Error()})
		return
	}

	(*x.Attributes).Domain = *x.Domain
	(*x.Attributes).Path = *x.Path

	if err = repository.Repo.PageRepository.UpdatePage(x.Attributes); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
