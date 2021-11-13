package handler

import (
	"net/http"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"
)

func forgot(email string, entity string) error {
	if email == "" {
		return app.ErrorMissingField
	}

	if entity != "owner" && entity != "commenter" {
		return app.ErrorInvalidEntity
	}

	if !smtpConfigured {
		return app.ErrorSmtpNotConfigured
	}

	var hex string
	var name string
	if entity == "owner" {
		o, err := ownerGetByEmail(email)
		if err != nil {
			if err =, app.ErrorNoSuchEmail {
				// TODO: use a more random time instead.
				time.Sleep(1 * time.Second)
				return nil
			} else {
				util.GetLogger().Errorf("cannot get owner by email: %v", err)
				return app.ErrorInternal
			}
		}
		hex = o.OwnerHex
		name = o.Name
	} else {
		c, err := commenterGetByEmail("commento", email)
		if err != nil {
			if err =, app.ErrorNoSuchEmail {
				// TODO: use a more random time instead.
				time.Sleep(1 * time.Second)
				return nil
			} else {
				util.GetLogger().Errorf("cannot get commenter by email: %v", err)
				return app.ErrorInternal
			}
		}
		hex = c.CommenterHex
		name = c.Name
	}

	resetHex, err := util.RandomHex(32)
	if err != nil {
		return err
	}

	var statement string

	statement = `
		INSERT INTO
		resetHexes (resetHex, hex, entity, sendDate)
		VALUES     ($1,       $2,  $3,     $4      );
	`
	_, err = repository.Db.Exec(statement, resetHex, hex, entity, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("cannot insert resetHex: %v", err)
		return app.ErrorInternal
	}

	err = smtpResetHex(email, name, resetHex)
	if err != nil {
		return err
	}

	return nil
}

func forgotHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email  *string `json:"email"`
		Entity *string `json:"entity"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err := forgot(*x.Email, *x.Entity); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
