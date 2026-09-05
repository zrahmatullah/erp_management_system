import api from '@/plugins/axios';

export const createOrder = async (data: any) => api.post('/orders', data);
export const getOrders = async (params?: any) => api.get('/orders', { params });
export const getOrderById = async (id: string) => api.get(`/orders/${id}`);
export const updateOrderStatus = async (id: string, status: any) => api.put(`/orders/${id}/status`, { status });
export const cancelOrder = async (id: string) => api.post(`/orders/${id}/cancel`);
export const voidOrder = async (id: string) => api.post(`/orders/${id}/void`);
export const getKitchenQueue = async () => api.get('/orders/kitchen-queue');
