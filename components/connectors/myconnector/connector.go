package myconnector

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type myConnector struct {
	cfg    *Config
	logger *zap.Logger
	next   consumer.Logs
}

var _ connector.Traces = (*myConnector)(nil)

func NewConnector(cfg *Config, logger *zap.Logger, next consumer.Logs) (connector.Traces, error) {
	return &myConnector{
		cfg:    cfg,
		logger: logger,
		next:   next,
	}, nil
}

func (r *myConnector) Start(ctx context.Context, host component.Host) error {
	r.logger.Info("Starting my connector")

	return nil
}

func (r *myConnector) Shutdown(ctx context.Context) error {
	r.logger.Info("Shutting down my connector")
	return nil
}

func (r *myConnector) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

func (r *myConnector) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	r.logger.Info("Consuming traces", zap.Any("traces", td))
	td.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes().PutStr("myconnector", "passou aqui")
	ld := plog.NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	lr := sl.LogRecords().AppendEmpty()
	lr.Body().SetStr("passou pelo nosso connector")
	lr.Attributes().PutStr("myconnector", "passou aqui")
	return r.next.ConsumeLogs(ctx, ld)
}
