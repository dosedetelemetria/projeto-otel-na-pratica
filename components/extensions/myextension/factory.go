package myextension

import (
	"context"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

func NewFactory() extension.Factory {
	extType, err := component.NewType("myextension")
	if err != nil {
		log.Fatal(err)
	}

	return extension.NewFactory(
		extType,
		createDefaultConfig,
		create,
		component.StabilityLevelDevelopment,
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Url: "http://dose.example.com",
	}
}

func create(_ context.Context, set extension.Settings, cfg component.Config) (extension.Extension, error) {
	return NewExtension(cfg.(*Config), set.Logger), nil
}
