.PHONY: build run clean docker-build docker-up docker-down

# Go related variables
GO=go
BINARY_NAME=play-go-api

# Docker related variables
DOCKER_COMPOSE=docker-compose

# Format Go code
fmt:
	$(GO) fmt ./...

# Run linter
lint:
	golangci-lint run ./...

# Run tests
test:
	$(GO) test -v ./...

# Default target
all: build

# Build binary
build: fmt
	$(GO) build -o $(BINARY_NAME) .

# Run the application locally
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)

# Build Docker image
docker-build:
	$(DOCKER_COMPOSE) build

# Start all Docker containers
docker-up:
	$(DOCKER_COMPOSE) up -d

# Stop and remove all Docker containers and volumes
docker-down:
	$(DOCKER_COMPOSE) down -v

# Start development environment
dev: docker-up
	$(GO) run main.go

# Show help
help:
	@echo "Available targets:"
	@echo "  fmt         - Format Go code"
	@echo "  lint        - Run linter"
	@echo "  test        - Run tests"
	@echo "  build       - Build binary"
	@echo "  run         - Run application locally"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-up   - Start Docker containers"
	@echo "  docker-down - Stop Docker containers"
	@echo "  dev         - Start development environment"
	@echo "  help        - Show this help message"
