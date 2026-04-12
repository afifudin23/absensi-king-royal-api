package handler

import (
	"net/http"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/service"
	"github.com/gin-gonic/gin"
)

type PayrollHandler struct {
	service service.PayrollService
}

func NewPayrollHandler(service service.PayrollService) *PayrollHandler {
	return &PayrollHandler{service: service}
}

func (h *PayrollHandler) GetAll(c *gin.Context) {
	payrolls, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollListResponse(payrolls)))
}

func (h *PayrollHandler) GetByID(c *gin.Context) {
	payrollID := c.Param("payroll_id")
	payroll, err := h.service.GetByID(c.Request.Context(), payrollID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollResponse(payroll)))
}

func (h *PayrollHandler) GenerateOne(c *gin.Context) {
	employeeID := c.Param("employee_id")
	payroll, err := h.service.GenerateOne(c.Request.Context(), employeeID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToPayrollResponse(payroll)))
}

func (h *PayrollHandler) GenerateAll(c *gin.Context) {
	payrolls, err := h.service.GenerateAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToPayrollListResponse(payrolls)))
}

func (h *PayrollHandler) Update(c *gin.Context) {
	var payload request.PayrollUpdateRequest
	payrollID := c.Param("payroll_id")

	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(common.ErrorValidation(err)))
		return
	}

	payroll, err := h.service.Update(c.Request.Context(), payrollID, payload)
	if err != nil {

		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollResponse(payroll)))
}

func (h *PayrollHandler) SendPayroll(c *gin.Context) {
	payrollID := c.Param("payroll_id")
	payroll, err := h.service.SendPayroll(c.Request.Context(), payrollID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollResponse(payroll)))
}
