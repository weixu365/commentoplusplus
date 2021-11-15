package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"

	"golang.org/x/crypto/bcrypt"
)

func commenterLogin(email string, password string) (string, error) {
	commenter, err := repository.Repo.CommenterRepository.GetActiveCommenterByEmail(email)
	if err != nil {
		return "", app.ErrorInvalidEmailPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(commenter.PasswordHash), []byte(password)); err != nil {
		// TODO: is this the only possible error?
		return "", app.ErrorInvalidEmailPassword
	}

	commenterToken, err := repository.Repo.CommenterRepository.CreateCommenterSessionToken(commenter.CommenterHex)
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

	email, err := repository.Repo.EmailRepository.GetEmail(commenter.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commenterToken": commenterToken, "commenter": commenter, "email": email})
}
