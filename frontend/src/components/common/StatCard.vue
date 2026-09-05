<template>
  <n-card class="rounded-lg shadow-sm">
    <div class="flex justify-between">
      <div>
        <p class="text-gray-500 text-sm mb-1">{{ title }}</p>
        <h3 class="text-2xl font-bold">{{ value }}</h3>
        <p class="text-xs mt-2" :class="trend === 'up' ? 'text-green-500' : (trend === 'down' ? 'text-red-500' : 'text-gray-400')">
          <span v-if="trend === 'up'">↑</span>
          <span v-else-if="trend === 'down'">↓</span>
          <span v-else>-</span>
          {{ trendValue }} vs kemarin
        </p>
      </div>
      <div :class="`w-10 h-10 rounded-full flex items-center justify-center bg-${color}-100 text-${color}-600`">
        <!-- icon slot -->
      </div>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { NCard } from 'naive-ui'

defineProps<{
  title: string
  value: string | number
  trend: 'up' | 'down' | 'neutral'
  trendValue: string
  color: string
}>()
</script>
