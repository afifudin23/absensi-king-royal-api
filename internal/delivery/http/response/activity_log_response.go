package response

import (
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type ActivityLogResponse struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"user_id"`
	UserName    *string   `json:"user_name"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	StatusCode  int       `json:"status_code"`
	IPAddress   string    `json:"ip_address"`
	LatencyMs   float64   `json:"latency_ms"`
	CreatedAt   time.Time `json:"created_at"`
}

type ActivityLogListResponse struct {
	Data  []ActivityLogResponse `json:"data"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}

func ToActivityLogResponse(a model.ActivityLog) ActivityLogResponse {
	return ActivityLogResponse{
		ID:          a.ID,
		UserID:      a.UserID,
		UserName:    a.UserName,
		Method:      a.Method,
		Path:        a.Path,
		Description: a.Description,
		StatusCode:  a.StatusCode,
		IPAddress:   a.IPAddress,
		LatencyMs:   a.LatencyMs,
		CreatedAt:   a.CreatedAt,
	}
}

func ToActivityLogListResponse(logs []model.ActivityLog, total int64, page, limit int) ActivityLogListResponse {
	items := make([]ActivityLogResponse, 0, len(logs))
	for _, l := range logs {
		items = append(items, ToActivityLogResponse(l))
	}
	return ActivityLogListResponse{Data: items, Total: total, Page: page, Limit: limit}
}
