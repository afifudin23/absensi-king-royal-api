package repository

import (
	"context"
	"fmt"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type PayrollSettingRepository interface {
	GetAll(ctx context.Context) ([]model.PayrollSetting, error)
	Create(ctx context.Context, payrollSetting *model.PayrollSetting) error
	Update(ctx context.Context, payrollSetting *model.PayrollSetting) error
	UpdateBulkByConfigKey(ctx context.Context, payrollSettings []model.PayrollSetting) ([]model.PayrollSetting, error)
	Delete(ctx context.Context, ids []string) (int, error)
}

type payrollSettingRepository struct {
	db *gorm.DB
}

func NewPayrollSettingRepository(db *gorm.DB) PayrollSettingRepository {
	return &payrollSettingRepository{db: db}
}

func (r *payrollSettingRepository) GetAll(ctx context.Context) ([]model.PayrollSetting, error) {
	var payrollSettings []model.PayrollSetting
	err := r.db.WithContext(ctx).
		Order("updated_at DESC").
		Find(&payrollSettings).
		Error
	return payrollSettings, err
}

func (r *payrollSettingRepository) Create(ctx context.Context, payrollSetting *model.PayrollSetting) error {
	return r.db.WithContext(ctx).Create(payrollSetting).Error
}

func (r *payrollSettingRepository) Update(ctx context.Context, payrollSetting *model.PayrollSetting) error {
	result := r.db.WithContext(ctx).Model(&model.PayrollSetting{}).
		Where("id = ?", payrollSetting.ID).
		Updates(payrollSetting)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *payrollSettingRepository) UpdateBulkByConfigKey(ctx context.Context, payrollSettings []model.PayrollSetting) ([]model.PayrollSetting, error) {
	var result []model.PayrollSetting

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range payrollSettings {
			updateResult := tx.Model(&model.PayrollSetting{}).
				Where("config_key = ?", item.ConfigKey).
				Update("value", item.Value)
			if updateResult.Error != nil {
				return updateResult.Error
			}
			if updateResult.RowsAffected == 0 {
				return fmt.Errorf("Payroll setting with config_key %s not found", item.ConfigKey)
			}

			var updated model.PayrollSetting
			if err := tx.Where("config_key = ?", item.ConfigKey).First(&updated).Error; err != nil {
				return err
			}

			result = append(result, updated)

		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *payrollSettingRepository) Delete(ctx context.Context, ids []string) (int, error) {
	result := r.db.WithContext(ctx).Delete(&model.PayrollSetting{}, "id IN ?", ids)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
