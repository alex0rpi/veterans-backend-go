# PHONY tells Make that the target is not a file name. It is used to avoid conflicts with files that may have the same name as the target.
# Each target is considered a file. If the target is not a file, it should be marked as .PHONY.

include .env
export

GOOSE_DIR := migrations

.PHONY: run build lint \
        docker-up docker-down docker-start docker-stop docker-logs \
        migrate-create prod-migrate-up prod-migrate-down prod-migrate-status \
        local-migrate-up local-migrate-down local-migrate-status

#! Go commands
run: 
	go run ./cmd/api
build:
	go build -o bin/api ./cmd/api
format:
# This will format all Go files in the current directory and its subdirectories.
	go fmt ./...
lint:
	golangci-lint run

#* DOCKER COMMANDS --------------------------------
docker-up:
	docker compose up -d
docker-down:
	docker compose down
docker-start:
	docker compose start
docker-stop:
	docker compose stop
docker-logs:
	docker compose logs -f postgres

#> MIGRATION COMMANDS --------------------------------
migrate-create:
	goose -s -dir $(GOOSE_DIR) create $(name) sql

# Production (Neon)
prod-migrate-up:
# 	goose -dir-env $(ENV_FILE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" up
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" up

prod-migrate-down:
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" down

prod-migrate-status:
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" status

# Local
local-migrate-up:
	goose -dir $(GOOSE_DIR) postgres "$$LOCAL_DATABASE_URL" up

local-migrate-down:
	goose -dir $(GOOSE_DIR) postgres "$$LOCAL_DATABASE_URL" down

local-migrate-status:
	goose -dir $(GOOSE_DIR) postgres "$$LOCAL_DATABASE_URL" status