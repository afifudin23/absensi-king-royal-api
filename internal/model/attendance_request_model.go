package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceRequestType string
type AttendanceRequestStatus string

const (
	AttendanceRequestTypeSick     AttendanceRequestType = "sick"
	AttendanceRequestTypeLeave    AttendanceRequestType = "leave"
	AttendanceRequestTypeExtraOff AttendanceRequestType = "extra_off"
	AttendanceRequestTypeOvertime AttendanceRequestType = "overtime"
)

const (
	AttendanceRequestStatusPending  AttendanceRequestStatus = "pending"
	AttendanceRequestStatusApproved AttendanceRequestStatus = "approved"
	AttendanceRequestStatusRejected AttendanceRequestStatus = "rejected"
)

type AttendanceRequest struct {
	ID     string                `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	UserID string                `gorm:"column:user_id;type:char(36);not null"`
	Type   AttendanceRequestType `gorm:"column:type;type:enum('sick', 'leave', 'extra_off', 'overtime');not null"`
	Status                 AttendanceRequestStatus `gorm:"column:status;type:enum('pending', 'approved', 'rejected');not null;default:pending"`
	StartDate              time.Time               `gorm:"column:start_date;type:date;not null"`
	EndDate                time.Time               `gorm:"column:end_date;type:date;not null"`
	RequestedOvertimeHours *int                    `gorm:"column:requested_overtime_hours;type:int;null"`
	Reason                 string                  `gorm:"column:reason;type:text;not null"`
	EvidenceFileID         *string                 `gorm:"column:evidence_file_id;type:char(36);null"`
	ReviewedBy *string    `gorm:"column:reviewed_by;type:char(36);null"`
	ReviewedAt *time.Time `gorm:"column:reviewed_at;type:datetime;null"`
	ReviewNote *string    `gorm:"column:review_note;type:text;null"`
	CreatedAt              time.Time               `gorm:"column:created_at"`
	UpdatedAt              time.Time               `gorm:"column:updated_at"`

	EvidenceFile *File `gorm:"foreignKey:EvidenceFileID;references:ID"`
	User         *User `gorm:"foreignKey:UserID;references:ID"`
}

func (ar *AttendanceRequest) BeforeCreate(tx *gorm.DB) error {
	if ar.ID == "" {
		ar.ID = uuid.NewString()
	}
	return nil
}

func (AttendanceRequest) TableName() string {
	return "attendance_requests"
}
