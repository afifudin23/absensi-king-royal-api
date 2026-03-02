package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerAttendanceRoutes(rg *gin.RouterGroup) {
	attendanceHandler := handler.NewAttendanceHandler()
	attendance := rg.Group("/attendance")

	attendance.Use(middleware.AuthMiddleware())

	{
		attendance.POST("/check-in", attendanceHandler.CheckIn)
		attendance.POST("/check-out", attendanceHandler.CheckOut)
		attendance.GET("/logs", attendanceHandler.GetLogs)
	}
}
