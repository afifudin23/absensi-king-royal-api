package model

import "time"

type AttendanceRequestType string
type AttendanceRequestStatus string

const (
	AttendanceRequestTypeSick       AttendanceRequestType = "sick"
	AttendanceRequestTypeLeave      AttendanceRequestType = "leave"
	AttendanceRequestTypeExtraOff   AttendanceRequestType = "extra_off"
	AttendanceRequestTypeOvertime   AttendanceRequestType = "overtime"
	AttendanceRequestTypeCorrection AttendanceRequestType = "correction"
)

const (
	AttendanceRequestStatusPending   AttendanceRequestStatus = "pending"
	AttendanceRequestStatusApproved  AttendanceRequestStatus = "approved"
	AttendanceRequestStatusRejected  AttendanceRequestStatus = "rejected"
	AttendanceRequestStatusCancelled AttendanceRequestStatus = "cancelled"
)

type AttendanceRequest struct {
	ID                       string                  `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	UserID                   string                  `gorm:"column:user_id;type:char(36);not null"`
	AttendanceID             *string                 `gorm:"column:attendance_id;type:char(36);null"`
	Type                     AttendanceRequestType   `gorm:"column:type;type:enum('sick', 'leave', 'extra_off', 'overtime', 'correction');not null"`
	Status                   AttendanceRequestStatus `gorm:"column:status;type:enum('pending', 'approved', 'rejected', 'cancelled');not null;default:pending"`
	StartDate                time.Time               `gorm:"column:start_date;type:date;not null"`
	EndDate                  time.Time               `gorm:"column:end_date;type:date;not null"`
	RequestedCheckInAt       *time.Time              `gorm:"column:requested_check_in_at;type:datetime;null"`
	RequestedCheckOutAt      *time.Time              `gorm:"column:requested_check_out_at;type:datetime;null"`
	RequestedOvertimeMinutes *int                    `gorm:"column:requested_overtime_minutes;type:int;null"`
	Reason                   string                  `gorm:"column:reason;type:text;not null"`
	EvidenceFileID           *string                 `gorm:"column:evidence_file_id;type:char(36);null"`
	ReviewedBy               *string                 `gorm:"column:reviewed_by;type:char(36);null"`
	ReviewedAt               *time.Time              `gorm:"column:reviewed_at;type:datetime;null"`
	ReviewNote               *string                 `gorm:"column:review_note;type:text;null"`
	CreatedAt                time.Time               `gorm:"column:created_at"`
	UpdatedAt                time.Time               `gorm:"column:updated_at"`

	EvidenceFile *File `gorm:"foreignKey:EvidenceFileID;references:ID"`
}

func (AttendanceRequest) TableName() string {
	return "attendance_requests"
}
