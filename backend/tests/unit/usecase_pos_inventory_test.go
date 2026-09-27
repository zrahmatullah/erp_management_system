package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"cafe-erp-system/backend/internal/domain"
	"cafe-erp-system/backend/internal/usecase/inventory"
	"cafe-erp-system/backend/internal/usecase/pos"
)

// MockOrderRepository implements domain.OrderRepository
type MockOrderRepository struct {
	orders map[uuid.UUID]*domain.Order
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{orders: make(map[uuid.UUID]*domain.Order)}
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	m.orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return o, nil
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	o, ok := m.orders[id]
	if !ok {
		return errors.New("order not found")
	}
	o.Status = status
	return nil
}

func (m *MockOrderRepository) List(ctx context.Context, branchID uuid.UUID) ([]domain.Order, error) {
	var list []domain.Order
	for _, o := range m.orders {
		if o.BranchID == branchID {
			list = append(list, *o)
		}
	}
	return list, nil
}

// MockProductRepository implements domain.ProductRepository
type MockProductRepository struct {
	products map[uuid.UUID]*domain.Product
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{products: make(map[uuid.UUID]*domain.Product)}
}

func (m *MockProductRepository) Create(ctx context.Context, p *domain.Product) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	m.products[p.ID] = p
	return nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return p, nil
}

func (m *MockProductRepository) Update(ctx context.Context, p *domain.Product) error {
	if _, ok := m.products[p.ID]; !ok {
		return errors.New("product not found")
	}
	m.products[p.ID] = p
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.products[id]; !ok {
		return errors.New("product not found")
	}
	delete(m.products, id)
	return nil
}

func (m *MockProductRepository) List(ctx context.Context, branchID uuid.UUID) ([]domain.Product, error) {
	var list []domain.Product
	for _, p := range m.products {
		if p.BranchID == branchID {
			list = append(list, *p)
		}
	}
	return list, nil
}

// MockInventoryRepository implements domain.InventoryRepository
type MockInventoryRepository struct {
	items     map[uuid.UUID]*domain.InventoryItem
	movements []*domain.StockMovement
}

func NewMockInventoryRepository() *MockInventoryRepository {
	return &MockInventoryRepository{
		items:     make(map[uuid.UUID]*domain.InventoryItem),
		movements: make([]*domain.StockMovement, 0),
	}
}

func (m *MockInventoryRepository) CreateItem(ctx context.Context, item *domain.InventoryItem) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	m.items[item.ID] = item
	return nil
}

func (m *MockInventoryRepository) GetItem(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, errors.New("inventory item not found")
	}
	return item, nil
}

func (m *MockInventoryRepository) ListItem(ctx context.Context, branchID uuid.UUID) ([]domain.InventoryItem, error) {
	var list []domain.InventoryItem
	for _, it := range m.items {
		if it.BranchID == branchID {
			list = append(list, *it)
		}
	}
	return list, nil
}

func (m *MockInventoryRepository) CreateMovement(ctx context.Context, movement *domain.StockMovement) error {
	m.movements = append(m.movements, movement)
	return nil
}

func TestPOS_OrderUsecase(t *testing.T) {
	orderRepo := NewMockOrderRepository()
	prodRepo := NewMockProductRepository()
	usecase := pos.NewOrderUsecase(orderRepo, prodRepo)

	t.Run("Create Dine-In Order sets KitchenReceived and Calculates 11% Tax", func(t *testing.T) {
		req := domain.Order{
			BranchID:       uuid.New(),
			OrderType:      domain.OrderDineIn,
			Subtotal:       100000,
			DiscountAmount: 10000,
			ServiceCharge:  5000,
		}

		created, err := usecase.CreateOrder(context.Background(), req)
		assert.NoError(t, err)
		assert.NotEmpty(t, created.OrderNumber)
		assert.Equal(t, domain.OrderKitchenReceived, created.Status)
		assert.Equal(t, 11000.0, created.TaxAmount)   // 11% of 100,000
		assert.Equal(t, 106000.0, created.TotalAmount) // 100,000 + 11,000 - 10,000 + 5,000
	})

	t.Run("Create Takeaway Order sets Pending Status", func(t *testing.T) {
		req := domain.Order{
			BranchID:       uuid.New(),
			OrderType:      domain.OrderTakeaway,
			Subtotal:       50000,
			DiscountAmount: 0,
		}

		created, err := usecase.CreateOrder(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, domain.OrderPending, created.Status)
		assert.Equal(t, 5500.0, created.TaxAmount)
		assert.Equal(t, 55500.0, created.TotalAmount)
	})

	t.Run("Get Existing Order", func(t *testing.T) {
		orderID := uuid.New()
		order := domain.Order{
			BaseEntity: domain.BaseEntity{ID: orderID},
			Status:     domain.OrderPending,
			Subtotal:   25000,
		}
		_ = orderRepo.Create(context.Background(), &order)

		fetched, err := usecase.GetOrder(context.Background(), orderID)
		assert.NoError(t, err)
		assert.Equal(t, orderID, fetched.ID)
	})

	t.Run("Get Nonexistent Order Returns Error", func(t *testing.T) {
		_, err := usecase.GetOrder(context.Background(), uuid.New())
		assert.Error(t, err)
	})

	t.Run("Valid Order Status Transitions", func(t *testing.T) {
		orderID := uuid.New()
		order := domain.Order{
			BaseEntity: domain.BaseEntity{ID: orderID},
			Status:     domain.OrderPending,
		}
		_ = orderRepo.Create(context.Background(), &order)

		// Pending -> KitchenReceived (valid)
		err := usecase.UpdateStatus(context.Background(), orderID, domain.OrderKitchenReceived)
		assert.NoError(t, err)

		// KitchenReceived -> InPreparation (valid)
		err = usecase.UpdateStatus(context.Background(), orderID, domain.OrderInPreparation)
		assert.NoError(t, err)

		// InPreparation -> Ready (valid)
		err = usecase.UpdateStatus(context.Background(), orderID, domain.OrderReady)
		assert.NoError(t, err)

		// Ready -> Served (valid)
		err = usecase.UpdateStatus(context.Background(), orderID, domain.OrderServed)
		assert.NoError(t, err)

		// Served -> Completed (valid)
		err = usecase.UpdateStatus(context.Background(), orderID, domain.OrderCompleted)
		assert.NoError(t, err)
	})

	t.Run("Invalid Order Status Transition Rejected", func(t *testing.T) {
		orderID := uuid.New()
		order := domain.Order{
			BaseEntity: domain.BaseEntity{ID: orderID},
			Status:     domain.OrderPending,
		}
		_ = orderRepo.Create(context.Background(), &order)

		// Pending directly to Completed (skips kitchen, prep, ready, serve) -> INVALID!
		err := usecase.UpdateStatus(context.Background(), orderID, domain.OrderCompleted)
		assert.Error(t, err)
		assert.Equal(t, "invalid status transition", err.Error())
	})
}

func TestPOS_ProductUsecase(t *testing.T) {
	prodRepo := NewMockProductRepository()
	usecase := pos.NewProductUsecase(prodRepo)
	branchID := uuid.New()

	t.Run("Create, Get, List, Update, and Delete Product", func(t *testing.T) {
		p := domain.Product{
			Name:      "Caramel Macchiato",
			SKU:       "COF-001",
			BasePrice: 38000,
			BranchID:  branchID,
			IsActive:  true,
		}

		created, err := usecase.CreateProduct(context.Background(), p)
		assert.NoError(t, err)
		assert.Equal(t, "Caramel Macchiato", created.Name)

		// GetProduct
		fetched, err := usecase.GetProduct(context.Background(), created.ID)
		assert.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)

		// ListProducts
		list, err := usecase.ListProducts(context.Background(), branchID)
		assert.NoError(t, err)
		assert.Len(t, list, 1)

		// UpdateProduct
		created.BasePrice = 42000
		updated, err := usecase.UpdateProduct(context.Background(), created.ID, created)
		assert.NoError(t, err)
		assert.Equal(t, 42000.0, updated.BasePrice)

		// DeleteProduct
		err = usecase.DeleteProduct(context.Background(), created.ID)
		assert.NoError(t, err)

		_, err = usecase.GetProduct(context.Background(), created.ID)
		assert.Error(t, err)
	})
}

func TestInventory_StockUsecase(t *testing.T) {
	invRepo := NewMockInventoryRepository()
	usecase := inventory.NewStockUsecase(invRepo)

	t.Run("CreateItem, GetItem, and RecordMovement", func(t *testing.T) {
		item := domain.InventoryItem{
			Name:        "Arabica Roasted Beans",
			SKU:         "RAW-001",
			Category:    "Coffee Beans",
			UOM:         "KG",
			MinStock:    10,
			MaxStock:    100,
			AverageCost: 120000,
			BranchID:    uuid.New(),
		}

		created, err := usecase.CreateItem(context.Background(), item)
		assert.NoError(t, err)
		assert.Equal(t, "Arabica Roasted Beans", created.Name)

		fetched, err := usecase.GetItem(context.Background(), created.ID)
		assert.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)

		movement := domain.StockMovement{
			ItemID:        created.ID,
			WarehouseID:   uuid.New(),
			Type:          domain.MovementOut,
			Quantity:      -5,
			ReferenceType: "pos_order",
			Remarks:       "POS Coffee sale consumption",
		}
		err = usecase.RecordMovement(context.Background(), movement)
		assert.NoError(t, err)
		assert.Len(t, invRepo.movements, 1)
	})

	t.Run("Get Nonexistent Item Returns Error", func(t *testing.T) {
		_, err := usecase.GetItem(context.Background(), uuid.New())
		assert.Error(t, err)
	})
}
