package handler

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/gin-gonic/gin"
)

type healthData struct {
	Status string `json:"status"`
}

// Health godoc
//
//	@Summary		Health check
//	@Description	Mengecek status koneksi database dan server
//	@Tags			General
//	@Produce		json
//	@Success		200	{object}	common.Response[healthData]
//	@Router			/health [get]
func Health(c *gin.Context) {
	c.JSON(200, common.SuccessResponse(healthData{Status: "ok"}))
}
