package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"
)

type DomainRepository interface {
	GetDomainByName(domainName string) (*model.Domain, error)

}

type DomainRepositoryPg struct {

}

func (r *DomainRepositoryPg) GetDomainByName(domainName string) (*model.Domain, error){
	domain := model.Domain{}
	statement := `
		SELECT *
		FROM domains
		WHERE canon($1) LIKE canon(domain);
	`
	
	if err := db.Get(&domain, statement, domainName); err != nil {
		return nil, app.ErrorNoSuchDomain
	}

	return &domain, nil
}
