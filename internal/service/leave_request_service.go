package service

import (
	"errors"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"gorm.io/gorm"
)

type LeaveRequestService interface {
	Create(data *model.LeaveRequest) error
	GetAll() ([]model.LeaveRequest, error)
	GetByID(id string) (*model.LeaveRequest, error)
	GetByUserID(userID string) ([]model.LeaveRequest, error)
	Update(data *model.LeaveRequest) error
	Delete(id string) error
}

type leaveRequestService struct {
	leaveRepo repository.LeaveRequestRepository
}

func NewLeaveRequestService() LeaveRequestService {
	return &leaveRequestService{leaveRepo: repository.NewLeaveRequestRepository()}
}

func (s *leaveRequestService) Create(data *model.LeaveRequest) error {
	if err := s.leaveRepo.Create(data); err != nil {
		return err
	}

	return nil
}

func (s *leaveRequestService) GetAll() ([]model.LeaveRequest, error) {
	return s.leaveRepo.GetAll()
}

func (s *leaveRequestService) GetByID(id string) (*model.LeaveRequest, error) {
	result, err := s.leaveRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("leave request not found")
		}
		return nil, err
	}
	return result, nil

}

func (s *leaveRequestService) GetByUserID(userID string) ([]model.LeaveRequest, error) {
	return s.leaveRepo.GetByUserID(userID)
}

func (s *leaveRequestService) Update(data *model.LeaveRequest) error {
	if err := s.leaveRepo.Update(data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NotFoundError("leave request not found")
		}
		return err
	}

	return nil
}

func (s *leaveRequestService) Delete(id string) error {
	err := s.leaveRepo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NotFoundError("leave request not found")
		}
		return err
	}
	return nil
}
