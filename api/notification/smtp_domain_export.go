package notification

import (
	"bytes"
	"os"
	"simple-commenting/app"
	"simple-commenting/util"
)

type domainExportPlugs struct {
	Origin    string
	Domain    string
	ExportHex string
}

func smtpDomainExport(to string, toName string, domain string, exportHex string) error {
	var body bytes.Buffer
	templates["domain-export"].Execute(&body, &domainExportPlugs{Origin: os.Getenv("ORIGIN"), ExportHex: exportHex})

	err := smtpSendMail(to, toName, "", "Commento Data Export", body.String())
	if err != nil {
		util.GetLogger().Errorf("cannot send data export email: %v", err)
		return app.ErrorCannotSendEmail
	}

	return nil
}
