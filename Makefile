.DEFAULT_GOAL := help

CONFIG_PATH ?= ./local/config.yaml

PLANS_ENDPOINT ?= `yq .subscriptions.plans_endpoint $(CONFIG_PATH)`
USERS_ENDPOINT ?= `yq .subscriptions.users_endpoint $(CONFIG_PATH)`
SUBSCRIPTIONS_ENDPOINT ?= `yq .payments.subscriptions_endpoint $(CONFIG_PATH)`
PAYMENTS_ENDPOINT ?= `yq .server.endpoint.http $(CONFIG_PATH)`

CMD_NATS_SERVER_START_BIN := nats-server -a 127.0.0.1 -js &
CMD_NATS_SERVER_STOP_BIN := killall nats-server

CMD_CREATE_PAYMENT_STREAM := nats stream create payments --subjects "payment.process" --storage memory --replicas 1 --retention=limits --discard=old --max-msgs 1_000_000 --max-msgs-per-subject 100_000 --max-bytes 4GiB --max-age 1d --max-msg-size 10MiB --dupe-window 2m --allow-rollup --no-deny-delete --no-deny-purge

CMD_ALL_IN_ONE_START := go run ./cmd/all-in-one/ -config $(CONFIG_PATH)

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

.PHONY: start
start: start-bin

.PHONY: start-bin
start-bin:
	@echo ;
	@echo Starting 'nats-server '
	$(CMD_NATS_SERVER_START_BIN)
	@echo ;
	@echo Creating Payment Stream
	$(CMD_CREATE_PAYMENT_STREAM)
	@echo ;
	@echo 'Run project (all in one)'
	$(CMD_ALL_IN_ONE_START)

.PHONY: stop
stop:
	@echo Stop nats-server
	$(CMD_NATS_SERVER_STOP_BIN)
	@echo ;


.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install github.com/nats-io/nats-server/v2@main
	@go install github.com/nats-io/natscli/nats@latest
	@go install github.com/goreleaser/goreleaser/v2@latest