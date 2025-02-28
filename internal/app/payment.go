// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	planhttp "github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/handler/http"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/model"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/store"
	storegorm "github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/store/gorm"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Payment struct {
	Handler  *planhttp.PaymentHandler
	Store    store.Payment
	natsConn *nats.Conn
	cctx     jetstream.ConsumeContext
}

func NewPayment(cfg *config.Payments) (*Payment, error) {
	ctx, span := telemetry.Tracer().Start(context.Background(), "NewPayment",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attribute.String("consumer", cfg.NATS.ConsumerName)),
	)
	defer span.End()

	span.SetAttributes(attribute.String("consumer", cfg.NATS.ConsumerName))

	span.AddEvent("connecting to database")
	db, err := gorm.Open(sqlite.Open(cfg.SQLLite.DSN))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	err = db.AutoMigrate(&model.Payment{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	nc, err := nats.Connect(cfg.NATS.Endpoint)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	stream, err := js.Stream(ctx, cfg.NATS.Stream)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// this is only relevant for the consumer
	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          cfg.NATS.ConsumerName,
		Durable:       cfg.NATS.ConsumerName,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	store := storegorm.NewPaymentStore(db)
	pmt := &Payment{
		Handler:  planhttp.NewPaymentHandler(store, js, cfg.NATS.Subject, cfg.SubscriptionsEndpoint),
		Store:    store,
		natsConn: nc,
	}

	pmt.cctx, err = cons.Consume(pmt.Handler.OnMessage)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.AddEvent("payment service created")
	return pmt, nil
}

func (a *Payment) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /payments", otelhttp.NewHandler(http.HandlerFunc(a.Handler.List), "GET /payments"))
	mux.Handle("POST /payments", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Create), "POST /payments"))
	mux.Handle("GET /payments/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Get), "GET /payments/{id}"))
	mux.Handle("PUT /payments/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Update), "PUT /payments/{id}"))
	mux.Handle("DELETE /payments/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Delete), "DELETE /payments/{id}"))
}

func (a *Payment) Shutdown() error {
	if a.cctx != nil {
		a.cctx.Drain()
	}
	return a.natsConn.Drain()
}
