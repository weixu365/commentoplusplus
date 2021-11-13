package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func domainUpdate(d model.Domain) error {
	if d.SsoProvider && d.SsoUrl == "" {
		return app.ErrorMissingField
	}

	statement := `
		UPDATE domains
		SET
			name=$2,
			state=$3,
			autoSpamFilter=$4,
			requireModeration=$5,
			requireIdentification=$6,
			moderateAllAnonymous=$7,
			emailNotificationPolicy=$8,
			commentoProvider=$9,
			googleProvider=$10,
			twitterProvider=$11,
			githubProvider=$12,
			gitlabProvider=$13,
			ssoProvider=$14,
			ssoUrl=$15,
			defaultSortPolicy=$16
		WHERE domain=$1;
	`

	_, err := repository.Db.Exec(statement,
		d.Domain,
		d.Name,
		d.State,
		d.AutoSpamFilter,
		d.RequireModeration,
		d.RequireIdentification,
		d.ModerateAllAnonymous,
		d.EmailNotificationPolicy,
		d.CommentoProvider,
		d.GoogleProvider,
		d.TwitterProvider,
		d.GithubProvider,
		d.GitlabProvider,
		d.SsoProvider,
		d.SsoUrl,
		d.DefaultSortPolicy)
	if err != nil {
		util.GetLogger().Errorf("cannot update non-moderators: %v", err)
		return app.ErrorInternal
	}

	return nil
}

func DomainUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string       `json:"ownerToken"`
		D          *model.Domain `json:"domain"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := util.DomainStrip((*x.D).Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotAuthorised.Error()})
		return
	}

	if err = domainUpdate(*x.D); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
