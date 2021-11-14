package handler

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func domainModeratorList(domainName string) (*[]model.Moderator, error) {
	return repository.Repo.DomainModeratorRepository.GetModeratorsForDomain(domainName)
}

func isDomainModerator(domain string, email string) (bool, error) {
	statement := `
		SELECT EXISTS (
			SELECT 1
			FROM moderators
			WHERE domain=$1 AND email=$2
		);
	`
	row := repository.Db.QueryRow(statement, domain, email)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		util.GetLogger().Errorf("cannot query if moderator: %v", err)
		return false, app.ErrorInternal
	}

	return exists, nil
}
