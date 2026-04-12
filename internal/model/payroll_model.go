package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PayrollStatus string

const (
	PayrollStatusUnsent PayrollStatus = "unsent"
	PayrollStatusSent   PayrollStatus = "sent"
	PayrollStatusFailed PayrollStatus = "failed"
)

type Payroll struct {
	ID                  string         `gorm:"column:id;type:char(36);primaryKey;default:(UUID())"`
	EmployeeID          string         `gorm:"column:employee_id;type:char(36);not null"`
	BasicSalary         float64        `gorm:"column:basic_salary;type:decimal(15,2);null;default:0"`
	PositionAllowance   float64        `gorm:"column:position_allowance;type:decimal(15,2);null;default:0"`
	OtherAllowance      float64        `gorm:"column:other_allowance;type:decimal(15,2);null;default:0"`
	OvertimeRate        float64        `gorm:"column:overtime_rate;type:decimal(15,2);null;default:0"`
	LoanDeduction       float64        `gorm:"column:loan_deduction;type:decimal(15,2);null;default:0"`
	AttendanceDeduction float64        `gorm:"column:attendance_deduction;type:decimal(15,2);null;default:0"`
	IncomeTax           float64        `gorm:"column:income_tax;type:decimal(15,2);null;default:0"`
	AdditionalData      datatypes.JSON `gorm:"column:additional_data;type:json;default:'{}'"`
	NetSalary           float64        `gorm:"column:net_salary;type:decimal(15,2);null;default:0"`
	Status              PayrollStatus  `gorm:"column:status;type:enum('unsent','sent','failed');default:'unsent'"`
	PDFPath             *string        `gorm:"column:pdf_path;type:text;null"`
	SentAt              *time.Time     `gorm:"column:sent_at;type:datetime;null"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
}

func (Payroll) TableName() string {
	return "payrolls"
}

func (p *Payroll) BeforeCreate(tx *gorm.DB) error {
	if len(p.AdditionalData) == 0 {
		p.AdditionalData = datatypes.JSON([]byte(`{}`))
	}
	return nil
}
