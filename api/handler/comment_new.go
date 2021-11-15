package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
	"strings"
	"time"
)

// Take `creationDate` as a param because comment import (from Disqus, for
// example) will require a custom time.
func commentNew(commenterHex string, domainName string, path string, parentHex string, markdown string, state string, creationDate time.Time) (string, error) {
	// path is allowed to be empty
	if commenterHex == "" || domainName == "" || parentHex == "" || markdown == "" || state == "" {
		return "", app.ErrorMissingField
	}

	page, err := repository.Repo.PageRepository.GetPageByPath(domainName, path)
	if err != nil {
		util.GetLogger().Errorf("cannot get page attributes: %v", err)
		return "", app.ErrorInternal
	}

	if page.IsLocked {
		return "", app.ErrorThreadLocked
	}

	commentHex, err := util.RandomHex(32)
	if err != nil {
		return "", err
	}

	html := util.MarkdownToHtml(markdown)

	if err = repository.Repo.PageRepository.CreatePage(domainName, path); err != nil {
		return "", err
	}

	comment := model.Comment{
		CommentHex:   commentHex,
		Domain:       domainName,
		Path:         path,
		CommenterHex: commenterHex,
		ParentHex:    parentHex,
		Markdown:     markdown,
		Html:         html,
		CreationDate: creationDate,
		State:        state,
	}

	if _, err = repository.Repo.CommentRepository.CreateComment(&comment); err != nil {
		return "", err
	}

	notification.NotificationHub.Broadcast <- []byte(domainName + path)

	return commentHex, nil
}

func CommentNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		AnonName       *string `json:"anonName"`
		Domain         *string `json:"domain"`
		Path           *string `json:"path"`
		ParentHex      *string `json:"parentHex"`
		Markdown       *string `json:"markdown"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := util.DomainStrip(*x.Domain)
	path := *x.Path

	d, err := domainGet(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if d.State == "frozen" {
		bodyMarshal(w, response{"success": false, "message": app.ErrorDomainFrozen.Error()})
		return
	}

	if d.RequireIdentification && *x.CommenterToken == "anonymous" {
		bodyMarshal(w, response{"success": false, "message": app.ErrorNotAuthorised.Error()})
		return
	}
	var state string

	var commenterHex, commenterName, commenterEmail, commenterLink string
	var isModerator bool

	if *x.CommenterToken == "anonymous" {
		commenterHex, commenterName, commenterEmail, commenterLink = "anonymous", "Anonymous", "", ""
		if util.IsSpam(*x.Domain, getIp(r), getUserAgent(r), "Anonymous", "", "", *x.Markdown) {
			state = "flagged"
		} else {
			// if given an anonName, add it to a new commenter entry
			if strings.TrimSpace(*x.AnonName) != "" {
				commenterHex, err = commenterNew("undefined", strings.TrimSpace(*x.AnonName), "undefined", "undefined", "anon", "undefined")
				if err != nil {
					bodyMarshal(w, response{"success": false, "message": err.Error()})
					return
				}
			}

			if d.ModerateAllAnonymous || d.RequireModeration {
				state = "unapproved"
			} else {
				state = "approved"
			}
		}
	} else {
		commenter, err := repository.Repo.CommenterRepository.GetCommenterByToken(*x.CommenterToken)
		if err != nil {
			bodyMarshal(w, response{"success": false, "message": err.Error()})
			return
		}
		commenterHex, commenterName, commenterEmail, commenterLink = commenter.CommenterHex, commenter.Name, commenter.Email, commenter.Link
		for _, mod := range *d.Moderators {
			if mod.Email == commenter.Email {
				isModerator = true
				break
			}
		}
	}

	if isModerator {
		state = "approved"
	} else if d.RequireModeration {
		state = "unapproved"
	} else if commenterHex == "anonymous" && d.ModerateAllAnonymous {
		state = "unapproved"
	} else if d.AutoSpamFilter && util.IsSpam(*x.Domain, getIp(r), getUserAgent(r), commenterName, commenterEmail, commenterLink, *x.Markdown) {
		state = "flagged"
	} else {
		state = "approved"
	}

	commentHex, err := commentNew(commenterHex, domain, path, *x.ParentHex, *x.Markdown, state, time.Now().UTC())
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// TODO: reuse html in commentNew and do only one markdown to HTML conversion?
	html := util.MarkdownToHtml(*x.Markdown)

	bodyMarshal(w, response{"success": true, "commentHex": commentHex, "state": state, "html": html})
	if notification.SmtpConfigured {
		go emailNotificationNew(d, path, commenterHex, commentHex, html, *x.ParentHex, state)
	}
}
