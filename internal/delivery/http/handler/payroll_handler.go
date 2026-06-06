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

type PayrollHandler struct {
	service service.PayrollService
}

func NewPayrollHandler(service service.PayrollService) *PayrollHandler {
	return &PayrollHandler{service: service}
}

// GetMyPayrolls godoc
//
//	@Summary		Payroll saya
//	@Description	Mendapatkan daftar slip gaji milik pengguna yang sedang login.
//	@Description
//	@Description	**Response 200:** Array slip gaji milik pengguna (gaji pokok, tunjangan, potongan, gaji bersih, status pengiriman)
//	@Tags			Payroll
//	@Produce		json
//	@Success		200	{object}	common.Response[[]response.PayrollResponse]
//	@Failure		401	{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payrolls/me [get]
func (h *PayrollHandler) GetMyPayrolls(c *gin.Context) {
	userID, ok := utils.GetCurrentUserID(c)
	if !ok {
		return
	}
	payrolls, err := h.service.GetMyPayrolls(c.Request.Context(), userID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollListResponse(payrolls)))
}

// GetAll godoc
//
//	@Summary		Semua payroll
//	@Description	Mendapatkan daftar semua slip gaji seluruh karyawan.
//	@Description
//	@Description	**Response 200:** Array data payroll (gaji pokok, tunjangan, potongan, gaji bersih, status, pdf path, dll)
//	@Tags			Payroll
//	@Produce		json
//	@Success		200	{object}	common.Response[[]response.PayrollResponse]
//	@Failure		401	{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payrolls [get]
func (h *PayrollHandler) GetAll(c *gin.Context) {
	payrolls, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollListResponse(payrolls)))
}

// GetByID godoc
//
//	@Summary		Detail payroll
//	@Description	Mendapatkan detail satu slip gaji berdasarkan ID.
//	@Description
//	@Description	**Response 200:** Data payroll lengkap (gaji pokok, semua tunjangan, semua potongan, gaji kotor, gaji bersih, status, pdf)
//	@Tags			Payroll
//	@Produce		json
//	@Param			payroll_id	path		string	true	"ID payroll"
//	@Success		200			{object}	common.Response[response.PayrollResponse]
//	@Failure		401			{object}	common.Response[any]
//	@Failure		404			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payrolls/{payroll_id} [get]
func (h *PayrollHandler) GetByID(c *gin.Context) {
	payrollID := c.Param("payroll_id")
	payroll, err := h.service.GetByID(c.Request.Context(), payrollID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollResponse(payroll)))
}

// GenerateOne godoc
//
//	@Summary		Generate payroll satu karyawan
//	@Description	Membuat slip gaji untuk satu karyawan berdasarkan ID karyawan.
//	@Description
//	@Description	**Path Param:**
//	@Description	- `employee_id` (UUID): ID karyawan yang akan dibuatkan payroll
//	@Description
//	@Description	**Behavior:**
//	@Description	- Menghitung gaji berdasarkan data profil karyawan (gaji pokok, tunjangan) dan konfigurasi payroll
//	@Description	- Memperhitungkan rekap absensi (lembur, potongan ketidakhadiran)
//	@Description	- File PDF slip gaji dibuat otomatis
//	@Description
//	@Description	**Response 201:** Data payroll yang baru digenerate
//	@Tags			Payroll
//	@Produce		json
//	@Param			employee_id	path		string	true	"ID karyawan"
//	@Success		201			{object}	common.Response[response.PayrollResponse]
//	@Failure		401			{object}	common.Response[any]
//	@Failure		404			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payrolls/generate/{employee_id} [post]
func (h *PayrollHandler) GenerateOne(c *gin.Context) {
	employeeID := c.Param("employee_id")
	payroll, err := h.service.GenerateOne(c.Request.Context(), employeeID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToPayrollResponse(payroll)))
}

// GenerateAll godoc
//
//	@Summary		Generate payroll semua karyawan
//	@Description	Membuat slip gaji untuk semua karyawan aktif secara sekaligus.
//	@Description
//	@Description	**Behavior:**
//	@Description	- Memproses semua karyawan dengan role `user` yang aktif
//	@Description	- Setiap payroll dihitung berdasarkan data profil dan rekap absensi masing-masing karyawan
//	@Description	- PDF slip gaji dibuat untuk setiap karyawan
//	@Description
//	@Description	**Response 201:** Array seluruh data payroll yang baru digenerate
//	@Tags			Payroll
//	@Produce		json
//	@Success		201	{object}	common.Response[[]response.PayrollResponse]
//	@Failure		401	{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payrolls/generate-all [post]
func (h *PayrollHandler) GenerateAll(c *gin.Context) {
	payrolls, err := h.service.GenerateAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response.ToPayrollListResponse(payrolls)))
}

// Update godoc
//
//	@Summary		Update payroll
//	@Description	Memperbarui komponen data slip gaji karyawan secara manual.
//	@Description
//	@Description	**Request Body (semua opsional):**
//	@Description	- `basic_salary` (number): Gaji pokok
//	@Description	- `position_allowance` (number): Tunjangan jabatan
//	@Description	- `other_allowance` (number): Tunjangan lain
//	@Description	- `overtime_rate` (number): Total nilai lembur
//	@Description	- `loan_deduction` (number): Potongan pinjaman
//	@Description	- `attendance_deduction` (number): Potongan absensi
//	@Description	- `income_tax` (number): Pajak penghasilan
//	@Description	- `additional_data` (JSON): Data tambahan
//	@Description
//	@Description	**Response 200:** Data payroll setelah diupdate (gaji kotor dan bersih dihitung ulang)
//	@Tags			Payroll
//	@Accept			json
//	@Produce		json
//	@Param			payroll_id	path		string						true	"ID payroll"
//	@Param			body		body		request.PayrollUpdateRequest	true	"Data yang diperbarui"
//	@Success		200			{object}	common.Response[response.PayrollResponse]
//	@Failure		400			{object}	common.Response[any]
//	@Failure		401			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payrolls/{payroll_id} [put]
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

// SendPayroll godoc
//
//	@Summary		Kirim payroll via email
//	@Description	Mengirimkan slip gaji dalam format PDF ke alamat email karyawan.
//	@Description
//	@Description	**Path Param:**
//	@Description	- `payroll_id` (UUID): ID payroll yang akan dikirim
//	@Description
//	@Description	**Behavior:**
//	@Description	- PDF slip gaji dikirim sebagai lampiran email ke alamat email karyawan
//	@Description	- Waktu pengiriman (`sent_at`) dicatat otomatis setelah berhasil dikirim
//	@Description	- Status payroll diperbarui menjadi `sent`
//	@Description
//	@Description	**Response 200:** Data payroll dengan `sent_at` yang telah diisi
//	@Tags			Payroll
//	@Produce		json
//	@Param			payroll_id	path		string	true	"ID payroll"
//	@Success		200			{object}	common.Response[response.PayrollResponse]
//	@Failure		401			{object}	common.Response[any]
//	@Failure		404			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payrolls/{payroll_id}/send [post]
func (h *PayrollHandler) SendPayroll(c *gin.Context) {
	payrollID := c.Param("payroll_id")
	payroll, err := h.service.SendPayroll(c.Request.Context(), payrollID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollResponse(payroll)))
}
