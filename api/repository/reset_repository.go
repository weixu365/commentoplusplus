package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"
	"time"

	"github.com/jmoiron/sqlx"
)

type ResetRepository interface {
	GetResetHex(resetHex string) (*model.ResetHex, error)
	DeleteResetHex(resetHex string) error
	CreateResetHex(hex, entity string) (string, error)
}

type ResetRepositoryPg struct {
	db *sqlx.DB
}

func (r *ResetRepositoryPg) GetResetHex(resetHex string) (*model.ResetHex, error) {
	if resetHex == "" {
		return nil, app.ErrorMissingField
	}

	hex := model.ResetHex{}
	statement := `
		SELECT hex, entity
		FROM resetHexes
		WHERE resetHex = $1;
	`

	if err := r.db.Get(&hex, statement, resetHex); err != nil {
		// TODO: is this the only error?
		return nil, app.ErrorNoSuchResetToken
	}

	return &hex, nil
}

func (r *ResetRepositoryPg) DeleteResetHex(resetHex string) error {
	if resetHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		DELETE FROM resetHexes
		WHERE resetHex = $1;
	`

	if _, err := r.db.Exec(statement, resetHex); err != nil {
		// TODO: is this the only error?
		return app.ErrorNoSuchResetToken
	}

	return nil
}

func (r *ResetRepositoryPg) CreateResetHex(hex, entity string) (string, error) {
	resetHex, err := util.RandomHex(32)
	if err != nil {
		return "", err
	}

	statement := `
		INSERT INTO
		resetHexes (resetHex, hex, entity, sendDate)
		VALUES     ($1,       $2,  $3,     $4      );
	`
	_, err = r.db.Exec(statement, resetHex, hex, entity, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert resetHex: %v", err)
		return "", err
	}

	return resetHex, nil
}
