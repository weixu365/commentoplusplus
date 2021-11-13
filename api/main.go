package main

import (
	"fmt"
	"os"
	"simple-commenting/app"
	"simple-commenting/handler"
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func exitIfError(err error) {
	if err != nil {
		fmt.Printf("fatal error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	util.GetLogger()
	exitIfError(versionPrint())
	exitIfError(app.ConfigParse())
	exitIfError(repository.DbConnect(5))
	exitIfError(repository.Migrate())
	exitIfError(notification.SmtpConfigure())
	exitIfError(notification.SmtpTemplatesLoad())
	exitIfError(handler.OauthConfigure())
	exitIfError(util.MarkdownRendererCreate())
	exitIfError(sigintCleanupSetup())
	exitIfError(versionCheckStart())
	exitIfError(domainExportCleanupBegin())
	exitIfError(viewsCleanupBegin())
	exitIfError(ssoTokenCleanupBegin())

	exitIfError(routesServe())
}
