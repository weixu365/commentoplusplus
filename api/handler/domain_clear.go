package handler

import (
	"net/http"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func domainClear(domain string) error {
	if domain == "" {
		return app.ErrorMissingField
	}

	statement := `
		DELETE FROM votes
		USING comments
		WHERE comments.commentHex = votes.commentHex AND canon($1) LIKE canon(comments.domain);
	`
	_, err := repository.Db.Exec(statement, domain)
	if err != nil {
		util.GetLogger().Errorf("cannot delete votes: %v", err)
		return app.ErrorInternal
	}

	statement = `
		DELETE FROM comments
		WHERE comments.domain = $1;
	`
	_, err = repository.Db.Exec(statement, domain)
	if err != nil {
		util.GetLogger().Errorf(statement, domain)
		return app.ErrorInternal
	}

	statement = `
		DELETE FROM pages
		WHERE pages.domain = $1;
	`
	_, err = repository.Db.Exec(statement, domain)
	if err != nil {
		util.GetLogger().Errorf(statement, domain)
		return app.ErrorInternal
	}

	return nil
}

func domainClearHandler(w http.ResponseWriter, r *http.Request) {
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

	domain := domainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err = domainClear(*x.Domain); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
