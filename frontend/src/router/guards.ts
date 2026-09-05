import { useAuthStore } from '@/stores/auth.store';
import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router';

export const requireAuth = (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  const authStore = useAuthStore();
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login');
  } else {
    next();
  }
};

export const requireRole = (roles: string[]) => {
  return (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
    const authStore = useAuthStore();
    if (roles.includes(authStore.currentRole)) {
      next();
    } else {
      next('/forbidden');
    }
  };
};
