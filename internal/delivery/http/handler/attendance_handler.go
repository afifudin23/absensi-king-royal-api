package handler

import (
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/afifudin23/absensi-king-royal-api/pkg/logger"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service service.AttendanceService
}

func NewAttendanceHandler() *AttendanceHandler {
	return &AttendanceHandler{service: service.NewAttendanceService()}
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	attendance, err := h.service.CheckIn(userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceResponse(*attendance)))
}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	attendance, err := h.service.CheckOut(userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceResponse(*attendance)))
}

func (h *AttendanceHandler) GetLogs(c *gin.Context) {
	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	logs, err := h.service.GetLogs(userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	logger.Info(
		c.Request.Context(),
		"attendance.handler",
		"attendance get logs",
		map[string]any{"user_id": userID},
	)
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceListResponse(logs)))
}
