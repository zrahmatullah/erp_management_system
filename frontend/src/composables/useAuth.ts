import { ref, computed } from 'vue';
import type { User } from '@/types/auth';

const user = ref<User | null>(null);
const token = ref<string | null>(localStorage.getItem('token'));

export const useAuth = () => {
  const isAuthenticated = computed(() => !!token.value);
  
  const login = (newToken: string, userData: User) => {
    token.value = newToken;
    user.value = userData;
    localStorage.setItem('token', newToken);
  };
  
  const logout = () => {
    token.value = null;
    user.value = null;
    localStorage.removeItem('token');
  };

  return { user, isAuthenticated, login, logout };
};
