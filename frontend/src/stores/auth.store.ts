import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User, Permission } from '@/types/auth';
import api from '@/plugins/axios';

export const useAuthStore = defineStore('auth', () => {
  const savedUser = localStorage.getItem('user');
  const user = ref<User | null>(savedUser ? JSON.parse(savedUser) : null);
  const token = ref<string | null>(localStorage.getItem('token'));
  const refreshTokenValue = ref<string | null>(localStorage.getItem('refreshToken'));
  const savedPerms = localStorage.getItem('permissions');
  const permissions = ref<Permission[]>(savedPerms ? JSON.parse(savedPerms) : []);

  const isAuthenticated = computed(() => !!token.value);
  const fullName = computed(() => user.value ? `${user.value.firstName} ${user.value.lastName}` : '');
  const currentRole = computed(() => user.value?.roleId || '');

  const hasPermission = (module: string, action: string) => {
    return permissions.value.some(p => p.module === module && p.action === action);
  };

  const login = (newToken: string, newRefreshToken: string, userData: User, perms: Permission[]) => {
    token.value = newToken;
    refreshTokenValue.value = newRefreshToken;
    user.value = userData;
    permissions.value = perms;
    localStorage.setItem('token', newToken);
    localStorage.setItem('refreshToken', newRefreshToken);
    localStorage.setItem('user', JSON.stringify(userData));
    localStorage.setItem('permissions', JSON.stringify(perms));
  };

  const logout = () => {
    token.value = null;
    refreshTokenValue.value = null;
    user.value = null;
    permissions.value = [];
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('user');
    localStorage.removeItem('permissions');
  };

  const setUser = (userData: User) => {
    user.value = userData;
  };

  const fetchProfile = async () => {
    if (!token.value) return;
    try {
      const res = await api.get('/auth/profile');
      const data = res.data?.data;
      if (data) {
        if (data.user) {
          user.value = data.user;
          localStorage.setItem('user', JSON.stringify(data.user));
        }
        if (data.permissions) {
          permissions.value = data.permissions;
          localStorage.setItem('permissions', JSON.stringify(data.permissions));
        }
      }
    } catch (err) {
      console.error('Failed to sync profile:', err);
    }
  };

  return {
    user, token, refreshTokenValue, permissions, isAuthenticated,
    fullName, currentRole, hasPermission, login, logout, setUser, fetchProfile
  };
});
