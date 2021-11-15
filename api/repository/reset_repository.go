package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"

	"github.com/jmoiron/sqlx"
)

type ResetRepository interface {
	GetResetHex(resetHex string) (*model.ResetHex, error)
	DeleteResetHex(resetHex string) error
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
