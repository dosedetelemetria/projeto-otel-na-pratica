package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitTraces() {
	exp, err := otlptracegrpc.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(trace.NewTracerProvider(trace.WithBatcher(exp)))
}
