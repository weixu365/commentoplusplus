package handler

import (
	"simple-commenting/model"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func emailNotificationModerator(d *model.Domain, path string, title string, commenterHex string, commentHex string, html string, state string) {
	if d.EmailNotificationPolicy == "none" {
		return
	}

	if d.EmailNotificationPolicy == "pending-moderation" && state == "approved" {
		return
	}

	var commenterName string
	var commenterEmail string
	if commenterHex == "anonymous" {
		commenterName = "Anonymous"
	} else {
		commenter, err := repository.Repo.CommenterRepository.GetCommenterByHex(commenterHex)
		if err != nil {
			util.GetLogger().Errorf("cannot get commenter to send email notification: %v", err)
			return
		}

		commenterName = commenter.Name
		commenterEmail = commenter.Email
	}

	kind := d.EmailNotificationPolicy
	if state != "approved" {
		kind = "pending-moderation"
	}

	for _, m := range *d.Moderators {
		// Do not email the commenting moderator their own comment.
		if commenterHex != "anonymous" && m.Email == commenterEmail {
			continue
		}

		email, err := repository.Repo.EmailRepository.GetEmail(m.Email)
		if err != nil {
			// No such email.
			continue
		}

		if !email.SendModeratorNotifications {
			continue
		}

		commenter, err := repository.Repo.CommenterRepository.GetCommenterByEmail1(m.Email)
		if err != nil {
			// The moderator has probably not created a commenter account.
			// We should only send emails to people who signed up, so skip.
			continue
		}

		if err := notification.SmtpEmailNotification(m.Email, commenter.Name, kind, d.Domain, path, commentHex, commenterName, title, html, email.UnsubscribeSecretHex); err != nil {
			util.GetLogger().Errorf("error sending email to %s: %v", m.Email, err)
			continue
		}
	}
}

func emailNotificationReply(d *model.Domain, path string, title string, commenterHex string, commentHex string, html string, parentHex string, state string) {
	// No reply notifications for root comments.
	if parentHex == "root" {
		return
	}

	// No reply notification emails for unapproved comments.
	if state != "approved" {
		return
	}

	statement := `
		SELECT commenterHex
		FROM comments
		WHERE commentHex = $1;
	`
	row := repository.Db.QueryRow(statement, parentHex)

	var parentCommenterHex string
	err := row.Scan(&parentCommenterHex)
	if err != nil {
		util.GetLogger().Errorf("cannot scan commenterHex and parentCommenterHex: %v", err)
		return
	}

	// No reply notification emails for anonymous users.
	if parentCommenterHex == "anonymous" {
		return
	}

	// No reply notification email for self replies.
	if parentCommenterHex == commenterHex {
		return
	}

	parentCommenter, err := repository.Repo.CommenterRepository.GetCommenterByHex(parentCommenterHex)
	if err != nil {
		util.GetLogger().Errorf("cannot get commenter to send email notification: %v", err)
		return
	}

	var commenterName string
	if commenterHex == "anonymous" {
		commenterName = "Anonymous"
	} else {
		commenter, err := repository.Repo.CommenterRepository.GetCommenterByHex(commenterHex)
		if err != nil {
			util.GetLogger().Errorf("cannot get commenter to send email notification: %v", err)
			return
		}
		commenterName = commenter.Name
	}

	parentCommenterEmail, err := repository.Repo.EmailRepository.GetEmail(parentCommenter.Email)
	if err != nil {
		// No such email.
		return
	}

	if !parentCommenterEmail.SendReplyNotifications {
		return
	}

	notification.SmtpEmailNotification(parentCommenter.Email, parentCommenter.Name, "reply", d.Domain, path, commentHex, commenterName, title, html, parentCommenterEmail.UnsubscribeSecretHex)
}

func emailNotificationNew(d *model.Domain, path string, commenterHex string, commentHex string, html string, parentHex string, state string) {
	page, err := repository.Repo.PageRepository.GetPageByPath(d.Domain, path)
	if err != nil {
		util.GetLogger().Errorf("cannot get page to send email notification: %v", err)
		return
	}

	if page.Title == "" {
		page.Title, err = pageTitleUpdate(d.Domain, path)
		if err != nil {
			// Not being able to update a page title isn't serious enough to skip an
			// email notification.
			page.Title = d.Domain
		}
	}

	emailNotificationModerator(d, path, page.Title, commenterHex, commentHex, html, state)
	emailNotificationReply(d, path, page.Title, commenterHex, commentHex, html, parentHex, state)
}
