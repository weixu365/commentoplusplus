package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func CreateLogger() *logrus.Logger {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.Out = os.Stdout

	return logger
}

func GetLogger() *logrus.Logger {
	if logger == nil {
		logger = CreateLogger()
	}

	return logger
}
