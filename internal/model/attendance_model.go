package model

import "time"

type Attendance struct {
	ID              string     `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	UserID          string     `gorm:"column:user_id;type:char(36);not null"`
	Date            time.Time  `gorm:"column:date;type:date;not null"`
	CheckInAt       *time.Time `gorm:"column:check_in_at;type:date;null"`
	CheckOutAt      *time.Time `gorm:"column:check_out_at;type:date;null"`
	CheckInFileID   *string    `gorm:"column:check_in_file_id;type:char(36);null"`
	CheckInFileURL  *string    `gorm:"column:check_in_file_url;type:text;null"`
	CheckOutFileID  *string    `gorm:"column:check_out_file_id;type:char(36);null"`
	CheckOutFileURL *string    `gorm:"column:check_out_file_url;type:text;null"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

func (Attendance) TableName() string {
	return "attendances"
}
