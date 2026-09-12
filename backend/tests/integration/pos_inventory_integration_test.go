package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"cafe-erp-system/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getBaseURL() string {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return strings.TrimRight(baseURL, "/")
}

// TestIntegration_POS_StockDeduction_TableStatus verifies the end-to-end integration between:
// POS Order Creation -> Table status Occupied -> Payment -> Table status Available -> Stock Deduction -> Stock Movements
func TestIntegration_POS_StockDeduction_TableStatus(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbpool.Close()

	// 1. Pick a sample test product with stock > 0
	var productID, productName string
	var initialStock float64
	err = dbpool.QueryRow(ctx, "SELECT id, name, stock FROM products WHERE is_active = true AND stock > 5 ORDER BY name LIMIT 1").Scan(&productID, &productName, &initialStock)
	if err != nil {
		t.Fatalf("No active product found for integration testing: %v", err)
	}
	t.Logf("Selected Product for Test: %s (ID: %s) with Initial Stock: %.2f", productName, productID, initialStock)

	// 2. Pick a test table and ensure it starts as 'available'
	var tableID, tableNumber string
	err = dbpool.QueryRow(ctx, "SELECT id, table_number FROM cafe_tables ORDER BY table_number LIMIT 1").Scan(&tableID, &tableNumber)
	if err != nil {
		t.Fatalf("No table found for testing: %v", err)
	}
	_, _ = dbpool.Exec(ctx, "UPDATE cafe_tables SET status = 'available' WHERE id = $1", tableID)

	client := &http.Client{Timeout: 8 * time.Second}

	// 3. Create POS Order for this table
	orderQty := 2
	orderPayload := map[string]interface{}{
		"order_type":    "dine_in",
		"table_number":  tableNumber,
		"customer_name": "Integration Test Customer",
		"notes":         "Regression Integration Test",
		"items": []map[string]interface{}{
			{
				"product_id": productID,
				"name":       productName,
				"quantity":   orderQty,
				"unit_price": 25000,
				"station":    "barista",
			},
		},
	}
	bodyBytes, _ := json.Marshal(orderPayload)
	reqOrder, _ := http.NewRequest(http.MethodPost, getBaseURL()+"/api/v1/pos/orders", bytes.NewBuffer(bodyBytes))
	reqOrder.Header.Set("Content-Type", "application/json")
	respOrder, err := client.Do(reqOrder)
	if err != nil {
		t.Fatalf("Failed to create POS order via API: %v", err)
	}
	defer respOrder.Body.Close()

	if respOrder.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 Created for POS order, got %d", respOrder.StatusCode)
	}

	var orderResp struct {
		Success     bool    `json:"success"`
		OrderID     string  `json:"order_id"`
		OrderNumber string  `json:"order_number"`
		QueueNumber string  `json:"queue_number"`
		Total       float64 `json:"total"`
	}
	if err := json.NewDecoder(respOrder.Body).Decode(&orderResp); err != nil {
		t.Fatalf("Failed to decode order response: %v", err)
	}
	t.Logf("Order Created: #%s, OrderID: %s, Queue: %s, Total: %.2f", orderResp.OrderNumber, orderResp.OrderID, orderResp.QueueNumber, orderResp.Total)

	// 4. Verify table status changed to 'occupied'
	var tableStatus string
	err = dbpool.QueryRow(ctx, "SELECT status FROM cafe_tables WHERE id = $1", tableID).Scan(&tableStatus)
	if err != nil {
		t.Fatalf("Failed to query table status: %v", err)
	}
	if tableStatus != "occupied" {
		t.Errorf("Expected table status 'occupied', got '%s'", tableStatus)
	} else {
		t.Logf("✅ Table %s correctly transitioned to 'occupied'", tableNumber)
	}

	// 5. Pay POS Order
	payPayload := map[string]interface{}{
		"payment_method": "cash",
		"amount_paid":    orderResp.Total + 10000,
		"total_amount":   orderResp.Total,
	}
	payBytes, _ := json.Marshal(payPayload)
	reqPay, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/pos/orders/%s/pay", getBaseURL(), orderResp.OrderID), bytes.NewBuffer(payBytes))
	reqPay.Header.Set("Content-Type", "application/json")
	respPay, err := client.Do(reqPay)
	if err != nil {
		t.Fatalf("Failed to execute payment via API: %v", err)
	}
	defer respPay.Body.Close()

	if respPay.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK for payment, got %d", respPay.StatusCode)
	}

	// 6. Verify table status reverted back to 'available'
	err = dbpool.QueryRow(ctx, "SELECT status FROM cafe_tables WHERE id = $1", tableID).Scan(&tableStatus)
	if err != nil {
		t.Fatalf("Failed to query table status after payment: %v", err)
	}
	if tableStatus != "available" {
		t.Errorf("Expected table status 'available' after payment, got '%s'", tableStatus)
	} else {
		t.Logf("✅ Table %s correctly released back to 'available'", tableNumber)
	}

	// 7. Verify product stock was reduced by exact ordered quantity
	var stockAfter float64
	err = dbpool.QueryRow(ctx, "SELECT stock FROM products WHERE id = $1", productID).Scan(&stockAfter)
	if err != nil {
		t.Fatalf("Failed to query product stock after payment: %v", err)
	}
	expectedStock := initialStock - float64(orderQty)
	if stockAfter != expectedStock {
		t.Errorf("Expected product stock %.2f, but got %.2f", expectedStock, stockAfter)
	} else {
		t.Logf("✅ Product stock correctly deducted: %.2f -> %.2f (-%d)", initialStock, stockAfter, orderQty)
	}

	// 8. Verify stock movement recorded in stock_movements ledger
	var movementCount int
	err = dbpool.QueryRow(ctx, `
		SELECT COUNT(*) FROM stock_movements 
		WHERE product_id = $1 AND reference_id = $2 AND type = 'out_pos_sales'`,
		productID, orderResp.OrderID).Scan(&movementCount)
	if err != nil {
		t.Fatalf("Failed to query stock_movements: %v", err)
	}
	if movementCount == 0 {
		t.Errorf("Expected at least 1 stock_movements record for order %s, found %d", orderResp.OrderID, movementCount)
	} else {
		t.Logf("✅ Stock movement record verified in Kartu Stok ledger (Count: %d)", movementCount)
	}
}

// TestIntegration_PO_Receipt_StockAddition verifies procurement flow:
// PO Created -> Status updated to 'received' -> Stock added -> Stock movement recorded in Kartu Stok
func TestIntegration_PO_Receipt_StockAddition(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbpool.Close()

	// 1. Pick active supplier
	var supplierID, supplierName string
	err = dbpool.QueryRow(ctx, "SELECT id, name FROM suppliers WHERE is_active = true LIMIT 1").Scan(&supplierID, &supplierName)
	if err != nil {
		t.Fatalf("No supplier found for testing: %v", err)
	}

	// 2. Pick inventory item
	var itemID, itemName string
	err = dbpool.QueryRow(ctx, "SELECT id, name FROM inventory_items WHERE is_active = true LIMIT 1").Scan(&itemID, &itemName)
	if err != nil {
		t.Fatalf("No inventory item found for testing: %v", err)
	}

	client := &http.Client{Timeout: 8 * time.Second}

	// 3. Create PO
	poQty := 5.0
	poPayload := map[string]interface{}{
		"supplier_id":   supplierID,
		"notes":         "Integration Test PO",
		"expected_date": time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		"items": []map[string]interface{}{
			{
				"inventory_item_id": itemID,
				"quantity":          poQty,
				"unit_price":        15000,
			},
		},
	}
	poBytes, _ := json.Marshal(poPayload)
	reqPO, _ := http.NewRequest(http.MethodPost, getBaseURL()+"/api/v1/inventory/purchase-orders", bytes.NewBuffer(poBytes))
	reqPO.Header.Set("Content-Type", "application/json")
	respPO, err := client.Do(reqPO)
	if err != nil {
		t.Fatalf("Failed to create PO via API: %v", err)
	}
	defer respPO.Body.Close()

	if respPO.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 Created for PO, got %d", respPO.StatusCode)
	}

	var poResp struct {
		Data struct {
			ID       string `json:"id"`
			PONumber string `json:"po_number"`
			Status   string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(respPO.Body).Decode(&poResp); err != nil {
		t.Fatalf("Failed to decode PO response: %v", err)
	}
	t.Logf("PO Created: %s (ID: %s)", poResp.Data.PONumber, poResp.Data.ID)

	// 4. Update PO status to 'received'
	statusPayload := map[string]string{"status": "received"}
	statusBytes, _ := json.Marshal(statusPayload)
	reqStatus, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v1/inventory/purchase-orders/%s/status", getBaseURL(), poResp.Data.ID), bytes.NewBuffer(statusBytes))
	reqStatus.Header.Set("Content-Type", "application/json")
	respStatus, err := client.Do(reqStatus)
	if err != nil {
		t.Fatalf("Failed to update PO status via API: %v", err)
	}
	defer respStatus.Body.Close()

	if respStatus.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK for PO status update, got %d", respStatus.StatusCode)
	}

	// 5. Verify stock movement recorded for PO receipt
	var poMovementCount int
	err = dbpool.QueryRow(ctx, `
		SELECT COUNT(*) FROM stock_movements 
		WHERE inventory_item_id = $1 AND reference_id = $2 AND type = 'in_purchase'`,
		itemID, poResp.Data.ID).Scan(&poMovementCount)
	if err != nil {
		t.Fatalf("Failed to query stock_movements for PO: %v", err)
	}
	if poMovementCount == 0 {
		t.Errorf("Expected stock movement with type 'in_purchase' for PO %s, found %d", poResp.Data.ID, poMovementCount)
	} else {
		t.Logf("✅ PO Receipt stock movement verified in Kartu Stok (Count: %d)", poMovementCount)
	}
}
