package response

import (
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type AttendanceRequestResponse struct {
	ID                     string                        `json:"id"`
	UserID                 string                        `json:"user_id"`
	EmployeeName           *string                       `json:"employee_name"`
	StartDate              string                        `json:"start_date"`
	EndDate                string                        `json:"end_date"`
	Reason                 string                        `json:"reason"`
	Type                   model.AttendanceRequestType   `json:"type"`
	Status                 model.AttendanceRequestStatus `json:"status"`
	EvidenceFileID         *string                       `json:"evidence_file_id"`
	EvidenceFileURL        *string                       `json:"evidence_file_url"`
	RequestedOvertimeHours *int                          `json:"requested_overtime_hours"`
	ReviewedBy *string    `json:"reviewed_by"`
	ReviewedAt *time.Time `json:"reviewed_at"`
	ReviewNote *string    `json:"review_note"`
	CreatedAt              time.Time                     `json:"created_at"`
	UpdatedAt              time.Time                     `json:"updated_at"`
}

func ToAttendanceRequestResponse(data model.AttendanceRequest) AttendanceRequestResponse {
	var employeeName *string
	if data.User != nil {
		employeeName = &data.User.FullName
	}
	return AttendanceRequestResponse{
		ID:                     data.ID,
		UserID:                 data.UserID,
		EmployeeName:           employeeName,
		StartDate:              data.StartDate.Format("2006-01-02"),
		EndDate:                data.EndDate.Format("2006-01-02"),
		Reason:                 data.Reason,
		Type:                   data.Type,
		Status:                 data.Status,
		EvidenceFileID:         data.EvidenceFileID,
		EvidenceFileURL:        toFileURLPtr(data.EvidenceFile),
		RequestedOvertimeHours: data.RequestedOvertimeHours,
		ReviewedBy:             data.ReviewedBy,
		ReviewedAt:             data.ReviewedAt,
		ReviewNote:             data.ReviewNote,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
	}
}

func ToAttendanceRequestListResponse(items []model.AttendanceRequest) []AttendanceRequestResponse {
	response := make([]AttendanceRequestResponse, 0, len(items))
	for _, item := range items {
		response = append(response, ToAttendanceRequestResponse(item))
	}
	return response
}
