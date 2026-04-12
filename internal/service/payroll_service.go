package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PayrollService interface {
	GetAll(ctx context.Context) ([]model.Payroll, error)
	GetByID(ctx context.Context, id string) (*model.Payroll, error)
	GenerateOne(ctx context.Context, employeeID string) (*model.Payroll, error)
	GenerateAll(ctx context.Context) ([]model.Payroll, error)
	Update(ctx context.Context, id string, payload request.PayrollUpdateRequest) (*model.Payroll, error)
}

type payrollService struct {
	payrollRepo        repository.PayrollRepository
	payrollSettingRepo repository.PayrollSettingRepository
	userRepo           repository.UserRepository
}

func NewPayrollService(payrollRepo repository.PayrollRepository, payrollSettingRepo repository.PayrollSettingRepository, userRepo repository.UserRepository) PayrollService {
	return &payrollService{payrollRepo: payrollRepo, payrollSettingRepo: payrollSettingRepo, userRepo: userRepo}
}

func (s *payrollService) generatePayrollPDF(
	ctx context.Context,
	payroll *model.Payroll,
	employee *model.User,
) (*string, error) {
	_ = ctx

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Payroll Slip")
	pdf.Ln(12)

	// Employee info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(50, 8, "Employee Name")
	pdf.Cell(5, 8, ":")
	pdf.Cell(0, 8, employee.FullName)
	pdf.Ln(8)

	pdf.Cell(50, 8, "Employee ID")
	pdf.Cell(5, 8, ":")
	pdf.Cell(0, 8, employee.ID)
	pdf.Ln(8)

	pdf.Cell(50, 8, "Payroll Date")
	pdf.Cell(5, 8, ":")
	pdf.Cell(0, 8, payroll.CreatedAt.Format("2006-01-02"))
	pdf.Ln(12)

	// Salary details
	pdf.SetFont("Arial", "B", 13)
	pdf.Cell(0, 8, "Salary Detail")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)

	addMoneyRow := func(label string, value float64) {
		pdf.Cell(70, 8, label)
		pdf.Cell(10, 8, ":")
		pdf.Cell(0, 8, fmt.Sprintf("Rp %.2f", value))
		pdf.Ln(8)
	}

	addMoneyRow("Basic Salary", payroll.BasicSalary)
	addMoneyRow("Position Allowance", payroll.PositionAllowance)
	addMoneyRow("Other Allowance", payroll.OtherAllowance)
	addMoneyRow("Overtime Rate", payroll.OvertimeRate)

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 12)
	addMoneyRow("Net Salary", payroll.NetSalary)

	// Output dir
	dir := "storage/payroll_pdfs"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("payroll_%s_%d.pdf", payroll.EmployeeID, time.Now().Unix())
	fullPath := filepath.Join(dir, filename)

	if err := pdf.OutputFileAndClose(fullPath); err != nil {
		return nil, err
	}

	return &fullPath, nil
}

func (s *payrollService) GetAll(ctx context.Context) ([]model.Payroll, error) {
	return s.payrollRepo.GetAll(ctx)
}

func (s *payrollService) GetByID(ctx context.Context, id string) (*model.Payroll, error) {
	return s.payrollRepo.GetByID(ctx, id)
}

func (s *payrollService) GenerateOne(ctx context.Context, employeeID string) (*model.Payroll, error) {
	employee, err := s.userRepo.GetByID(ctx, employeeID, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("User not found")
		}
		return nil, err
	}
	log.Print(employee)
	if employee.Profile == nil {
		return nil, common.BadRequestError("User profile is incomplete")
	}

	basicSalary := 0.0
	if employee.Profile.BasicSalary != nil {
		basicSalary = *employee.Profile.BasicSalary
	}

	positionAllowance := 0.0
	if employee.Profile.PositionAllowance != nil {
		positionAllowance = *employee.Profile.PositionAllowance
	}

	otherAllowance := 0.0
	if employee.Profile.OtherAllowance != nil {
		otherAllowance = *employee.Profile.OtherAllowance
	}

	payrollSetting, err := s.payrollSettingRepo.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	additionalData := make(map[string]interface{})
	overtimeRate := 0.0

	for _, setting := range payrollSetting {
		if setting.ConfigKey == "hourly_overtime_rate" {
			overtimeRate = setting.Value * 2
		} else {
			additionalData[setting.ConfigKey] = setting.Value
		}
	}

	dataBytes, err := json.Marshal(additionalData)
	if err != nil {
		return nil, err
	}

	payroll := &model.Payroll{
		EmployeeID:        employee.ID,
		BasicSalary:       basicSalary,
		PositionAllowance: positionAllowance,
		OtherAllowance:    otherAllowance,
		OvertimeRate:      overtimeRate,
		Status:            model.PayrollStatusUnsent,
		NetSalary:         basicSalary + positionAllowance + otherAllowance + overtimeRate,
		AdditionalData:    datatypes.JSON(dataBytes),
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	localNow := time.Now().In(loc)
	year, month, _ := localNow.Date()
	startLocal := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	endLocal := startLocal.AddDate(0, 1, 0)

	existing, err := s.payrollRepo.GetByEmployeeIDAndCreatedAtRange(ctx, employee.ID, startLocal.UTC(), endLocal.UTC())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var saved *model.Payroll

	if err == nil && existing != nil {
		payroll.ID = existing.ID
		payroll.Status = existing.Status
		payroll.SentAt = existing.SentAt

		saved, err = s.payrollRepo.Update(ctx, payroll)
		if err != nil {
			return nil, err
		}
	} else {
		saved, err = s.payrollRepo.GenerateOne(ctx, payroll)
		if err != nil {
			return nil, err
		}
	}

	pdfPath, err := s.generatePayrollPDF(ctx, saved, employee)
	if err != nil {
		return nil, err
	}

	saved.PDFPath = pdfPath
	saved, err = s.payrollRepo.Update(ctx, saved)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (s *payrollService) GenerateAll(ctx context.Context) ([]model.Payroll, error) {
	employees, err := s.userRepo.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	payrollSetting, err := s.payrollSettingRepo.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	additionalData := make(map[string]interface{})
	overtimeRate := 0.0
	for _, setting := range payrollSetting {
		if setting.ConfigKey == "hourly_overtime_rate" {
			overtimeRate = setting.Value * 2
		} else {
			additionalData[setting.ConfigKey] = setting.Value
		}
	}

	dataBytes, err := json.Marshal(additionalData)
	if err != nil {
		return nil, err
	}
	additionalJSON := datatypes.JSON(dataBytes)

	loc, _ := time.LoadLocation("Asia/Jakarta") // WIB
	localNow := time.Now().In(loc)
	year, month, _ := localNow.Date()
	startLocal := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	endLocal := startLocal.AddDate(0, 1, 0)
	startUTC := startLocal.UTC()
	endUTC := endLocal.UTC()

	employeeIDs := make([]string, 0, len(employees))
	for _, employee := range employees {
		employeeIDs = append(employeeIDs, employee.ID)
	}

	existingPayrolls, err := s.payrollRepo.GetByEmployeeIDsAndCreatedAtRange(ctx, employeeIDs, startUTC, endUTC)
	if err != nil {
		return nil, err
	}
	existingByEmployeeID := make(map[string]model.Payroll, len(existingPayrolls))
	for _, p := range existingPayrolls {
		existingByEmployeeID[p.EmployeeID] = p
	}

	toCreate := make([]model.Payroll, 0)
	toUpdate := make([]*model.Payroll, 0)

	for _, employee := range employees {
		log.Print(employee.Profile)
		if employee.Profile == nil {
			return nil, common.BadRequestError("User profile is incomplete")
		}

		basicSalary := 0.0
		if employee.Profile.BasicSalary != nil {
			basicSalary = *employee.Profile.BasicSalary
		}
		positionAllowance := 0.0
		if employee.Profile.PositionAllowance != nil {
			positionAllowance = *employee.Profile.PositionAllowance
		}
		otherAllowance := 0.0
		if employee.Profile.OtherAllowance != nil {
			otherAllowance = *employee.Profile.OtherAllowance
		}

		payroll := &model.Payroll{
			EmployeeID:        employee.ID,
			BasicSalary:       basicSalary,
			PositionAllowance: positionAllowance,
			OtherAllowance:    otherAllowance,
			OvertimeRate:      overtimeRate,
			Status:            model.PayrollStatusUnsent,
			NetSalary:         basicSalary + positionAllowance + otherAllowance + overtimeRate,
			AdditionalData:    additionalJSON,
		}

		if existing, ok := existingByEmployeeID[employee.ID]; ok {
			payroll.ID = existing.ID
			payroll.Status = existing.Status
			payroll.SentAt = existing.SentAt
			toUpdate = append(toUpdate, payroll)
			continue
		}

		toCreate = append(toCreate, *payroll)
	}

	if err := s.payrollRepo.GenerateMany(ctx, toCreate); err != nil {
		return nil, err
	}
	for _, payroll := range toUpdate {
		if _, err := s.payrollRepo.Update(ctx, payroll); err != nil {
			return nil, err
		}
	}

	return s.payrollRepo.GenerateAll(ctx)
}

func (s *payrollService) Update(ctx context.Context, id string, payload request.PayrollUpdateRequest) (*model.Payroll, error) {
	existing, err := s.payrollRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Payroll not found")
		}
		return nil, err
	}

	if payload.BasicSalary != nil {
		existing.BasicSalary = *payload.BasicSalary
	}
	if payload.PositionAllowance != nil {
		existing.PositionAllowance = *payload.PositionAllowance
	}
	if payload.OtherAllowance != nil {
		existing.OtherAllowance = *payload.OtherAllowance
	}
	if payload.OvertimeRate != nil {
		existing.OvertimeRate = *payload.OvertimeRate
	}
	if payload.LoanDeduction != nil {
		existing.LoanDeduction = *payload.LoanDeduction
	}
	if payload.AttendanceDeduction != nil {
		existing.AttendanceDeduction = *payload.AttendanceDeduction
	}
	if payload.IncomeTax != nil {
		existing.IncomeTax = *payload.IncomeTax
	}

	if payload.AdditionalData != nil {
		raw := strings.TrimSpace(*payload.AdditionalData)
		if raw == "" {
			existing.AdditionalData = datatypes.JSON([]byte(`{}`))
		} else {
			b := []byte(raw)
			if !json.Valid(b) {
				return nil, common.BadRequestError("additional_data must be valid JSON")
			}
			existing.AdditionalData = datatypes.JSON(b)
		}
	}

	// Recalculate net salary after updates.
	gross := existing.BasicSalary + existing.PositionAllowance + existing.OtherAllowance + existing.OvertimeRate
	existing.NetSalary = gross - existing.LoanDeduction - existing.AttendanceDeduction - existing.IncomeTax

	updated, err := s.payrollRepo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Payroll not found")
		}
		return nil, err
	}
	return updated, nil
}
