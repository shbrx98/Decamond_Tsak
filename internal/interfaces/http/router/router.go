package router

import (
    "net/http"

    "github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"

    appAuth "github.com/shbrx98/Decamond_Tsak/internal/application/auth"
    appUser "github.com/shbrx98/Decamond_Tsak/internal/application/user"
    "github.com/shbrx98/Decamond_Tsak/internal/infrastructure/config"
    "github.com/shbrx98/Decamond_Tsak/internal/interfaces/http/handlers"
    "github.com/shbrx98/Decamond_Tsak/internal/interfaces/http/middleware"
    "github.com/shbrx98/Decamond_Tsak/internal/pkg/logger"
)

func NewRouter(
    requestOTP *appAuth.RequestOTPUseCase,
    verifyOTP *appAuth.VerifyOTPUseCase,
    getUser *appUser.GetUserUseCase,
    listUsers *appUser.ListUsersUseCase,
    cfg *config.Config,
    log logger.Logger,
) *gin.Engine {
    if cfg.Server.Env == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.New()
    r.Use(gin.Recovery())
    // Optional: structured request logging middleware if desired

    // Health
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    // Swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Handlers
    authHandler := handlers.NewAuthHandler(requestOTP, verifyOTP, log)
    userHandler := handlers.NewUserHandler(getUser, listUsers, log)

    api := r.Group("/api/v1")
    {
        auth := api.Group("/auth")
        auth.POST("/otp/request", authHandler.RequestOTP)
        auth.POST("/otp/verify", authHandler.VerifyOTP)

        // protected
        jwtCfg := middleware.JWTConfig{
            Secret:   cfg.JWT.Secret,
            Issuer:   cfg.JWT.Issuer,
            Audience: cfg.JWT.Audience,
        }
        users := api.Group("/users", middleware.AuthMiddleware(jwtCfg))
        users.GET("", userHandler.List)
        users.GET("/:id", userHandler.Get)
    }

    return r
}