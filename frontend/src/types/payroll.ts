export enum PayrollStatus { DRAFT, PROCESSED, PAID }
export interface PayrollPeriod { id: string; startDate: string; endDate: string; }
export interface Payslip { id: string; employeeId: string; status: PayrollStatus; netAmount: number; }
export interface PayrollItem { id: string; payslipId: string; amount: number; }
export interface PayrollDeduction { id: string; payslipId: string; amount: number; }
export interface PayrollFilters { periodId?: string; }
