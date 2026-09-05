package hris

import (
	"context"
	"time"

	"cafe-erp-system/backend/internal/domain"
)

type PayrollUsecaseImpl struct {
	repo domain.PayrollRepository
}

func NewPayrollUsecase(repo domain.PayrollRepository) domain.PayrollUsecase {
	return &PayrollUsecaseImpl{repo: repo}
}

func (u *PayrollUsecaseImpl) GeneratePayroll(ctx context.Context, name string, start, end time.Time) (domain.PayrollPeriod, error) {
	period := domain.PayrollPeriod{
		Name:      name,
		StartDate: start,
		EndDate:   end,
		Status:    domain.PayrollDraft,
	}

	err := u.repo.CreatePeriod(ctx, &period)
	return period, err
}
