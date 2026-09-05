import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User, Permission } from '@/types/auth';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('token'));
  const refreshTokenValue = ref<string | null>(localStorage.getItem('refreshToken'));
  const permissions = ref<Permission[]>([]);

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
  };

  const logout = () => {
    token.value = null;
    refreshTokenValue.value = null;
    user.value = null;
    permissions.value = [];
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
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
