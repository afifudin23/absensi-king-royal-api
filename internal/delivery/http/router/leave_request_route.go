package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

func registerLeaveRequestRoutes(rg *gin.RouterGroup) {
	db := config.GetDB()
	leaveRepo := repository.NewLeaveRequestRepository(db)
	fileRepo := repository.NewFileRepository(db)
	leaveService := service.NewLeaveRequestService(leaveRepo, fileRepo)
	leaveRequestHandler := handler.NewLeaveRequestHandler(leaveService)
	leaveRequest := rg.Group("/leave-requests")
	
	leaveRequest.Use(middleware.AuthMiddleware())
	{
		leaveRequest.GET("", leaveRequestHandler.GetAll)
		leaveRequest.POST("", leaveRequestHandler.Create)
		leaveRequest.GET("/me", leaveRequestHandler.GetByUserID)
		leaveRequest.GET("/:leave_id", leaveRequestHandler.GetByID)
		leaveRequest.PUT("/:leave_id", leaveRequestHandler.Update)
		leaveRequest.DELETE("/:leave_id", leaveRequestHandler.Delete)
	}

}
