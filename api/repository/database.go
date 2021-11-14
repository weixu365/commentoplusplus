package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var Db *sql.DB

var db *sqlx.DB

type Repositories struct {
	domainRepository DomainRepository
	domainModeratorRepository DomainModeratorRepository
}

var Repo *Repositories

