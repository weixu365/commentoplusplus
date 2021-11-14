package repository

import (
	"os"
	"simple-commenting/util"
	"time"

	"github.com/jmoiron/sqlx"
)

type LogRepository interface {
	LogDomainViewRecord(domain string, commenterHex string)
}

type LogRepositoryPg struct {
	db *sqlx.DB
}

func (r *LogRepositoryPg) LogDomainViewRecord(domain string, commenterHex string) {
	if os.Getenv("ENABLE_LOGGING") != "false" && os.Getenv("ENABLE_LOGGING") != "" {
		statement := `
			INSERT INTO
			views  (domain, commenterHex, viewDate)
			VALUES ($1,     $2,           $3      );
		`
		_, err := r.db.Exec(statement, domain, commenterHex, time.Now().UTC())

		if err != nil {
			util.GetLogger().Warningf("cannot insert views: %v", err)
		}
	}
}
