package pos

import (
	"context"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
)

type ProductUsecaseImpl struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &ProductUsecaseImpl{productRepo: productRepo}
}

func (u *ProductUsecaseImpl) CreateProduct(ctx context.Context, req domain.Product) (domain.Product, error) {
	err := u.productRepo.Create(ctx, &req)
	return req, err
}

func (u *ProductUsecaseImpl) GetProduct(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	p, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}
	return *p, nil
}

func (u *ProductUsecaseImpl) ListProducts(ctx context.Context, branchID uuid.UUID) ([]domain.Product, error) {
	return u.productRepo.List(ctx, branchID)
}

func (u *ProductUsecaseImpl) UpdateProduct(ctx context.Context, id uuid.UUID, req domain.Product) (domain.Product, error) {
	req.ID = id
	err := u.productRepo.Update(ctx, &req)
	return req, err
}

func (u *ProductUsecaseImpl) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return u.productRepo.Delete(ctx, id)
}
