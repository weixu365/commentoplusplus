package handler

import (
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func domainOwnershipVerify(ownerHex string, domain string) (bool, error) {
	if ownerHex == "" || domain == "" {
		return false, app.ErrorMissingField
	}

	statement := `
		SELECT EXISTS (
			SELECT 1
			FROM domains
			WHERE ownerHex=$1 AND domain=$2
		);
	`
	row := repository.Db.QueryRow(statement, ownerHex, domain)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		util.GetLogger().Errorf("cannot query if domain owner: %v", err)
		return false, app.ErrorInternal
	}

	return exists, nil
}
