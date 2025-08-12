.PHONY: help build run test clean docker-up docker-down swagger migrate setup dev prod

# Variables
APP_NAME=otp-auth-service
DOCKER_COMPOSE=docker-compose
GO=go
SWAG=swag
GOOSE=goose
GOLANGCI=golangci-lint
AIR=air

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Show this help
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

setup: ## Initial project setup
	@echo "$(GREEN)Setting up project...$(NC)"
	@$(GO) mod download
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@$(GO) install github.com/cosmtrek/air@latest
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "$(GREEN)✔ Setup complete!$(NC)"

install: ## Install dependencies
	@echo "$(GREEN)Installing dependencies...$(NC)"
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "$(GREEN)✔ Dependencies installed!$(NC)"

build: ## Build the application
	@echo "$(GREEN)Building application...$(NC)"
	@$(GO) build -o bin/$(APP_NAME) cmd/api/main.go
	@echo "$(GREEN)✔ Build complete! Binary: bin/$(APP_NAME)$(NC)"

run: ## Run the application locally
	@echo "$(GREEN)Starting application...$(NC)"
	@$(GO) run cmd/api/main.go

dev: ## Run in development mode with hot reload
	@echo "$(GREEN)Starting development server with hot reload...$(NC)"
	@$(AIR)

test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	@$(GO) test -v -cover -race ./...

test-coverage: ## Run tests with coverage report
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✔ Coverage report generated: coverage.html$(NC)"

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	@$(GO) test -v -tags=integration ./tests/integration/...

test-e2e: ## Run end-to-end tests
	@echo "$(GREEN)Running e2e tests...$(NC)"
	@$(GO) test -v -tags=e2e ./tests/e2e/...

benchmark: ## Run benchmarks
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@$(GO) test -bench=. -benchmem ./...

swagger: ## Generate swagger documentation
	@echo "$(GREEN)Generating Swagger documentation...$(NC)"
	@$(SWAG) init -g cmd/api/main.go
	@echo "$(GREEN)✔ Swagger docs generated!$(NC)"

migrate-up: ## Run database migrations
	@echo "$(GREEN)Running migrations...$(NC)"
	@$(GOOSE) -dir migrations postgres "${DB_DSN}" up
	@echo "$(GREEN)✔ Migrations complete!$(NC)"

migrate-down: ## Rollback database migrations
	@echo "$(YELLOW)Rolling back migrations...$(NC)"
	@$(GOOSE) -dir migrations postgres "${DB_DSN}" down
	@echo "$(GREEN)✔ Rollback complete!$(NC)"

migrate-status: ## Check migration status
	@$(GOOSE) -dir migrations postgres "${DB_DSN}" status

docker-build: ## Build docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(GREEN)✔ Docker image built: $(APP_NAME):latest$(NC)"

docker-up: ## Start all services with docker-compose
	@echo "$(GREEN)Starting Docker services...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✔ Services started!$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Swagger: http://localhost:8080/swagger/index.html$(NC)"

docker-down: ## Stop all services
	@echo "$(YELLOW)Stopping Docker services...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✔ Services stopped!$(NC)"

docker-logs: ## Show service logs
	@$(DOCKER_COMPOSE) logs -f

docker-clean: ## Clean docker resources
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	@$(DOCKER_COMPOSE) down -v
	@docker system prune -f
	@echo "$(GREEN)✔ Docker resources cleaned!$(NC)"

lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	@$(GOLANGCI) run ./...
	@echo "$(GREEN)✔ Linting complete!$(NC)"

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GO) fmt ./...
	@goimports -w .
	@echo "$(GREEN)✔ Code formatted!$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf bin/ tmp/ vendor/ coverage.* *.out
	@echo "$(GREEN)✔ Clean complete!$(NC)"

check: lint test ## Run lint and tests
	@echo "$(GREEN)✔ All checks passed!$(NC)"

ci: check build ## Run CI pipeline
	@echo "$(GREEN)✔ CI pipeline complete!$(NC)"

prod: docker-build docker-up ## Deploy to production
	@echo "$(GREEN)✔ Production deployment complete!$(NC)"

.DEFAULT_GOAL := help
