package repository

import (
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	GetByUserAndDate(userID string, date time.Time) (*model.Attendance, error)
	Create(attendance *model.Attendance) error
	Update(attendance *model.Attendance) error
	GetLogsByUserID(userID string) ([]model.Attendance, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository() AttendanceRepository {
	return &attendanceRepository{db: config.GetDB()}
}

func (r *attendanceRepository) GetByUserAndDate(userID string, date time.Time) (*model.Attendance, error) {
	var attendance model.Attendance
	err := r.db.
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		Take(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) Create(attendance *model.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) Update(attendance *model.Attendance) error {
	return r.db.Save(attendance).Error
}

func (r *attendanceRepository) GetLogsByUserID(userID string) ([]model.Attendance, error) {
	var logs []model.Attendance
	err := r.db.
		Where("user_id = ?", userID).
		Order("date DESC, created_at DESC").
		Find(&logs).Error
	return logs, err
}
