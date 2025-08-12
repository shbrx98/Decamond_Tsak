package middleware

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
    Secret   string
    Issuer   string
    Audience string
}

func AuthMiddleware(cfg JWTConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "MISSING_TOKEN", "message": "Authorization header is required"})
            c.Abort()
            return
        }
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_TOKEN_FORMAT", "message": "Use: Bearer <token>"})
            c.Abort()
            return
        }
        tokenString = parts[1]

        type Claims struct {
            Phone string `json:"phone"`
            jwt.RegisteredClaims
        }

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(cfg.Secret), nil
        },
            jwt.WithValidMethods([]string{"HS256"}),
            jwt.WithIssuer(cfg.Issuer),
            jwt.WithAudience(cfg.Audience),
            jwt.WithLeeway(1*time.Minute),
        )
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_TOKEN", "message": "Invalid or expired token"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(*Claims); ok {
            c.Set("user_id", claims.Subject)
            c.Set("phone", claims.Phone)
        }

        c.Next()
    }
}