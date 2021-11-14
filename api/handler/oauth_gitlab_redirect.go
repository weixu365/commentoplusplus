package handler

import (
	"fmt"
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func GitlabRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if gitlabConfig == nil {
		util.GetLogger().Errorf("gitlab oauth access attempt without configuration")
		fmt.Fprintf(w, "error: this website has not configured gitlab OAuth")
		return
	}

	commenterToken := r.FormValue("commenterToken")

	_, err := repository.Repo.CommenterRepository.GetCommenterByToken(commenterToken)
	if err != nil && err != app.ErrorNoSuchToken {
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	url := gitlabConfig.AuthCodeURL(commenterToken)
	http.Redirect(w, r, url, http.StatusFound)
}
