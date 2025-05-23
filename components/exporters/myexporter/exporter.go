package myexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type myExporter struct {
	cfg    *Config
	logger *zap.Logger
}

var _ exporter.Traces = (*myExporter)(nil)

func NewExporter(cfg *Config, logger *zap.Logger) (exporter.Traces, error) {
	return &myExporter{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (r *myExporter) Start(ctx context.Context, host component.Host) error {
	r.logger.Info("Starting my exporter")

	return nil
}

func (r *myExporter) Shutdown(ctx context.Context) error {
	r.logger.Info("Shutting down my exporter")
	return nil
}

func (r *myExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

func (r *myExporter) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	td.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes().PutStr("myexporter", "passou aqui")
	r.logger.Info("Consuming traces", zap.Int("spans", td.SpanCount()))
	return nil
}
