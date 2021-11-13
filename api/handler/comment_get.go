package handler

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
)

var commentsRowColumns = `
	comments.commentHex,
	comments.commenterHex,
	comments.markdown,
	comments.html,
	comments.parentHex,
	comments.score,
	comments.state,
	comments.deleted,
	comments.creationDate
`

func commentsRowScan(s repository.SqlScanner, c *model.Comment) error {
	return s.Scan(
		&c.CommentHex,
		&c.CommenterHex,
		&c.Markdown,
		&c.Html,
		&c.ParentHex,
		&c.Score,
		&c.State,
		&c.Deleted,
		&c.CreationDate,
	)
}

func commentGetByCommentHex(commentHex string) (model.Comment, error) {
	if commentHex == "" {
		return model.Comment{}, app.ErrorMissingField
	}

	statement := `
		SELECT ` + commentsRowColumns + `
		FROM comments
		WHERE comments.commentHex = $1;
	`
	row := repository.Db.QueryRow(statement, commentHex)

	var c model.Comment
	if err := commentsRowScan(row, &c); err != nil {
		// TODO: is this the only error?
		return c, app.ErrorNoSuchComment
	}

	return c, nil
}
