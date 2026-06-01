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
