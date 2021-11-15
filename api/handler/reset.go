package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"

	"golang.org/x/crypto/bcrypt"
)

func reset(resetHex string, password string) (string, error) {
	if resetHex == "" || password == "" {
		return "", app.ErrorMissingField
	}

	hex, err := repository.Repo.ResetRepository.GetResetHex(resetHex)
	if err != nil {
		// TODO: is this the only error?
		return "", app.ErrorNoSuchResetToken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		util.GetLogger().Errorf("cannot generate hash from password: %v\n", err)
		return "", app.ErrorInternal
	}

	if hex.Entity == "owner" {
		err = repository.Repo.OwnerRepository.UpdatePassword(string(passwordHash), hex.Hex)
	} else {
		err = repository.Repo.CommenterRepository.UpdateCommenterPassword(string(passwordHash), hex.Hex)
	}

	if err != nil {
		util.GetLogger().Errorf("cannot change %s's password: %v\n", hex.Entity, err)
		return "", app.ErrorInternal
	}

	err = repository.Repo.ResetRepository.DeleteResetHex(resetHex)
	if err != nil {
		util.GetLogger().Warningf("cannot remove resetHex: %v\n", err)
	}

	return hex.Entity, nil
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ResetHex *string `json:"resetHex"`
		Password *string `json:"password"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	entity, err := reset(*x.ResetHex, *x.Password)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "entity": entity})
}
