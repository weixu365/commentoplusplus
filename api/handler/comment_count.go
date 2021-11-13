package handler

import (
	"net/http"
	"simple-commenting/app"
	"simple-commenting/repository"
	"simple-commenting/util"

	"github.com/lib/pq"
)

func commentCount(domain string, paths []string) (map[string]int, error) {
	commentCounts := map[string]int{}

	if domain == "" {
		return nil, app.ErrorMissingField
	}

	if len(paths) == 0 {
		return nil, app.ErrorEmptyPaths
	}

	statement := `
		SELECT path, commentCount
		FROM pages
		WHERE canon($1) LIKE canon(domain) AND path = ANY($2);
	`
	rows, err := repository.Db.Query(statement, domain, pq.Array(paths))
	if err != nil {
		util.GetLogger().Errorf("cannot get comments: %v", err)
		return nil, app.ErrorInternal
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		var commentCount int
		if err = rows.Scan(&path, &commentCount); err != nil {
			util.GetLogger().Errorf("cannot scan path and commentCount: %v", err)
			return nil, app.ErrorInternal
		}

		commentCounts[path] = commentCount
	}

	return commentCounts, nil
}

func CommentCountHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Domain *string   `json:"domain"`
		Paths  *[]string `json:"paths"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := util.DomainStrip(*x.Domain)

	commentCounts, err := commentCount(domain, *x.Paths)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "commentCounts": commentCounts})
}
