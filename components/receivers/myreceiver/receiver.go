package myreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

type myReceiver struct {
	cfg    *Config
	logger *zap.Logger
	next   consumer.Traces
	ticker *time.Ticker
}

var _ receiver.Traces = (*myReceiver)(nil)

func NewReceiver(cfg *Config, logger *zap.Logger, next consumer.Traces) (receiver.Traces, error) {
	return &myReceiver{
		cfg:    cfg,
		logger: logger,
		next:   next,
		ticker: time.NewTicker(1 * time.Second),
	}, nil
}

func (r *myReceiver) Start(ctx context.Context, host component.Host) error {
	r.logger.Info("Starting my receiver")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-r.ticker.C:
				r.logger.Info("Sending trace")
				r.next.ConsumeTraces(ctx, ptrace.NewTraces())
			}
		}
	}()

	return nil
}

func (r *myReceiver) Shutdown(ctx context.Context) error {
	r.logger.Info("Shutting down my receiver")
	r.ticker.Stop()
	return nil
}
