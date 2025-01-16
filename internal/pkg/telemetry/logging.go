// Package telemetry provides utilities for logging and telemetry setup across the project.
package telemetry

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

// InitLogger initializes a global logger instance with the provided configuration.
// It should be called once at the application startup.
func InitLogger(developmentMode bool) *zap.Logger {
	var err error
	if developmentMode {
		logger, err = zap.NewDevelopment() // Development-friendly logger
	} else {
		logger, err = zap.NewProduction() // Production-optimized logger
	}

	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	return logger
}

// GetLogger retrieves the global logger instance.
// Make sure InitLogger is called before using this function.
func GetLogger() *zap.Logger {
	if logger == nil {
		log.Fatalf("logger not initialized: call InitLogger first")
	}
	return logger
}
