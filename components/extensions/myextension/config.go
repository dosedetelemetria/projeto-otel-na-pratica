package myextension

import "go.opentelemetry.io/collector/component"

var _ component.Config = (*Config)(nil)

type Config struct {
	Url string `mapstructure:"url"`
}

func (c *Config) Validate() error {
	return nil
}
