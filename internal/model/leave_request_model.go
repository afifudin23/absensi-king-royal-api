package model

import "time"

type LeaveRequestType string
type LeaveRequestStatus string

const (
	LeaveTypeSick     LeaveRequestType = "sick"
	LeaveTypeExtraOff LeaveRequestType = "extra_off"
	LeaveTypeOvertime LeaveRequestType = "overtime"
	LeaveTypeLeave    LeaveRequestType = "leave"
)

const (
	LeaveRequestStatusPending  LeaveRequestStatus = "pending"
	LeaveRequestStatusApproved LeaveRequestStatus = "approved"
	LeaveRequestStatusRejected LeaveRequestStatus = "rejected"
)

type LeaveRequest struct {
	ID        string             `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	UserID    string             `gorm:"column:user_id;type:char(36);not null"`
	StartDate time.Time          `gorm:"column:start_date;type:date;not null"`
	EndDate   time.Time          `gorm:"column:end_date;type:date;not null"`
	Reason    string             `gorm:"column:reason;type:text;not null"`
	Type      LeaveRequestType   `gorm:"column:type;type:enum('sick', 'extra_off', 'overtime', 'leave');not null"`
	Status    LeaveRequestStatus `gorm:"column:status;type:enum('pending', 'approved', 'rejected');not null;default:pending"`

	// Optional fields
	EvidenceFileID  *string  `gorm:"column:evidence_file_id;type:char(36);null"`
	EvidenceFileURL *string  `gorm:"column:evidence_file_url;type:text;null"`
	OvertimeHours   *float64 `gorm:"column:overtime_hours;type:decimal(5,2);null"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (LeaveRequest) TableName() string {
	return "leave_requests"
}
