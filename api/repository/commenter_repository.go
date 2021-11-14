package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)


type CommenterRepository interface {
	CreateCommenterToken() (string, error)
	GetCommenterHex(commenterToken string) (string, error)
	UpdateCommenterSession(commenterToken string, commenterHex string) error
	CreateCommenter(email string, name string, link string, photo string, provider string, password string) (string, error)
	GetCommenterByEmail(provider string, email string) (*model.Commenter, error)
	GetCommenterByHex(commenterHex string) (*model.Commenter, error)
	GetCommenterByToken(commenterToken string) (*model.Commenter, error)
}

type CommenterRepositoryPg struct {
	db *sqlx.DB
}

func (r CommenterRepositoryPg) CreateCommenterToken() (string, error) {
	commenterToken, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("cannot create commenterToken: %v", err)
		return "", app.ErrorInternal
	}

	statement := `
		INSERT INTO
		commenterSessions (commenterToken, creationDate)
		VALUES            ($1,             $2          );
	`
	_, err = r.db.Exec(statement, commenterToken, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert new commenterToken: %v", err)
		return "", app.ErrorInternal
	}

	return commenterToken, nil
}

func (r CommenterRepositoryPg) UpdateCommenterSession(commenterToken string, commenterHex string) error {
	if commenterToken == "" || commenterHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		UPDATE commenterSessions
		SET commenterHex = $2
		WHERE commenterToken = $1;
	`
	_, err := r.db.Exec(statement, commenterToken, commenterHex)
	if err != nil {
		util.GetLogger().Errorf("error updating commenterHex: %v", err)
		return app.ErrorInternal
	}

	return nil
}


func (r CommenterRepositoryPg) GetCommenterHex(commenterToken string) (string, error) {
	statement := `
		SELECT commenterHex
		FROM commenterSessions
		WHERE commenterToken = $1;
	`
	row := r.db.QueryRow(statement, commenterToken)
	
	var commenterHex string
	if err := row.Scan(&commenterHex); err != nil {
		return "", err
	}

	return commenterHex, nil
}

func (r CommenterRepositoryPg) CreateCommenter(email string, name string, link string, photo string, provider string, password string) (string, error) {
	commenterHex, err := util.RandomHex(32)
	if err != nil {
		return "", app.ErrorInternal
	}

	passwordHash := []byte{}
	if password != "" {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			util.GetLogger().Errorf("cannot generate hash from password: %v\n", err)
			return "", app.ErrorInternal
		}
	}

	if provider == "anon" {
		passwordHash = []byte{}
	}

	statement := `
		INSERT INTO
		commenters (commenterHex, email, name, link, photo, provider, passwordHash, joinDate)
		VALUES     ($1,           $2,    $3,   $4,   $5,    $6,       $7,           $8      );
	`
	_, err = r.db.Exec(statement, commenterHex, email, name, link, photo, provider, string(passwordHash), time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert commenter: %v", err)
		return "", app.ErrorInternal
	}

	return commenterHex, nil
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

	if err := r.db.Unsafe().Get(&commenter, statement, commenterToken); err != nil {
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
	if err := r.db.Unsafe().Get(&commenter, statement, commenterHex); err != nil {
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

	if err := r.db.Unsafe().Get(&commenter, statement, email, provider); err != nil {
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
