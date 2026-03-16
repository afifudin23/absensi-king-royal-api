package response

import (
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
)

type FileResponse struct {
	ID         string         `json:"id"`
	FileURL    string         `json:"file_url"`
	FileName   string         `json:"file_name"`
	FileSize   int64          `json:"file_size"`
	FilePath   string         `json:"file_path"`
	MimeType   string         `json:"mime_type"`
	UploadedBy string         `json:"uploaded_by"`
	Type       model.FileType `json:"type"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
}

func ToFileResponse(file model.File) FileResponse {
	return FileResponse{
		ID:         file.ID,
		FileURL:    file.FileURL,
		FileName:   file.FileName,
		FileSize:   file.FileSize,
		FilePath:   file.FilePath,
		MimeType:   file.MimeType,
		UploadedBy: file.UploadedBy,
		Type:       file.Type,
		CreatedAt:  file.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  file.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
