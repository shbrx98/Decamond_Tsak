package dto

import (
    "time"

    domUser "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
)

type RequestOTPRequest struct {
    Phone string `json:"phone" binding:"required" example:"+989121234567"`
}

type VerifyOTPRequest struct {
    Phone string `json:"phone" binding:"required" example:"+989121234567"`
    OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
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


type PaginatedUsersResponse struct {
    Users      []*UserResponse `json:"users"`
    Total      int             `json:"total"`
    Page       int             `json:"page"`
    PerPage    int             `json:"per_page"`
    TotalPages int             `json:"total_pages"`
}

func PaginatedUsersFromDomain(res *domUser.ListResult) *PaginatedUsersResponse {
    out := make([]*UserResponse, 0, len(res.Users))
    for _, u := range res.Users {
        out = append(out, UserFromDomain(u))
    }
    return &PaginatedUsersResponse{
        Users:      out,
        Total:      res.Total,
        Page:       res.Page,
        PerPage:    res.PerPage,
        TotalPages: res.TotalPages,
    }
}