package handler

import (
	"os"
	"simple-commenting/util"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubConfig *oauth2.Config

func githubOauthConfigure() error {
	githubConfig = nil
	if os.Getenv("GITHUB_KEY") == "" && os.Getenv("GITHUB_SECRET") == "" {
		return nil
	}

	if os.Getenv("GITHUB_KEY") == "" {
		util.GetLogger().Errorf("COMMENTO_GITHUB_KEY not configured, but COMMENTO_GITHUB_SECRET is set")
		return app.ErrorOauthMisconfigured
	}

	if os.Getenv("GITHUB_SECRET") == "" {
		util.GetLogger().Errorf("COMMENTO_GITHUB_SECRET not configured, but COMMENTO_GITHUB_KEY is set")
		return app.ErrorOauthMisconfigured
	}

	util.GetLogger().Infof("loading github OAuth config")

	githubConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("ORIGIN") + "/api/oauth/github/callback",
		ClientID:     os.Getenv("GITHUB_KEY"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
		Scopes: []string{
			"read:user",
			"user:email",
		},
		Endpoint: github.Endpoint,
	}

	githubConfigured = true

	return nil
}
