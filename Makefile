.DEFAULT_GOAL := help

CONFIG_PATH ?= ./local/config.yaml

PLANS_ENDPOINT ?= `yq .subscriptions.plans_endpoint $(CONFIG_PATH)`
USERS_ENDPOINT ?= `yq .subscriptions.users_endpoint $(CONFIG_PATH)`
SUBSCRIPTIONS_ENDPOINT ?= `yq .payments.subscriptions_endpoint $(CONFIG_PATH)`
PAYMENTS_ENDPOINT ?= `yq .server.endpoint.http $(CONFIG_PATH)`

# CMD_START_NATS_SERVER := nats-server -a 127.0.0.1 -js
CMD_NATS_SERVER_START := docker run -d --name nats-server -p 4222:4222 nats:latest --jetstream
CMD_NATS_SERVER_STOP := docker stop nats-server && docker rm nats-server

CMD_CREATE_PAYMENT_STREAM := nats stream create payments --subjects "payment.process" --storage memory --replicas 1 --retention=limits --discard=old --max-msgs 1_000_000 --max-msgs-per-subject 100_000 --max-bytes 4GiB --max-age 1d --max-msg-size 10MiB --dupe-window 2m --allow-rollup --no-deny-delete --no-deny-purge

CMD_ALL_IN_ONE_START := go run ./cmd/all-in-one/ -config $(CONFIG_PATH)

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

.PHONY: start
start: start-docker

.PHONY: start-docker
start-docker:
	@echo Starting 'nats-server '
	$(CMD_NATS_SERVER_START)
	@echo Creating Payment Stream
	$(CMD_CREATE_PAYMENT_STREAM)
	@echo 'Run project (all in one)'
	$(CMD_ALL_IN_ONE_START)
	@echo Systems are running...

.PHONY: stop
stop:
	@echo Stop nats-server
	$(CMD_NATS_SERVER_STOP)

.PHONY: test-load
test-load:
# https://github.com/grafana/k6/releases
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