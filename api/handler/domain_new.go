package handler

import (
	"net/http"
	"os"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
	"strings"
)

func domainNew(ownerHex string, name string, domain string) error {
	if ownerHex == "" || name == "" || domain == "" {
		return app.ErrorMissingField
	}

	if strings.Contains(domain, "/") {
		return app.ErrorInvalidDomain
	}

	// if asked to disable wildcards, then don't allow them...
	if os.Getenv("ENABLE_WILDCARDS") == "false" {
		if strings.Contains(domain, "%") || strings.Contains(domain, "_") {
			return app.ErrorInvalidDomain
		}
	}

	return repository.Repo.DomainRepository.CreateDomain(ownerHex, name, domain)
}

func DomainNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Name       *string `json:"name"`
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

	if err = domainNew(o.OwnerHex, *x.Name, domain); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err = domainModeratorNew(domain, o.Email); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "domain": domain})
}
