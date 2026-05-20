package handler

import (
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AttendanceRequestHandler struct {
	service service.AttendanceRequestService
}

func NewAttendanceRequestHandler(attendanceRequestService service.AttendanceRequestService) *AttendanceRequestHandler {
	return &AttendanceRequestHandler{service: attendanceRequestService}
}

func (h *AttendanceRequestHandler) Create(c *gin.Context) {
	var payload request.AttendanceRequestCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	data, err := h.service.Create(c.Request.Context(), userID, payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(common.ToSuccessResponse(data.ID)))
}

func (h *AttendanceRequestHandler) GetAll(c *gin.Context) {
	items, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceRequestListResponse(items)))
}

func (h *AttendanceRequestHandler) GetByID(c *gin.Context) {
	attendanceRequestID := c.Param("attendance_request_id")
	item, err := h.service.GetByID(c.Request.Context(), attendanceRequestID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceRequestResponse(*item)))
}

func (h *AttendanceRequestHandler) GetByUserID(c *gin.Context) {
	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}
	items, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceRequestListResponse(items)))
}

func (h *AttendanceRequestHandler) Update(c *gin.Context) {
	var payload request.AttendanceRequestUpdateRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	attendanceRequestID := c.Param("attendance_request_id")
	data, err := h.service.Update(c.Request.Context(), userID, attendanceRequestID, payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(data.ID)))
}

func (h *AttendanceRequestHandler) UpdateStatus(c *gin.Context) {
	var payload request.AttendanceRequestUpdateStatusRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	reviewerID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	attendanceRequestID := c.Param("attendance_request_id")
	data, err := h.service.UpdateStatus(c.Request.Context(), reviewerID, attendanceRequestID, payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(data.ID)))
}

func (h *AttendanceRequestHandler) Delete(c *gin.Context) {
	attendanceRequestID := c.Param("attendance_request_id")
	err := h.service.Delete(c.Request.Context(), attendanceRequestID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(attendanceRequestID)))
}
