package notification

import (
	"fmt"
	"os"
	"simple-commenting/app"
	"simple-commenting/util"
	"text/template"
)

var headerTemplate *template.Template

type headerPlugs struct {
	FromAddress string
	ToName      string
	ToAddress   string
	Subject     string
}

var templates map[string]*template.Template

func smtpTemplatesLoad() error {
	var err error
	headerTemplate, err = template.New("header").Parse(`MIME-Version: 1.0
From: Commento <{{.FromAddress}}>
To: {{.ToName}} <{{.ToAddress}}>
Content-Type: text/plain; charset=UTF-8
Subject: {{.Subject}}

`)
	if err != nil {
		util.GetLogger().Errorf("cannot parse header template: %v", err)
		return app.ErrorMalformedTemplate
	}

	names := []string{
		"confirm-hex",
		"reset-hex",
		"domain-export",
		"domain-export-error",
	}

	templates = make(map[string]*template.Template)

	util.GetLogger().Infof("loading templates: %v", names)
	for _, name := range names {
		var err error
		templates[name] = template.New(name)
		templates[name], err = template.ParseFiles(fmt.Sprintf("%s/templates/%s.txt", os.Getenv("STATIC"), name))
		if err != nil {
			util.GetLogger().Errorf("cannot parse %s/templates/%s.txt: %v", os.Getenv("STATIC"), name, err)
			return app.ErrorMalformedTemplate
		}
	}

	return nil
}
