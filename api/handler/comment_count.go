package handler

import (
	"net/http"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func CommentCountHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Domain *string   `json:"domain"`
		Paths  *[]string `json:"paths"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := util.DomainStrip(*x.Domain)

	commentCounts, err := repository.Repo.CommentRepository.GetCommentsCount(domain, *x.Paths)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commentCounts": commentCounts})
}
