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
		INSERT INTO users (id, username, email, password_hash, full_name, phone, avatar_url, is_active, branch_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at`

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.FullName,
		user.Phone, user.AvatarURL, user.IsActive, user.BranchID,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, COALESCE(phone, ''), COALESCE(avatar_url, ''), is_active, branch_id, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Phone, &user.AvatarURL, &user.IsActive, &user.BranchID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, COALESCE(phone, ''), COALESCE(avatar_url, ''), is_active, branch_id, created_at, updated_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Phone, &user.AvatarURL, &user.IsActive, &user.BranchID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, COALESCE(phone, ''), COALESCE(avatar_url, ''), is_active, branch_id, created_at, updated_at, deleted_at FROM users WHERE username = $1 AND deleted_at IS NULL`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Phone, &user.AvatarURL, &user.IsActive, &user.BranchID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET username=$1, email=$2, full_name=$3, phone=$4, avatar_url=$5, is_active=$6, branch_id=$7, updated_at=NOW()
		WHERE id=$8 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.FullName, user.Phone, user.AvatarURL, user.IsActive, user.BranchID, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	query := `SELECT id, username, email, full_name, COALESCE(phone, ''), is_active, branch_id, created_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.IsActive, &u.BranchID, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
