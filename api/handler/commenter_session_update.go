package handler

import (
	"simple-commenting/repository"
	"simple-commenting/util"
)

func commenterSessionUpdate(commenterToken string, commenterHex string) error {
	if commenterToken == "" || commenterHex == "" {
		return errorMissingField
	}

	statement := `
		UPDATE commenterSessions
		SET commenterHex = $2
		WHERE commenterToken = $1;
	`
	_, err := repository.Db.Exec(statement, commenterToken, commenterHex)
	if err != nil {
		util.GetLogger().Errorf("error updating commenterHex: %v", err)
		return errorInternal
	}

	return nil
}
