package handler

import (
	"os"
	"simple-commenting/repository"
	"simple-commenting/util"
)

func pageTitleUpdate(domain string, path string) (string, error) {
	ssl := os.Getenv("SSL")
	pre := ""
	if ssl == "true" {
		pre = "https://"
	} else {
		pre = "http://"
	}
	title, err := util.HtmlTitleGet(pre + domain + path)
	if err != nil {
		// This could fail due to a variety of reasons that we can't control such
		// as the user's URL 404 or something, so let's not pollute the error log
		// with messages. Just use a sane title. Maybe we'll have the ability to
		// retry later.
		util.GetLogger().Errorf("%v", err)
		title = domain
	}

	err = repository.Repo.PageRepository.UpdatePageTitle(domain, path, title)
	if err != nil {
		util.GetLogger().Errorf("cannot update pages table with title: %v", err)
		return "", err
	}

	return title, nil
}
