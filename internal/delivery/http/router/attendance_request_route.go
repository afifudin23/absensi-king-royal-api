package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

func registerAttendanceRequestRoutes(rg *gin.RouterGroup) {
	db := config.GetDB()
	attendanceRequestRepo := repository.NewAttendanceRequestRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	fileRepo := repository.NewFileRepository(db)
	attendanceRequestService := service.NewAttendanceRequestService(attendanceRequestRepo, attendanceRepo, fileRepo)
	attendanceRequestHandler := handler.NewAttendanceRequestHandler(attendanceRequestService)
	attendanceRequest := rg.Group("/attendance-requests")

	attendanceRequest.Use(middleware.AuthMiddleware())
	{
		attendanceRequest.GET("", attendanceRequestHandler.GetAll)
		attendanceRequest.POST("", attendanceRequestHandler.Create)
		attendanceRequest.GET("/me", attendanceRequestHandler.GetByUserID)
		attendanceRequest.GET("/:attendance_request_id", attendanceRequestHandler.GetByID)
		attendanceRequest.PUT("/:attendance_request_id", attendanceRequestHandler.Update)
		attendanceRequest.PATCH("/:attendance_request_id/status", middleware.AdminOnly(), attendanceRequestHandler.UpdateStatus)
		attendanceRequest.DELETE("/:attendance_request_id", attendanceRequestHandler.Delete)
	}
}
