package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EmployeeStatus string
const (
	EmpActive     EmployeeStatus = "active"
	EmpProbation  EmployeeStatus = "probation"
	EmpContract   EmployeeStatus = "contract"
	EmpResigned   EmployeeStatus = "resigned"
	EmpTerminated EmployeeStatus = "terminated"
)

type Employee struct {
	BaseEntity
	UserID          *uuid.UUID     `json:"user_id"`
	EmployeeCode    string         `json:"employee_code"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	DOB             time.Time      `json:"dob"`
	Gender          string         `json:"gender"`
	NationalID      string         `json:"national_id"`
	TaxID           string         `json:"tax_id"`
	MaritalStatus   string         `json:"marital_status"`
	Address         string         `json:"address"`
	Phone           string         `json:"phone"`
	Email           string         `json:"email"`
	BankName        string         `json:"bank_name"`
	BankAccount     string         `json:"bank_account"`
	BankAccountName string         `json:"bank_account_name"`
	JoinDate        time.Time      `json:"join_date"`
	EndDate         *time.Time     `json:"end_date"`
	Status          EmployeeStatus `json:"status"`
	DepartmentID    uuid.UUID      `json:"department_id"`
	PositionID      uuid.UUID      `json:"position_id"`
	BranchID        uuid.UUID      `json:"branch_id"`
	Photo           string         `json:"photo"`
}

type Department struct {
	BaseEntity
	Name        string     `json:"name"`
	ManagerID   *uuid.UUID `json:"manager_id"`
	Description string     `json:"description"`
}

type Position struct {
	BaseEntity
	DepartmentID uuid.UUID `json:"department_id"`
	Title        string    `json:"title"`
	Level        int       `json:"level"`
	BaseSalary   float64   `json:"base_salary"`
}

type EmploymentHistory struct {
	BaseEntity
	EmployeeID   uuid.UUID  `json:"employee_id"`
	PositionID   uuid.UUID  `json:"position_id"`
	DepartmentID uuid.UUID  `json:"department_id"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Status       string     `json:"status"`
	Notes        string     `json:"notes"`
}

type EmployeeRepository interface {
	Create(ctx context.Context, emp *Employee) error
	GetByID(ctx context.Context, id uuid.UUID) (*Employee, error)
	List(ctx context.Context, branchID uuid.UUID) ([]Employee, error)
}

type EmployeeUsecase interface {
	Create(ctx context.Context, emp Employee) (Employee, error)
	GetByID(ctx context.Context, id uuid.UUID) (Employee, error)
}
