import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User, Permission } from '@/types/auth';

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
    // API call placeholder if needed
  };

  return {
    user, token, refreshTokenValue, permissions, isAuthenticated,
    fullName, currentRole, hasPermission, login, logout, setUser, fetchProfile
  };
});
