package repository

import (
	"database/sql"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"
	"time"

	"github.com/jmoiron/sqlx"
)

type DomainModeratorRepositoryPg struct {
	db *sqlx.DB
}

func (r *DomainModeratorRepositoryPg) GetModeratorsForDomain(domainName string) (*[]model.Moderator, error) {
	moderators := []model.Moderator{}
	statement := `
		SELECT email, addDate
		FROM moderators
		WHERE domain=$1;
	`

	if err := r.db.Select(&moderators, statement, domainName); err != nil {
		return nil, app.ErrorInternal
	}

	return &moderators, nil
}

func (r *DomainModeratorRepositoryPg) CreateModerator(domain string, email string) error {
	statement := `
		INSERT INTO
		moderators (domain, email, addDate)
		VALUES     ($1,     $2,    $3     );
	`
	_, err := r.db.Exec(statement, domain, email, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert new moderator: %v", err)
		return app.ErrorInternal
	}

	return nil
}

func (r *DomainModeratorRepositoryPg) IsDomainModerator(domain string, email string) (bool, error) {
	moderator := model.Moderator{}

	statement := `
		SELECT *
		FROM moderators
		WHERE domain=$1 AND email=$2
	`
	err := r.db.Get(&moderator, statement, domain, email)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
