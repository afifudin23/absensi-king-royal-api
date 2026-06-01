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

type AttendanceRequestService interface {
	Create(ctx context.Context, userID string, payload request.AttendanceRequestCreateRequest) (*model.AttendanceRequest, error)
	GetAll(ctx context.Context, filter *repository.AttendanceRequestFilter) ([]model.AttendanceRequest, error)
	GetByID(ctx context.Context, id string) (*model.AttendanceRequest, error)
	GetByUserID(ctx context.Context, userID string) ([]model.AttendanceRequest, error)
	Update(ctx context.Context, userID string, id string, payload request.AttendanceRequestUpdateRequest) (*model.AttendanceRequest, error)
	UpdateStatus(ctx context.Context, reviewerID string, id string, payload request.AttendanceRequestUpdateStatusRequest) (*model.AttendanceRequest, error)
	Delete(ctx context.Context, ids []string) error
}

type attendanceRequestService struct {
	attendanceRequestRepo repository.AttendanceRequestRepository
	attendanceRepo        repository.AttendanceRepository
	fileRepo              repository.FileRepository
}

func NewAttendanceRequestService(attendanceRequestRepo repository.AttendanceRequestRepository, attendanceRepo repository.AttendanceRepository, fileRepo repository.FileRepository) AttendanceRequestService {
	return &attendanceRequestService{
		attendanceRequestRepo: attendanceRequestRepo,
		attendanceRepo:        attendanceRepo,
		fileRepo:              fileRepo,
	}
}

func (s *attendanceRequestService) Create(ctx context.Context, userID string, payload request.AttendanceRequestCreateRequest) (*model.AttendanceRequest, error) {
	startDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		return nil, common.BadRequestError("Start date must be in YYYY-MM-DD format")
	}

	endDate, err := time.Parse("2006-01-02", payload.EndDate)
	if err != nil {
		return nil, common.BadRequestError("End date must be in YYYY-MM-DD format")
	}

	if endDate.Before(startDate) {
		return nil, common.BadRequestError("End date must be greater than or equal to start date")
	}

	overlapping, err := s.attendanceRequestRepo.HasOverlappingRequest(ctx, userID, startDate, endDate, nil)
	if err != nil {
		return nil, err
	}
	if overlapping {
		return nil, common.BadRequestError("Sudah ada pengajuan izin yang pending atau approved pada tanggal tersebut")
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
	}

	data := &model.AttendanceRequest{
		UserID:                 userID,
		StartDate:              startDate,
		EndDate:                endDate,
		Reason:                 payload.Reason,
		Type:                   payload.Type,
		EvidenceFileID:         payload.EvidenceFileID,
		RequestedOvertimeHours: payload.RequestedOvertimeHours,
		Status:                 model.AttendanceRequestStatusPending,
	}

	if err := s.attendanceRequestRepo.Create(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *attendanceRequestService) GetAll(ctx context.Context, filter *repository.AttendanceRequestFilter) ([]model.AttendanceRequest, error) {
	return s.attendanceRequestRepo.GetAll(ctx, true, filter)
}

func (s *attendanceRequestService) GetByID(ctx context.Context, id string) (*model.AttendanceRequest, error) {
	result, err := s.attendanceRequestRepo.GetByID(ctx, id, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Attendance request not found")
		}
		return nil, err
	}
	return result, nil
}

func (s *attendanceRequestService) GetByUserID(ctx context.Context, userID string) ([]model.AttendanceRequest, error) {
	return s.attendanceRequestRepo.GetByUserID(ctx, userID)
}

func (s *attendanceRequestService) Update(ctx context.Context, userID string, id string, payload request.AttendanceRequestUpdateRequest) (*model.AttendanceRequest, error) {
	existing, err := s.attendanceRequestRepo.GetByID(ctx, id, false)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.NotFoundError("Attendance request not found")
		}
		return nil, err
	}
	if existing.UserID != userID {
		return nil, common.ForbiddenError("Attendance request does not belong to current user")
	}

	data := &model.AttendanceRequest{ID: id}

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

		data.EvidenceFileID = payload.EvidenceFileID
	}

	if payload.RequestedOvertimeHours != nil {
		data.RequestedOvertimeHours = payload.RequestedOvertimeHours
	}

	if err := s.attendanceRequestRepo.Update(ctx, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Attendance request not found")
		}
		return nil, err
	}

	return existing, nil
}

func (s *attendanceRequestService) UpdateStatus(ctx context.Context, reviewerID string, id string, payload request.AttendanceRequestUpdateStatusRequest) (*model.AttendanceRequest, error) {
	existing, err := s.attendanceRequestRepo.GetByID(ctx, id, false)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.NotFoundError("Attendance request not found")
		}
		return nil, err
	}

	now := time.Now()

	existing.Status = payload.Status
	existing.ReviewedBy = &reviewerID
	existing.ReviewedAt = &now

	if err := s.attendanceRequestRepo.Update(ctx, existing); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Attendance request not found")
		}
		return nil, err
	}

	if payload.Status == model.AttendanceRequestStatusApproved {
		if err := s.applyApprovedRequestToAttendance(ctx, existing, reviewerID); err != nil {
			return nil, err
		}
	}

	return existing, nil
}

func (s *attendanceRequestService) Delete(ctx context.Context, ids []string) error {
	err := s.attendanceRequestRepo.Delete(ctx, ids)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NotFoundError("Attendance request not found")
		}
		return err
	}
	return nil
}

func (s *attendanceRequestService) applyApprovedRequestToAttendance(ctx context.Context, req *model.AttendanceRequest, reviewerID string) error {
	switch req.Type {
	case model.AttendanceRequestTypeSick, model.AttendanceRequestTypeLeave, model.AttendanceRequestTypeExtraOff:
		status := model.AttendanceStatus(req.Type)
		note := req.Reason

		for d := req.StartDate; !d.After(req.EndDate); d = d.AddDate(0, 0, 1) {
			day := startOfDay(d)
			if err := s.upsertAttendanceForDay(ctx, req.UserID, day, func(a *model.Attendance) {
				a.Status = status
				a.Note = &note
				a.CheckInAt = nil
				a.CheckOutAt = nil
				a.CheckInFileID = nil
				a.CheckOutFileID = nil
				a.EvidenceFileID = req.EvidenceFileID
				a.Source = model.AttendanceSourceApprovedRequest
				a.UpdatedBy = &reviewerID
			}); err != nil {
				return err
			}
		}
		return nil

	case model.AttendanceRequestTypeOvertime:
		// Isi overtime_hours ke attendance saat request diapprove
		for d := req.StartDate; !d.After(req.EndDate); d = d.AddDate(0, 0, 1) {
			day := startOfDay(d)
			if err := s.upsertAttendanceForDay(ctx, req.UserID, day, func(a *model.Attendance) {
				a.OvertimeHours = req.RequestedOvertimeHours
				a.Source = model.AttendanceSourceApprovedRequest
				a.UpdatedBy = &reviewerID
			}); err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}

func (s *attendanceRequestService) upsertAttendanceForDay(ctx context.Context, userID string, day time.Time, apply func(a *model.Attendance)) error {
	attendance, err := s.attendanceRepo.GetByUserAndDate(ctx, userID, day)
	if err != nil {
		if !isNotFoundError(err) {
			return err
		}

		attendance = &model.Attendance{
			UserID: userID,
			Date:   day,
		}
		apply(attendance)

		if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
			// Race protection: if a duplicate row is inserted concurrently, re-fetch it.
			if isDuplicateError(err) {
				attendance, err = s.attendanceRepo.GetByUserAndDate(ctx, userID, day)
				if err != nil {
					return err
				}
				apply(attendance)
				return s.attendanceRepo.Update(ctx, attendance)
			}
			return err
		}
		return nil
	}

	apply(attendance)
	return s.attendanceRepo.Update(ctx, attendance)
}
