-include .env

MIGRATE ?= migrate
MIGRATIONS_PATH ?= ./cmd/migrate/migrations
MIGRATIONS_SOURCE ?= file://$(MIGRATIONS_PATH)
APP_BIN ?= ./bin/main
export DB_ADDR

.PHONY: help api build run test tidy docker-up docker-down migration-create migration-up migration-down migration-down-all migration-force migration-version migration-goto migration-drop

help:
	@echo "Available targets:"
	@echo "  make api                           Run API server (go run ./cmd/api)"
	@echo "  make build                         Build API binary to ./bin/main"
	@echo "  make run                           Run built binary"
	@echo "  make test                          Run go tests"
	@echo "  make tidy                          Run go mod tidy"
	@echo "  make docker-up                     Start docker compose services"
	@echo "  make docker-down                   Stop docker compose services"
	@echo "  make migration-create name=...     Create a new sequential SQL migration"
	@echo "  make migration-up                  Apply all up migrations"
	@echo "  make migration-down [n=1]          Roll back n migrations (default 1)"
	@echo "  make migration-down-all            Roll back all migrations"
	@echo "  make migration-force version=...   Force migration version"
	@echo "  make migration-version             Show current migration version"
	@echo "  make migration-goto version=...    Migrate to a specific version"
	@echo "  make migration-drop                Drop everything in DB (asks confirmation)"

api:
	go run ./cmd/api

seed:
	go run ./cmd/migrate/seed/main.go
	
build:
	go build -o $(APP_BIN) ./cmd/api

run:
	$(APP_BIN)

test:
	go test -v ./...

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down

gen-docs:
	@swag init -g main.go -d cmd/api,internals && swag fmt

migration-create:
	$(if $(strip $(name)),,$(error name is required: make migration-create name=create_users))
	$(MIGRATE) create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)

migration-up:
	$(MIGRATE) -source $(MIGRATIONS_SOURCE) -database "$(DB_ADDR)" up

migration-down:
	$(MIGRATE) -source $(MIGRATIONS_SOURCE) -database "$(DB_ADDR)" down $(or $(n),1)

migration-down-all:
	$(MIGRATE) -source $(MIGRATIONS_SOURCE) -database "$(DB_ADDR)" down -all

migration-force:
	$(if $(strip $(version)),,$(error version is required: make migration-force version=1))
	$(MIGRATE) -source $(MIGRATIONS_SOURCE) -database "$(DB_ADDR)" force $(version)

migration-version:
	$(MIGRATE) -source $(MIGRATIONS_SOURCE) -database "$(DB_ADDR)" version

migration-goto:
	$(if $(strip $(version)),,$(error version is required: make migration-goto version=1))
	$(MIGRATE) -source $(MIGRATIONS_SOURCE) -database "$(DB_ADDR)" goto $(version)

migration-drop:
	$(MIGRATE) -source $(MIGRATIONS_SOURCE) -database "$(DB_ADDR)" drop
 