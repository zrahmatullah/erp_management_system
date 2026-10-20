import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { UIVersion, UI_VERSION_STORAGE_KEY, getInitialUIVersion } from '@/config/uiVersion';

export const useUiVersionStore = defineStore('uiVersion', () => {
  const currentVersion = ref<UIVersion>(getInitialUIVersion());

  const isV2 = computed(() => currentVersion.value === 'v2');
  const isV1 = computed(() => currentVersion.value === 'v1');

  const applyDOMVersion = (version: UIVersion) => {
    if (typeof document !== 'undefined') {
      document.documentElement.setAttribute('data-ui-version', version);
      if (version === 'v2') {
        document.documentElement.classList.add('ui-v2');
        document.documentElement.classList.remove('ui-v1');
      } else {
        document.documentElement.classList.add('ui-v1');
        document.documentElement.classList.remove('ui-v2');
      }
    }
  };

  const setVersion = (newVersion: UIVersion) => {
    currentVersion.value = newVersion;
    try {
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem(UI_VERSION_STORAGE_KEY, newVersion);
      }
    } catch (err) {
      console.warn('[UIVersionStore] Failed to write to localStorage:', err);
    }
    applyDOMVersion(newVersion);
  };

  const toggleVersion = () => {
    setVersion(currentVersion.value === 'v2' ? 'v1' : 'v2');
  };

  const resetToEnvDefault = () => {
    try {
      if (typeof localStorage !== 'undefined') {
        localStorage.removeItem(UI_VERSION_STORAGE_KEY);
      }
    } catch (err) {
      console.warn('[UIVersionStore] Failed to clear localStorage:', err);
    }
    const envVer = (import.meta.env.VITE_UI_VERSION === 'v1' ? 'v1' : 'v2') as UIVersion;
    currentVersion.value = envVer;
    applyDOMVersion(envVer);
  };

  // Initial DOM application
  applyDOMVersion(currentVersion.value);

  return {
    currentVersion,
    isV2,
    isV1,
    setVersion,
    toggleVersion,
    resetToEnvDefault
  };
});

