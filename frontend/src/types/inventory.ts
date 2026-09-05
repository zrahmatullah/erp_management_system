export interface InventoryItem { id: string; name: string; currentStock: number; }
export interface Warehouse { id: string; name: string; }
export interface WarehouseLocation { id: string; warehouseId: string; name: string; }
export enum StockMovementType { IN, OUT, TRANSFER, ADJUSTMENT }
export interface StockMovement { id: string; type: StockMovementType; quantity: number; }
export interface Supplier { id: string; name: string; }
export interface SupplierContact { id: string; supplierId: string; phone: string; }
export enum POStatus { DRAFT, SUBMITTED, APPROVED, RECEIVED, CANCELLED }
export interface PurchaseOrder { id: string; status: POStatus; supplierId: string; }
export interface PurchaseOrderItem { id: string; poId: string; quantity: number; }
export interface GoodsReceipt { id: string; poId: string; }
export interface GoodsReceiptItem { id: string; receiptId: string; quantity: number; }
export enum OpnameStatus { DRAFT, IN_PROGRESS, COMPLETED }
export interface StockOpname { id: string; status: OpnameStatus; }
export interface StockOpnameItem { id: string; expected: number; actual: number; }
export interface InventoryFilters { search?: string; }
export interface POFilters { status?: POStatus; }
