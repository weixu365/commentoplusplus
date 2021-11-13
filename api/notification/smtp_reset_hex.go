package notification

import (
	"bytes"
	"os"
	"simple-commenting/app"
	"simple-commenting/util"
)

type resetHexPlugs struct {
	Origin   string
	ResetHex string
}

func smtpResetHex(to string, toName string, resetHex string) error {
	var body bytes.Buffer
	templates["reset-hex"].Execute(&body, &resetHexPlugs{Origin: os.Getenv("ORIGIN"), ResetHex: resetHex})

	err := smtpSendMail(to, toName, "", "Reset your password", body.String())
	if err != nil {
		util.GetLogger().Errorf("cannot send reset email: %v", err)
		return app.ErrorCannotSendEmail
	}

	return nil
}
