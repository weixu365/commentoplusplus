package handler

import (
	"net/http"
	"simple-commenting/app"
)

func ownerSelfHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err == app.ErrorNoSuchToken {
		bodyMarshal(w, response{"success": true, "loggedIn": false})
		return
	}

	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "loggedIn": true, "owner": o})
}
