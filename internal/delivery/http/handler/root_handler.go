package handler

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/gin-gonic/gin"
)

type rootData struct {
	Message    string `json:"message"`
	Status     string `json:"status"`
	AppVersion string `json:"app_version"`
}

func Root(c *gin.Context) {
	payload := rootData{
		Message:    "Welcome to the Absensi King Royal API",
		Status:     "ok",
		AppVersion: config.AppVersion,
	}

	c.JSON(200, common.SuccessResponse(payload))
}
