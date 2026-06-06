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

// GetAll godoc
//
//	@Summary		Daftar pengaturan payroll
//	@Description	Mendapatkan semua komponen konfigurasi penggajian yang tersimpan.
//	@Description
//	@Description	**Response 200:** Array konfigurasi (id, config_name, config_key, value, tanggal dibuat/diupdate)
//	@Description
//	@Description	**Contoh komponen:** Tarif lembur, potongan absensi, pajak penghasilan, tunjangan default, dll
//	@Tags			Payroll Settings
//	@Produce		json
//	@Success		200	{object}	common.Response[[]response.PayrollSettingResponse]
//	@Failure		401	{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payroll-settings [get]
func (h *PayrollSettingHandler) GetAll(c *gin.Context) {
	payrollSettings, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToPayrollSettingListResponse(payrollSettings)))

}

// Create godoc
//
//	@Summary		Buat pengaturan payroll
//	@Description	Menambahkan komponen konfigurasi penggajian baru.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `config_name` (string, required, min 3 karakter): Nama komponen (contoh: "Tarif Lembur Per Jam")
//	@Description	- `value` (number, required): Nilai komponen (contoh: 50000)
//	@Description
//	@Description	**Behavior:**
//	@Description	- `config_key` dibuat otomatis dari `config_name` (snake_case)
//	@Description	- Nama komponen harus unik, jika sudah ada akan mengembalikan 400
//	@Description
//	@Description	**Response 201:** Data komponen yang baru dibuat (id, config_name, config_key, value)
//	@Tags			Payroll Settings
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.PayrollSettingRequest	true	"Data komponen penggajian"
//	@Success		201		{object}	common.Response[response.PayrollSettingResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payroll-settings [post]
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

// Update godoc
//
//	@Summary		Update pengaturan payroll
//	@Description	Memperbarui satu komponen konfigurasi penggajian berdasarkan ID.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `config_name` (string, required): Nama komponen baru
//	@Description	- `value` (number, required): Nilai komponen baru
//	@Description
//	@Description	**Response 201:** Data komponen setelah diupdate
//	@Tags			Payroll Settings
//	@Accept			json
//	@Produce		json
//	@Param			payroll_id	path		string							true	"ID pengaturan payroll"
//	@Param			body		body		request.PayrollSettingRequest	true	"Data yang diperbarui"
//	@Success		201			{object}	common.Response[response.PayrollSettingResponse]
//	@Failure		400			{object}	common.Response[any]
//	@Failure		401			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payroll-settings/{payroll_id} [patch]
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

// UpdateBulk godoc
//
//	@Summary		Bulk update pengaturan payroll
//	@Description	Memperbarui beberapa komponen konfigurasi penggajian sekaligus dalam satu request.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `settings` (array, required, min 1): Daftar komponen yang akan diupdate
//	@Description	  - `config_key` (string, required): Key komponen yang akan diupdate (harus sudah ada)
//	@Description	  - `value` (number, required): Nilai baru
//	@Description
//	@Description	**Response 201:** Array seluruh komponen setelah diupdate
//	@Tags			Payroll Settings
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.PayrollSettingUpdateBulkRequest	true	"Daftar komponen yang diperbarui"
//	@Success		201		{object}	common.Response[[]response.PayrollSettingResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payroll-settings/bulk [put]
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

// Delete godoc
//
//	@Summary		Hapus pengaturan payroll
//	@Description	Menghapus satu atau beberapa komponen konfigurasi penggajian berdasarkan daftar ID.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `ids` (array of UUID, required, min 1): Daftar ID komponen yang akan dihapus
//	@Description
//	@Description	**Response 200:** `total` (jumlah yang dihapus), `deleted_count`, `skipped_count`
//	@Tags			Payroll Settings
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.PayrollSettingIdsRequest	true	"Daftar ID yang akan dihapus"
//	@Success		200		{object}	common.Response[response.PayrollSettingDeleteResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/payroll-settings [delete]
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
