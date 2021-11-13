package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
)

var emailsRowColumns = `
	emails.email,
	emails.unsubscribeSecretHex,
	emails.lastEmailNotificationDate,
	emails.sendReplyNotifications,
	emails.sendModeratorNotifications
`

func emailsRowScan(s repository.SqlScanner, e *model.Email) error {
	return s.Scan(
		&e.Email,
		&e.UnsubscribeSecretHex,
		&e.LastEmailNotificationDate,
		&e.SendReplyNotifications,
		&e.SendModeratorNotifications,
	)
}

func emailGet(em string) (model.Email, error) {
	statement := `
		SELECT ` + emailsRowColumns + `
		FROM emails
		WHERE email = $1;
	`
	row := repository.Db.QueryRow(statement, em)

	var e model.Email
	if err := emailsRowScan(row, &e); err != nil {
		// TODO: is this the only error?
		return e, app.ErrorNoSuchEmail
	}

	return e, nil
}

func emailGetByUnsubscribeSecretHex(unsubscribeSecretHex string) (model.Email, error) {
	statement := `
		SELECT ` + emailsRowColumns + `
		FROM emails
		WHERE unsubscribeSecretHex = $1;
	`
	row := repository.Db.QueryRow(statement, unsubscribeSecretHex)

	e := model.Email{}
	if err := emailsRowScan(row, &e); err != nil {
		// TODO: is this the only error?
		return e, app.ErrorNoSuchUnsubscribeSecretHex
	}

	return e, nil
}

func EmailGetHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UnsubscribeSecretHex *string `json:"unsubscribeSecretHex"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	e, err := emailGetByUnsubscribeSecretHex(*x.UnsubscribeSecretHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "email": e})
}
