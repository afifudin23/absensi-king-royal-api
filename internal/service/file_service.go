package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
)

type FileService interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader, fileType model.FileType, userID string) (*model.File, error)
	Delete(ctx context.Context, id string) error
}

type fileService struct {
	fileRepo repository.FileRepository
}

func NewFileService(fileRepo repository.FileRepository) FileService {
	return &fileService{fileRepo: fileRepo}
}

func folderNameForFileType(fileType model.FileType) string {
	// Keep the DB enum `files.type` as-is for validation, but group leave evidences
	// into a single folder for storage/URLs.
	switch fileType {
	case model.FileTypeSick, model.FileTypeExtraOff, model.FileTypeOvertime, model.FileTypeLeave:
		return "evidence"
	default:
		return string(fileType)
	}
}

func (s *fileService) Upload(ctx context.Context, fileHeader *multipart.FileHeader, fileType model.FileType, userID string) (*model.File, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	newFileName := fmt.Sprintf("%s%s", time.Now().Format("20060102_150405"), ext)

	dateFolder := time.Now().Format("2006-01-02")
	folderName := folderNameForFileType(fileType)
	folderPath := filepath.Join("files", folderName, dateFolder)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return nil, err
	}

	filePath := filepath.Join(folderPath, newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	fileURL := fmt.Sprintf("%s/%s/%s", config.GetEnv().ServerBaseURL, "files", path.Join(folderName, dateFolder, newFileName))

	image := &model.File{
		FileName:   newFileName,
		MimeType:   fileHeader.Header.Get("Content-Type"),
		FileSize:   fileHeader.Size,
		FilePath:   filePath,
		FileURL:    fileURL,
		UploadedBy: userID,
		Type:       fileType,
	}

	if err := s.fileRepo.Create(ctx, image); err != nil {
		return nil, err
	}

	return image, nil
}

func (s *fileService) Delete(ctx context.Context, id string) error {
	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		return nil
	}

	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return s.fileRepo.Delete(ctx, id)
}
