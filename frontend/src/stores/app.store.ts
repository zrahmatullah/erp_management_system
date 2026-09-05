import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false);
  const theme = ref<'light' | 'dark'>('light');
  const locale = ref('id-ID');
  const breadcrumbs = ref<{label: string, path?: string}[]>([]);
  const currentBranch = ref<string | null>(null);

  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value;
  };

  const setTheme = (newTheme: 'light' | 'dark') => {
    theme.value = newTheme;
  };

  const setBreadcrumbs = (newBreadcrumbs: {label: string, path?: string}[]) => {
    breadcrumbs.value = newBreadcrumbs;
  };

  const setCurrentBranch = (branchId: string) => {
    currentBranch.value = branchId;
  };

  return {
    sidebarCollapsed, theme, locale, breadcrumbs, currentBranch,
    toggleSidebar, setTheme, setBreadcrumbs, setCurrentBranch
  };
});
