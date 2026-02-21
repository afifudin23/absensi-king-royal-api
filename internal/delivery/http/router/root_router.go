package router

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func registerRootRoutes(rg *gin.RouterGroup) {
	rg.GET("/", handler.Root)
}
