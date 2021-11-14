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
	PageRepository            PageRepository
	CommenterRepository       CommenterRepository
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

type PageRepository interface {
	CreatePage(domainName string, path string) error
	GetPageByPath(domainName string, path string) (*model.Page, error)
}

type CommenterRepository interface {
	GetCommenterByEmail(provider string, email string) (*model.Commenter, error)
	GetCommenterByHex(commenterHex string) (*model.Commenter, error)
	GetCommenterByToken(commenterToken string) (*model.Commenter, error)
}
