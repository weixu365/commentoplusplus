package util

import (
	"os"

	"github.com/adtac/go-akismet/akismet"
)

func isSpam(domain string, userIp string, userAgent string, name string, email string, url string, markdown string) bool {
	akismetKey := os.Getenv("AKISMET_KEY")
	if akismetKey == "" {
		return false
	}

	res, err := akismet.Check(&akismet.Comment{
		Blog:               domain,
		UserIP:             userIp,
		UserAgent:          userAgent,
		CommentType:        "comment",
		CommentAuthor:      name,
		CommentAuthorEmail: email,
		CommentAuthorURL:   url,
		CommentContent:     markdown,
	}, akismetKey)

	if err != nil {
		util.GetLogger().Errorf("error: cannot validate commenet using Akismet: %v", err)
		return true
	}

	return res
}
