package handler

import (
	"fmt"
	"net/http"
	"os"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func ownerConfirmHex(confirmHex string) error {
	if confirmHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		UPDATE owners
		SET confirmedEmail=true
		WHERE ownerHex IN (
			SELECT ownerHex FROM ownerConfirmHexes
			WHERE confirmHex=$1
		);
	`
	res, err := repository.Db.Exec(statement, confirmHex)
	if err != nil {
		util.GetLogger().Errorf("cannot mark user's confirmedEmail as true: %v\n", err)
		return app.ErrorInternal
	}

	count, err := res.RowsAffected()
	if err != nil {
		util.GetLogger().Errorf("cannot count rows affected: %v\n", err)
		return app.ErrorInternal
	}

	if count == 0 {
		return app.ErrorNoSuchConfirmationToken
	}

	statement = `
		DELETE FROM ownerConfirmHexes
		WHERE confirmHex=$1;
	`
	_, err = repository.Db.Exec(statement, confirmHex)
	if err != nil {
		util.GetLogger().Warningf("cannot remove confirmation token: %v\n", err)
		// Don't return an error because this is not critical.
	}

	return nil
}

func OwnerConfirmHexHandler(w http.ResponseWriter, r *http.Request) {
	if confirmHex := r.FormValue("confirmHex"); confirmHex != "" {
		if err := ownerConfirmHex(confirmHex); err == nil {
			http.Redirect(w, r, fmt.Sprintf("%s/login?confirmed=true", os.Getenv("ORIGIN")), http.StatusTemporaryRedirect)
			return
		}
	}

	// TODO: include error message in the URL
	http.Redirect(w, r, fmt.Sprintf("%s/login?confirmed=false", os.Getenv("ORIGIN")), http.StatusTemporaryRedirect)
}
