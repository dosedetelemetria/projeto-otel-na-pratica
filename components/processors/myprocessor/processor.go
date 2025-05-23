package myprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
)

type myProcessor struct {
	cfg    *Config
	logger *zap.Logger
	next   consumer.Traces
}

var _ processor.Traces = (*myProcessor)(nil)

func NewProcessor(cfg *Config, logger *zap.Logger, next consumer.Traces) (processor.Traces, error) {
	return &myProcessor{
		cfg:    cfg,
		logger: logger,
		next:   next,
	}, nil
}

func (r *myProcessor) Start(ctx context.Context, host component.Host) error {
	r.logger.Info("Starting my processor")

	return nil
}

func (r *myProcessor) Shutdown(ctx context.Context) error {
	r.logger.Info("Shutting down my processor")
	return nil
}

func (r *myProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

func (r *myProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	r.logger.Info("Consuming traces", zap.Any("traces", td))
	td.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes().PutStr("myprocessor", "passou aqui")
	return r.next.ConsumeTraces(ctx, td)
}
