import api from '@/plugins/axios';

export const getRoles = async () => api.get('/roles');
export const createRole = async (data: any) => api.post('/roles', data);
export const updateRole = async (id: string, data: any) => api.put(`/roles/${id}`, data);
export const deleteRole = async (id: string) => api.delete(`/roles/${id}`);
export const getPermissions = async () => api.get('/permissions');
export const assignPermissions = async (id: string, data: any) => api.post(`/roles/${id}/permissions`, data);
