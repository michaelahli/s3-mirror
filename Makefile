.PHONY: build test clean release snapshot help

# Variables
BINARY_NAME=s3-mirror
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/s3-mirror

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

coverage: test ## Show test coverage
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -rf dist/

install: build ## Install binary to $GOPATH/bin
	mv $(BINARY_NAME) $(GOPATH)/bin/

snapshot: ## Create a snapshot release (local testing)
	goreleaser release --snapshot --clean

release: ## Create a release (requires git tag)
	goreleaser release --clean

lint: ## Run linters
	golangci-lint run ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

tidy: ## Tidy go modules
	go mod tidy

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all checks

run: build ## Build and run with example config
	./$(BINARY_NAME) -config config.yaml -dry-run -verbose

.DEFAULT_GOAL := help
