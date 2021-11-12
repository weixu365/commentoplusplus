package crons

import (
	"simple-commenting/util"
	"time"
)

func domainExportCleanupBegin() error {
	go func() {
		for {
			statement := `
				DELETE FROM exports
				WHERE creationDate < $1;
			`
			_, err := db.Exec(statement, time.Now().UTC().AddDate(0, 0, -7))
			if err != nil {
				util.GetLogger().Errorf("error cleaning up export rows: %v", err)
				return
			}

			time.Sleep(2 * time.Hour)
		}
	}()

	return nil
}
