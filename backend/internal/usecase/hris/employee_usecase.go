package hris

import (
	"context"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
)

type EmployeeUsecaseImpl struct {
	repo domain.EmployeeRepository
}

func NewEmployeeUsecase(repo domain.EmployeeRepository) domain.EmployeeUsecase {
	return &EmployeeUsecaseImpl{repo: repo}
}

func (u *EmployeeUsecaseImpl) Create(ctx context.Context, emp domain.Employee) (domain.Employee, error) {
	emp.Status = domain.EmpActive
	err := u.repo.Create(ctx, &emp)
	return emp, err
}

func (u *EmployeeUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (domain.Employee, error) {
	emp, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Employee{}, err
	}
	return *emp, nil
}
