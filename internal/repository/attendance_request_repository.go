package repository

import (
	"context"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type AttendanceRequestRepository interface {
	Create(ctx context.Context, attendanceRequest *model.AttendanceRequest) error
	GetAll(ctx context.Context, loadFile bool) ([]model.AttendanceRequest, error)
	GetByID(ctx context.Context, id string, loadFile bool) (*model.AttendanceRequest, error)
	GetByUserID(ctx context.Context, userID string) ([]model.AttendanceRequest, error)
	Update(ctx context.Context, attendanceRequest *model.AttendanceRequest) error
	Delete(ctx context.Context, id string) error
}

type attendanceRequestRepository struct {
	db *gorm.DB
}

func NewAttendanceRequestRepository(db *gorm.DB) AttendanceRequestRepository {
	return &attendanceRequestRepository{db: db}
}

func (r *attendanceRequestRepository) Create(ctx context.Context, attendanceRequest *model.AttendanceRequest) error {
	return r.db.WithContext(ctx).Create(attendanceRequest).Error
}

func (r *attendanceRequestRepository) GetAll(ctx context.Context, loadFile bool) ([]model.AttendanceRequest, error) {
	var attendanceRequests []model.AttendanceRequest

	db := r.db.WithContext(ctx)

	if loadFile {
		db = db.Preload("EvidenceFile")
	}

	err := db.
		Order("created_at desc").
		Find(&attendanceRequests).
		Error
	return attendanceRequests, err
}

func (r *attendanceRequestRepository) GetByID(ctx context.Context, id string, loadFile bool) (*model.AttendanceRequest, error) {
	var attendanceRequest model.AttendanceRequest

	db := r.db.WithContext(ctx)

	if loadFile {
		db = db.Preload("EvidenceFile")
	}

	err := db.Where("id = ?", id).First(&attendanceRequest).Error
	if err != nil {
		return nil, err
	}
	return &attendanceRequest, nil
}

func (r *attendanceRequestRepository) GetByUserID(ctx context.Context, userID string) ([]model.AttendanceRequest, error) {
	var attendanceRequests []model.AttendanceRequest
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&attendanceRequests)
	if result.Error != nil {
		return nil, result.Error
	}
	return attendanceRequests, nil
}

func (r *attendanceRequestRepository) Update(ctx context.Context, attendanceRequest *model.AttendanceRequest) error {
	result := r.db.WithContext(ctx).Model(&model.AttendanceRequest{}).
		Where("id = ?", attendanceRequest.ID).
		Updates(attendanceRequest)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *attendanceRequestRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.AttendanceRequest{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
