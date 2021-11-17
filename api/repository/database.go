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
	CommentRepository         CommentRepository
	LogRepository             LogRepository
	OwnerRepository           OwnerRepository
	ResetRepository           ResetRepository
	StatisticsRepository      StatisticsRepository
}

type DomainRepository interface {
	CreateDomain(ownerHex string, name string, domain string) error
	UpdateDomain(d *model.Domain) error
	GetDomainByName(domainName string) (*model.Domain, error)
	ListDomain(ownerHex string) ([]*model.Domain, error)
}

type DomainModeratorRepository interface {
	CreateModerator(domain string, email string) error
	GetModeratorsForDomain(domainName string) (*[]model.Moderator, error)
	IsDomainModerator(domain string, email string) (bool, error)
}

type EmailRepository interface {
	CreateEmail(emailAddress string) error
	UpdateEmail(e *model.Email) error
	GetEmail(emailAddress string) (*model.Email, error)
	GetByUnsubscribeSecretHex(unsubscribeSecretHex string) (*model.Email, error)
}

type PageRepository interface {
	CreatePage(domainName string, path string) error
	UpdatePage(p *model.Page) error
	UpdatePageTitle(domain, path, title string) error
	GetPageByPath(domainName string, path string) (*model.Page, error)
}
