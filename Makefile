.PHONY: help build test run clean docker-build docker-run helm-lint helm-install fmt vet deps

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Go related targets
build: ## Build the Go application
	go build -o bin/snapshots-api ./cmd/server

test: ## Run all tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Generate test coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

run: ## Run the application locally
	go run ./cmd/server

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Code quality targets
fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

deps: ## Download and tidy dependencies
	go mod download
	go mod tidy

# Docker targets
docker-build: ## Build Docker image
	docker build -t snapshots-api:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --rm snapshots-api:latest

# Helm targets
helm-lint: ## Lint Helm chart
	helm lint charts/snapshots-api

helm-template: ## Render Helm templates
	helm template test-release charts/snapshots-api

helm-install: ## Install Helm chart (requires kubectl context)
	helm install snapshots-api charts/snapshots-api

helm-upgrade: ## Upgrade Helm chart (requires kubectl context)
	helm upgrade snapshots-api charts/snapshots-api

helm-uninstall: ## Uninstall Helm chart (requires kubectl context)
	helm uninstall snapshots-api

# Development targets
dev-setup: deps ## Setup development environment
	@echo "Development environment setup complete"

dev-test: fmt vet test ## Run development tests (format, vet, test)

ci-test: deps test helm-lint ## Run CI-like tests locally

# API testing targets
test-api: ## Test API endpoints (requires running server)
	@echo "Testing health endpoint..."
	curl -f http://localhost:8080/health || echo "Health check failed"
	@echo "\nTesting readiness endpoint..."
	curl -f http://localhost:8080/ready || echo "Readiness check failed"
	@echo "\nTesting mainnet snapshots..."
	curl -f "http://localhost:8080/?network=mainnet" || echo "Mainnet API failed"
	@echo "\nTesting testnet snapshots..."
	curl -f "http://localhost:8080/?network=testnet" || echo "Testnet API failed"
	@echo "\nTesting devnet snapshots..."
	curl -f "http://localhost:8080/?network=devnet" || echo "Devnet API failed" 