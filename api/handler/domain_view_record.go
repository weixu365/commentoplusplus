package handler

import (
	"os"
	"simple-commenting/util"
	"time"
)

func domainViewRecord(domain string, commenterHex string) {
	if os.Getenv("ENABLE_LOGGING") != "false" && os.Getenv("ENABLE_LOGGING") != "" {
		statement := `
			INSERT INTO
			views  (domain, commenterHex, viewDate)
			VALUES ($1,     $2,           $3      );
		`
		_, err := db.Exec(statement, domain, commenterHex, time.Now().UTC())
		if err != nil {
			util.GetLogger().Warningf("cannot insert views: %v", err)
		}
	}
}
