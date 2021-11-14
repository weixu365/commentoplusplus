package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"

	"github.com/jmoiron/sqlx"
)

type CommenterRepositoryPg struct {
	db *sqlx.DB
}

func (r *CommenterRepositoryPg) GetCommenterByToken(commenterToken string) (*model.Commenter, error) {
	if commenterToken == "" {
		return nil, app.ErrorMissingField
	}

	commenter := model.Commenter{}
	statement := `
		SELECT *
		FROM commenterSessions
		JOIN commenters ON commenterSessions.commenterHex = commenters.commenterHex
		WHERE commenterToken = $1;
	`

	if err := r.db.Get(&commenter, statement, commenterToken); err != nil {
		// TODO: is this the only error?
		return nil, app.ErrorNoSuchToken
	}

	if commenter.CommenterHex == "none" {
		return nil, app.ErrorNoSuchToken
	}

	if commenter.Deleted {
		commenter.Email = "undefined"
		commenter.Name = "[deleted]"
		commenter.Link = "undefined"
		commenter.Photo = "undefined"
	}

	return &commenter, nil
}

func (r *CommenterRepositoryPg) GetCommenterByHex(commenterHex string) (*model.Commenter, error) {
	if commenterHex == "" {
		return nil, app.ErrorMissingField
	}

	commenter := model.Commenter{}
	statement := `
		SELECT *
		FROM commenters
		WHERE commenterHex = $1;
	`
	if err := r.db.Get(&commenter, statement, commenterHex); err != nil {
		// TODO: is this the only error?
		return nil, app.ErrorNoSuchCommenter
	}

	if commenter.Deleted {
		commenter.Email = "undefined"
		commenter.Name = "[deleted]"
		commenter.Link = "undefined"
		commenter.Photo = "undefined"
	}

	return &commenter, nil
}

func (r *CommenterRepositoryPg) GetCommenterByEmail(provider string, email string) (*model.Commenter, error) {
	if provider == "" || email == "" {
		return nil, app.ErrorMissingField
	}

	commenter := model.Commenter{}
	statement := `
		SELECT *
		FROM commenters
		WHERE email = $1 AND provider = $2;
	`

	if err := r.db.Get(&commenter, statement, email, provider); err != nil {
		// TODO: is this the only error?
		return nil, app.ErrorNoSuchCommenter
	}

	if commenter.Deleted {
		commenter.Email = "undefined"
		commenter.Name = "[deleted]"
		commenter.Link = "undefined"
		commenter.Photo = "undefined"
	}

	return &commenter, nil
}
