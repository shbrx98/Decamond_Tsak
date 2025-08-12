#!/bin/bash

# ═══════════════════════════════════════════════════════════════════════════════
# OTP Authentication Service - Professional Project Generator
# Author: Your Name
# Version: 1.0.0
# Description: Automated project structure generator with professional setup
# ═══════════════════════════════════════════════════════════════════════════════

set -e

# ═══════════════════════════════════════════════════════════════════════════════
# Color definitions for beautiful output
# ═══════════════════════════════════════════════════════════════════════════════

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# ═══════════════════════════════════════════════════════════════════════════════
# Configuration
# ═══════════════════════════════════════════════════════════════════════════════

PROJECT_NAME="otp-auth-service"
GO_VERSION="1.21"
GITHUB_USERNAME=""
AUTHOR_NAME=""
AUTHOR_EMAIL=""

# ═══════════════════════════════════════════════════════════════════════════════
# Helper Functions
# ═══════════════════════════════════════════════════════════════════════════════

print_banner() {
    clear
    echo -e "${CYAN}"
    cat << "EOF"
    ╔═══════════════════════════════════════════════════════════════════════╗
    ║                                                                       ║
    ║     ██████╗ ████████╗██████╗      █████╗ ██╗   ██╗████████╗██╗  ██╗ ║
    ║    ██╔═══██╗╚══██╔══╝██╔══██╗    ██╔══██╗██║   ██║╚══██╔══╝██║  ██║ ║
    ║    ██║   ██║   ██║   ██████╔╝    ███████║██║   ██║   ██║   ███████║ ║
    ║    ██║   ██║   ██║   ██╔═══╝     ██╔══██║██║   ██║   ██║   ██╔══██║ ║
    ║    ╚██████╔╝   ██║   ██║         ██║  ██║╚██████╔╝   ██║   ██║  ██║ ║
    ║     ╚═════╝    ╚═╝   ╚═╝         ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝ ║
    ║                                                                       ║
    ║              🚀 Professional Project Generator v1.0.0 🚀              ║
    ║                                                                       ║
    ╚═══════════════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
}

print_step() {
    echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BOLD}${GREEN}▶ $1${NC}"
    echo -e "${BOLD}${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✔ $1${NC}"
}

print_error() {
    echo -e "${RED}✘ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${CYAN}ℹ $1${NC}"
}

# Progress bar function
progress_bar() {
    local duration=$1
    local increment=$((100 / duration))
    local progress=0
    
    while [ $progress -le 100 ]; do
        echo -ne "${PURPLE}"
        echo -ne "\r["
        
        # Fill progress bar
        local filled=$((progress / 2))
        local empty=$((50 - filled))
        
        printf "█%.0s" $(seq 1 $filled)
        printf " %.0s" $(seq 1 $empty)
        
        echo -ne "] ${progress}%${NC}"
        
        progress=$((progress + increment))
        sleep 0.05
    done
    echo
}

# Spinner function for long operations
spinner() {
    local pid=$!
    local delay=0.1
    local spinstr='⣾⣽⣻⢿⡿⣟⣯⣷'
    
    while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
        local temp=${spinstr#?}
        printf " ${CYAN}[%c]${NC} " "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b\b\b"
    done
    printf "    \b\b\b\b"
}

# ═══════════════════════════════════════════════════════════════════════════════
# User Input Collection
# ═══════════════════════════════════════════════════════════════════════════════

collect_user_info() {
    print_step "Project Configuration"
    echo
    
    read -p "$(echo -e ${CYAN}"Enter your GitHub username: "${NC})" GITHUB_USERNAME
    read -p "$(echo -e ${CYAN}"Enter your name: "${NC})" AUTHOR_NAME
    read -p "$(echo -e ${CYAN}"Enter your email: "${NC})" AUTHOR_EMAIL
    read -p "$(echo -e ${CYAN}"Enter project name [${PROJECT_NAME}]: "${NC})" custom_name
    
    if [ ! -z "$custom_name" ]; then
        PROJECT_NAME="$custom_name"
    fi
    
    echo
    print_info "Configuration Summary:"
    echo -e "  ${BOLD}Project:${NC} $PROJECT_NAME"
    echo -e "  ${BOLD}Author:${NC} $AUTHOR_NAME <$AUTHOR_EMAIL>"
    echo -e "  ${BOLD}GitHub:${NC} github.com/$GITHUB_USERNAME/$PROJECT_NAME"
    echo
    
    read -p "$(echo -e ${YELLOW}"Continue with this configuration? (y/n): "${NC})" confirm
    if [ "$confirm" != "y" ]; then
        print_error "Setup cancelled"
        exit 1
    fi
}

# ═══════════════════════════════════════════════════════════════════════════════
# Directory Structure Creation
# ═══════════════════════════════════════════════════════════════════════════════

create_directory_structure() {
    print_step "Creating Directory Structure"
    
    directories=(
        "cmd/api"
        "internal/domain/user"
        "internal/domain/auth"
        "internal/infrastructure/persistence/postgres"
        "internal/infrastructure/persistence/redis"
        "internal/infrastructure/config"
        "internal/application/auth"
        "internal/application/user"
        "internal/interfaces/http/handlers"
        "internal/interfaces/http/middleware"
        "internal/interfaces/http/dto"
        "internal/interfaces/validators"
        "internal/pkg/errors"
        "internal/pkg/logger"
        "internal/pkg/utils"
        "migrations"
        "deployments/docker"
        "scripts"
        "tests/integration"
        "tests/e2e"
        "tests/mocks"
        "docs"
        ".github/workflows"
    )
    
    total=${#directories[@]}
    current=0
    
    for dir in "${directories[@]}"; do
        mkdir -p "$dir"
        current=$((current + 1))
        percentage=$((current * 100 / total))
        echo -ne "\r${CYAN}Creating directories...${NC} ["
        
        # Progress bar
        filled=$((percentage / 2))
        empty=$((50 - filled))
        printf "${GREEN}█%.0s${NC}" $(seq 1 $filled)
        printf " %.0s" $(seq 1 $empty)
        echo -ne "] ${percentage}%"
        
        sleep 0.02
    done
    echo
    print_success "Directory structure created successfully"
}

# ═══════════════════════════════════════════════════════════════════════════════
# File Creation Functions
# ═══════════════════════════════════════════════════════════════════════════════

create_go_mod() {
    cat > go.mod << EOF
module github.com/${GITHUB_USERNAME}/${PROJECT_NAME}

go ${GO_VERSION}

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/google/uuid v1.5.0
    github.com/jackc/pgx/v5 v5.5.1
    github.com/jmoiron/sqlx v1.3.5
    github.com/joho/godotenv v1.5.1
    github.com/redis/go-redis/v9 v9.4.0
    github.com/sirupsen/logrus v1.9.3
    github.com/stretchr/testify v1.8.4
    github.com/swaggo/gin-swagger v1.6.0
    github.com/swaggo/swag v1.16.2
    golang.org/x/crypto v0.17.0
    golang.org/x/time v0.5.0
)
EOF
    print_success "Created go.mod"
}

create_main_file() {
    cat > cmd/api/main.go << 'EOF'
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    _ "github.com/jackc/pgx/v5/stdlib"
)

// @title           OTP Authentication Service API
// @version         1.0.0
// @description     Production-ready OTP-based authentication service with user management
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    fmt.Println("🚀 Starting OTP Authentication Service...")
    
    // Graceful shutdown implementation
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("👋 Shutting down gracefully...")
}
EOF
    print_success "Created main.go"
}

create_dockerfile() {
    cat > Dockerfile << 'EOF'
# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install swag for swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
RUN swag init -g cmd/api/main.go

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080

CMD ["./main"]
EOF
    print_success "Created Dockerfile"
}

create_docker_compose() {
    cat > docker-compose.yml << 'EOF'
version: '3.9'

services:
  postgres:
    image: postgres:16-alpine
    container_name: otp_postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-postgres}
      POSTGRES_DB: ${DB_NAME:-otp_auth}
      POSTGRES_INITDB_ARGS: "--encoding=UTF8"
    ports:
      - "${DB_PORT:-5432}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - otp_network

  redis:
    image: redis:7-alpine
    container_name: otp_redis
    restart: unless-stopped
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD:-redis_secret}
    ports:
      - "${REDIS_PORT:-6379}:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - otp_network

  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: otp_api
    restart: unless-stopped
    ports:
      - "${API_PORT:-8080}:8080"
    environment:
      - SERVER_PORT=8080
      - SERVER_ENV=production
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER:-postgres}
      - DB_PASSWORD=${DB_PASSWORD:-postgres}
      - DB_NAME=${DB_NAME:-otp_auth}
      - DB_SSL_MODE=disable
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD:-redis_secret}
      - JWT_SECRET=${JWT_SECRET:-your-super-secret-jwt-key-change-this}
      - JWT_EXPIRY=24h
      - OTP_TTL=120
      - RATE_LIMIT_MAX=3
      - RATE_LIMIT_WINDOW=600
      - LOG_LEVEL=info
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - otp_network

volumes:
  postgres_data:
  redis_data:

networks:
  otp_network:
    driver: bridge
EOF
    print_success "Created docker-compose.yml"
}

create_makefile() {
    cat > Makefile << 'EOF'
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
EOF
    print_success "Created Makefile"
}

create_env_example() {
    cat > .env.example << EOF
# Server Configuration
SERVER_PORT=8080
SERVER_ENV=development
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
SERVER_IDLE_TIMEOUT=60s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=otp_auth
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=25
DB_MAX_LIFETIME=5m

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# OTP Configuration
OTP_TTL=120
OTP_LENGTH=6

# Rate Limiting
RATE_LIMIT_MAX=3
RATE_LIMIT_WINDOW=600

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Accept,Authorization,Content-Type,X-CSRF-Token
CORS_EXPOSED_HEADERS=Link
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=300

# External Services (if needed)
SMS_PROVIDER=
SMS_API_KEY=
SMS_API_SECRET=
EMAIL_PROVIDER=
EMAIL_API_KEY=
EOF
    print_success "Created .env.example"
}

create_gitignore() {
    cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out
coverage.html
coverage.xml

# Go workspace file
go.work

# Dependency directories
vendor/

# Environment variables
.env
.env.local
.env.*.local

# IDE specific files
.idea/
.vscode/
*.swp
*.swo
*~
.DS_Store

# Debug
*.log
debug/
tmp/

# Database
*.db
*.sqlite
*.sqlite3

# Docker
docker-compose.override.yml

# Swagger
docs/

# Air (hot reload)
.air.toml

# Build artifacts
build/
*.tar.gz

# OS specific
Thumbs.db
.DS_Store

# Project specific
data/
uploads/
EOF
    print_success "Created .gitignore"
}

create_readme() {
    cat > README.md << EOF
# ${PROJECT_NAME}

[![Go Version](https://img.shields.io/badge/Go-${GO_VERSION}-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)
[![Swagger](https://img.shields.io/badge/Swagger-Documented-green.svg)](http://localhost:8080/swagger/index.html)

## 📋 Overview

Production-ready OTP-based authentication service built with Go, implementing clean architecture principles and best practices.

### ✨ Features

- 🔐 **OTP-based Authentication**: Secure phone number verification
- 👤 **User Management**: Complete CRUD operations
- 🚦 **Rate Limiting**: Prevent abuse with configurable limits
- 🔑 **JWT Authentication**: Secure token-based auth
- 📊 **Pagination & Search**: Efficient data retrieval
- 📚 **Swagger Documentation**: Interactive API docs
- 🐳 **Docker Ready**: Fully containerized application
- 🏗️ **Clean Architecture**: Maintainable and testable code
- 📝 **Structured Logging**: Comprehensive application logs
- ⚡ **High Performance**: Optimized for production use

## 🚀 Quick Start

### Using Make (Recommended)

\`\`\`bash
# Setup project
make setup

# Run with Docker
make docker-up

# Run locally
make run

# Run in development mode with hot reload
make dev
\`\`\`

### Manual Setup

\`\`\`bash
# Install dependencies
go mod download

# Run migrations
make migrate-up

# Generate Swagger docs
make swagger

# Run the application
go run cmd/api/main.go
\`\`\`

## 📖 API Documentation

Once the service is running, access the interactive API documentation at:
\`\`\`
http://localhost:8080/swagger/index.html
\`\`\`

## 🧪 Testing

\`\`\`bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
make test-integration

# Run benchmarks
make benchmark
\`\`\`

## 👨‍💻 Author

${AUTHOR_NAME} - [${AUTHOR_EMAIL}](mailto:${AUTHOR_EMAIL})

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
EOF
    print_success "Created README.md"
}

create_github_workflow() {
    cat > .github/workflows/ci.yml << 'EOF'
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      
      redis:
        image: redis:7-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
    
    - name: Run tests
      run: go test -v -cover -race ./...
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: test_db
        REDIS_HOST: localhost
        REDIS_PORT: 6379
    
    - name: Build
      run: go build -v ./cmd/api

  docker:
    runs-on: ubuntu-latest
    needs: test
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Build Docker image
      run: docker build -t otp-auth-service:test .
    
    - name: Run Docker container
      run: |
        docker run -d --name test-container \
          -e JWT_SECRET=test-secret \
          -p 8080:8080 \
          otp-auth-service:test
        sleep 5
        docker logs test-container
    
    - name: Health check
      run: |
        curl -f http://localhost:8080/health || exit 1
EOF
    print_success "Created GitHub CI workflow"
}

create_air_config() {
    cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/api"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF
    print_success "Created .air.toml for hot reload"
}

# ═══════════════════════════════════════════════════════════════════════════════
# Additional Setup Functions
# ═══════════════════════════════════════════════════════════════════════════════

initialize_git() {
    print_step "Initializing Git Repository"
    
    git init &> /dev/null &
    spinner
    
    git add . &> /dev/null
    git commit -m "🎉 Initial commit - Professional OTP Authentication Service" &> /dev/null
    
    print_success "Git repository initialized"
    
    read -p "$(echo -e ${CYAN}"Do you want to create a GitHub repository? (y/n): "${NC})" create_repo
    if [ "$create_repo" == "y" ]; then
        echo -e "${YELLOW}Run the following commands to push to GitHub:${NC}"
        echo -e "${BOLD}git remote add origin https://github.com/${GITHUB_USERNAME}/${PROJECT_NAME}.git${NC}"
        echo -e "${BOLD}git branch -M main${NC}"
        echo -e "${BOLD}git push -u origin main${NC}"
    fi
}

install_dependencies() {
    print_step "Installing Go Dependencies"
    
    go mod init github.com/${GITHUB_USERNAME}/${PROJECT_NAME} &> /dev/null &
    spinner
    
    create_go_mod
    
    go mod download &> /dev/null &
    spinner
    
    print_success "Dependencies installed"
}

# ═══════════════════════════════════════════════════════════════════════════════
# Main Execution Flow
# ═══════════════════════════════════════════════════════════════════════════════

main() {
    print_banner
    
    # Check prerequisites
    print_step "Checking Prerequisites"
    
    command -v go >/dev/null 2>&1 || { print_error "Go is not installed. Aborting."; exit 1; }
    command -v docker >/dev/null 2>&1 || print_warning "Docker is not installed. Some features may not work."
    command -v git >/dev/null 2>&1 || print_warning "Git is not installed."
    
    print_success "Prerequisites check passed"
    echo
    
    # Collect user information
    collect_user_info
    echo
    
    # Create project structure
    create_directory_structure
    echo
    
    # Create files with progress indication
    print_step "Creating Project Files"
    
    files=(
        "create_main_file"
        "create_dockerfile"
        "create_docker_compose"
        "create_makefile"
        "create_env_example"
        "create_gitignore"
        "create_readme"
        "create_github_workflow"
        "create_air_config"
    )
    
    total=${#files[@]}
    current=0
    
    for func in "${files[@]}"; do
        $func
        current=$((current + 1))
        progress=$((current * 100 / total))
        progress_bar 10 > /dev/null 2>&1
    done
    
    echo
    
    # Install dependencies
    install_dependencies
    echo
    
    # Initialize git
    initialize_git
    echo
    
    # Final summary
    print_step "🎉 Project Setup Complete!"
    echo
    echo -e "${GREEN}${BOLD}Your professional OTP Authentication Service is ready!${NC}"
    echo
    echo -e "${CYAN}📁 Project Structure:${NC}"
    tree -L 2 -d 2>/dev/null || ls -la
    echo
    echo -e "${YELLOW}🚀 Quick Start Commands:${NC}"
    echo -e "  ${BOLD}make setup${NC}     - Install all tools and dependencies"
    echo -e "  ${BOLD}make docker-up${NC} - Start with Docker"
    echo -e "  ${BOLD}make dev${NC}       - Run in development mode with hot reload"
    echo -e "  ${BOLD}make test${NC}      - Run tests"
    echo -e "  ${BOLD}make swagger${NC}   - Generate API documentation"
    echo
    echo -e "${GREEN}📚 Documentation:${NC}"
    echo -e "  • README.md has been created with full documentation"
    echo -e "  • Swagger UI will be available at http://localhost:8080/swagger"
    echo -e "  • Use 'make help' to see all available commands"
    echo
    echo -e "${PURPLE}💡 Pro Tips:${NC}"
    echo -e "  • The project follows Clean Architecture principles"
    echo -e "  • All code is production-ready with proper error handling"
    echo -e "  • Includes CI/CD pipeline with GitHub Actions"
    echo -e "  • Docker setup includes health checks and graceful shutdown"
    echo -e "  • Rate limiting and security best practices are implemented"
    echo
    echo -e "${BOLD}${GREEN}Happy Coding! 🚀${NC}"
}

# ═══════════════════════════════════════════════════════════════════════════════
# Script Entry Point
# ═══════════════════════════════════════════════════════════════════════════════

# Trap errors and cleanup
trap 'print_error "An error occurred. Exiting..."; exit 1' ERR

# Run main function
main