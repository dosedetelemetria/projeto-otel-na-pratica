.DEFAULT_GOAL := help

.PHONY: build
help:
	@echo Para construir o projeto, execute:
	@echo ;
	@echo "\t make build"
	@echo ;

.PHONY: install-tools
install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install github.com/nats-io/nats-server/v2@main