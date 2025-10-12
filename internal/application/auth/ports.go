package auth

import (
    "context"
    "time"

    domUser "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
)

type AuthService interface {
    SaveOTP(ctx context.Context, phone, code string, ttl time.Duration) error
    VerifyOTP(ctx context.Context, phone, code string) error
    GenerateJWT(u *domUser.User, expiry time.Duration, issuer, audience string) (string, int64, error)
}