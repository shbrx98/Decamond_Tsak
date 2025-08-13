package auth

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "crypto/subtle"
    "encoding/hex"
    "regexp"
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
    store          otpStore
    userRepo       domUser.Repository
    jwtSecret      string
    otpPepper      []byte
    verifyAttempts int
    verifyWindow   time.Duration
}

func NewService(store otpStore, userRepo domUser.Repository, jwtSecret string, otpPepper string) *Service {
    return &Service{
        store:          store,
        userRepo:       userRepo,
        jwtSecret:      jwtSecret,
        otpPepper:      []byte(otpPepper),
        verifyAttempts: 5,
        verifyWindow:   10 * time.Minute,
    }
}

func (s *Service) StoreOTP(ctx context.Context, otp *OTP) error {
    
    return s.store.Save(ctx, otp.Phone, otp.Code, time.Until(otp.ExpiresAt))
}

func (s *Service) VerifyOTP(ctx context.Context, phone, code string) error {
    
    attempts, err := s.store.IncAttempts(ctx, phone, s.verifyWindow)
    if err != nil {
        return appErr.NewInternalError("attempts counter error")
    }
    if attempts > int64(s.verifyAttempts) {
        return appErr.NewRateLimitError("too many verification attempts")
    }

    expected, err := s.store.GetHash(ctx, phone)
    if err != nil || expected == "" {
        return appErr.NewUnauthorizedError("invalid or expired OTP")
    }

    
    if looksSHA256Hex(expected) && len(s.otpPepper) > 0 {
        computed := s.hash(phone, code)
        if subtle.ConstantTimeCompare([]byte(expected), []byte(computed)) != 1 {
            return appErr.NewUnauthorizedError("invalid or expired OTP")
        }
    } else {
         
        if subtle.ConstantTimeCompare([]byte(expected), []byte(code)) != 1 {
            return appErr.NewUnauthorizedError("invalid or expired OTP")
        }
    }

    _ = s.store.Delete(ctx, phone)
    _ = s.store.ResetAttempts(ctx, phone)
    return nil
}

func (s *Service) GenerateJWT(u *domUser.User, expiry time.Duration, issuer, audience string) (string, int64, error) {
    return GenerateJWT(s.jwtSecret, u, expiry, issuer, audience)
}

func (s *Service) hash(phone, code string) string {
    mac := hmac.New(sha256.New, s.otpPepper)
    mac.Write([]byte(phone))
    mac.Write([]byte("|"))
    mac.Write([]byte(code))
    return hex.EncodeToString(mac.Sum(nil))
}

var shaHexRe = regexp.MustCompile(`^[0-9a-f]{64}$`)
func looksSHA256Hex(s string) bool { return shaHexRe.MatchString(s) }