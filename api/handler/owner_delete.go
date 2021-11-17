package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func ownerDelete(ownerHex string, deleteDomains bool) error {
	domains, err := repository.Repo.DomainRepository.ListDomain(ownerHex)
	if err != nil {
		return err
	}

	if len(domains) > 0 {
		if !deleteDomains {
			return app.ErrorCannotDeleteOwnerWithActiveDomains
		}
		for _, d := range domains {
			if err := domainDelete(d.Domain); err != nil {
				return err
			}
		}
	}

	statement := `
		DELETE FROM owners
		WHERE ownerHex = $1;
	`
	_, err = repository.Db.Exec(statement, ownerHex)
	if err != nil {
		return app.ErrorNoSuchOwner
	}

	statement = `
		DELETE FROM ownersessions
		WHERE ownerHex = $1;
	`
	_, err = repository.Db.Exec(statement, ownerHex)
	if err != nil {
		util.GetLogger().Errorf("cannot delete from ownersessions: %v", err)
		return app.ErrorInternal
	}

	statement = `
		DELETE FROM resethexes
		WHERE hex = $1;
	`
	_, err = repository.Db.Exec(statement, ownerHex)
	if err != nil {
		util.GetLogger().Errorf("cannot delete from resethexes: %v", err)
		return app.ErrorInternal
	}

	return nil
}

func OwnerDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = ownerDelete(o.OwnerHex, false); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
