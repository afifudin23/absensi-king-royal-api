package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(rg *gin.RouterGroup) {
	authHandler := handler.NewAuthHandler()

	auth := rg.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
}
