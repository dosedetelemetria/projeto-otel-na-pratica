package myextension

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
)

type myExtension struct {
	cfg    *Config
	logger *zap.Logger
}

func NewExtension(cfg *Config, logger *zap.Logger) extension.Extension {
	return &myExtension{
		cfg:    cfg,
		logger: logger,
	}
}

func (e *myExtension) Start(ctx context.Context, host component.Host) error {
	e.logger.Info("Starting my extension", zap.String("url", e.cfg.Url))
	return nil
}

func (e *myExtension) Shutdown(ctx context.Context) error {
	e.logger.Info("Shutting down my extension")
	return nil
}
