package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StockMovementType string
const (
	MovementIn         StockMovementType = "in"
	MovementOut        StockMovementType = "out"
	MovementTransfer   StockMovementType = "transfer"
	MovementAdjustment StockMovementType = "adjustment"
	MovementProduction StockMovementType = "production"
)

type POStatus string
const (
	PODraft             POStatus = "draft"
	POSubmitted         POStatus = "submitted"
	POManagerApproved   POStatus = "manager_approved"
	POOwnerApproved     POStatus = "owner_approved"
	PORejected          POStatus = "rejected"
	POSent              POStatus = "sent"
	POPartiallyReceived POStatus = "partially_received"
	POReceived          POStatus = "received"
	POInvoiced          POStatus = "invoiced"
	POPaid              POStatus = "paid"
	POCancelled         POStatus = "cancelled"
)

type InventoryItem struct {
	BaseEntity
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	UOM         string    `json:"uom"`
	MinStock    float64   `json:"min_stock"`
	MaxStock    float64   `json:"max_stock"`
	AverageCost float64   `json:"average_cost"`
	BranchID    uuid.UUID `json:"branch_id"`
}

type Warehouse struct {
	BaseEntity
	BranchID uuid.UUID `json:"branch_id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
}

type WarehouseLocation struct {
	BaseEntity
	WarehouseID uuid.UUID `json:"warehouse_id"`
	RackCode    string    `json:"rack_code"`
	Description string    `json:"description"`
}

type StockMovement struct {
	BaseEntity
	ItemID        uuid.UUID         `json:"item_id"`
	WarehouseID   uuid.UUID         `json:"warehouse_id"`
	Type          StockMovementType `json:"type"`
	Quantity      float64           `json:"quantity"`
	ReferenceType string            `json:"reference_type"`
	ReferenceID   *uuid.UUID        `json:"reference_id"`
	Remarks       string            `json:"remarks"`
	CreatedBy     uuid.UUID         `json:"created_by"`
}

type Supplier struct {
	BaseEntity
	Name         string `json:"name"`
	Address      string `json:"address"`
	TaxID        string `json:"tax_id"`
	PaymentTerms string `json:"payment_terms"`
	Rating       int    `json:"rating"`
	IsActive     bool   `json:"is_active"`
}

type SupplierContact struct {
	BaseEntity
	SupplierID uuid.UUID `json:"supplier_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	IsPrimary  bool      `json:"is_primary"`
}

type PurchaseOrder struct {
	BaseEntity
	PONumber         string     `json:"po_number"`
	BranchID         uuid.UUID  `json:"branch_id"`
	SupplierID       uuid.UUID  `json:"supplier_id"`
	Status           POStatus   `json:"status"`
	TotalAmount      float64    `json:"total_amount"`
	ExpectedDelivery time.Time  `json:"expected_delivery"`
	Notes            string     `json:"notes"`
	CreatedBy        uuid.UUID  `json:"created_by"`
	ApprovedBy       *uuid.UUID `json:"approved_by"`
}

type PurchaseOrderItem struct {
	BaseEntity
	POID        uuid.UUID `json:"po_id"`
	ItemID      uuid.UUID `json:"item_id"`
	Quantity    float64   `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	ReceivedQty float64   `json:"received_qty"`
}

type GoodsReceipt struct {
	BaseEntity
	GRNumber     string    `json:"gr_number"`
	POID         uuid.UUID `json:"po_id"`
	WarehouseID  uuid.UUID `json:"warehouse_id"`
	ReceivedDate time.Time `json:"received_date"`
	ReceivedBy   uuid.UUID `json:"received_by"`
	Status       string    `json:"status"`
	Notes        string    `json:"notes"`
}

type GoodsReceiptItem struct {
	BaseEntity
	GRID             uuid.UUID `json:"gr_id"`
	POItemID         uuid.UUID `json:"po_item_id"`
	ReceivedQuantity float64   `json:"received_quantity"`
	Condition        string    `json:"condition"`
	Notes            string    `json:"notes"`
}

type StockOpname struct {
	BaseEntity
	WarehouseID   uuid.UUID  `json:"warehouse_id"`
	OpnameNumber  string     `json:"opname_number"`
	Status        string     `json:"status"` // scheduled, in_progress, pending_review, approved, rejected
	ScheduledDate time.Time  `json:"scheduled_date"`
	CompletedDate *time.Time `json:"completed_date"`
	ConductedBy   uuid.UUID  `json:"conducted_by"`
	ApprovedBy    *uuid.UUID `json:"approved_by"`
}

type StockOpnameItem struct {
	BaseEntity
	OpnameID         uuid.UUID `json:"opname_id"`
	ItemID           uuid.UUID `json:"item_id"`
	SystemQuantity   float64   `json:"system_quantity"`
	PhysicalQuantity float64   `json:"physical_quantity"`
	Variance         float64   `json:"variance"`
	Reason           string    `json:"reason"`
}

type InventoryRepository interface {
	CreateItem(ctx context.Context, item *InventoryItem) error
	GetItem(ctx context.Context, id uuid.UUID) (*InventoryItem, error)
	ListItem(ctx context.Context, branchID uuid.UUID) ([]InventoryItem, error)
	CreateMovement(ctx context.Context, movement *StockMovement) error
}

type InventoryUsecase interface {
	CreateItem(ctx context.Context, item InventoryItem) (InventoryItem, error)
	GetItem(ctx context.Context, id uuid.UUID) (InventoryItem, error)
	RecordMovement(ctx context.Context, movement StockMovement) error
}
