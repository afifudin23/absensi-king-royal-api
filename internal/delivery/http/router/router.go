package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/afifudin23/absensi-king-royal-api/docs"
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
		
		// Set Swagger Host dynamically so it works from both localhost and ngrok
		docs.SwaggerInfo.Host = c.Request.Host
		if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
			docs.SwaggerInfo.Schemes = []string{"https"}
		} else {
			docs.SwaggerInfo.Schemes = []string{"http", "https"}
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
