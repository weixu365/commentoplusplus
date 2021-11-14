package repository

import (
	"database/sql"
	"simple-commenting/model"

	"github.com/jmoiron/sqlx"
)

var Db *sql.DB

var db *sqlx.DB

var Repo *Repositories

type Repositories struct {
	DomainRepository          DomainRepository
	DomainModeratorRepository DomainModeratorRepository
	EmailRepository           EmailRepository
}

type DomainRepository interface {
	GetDomainByName(domainName string) (*model.Domain, error)
}

type DomainModeratorRepository interface {
	CreateModerator(domain string, email string) error
	GetModeratorsForDomain(domainName string) (*[]model.Moderator, error)
	IsDomainModerator(domain string, email string) (bool, error)
}

type EmailRepository interface {
	CreateEmail(emailAddress string) error
}
