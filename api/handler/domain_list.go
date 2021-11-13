package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func domainList(ownerHex string) ([]model.Domain, error) {
	if ownerHex == "" {
		return []model.Domain{}, app.ErrorMissingField
	}

	statement := `
		SELECT ` + domainsRowColumns + `
		FROM domains
		WHERE ownerHex=$1;
	`
	rows, err := repository.Db.Query(statement, ownerHex)
	if err != nil {
		util.GetLogger().Errorf("cannot query domains: %v", err)
		return nil, app.ErrorInternal
	}
	defer rows.Close()

	domains := []model.Domain{}
	for rows.Next() {
		var d model.Domain
		if err = domainsRowScan(rows, &d); err != nil {
			util.GetLogger().Errorf("cannot Scan domain: %v", err)
			return nil, app.ErrorInternal
		}

		d.Moderators, err = domainModeratorList(d.Domain)
		if err != nil {
			return []model.Domain{}, err
		}

		domains = append(domains, d)
	}

	return domains, rows.Err()
}

func DomainListHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
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

	domains, err := domainList(o.OwnerHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{
		"success": true,
		"domains": domains,
		"configuredOauths": map[string]bool{
			"google":  googleConfigured,
			"twitter": twitterConfigured,
			"github":  githubConfigured,
			"gitlab":  gitlabConfigured,
		},
	})
}
