package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(rg *gin.RouterGroup) {
	db := config.GetDB()
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	auth := rg.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
}
