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

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/request"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/response/common"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/internal/repository"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PayrollService interface {
	GetAll(ctx context.Context) ([]model.Payroll, error)
	GetMyPayrolls(ctx context.Context, userID string) ([]model.Payroll, error)
	GetByID(ctx context.Context, id string) (*model.Payroll, error)
	GenerateOne(ctx context.Context, employeeID string) (*model.Payroll, error)
	GenerateAll(ctx context.Context) ([]model.Payroll, error)
	Update(ctx context.Context, id string, payload request.PayrollUpdateRequest) (*model.Payroll, error)
	SendPayroll(ctx context.Context, id string) (*model.Payroll, error)
}

type payrollService struct {
	payrollRepo        repository.PayrollRepository
	payrollSettingRepo repository.PayrollSettingRepository
	userRepo           repository.UserRepository
}

func NewPayrollService(payrollRepo repository.PayrollRepository, payrollSettingRepo repository.PayrollSettingRepository, userRepo repository.UserRepository) PayrollService {
	return &payrollService{payrollRepo: payrollRepo, payrollSettingRepo: payrollSettingRepo, userRepo: userRepo}
}

func bpjsDeduction(additionalData datatypes.JSON) float64 {
	var extra struct {
		BPJSHealth float64 `json:"bpjs_health_rate"`
		BPJSJHT    float64 `json:"bpjs_employment_jht_rate"`
		BPJSJP     float64 `json:"bpjs_employment_jp_rate"`
	}
	_ = json.Unmarshal(additionalData, &extra)
	return extra.BPJSHealth + extra.BPJSJHT + extra.BPJSJP
}

func formatIDR(amount float64) string {
	intPart := int64(amount)
	s := fmt.Sprintf("%d", intPart)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	return "Rp " + result
}

func slipNum(amount float64) string {
	if amount == 0 {
		return "-"
	}
	intPart := int64(amount)
	s := fmt.Sprintf("%d", intPart)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	return result
}

func (s *payrollService) generatePayrollPDF(
	payroll *model.Payroll,
	employee *model.User,
) (*string, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t := time.Now().In(loc)

	var extra struct {
		BPJSHealth float64 `json:"bpjs_health_rate"`
		BPJSJHT    float64 `json:"bpjs_employment_jht_rate"`
		BPJSJP     float64 `json:"bpjs_employment_jp_rate"`
	}
	_ = json.Unmarshal(payroll.AdditionalData, &extra)

	period := payroll.CreatedAt.In(loc).Format("2006/01")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// ── HEADER ──────────────────────────────────────────────
	pdf.Image("assets/images/logo.jpg", 15, 10, 22, 0, false, "", 0, "")

	pdf.SetFont("Arial", "B", 11)
	pdf.SetXY(40, 13)
	pdf.Cell(60, 6, "PT. King Royal Sejahtera")
	pdf.SetFont("Arial", "B", 22)
	pdf.SetXY(100, 12)
	pdf.CellFormat(95, 10, "Slip Gaji", "", 0, "R", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(40, 20)
	pdf.Cell(60, 5, "Ahmad Yani No 134 Brebes")
	pdf.SetXY(40, 25)
	pdf.Cell(60, 5, "Brebes, Indonesia")

	// ── EMPLOYEE INFO ────────────────────────────────────────
	pdf.SetDrawColor(0, 0, 0)
	pdf.Line(15, 36, 195, 36)

	nameNIK := employee.FullName
	deptJob := ""
	joinedAt := "-"
	bankInfo := "-"

	if employee.Profile != nil {
		if employee.Profile.EmployeeCode != nil {
			nameNIK += " (" + *employee.Profile.EmployeeCode + ")"
		}
		parts := []string{}
		if employee.Profile.Department != nil {
			parts = append(parts, *employee.Profile.Department)
		}
		if employee.Profile.Position != nil {
			parts = append(parts, *employee.Profile.Position)
		}
		if len(parts) > 0 {
			deptJob = parts[0]
			if len(parts) > 1 {
				deptJob += " / " + parts[1]
			}
		}
		if employee.Profile.JoinedAt != nil {
			joinedAt = employee.Profile.JoinedAt.In(loc).Format("02 January 2006")
		}
		if employee.Profile.BankAccountNumber != nil {
			bankInfo = *employee.Profile.BankAccountNumber + " (" + employee.FullName + ")"
		}
	}

	pdf.SetFont("Arial", "", 9)
	infoY := 39.0
	rowH := 5.5
	infoLeft := func(label, val string, y float64) {
		pdf.SetXY(15, y)
		pdf.CellFormat(42, rowH, label, "", 0, "L", false, 0, "")
		pdf.SetXY(57, y)
		pdf.CellFormat(65, rowH, val, "", 0, "L", false, 0, "")
	}
	infoLeft("Nama / NIK", nameNIK, infoY)
	pdf.SetXY(130, infoY)
	pdf.CellFormat(25, rowH, "Periode Gaji", "", 0, "L", false, 0, "")
	pdf.SetXY(162, infoY)
	pdf.CellFormat(33, rowH, period, "", 0, "L", false, 0, "")

	infoLeft("Dept / Jabatan", deptJob, infoY+rowH)
	infoLeft("Tgl Mulai Bekerja", joinedAt, infoY+rowH*2)

	tableY := infoY + rowH*3 + 4
	pdf.Line(15, tableY, 195, tableY)
	tableY += 3

	// ── INCOME / DEDUCTIONS TABLE ────────────────────────────
	// Left col : label x=15..72, value x=72..100 (right-aligned)
	// Divider  : x=101
	// Right col: label x=103..160, value x=160..193 (right-aligned)
	const (
		lLblX   = 15.0
		lLblW   = 57.0
		lValX   = 72.0
		lValW   = 28.0 // ends at 100
		divLine = 101.0
		rLblX   = 103.0
		rLblW   = 57.0
		rValX   = 160.0
		rValW   = 33.0 // ends at 193
		tblRowH = 6.5
	)

	// Header row
	pdf.SetFillColor(230, 230, 230)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(15, tableY)
	pdf.CellFormat(86, 7, "  Pendapatan", "1", 0, "L", true, 0, "")
	pdf.SetXY(103, tableY)
	pdf.CellFormat(92, 7, "  Potongan", "1", 0, "L", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	tableY += 8

	type slipRow struct {
		iLabel string
		iVal   float64
		dLabel string
		dVal   float64
	}
	rows := []slipRow{
		{"Gaji Pokok", payroll.BasicSalary, "Potongan Pinjaman Karyawan", payroll.LoanDeduction},
		{"Lembur", payroll.OvertimeRate, "Potongan Absen", payroll.AttendanceDeduction},
		{"Tunjangan Jabatan", payroll.PositionAllowance, "Potongan BPJS Kesehatan", extra.BPJSHealth},
		{"Tunjangan Lain", payroll.OtherAllowance, "Potongan BPJS TK JHT", extra.BPJSJHT},
		{"", 0, "Potongan BPJS TK JP", extra.BPJSJP},
		{"", 0, "Potongan Pajak PPh21", payroll.IncomeTax},
	}

	startTableY := tableY
	pdf.SetFont("Arial", "", 9)
	for _, r := range rows {
		if r.iLabel != "" {
			pdf.SetXY(lLblX, tableY)
			pdf.CellFormat(lLblW, tblRowH, r.iLabel, "", 0, "L", false, 0, "")
			pdf.SetXY(lValX, tableY)
			pdf.CellFormat(lValW, tblRowH, slipNum(r.iVal), "", 0, "R", false, 0, "")
		}
		pdf.SetXY(rLblX, tableY)
		pdf.CellFormat(rLblW, tblRowH, r.dLabel, "", 0, "L", false, 0, "")
		pdf.SetXY(rValX, tableY)
		pdf.CellFormat(rValW, tblRowH, slipNum(r.dVal), "", 0, "R", false, 0, "")
		tableY += tblRowH
	}

	// Vertical divider between columns
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(divLine, startTableY, divLine, tableY)

	// Separator + total row
	tableY += 1
	pdf.SetDrawColor(0, 0, 0)
	pdf.Line(15, tableY, 195, tableY)
	tableY += 2

	totalIncome := payroll.BasicSalary + payroll.PositionAllowance + payroll.OtherAllowance + payroll.OvertimeRate
	totalDeduct := payroll.LoanDeduction + payroll.AttendanceDeduction + extra.BPJSHealth + extra.BPJSJHT + extra.BPJSJP + payroll.IncomeTax

	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(lLblX, tableY)
	pdf.CellFormat(lLblW, 7, "Total Pendapatan", "", 0, "L", false, 0, "")
	pdf.SetXY(lValX, tableY)
	pdf.CellFormat(lValW, 7, slipNum(totalIncome), "", 0, "R", false, 0, "")
	pdf.SetXY(rLblX, tableY)
	pdf.CellFormat(rLblW, 7, "Total Potongan", "", 0, "L", false, 0, "")
	pdf.SetXY(rValX, tableY)
	pdf.CellFormat(rValW, 7, slipNum(totalDeduct), "", 0, "R", false, 0, "")
	tableY += 10

	// ── PAYMENT INFO + NET BOX ───────────────────────────────
	netPenerimaan := totalIncome - totalDeduct
	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(15, tableY)
	pdf.MultiCell(85, 5, "Pembayaran gaji telah dilakukan oleh perusahaan\nSecara transfer ke rekening karyawan\n"+bankInfo, "", "L", false)

	boxX := 110.0
	boxY := tableY
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(boxX, boxY)
	pdf.CellFormat(80, 7, "Total Penerimaan Bulan Ini", "TLR", 0, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 20)
	pdf.SetXY(boxX, boxY+7)
	pdf.CellFormat(80, 14, slipNum(netPenerimaan), "LRB", 0, "C", false, 0, "")

	// ── FOOTER ───────────────────────────────────────────────
	footerY := tableY + 30
	fX := 130.0
	fW := 65.0 // ends at 195
	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(fX, footerY)
	pdf.CellFormat(fW, 5, "Brebes, "+t.Format("02 January 2006"), "", 0, "R", false, 0, "")
	pdf.SetXY(fX, footerY+6)
	pdf.CellFormat(fW, 5, "HR-Chief Accounting", "", 0, "R", false, 0, "")
	pdf.Image("assets/images/ttd.png", fX+fW-32, footerY+12, 30, 0, false, "", 0, "")
	pdf.SetXY(fX, footerY+26)
	pdf.CellFormat(fW, 5, "Rizki Aenun Nisa, S.M", "", 0, "R", false, 0, "")

	// ── OUTPUT ───────────────────────────────────────────────
	dir := "files/payrolls"
	yearMonth := t.Format("2006-01")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s_%s.pdf", payroll.EmployeeID, yearMonth)
	fullPath := filepath.Join(dir, filename)

	if err := pdf.OutputFileAndClose(fullPath); err != nil {
		return nil, err
	}

	urlPath := dir + "/" + filename
	return &urlPath, nil
}

func (s *payrollService) GetAll(ctx context.Context) ([]model.Payroll, error) {
	return s.payrollRepo.GetAll(ctx)
}

func (s *payrollService) GetMyPayrolls(ctx context.Context, userID string) ([]model.Payroll, error) {
	return s.payrollRepo.GetByEmployeeID(ctx, userID)
}

func (s *payrollService) GetByID(ctx context.Context, id string) (*model.Payroll, error) {
	payroll, err := s.payrollRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Payroll not found")
		}
		return nil, err
	}
	return payroll, nil
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

	gross := basicSalary + positionAllowance + otherAllowance + overtimeRate
	payroll := &model.Payroll{
		EmployeeID:        employee.ID,
		BasicSalary:       basicSalary,
		PositionAllowance: positionAllowance,
		OtherAllowance:    otherAllowance,
		OvertimeRate:      overtimeRate,
		Status:            model.PayrollStatusUnsent,
		GrossSalary:       gross,
		NetSalary:         gross - bpjsDeduction(datatypes.JSON(dataBytes)),
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
		payroll.CreatedAt = existing.CreatedAt

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

	pdfPath, err := s.generatePayrollPDF(saved, employee)
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
	employees, err := s.userRepo.GetAll(ctx, true, nil)
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
	employeeMap := make(map[string]*model.User, len(employees))

	for i := range employees {
		employeeIDs = append(employeeIDs, employees[i].ID)
		employeeMap[employees[i].ID] = &employees[i]
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

	for i := range employees {
		employee := employees[i]
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

		empGross := basicSalary + positionAllowance + otherAllowance + overtimeRate
		payroll := &model.Payroll{
			EmployeeID:        employee.ID,
			BasicSalary:       basicSalary,
			PositionAllowance: positionAllowance,
			OtherAllowance:    otherAllowance,
			OvertimeRate:      overtimeRate,
			Status:            model.PayrollStatusUnsent,
			GrossSalary:       empGross,
			NetSalary:         empGross - bpjsDeduction(additionalJSON),
			AdditionalData:    additionalJSON,
		}

		if existing, ok := existingByEmployeeID[employee.ID]; ok {
			payroll.ID = existing.ID
			payroll.Status = existing.Status
			payroll.SentAt = existing.SentAt
			payroll.PDFPath = existing.PDFPath
			payroll.CreatedAt = existing.CreatedAt
			payroll.UpdatedAt = existing.UpdatedAt
			toUpdate = append(toUpdate, payroll)
			continue
		}

		toCreate = append(toCreate, *payroll)
	}

	if len(toCreate) > 0 {
		if err := s.payrollRepo.GenerateMany(ctx, toCreate); err != nil {
			return nil, err
		}
	}
	if len(toUpdate) > 0 {
		for _, payroll := range toUpdate {
			if _, err := s.payrollRepo.Update(ctx, payroll); err != nil {
				return nil, err
			}
		}
	}

	savedPayrolls, err := s.payrollRepo.GetByEmployeeIDsAndCreatedAtRange(ctx, employeeIDs, startUTC, endUTC)
	if err != nil {
		return nil, err
	}

	for i := range savedPayrolls {
		employee, ok := employeeMap[savedPayrolls[i].EmployeeID]
		if !ok {
			return nil, common.NotFoundError("User not found for payroll with employee ID: " + savedPayrolls[i].EmployeeID)
		}

		pdfPath, err := s.generatePayrollPDF(&savedPayrolls[i], employee)
		if err != nil {
			return nil, err
		}

		savedPayrolls[i].PDFPath = pdfPath
		updated, err := s.payrollRepo.Update(ctx, &savedPayrolls[i])
		if err != nil {
			return nil, err
		}

		savedPayrolls[i] = *updated
	}

	return savedPayrolls, nil
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

	// Recalculate gross and net salary after updates.
	gross := existing.BasicSalary + existing.PositionAllowance + existing.OtherAllowance + existing.OvertimeRate
	existing.GrossSalary = gross
	existing.NetSalary = gross - existing.LoanDeduction - existing.AttendanceDeduction - existing.IncomeTax - bpjsDeduction(existing.AdditionalData)

	updated, err := s.payrollRepo.Update(ctx, existing)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Payroll not found")
		}
		return nil, err
	}
	return updated, nil
}

func (s *payrollService) SendPayroll(ctx context.Context, id string) (*model.Payroll, error) {
	payroll, err := s.payrollRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("Payroll not found")
		}
		return nil, err
	}

	employee, err := s.userRepo.GetByID(ctx, payroll.EmployeeID, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NotFoundError("User not found")
		}
		return nil, err
	}

	if employee.Profile == nil {
		return nil, common.BadRequestError("User profile is incomplete")
	}

	if employee.Email == "" {
		return nil, common.BadRequestError("User email is empty")
	}

	if !strings.Contains(employee.Email, "@") {
		return nil, common.BadRequestError("User email is not valid, please update user email first")
	}

	pdfPath, err := s.generatePayrollPDF(payroll, employee)
	if err != nil {
		return nil, err
	}
	payroll.PDFPath = pdfPath
	payroll, err = s.payrollRepo.Update(ctx, payroll)
	if err != nil {
		return nil, err
	}

	cleanPDFPath := strings.ReplaceAll(*payroll.PDFPath, "\\", "/")
	payrollUrl := fmt.Sprintf("%s/%s", config.GetEnv().ServerBaseURL, cleanPDFPath)
	err = utils.SendEmail(utils.EmailParams{
		FromName:   config.GetEnv().SMTPFromName,
		FromEmail:  config.GetEnv().SMTPFromEmail,
		Password:   config.GetEnv().SMTPPassword,
		Host:       config.GetEnv().SMTPHost,
		Port:       config.GetEnv().SMTPPort,
		Encryption: utils.EncryptionType(config.GetEnv().SMTPEncryption),
		ToName:     employee.FullName,
		ToEmail:    employee.Email,
		Subject:    fmt.Sprintf("Payroll Slip - %s", time.Now().In(time.FixedZone("WIB", 7*3600)).Format("January 2006")),
		Template:   "templates/payroll.html",
		Data: map[string]any{
			"EmployeeName": employee.FullName,
			"PayrollDate":  time.Now().In(time.FixedZone("WIB", 7*3600)).Format("02 January 2006"),
			"NetSalary":    payroll.NetSalary,
			"PayrollURL":   payrollUrl,
		},
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	payroll.Status = model.PayrollStatusSent
	payroll.SentAt = &now

	updated, err := s.payrollRepo.Update(ctx, payroll)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
