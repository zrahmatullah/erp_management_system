package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type PayrollRepository struct {
	db *pgxpool.Pool
}

func NewPayrollRepository(db *pgxpool.Pool) domain.PayrollRepository {
	return &PayrollRepository{db: db}
}

func (r *PayrollRepository) CreatePeriod(ctx context.Context, p *domain.PayrollPeriod) error {
	query := `
		INSERT INTO payroll_periods (id, name, start_date, end_date, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at`

	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		p.ID, p.Name, p.StartDate, p.EndDate, p.Status, p.CreatedBy,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *PayrollRepository) CreatePayslip(ctx context.Context, ps *domain.Payslip) error {
	query := `
		INSERT INTO payslips (id, payroll_period_id, employee_id, basic_salary, total_allowances, total_overtime, total_bonus, gross_income, total_deductions, tax_amount, net_salary, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
		RETURNING created_at, updated_at`

	if ps.ID == uuid.Nil {
		ps.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		ps.ID, ps.PayrollPeriodID, ps.EmployeeID, ps.BasicSalary, ps.TotalAllowances, ps.TotalOvertime, ps.TotalBonus, ps.GrossIncome, ps.TotalDeductions, ps.TaxAmount, ps.NetSalary, ps.Status,
	).Scan(&ps.CreatedAt, &ps.UpdatedAt)
}
