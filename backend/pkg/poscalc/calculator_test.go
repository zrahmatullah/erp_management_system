package poscalc

import (
	"testing"
)

func TestLineItemCalculation(t *testing.T) {
	tests := []struct {
		name      string
		input     OrderItemInput
		wantGross float64
		wantDisc  float64
		wantNet   float64
		wantErr   error
	}{
		{
			name: "Standard Item Without Discount",
			input: OrderItemInput{
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
			input: OrderItemInput{
				ProductID:     "prod-2",
				Name:          "Cafe Latte",
				Quantity:      2,
				UnitPrice:     35000,
				DiscountType:  DiscountTypePercentage,
				DiscountValue: 20,
			},
			wantGross: 70000,
			wantDisc:  14000,
			wantNet:   56000,
			wantErr:   nil,
		},
		{
			name: "Item with Fixed Discount",
			input: OrderItemInput{
				ProductID:     "prod-3",
				Name:          "Croissant",
				Quantity:      3,
				UnitPrice:     20000,
				DiscountType:  DiscountTypeFixed,
				DiscountValue: 15000,
			},
			wantGross: 60000,
			wantDisc:  15000,
			wantNet:   45000,
			wantErr:   nil,
		},
		{
			name: "Fixed Discount Capped at Item Gross Total",
			input: OrderItemInput{
				ProductID:     "prod-4",
				Name:          "Mineral Water",
				Quantity:      1,
				UnitPrice:     10000,
				DiscountType:  DiscountTypeFixed,
				DiscountValue: 15000,
			},
			wantGross: 10000,
			wantDisc:  10000,
			wantNet:   0,
			wantErr:   nil,
		},
		{
			name: "Invalid Zero Quantity",
			input: OrderItemInput{
				ProductID: "prod-5",
				Quantity:  0,
				UnitPrice: 10000,
			},
			wantErr: ErrInvalidQuantity,
		},
		{
			name: "Invalid Negative Quantity",
			input: OrderItemInput{
				ProductID: "prod-6",
				Quantity:  -2,
				UnitPrice: 10000,
			},
			wantErr: ErrInvalidQuantity,
		},
		{
			name: "Invalid Negative Unit Price",
			input: OrderItemInput{
				ProductID: "prod-7",
				Quantity:  1,
				UnitPrice: -5000,
			},
			wantErr: ErrInvalidUnitPrice,
		},
		{
			name: "Invalid Discount Percentage > 100",
			input: OrderItemInput{
				ProductID:     "prod-8",
				Quantity:      1,
				UnitPrice:     50000,
				DiscountType:  DiscountTypePercentage,
				DiscountValue: 120,
			},
			wantErr: ErrInvalidDiscount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CalculateOrderItem(tt.input)
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

func TestOrderCalculation_Scenarios(t *testing.T) {
	t.Run("Standard Dine-In Order with 10% PB1 Tax", func(t *testing.T) {
		items := []OrderItemInput{
			{ProductID: "1", Name: "Nasi Goreng", Quantity: 2, UnitPrice: 35000},
			{ProductID: "2", Name: "Es Teh Manis", Quantity: 2, UnitPrice: 8000},
		}
		cfg := DefaultTaxConfig()

		res, err := CalculateOrder(items, nil, cfg)
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
		items := []OrderItemInput{
			{ProductID: "1", Name: "Sirloin Steak", Quantity: 2, UnitPrice: 120000},
		}
		promo := &PromotionInput{
			Type:        DiscountTypePercentage,
			Value:       50,
			MinOrder:    100000,
			MaxDiscount: 50000,
		}
		cfg := DefaultTaxConfig()

		res, err := CalculateOrder(items, promo, cfg)
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
		items := []OrderItemInput{
			{ProductID: "1", Name: "Americano", Quantity: 1, UnitPrice: 28000},
		}
		promo := &PromotionInput{
			Type:        DiscountTypeFixed,
			Value:       10000,
			MinOrder:    50000,
			MaxDiscount: 10000,
		}
		cfg := DefaultTaxConfig()

		res, err := CalculateOrder(items, promo, cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.OrderDiscount != 0 {
			t.Errorf("OrderDiscount = %v; want 0 (min order not met)", res.OrderDiscount)
		}
		if res.TotalAmount != 30800 {
			t.Errorf("TotalAmount = %v; want 30800", res.TotalAmount)
		}
	})

	t.Run("Order with Service Charge and 11% PPN", func(t *testing.T) {
		items := []OrderItemInput{
			{ProductID: "1", Name: "Premium Buffet", Quantity: 1, UnitPrice: 100000},
		}
		cfg := TaxConfig{
			TaxRate:             0.11,
			ServiceChargeRate:   0.05,
			TaxAppliedAfterDisc: true,
		}

		res, err := CalculateOrder(items, nil, cfg)
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
		items := []OrderItemInput{
			{ProductID: "1", Name: "Birthday Coffee", Quantity: 1, UnitPrice: 35000},
		}
		promo := &PromotionInput{
			Type:  DiscountTypePercentage,
			Value: 100,
		}
		cfg := DefaultTaxConfig()

		res, err := CalculateOrder(items, promo, cfg)
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
		_, err := CalculateOrder([]OrderItemInput{}, nil, DefaultTaxConfig())
		if err != ErrEmptyOrder {
			t.Fatalf("expected ErrEmptyOrder, got %v", err)
		}
	})
}

func TestChangeCalculation(t *testing.T) {
	tests := []struct {
		name       string
		totalDue   float64
		amountPaid float64
		wantChange float64
		wantErr    error
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
			wantErr:    ErrInsufficientAmount,
		},
		{
			name:       "Negative Payment Amount Rejection",
			totalDue:   50000,
			amountPaid: -10000,
			wantChange: 0,
			wantErr:    ErrNegativePayment,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change, err := CalculateChange(tt.totalDue, tt.amountPaid)
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
