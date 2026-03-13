package repository

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type LeaveRequestRepository interface {
	Create(leaveRequest *model.LeaveRequest) error
	GetAll() ([]model.LeaveRequest, error)
	GetByID(id string) (*model.LeaveRequest, error)
	GetByUserID(userID string) ([]model.LeaveRequest, error)
	Update(leaveRequest *model.LeaveRequest) error
	Delete(id string) error
}

type leaveRequestRepository struct {
	db *gorm.DB
}

func NewLeaveRequestRepository() LeaveRequestRepository {
	return &leaveRequestRepository{db: config.GetDB()}
}

func (r *leaveRequestRepository) Create(leaveRequest *model.LeaveRequest) error {
	return r.db.Create(&leaveRequest).Error
}

func (r *leaveRequestRepository) GetAll() ([]model.LeaveRequest, error) {
	var leaveRequests []model.LeaveRequest
	err := r.db.Order("created_at desc").Find(&leaveRequests).Error
	return leaveRequests, err
}

func (r *leaveRequestRepository) GetByID(id string) (*model.LeaveRequest, error) {
	var leaveRequest model.LeaveRequest
	err := r.db.Where("id = ?", id).First(&leaveRequest).Error
	if err != nil {
		return nil, err
	}
	return &leaveRequest, nil
}

func (r *leaveRequestRepository) GetByUserID(userID string) ([]model.LeaveRequest, error) {
	var leaveRequests []model.LeaveRequest
	result := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&leaveRequests)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return leaveRequests, nil
}

func (r *leaveRequestRepository) Update(leaveRequest *model.LeaveRequest) error {
	result := r.db.Model(&model.LeaveRequest{}).
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

func (r *leaveRequestRepository) Delete(id string) error {
	result := r.db.Delete(&model.LeaveRequest{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
