.PHONY: help build run test clean deps tidy fmt lint docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  deps         - Download dependencies"
	@echo "  tidy         - Tidy go.mod"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"

# Download dependencies
deps:
	go mod download

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

# Code check
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config=../../.golangci.yml ./...; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run --config=../../.golangci.yml ./...; \
	fi

# Build application
build:
	go build -o bin/shortcode cmd/server/main.go

# Run application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Additional quality check targets
.PHONY: fmt-check vet test-unit check-quality

fmt-check: ## Check if code is formatted
	@echo "Checking code formatting..."
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "The following files are not formatted:"; \
		gofmt -s -l .; \
		exit 1; \
	fi
	@echo "All files are properly formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

test-unit: ## Run unit tests with coverage
	@echo "Running unit tests..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@if [ -f coverage.out ]; then \
		echo ""; \
		echo "Coverage report:"; \
		go tool cover -func=coverage.out | tail -n 1; \
	fi

check-quality: fmt-check vet lint test-unit ## Run all quality checks
	@echo "All quality checks passed!"

# Build Docker image
docker-build:
	docker build -t shortcode-service:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env shortcode-service:latest
