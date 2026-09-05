import api from '@/plugins/axios';

export const getUsers = async (params?: any) => api.get('/users', { params });
export const getUserById = async (id: string) => api.get(`/users/${id}`);
export const createUser = async (data: any) => api.post('/users', data);
export const updateUser = async (id: string, data: any) => api.put(`/users/${id}`, data);
export const deleteUser = async (id: string) => api.delete(`/users/${id}`);
export const assignRole = async (id: string, data: any) => api.post(`/users/${id}/role`, data);
