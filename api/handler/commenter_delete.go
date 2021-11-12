package handler

import (
	"net/http"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func commenterDelete(commenterHex string) error {
	if commenterHex == "" {
		return app.ErrorMissingField
	}

	statement := `
		UPDATE commenters
		SET deleted=true
		WHERE commenterHex = $1;
	`
	_, err := repository.Db.Exec(statement, commenterHex)
	if err != nil {
		util.GetLogger().Errorf("cannot delete commenter: %v", err)
		return app.ErrorInternal
	}

	return nil
}

func commenterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if err = commenterDelete(c.CommenterHex); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
