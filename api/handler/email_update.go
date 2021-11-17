package handler

import (
	"net/http"
	"simple-commenting/model"
	"simple-commenting/repository"
)

func EmailUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email *model.Email `json:"email"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := repository.Repo.EmailRepository.UpdateEmail(x.Email); err != nil {
		bodyMarshal(w, response{"success": true, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
