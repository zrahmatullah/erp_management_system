package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type SupplierRepository struct {
	db *pgxpool.Pool
}

func NewSupplierRepository(db *pgxpool.Pool) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) Create(ctx context.Context, s *domain.Supplier) error {
	query := `
		INSERT INTO suppliers (id, name, address, tax_id, payment_terms, rating, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING created_at, updated_at`

	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		s.ID, s.Name, s.Address, s.TaxID, s.PaymentTerms, s.Rating, s.IsActive,
	).Scan(&s.CreatedAt, &s.UpdatedAt)
}

func (r *SupplierRepository) List(ctx context.Context) ([]domain.Supplier, error) {
	query := `SELECT id, name, address, tax_id, payment_terms, rating, is_active FROM suppliers WHERE deleted_at IS NULL ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []domain.Supplier
	for rows.Next() {
		var s domain.Supplier
		if err := rows.Scan(&s.ID, &s.Name, &s.Address, &s.TaxID, &s.PaymentTerms, &s.Rating, &s.IsActive); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}
	return suppliers, nil
}
