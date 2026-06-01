package request

import (
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type AttendanceRequestCreateRequest struct {
	StartDate              string                      `json:"start_date" binding:"required"`
	EndDate                string                      `json:"end_date" binding:"required"`
	Reason                 string                      `json:"reason" binding:"required"`
	Type                   model.AttendanceRequestType `json:"type" binding:"required,oneof=sick leave extra_off overtime"`
	EvidenceFileID         *string                     `json:"evidence_file_id"`
	RequestedOvertimeHours *int                        `json:"requested_overtime_hours"`
}

func (rq *AttendanceRequestCreateRequest) Normalize() {
	rq.Reason = strings.TrimSpace(rq.Reason)
	rq.Type = model.AttendanceRequestType(strings.ToLower(strings.TrimSpace(string(rq.Type))))
	normalizeOptionalString(&rq.EvidenceFileID, false)
}

type AttendanceRequestUpdateRequest struct {
	StartDate              *string                      `json:"start_date"`
	EndDate                *string                      `json:"end_date"`
	Reason                 *string                      `json:"reason"`
	Type                   *model.AttendanceRequestType `json:"type" binding:"omitempty,oneof=sick leave extra_off overtime"`
	EvidenceFileID         *string                      `json:"evidence_file_id"`
	RequestedOvertimeHours *int                         `json:"requested_overtime_hours"`
}

func (rq *AttendanceRequestUpdateRequest) Normalize() {
	normalizeOptionalString(&rq.Reason, false)
	normalizeOptionalString(&rq.EvidenceFileID, false)
	if rq.Type != nil {
		normalized := model.AttendanceRequestType(strings.ToLower(strings.TrimSpace(string(*rq.Type))))
		rq.Type = &normalized
	}
}

type AttendanceRequestUpdateStatusRequest struct {
	Status model.AttendanceRequestStatus `json:"status" binding:"required,oneof=approved rejected"`
}

type AttendanceRequestBulkDeleteRequest struct {
    IDs []string `json:"ids" binding:"required,min=1"`
}
