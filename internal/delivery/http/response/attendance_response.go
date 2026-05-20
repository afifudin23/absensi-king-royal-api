package response

import (
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type AttendanceResponse struct {
	ID              string                 `json:"id"`
	UserID          string                 `json:"user_id"`
	Status          model.AttendanceStatus `json:"status"`
	Date            string                 `json:"date"`
	CheckInAt       *string                `json:"check_in_at"`
	CheckOutAt      *string                `json:"check_out_at"`
	CheckInFileID   *string                `json:"check_in_file_id"`
	CheckInFileURL  *string                `json:"check_in_file_url"`
	CheckOutFileID  *string                `json:"check_out_file_id"`
	CheckOutFileURL *string                `json:"check_out_file_url"`
	Note            *string                `json:"note"`
	Source          model.AttendanceSource `json:"source"`
	UpdatedBy       *string                `json:"updated_by"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

func ToAttendanceResponse(attendance model.Attendance) AttendanceResponse {
	return AttendanceResponse{
		ID:              attendance.ID,
		UserID:          attendance.UserID,
		Status:          attendance.Status,
		Date:            attendance.Date.Format("2006-01-02"),
		CheckInAt:       toTimeStringPtr(attendance.CheckInAt),
		CheckOutAt:      toTimeStringPtr(attendance.CheckOutAt),
		CheckInFileID:   attendance.CheckInFileID,
		CheckInFileURL:  toFileURLPtr(attendance.CheckInFile),
		CheckOutFileID:  attendance.CheckOutFileID,
		CheckOutFileURL: toFileURLPtr(attendance.CheckOutFile),
		Note:            attendance.Note,
		Source:          attendance.Source,
		UpdatedBy:       attendance.UpdatedBy,
		CreatedAt:       attendance.CreatedAt,
		UpdatedAt:       attendance.UpdatedAt,
	}
}

func ToAttendanceListResponse(attendances []model.Attendance) []AttendanceResponse {
	response := make([]AttendanceResponse, 0, len(attendances))
	for _, attendance := range attendances {
		response = append(response, ToAttendanceResponse(attendance))
	}
	return response
}

func toTimeStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	v := t.Format(time.RFC3339)
	return &v
}

func toFileURLPtr(file *model.File) *string {
	if file == nil {
		return nil
	}
	if file.FileURL == "" {
		return nil
	}
	v := file.FileURL
	return &v
}
