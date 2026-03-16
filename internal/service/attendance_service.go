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

type AttendanceService interface {
	CheckIn(ctx context.Context, userID string, payload request.AttendanceRequest) (*model.Attendance, error)
	CheckOut(ctx context.Context, userID string, payload request.AttendanceRequest) (*model.Attendance, error)
	GetLogs(ctx context.Context, userID string) ([]model.Attendance, error)
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
			return nil, common.BadRequestError("File not found")
		}
		return nil, err
	}
	if file.Type != model.FileTypeCheckIn {
		return nil, common.BadRequestError("Invalid file type for check-in")
	}
	if file.UploadedBy != userID {
		return nil, common.ForbiddenError("File does not belong to current user")
	}
	fileURL := file.FileURL

	attendance, err := s.attendanceRepo.GetByUserAndDate(ctx, userID, today)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		attendance = &model.Attendance{
			UserID:         userID,
			Date:           today,
			CheckInAt:      &now,
			CheckInFileID:  &payload.FileID,
			CheckInFileURL: &fileURL,
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

	attendance.CheckInAt = &now
	attendance.CheckInFileID = &payload.FileID
	attendance.CheckInFileURL = &fileURL
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
	fileURL := file.FileURL

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

	attendance.CheckOutAt = &now
	attendance.CheckOutFileID = &payload.FileID
	attendance.CheckOutFileURL = &fileURL
	if err := s.attendanceRepo.Update(ctx, attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *attendanceService) GetLogs(ctx context.Context, userID string) ([]model.Attendance, error) {
	return s.attendanceRepo.GetLogsByUserID(ctx, userID)
}

func startOfDay(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta") // WIB
	local := t.In(loc)
	// Store date-only values in UTC so the SQL driver doesn't shift the day when converting time zones.
	y, m, d := local.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
