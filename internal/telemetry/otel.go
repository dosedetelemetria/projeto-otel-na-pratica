package telemetry

import (
	"context"
	"os"

	"go.opentelemetry.io/contrib/otelconf/v0.3.0"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const Scope = "projeto-otel-na-pratica"

func Setup(ctx context.Context, confFlag string) (func(context.Context) error, error) {
	b, err := os.ReadFile(confFlag)
	if err != nil {
		return nil, err
	}

	// interpolate the environment variables
	b = []byte(os.ExpandEnv(string(b)))

	// parse the config
	conf, err := otelconf.ParseYAML(b)
	if err != nil {
		return nil, err
	}
	sdk, err := otelconf.NewSDK(otelconf.WithContext(ctx), otelconf.WithOpenTelemetryConfiguration(*conf))
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	otel.SetTracerProvider(sdk.TracerProvider())
	otel.SetMeterProvider(sdk.MeterProvider())
	global.SetLoggerProvider(sdk.LoggerProvider())
	return sdk.Shutdown, nil
}

func Tracer() trace.Tracer {
	return otel.Tracer(Scope)
}

func Meter() metric.Meter {
	return otel.Meter(Scope)
}
