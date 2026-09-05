package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type PurchaseRepository struct {
	db *pgxpool.Pool
}

func NewPurchaseRepository(db *pgxpool.Pool) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) CreatePO(ctx context.Context, po *domain.PurchaseOrder) error {
	query := `
		INSERT INTO purchase_orders (id, po_number, branch_id, supplier_id, status, total_amount, expected_delivery, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at`

	if po.ID == uuid.Nil {
		po.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		po.ID, po.PONumber, po.BranchID, po.SupplierID, po.Status, po.TotalAmount, po.ExpectedDelivery, po.Notes, po.CreatedBy,
	).Scan(&po.CreatedAt, &po.UpdatedAt)
}

func (r *PurchaseRepository) UpdatePOStatus(ctx context.Context, id uuid.UUID, status domain.POStatus) error {
	query := `UPDATE purchase_orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}
