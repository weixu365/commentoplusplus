package handler

import (
	"net/http"
	"simple-commenting/repository"
)

func CommenterTokenNewHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken, err := repository.Repo.CommenterRepository.CreateCommenterToken()
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commenterToken": commenterToken})
}
