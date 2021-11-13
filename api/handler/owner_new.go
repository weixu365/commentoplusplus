package handler

import (
	"net/http"
	"os"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func ownerNew(email string, name string, password string) (string, error) {
	if email == "" || name == "" || password == "" {
		return "", app.ErrorMissingField
	}

	if os.Getenv("FORBID_NEW_OWNERS") == "true" {
		return "", app.ErrorNewOwnerForbidden
	}

	if _, err := ownerGetByEmail(email); err == nil {
		return "", app.ErrorrrorEmailAlreadyExists
	}

	if err := EmailNew(email); err != nil {
		return "", app.ErrorInternal
	}

	ownerHex, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("cannot generate ownerHex: %v", err)
		return "", app.ErrorInternal
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		util.GetLogger().Errorf("cannot generate hash from password: %v\n", err)
		return "", app.ErrorInternal
	}

	statement := `
		INSERT INTO
		owners (ownerHex, email, name, passwordHash, joinDate, confirmedEmail)
		VALUES ($1,       $2,    $3,   $4,           $5,       $6            );
	`
	_, err = repository.Db.Exec(statement, ownerHex, email, name, string(passwordHash), time.Now().UTC(), !smtpConfigured)
	if err != nil {
		// TODO: Make sure `err` is actually about conflicting UNIQUE, and not some
		// other error. If it is something else, we should probably return `errorInternal`.
		return "", app.ErrorEmailAlreadyExists
	}

	if smtpConfigured {
		confirmHex, err := util.RandomHex(32)
		if err != nil {
			util.GetLogger().Errorf("cannot generate confirmHex: %v", err)
			return "", app.ErrorInternal
		}

		statement = `
			INSERT INTO
			ownerConfirmHexes (confirmHex, ownerHex, sendDate)
			VALUES            ($1,         $2,       $3      );
		`
		_, err = repository.Db.Exec(statement, confirmHex, ownerHex, time.Now().UTC())
		if err != nil {
			util.GetLogger().Errorf("cannot insert confirmHex: %v\n", err)
			return "", app.ErrorInternal
		}

		if err = smtpOwnerConfirmHex(email, name, confirmHex); err != nil {
			return "", err
		}
	}

	return ownerHex, nil
}

func ownerNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Name     *string `json:"name"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if _, err := ownerNew(*x.Email, *x.Name, *x.Password); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// Errors in creating a commenter account should not hold this up.
	_, _ = commenterNew(*x.Email, *x.Name, "undefined", "undefined", "commento", *x.Password)

	bodyMarshal(w, response{"success": true, "confirmEmail": smtpConfigured})
}
