package main

import "simple-commenting/util"

func main() {
	exitIfError(util.GetLogger())
	exitIfError(versionPrint())
	exitIfError(configParse())
	exitIfError(dbConnect(5))
	exitIfError(migrate())
	exitIfError(smtpConfigure())
	exitIfError(smtpTemplatesLoad())
	exitIfError(oauthConfigure())
	exitIfError(markdownRendererCreate())
	exitIfError(sigintCleanupSetup())
	exitIfError(versionCheckStart())
	exitIfError(domainExportCleanupBegin())
	exitIfError(viewsCleanupBegin())
	exitIfError(ssoTokenCleanupBegin())

	exitIfError(routesServe())
}
