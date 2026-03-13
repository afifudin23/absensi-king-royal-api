package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.StructuredLoggingMiddleware(),
		middleware.RecoveryMiddleware(),
		middleware.ErrorMiddleware(),
	)

	api := r.Group("/api")
	v1 := api.Group("/v1")

	registerHealthRoutes(v1)
	registerRootRoutes(v1)
	registerAuthRoutes(v1)
	registerUserRouter(v1)
	registerAttendanceRoutes(v1)
	registerLeaveRequestRoutes(v1)

	return r
}
