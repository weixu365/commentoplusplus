package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func domainDelete(domain string) error {
	if domain == "" {
		return app.ErrorMissingField
	}

	statement := `
		DELETE FROM domains
		WHERE domain = $1;
	`
	_, err := repository.Db.Exec(statement, domain)
	if err != nil {
		return app.ErrorNoSuchDomain
	}

	statement = `
		DELETE FROM views
		WHERE views.domain = $1;
	`
	_, err = repository.Db.Exec(statement, domain)
	if err != nil {
		util.GetLogger().Errorf("cannot delete domain from views: %v", err)
		return app.ErrorInternal
	}

	statement = `
		DELETE FROM moderators
		WHERE moderators.domain = $1;
	`
	_, err = repository.Db.Exec(statement, domain)
	if err != nil {
		util.GetLogger().Errorf("cannot delete domain from moderators: %v", err)
		return app.ErrorInternal
	}

	statement = `
		DELETE FROM ssotokens
		WHERE ssotokens.domain = $1;
	`
	_, err = repository.Db.Exec(statement, domain)
	if err != nil {
		util.GetLogger().Errorf("cannot delete domain from ssotokens: %v", err)
		return app.ErrorInternal
	}

	// comments, votes, and pages are handled by domainClear
	if err = domainClear(domain); err != nil {
		util.GetLogger().Errorf("cannot clear domain: %v", err)
		return app.ErrorInternal
	}

	return nil
}

func DomainDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = domainDelete(*x.Domain); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
