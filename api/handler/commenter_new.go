package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func commenterNew(email string, name string, link string, photo string, provider string, password string) (string, error) {
	if email == "" || name == "" || link == "" || photo == "" || provider == "" {
		return "", app.ErrorMissingField
	}

	if provider == "commento" && password == "" {
		return "", app.ErrorMissingField
	}

	// See utils_sanitise.go's documentation on IsHttpsUrl. This is not a URL
	// validator, just an XSS preventor.
	// TODO: reject URLs instead of malforming them.
	if link != "undefined" && !util.IsHttpsUrl(link) {
		link = "https://" + link
	}

	if provider != "anon" {
		if _, err := repository.Repo.CommenterRepository.GetCommenterByEmail(provider, email); err == nil {
			return "", app.ErrorEmailAlreadyExists
		}

		if err := repository.Repo.EmailRepository.CreateEmail(email); err != nil {
			return "", app.ErrorInternal
		}
	}

	return repository.Repo.CommenterRepository.CreateCommenter(email, name, link, photo, provider, password)
}

func CommenterNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    *string `json:"email"`
		Name     *string `json:"name"`
		Website  *string `json:"website"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// TODO: add gravatar?
	// TODO: email confirmation if provider = commento?
	// TODO: email confirmation if provider = commento?
	if *x.Website == "" {
		*x.Website = "undefined"
	}

	if _, err := commenterNew(*x.Email, *x.Name, *x.Website, "undefined", "commento", *x.Password); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "confirmEmail": notification.SmtpConfigured})
}
