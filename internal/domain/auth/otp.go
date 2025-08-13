package auth

import (
    "crypto/rand"
    "fmt"
    "time"
)

type OTP struct {
    Code      string
    Phone     string
    ExpiresAt time.Time
}

func NewOTP(phone string, ttl time.Duration) *OTP {
    return &OTP{
        Code:      generateOTPCode(),
        Phone:     phone,
        ExpiresAt: time.Now().Add(ttl),
    }
}

func generateOTPCode() string {
    b := make([]byte, 3)
    if _, err := rand.Read(b); err != nil {
        return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
    }
    return fmt.Sprintf("%06d", (int(b[0])<<16|int(b[1])<<8|int(b[2]))%1000000)
}