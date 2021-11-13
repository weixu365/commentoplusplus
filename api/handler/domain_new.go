package handler

import (
	"net/http"
	"os"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
	"strings"
	"time"
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

	// test if domain already exists
	statement := `
		SELECT COUNT(*) FROM
		domains WHERE
		canon(regexp_replace($1, '[%]', '')) LIKE canon(domain) OR canon(domain) LIKE canon($1);
	`
	row := repository.Db.QueryRow(statement, domain)
	var err error
	var count int

	if err = row.Scan(&count); err != nil {
		return app.ErrorInvalidDomain
	}

	if count > 0 {
		return app.ErrorDomainAlreadyExists
	}

	statement = `
		INSERT INTO
		domains (ownerHex, name, domain, creationDate)
		VALUES  ($1,       $2,   $3,     $4          );
	`
	_, err = repository.Db.Exec(statement, ownerHex, name, domain, time.Now().UTC())
	if err != nil {
		// TODO: This should not happen given the above check, so this is likely not the error. Be more informative?
		return app.ErrorDomainAlreadyExists
	}

	return nil
}

func domainNewHandler(w http.ResponseWriter, r *http.Request) {
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
