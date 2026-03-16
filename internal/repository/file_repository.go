package repository

import (
	"context"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(ctx context.Context, image *model.File) error
	GetByID(ctx context.Context, id string) (*model.File, error)
	Delete(ctx context.Context, id string) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, image *model.File) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *fileRepository) GetByID(ctx context.Context, id string) (*model.File, error) {
	var file model.File
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.File{}, "id = ?", id).Error
}
