package service

import (
	"context"
	"strings"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
)

type PayrollSettingService interface {
	GetAll(ctx context.Context) ([]model.PayrollSetting, error)
	Create(ctx context.Context, payload request.PayrollSettingRequest) (*model.PayrollSetting, error)
	Update(ctx context.Context, id string, payload request.PayrollSettingRequest) (*model.PayrollSetting, error)
	UpdateBulk(ctx context.Context, payload []request.PayrollSettingByKeyRequest) ([]model.PayrollSetting, error)
	Delete(ctx context.Context, id string) error
}

type payrollSettingService struct {
	payrollSettingRepo repository.PayrollSettingRepository
}

func NewPayrollSettingService(payrollSettingRepo repository.PayrollSettingRepository) PayrollSettingService {
	return &payrollSettingService{payrollSettingRepo: payrollSettingRepo}
}

func (s *payrollSettingService) GetAll(ctx context.Context) ([]model.PayrollSetting, error) {
	return s.payrollSettingRepo.GetAll(ctx)
}

func (s *payrollSettingService) Create(ctx context.Context, payload request.PayrollSettingRequest) (*model.PayrollSetting, error) {
	data := model.PayrollSetting{
		ConfigName: payload.ConfigName,
		ConfigKey:  normalizeConfigKey(payload.ConfigName),
		Value:      payload.Value,
	}
	if err := s.payrollSettingRepo.Create(ctx, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *payrollSettingService) Update(ctx context.Context, id string, payload request.PayrollSettingRequest) (*model.PayrollSetting, error) {
	data := model.PayrollSetting{
		ID:         id,
		ConfigName: payload.ConfigName,
		ConfigKey:  normalizeConfigKey(payload.ConfigName),
		Value:      payload.Value,
	}

	if err := s.payrollSettingRepo.Update(ctx, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *payrollSettingService) UpdateBulk(ctx context.Context, payload []request.PayrollSettingByKeyRequest) ([]model.PayrollSetting, error) {
	var payrollSettings []model.PayrollSetting
	for _, item := range payload {
		payrollSettings = append(payrollSettings, model.PayrollSetting{
			ConfigKey: item.ConfigKey,
			Value:     item.Value,
		})

	}
	updatedPayrollSettings, err := s.payrollSettingRepo.UpdateBulkByConfigKey(ctx, payrollSettings)
	if err != nil {
		if isNotFoundError(err) {
			return nil, common.NotFoundError(err.Error())
		}
		return nil, err
	}
	return updatedPayrollSettings, nil
}

func (s *payrollSettingService) Delete(ctx context.Context, id string) error {
	return s.payrollSettingRepo.Delete(ctx, id)
}

func normalizeConfigKey(key string) string {
	key = strings.ToLower(key)
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, " ", "_")
	return key
}
