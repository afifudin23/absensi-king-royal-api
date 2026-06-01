package repository

import (
	"context"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type AttendanceRecapRow struct {
	model.Attendance
	EmployeeName    string  `gorm:"column:employee_name"`
	CheckInFileURL  *string `gorm:"column:check_in_file_url"`
	CheckOutFileURL *string `gorm:"column:check_out_file_url"`
	EvidenceFileURL *string `gorm:"column:evidence_file_url"`
}

type AttendanceRepository interface {
	GetByUserAndDate(ctx context.Context, userID string, date time.Time) (*model.Attendance, error)
	Create(ctx context.Context, attendance *model.Attendance) error
	Update(ctx context.Context, attendance *model.Attendance) error
	GetLogsByUserID(ctx context.Context, userID string, startDate, endDate *time.Time) ([]model.Attendance, error)
	GetByID(ctx context.Context, id string) (*model.Attendance, error)
	GetByMonthWithUser(ctx context.Context, month, year int) ([]AttendanceRecapRow, error)
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

func (r *attendanceRepository) GetLogsByUserID(ctx context.Context, userID string, startDate, endDate *time.Time) ([]model.Attendance, error) {
	var logs []model.Attendance
	q := r.db.WithContext(ctx).
		Preload("CheckInFile").
		Preload("CheckOutFile").
		Preload("EvidenceFile").
		Where("user_id = ?", userID)
	if startDate != nil {
		q = q.Where("date >= ?", startDate.Format("2006-01-02"))
	}
	if endDate != nil {
		q = q.Where("date <= ?", endDate.Format("2006-01-02"))
	}
	err := q.Order("date DESC, created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *attendanceRepository) GetByMonthWithUser(ctx context.Context, month, year int) ([]AttendanceRecapRow, error) {
	var rows []AttendanceRecapRow
	err := r.db.WithContext(ctx).
		Model(&model.Attendance{}).
		Select("attendances.*, users.full_name AS employee_name, f1.file_url AS check_in_file_url, f2.file_url AS check_out_file_url, f3.file_url AS evidence_file_url").
		Joins("JOIN users ON users.id = attendances.user_id").
		Joins("LEFT JOIN files f1 ON f1.id = attendances.check_in_file_id").
		Joins("LEFT JOIN files f2 ON f2.id = attendances.check_out_file_id").
		Joins("LEFT JOIN files f3 ON f3.id = attendances.evidence_file_id").
		Where("MONTH(attendances.date) = ? AND YEAR(attendances.date) = ?", month, year).
		Order("attendances.user_id ASC, attendances.date ASC").
		Scan(&rows).Error
	return rows, err
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
