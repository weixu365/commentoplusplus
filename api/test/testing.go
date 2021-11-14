package test

import (
	"simple-commenting/notification"
	"simple-commenting/repository"
	"simple-commenting/util"
	"testing"

	"github.com/op/go-logging"
)

func FailTestOnError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("failed test: %v", err)
	}
}

var setupComplete bool

func SetupTestEnv() error {
	repository.SetupTestRepo()

	if !setupComplete {
		setupComplete = true

		util.GetLogger()

		// Print messages to console only if verbose. Sounds like a good idea to
		// keep the console clean on `go test`.
		if !testing.Verbose() {
			logging.SetLevel(logging.CRITICAL, "")
		}

		if err := util.MarkdownRendererCreate(); err != nil {
			return err
		}
	}

	notification.NotificationHub = notification.NewHub()

	return nil
}
