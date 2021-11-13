package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func commentEdit(commentHex string, markdown string, url string) (string, error) {
	if commentHex == "" {
		return "", app.ErrorMissingField
	}

	html := util.MarkdownToHtml(markdown)

	statement := `
		UPDATE comments
		SET markdown = $2, html = $3
		WHERE commentHex=$1;
	`
	_, err := repository.Db.Exec(statement, commentHex, markdown, html)

	if err != nil {
		// TODO: make sure this is the error is actually non-existant commentHex
		return "", app.ErrorNoSuchComment
	}

	notification.NotificationHub.Broadcast <- []byte(url)

	return html, nil
}

func CommentEditHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		CommentHex     *string `json:"commentHex"`
		Markdown       *string `json:"markdown"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain, path, err := commentDomainPathGet(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	cm, err := commentGetByCommentHex(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if cm.CommenterHex != c.CommenterHex {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotAuthorised.Error()})
		return
	}

	html, err := commentEdit(*x.CommentHex, *x.Markdown, domain+path)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "html": html})
}
