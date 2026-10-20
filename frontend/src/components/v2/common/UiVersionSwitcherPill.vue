<template>
  <div class="relative inline-flex items-center">
    <button
      type="button"
      @click="uiVersionStore.toggleVersion"
      :title="tooltipText"
      class="group inline-flex items-center gap-1.5 px-2.5 py-1 text-[11px] font-semibold rounded-full transition-all duration-200 cursor-pointer border select-none"
      :class="uiVersionStore.isV2 
        ? 'bg-blue-50/80 hover:bg-blue-100 text-blue-700 border-blue-200/80 dark:bg-blue-950/40 dark:hover:bg-blue-900/50 dark:text-blue-300 dark:border-blue-800/60 shadow-xs' 
        : 'bg-amber-50 hover:bg-amber-100 text-amber-800 border-amber-300 dark:bg-amber-950/40 dark:text-amber-300 dark:border-amber-800 shadow-xs'"
    >
      <!-- Sparkle / Modern Icon for V2 -->
      <Sparkles v-if="uiVersionStore.isV2" class="w-3.5 h-3.5 text-blue-600 dark:text-blue-400 group-hover:rotate-12 transition-transform" />
      <!-- History / Classic Icon for V1 -->
      <History v-else class="w-3.5 h-3.5 text-amber-600 dark:text-amber-400 group-hover:-rotate-45 transition-transform" />

      <span class="tracking-tight font-medium">
        {{ uiVersionStore.isV2 ? 'UI V2 (Modern)' : 'UI V1 (Classic)' }}
      </span>

      <span 
        class="text-[9px] px-1 py-0.2 rounded font-mono uppercase tracking-wider"
        :class="uiVersionStore.isV2 
          ? 'bg-blue-200/70 text-blue-800 dark:bg-blue-900 dark:text-blue-200' 
          : 'bg-amber-200/70 text-amber-900 dark:bg-amber-900 dark:text-amber-200'"
      >
        {{ uiVersionStore.isV2 ? 'Switch V1' : 'Switch V2' }}
      </span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Sparkles, History } from 'lucide-vue-next'
import { useUiVersionStore } from '@/stores/uiVersion.store'

const uiVersionStore = useUiVersionStore()

const tooltipText = computed(() => {
  return uiVersionStore.isV2
    ? 'Klik untuk Instant Rollback ke Tampilan UI V1 (Classic)'
    : 'Klik untuk beralih ke Tampilan UI V2 (Modern & Minimalis)'
})
</script>

