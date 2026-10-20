<template>
  <div class="page-header-v2 mb-6 sm:mb-8">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <!-- Title & Description -->
      <div class="min-w-0">
        <div class="flex items-center gap-2.5 flex-wrap">
          <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-50">
            {{ title }}
          </h1>
          <span 
            v-if="badge"
            class="text-[11px] font-semibold px-2.5 py-0.5 rounded-full bg-blue-50 text-blue-700 border border-blue-200/80 dark:bg-blue-950/40 dark:text-blue-300 dark:border-blue-800/60"
          >
            {{ badge }}
          </span>
        </div>
        <p v-if="description" class="text-xs sm:text-sm text-slate-500 dark:text-slate-400 mt-1 max-w-3xl leading-relaxed">
          {{ description }}
        </p>
      </div>

      <!-- Action Buttons Slot -->
      <div v-if="$slots.actions" class="flex items-center gap-2.5 shrink-0 flex-wrap">
        <slot name="actions" />
      </div>
    </div>

    <!-- Optional Bottom Tabs Slot -->
    <div v-if="$slots.tabs" class="mt-4 pt-1 border-b border-slate-200/70 dark:border-slate-800">
      <slot name="tabs" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  badge?: string
}>()
</script>

