package response

import "github.com/afifudin23/absensi-king-royal-api/internal/model"

type LeaveResponse struct {
	ID              string                   `json:"id"`
	UserID          string                   `json:"user_id"`
	StartDate       string                   `json:"start_date"`
	EndDate         string                   `json:"end_date"`
	Reason          string                   `json:"reason"`
	Type            model.LeaveRequestType   `json:"type"`
	Status          model.LeaveRequestStatus `json:"status"`
	EvidenceFileID  *string                  `json:"evidence_file_id"`
	EvidenceFileURL *string                  `json:"evidence_file_url"`
	OvertimeHours   *float64                 `json:"overtime_hours"`
	CreatedAt       string                   `json:"created_at"`
	UpdatedAt       string                   `json:"updated_at"`
}

func ToLeaveResponse(leave model.LeaveRequest) LeaveResponse {
	return LeaveResponse{
		ID:              leave.ID,
		UserID:          leave.UserID,
		StartDate:       leave.StartDate.Format("2006-01-02"),
		EndDate:         leave.EndDate.Format("2006-01-02"),
		Reason:          leave.Reason,
		Type:            leave.Type,
		Status:          leave.Status,
		EvidenceFileID:  leave.EvidenceFileID,
		EvidenceFileURL: leave.EvidenceFileURL,
		OvertimeHours:   leave.OvertimeHours,
		CreatedAt:       leave.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       leave.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToLeaveListResponse(leaves []model.LeaveRequest) []LeaveResponse {
	response := make([]LeaveResponse, 0, len(leaves))
	for _, leave := range leaves {
		response = append(response, ToLeaveResponse(leave))
	}
	return response
}
