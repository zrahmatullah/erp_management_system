package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) domain.EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(ctx context.Context, e *domain.Employee) error {
	query := `
		INSERT INTO employees (id, user_id, employee_code, first_name, last_name, dob, gender, national_id, tax_id, marital_status, address, phone, email, bank_name, bank_account, bank_account_name, join_date, status, department_id, position_id, branch_id, photo, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, NOW(), NOW())
		RETURNING created_at, updated_at`

	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		e.ID, e.UserID, e.EmployeeCode, e.FirstName, e.LastName, e.DOB, e.Gender, e.NationalID, e.TaxID, e.MaritalStatus, e.Address, e.Phone, e.Email, e.BankName, e.BankAccount, e.BankAccountName, e.JoinDate, e.Status, e.DepartmentID, e.PositionID, e.BranchID, e.Photo,
	).Scan(&e.CreatedAt, &e.UpdatedAt)
}

func (r *EmployeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	query := `SELECT id, employee_code, first_name, last_name, phone, email, status, department_id, position_id, branch_id FROM employees WHERE id = $1 AND deleted_at IS NULL`
	e := &domain.Employee{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.EmployeeCode, &e.FirstName, &e.LastName, &e.Phone, &e.Email, &e.Status, &e.DepartmentID, &e.PositionID, &e.BranchID,
	)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *EmployeeRepository) List(ctx context.Context, branchID uuid.UUID) ([]domain.Employee, error) {
	query := `SELECT id, employee_code, first_name, last_name, phone, email, status FROM employees WHERE branch_id = $1 AND deleted_at IS NULL ORDER BY first_name ASC`
	rows, err := r.db.Query(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emps []domain.Employee
	for rows.Next() {
		var e domain.Employee
		if err := rows.Scan(&e.ID, &e.EmployeeCode, &e.FirstName, &e.LastName, &e.Phone, &e.Email, &e.Status); err != nil {
			return nil, err
		}
		emps = append(emps, e)
	}
	return emps, nil
}
