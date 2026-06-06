package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.StructuredLoggingMiddleware(),
		middleware.RecoveryMiddleware(),
		middleware.ErrorMiddleware(),
	)

	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/docs/index.html")
	})
	r.GET("/docs/*any", func(c *gin.Context) {
		if c.Param("any") == "/" {
			c.Redirect(301, "/docs/index.html")
			return
		}
		swaggerHandler(c)
	})
	registerRootRoutes(r.Group(""))
	registerHealthRoutes(r.Group(""))

	api := r.Group("/api")
	v1 := api.Group("/v1")

	registerAuthRoutes(v1)
	registerUserRouter(v1)
	registerAttendanceRoutes(v1)
	registerAttendanceRequestRoutes(v1)
	registerFileRoutes(v1)
	registerPayrollSetting(v1)
	registerPayroll(v1)
	registerActivityLogRoutes(v1)

	return r
}
