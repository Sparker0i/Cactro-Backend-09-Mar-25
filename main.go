package main

import (
	"fmt"

	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/config"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/api/routes"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Configure logging
	configureLogging(cfg.LogLevel)

	// Setup router
	router := routes.SetupRouter(cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logrus.WithField("address", addr).Info("Starting server")
	if err := router.Run(addr); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}

// configureLogging sets up the logging configuration
func configureLogging(level string) {
	// Configure log format
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
	})

	// Set log level
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.WithField("level", level).Info("Log level set")
}
