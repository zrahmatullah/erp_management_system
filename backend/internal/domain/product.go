package domain

import (
	"context"

	"github.com/google/uuid"
)

type Category struct {
	BaseEntity
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id"`
	ImageURL  string     `json:"image_url"`
	SortOrder int        `json:"sort_order"`
	IsActive  bool       `json:"is_active"`
}

type Product struct {
	BaseEntity
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `json:"is_active"`
	BranchID    uuid.UUID `json:"branch_id"`
}

type ProductVariant struct {
	BaseEntity
	ProductID       uuid.UUID `json:"product_id"`
	Name            string    `json:"name"`
	AdditionalPrice float64   `json:"additional_price"`
	SKU             string    `json:"sku"`
}

type ProductImage struct {
	BaseEntity
	ProductID uuid.UUID `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary"`
	SortOrder int       `json:"sort_order"`
}

type Recipe struct {
	BaseEntity
	ProductID    uuid.UUID  `json:"product_id"`
	VariantID    *uuid.UUID `json:"variant_id"`
	YieldAmount  float64    `json:"yield_amount"`
	Instructions string     `json:"instructions"`
}

type RecipeItem struct {
	BaseEntity
	RecipeID        uuid.UUID `json:"recipe_id"`
	InventoryItemID uuid.UUID `json:"inventory_item_id"`
	Quantity        float64   `json:"quantity"`
	UOM             string    `json:"uom"`
}

type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	List(ctx context.Context, branchID uuid.UUID) ([]Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req Product) (Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (Product, error)
	ListProducts(ctx context.Context, branchID uuid.UUID) ([]Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req Product) (Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}
