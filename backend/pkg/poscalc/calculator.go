package poscalc

import (
	"errors"
	"math"
)

var (
	ErrEmptyOrder         = errors.New("order must contain at least one item")
	ErrInvalidQuantity    = errors.New("item quantity must be greater than zero")
	ErrInvalidUnitPrice   = errors.New("item unit price cannot be negative")
	ErrInvalidDiscount    = errors.New("discount cannot be negative or exceed 100%")
	ErrInsufficientAmount = errors.New("amount paid is less than the total amount due")
	ErrNegativePayment    = errors.New("amount paid cannot be negative")
)

// DiscountType represents percentage or fixed amount discounts.
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
)

// OrderItemInput defines the line item payload for POS calculations.
type OrderItemInput struct {
	ProductID        string
	Name             string
	Quantity         int
	UnitPrice        float64
	DiscountType     DiscountType
	DiscountValue    float64
}

// PromotionInput defines store-wide or voucher-based promotions.
type PromotionInput struct {
	Type        DiscountType // percentage or fixed
	Value       float64
	MinOrder    float64
	MaxDiscount float64
}

// TaxConfig defines tax and service charges (PB1 / PPN / Service).
type TaxConfig struct {
	TaxRate            float64 // e.g. 0.10 for 10% PB1
	ServiceChargeRate  float64 // e.g. 0.05 for 5% Service Charge
	TaxAppliedAfterDisc bool    // true if tax applies to net amount after discounts
}

// DefaultTaxConfig provides standard Indonesian Restaurant Tax (PB1 10%).
func DefaultTaxConfig() TaxConfig {
	return TaxConfig{
		TaxRate:            0.10,
		ServiceChargeRate:  0.0,
		TaxAppliedAfterDisc: true,
	}
}

// OrderItemResult holds computed financial figures for a single line item.
type OrderItemResult struct {
	ProductID      string  `json:"product_id"`
	Name           string  `json:"name"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	GrossAmount    float64 `json:"gross_amount"`
	DiscountAmount float64 `json:"discount_amount"`
	NetAmount      float64 `json:"net_amount"`
}

// CalculationResult provides full financial transparency of the order.
type CalculationResult struct {
	Items              []OrderItemResult `json:"items"`
	GrossSubtotal      float64           `json:"gross_subtotal"`
	TotalItemDiscounts float64           `json:"total_item_discounts"`
	NetSubtotal        float64           `json:"net_subtotal"`
	OrderDiscount      float64           `json:"order_discount"`
	TotalDiscounts     float64           `json:"total_discounts"`
	TaxableBase        float64           `json:"taxable_base"`
	ServiceCharge      float64           `json:"service_charge"`
	TaxAmount          float64           `json:"tax_amount"`
	TotalAmount        float64           `json:"total_amount"`
}

// RoundIDR rounds a monetary value to the nearest integer Rupiah.
func RoundIDR(val float64) float64 {
	return math.Round(val)
}

// CalculateOrderItem computes gross, discount, and net for a single line item.
func CalculateOrderItem(item OrderItemInput) (OrderItemResult, error) {
	if item.Quantity <= 0 {
		return OrderItemResult{}, ErrInvalidQuantity
	}
	if item.UnitPrice < 0 {
		return OrderItemResult{}, ErrInvalidUnitPrice
	}

	gross := float64(item.Quantity) * item.UnitPrice
	var discount float64

	switch item.DiscountType {
	case DiscountTypePercentage:
		if item.DiscountValue < 0 || item.DiscountValue > 100 {
			return OrderItemResult{}, ErrInvalidDiscount
		}
		discount = gross * (item.DiscountValue / 100.0)
	case DiscountTypeFixed:
		if item.DiscountValue < 0 {
			return OrderItemResult{}, ErrInvalidDiscount
		}
		discount = item.DiscountValue
		if discount > gross {
			discount = gross // Cap discount to line item total
		}
	}

	discount = RoundIDR(discount)
	net := gross - discount
	if net < 0 {
		net = 0
	}

	return OrderItemResult{
		ProductID:      item.ProductID,
		Name:           item.Name,
		Quantity:       item.Quantity,
		UnitPrice:      item.UnitPrice,
		GrossAmount:    gross,
		DiscountAmount: discount,
		NetAmount:      net,
	}, nil
}

// CalculateOrder calculates complete order totals including items, vouchers, taxes, and service charges.
func CalculateOrder(items []OrderItemInput, promo *PromotionInput, cfg TaxConfig) (CalculationResult, error) {
	if len(items) == 0 {
		return CalculationResult{}, ErrEmptyOrder
	}

	var computedItems []OrderItemResult
	var grossSubtotal, totalItemDiscounts float64

	for _, it := range items {
		res, err := CalculateOrderItem(it)
		if err != nil {
			return CalculationResult{}, err
		}
		computedItems = append(computedItems, res)
		grossSubtotal += res.GrossAmount
		totalItemDiscounts += res.DiscountAmount
	}

	netSubtotal := grossSubtotal - totalItemDiscounts
	if netSubtotal < 0 {
		netSubtotal = 0
	}

	// Order-level promotion / voucher
	var orderDiscount float64
	if promo != nil && promo.Value > 0 {
		if netSubtotal >= promo.MinOrder {
			switch promo.Type {
			case DiscountTypePercentage:
				if promo.Value > 0 && promo.Value <= 100 {
					orderDiscount = netSubtotal * (promo.Value / 100.0)
					if promo.MaxDiscount > 0 && orderDiscount > promo.MaxDiscount {
						orderDiscount = promo.MaxDiscount
					}
				}
			case DiscountTypeFixed:
				orderDiscount = promo.Value
				if orderDiscount > netSubtotal {
					orderDiscount = netSubtotal
				}
			}
			orderDiscount = RoundIDR(orderDiscount)
		}
	}

	totalDiscounts := totalItemDiscounts + orderDiscount
	if totalDiscounts > grossSubtotal {
		totalDiscounts = grossSubtotal
	}

	taxableBase := netSubtotal - orderDiscount
	if taxableBase < 0 {
		taxableBase = 0
	}

	// Service charge (applied to taxable base)
	var serviceCharge float64
	if cfg.ServiceChargeRate > 0 {
		serviceCharge = RoundIDR(taxableBase * cfg.ServiceChargeRate)
	}

	// Tax calculation (PB1 / PPN)
	var taxBase float64
	if cfg.TaxAppliedAfterDisc {
		taxBase = taxableBase + serviceCharge
	} else {
		taxBase = grossSubtotal + serviceCharge
	}

	var taxAmount float64
	if cfg.TaxRate > 0 {
		taxAmount = RoundIDR(taxBase * cfg.TaxRate)
	}

	totalAmount := RoundIDR(taxableBase + serviceCharge + taxAmount)

	return CalculationResult{
		Items:              computedItems,
		GrossSubtotal:      grossSubtotal,
		TotalItemDiscounts: totalItemDiscounts,
		NetSubtotal:        netSubtotal,
		OrderDiscount:      orderDiscount,
		TotalDiscounts:     totalDiscounts,
		TaxableBase:        taxableBase,
		ServiceCharge:      serviceCharge,
		TaxAmount:          taxAmount,
		TotalAmount:        totalAmount,
	}, nil
}

// CalculateChange determines the change due to the customer or flags underpayment.
func CalculateChange(totalAmount, amountPaid float64) (changeDue float64, err error) {
	if amountPaid < 0 {
		return 0, ErrNegativePayment
	}
	if amountPaid < totalAmount {
		return 0, ErrInsufficientAmount
	}
	return RoundIDR(amountPaid - totalAmount), nil
}
