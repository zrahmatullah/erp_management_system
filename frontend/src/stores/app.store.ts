import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false);

  // Read saved theme from localStorage, default to dark
  const savedTheme = typeof localStorage !== 'undefined' ? localStorage.getItem('cafe_erp_theme') as 'light' | 'dark' | null : null;
  const initialTheme: 'light' | 'dark' = savedTheme ? savedTheme : 'dark';

  const theme = ref<'light' | 'dark'>(initialTheme);
  const locale = ref('id-ID');
  const breadcrumbs = ref<{label: string, path?: string}[]>([]);
  const currentBranch = ref<string | null>(null);

  const applyTheme = (targetTheme: 'light' | 'dark') => {
    theme.value = targetTheme;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('cafe_erp_theme', targetTheme);
    }
    if (typeof document !== 'undefined') {
      if (targetTheme === 'dark') {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    }
  };

  // Sync initially
  applyTheme(initialTheme);

  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value;
  };

  const setTheme = (newTheme: 'light' | 'dark') => {
    applyTheme(newTheme);
  };

  const toggleTheme = () => {
    applyTheme(theme.value === 'dark' ? 'light' : 'dark');
  };

  const setBreadcrumbs = (newBreadcrumbs: {label: string, path?: string}[]) => {
    breadcrumbs.value = newBreadcrumbs;
  };

  const setCurrentBranch = (branchId: string) => {
    currentBranch.value = branchId;
  };

  return {
    sidebarCollapsed, theme, locale, breadcrumbs, currentBranch,
    toggleSidebar, setTheme, toggleTheme, setBreadcrumbs, setCurrentBranch
  };
});

