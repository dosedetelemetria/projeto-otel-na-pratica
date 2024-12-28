COMPONENT_ID ?= api

.DEFAULT_GOAL := build

.PHONY: help
help:
	@echo ;
	@echo Tarefas comuns relacionadas à geração da API gRPC
	@echo ;

.PHONY: build
build: protoc

.PHONY: protoc
protoc:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./plan.proto

