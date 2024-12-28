COMPONENT_ID ?= payments

BUILD_DATE := `date -u +'%Y%m%dT%H%M%S'`
LAST_SHA := `git log --format=%h -n 1`

IMAGE_NAME ?= otelnapratica/$(COMPONENT_ID)
IMAGE_TAG ?= $(LAST_SHA).$(BUILD_DATE)

.DEFAULT_GOAL := all

.PHONY: help
help:
	@echo ;
	@echo Aplicação contendo todos os serviços do sistema
	@echo ;

.PHONY: all
all: lint test vulncheck build

.PHONY: lint
lint:
	@golangci-lint --timeout 20s run ./... --show-stats

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
build:
	@go build -o ./$(COMPONENT_ID) .

