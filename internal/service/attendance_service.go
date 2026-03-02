package service

import (
	"errors"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"gorm.io/gorm"
)

type AttendanceService interface {
	CheckIn(userID string) (*model.Attendance, error)
	CheckOut(userID string) (*model.Attendance, error)
	GetLogs(userID string) ([]model.Attendance, error)
}

type attendanceService struct {
	attendanceRepo repository.AttendanceRepository
}

func NewAttendanceService() AttendanceService {
	return &attendanceService{attendanceRepo: repository.NewAttendanceRepository()}
}

func (s *attendanceService) CheckIn(userID string) (*model.Attendance, error) {
	now := time.Now()
	today := startOfDay(now)

	attendance, err := s.attendanceRepo.GetByUserAndDate(userID, today)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		attendance = &model.Attendance{
			UserID:    userID,
			Date:      today,
			CheckInAt: &now,
		}
		if err := s.attendanceRepo.Create(attendance); err != nil {
			return nil, err
		}
		return attendance, nil
	}

	if attendance.CheckInAt != nil {
		return nil, common.BadRequestError("You have already checked in today")
	}

	attendance.CheckInAt = &now
	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *attendanceService) CheckOut(userID string) (*model.Attendance, error) {
	now := time.Now()
	today := startOfDay(now)

	attendance, err := s.attendanceRepo.GetByUserAndDate(userID, today)
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
	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *attendanceService) GetLogs(userID string) ([]model.Attendance, error) {
	return s.attendanceRepo.GetLogsByUserID(userID)
}

func startOfDay(t time.Time) time.Time {
	local := t.In(time.Local)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, local.Location())
}
