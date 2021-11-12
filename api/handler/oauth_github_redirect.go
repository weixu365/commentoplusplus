package handler

import (
	"fmt"
	"net/http"
	"simple-commenting/app"
	"simple-commenting/util"
)

func githubRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if githubConfig == nil {
		util.GetLogger().Errorf("github oauth access attempt without configuration")
		fmt.Fprintf(w, "error: this website has not configured github OAuth")
		return
	}

	commenterToken := r.FormValue("commenterToken")

	_, err := commenterGetByCommenterToken(commenterToken)
	if err != nil && err != app.ErrorNoSuchToken {
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	url := githubConfig.AuthCodeURL(commenterToken)
	http.Redirect(w, r, url, http.StatusFound)
}
