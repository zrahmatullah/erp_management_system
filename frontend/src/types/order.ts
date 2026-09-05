export enum TableStatus { AVAILABLE, OCCUPIED, RESERVED }
export interface Table { id: string; number: string; status: TableStatus; }
export interface TableZone { id: string; name: string; }
export enum OrderStatus { PENDING, PREPARING, SERVED, COMPLETED, CANCELLED }
export enum OrderType { DINE_IN, TAKEAWAY, DELIVERY }
export interface Order { id: string; status: OrderStatus; type: OrderType; totalAmount: number; }
export interface OrderItem { id: string; productId: string; quantity: number; price: number; }
export interface OrderItemModifier { id: string; variantId: string; }
export interface Payment { id: string; amount: number; methodId: string; }
export interface PaymentMethod { id: string; name: string; }
export interface Promotion { id: string; name: string; discountPercentage: number; }
export interface Voucher { id: string; code: string; discountAmount: number; }
export interface OrderCreateRequest { type: OrderType; items: any[]; }
export interface CartItem { productId: string; quantity: number; }
export interface OrderFilters { status?: OrderStatus; date?: string; }
