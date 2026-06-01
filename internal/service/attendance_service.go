package service

import (
	"context"
	"errors"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"gorm.io/gorm"
)

type AttendanceRecapDaily struct {
	AttendanceID    string                 `json:"attendance_id"`
	Date            string                 `json:"date"`
	Status          model.AttendanceStatus `json:"status"`
	CheckInAt       *string                `json:"check_in_at"`
	CheckOutAt      *string                `json:"check_out_at"`
	OvertimeHours   *int                   `json:"overtime_hours"`
	Note            *string                `json:"note"`
	CheckInFileID   *string                `json:"check_in_file_id"`
	CheckOutFileID  *string                `json:"check_out_file_id"`
	Source          model.AttendanceSource `json:"source"`
	CheckInFileURL  *string                `json:"check_in_file_url"`
	CheckOutFileURL *string                `json:"check_out_file_url"`
	EvidenceFileURL *string                `json:"evidence_file_url"`
}

type AttendanceRecapData struct {
	UserID             string                 `json:"user_id"`
	EmployeeName       string                 `json:"employee_name"`
	Month              int                    `json:"month"`
	Year               int                    `json:"year"`
	TotalPresent       int                    `json:"total_present"`
	TotalOff           int                    `json:"total_off"`
	TotalSick          int                    `json:"total_sick"`
	TotalExtraOff      int                    `json:"total_extra_off"`
	TotalAbsent        int                    `json:"total_absent"`
	TotalLeave         int                    `json:"total_leave"`
	TotalOvertimeHours int                    `json:"total_overtime_hours"`
	DailyDetails       []AttendanceRecapDaily `json:"daily_details"`
}

type AttendanceService interface {
	CheckIn(ctx context.Context, userID string, payload request.AttendanceRequest) (*model.Attendance, error)
	CheckOut(ctx context.Context, userID string, payload request.AttendanceRequest) (*model.Attendance, error)
	GetLogs(ctx context.Context, userID string, startDate, endDate *time.Time) ([]model.Attendance, error)
	Update(ctx context.Context, updaterID string, id string, payload request.AttendanceUpdateRequest) (*model.Attendance, error)
	GetRecap(ctx context.Context, month, year int) ([]AttendanceRecapData, error)
}

type attendanceService struct {
	attendanceRepo repository.AttendanceRepository
	fileRepo       repository.FileRepository
}

func NewAttendanceService(attendanceRepo repository.AttendanceRepository, fileRepo repository.FileRepository) AttendanceService {
	return &attendanceService{attendanceRepo: attendanceRepo, fileRepo: fileRepo}
}

func (s *attendanceService) CheckIn(ctx context.Context, userID string, payload request.AttendanceRequest) (*model.Attendance, error) {
	now := time.Now()
	today := startOfDay(now)

	// Validate file_id up-front to avoid foreign key 1452 (and return a proper 4xx).
	file, err := s.fileRepo.GetByID(ctx, payload.FileID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.BadRequestError("Invalid file_id")
		}
		return nil, err
	}
	if file.Type != model.FileTypeCheckIn {
		return nil, common.BadRequestError("Invalid file type for check-in")
	}
	if file.UploadedBy != userID {
		return nil, common.ForbiddenError("File does not belong to current user")
	}

	attendance, err := s.attendanceRepo.GetByUserAndDate(ctx, userID, today)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		attendance = &model.Attendance{
			UserID:        userID,
			Status:        model.AttendanceStatusPresent,
			Date:          today,
			CheckInAt:     &now,
			CheckInFileID: &payload.FileID,
			Source:        model.AttendanceSourceSelfService,
		}
		if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
			// Race protection: if a duplicate row is inserted concurrently, re-fetch it.
			if isDuplicateError(err) {
				attendance, err = s.attendanceRepo.GetByUserAndDate(ctx, userID, today)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	attendance.Status = model.AttendanceStatusPresent
	attendance.CheckInAt = &now
	attendance.CheckInFileID = &payload.FileID
	attendance.Source = model.AttendanceSourceSelfService
	if err := s.attendanceRepo.Update(ctx, attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *attendanceService) CheckOut(ctx context.Context, userID string, payload request.AttendanceRequest) (*model.Attendance, error) {
	now := time.Now()
	today := startOfDay(now)

	// Validate file_id up-front to avoid foreign key 1452 (and return a proper 4xx).
	file, err := s.fileRepo.GetByID(ctx, payload.FileID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.BadRequestError("Invalid file_id")
		}
		return nil, err
	}
	if file.Type != model.FileTypeCheckOut {
		return nil, common.BadRequestError("Invalid file type for check-out")
	}
	if file.UploadedBy != userID {
		return nil, common.ForbiddenError("File does not belong to current user")
	}

	attendance, err := s.attendanceRepo.GetByUserAndDate(ctx, userID, today)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.BadRequestError("You must check in before checking out")
		}
		return nil, err
	}

	if attendance.CheckInAt == nil {
		return nil, common.BadRequestError("You must check in before checking out")
	}

	if attendance.CheckOutAt != nil {
		return nil, common.BadRequestError("You have already checked out today")
	}

	if !now.After(*attendance.CheckInAt) {
		return nil, common.BadRequestError("Check-out time must be after check-in time")
	}

	attendance.Status = model.AttendanceStatusPresent
	attendance.CheckOutAt = &now
	attendance.CheckOutFileID = &payload.FileID
	attendance.CheckOutFile = file
	attendance.Source = model.AttendanceSourceSelfService
	if err := s.attendanceRepo.Update(ctx, attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *attendanceService) GetLogs(ctx context.Context, userID string, startDate, endDate *time.Time) ([]model.Attendance, error) {
	return s.attendanceRepo.GetLogsByUserID(ctx, userID, startDate, endDate)
}

func (s *attendanceService) GetRecap(ctx context.Context, month, year int) ([]AttendanceRecapData, error) {
	rows, err := s.attendanceRepo.GetByMonthWithUser(ctx, month, year)
	if err != nil {
		return nil, err
	}

	grouped := map[string]*AttendanceRecapData{}
	order := []string{}

	for _, row := range rows {
		d, exists := grouped[row.UserID]
		if !exists {
			d = &AttendanceRecapData{
				UserID:       row.UserID,
				EmployeeName: row.EmployeeName,
				Month:        month,
				Year:         year,
			}
			grouped[row.UserID] = d
			order = append(order, row.UserID)
		}

		switch row.Status {
		case model.AttendanceStatusPresent:
			d.TotalPresent++
		case model.AttendanceStatusOff:
			d.TotalOff++
		case model.AttendanceStatusSick:
			d.TotalSick++
		case model.AttendanceStatusExtraOff:
			d.TotalExtraOff++
		case model.AttendanceStatusAbsent:
			d.TotalAbsent++
		case model.AttendanceStatusLeave:
			d.TotalLeave++
		}

		if row.OvertimeHours != nil {
			d.TotalOvertimeHours += *row.OvertimeHours
		}

		daily := AttendanceRecapDaily{
			AttendanceID:    row.ID,
			Date:            row.Date.Format("2006-01-02"),
			Status:          row.Status,
			OvertimeHours:   row.OvertimeHours,
			Note:            row.Note,
			CheckInFileID:   row.CheckInFileID,
			CheckOutFileID:  row.CheckOutFileID,
			Source:          row.Source,
			CheckInFileURL:  row.CheckInFileURL,
			CheckOutFileURL: row.CheckOutFileURL,
			EvidenceFileURL: row.EvidenceFileURL,
		}
		if row.CheckInAt != nil {
			v := row.CheckInAt.Format("15:04")
			daily.CheckInAt = &v
		}
		if row.CheckOutAt != nil {
			v := row.CheckOutAt.Format("15:04")
			daily.CheckOutAt = &v
		}
		d.DailyDetails = append(d.DailyDetails, daily)
	}

	result := make([]AttendanceRecapData, 0, len(order))
	for _, userID := range order {
		result = append(result, *grouped[userID])
	}
	return result, nil
}

func (s *attendanceService) Update(ctx context.Context, updaterID string, id string, payload request.AttendanceUpdateRequest) (*model.Attendance, error) {
	existing, err := s.attendanceRepo.GetByID(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.NotFoundError("Attendance not found")
		}
		return nil, err
	}

	if payload.CheckInAt != nil && *payload.CheckInAt != "" {
		checkInTime, err := CombineDateAndHHMM(existing.Date, *payload.CheckInAt)
		if err != nil {
			return nil, err
		}
		existing.CheckInAt = checkInTime
	}

	if payload.Status != nil {
		existing.Status = *payload.Status
	}

	if payload.CheckOutAt != nil && *payload.CheckOutAt != "" {
		checkOut, err := CombineDateAndHHMM(existing.Date, *payload.CheckOutAt)
		if err != nil {
			return nil, err
		}
		existing.CheckOutAt = checkOut
	}
	if payload.Note != nil {
		existing.Note = payload.Note
	}
	if payload.OvertimeHours != nil {
		existing.OvertimeHours = payload.OvertimeHours
	}
	if payload.EvidenceFileID != nil {
		// Validasi file harus ada dan milik sistem (bukan user tertentu)
		_, err := s.fileRepo.GetByID(ctx, *payload.EvidenceFileID)
		if err != nil {
			if isNotFoundError(err) {
				return nil, common.BadRequestError("Invalid evidence_file_id")
			}
			return nil, err
		}
		existing.EvidenceFileID = payload.EvidenceFileID
	}
	existing.Source = model.AttendanceSourceAdminEdit
	existing.UpdatedBy = &updaterID
	if err := s.attendanceRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func CombineDateAndHHMM(date time.Time, hhmm string) (*time.Time, error) {
	parsed, err := time.Parse("15:04", hhmm)
	if err != nil {
		return nil, err
	}

	result := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		parsed.Hour(),
		parsed.Minute(),
		0,
		0,
		time.Local,
	)

	return &result, nil
}

func startOfDay(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta") // WIB
	local := t.In(loc)
	// Store date-only values in UTC so the SQL driver doesn't shift the day when converting time zones.
	y, m, d := local.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
