package handler

import (
	"net/http"
	"simple-commenting/repository"
)

func CommenterSelfHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	commenter, err := repository.Repo.CommenterRepository.GetCommenterByToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	email, err := repository.Repo.EmailRepository.GetEmail(commenter.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commenter": commenter, "email": email})
}
