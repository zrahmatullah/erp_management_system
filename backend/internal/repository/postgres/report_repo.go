package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

// Basic implementation of report fetching
func (r *ReportRepository) GetSalesReport(ctx context.Context, branchID uuid.UUID, startDate, endDate time.Time) (float64, error) {
	query := `SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE branch_id = $1 AND status = 'completed' AND created_at >= $2 AND created_at <= $3 AND deleted_at IS NULL`
	var total float64
	err := r.db.QueryRow(ctx, query, branchID, startDate, endDate).Scan(&total)
	return total, err
}

func (r *ReportRepository) GetDashboardMetrics(ctx context.Context, branchID uuid.UUID) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	var dailySales float64
	err := r.db.QueryRow(ctx, `SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE branch_id = $1 AND status = 'completed' AND created_at >= $2`, branchID, startOfDay).Scan(&dailySales)
	if err != nil {
		return nil, err
	}
	
	metrics["daily_sales"] = dailySales
	return metrics, nil
}
