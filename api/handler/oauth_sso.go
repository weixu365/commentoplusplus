package handler

import (
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"
)

type ssoPayload struct {
	Domain string `json:"domain"`
	Token  string `json:"token"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Link   string `json:"link"`
	Photo  string `json:"photo"`
}

func ssoTokenNew(domain string, commenterToken string) (string, error) {
	token, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("error generating SSO token hex: %v", err)
		return "", app.ErrorInternal
	}

	statement := `
		INSERT INTO
		ssoTokens (token, domain, commenterToken, creationDate)
		VALUES    ($1,    $2,     $3,             $4          );
	`
	_, err = repository.Db.Exec(statement, token, domain, commenterToken, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("error inserting SSO token: %v", err)
		return "", app.ErrorInternal
	}

	return token, nil
}

func ssoTokenExtract(token string) (string, string, error) {
	statement := `
		SELECT domain, commenterToken
		FROM ssoTokens
		WHERE token = $1;
	`
	row := repository.Db.QueryRow(statement, token)

	var domain string
	var commenterToken string
	if err := row.Scan(&domain, &commenterToken); err != nil {
		return "", "", app.ErrorNoSuchToken
	}

	statement = `
		DELETE FROM ssoTokens
		WHERE token = $1;
	`
	if _, err := repository.Db.Exec(statement, token); err != nil {
		util.GetLogger().Errorf("cannot delete SSO token after usage: %v", err)
		return "", "", app.ErrorInternal
	}

	return domain, commenterToken, nil
}
