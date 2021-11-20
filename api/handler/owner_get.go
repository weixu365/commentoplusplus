package handler

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/util"
)

var ownersRowColumns string = `
	owners.ownerHex,
	owners.email,
	owners.name,
	owners.confirmedEmail,
	owners.joinDate
`

func ownersRowScan(s repository.SqlScanner, o *model.Owner) error {
	return s.Scan(
		&o.OwnerHex,
		&o.Email,
		&o.Name,
		&o.ConfirmedEmail,
		&o.JoinDate,
	)
}

func ownerGetByOwnerToken(ownerToken string) (model.Owner, error) {
	if ownerToken == "" {
		return model.Owner{}, app.ErrorMissingField
	}

	statement := `
		SELECT ` + ownersRowColumns + `
		FROM owners
		WHERE owners.ownerHex IN (
			SELECT ownerSessions.ownerHex FROM ownerSessions
			WHERE ownerSessions.ownerToken = $1
		);
	`
	row := repository.Db.QueryRow(statement, ownerToken)

	var o model.Owner
	if err := ownersRowScan(row, &o); err != nil {
		util.GetLogger().Errorf("cannot scan owner: %v\n", err)
		return model.Owner{}, app.ErrorInternal
	}

	return o, nil
}

func ownerGetByOwnerHex(ownerHex string) (model.Owner, error) {
	if ownerHex == "" {
		return model.Owner{}, app.ErrorMissingField
	}

	statement := `
		SELECT ` + ownersRowColumns + `
		FROM owners
		WHERE ownerHex = $1;
	`
	row := repository.Db.QueryRow(statement, ownerHex)

	var o model.Owner
	if err := ownersRowScan(row, &o); err != nil {
		util.GetLogger().Errorf("cannot scan owner: %v\n", err)
		return model.Owner{}, app.ErrorInternal
	}

	return o, nil
}
