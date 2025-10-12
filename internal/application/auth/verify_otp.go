package auth

import (
    "context"
    "strings"
    "time"

    domUser "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/utils"
)

type VerifyOTPUseCase struct {
    authService AuthService
    userService *domUser.Service
    logger      logger.Logger

    jwtExpiry   time.Duration
    jwtIssuer   string
    jwtAudience string
}

type VerifyResult struct {
    Token     string
    ExpiresIn int64
    User      *domUser.User
}

func NewVerifyOTPUseCase(
    authService AuthService,
    userService *domUser.Service,
    log logger.Logger,
    jwtExpiry time.Duration,
    issuer string,
    audience string,
) *VerifyOTPUseCase {
    return &VerifyOTPUseCase{
        authService: authService,
        userService: userService,
        logger:      log,
        jwtExpiry:   jwtExpiry,
        jwtIssuer:   issuer,
        jwtAudience: audience,
    }
}

func (uc *VerifyOTPUseCase) Execute(ctx context.Context, phone, code string) (*VerifyResult, error) {
    
    cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    
    normalizedPhone, err := utils.NormalizePhoneNumber(strings.TrimSpace(phone))
    if err != nil {
        uc.logger.WithField("phone", phone).Warn("invalid phone format")
        return nil, appErr.NewValidationError("invalid phone number format")
    }

    
    code = strings.TrimSpace(code)
    if len(code) != 6 || !allDigits(code) {
        return nil, appErr.NewValidationError("otp must be 6 digits")
    }

    
    if err := uc.authService.VerifyOTP(cctx, normalizedPhone, code); err != nil {
        
        uc.logger.WithFields(map[string]interface{}{
            "phone": normalizedPhone,
        }).WithError(err).Warn("otp verification failed")
        return nil, err
    }

    
    u, err := uc.userService.GetOrCreateByPhone(cctx, normalizedPhone)
    if err != nil {
        uc.logger.WithError(err).Error("user persistence failed")
        return nil, appErr.Wrap(err, "user persistence")
    }

    
    token, exp, err := uc.authService.GenerateJWT(u, uc.jwtExpiry, uc.jwtIssuer, uc.jwtAudience)
    if err != nil {
        uc.logger.WithError(err).Error("token generation failed")
        return nil, appErr.NewInternalError("token generation failed")
    }

    return &VerifyResult{
        Token:     token,
        ExpiresIn: exp,
        User:      u,
    }, nil
}

func allDigits(s string) bool {
    for _, r := range s {
        if r < '0' || r > '9' {
            return false
        }
    }
    return true
}