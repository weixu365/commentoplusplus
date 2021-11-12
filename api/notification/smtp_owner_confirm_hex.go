package notification

import (
	"bytes"
	"os"
	"simple-commenting/util"
)

type ownerConfirmHexPlugs struct {
	Origin     string
	ConfirmHex string
}

func smtpOwnerConfirmHex(to string, toName string, confirmHex string) error {
	var body bytes.Buffer
	templates["confirm-hex"].Execute(&body, &ownerConfirmHexPlugs{Origin: os.Getenv("ORIGIN"), ConfirmHex: confirmHex})

	err := smtpSendMail(to, toName, "", "Please confirm your email address", body.String())
	if err != nil {
		util.GetLogger().Errorf("cannot send confirmation email: %v", err)
		return errorCannotSendEmail
	}

	return nil
}
