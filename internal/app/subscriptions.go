// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	subscriptionhttp "github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/handler/http"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/store"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/store/memory"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Subscription struct {
	Handler *subscriptionhttp.SubscriptionHandler
	Store   store.Subscription
}

func NewSubscription(cfg *config.Subscriptions) *Subscription {
	store := memory.NewSubscriptionStore()
	return &Subscription{
		Handler: subscriptionhttp.NewSubscriptionHandler(store, cfg.UsersEndpoint, cfg.PlansEndpoint),
		Store:   store,
	}
}

func (a *Subscription) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /subscriptions", otelhttp.NewHandler(http.HandlerFunc(a.Handler.List), "GET /subscriptions"))
	mux.Handle("POST /subscriptions", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Create), "POST /subscriptions"))
	mux.Handle("GET /subscriptions/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Get), "GET /subscriptions/{id}"))
	mux.Handle("PUT /subscriptions/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Update), "PUT /subscriptions/{id}"))
	mux.Handle("DELETE /subscriptions/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Delete), "DELETE /subscriptions/{id}"))
}
