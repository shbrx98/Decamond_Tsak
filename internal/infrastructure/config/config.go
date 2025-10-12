package config

import (
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type ServerConfig struct {
    Port int
    Env  string
}
type DatabaseConfig struct {
    Host            string
    Port            int
    User            string
    Password        string
    Name            string
    SSLMode         string
    DSN             string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}
type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
    Addr     string
}
type JWTConfig struct {
    Secret   string
    Expiry   time.Duration
    Issuer   string
    Audience string
}
type OTPConfig struct {
    TTL    time.Duration
    Pepper string // used to HMAC OTPs at rest
}
type RateLimitConfig struct {
    Max    int
    Window time.Duration
}

type Config struct {
    Server    ServerConfig
    Database  DatabaseConfig
    Redis     RedisConfig
    JWT       JWTConfig
    OTP       OTPConfig
    RateLimit RateLimitConfig
}

func Load() (*Config, error) {
    _ = godotenv.Load()

    cfg := &Config{
        Server: ServerConfig{
            Port: getInt("SERVER_PORT", 8082),
            Env:  get("SERVER_ENV", "development"),
        },
        Database: DatabaseConfig{
            Host:            get("DB_HOST", "localhost"),
            Port:            getInt("DB_PORT", 5432),
            User:            get("DB_USER", "postgres"),
            Password:        get("DB_PASSWORD", "Hh@12345"),
            Name:            get("DB_NAME", "otp_auth"),
            SSLMode:         get("DB_SSL_MODE", "disable"),
            MaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 20),
            MaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 10),
            ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
        },
        Redis: RedisConfig{
            Host:     get("REDIS_HOST", "localhost"),
            Port:     getInt("REDIS_PORT", 6379),
            Password: get("REDIS_PASSWORD", ""),
            DB:       getInt("REDIS_DB", 0),
        },
        JWT: JWTConfig{
            Secret:   get("JWT_SECRET", ""),
            Expiry:   getDuration("JWT_EXPIRY", 24*time.Hour),
            Issuer:   get("JWT_ISSUER", "otp-auth"),
            Audience: get("JWT_AUDIENCE", "otp-clients"),
        },
        OTP: OTPConfig{
            TTL:    getDuration("OTP_TTL", 2*time.Minute),
            Pepper: get("OTP_STORE_PEPPER", get("JWT_SECRET", "change-me")),
        },
        RateLimit: RateLimitConfig{
            Max:    getInt("RATE_LIMIT_MAX", 3),
            Window: getDuration("RATE_LIMIT_WINDOW", 10*time.Minute),
        },
    }
    cfg.Database.DSN = "postgres://" + cfg.Database.User + ":" + cfg.Database.Password +
        "@" + cfg.Database.Host + ":" + strconv.Itoa(cfg.Database.Port) + "/" + cfg.Database.Name +
        "?sslmode=" + cfg.Database.SSLMode

    cfg.Redis.Addr = cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port)
    return cfg, nil
}

func get(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}
func getInt(key string, def int) int {
    if v := os.Getenv(key); v != "" {
        if i, err := strconv.Atoi(v); err == nil {
            return i
        }
    }
    return def
}
func getDuration(key string, def time.Duration) time.Duration {
    if v := os.Getenv(key); v != "" {
        if d, err := time.ParseDuration(v); err == nil {
            return d
        }
    }
    return def
}