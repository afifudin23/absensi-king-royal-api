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

type LeaveRequestHandler struct {
	service service.LeaveRequestService
}

func NewLeaveRequestHandler(leaveRequestService service.LeaveRequestService) *LeaveRequestHandler {
	return &LeaveRequestHandler{service: leaveRequestService}
}

func (h *LeaveRequestHandler) Create(c *gin.Context) {
	var payload request.LeaveRequestCreateRequest
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

func (h *LeaveRequestHandler) GetAll(c *gin.Context) {
	leaves, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToLeaveListResponse(leaves)))
}

func (h *LeaveRequestHandler) GetByID(c *gin.Context) {
	leaveID := c.Param("leave_id")
	leave, err := h.service.GetByID(c.Request.Context(), leaveID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToLeaveResponse(*leave)))

}

func (h *LeaveRequestHandler) GetByUserID(c *gin.Context) {
	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}
	leaves, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToLeaveListResponse(leaves)))
}

func (h *LeaveRequestHandler) Update(c *gin.Context) {
	var payload request.LeaveRequestUpdateRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}
	payload.Normalize()

	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}

	leaveID := c.Param("leave_id")
	data, err := h.service.Update(c.Request.Context(), userID, leaveID, payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(data.ID)))
}

func (h *LeaveRequestHandler) Delete(c *gin.Context) {
	leaveID := c.Param("leave_id")
	err := h.service.Delete(c.Request.Context(), leaveID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(leaveID)))
}
