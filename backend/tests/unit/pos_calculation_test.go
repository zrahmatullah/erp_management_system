package unit_test

import (
	"testing"

	"cafe-erp-system/backend/pkg/poscalc"
)

// TestPOS_LineItemCalculation tests item-level computations, discounts, and validations.
func TestPOS_LineItemCalculation(t *testing.T) {
	tests := []struct {
		name        string
		input       poscalc.OrderItemInput
		wantGross   float64
		wantDisc    float64
		wantNet     float64
		wantErr     error
	}{
		{
			name: "Standard Item Without Discount",
			input: poscalc.OrderItemInput{
				ProductID: "prod-1",
				Name:      "Espresso",
				Quantity:  2,
				UnitPrice: 25000,
			},
			wantGross: 50000,
			wantDisc:  0,
			wantNet:   50000,
			wantErr:   nil,
		},
		{
			name: "Item with 20% Percentage Discount",
			input: poscalc.OrderItemInput{
				ProductID:     "prod-2",
				Name:          "Cafe Latte",
				Quantity:      2,
				UnitPrice:     35000,
				DiscountType:  poscalc.DiscountTypePercentage,
				DiscountValue: 20,
			},
			wantGross: 70000,
			wantDisc:  14000,
			wantNet:   56000,
			wantErr:   nil,
		},
		{
			name: "Item with Fixed Discount",
			input: poscalc.OrderItemInput{
				ProductID:     "prod-3",
				Name:          "Croissant",
				Quantity:      3,
				UnitPrice:     20000,
				DiscountType:  poscalc.DiscountTypeFixed,
				DiscountValue: 15000,
			},
			wantGross: 60000,
			wantDisc:  15000,
			wantNet:   45000,
			wantErr:   nil,
		},
		{
			name: "Fixed Discount Capped at Item Gross Total",
			input: poscalc.OrderItemInput{
				ProductID:     "prod-4",
				Name:          "Mineral Water",
				Quantity:      1,
				UnitPrice:     10000,
				DiscountType:  poscalc.DiscountTypeFixed,
				DiscountValue: 15000, // Exceeds gross
			},
			wantGross: 10000,
			wantDisc:  10000, // Capped
			wantNet:   0,
			wantErr:   nil,
		},
		{
			name: "Invalid Zero Quantity",
			input: poscalc.OrderItemInput{
				ProductID: "prod-5",
				Quantity:  0,
				UnitPrice: 10000,
			},
			wantErr: poscalc.ErrInvalidQuantity,
		},
		{
			name: "Invalid Negative Quantity",
			input: poscalc.OrderItemInput{
				ProductID: "prod-6",
				Quantity:  -2,
				UnitPrice: 10000,
			},
			wantErr: poscalc.ErrInvalidQuantity,
		},
		{
			name: "Invalid Negative Unit Price",
			input: poscalc.OrderItemInput{
				ProductID: "prod-7",
				Quantity:  1,
				UnitPrice: -5000,
			},
			wantErr: poscalc.ErrInvalidUnitPrice,
		},
		{
			name: "Invalid Discount Percentage > 100",
			input: poscalc.OrderItemInput{
				ProductID:     "prod-8",
				Quantity:      1,
				UnitPrice:     50000,
				DiscountType:  poscalc.DiscountTypePercentage,
				DiscountValue: 120,
			},
			wantErr: poscalc.ErrInvalidDiscount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := poscalc.CalculateOrderItem(tt.input)
			if tt.wantErr != nil {
				if err == nil || err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.GrossAmount != tt.wantGross {
				t.Errorf("GrossAmount = %v; want %v", res.GrossAmount, tt.wantGross)
			}
			if res.DiscountAmount != tt.wantDisc {
				t.Errorf("DiscountAmount = %v; want %v", res.DiscountAmount, tt.wantDisc)
			}
			if res.NetAmount != tt.wantNet {
				t.Errorf("NetAmount = %v; want %v", res.NetAmount, tt.wantNet)
			}
		})
	}
}

// TestPOS_OrderCalculation_Scenarios tests complete order totals including vouchers, tax (PB1), and service charge.
func TestPOS_OrderCalculation_Scenarios(t *testing.T) {
	t.Run("Standard Dine-In Order with 10% PB1 Tax", func(t *testing.T) {
		items := []poscalc.OrderItemInput{
			{ProductID: "1", Name: "Nasi Goreng", Quantity: 2, UnitPrice: 35000},
			{ProductID: "2", Name: "Es Teh Manis", Quantity: 2, UnitPrice: 8000},
		}
		// Gross: (2*35000) + (2*8000) = 70000 + 16000 = 86000
		// Tax (10%): 8600
		// Total: 94600
		cfg := poscalc.DefaultTaxConfig()

		res, err := poscalc.CalculateOrder(items, nil, cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.GrossSubtotal != 86000 {
			t.Errorf("GrossSubtotal = %v; want 86000", res.GrossSubtotal)
		}
		if res.TaxAmount != 8600 {
			t.Errorf("TaxAmount = %v; want 8600", res.TaxAmount)
		}
		if res.TotalAmount != 94600 {
			t.Errorf("TotalAmount = %v; want 94600", res.TotalAmount)
		}
	})

	t.Run("Order with Voucher Promotion (Percentage with Cap)", func(t *testing.T) {
		items := []poscalc.OrderItemInput{
			{ProductID: "1", Name: "Sirloin Steak", Quantity: 2, UnitPrice: 120000}, // 240,000
		}
		// Promo: 50% discount, MinOrder: 100,000, MaxDiscount: 50,000
		// Calculated discount would be 120,000, but capped at 50,000
		// Taxable Base: 240,000 - 50,000 = 190,000
		// Tax (10%): 19,000
		// Total: 209,000
		promo := &poscalc.PromotionInput{
			Type:        poscalc.DiscountTypePercentage,
			Value:       50,
			MinOrder:    100000,
			MaxDiscount: 50000,
		}
		cfg := poscalc.DefaultTaxConfig()

		res, err := poscalc.CalculateOrder(items, promo, cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.OrderDiscount != 50000 {
			t.Errorf("OrderDiscount = %v; want 50000 (capped)", res.OrderDiscount)
		}
		if res.TaxableBase != 190000 {
			t.Errorf("TaxableBase = %v; want 190000", res.TaxableBase)
		}
		if res.TaxAmount != 19000 {
			t.Errorf("TaxAmount = %v; want 19000", res.TaxAmount)
		}
		if res.TotalAmount != 209000 {
			t.Errorf("TotalAmount = %v; want 209000", res.TotalAmount)
		}
	})

	t.Run("Order with Voucher Below Minimum Order Threshold", func(t *testing.T) {
		items := []poscalc.OrderItemInput{
			{ProductID: "1", Name: "Americano", Quantity: 1, UnitPrice: 28000},
		}
		// Promo requires min 50,000 order; current subtotal is 28,000 -> Promo NOT applied
		promo := &poscalc.PromotionInput{
			Type:        poscalc.DiscountTypeFixed,
			Value:       10000,
			MinOrder:    50000,
			MaxDiscount: 10000,
		}
		cfg := poscalc.DefaultTaxConfig()

		res, err := poscalc.CalculateOrder(items, promo, cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.OrderDiscount != 0 {
			t.Errorf("OrderDiscount = %v; want 0 (min order not met)", res.OrderDiscount)
		}
		if res.TotalAmount != 30800 { // 28000 + 2800 tax
			t.Errorf("TotalAmount = %v; want 30800", res.TotalAmount)
		}
	})

	t.Run("Order with Service Charge and 11% PPN", func(t *testing.T) {
		items := []poscalc.OrderItemInput{
			{ProductID: "1", Name: "Premium Buffet", Quantity: 1, UnitPrice: 100000},
		}
		// Subtotal: 100,000
		// Service Charge (5%): 5,000
		// Tax Base (Subtotal + Service): 105,000
		// Tax (11%): 11,550
		// Total: 100,000 + 5,000 + 11,550 = 116,550
		cfg := poscalc.TaxConfig{
			TaxRate:             0.11,
			ServiceChargeRate:   0.05,
			TaxAppliedAfterDisc: true,
		}

		res, err := poscalc.CalculateOrder(items, nil, cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.ServiceCharge != 5000 {
			t.Errorf("ServiceCharge = %v; want 5000", res.ServiceCharge)
		}
		if res.TaxAmount != 11550 {
			t.Errorf("TaxAmount = %v; want 11550", res.TaxAmount)
		}
		if res.TotalAmount != 116550 {
			t.Errorf("TotalAmount = %v; want 116550", res.TotalAmount)
		}
	})

	t.Run("Edge Case: 100% Free Promotion", func(t *testing.T) {
		items := []poscalc.OrderItemInput{
			{ProductID: "1", Name: "Birthday Coffee", Quantity: 1, UnitPrice: 35000},
		}
		promo := &poscalc.PromotionInput{
			Type:  poscalc.DiscountTypePercentage,
			Value: 100,
		}
		cfg := poscalc.DefaultTaxConfig()

		res, err := poscalc.CalculateOrder(items, promo, cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.TotalAmount != 0 {
			t.Errorf("TotalAmount = %v; want 0 for 100%% discount", res.TotalAmount)
		}
		if res.TaxAmount != 0 {
			t.Errorf("TaxAmount = %v; want 0", res.TaxAmount)
		}
	})

	t.Run("Edge Case: Empty Order Items", func(t *testing.T) {
		_, err := poscalc.CalculateOrder([]poscalc.OrderItemInput{}, nil, poscalc.DefaultTaxConfig())
		if err != poscalc.ErrEmptyOrder {
			t.Fatalf("expected ErrEmptyOrder, got %v", err)
		}
	})
}

// TestPOS_ChangeCalculation tests payment reconciliation, exact change, and underpayment prevention.
func TestPOS_ChangeCalculation(t *testing.T) {
	tests := []struct {
		name        string
		totalDue    float64
		amountPaid  float64
		wantChange  float64
		wantErr     error
	}{
		{
			name:       "Exact Cash Payment",
			totalDue:   94600,
			amountPaid: 94600,
			wantChange: 0,
			wantErr:    nil,
		},
		{
			name:       "Excess Cash Payment with Change",
			totalDue:   94600,
			amountPaid: 100000,
			wantChange: 5400,
			wantErr:    nil,
		},
		{
			name:       "Underpayment Rejection",
			totalDue:   94600,
			amountPaid: 90000,
			wantChange: 0,
			wantErr:    poscalc.ErrInsufficientAmount,
		},
		{
			name:       "Negative Payment Amount Rejection",
			totalDue:   50000,
			amountPaid: -10000,
			wantChange: 0,
			wantErr:    poscalc.ErrNegativePayment,
		},
		{
			name:       "Large Transaction Payment",
			totalDue:   1250000,
			amountPaid: 1500000,
			wantChange: 250000,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change, err := poscalc.CalculateChange(tt.totalDue, tt.amountPaid)
			if tt.wantErr != nil {
				if err == nil || err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if change != tt.wantChange {
				t.Errorf("changeDue = %v; want %v", change, tt.wantChange)
			}
		})
	}
}
