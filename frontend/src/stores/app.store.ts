import { defineStore } from 'pinia';
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

  // Read saved theme from storage, default to 'system'
  const rawSaved = safeGetStorage(STORAGE_KEY);
  const initialMode: ThemeMode = (rawSaved === 'light' || rawSaved === 'dark' || rawSaved === 'system')
    ? rawSaved
    : 'system';

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
  const breadcrumbs = ref<{ label: string; path?: string }[]>([]);
  const currentBranch = ref<string | null>(null);

  // Apply resolved theme class, data-theme attribute & color-scheme to document root
  const applyDOMTheme = (targetTheme: 'light' | 'dark') => {
    if (typeof document !== 'undefined') {
      const root = document.documentElement;
      if (targetTheme === 'dark') {
        root.classList.add('dark');
      } else {
        root.classList.remove('dark');
      }
      root.setAttribute('data-theme', targetTheme);
      root.style.colorScheme = targetTheme;
    }
  };

  // Set new theme mode ('light', 'dark', or 'system')
  const setTheme = (newMode: ThemeMode) => {
    themeMode.value = newMode;
    safeSetStorage(STORAGE_KEY, newMode);
    applyDOMTheme(resolvedTheme.value);
  };

  // Simple toggle between light and dark
  const toggleTheme = () => {
    const nextTheme = resolvedTheme.value === 'dark' ? 'light' : 'dark';
    setTheme(nextTheme);
  };

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

  const setBreadcrumbs = (newBreadcrumbs: { label: string; path?: string }[]) => {
    breadcrumbs.value = newBreadcrumbs;
  };

  const setCurrentBranch = (branchId: string) => {
    currentBranch.value = branchId;
  };

  return {
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
