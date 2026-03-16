package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

func registerAttendanceRoutes(rg *gin.RouterGroup) {
	db := config.GetDB()
	attendanceRepo := repository.NewAttendanceRepository(db)
	fileRepo := repository.NewFileRepository(db)
	attendanceService := service.NewAttendanceService(attendanceRepo, fileRepo)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	attendance := rg.Group("/attendance")

	attendance.Use(middleware.AuthMiddleware())

	{
		attendance.POST("/check-in", attendanceHandler.CheckIn)
		attendance.POST("/check-out", attendanceHandler.CheckOut)
		attendance.GET("/logs", attendanceHandler.GetLogs)
	}
}
