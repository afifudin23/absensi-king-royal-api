package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/gin-gonic/gin"
)

func registerActivityLogRoutes(rg *gin.RouterGroup) {
	repo := repository.NewActivityLogRepository(config.GetDB())
	h := handler.NewActivityLogHandler(repo)

	group := rg.Group("/activity-logs")
	group.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		group.GET("", h.GetAll)
	}
}
