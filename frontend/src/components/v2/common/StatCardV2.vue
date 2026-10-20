<template>
  <div 
    class="stat-card-v2 p-5 rounded-2xl border transition-all duration-200 hover:shadow-xs group relative overflow-hidden"
    style="background-color: var(--bg-primary); border-color: var(--border-color);"
  >
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0 flex-1">
        <p class="text-xs font-semibold text-slate-500 dark:text-slate-400 tracking-wide uppercase truncate">
          {{ title }}
        </p>
        <div class="mt-2 flex items-baseline gap-2">
          <h3 class="text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-50 truncate" :title="String(value)">
            {{ value }}
          </h3>
          <span v-if="unit" class="text-xs font-medium text-slate-400">
            {{ unit }}
          </span>
        </div>

        <!-- Trend / Context Pill -->
        <div v-if="trendValue || subtitle" class="mt-3 flex items-center gap-2 text-xs">
          <span 
            v-if="trendValue"
            class="inline-flex items-center gap-1 font-semibold px-2 py-0.5 rounded-full text-[11px]"
            :class="trendClass"
          >
            <TrendingUp v-if="trend === 'up'" class="w-3 h-3" />
            <TrendingDown v-else-if="trend === 'down'" class="w-3 h-3" />
            <Minus v-else class="w-3 h-3" />
            {{ trendValue }}
          </span>
          <span v-if="subtitle" class="text-slate-400 truncate text-[11px]">
            {{ subtitle }}
          </span>
        </div>
      </div>

      <!-- Icon Container -->
      <div 
        class="w-11 h-11 rounded-xl flex items-center justify-center shrink-0 transition-transform group-hover:scale-105"
        :class="iconBgClass"
      >
        <slot name="icon">
          <component :is="icon" v-if="icon" class="w-5 h-5" :class="iconColorClass" />
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, TrendingDown, Minus } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  title: string
  value: string | number
  unit?: string
  trend?: 'up' | 'down' | 'neutral'
  trendValue?: string
  subtitle?: string
  icon?: any
  color?: 'blue' | 'emerald' | 'amber' | 'rose' | 'indigo' | 'purple' | string
}>(), {
  trend: 'neutral',
  color: 'blue'
})

const trendClass = computed(() => {
  if (props.trend === 'up') return 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/50 dark:text-emerald-300'
  if (props.trend === 'down') return 'bg-rose-50 text-rose-700 dark:bg-rose-950/50 dark:text-rose-300'
  return 'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-300'
})

const iconBgClass = computed(() => {
  switch (props.color) {
    case 'emerald':
    case 'green':
      return 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400'
    case 'amber':
    case 'yellow':
      return 'bg-amber-50 text-amber-600 dark:bg-amber-950/60 dark:text-amber-400'
    case 'rose':
    case 'red':
      return 'bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400'
    case 'purple':
      return 'bg-purple-50 text-purple-600 dark:bg-purple-950/60 dark:text-purple-400'
    case 'indigo':
      return 'bg-indigo-50 text-indigo-600 dark:bg-indigo-950/60 dark:text-indigo-400'
    case 'blue':
    default:
      return 'bg-blue-50 text-blue-600 dark:bg-blue-950/60 dark:text-blue-400'
  }
})

const iconColorClass = computed(() => {
  switch (props.color) {
    case 'emerald':
    case 'green':
      return 'text-emerald-600 dark:text-emerald-400'
    case 'amber':
    case 'yellow':
      return 'text-amber-600 dark:text-amber-400'
    case 'rose':
    case 'red':
      return 'text-rose-600 dark:text-rose-400'
    case 'purple':
      return 'text-purple-600 dark:text-purple-400'
    case 'indigo':
      return 'text-indigo-600 dark:text-indigo-400'
    case 'blue':
    default:
      return 'text-blue-600 dark:text-blue-400'
  }
})
</script>

