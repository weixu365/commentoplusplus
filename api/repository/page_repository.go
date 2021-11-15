package repository

import (
	"database/sql"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"

	"github.com/jmoiron/sqlx"
)

type PageRepositoryPg struct {
	db *sqlx.DB
}

func (r *PageRepositoryPg) CreatePage(domainName string, path string) error {
	// path can be empty
	if domainName == "" {
		return app.ErrorMissingField
	}

	statement := `
		INSERT INTO
		pages  (domain, path)
		VALUES ($1,     $2  )
		ON CONFLICT DO NOTHING;
	`
	_, err := r.db.Exec(statement, domainName, path)
	if err != nil {
		util.GetLogger().Errorf("error inserting new page: %v", err)
		return app.ErrorInternal
	}

	return nil
}

func (r *PageRepositoryPg) UpdatePage(p *model.Page) error {
	if p.Domain == "" {
		return app.ErrorMissingField
	}

	// fields to not update:
	//   commentCount
	statement := `
		INSERT INTO
		pages  (domain, path, isLocked, stickyCommentHex)
		VALUES ($1,     $2,   $3,       $4              )
		ON CONFLICT (domain, path) DO
			UPDATE SET isLocked = $3, stickyCommentHex = $4;
	`
	_, err := r.db.Exec(statement, p.Domain, p.Path, p.IsLocked, p.StickyCommentHex)
	if err != nil {
		util.GetLogger().Errorf("error setting page attributes: %v", err)
		return app.ErrorInternal
	}

	return nil
}

func (r *PageRepositoryPg) UpdatePageTitle(domain, path, title string) error {
	statement := `
		UPDATE pages
		SET title = $3
		WHERE canon($1) LIKE canon(domain) AND path = $2;
	`
	_, err ：= r.db.Exec(statement, domain, path, title)
	if err != nil {
		util.GetLogger().Errorf("cannot update pages table with title: %v", err)
		return err
	}

	return nil
}

func (r *PageRepositoryPg) GetPageByPath(domainName string, path string) (*model.Page, error) {
	// path can be empty
	if domainName == "" {
		return nil, app.ErrorMissingField
	}

	page := model.Page{Domain: domainName, Path: path}
	statement := `
		SELECT isLocked, commentCount, stickyCommentHex, title
		FROM pages
		WHERE canon($1) LIKE canon(domain) AND path=$2;
	`
	if err := r.db.Get(&page, statement, domainName, path); err != nil {
		if err == sql.ErrNoRows {
			// If there haven't been any comments, there won't be a record for this
			// page. The sane thing to do is return defaults.
			// TODO: the defaults are hard-coded in two places: here and the schema
			page.IsLocked = false
			page.CommentCount = 0
			page.StickyCommentHex = "none"
			page.Title = ""
		} else {
			util.GetLogger().Errorf("error scanning page: %v", err)
			return nil, app.ErrorInternal
		}
	}

	return &page, nil
}
