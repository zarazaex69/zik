.PHONY: help build run test test-unit test-integration cover lint clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the server binary
	@echo "Building zik-server..."
	@go build -o bin/zik-server ./cmd/zik-server

run: ## Run the server locally
	@echo "Starting zik-server..."
	@go run ./cmd/zik-server/main.go

test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race ./...

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -race -short ./...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v -race -run Integration ./test/integration/...

cover: ## Generate test coverage report
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ coverage.out coverage.html
	@go clean

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
