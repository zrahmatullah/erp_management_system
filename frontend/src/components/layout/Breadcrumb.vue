<template>
  <n-breadcrumb>
    <n-breadcrumb-item href="/dashboard">
      Home
    </n-breadcrumb-item>
    <n-breadcrumb-item v-for="(crumb, idx) in appStore.breadcrumbs" :key="idx" :href="crumb.path">
      {{ crumb.label }}
    </n-breadcrumb-item>
  </n-breadcrumb>
</template>

<script setup lang="ts">
import { NBreadcrumb, NBreadcrumbItem } from 'naive-ui'
import { useAppStore } from '@/stores/app.store'
const appStore = useAppStore()
</script>
