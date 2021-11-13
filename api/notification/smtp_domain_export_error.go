package notification

import (
	"bytes"
	"os"
	"simple-commenting/app"
	"simple-commenting/util"
)

type domainExportErrorPlugs struct {
	Origin string
	Domain string
}

func SmtpDomainExportError(to string, toName string, domain string) error {
	var body bytes.Buffer
	templates["data-export-error"].Execute(&body, &domainExportPlugs{Origin: os.Getenv("ORIGIN")})

	err := smtpSendMail(to, toName, "", "Commento Data Export", body.String())
	if err != nil {
		util.GetLogger().Errorf("cannot send data export error email: %v", err)
		return app.ErrorCannotSendEmail
	}

	return nil
}
