package handler

import (
	"net/http"
	"strconv"
	"time"

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

	var startDate, endDate *time.Time
	if s := c.Query("start_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			common.ErrorHandler(c, common.BadRequestError("Invalid start_date format, use YYYY-MM-DD"))
			return
		}
		startDate = &t
	}
	if e := c.Query("end_date"); e != "" {
		t, err := time.Parse("2006-01-02", e)
		if err != nil {
			common.ErrorHandler(c, common.BadRequestError("Invalid end_date format, use YYYY-MM-DD"))
			return
		}
		endDate = &t
	}

	logs, err := h.service.GetLogs(c.Request.Context(), userID, startDate, endDate)
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

func (h *AttendanceHandler) GetRecap(c *gin.Context) {
	monthStr := c.DefaultQuery("month", "")
	yearStr := c.DefaultQuery("year", "")

	if monthStr == "" || yearStr == "" {
		common.ErrorHandler(c, common.BadRequestError("month and year query params are required"))
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		common.ErrorHandler(c, common.BadRequestError("Invalid month, must be 1-12"))
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		common.ErrorHandler(c, common.BadRequestError("Invalid year"))
		return
	}

	recap, err := h.service.GetRecap(c.Request.Context(), month, year)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(recap))
}

func (h *AttendanceHandler) Update(c *gin.Context) {
	var payload request.AttendanceUpdateRequest
	attendanceID := c.Param("attendance_id")

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	updaterID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	attendance, err := h.service.Update(c.Request.Context(), updaterID, attendanceID, payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(attendance.ID)))
}
