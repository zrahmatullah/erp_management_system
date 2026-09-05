package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) domain.OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, branch_id, order_number, table_id, customer_id, waiter_id, cashier_id, status, order_type, subtotal, tax_amount, discount_amount, service_charge, total_amount, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW())
		RETURNING created_at, updated_at`

	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		order.ID, order.BranchID, order.OrderNumber, order.TableID, order.CustomerID, order.WaiterID,
		order.CashierID, order.Status, order.OrderType, order.Subtotal, order.TaxAmount,
		order.DiscountAmount, order.ServiceCharge, order.TotalAmount, order.Notes,
	).Scan(&order.CreatedAt, &order.UpdatedAt)
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `SELECT id, branch_id, order_number, status, order_type, total_amount, created_at FROM orders WHERE id = $1 AND deleted_at IS NULL`
	order := &domain.Order{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID, &order.BranchID, &order.OrderNumber, &order.Status, &order.OrderType, &order.TotalAmount, &order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *OrderRepository) List(ctx context.Context, branchID uuid.UUID) ([]domain.Order, error) {
	query := `SELECT id, branch_id, order_number, status, order_type, total_amount, created_at FROM orders WHERE branch_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 100`
	rows, err := r.db.Query(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.BranchID, &o.OrderNumber, &o.Status, &o.OrderType, &o.TotalAmount, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
