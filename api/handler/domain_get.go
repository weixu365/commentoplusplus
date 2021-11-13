package handler

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
)

var domainsRowColumns = `
	domains.domain,
	domains.ownerHex,
	domains.name,
	domains.creationDate,
	domains.state,
	domains.importedComments,
	domains.autoSpamFilter,
	domains.requireModeration,
	domains.requireIdentification,
	domains.moderateAllAnonymous,
	domains.emailNotificationPolicy,
	domains.commentoProvider,
	domains.googleProvider,
	domains.twitterProvider,
	domains.githubProvider,
	domains.gitlabProvider,
	domains.ssoProvider,
	domains.ssoSecret,
	domains.ssoUrl,
	domains.defaultSortPolicy
`

func domainsRowScan(s repository.SqlScanner, d *model.Domain) error {
	return s.Scan(
		&d.Domain,
		&d.OwnerHex,
		&d.Name,
		&d.CreationDate,
		&d.State,
		&d.ImportedComments,
		&d.AutoSpamFilter,
		&d.RequireModeration,
		&d.RequireIdentification,
		&d.ModerateAllAnonymous,
		&d.EmailNotificationPolicy,
		&d.CommentoProvider,
		&d.GoogleProvider,
		&d.TwitterProvider,
		&d.GithubProvider,
		&d.GitlabProvider,
		&d.SsoProvider,
		&d.SsoSecret,
		&d.SsoUrl,
		&d.DefaultSortPolicy,
	)
}

func domainGet(dmn string) (model.Domain, error) {
	if dmn == "" {
		return model.Domain{}, app.ErrorMissingField
	}

	statement := `
		SELECT ` + domainsRowColumns + `
		FROM domains
		WHERE canon($1) LIKE canon(domain);
	`
	row := repository.Db.QueryRow(statement, dmn)

	var err error
	d := model.Domain{}
	if err = domainsRowScan(row, &d); err != nil {
		return d, app.ErrorNoSuchDomain
	}

	d.Moderators, err = domainModeratorList(d.Domain)
	if err != nil {
		return model.Domain{}, err
	}

	return d, nil
}
