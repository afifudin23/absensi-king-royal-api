package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

func registerPayroll(rg *gin.RouterGroup) {
	db := config.GetDB()
	payrollRepo := repository.NewPayrollRepository(db)
	userRepo := repository.NewUserRepository(db)
	payrollSettingRepo := repository.NewPayrollSettingRepository(db)
	payrollService := service.NewPayrollService(payrollRepo, payrollSettingRepo, userRepo)
	payrollHandler := handler.NewPayrollHandler(payrollService)

	payroll := rg.Group("/payrolls")
	payroll.Use(middleware.AuthMiddleware())
	{
		payroll.GET("", payrollHandler.GetAll)
		payroll.GET("/:payroll_id", payrollHandler.GetByID)
		payroll.POST("/generate/:employee_id", payrollHandler.GenerateOne)
		payroll.POST("/generate-all", payrollHandler.GenerateAll)
		payroll.PUT("/:payroll_id", payrollHandler.Update)
	}
}
