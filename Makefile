LOCAL_BIN := $(CURDIR)/bin
PATH := $(PATH):$(PWD)/bin
MIGRATIONS_DIR := $(CURDIR)/migrations

.PHONY: bin-deps
bin-deps:
	$(info Installing binary dependencies...)
	@mkdir -p $(LOCAL_BIN)
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15 && \
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: gen-proto-protoc
gen-proto-protoc:
	protoc -I . -I api \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=./pkg --go_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative \
	api/service.proto \
	easyp-demo-service/api/service.proto

.PHONY: generate
generate:
	@$(LOCAL_BIN)/easyp generate

.PHONY: lint
lint:
	@$(LOCAL_BIN)/easyp lint --path api

.PHONY: breaking
breaking:
	@$(LOCAL_BIN)/easyp breaking --against main --path api

.PHONY: migration
migration:
	@mkdir -p $(MIGRATIONS_DIR)
	@read -p "Migration name: " migration_name && \
	$(LOCAL_BIN)/goose -dir $(MIGRATIONS_DIR) create "$$migration_name" sql

.PHONY: start-infra
start-infra:
	$(info Starting Docker infrastructure and running migrations...)
	@docker compose --env-file infra.env up -d
	@./migration.sh

.PHONY: stop-infra
stop-infra:
	$(info Stopping Docker infrastructure...)
	@docker compose --env-file infra.env stop

.PHONY: clear-infra
clear-infra:
	$(info Stopping Docker infrastructure and removing volumes...)
	@docker compose --env-file infra.env down -v

.PHONY: docker-up
docker-up:
	@docker compose --env-file infra.env up -d

.PHONY: docker-down
docker-down:
	$(info Stopping Docker containers...)
	@docker compose --env-file infra.env down

.PHONY: run
run:
	@echo "Running application locally..."
	@go run ./cmd/grpc-easyp