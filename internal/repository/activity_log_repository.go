package repository

import (
	"context"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type ActivityLogFilter struct {
	UserID     string
	Method     string
	StatusCode int
	Search     string
	Page       int
	Limit      int
}

type ActivityLogRepository interface {
	Create(ctx context.Context, log *model.ActivityLog) error
	GetAll(ctx context.Context, filter ActivityLogFilter) ([]model.ActivityLog, int64, error)
}

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, log *model.ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *activityLogRepository) GetAll(ctx context.Context, filter ActivityLogFilter) ([]model.ActivityLog, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.ActivityLog{})

	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Method != "" {
		query = query.Where("method = ?", filter.Method)
	}
	if filter.StatusCode > 0 {
		query = query.Where("status_code = ?", filter.StatusCode)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("description LIKE ? OR path LIKE ? OR user_name LIKE ?", like, like, like)
	}

	var total int64
	query.Count(&total)

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var logs []model.ActivityLog
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}
