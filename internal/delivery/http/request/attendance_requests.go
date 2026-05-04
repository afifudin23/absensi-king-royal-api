package request

import (
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type AttendanceRequestCreateRequest struct {
	StartDate                string                      `json:"start_date" binding:"required"`
	EndDate                  string                      `json:"end_date" binding:"required"`
	Reason                   string                      `json:"reason" binding:"required"`
	Type                     model.AttendanceRequestType `json:"type" binding:"required,oneof=sick leave extra_off overtime correction"`
	AttendanceID             *string                     `json:"attendance_id"`
	EvidenceFileID           *string                     `json:"evidence_file_id"`
	RequestedCheckInAt       *string                     `json:"requested_check_in_at"`
	RequestedCheckOutAt      *string                     `json:"requested_check_out_at"`
	RequestedOvertimeMinutes *int                        `json:"requested_overtime_minutes"`
}

func (rq *AttendanceRequestCreateRequest) Normalize() {
	rq.Reason = strings.TrimSpace(rq.Reason)
	rq.Type = model.AttendanceRequestType(strings.ToLower(strings.TrimSpace(string(rq.Type))))
	normalizeOptionalString(&rq.AttendanceID, false)
	normalizeOptionalString(&rq.EvidenceFileID, false)
	normalizeOptionalString(&rq.RequestedCheckInAt, false)
	normalizeOptionalString(&rq.RequestedCheckOutAt, false)
}

type AttendanceRequestUpdateRequest struct {
	StartDate                *string                      `json:"start_date"`
	EndDate                  *string                      `json:"end_date"`
	Reason                   *string                      `json:"reason"`
	Type                     *model.AttendanceRequestType `json:"type" binding:"omitempty,oneof=sick leave extra_off overtime correction"`
	AttendanceID             *string                      `json:"attendance_id"`
	EvidenceFileID           *string                      `json:"evidence_file_id"`
	RequestedCheckInAt       *string                      `json:"requested_check_in_at"`
	RequestedCheckOutAt      *string                      `json:"requested_check_out_at"`
	RequestedOvertimeMinutes *int                         `json:"requested_overtime_minutes"`
}

func (rq *AttendanceRequestUpdateRequest) Normalize() {
	normalizeOptionalString(&rq.Reason, false)
	normalizeOptionalString(&rq.AttendanceID, false)
	normalizeOptionalString(&rq.EvidenceFileID, false)
	normalizeOptionalString(&rq.RequestedCheckInAt, false)
	normalizeOptionalString(&rq.RequestedCheckOutAt, false)
	if rq.Type != nil {
		normalized := model.AttendanceRequestType(strings.ToLower(strings.TrimSpace(string(*rq.Type))))
		rq.Type = &normalized
	}
}

type AttendanceRequestUpdateStatusRequest struct {
	Status model.AttendanceRequestStatus `json:"status" binding:"required,oneof=approved rejected cancelled"`
}
