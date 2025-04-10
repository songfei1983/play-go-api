.PHONY: help test build run clean integration-test up down

# Make help the default target
.DEFAULT_GOAL := help

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput sgr0)

TARGET_MAX_CHAR_NUM=20

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
	go build -o bin/app cmd/main.go

## Run the application
run: build-app
	./bin/app

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
