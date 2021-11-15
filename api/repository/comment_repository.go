package repository

import (
	"database/sql"
	"os"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"
	"time"

	"github.com/jmoiron/sqlx"
)

type CommentRepository interface {
	CreateComment(comment *model.Comment) (*model.Comment, error)
	UpdateComment(commentHex string, markdown string, html string) error
	GetByCommentHex(commentHex string) (*model.Comment, error)
	GetCommentDomainPath(commentHex string) (string, string, error)
	ApproveComment(commentHex string, url string) error
	VoteComment(commenterHex string, commentHex string, direction int, url string) error
	DeleteComment(commentHex string, deleterHex string, domain string, path string) error
	ListComments(commenterHex string, domain string, path string, includeUnapproved bool) ([]*model.Comment, map[string]*model.Commenter, error)
	GetCommentsCount(domain string, paths []string) (map[string]int, error)
}

type CommentRepositoryPg struct {
	db *sqlx.DB
}

func (r *CommentRepositoryPg) CreateComment(comment *model.Comment) (*model.Comment, error) {
	statement := `
		INSERT INTO
		comments (commentHex, domain, path, commenterHex, parentHex, markdown, html, creationDate, state)
		VALUES   (:commentHex,:domain, :path, :commenterHex, :parentHex, :markdown, :html, :creationDate, :state);
	`
	//TODO: check if need `db:"first_name"` in model struct
	_, err := r.db.NamedExec(statement, comment)
	if err != nil {
		util.GetLogger().Errorf("cannot insert comment: %v", err)
		return nil, err
	}

	err = r.UpdateCommentCount(comment.Domain, comment.Path)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepositoryPg) UpdateComment(commentHex, markdown, html string) error {
	if commentHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		UPDATE comments
		SET markdown = $2, html = $3
		WHERE commentHex=$1;
	`
	_, err := r.db.Exec(statement, commentHex, markdown, html)

	if err != nil {
		// TODO: make sure this is the error is actually non-existant commentHex
		return app.ErrorNoSuchComment
	}

	return nil
}

func (r *CommentRepositoryPg) UpdateCommentCount(domainName string, path string) error {
	if domainName == "" {
		return app.ErrorMissingField
	}

	statement := `
		UPDATE pages
		set commentCount = (
			select count(*) from comments 
		)
		WHERE canon($1) LIKE canon(domain) AND path = $2
	`

	_, err := r.db.Exec(statement, domainName, path)
	if err != nil {
		util.GetLogger().Errorf("Failed to update cannot update comments count: %v", err)
		return err
	}

	return nil
}

func (r *CommentRepositoryPg) GetCommentsCount(domainName string, paths []string) (map[string]int, error) {
	commentCounts := map[string]int{}

	if domainName == "" {
		return nil, app.ErrorMissingField
	}

	if len(paths) == 0 {
		return nil, app.ErrorEmptyPaths
	}

	statement := `
		SELECT path, commentCount
		FROM pages
		WHERE canon(:domain) LIKE canon(domain) AND path in (:paths);
	`
	query, args, err := sqlx.Named(statement, map[string]interface{}{
		"domain": domainName,
		"paths":  paths,
	})
	if err != nil {
		util.GetLogger().Errorf("Failed to create named query when get count of comments: %v", err)
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		util.GetLogger().Errorf("Failed to bind parameters to variables when get count of comments: %v", err)
		return nil, err
	}

	query = r.db.Rebind(query)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		util.GetLogger().Errorf("cannot get comments: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		var commentCount int
		if err = rows.Scan(&path, &commentCount); err != nil {
			util.GetLogger().Errorf("cannot scan path and commentCount: %v", err)
			return nil, err
		}

		commentCounts[path] = commentCount
	}

	return commentCounts, nil
}

func (r *CommentRepositoryPg) GetByCommentHex(commentHex string) (*model.Comment, error) {
	if commentHex == "" {
		return nil, app.ErrorMissingField
	}
	comment := model.Comment{}

	statement := `
		SELECT *
		FROM comments
		WHERE comments.commentHex = $1;
	`

	if err := r.db.Unsafe().Get(&comment, statement, commentHex); err != nil {
		// TODO: is this the only error?
		return nil, app.ErrorNoSuchComment
	}

	return &comment, nil
}

func (r *CommentRepositoryPg) GetCommentDomainPath(commentHex string) (string, string, error) {
	if commentHex == "" {
		return "", "", app.ErrorMissingField
	}

	statement := `
		SELECT domain, path
		FROM comments
		WHERE commentHex = $1;
	`
	row := r.db.QueryRow(statement, commentHex)

	var domain string
	var path string
	var err error
	if err = row.Scan(&domain, &path); err != nil {
		return "", "", app.ErrorNoSuchDomain
	}

	return domain, path, nil
}

func (r *CommentRepositoryPg) ApproveComment(commentHex string, url string) error {
	if commentHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		UPDATE comments
		SET state = 'approved'
		WHERE commentHex = $1;
	`

	_, err := r.db.Exec(statement, commentHex)
	if err != nil {
		util.GetLogger().Errorf("cannot approve comment: %v", err)
		return err
	}

	return nil
}

func (r *CommentRepositoryPg) VoteComment(commenterHex string, commentHex string, direction int, url string) error {
	if commentHex == "" || commenterHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		SELECT commenterHex
		FROM comments
		WHERE commentHex = $1;
	`
	row := r.db.QueryRow(statement, commentHex)

	var authorHex string
	if err := row.Scan(&authorHex); err != nil {
		util.GetLogger().Errorf("error selecting authorHex for vote")
		return err
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
	_, err := r.db.Exec(statement, commentHex, commenterHex, direction, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("error inserting/updating votes: %v", err)
		return err
	}

	return nil
}

func (r *CommentRepositoryPg) DeleteComment(commentHex string, deleterHex string, domain string, path string) error {
	if commentHex == "" || deleterHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		UPDATE comments
		SET
			deleted = true,
			markdown = '[deleted]',
			html = '[deleted]',
			commenterHex = 'anonymous',
			deleterHex = $2,
			deletionDate = $3
		WHERE commentHex = $1;
	`
	_, err := r.db.Exec(statement, commentHex, deleterHex, time.Now().UTC())

	if err != nil {
		// TODO: make sure this is the error is actually non-existant commentHex
		return app.ErrorNoSuchComment
	}

	// Since we're no longer actually deleting comments, we are no longer running the trigger function!
	statement = `
		UPDATE pages
		SET commentCount = commentCount - 1
		WHERE canon($1) LIKE canon(domain) AND path = $2;
	`
	_, err = r.db.Exec(statement, domain, path)

	if err != nil {
		return app.ErrorNoSuchComment
	}

	return nil
}

func (r *CommentRepositoryPg) ListComments(commenterHex string, domain string, path string, includeUnapproved bool) ([]*model.Comment, map[string]*model.Commenter, error) {
	//TODO: use join instead of 2N + 1 queries
	// path can be empty
	if commenterHex == "" || domain == "" {
		return nil, nil, app.ErrorMissingField
	}

	statement := `
		SELECT
			commentHex,
			commenterHex,
			markdown,
			html,
			parentHex,
			score,
			state,
			deleted,
			creationDate
		FROM comments
		WHERE
			canon($1) LIKE canon(comments.domain) AND
			comments.path = $2
	`

	if !includeUnapproved {
		if commenterHex == "anonymous" {
			statement += `AND state = 'approved'`
		} else {
			statement += `AND (state = 'approved' OR commenterHex = $3)`
		}
	}

	statement += `;`

	var rows *sql.Rows
	var err error

	if !includeUnapproved && commenterHex != "anonymous" {
		rows, err = r.db.Query(statement, domain, path, commenterHex)
	} else {
		rows, err = r.db.Query(statement, domain, path)
	}

	if err != nil {
		util.GetLogger().Errorf("cannot get comments: %v", err)
		return nil, nil, err
	}
	defer rows.Close()

	commenters := make(map[string]*model.Commenter)
	commenters["anonymous"] = &model.Commenter{CommenterHex: "anonymous", Email: "undefined", Name: "Anonymous", Link: "undefined", Photo: "undefined", Provider: "undefined"}

	comments := []*model.Comment{}
	for rows.Next() {
		comment := model.Comment{}
		if err = rows.Scan(
			&comment.CommentHex,
			&comment.CommenterHex,
			&comment.Markdown,
			&comment.Html,
			&comment.ParentHex,
			&comment.Score,
			&comment.State,
			&comment.Deleted,
			&comment.CreationDate); err != nil {
			return nil, nil, err
		}

		if commenterHex != "anonymous" {
			statement = `
				SELECT direction
				FROM votes
				WHERE commentHex=$1 AND commenterHex=$2;
			`
			row := r.db.QueryRow(statement, comment.CommentHex, commenterHex)

			if err = row.Scan(&comment.Direction); err != nil {
				// TODO: is the only error here that there is no such entry?
				comment.Direction = 0
			}
		}

		if commenterHex != comment.CommenterHex {
			comment.Markdown = ""
		}

		if !includeUnapproved {
			comment.State = ""
		}

		comments = append(comments, &comment)

		if _, ok := commenters[comment.CommenterHex]; !ok {
			commenters[comment.CommenterHex], err = Repo.CommenterRepository.GetCommenterByHex(comment.CommenterHex)
			if err != nil {
				util.GetLogger().Errorf("cannot retrieve commenter: %v", err)
				return nil, nil, err
			}
		}
	}

	return comments, commenters, nil
}

func (r *CommentRepositoryPg) LogDomainViewRecord(domain string, commenterHex string) {
	if os.Getenv("ENABLE_LOGGING") != "false" && os.Getenv("ENABLE_LOGGING") != "" {
		statement := `
			INSERT INTO
			views  (domain, commenterHex, viewDate)
			VALUES ($1,     $2,           $3      );
		`
		_, err := r.db.Exec(statement, domain, commenterHex, time.Now().UTC())

		if err != nil {
			util.GetLogger().Warningf("cannot insert views: %v", err)
		}
	}
}
