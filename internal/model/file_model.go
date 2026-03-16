package model

import "time"

type FileType string

const (
	FileTypeCheckIn        FileType = "check_in"
	FileTypeCheckOut       FileType = "check_out"
	FileTypeProfilePicture FileType = "profile_picture"
	FileTypeSick           FileType = "sick"
	FileTypeExtraOff       FileType = "extra_off"
	FileTypeOvertime       FileType = "overtime"
	FileTypeLeave          FileType = "leave"
)

type File struct {
	ID         string    `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	FileName   string    `gorm:"column:file_name;type:varchar(255);not null"`
	MimeType   string    `gorm:"column:mime_type;type:varchar(100);not null"`
	FileSize   int64     `gorm:"column:file_size;type:bigint;not null"`
	FilePath   string    `gorm:"column:file_path;type:text;not null"`
	FileURL    string    `gorm:"column:file_url;type:text;not null"`
	UploadedBy string    `gorm:"column:uploaded_by;type:char(36);not null"`
	Type       FileType  `gorm:"column:type;type:enum('check_in', 'check_out', 'profile_picture', 'sick', 'extra_off', 'overtime', 'leave');not null"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (File) TableName() string {
	return "files"
}
