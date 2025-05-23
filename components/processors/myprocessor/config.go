package myprocessor

import "go.opentelemetry.io/collector/component"

var _ component.Config = (*Config)(nil)

type Config struct {
}

func (c *Config) Validate() error {
	return nil
}
