package inventory

import (
	"context"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
)

type StockUsecaseImpl struct {
	repo domain.InventoryRepository
}

func NewStockUsecase(repo domain.InventoryRepository) domain.InventoryUsecase {
	return &StockUsecaseImpl{repo: repo}
}

func (u *StockUsecaseImpl) CreateItem(ctx context.Context, item domain.InventoryItem) (domain.InventoryItem, error) {
	err := u.repo.CreateItem(ctx, &item)
	return item, err
}

func (u *StockUsecaseImpl) GetItem(ctx context.Context, id uuid.UUID) (domain.InventoryItem, error) {
	item, err := u.repo.GetItem(ctx, id)
	if err != nil {
		return domain.InventoryItem{}, err
	}
	return *item, nil
}

func (u *StockUsecaseImpl) RecordMovement(ctx context.Context, movement domain.StockMovement) error {
	return u.repo.CreateMovement(ctx, &movement)
}
