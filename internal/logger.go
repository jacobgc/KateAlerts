package internal

import (
	"github.com/gridwise/zapdriver"
	"go.uber.org/zap"
	"os"
)

func GetLogger() *zap.Logger {
	gcpProject := os.Getenv("GCP_PROJECT")
	env := os.Getenv("env")
	var logger *zap.Logger
	var err error
	// Assume running on Google Cloud
	if gcpProject != "" {
		logger, err = zapdriver.NewProduction()
	} else if env == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	return logger
}
