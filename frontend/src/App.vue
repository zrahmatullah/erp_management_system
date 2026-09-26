<template>
  <n-config-provider :theme="appStore.theme === 'dark' ? darkTheme : null" :theme-overrides="themeOverrides">
    <n-notification-provider>
      <n-message-provider>
        <n-dialog-provider>
          <router-view />
          <ToastContainer />
          <GlobalDialogContainer />
        </n-dialog-provider>
      </n-message-provider>
    </n-notification-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NNotificationProvider, NMessageProvider, NDialogProvider, darkTheme, type GlobalThemeOverrides } from 'naive-ui'
import ToastContainer from '@/components/common/ToastContainer.vue'
import GlobalDialogContainer from '@/components/common/GlobalDialogContainer.vue'
import { useAppStore } from '@/stores/app.store'

const appStore = useAppStore()

const themeOverrides = computed<GlobalThemeOverrides>(() => {
  const isDark = appStore.theme === 'dark'
  return {
    common: {
      fontSize: '14.5px',
      fontSizeMedium: '14.5px',
      fontSizeLarge: '16px',
      heightMedium: '38px',
      heightLarge: '42px',
      borderRadius: '12px',
      primaryColor: isDark ? '#3B82F6' : '#2563EB',
      primaryColorHover: isDark ? '#60A5FA' : '#1D4ED8',
      primaryColorPressed: isDark ? '#1D4ED8' : '#1E40AF',
      bodyColor: isDark ? '#020617' : '#f8fafc',
      cardColor: isDark ? '#0f172a' : '#ffffff',
      modalColor: isDark ? '#0f172a' : '#ffffff',
      popoverColor: isDark ? '#1e293b' : '#ffffff',
      tableColor: isDark ? '#0f172a' : '#ffffff',
      textColorBase: isDark ? '#f8fafc' : '#0f172a',
      borderColor: isDark ? '#334155' : '#e2e8f0',
      dividerColor: isDark ? '#334155' : '#e2e8f0'
    }
  }
})
</script>
