package handler

import (
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
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

func NewAttendanceHandler(attendanceService service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: attendanceService}
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	var payload request.AttendanceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	attendance, err := h.service.CheckIn(c.Request.Context(), userID, payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceResponse(*attendance)))
}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	var payload request.AttendanceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	attendance, err := h.service.CheckOut(c.Request.Context(), userID, payload)
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

	logs, err := h.service.GetLogs(c.Request.Context(), userID)
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
