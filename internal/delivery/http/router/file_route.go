package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

func registerFileRoutes(rg *gin.RouterGroup) {
	db := config.GetDB()
	fileRepo := repository.NewFileRepository(db)
	userRepo := repository.NewUserRepository(db)
	fileService := service.NewFileService(fileRepo, config.GetEnv().ServerBaseURL)
	userService := service.NewUserService(userRepo, fileRepo)
	fileHandler := handler.NewFileHandler(fileService, userService)
	router := rg.Group("/files")

	router.Use(middleware.AuthMiddleware())
	{
		router.POST("", fileHandler.Upload)
		router.DELETE("/:file_id", fileHandler.Delete)
	}
}
