VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY := skillsmith

.PHONY: build test lint fmt clean mod-tidy coverage help install

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/skillsmith

install: build ## Install to GOPATH/bin
	go install -ldflags "-X main.version=$(VERSION)" ./cmd/skillsmith

test: ## Run tests
	go test -race ./...

coverage: ## Run tests with coverage
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run --timeout=5m

fmt: ## Format code
	golangci-lint fmt

clean: ## Clean build artifacts
	rm -f coverage.out coverage.html $(BINARY)

mod-tidy: ## Tidy Go modules
	go mod tidy
