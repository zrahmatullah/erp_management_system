package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type AttendanceRepository struct {
	db *pgxpool.Pool
}

func NewAttendanceRepository(db *pgxpool.Pool) domain.AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) ClockIn(ctx context.Context, att *domain.AttendanceRecord) error {
	query := `
		INSERT INTO attendance_records (id, employee_id, date, clock_in, geo_location_in, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at`

	if att.ID == uuid.Nil {
		att.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		att.ID, att.EmployeeID, att.Date, att.ClockIn, att.GeoLocationIn, att.Status,
	).Scan(&att.CreatedAt, &att.UpdatedAt)
}

func (r *AttendanceRepository) ClockOut(ctx context.Context, id uuid.UUID, outTime time.Time, geo string) error {
	query := `UPDATE attendance_records SET clock_out = $1, geo_location_out = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(ctx, query, outTime, geo, id)
	return err
}

func (r *AttendanceRepository) GetByDateRange(ctx context.Context, empID uuid.UUID, start, end time.Time) ([]domain.AttendanceRecord, error) {
	query := `SELECT id, employee_id, date, clock_in, clock_out, status, overtime_hours, late_minutes FROM attendance_records WHERE employee_id = $1 AND date >= $2 AND date <= $3 ORDER BY date ASC`
	rows, err := r.db.Query(ctx, query, empID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.AttendanceRecord
	for rows.Next() {
		var a domain.AttendanceRecord
		if err := rows.Scan(&a.ID, &a.EmployeeID, &a.Date, &a.ClockIn, &a.ClockOut, &a.Status, &a.OvertimeHours, &a.LateMinutes); err != nil {
			return nil, err
		}
		records = append(records, a)
	}
	return records, nil
}
