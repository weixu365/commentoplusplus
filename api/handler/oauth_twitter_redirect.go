package handler

import (
	"fmt"
	"net/http"
	"os"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func TwitterRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if twitterClient == nil {
		util.GetLogger().Errorf("twitter oauth access attempt without configuration")
		fmt.Fprintf(w, "error: this website has not configured twitter OAuth")
		return
	}

	commenterToken := r.FormValue("commenterToken")

	_, err := repository.Repo.CommenterRepository.GetCommenterByToken(commenterToken)
	if err != nil && err != app.ErrorNoSuchToken {
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	cred, err := twitterClient.RequestTemporaryCredentials(nil, os.Getenv("ORIGIN")+"/api/oauth/twitter/callback", nil)
	if err != nil {
		util.GetLogger().Errorf("cannot get temporary twitter credentials: %v", err)
		fmt.Fprintf(w, "error: %v", app.ErrorInternal.Error())
		return
	}

	twitterCredMapLock.Lock()
	twitterCredMap[cred.Token] = twitterOauthState{
		CommenterToken: commenterToken,
		Cred:           cred,
	}
	twitterCredMapLock.Unlock()

	http.Redirect(w, r, twitterClient.AuthorizationURL(cred, nil), http.StatusFound)
}
