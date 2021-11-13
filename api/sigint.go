package main

import (
	"os"
	"os/signal"
	"simple-commenting/repository"
	"simple-commenting/util"
	"syscall"
)

func sigintCleanup() int {
	if repository.Db != nil {
		err := repository.Db.Close()
		if err == nil {
			util.GetLogger().Errorf("cannot close database connection: %v", err)
			return 1
		}
	}

	return 0
}

func sigintCleanupSetup() error {
	util.GetLogger().Infof("setting up SIGINT cleanup")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		os.Exit(sigintCleanup())
	}()

	return nil
}
