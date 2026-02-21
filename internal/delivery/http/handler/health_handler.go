package handler

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/gin-gonic/gin"
)

type healthData struct {
	Status string `json:"status"`
}

func Health(c *gin.Context) {
	c.JSON(200, common.SuccessResponse(healthData{Status: "ok"}))
}
