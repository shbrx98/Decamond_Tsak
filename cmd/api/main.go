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
    approuter "github.com/shbrx98/Decamond_Tsak/internal/interfaces/http/router"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"

    _ "github.com/shbrx98/Decamond_Tsak/docs" // swagger docs
)
// @title           Decamond OTP Authentication Service API
// @version         1.0.0
// @description     OTP-based authentication service (Clean Architecture + DDD)
// @schemes         http
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
    log := logger.New()
    log.Info("Starting OTP Authentication Service...")

    // Load config
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // DB connect with retry
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

    // Migrations (for dev; use golang-migrate in prod)
    if err := runMigrations(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Redis
    redisClient := redisInfra.NewClient(cfg.Redis)
    defer redisClient.Close()
    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        log.Fatalf("Redis ping error: %v", err)
    }

    // Infra & repositories
    userRepo := pgRepo.NewUserRepository(db)
    otpStore := redisInfra.NewOTPStore(redisClient, cfg.OTP.Pepper)
    rateLimiter := redisInfra.NewRateLimiter(redisClient)

    // Domain services
    userService := userDomain.NewService(userRepo)
    
    authDomainSvc := authDomain.NewService(otpStore, userRepo, cfg.JWT.Secret, cfg.OTP.Pepper)

    // Adapter: domain auth -> application port
    authAppPort := authSvcAdapter{svc: authDomainSvc}

    // Use-cases
    requestOTPUseCase := appAuth.NewRequestOTPUseCase(authAppPort, rateLimiter, log)
    verifyOTPUseCase := appAuth.NewVerifyOTPUseCase(authAppPort, userService, log, cfg.JWT.Expiry, cfg.JWT.Issuer, cfg.JWT.Audience)
    getUserUseCase := appUser.NewGetUserUseCase(userService, log)
    listUsersUseCase := appUser.NewListUsersUseCase(userService, log)

    // Router
    routerEngine := approuter.NewRouter(
        requestOTPUseCase,
        verifyOTPUseCase,
        getUserUseCase,
        listUsersUseCase,
        cfg,
        log,
    )

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
        Handler:      routerEngine,
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
    log.Infof("Swagger: http://localhost:%d/swagger/index.html", cfg.Server.Port)
    log.Infof("Health:  http://localhost:%d/health", cfg.Server.Port)


    // Graceful shutdown
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
        db, err = sqlx.Connect("pgx", cfg.Database.DSN)
        if err == nil {
            return db, nil
        }
        log.Warnf("DB connect failed (attempt %d/%d): %v", i+1, maxRetries, err)
        time.Sleep(2 * time.Second)
    }
    return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}

func runMigrations(db *sqlx.DB) error {
    migrations := []string{
        `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
        `CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            phone VARCHAR(20) NOT NULL UNIQUE,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        );`,
        `CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);`,
        `CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);`,
    }
    for _, m := range migrations {
        if _, err := db.Exec(m); err != nil {
            return fmt.Errorf("migration failed: %w", err)
        }
    }
    return nil
}


type authSvcAdapter struct {
    svc *authDomain.Service
}

func (a authSvcAdapter) SaveOTP(ctx context.Context, phone, code string, ttl time.Duration) error {
    otp := &authDomain.OTP{
        Code:      code,
        Phone:     phone,
        ExpiresAt: time.Now().Add(ttl),
    }
    return a.svc.StoreOTP(ctx, otp)
}

func (a authSvcAdapter) VerifyOTP(ctx context.Context, phone, code string) error {
    return a.svc.VerifyOTP(ctx, phone, code)
}

func (a authSvcAdapter) GenerateJWT(u *userDomain.User, expiry time.Duration, issuer, audience string) (string, int64, error) {
    return a.svc.GenerateJWT(u, expiry, issuer, audience)
}