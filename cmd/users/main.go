// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"go.uber.org/zap"
)

func main() {
	// Cria um logger de desenvolvimento
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to initialize logger")
	}
	defer logger.Sync()

	logger.Info("Starting the application...")

	// Parse da flag de configuração
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	if *configFlag == "" {
		logger.Warn("Config file not provided. Using default values.", zap.String("flag", "config"))
	}

	logger.Info("Loading the config file", zap.String("file", *configFlag))
	c, err := config.LoadConfig(*configFlag)
	if err != nil {
		logger.Fatal("Error loading the config file", zap.String("file", *configFlag), zap.Error(err))
	}

	// Criação da aplicação
	logger.Info("Initializing the application...")
	a := app.NewUser(&c.Users)

	// Registro das rotas
	logger.Info("Registering the routes...")
	a.RegisterRoutes(http.DefaultServeMux)

	// Inicialização do servidor
	logger.Info("Starting HTTP server...", zap.String("endpoint", c.Server.Endpoint.HTTP))
	err = http.ListenAndServe(c.Server.Endpoint.HTTP, http.DefaultServeMux)
	if err != nil {
		logger.Fatal("Error starting the server", zap.String("endpoint", c.Server.Endpoint.HTTP), zap.Error(err))
	}
}
