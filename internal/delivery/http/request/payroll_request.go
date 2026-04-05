package request

type PayrollUpdateRequest struct {
	BasicSalary         *float64 `json:"basic_salary" binding:"omitempty"`
	PositionAllowance   *float64 `json:"position_allowance" binding:"omitempty"`
	OtherAllowance      *float64 `json:"other_allowance" binding:"omitempty"`
	OvertimeRate        *float64 `json:"overtime_rate" binding:"omitempty"`
	LoanDeduction       *float64 `json:"loan_deduction" binding:"omitempty"`
	AttendanceDeduction *float64 `json:"attendance_deduction" binding:"omitempty"`
	IncomeTax           *float64 `json:"income_tax" binding:"omitempty"`
	AdditionalData      *string  `json:"additional_data" binding:"omitempty"`
}
