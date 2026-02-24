package router

import (
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	v1 := api.Group("/v1")

	registerHealthRoutes(v1)
	registerRootRoutes(v1)
	registerAuthRoutes(v1)

	return r
}
