package myprocessor

import (
	"context"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

func NewFactory() processor.Factory {
	extType, err := component.NewType("myprocessor")
	if err != nil {
		log.Fatal(err)
	}

	return processor.NewFactory(
		extType,
		createDefaultConfig,
		processor.WithTraces(createTraces, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTraces(_ context.Context, set processor.Settings, cfg component.Config, next consumer.Traces) (processor.Traces, error) {
	return NewProcessor(cfg.(*Config), set.Logger, next)
}
