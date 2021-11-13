package crons

import (
	"simple-commenting/repository"
	"simple-commenting/util"
	"time"
)

func ViewsCleanupBegin() error {
	go func() {
		for {
			statement := `
				DELETE FROM views
				WHERE viewDate < $1;
			`
			_, err := repository.Db.Exec(statement, time.Now().UTC().AddDate(0, 0, -45))
			if err != nil {
				util.GetLogger().Errorf("error cleaning up views: %v", err)
				return
			}

			time.Sleep(24 * time.Hour)
		}
	}()

	return nil
}
