package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"cafe-erp-system/backend/internal/delivery/http/handler"
	appMiddleware "cafe-erp-system/backend/internal/delivery/http/middleware"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	masterHandler *handler.MasterHandler,
	opHandler *handler.OperationalHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(appMiddleware.Recovery)
	r.Use(appMiddleware.RequestLogger)
	r.Use(appMiddleware.CORS())

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"system":  "Cafe ERP System",
			"version": "1.0.0",
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Auth Routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.RefreshToken)
		})

		// Dashboard Stats
		r.Get("/dashboard/stats", opHandler.GetDashboardStats)

		// POS & Orders
		r.Route("/pos", func(r chi.Router) {
			r.Get("/products", opHandler.GetPOSProducts)
			r.Get("/tables", opHandler.GetPOSTables)
			r.Get("/tables/{id}/order", opHandler.GetTableOrderDetail)
			r.Put("/tables/{id}/status", opHandler.UpdateTableStatus)
			r.Get("/orders", opHandler.GetPOSOrders)
			r.Get("/orders/active", opHandler.GetPOSOrdersActive)
			r.Post("/orders", opHandler.CreatePOSOrder)
			r.Post("/orders/{id}/pay", opHandler.PayOrder)
			r.Get("/takeaways", opHandler.GetTakeawayOrders)
			r.Get("/transactions", opHandler.GetPOSTransactions)
		})

		// Kitchen Display System
		r.Route("/kds", func(r chi.Router) {
			r.Get("/tickets", opHandler.GetKDSTickets)
			r.Put("/items/{id}/status", opHandler.UpdateKDSStatus)
		})

		// Inventory Operations
		r.Route("/inventory", func(r chi.Router) {
			r.Get("/stocks", opHandler.GetInventoryStocks)
			r.Get("/purchase-orders", opHandler.GetPurchaseOrders)
			r.Get("/opnames", opHandler.GetStockOpnames)
		})

		// HRIS Operations
		r.Route("/hris", func(r chi.Router) {
			r.Get("/employees", opHandler.GetEmployees)
			r.Get("/attendances", opHandler.GetAttendances)
			r.Get("/leaves", opHandler.GetLeaves)
			r.Get("/payrolls", opHandler.GetPayrolls)
			r.Get("/schedules", opHandler.GetShiftSchedules)
		})

		// Finance Operations
		r.Route("/finance", func(r chi.Router) {
			r.Get("/overview", opHandler.GetFinanceOverview)
			r.Get("/journals", opHandler.GetJournalEntries)
		})

		// Master Data Routes (Super Admin CRUD)
		r.Route("/master", func(r chi.Router) {
			// Branches
			r.Get("/branches", masterHandler.ListBranches)
			r.Post("/branches", masterHandler.CreateBranch)
			r.Put("/branches/{id}", masterHandler.UpdateBranch)
			r.Delete("/branches/{id}", masterHandler.DeleteBranch)

			// Users
			r.Get("/users", masterHandler.ListUsers)
			r.Post("/users", masterHandler.CreateUser)
			r.Put("/users/{id}", masterHandler.UpdateUser)
			r.Delete("/users/{id}", masterHandler.DeleteUser)

			// Roles & Permissions
			r.Get("/roles", masterHandler.ListRoles)
			r.Get("/roles/{id}/permissions", masterHandler.GetRolePermissions)
			r.Put("/roles/{id}/permissions", masterHandler.UpdateRolePermissions)

			// Menu Categories
			r.Get("/categories", masterHandler.ListCategories)
			r.Post("/categories", masterHandler.CreateCategory)
			r.Put("/categories/{id}", masterHandler.UpdateCategory)
			r.Delete("/categories/{id}", masterHandler.DeleteCategory)

			// Products
			r.Get("/products", masterHandler.ListProducts)
			r.Post("/products", masterHandler.CreateProduct)
			r.Put("/products/{id}", masterHandler.UpdateProduct)
			r.Delete("/products/{id}", masterHandler.DeleteProduct)

			// Inventory Items
			r.Get("/inventory-items", masterHandler.ListInventoryItems)
			r.Post("/inventory-items", masterHandler.CreateInventoryItem)
			r.Put("/inventory-items/{id}", masterHandler.UpdateInventoryItem)
			r.Delete("/inventory-items/{id}", masterHandler.DeleteInventoryItem)

			// Tables
			r.Get("/tables", masterHandler.ListTables)
			r.Post("/tables", masterHandler.CreateTable)
			r.Put("/tables/{id}", masterHandler.UpdateTable)
			r.Delete("/tables/{id}", masterHandler.DeleteTable)

			// Chart of Accounts
			r.Get("/accounts", masterHandler.ListAccounts)
			r.Post("/accounts", masterHandler.CreateAccount)
			r.Put("/accounts/{id}", masterHandler.UpdateAccount)
			r.Delete("/accounts/{id}", masterHandler.DeleteAccount)

			// Other Masters
			r.Get("/suppliers", masterHandler.ListSuppliers)
			r.Get("/warehouses", masterHandler.ListWarehouses)
			r.Get("/shifts", masterHandler.ListShifts)
			r.Get("/departments", masterHandler.ListDepartments)
			r.Get("/positions", masterHandler.ListPositions)
		})
	})

	return r
}
