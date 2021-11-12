package crons

import (
	"simple-commenting/util"
	"time"
)

func ssoTokenCleanupBegin() error {
	go func() {
		for {
			statement := `
				DELETE FROM ssoTokens
				WHERE creationDate < $1;
			`
			_, err := db.Exec(statement, time.Now().UTC().Add(time.Duration(-10)*time.Minute))
			if err != nil {
				util.GetLogger().Errorf("error cleaning up export rows: %v", err)
				return
			}

			time.Sleep(10 * time.Minute)
		}
	}()

	return nil
}
