package myconnector

import (
	"context"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
)

func NewFactory() connector.Factory {
	extType, err := component.NewType("myconnector")
	if err != nil {
		log.Fatal(err)
	}

	return connector.NewFactory(
		extType,
		createDefaultConfig,
		connector.WithTracesToLogs(createTracesToLogs, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTracesToLogs(_ context.Context, set connector.Settings, cfg component.Config, next consumer.Logs) (connector.Traces, error) {
	return NewConnector(cfg.(*Config), set.Logger, next)
}
