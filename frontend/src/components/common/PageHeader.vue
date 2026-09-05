<template>
  <div class="mb-6">
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-4">
        <n-button v-if="showBack" circle @click="$router.back()">
          ←
        </n-button>
        <div>
          <h1 class="text-2xl font-bold m-0">{{ title }}</h1>
          <p v-if="subtitle" class="text-gray-500 m-0 text-sm">{{ subtitle }}</p>
        </div>
      </div>
      <div>
        <slot></slot>
      </div>
    </div>
    <n-divider class="!m-0" />
  </div>
</template>

<script setup lang="ts">
import { NButton, NDivider } from 'naive-ui'

defineProps<{
  title: string
  subtitle?: string
  showBack?: boolean
}>()
</script>
