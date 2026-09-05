package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string
const (
	OrderPending           OrderStatus = "pending"
	OrderKitchenReceived   OrderStatus = "kitchen_received"
	OrderInPreparation     OrderStatus = "in_preparation"
	OrderReady             OrderStatus = "ready"
	OrderServed            OrderStatus = "served"
	OrderPayment           OrderStatus = "payment"
	OrderCompleted         OrderStatus = "completed"
	OrderCancelled         OrderStatus = "cancelled"
	OrderVoided            OrderStatus = "voided"
)

type OrderType string
const (
	OrderDineIn   OrderType = "dine_in"
	OrderTakeaway OrderType = "takeaway"
	OrderDelivery OrderType = "delivery"
)

type OrderItemStatus string
const (
	ItemPending   OrderItemStatus = "pending"
	ItemPreparing OrderItemStatus = "preparing"
	ItemReady     OrderItemStatus = "ready"
	ItemServed    OrderItemStatus = "served"
	ItemCancelled OrderItemStatus = "cancelled"
)

type Table struct {
	BaseEntity
	BranchID    uuid.UUID `json:"branch_id"`
	ZoneID      uuid.UUID `json:"zone_id"`
	TableNumber string    `json:"table_number"`
	Capacity    int       `json:"capacity"`
	Status      string    `json:"status"`
	XPos        int       `json:"x_pos"`
	YPos        int       `json:"y_pos"`
}

type TableZone struct {
	BaseEntity
	BranchID uuid.UUID `json:"branch_id"`
	Name     string    `json:"name"`
}

type Order struct {
	BaseEntity
	BranchID       uuid.UUID   `json:"branch_id"`
	OrderNumber    string      `json:"order_number"`
	TableID        *uuid.UUID  `json:"table_id"`
	CustomerID     *uuid.UUID  `json:"customer_id"`
	WaiterID       *uuid.UUID  `json:"waiter_id"`
	CashierID      *uuid.UUID  `json:"cashier_id"`
	Status         OrderStatus `json:"status"`
	OrderType      OrderType   `json:"order_type"`
	Subtotal       float64     `json:"subtotal"`
	TaxAmount      float64     `json:"tax_amount"`
	DiscountAmount float64     `json:"discount_amount"`
	ServiceCharge  float64     `json:"service_charge"`
	TotalAmount    float64     `json:"total_amount"`
	Notes          string      `json:"notes"`
}

type OrderItem struct {
	BaseEntity
	OrderID    uuid.UUID       `json:"order_id"`
	ProductID  uuid.UUID       `json:"product_id"`
	VariantID  *uuid.UUID      `json:"variant_id"`
	Quantity   int             `json:"quantity"`
	UnitPrice  float64         `json:"unit_price"`
	TotalPrice float64         `json:"total_price"`
	Status     OrderItemStatus `json:"status"`
	Notes      string          `json:"notes"`
}

type OrderItemModifier struct {
	BaseEntity
	OrderItemID     uuid.UUID `json:"order_item_id"`
	ModifierName    string    `json:"modifier_name"`
	AdditionalPrice float64   `json:"additional_price"`
}

type Payment struct {
	BaseEntity
	OrderID         uuid.UUID `json:"order_id"`
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	ReferenceNumber string    `json:"reference_number"`
	ProcessedBy     uuid.UUID `json:"processed_by"`
}

type PaymentMethod struct {
	BaseEntity
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	MDRPercentage float64 `json:"mdr_percentage"`
	IsActive      bool    `json:"is_active"`
}

type Promotion struct {
	BaseEntity
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Value       float64   `json:"value"`
	MinOrder    float64   `json:"min_order"`
	MaxDiscount float64   `json:"max_discount"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsActive    bool      `json:"is_active"`
}

type Voucher struct {
	BaseEntity
	Code         string    `json:"code"`
	PromotionID  uuid.UUID `json:"promotion_id"`
	MaxUsage     int       `json:"max_usage"`
	CurrentUsage int       `json:"current_usage"`
	IsActive     bool      `json:"is_active"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status OrderStatus) error
	List(ctx context.Context, branchID uuid.UUID) ([]Order, error)
}

type OrderUsecase interface {
	CreateOrder(ctx context.Context, req Order) (Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status OrderStatus) error
}
