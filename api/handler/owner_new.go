package handler

import (
	"net/http"
	"os"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"

	"golang.org/x/crypto/bcrypt"
)

func ownerNew(email string, name string, password string) (string, error) {
	if email == "" || name == "" || password == "" {
		return "", app.ErrorMissingField
	}

	if os.Getenv("FORBID_NEW_OWNERS") == "true" {
		return "", app.ErrorNewOwnerForbidden
	}

	if _, err := repository.Repo.OwnerRepository.GetByEmail(email); err == nil {
		return "", app.ErrorEmailAlreadyExists
	}

	if err := repository.Repo.EmailRepository.CreateEmail(email); err != nil {
		return "", app.ErrorInternal
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		util.GetLogger().Errorf("cannot generate hash from password: %v\n", err)
		return "", app.ErrorInternal
	}

	owner := model.Owner{
		Email:          email,
		Name:           name,
		PasswordHash:   string(passwordHash),
		ConfirmedEmail: !notification.SmtpConfigured,
	}

	createdOwner, err := repository.Repo.OwnerRepository.CreateOwner(&owner)
	if err != nil {
		util.GetLogger().Errorf("cannot generate hash from password: %v\n", err)
		return "", err
	}

	if notification.SmtpConfigured {
		confirmHex, err := repository.Repo.OwnerRepository.CreateOwnerConfirmHex(createdOwner.OwnerHex)
		if err = notification.SmtpOwnerConfirmHex(email, name, confirmHex); err != nil {
			return "", err
		}
	}

	return createdOwner.OwnerHex, nil
}

func OwnerNewHandler(w http.ResponseWriter, r *http.Request) {
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

	bodyMarshal(w, response{"success": true, "confirmEmail": notification.SmtpConfigured})
}
