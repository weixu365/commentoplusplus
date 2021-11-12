package app

import (
	"os"
	"path/filepath"
	"simple-commenting/util"
	"strconv"
	"strings"
)

func configParse() error {
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		util.GetLogger().Errorf("cannot load binary path: %v", err)
		return err
	}

	defaults := map[string]string{
		"CONFIG_FILE": "",

		"POSTGRES": "postgres://postgres:postgres@localhost/commento?sslmode=disable",

		// PostgreSQL recommends max_connections in the order of hundreds. The default
		// is 100, so let's use half that and leave the other half for other services.
		// Ideally, you'd be setting this to a much higher number (for example, at the
		// time of writing, commento.io uses 600). See https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
		"MAX_IDLE_PG_CONNECTIONS": "50",

		"BIND_ADDRESS": "127.0.0.1",
		"PORT":         "8080",
		"ORIGIN":       "",

		"SSL":      "false",
		"SSL_CERT": "",
		"SSL_KEY":  "",

		"CDN_PREFIX": "",

		"FORBID_NEW_OWNERS": "false",

		"STATIC": binPath,

		"GZIP_STATIC": "false",

		"SMTP_USERNAME":         "",
		"SMTP_PASSWORD":         "",
		"SMTP_HOST":             "",
		"SMTP_PORT":             "",
		"SMTP_FROM_ADDRESS":     "",
		"SMTP_SKIP_HOST_VERIFY": "false",

		"AKISMET_KEY": "",

		"GOOGLE_KEY":    "",
		"GOOGLE_SECRET": "",

		"GITHUB_KEY":    "",
		"GITHUB_SECRET": "",

		"TWITTER_KEY":    "",
		"TWITTER_SECRET": "",

		"GITLAB_KEY":    "",
		"GITLAB_SECRET": "",
		"GITLAB_URL":    "https://gitlab.com",

		"ENABLE_LOGGING":   "false",
		"ENABLE_WILDCARDS": "true",
	}

	if os.Getenv("COMMENTO_CONFIG_FILE") != "" {
		if err := configFileLoad(os.Getenv("COMMENTO_CONFIG_FILE")); err != nil {
			return err
		}
	}

	for key, value := range defaults {
		if os.Getenv("COMMENTO_"+key) == "" {
			os.Setenv(key, value)
		} else {
			os.Setenv(key, os.Getenv("COMMENTO_"+key))
		}
	}

	// Mandatory config parameters
	for _, env := range []string{"POSTGRES", "PORT", "ORIGIN", "FORBID_NEW_OWNERS", "MAX_IDLE_PG_CONNECTIONS"} {
		if os.Getenv(env) == "" {
			util.GetLogger().Errorf("missing COMMENTO_%s environment variable", env)
			return ErrorMissingConfig
		}
	}

	os.Setenv("ORIGIN", strings.TrimSuffix(os.Getenv("ORIGIN"), "/"))
	os.Setenv("ORIGIN", util.AddHttpIfAbsent(os.Getenv("ORIGIN")))

	if os.Getenv("CDN_PREFIX") == "" {
		os.Setenv("CDN_PREFIX", os.Getenv("ORIGIN"))
	}

	os.Setenv("CDN_PREFIX", strings.TrimSuffix(os.Getenv("CDN_PREFIX"), "/"))
	os.Setenv("CDN_PREFIX", util.AddHttpIfAbsent(os.Getenv("CDN_PREFIX")))

	if os.Getenv("FORBID_NEW_OWNERS") != "true" && os.Getenv("FORBID_NEW_OWNERS") != "false" {
		util.GetLogger().Errorf("COMMENTO_FORBID_NEW_OWNERS neither 'true' nor 'false'")
		return ErrorInvalidConfigValue
	}

	static := os.Getenv("STATIC")
	for strings.HasSuffix(static, "/") {
		static = static[0 : len(static)-1]
	}

	file, err := os.Stat(static)
	if err != nil {
		util.GetLogger().Errorf("cannot load %s: %v", static, err)
		return err
	}

	if !file.IsDir() {
		util.GetLogger().Errorf("COMMENTO_STATIC=%s is not a directory", static)
		return ErrorNotADirectory
	}

	os.Setenv("STATIC", static)

	if num, err := strconv.Atoi(os.Getenv("MAX_IDLE_PG_CONNECTIONS")); err != nil {
		util.GetLogger().Errorf("invalid COMMENTO_MAX_IDLE_PG_CONNECTIONS: %v", err)
		return ErrorInvalidConfigValue
	} else if num <= 0 {
		util.GetLogger().Errorf("COMMENTO_MAX_IDLE_PG_CONNECTIONS should be a positive integer")
		return ErrorInvalidConfigValue
	}

	return nil
}
