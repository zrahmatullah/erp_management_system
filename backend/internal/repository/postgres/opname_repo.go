package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type OpnameRepository struct {
	db *pgxpool.Pool
}

func NewOpnameRepository(db *pgxpool.Pool) *OpnameRepository {
	return &OpnameRepository{db: db}
}

func (r *OpnameRepository) Create(ctx context.Context, op *domain.StockOpname) error {
	query := `
		INSERT INTO stock_opname (id, warehouse_id, opname_number, status, scheduled_date, conducted_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at`

	if op.ID == uuid.Nil {
		op.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		op.ID, op.WarehouseID, op.OpnameNumber, op.Status, op.ScheduledDate, op.ConductedBy,
	).Scan(&op.CreatedAt, &op.UpdatedAt)
}
