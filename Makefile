# PHONY tells Make that the target is not a file name. It is used to avoid conflicts with files that may have the same name as the target.
# Each target is considered a file. If the target is not a file, it should be marked as .PHONY.

include .env
export

GOOSE_DIR := migrations

.PHONY: run build lint \
        docker-up docker-down docker-start docker-stop docker-logs \
        migrate-create migrate-up migrate-down migrate-status

#! Go commands
run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

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

migrate-up:
# 	goose -dir-env $(ENV_FILE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" up
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" up

migrate-down:
# 	goose -env $(ENV_FILE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" down
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" down

migrate-status:
# 	goose -dir-env $(ENV_FILE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" status
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" status