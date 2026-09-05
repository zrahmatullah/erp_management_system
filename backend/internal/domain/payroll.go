package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PayrollStatus string
const (
	PayrollDraft           PayrollStatus = "draft"
	PayrollCalculating     PayrollStatus = "calculating"
	PayrollReview          PayrollStatus = "review"
	PayrollManagerApproved PayrollStatus = "manager_approved"
	PayrollFinanceApproved PayrollStatus = "finance_approved"
	PayrollProcessed       PayrollStatus = "processed"
	PayrollPaid            PayrollStatus = "paid"
)

type PayrollPeriod struct {
	BaseEntity
	Name       string        `json:"name"`
	StartDate  time.Time     `json:"start_date"`
	EndDate    time.Time     `json:"end_date"`
	Status     PayrollStatus `json:"status"`
	CreatedBy  uuid.UUID     `json:"created_by"`
	ApprovedBy *uuid.UUID    `json:"approved_by"`
}

type Payslip struct {
	BaseEntity
	PayrollPeriodID uuid.UUID `json:"payroll_period_id"`
	EmployeeID      uuid.UUID `json:"employee_id"`
	BasicSalary     float64   `json:"basic_salary"`
	TotalAllowances float64   `json:"total_allowances"`
	TotalOvertime   float64   `json:"total_overtime"`
	TotalBonus      float64   `json:"total_bonus"`
	GrossIncome     float64   `json:"gross_income"`
	TotalDeductions float64   `json:"total_deductions"`
	TaxAmount       float64   `json:"tax_amount"`
	NetSalary       float64   `json:"net_salary"`
	Status          string    `json:"status"`
}

type PayrollItem struct {
	BaseEntity
	PayslipID   uuid.UUID `json:"payslip_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
}

type PayrollDeduction struct {
	BaseEntity
	PayslipID   uuid.UUID `json:"payslip_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
}

type PayrollRepository interface {
	CreatePeriod(ctx context.Context, period *PayrollPeriod) error
	CreatePayslip(ctx context.Context, payslip *Payslip) error
}

type PayrollUsecase interface {
	GeneratePayroll(ctx context.Context, name string, start, end time.Time) (PayrollPeriod, error)
}
