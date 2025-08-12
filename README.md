# Decamond_otp-auth-service

[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://go.dev/)
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

```bash
# Setup project
make setup

# Run with Docker
make docker-up

# Run locally
make run

# Run in development mode with hot reload
make dev
```

### Manual Setup

```bash
# Install dependencies
go mod download

# Run migrations
make migrate-up

# Generate Swagger docs
make swagger

# Run the application
go run cmd/api/main.go
```

## 📖 API Documentation

Once the service is running, access the interactive API documentation at:

```text
http://localhost:8080/swagger/index.html
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
make test-integration

# Run benchmarks
make benchmark
```

## 👨‍💻 Author

hossein razavi - [](mailto:)

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
