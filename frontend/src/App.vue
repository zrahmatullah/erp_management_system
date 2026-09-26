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
import { NConfigProvider, NNotificationProvider, NMessageProvider, NDialogProvider, darkTheme, type GlobalThemeOverrides } from 'naive-ui'
import ToastContainer from '@/components/common/ToastContainer.vue'
import GlobalDialogContainer from '@/components/common/GlobalDialogContainer.vue'
import { useAppStore } from '@/stores/app.store'

const appStore = useAppStore()

const themeOverrides: GlobalThemeOverrides = {
  common: {
    fontSize: '14.5px',
    fontSizeMedium: '14.5px',
    fontSizeLarge: '16px',
    heightMedium: '38px',
    heightLarge: '42px',
    borderRadius: '12px',
    primaryColor: '#2563EB',
    primaryColorHover: '#1D4ED8',
    primaryColorPressed: '#1E40AF'
  }
}
</script>
