package handler

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
)

var commentersRowColumns string = `
	commenters.commenterHex,
	commenters.email,
	commenters.name,
	commenters.link,
	commenters.photo,
	commenters.provider,
	commenters.joinDate,
	commenters.deleted
`

func commentersRowScan(s repository.SqlScanner, c *model.Commenter) error {
	return s.Scan(
		&c.CommenterHex,
		&c.Email,
		&c.Name,
		&c.Link,
		&c.Photo,
		&c.Provider,
		&c.JoinDate,
		&c.Deleted,
	)
}

func commenterGetByHex(commenterHex string) (model.Commenter, error) {
	if commenterHex == "" {
		return model.Commenter{}, app.ErrorMissingField
	}

	statement := `
		SELECT ` + commentersRowColumns + `
		FROM commenters
		WHERE commenterHex = $1;
	`
	row := repository.Db.QueryRow(statement, commenterHex)

	var c model.Commenter
	if err := commentersRowScan(row, &c); err != nil {
		// TODO: is this the only error?
		return model.Commenter{}, app.ErrorNoSuchCommenter
	}

	if c.Deleted == true {
		c.Email = "undefined"
		c.Name = "[deleted]"
		c.Link = "undefined"
		c.Photo = "undefined"
	}

	return c, nil
}

func commenterGetByEmail(provider string, email string) (model.Commenter, error) {
	if provider == "" || email == "" {
		return model.Commenter{}, app.ErrorMissingField
	}

	statement := `
		SELECT ` + commentersRowColumns + `
		FROM commenters
		WHERE email = $1 AND provider = $2;
	`
	row := repository.Db.QueryRow(statement, email, provider)

	var c model.Commenter
	if err := commentersRowScan(row, &c); err != nil {
		// TODO: is this the only error?
		return model.Commenter{}, app.ErrorNoSuchCommenter
	}

	if c.Deleted == true {
		c.Email = "undefined"
		c.Name = "[deleted]"
		c.Link = "undefined"
		c.Photo = "undefined"
	}

	return c, nil
}

func commenterGetByCommenterToken(commenterToken string) (model.Commenter, error) {
	if commenterToken == "" {
		return model.Commenter{}, app.ErrorMissingField
	}

	statement := `
		SELECT ` + commentersRowColumns + `
		FROM commenterSessions
		JOIN commenters ON commenterSessions.commenterHex = commenters.commenterHex
		WHERE commenterToken = $1;
	`
	row := repository.Db.QueryRow(statement, commenterToken)

	var c model.Commenter
	if err := commentersRowScan(row, &c); err != nil {
		// TODO: is this the only error?
		return model.Commenter{}, app.ErrorNoSuchToken
	}

	if c.CommenterHex == "none" {
		return model.Commenter{}, app.ErrorNoSuchToken
	}

	if c.Deleted == true {
		c.Email = "undefined"
		c.Name = "[deleted]"
		c.Link = "undefined"
		c.Photo = "undefined"
	}

	return c, nil
}
