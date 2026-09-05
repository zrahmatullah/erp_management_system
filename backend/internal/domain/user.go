package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	BaseEntity
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Phone        string     `json:"phone"`
	AvatarURL    string     `json:"avatar_url"`
	IsActive     bool       `json:"is_active"`
	BranchID     *uuid.UUID `json:"branch_id"`
}

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Module      string    `json:"module"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
}

type UserRole struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
}

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]User, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]Role, error)
	AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]Permission, error)
}

type AuthUsecase interface {
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)
	Register(ctx context.Context, req RegisterRequest) (User, error)
	RefreshToken(ctx context.Context, token string) (LoginResponse, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type UserUsecase interface {
	Create(ctx context.Context, req RegisterRequest) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	Update(ctx context.Context, id uuid.UUID, req RegisterRequest) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]User, int, error)
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
	GetUsersByBranch(ctx context.Context, branchID uuid.UUID) ([]User, error)
}
