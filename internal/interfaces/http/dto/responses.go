package dto

import (
    "time"

    domUser "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
)

type RequestOTPRequest struct {
    Phone string `json:"phone" binding:"required"`
}

type VerifyOTPRequest struct {
    Phone string `json:"phone" binding:"required"`
    OTP   string `json:"otp" binding:"required,len=6"`
}

type MessageResponse struct {
    Message string `json:"message"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Phone     string    `json:"phone"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
    Token     string        `json:"token"`
    ExpiresIn int64         `json:"expires_in"`
    User      *UserResponse `json:"user"`
}

func UserFromDomain(u *domUser.User) *UserResponse {
    return &UserResponse{
        ID:        u.ID.String(),
        Phone:     u.Phone,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}