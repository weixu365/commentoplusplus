package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func pageUpdate(p model.Page) error {
	if p.Domain == "" {
		return app.ErrorMissingField
	}

	// fields to not update:
	//   commentCount
	statement := `
		INSERT INTO
		pages  (domain, path, isLocked, stickyCommentHex)
		VALUES ($1,     $2,   $3,       $4              )
		ON CONFLICT (domain, path) DO
			UPDATE SET isLocked = $3, stickyCommentHex = $4;
	`
	_, err := repository.Db.Exec(statement, p.Domain, p.Path, p.IsLocked, p.StickyCommentHex)
	if err != nil {
		util.GetLogger().Errorf("error setting page attributes: %v", err)
		return app.ErrorInternal
	}

	return nil
}

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

	if err = pageUpdate(*x.Attributes); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
