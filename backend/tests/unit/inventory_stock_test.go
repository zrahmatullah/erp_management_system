package unit_test

import (
	"testing"

	"cafe-erp-system/backend/pkg/invcalc"
)

// TestInventory_StockStatusClassification tests inventory classification levels and alert flags.
func TestInventory_StockStatusClassification(t *testing.T) {
	tests := []struct {
		name         string
		currentStock float64
		minStock     float64
		maxStock     float64
		wantStatus   invcalc.StockStatus
		wantAlert    bool
	}{
		{
			name:         "Negative Stock -> Out of Stock with Alert",
			currentStock: -5,
			minStock:     10,
			maxStock:     100,
			wantStatus:   invcalc.StatusOutOfStock,
			wantAlert:    true,
		},
		{
			name:         "Zero Stock -> Out of Stock with Alert",
			currentStock: 0,
			minStock:     10,
			maxStock:     100,
			wantStatus:   invcalc.StatusOutOfStock,
			wantAlert:    true,
		},
		{
			name:         "Stock Exactly at Min Threshold -> Low Stock with Alert",
			currentStock: 10,
			minStock:     10,
			maxStock:     100,
			wantStatus:   invcalc.StatusLowStock,
			wantAlert:    true,
		},
		{
			name:         "Stock Below Min Threshold -> Low Stock with Alert",
			currentStock: 4.5,
			minStock:     10,
			maxStock:     100,
			wantStatus:   invcalc.StatusLowStock,
			wantAlert:    true,
		},
		{
			name:         "Healthy Inventory -> Normal Level without Alert",
			currentStock: 50,
			minStock:     10,
			maxStock:     100,
			wantStatus:   invcalc.StatusNormal,
			wantAlert:    false,
		},
		{
			name:         "Stock Exceeding Max Storage -> Overstocked Level",
			currentStock: 120,
			minStock:     10,
			maxStock:     100,
			wantStatus:   invcalc.StatusOverstocked,
			wantAlert:    false,
		},
		{
			name:         "Unbounded Max Stock (0) -> Normal Level",
			currentStock: 250,
			minStock:     10,
			maxStock:     0,
			wantStatus:   invcalc.StatusNormal,
			wantAlert:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, alert := invcalc.ClassifyStockStatus(tt.currentStock, tt.minStock, tt.maxStock)
			if status != tt.wantStatus {
				t.Errorf("status = %v; want %v", status, tt.wantStatus)
			}
			if alert != tt.wantAlert {
				t.Errorf("alert = %v; want %v", alert, tt.wantAlert)
			}
		})
	}
}

// TestInventory_StockDeduction tests stock depletion, negative stock prevention, and fractional quantities.
func TestInventory_StockDeduction(t *testing.T) {
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
			wantStock:     10, // Stock unaltered
			wantErr:       invcalc.ErrInsufficientStock,
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
			name:          "Fractional Ingredient Deduction (e.g. 18g coffee beans)",
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
			wantErr:       invcalc.ErrInvalidQuantity,
		},
		{
			name:          "Invalid Negative Quantity Deduction",
			initialStock:  50,
			deductQty:     -10,
			allowNegative: false,
			wantStock:     50,
			wantErr:       invcalc.ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newStock, err := invcalc.DeductStock(tt.initialStock, tt.deductQty, tt.allowNegative)
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

// TestInventory_RecipeConsumption_BOM tests Bill of Materials aggregation across sold POS items.
func TestInventory_RecipeConsumption_BOM(t *testing.T) {
	// Setup recipes:
	// "prod-latte": 0.018 kg Espresso Beans, 0.200 L Fresh Milk
	// "prod-cappuccino": 0.018 kg Espresso Beans, 0.150 L Fresh Milk
	// "prod-americano": 0.018 kg Espresso Beans
	recipes := map[string][]invcalc.RecipeIngredient{
		"prod-latte": {
			{InventoryItemID: "mat-coffee-beans", QuantityRequired: 0.018, UOM: "kg"},
			{InventoryItemID: "mat-fresh-milk", QuantityRequired: 0.200, UOM: "liter"},
		},
		"prod-cappuccino": {
			{InventoryItemID: "mat-coffee-beans", QuantityRequired: 0.018, UOM: "kg"},
			{InventoryItemID: "mat-fresh-milk", QuantityRequired: 0.150, UOM: "liter"},
		},
		"prod-americano": {
			{InventoryItemID: "mat-coffee-beans", QuantityRequired: 0.018, UOM: "kg"},
		},
	}

	// Sales Order: 2 Lattes, 3 Cappuccinos, 5 Americanos
	// Expected Coffee Beans: (2 * 0.018) + (3 * 0.018) + (5 * 0.018) = 10 * 0.018 = 0.180 kg
	// Expected Fresh Milk: (2 * 0.200) + (3 * 0.150) = 0.400 + 0.450 = 0.850 L
	sales := []invcalc.ProductSaleItem{
		{ProductID: "prod-latte", Quantity: 2},
		{ProductID: "prod-cappuccino", Quantity: 3},
		{ProductID: "prod-americano", Quantity: 5},
	}

	consumed, err := invcalc.CalculateRecipeConsumption(sales, recipes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedMap := map[string]float64{
		"mat-coffee-beans": 0.18,
		"mat-fresh-milk":   0.85,
	}

	if len(consumed) != len(expectedMap) {
		t.Fatalf("consumed items length = %d; want %d", len(consumed), len(expectedMap))
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

// TestInventory_StockOpnameVariance tests physical audit reconciliation and discrepancy tagging.
func TestInventory_StockOpnameVariance(t *testing.T) {
	tests := []struct {
		name         string
		systemQty    float64
		physicalQty  float64
		unitCost     float64
		wantVariance float64
		wantVal      float64
		wantType     invcalc.DiscrepancyType
	}{
		{
			name:         "Balanced Stock Opname",
			systemQty:    100,
			physicalQty:  100,
			unitCost:     15000,
			wantVariance: 0,
			wantVal:      0,
			wantType:     invcalc.DiscrepancyBalanced,
		},
		{
			name:         "Shrinkage / Physical Shortage",
			systemQty:    50,
			physicalQty:  47,
			unitCost:     20000,
			wantVariance: -3,
			wantVal:      -60000,
			wantType:     invcalc.DiscrepancyShrinkage,
		},
		{
			name:         "Surplus / Found Inventory",
			systemQty:    30,
			physicalQty:  35,
			unitCost:     10000,
			wantVariance: 5,
			wantVal:      50000,
			wantType:     invcalc.DiscrepancySurplus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := invcalc.CalculateStockVariance(tt.systemQty, tt.physicalQty, tt.unitCost)
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
