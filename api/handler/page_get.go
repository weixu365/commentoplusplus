package handler

import (
	"database/sql"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func pageGet(domain string, path string) (model.Page, error) {
	// path can be empty
	if domain == "" {
		return model.Page{}, app.ErrorMissingField
	}

	statement := `
		SELECT isLocked, commentCount, stickyCommentHex, title
		FROM pages
		WHERE canon($1) LIKE canon(domain) AND path=$2;
	`
	row := repository.Db.QueryRow(statement, domain, path)

	p := model.Page{Domain: domain, Path: path}
	if err := row.Scan(&p.IsLocked, &p.CommentCount, &p.StickyCommentHex, &p.Title); err != nil {
		if err == sql.ErrNoRows {
			// If there haven't been any comments, there won't be a record for this
			// page. The sane thing to do is return defaults.
			// TODO: the defaults are hard-coded in two places: here and the schema
			p.IsLocked = false
			p.CommentCount = 0
			p.StickyCommentHex = "none"
			p.Title = ""
		} else {
			util.GetLogger().Errorf("error scanning page: %v", err)
			return model.Page{}, app.ErrorInternal
		}
	}

	return p, nil
}
