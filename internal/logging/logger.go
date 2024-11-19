package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}

