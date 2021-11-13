package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"simple-commenting/app"
	"simple-commenting/util"
)

type commentoExportV1 struct {
	Version    int         `json:"version"`
	Comments   []comment   `json:"comments"`
	Commenters []commenter `json:"commenters"`
}

func domainImportCommento(domain string, url string) (int, error) {
	if domain == "" || url == "" {
		return 0, app.ErrorMissingField
	}

	resp, err := http.Get(url)
	if err != nil {
		util.Get, app.Error).Errorf("cannot get url: %v", err)
		return 0, app.ErrorCannotDownloadCommento
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.GetLogger().Errorf("cannot read body: %v", err)
		return 0, app.ErrorCannotDownloadCommento
	}

	zr, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		util.GetLogger().Errorf("cannot create gzip reader: %v", err)
		return 0, app.ErrorInternal
	}

	contents, err := ioutil.ReadAll(zr)
	if err != nil {
		util.GetLogger().Errorf("cannot read gzip contents uncompressed: %v", err)
		return 0, app.ErrorInternal
	}

	var data commentoExportV1
	if err := json.Unmarshal(contents, &data); err != nil {
		util.GetLogger().Errorf("cannot unmarshal JSON at %s: %v", url, err)
		return 0, app.ErrorInternal
	}

	if data.Version != 1 {
		util.GetLogger().Errorf("invalid data version (got %d, want 1): %v", data.Version, err)
		return 0, app.ErrorUnsupportedCommentoImportVersion
	}

	// Check if imported commentedHex or email exists, creating a map of
	// commenterHex (old hex, new hex)
	commenterHex := map[string]string{"anonymous": "anonymous"}
	for _, commenter := range data.Commenters {
		c, err := commenterGetByEmail("commento", commenter.Email)
		if err != nil && err != app.ErrorNoSuchCommenter {
			util.GetLogger().Errorf("cannot get commenter by email: %v", err)
			return 0, app.ErrorInternal
		}

		if err == nil {
			commenterHex[commenter.CommenterHex] = c.CommenterHex
			continue
		}

		randomPassword, err := util.RandomHex(32)
		if err != nil {
			util.GetLogger().Errorf("cannot generate random password for new commenter: %v", err)
			return 0, app.ErrorInternal
		}

		commenterHex[commenter.CommenterHex], err = commenterNew(commenter.Email,
			commenter.Name, commenter.Link, commenter.Photo, "commento", randomPassword)
		if err != nil {
			return 0, err
		}
	}

	// Create a map of (parent hex, comments)
	comments := make(map[string][]comment)
	for _, comment := range data.Comments {
		parentHex := comment.ParentHex
		comments[parentHex] = append(comments[parentHex], comment)
	}

	// Import comments, creating a map of comment hex (old hex, new hex)
	commentHex := map[string]string{"root": "root"}
	numImported := 0
	keys := []string{"root"}
	for i := 0; i < len(keys); i++ {
		for _, comment := range comments[keys[i]] {
			cHex, ok := commenterHex[comment.CommenterHex]
			if !ok {
				util.GetLogger().Errorf("cannot get commenter: %v", err)
				return numImported, app.ErrorInternal
			}
			parentHex, ok := commentHex[comment.ParentHex]
			if !ok {
				util.GetLogger().Errorf("cannot get parent comment: %v", err)
				return numImported, app.ErrorInternal
			}

			hex, err := commentNew(
				cHex,
				domain,
				comment.Path,
				parentHex,
				comment.Markdown,
				comment.State,
				comment.CreationDate)
			if err != nil {
				return numImported, err
			}
			commentHex[comment.CommentHex] = hex
			numImported++
			keys = append(keys, comment.CommentHex)
		}
	}

	return numImported, nil
}

func domainImportCommentoHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
		URL        *string `json:"url"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	numImported, err := domainImportCommento(domain, *x.URL)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "numImported": numImported})
}
