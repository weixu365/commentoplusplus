package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"

	"golang.org/x/crypto/bcrypt"
)

func reset(resetHex string, password string) (string, error) {
	if resetHex == "" || password == "" {
		return "", app.ErrorMissingField
	}

	statement := `
		SELECT hex, entity
		FROM resetHexes
		WHERE resetHex = $1;
	`
	row := repository.Db.QueryRow(statement, resetHex)

	var hex string
	var entity string
	if err := row.Scan(&hex, &entity); err != nil {
		// TODO: is this the only error?
		return "", app.ErrorNoSuchResetToken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		util.GetLogger().Errorf("cannot generate hash from password: %v\n", err)
		return "", app.ErrorInternal
	}

	if entity == "owner" {
		statement = `
			UPDATE owners SET passwordHash = $1, confirmedEmail=true
			WHERE ownerHex = $2;
		`
	} else {
		statement = `
			UPDATE commenters SET passwordHash = $1
			WHERE commenterHex = $2;
		`
	}

	_, err = repository.Db.Exec(statement, string(passwordHash), hex)
	if err != nil {
		util.GetLogger().Errorf("cannot change %s's password: %v\n", entity, err)
		return "", app.ErrorInternal
	}

	statement = `
		DELETE FROM resetHexes
		WHERE resetHex = $1;
	`
	_, err = repository.Db.Exec(statement, resetHex)
	if err != nil {
		util.GetLogger().Warningf("cannot remove resetHex: %v\n", err)
	}

	return entity, nil
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ResetHex *string `json:"resetHex"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	entity, err := reset(*x.ResetHex, *x.Password)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "entity": entity})
}
