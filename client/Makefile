.PHONY: build run test clean install help

# Variable definitions
BINARY_NAME=shortcode-client
BUILD_DIR=bin
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

# Default target URL
URL ?= http://localhost

help: ## Show help information
	@echo "Short link service client tool"
	@echo ""
	@echo "Available commands:"
	@echo "  make build          - Build client"
	@echo "  make install        - Install dependencies"
	@echo "  make test           - Run complete test suite"
	@echo "  make test-verbose   - Run tests (verbose output)"
	@echo "  make create         - Create short link example"
	@echo "  make clean          - Clean build files"
	@echo "Usage examples:"
	@echo "  make test URL=http://localhost"
	@echo "  make create URL=https://github.com/lincyaw CODE=mycode"

install: ## Install Go dependencies
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

build: ## Build client
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

test: build ## Run test suite
	@echo "Running test suite..."
	./$(BUILD_DIR)/$(BINARY_NAME) test -u $(URL)

test-verbose: build ## Run test suite (verbose output)
	@echo "Running test suite (verbose mode)..."
	./$(BUILD_DIR)/$(BINARY_NAME) test -u $(URL) -v

create: build ## Create short link example
	@echo "Creating short link example..."
	./$(BUILD_DIR)/$(BINARY_NAME) create -u $(URL) -l "https://github.com/lincyaw/tools" $(if $(CODE),-c $(CODE))

stats: build ## Get statistics
	@if [ -z "$(CODE)" ]; then \
		echo "Error: Please provide CODE parameter"; \
		echo "Example: make stats CODE=test123"; \
		exit 1; \
	fi
	./$(BUILD_DIR)/$(BINARY_NAME) stats $(CODE) -u $(URL)

delete: build ## Delete short link
	@if [ -z "$(CODE)" ]; then \
		echo "Error: Please provide CODE parameter"; \
		echo "Example: make delete CODE=test123"; \
		exit 1; \
	fi
	./$(BUILD_DIR)/$(BINARY_NAME) delete $(CODE) -u $(URL)

redirect: build ## Test redirect
	@if [ -z "$(CODE)" ]; then \
		echo "Error: Please provide CODE parameter"; \
		echo "Example: make redirect CODE=test123"; \
		exit 1; \
	fi
	./$(BUILD_DIR)/$(BINARY_NAME) redirect $(CODE) -u $(URL)

clean: ## Clean build files
	@echo "Cleaning build files..."
	rm -rf $(BUILD_DIR)
	@echo "Cleaning completed"

run-help: build ## Show client help
	./$(BUILD_DIR)/$(BINARY_NAME) --help

# Code quality targets
.PHONY: lint fmt vet test-unit check-quality

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config=../.golangci.yml ./...; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run --config=../.golangci.yml ./...; \
	fi

fmt: ## Format code with gofmt
	@echo "Formatting code..."
	@gofmt -s -w .
	@echo "Code formatted"

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
	$(GO) vet ./...

test-unit: ## Run unit tests with coverage
	@echo "Running unit tests..."
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@if [ -f coverage.out ]; then \
		echo ""; \
		echo "Coverage report:"; \
		$(GO) tool cover -func=coverage.out | tail -n 1; \
	fi

check-quality: fmt-check vet lint test-unit ## Run all quality checks
	@echo "All quality checks passed!"
