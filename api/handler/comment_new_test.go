package handler

import (
	"simple-commenting/repository"
	"simple-commenting/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommentNewBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if _, err := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC()); err != nil {
		t.Errorf("unexpected error creating new comment: %v", err)
		return
	}
}

func TestCommentNewEmpty(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if _, err := commentNew("temp-commenter-hex", "example.com", "", "root", "**foo**", "approved", time.Now().UTC()); err != nil {
		t.Errorf("empty path not allowed: %v", err)
		return
	}

	if _, err := commentNew("temp-commenter-hex", "", "", "root", "**foo**", "approved", time.Now().UTC()); err == nil {
		t.Errorf("expected error not found creatingn new comment with empty domain")
		return
	}

	if _, err := commentNew("", "", "", "", "", "", time.Now().UTC()); err == nil {
		t.Errorf("expected error not found creatingn new comment with empty everything")
		return
	}
}

func TestCommentNewUpvoted(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	commentHex, _ := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	statement := `
		SELECT score
		FROM comments
		WHERE commentHex = $1;
	`
	row := repository.Db.QueryRow(statement, commentHex)

	var score int
	if err := row.Scan(&score); err != nil {
		t.Errorf("error scanning score from comments table: %v", err)
		return
	}

	if score != 0 {
		t.Errorf("expected comment to be at 0 points")
		return
	}
}

func TestCommentNewThreadLocked(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	repository.Repo.PageRepository.CreatePage("example.com", "/path.html")
	p, _ := pageGet("example.com", "/path.html")
	p.IsLocked = true
	pageUpdate(p)

	_, err := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())
	if err == nil {
		t.Errorf("expected error not found creating a new comment on a locked thread")
		return
	}
}

func TestCommentDomainPathGetBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	commentHex, _ := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	domain, path, err := repository.Repo.CommentRepository.GetCommentDomainPath(commentHex)
	assert.NoError(t, err)

	assert.Equal(t, domain, "example.com")
	assert.Equal(t, path, "/path.html")
}
