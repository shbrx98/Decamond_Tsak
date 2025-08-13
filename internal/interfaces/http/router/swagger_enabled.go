//go:build swagger

package router

import (
    "github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"

    docs "github.com/shbrx98/Decamond_Tsak/docs"
)

func attachSwagger(r *gin.Engine) {
    docs.SwaggerInfo.Title = "DECAMOND OTP Authentication Service API"
    docs.SwaggerInfo.Version = "1.0.0"
    docs.SwaggerInfo.BasePath = "/api/v1"

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}