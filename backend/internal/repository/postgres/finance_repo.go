package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type FinanceRepository struct {
	db *pgxpool.Pool
}

func NewFinanceRepository(db *pgxpool.Pool) domain.FinanceRepository {
	return &FinanceRepository{db: db}
}

func (r *FinanceRepository) CreateCOA(ctx context.Context, coa *domain.ChartOfAccount) error {
	query := `
		INSERT INTO chart_of_accounts (id, account_code, name, type, parent_id, is_active, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING created_at, updated_at`

	if coa.ID == uuid.Nil {
		coa.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		coa.ID, coa.AccountCode, coa.Name, coa.Type, coa.ParentID, coa.IsActive, coa.Description,
	).Scan(&coa.CreatedAt, &coa.UpdatedAt)
}

func (r *FinanceRepository) CreateJournalEntry(ctx context.Context, entry *domain.JournalEntry) error {
	query := `
		INSERT INTO journal_entries (id, entry_number, date, description, source_module, source_id, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING created_at, updated_at`

	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		entry.ID, entry.EntryNumber, entry.Date, entry.Description, entry.SourceModule, entry.SourceID, entry.Status, entry.CreatedBy,
	).Scan(&entry.CreatedAt, &entry.UpdatedAt)
}

func (r *FinanceRepository) CreateExpense(ctx context.Context, expense *domain.Expense) error {
	query := `
		INSERT INTO expenses (id, branch_id, category_id, employee_id, amount, date, description, receipt_url, status, submitted_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING created_at, updated_at`

	if expense.ID == uuid.Nil {
		expense.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		expense.ID, expense.BranchID, expense.CategoryID, expense.EmployeeID, expense.Amount, expense.Date, expense.Description, expense.ReceiptURL, expense.Status, expense.SubmittedBy,
	).Scan(&expense.CreatedAt, &expense.UpdatedAt)
}
