package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

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
		       COALESCE(b.id::text, ''), COALESCE(b.name, 'All Branches'),
		       COALESCE(e.id::text, ''), COALESCE(e.nik, ''), COALESCE(concat(e.first_name, ' ', e.last_name), ''),
		       COALESCE(p.title, '')
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN branches b ON u.branch_id = b.id
		LEFT JOIN employees e ON e.user_id = u.id AND e.deleted_at IS NULL
		LEFT JOIN positions p ON e.position_id = p.id
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
		var empID, empNik, empName, empPos string
		var isActive bool
		if err := rows.Scan(
			&id, &username, &email, &fullName, &phone, &isActive,
			&roleID, &roleName, &branchID, &branchName,
			&empID, &empNik, &empName, &empPos,
		); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "username": username, "email": email, "full_name": fullName,
				"phone": phone, "is_active": isActive, "role_id": roleID, "role_name": roleName,
				"branch_id": branchID, "branch_name": branchName,
				"employee_id": empID, "employee_nik": empNik,
				"employee_name": empName, "employee_position": empPos,
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
		Username   string  `json:"username"`
		Email      string  `json:"email"`
		Password   string  `json:"password"`
		FullName   string  `json:"full_name"`
		Phone      string  `json:"phone"`
		RoleID     *string `json:"role_id"`
		BranchID   *string `json:"branch_id"`
		EmployeeID *string `json:"employee_id"`
		PINCode    string  `json:"pin_code"`
		IsActive   bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if body.Email == "" {
		writeError(w, http.StatusBadRequest, "Email wajib diisi")
		return
	}

	if body.Username == "" {
		parts := strings.Split(body.Email, "@")
		body.Username = parts[0]
	}

	passwordHash := "$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK" // Default Admin@123
	if strings.TrimSpace(body.Password) != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err == nil {
			passwordHash = string(hashed)
		}
	}

	newID := uuid.New()
	query := `INSERT INTO users (id, username, email, password_hash, full_name, phone, role_id, branch_id, pin_code, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()) RETURNING id`
	err := h.db.QueryRow(r.Context(), query, newID, body.Username, body.Email, passwordHash, body.FullName, body.Phone, body.RoleID, body.BranchID, body.PINCode, body.IsActive).Scan(&newID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Link user to employee if employee_id provided
	if body.EmployeeID != nil && *body.EmployeeID != "" {
		_, _ = h.db.Exec(r.Context(), `UPDATE employees SET user_id = $1 WHERE id = $2`, newID, *body.EmployeeID)
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "User berhasil dibuat dan ditautkan ke sistem"})
}

func (h *MasterHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Username   string  `json:"username"`
		Email      string  `json:"email"`
		Password   string  `json:"password"`
		FullName   string  `json:"full_name"`
		Phone      string  `json:"phone"`
		RoleID     *string `json:"role_id"`
		BranchID   *string `json:"branch_id"`
		EmployeeID *string `json:"employee_id"`
		PINCode    string  `json:"pin_code"`
		IsActive   bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if strings.TrimSpace(body.Password) != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err == nil {
			_, _ = h.db.Exec(r.Context(), `UPDATE users SET password_hash=$1 WHERE id=$2`, string(hashed), id)
		}
	}

	query := `UPDATE users SET username=$1, email=$2, full_name=$3, phone=$4, role_id=$5, branch_id=$6, pin_code=$7, is_active=$8, updated_at=NOW()
			  WHERE id=$9 AND deleted_at IS NULL`
	_, err := h.db.Exec(r.Context(), query, body.Username, body.Email, body.FullName, body.Phone, body.RoleID, body.BranchID, body.PINCode, body.IsActive, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})

	// Handle employee linkage
	if body.EmployeeID != nil {
		// Clear previous linkage for this user
		_, _ = h.db.Exec(r.Context(), `UPDATE employees SET user_id=NULL WHERE user_id=$1`, id)
		if *body.EmployeeID != "" {
			_, _ = h.db.Exec(r.Context(), `UPDATE employees SET user_id=$1 WHERE id=$2`, id, *body.EmployeeID)
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User berhasil diperbarui"})
}

func (h *MasterHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// Unlink employee first
	_, _ = h.db.Exec(r.Context(), `UPDATE employees SET user_id=NULL WHERE user_id=$1`, id)
	_, err := h.db.Exec(r.Context(), `UPDATE users SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
	writeJSON(w, http.StatusOK, map[string]string{"message": "User berhasil dinonaktifkan/dihapus"})
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

func (h *MasterHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if strings.TrimSpace(body.Name) == "" {
		writeError(w, http.StatusBadRequest, "Nama peran (role) wajib diisi")
		return
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())`,
		newID, strings.TrimSpace(body.Name), body.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Role created successfully"})
}

func (h *MasterHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if strings.TrimSpace(body.Name) == "" {
		writeError(w, http.StatusBadRequest, "Nama peran (role) wajib diisi")
		return
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE roles SET name=$1, description=$2, updated_at=NOW()
		WHERE id=$3 AND deleted_at IS NULL`,
		strings.TrimSpace(body.Name), body.Description, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Role updated successfully"})
}

func (h *MasterHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE roles SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Role deleted successfully"})
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
		WITH ingredient_stocks AS (
			SELECT inventory_item_id, COALESCE(SUM(quantity), 0) as total_qty
			FROM inventory_stocks
			GROUP BY inventory_item_id
		),
		recipe_agg AS (
			SELECT 
				pr.product_id,
				COUNT(pr.id) as total_ingredients,
				COALESCE(SUM(pr.quantity_required * COALESCE(ii.average_cost, 0)), 0) as total_cogs,
				COALESCE(MIN(FLOOR(COALESCE(ist.total_qty, 0) / NULLIF(pr.quantity_required, 0))), 0) AS max_cookable,
				(
					SELECT ii2.name || ' (Sisa ' || TRIM(TO_CHAR(COALESCE(ist2.total_qty, 0), 'FM999999990.00')) || ' ' || pr2.uom || ')'
					FROM product_recipes pr2
					JOIN inventory_items ii2 ON pr2.inventory_item_id = ii2.id
					LEFT JOIN ingredient_stocks ist2 ON pr2.inventory_item_id = ist2.inventory_item_id
					WHERE pr2.product_id = pr.product_id
					ORDER BY (COALESCE(ist2.total_qty, 0) / NULLIF(pr2.quantity_required, 0)) ASC
					LIMIT 1
				) as bottleneck_item
			FROM product_recipes pr
			JOIN inventory_items ii ON pr.inventory_item_id = ii.id
			LEFT JOIN ingredient_stocks ist ON pr.inventory_item_id = ist.inventory_item_id
			WHERE pr.deleted_at IS NULL
			GROUP BY pr.product_id
		)
		SELECT 
			p.id, p.category_id, c.name as category_name, p.name, p.sku, 
			COALESCE(p.description, ''), p.base_price, p.target_station, 
			COALESCE(p.image_url, ''), p.is_active, p.created_at,
			COALESCE(p.stock, 50) as base_stock,
			COALESCE(p.min_stock, 5) as min_stock,
			ra.total_ingredients,
			ra.total_cogs,
			ra.max_cookable,
			COALESCE(ra.bottleneck_item, '') as bottleneck_ingredient
		FROM products p
		LEFT JOIN menu_categories c ON p.category_id = c.id
		LEFT JOIN recipe_agg ra ON p.id = ra.product_id
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
		var id, catID, catName, name, sku, desc, station, img, bottleneck string
		var price, baseStock, minStock float64
		var isActive bool
		var createdAt interface{}
		var totalIngredients *int
		var totalCOGS, maxCookable *float64

		if err := rows.Scan(&id, &catID, &catName, &name, &sku, &desc, &price, &station, &img, &isActive, &createdAt,
			&baseStock, &minStock, &totalIngredients, &totalCOGS, &maxCookable, &bottleneck); err == nil {

			hasRecipe := totalIngredients != nil && *totalIngredients > 0
			effectiveStock := baseStock
			if hasRecipe && maxCookable != nil {
				effectiveStock = *maxCookable
			}

			cogs := 0.0
			if totalCOGS != nil {
				cogs = *totalCOGS
			}

			grossProfit := price - cogs
			marginPercent := 0.0
			if price > 0 {
				marginPercent = (grossProfit / price) * 100
			}

			ingCount := 0
			if totalIngredients != nil {
				ingCount = *totalIngredients
			}

			list = append(list, map[string]interface{}{
				"id":                    id,
				"category_id":           catID,
				"category_name":         catName,
				"name":                  name,
				"sku":                   sku,
				"description":           desc,
				"base_price":            price,
				"price":                 price,
				"target_station":        station,
				"image_url":             img,
				"is_active":             isActive,
				"stock":                 effectiveStock,
				"base_stock":            baseStock,
				"min_stock":             minStock,
				"has_recipe":            hasRecipe,
				"total_ingredients":     ingCount,
				"cogs":                  cogs,
				"gross_profit":          grossProfit,
				"margin_percent":        marginPercent,
				"is_out_of_stock":       effectiveStock <= 0,
				"bottleneck_ingredient": bottleneck,
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
	if strings.TrimSpace(body.Name) == "" {
		writeError(w, http.StatusBadRequest, "Nama menu wajib diisi")
		return
	}
	if strings.TrimSpace(body.TargetStation) == "" {
		body.TargetStation = "barista"
	}
	if strings.TrimSpace(body.SKU) == "" {
		body.SKU = "PRD-" + strings.ToUpper(uuid.New().String()[:6])
	}

	var catID uuid.UUID
	var err error
	if strings.TrimSpace(body.CategoryID) != "" {
		catID, _ = uuid.Parse(body.CategoryID)
	}
	if catID == uuid.Nil {
		_ = h.db.QueryRow(r.Context(), `SELECT id FROM menu_categories WHERE is_active=true AND deleted_at IS NULL ORDER BY sort_order ASC LIMIT 1`).Scan(&catID)
	}

	newID := uuid.New()
	_, err = h.db.Exec(r.Context(), `
		INSERT INTO products (id, category_id, name, sku, description, base_price, target_station, image_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`,
		newID, catID, body.Name, body.SKU, body.Description, body.BasePrice, body.TargetStation, body.ImageURL, body.IsActive)
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
	var catID uuid.UUID
	if strings.TrimSpace(body.CategoryID) != "" {
		catID, _ = uuid.Parse(body.CategoryID)
	}
	if catID == uuid.Nil {
		_ = h.db.QueryRow(r.Context(), `SELECT category_id FROM products WHERE id=$1`, id).Scan(&catID)
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE products SET category_id=$1, name=$2, sku=$3, description=$4, base_price=$5, target_station=$6, image_url=$7, is_active=$8, updated_at=NOW()
		WHERE id=$9 AND deleted_at IS NULL`,
		catID, body.Name, body.SKU, body.Description, body.BasePrice, body.TargetStation, body.ImageURL, body.IsActive, id)
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
// 5b. MAPPING RESEP & BAHAN BAKU (BILL OF MATERIALS / BOM)
// -----------------------------------------------------------------------------

// GetProductRecipe returns the recipe (BOM) for a specific product including raw material costs and available cookable portions
func (h *MasterHandler) GetProductRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	var prodID, name, sku, img string
	var basePrice float64
	err := h.db.QueryRow(ctx, `
		SELECT id, name, sku, base_price, COALESCE(image_url, '')
		FROM products
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(&prodID, &name, &sku, &basePrice, &img)
	if err != nil {
		writeError(w, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	query := `
		SELECT 
			pr.id, pr.inventory_item_id, ii.name, ii.sku, pr.quantity_required, pr.uom,
			COALESCE(ii.average_cost, 0) as unit_cost,
			(pr.quantity_required * COALESCE(ii.average_cost, 0)) as cost_subtotal,
			COALESCE(st.total_stock, 0) as current_stock,
			COALESCE(pr.instructions, '') as instructions
		FROM product_recipes pr
		JOIN inventory_items ii ON pr.inventory_item_id = ii.id
		LEFT JOIN (
			SELECT inventory_item_id, COALESCE(SUM(quantity), 0) as total_stock
			FROM inventory_stocks
			GROUP BY inventory_item_id
		) st ON pr.inventory_item_id = st.inventory_item_id
		WHERE pr.product_id = $1 AND pr.deleted_at IS NULL
		ORDER BY ii.name ASC`

	rows, err := h.db.Query(ctx, query, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []map[string]interface{}
	var totalCOGS float64
	var minCookable *float64

	for rows.Next() {
		var rID, invID, iName, iSKU, uom, instructions string
		var reqQty, unitCost, costSubtotal, curStock float64
		if err := rows.Scan(&rID, &invID, &iName, &iSKU, &reqQty, &uom, &unitCost, &costSubtotal, &curStock, &instructions); err == nil {
			totalCOGS += costSubtotal
			canMake := 0.0
			if reqQty > 0 {
				canMake = math.Floor(curStock / reqQty)
			}
			if minCookable == nil || canMake < *minCookable {
				canMakeCopy := canMake
				minCookable = &canMakeCopy
			}

			items = append(items, map[string]interface{}{
				"id":                rID,
				"inventory_item_id": invID,
				"item_name":         iName,
				"sku":               iSKU,
				"quantity_required": reqQty,
				"uom":               uom,
				"unit_cost":         unitCost,
				"cost_subtotal":     costSubtotal,
				"current_stock":     curStock,
				"can_make_portions": canMake,
				"instructions":      instructions,
			})
		}
	}
	if items == nil {
		items = []map[string]interface{}{}
	}

	cookablePortions := 0.0
	if minCookable != nil {
		cookablePortions = *minCookable
	}

	grossProfit := basePrice - totalCOGS
	var marginPercent float64
	if basePrice > 0 {
		marginPercent = (grossProfit / basePrice) * 100
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"product_id":           prodID,
		"product_name":         name,
		"sku":                  sku,
		"base_price":           basePrice,
		"image_url":            img,
		"items":                items,
		"total_cogs":           totalCOGS,
		"gross_profit":         grossProfit,
		"gross_margin_percent": marginPercent,
		"cookable_portions":    cookablePortions,
	})
}

// UpdateProductRecipe updates or creates recipe ingredients for a product
func (h *MasterHandler) UpdateProductRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	var body struct {
		Items []struct {
			InventoryItemID  string  `json:"inventory_item_id"`
			QuantityRequired float64 `json:"quantity_required"`
			UOM              string  `json:"uom"`
			Instructions     string  `json:"instructions"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Payload resep tidak valid")
		return
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(ctx)

	// Clear existing recipe items for this product
	_, err = tx.Exec(ctx, "DELETE FROM product_recipes WHERE product_id = $1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Insert new items
	for _, it := range body.Items {
		if it.InventoryItemID == "" || it.QuantityRequired <= 0 {
			continue
		}
		if it.UOM == "" {
			_ = tx.QueryRow(ctx, "SELECT uom FROM inventory_items WHERE id = $1", it.InventoryItemID).Scan(&it.UOM)
			if it.UOM == "" {
				it.UOM = "pcs"
			}
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO product_recipes (id, product_id, inventory_item_id, quantity_required, uom, instructions, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
			uuid.New().String(), id, it.InventoryItemID, it.QuantityRequired, it.UOM, it.Instructions)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Mapping resep dan takaran bahan berhasil disimpan",
	})
}

// ListRecipes returns an overview of all menu products with their recipe mapping status, COGS, and cookable portions
func (h *MasterHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := `
		WITH ingredient_stocks AS (
			SELECT inventory_item_id, COALESCE(SUM(quantity), 0) as total_qty
			FROM inventory_stocks
			GROUP BY inventory_item_id
		),
		recipe_agg AS (
			SELECT 
				pr.product_id,
				COUNT(pr.id) as total_ingredients,
				COALESCE(SUM(pr.quantity_required * COALESCE(ii.average_cost, 0)), 0) as total_cogs,
				COALESCE(MIN(FLOOR(COALESCE(ist.total_qty, 0) / NULLIF(pr.quantity_required, 0))), 0) AS max_cookable,
				(
					SELECT ii2.name || ' (Sisa ' || TRIM(TO_CHAR(COALESCE(ist2.total_qty, 0), 'FM999999990.00')) || ' ' || pr2.uom || ')'
					FROM product_recipes pr2
					JOIN inventory_items ii2 ON pr2.inventory_item_id = ii2.id
					LEFT JOIN ingredient_stocks ist2 ON pr2.inventory_item_id = ist2.inventory_item_id
					WHERE pr2.product_id = pr.product_id
					ORDER BY (COALESCE(ist2.total_qty, 0) / NULLIF(pr2.quantity_required, 0)) ASC
					LIMIT 1
				) as bottleneck_item
			FROM product_recipes pr
			JOIN inventory_items ii ON pr.inventory_item_id = ii.id
			LEFT JOIN ingredient_stocks ist ON pr.inventory_item_id = ist.inventory_item_id
			WHERE pr.deleted_at IS NULL
			GROUP BY pr.product_id
		)
		SELECT 
			p.id, p.name, p.sku, COALESCE(c.name, 'Uncategorized') as category_name,
			p.base_price, COALESCE(p.image_url, ''), p.is_active,
			COALESCE(ra.total_ingredients, 0) as total_ingredients,
			COALESCE(ra.total_cogs, 0) as total_cogs,
			COALESCE(ra.max_cookable, 0) as max_cookable,
			COALESCE(ra.bottleneck_item, '') as bottleneck_ingredient
		FROM products p
		LEFT JOIN menu_categories c ON p.category_id = c.id
		LEFT JOIN recipe_agg ra ON p.id = ra.product_id
		WHERE p.deleted_at IS NULL
		ORDER BY p.name ASC`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, name, sku, catName, img, bottleneck string
		var basePrice, totalCOGS, maxCookable float64
		var totalIngredients int
		var isActive bool

		if err := rows.Scan(&id, &name, &sku, &catName, &basePrice, &img, &isActive,
			&totalIngredients, &totalCOGS, &maxCookable, &bottleneck); err == nil {

			grossProfit := basePrice - totalCOGS
			marginPercent := 0.0
			if basePrice > 0 {
				marginPercent = (grossProfit / basePrice) * 100
			}

			list = append(list, map[string]interface{}{
				"product_id":            id,
				"product_name":          name,
				"sku":                   sku,
				"category":              catName,
				"base_price":            basePrice,
				"image_url":             img,
				"is_active":             isActive,
				"has_recipe":            totalIngredients > 0,
				"total_ingredients":     totalIngredients,
				"cogs":                  totalCOGS,
				"gross_profit":          grossProfit,
				"margin_percent":        marginPercent,
				"cookable_portions":     maxCookable,
				"bottleneck_ingredient": bottleneck,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}

	writeJSON(w, http.StatusOK, list)
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
func (h *MasterHandler) ListTableZones(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT z.id, z.branch_id, COALESCE(b.name, 'All Branches') as branch_name, z.name, COALESCE(z.description,'')
		FROM table_zones z
		LEFT JOIN branches b ON z.branch_id = b.id
		WHERE z.deleted_at IS NULL
		ORDER BY z.name ASC`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, bID, bName, name, desc string
		if err := rows.Scan(&id, &bID, &bName, &name, &desc); err == nil {
			list = append(list, map[string]interface{}{
				"id": id, "branch_id": bID, "branch_name": bName, "name": name, "description": desc,
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *MasterHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT t.id, 
		       COALESCE(t.branch_id::text, ''), 
		       COALESCE(b.name, 'All Branches'), 
		       COALESCE(t.zone_id::text, ''), 
		       COALESCE(z.name, 'No Zone'),
		       t.table_number, 
		       COALESCE(t.capacity, 4), 
		       COALESCE(t.status, 'available'), 
		       COALESCE(t.pos_x, 0), 
		       COALESCE(t.pos_y, 0)
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
		Zone        string `json:"zone"`
		TableNumber string `json:"table_number"`
		Capacity    int    `json:"capacity"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if strings.TrimSpace(body.TableNumber) == "" {
		writeError(w, http.StatusBadRequest, "Nomor meja wajib diisi")
		return
	}
	if body.Capacity <= 0 {
		body.Capacity = 4
	}
	if body.Status == "" {
		body.Status = "available"
	}

	// Validate / Fallback BranchID
	var branchUUID *uuid.UUID
	if bID, err := uuid.Parse(body.BranchID); err == nil {
		branchUUID = &bID
	} else {
		var defaultBranchID uuid.UUID
		err := h.db.QueryRow(r.Context(), `SELECT id FROM branches WHERE is_active = true AND deleted_at IS NULL ORDER BY created_at ASC LIMIT 1`).Scan(&defaultBranchID)
		if err != nil {
			_ = h.db.QueryRow(r.Context(), `SELECT id FROM branches WHERE deleted_at IS NULL ORDER BY created_at ASC LIMIT 1`).Scan(&defaultBranchID)
		}
		if defaultBranchID != uuid.Nil {
			branchUUID = &defaultBranchID
		}
	}
	if branchUUID == nil {
		writeError(w, http.StatusBadRequest, "Cabang tidak valid dan belum ada cabang aktif")
		return
	}

	// Validate / Fallback ZoneID
	var zoneUUID *uuid.UUID
	if zID, err := uuid.Parse(body.ZoneID); err == nil {
		zoneUUID = &zID
	} else if body.Zone != "" {
		var foundZoneID uuid.UUID
		err := h.db.QueryRow(r.Context(), `SELECT id FROM table_zones WHERE name ILIKE '%' || $1 || '%' AND deleted_at IS NULL LIMIT 1`, body.Zone).Scan(&foundZoneID)
		if err == nil {
			zoneUUID = &foundZoneID
		}
	}

	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO cafe_tables (id, branch_id, zone_id, table_number, capacity, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (branch_id, table_number) DO UPDATE
		SET capacity = EXCLUDED.capacity, status = EXCLUDED.status, zone_id = EXCLUDED.zone_id, deleted_at = NULL, updated_at = NOW()`,
		newID, *branchUUID, zoneUUID, body.TableNumber, body.Capacity, body.Status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Table created successfully"})
}

func (h *MasterHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		BranchID    string `json:"branch_id"`
		ZoneID      string `json:"zone_id"`
		Zone        string `json:"zone"`
		TableNumber string `json:"table_number"`
		Capacity    int    `json:"capacity"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}
	if strings.TrimSpace(body.TableNumber) == "" {
		writeError(w, http.StatusBadRequest, "Nomor meja wajib diisi")
		return
	}
	if body.Capacity <= 0 {
		body.Capacity = 4
	}
	if body.Status == "" {
		body.Status = "available"
	}

	var branchUUID *uuid.UUID
	if bID, err := uuid.Parse(body.BranchID); err == nil {
		branchUUID = &bID
	} else {
		var curBranchID uuid.UUID
		_ = h.db.QueryRow(r.Context(), `SELECT branch_id FROM cafe_tables WHERE id = $1`, id).Scan(&curBranchID)
		if curBranchID != uuid.Nil {
			branchUUID = &curBranchID
		}
	}

	var zoneUUID *uuid.UUID
	if zID, err := uuid.Parse(body.ZoneID); err == nil {
		zoneUUID = &zID
	} else if body.Zone != "" {
		var foundZoneID uuid.UUID
		err := h.db.QueryRow(r.Context(), `SELECT id FROM table_zones WHERE name ILIKE '%' || $1 || '%' AND deleted_at IS NULL LIMIT 1`, body.Zone).Scan(&foundZoneID)
		if err == nil {
			zoneUUID = &foundZoneID
		}
	}

	_, err := h.db.Exec(r.Context(), `
		UPDATE cafe_tables SET branch_id=COALESCE($1, branch_id), zone_id=$2, table_number=$3, capacity=$4, status=$5, updated_at=NOW()
		WHERE id=$6 AND deleted_at IS NULL`,
		branchUUID, zoneUUID, body.TableNumber, body.Capacity, body.Status, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Table updated successfully"})
}

func (h *MasterHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE cafe_tables SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Table deleted successfully"})
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

// -----------------------------------------------------------------------------
// 10. CRUD FOR SUPPLIERS, WAREHOUSES, SHIFTS, DEPARTMENTS, POSITIONS
// -----------------------------------------------------------------------------
func (h *MasterHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		ContactPerson string `json:"contact_person"`
		Phone         string `json:"phone"`
		Email         string `json:"email"`
		Address       string `json:"address"`
		IsActive      bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if strings.TrimSpace(body.Name) == "" {
		writeError(w, http.StatusBadRequest, "Nama supplier wajib diisi")
		return
	}
	if strings.TrimSpace(body.Code) == "" {
		body.Code = "SUP-" + strings.ToUpper(uuid.New().String()[:6])
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO suppliers (id, code, name, contact_person, phone, email, address, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`,
		newID, body.Code, body.Name, body.ContactPerson, body.Phone, body.Email, body.Address, body.IsActive)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Supplier created successfully"})
}

func (h *MasterHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		ContactPerson string `json:"contact_person"`
		Phone         string `json:"phone"`
		Email         string `json:"email"`
		Address       string `json:"address"`
		IsActive      bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE suppliers SET code=$1, name=$2, contact_person=$3, phone=$4, email=$5, address=$6, is_active=$7, updated_at=NOW()
		WHERE id=$8 AND deleted_at IS NULL`,
		body.Code, body.Name, body.ContactPerson, body.Phone, body.Email, body.Address, body.IsActive, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Supplier updated successfully"})
}

func (h *MasterHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE suppliers SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Supplier deleted successfully"})
}

func (h *MasterHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BranchID string `json:"branch_id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Address  string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if strings.TrimSpace(body.Name) == "" {
		writeError(w, http.StatusBadRequest, "Nama gudang wajib diisi")
		return
	}
	var branchUUID uuid.UUID
	if bID, err := uuid.Parse(body.BranchID); err == nil {
		branchUUID = bID
	} else {
		_ = h.db.QueryRow(r.Context(), `SELECT id FROM branches WHERE is_active=true AND deleted_at IS NULL ORDER BY created_at ASC LIMIT 1`).Scan(&branchUUID)
	}
	if branchUUID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "Cabang tidak valid")
		return
	}
	if body.Type == "" {
		body.Type = "main"
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO warehouses (id, branch_id, name, type, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		newID, branchUUID, body.Name, body.Type, body.Address)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Warehouse created successfully"})
}

func (h *MasterHandler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		BranchID string `json:"branch_id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Address  string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	var branchUUID *uuid.UUID
	if bID, err := uuid.Parse(body.BranchID); err == nil {
		branchUUID = &bID
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE warehouses SET branch_id=COALESCE($1, branch_id), name=$2, type=$3, address=$4, updated_at=NOW()
		WHERE id=$5 AND deleted_at IS NULL`,
		branchUUID, body.Name, body.Type, body.Address, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Warehouse updated successfully"})
}

func (h *MasterHandler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE warehouses SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Warehouse deleted successfully"})
}

func (h *MasterHandler) CreateShift(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BranchID  string `json:"branch_id"`
		Name      string `json:"name"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Color     string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if strings.TrimSpace(body.Name) == "" {
		writeError(w, http.StatusBadRequest, "Nama shift wajib diisi")
		return
	}
	var branchUUID uuid.UUID
	if bID, err := uuid.Parse(body.BranchID); err == nil {
		branchUUID = bID
	} else {
		_ = h.db.QueryRow(r.Context(), `SELECT id FROM branches WHERE is_active=true AND deleted_at IS NULL ORDER BY created_at ASC LIMIT 1`).Scan(&branchUUID)
	}
	if branchUUID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "Cabang tidak valid")
		return
	}
	if body.StartTime == "" {
		body.StartTime = "08:00:00"
	}
	if body.EndTime == "" {
		body.EndTime = "17:00:00"
	}
	if body.Color == "" {
		body.Color = "#2563eb"
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO work_shifts (id, branch_id, name, start_time, end_time, color, created_at, updated_at)
		VALUES ($1, $2, $3, $4::time, $5::time, $6, NOW(), NOW())`,
		newID, branchUUID, body.Name, body.StartTime, body.EndTime, body.Color)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Shift created successfully"})
}

func (h *MasterHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		BranchID  string `json:"branch_id"`
		Name      string `json:"name"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Color     string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	var branchUUID *uuid.UUID
	if bID, err := uuid.Parse(body.BranchID); err == nil {
		branchUUID = &bID
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE work_shifts SET branch_id=COALESCE($1, branch_id), name=$2, start_time=$3::time, end_time=$4::time, color=$5, updated_at=NOW()
		WHERE id=$6 AND deleted_at IS NULL`,
		branchUUID, body.Name, body.StartTime, body.EndTime, body.Color, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Shift updated successfully"})
}

func (h *MasterHandler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE work_shifts SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Shift deleted successfully"})
}

func (h *MasterHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if strings.TrimSpace(body.Name) == "" {
		writeError(w, http.StatusBadRequest, "Nama departemen wajib diisi")
		return
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO departments (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())`,
		newID, strings.TrimSpace(body.Name), body.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Department created successfully"})
}

func (h *MasterHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE departments SET name=$1, description=$2, updated_at=NOW()
		WHERE id=$3 AND deleted_at IS NULL`,
		strings.TrimSpace(body.Name), body.Description, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Department updated successfully"})
}

func (h *MasterHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE departments SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Department deleted successfully"})
}

func (h *MasterHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DepartmentID string  `json:"department_id"`
		Title        string  `json:"title"`
		Level        int     `json:"level"`
		BaseSalary   float64 `json:"base_salary"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if strings.TrimSpace(body.Title) == "" {
		writeError(w, http.StatusBadRequest, "Nama posisi/jabatan wajib diisi")
		return
	}
	var deptUUID uuid.UUID
	if dID, err := uuid.Parse(body.DepartmentID); err == nil {
		deptUUID = dID
	} else {
		_ = h.db.QueryRow(r.Context(), `SELECT id FROM departments WHERE deleted_at IS NULL ORDER BY created_at ASC LIMIT 1`).Scan(&deptUUID)
	}
	if deptUUID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "Departemen tidak valid")
		return
	}
	if body.Level <= 0 {
		body.Level = 1
	}
	newID := uuid.New()
	_, err := h.db.Exec(r.Context(), `
		INSERT INTO positions (id, department_id, title, level, base_salary, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		newID, deptUUID, body.Title, body.Level, body.BaseSalary)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": newID, "message": "Position created successfully"})
}

func (h *MasterHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		DepartmentID string  `json:"department_id"`
		Title        string  `json:"title"`
		Level        int     `json:"level"`
		BaseSalary   float64 `json:"base_salary"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	var deptUUID *uuid.UUID
	if dID, err := uuid.Parse(body.DepartmentID); err == nil {
		deptUUID = &dID
	}
	_, err := h.db.Exec(r.Context(), `
		UPDATE positions SET department_id=COALESCE($1, department_id), title=$2, level=$3, base_salary=$4, updated_at=NOW()
		WHERE id=$5 AND deleted_at IS NULL`,
		deptUUID, body.Title, body.Level, body.BaseSalary, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Position updated successfully"})
}

func (h *MasterHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `UPDATE positions SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Position deleted successfully"})
}
