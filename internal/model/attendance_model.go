package model

import "time"

type AttendanceStatus string
type AttendanceSource string

const (
	AttendanceStatusPresent  AttendanceStatus = "present"
	AttendanceStatusOff      AttendanceStatus = "off"
	AttendanceStatusSick     AttendanceStatus = "sick"
	AttendanceStatusExtraOff AttendanceStatus = "extra_off"
	AttendanceStatusAbsent   AttendanceStatus = "absent"
	AttendanceStatusLeave    AttendanceStatus = "leave"
)

const (
	AttendanceSourceSelfService     AttendanceSource = "self_service"
	AttendanceSourceAdminEdit       AttendanceSource = "admin_edit"
	AttendanceSourceApprovedRequest AttendanceSource = "approved_request"
	AttendanceSourceSystem          AttendanceSource = "system"
)

type Attendance struct {
	ID             string           `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	UserID         string           `gorm:"column:user_id;type:char(36);not null"`
	Status         AttendanceStatus `gorm:"column:status;type:enum('present', 'off', 'sick', 'extra_off', 'absent', 'leave');not null;default:present"`
	Date           time.Time        `gorm:"column:date;type:date;not null"`
	CheckInAt      *time.Time       `gorm:"column:check_in_at;type:date;null"`
	CheckOutAt     *time.Time       `gorm:"column:check_out_at;type:date;null"`
	CheckInFileID  *string          `gorm:"column:check_in_file_id;type:char(36);null"`
	CheckOutFileID *string          `gorm:"column:check_out_file_id;type:char(36);null"`
	Note           *string          `gorm:"column:note;type:text;null"`
	Source         AttendanceSource `gorm:"column:source;type:enum('self_service', 'admin_edit', 'approved_request', 'system');not null;default:self_service"`
	UpdatedBy      *string          `gorm:"column:updated_by;type:char(36);null"`
	CreatedAt      time.Time        `gorm:"column:created_at"`
	UpdatedAt      time.Time        `gorm:"column:updated_at"`

	CheckInFile  *File `gorm:"foreignKey:CheckInFileID;references:ID"`
	CheckOutFile *File `gorm:"foreignKey:CheckOutFileID;references:ID"`
}

func (Attendance) TableName() string {
	return "attendances"
}
