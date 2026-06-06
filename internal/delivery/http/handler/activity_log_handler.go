package handler

import (
	"net/http"
	"strconv"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/gin-gonic/gin"
)

type ActivityLogHandler struct {
	repo repository.ActivityLogRepository
}

func NewActivityLogHandler(repo repository.ActivityLogRepository) *ActivityLogHandler {
	return &ActivityLogHandler{repo: repo}
}

// GetAll godoc
//
//	@Summary		Daftar activity log (Admin)
//	@Description	Mendapatkan log aktivitas seluruh pengguna dengan pagination dan filter. Akses hanya untuk admin.
//	@Description
//	@Description	**Query Params (opsional):**
//	@Description	- `page` (integer, default 1): Nomor halaman
//	@Description	- `limit` (integer, default 20): Jumlah data per halaman
//	@Description	- `user_id` (UUID): Filter log berdasarkan pengguna tertentu
//	@Description	- `method` (string): Filter berdasarkan HTTP method — `POST`, `PUT`, `PATCH`, `DELETE`
//	@Description	- `status_code` (integer): Filter berdasarkan kode HTTP response (contoh: 200, 400, 500)
//	@Description	- `search` (string): Cari berdasarkan path atau deskripsi aksi
//	@Description
//	@Description	**Behavior:** Hanya mencatat request non-GET (perubahan data). GET tidak dilog.
//	@Description
//	@Description	**Response 200:** `data` (array log), `total` (total semua), `page`, `limit`
//	@Tags			Activity Logs
//	@Produce		json
//	@Param			page		query		int		false	"Nomor halaman (default: 1)"
//	@Param			limit		query		int		false	"Jumlah data per halaman (default: 20)"
//	@Param			user_id		query		string	false	"Filter berdasarkan user ID"
//	@Param			method		query		string	false	"Filter berdasarkan HTTP method (GET/POST/PUT/PATCH/DELETE)"
//	@Param			status_code	query		int		false	"Filter berdasarkan status code HTTP"
//	@Param			search		query		string	false	"Cari berdasarkan path atau deskripsi"
//	@Success		200			{object}	common.Response[response.ActivityLogListResponse]
//	@Failure		401			{object}	common.Response[any]
//	@Failure		403			{object}	common.Response[any]
//	@Security		BearerAuth
//	@Router			/activity-logs [get]
func (h *ActivityLogHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	statusCode, _ := strconv.Atoi(c.Query("status_code"))

	filter := repository.ActivityLogFilter{
		UserID:     c.Query("user_id"),
		Method:     c.Query("method"),
		StatusCode: statusCode,
		Search:     c.Query("search"),
		Page:       page,
		Limit:      limit,
	}

	logs, total, err := h.repo.GetAll(c.Request.Context(), filter)
	if err != nil {
		common.ErrorHandler(c, common.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(
		response.ToActivityLogListResponse(logs, total, page, limit),
	))
}
