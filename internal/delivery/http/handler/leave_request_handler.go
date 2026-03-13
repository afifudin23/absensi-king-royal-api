package handler

import (
	"net/http"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type LeaveRequestHandler struct {
	service service.LeaveRequestService
}

func NewLeaveRequestHandler() *LeaveRequestHandler {
	return &LeaveRequestHandler{service: service.NewLeaveRequestService()}
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
	startDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		common.ErrorHandler(c, common.BadRequestError("start_date must be in YYYY-MM-DD format"))
		return
	}

	endDate, err := time.Parse("2006-01-02", payload.EndDate)
	if err != nil {
		common.ErrorHandler(c, common.BadRequestError("end_date must be in YYYY-MM-DD format"))
		return
	}
	data := &model.LeaveRequest{
		UserID:           userID,
		StartDate:        startDate,
		EndDate:          endDate,
		Reason:           payload.Reason,
		Type:             payload.Type,
		EvidenceURL:      payload.EvidenceURL,
		EvidencePublicID: payload.EvidencePublicID,
		OvertimeHours:    payload.OvertimeHours,
		Status:           model.LeaveRequestStatusPending,
	}
	err = h.service.Create(data)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(common.ToSuccessResponse(data.ID)))
}

func (h *LeaveRequestHandler) GetAll(c *gin.Context) {
	leaves, err := h.service.GetAll()
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToLeaveListResponse(leaves)))
}

func (h *LeaveRequestHandler) GetByID(c *gin.Context) {
	leaveID := c.Param("leave_id")
	leave, err := h.service.GetByID(leaveID)
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
	leaves, err := h.service.GetByUserID(userID)
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

	leaveID := c.Param("leave_id")
	data := &model.LeaveRequest{ID: leaveID}

	if payload.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *payload.StartDate)
		if err != nil {
			common.ErrorHandler(c, common.BadRequestError("start_date must be in YYYY-MM-DD format"))
			return
		}
		data.StartDate = startDate
	}

	if payload.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *payload.EndDate)
		if err != nil {
			common.ErrorHandler(c, common.BadRequestError("end_date must be in YYYY-MM-DD format"))
			return
		}
		data.EndDate = endDate
	}

	if payload.Reason != nil {
		data.Reason = *payload.Reason
	}

	if payload.Type != nil {
		data.Type = *payload.Type
	}

	if payload.EvidenceURL != nil {
		data.EvidenceURL = payload.EvidenceURL
	}

	if payload.EvidencePublicID != nil {
		data.EvidencePublicID = payload.EvidencePublicID
	}

	if payload.OvertimeHours != nil {
		data.OvertimeHours = payload.OvertimeHours
	}

	err := h.service.Update(data)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(data.ID)))
}

func (h *LeaveRequestHandler) Delete(c *gin.Context) {
	leaveID := c.Param("leave_id")
	err := h.service.Delete(leaveID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToSuccessResponse(leaveID)))
}
