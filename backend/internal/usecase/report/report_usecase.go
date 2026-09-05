package report

import (
	"context"
	"time"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/repository/postgres"
)

type ReportUsecaseImpl struct {
	repo *postgres.ReportRepository
}

func NewReportUsecase(repo *postgres.ReportRepository) *ReportUsecaseImpl {
	return &ReportUsecaseImpl{repo: repo}
}

func (u *ReportUsecaseImpl) GetSalesReport(ctx context.Context, branchID uuid.UUID, start, end time.Time) (float64, error) {
	return u.repo.GetSalesReport(ctx, branchID, start, end)
}
