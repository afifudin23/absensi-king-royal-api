package repository

import (
	"context"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type LeaveRequestRepository interface {
	Create(ctx context.Context, leaveRequest *model.LeaveRequest) error
	GetAll(ctx context.Context) ([]model.LeaveRequest, error)
	GetByID(ctx context.Context, id string) (*model.LeaveRequest, error)
	GetByUserID(ctx context.Context, userID string) ([]model.LeaveRequest, error)
	Update(ctx context.Context, leaveRequest *model.LeaveRequest) error
	Delete(ctx context.Context, id string) error
}

type leaveRequestRepository struct {
	db *gorm.DB
}

func NewLeaveRequestRepository(db *gorm.DB) LeaveRequestRepository {
	return &leaveRequestRepository{db: db}
}

func (r *leaveRequestRepository) Create(ctx context.Context, leaveRequest *model.LeaveRequest) error {
	return r.db.WithContext(ctx).Create(leaveRequest).Error
}

func (r *leaveRequestRepository) GetAll(ctx context.Context) ([]model.LeaveRequest, error) {
	var leaveRequests []model.LeaveRequest
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&leaveRequests).Error
	return leaveRequests, err
}

func (r *leaveRequestRepository) GetByID(ctx context.Context, id string) (*model.LeaveRequest, error) {
	var leaveRequest model.LeaveRequest
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&leaveRequest).Error
	if err != nil {
		return nil, err
	}
	return &leaveRequest, nil
}

func (r *leaveRequestRepository) GetByUserID(ctx context.Context, userID string) ([]model.LeaveRequest, error) {
	var leaveRequests []model.LeaveRequest
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&leaveRequests)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return leaveRequests, nil
}

func (r *leaveRequestRepository) Update(ctx context.Context, leaveRequest *model.LeaveRequest) error {
	result := r.db.WithContext(ctx).Model(&model.LeaveRequest{}).
		Where("id = ?", leaveRequest.ID).
		Updates(leaveRequest)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *leaveRequestRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.LeaveRequest{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
