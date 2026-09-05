export interface SalesReportData { date: string; revenue: number; }
export interface SalesReportFilters { startDate: string; endDate: string; }
export interface FinancialReportData { category: string; amount: number; }
export interface InventoryReportData { item: string; stock: number; }
export interface HRReportData { metric: string; value: number; }
export interface DashboardData { revenue: number; }
export interface DashboardMetrics { totalRevenue: number; orders: number; }
