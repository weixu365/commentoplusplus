package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func SsoRedirectHandler(w http.ResponseWriter, r *http.Request) {
	commenterToken := r.FormValue("commenterToken")
	domainName := r.Header.Get("Referer")

	if commenterToken == "" {
		fmt.Fprintf(w, "Error: %s\n", app.ErrorMissingField.Error())
		return
	}

	domainName = util.DomainStrip(domainName)
	if domainName == "" {
		fmt.Fprintf(w, "Error: No Referer header found in request\n")
		return
	}

	_, err := repository.Repo.CommenterRepository.GetCommenterByToken(commenterToken)
	if err != nil && err != app.ErrorNoSuchToken {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	domain, err := domainGet(domainName)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", app.ErrorNoSuchDomain.Error())
		return
	}

	if !domain.SsoProvider {
		fmt.Fprintf(w, "Error: SSO not configured for %s\n", domainName)
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

	token, err := ssoTokenNew(domainName, commenterToken)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		util.GetLogger().Errorf("cannot decode hex token: %v", err)
		fmt.Fprintf(w, "Error: %s\n", app.ErrorInternal.Error())
		return
	}

	h := hmac.New(sha256.New, key)
	h.Write(tokenBytes)
	signature := hex.EncodeToString(h.Sum(nil))

	u, err := url.Parse(domain.SsoUrl)
	if err != nil {
		// this should really not be happening; we're checking if the
		// passed URL is valid at domain update
		util.GetLogger().Errorf("cannot parse URL: %v", err)
		fmt.Fprintf(w, "Error: %s\n", app.ErrorInternal.Error())
		return
	}

	q := u.Query()
	q.Set("token", token)
	q.Set("hmac", signature)
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusFound)
}
