package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/notification"
	"simple-commenting/repository"
)

func CommentDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	comment, err := repository.Repo.CommentRepository.GetByCommentHex(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domainName, path, err := repository.Repo.CommentRepository.GetCommentDomainPath(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	isModerator, err := isDomainModerator(domainName, commenter.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isModerator && comment.CommenterHex != commenter.CommenterHex {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotModerator.Error()})
		return
	}

	if err = repository.Repo.CommentRepository.DeleteComment(*x.CommentHex, *x.CommenterToken, domainName, path); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	notification.NotificationHub.Broadcast <- []byte(domainName + path)

	bodyMarshal(w, response{"success": true})
}

func CommentOwnerDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = repository.Repo.CommentRepository.DeleteComment(*x.CommentHex, *x.OwnerToken, domain, path); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	notification.NotificationHub.Broadcast <- []byte(domain + path)

	bodyMarshal(w, response{"success": true})
}
