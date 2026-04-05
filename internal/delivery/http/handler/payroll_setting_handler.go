package handler

import (
	"errors"
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

type PayrollSettingHandler struct {
	service service.PayrollSettingService
}

func NewPayrollSettingHandler(payrollSettingService service.PayrollSettingService) *PayrollSettingHandler {
	return &PayrollSettingHandler{service: payrollSettingService}
}

func (h *PayrollSettingHandler) GetAll(c *gin.Context) {
	payrollSettings, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollSettingListResponse(payrollSettings)))

}

func (h *PayrollSettingHandler) Create(c *gin.Context) {
	var payload request.PayrollSettingRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	payrollSetting, err := h.service.Create(c.Request.Context(), payload)
	if err != nil {
		if errors.Is(err, service.ErrPayrollSettingAlreadyExists) {
			common.ErrorHandler(c, common.BadRequestError(err.Error()))
			return
		}
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToPayrollSettingResponse(*payrollSetting)))
}

func (h *PayrollSettingHandler) Update(c *gin.Context) {
	var payload request.PayrollSettingRequest
	payrollID := c.Param("payroll_id")

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	payrollSetting, err := h.service.Update(c.Request.Context(), payrollID, payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToPayrollSettingResponse(*payrollSetting)))
}

func (h *PayrollSettingHandler) UpdateBulk(c *gin.Context) {
	var payload request.PayrollSettingUpdateBulkRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	payrollSettings, err := h.service.UpdateBulk(c.Request.Context(), payload.Settings)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToPayrollSettingListResponse(payrollSettings)))
}

func (h *PayrollSettingHandler) Delete(c *gin.Context) {
	var payload request.PayrollSettingIdsRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	deletedCount, err := h.service.Delete(c.Request.Context(), payload)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	// Semantics:
	// - total = jumlah data yang berhasil dihapus (kalau tidak ada yang match: 0)
	// - skipped_count tidak dipakai di payroll settings delete (0)
	total := deletedCount
	c.JSON(
		http.StatusOK,
		common.SuccessResponse(response.ToPayrollSettingDeleteResponse(total, deletedCount, 0)),
	)
}
