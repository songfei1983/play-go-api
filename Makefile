.PHONY: help test build run clean integration-test up down run-local set-env

# Make help the default target
.DEFAULT_GOAL := help

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput sgr0)

TARGET_MAX_CHAR_NUM=20

# Environment variables
export DB_HOST ?= localhost
export DB_PORT ?= 3306
export DB_USER ?= app_user
export DB_PASSWORD ?= app_password
export DB_NAME ?= app_db
export REDIS_HOST ?= localhost
export REDIS_PORT ?= 6379
export JWT_SECRET ?= your-secret-key
export PORT ?= 8080

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Build the application
build-app:
	go build -o bin/app cmd/api/main.go

## Run the application
run: build-app
	./bin/app

## Run the application locally with environment variables
run-local: set-env
	go run cmd/api/main.go

## Set environment variables for local development
set-env:
	@echo "Setting up environment variables..."
	@echo "DB_HOST=$(DB_HOST)"
	@echo "DB_PORT=$(DB_PORT)"
	@echo "DB_USER=$(DB_USER)"
	@echo "DB_NAME=$(DB_NAME)"
	@echo "REDIS_HOST=$(REDIS_HOST)"
	@echo "REDIS_PORT=$(REDIS_PORT)"
	@echo "PORT=$(PORT)"

## Run all tests
test:
	go test -v ./...

## Run integration tests
integration-test:
	./test/integration_test.sh

## Clean build artifacts
clean:
	rm -rf bin/
	go clean

## Build docker image
build:
	docker-compose build

## Start Docker containers
up: build
	docker-compose up -d

## Stop Docker containers
down:
	docker-compose down
