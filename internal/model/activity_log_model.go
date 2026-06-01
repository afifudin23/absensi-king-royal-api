package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityLog struct {
	ID          string    `gorm:"column:id;type:char(36);primaryKey"`
	UserID      *string   `gorm:"column:user_id;type:char(36);null"`
	UserName    *string   `gorm:"column:user_name;type:varchar(255);null"`
	Method      string    `gorm:"column:method;type:varchar(10);not null"`
	Path        string    `gorm:"column:path;type:varchar(500);not null"`
	Description string    `gorm:"column:description;type:varchar(255);not null"`
	StatusCode  int       `gorm:"column:status_code;type:int;not null"`
	IPAddress   string    `gorm:"column:ip_address;type:varchar(100);not null"`
	LatencyMs   float64   `gorm:"column:latency_ms;type:decimal(10,3);not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

func (a *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	return nil
}
