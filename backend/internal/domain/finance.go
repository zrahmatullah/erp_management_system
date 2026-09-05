package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountType string
const (
	AccAsset     AccountType = "asset"
	AccLiability AccountType = "liability"
	AccEquity    AccountType = "equity"
	AccRevenue   AccountType = "revenue"
	AccExpense   AccountType = "expense"
)

type ChartOfAccount struct {
	BaseEntity
	AccountCode string      `json:"account_code"`
	Name        string      `json:"name"`
	Type        AccountType `json:"type"`
	ParentID    *uuid.UUID  `json:"parent_id"`
	IsActive    bool        `json:"is_active"`
	Description string      `json:"description"`
}

type JournalEntry struct {
	BaseEntity
	EntryNumber  string     `json:"entry_number"`
	Date         time.Time  `json:"date"`
	Description  string     `json:"description"`
	SourceModule string     `json:"source_module"`
	SourceID     *uuid.UUID `json:"source_id"`
	Status       string     `json:"status"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	ApprovedBy   *uuid.UUID `json:"approved_by"`
}

type JournalEntryLine struct {
	BaseEntity
	JournalEntryID uuid.UUID  `json:"journal_entry_id"`
	AccountID      uuid.UUID  `json:"account_id"`
	Description    string     `json:"description"`
	Debit          float64    `json:"debit"`
	Credit         float64    `json:"credit"`
	BranchID       *uuid.UUID `json:"branch_id"`
}

type ExpenseStatus string
const (
	ExpDraft           ExpenseStatus = "draft"
	ExpSubmitted       ExpenseStatus = "submitted"
	ExpManagerApproved ExpenseStatus = "manager_approved"
	ExpFinanceApproved ExpenseStatus = "finance_approved"
	ExpRejected        ExpenseStatus = "rejected"
	ExpReimbursed      ExpenseStatus = "reimbursed"
)

type Expense struct {
	BaseEntity
	BranchID    uuid.UUID     `json:"branch_id"`
	CategoryID  uuid.UUID     `json:"category_id"`
	EmployeeID  uuid.UUID     `json:"employee_id"`
	Amount      float64       `json:"amount"`
	Date        time.Time     `json:"date"`
	Description string        `json:"description"`
	ReceiptURL  string        `json:"receipt_url"`
	Status      ExpenseStatus `json:"status"`
	SubmittedBy uuid.UUID     `json:"submitted_by"`
}

type ExpenseCategory struct {
	BaseEntity
	Name        string    `json:"name"`
	AccountID   uuid.UUID `json:"account_id"`
	Description string    `json:"description"`
}

type ExpenseApproval struct {
	BaseEntity
	ExpenseID  uuid.UUID `json:"expense_id"`
	ApproverID uuid.UUID `json:"approver_id"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
	ActionDate time.Time `json:"action_date"`
}

type Budget struct {
	BaseEntity
	BranchID     uuid.UUID `json:"branch_id"`
	DepartmentID uuid.UUID `json:"department_id"`
	Year         int       `json:"year"`
	Month        int       `json:"month"`
	Status       string    `json:"status"`
	CreatedBy    uuid.UUID `json:"created_by"`
}

type BudgetItem struct {
	BaseEntity
	BudgetID      uuid.UUID `json:"budget_id"`
	AccountID     uuid.UUID `json:"account_id"`
	PlannedAmount float64   `json:"planned_amount"`
	ActualAmount  float64   `json:"actual_amount"`
}

type TaxSetting struct {
	BaseEntity
	TaxType   string    `json:"tax_type"`
	Rate      float64   `json:"rate"`
	AccountID uuid.UUID `json:"account_id"`
	IsActive  bool      `json:"is_active"`
}

type TaxInvoice struct {
	BaseEntity
	InvoiceNumber string    `json:"invoice_number"`
	OrderID       uuid.UUID `json:"order_id"`
	TaxableAmount float64   `json:"taxable_amount"`
	TaxAmount     float64   `json:"tax_amount"`
	IssuedDate    time.Time `json:"issued_date"`
}

type FinanceRepository interface {
	CreateCOA(ctx context.Context, coa *ChartOfAccount) error
	CreateJournalEntry(ctx context.Context, entry *JournalEntry) error
	CreateExpense(ctx context.Context, expense *Expense) error
}

type FinanceUsecase interface {
	SubmitExpense(ctx context.Context, req Expense) (Expense, error)
	ApproveExpense(ctx context.Context, expenseID, approverID uuid.UUID, status ExpenseStatus) error
}
