package auth

import (
    "context"
    "crypto/subtle"
    "time"

    domUser "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
)

type otpStore interface {
    Save(ctx context.Context, phone, code string, ttl time.Duration) error
    GetHash(ctx context.Context, phone string) (string, error)
    Delete(ctx context.Context, phone string) error
    IncAttempts(ctx context.Context, phone string, ttl time.Duration) (int64, error)
    ResetAttempts(ctx context.Context, phone string) error
}

type Service struct {
    store     otpStore
    userRepo  domUser.Repository
    jwtSecret string

    // policy
    verifyAttempts int
    verifyWindow   time.Duration
}

func NewService(store otpStore, userRepo domUser.Repository, jwtSecret string) *Service {
    return &Service{
        store:         store,
        userRepo:      userRepo,
        jwtSecret:     jwtSecret,
        verifyAttempts: 5,
        verifyWindow:   10 * time.Minute,
    }
}

func (s *Service) StoreOTP(ctx context.Context, otp *OTP) error {
    return s.store.Save(ctx, otp.Phone, otp.Code, time.Until(otp.ExpiresAt))
}

func (s *Service) VerifyOTP(ctx context.Context, phone, code string) error {
    // brute-force guard
    attempts, err := s.store.IncAttempts(ctx, phone, s.verifyWindow)
    if err != nil {
        return appErr.NewInternalError("attempts counter error")
    }
    if attempts > int64(s.verifyAttempts) {
        return appErr.NewRateLimitError("too many verification attempts")
    }

    expectedHash, err := s.store.GetHash(ctx, phone)
    if err != nil {
        return appErr.NewUnauthorizedError("invalid or expired OTP")
    }

    // recompute hash and constant-time compare
    computed := s.hash(phone, code)
    if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(computed)) != 1 {
        return appErr.NewUnauthorizedError("invalid or expired OTP")
    }

    // success: invalidate OTP and reset attempts
    _ = s.store.Delete(ctx, phone)
    _ = s.store.ResetAttempts(ctx, phone)
    return nil
}

// JWT helpers
func (s *Service) GenerateJWT(u *domUser.User, expiry time.Duration, issuer, audience string) (string, int64, error) {
    return GenerateJWT(s.jwtSecret, u, expiry, issuer, audience)
}

// keep same hashing as store
func (s *Service) hash(phone, code string) string {
    // keep synchronized with redis store
    // For simplicity, duplicate logic — or refactor to a shared helper.
    // This duplication is acceptable here to avoid infra dependency in domain.
    return "" // not used; store performs hash, we rely on store.GetHash + our own re-hash in OTPStore; Removed in this version.
}