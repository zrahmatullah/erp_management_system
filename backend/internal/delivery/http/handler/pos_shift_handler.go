package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/delivery/http/middleware"
	"cafe-erp-system/backend/pkg/response"
)

type POSShiftHandler struct {
	db *pgxpool.Pool
}

func NewPOSShiftHandler(db *pgxpool.Pool) *POSShiftHandler {
	return &POSShiftHandler{db: db}
}

type CashMovementRequest struct {
	MovementType string  `json:"movement_type"` // "cash_drop", "paid_out", "cash_in"
	Amount       float64 `json:"amount"`
	Reason       string  `json:"reason"`
	AuthorizedBy *string `json:"authorized_by,omitempty"`
}

type OpenShiftRequest struct {
	BranchID         *string `json:"branch_id,omitempty"`
	ShiftName        string  `json:"shift_name"`
	OpeningCashFloat float64 `json:"opening_cash_float"`
	Notes            string  `json:"notes,omitempty"`
}

type CloseShiftRequest struct {
	ActualCashCounted float64 `json:"actual_cash_counted"`
	SupervisorID      *string `json:"supervisor_id,omitempty"`
	Notes             string  `json:"notes,omitempty"`
}

// GetCurrentShift returns the active open shift for the current cashier or branch
func (h *POSShiftHandler) GetCurrentShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	branchID := r.URL.Query().Get("branch_id")
	cashierID := ""

	if claims, err := middleware.GetUserFromContext(ctx); err == nil && claims != nil {
		cashierID = claims.UserID.String()
		if branchID == "" && claims.BranchID != nil {
			branchID = claims.BranchID.String()
		}
	}

	query := `
		SELECT s.id, s.branch_id, b.name as branch_name, s.cashier_id, u.name as cashier_name,
		       s.shift_name, s.opening_cash_float, s.actual_cash_counted, s.expected_cash_total,
		       s.cash_difference, s.total_cash_sales, s.total_non_cash_sales, s.total_cash_drops,
		       s.total_refunds, s.status, s.opened_at, s.notes
		FROM cashier_shifts s
		JOIN branches b ON s.branch_id = b.id
		JOIN users u ON s.cashier_id = u.id
		WHERE s.status = 'open'
	`
	var args []interface{}
	argCount := 1

	if cashierID != "" {
		query += fmt.Sprintf(" AND (s.cashier_id = $%d", argCount)
		args = append(args, cashierID)
		argCount++
		if branchID != "" {
			query += fmt.Sprintf(" OR s.branch_id = $%d)", argCount)
			args = append(args, branchID)
			argCount++
		} else {
			query += ")"
		}
	} else if branchID != "" {
		query += fmt.Sprintf(" AND s.branch_id = $%d", argCount)
		args = append(args, branchID)
		argCount++
	}

	query += " ORDER BY s.opened_at DESC LIMIT 1"

	var shift struct {
		ID                  string                   `json:"id"`
		BranchID            string                   `json:"branch_id"`
		BranchName          string                   `json:"branch_name"`
		CashierID           string                   `json:"cashier_id"`
		CashierName         string                   `json:"cashier_name"`
		ShiftName           string                   `json:"shift_name"`
		OpeningCashFloat    float64                  `json:"opening_cash_float"`
		ActualCashCounted   float64                  `json:"actual_cash_counted"`
		ExpectedCashTotal   float64                  `json:"expected_cash_total"`
		CashDifference      float64                  `json:"cash_difference"`
		TotalCashSales      float64                  `json:"total_cash_sales"`
		TotalNonCashSales   float64                  `json:"total_non_cash_sales"`
		TotalCashDrops      float64                  `json:"total_cash_drops"`
		TotalRefunds        float64                  `json:"total_refunds"`
		Status              string                   `json:"status"`
		OpenedAt            time.Time                `json:"opened_at"`
		Notes               *string                  `json:"notes,omitempty"`
		ExpectedCurrentCash float64                  `json:"expected_current_cash"`
		RecentMovements     []map[string]interface{} `json:"recent_movements"`
	}

	err := h.db.QueryRow(ctx, query, args...).Scan(
		&shift.ID, &shift.BranchID, &shift.BranchName, &shift.CashierID, &shift.CashierName,
		&shift.ShiftName, &shift.OpeningCashFloat, &shift.ActualCashCounted, &shift.ExpectedCashTotal,
		&shift.CashDifference, &shift.TotalCashSales, &shift.TotalNonCashSales, &shift.TotalCashDrops,
		&shift.TotalRefunds, &shift.Status, &shift.OpenedAt, &shift.Notes,
	)

	if err != nil {
		response.Success(w, "Tidak ada shift aktif", nil)
		return
	}

	shift.ExpectedCurrentCash = shift.OpeningCashFloat + shift.TotalCashSales - shift.TotalCashDrops - shift.TotalRefunds

	rows, err := h.db.Query(ctx, `
		SELECT m.id, m.movement_type, m.amount, m.reason, m.created_at, u.name as authorized_by_name
		FROM cashier_shift_movements m
		LEFT JOIN users u ON m.authorized_by = u.id
		WHERE m.shift_id = $1
		ORDER BY m.created_at DESC LIMIT 10`, shift.ID)
	if err == nil {
		defer rows.Close()
		shift.RecentMovements = make([]map[string]interface{}, 0)
		for rows.Next() {
			var movID, movType, reason string
			var amount float64
			var createdAt time.Time
			var authName *string
			if err := rows.Scan(&movID, &movType, &amount, &reason, &createdAt, &authName); err == nil {
				shift.RecentMovements = append(shift.RecentMovements, map[string]interface{}{
					"id":                 movID,
					"movement_type":      movType,
					"amount":             amount,
					"reason":             reason,
					"created_at":         createdAt,
					"authorized_by_name": authName,
				})
			}
		}
	}

	response.Success(w, "Shift kasir aktif ditemukan", shift)
}

// OpenShift opens a new cashier shift with an opening cash float
func (h *POSShiftHandler) OpenShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req OpenShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.OpeningCashFloat < 0 {
		response.Error(w, http.StatusBadRequest, "Modal awal kasir tidak boleh bernilai negatif")
		return
	}
	if req.ShiftName == "" {
		req.ShiftName = "Shift Pagi"
	}

	cashierID := ""
	if claims, err := middleware.GetUserFromContext(ctx); err == nil && claims != nil {
		cashierID = claims.UserID.String()
		if req.BranchID == nil && claims.BranchID != nil {
			bid := claims.BranchID.String()
			req.BranchID = &bid
		}
	}

	if req.BranchID == nil || *req.BranchID == "" {
		var firstBranch string
		if err := h.db.QueryRow(ctx, "SELECT id FROM branches ORDER BY is_active DESC LIMIT 1").Scan(&firstBranch); err == nil {
			req.BranchID = &firstBranch
		} else {
			defaultB := "b1111111-0000-0000-0000-000000000001"
			req.BranchID = &defaultB
		}
	}

	if cashierID == "" {
		_ = h.db.QueryRow(ctx, "SELECT id FROM users WHERE is_active = true ORDER BY created_at ASC LIMIT 1").Scan(&cashierID)
	}

	var existingOpenID string
	err := h.db.QueryRow(ctx, `
		SELECT id FROM cashier_shifts 
		WHERE cashier_id = $1 AND status = 'open' 
		LIMIT 1`, cashierID).Scan(&existingOpenID)
	if err == nil {
		response.Error(w, http.StatusBadRequest, "Kasir masih memiliki shift aktif yang belum ditutup. Harap tutup shift sebelumnya.")
		return
	}

	shiftID := uuid.New().String()
	_, err = h.db.Exec(ctx, `
		INSERT INTO cashier_shifts (
			id, branch_id, cashier_id, shift_name, opening_cash_float,
			status, notes, opened_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, 'open', $6, NOW(), NOW(), NOW())`,
		shiftID, *req.BranchID, cashierID, req.ShiftName, req.OpeningCashFloat, req.Notes)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal membuka shift kasir: "+err.Error())
		return
	}

	response.Success(w, "Shift kasir berhasil dibuka", map[string]interface{}{
		"id":                 shiftID,
		"branch_id":          *req.BranchID,
		"cashier_id":         cashierID,
		"shift_name":         req.ShiftName,
		"opening_cash_float": req.OpeningCashFloat,
		"status":             "open",
		"opened_at":          time.Now(),
	})
}

// RecordCashMovement records cash drop, paid out, or cash in
func (h *POSShiftHandler) RecordCashMovement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shiftID := chi.URLParam(r, "id")

	var req CashMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Amount <= 0 {
		response.Error(w, http.StatusBadRequest, "Nominal mutasi kas harus lebih besar dari 0")
		return
	}
	if req.MovementType == "" {
		req.MovementType = "cash_drop"
	}
	if req.Reason == "" {
		req.Reason = "Mutasi kasir operasional"
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	var status string
	err = tx.QueryRow(ctx, "SELECT status FROM cashier_shifts WHERE id = $1", shiftID).Scan(&status)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Shift kasir tidak ditemukan")
		return
	}
	if status != "open" {
		response.Error(w, http.StatusBadRequest, "Hanya shift berstatus 'open' yang dapat mencatat mutasi kas")
		return
	}

	movID := uuid.New().String()
	_, err = tx.Exec(ctx, `
		INSERT INTO cashier_shift_movements (id, shift_id, movement_type, amount, reason, authorized_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
		movID, shiftID, req.MovementType, req.Amount, req.Reason, req.AuthorizedBy)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal mencatat mutasi kas: "+err.Error())
		return
	}

	if req.MovementType == "cash_drop" {
		_, _ = tx.Exec(ctx, `
			UPDATE cashier_shifts 
			SET total_cash_drops = total_cash_drops + $1, updated_at = NOW() 
			WHERE id = $2`, req.Amount, shiftID)
	}

	if err := tx.Commit(ctx); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "Mutasi kas kasir berhasil dicatat", map[string]interface{}{
		"id":            movID,
		"shift_id":      shiftID,
		"movement_type": req.MovementType,
		"amount":        req.Amount,
		"reason":        req.Reason,
	})
}

// CloseShift performs blind closing reconciliation and calculates discrepancy
func (h *POSShiftHandler) CloseShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shiftID := chi.URLParam(r, "id")

	var req CloseShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.ActualCashCounted < 0 {
		response.Error(w, http.StatusBadRequest, "Uang fisik kasir tidak boleh bernilai negatif")
		return
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	var openingFloat, totalRefunds float64
	var status string
	err = tx.QueryRow(ctx, `
		SELECT status, opening_cash_float, total_refunds 
		FROM cashier_shifts WHERE id = $1 FOR UPDATE`, shiftID).
		Scan(&status, &openingFloat, &totalRefunds)

	if err != nil {
		response.Error(w, http.StatusNotFound, "Shift kasir tidak ditemukan")
		return
	}
	if status != "open" {
		response.Error(w, http.StatusBadRequest, "Shift ini sudah ditutup sebelumnya")
		return
	}

	var cashSales, nonCashSales float64
	_ = tx.QueryRow(ctx, `
		SELECT 
			COALESCE(SUM(CASE WHEN LOWER(payment_method) = 'cash' THEN amount_paid - change_due ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN LOWER(payment_method) != 'cash' THEN amount_paid ELSE 0 END), 0)
		FROM payments
		WHERE shift_id = $1 AND status = 'completed'`, shiftID).Scan(&cashSales, &nonCashSales)

	var totalDrops float64
	_ = tx.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0) 
		FROM cashier_shift_movements 
		WHERE shift_id = $1 AND movement_type = 'cash_drop'`, shiftID).Scan(&totalDrops)

	expectedCash := openingFloat + cashSales - totalDrops - totalRefunds
	cashDifference := req.ActualCashCounted - expectedCash

	_, err = tx.Exec(ctx, `
		UPDATE cashier_shifts SET
			actual_cash_counted = $1,
			expected_cash_total = $2,
			cash_difference = $3,
			total_cash_sales = $4,
			total_non_cash_sales = $5,
			total_cash_drops = $6,
			status = 'closed',
			closed_at = NOW(),
			supervisor_id = $7,
			notes = $8,
			updated_at = NOW()
		WHERE id = $9`,
		req.ActualCashCounted, expectedCash, cashDifference,
		cashSales, nonCashSales, totalDrops,
		req.SupervisorID, req.Notes, shiftID)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal menutup shift: "+err.Error())
		return
	}

	if err := tx.Commit(ctx); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "Shift kasir berhasil ditutup (Blind Closing Selesai)", map[string]interface{}{
		"shift_id":             shiftID,
		"status":               "closed",
		"opening_cash_float":   openingFloat,
		"total_cash_sales":     cashSales,
		"total_non_cash_sales": nonCashSales,
		"total_cash_drops":     totalDrops,
		"expected_cash_total":  expectedCash,
		"actual_cash_counted":  req.ActualCashCounted,
		"cash_difference":      cashDifference,
		"closed_at":            time.Now(),
	})
}

// GetShiftSummary returns full Z-Report details for receipt printing
func (h *POSShiftHandler) GetShiftSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shiftID := chi.URLParam(r, "id")

	query := `
		SELECT s.id, s.branch_id, b.name as branch_name, s.cashier_id, u.name as cashier_name,
		       s.supervisor_id, su.name as supervisor_name, s.shift_name,
		       s.opening_cash_float, s.actual_cash_counted, s.expected_cash_total, s.cash_difference,
		       s.total_cash_sales, s.total_non_cash_sales, s.total_cash_drops, s.total_refunds,
		       s.status, s.opened_at, s.closed_at, s.notes
		FROM cashier_shifts s
		JOIN branches b ON s.branch_id = b.id
		JOIN users u ON s.cashier_id = u.id
		LEFT JOIN users su ON s.supervisor_id = su.id
		WHERE s.id = $1`

	var summary struct {
		ID                string     `json:"id"`
		BranchID          string     `json:"branch_id"`
		BranchName        string     `json:"branch_name"`
		CashierID         string     `json:"cashier_id"`
		CashierName       string     `json:"cashier_name"`
		SupervisorID      *string    `json:"supervisor_id,omitempty"`
		SupervisorName    *string    `json:"supervisor_name,omitempty"`
		ShiftName         string     `json:"shift_name"`
		OpeningCashFloat  float64    `json:"opening_cash_float"`
		ActualCashCounted float64    `json:"actual_cash_counted"`
		ExpectedCashTotal float64    `json:"expected_cash_total"`
		CashDifference    float64    `json:"cash_difference"`
		TotalCashSales    float64    `json:"total_cash_sales"`
		TotalNonCashSales float64    `json:"total_non_cash_sales"`
		TotalCashDrops    float64    `json:"total_cash_drops"`
		TotalRefunds      float64    `json:"total_refunds"`
		Status            string     `json:"status"`
		OpenedAt          time.Time  `json:"opened_at"`
		ClosedAt          *time.Time `json:"closed_at,omitempty"`
		Notes             *string    `json:"notes,omitempty"`
	}

	err := h.db.QueryRow(ctx, query, shiftID).Scan(
		&summary.ID, &summary.BranchID, &summary.BranchName, &summary.CashierID, &summary.CashierName,
		&summary.SupervisorID, &summary.SupervisorName, &summary.ShiftName,
		&summary.OpeningCashFloat, &summary.ActualCashCounted, &summary.ExpectedCashTotal, &summary.CashDifference,
		&summary.TotalCashSales, &summary.TotalNonCashSales, &summary.TotalCashDrops, &summary.TotalRefunds,
		&summary.Status, &summary.OpenedAt, &summary.ClosedAt, &summary.Notes,
	)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Shift kasir tidak ditemukan")
		return
	}

	rows, err := h.db.Query(ctx, `
		SELECT payment_method, COUNT(*), COALESCE(SUM(amount_paid - change_due), 0)
		FROM payments
		WHERE shift_id = $1 AND status = 'completed'
		GROUP BY payment_method`, shiftID)
	paymentBreakdown := make([]map[string]interface{}, 0)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var method string
			var count int
			var total float64
			if err := rows.Scan(&method, &count, &total); err == nil {
				paymentBreakdown = append(paymentBreakdown, map[string]interface{}{
					"method": strings.ToUpper(method),
					"count":  count,
					"total":  total,
				})
			}
		}
	}

	movRows, err := h.db.Query(ctx, `
		SELECT m.id, m.movement_type, m.amount, m.reason, m.created_at, u.name as authorized_by_name
		FROM cashier_shift_movements m
		LEFT JOIN users u ON m.authorized_by = u.id
		WHERE m.shift_id = $1
		ORDER BY m.created_at ASC`, shiftID)
	movements := make([]map[string]interface{}, 0)
	if err == nil {
		defer movRows.Close()
		for movRows.Next() {
			var mID, mType, reason string
			var amount float64
			var cAt time.Time
			var authName *string
			if err := movRows.Scan(&mID, &mType, &amount, &reason, &cAt, &authName); err == nil {
				movements = append(movements, map[string]interface{}{
					"id":                 mID,
					"movement_type":      mType,
					"amount":             amount,
					"reason":             reason,
					"created_at":         cAt,
					"authorized_by_name": authName,
				})
			}
		}
	}

	response.Success(w, "Z-Report Ringkasan Shift Kasir", map[string]interface{}{
		"shift":             summary,
		"payment_breakdown": paymentBreakdown,
		"cash_movements":    movements,
		"gross_total_sales": summary.TotalCashSales + summary.TotalNonCashSales,
	})
}

// ListShifts returns paginated list of past cashier shifts
func (h *POSShiftHandler) ListShifts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	branchID := r.URL.Query().Get("branch_id")
	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	baseQuery := `
		FROM cashier_shifts s
		JOIN branches b ON s.branch_id = b.id
		JOIN users u ON s.cashier_id = u.id
		WHERE 1=1`
	var args []interface{}
	argCount := 1

	if branchID != "" {
		baseQuery += fmt.Sprintf(" AND s.branch_id = $%d", argCount)
		args = append(args, branchID)
		argCount++
	}
	if status != "" {
		baseQuery += fmt.Sprintf(" AND s.status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	var total int
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total)

	selectQuery := fmt.Sprintf(`
		SELECT s.id, s.branch_id, b.name as branch_name, s.cashier_id, u.name as cashier_name,
		       s.shift_name, s.opening_cash_float, s.actual_cash_counted, s.expected_cash_total,
		       s.cash_difference, s.total_cash_sales, s.total_non_cash_sales, s.total_cash_drops,
		       s.status, s.opened_at, s.closed_at
		%s
		ORDER BY s.opened_at DESC
		LIMIT $%d OFFSET $%d`, baseQuery, argCount, argCount+1)
	args = append(args, perPage, offset)

	rows, err := h.db.Query(ctx, selectQuery, args...)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal mengambil daftar shift: "+err.Error())
		return
	}
	defer rows.Close()

	type ShiftItem struct {
		ID                string     `json:"id"`
		BranchID          string     `json:"branch_id"`
		BranchName        string     `json:"branch_name"`
		CashierID         string     `json:"cashier_id"`
		CashierName       string     `json:"cashier_name"`
		ShiftName         string     `json:"shift_name"`
		OpeningCashFloat  float64    `json:"opening_cash_float"`
		ActualCashCounted float64    `json:"actual_cash_counted"`
		ExpectedCashTotal float64    `json:"expected_cash_total"`
		CashDifference    float64    `json:"cash_difference"`
		TotalCashSales    float64    `json:"total_cash_sales"`
		TotalNonCashSales float64    `json:"total_non_cash_sales"`
		TotalCashDrops    float64    `json:"total_cash_drops"`
		Status            string     `json:"status"`
		OpenedAt          time.Time  `json:"opened_at"`
		ClosedAt          *time.Time `json:"closed_at,omitempty"`
	}

	shifts := make([]ShiftItem, 0)
	for rows.Next() {
		var s ShiftItem
		if err := rows.Scan(
			&s.ID, &s.BranchID, &s.BranchName, &s.CashierID, &s.CashierName,
			&s.ShiftName, &s.OpeningCashFloat, &s.ActualCashCounted, &s.ExpectedCashTotal,
			&s.CashDifference, &s.TotalCashSales, &s.TotalNonCashSales, &s.TotalCashDrops,
			&s.Status, &s.OpenedAt, &s.ClosedAt,
		); err == nil {
			shifts = append(shifts, s)
		}
	}

	totalPages := (total + perPage - 1) / perPage
	response.Paginated(w, "Daftar shift kasir berhasil diambil", shifts, response.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}
