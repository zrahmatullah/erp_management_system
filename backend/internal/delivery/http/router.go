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
	p2pHandler *handler.P2PHandler,
	hrisHandler *handler.HRISHandler,
	reportHandler *handler.ReportHandler,
	posShiftHandler *handler.POSShiftHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(appMiddleware.Recovery)
	r.Use(appMiddleware.SecurityHeaders)
	r.Use(appMiddleware.RequestSizeLimit(10 << 20)) // OWASP A04/A05: 10MB ceiling to prevent buffer overflow & DoS
	r.Use(appMiddleware.RequestLogger)
	r.Use(appMiddleware.CORS())
	r.Use(appMiddleware.RateLimit) // OWASP A04: Global rate limit

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"system":  "Cafe ERP System",
			"version": "1.0.0",
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Auth Routes (OWASP A07: Brute-Force Rate Limited)
		r.Route("/auth", func(r chi.Router) {
			r.With(appMiddleware.AuthRateLimit).Post("/login", authHandler.Login)
			r.With(appMiddleware.AuthRateLimit).Post("/refresh", authHandler.RefreshToken)
			r.Post("/logout", authHandler.Logout)
			r.Post("/forgot-password", authHandler.ForgotPassword)
			r.Post("/reset-password", authHandler.ResetPassword)
			r.With(appMiddleware.JWTAuth).Get("/profile", authHandler.GetProfile)
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
			r.Get("/orders/{id}", opHandler.GetPOSOrderByID)
			r.Put("/orders/{id}/status", opHandler.UpdatePOSOrderStatus)
			r.Post("/orders/{id}/cancel", opHandler.CancelPOSOrder)
			r.Post("/orders/{id}/void", opHandler.VoidPOSOrder)
			r.Post("/orders", opHandler.CreatePOSOrder)
			r.Post("/orders/{id}/pay", opHandler.PayOrder)
			r.Get("/takeaways", opHandler.GetTakeawayOrders)
			r.Get("/transactions", opHandler.GetPOSTransactions)

			// Cashier Shifts & Blind Cash Balancing (F-01)
			r.Route("/shifts", func(r chi.Router) {
				r.Get("/current", posShiftHandler.GetCurrentShift)
				r.Post("/open", posShiftHandler.OpenShift)
				r.Post("/{id}/cash-movement", posShiftHandler.RecordCashMovement)
				r.Post("/{id}/close", posShiftHandler.CloseShift)
				r.Get("/{id}/summary", posShiftHandler.GetShiftSummary)
				r.Get("/", posShiftHandler.ListShifts)
			})
		})

		// Orders Service Aliases (Compatible with order.service.ts)
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", opHandler.GetPOSOrders)
			r.Post("/", opHandler.CreatePOSOrder)
			r.Get("/active", opHandler.GetPOSOrdersActive)
			r.Get("/kitchen-queue", opHandler.GetKDSTickets)
			r.Get("/{id}", opHandler.GetPOSOrderByID)
			r.Put("/{id}/status", opHandler.UpdatePOSOrderStatus)
			r.Post("/{id}/cancel", opHandler.CancelPOSOrder)
			r.Post("/{id}/void", opHandler.VoidPOSOrder)
			r.Post("/{id}/pay", opHandler.PayOrder)
		})

		// Roles & Permissions Service Aliases (Compatible with role.service.ts)
		r.Route("/roles", func(r chi.Router) {
			r.Use(appMiddleware.JWTAuth)
			r.Get("/", masterHandler.ListRoles)
			r.Post("/", masterHandler.CreateRole)
			r.Put("/{id}", masterHandler.UpdateRole)
			r.Delete("/{id}", masterHandler.DeleteRole)
			r.Get("/{id}/permissions", masterHandler.GetRolePermissions)
			r.Put("/{id}/permissions", masterHandler.UpdateRolePermissions)
			r.Post("/{id}/permissions", masterHandler.UpdateRolePermissions)
		})
		r.With(appMiddleware.JWTAuth).Get("/permissions", masterHandler.ListAllPermissions)

		// Users Service Aliases (Compatible with user.service.ts)
		r.Route("/users", func(r chi.Router) {
			r.Use(appMiddleware.JWTAuth)
			r.Get("/", masterHandler.ListUsers)
			r.Post("/", masterHandler.CreateUser)
			r.Get("/{id}", masterHandler.GetUserByID)
			r.Put("/{id}", masterHandler.UpdateUser)
			r.Delete("/{id}", masterHandler.DeleteUser)
			r.Post("/{id}/role", masterHandler.AssignUserRole)
		})

		// Kitchen Display System
		r.Route("/kds", func(r chi.Router) {
			r.Get("/tickets", opHandler.GetKDSTickets)
			r.Put("/items/{id}/status", opHandler.UpdateKDSStatus)
		})

		// Inventory Operations
		// Inventory & Procure-to-Pay (P2P) Operations
		r.Route("/inventory", func(r chi.Router) {
			r.Get("/stocks", opHandler.GetInventoryStocks)
			r.With(appMiddleware.HeavyLimitMiddleware).Get("/stock-movements", opHandler.GetStockMovements)
			r.Get("/purchase-orders", opHandler.GetPurchaseOrders)
			r.Post("/purchase-orders", opHandler.CreatePurchaseOrder)
			r.Get("/purchase-orders/{id}", opHandler.GetPurchaseOrderDetail)
			r.Put("/purchase-orders/{id}/status", opHandler.UpdatePurchaseOrderStatus)
			r.With(appMiddleware.HeavyLimitMiddleware).Get("/opnames", opHandler.GetStockOpnames)
			r.With(appMiddleware.HeavyLimitMiddleware).Post("/opnames", opHandler.CreateStockOpname)

			// P2P Full Cycle
			r.Get("/pr", p2pHandler.GetPurchaseRequisitions)
			r.Post("/pr", p2pHandler.CreatePurchaseRequisition)
			r.Post("/pr/from-menu", p2pHandler.CreatePRFromDepletedMenu)
			r.Put("/pr/{id}/status", p2pHandler.UpdatePurchaseRequisitionStatus)
			r.Post("/pr/{id}/convert-to-po", p2pHandler.ConvertPRToPO)
			r.Post("/pr/{id}/convert-po", p2pHandler.ConvertPRToPO)
			r.Get("/grn", p2pHandler.GetGoodsReceiptNotes)
			r.Post("/grn", p2pHandler.CreateGoodsReceiptNote)
			r.Get("/invoices", p2pHandler.GetVendorInvoices)
			r.Post("/invoices", p2pHandler.CreateVendorInvoice)
			r.Get("/vendor-payments", p2pHandler.GetVendorPayments)
			r.Post("/vendor-payments", p2pHandler.CreateVendorPayment)
		})

		// HRIS Operations (Protected with JWTAuth)
		r.Route("/hris", func(r chi.Router) {
			r.Use(appMiddleware.JWTAuth)

			// Employee Self-Service Portal
			r.Get("/my-portal", hrisHandler.GetMyEmployeePortal)

			// Employees
			r.Get("/employees", hrisHandler.GetEmployees)
			r.Post("/employees", hrisHandler.CreateEmployee)
			r.Put("/employees/{id}", hrisHandler.UpdateEmployee)
			r.Delete("/employees/{id}", hrisHandler.DeleteEmployee)

			// Attendances
			r.Get("/attendances", hrisHandler.GetAttendances)
			r.Get("/attendances/today-status", hrisHandler.GetTodayAttendanceStatus)
			r.Post("/attendances/clock-in", hrisHandler.ClockIn)
			r.Post("/attendances/clock-out", hrisHandler.ClockOut)

			// Leaves & Approval
			r.Get("/leaves", hrisHandler.GetLeaves)
			r.Post("/leaves", hrisHandler.CreateLeave)
			r.Put("/leaves/{id}/status", hrisHandler.UpdateLeaveStatus)

			// Shift Schedules
			r.Get("/schedules", hrisHandler.GetShiftSchedules)
			r.Post("/schedules", hrisHandler.SetShiftSchedule)

			// Payroll (Resource Intensive: Protected with HeavyLimitMiddleware)
			r.Get("/payrolls", hrisHandler.GetPayrolls)
			r.With(appMiddleware.HeavyLimitMiddleware).Post("/payrolls/run", hrisHandler.RunPayroll)
			r.With(appMiddleware.HeavyLimitMiddleware).Post("/payrolls/batch-approve", hrisHandler.BatchApprovePayroll)
		})

		// Finance Operations
		r.Route("/finance", func(r chi.Router) {
			r.Get("/overview", opHandler.GetFinanceOverview)
			r.Get("/journals", opHandler.GetJournalEntries)
			r.Post("/journals", opHandler.CreateJournalEntry)
			r.Get("/journals/{id}/lines", opHandler.GetJournalEntryLines)

			// Chart of Accounts
			r.Get("/coa", opHandler.GetChartOfAccounts)
			r.Post("/coa", opHandler.CreateChartOfAccount)

			// Expenses & Petty Cash
			r.Get("/expenses", opHandler.GetExpenses)
			r.Post("/expenses", opHandler.CreateExpense)
			r.Put("/expenses/{id}/status", opHandler.UpdateExpenseStatus)

			// Tax Management
			r.Get("/tax/summary", opHandler.GetTaxSummary)

			// Budgets & Cost Control
			r.Get("/budgets", opHandler.GetBudgets)
			r.Post("/budgets", opHandler.SetBudget)

			// Bank Reconciliation
			r.Get("/reconciliation", opHandler.GetBankReconciliation)
			r.With(appMiddleware.HeavyLimitMiddleware).Post("/reconciliation/match", opHandler.MatchBankReconciliation)
		})

		// Reports & Analytics Routes (Protected with HeavyLimitMiddleware)
		r.Route("/reports", func(r chi.Router) {
			r.Use(appMiddleware.HeavyLimitMiddleware)
			r.Get("/sales", reportHandler.GetSalesReport)
			r.Get("/financial", reportHandler.GetFinancialReport)
			r.Get("/inventory", reportHandler.GetInventoryReport)
			r.Get("/hr", reportHandler.GetHRReport)
			r.Post("/custom", reportHandler.GenerateCustomReport)
			r.Get("/schedules", reportHandler.GetReportSchedules)
			r.Post("/schedules", reportHandler.CreateReportSchedule)
			r.Delete("/schedules/{id}", reportHandler.DeleteReportSchedule)
			r.Get("/templates", reportHandler.GetCustomReportTemplates)
			r.Get("/shifts/audit", posShiftHandler.ListShifts)
		})

		// Master Data Routes (Super Admin CRUD - OWASP A01 Protected)
		r.Route("/master", func(r chi.Router) {
			r.Use(appMiddleware.JWTAuth)
			// Branches
			r.Get("/branches", masterHandler.ListBranches)
			r.Post("/branches", masterHandler.CreateBranch)
			r.Put("/branches/{id}", masterHandler.UpdateBranch)
			r.Delete("/branches/{id}", masterHandler.DeleteBranch)

			// Users
			r.Get("/users", masterHandler.ListUsers)
			r.Post("/users", masterHandler.CreateUser)
			r.Get("/users/{id}", masterHandler.GetUserByID)
			r.Put("/users/{id}", masterHandler.UpdateUser)
			r.Delete("/users/{id}", masterHandler.DeleteUser)
			r.Post("/users/{id}/role", masterHandler.AssignUserRole)

			// Roles & Permissions
			r.Get("/roles", masterHandler.ListRoles)
			r.Post("/roles", masterHandler.CreateRole)
			r.Put("/roles/{id}", masterHandler.UpdateRole)
			r.Delete("/roles/{id}", masterHandler.DeleteRole)
			r.Get("/roles/{id}/permissions", masterHandler.GetRolePermissions)
			r.Put("/roles/{id}/permissions", masterHandler.UpdateRolePermissions)
			r.Post("/roles/{id}/permissions", masterHandler.UpdateRolePermissions)
			r.Get("/permissions", masterHandler.ListAllPermissions)

			// Menu Categories
			r.Get("/categories", masterHandler.ListCategories)
			r.Post("/categories", masterHandler.CreateCategory)
			r.Put("/categories/{id}", masterHandler.UpdateCategory)
			r.Delete("/categories/{id}", masterHandler.DeleteCategory)

			// Products & Recipes (BOM)
			r.Get("/products", masterHandler.ListProducts)
			r.Post("/products", masterHandler.CreateProduct)
			r.Put("/products/{id}", masterHandler.UpdateProduct)
			r.Delete("/products/{id}", masterHandler.DeleteProduct)
			r.Get("/products/{id}/recipe", masterHandler.GetProductRecipe)
			r.Put("/products/{id}/recipe", masterHandler.UpdateProductRecipe)
			r.Get("/recipes", masterHandler.ListRecipes)

			// Inventory Items
			r.Get("/inventory-items", masterHandler.ListInventoryItems)
			r.Post("/inventory-items", masterHandler.CreateInventoryItem)
			r.Put("/inventory-items/{id}", masterHandler.UpdateInventoryItem)
			r.Delete("/inventory-items/{id}", masterHandler.DeleteInventoryItem)

			// Tables & Zones
			r.Get("/tables", masterHandler.ListTables)
			r.Post("/tables", masterHandler.CreateTable)
			r.Put("/tables/{id}", masterHandler.UpdateTable)
			r.Delete("/tables/{id}", masterHandler.DeleteTable)
			r.Get("/zones", masterHandler.ListTableZones)

			// Chart of Accounts
			r.Get("/accounts", masterHandler.ListAccounts)
			r.Post("/accounts", masterHandler.CreateAccount)
			r.Put("/accounts/{id}", masterHandler.UpdateAccount)
			r.Delete("/accounts/{id}", masterHandler.DeleteAccount)

			// Suppliers
			r.Get("/suppliers", masterHandler.ListSuppliers)
			r.Post("/suppliers", masterHandler.CreateSupplier)
			r.Put("/suppliers/{id}", masterHandler.UpdateSupplier)
			r.Delete("/suppliers/{id}", masterHandler.DeleteSupplier)

			// Warehouses
			r.Get("/warehouses", masterHandler.ListWarehouses)
			r.Post("/warehouses", masterHandler.CreateWarehouse)
			r.Put("/warehouses/{id}", masterHandler.UpdateWarehouse)
			r.Delete("/warehouses/{id}", masterHandler.DeleteWarehouse)

			// Shifts
			r.Get("/shifts", masterHandler.ListShifts)
			r.Post("/shifts", masterHandler.CreateShift)
			r.Put("/shifts/{id}", masterHandler.UpdateShift)
			r.Delete("/shifts/{id}", masterHandler.DeleteShift)

			// Departments
			r.Get("/departments", masterHandler.ListDepartments)
			r.Post("/departments", masterHandler.CreateDepartment)
			r.Put("/departments/{id}", masterHandler.UpdateDepartment)
			r.Delete("/departments/{id}", masterHandler.DeleteDepartment)

			// Positions
			r.Get("/positions", masterHandler.ListPositions)
			r.Post("/positions", masterHandler.CreatePosition)
			r.Put("/positions/{id}", masterHandler.UpdatePosition)
			r.Delete("/positions/{id}", masterHandler.DeletePosition)

			// Company Profile
			r.Get("/company-profile", p2pHandler.GetCompanyProfile)
			r.Put("/company-profile", p2pHandler.UpdateCompanyProfile)
		})
	})

	return r
}
