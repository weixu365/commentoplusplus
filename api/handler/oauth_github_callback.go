package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"

	"golang.org/x/oauth2"
)

func githubGetPrimaryEmail(accessToken string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Add("Authorization", "token "+accessToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", app.ErrorCannotReadResponse
	}

	user := []map[string]interface{}{}
	if err := json.Unmarshal(contents, &user); err != nil {
		util.GetLogger().Errorf("error unmarshaling github user: %v", err)
		return "", app.ErrorInternal
	}

	nonPrimaryEmail := ""
	for _, email := range user {
		nonPrimaryEmail = email["email"].(string)
		if email["primary"].(bool) {
			return email["email"].(string), nil
		}
	}

	return nonPrimaryEmail, nil
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken := r.FormValue("state")
	code := r.FormValue("code")

	_, err := repository.Repo.CommenterRepository.GetCommenterByToken(commenterToken)
	if err != nil && err != app.ErrorNoSuchToken {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	token, err := githubConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	email, err := githubGetPrimaryEmail(token.AccessToken)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Authorization", "token "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", app.ErrorCannotReadResponse.Error())
		return
	}

	user := make(map[string]interface{})
	if err := json.Unmarshal(contents, &user); err != nil {
		fmt.Fprintf(w, "Error: %s", app.ErrorInternal.Error())
		return
	}

	if email == "" {
		if user["email"] == nil {
			fmt.Fprintf(w, "Error: no email address returned by Github")
			return
		}

		email = user["email"].(string)
	}

	name := user["login"].(string)
	if user["name"] != nil {
		name = user["name"].(string)
	}

	link := "undefined"
	if user["html_url"] != nil {
		link = user["html_url"].(string)
	}

	photo := "undefined"
	if user["avatar_url"] != nil {
		photo = user["avatar_url"].(string)
	}

	c, err := repository.Repo.CommenterRepository.GetCommenterByEmail("github", email)
	if err != nil && err != app.ErrorNoSuchCommenter {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	var commenterHex string

	if err == app.ErrorNoSuchCommenter {
		commenterHex, err = commenterNew(email, name, link, photo, "github", "")
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
	} else {
		if err = commenterUpdate(c.CommenterHex, email, name, link, photo, "github"); err != nil {
			util.GetLogger().Warningf("cannot update commenter: %s", err)
			// not a serious enough to exit with an error
		}

		commenterHex = c.CommenterHex
	}

	if err := repository.Repo.CommenterRepository.(commenterToken, commenterHex); err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	fmt.Fprintf(w, "<html><script>window.parent.close()</script></html>")
}
