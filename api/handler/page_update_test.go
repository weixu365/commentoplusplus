package handler

import (
	"simple-commenting/model"
	"simple-commenting/repository"
	"simple-commenting/test"
	"testing"
	"time"
)

func TestPageUpdateBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")

	commentNew(commenterHex, "example.com", "/path.html", "root", "**foo**", "unapproved", time.Now().UTC())

	page, _ := repository.Repo.PageRepository.GetPageByPath("example.com", "/path.html")
	if page.IsLocked != false {
		t.Errorf("expected IsLocked=false got %v", page.IsLocked)
		return
	}

	page.IsLocked = true

	if err := repository.Repo.PageRepository.UpdatePage(page); err != nil {
		t.Errorf("unexpected error updating page: %v", err)
		return
	}

	page, _ = repository.Repo.PageRepository.GetPageByPath("example.com", "/path.html")
	if page.IsLocked != true {
		t.Errorf("expected IsLocked=true got %v", page.IsLocked)
		return
	}
}

func TestPageUpdateEmpty(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	page := model.Page{Domain: "", Path: "", IsLocked: false}
	if err := repository.Repo.PageRepository.UpdatePage(&page); err == nil {
		t.Errorf("expected error not found updating page with empty everything")
		return
	}
}
