package handler

import (
	"os"
	"simple-commenting/util"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/gitlab"
)

var gitlabConfig *oauth2.Config

func gitlabOauthConfigure() error {
	gitlabConfig = nil
	if os.Getenv("GITLAB_KEY") == "" && os.Getenv("GITLAB_SECRET") == "" {
		return nil
	}

	if os.Getenv("GITLAB_KEY") == "" {
		util.GetLogger().Errorf("COMMENTO_GITLAB_KEY not configured, but COMMENTO_GITLAB_SECRET is set")
		return app.ErrorOauthMisconfigured
	}

	if os.Getenv("GITLAB_SECRET") == "" {
		util.GetLogger().Errorf("COMMENTO_GITLAB_SECRET not configured, but COMMENTO_GITLAB_KEY is set")
		return app.ErrorOauthMisconfigured
	}

	util.GetLogger().Infof("loading gitlab OAuth config")

	gitlabConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("ORIGIN") + "/api/oauth/gitlab/callback",
		ClientID:     os.Getenv("GITLAB_KEY"),
		ClientSecret: os.Getenv("GITLAB_SECRET"),
		Scopes: []string{
			"read_user",
		},
		Endpoint: gitlab.Endpoint,
	}
	gitlabConfig.Endpoint.AuthURL = os.Getenv("GITLAB_URL") + "/oauth/authorize"
	gitlabConfig.Endpoint.TokenURL = os.Getenv("GITLAB_URL") + "/oauth/token"

	gitlabConfigured = true

	return nil
}
