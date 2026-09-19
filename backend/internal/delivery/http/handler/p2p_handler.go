package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type P2PHandler struct {
	db *pgxpool.Pool
}

func NewP2PHandler(db *pgxpool.Pool) *P2PHandler {
	return &P2PHandler{db: db}
}

// -----------------------------------------------------------------------------
// 1. PURCHASE REQUISITIONS (PR)
// -----------------------------------------------------------------------------

func (h *P2PHandler) GetPurchaseRequisitions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		SELECT 
			pr.id, pr.pr_number, pr.department, pr.status, pr.required_date, 
			COALESCE(pr.notes, '') as notes, pr.total_estimated_cost, pr.created_at,
			COALESCE(u.full_name, 'Staff') as requester_name,
			COALESCE(
				json_agg(
					json_build_object(
						'id', pri.id,
						'inventory_item_id', pri.inventory_item_id,
						'name', ii.name,
						'uom', ii.uom,
						'quantity', pri.quantity,
						'estimated_unit_price', pri.estimated_unit_price,
						'estimated_total_price', pri.estimated_total_price,
						'notes', COALESCE(pri.notes, '')
					)
				) FILTER (WHERE pri.id IS NOT NULL), '[]'
			) as items
		FROM purchase_requisitions pr
		LEFT JOIN users u ON pr.requested_by = u.id
		LEFT JOIN purchase_requisition_items pri ON pr.id = pri.purchase_requisition_id
		LEFT JOIN inventory_items ii ON pri.inventory_item_id = ii.id
		WHERE pr.deleted_at IS NULL
		GROUP BY pr.id, u.full_name
		ORDER BY pr.created_at DESC`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, prNumber, department, status, notes, requesterName, itemsJSON string
		var reqDate time.Time
		var createdAt time.Time
		var totalEst float64

		if err := rows.Scan(&id, &prNumber, &department, &status, &reqDate, &notes, &totalEst, &createdAt, &requesterName, &itemsJSON); err == nil {
			var items []map[string]interface{}
			_ = json.Unmarshal([]byte(itemsJSON), &items)
			if items == nil {
				items = []map[string]interface{}{}
			}

			list = append(list, map[string]interface{}{
				"id":                   id,
				"pr_number":            prNumber,
				"department":           department,
				"status":               status,
				"required_date":        reqDate.Format("2006-01-02"),
				"notes":                notes,
				"total_estimated_cost": totalEst,
				"requester_name":       requesterName,
				"created_at":           createdAt.Format("2006-01-02 15:04"),
				"items":                items,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

func (h *P2PHandler) CreatePurchaseRequisition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		Department   string `json:"department"`
		RequiredDate string `json:"required_date"`
		Notes        string `json:"notes"`
		Items        []struct {
			InventoryItemID    string  `json:"inventory_item_id"`
			Quantity           float64 `json:"quantity"`
			EstimatedUnitPrice float64 `json:"estimated_unit_price"`
			Notes              string  `json:"notes"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(body.Items) == 0 {
		writeError(w, http.StatusBadRequest, "Setidaknya harus ada 1 item dalam permintaan pembelian")
		return
	}

	if body.Department == "" {
		body.Department = "Kitchen & Bar"
	}

	reqDate := time.Now().Add(3 * 24 * time.Hour)
	if body.RequiredDate != "" {
		if t, err := time.Parse("2006-01-02", body.RequiredDate); err == nil {
			reqDate = t
		}
	}

	var branchID string
	_ = h.db.QueryRow(ctx, "SELECT id FROM branches LIMIT 1").Scan(&branchID)

	prID := uuid.New().String()
	prNumber := fmt.Sprintf("PR-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	var totalEst float64
	for _, it := range body.Items {
		totalEst += it.Quantity * it.EstimatedUnitPrice
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO purchase_requisitions (id, branch_id, pr_number, department, status, required_date, notes, total_estimated_cost, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'pending_approval', $5, $6, $7, NOW(), NOW())`,
		prID, branchID, prNumber, body.Department, reqDate, body.Notes, totalEst)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, it := range body.Items {
		estTotal := it.Quantity * it.EstimatedUnitPrice
		_, err = tx.Exec(ctx, `
			INSERT INTO purchase_requisition_items (id, purchase_requisition_id, inventory_item_id, quantity, estimated_unit_price, estimated_total_price, notes, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
			uuid.New().String(), prID, it.InventoryItemID, it.Quantity, it.EstimatedUnitPrice, estTotal, it.Notes)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":   true,
		"message":   fmt.Sprintf("Purchase Requisition %s berhasil dibuat", prNumber),
		"pr_id":     prID,
		"pr_number": prNumber,
		"status":    "pending_approval",
	})
}

func (h *P2PHandler) UpdatePurchaseRequisitionStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var body struct {
		Status string `json:"status"` // 'approved', 'rejected'
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	_, err := h.db.Exec(ctx, "UPDATE purchase_requisitions SET status = $1, approved_at = NOW(), updated_at = NOW() WHERE id = $2", body.Status, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Status Purchase Requisition berhasil diubah menjadi %s", body.Status),
	})
}

func (h *P2PHandler) ConvertPRToPO(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	prID := chi.URLParam(r, "id")

	var body struct {
		SupplierID   string  `json:"supplier_id"`
		TaxType      string  `json:"tax_type"`      // 'include', 'exclude', 'non_pkp'
		PaymentTerms string  `json:"payment_terms"` // 'cod', 'net_14', 'net_30'
		ExpectedDate string  `json:"expected_date"`
		Notes        string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if body.SupplierID == "" {
		writeError(w, http.StatusBadRequest, "Supplier wajib dipilih untuk menerbitkan PO")
		return
	}

	// Verify PR status is approved
	var prStatus, prNumber, branchID string
	err := h.db.QueryRow(ctx, "SELECT status, pr_number, branch_id FROM purchase_requisitions WHERE id = $1", prID).Scan(&prStatus, &prNumber, &branchID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Purchase Requisition tidak ditemukan")
		return
	}

	if prStatus != "approved" && prStatus != "pending_approval" {
		writeError(w, http.StatusBadRequest, "Hanya PR berstatus approved atau pending yang dapat diubah ke PO")
		return
	}

	// Fetch items from PR
	itemRows, err := h.db.Query(ctx, "SELECT inventory_item_id, quantity, estimated_unit_price FROM purchase_requisition_items WHERE purchase_requisition_id = $1", prID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer itemRows.Close()

	type prItem struct {
		itemID    string
		qty       float64
		unitPrice float64
	}
	var items []prItem
	var dppAmount float64
	for itemRows.Next() {
		var it prItem
		if err := itemRows.Scan(&it.itemID, &it.qty, &it.unitPrice); err == nil {
			dppAmount += it.qty * it.unitPrice
			items = append(items, it)
		}
	}

	if len(items) == 0 {
		writeError(w, http.StatusBadRequest, "PR tidak memiliki item barang")
		return
	}

	// Calculate Tax (PPN 11%)
	var taxRate float64 = 11.00
	var taxAmount float64 = 0.00
	var totalAmount float64 = dppAmount

	if body.TaxType == "include" {
		dppAmount = totalAmount / 1.11
		taxAmount = totalAmount - dppAmount
	} else if body.TaxType == "exclude" || body.TaxType == "" {
		body.TaxType = "exclude"
		taxAmount = dppAmount * 0.11
		totalAmount = dppAmount + taxAmount
	} else {
		body.TaxType = "non_pkp"
		taxRate = 0
		taxAmount = 0
		totalAmount = dppAmount
	}

	expDate := time.Now().Add(7 * 24 * time.Hour)
	if body.ExpectedDate != "" {
		if t, err := time.Parse("2006-01-02", body.ExpectedDate); err == nil {
			expDate = t
		}
	}

	poID := uuid.New().String()
	poNumber := fmt.Sprintf("PO-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	poNotes := fmt.Sprintf("Generated from %s. %s", prNumber, body.Notes)
	_, err = tx.Exec(ctx, `
		INSERT INTO purchase_orders (id, branch_id, supplier_id, pr_id, po_number, status, dpp_amount, tax_type, tax_rate, tax_amount, total_amount, payment_terms, notes, ordered_date, expected_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 'sent', $6, $7, $8, $9, $10, $11, $12, CURRENT_DATE, $13, NOW(), NOW())`,
		poID, branchID, body.SupplierID, prID, poNumber, dppAmount, body.TaxType, taxRate, taxAmount, totalAmount, body.PaymentTerms, poNotes, expDate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, it := range items {
		itemTotal := it.qty * it.unitPrice
		_, err = tx.Exec(ctx, `
			INSERT INTO purchase_order_items (id, purchase_order_id, inventory_item_id, quantity, unit_price, total_price, quantity_received, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, 0, NOW())`,
			uuid.New().String(), poID, it.itemID, it.qty, it.unitPrice, itemTotal)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Update PR status
	_, _ = tx.Exec(ctx, "UPDATE purchase_requisitions SET status = 'converted_to_po', updated_at = NOW() WHERE id = $1", prID)

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":      true,
		"message":      fmt.Sprintf("PR %s berhasil diubah menjadi PO resmi %s (Status: Sent ke Supplier)", prNumber, poNumber),
		"po_id":        poID,
		"po_number":    poNumber,
		"total_amount": totalAmount,
	})
}

// -----------------------------------------------------------------------------
// 2. GOODS RECEIPT NOTE (GRN) & ACCRUAL ACCOUNTING
// -----------------------------------------------------------------------------

func (h *P2PHandler) GetGoodsReceiptNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		SELECT 
			grn.id, grn.grn_number, grn.delivery_order_number, grn.received_date, grn.status,
			COALESCE(grn.notes, ''), grn.created_at,
			po.po_number, s.name as supplier_name, w.name as warehouse_name,
			COALESCE(
				json_agg(
					json_build_object(
						'id', grni.id,
						'inventory_item_id', grni.inventory_item_id,
						'item_name', ii.name,
						'uom', ii.uom,
						'quantity_ordered', grni.quantity_ordered,
						'quantity_received', grni.quantity_received,
						'quantity_rejected', grni.quantity_rejected,
						'condition_status', grni.condition_status,
						'notes', COALESCE(grni.notes, '')
					)
				) FILTER (WHERE grni.id IS NOT NULL), '[]'
			) as items
		FROM goods_receipt_notes grn
		JOIN purchase_orders po ON grn.purchase_order_id = po.id
		JOIN suppliers s ON po.supplier_id = s.id
		JOIN warehouses w ON grn.warehouse_id = w.id
		LEFT JOIN goods_receipt_note_items grni ON grn.id = grni.grn_id
		LEFT JOIN inventory_items ii ON grni.inventory_item_id = ii.id
		WHERE grn.deleted_at IS NULL
		GROUP BY grn.id, po.po_number, s.name, w.name
		ORDER BY grn.received_date DESC, grn.created_at DESC`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, grnNumber, doNumber, status, notes, poNumber, supplierName, warehouseName, itemsJSON string
		var recDate, createdAt time.Time

		if err := rows.Scan(&id, &grnNumber, &doNumber, &recDate, &status, &notes, &createdAt, &poNumber, &supplierName, &warehouseName, &itemsJSON); err == nil {
			var items []map[string]interface{}
			_ = json.Unmarshal([]byte(itemsJSON), &items)
			if items == nil {
				items = []map[string]interface{}{}
			}

			list = append(list, map[string]interface{}{
				"id":                    id,
				"grn_number":            grnNumber,
				"delivery_order_number": doNumber,
				"received_date":         recDate.Format("2006-01-02"),
				"status":                status,
				"notes":                 notes,
				"po_number":             poNumber,
				"supplier_name":         supplierName,
				"warehouse_name":        warehouseName,
				"created_at":            createdAt.Format("2006-01-02 15:04"),
				"items":                 items,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

func (h *P2PHandler) CreateGoodsReceiptNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		PurchaseOrderID     string `json:"purchase_order_id"`
		DeliveryOrderNumber string `json:"delivery_order_number"`
		WarehouseID         string `json:"warehouse_id"`
		ReceivedDate        string `json:"received_date"`
		Notes               string `json:"notes"`
		Items               []struct {
			POItemID         string  `json:"po_item_id"`
			InventoryItemID  string  `json:"inventory_item_id"`
			QuantityOrdered  float64 `json:"quantity_ordered"`
			QuantityReceived float64 `json:"quantity_received"`
			QuantityRejected float64 `json:"quantity_rejected"`
			ConditionStatus  string  `json:"condition_status"` // 'good', 'damaged', 'expired'
			Notes            string  `json:"notes"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if body.PurchaseOrderID == "" || body.DeliveryOrderNumber == "" || len(body.Items) == 0 {
		writeError(w, http.StatusBadRequest, "Nomor PO, Nomor Surat Jalan (DO), dan item penerimaan wajib diisi")
		return
	}

	var poNumber, branchID string
	var totalPOAmount float64
	err := h.db.QueryRow(ctx, "SELECT po_number, branch_id, total_amount FROM purchase_orders WHERE id = $1", body.PurchaseOrderID).Scan(&poNumber, &branchID, &totalPOAmount)
	if err != nil {
		writeError(w, http.StatusNotFound, "Dokumen Purchase Order tidak ditemukan")
		return
	}

	if body.WarehouseID == "" {
		_ = h.db.QueryRow(ctx, "SELECT id FROM warehouses WHERE branch_id = $1 ORDER BY type = 'main' DESC LIMIT 1", branchID).Scan(&body.WarehouseID)
	}

	recDate := time.Now()
	if body.ReceivedDate != "" {
		if t, err := time.Parse("2006-01-02", body.ReceivedDate); err == nil {
			recDate = t
		}
	}

	grnID := uuid.New().String()
	grnNumber := fmt.Sprintf("GRN-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	// 1. Insert GRN
	_, err = tx.Exec(ctx, `
		INSERT INTO goods_receipt_notes (id, branch_id, warehouse_id, purchase_order_id, grn_number, delivery_order_number, received_date, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'received_full', $8, NOW(), NOW())`,
		grnID, branchID, body.WarehouseID, body.PurchaseOrderID, grnNumber, body.DeliveryOrderNumber, recDate, body.Notes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var totalReceivedValue float64

	// 2. Process Items, Stock Increment, and Kartu Stok
	for _, it := range body.Items {
		cond := it.ConditionStatus
		if cond == "" {
			cond = "good"
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO goods_receipt_note_items (id, grn_id, po_item_id, inventory_item_id, quantity_ordered, quantity_received, quantity_rejected, condition_status, notes, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())`,
			uuid.New().String(), grnID, it.POItemID, it.InventoryItemID, it.QuantityOrdered, it.QuantityReceived, it.QuantityRejected, cond, it.Notes)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Update quantity_received on PO Item if matched
		if it.POItemID != "" {
			_, _ = tx.Exec(ctx, "UPDATE purchase_order_items SET quantity_received = quantity_received + $1 WHERE id = $2", it.QuantityReceived, it.POItemID)
		}

		// Upsert into inventory_stocks
		var newStock float64
		var exists bool
		_ = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM inventory_stocks WHERE inventory_item_id = $1 AND warehouse_id = $2)", it.InventoryItemID, body.WarehouseID).Scan(&exists)
		if exists {
			_ = tx.QueryRow(ctx, `
				UPDATE inventory_stocks
				SET quantity = quantity + $1, updated_at = NOW()
				WHERE inventory_item_id = $2 AND warehouse_id = $3
				RETURNING quantity`, it.QuantityReceived, it.InventoryItemID, body.WarehouseID).Scan(&newStock)
		} else {
			newStock = it.QuantityReceived
			_, _ = tx.Exec(ctx, `
				INSERT INTO inventory_stocks (id, inventory_item_id, warehouse_id, quantity, created_at, updated_at)
				VALUES ($1, $2, $3, $4, NOW(), NOW())`,
				uuid.New().String(), it.InventoryItemID, body.WarehouseID, newStock)
		}

		// Look up unit price for received inventory item
		var itemName string
		var avgCost float64
		_ = tx.QueryRow(ctx, "SELECT name, average_cost FROM inventory_items WHERE id = $1", it.InventoryItemID).Scan(&itemName, &avgCost)

		totalReceivedValue += it.QuantityReceived * avgCost

		// Record in stock_movements (Kartu Stok)
		movementRemark := fmt.Sprintf("Penerimaan GRN #%s (DO: %s, PO: #%s) - %s (+%.2f %s)", grnNumber, body.DeliveryOrderNumber, poNumber, itemName, it.QuantityReceived, cond)
		_, _ = tx.Exec(ctx, `
			INSERT INTO stock_movements (id, inventory_item_id, warehouse_id, type, quantity, balance_after, reference_id, reference_type, remarks, created_at)
			VALUES ($1, $2, $3, 'in_purchase', $4, $5, $6, 'goods_receipt', $7, NOW())`,
			uuid.New().String(), it.InventoryItemID, body.WarehouseID, it.QuantityReceived, newStock, grnID, movementRemark)
	}

	// 3. Update PO status to received
	_, _ = tx.Exec(ctx, "UPDATE purchase_orders SET status = 'received', received_date = NOW(), updated_at = NOW() WHERE id = $1", body.PurchaseOrderID)

	// 4. Create Accrual Journal Entry:
	// [Debit] Persediaan Bahan Baku (11301)
	// [Credit] Hutang Belum Difakturkan / Unbilled AP (21201)
	if totalReceivedValue <= 0 {
		totalReceivedValue = totalPOAmount
	}

	var accPersediaanID, accUnbilledAPID string
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE code = '11301' LIMIT 1").Scan(&accPersediaanID)
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE code = '21201' LIMIT 1").Scan(&accUnbilledAPID)

	if accPersediaanID != "" && accUnbilledAPID != "" {
		journalID := uuid.New().String()
		journalRef := fmt.Sprintf("JRN-GRN-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)
		journalDesc := fmt.Sprintf("Akrual Penerimaan Barang GRN #%s (PO #%s, DO #%s)", grnNumber, poNumber, body.DeliveryOrderNumber)

		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entries (id, branch_id, reference_number, entry_date, description, status, created_at, updated_at)
			VALUES ($1, $2, $3, CURRENT_DATE, $4, 'posted', NOW(), NOW())`,
			journalID, branchID, journalRef, journalDesc)

		// Debit Persediaan
		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
			VALUES ($1, $2, $3, 'Persediaan Bahan Masuk', $4, 0.00, NOW())`,
			uuid.New().String(), journalID, accPersediaanID, totalReceivedValue)

		// Credit Unbilled AP
		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
			VALUES ($1, $2, $3, 'Akrual Hutang Belum Difakturkan', 0.00, $4, NOW())`,
			uuid.New().String(), journalID, accUnbilledAPID, totalReceivedValue)
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":       true,
		"message":       fmt.Sprintf("Penerimaan Barang GRN %s berhasil disimpan, stok gudang bertambah & jurnal akrual tercatat", grnNumber),
		"grn_id":        grnID,
		"grn_number":    grnNumber,
		"accrued_value": totalReceivedValue,
	})
}

// -----------------------------------------------------------------------------
// 3. VENDOR INVOICES & 3-WAY MATCHING (FAKTUR PAJAK & AP RECOGNITION)
// -----------------------------------------------------------------------------

func (h *P2PHandler) GetVendorInvoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		SELECT 
			vi.id, vi.invoice_number, COALESCE(vi.tax_invoice_number, '-') as tax_invoice,
			vi.invoice_date, vi.due_date, vi.subtotal, vi.tax_amount, vi.total_amount,
			vi.paid_amount, vi.status, COALESCE(vi.notes, ''), vi.created_at,
			po.po_number, s.name as supplier_name
		FROM vendor_invoices vi
		JOIN purchase_orders po ON vi.purchase_order_id = po.id
		JOIN suppliers s ON po.supplier_id = s.id
		WHERE vi.deleted_at IS NULL
		ORDER BY vi.created_at DESC`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, invNum, taxInv, status, notes, poNum, supplierName string
		var invDate, dueDate, createdAt time.Time
		var subtotal, taxAmt, totalAmt, paidAmt float64

		if err := rows.Scan(&id, &invNum, &taxInv, &invDate, &dueDate, &subtotal, &taxAmt, &totalAmt, &paidAmt, &status, &notes, &createdAt, &poNum, &supplierName); err == nil {
			list = append(list, map[string]interface{}{
				"id":                 id,
				"invoice_number":     invNum,
				"tax_invoice_number": taxInv,
				"invoice_date":       invDate.Format("2006-01-02"),
				"due_date":           dueDate.Format("2006-01-02"),
				"subtotal":           subtotal,
				"tax_amount":         taxAmt,
				"total_amount":       totalAmt,
				"paid_amount":        paidAmt,
				"remaining_amount":   totalAmt - paidAmt,
				"status":             status,
				"notes":              notes,
				"po_number":          poNum,
				"supplier_name":      supplierName,
				"created_at":         createdAt.Format("2006-01-02 15:04"),
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

func (h *P2PHandler) CreateVendorInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		PurchaseOrderID  string  `json:"purchase_order_id"`
		GRNID            string  `json:"grn_id"`
		InvoiceNumber    string  `json:"invoice_number"`
		TaxInvoiceNumber string  `json:"tax_invoice_number"`
		InvoiceDate      string  `json:"invoice_date"`
		DueDate          string  `json:"due_date"`
		Subtotal         float64 `json:"subtotal"`
		TaxAmount        float64 `json:"tax_amount"`
		TotalAmount      float64 `json:"total_amount"`
		Notes            string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if body.PurchaseOrderID == "" || body.InvoiceNumber == "" {
		writeError(w, http.StatusBadRequest, "Nomor PO dan Nomor Invoice Vendor wajib diisi")
		return
	}

	var poNumber, branchID string
	var poTotal, poDPP, poTax float64
	err := h.db.QueryRow(ctx, "SELECT po_number, branch_id, total_amount, dpp_amount, tax_amount FROM purchase_orders WHERE id = $1", body.PurchaseOrderID).Scan(&poNumber, &branchID, &poTotal, &poDPP, &poTax)
	if err != nil {
		writeError(w, http.StatusNotFound, "Purchase Order tidak ditemukan")
		return
	}

	if body.TotalAmount <= 0 {
		body.TotalAmount = poTotal
		body.Subtotal = poDPP
		body.TaxAmount = poTax
	}

	invDate := time.Now()
	if body.InvoiceDate != "" {
		if t, err := time.Parse("2006-01-02", body.InvoiceDate); err == nil {
			invDate = t
		}
	}

	dueDate := invDate.Add(30 * 24 * time.Hour)
	if body.DueDate != "" {
		if t, err := time.Parse("2006-01-02", body.DueDate); err == nil {
			dueDate = t
		}
	}

	invID := uuid.New().String()

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	var grnRef *string
	if body.GRNID != "" {
		grnRef = &body.GRNID
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO vendor_invoices (id, branch_id, purchase_order_id, grn_id, invoice_number, tax_invoice_number, invoice_date, due_date, subtotal, tax_amount, total_amount, paid_amount, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 0, 'unpaid', $12, NOW(), NOW())`,
		invID, branchID, body.PurchaseOrderID, grnRef, body.InvoiceNumber, body.TaxInvoiceNumber, invDate, dueDate, body.Subtotal, body.TaxAmount, body.TotalAmount, body.Notes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 3-Way Matching Journal Recognition:
	// [Debit] Hutang Belum Difakturkan (21201) (Menutup akrual GRN)
	// [Debit] PPN Masukan (11401) (Klaim Pajak Masukan)
	// [Credit] Hutang Usaha / AP (21101) (Pengakuan Hutang Resmi ke Vendor)
	var accUnbilledAPID, accPPNMasukanID, accHutangUsahaID string
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE code = '21201' LIMIT 1").Scan(&accUnbilledAPID)
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE code = '11401' LIMIT 1").Scan(&accPPNMasukanID)
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE code = '21101' LIMIT 1").Scan(&accHutangUsahaID)

	if accUnbilledAPID != "" && accHutangUsahaID != "" {
		journalID := uuid.New().String()
		journalRef := fmt.Sprintf("JRN-INV-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)
		journalDesc := fmt.Sprintf("Pengakuan Faktur Vendor #%s (PO #%s, Faktur Pajak: %s)", body.InvoiceNumber, poNumber, body.TaxInvoiceNumber)

		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entries (id, branch_id, reference_number, entry_date, description, status, created_at, updated_at)
			VALUES ($1, $2, $3, CURRENT_DATE, $4, 'posted', NOW(), NOW())`,
			journalID, branchID, journalRef, journalDesc)

		// Debit Unbilled AP
		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
			VALUES ($1, $2, $3, 'Penutupan Hutang Belum Difakturkan', $4, 0.00, NOW())`,
			uuid.New().String(), journalID, accUnbilledAPID, body.Subtotal)

		// Debit PPN Masukan jika ada
		if body.TaxAmount > 0 && accPPNMasukanID != "" {
			_, _ = tx.Exec(ctx, `
				INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
				VALUES ($1, $2, $3, 'Klaim PPN Masukan 11%', $4, 0.00, NOW())`,
				uuid.New().String(), journalID, accPPNMasukanID, body.TaxAmount)
		}

		// Credit Hutang Usaha (AP Total)
		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
			VALUES ($1, $2, $3, 'Pengakuan Hutang Dagang Supplier', 0.00, $4, NOW())`,
			uuid.New().String(), journalID, accHutangUsahaID, body.TotalAmount)
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":        true,
		"message":        fmt.Sprintf("Faktur Vendor %s berhasil diverifikasi (3-Way Match) & Hutang Usaha resmi tercatat", body.InvoiceNumber),
		"invoice_id":     invID,
		"invoice_number": body.InvoiceNumber,
		"total_amount":   body.TotalAmount,
	})
}

// -----------------------------------------------------------------------------
// 4. VENDOR PAYMENTS (PELUNASAN AP & PENUTUPAN PO)
// -----------------------------------------------------------------------------

func (h *P2PHandler) GetVendorPayments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		SELECT 
			vp.id, vp.payment_number, vp.payment_date, vp.payment_method, vp.amount_paid,
			COALESCE(vp.reference_number, '-') as ref_num, COALESCE(vp.notes, ''), vp.created_at,
			vi.invoice_number, po.po_number, s.name as supplier_name
		FROM vendor_payments vp
		JOIN vendor_invoices vi ON vp.vendor_invoice_id = vi.id
		JOIN purchase_orders po ON vi.purchase_order_id = po.id
		JOIN suppliers s ON po.supplier_id = s.id
		ORDER BY vp.payment_date DESC, vp.created_at DESC`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, payNum, method, refNum, notes, invNum, poNum, supplierName string
		var payDate, createdAt time.Time
		var amtPaid float64

		if err := rows.Scan(&id, &payNum, &payDate, &method, &amtPaid, &refNum, &notes, &createdAt, &invNum, &poNum, &supplierName); err == nil {
			list = append(list, map[string]interface{}{
				"id":               id,
				"payment_number":   payNum,
				"payment_date":     payDate.Format("2006-01-02"),
				"payment_method":   method,
				"amount_paid":      amtPaid,
				"reference_number": refNum,
				"notes":            notes,
				"invoice_number":   invNum,
				"po_number":        poNum,
				"supplier_name":    supplierName,
				"created_at":       createdAt.Format("2006-01-02 15:04"),
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

func (h *P2PHandler) CreateVendorPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		VendorInvoiceID string  `json:"vendor_invoice_id"`
		PaymentMethod   string  `json:"payment_method"` // 'bank_transfer', 'cash'
		AmountPaid      float64 `json:"amount_paid"`
		ReferenceNumber string  `json:"reference_number"`
		Notes           string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if body.VendorInvoiceID == "" || body.AmountPaid <= 0 {
		writeError(w, http.StatusBadRequest, "Pilih faktur vendor dan isi nominal pembayaran yang valid")
		return
	}

	var invoiceNum, poID, branchID string
	var totalAmt, alreadyPaid float64
	err := h.db.QueryRow(ctx, `
		SELECT vi.invoice_number, vi.purchase_order_id, vi.branch_id, vi.total_amount, vi.paid_amount
		FROM vendor_invoices vi 
		WHERE vi.id = $1`, body.VendorInvoiceID).Scan(&invoiceNum, &poID, &branchID, &totalAmt, &alreadyPaid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Faktur Vendor tidak ditemukan")
		return
	}

	payID := uuid.New().String()
	payNumber := fmt.Sprintf("PAY-AP-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	method := body.PaymentMethod
	if method == "" {
		method = "bank_transfer"
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	// Determine cash/bank account
	var bankAccountCode string = "11102"
	if method == "cash" {
		bankAccountCode = "11101"
	}
	var paymentAccID string
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE code = $1 LIMIT 1", bankAccountCode).Scan(&paymentAccID)

	var accRef *string
	if paymentAccID != "" {
		accRef = &paymentAccID
	}

	// 1. Insert Payment
	_, err = tx.Exec(ctx, `
		INSERT INTO vendor_payments (id, branch_id, vendor_invoice_id, payment_number, payment_date, payment_method, payment_account_id, amount_paid, reference_number, notes, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_DATE, $5, $6, $7, $8, $9, NOW())`,
		payID, branchID, body.VendorInvoiceID, payNumber, method, accRef, body.AmountPaid, body.ReferenceNumber, body.Notes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 2. Update Invoice Paid Amount and Status
	newPaid := alreadyPaid + body.AmountPaid
	newStatus := "partially_paid"
	if newPaid >= totalAmt {
		newStatus = "paid"
		// If invoice is fully paid, close the purchase order!
		_, _ = tx.Exec(ctx, "UPDATE purchase_orders SET status = 'completed', updated_at = NOW() WHERE id = $1", poID)
	}

	_, err = tx.Exec(ctx, "UPDATE vendor_invoices SET paid_amount = $1, status = $2, updated_at = NOW() WHERE id = $3", newPaid, newStatus, body.VendorInvoiceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. AP Settlement Journal:
	// [Debit] Hutang Usaha / Accounts Payable (21101)
	// [Credit] Bank BCA Operasional (11102) atau Kas (11101)
	var accHutangUsahaID string
	_ = tx.QueryRow(ctx, "SELECT id FROM chart_of_accounts WHERE code = '21101' LIMIT 1").Scan(&accHutangUsahaID)

	if accHutangUsahaID != "" && paymentAccID != "" {
		journalID := uuid.New().String()
		journalRef := fmt.Sprintf("JRN-PAY-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)
		journalDesc := fmt.Sprintf("Pelunasan Hutang Vendor (No Ref %s, Inv #%s)", payNumber, invoiceNum)

		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entries (id, branch_id, reference_number, entry_date, description, status, created_at, updated_at)
			VALUES ($1, $2, $3, CURRENT_DATE, $4, 'posted', NOW(), NOW())`,
			journalID, branchID, journalRef, journalDesc)

		// Debit Hutang Usaha
		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
			VALUES ($1, $2, $3, 'Pelunasan Hutang Dagang AP', $4, 0.00, NOW())`,
			uuid.New().String(), journalID, accHutangUsahaID, body.AmountPaid)

		// Credit Kas / Bank
		_, _ = tx.Exec(ctx, `
			INSERT INTO journal_entry_lines (id, journal_entry_id, account_id, description, debit, credit, created_at)
			VALUES ($1, $2, $3, 'Pengeluaran Kas/Bank', 0.00, $4, NOW())`,
			uuid.New().String(), journalID, paymentAccID, body.AmountPaid)
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success":        true,
		"message":        fmt.Sprintf("Pembayaran %s sebesar Rp %.0f berhasil diproses. Status Invoice: %s", payNumber, body.AmountPaid, newStatus),
		"payment_number": payNumber,
		"invoice_status": newStatus,
	})
}

// -----------------------------------------------------------------------------
// 5. COMPANY PROFILE ENDPOINT (SETTINGS)
// -----------------------------------------------------------------------------

func (h *P2PHandler) GetCompanyProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var name, code, address, phone, email string
	err := h.db.QueryRow(ctx, "SELECT name, code, address, phone, email FROM branches ORDER BY created_at ASC LIMIT 1").Scan(&name, &code, &address, &phone, &email)
	if err != nil {
		name = "Cafe ERP Indonesia"
		code = "HQ-01"
		address = "Jl. Senopati Raya No. 45, Jakarta Selatan"
		phone = "021-7201234"
		email = "contact@cafe-erp.com"
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"company_name": name,
			"code":         code,
			"npwp":         "01.234.567.8-901.000",
			"email":        email,
			"phone":        phone,
			"address":      address,
			"currency":     "IDR",
			"tax_rate":     11.00,
		},
	})
}

func (h *P2PHandler) UpdateCompanyProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		CompanyName string `json:"company_name"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	_, _ = h.db.Exec(ctx, `
		UPDATE branches 
		SET name = COALESCE(NULLIF($1, ''), name), 
		    address = COALESCE(NULLIF($2, ''), address), 
		    phone = COALESCE(NULLIF($3, ''), phone), 
		    email = COALESCE(NULLIF($4, ''), email),
		    updated_at = NOW()
		WHERE id IN (SELECT id FROM branches ORDER BY created_at ASC LIMIT 1)`,
		body.CompanyName, body.Address, body.Phone, body.Email)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profil perusahaan berhasil diperbarui",
	})
}
