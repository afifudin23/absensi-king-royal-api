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

// CheckIn godoc
//
//	@Summary		Check-in
//	@Description	Melakukan check-in absensi harian dengan menyertakan foto bukti kehadiran.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `file_id` (UUID, required): ID file foto check-in yang sudah diupload sebelumnya via `POST /files`
//	@Description
//	@Description	**Validasi & Behavior:**
//	@Description	- Hanya bisa check-in sekali per hari
//	@Description	- Jika sudah check-in hari ini, request akan ditolak
//	@Description	- Waktu check-in dicatat otomatis sesuai waktu server
//	@Description
//	@Description	**Response 200:** Data absensi hari ini (id, status, tanggal, jam check-in, URL foto)
//	@Tags			Attendance
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.AttendanceRequest	true	"ID file foto check-in"
//	@Success		200		{object}	common.Response[response.AttendanceResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance/check-in [post]
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

// CheckOut godoc
//
//	@Summary		Check-out
//	@Description	Melakukan check-out absensi harian dengan menyertakan foto bukti.
//	@Description
//	@Description	**Request Body:**
//	@Description	- `file_id` (UUID, required): ID file foto check-out yang sudah diupload via `POST /files`
//	@Description
//	@Description	**Validasi & Behavior:**
//	@Description	- Harus sudah check-in terlebih dahulu di hari yang sama
//	@Description	- Hanya bisa check-out sekali per hari
//	@Description	- Waktu check-out dicatat otomatis sesuai waktu server
//	@Description
//	@Description	**Response 200:** Data absensi hari ini (id, status, jam check-in, jam check-out, URL foto)
//	@Tags			Attendance
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.AttendanceRequest	true	"ID file foto check-out"
//	@Success		200		{object}	common.Response[response.AttendanceResponse]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance/check-out [post]
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

// GetLogs godoc
//
//	@Summary		Log absensi saya
//	@Description	Mendapatkan riwayat absensi milik pengguna yang sedang login.
//	@Description
//	@Description	**Query Params (opsional):**
//	@Description	- `start_date` (string, format YYYY-MM-DD): Filter dari tanggal ini
//	@Description	- `end_date` (string, format YYYY-MM-DD): Filter sampai tanggal ini
//	@Description
//	@Description	**Response 200:** Array data absensi (tanggal, status, jam masuk/keluar, foto, lembur, dll)
//	@Tags			Attendance
//	@Produce		json
//	@Param			start_date	query		string	false	"Tanggal mulai (YYYY-MM-DD)"
//	@Param			end_date	query		string	false	"Tanggal selesai (YYYY-MM-DD)"
//	@Success		200			{object}	common.Response[[]response.AttendanceResponse]
//	@Failure		400			{object}	common.Response[any]
//	@Failure		401			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance/logs [get]
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

// GetRecap godoc
//
//	@Summary		Rekap absensi (Admin)
//	@Description	Mendapatkan rekap absensi seluruh karyawan dalam satu bulan. Akses hanya untuk admin.
//	@Description
//	@Description	**Query Params (required):**
//	@Description	- `month` (integer, 1–12): Bulan yang ingin dilihat
//	@Description	- `year` (integer, min 2000): Tahun yang ingin dilihat
//	@Description
//	@Description	**Response 200:** Rekap absensi per karyawan (hadir, sakit, cuti, absen, lembur, dll)
//	@Description
//	@Description	**Akses:** Hanya admin
//	@Tags			Attendance
//	@Produce		json
//	@Param			month	query		int	true	"Bulan (1-12)"
//	@Param			year	query		int	true	"Tahun (contoh: 2024)"
//	@Success		200		{object}	common.Response[any]
//	@Failure		400		{object}	common.Response[any]
//	@Failure		401		{object}	common.Response[any]
//	@Failure		403		{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance/recap [get]
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

// Update godoc
//
//	@Summary		Update absensi (Admin)
//	@Description	Admin memperbarui data absensi karyawan tertentu. Akses hanya untuk admin.
//	@Description
//	@Description	**Path Param:**
//	@Description	- `attendance_id` (UUID): ID record absensi yang akan diupdate
//	@Description
//	@Description	**Request Body (semua opsional):**
//	@Description	- `status` (string): Status absensi — `present`, `off`, `sick`, `extra_off`, `absent`, `leave`
//	@Description	- `check_in_at` (string, RFC3339): Jam check-in manual
//	@Description	- `check_out_at` (string, RFC3339): Jam check-out manual
//	@Description	- `note` (string): Catatan admin
//	@Description	- `overtime_hours` (integer): Jumlah jam lembur
//	@Description	- `evidence_file_id` (UUID): ID file bukti
//	@Description
//	@Description	**Response 200:** ID absensi yang diupdate
//	@Tags			Attendance
//	@Accept			json
//	@Produce		json
//	@Param			attendance_id	path		string							true	"ID absensi"
//	@Param			body			body		request.AttendanceUpdateRequest	true	"Data absensi yang diperbarui"
//	@Success		200				{object}	common.Response[common.ActionSuccessResponse]
//	@Failure		400				{object}	common.Response[any]
//	@Failure		401				{object}	common.Response[any]
//	@Failure		403				{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/attendance/{attendance_id} [patch]
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
