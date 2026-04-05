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

type LeaveRequestService interface {
	Create(ctx context.Context, userID string, payload request.LeaveRequestCreateRequest) (*model.LeaveRequest, error)
	GetAll(ctx context.Context) ([]model.LeaveRequest, error)
	GetByID(ctx context.Context, id string) (*model.LeaveRequest, error)
	GetByUserID(ctx context.Context, userID string) ([]model.LeaveRequest, error)
	Update(ctx context.Context, userID string, id string, payload request.LeaveRequestUpdateRequest) (*model.LeaveRequest, error)
	UpdateStatus(ctx context.Context, id string, payload request.LeaveRequestUpdateStatusRequest) (*model.LeaveRequest, error)
	Delete(ctx context.Context, id string) error
}

type leaveRequestService struct {
	leaveRepo repository.LeaveRequestRepository
	fileRepo  repository.FileRepository
}

func NewLeaveRequestService(leaveRepo repository.LeaveRequestRepository, fileRepo repository.FileRepository) LeaveRequestService {
	return &leaveRequestService{leaveRepo: leaveRepo, fileRepo: fileRepo}
}

func (s *leaveRequestService) Create(ctx context.Context, userID string, payload request.LeaveRequestCreateRequest) (*model.LeaveRequest, error) {

	startDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		return nil, common.BadRequestError("Start date must be in YYYY-MM-DD format")
	}

	endDate, err := time.Parse("2006-01-02", payload.EndDate)
	if err != nil {
		return nil, common.BadRequestError("End date must be in YYYY-MM-DD format")
	}

	var evidenceFileURL *string
	if payload.EvidenceFileID != nil {
		file, err := s.fileRepo.GetByID(ctx, *payload.EvidenceFileID)
		if err != nil {
			if isNotFoundError(err) {
				return nil, common.BadRequestError("Invalid file_id")
			}
			return nil, err
		}

		if file.UploadedBy != userID {
			return nil, common.ForbiddenError("File does not belong to current user")
		}

		url := file.FileURL
		evidenceFileURL = &url
	}

	data := &model.LeaveRequest{
		UserID:          userID,
		StartDate:       startDate,
		EndDate:         endDate,
		Reason:          payload.Reason,
		Type:            payload.Type,
		EvidenceFileID:  payload.EvidenceFileID,
		EvidenceFileURL: evidenceFileURL,
		OvertimeHours:   payload.OvertimeHours,
		Status:          model.LeaveRequestStatusPending,
	}

	if err := s.leaveRepo.Create(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *leaveRequestService) GetAll(ctx context.Context) ([]model.LeaveRequest, error) {
	return s.leaveRepo.GetAll(ctx)
}

func (s *leaveRequestService) GetByID(ctx context.Context, id string) (*model.LeaveRequest, error) {
	result, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Leave request not found")
		}
		return nil, err
	}
	return result, nil

}

func (s *leaveRequestService) GetByUserID(ctx context.Context, userID string) ([]model.LeaveRequest, error) {
	return s.leaveRepo.GetByUserID(ctx, userID)
}

func (s *leaveRequestService) Update(ctx context.Context, userID string, id string, payload request.LeaveRequestUpdateRequest) (*model.LeaveRequest, error) {
	existing, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.NotFoundError("Leave request not found")
		}
		return nil, err
	}
	if existing.UserID != userID {
		return nil, common.ForbiddenError("Leave request does not belong to current user")
	}

	data := &model.LeaveRequest{ID: id}

	// If type is being updated but file_id isn't being replaced, validate the existing file still matches the type.
	if payload.Type != nil && payload.EvidenceFileID == nil && existing.EvidenceFileID != nil {
		file, err := s.fileRepo.GetByID(ctx, *existing.EvidenceFileID)
		if err != nil {
			if isNotFoundError(err) {
				return nil, common.BadRequestError("Invalid file_id")
			}
			return nil, err
		}
		if file.UploadedBy != userID {
			return nil, common.ForbiddenError("File does not belong to current user")
		}
	}

	if payload.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *payload.StartDate)
		if err != nil {
			return nil, common.BadRequestError("Start date must be in YYYY-MM-DD format")
		}
		data.StartDate = startDate
	}

	if payload.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *payload.EndDate)
		if err != nil {
			return nil, common.BadRequestError("End date must be in YYYY-MM-DD format")
		}
		data.EndDate = endDate
	}

	if payload.Reason != nil {
		data.Reason = *payload.Reason
	}

	if payload.Type != nil {
		data.Type = *payload.Type
	}

	if payload.EvidenceFileID != nil {
		file, err := s.fileRepo.GetByID(ctx, *payload.EvidenceFileID)
		if err != nil {
			if isNotFoundError(err) {
				return nil, common.BadRequestError("Invalid file_id")
			}
			return nil, err
		}
		if file.UploadedBy != userID {
			return nil, common.ForbiddenError("File does not belong to current user")
		}

		fileURL := file.FileURL
		data.EvidenceFileID = payload.EvidenceFileID
		data.EvidenceFileURL = &fileURL
	}

	if payload.OvertimeHours != nil {
		data.OvertimeHours = payload.OvertimeHours
	}

	if err := s.leaveRepo.Update(ctx, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Leave request not found")
		}
		return nil, err
	}

	updated, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Leave request not found")
		}
		return nil, err
	}

	return updated, nil
}

func (s *leaveRequestService) UpdateStatus(ctx context.Context, id string, payload request.LeaveRequestUpdateStatusRequest) (*model.LeaveRequest, error) {
	existing, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Leave request not found")
		}
		return nil, err
	}

	data := &model.LeaveRequest{
		ID:     id,
		Status: payload.Status,
	}

	if err := s.leaveRepo.Update(ctx, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Leave request not found")
		}
		return nil, err
	}

	existing.Status = payload.Status
	return existing, nil
}

func (s *leaveRequestService) Delete(ctx context.Context, id string) error {
	err := s.leaveRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NotFoundError("Leave request not found")
		}
		return err
	}
	return nil
}
