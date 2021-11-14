package repository

import (
	"fmt"
	"os"
	"path"
	"simple-commenting/util"
	"testing"

	"github.com/op/go-logging"
)

func SetupTestRepo() {
	if !setupComplete {
		setupComplete = true

		util.GetLogger()

		// Print messages to console only if verbose. Sounds like a good idea to
		// keep the console clean on `go test`.
		if !testing.Verbose() {
			logging.SetLevel(logging.CRITICAL, "")
		}

		if err := setupTestDatabase(FindRootFolder()); err != nil {
			panic(err)
		}
	}

	if err := clearTables(); err != nil {
		panic(err)
	}

	var err error
	if Repo, err = NewPostgresqlRepositories(os.Getenv("POSTGRES")); err != nil {
		panic(err)
	}
}

func getPublicTables() ([]string, error) {
	statement := `
		SELECT tablename
		FROM pg_tables
		WHERE schemaname='public';
	`
	rows, err := Db.Query(statement)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot query public tables: %v", err)
		return []string{}, err
	}

	defer rows.Close()

	tables := []string{}
	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			fmt.Fprintf(os.Stderr, "cannot scan table name: %v", err)
			return []string{}, err
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func DropTables() error {
	tables, err := getPublicTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		if table != "migrations" {
			_, err = Db.Exec(fmt.Sprintf("DROP TABLE %s;", table))
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot drop %s: %v", table, err)
				return err
			}
		}
	}

	return nil
}

func setupTestDatabase(rootPath string) error {
	if os.Getenv("COMMENTO_POSTGRES") != "" {
		// set it manually because we need to use commento_test, not commento, by mistake
		os.Setenv("POSTGRES", os.Getenv("COMMENTO_POSTGRES"))
	} else {
		os.Setenv("POSTGRES", "postgres://postgres:postgres@localhost/commento_test?sslmode=disable")
	}

	if err := DbConnect(0); err != nil {
		return err
	}

	if err := DropTables(); err != nil {
		return err
	}

	if err := MigrateFromDir(path.Join(rootPath, "db")); err != nil {
		return err
	}

	return nil
}

func clearTables() error {
	tables, err := getPublicTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		_, err = Db.Exec(fmt.Sprintf("DELETE FROM %s;", table))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot clear %s: %v", table, err)
			return err
		}
	}

	return nil
}

var setupComplete bool

func FindRootFolder() string {
	currentFolder, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	for {
		fi, err := os.Stat(path.Join(currentFolder, "db/Makefile"))

		if err == nil {
			switch mode := fi.Mode(); {
			case mode.IsRegular():
				return currentFolder
			}
		}

		parentFolder := path.Dir(currentFolder)

		if parentFolder == "/" || parentFolder == currentFolder {
			panic("Couldn't find project root folder, please check if the file path is correct")
		}

		currentFolder = parentFolder
	}
}
