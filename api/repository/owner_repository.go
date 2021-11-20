package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"
	"time"

	"github.com/jmoiron/sqlx"
)

type OwnerRepository interface {
	CreateOwner(owner *model.Owner) (*model.Owner, error)
	CreateOwnerConfirmHex(ownerHex string) (string, error)
	UpdatePassword(passwordHash, ownerHex string) error
	GetByEmail(email string) (*model.Owner, error)
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

func (r *OwnerRepositoryPg) CreateOwner(owner *model.Owner) (*model.Owner, error) {
	ownerHex, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("cannot generate ownerHex: %v", err)
		return nil, err
	}
	owner.OwnerHex = ownerHex
	owner.JoinDate = time.Now().UTC()

	statement := `
		INSERT INTO
		owners (ownerHex, email, name, passwordHash, joinDate, confirmedEmail)
		VALUES (:OwnerHex,  :Email, :Name, :PasswordHash, :JoinDate, :ConfirmedEmail);
	`
	_, err = r.db.NamedExec(statement, owner)
	if err != nil {
		// TODO: Make sure `err` is actually about conflicting UNIQUE, and not some
		// other error. If it is something else, we should probably return `errorInternal`.
		return nil, app.ErrorEmailAlreadyExists
	}

	return owner, nil
}

func (r *OwnerRepositoryPg) CreateOwnerConfirmHex(ownerHex string) (string, error) {
	confirmHex, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("cannot generate confirmHex: %v", err)
		return "", err
	}

	statement := `
			INSERT INTO
			ownerConfirmHexes (confirmHex, ownerHex, sendDate)
			VALUES            ($1,         $2,       $3      );
		`
	_, err = r.db.Exec(statement, confirmHex, ownerHex, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert confirmHex: %v\n", err)
		return "", err
	}

	return confirmHex, nil
}

func (r *OwnerRepositoryPg) GetByEmail(email string) (*model.Owner, error) {
	if email == "" {
		return nil, app.ErrorMissingField
	}

	owner := model.Owner{}
	statement := `
		SELECT *
		FROM owners
		WHERE email=$1;
	`

	if err := r.db.Get(&owner, statement, email); err != nil {
		// TODO: Make sure this is actually no such email.
		return nil, app.ErrorNoSuchEmail
	}

	return &owner, nil
}
