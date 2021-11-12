package handler

import (
	"simple-commenting/app"
	"simple-commenting/repository"
)

func commentDomainPathGet(commentHex string) (string, string, error) {
	if commentHex == "" {
		return "", "", app.ErrorMissingField
	}

	statement := `
		SELECT domain, path
		FROM comments
		WHERE commentHex = $1;
	`
	row := repository.Db.QueryRow(statement, commentHex)

	var domain string
	var path string
	var err error
	if err = row.Scan(&domain, &path); err != nil {
		return "", "", app.ErrorNoSuchDomain
	}

	return domain, path, nil
}
