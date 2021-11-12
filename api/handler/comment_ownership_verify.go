package handler

import (
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func commentOwnershipVerify(commenterHex string, commentHex string) (bool, error) {
	if commenterHex == "" || commentHex == "" {
		return false, app.ErrorMissingField
	}

	statement := `
		SELECT EXISTS (
			SELECT 1
			FROM comments
			WHERE commenterHex=$1 AND commentHex=$2
		);
	`
	row := repository.Db.QueryRow(statement, commenterHex, commentHex)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		util.GetLogger().Errorf("cannot query if comment owner: %v", err)
		return false, app.ErrorInternal
	}

	return exists, nil
}
