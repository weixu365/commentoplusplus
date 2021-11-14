package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/notification"
	"simple-commenting/repository"
)

func CommentVoteHandler(w http.ResponseWriter, r *http.Request) {
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

	commenter, err := repository.Repo.CommenterRepository.GetCommenterByToken(*x.CommenterToken)
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

	domain, path, err := repository.Repo.CommentRepository.GetCommentDomainPath(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := repository.Repo.CommentRepository.VoteComment(commenter.CommenterHex, *x.CommentHex, direction, domain+path); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	notification.NotificationHub.Broadcast <- []byte(domain + path)

	bodyMarshal(w, response{"success": true})
}
