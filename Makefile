
install-tools: ## Install required development tools
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully!"

lint: client-lint server-lint ## Run linters on all modules
	@echo "✓ All modules linted successfully"

client-lint: ## Lint client module
	@echo "Linting client module..."
	@cd client && golangci-lint run --config=../.golangci.yml ./...

server-lint: ## Lint server module
	@echo "Linting server module..."
	@cd services/shortcode && golangci-lint run --config=../../.golangci.yml ./...

fmt: ## Format all Go code
	@echo "Formatting client code..."
	@cd client && gofmt -s -w .
	@echo "Formatting server code..."
	@cd services/shortcode && gofmt -s -w .
	@echo "✓ All code formatted"

fmt-check: ## Check if all code is formatted
	@echo "Checking client code formatting..."
	@cd client && if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Client files not formatted:"; \
		gofmt -s -l .; \
		exit 1; \
	fi
	@echo "Checking server code formatting..."
	@cd services/shortcode && if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Server files not formatted:"; \
		gofmt -s -l .; \
		exit 1; \
	fi
	@echo "✓ All code is properly formatted"

vet: ## Run go vet on all modules
	@echo "Running go vet on client..."
	@cd client && go vet ./...
	@echo "Running go vet on server..."
	@cd services/shortcode && go vet ./...
	@echo "✓ Go vet passed"

test: client-test server-test ## Run all tests
	@echo "✓ All tests passed"

client-test: ## Test client module
	@echo "Testing client module..."
	@cd client && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

server-test: ## Test server module
	@echo "Testing server module..."
	@cd services/shortcode && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

all: fmt-check vet lint test ## Run all quality checks (simulates CI)
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "✓ All quality checks passed!"
	@echo "  - Code formatting: OK"
	@echo "  - Go vet: OK"
	@echo "  - Linters: OK"
	@echo "  - Tests: OK"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Docker targets
.PHONY: docker-up docker-down docker-logs docker-build docker-push docker-dev-up docker-dev-down

docker-up: ## Start all services with docker-compose (using DockerHub images)
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## View logs from all services
	docker-compose logs -f

docker-build: ## Build Docker images locally
	@./scripts/build-images.sh $(DOCKER_USERNAME)

docker-push: ## Build and push Docker images to DockerHub
	@./scripts/build-and-push.sh $(DOCKER_USERNAME)

docker-dev-up: ## Start all services with docker-compose (building locally)
	docker-compose -f docker-compose.dev.yml up -d

docker-dev-down: ## Stop all development services
	docker-compose -f docker-compose.dev.yml down
