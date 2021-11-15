package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"
	"time"

	"github.com/jmoiron/sqlx"
)

type EmailRepositoryPg struct {
	db *sqlx.DB
}

func (r *EmailRepositoryPg) CreateEmail(emailAddress string) error {
	unsubscribeSecretHex, err := util.RandomHex(32)
	if err != nil {
		return err
	}
	statement := `
		INSERT INTO
		emails (email, unsubscribeSecretHex, lastEmailNotificationDate)
		VALUES ($1,    $2,                   $3                       )
		ON CONFLICT DO NOTHING;
	`
	_, err = r.db.Exec(statement, emailAddress, unsubscribeSecretHex, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert email into emails: %v", err)
		return err
	}

	return nil
}

func (r *EmailRepositoryPg) GetEmail(emailAddress string) (*model.Email, error) {
	email := model.Email{}
	statement := `
		SELECT *
		FROM emails
		WHERE email = $1;
	`

	if err := r.db.Get(&email, statement, emailAddress); err != nil {
		// TODO: is this the only error?
		return nil, app.ErrorNoSuchEmail
	}

	return &email, nil
}

func (r *EmailRepositoryPg) GetByUnsubscribeSecretHex(unsubscribeSecretHex string) (*model.Email, error) {
	email := model.Email{}
	statement := `
		SELECT *
		FROM emails
		WHERE unsubscribeSecretHex = $1;
	`

	if err := r.db.Get(&email, statement, unsubscribeSecretHex); err != nil {
		// TODO: is this the only error?
		return nil, app.ErrorNoSuchUnsubscribeSecretHex
	}

	return &email, nil
}
