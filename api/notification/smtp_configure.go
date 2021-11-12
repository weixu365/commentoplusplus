package notification

import (
	"net/smtp"
	"os"
	"simple-commenting/util"
)

var smtpConfigured bool
var smtpAuth smtp.Auth

func smtpConfigure() error {
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if host == "" || port == "" {
		util.GetLogger().Warningf("smtp not configured, no emails will be sent")
		smtpConfigured = false
		return nil
	}

	if os.Getenv("SMTP_FROM_ADDRESS") == "" {
		util.GetLogger().Errorf("COMMENTO_SMTP_FROM_ADDRESS not set")
		smtpConfigured = false
		return app.ErrorMissingSmtpAddress
	}

	util.GetLogger().Infof("configuring smtp: %s", host)
	if username == "" || password == "" {
		util.GetLogger().Warningf("no SMTP username/password set, Commento will assume they aren't required")
	} else {
		smtpAuth = smtp.PlainAuth("", username, password, host)
	}
	smtpConfigured = true
	return nil
}
