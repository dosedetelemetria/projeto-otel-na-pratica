// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/model"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/store"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// PaymentHandler is an HTTP handler that performs CRUD operations for model.Payment using a store.Payment
type PaymentHandler struct {
	store                 store.Payment
	js                    jetstream.JetStream
	jsSubject             string
	subscriptionsEndpoint string

	cartSizeHistogram metric.Float64Histogram
	numPayments       metric.Int64Counter
}

// NewPaymentHandler returns a new PaymentHandler
func NewPaymentHandler(store store.Payment, js jetstream.JetStream, jsSubject string, subscriptionsEndpoint string) *PaymentHandler {
	cartSizeHistogram := internal.Must(telemetry.Meter().Float64Histogram("payment.cart_size",
		metric.WithDescription("The size of the cart in a payment."),
		metric.WithUnit("{price}"),
	))

	numPayments := internal.Must(telemetry.Meter().Int64Counter("payment.num_payments",
		metric.WithDescription("The number of payments."),
		metric.WithUnit("{payment}"),
	))

	return &PaymentHandler{
		store:                 store,
		js:                    js,
		jsSubject:             jsSubject,
		subscriptionsEndpoint: subscriptionsEndpoint,

		cartSizeHistogram: cartSizeHistogram,
		numPayments:       numPayments,
	}
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("tenant", r.Header.Get("Tenant")))

	payments, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(payments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payment model.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tenant := r.Header.Get("Tenant")

	h.cartSizeHistogram.Record(r.Context(),
		payment.Amount,
		metric.WithAttributeSet(attribute.NewSet(attribute.String("tenant", tenant))),
	)

	h.numPayments.Add(r.Context(), 1)

	// Check if subscription exists
	sub, err := otelhttp.Get(r.Context(), h.subscriptionsEndpoint+"/"+payment.SubscriptionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer sub.Body.Close()
	if sub.StatusCode != http.StatusOK {
		http.Error(w, "Subscription not found", http.StatusBadRequest)
		return
	}

	payload, err := json.Marshal(payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hs := nats.Header{
		"tenant": []string{"dose"},
	}

	otel.GetTextMapPropagator().Inject(r.Context(), carrier{hs})

	_, err = h.js.PublishMsgAsync(&nats.Msg{
		Subject: h.jsSubject,
		Data:    payload,
		Header:  hs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	payment, err := h.store.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if payment == nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PaymentHandler) Update(w http.ResponseWriter, r *http.Request) {
	payment := &model.Payment{}
	if err := json.NewDecoder(r.Body).Decode(payment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := h.store.Update(r.Context(), payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.store.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PaymentHandler) OnMessage(msg jetstream.Msg) {
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier{msg.Headers()})
	ctx, span := otel.Tracer("payments").Start(ctx, "onMessage")
	defer span.End()

	payment := &model.Payment{}
	err := json.Unmarshal(msg.Data(), payment)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	_, err = h.store.Create(ctx, payment)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	err = msg.Ack()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
}

type carrier struct {
	nats.Header
}

func (c carrier) Keys() []string {
	keys := make([]string, 0, len(c.Header))
	for k := range c.Header {
		keys = append(keys, k)
	}
	return keys
}
