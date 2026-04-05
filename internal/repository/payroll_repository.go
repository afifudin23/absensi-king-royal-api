package repository

import (
	"context"
	"time"

	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"gorm.io/gorm"
)

type PayrollRepository interface {
	GetAll(ctx context.Context) ([]model.Payroll, error)
	GetByID(ctx context.Context, id string) (*model.Payroll, error)
	GetByEmployeeIDAndCreatedAtRange(ctx context.Context, employeeID string, start time.Time, end time.Time) (*model.Payroll, error)
	GetByEmployeeIDsAndCreatedAtRange(ctx context.Context, employeeIDs []string, start time.Time, end time.Time) ([]model.Payroll, error)
	GenerateOne(ctx context.Context, payroll *model.Payroll) (*model.Payroll, error)
	GenerateMany(ctx context.Context, payrolls []model.Payroll) error
	GenerateAll(ctx context.Context) ([]model.Payroll, error)
	Update(ctx context.Context, payroll *model.Payroll) (*model.Payroll, error)
}

type payrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &payrollRepository{db: db}
}

func (r *payrollRepository) GetAll(ctx context.Context) ([]model.Payroll, error) {
	var payrolls []model.Payroll
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&payrolls).Error
	return payrolls, err
}

func (r *payrollRepository) GetByID(ctx context.Context, id string) (*model.Payroll, error) {
	var payroll model.Payroll
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payroll).Error
	if err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (r *payrollRepository) GetByEmployeeIDAndCreatedAtRange(ctx context.Context, employeeID string, start time.Time, end time.Time) (*model.Payroll, error) {
	var payroll model.Payroll
	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND created_at >= ? AND created_at < ?", employeeID, start, end).
		Order("created_at asc").
		First(&payroll).
		Error
	if err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (r *payrollRepository) GetByEmployeeIDsAndCreatedAtRange(ctx context.Context, employeeIDs []string, start time.Time, end time.Time) ([]model.Payroll, error) {
	if len(employeeIDs) == 0 {
		return []model.Payroll{}, nil
	}

	var payrolls []model.Payroll
	err := r.db.WithContext(ctx).
		Where("employee_id IN ? AND created_at >= ? AND created_at < ?", employeeIDs, start, end).
		Order("created_at asc").
		Find(&payrolls).
		Error
	if err != nil {
		return nil, err
	}
	return payrolls, nil
}

func (r *payrollRepository) GenerateOne(ctx context.Context, payroll *model.Payroll) (*model.Payroll, error) {
	err := r.db.WithContext(ctx).Create(payroll).Error
	if err != nil {
		return nil, err
	}
	return payroll, nil
}

func (r *payrollRepository) GenerateMany(ctx context.Context, payrolls []model.Payroll) error {
	if len(payrolls) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&payrolls).Error
}

func (r *payrollRepository) GenerateAll(ctx context.Context) ([]model.Payroll, error) {
	var payrolls []model.Payroll
	err := r.db.WithContext(ctx).Find(&payrolls).Error
	if err != nil {
		return nil, err
	}
	return payrolls, nil
}

func (r *payrollRepository) Update(ctx context.Context, payroll *model.Payroll) (*model.Payroll, error) {
	result := r.db.WithContext(ctx).
		Model(&model.Payroll{}).
		Where("id = ?", payroll.ID).
		Updates(map[string]any{
			"employee_id":          payroll.EmployeeID,
			"basic_salary":         payroll.BasicSalary,
			"position_allowance":   payroll.PositionAllowance,
			"other_allowance":      payroll.OtherAllowance,
			"overtime_rate":        payroll.OvertimeRate,
			"loan_deduction":       payroll.LoanDeduction,
			"attendance_deduction": payroll.AttendanceDeduction,
			"income_tax":           payroll.IncomeTax,
			"additional_data":      payroll.AdditionalData,
			"net_salary":           payroll.NetSalary,
			"status":               payroll.Status,
			"sent_at":              payroll.SentAt,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return r.GetByID(ctx, payroll.ID)
}
