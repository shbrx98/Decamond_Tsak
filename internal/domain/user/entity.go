package user

import (
    "time"

    "github.com/google/uuid"
    appErr "github.com/shbrx98/Decamond_Tsak/internal/pkg/errors"
)

type User struct {
    ID        uuid.UUID `json:"id" db:"id"`
    Phone     string    `json:"phone" db:"phone"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

var ErrInvalidPhone = appErr.NewValidationError("invalid phone")

func NewUser(phone string) *User {
    now := time.Now()
    return &User{
        ID:        uuid.New(),
        Phone:     phone,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func (u *User) Validate() error {
    if u.Phone == "" {
        return ErrInvalidPhone
    }
    return nil
}


func (u *User) Touch() { u.UpdatedAt = time.Now() }