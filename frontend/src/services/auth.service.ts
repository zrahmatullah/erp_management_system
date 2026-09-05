import api from '@/plugins/axios';
import type { LoginResponse } from '@/types/auth';

export const login = async (data: any) => api.post<LoginResponse>('/auth/login', data);
export const logout = async () => api.post('/auth/logout');
export const refreshToken = async (data: any) => api.post('/auth/refresh', data);
export const forgotPassword = async (data: any) => api.post('/auth/forgot-password', data);
export const resetPassword = async (data: any) => api.post('/auth/reset-password', data);
export const getProfile = async () => api.get('/auth/profile');
