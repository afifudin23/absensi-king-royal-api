package request

import "strings"

type AttendanceRequest struct {
	// FileID is the `files.id` that has been uploaded previously.
	FileID string `json:"file_id" binding:"required"`
}

func (rq *AttendanceRequest) Normalize() {
	rq.FileID = strings.TrimSpace(rq.FileID)
}

type AttendanceUpdateRequest struct {
	CheckInAt  *string `json:"check_in_at"`
	CheckOutAt *string `json:"check_out_at"`
}
