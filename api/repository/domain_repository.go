package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"

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
