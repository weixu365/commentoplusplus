package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func ownerLogin(email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", app.ErrorMissingField
	}

	statement := `
		SELECT ownerHex, confirmedEmail, passwordHash
		FROM owners
		WHERE email=$1;
	`
	row := repository.Db.QueryRow(statement, email)

	var ownerHex string
	var confirmedEmail bool
	var passwordHash string
	if err := row.Scan(&ownerHex, &confirmedEmail, &passwordHash); err != nil {
		return "", app.ErrorInvalidEmailPassword
	}

	if !confirmedEmail {
		return "", app.ErrorUnconfirmedEmail
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		// TODO: is this the only possible error?
		return "", app.ErrorInvalidEmailPassword
	}

	ownerToken, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("cannot create ownerToken: %v", err)
		return "", app.ErrorInternal
	}

	statement = `
		INSERT INTO
		ownerSessions (ownerToken, ownerHex, loginDate)
		VALUES        ($1,         $2,       $3       );
	`
	_, err = repository.Db.Exec(statement, ownerToken, ownerHex, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert ownerSession: %v\n", err)
		return "", app.ErrorInternal
	}

	return ownerToken, nil
}

func OwnerLoginHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	ownerToken, err := ownerLogin(*x.Email, *x.Password)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "ownerToken": ownerToken})
}
