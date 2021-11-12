package handler

import (
	"simple-commenting/repository"
	"simple-commenting/util"
)

func pageNew(domain string, path string) error {
	// path can be empty
	if domain == "" {
		return errorMissingField
	}

	statement := `
		INSERT INTO
		pages  (domain, path)
		VALUES ($1,     $2  )
		ON CONFLICT DO NOTHING;
	`
	_, err := repository.Db.Exec(statement, domain, path)
	if err != nil {
		util.GetLogger().Errorf("error inserting new page: %v", err)
		return errorInternal
	}

	return nil
}
