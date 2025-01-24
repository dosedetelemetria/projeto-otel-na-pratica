package telemetry

import (
	"context"
	"os"

	config "go.opentelemetry.io/contrib/config/v0.3.0"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
)

func Setup(ctx context.Context, confFlag string) (func(context.Context) error, error) {
	b, err := os.ReadFile(confFlag)
	if err != nil {
		return nil, err
	}

	// interpolate the environment variables
	b = []byte(os.ExpandEnv(string(b)))

	// parse the config
	conf, err := config.ParseYAML(b)
	if err != nil {
		return nil, err
	}
	sdk, err := config.NewSDK(config.WithContext(ctx), config.WithOpenTelemetryConfiguration(*conf))
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(sdk.TracerProvider())
	otel.SetMeterProvider(sdk.MeterProvider())
	global.SetLoggerProvider(sdk.LoggerProvider())
	return sdk.Shutdown, nil
}
