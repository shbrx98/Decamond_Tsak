package auth

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
    domUser "github.com/shbrx98/Decamond_Tsak/internal/domain/user"
)

func GenerateJWT(secret string, u *domUser.User, expiry time.Duration, issuer, audience string) (string, int64, error) {
    expiresAt := time.Now().Add(expiry)
    claims := jwt.RegisteredClaims{
        Subject:   u.ID.String(),
        ExpiresAt: jwt.NewNumericDate(expiresAt),
        Issuer:    issuer,
        Audience:  []string{audience},
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    // custom claims
    type Claims struct {
        Phone string `json:"phone"`
        jwt.RegisteredClaims
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
        Phone:            u.Phone,
        RegisteredClaims: claims,
    })
    signed, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", 0, err
    }
    return signed, int64(expiry.Seconds()), nil
}