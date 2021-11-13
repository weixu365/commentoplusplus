package handler

import (
	"encoding/json"
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"
)

func domainExportBeginError(email string, toName string, domain string, err error) {
	// we're not using err at the moment because it's all errorInternal
	if err2 := notification.SmtpDomainExportError(email, toName, domain); err2 != nil {
		util.GetLogger().Errorf("cannot send domain export error email for %s: %v", domain, err2)
		return
	}
}

func domainExportBegin(email string, toName string, domain string) {
	e := commentoExportV1{Version: 1, Comments: []model.Comment{}, Commenters: []model.Commenter{}}

	statement := `
		SELECT commentHex, domain, path, commenterHex, markdown, parentHex, score, state, creationDate
		FROM comments
		WHERE domain = $1;
	`
	rows1, err := repository.Db.Query(statement, domain)
	if err != nil {
		util.GetLogger().Errorf("cannot select comments while exporting %s: %v", domain, err)
		domainExportBeginError(email, toName, domain, app.ErrorInternal)
		return
	}
	defer rows1.Close()

	for rows1.Next() {
		c := model.Comment{}
		if err = rows1.Scan(&c.CommentHex, &c.Domain, &c.Path, &c.CommenterHex, &c.Markdown, &c.ParentHex, &c.Score, &c.State, &c.CreationDate); err != nil {
			util.GetLogger().Errorf("cannot scan comment while exporting %s: %v", domain, err)
			domainExportBeginError(email, toName, domain, app.ErrorInternal)
			return
		}

		e.Comments = append(e.Comments, c)
	}

	statement = `
		SELECT commenters.commenterHex, commenters.email, commenters.name, commenters.link, commenters.photo, commenters.provider, commenters.joinDate
		FROM commenters, comments
		WHERE comments.domain = $1 AND commenters.commenterHex = comments.commenterHex;
	`
	rows2, err := repository.Db.Query(statement, domain)
	if err != nil {
		util.GetLogger().Errorf("cannot select commenters while exporting %s: %v", domain, err)
		domainExportBeginError(email, toName, domain, app.ErrorInternal)
		return
	}
	defer rows2.Close()

	for rows2.Next() {
		c := model.Commenter{}
		if err := rows2.Scan(&c.CommenterHex, &c.Email, &c.Name, &c.Link, &c.Photo, &c.Provider, &c.JoinDate); err != nil {
			util.GetLogger().Errorf("cannot scan commenter while exporting %s: %v", domain, err)
			domainExportBeginError(email, toName, domain, app.ErrorInternal)
			return
		}

		e.Commenters = append(e.Commenters, c)
	}

	je, err := json.Marshal(e)
	if err != nil {
		util.GetLogger().Errorf("cannot marshall JSON while exporting %s: %v", domain, err)
		domainExportBeginError(email, toName, domain, app.ErrorInternal)
		return
	}

	gje, err := util.GzipStatic(je)
	if err != nil {
		util.GetLogger().Errorf("cannot gzip JSON while exporting %s: %v", domain, err)
		domainExportBeginError(email, toName, domain, app.ErrorInternal)
		return
	}

	exportHex, err := util.RandomHex(32)
	if err != nil {
		util.GetLogger().Errorf("cannot generate exportHex while exporting %s: %v", domain, err)
		domainExportBeginError(email, toName, domain, app.ErrorInternal)
		return
	}

	statement = `
		INSERT INTO
		exports (exportHex, binData, domain, creationDate)
		VALUES  ($1,        $2,      $3    , $4          );
	`
	_, err = repository.Db.Exec(statement, exportHex, gje, domain, time.Now().UTC())
	if err != nil {
		util.GetLogger().Errorf("error inserting expiry binary data while exporting %s: %v", domain, err)
		domainExportBeginError(email, toName, domain, app.ErrorInternal)
		return
	}

	err = notification.SmtpDomainExport(email, toName, domain, exportHex)
	if err != nil {
		util.GetLogger().Errorf("error sending data export email for %s: %v", domain, err)
		return
	}
}

func DomainExportBeginHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !notification.SmtpConfigured {
		bodyMarshal(w, response{"success": false, "message": app.ErrorSmtpNotConfigured.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	isOwner, err := domainOwnershipVerify(o.OwnerHex, *x.Domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotAuthorised.Error()})
		return
	}

	go domainExportBegin(o.Email, o.Name, *x.Domain)

	bodyMarshal(w, response{"success": true})
}
