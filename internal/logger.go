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
	// Assume running on Google Cloud
	if gcpProject != "" {
		logger, _ = zapdriver.NewProduction()
		logger.Info("Using zapdriver logger")
		return logger
	}
	if env == "development" {
		logger, _ = zap.NewDevelopment()
		logger.Info("Using development logger")
		return logger
	} else {
		logger, _ = zap.NewProduction()
		logger.Info("Using production logger")
		return logger
	}

	return logger
}
