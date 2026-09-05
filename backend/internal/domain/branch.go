package domain

import (
	"context"

	"github.com/google/uuid"
)

type Branch struct {
	BaseEntity
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Address   string     `json:"address"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	ManagerID *uuid.UUID `json:"manager_id"`
	IsActive  bool       `json:"is_active"`
}

type BranchSetting struct {
	BranchID     uuid.UUID `json:"branch_id"`
	SettingKey   string    `json:"setting_key"`
	SettingValue string    `json:"setting_value"` // JSON string
}

type BranchRepository interface {
	Create(ctx context.Context, branch *Branch) error
	GetByID(ctx context.Context, id uuid.UUID) (*Branch, error)
	Update(ctx context.Context, branch *Branch) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]Branch, error)
}

type BranchUsecase interface {
	Create(ctx context.Context, branch Branch) (Branch, error)
	GetByID(ctx context.Context, id uuid.UUID) (Branch, error)
	Update(ctx context.Context, id uuid.UUID, branch Branch) (Branch, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]Branch, error)
}
