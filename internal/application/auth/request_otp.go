package auth

import (
    "context"
    "crypto/rand"
    "fmt"
    "time"

    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/utils"
)

type RateLimiter interface {
    CheckLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}

type RequestOTPUseCase struct {
    authService AuthService
    rateLimiter RateLimiter
    logger      logger.Logger

    otpTTL time.Duration
}

func NewRequestOTPUseCase(authService AuthService, rateLimiter RateLimiter, log logger.Logger) *RequestOTPUseCase {
    return &RequestOTPUseCase{
        authService: authService,
        rateLimiter: rateLimiter,
        logger:      log,
        otpTTL:      2 * time.Minute, 
    }
}

func (uc *RequestOTPUseCase) Execute(ctx context.Context, phone string) error {
    normalizedPhone, err := utils.NormalizePhoneNumber(phone)
    if err != nil {
        uc.logger.WithField("phone", phone).Error("Invalid phone format")
        return appErr.NewValidationError("invalid phone number format")
    }

    allowed, err := uc.rateLimiter.CheckLimit(
        ctx,
        fmt.Sprintf("otp:%s", normalizedPhone),
        3,
        10*time.Minute,
    )
    if err != nil {
        uc.logger.WithError(err).Error("Rate limiter error")
        return appErr.NewInternalError("rate limiting check failed")
    }
    if !allowed {
        uc.logger.WithField("phone", normalizedPhone).Warn("Rate limit exceeded")
        return appErr.NewRateLimitError("too many OTP requests, please try again later")
    }

    code := generateOTPCode()
    if err := uc.authService.SaveOTP(ctx, normalizedPhone, code, uc.otpTTL); err != nil {
        uc.logger.WithError(err).Error("Failed to store OTP")
        return appErr.NewInternalError("failed to generate OTP")
    }

    
    uc.logger.WithFields(map[string]interface{}{
        "phone":      normalizedPhone,
        "otp":        code,
        "expires_in": uc.otpTTL.String(),
    }).Info("OTP Generated")

    return nil
}

func generateOTPCode() string {
    b := make([]byte, 3)
    if _, err := rand.Read(b); err != nil {
        return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
    }
    return fmt.Sprintf("%06d", (int(b[0])<<16|int(b[1])<<8|int(b[2]))%1000000)
}