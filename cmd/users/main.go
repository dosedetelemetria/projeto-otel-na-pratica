// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	texp, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint("localhost:4318"))
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.Default()),
		sdktrace.WithBatcher(texp),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))),
		sdktrace.WithRawSpanLimits(sdktrace.SpanLimits{
			AttributeValueLengthLimit: 1000,
			EventCountLimit:           1000,
			LinkCountLimit:            1000,
		}),
	)

	mexp, err := otlpmetrichttp.New(context.Background(), otlpmetrichttp.WithInsecure(), otlpmetrichttp.WithEndpoint("localhost:4318"))
	if err != nil {
		log.Fatal(err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource.Default()),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(mexp)),
		sdkmetric.WithView(
			sdkmetric.NewView(
				sdkmetric.Instrument{
					Name: "payment.cart_size",
					Scope: instrumentation.Scope{
						Name: telemetry.Scope,
					},
				},
				sdkmetric.Stream{
					Name:        "payment.cart_sum",
					Aggregation: sdkmetric.AggregationSum{},
				},
			),
		),
		sdkmetric.WithExemplarFilter(func(ctx context.Context) bool {
			return false
		}),
	)

	lexp, err := otlploghttp.New(context.Background(), otlploghttp.WithInsecure(), otlploghttp.WithEndpoint("localhost:4318"))
	if err != nil {
		log.Fatal(err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(resource.Default()),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(lexp)),
	)

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	global.SetLoggerProvider(lp)

	// producao
	tr := otel.Tracer("github.com/dosedetelemetria/projeto-otel-na-pratica/users")
	ctx, span := tr.Start(context.Background(), "user.main")
	defer span.End()

	span.AddEvent("user.main")

	c, _ := config.LoadConfig(*configFlag)

	a := app.NewUser(ctx, &c.Users)
	a.RegisterRoutes(http.DefaultServeMux)
	_ = http.ListenAndServe(c.Server.Endpoint.HTTP, http.DefaultServeMux)
}
