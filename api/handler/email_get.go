package handler

import (
	"net/http"
	"simple-commenting/repository"
)

func EmailGetHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UnsubscribeSecretHex *string `json:"unsubscribeSecretHex"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	email, err := repository.Repo.EmailRepository.GetByUnsubscribeSecretHex(*x.UnsubscribeSecretHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "email": email})
}
