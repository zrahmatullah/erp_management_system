package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/pkg/logger"
)

type MasterHandler struct {
	db *pgxpool.Pool
}

func NewMasterHandler(db *pgxpool.Pool) *MasterHandler {
	return &MasterHandler{db: db}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	if status >= 500 {
		logger.Log.Error().Str("raw_error", msg).Int("status", status).Msg("Internal server error")
		writeJSON(w, status, map[string]string{"error": "Terjadi kesalahan internal pada server"})
		return
	}
	writeJSON(w, status, map[string]string{"error": msg})
}

// -----------------------------------------------------------------------------
// 1. CABANG / BRANCHES
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListBranches(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT id, code, name, COALESCE(address,''), COALESCE(phone,''), COALESCE(email,''), is_active, created_at 
		FROM branches WHERE deleted_at IS NULL ORDER BY created_at ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, code, name, address, phone, email string
		var isActive bool
		var createdAt interface{}
		if err := rows.Scan(&id, &code, &name, &address, &phone, &email, &isActive, &createdAt); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "code": code, "name": name, "address": address,
				"phone": phone, "email": email, "is_active": isActive, "created_at": createdAt,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code     string `json:"code"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		IsActive bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	newID := uuid.New()
	query := `INSERT INTO branches (id, code, name, address, phone, email, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id`
	err := h.db.QueryRow(r.Context(), query, newID, body.Code, body.Name, body.Address, body.Phone, body.Email, body.IsActive).Scan(&newID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Branch created successfully"})
}

func (h *MasterHandler) UpdateBranch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Code     string `json:"code"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		IsActive bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	query := `UPDATE branches SET code=$1, name=$2, address=$3, phone=$4, email=$5, is_active=$6, updated_at=NOW()
			  WHERE id=$7 AND deleted_at IS NULL`
	_, err := h.db.Exec(r.Context(), query, body.Code, body.Name, body.Address, body.Phone, body.Email, body.IsActive, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Branch updated successfully"})
}

func (h *MasterHandler) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE branches SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Branch deleted successfully"})
}

// -----------------------------------------------------------------------------
// 2. PENGGUNA / USERS
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT u.id, u.username, u.email, u.full_name, COALESCE(u.phone, ''), u.is_active, 
		       COALESCE(r.id::text, ''), COALESCE(r.name, 'No Role'),
		       COALESCE(b.id::text, ''), COALESCE(b.name, 'All Branches')
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN branches b ON u.branch_id = b.id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at ASC`
	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, username, email, fullName, phone, roleID, roleName, branchID, branchName string
		var isActive bool
		if err := rows.Scan(&id, &username, &email, &fullName, &phone, &isActive, &roleID, &roleName, &branchID, &branchName); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "username": username, "email": email, "full_name": fullName,
				"phone": phone, "is_active": isActive, "role_id": roleID, "role_name": roleName,
				"branch_id": branchID, "branch_name": branchName,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		FullName string  `json:"full_name"`
		Phone    string  `json:"phone"`
		RoleID   *string `json:"role_id"`
		BranchID *string `json:"branch_id"`
		PINCode  string  `json:"pin_code"`
		IsActive bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	defaultPasswordHash := "$2a$12$Z0o2.J.YJ1oN2w3K6T0buexM.iP9Lp0nLp5R6m3p3mK5C4R7D0QjO" // Admin@123
	newID := uuid.New()
	query := `INSERT INTO users (id, username, email, password_hash, full_name, phone, role_id, branch_id, pin_code, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()) RETURNING id`
	err := h.db.QueryRow(r.Context(), query, newID, body.Username, body.Email, defaultPasswordHash, body.FullName, body.Phone, body.RoleID, body.BranchID, body.PINCode, body.IsActive).Scan(&newID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "User created successfully"})
}

func (h *MasterHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		FullName string  `json:"full_name"`
		Phone    string  `json:"phone"`
		RoleID   *string `json:"role_id"`
		BranchID *string `json:"branch_id"`
		PINCode  string  `json:"pin_code"`
		IsActive bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	query := `UPDATE users SET username=$1, email=$2, full_name=$3, phone=$4, role_id=$5, branch_id=$6, pin_code=$7, is_active=$8, updated_at=NOW()
			  WHERE id=$9 AND deleted_at IS NULL`
	_, err := h.db.Exec(r.Context(), query, body.Username, body.Email, body.FullName, body.Phone, body.RoleID, body.BranchID, body.PINCode, body.IsActive, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (h *MasterHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE users SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// -----------------------------------------------------------------------------
// 3. PERAN & MATRIKS HAK AKSES (ROLES & PERMISSIONS)
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT r.id, r.name, COALESCE(r.description, ''), COUNT(u.id) as user_count
		FROM roles r
		LEFT JOIN users u ON u.role_id = r.id AND u.deleted_at IS NULL
		WHERE r.deleted_at IS NULL
		GROUP BY r.id, r.name, r.description
		ORDER BY r.created_at ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, name, desc string
		var count int
		if err := rows.Scan(&id, &name, &desc, &count); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "name": name, "description": desc, "user_count": count,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")

	// Get all permissions
	permRows, err := h.db.Query(r.Context(), `SELECT id, module, action, COALESCE(description,'') FROM permissions ORDER BY module, action`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer permRows.Close()

	type PermItem struct {
		ID          string `json:"id"`
		Module      string `json:"module"`
		Action      string `json:"action"`
		Description string `json:"description"`
		Assigned    bool   `json:"assigned"`
	}
	var allPerms []PermItem
	permMap := make(map[string]int)

	for permRows.Next() {
		var p PermItem
		if err := permRows.Scan(&p.ID, &p.Module, &p.Action, &p.Description); err == nil {
			permMap[p.ID] = len(allPerms)
			allPerms = append(allPerms, p)
		}
	}

	// Get assigned permissions for this role
	assignedRows, err := h.db.Query(r.Context(), `SELECT permission_id FROM role_permissions WHERE role_id = $1`, roleID)
	if err == nil {
		defer assignedRows.Close()
		for assignedRows.Next() {
			var pid string
			if err := assignedRows.Scan(&pid); err == nil {
				if idx, ok := permMap[pid]; ok {
					allPerms[idx].Assigned = true
				}
			}
		}
	}

	writeJSON(w, http.StatusOK, allPerms)
}

func (h *MasterHandler) UpdateRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	var body struct {
		PermissionIDs []string `json:"permission_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	_, err = tx.Exec(r.Context(), `DELETE FROM role_permissions WHERE role_id = $1`, roleID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, pid := range body.PermissionIDs {
		_, err = tx.Exec(r.Context(), `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, roleID, pid)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	_ = tx.Commit(r.Context())
	writeJSON(w, http.StatusOK, map[string]string{"message": "Role permissions updated successfully"})
}

// -----------------------------------------------------------------------------
// 4. KATEGORI MENU (MENU CATEGORIES)
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT id, name, slug, COALESCE(icon,''), sort_order, is_active 
		FROM menu_categories WHERE deleted_at IS NULL ORDER BY sort_order ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, name, slug, icon string
		var sort int
		var isActive bool
		if err := rows.Scan(&id, &name, &slug, &icon, &sort, &isActive); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "name": name, "slug": slug, "icon": icon, "sort_order": sort, "is_active": isActive,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO menu_categories (id, name, slug, icon, sort_order, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		newID, body.Name, body.Slug, body.Icon, body.SortOrder, body.IsActive)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Category created"})
}

func (h *MasterHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE menu_categories SET name=$1, slug=$2, icon=$3, sort_order=$4, is_active=$5, updated_at=NOW()
		WHERE id=$6 AND deleted_at IS NULL`,
		body.Name, body.Slug, body.Icon, body.SortOrder, body.IsActive, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Category updated"})
}

func (h *MasterHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE menu_categories SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Category deleted"})
}

// -----------------------------------------------------------------------------
// 5. PRODUK & MENU (PRODUCTS & VARIANTS)
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT p.id, p.category_id, c.name as category_name, p.name, p.sku, 
		       COALESCE(p.description, ''), p.base_price, p.target_station, 
		       COALESCE(p.image_url, ''), p.is_active, p.created_at
		FROM products p
		LEFT JOIN menu_categories c ON p.category_id = c.id
		WHERE p.deleted_at IS NULL
		ORDER BY p.name ASC`
	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, catID, catName, name, sku, desc, station, img string
		var price float64
		var isActive bool
		var createdAt interface{}
		if err := rows.Scan(&id, &catID, &catName, &name, &sku, &desc, &price, &station, &img, &isActive, &createdAt); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "category_id": catID, "category_name": catName,
				"name": name, "sku": sku, "description": desc, "base_price": price,
				"target_station": station, "image_url": img, "is_active": isActive,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CategoryID    string  `json:"category_id"`
		Name          string  `json:"name"`
		SKU           string  `json:"sku"`
		Description   string  `json:"description"`
		BasePrice     float64 `json:"base_price"`
		TargetStation string  `json:"target_station"`
		ImageURL      string  `json:"image_url"`
		IsActive      bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO products (id, category_id, name, sku, description, base_price, target_station, image_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`,
		newID, body.CategoryID, body.Name, body.SKU, body.Description, body.BasePrice, body.TargetStation, body.ImageURL, body.IsActive)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Product created"})
}

func (h *MasterHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		CategoryID    string  `json:"category_id"`
		Name          string  `json:"name"`
		SKU           string  `json:"sku"`
		Description   string  `json:"description"`
		BasePrice     float64 `json:"base_price"`
		TargetStation string  `json:"target_station"`
		ImageURL      string  `json:"image_url"`
		IsActive      bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE products SET category_id=$1, name=$2, sku=$3, description=$4, base_price=$5, target_station=$6, image_url=$7, is_active=$8, updated_at=NOW()
		WHERE id=$9 AND deleted_at IS NULL`,
		body.CategoryID, body.Name, body.SKU, body.Description, body.BasePrice, body.TargetStation, body.ImageURL, body.IsActive, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Product updated"})
}

func (h *MasterHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE products SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Product deleted"})
}

// -----------------------------------------------------------------------------
// 6. BAHAN BAKU / INVENTORI (INVENTORY ITEMS)
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListInventoryItems(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT i.id, i.sku, i.name, i.category, i.uom, i.min_stock, i.max_stock, i.average_cost, i.is_active,
		       COALESCE(SUM(s.quantity), 0) as current_stock
		FROM inventory_items i
		LEFT JOIN inventory_stocks s ON i.id = s.inventory_item_id
		WHERE i.deleted_at IS NULL
		GROUP BY i.id, i.sku, i.name, i.category, i.uom, i.min_stock, i.max_stock, i.average_cost, i.is_active
		ORDER BY i.name ASC`
	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, sku, name, cat, uom string
		var minStock, maxStock, avgCost, currentStock float64
		var isActive bool
		if err := rows.Scan(&id, &sku, &name, &cat, &uom, &minStock, &maxStock, &avgCost, &isActive, &currentStock); err == nil {
			status := "Normal"
			if currentStock <= minStock {
				status = "Low Stock"
			} else if maxStock > 0 && currentStock > maxStock {
				status = "Overstocked"
			}

			list = append(list, map[string]interface{}{
				"id": id, "sku": sku, "name": name, "category": cat, "uom": uom,
				"min_stock": minStock, "max_stock": maxStock, "average_cost": avgCost,
				"current_stock": currentStock, "status": status, "is_active": isActive,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SKU         string  `json:"sku"`
		Name        string  `json:"name"`
		Category    string  `json:"category"`
		UOM         string  `json:"uom"`
		MinStock    float64 `json:"min_stock"`
		MaxStock    float64 `json:"max_stock"`
		AverageCost float64 `json:"average_cost"`
		IsActive    bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO inventory_items (id, sku, name, category, uom, min_stock, max_stock, average_cost, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`,
		newID, body.SKU, body.Name, body.Category, body.UOM, body.MinStock, body.MaxStock, body.AverageCost, body.IsActive)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Item created"})
}

func (h *MasterHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		SKU         string  `json:"sku"`
		Name        string  `json:"name"`
		Category    string  `json:"category"`
		UOM         string  `json:"uom"`
		MinStock    float64 `json:"min_stock"`
		MaxStock    float64 `json:"max_stock"`
		AverageCost float64 `json:"average_cost"`
		IsActive    bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE inventory_items SET sku=$1, name=$2, category=$3, uom=$4, min_stock=$5, max_stock=$6, average_cost=$7, is_active=$8, updated_at=NOW()
		WHERE id=$9 AND deleted_at IS NULL`,
		body.SKU, body.Name, body.Category, body.UOM, body.MinStock, body.MaxStock, body.AverageCost, body.IsActive, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Item updated"})
}

func (h *MasterHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE inventory_items SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Item deleted"})
}

// -----------------------------------------------------------------------------
// 7. MEJA & ZONASI (TABLES & ZONES)
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT t.id, t.branch_id, b.name as branch_name, t.zone_id, COALESCE(z.name, 'No Zone') as zone_name,
		       t.table_number, t.capacity, t.status, t.pos_x, t.pos_y
		FROM cafe_tables t
		LEFT JOIN branches b ON t.branch_id = b.id
		LEFT JOIN table_zones z ON t.zone_id = z.id
		WHERE t.deleted_at IS NULL
		ORDER BY t.table_number ASC`
	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, bID, bName, zID, zName, tableNum, status string
		var cap, x, y int
		if err := rows.Scan(&id, &bID, &bName, &zID, &zName, &tableNum, &cap, &status, &x, &y); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "branch_id": bID, "branch_name": bName, "zone_id": zID, "zone_name": zName,
				"table_number": tableNum, "capacity": cap, "status": status, "pos_x": x, "pos_y": y,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BranchID    string `json:"branch_id"`
		ZoneID      string `json:"zone_id"`
		TableNumber string `json:"table_number"`
		Capacity    int    `json:"capacity"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if body.Status == "" {
		body.Status = "available"
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO cafe_tables (id, branch_id, zone_id, table_number, capacity, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		newID, body.BranchID, body.ZoneID, body.TableNumber, body.Capacity, body.Status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Table created"})
}

func (h *MasterHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		BranchID    string `json:"branch_id"`
		ZoneID      string `json:"zone_id"`
		TableNumber string `json:"table_number"`
		Capacity    int    `json:"capacity"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE cafe_tables SET branch_id=$1, zone_id=$2, table_number=$3, capacity=$4, status=$5, updated_at=NOW()
		WHERE id=$6 AND deleted_at IS NULL`,
		body.BranchID, body.ZoneID, body.TableNumber, body.Capacity, body.Status, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Table updated"})
}

func (h *MasterHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE cafe_tables SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Table deleted"})
}

// -----------------------------------------------------------------------------
// 8. BAGAN AKUN / CHART OF ACCOUNTS
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT id, code, name, account_type, COALESCE(description,''), is_active 
		FROM chart_of_accounts WHERE deleted_at IS NULL ORDER BY code ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, code, name, accType, desc string
		var isActive bool
		if err := rows.Scan(&id, &code, &name, &accType, &desc, &isActive); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "code": code, "name": name, "account_type": accType, "description": desc, "is_active": isActive,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		AccountType string `json:"account_type"`
		Description string `json:"description"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO chart_of_accounts (id, code, name, account_type, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		newID, body.Code, body.Name, body.AccountType, body.Description, body.IsActive)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Account created"})
}

func (h *MasterHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		AccountType string `json:"account_type"`
		Description string `json:"description"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE chart_of_accounts SET code=$1, name=$2, account_type=$3, description=$4, is_active=$5, updated_at=NOW()
		WHERE id=$6 AND deleted_at IS NULL`,
		body.Code, body.Name, body.AccountType, body.Description, body.IsActive, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Account updated"})
}

func (h *MasterHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE chart_of_accounts SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Account deleted"})
}

// -----------------------------------------------------------------------------
// 9. SUPPLIERS, WAREHOUSES, SHIFTS, DEPARTMENTS, POSITIONS
// -----------------------------------------------------------------------------
func (h *MasterHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT id, code, name, COALESCE(contact_person,''), COALESCE(phone,''), COALESCE(email,''), COALESCE(address,''), is_active FROM suppliers WHERE deleted_at IS NULL ORDER BY name ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id, code, name, contact, phone, email, addr string
		var isActive bool
		if err := rows.Scan(&id, &code, &name, &contact, &phone, &email, &addr, &isActive); err == nil {
			list = append(list, map[string]interface{}{"id": id, "code": code, "name": name, "contact_person": contact, "phone": phone, "email": email, "address": addr, "is_active": isActive})
		}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) ListWarehouses(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT w.id, w.branch_id, b.name as branch_name, w.name, w.type, COALESCE(w.address,'') FROM warehouses w LEFT JOIN branches b ON w.branch_id = b.id WHERE w.deleted_at IS NULL ORDER BY w.name ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id, bID, bName, name, wType, addr string
		if err := rows.Scan(&id, &bID, &bName, &name, &wType, &addr); err == nil {
			list = append(list, map[string]interface{}{"id": id, "branch_id": bID, "branch_name": bName, "name": name, "type": wType, "address": addr})
		}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) ListShifts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT s.id, s.branch_id, s.name, s.start_time::text, s.end_time::text, s.color FROM work_shifts s WHERE s.deleted_at IS NULL ORDER BY s.start_time ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id, bID, name, start, end, color string
		if err := rows.Scan(&id, &bID, &name, &start, &end, &color); err == nil {
			list = append(list, map[string]interface{}{"id": id, "branch_id": bID, "name": name, "start_time": start, "end_time": end, "color": color})
		}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT id, name, COALESCE(description,'') FROM departments WHERE deleted_at IS NULL ORDER BY name ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id, name, desc string
		if err := rows.Scan(&id, &name, &desc); err == nil {
			list = append(list, map[string]interface{}{"id": id, "name": name, "description": desc})
		}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) ListPositions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT p.id, p.department_id, d.name, p.title, p.level, p.base_salary FROM positions p LEFT JOIN departments d ON p.department_id = d.id WHERE p.deleted_at IS NULL ORDER BY p.title ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id, deptID, deptName, title string
		var level int
		var salary float64
		if err := rows.Scan(&id, &deptID, &deptName, &title, &level, &salary); err == nil {
			list = append(list, map[string]interface{}{"id": id, "department_id": deptID, "department_name": deptName, "title": title, "level": level, "base_salary": salary})
		}
	}
	writeJSON(w, http.StatusOK, list)
}
