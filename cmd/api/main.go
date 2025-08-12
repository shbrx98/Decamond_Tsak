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
