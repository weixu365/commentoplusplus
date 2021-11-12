package main

import "simple-commenting/util"

func main() {
	exitIfError(util.GetLogger())
	exitIfError(versionPrint())
	exitIfError(configParse())
	exitIfError(DbConnect(5))
	exitIfError(Migrate())
	exitIfError(smtpConfigure())
	exitIfError(smtpTemplatesLoad())
	exitIfError(oauthConfigure())
	exitIfError(MarkdownRendererCreate())
	exitIfError(sigintCleanupSetup())
	exitIfError(versionCheckStart())
	exitIfError(domainExportCleanupBegin())
	exitIfError(viewsCleanupBegin())
	exitIfError(ssoTokenCleanupBegin())

	exitIfError(routesServe())
}
