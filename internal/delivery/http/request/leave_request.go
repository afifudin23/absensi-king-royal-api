package request

import (
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type LeaveRequestCreateRequest struct {
	StartDate string                 `json:"start_date" binding:"required"`
	EndDate   string                 `json:"end_date" binding:"required"`
	Reason    string                 `json:"reason" binding:"required"`
	Type      model.LeaveRequestType `json:"type" binding:"required,oneof=sick extra_off overtime leave"`

	EvidenceFileID *string  `json:"evidence_file_id"`
	OvertimeHours  *float64 `json:"overtime_hours"`
}

func (rq *LeaveRequestCreateRequest) Normalize() {
	rq.Reason = strings.TrimSpace(rq.Reason)
	rq.Type = model.LeaveRequestType(strings.ToLower(strings.TrimSpace(string(rq.Type))))
}

type LeaveRequestUpdateRequest struct {
	StartDate      *string                 `json:"start_date"`
	EndDate        *string                 `json:"end_date"`
	Reason         *string                 `json:"reason"`
	Type           *model.LeaveRequestType `json:"type" binding:"omitempty,oneof=sick extra_off overtime leave"`
	EvidenceFileID *string                 `json:"evidence_file_id"`
	OvertimeHours  *float64                `json:"overtime_hours"`
}

func (rq *LeaveRequestUpdateRequest) Normalize() {
	normalizeOptionalString(&rq.Reason, false)
	if rq.Type != nil {
		normalized := model.LeaveRequestType(strings.ToLower(strings.TrimSpace(string(*rq.Type))))
		rq.Type = &normalized
	}
}


type LeaveRequestUpdateStatusRequest struct {
	Status model.LeaveRequestStatus `json:"status" binding:"required,oneof=approved rejected"`
}