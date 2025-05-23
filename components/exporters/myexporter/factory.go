package myexporter

import (
	"context"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

func NewFactory() exporter.Factory {
	extType, err := component.NewType("myexporter")
	if err != nil {
		log.Fatal(err)
	}

	return exporter.NewFactory(
		extType,
		createDefaultConfig,
		exporter.WithTraces(createTraces, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTraces(_ context.Context, set exporter.Settings, cfg component.Config) (exporter.Traces, error) {
	return NewExporter(cfg.(*Config), set.Logger)
}
