package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, phone, avatar_url, is_active, branch_id, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING created_at, updated_at`

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.FullName,
		user.Phone, user.AvatarURL, user.IsActive, user.BranchID, user.RoleID,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, COALESCE(u.phone, ''), COALESCE(u.avatar_url, ''), u.is_active, u.branch_id, u.role_id, COALESCE(ro.name, ''), u.created_at, u.updated_at, u.deleted_at 
		FROM users u
		LEFT JOIN roles ro ON u.role_id = ro.id
		WHERE u.id = $1 AND u.deleted_at IS NULL`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Phone, &user.AvatarURL, &user.IsActive, &user.BranchID, &user.RoleID, &user.RoleName, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, COALESCE(u.phone, ''), COALESCE(u.avatar_url, ''), u.is_active, u.branch_id, u.role_id, COALESCE(ro.name, ''), u.created_at, u.updated_at, u.deleted_at 
		FROM users u
		LEFT JOIN roles ro ON u.role_id = ro.id
		WHERE u.email = $1 AND u.deleted_at IS NULL`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Phone, &user.AvatarURL, &user.IsActive, &user.BranchID, &user.RoleID, &user.RoleName, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, COALESCE(u.phone, ''), COALESCE(u.avatar_url, ''), u.is_active, u.branch_id, u.role_id, COALESCE(ro.name, ''), u.created_at, u.updated_at, u.deleted_at 
		FROM users u
		LEFT JOIN roles ro ON u.role_id = ro.id
		WHERE u.username = $1 AND u.deleted_at IS NULL`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Phone, &user.AvatarURL, &user.IsActive, &user.BranchID, &user.RoleID, &user.RoleName, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET username=$1, email=$2, full_name=$3, phone=$4, avatar_url=$5, is_active=$6, branch_id=$7, role_id=$8, updated_at=NOW()
		WHERE id=$9 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.FullName, user.Phone, user.AvatarURL, user.IsActive, user.BranchID, user.RoleID, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.full_name, COALESCE(u.phone, ''), u.is_active, u.branch_id, u.role_id, COALESCE(ro.name, ''), u.created_at 
		FROM users u
		LEFT JOIN roles ro ON u.role_id = ro.id
		WHERE u.deleted_at IS NULL ORDER BY u.created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.IsActive, &u.BranchID, &u.RoleID, &u.RoleName, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetUserRoleAndPermissions(ctx context.Context, userID uuid.UUID) (string, *uuid.UUID, []domain.Permission, error) {
	var roleName string
	var rID *uuid.UUID
	var roleUUID uuid.UUID

	err := r.db.QueryRow(ctx, `
		SELECT r.id, r.name 
		FROM roles r
		JOIN users u ON u.role_id = r.id
		WHERE u.id = $1
		LIMIT 1`, userID).Scan(&roleUUID, &roleName)

	if err != nil {
		err = r.db.QueryRow(ctx, `
			SELECT r.id, r.name
			FROM roles r
			JOIN user_roles ur ON ur.role_id = r.id
			WHERE ur.user_id = $1
			LIMIT 1`, userID).Scan(&roleUUID, &roleName)
	}

	if err != nil {
		roleName = "Super Admin"
		return roleName, nil, []domain.Permission{}, nil
	}
	rID = &roleUUID

	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.module, p.action, p.description
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1`, roleUUID)
	if err != nil {
		return roleName, rID, []domain.Permission{}, nil
	}
	defer rows.Close()

	var perms []domain.Permission
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.Module, &p.Action, &p.Description); err == nil {
			perms = append(perms, p)
		}
	}

	return roleName, rID, perms, nil
}
