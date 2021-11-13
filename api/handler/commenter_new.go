package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func commenterNew(email string, name string, link string, photo string, provider string, password string) (string, error) {
	if email == "" || name == "" || link == "" || photo == "" || provider == "" {
		return "", app.ErrorMissingField
	}

	if provider == "commento" && password == "" {
		return "", app.ErrorMissingField
	}

	// See utils_sanitise.go's documentation on IsHttpsUrl. This is not a URL
	// validator, just an XSS preventor.
	// TODO: reject URLs instead of malforming them.
	if link != "undefined" && !util.IsHttpsUrl(link) {
		link = "https://" + link
	}

	if provider != "anon" {
		if _, err := commenterGetByEmail(provider, email); err == nil {
			return "", app.ErrorEmailAlreadyExists
		}

		if err := repository.EmailNew(email); err != nil {
			return "", app.ErrorInternal
		}
	}

	commenterHex, err := util.RandomHex(32)
	if err != nil {
		return "", app.ErrorInternal
	}

	passwordHash := []byte{}
	if password != "" {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			util.GetLogger().Errorf("cannot generate hash from password: %v\n", err)
			return "", app.ErrorInternal
		}
	}
	if provider == "anon" {
		passwordHash = []byte{}
	}

	statement := `
		INSERT INTO
		commenters (commenterHex, email, name, link, photo, provider, passwordHash, joinDate)
		VALUES     ($1,           $2,    $3,   $4,   $5,    $6,       $7,           $8      );
	`
	_, err = repository.Db.Exec(statement, commenterHex, email, name, link, photo, provider, string(passwordHash), time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert commenter: %v", err)
		return "", app.ErrorInternal
	}

	return commenterHex, nil
}

func commenterNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Name     *string `json:"name"`
		Website  *string `json:"website"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// TODO: add gravatar?
	// TODO: email confirmation if provider = commento?
	// TODO: email confirmation if provider = commento?
	if *x.Website == "" {
		*x.Website = "undefined"
	}

	if _, err := commenterNew(*x.Email, *x.Name, *x.Website, "undefined", "commento", *x.Password); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "confirmEmail": notification.SmtpConfigured})
}
