package repository

import (
	"simple-commenting/util"

	"github.com/jmoiron/sqlx"
)

func NewPostgresqlRepositories(dataSourceName string) (*Repositories, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		util.GetLogger().Errorf("sqlx: cannot open connection to postgres: %v", err)
		return nil, err
	}

	return &Repositories{
		DomainRepository:          &DomainRepositoryPg{db: db},
		DomainModeratorRepository: &DomainModeratorRepositoryPg{db: db},
		EmailRepository:           &EmailRepositoryPg{db: db},
		PageRepository:            &PageRepositoryPg{db: db},
	}, nil
}
