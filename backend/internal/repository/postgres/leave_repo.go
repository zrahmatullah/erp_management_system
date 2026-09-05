package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type LeaveRepository struct {
	db *pgxpool.Pool
}

func NewLeaveRepository(db *pgxpool.Pool) domain.LeaveRepository {
	return &LeaveRepository{db: db}
}

func (r *LeaveRepository) CreateRequest(ctx context.Context, req *domain.LeaveRequest) error {
	query := `
		INSERT INTO leave_requests (id, employee_id, leave_type_id, start_date, end_date, total_days, reason, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING created_at, updated_at`

	if req.ID == uuid.Nil {
		req.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		req.ID, req.EmployeeID, req.LeaveTypeID, req.StartDate, req.EndDate, req.TotalDays, req.Reason, req.Status,
	).Scan(&req.CreatedAt, &req.UpdatedAt)
}

func (r *LeaveRepository) UpdateRequest(ctx context.Context, req *domain.LeaveRequest) error {
	query := `UPDATE leave_requests SET status = $1, approved_by = $2, approver_notes = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.db.Exec(ctx, query, req.Status, req.ApprovedBy, req.ApproverNotes, req.ID)
	return err
}

func (r *LeaveRepository) GetBalance(ctx context.Context, empID uuid.UUID, year int) ([]domain.LeaveBalance, error) {
	query := `SELECT id, employee_id, leave_type_id, year, allocated, used, remaining FROM leave_balances WHERE employee_id = $1 AND year = $2`
	rows, err := r.db.Query(ctx, query, empID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []domain.LeaveBalance
	for rows.Next() {
		var b domain.LeaveBalance
		if err := rows.Scan(&b.ID, &b.EmployeeID, &b.LeaveTypeID, &b.Year, &b.Allocated, &b.Used, &b.Remaining); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}
