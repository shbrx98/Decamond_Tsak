package auth

import (
    "context"

    domAuth "github.com/shbrx98/Decamond_Tsak/internal/domain/auth"
    domUser "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/utils"
    "time"
)

type VerifyOTPUseCase struct {
    authService *domAuth.Service
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
    authService *domAuth.Service,
    userService *domUser.Service,
    logger logger.Logger,
    jwtExpiry time.Duration,
    issuer string,
    audience string,
) *VerifyOTPUseCase {
    return &VerifyOTPUseCase{
        authService: authService,
        userService: userService,
        logger:      logger,
        jwtExpiry:   jwtExpiry,
        jwtIssuer:   issuer,
        jwtAudience: audience,
    }
}

func (uc *VerifyOTPUseCase) Execute(ctx context.Context, phone, code string) (*VerifyResult, error) {
    normalizedPhone, err := utils.NormalizePhoneNumber(phone)
    if err != nil {
        return nil, appErr.NewValidationError("invalid phone number format")
    }

    if err := uc.authService.VerifyOTP(ctx, normalizedPhone, code); err != nil {
        return nil, err
    }

    u, err := uc.userService.GetOrCreateByPhone(ctx, normalizedPhone)
    if err != nil {
        return nil, appErr.Wrap(err, "user persistence")
    }

    token, exp, err := uc.authService.GenerateJWT(u, uc.jwtExpiry, uc.jwtIssuer, uc.jwtAudience)
    if err != nil {
        return nil, appErr.NewInternalError("token generation failed")
    }

    return &VerifyResult{
        Token:     token,
        ExpiresIn: exp,
        User:      u,
    }, nil
}