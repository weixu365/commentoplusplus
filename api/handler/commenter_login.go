package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func commenterLogin(email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", app.ErrorMissingField
	}

	statement := `
		SELECT commenterHex, passwordHash
		FROM commenters
		WHERE email = $1 AND provider = 'commento' AND deleted=false;
	`
	row := repository.Db.QueryRow(statement, email)

	var commenterHex string
	var passwordHash string
	if err := row.Scan(&commenterHex, &passwordHash); err != nil {
		return "", app.ErrorInvalidEmailPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		// TODO: is this the only possible error?
		return "", app.ErrorInvalidEmailPassword
	}

	commenterToken, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("cannot create commenterToken: %v", err)
		return "", app.ErrorInternal
	}

	statement = `
		INSERT INTO
		commenterSessions (commenterToken, commenterHex, creationDate)
		VALUES            ($1,             $2,           $3          );
	`
	_, err = repository.Db.Exec(statement, commenterToken, commenterHex, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert commenterToken token: %v\n", err)
		return "", app.ErrorInternal
	}

	return commenterToken, nil
}

func CommenterLoginHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	commenterToken, err := commenterLogin(*x.Email, *x.Password)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// TODO: modify commenterLogin to directly return commenter?
	commenter, err := repository.Repo.CommenterRepository.GetCommenterByToken(commenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	e, err := emailGet(commenter.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commenterToken": commenterToken, "commenter": commenter, "email": e})
}
