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
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	ctx := context.Background()
	db, err := gorm.Open(sqlite.Open(cfg.SQLLite.DSN))
	if err != nil {
		return nil, err
	}
	_ = db.AutoMigrate(&model.Payment{})

	nc, err := nats.Connect(cfg.NATS.Endpoint)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	stream, err := js.Stream(ctx, cfg.NATS.Stream)
	if err != nil {
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
		return nil, err
	}

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
