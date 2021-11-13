package notification

import (
	"net/smtp"
	"os"
	"simple-commenting/app"
	"simple-commenting/util"
)

var SmtpConfigured bool
var smtpAuth smtp.Auth

func SmtpConfigure() error {
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if host == "" || port == "" {
		util.GetLogger().Warningf("smtp not configured, no emails will be sent")
		SmtpConfigured = false
		return nil
	}

	if os.Getenv("SMTP_FROM_ADDRESS") == "" {
		util.GetLogger().Errorf("COMMENTO_SMTP_FROM_ADDRESS not set")
		SmtpConfigured = false
		return app.ErrorMissingSmtpAddress
	}

	util.GetLogger().Infof("configuring smtp: %s", host)
	if username == "" || password == "" {
		util.GetLogger().Warningf("no SMTP username/password set, Commento will assume they aren't required")
	} else {
		smtpAuth = smtp.PlainAuth("", username, password, host)
	}
	SmtpConfigured = true
	return nil
}
