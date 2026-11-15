<template>
  <div class="relative inline-flex items-center">
    <button
      type="button"
      @click="uiVersionStore.toggleVersion"
      :title="tooltipText"
      class="h-9 group inline-flex items-center gap-2 px-3 text-xs font-semibold rounded-xl border transition-all duration-200 cursor-pointer shadow-xs select-none"
      :class="uiVersionStore.isV2 
        ? 'bg-blue-50/90 hover:bg-blue-100/90 text-[#4361EE] dark:text-blue-300 border-blue-200/90 dark:border-blue-800/70 dark:bg-blue-950/40 dark:hover:bg-blue-900/50' 
        : 'bg-amber-50/90 hover:bg-amber-100/90 text-amber-800 dark:text-amber-300 border-amber-300/90 dark:border-amber-800/70 dark:bg-amber-950/40 dark:hover:bg-amber-900/50'"
    >
      <!-- Sparkle / Modern Icon for V2 -->
      <Sparkles v-if="uiVersionStore.isV2" class="w-3.5 h-3.5 text-[#4361EE] dark:text-blue-400 group-hover:rotate-12 transition-transform shrink-0" />
      <!-- History / Classic Icon for V1 -->
      <History v-else class="w-3.5 h-3.5 text-amber-600 dark:text-amber-400 group-hover:-rotate-45 transition-transform shrink-0" />

      <span class="tracking-tight font-bold">
        {{ uiVersionStore.isV2 ? 'UI V2 (Modern)' : 'UI V1 (Classic)' }}
      </span>

      <!-- Action pill badge inside -->
      <span 
        class="text-[10px] px-1.5 py-0.5 rounded-md font-mono font-bold uppercase tracking-wider shrink-0 transition-colors"
        :class="uiVersionStore.isV2 
          ? 'bg-blue-200/80 text-blue-900 dark:bg-blue-900/80 dark:text-blue-200 group-hover:bg-blue-300/80' 
          : 'bg-amber-200/80 text-amber-950 dark:bg-amber-900/80 dark:text-amber-200 group-hover:bg-amber-300/80'"
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
