package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"simple-commenting/app"
	"simple-commenting/util"
	"time"
)

func versionPrint() error {
	util.GetLogger().Infof("starting Commento %s", app.Version)
	return nil
}

func versionCheckStart() error {
	go func() {
		printedError := false
		errorCount := 0
		latestSeen := ""

		for {
			time.Sleep(5 * time.Minute)

			data := url.Values{
				"version": {app.Version},
			}

			resp, err := http.Post("https://version.commento.io/api/check", "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
			if err != nil {
				errorCount++
				// print the error only once; we don't want to spam the logs with this
				// every five minutes
				if !printedError && errorCount > 5 {
					util.GetLogger().Errorf("error checking version: %v", err)
					printedError = true
				}
				continue
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errorCount++
				if !printedError && errorCount > 5 {
					util.GetLogger().Errorf("error reading body: %s", err)
					printedError = true
				}
				continue
			}

			type response struct {
				Success   bool   `json:"success"`
				Message   string `json:"message"`
				Latest    string `json:"latest"`
				NewUpdate bool   `json:"newUpdate"`
			}

			r := response{}
			json.Unmarshal(body, &r)
			if r.Success == false {
				errorCount++
				if !printedError && errorCount > 5 {
					util.GetLogger().Errorf("error checking version: %s", r.Message)
					printedError = true
				}
				continue
			}

			if r.NewUpdate && r.Latest != latestSeen {
				util.GetLogger().Infof("New update available! Latest version: %s", r.Latest)
				latestSeen = r.Latest
			}

			errorCount = 0
			printedError = false
		}
	}()

	return nil
}
