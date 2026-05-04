package repository

import (
	"context"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	GetByUserAndDate(ctx context.Context, userID string, date time.Time) (*model.Attendance, error)
	Create(ctx context.Context, attendance *model.Attendance) error
	Update(ctx context.Context, attendance *model.Attendance) error
	GetLogsByUserID(ctx context.Context, userID string) ([]model.Attendance, error)
	GetByID(ctx context.Context, id string) (*model.Attendance, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) GetByUserAndDate(ctx context.Context, userID string, date time.Time) (*model.Attendance, error) {
	var attendance model.Attendance
	err := r.db.WithContext(ctx).
		Preload("CheckInFile").
		Preload("CheckOutFile").
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		Take(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) Create(ctx context.Context, attendance *model.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

func (r *attendanceRepository) Update(ctx context.Context, attendance *model.Attendance) error {
	return r.db.WithContext(ctx).Save(attendance).Error
}

func (r *attendanceRepository) GetLogsByUserID(ctx context.Context, userID string) ([]model.Attendance, error) {
	var logs []model.Attendance
	err := r.db.WithContext(ctx).
		Preload("CheckInFile").
		Preload("CheckOutFile").
		Where("user_id = ?", userID).
		Order("date DESC, created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *attendanceRepository) GetByID(ctx context.Context, id string) (*model.Attendance, error) {
	var attendance model.Attendance
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Take(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}
