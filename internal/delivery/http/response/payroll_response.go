package response

import (
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/datatypes"
)

type PayrollResponse struct {
	ID                  string         `json:"id"`
	EmployeeID          string         `json:"employee_id"`
	BasicSalary         float64        `json:"basic_salary"`
	PositionAllowance   float64        `json:"position_allowance"`
	OtherAllowance      float64        `json:"other_allowance"`
	OvertimeRate        float64        `json:"overtime_rate"`
	LoanDeduction       float64        `json:"loan_deduction"`
	AttendanceDeduction float64        `json:"attendance_deduction"`
	IncomeTax           float64        `json:"income_tax"`
	AdditionalData      datatypes.JSON `json:"additional_data"`
	NetSalary           float64        `json:"net_salary"`
	Status              string         `json:"status"`
	PDFPath             *string        `json:"pdf_path"`
	SentAt              *time.Time     `json:"sent_at"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

func ToPayrollResponse(data *model.Payroll) *PayrollResponse {
	if data == nil {
		return nil
	}

	return &PayrollResponse{
		ID:                  data.ID,
		EmployeeID:          data.EmployeeID,
		BasicSalary:         data.BasicSalary,
		PositionAllowance:   data.PositionAllowance,
		OtherAllowance:      data.OtherAllowance,
		OvertimeRate:        data.OvertimeRate,
		LoanDeduction:       data.LoanDeduction,
		AttendanceDeduction: data.AttendanceDeduction,
		IncomeTax:           data.IncomeTax,
		AdditionalData:      data.AdditionalData,
		NetSalary:           data.NetSalary,
		PDFPath:             data.PDFPath,
		Status:              string(data.Status),
		SentAt:              data.SentAt,
		CreatedAt:           data.CreatedAt,
		UpdatedAt:           data.UpdatedAt,
	}
}

func ToPayrollListResponse(data []model.Payroll) []PayrollResponse {
	response := make([]PayrollResponse, 0, len(data))
	for _, item := range data {
		response = append(response, *ToPayrollResponse(&item))
	}

	return response
}
