# PHONY tells Make that the target is not a file name. It is used to avoid conflicts with files that may have the same name as the target.
# Each target is considered a file. If the target is not a file, it should be marked as .PHONY.

include .env
export

GOOSE_DIR := migrations

#Go commands
.PHONY: run
run:
	go run ./cmd/api

.PHONY: build
build:
	go build -o bin/api ./cmd/api

.PHONY: lint
lint:
	golangci-lint run

#* DOCKER COMMANDS --------------------------------
.PHONY: docker-up docker-down
docker-up:
	docker compose up -d
docker-down:
	docker compose down

.PHONY: docker-start docker-stop docker-logs
docker-start:
	docker compose start
docker-stop:
	docker compose stop
docker-logs:
	docker compose logs -f postgres

# MIGRATION COMMANDS --------------------------------
.PHONY: migrate-create
migrate-create:
	goose -s -dir $(GOOSE_DIR) create $(name) sql

.PHONY: migrate-up
migrate-up:
# 	goose -dir-env $(ENV_FILE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" up
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" up

.PHONY: migrate-down
migrate-down:
# 	goose -env $(ENV_FILE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" down
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" down

.PHONY: migrate-status
migrate-status:
# 	goose -dir-env $(ENV_FILE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" status
	goose -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" status