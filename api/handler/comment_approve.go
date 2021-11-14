package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/notification"
	"simple-commenting/repository"
)

func CommentApproveHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		CommentHex     *string `json:"commentHex"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	commenter, err := repository.Repo.CommenterRepository.GetCommenterByToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain, path, err := repository.Repo.CommentRepository.GetCommentDomainPath(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	isModerator, err := isDomainModerator(domain, commenter.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isModerator {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotModerator.Error()})
		return
	}

	if err = repository.Repo.CommentRepository.ApproveComment(*x.CommentHex, domain+path); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	notification.NotificationHub.Broadcast <- []byte(domain + path)

	bodyMarshal(w, response{"success": true})
}

func CommentOwnerApproveHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		CommentHex *string `json:"commentHex"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain, path, err := repository.Repo.CommentRepository.GetCommentDomainPath(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotAuthorised.Error()})
		return
	}

	if err = repository.Repo.CommentRepository.ApproveComment(*x.CommentHex, domain+path); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	notification.NotificationHub.Broadcast <- []byte(domain + path)

	bodyMarshal(w, response{"success": true})
}
