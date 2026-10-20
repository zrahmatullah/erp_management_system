export type UIVersion = 'v1' | 'v2';

export const UI_VERSION_STORAGE_KEY = 'cafe_erp_ui_version';

/**
 * Resolves the initial UI version using the following priority:
 * 1. Explicit user override stored in localStorage (`cafe_erp_ui_version`)
 * 2. Build-time / Environment variable (`VITE_UI_VERSION`)
 * 3. Default fallback to 'v2'
 */
export function getInitialUIVersion(): UIVersion {
  try {
    if (typeof localStorage !== 'undefined') {
      const stored = localStorage.getItem(UI_VERSION_STORAGE_KEY);
      if (stored === 'v1' || stored === 'v2') {
        return stored;
      }
    }
  } catch (err) {
    console.warn('[UIVersionConfig] localStorage is unavailable:', err);
  }

  const envVersion = import.meta.env.VITE_UI_VERSION;
  if (envVersion === 'v1' || envVersion === 'v2') {
    return envVersion;
  }

  return 'v2';
}

