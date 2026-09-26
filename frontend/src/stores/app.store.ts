import { defineStore } from 'pinia';
import { ref } from 'vue';
import { ref, computed } from 'vue';

export type ThemeMode = 'light' | 'dark' | 'system';

const STORAGE_KEY = 'cafe_erp_theme';

// Safe localStorage access handling private browsing & quota exceptions
function safeGetStorage(key: string): string | null {
  try {
    if (typeof localStorage !== 'undefined') {
      return localStorage.getItem(key);
    }
  } catch (err) {
    console.warn('[ThemeStore] localStorage is not accessible (e.g. private browsing):', err);
  }
  return null;
}

function safeSetStorage(key: string, value: string): void {
  try {
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(key, value);
    }
  } catch (err) {
    console.warn('[ThemeStore] Unable to save to localStorage:', err);
  }
}

// Media query helper for system preference
function getSystemTheme(): 'light' | 'dark' {
  if (typeof window !== 'undefined' && typeof window.matchMedia === 'function') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }
  return 'dark'; // Fallback default
}

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false);

  // Read saved theme from localStorage, default to dark
  const savedTheme = typeof localStorage !== 'undefined' ? localStorage.getItem('cafe_erp_theme') as 'light' | 'dark' | null : null;
  const initialTheme: 'light' | 'dark' = savedTheme ? savedTheme : 'dark';
  // Read saved theme from storage, default to 'system'
  const rawSaved = safeGetStorage(STORAGE_KEY);
  const initialMode: ThemeMode = (rawSaved === 'light' || rawSaved === 'dark' || rawSaved === 'system')
    ? rawSaved
    : 'system';

  const theme = ref<'light' | 'dark'>(initialTheme);
  const themeMode = ref<ThemeMode>(initialMode);
  const systemPreference = ref<'light' | 'dark'>(getSystemTheme());

  // Resolved active theme ('light' or 'dark') applied to DOM
  const resolvedTheme = computed<'light' | 'dark'>(() => {
    if (themeMode.value === 'system') {
      return systemPreference.value;
    }
    return themeMode.value;
  });

  // For backward compatibility across all existing components expecting `theme`
  const theme = computed<'light' | 'dark'>({
    get: () => resolvedTheme.value,
    set: (val: 'light' | 'dark') => setTheme(val)
  });

  const locale = ref('id-ID');
  const breadcrumbs = ref<{label: string, path?: string}[]>([]);
  const breadcrumbs = ref<{ label: string; path?: string }[]>([]);
  const currentBranch = ref<string | null>(null);

  const applyTheme = (targetTheme: 'light' | 'dark') => {
    theme.value = targetTheme;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('cafe_erp_theme', targetTheme);
    }
  // Apply resolved theme class & color-scheme to document root
  const applyDOMTheme = (targetTheme: 'light' | 'dark') => {
    if (typeof document !== 'undefined') {
      const root = document.documentElement;
      if (targetTheme === 'dark') {
        document.documentElement.classList.add('dark');
        root.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
        root.classList.remove('dark');
      }
      root.style.colorScheme = targetTheme;
    }
  };

  // Sync initially
  applyTheme(initialTheme);
  // Set new theme mode ('light', 'dark', or 'system')
  const setTheme = (newMode: ThemeMode) => {
    themeMode.value = newMode;
    safeSetStorage(STORAGE_KEY, newMode);
    applyDOMTheme(resolvedTheme.value);
  };

  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value;
  // Simple toggle between light and dark
  const toggleTheme = () => {
    const nextTheme = resolvedTheme.value === 'dark' ? 'light' : 'dark';
    setTheme(nextTheme);
  };

  const setTheme = (newTheme: 'light' | 'dark') => {
    applyTheme(newTheme);
  // 3-way cyclic toggle (light -> dark -> system -> light)
  const cycleTheme = () => {
    if (themeMode.value === 'light') {
      setTheme('dark');
    } else if (themeMode.value === 'dark') {
      setTheme('system');
    } else {
      setTheme('light');
    }
  };

  const toggleTheme = () => {
    applyTheme(theme.value === 'dark' ? 'light' : 'dark');
  // Real-time listener for OS preference changes
  if (typeof window !== 'undefined' && typeof window.matchMedia === 'function') {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleSystemChange = (e: MediaQueryListEvent) => {
      systemPreference.value = e.matches ? 'dark' : 'light';
      if (themeMode.value === 'system') {
        applyDOMTheme(systemPreference.value);
      }
    };

    if (typeof mediaQuery.addEventListener === 'function') {
      mediaQuery.addEventListener('change', handleSystemChange);
    } else if (typeof (mediaQuery as any).addListener === 'function') {
      // Legacy Safari / older browsers
      (mediaQuery as any).addListener(handleSystemChange);
    }
  }

  // Initial sync
  applyDOMTheme(resolvedTheme.value);

  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value;
  };

  const setBreadcrumbs = (newBreadcrumbs: {label: string, path?: string}[]) => {
  const setBreadcrumbs = (newBreadcrumbs: { label: string; path?: string }[]) => {
    breadcrumbs.value = newBreadcrumbs;
  };

  const setCurrentBranch = (branchId: string) => {
    currentBranch.value = branchId;
  };

  return {
    sidebarCollapsed, theme, locale, breadcrumbs, currentBranch,
    toggleSidebar, setTheme, toggleTheme, setBreadcrumbs, setCurrentBranch
    sidebarCollapsed,
    themeMode,
    resolvedTheme,
    theme,
    locale,
    breadcrumbs,
    currentBranch,
    toggleSidebar,
    setTheme,
    toggleTheme,
    cycleTheme,
    setBreadcrumbs,
    setCurrentBranch
  };
});

