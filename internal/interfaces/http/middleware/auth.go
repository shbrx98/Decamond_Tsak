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

    AcceptWithoutBearer bool   
    AcceptFromCookie     bool   
    CookieName           string 
    AcceptFromQuery      bool   
    QueryParam           string 
}

func AuthMiddleware(cfg JWTConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := extractToken(c, cfg)
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error":   "MISSING_TOKEN",
                "message": "Authorization header is required (Use: Bearer <token>)",
            })
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
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
            c.JSON(http.StatusUnauthorized, gin.H{
                "error":   "INVALID_TOKEN",
                "message": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("user_id", claims["sub"])
            c.Set("phone", claims["phone"])
        }
        c.Next()
    }
}

func extractToken(c *gin.Context, cfg JWTConfig) string {
    // 1) Authorization header
    hdr := strings.TrimSpace(c.GetHeader("Authorization"))
    if hdr != "" {
        parts := strings.Fields(hdr)
        if len(parts) >= 2 && strings.EqualFold(parts[0], "Bearer") {
            return strings.TrimSpace(parts[1])
        }
        if cfg.AcceptWithoutBearer {
         
            return hdr
        }
     
        return ""
    }

    // 2) Cookie
    if cfg.AcceptFromCookie && cfg.CookieName != "" {
        if tok, err := c.Cookie(cfg.CookieName); err == nil && strings.TrimSpace(tok) != "" {
            return strings.TrimSpace(tok)
        }
    }


    if cfg.AcceptFromQuery && cfg.QueryParam != "" {
        if tok := strings.TrimSpace(c.Query(cfg.QueryParam)); tok != "" {
            return tok
        }
    }

    return ""
}