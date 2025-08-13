package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    appAuth "github.com/shbrx98/Decamond_Tsak/internal/application/auth"
    "github.com/shbrx98/Decamond_Tsak/internal/interfaces/http/dto"
    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
)

type AuthHandler struct {
    requestOTP *appAuth.RequestOTPUseCase
    verifyOTP  *appAuth.VerifyOTPUseCase
    logger     logger.Logger
}

func NewAuthHandler(requestOTP *appAuth.RequestOTPUseCase, verifyOTP *appAuth.VerifyOTPUseCase, l logger.Logger) *AuthHandler {
    return &AuthHandler{requestOTP: requestOTP, verifyOTP: verifyOTP, logger: l}
}

// RequestOTP godoc
// @Summary      Request OTP
// @Description  Generate OTP for a phone number
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RequestOTPRequest  true  "Phone number"
// @Success      200      {object}  dto.MessageResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      429      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /auth/otp/request [post]
func (h *AuthHandler) RequestOTP(c *gin.Context) {
    var req dto.RequestOTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.WithError(err).Error("Invalid request body")
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "INVALID_REQUEST", Message: "Invalid request body"})
        return
    }
    if err := h.requestOTP.Execute(c.Request.Context(), req.Phone); err != nil {
        h.handleError(c, err)
        return
    }
    c.JSON(http.StatusOK, dto.MessageResponse{Message: "OTP sent successfully. Check server logs (dev)."})
}

// VerifyOTP godoc
// @Summary      Verify OTP
// @Description  Verify OTP and return JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      dto.VerifyOTPRequest  true  "Phone and OTP"
// @Success      200      {object}  dto.AuthResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      401      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /auth/otp/verify [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
    var req dto.VerifyOTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.WithError(err).Error("Invalid request body")
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "INVALID_REQUEST", Message: "Invalid request body"})
        return
    }
    result, err := h.verifyOTP.Execute(c.Request.Context(), req.Phone, req.OTP)
    if err != nil {
        h.handleError(c, err)
        return
    }
    c.JSON(http.StatusOK, dto.AuthResponse{
        Token:     result.Token,
        ExpiresIn: result.ExpiresIn,
        User: &dto.UserResponse{
            ID:        result.User.ID.String(),
            Phone:     result.User.Phone,
            CreatedAt: result.User.CreatedAt,
            UpdatedAt: result.User.UpdatedAt,
        },
    })
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
    switch e := err.(type) {
    case *appErr.ValidationError:
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "VALIDATION_ERROR", Message: e.Error()})
    case *appErr.NotFoundError:
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "NOT_FOUND", Message: e.Error()})
    case *appErr.UnauthorizedError:
        c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "UNAUTHORIZED", Message: e.Error()})
    case *appErr.RateLimitError:
        c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: "RATE_LIMIT_EXCEEDED", Message: e.Error()})
    default:
        h.logger.WithError(err).Error("Internal server error")
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "INTERNAL_ERROR", Message: "An internal error occurred"})
    }
}