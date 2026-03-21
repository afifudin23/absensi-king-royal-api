package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

func registerPayrollSetting(rg *gin.RouterGroup) {
	db := config.GetDB()
	payrollSettingRepo := repository.NewPayrollSettingRepository(db)
	payrollSettingService := service.NewPayrollSettingService(payrollSettingRepo)
	payrollSettingHandler := handler.NewPayrollSettingHandler(payrollSettingService)
	payrollSetting := rg.Group("/payroll-settings")

	{
		payrollSetting.GET("", payrollSettingHandler.GetAll)
		payrollSetting.POST("", payrollSettingHandler.Create)
		payrollSetting.PUT("/bulk", payrollSettingHandler.UpdateBulk)
		payrollSetting.PATCH("/:payroll_id", payrollSettingHandler.Update)
		payrollSetting.DELETE("/:payroll_id", payrollSettingHandler.Delete)
	}

}
