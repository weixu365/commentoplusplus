package handler

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func domainModeratorList(domain string) ([]model.Moderator, error) {
	statement := `
		SELECT email, addDate
		FROM moderators
		WHERE domain=$1;
	`
	rows, err := repository.Db.Query(statement, domain)
	if err != nil {
		util.GetLogger().Errorf("cannot get moderators: %v", err)
		return nil, app.ErrorInternal
	}
	defer rows.Close()

	moderators := []model.Moderator{}
	for rows.Next() {
		m := model.Moderator{}
		if err = rows.Scan(&m.Email, &m.AddDate); err != nil {
			util.GetLogger().Errorf("cannot Scan moderator: %v", err)
			return nil, app.ErrorInternal
		}

		moderators = append(moderators, m)
	}

	return moderators, nil
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
