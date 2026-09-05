package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) domain.ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (id, category_id, name, sku, description, base_price, image_url, is_active, branch_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at`

	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		p.ID, p.CategoryID, p.Name, p.SKU, p.Description, p.BasePrice, p.ImageURL, p.IsActive, p.BranchID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `SELECT id, category_id, name, sku, description, base_price, image_url, is_active, branch_id, created_at, updated_at FROM products WHERE id = $1 AND deleted_at IS NULL`
	p := &domain.Product{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CategoryID, &p.Name, &p.SKU, &p.Description, &p.BasePrice, &p.ImageURL, &p.IsActive, &p.BranchID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepository) List(ctx context.Context, branchID uuid.UUID) ([]domain.Product, error) {
	query := `SELECT id, category_id, name, sku, description, base_price, image_url, is_active, branch_id FROM products WHERE branch_id = $1 AND deleted_at IS NULL ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.SKU, &p.Description, &p.BasePrice, &p.ImageURL, &p.IsActive, &p.BranchID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	query := `
		UPDATE products 
		SET category_id=$1, name=$2, sku=$3, description=$4, base_price=$5, image_url=$6, is_active=$7, updated_at=NOW()
		WHERE id=$8 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, p.CategoryID, p.Name, p.SKU, p.Description, p.BasePrice, p.ImageURL, p.IsActive, p.ID)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET deleted_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
