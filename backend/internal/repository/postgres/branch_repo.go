package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type BranchRepository struct {
	db *pgxpool.Pool
}

func NewBranchRepository(db *pgxpool.Pool) domain.BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	query := `
		INSERT INTO branches (id, name, code, address, phone, email, manager_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING created_at, updated_at`

	if branch.ID == uuid.Nil {
		branch.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		branch.ID, branch.Name, branch.Code, branch.Address,
		branch.Phone, branch.Email, branch.ManagerID, branch.IsActive,
	).Scan(&branch.CreatedAt, &branch.UpdatedAt)
}

func (r *BranchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	query := `SELECT id, name, code, address, phone, email, manager_id, is_active, created_at, updated_at FROM branches WHERE id = $1 AND deleted_at IS NULL`
	branch := &domain.Branch{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&branch.ID, &branch.Name, &branch.Code, &branch.Address,
		&branch.Phone, &branch.Email, &branch.ManagerID, &branch.IsActive,
		&branch.CreatedAt, &branch.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func (r *BranchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	query := `
		UPDATE branches 
		SET name=$1, code=$2, address=$3, phone=$4, email=$5, manager_id=$6, is_active=$7, updated_at=NOW()
		WHERE id=$8 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, branch.Name, branch.Code, branch.Address, branch.Phone, branch.Email, branch.ManagerID, branch.IsActive, branch.ID)
	return err
}

func (r *BranchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE branches SET deleted_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *BranchRepository) List(ctx context.Context) ([]domain.Branch, error) {
	query := `SELECT id, name, code, address, phone, email, is_active FROM branches WHERE deleted_at IS NULL ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []domain.Branch
	for rows.Next() {
		var b domain.Branch
		if err := rows.Scan(&b.ID, &b.Name, &b.Code, &b.Address, &b.Phone, &b.Email, &b.IsActive); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}
	return branches, nil
}
