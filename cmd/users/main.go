// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/telemetry"
)

func main() {
	// Cria um logger de desenvolvimento
	telemetry.InitLogger(false)
	logger := telemetry.GetLogger()
	defer logger.Sync()

	logger.Info("Starting the application...")

	// Parse da flag de configuração
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	if *configFlag == "" {
		logger.Warn("Config file not provided. Using default values.")
	}

	logger.Info("Loading the config file")
	c, err := config.LoadConfig(*configFlag)
	if err != nil {
		logger.Fatal("Error loading the config file" + err.Error())
	}

	// Criação da aplicação
	logger.Info("Initializing the application...")
	a := app.NewUser(&c.Users)

	// Registro das rotas
	logger.Info("Registering the routes...")
	a.RegisterRoutes(http.DefaultServeMux)

	// Inicialização do servidor
	logger.Info("Starting HTTP server..." + c.Server.Endpoint.HTTP)
	err = http.ListenAndServe(c.Server.Endpoint.HTTP, http.DefaultServeMux)
	if err != nil {
		logger.Fatal("Error starting the server" + err.Error())
	}
}
