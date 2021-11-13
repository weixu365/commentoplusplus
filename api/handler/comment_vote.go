package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"
)

func commentVote(commenterHex string, commentHex string, direction int, url string) error {
	if commentHex == "" || commenterHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		SELECT commenterHex
		FROM comments
		WHERE commentHex = $1;
	`
	row := repository.Db.QueryRow(statement, commentHex)

	var authorHex string
	if err := row.Scan(&authorHex); err != nil {
		util.GetLogger().Errorf("error selecting authorHex for vote")
		return app.ErrorInternal
	}

	if authorHex == commenterHex {
		return app.ErrorSelfVote
	}

	statement = `
		INSERT INTO
		votes  (commentHex, commenterHex, direction, voteDate)
		VALUES ($1,         $2,           $3,        $4      )
		ON CONFLICT (commentHex, commenterHex) DO
		UPDATE SET direction = $3;
	`
	_, err := repository.Db.Exec(statement, commentHex, commenterHex, direction, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("error inserting/updating votes: %v", err)
		return app.ErrorInternal
	}

	notification.NotificationHub.Broadcast <- []byte(url)

	return nil
}

func commentVoteHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		CommentHex     *string `json:"commentHex"`
		Direction      *int    `json:"direction"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if *x.CommenterToken == "anonymous" {
		bodyMarshal(w, response{"success": false, "message": app.ErrorUnauthorisedVote.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	direction := 0
	if *x.Direction > 0 {
		direction = 1
	} else if *x.Direction < 0 {
		direction = -1
	}

	domain, path, err := commentDomainPathGet(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := commentVote(c.CommenterHex, *x.CommentHex, direction, domain+path); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
