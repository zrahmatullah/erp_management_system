package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) domain.RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `INSERT INTO roles (id, name, description, created_at) VALUES ($1, $2, $3, NOW()) RETURNING created_at`
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	return r.db.QueryRow(ctx, query, role.ID, role.Name, role.Description).Scan(&role.CreatedAt)
}

func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles WHERE id = $1`
	role := &domain.Role{}
	err := r.db.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `UPDATE roles SET name=$1, description=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, role.Name, role.Description, role.ID)
	return err
}

func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM roles WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *RoleRepository) List(ctx context.Context) ([]domain.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, roleID, permissionID)
	return err
}

func (r *RoleRepository) GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error) {
	query := `
		SELECT p.id, p.module, p.action, p.description 
		FROM permissions p 
		INNER JOIN role_permissions rp ON p.id = rp.permission_id 
		WHERE rp.role_id = $1`
	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []domain.Permission
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.Module, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}
