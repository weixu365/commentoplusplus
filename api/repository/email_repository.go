package repository

import (
	"simple-commenting/app"
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
		return app.ErrorInternal
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
		return app.ErrorInternal
	}

	return nil
}
