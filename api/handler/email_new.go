package handler

import (
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"
)

func emailNew(email string) error {
	unsubscribeSecretHex, err := randomHex(32)
	if err != nil {
		return app.ErrorInternal
	}

	statement := `
		INSERT INTO
		emails (email, unsubscribeSecretHex, lastEmailNotificationDate)
		VALUES ($1,    $2,                   $3                       )
		ON CONFLICT DO NOTHING;
	`
	_, err = repository.Db.Exec(statement, email, unsubscribeSecretHex, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert email into emails: %v", err)
		return app.ErrorInternal
	}

	return nil
}
