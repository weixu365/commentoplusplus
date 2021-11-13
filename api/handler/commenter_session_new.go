package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"
)

func commenterTokenNew() (string, error) {
	commenterToken, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("cannot create commenterToken: %v", err)
		return "", app.ErrorInternal
	}

	statement := `
		INSERT INTO
		commenterSessions (commenterToken, creationDate)
		VALUES            ($1,             $2          );
	`
	_, err = repository.Db.Exec(statement, commenterToken, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert new commenterToken: %v", err)
		return "", app.ErrorInternal
	}

	return commenterToken, nil
}

func CommenterTokenNewHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken, err := commenterTokenNew()
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commenterToken": commenterToken})
}
