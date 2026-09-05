package finance

import (
	"context"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
)

type FinanceUsecaseImpl struct {
	repo domain.FinanceRepository
}

func NewFinanceUsecase(repo domain.FinanceRepository) domain.FinanceUsecase {
	return &FinanceUsecaseImpl{repo: repo}
}

func (u *FinanceUsecaseImpl) SubmitExpense(ctx context.Context, req domain.Expense) (domain.Expense, error) {
	req.Status = domain.ExpSubmitted
	err := u.repo.CreateExpense(ctx, &req)
	return req, err
}

func (u *FinanceUsecaseImpl) ApproveExpense(ctx context.Context, expenseID, approverID uuid.UUID, status domain.ExpenseStatus) error {
	// Omitting for brevity
	return nil
}
