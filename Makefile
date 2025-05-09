.DEFAULT_GOAL := help

CONFIG_PATH ?= ./local/config.yaml

PLANS_ENDPOINT ?= `yq .subscriptions.plans_endpoint $(CONFIG_PATH)`
USERS_ENDPOINT ?= `yq .subscriptions.users_endpoint $(CONFIG_PATH)`
SUBSCRIPTIONS_ENDPOINT ?= `yq .payments.subscriptions_endpoint $(CONFIG_PATH)`
PAYMENTS_ENDPOINT ?= `yq .server.endpoint.http $(CONFIG_PATH)`

CMD_K6_RUN ?= K6_WEB_DASHBOARD=true k6 run -e PLANS_ENDPOINT=$(PLANS_ENDPOINT) -e USERS_ENDPOINT=$(USERS_ENDPOINT) -e SUBSCRIPTIONS_ENDPOINT=$(SUBSCRIPTIONS_ENDPOINT) -e PAYMENTS_ENDPOINT=$(PAYMENTS_ENDPOINT) ./tests/k6/k6s.js

.PHONY: help
help:
	@echo Para construir o projeto, execute:
	@echo ;
	@echo "\t make build"
	@echo ;

.PHONY: lint
lint:
	@golangci-lint --timeout 30s run ./... --show-stats

.PHONY: test
test:
	@go test -v ./...

.PHONY: vulncheck
vulncheck:
	@govulncheck ./...

.PHONY: vulnfix
vulnfix:
	@go get -t -u

.PHONY: build
build: build-users build-payments build-plans build-subscriptions

.PHONY: build-users
build-users:
	@goreleaser build -f=.goreleaser.yaml --snapshot --clean --single-target --id users

.PHONY: build-payments
build-payments:
	@goreleaser build -f=.goreleaser.yaml --snapshot --clean --single-target --id payments

.PHONY: build-subscriptions
build-subscriptions:
	@goreleaser build -f=.goreleaser.yaml --snapshot --clean --single-target --id subscriptions

.PHONY: build-plans
build-plans:
	@goreleaser build -f=.goreleaser.yaml --snapshot --clean --single-target --id plans

.PHONY: test-load
test-load:
	@echo Run K6 to test services
	@echo "$(CMD_K6_RUN)"
	$(CMD_K6_RUN)

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install github.com/nats-io/nats-server/v2@main
	@go install github.com/nats-io/natscli/nats@latest
	@go install github.com/goreleaser/goreleaser/v2@latest