package handler

import (
	"net/http"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
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

// Create godoc
//
//	@Summary		Buat pengajuan absensi
//	@Description	Membuat pengajuan kehadiran baru seperti cuti, sakit, lembur, atau libur tambahan.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `start_date` (string, required, format YYYY-MM-DD): Tanggal mulai pengajuan
//	@Description	- `end_date` (string, required, format YYYY-MM-DD): Tanggal selesai pengajuan
//	@Description	- `reason` (string, required): Alasan pengajuan
//	@Description	- `type` (string, required): Tipe — `sick`, `leave`, `extra_off`, `overtime`
//	@Description	- `evidence_file_id` (UUID, opsional): ID file bukti (contoh: surat dokter untuk sakit)
//	@Description	- `requested_overtime_hours` (integer, opsional): Jumlah jam lembur (khusus tipe `overtime`)
//	@Description
//	@Description	**Response 201:** ID pengajuan yang baru dibuat
//	@Tags			Attendance Requests
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.AttendanceRequestCreateRequest	true	"Data pengajuan"
//	@Success		201		{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance-requests [post]
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

// GetAll godoc
//
//	@Summary		Semua pengajuan absensi (Admin)
//	@Description	Mendapatkan daftar semua pengajuan absensi dari seluruh karyawan. Akses hanya untuk admin.
//	@Description
//	@Description	**Query Params (opsional):**
//	@Description	- `status` (string): Filter status — `pending`, `approved`, `rejected`
//	@Description	- `type` (string): Filter tipe — `sick`, `leave`, `extra_off`, `overtime`
//	@Description	- `start_date` (string, YYYY-MM-DD): Filter dari tanggal ini
//	@Description	- `end_date` (string, YYYY-MM-DD): Filter sampai tanggal ini
//	@Description
//	@Description	**Response 200:** Array pengajuan beserta nama karyawan, status review, dan bukti
//	@Tags			Attendance Requests
//	@Produce		json
//	@Param			status		query		string	false	"Filter status (pending/approved/rejected)"
//	@Param			type		query		string	false	"Filter tipe (sick/leave/extra_off/overtime)"
//	@Param			start_date	query		string	false	"Tanggal mulai (YYYY-MM-DD)"
//	@Param			end_date	query		string	false	"Tanggal selesai (YYYY-MM-DD)"
//	@Success		200			{object}	common.Response[[]response.AttendanceRequestResponse]
//	@Failure		401			{object}	common.Response[any]
//	@Failure		403			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance-requests [get]
func (h *AttendanceRequestHandler) GetAll(c *gin.Context) {
	filter := &repository.AttendanceRequestFilter{
		Status: c.Query("status"),
		Type:   c.Query("type"),
	}
	if s := c.Query("start_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			common.ErrorHandler(c, common.BadRequestError("Invalid start_date format, use YYYY-MM-DD"))
			return
		}
		filter.StartDate = &t
	}
	if e := c.Query("end_date"); e != "" {
		t, err := time.Parse("2006-01-02", e)
		if err != nil {
			common.ErrorHandler(c, common.BadRequestError("Invalid end_date format, use YYYY-MM-DD"))
			return
		}
		filter.EndDate = &t
	}

	items, err := h.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceRequestListResponse(items)))
}

// GetByID godoc
//
//	@Summary		Detail pengajuan absensi
//	@Description	Mendapatkan detail satu pengajuan absensi berdasarkan ID.
//	@Description
//	@Description	**Path Param:**
//	@Description	- `attendance_request_id` (UUID): ID pengajuan yang ingin dilihat
//	@Description
//	@Description	**Response 200:** Data lengkap pengajuan (status, tipe, tanggal, alasan, bukti, reviewer, catatan review)
//	@Tags			Attendance Requests
//	@Produce		json
//	@Param			attendance_request_id	path		string	true	"ID pengajuan"
//	@Success		200						{object}	common.Response[response.AttendanceRequestResponse]
//	@Failure		401						{object}	common.Response[any]
//	@Failure		404						{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance-requests/{attendance_request_id} [get]
func (h *AttendanceRequestHandler) GetByID(c *gin.Context) {
	attendanceRequestID := c.Param("attendance_request_id")
	item, err := h.service.GetByID(c.Request.Context(), attendanceRequestID)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, common.SuccessResponse(response.ToAttendanceRequestResponse(*item)))
}

// GetByUserID godoc
//
//	@Summary		Pengajuan absensi saya
//	@Description	Mendapatkan daftar semua pengajuan absensi milik pengguna yang sedang login.
//	@Description
//	@Description	**Response 200:** Array pengajuan milik pengguna (status, tipe, tanggal, alasan, bukti, info review)
//	@Tags			Attendance Requests
//	@Produce		json
//	@Success		200	{object}	common.Response[[]response.AttendanceRequestResponse]
//	@Failure		401	{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance-requests/me [get]
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

// Update godoc
//
//	@Summary		Update pengajuan absensi
//	@Description	Memperbarui pengajuan absensi yang masih berstatus pending. Pengajuan yang sudah diapprove/reject tidak bisa diubah.
//	@Description
//	@Description	**Request Body (semua opsional):**
//	@Description	- `start_date` (string, YYYY-MM-DD): Tanggal mulai baru
//	@Description	- `end_date` (string, YYYY-MM-DD): Tanggal selesai baru
//	@Description	- `reason` (string): Alasan baru
//	@Description	- `type` (string): Tipe baru — `sick`, `leave`, `extra_off`, `overtime`
//	@Description	- `evidence_file_id` (UUID): File bukti baru
//	@Description	- `requested_overtime_hours` (integer): Jam lembur baru
//	@Description
//	@Description	**Response 200:** ID pengajuan yang diupdate
//	@Tags			Attendance Requests
//	@Accept			json
//	@Produce		json
//	@Param			attendance_request_id	path		string									true	"ID pengajuan"
//	@Param			body					body		request.AttendanceRequestUpdateRequest	true	"Data yang diperbarui"
//	@Success		200						{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		400						{object}	common.Response[any]
//	@Failure		401						{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance-requests/{attendance_request_id} [put]
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

// UpdateStatus godoc
//
//	@Summary		Setujui/tolak pengajuan (Admin)
//	@Description	Admin menyetujui atau menolak pengajuan absensi karyawan. Akses hanya untuk admin.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `status` (string, required): Status baru — `approved` atau `rejected`
//	@Description
//	@Description	**Behavior:**
//	@Description	- Jika diapprove, data absensi karyawan pada tanggal yang bersangkutan akan diperbarui otomatis
//	@Description	- Reviewer dan waktu review dicatat otomatis
//	@Description
//	@Description	**Response 200:** ID pengajuan yang diupdate statusnya
//	@Tags			Attendance Requests
//	@Accept			json
//	@Produce		json
//	@Param			attendance_request_id	path		string											true	"ID pengajuan"
//	@Param			body					body		request.AttendanceRequestUpdateStatusRequest	true	"Status baru (approved/rejected)"
//	@Success		200						{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		400						{object}	common.Response[any]
//	@Failure		401						{object}	common.Response[any]
//	@Failure		403						{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance-requests/{attendance_request_id}/status [patch]
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

// Delete godoc
//
//	@Summary		Hapus pengajuan absensi
//	@Description	Menghapus satu atau beberapa pengajuan absensi sekaligus berdasarkan daftar ID.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `ids` (array of UUID, required, min 1 item): Daftar ID pengajuan yang akan dihapus
//	@Description
//	@Description	**Response 200:** Jumlah data yang berhasil dihapus
//	@Tags			Attendance Requests
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.AttendanceRequestBulkDeleteRequest	true	"Daftar ID pengajuan yang akan dihapus"
//	@Success		200		{object}	common.Response[common.DeleteSuccessResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance-requests [delete]
func (h *AttendanceRequestHandler) Delete(c *gin.Context) {
	var payload request.AttendanceRequestBulkDeleteRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		common.ErrorHandler(c, common.ValidationError(err))
		return
	}
	err := h.service.Delete(c.Request.Context(), payload.IDs)
	if err != nil {
		common.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(common.ToDeleteSuccessResponse(len(payload.IDs))))
}
