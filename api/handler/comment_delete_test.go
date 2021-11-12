package handler

import (
	"simple-commenting/test"
	"testing"
	"time"
)

func TestCommentDeleteBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	commenterHex := "temp-commenter-hex"
	commentHex, _ := commentNew(commenterHex, "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())
	commentNew(commenterHex, "example.com", "/path.html", commentHex, "**bar**", "approved", time.Now().UTC())

	if err := commentDelete(commentHex, commenterHex, "example.com", "/path.html"); err != nil {
		t.Errorf("unexpected error deleting comment: %v", err)
		return
	}

	c, _, _ := commentList(commenterHex, "example.com", "/path.html", false)

	if len(c) != 0 {
		t.Errorf("expected no comments found %d comments", len(c))
		return
	}
}

func TestCommentDeleteEmpty(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if err := commentDelete("", "test-commenter-hex", "anydomain.com", "/path.html"); err == nil {
		t.Errorf("expected error deleting comment with empty commentHex")
		return
	}
}
