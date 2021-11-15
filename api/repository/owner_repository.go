package repository

import (
	"simple-commenting/app"

	"github.com/jmoiron/sqlx"
)

type OwnerRepository interface {
	UpdatePassword(passwordHash, ownerHex string) error
}

type OwnerRepositoryPg struct {
	db *sqlx.DB
}

func (r *OwnerRepositoryPg) UpdatePassword(passwordHash, ownerHex string) error {
	if ownerHex == "" || passwordHash == "" {
		return app.ErrorMissingField
	}

	statement := `
			UPDATE owners SET passwordHash = $1, confirmedEmail=true
			WHERE ownerHex = $2;
		`

	if _, err := r.db.Exec(statement, passwordHash, ownerHex); err != nil {
		// TODO: is this the only error?
		return app.ErrorNoSuchResetToken
	}

	return nil
}

func (r *OwnerRepositoryPg) DeleteResetHex(resetHex string) error {
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
