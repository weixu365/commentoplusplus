package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"
)

type DomainModeratorRepository interface {
	GetModeratorsForDomain(domainName string) (*[]model.Moderator, error)

}

type DomainModeratorRepositoryPg struct {

}

func (r *DomainModeratorRepositoryPg) GetModeratorsForDomain(domainName string) (*[]model.Moderator, error){
	moderators := []model.Moderator{}
	statement := `
		SELECT email, addDate
		FROM moderators
		WHERE domain=$1;
	`
	
	if err := db.Select(&moderators, statement, domainName); err != nil {
		return nil, app.ErrorInternal
	}

	return &moderators, nil
}
