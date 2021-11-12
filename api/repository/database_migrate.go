package repository

import (
	"io/ioutil"
	"os"
	"simple-commenting/util"
	"strings"
)

var goMigrations = map[string](func() error){
	"20190213033530-email-notifications.sql": MigrateEmails,
}

func Migrate() error {
	return MigrateFromDir(os.Getenv("STATIC") + "/db")
}

func MigrateFromDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		util.GetLogger().Errorf("cannot read directory for migrations: %v", err)
		return err
	}

	statement := `
		SELECT filename
		FROM migrations;
	`
	rows, err := repository.Db.Query(statement)
	if err != nil {
		util.GetLogger().Errorf("cannot query migrations: %v", err)
		return err
	}

	defer rows.Close()

	filenames := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err = rows.Scan(&filename); err != nil {
			util.GetLogger().Errorf("cannot scan filename: %v", err)
			return err
		}

		filenames[filename] = true
		util.GetLogger().Infof("Found applied db script: %s", filename)
	}

	util.GetLogger().Infof("%d migrations already installed, looking for more", len(filenames))

	completed := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			if !filenames[file.Name()] {
				f := dir + string(os.PathSeparator) + file.Name()
				contents, err := ioutil.ReadFile(f)
				if err != nil {
					util.GetLogger().Errorf("cannot read file %s: %v", file.Name(), err)
					return err
				}

				if _, err = repository.Db.Exec(string(contents)); err != nil {
					util.GetLogger().Errorf("cannot execute the SQL in %s: %v", f, err)
					return err
				}

				statement = `
					INSERT INTO
					migrations (filename)
					VALUES     ($1      );
				`
				_, err = repository.Db.Exec(statement, file.Name())
				if err != nil {
					util.GetLogger().Errorf("cannot insert filename into the migrations table: %v", err)
					return err
				}

				if fn, ok := goMigrations[file.Name()]; ok {
					if err = fn(); err != nil {
						util.GetLogger().Errorf("cannot execute Go migration associated with SQL %s: %v", f, err)
						return err
					}
				}

				completed++
			}
		}
	}

	if completed > 0 {
		util.GetLogger().Infof("%d new migrations completed (%d total)", completed, len(filenames)+completed)
	} else {
		util.GetLogger().Infof("none found")
	}

	return nil
}
