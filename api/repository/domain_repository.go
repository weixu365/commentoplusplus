package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"
	"time"

	"github.com/jmoiron/sqlx"
)

type DomainRepositoryPg struct {
	db *sqlx.DB
}

func (r *DomainRepositoryPg) GetDomainByName(domainName string) (*model.Domain, error) {
	domain := model.Domain{}
	statement := `
		SELECT *
		FROM domains
		WHERE canon($1) LIKE canon(domain);
	`

	if err := r.db.Get(&domain, statement, domainName); err != nil {
		return nil, app.ErrorNoSuchDomain
	}

	return &domain, nil
}

func (r *DomainRepositoryPg) ListDomain(ownerHex string) ([]*model.Domain, error) {
	if ownerHex == "" {
		return nil, app.ErrorMissingField
	}

	domains := []*model.Domain{}
	statement := `
		SELECT *
		FROM domains
		WHERE ownerHex=$1;
	`
	err := r.db.Select(&domains, statement, ownerHex)
	if err != nil {
		util.GetLogger().Errorf("cannot query domains: %v", err)
		return nil, err
	}

	for _, domain := range domains {
		(*domain).Moderators, err = Repo.DomainModeratorRepository.GetModeratorsForDomain((*domain).Domain)
		if err != nil {
			return nil, err
		}
	}

	return domains, nil
}

func (r *DomainRepositoryPg) CreateDomain(ownerHex string, name string, domain string) error {
	// test if domain already exists
	statement := `
		SELECT COUNT(*) FROM
		domains WHERE
		canon(regexp_replace($1, '[%]', '')) LIKE canon(domain) OR canon(domain) LIKE canon($1);
	`
	row := r.db.QueryRow(statement, domain)
	var err error
	var count int

	if err = row.Scan(&count); err != nil {
		return app.ErrorInvalidDomain
	}

	if count > 0 {
		return app.ErrorDomainAlreadyExists
	}

	statement = `
		INSERT INTO
		domains (ownerHex, name, domain, creationDate)
		VALUES  ($1,       $2,   $3,     $4          );
	`
	_, err = r.db.Exec(statement, ownerHex, name, domain, time.Now().UTC())
	if err != nil {
		// TODO: This should not happen given the above check, so this is likely not the error. Be more informative?
		return app.ErrorDomainAlreadyExists
	}

	return nil
}

func (r *DomainRepositoryPg) UpdateDomain(d *model.Domain) error {
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

	_, err := r.db.Exec(statement,
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
		return err
	}

	return nil
}
