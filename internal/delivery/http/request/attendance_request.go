package request

import (
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type AttendanceRequest struct {
	// FileID is the `files.id` that has been uploaded previously.
	FileID string `json:"file_id" binding:"required"`
}

func (rq *AttendanceRequest) Normalize() {
	rq.FileID = strings.TrimSpace(rq.FileID)
}

type AttendanceUpdateRequest struct {
	Status     *model.AttendanceStatus `json:"status" binding:"omitempty,oneof=present off sick extra_off absent leave"`
	CheckInAt  *string                 `json:"check_in_at"`
	CheckOutAt *string                 `json:"check_out_at"`
	Note       *string                 `json:"note"`
}

func (rq *AttendanceUpdateRequest) Normalize() {
	if rq.Status != nil {
		normalized := model.AttendanceStatus(strings.ToLower(strings.TrimSpace(string(*rq.Status))))
		rq.Status = &normalized
	}
	if rq.Note != nil {
		trimmed := strings.TrimSpace(*rq.Note)
		rq.Note = &trimmed
	}
}
