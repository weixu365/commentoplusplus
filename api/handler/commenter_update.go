package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func commenterUpdate(commenterHex string, email string, name string, link string, photo string, provider string) error {
	if email == "" || name == "" || provider == "" {
		return app.ErrorMissingField
	}

	// See utils_sanitise.go's documentation on util.IsHttpsUrl. This is not a URL
	// validator, just an XSS preventor.
	// TODO: reject URLs instead of malforming them.
	if link == "" {
		link = "undefined"
	} else if link != "undefined" && !util.IsHttpsUrl(link) {
		link = "https://" + link
	}

	if photo == "" {
		photo = "undefined"
	} else if photo != "undefined" && !util.IsHttpsUrl(photo) {
		photo = "https://" + photo
	}

	// reserved "name"
	if name == "[deleted]" {
		return app.ErrorReservedName
	}

	err := repository.Repo.CommenterRepository.UpdateCommenter(&model.Commenter{
		CommenterHex: commenterHex,
		Provider:     provider,
		Email:        email,
		Name:         name,
		Link:         link,
		Photo:        photo,
	})
	if err != nil {
		util.GetLogger().Errorf("cannot update commenter: %v", err)
		return err
	}

	return nil
}

func CommenterUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		Name           *string `json:"name"`
		Email          *string `json:"email"`
		Link           *string `json:"link"`
		Photo          *string `json:"photo"`
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

	if commenter.Provider != "commento" {
		bodyMarshal(w, response{"success": false, "message": app.ErrorCannotUpdateOauthProfile.Error()})
		return
	}

	*x.Email = commenter.Email

	if err = commenterUpdate(commenter.CommenterHex, *x.Email, *x.Name, *x.Link, *x.Photo, commenter.Provider); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
