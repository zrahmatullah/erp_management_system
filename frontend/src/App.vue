<template>
  <n-config-provider>
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
import { NConfigProvider, NNotificationProvider, NMessageProvider, NDialogProvider } from 'naive-ui'
import ToastContainer from '@/components/common/ToastContainer.vue'
import GlobalDialogContainer from '@/components/common/GlobalDialogContainer.vue'
</script>
