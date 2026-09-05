package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) domain.InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) CreateItem(ctx context.Context, item *domain.InventoryItem) error {
	query := `
		INSERT INTO inventory_items (id, sku, name, category, uom, min_stock, max_stock, average_cost, branch_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at`

	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		item.ID, item.SKU, item.Name, item.Category, item.UOM, item.MinStock, item.MaxStock, item.AverageCost, item.BranchID,
	).Scan(&item.CreatedAt, &item.UpdatedAt)
}

func (r *InventoryRepository) GetItem(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error) {
	query := `SELECT id, sku, name, category, uom, min_stock, max_stock, average_cost, branch_id, created_at, updated_at FROM inventory_items WHERE id = $1 AND deleted_at IS NULL`
	item := &domain.InventoryItem{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.SKU, &item.Name, &item.Category, &item.UOM, &item.MinStock, &item.MaxStock, &item.AverageCost, &item.BranchID, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *InventoryRepository) ListItem(ctx context.Context, branchID uuid.UUID) ([]domain.InventoryItem, error) {
	query := `SELECT id, sku, name, category, uom, min_stock, max_stock, average_cost, branch_id FROM inventory_items WHERE branch_id = $1 AND deleted_at IS NULL ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.InventoryItem
	for rows.Next() {
		var i domain.InventoryItem
		if err := rows.Scan(&i.ID, &i.SKU, &i.Name, &i.Category, &i.UOM, &i.MinStock, &i.MaxStock, &i.AverageCost, &i.BranchID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (r *InventoryRepository) CreateMovement(ctx context.Context, m *domain.StockMovement) error {
	query := `
		INSERT INTO stock_movements (id, item_id, warehouse_id, type, quantity, reference_type, reference_id, remarks, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`

	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query, m.ID, m.ItemID, m.WarehouseID, m.Type, m.Quantity, m.ReferenceType, m.ReferenceID, m.Remarks, m.CreatedBy)
	return err
}
