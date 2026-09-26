package invcalc

import (
	"testing"
)

func TestStockStatusClassification(t *testing.T) {
	tests := []struct {
		name         string
		currentStock float64
		minStock     float64
		maxStock     float64
		wantStatus   StockStatus
		wantAlert    bool
	}{
		{
			name:         "Negative Stock -> Out of Stock with Alert",
			currentStock: -5,
			minStock:     10,
			maxStock:     100,
			wantStatus:   StatusOutOfStock,
			wantAlert:    true,
		},
		{
			name:         "Zero Stock -> Out of Stock with Alert",
			currentStock: 0,
			minStock:     10,
			maxStock:     100,
			wantStatus:   StatusOutOfStock,
			wantAlert:    true,
		},
		{
			name:         "Stock Exactly at Min Threshold -> Low Stock with Alert",
			currentStock: 10,
			minStock:     10,
			maxStock:     100,
			wantStatus:   StatusLowStock,
			wantAlert:    true,
		},
		{
			name:         "Stock Below Min Threshold -> Low Stock with Alert",
			currentStock: 4.5,
			minStock:     10,
			maxStock:     100,
			wantStatus:   StatusLowStock,
			wantAlert:    true,
		},
		{
			name:         "Healthy Inventory -> Normal Level without Alert",
			currentStock: 50,
			minStock:     10,
			maxStock:     100,
			wantStatus:   StatusNormal,
			wantAlert:    false,
		},
		{
			name:         "Stock Exceeding Max Storage -> Overstocked Level",
			currentStock: 120,
			minStock:     10,
			maxStock:     100,
			wantStatus:   StatusOverstocked,
			wantAlert:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, alert := ClassifyStockStatus(tt.currentStock, tt.minStock, tt.maxStock)
			if status != tt.wantStatus {
				t.Errorf("status = %v; want %v", status, tt.wantStatus)
			}
			if alert != tt.wantAlert {
				t.Errorf("alert = %v; want %v", alert, tt.wantAlert)
			}
		})
	}
}

func TestStockDeduction(t *testing.T) {
	tests := []struct {
		name          string
		initialStock  float64
		deductQty     float64
		allowNegative bool
		wantStock     float64
		wantErr       error
	}{
		{
			name:          "Standard Deduction Happy Path",
			initialStock:  100,
			deductQty:     25,
			allowNegative: false,
			wantStock:     75,
			wantErr:       nil,
		},
		{
			name:          "Exact Depletion to Zero",
			initialStock:  50,
			deductQty:     50,
			allowNegative: false,
			wantStock:     0,
			wantErr:       nil,
		},
		{
			name:          "Minus Stock Prevention (Strict No-Minus)",
			initialStock:  10,
			deductQty:     15,
			allowNegative: false,
			wantStock:     10,
			wantErr:       ErrInsufficientStock,
		},
		{
			name:          "Minus Stock Allowed by Overdraft Policy",
			initialStock:  10,
			deductQty:     15,
			allowNegative: true,
			wantStock:     -5,
			wantErr:       nil,
		},
		{
			name:          "Fractional Ingredient Deduction",
			initialStock:  1.000,
			deductQty:     0.018,
			allowNegative: false,
			wantStock:     0.982,
			wantErr:       nil,
		},
		{
			name:          "Invalid Zero Quantity Deduction",
			initialStock:  50,
			deductQty:     0,
			allowNegative: false,
			wantStock:     50,
			wantErr:       ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newStock, err := DeductStock(tt.initialStock, tt.deductQty, tt.allowNegative)
			if tt.wantErr != nil {
				if err == nil || err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if newStock != tt.wantStock {
				t.Errorf("newStock = %v; want %v", newStock, tt.wantStock)
			}
		})
	}
}

func TestRecipeConsumption_BOM(t *testing.T) {
	recipes := map[string][]RecipeIngredient{
		"prod-latte": {
			{InventoryItemID: "mat-coffee-beans", QuantityRequired: 0.018, UOM: "kg"},
			{InventoryItemID: "mat-fresh-milk", QuantityRequired: 0.200, UOM: "liter"},
		},
		"prod-cappuccino": {
			{InventoryItemID: "mat-coffee-beans", QuantityRequired: 0.018, UOM: "kg"},
			{InventoryItemID: "mat-fresh-milk", QuantityRequired: 0.150, UOM: "liter"},
		},
	}

	sales := []ProductSaleItem{
		{ProductID: "prod-latte", Quantity: 2},
		{ProductID: "prod-cappuccino", Quantity: 3},
	}

	consumed, err := CalculateRecipeConsumption(sales, recipes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedMap := map[string]float64{
		"mat-coffee-beans": 0.09, // (2*0.018) + (3*0.018) = 0.090
		"mat-fresh-milk":   0.85, // (2*0.200) + (3*0.150) = 0.850
	}

	for _, item := range consumed {
		expectedQty, exists := expectedMap[item.InventoryItemID]
		if !exists {
			t.Errorf("unexpected material consumed: %s", item.InventoryItemID)
			continue
		}
		if item.TotalRequired != expectedQty {
			t.Errorf("material %s total = %v; want %v", item.InventoryItemID, item.TotalRequired, expectedQty)
		}
	}
}

func TestStockOpnameVariance(t *testing.T) {
	tests := []struct {
		name         string
		systemQty    float64
		physicalQty  float64
		unitCost     float64
		wantVariance float64
		wantVal      float64
		wantType     DiscrepancyType
	}{
		{
			name:         "Balanced Stock Opname",
			systemQty:    100,
			physicalQty:  100,
			unitCost:     15000,
			wantVariance: 0,
			wantVal:      0,
			wantType:     DiscrepancyBalanced,
		},
		{
			name:         "Shrinkage / Physical Shortage",
			systemQty:    50,
			physicalQty:  47,
			unitCost:     20000,
			wantVariance: -3,
			wantVal:      -60000,
			wantType:     DiscrepancyShrinkage,
		},
		{
			name:         "Surplus / Found Inventory",
			systemQty:    30,
			physicalQty:  35,
			unitCost:     10000,
			wantVariance: 5,
			wantVal:      50000,
			wantType:     DiscrepancySurplus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CalculateStockVariance(tt.systemQty, tt.physicalQty, tt.unitCost)
			if res.VarianceQuantity != tt.wantVariance {
				t.Errorf("VarianceQuantity = %v; want %v", res.VarianceQuantity, tt.wantVariance)
			}
			if res.VarianceValue != tt.wantVal {
				t.Errorf("VarianceValue = %v; want %v", res.VarianceValue, tt.wantVal)
			}
			if res.Discrepancy != tt.wantType {
				t.Errorf("Discrepancy = %v; want %v", res.Discrepancy, tt.wantType)
			}
		})
	}
}
