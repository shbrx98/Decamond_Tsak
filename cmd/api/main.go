package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v5/stdlib" // register "pgx" driver

    appAuth "github.com/shbrx98/Decamond_Tsak/internal/application/auth"
    appUser "github.com/shbrx98/Decamond_Tsak/internal/application/user"
    authDomain "github.com/shbrx98/Decamond_Tsak/internal/domain/auth"
    userDomain "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
    "github.com/shbrx98/Decamond_Tsak/internal/infrastructure/config"
    pgRepo "github.com/shbrx98/Decamond_Tsak/internal/infrastructure/persistence/postgres"
    redisInfra "github.com/shbrx98/Decamond_Tsak/internal/infrastructure/persistence/redis"
    httpRouter "github.com/shbrx98/Decamond_Tsak/internal/interfaces/http/router"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"

    _ "github.com/shbrx98/Decamond_Tsak/docs"
)

func main() {
    log := logger.New()
    log.Info("Starting OTP Authentication Service...")

    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    db, err := connectWithRetry(cfg, log, 10)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Pool tuning
    db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

    if err := db.Ping(); err != nil {
        log.Fatalf("DB ping error: %v", err)
    }

    if err := runMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Redis
    redisClient := redisInfra.NewClient(cfg.Redis)
    defer redisClient.Close()
    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        log.Fatalf("Redis ping error: %v", err)
    }

    // Repos and infra
    userRepo := pgRepo.NewUserRepository(db)
    otpStore := redisInfra.NewOTPStore(redisClient, cfg.OTP.Pepper)
    rateLimiter := redisInfra.NewRateLimiter(redisClient)

    // Domain services
    userService := userDomain.NewService(userRepo)
    authService := authDomain.NewService(otpStore, userRepo, cfg.JWT.Secret)

    // Use-cases
    requestOTPUseCase := appAuth.NewRequestOTPUseCase(authService, rateLimiter, log)
    verifyOTPUseCase := appAuth.NewVerifyOTPUseCase(authService, userService, log, cfg.JWT.Expiry, cfg.JWT.Issuer, cfg.JWT.Audience)
    getUserUseCase := appUser.NewGetUserUseCase(userService, log)
    listUsersUseCase := appUser.NewListUsersUseCase(userService, log)

    // Router
    httpRouter := httpRouter.NewRouter(
        requestOTPUseCase,
        verifyOTPUseCase,
        getUserUseCase,
        listUsersUseCase,
        cfg,
        log,
    )

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
        Handler:      httpRouter,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    log.Infof("Server started on port %d", cfg.Server.Port)
    log.Info("Swagger: http://localhost:8080/swagger/index.html")

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Info("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    log.Info("Server exited")
}

func connectWithRetry(cfg *config.Config, log logger.Logger, maxRetries int) (*sqlx.DB, error) {
    var db *sqlx.DB
    var err error
    for i := 0; i < maxRetries; i++ {
        db, err = sqlx.Connect("pgx", cfg.Database.DSN) // use pgx driver
        if err == nil {
            return db, nil
        }
        log.Warnf("DB connect failed (attempt %d/%d): %v", i+1, maxRetries, err)
        time.Sleep(2 * time.Second)
    }
    return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}