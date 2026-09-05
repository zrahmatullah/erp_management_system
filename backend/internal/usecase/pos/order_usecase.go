package pos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
)

type OrderUsecaseImpl struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderUsecase(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) domain.OrderUsecase {
	return &OrderUsecaseImpl{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *OrderUsecaseImpl) CreateOrder(ctx context.Context, req domain.Order) (domain.Order, error) {
	req.OrderNumber = fmt.Sprintf("ORD-%s-%s", time.Now().Format("20060102"), uuid.New().String()[:6])
	
	if req.OrderType == domain.OrderDineIn {
		req.Status = domain.OrderKitchenReceived
	} else {
		req.Status = domain.OrderPending
	}
	
	req.TaxAmount = req.Subtotal * 0.11
	req.TotalAmount = req.Subtotal + req.TaxAmount - req.DiscountAmount + req.ServiceCharge
	
	err := u.orderRepo.Create(ctx, &req)
	return req, err
}

func (u *OrderUsecaseImpl) GetOrder(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Order{}, err
	}
	return *order, nil
}

func (u *OrderUsecaseImpl) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus domain.OrderStatus) error {
	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	validTransitions := map[domain.OrderStatus][]domain.OrderStatus{
		domain.OrderPending:         {domain.OrderKitchenReceived, domain.OrderCancelled},
		domain.OrderKitchenReceived: {domain.OrderInPreparation, domain.OrderCancelled},
		domain.OrderInPreparation:   {domain.OrderReady, domain.OrderCancelled},
		domain.OrderReady:           {domain.OrderServed, domain.OrderCancelled},
		domain.OrderServed:          {domain.OrderPayment, domain.OrderCompleted},
		domain.OrderPayment:         {domain.OrderCompleted},
		domain.OrderCompleted:       {domain.OrderVoided},
	}

	isValid := false
	if transitions, ok := validTransitions[order.Status]; ok {
		for _, s := range transitions {
			if s == newStatus {
				isValid = true
				break
			}
		}
	}

	if !isValid {
		return errors.New("invalid status transition")
	}

	return u.orderRepo.UpdateStatus(ctx, id, newStatus)
}
