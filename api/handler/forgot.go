package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/notification"
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

	if !notification.SmtpConfigured {
		return app.ErrorSmtpNotConfigured
	}

	var hex string
	var name string
	if entity == "owner" {
		owner, err := repository.Repo.OwnerRepository.GetByEmail(email)
		if err != nil {
			if err == app.ErrorNoSuchEmail {
				// TODO: use a more random time instead.
				time.Sleep(1 * time.Second)
				return nil
			} else {
				util.GetLogger().Errorf("cannot get owner by email: %v", err)
				return app.ErrorInternal
			}
		}
		hex = owner.OwnerHex
		name = owner.Name
	} else {
		c, err := repository.Repo.CommenterRepository.GetCommenterByEmail("commento", email)
		if err != nil {
			if err == app.ErrorNoSuchEmail {
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

	resetHex, err := repository.Repo.ResetRepository.CreateResetHex(hex, entity)
	if err != nil {
		return err
	}

	err = notification.SmtpResetHex(email, name, resetHex)
	if err != nil {
		return err
	}

	return nil
}

func ForgotHandler(w http.ResponseWriter, r *http.Request) {
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
