package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func SsoCallbackHandler(w http.ResponseWriter, r *http.Request) {
	payloadHex := r.FormValue("payload")
	signature := r.FormValue("hmac")

	payloadBytes, err := hex.DecodeString(payloadHex)
	if err != nil {
		fmt.Fprintf(w, "Error: invalid JSON payload hex encoding: %s\n", err.Error())
		return
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		fmt.Fprintf(w, "Error: invalid HMAC signature hex encoding: %s\n", err.Error())
		return
	}

	payload := ssoPayload{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		fmt.Fprintf(w, "Error: cannot unmarshal JSON payload: %s\n", err.Error())
		return
	}

	if payload.Token == "" || payload.Email == "" || payload.Name == "" {
		fmt.Fprintf(w, "Error: %s\n", app.ErrorMissingField.Error())
		return
	}

	if payload.Link == "" {
		payload.Link = "undefined"
	}

	if payload.Photo == "" {
		payload.Photo = "undefined"
	}

	domainName, commenterToken, err := ssoTokenExtract(payload.Token)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	domain, err := domainGet(domainName)
	if err != nil {
		if err == app.ErrorNoSuchDomain {
			fmt.Fprintf(w, "Error: %s\n", err.Error())
		} else {
			util.GetLogger().Errorf("cannot get domain for SSO: %v", err)
			fmt.Fprintf(w, "Error: %s\n", app.ErrorInternal.Error())
		}
		return
	}

	if domain.SsoSecret == "" || domain.SsoUrl == "" {
		fmt.Fprintf(w, "Error: %s\n", app.ErrorMissingConfig.Error())
		return
	}

	key, err := hex.DecodeString(domain.SsoSecret)
	if err != nil {
		util.GetLogger().Errorf("cannot decode SSO secret as hex: %v", err)
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	h := hmac.New(sha256.New, key)
	h.Write(payloadBytes)
	expectedSignatureBytes := h.Sum(nil)
	if !hmac.Equal(expectedSignatureBytes, signatureBytes) {
		fmt.Fprintf(w, "Error: HMAC signature verification failed\n")
		return
	}

	_, err = repository.Repo.CommenterRepository.GetCommenterByToken(commenterToken)
	if err != nil && err != app.ErrorNoSuchToken {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	c, err := repository.Repo.CommenterRepository.GetCommenterByEmail("sso:"+domainName, payload.Email)
	if err != nil && err != app.ErrorNoSuchCommenter {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	var commenterHex string

	if err == app.ErrorNoSuchCommenter {
		commenterHex, err = commenterNew(payload.Email, payload.Name, payload.Link, payload.Photo, "sso:"+domainName, "")
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
	} else {
		if err = commenterUpdate(c.CommenterHex, payload.Email, payload.Name, payload.Link, payload.Photo, "sso:"+domainName); err != nil {
			util.GetLogger().Warningf("cannot update commenter: %s", err)
			// not a serious enough to exit with an error
		}

		commenterHex = c.CommenterHex
	}

	if err = repository.Repo.CommenterRepository.UpdateCommenterSession(commenterToken, commenterHex); err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	fmt.Fprintf(w, "<html><script>window.parent.close()</script></html>")
}
