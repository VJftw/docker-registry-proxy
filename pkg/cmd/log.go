package cmd

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// Logger represents the command line logger
var Logger *zap.Logger
var once sync.Once

func init() {
	once.Do(func() {
		var config zap.Config
		if strings.ToLower(os.Getenv("DEBUG")) == "true" {
			config = zap.NewDevelopmentConfig()
		} else {
			config = zap.NewProductionConfig()
		}
		Logger, _ = config.Build()
	})
}
