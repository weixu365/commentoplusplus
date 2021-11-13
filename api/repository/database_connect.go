package repository

import (
	"database/sql"
	"net/url"
	"os"
	"simple-commenting/util"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func DbConnect(retriesLeft int) error {
	con := os.Getenv("POSTGRES")
	u, err := url.Parse(con)
	if err != nil {
		util.GetLogger().Errorf("invalid postgres connection URI: %v", err)
		return err
	}
	u.User = url.UserPassword(u.User.Username(), "redacted")
	util.GetLogger().Infof("opening connection to postgres: %s", u.String())

	db, err = sql.Open("postgres", con)
	if err != nil {
		util.GetLogger().Errorf("cannot open connection to postgres: %v", err)
		return err
	}

	err = Db.Ping()
	if err != nil {
		if retriesLeft > 0 {
			util.GetLogger().Errorf("cannot talk to postgres, retrying in 10 seconds (%d attempts left): %v", retriesLeft-1, err)
			time.Sleep(10 * time.Second)
			return DbConnect(retriesLeft - 1)
		} else {
			util.GetLogger().Errorf("cannot talk to postgres, last attempt failed: %v", err)
			return err
		}
	}

	statement := `
		CREATE TABLE IF NOT EXISTS migrations (
			filename TEXT NOT NULL UNIQUE
		);
	`
	_, err = Db.Exec(statement)
	if err != nil {
		util.GetLogger().Errorf("cannot create migrations table: %v", err)
		return err
	}

	maxIdleConnections, err := strconv.Atoi(os.Getenv("MAX_IDLE_PG_CONNECTIONS"))
	if err != nil {
		util.GetLogger().Warningf("cannot parse COMMENTO_MAX_IDLE_PG_CONNECTIONS: %v", err)
		maxIdleConnections = 50
	}

	db.SetMaxIdleConns(maxIdleConnections)

	return nil
}
