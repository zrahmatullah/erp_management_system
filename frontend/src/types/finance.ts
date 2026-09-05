export enum AccountType { ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE }
export interface ChartOfAccount { id: string; name: string; type: AccountType; }
export interface JournalEntry { id: string; date: string; reference: string; }
export interface JournalEntryLine { id: string; entryId: string; accountId: string; debit: number; credit: number; }
export enum ExpenseStatus { DRAFT, SUBMITTED, APPROVED, REJECTED }
export interface Expense { id: string; amount: number; status: ExpenseStatus; }
export interface ExpenseCategory { id: string; name: string; }
export interface ExpenseApproval { id: string; expenseId: string; }
export interface Budget { id: string; year: number; }
export interface BudgetItem { id: string; budgetId: string; amount: number; }
export interface TaxSetting { id: string; rate: number; }
export interface TaxInvoice { id: string; amount: number; }
