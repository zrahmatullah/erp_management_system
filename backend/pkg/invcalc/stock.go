package invcalc

import (
	"errors"
	"math"
)

var (
	ErrInvalidQuantity    = errors.New("deduction quantity must be greater than zero")
	ErrInsufficientStock  = errors.New("insufficient inventory stock for deduction")
	ErrNegativeStockParam = errors.New("min or max stock parameters cannot be negative")
)

type StockStatus string

const (
	StatusOutOfStock  StockStatus = "out_of_stock"
	StatusLowStock    StockStatus = "low_stock"
	StatusNormal      StockStatus = "normal"
	StatusOverstocked StockStatus = "overstocked"
)

// ClassifyStockStatus evaluates inventory levels against threshold boundaries.
func ClassifyStockStatus(currentStock, minStock, maxStock float64) (StockStatus, bool) {
	// Returns (status, isLowStockAlert)
	if currentStock <= 0 {
		return StatusOutOfStock, true
	}
	if currentStock <= minStock {
		return StatusLowStock, true
	}
	if maxStock > 0 && currentStock > maxStock {
		return StatusOverstocked, false
	}
	return StatusNormal, false
}

// DeductStock validates and applies stock deductions, preventing negative stocks
// unless allowNegative is explicitly enabled for specific business flows.
func DeductStock(currentStock, quantity float64, allowNegative bool) (float64, error) {
	if quantity <= 0 {
		return currentStock, ErrInvalidQuantity
	}

	if currentStock < quantity && !allowNegative {
		return currentStock, ErrInsufficientStock
	}

	newStock := currentStock - quantity
	// Precision cleanup for floating point representation (e.g. 0.05 kg)
	newStock = math.Round(newStock*10000) / 10000

	if newStock < 0 && !allowNegative {
		newStock = 0
	}

	return newStock, nil
}

// ProductSaleItem represents a sold POS product.
type ProductSaleItem struct {
	ProductID string
	Quantity  int
}

// RecipeIngredient represents a raw material required for a product.
type RecipeIngredient struct {
	InventoryItemID  string
	QuantityRequired float64 // Quantity of raw material per 1 unit of product
	UOM              string
}

// MaterialConsumption aggregates required raw materials across multiple sold items.
type MaterialConsumption struct {
	InventoryItemID string  `json:"inventory_item_id"`
	TotalRequired   float64 `json:"total_required"`
}

// CalculateRecipeConsumption multiplies sold items with their respective recipes
// and aggregates total material consumption for multi-ingredient deductions.
func CalculateRecipeConsumption(sales []ProductSaleItem, recipes map[string][]RecipeIngredient) ([]MaterialConsumption, error) {
	totals := make(map[string]float64)

	for _, sale := range sales {
		if sale.Quantity <= 0 {
			continue
		}

		ingredients, hasRecipe := recipes[sale.ProductID]
		if !hasRecipe {
			continue
		}

		for _, ing := range ingredients {
			if ing.QuantityRequired <= 0 {
				continue
			}
			needed := float64(sale.Quantity) * ing.QuantityRequired
			totals[ing.InventoryItemID] += needed
		}
	}

	var results []MaterialConsumption
	for itemID, qty := range totals {
		results = append(results, MaterialConsumption{
			InventoryItemID: itemID,
			TotalRequired:   math.Round(qty*10000) / 10000,
		})
	}

	return results, nil
}

// DiscrepancyType identifies stock opname reconciliation results.
type DiscrepancyType string

const (
	DiscrepancyBalanced  DiscrepancyType = "balanced"
	DiscrepancySurplus   DiscrepancyType = "surplus"
	DiscrepancyShrinkage DiscrepancyType = "shrinkage"
)

// StockVarianceResult holds physical vs system stock reconciliation details.
type StockVarianceResult struct {
	SystemQuantity   float64         `json:"system_quantity"`
	PhysicalQuantity float64         `json:"physical_quantity"`
	VarianceQuantity float64         `json:"variance_quantity"`
	VarianceValue    float64         `json:"variance_value"`
	Discrepancy      DiscrepancyType `json:"discrepancy"`
}

// CalculateStockVariance computes discrepancy between physical audit and digital records.
func CalculateStockVariance(systemQty, physicalQty, unitCost float64) StockVarianceResult {
	variance := physicalQty - systemQty
	variance = math.Round(variance*10000) / 10000

	val := variance * unitCost
	val = math.Round(val)

	var discType DiscrepancyType
	if variance == 0 {
		discType = DiscrepancyBalanced
	} else if variance > 0 {
		discType = DiscrepancySurplus
	} else {
		discType = DiscrepancyShrinkage
	}

	return StockVarianceResult{
		SystemQuantity:   systemQty,
		PhysicalQuantity: physicalQty,
		VarianceQuantity: variance,
		VarianceValue:    val,
		Discrepancy:      discType,
	}
}
