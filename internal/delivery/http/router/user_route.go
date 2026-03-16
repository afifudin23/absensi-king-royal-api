package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

func registerUserRouter(rg *gin.RouterGroup) {
	db := config.GetDB()
	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)
	userService := service.NewUserService(userRepo, fileRepo)
	userHandler := handler.NewUserHandler(userService)
	users := rg.Group("/users")

	users.Use(middleware.AuthMiddleware())

	{
		users.GET("", userHandler.GetAllUsers)
		users.GET("/me", userHandler.GetMyProfile)
		users.PUT("/me", userHandler.UpdateMyProfile)
		users.POST("", userHandler.CreateUser)
		users.GET("/:user_id", userHandler.GetUserByID)
		users.PUT("/:user_id", userHandler.UpdateUser)
		users.DELETE("/:user_id", userHandler.DeleteUser)
	}
}
