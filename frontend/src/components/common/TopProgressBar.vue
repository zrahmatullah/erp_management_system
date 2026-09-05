<template>
  <div
    v-if="loadingStore.isLoading || loadingStore.progress > 0"
    class="fixed top-0 left-0 right-0 z-9999 pointer-events-none"
  >
    <!-- Top Progress Line -->
    <div
      class="h-[3px] bg-gradient-to-r from-blue-600 via-sky-400 to-indigo-600 shadow-sm shadow-blue-500/50 transition-all duration-200 ease-out"
      :style="{
        width: `${loadingStore.progress}%`,
        opacity: loadingStore.isLoading ? 1 : 0
      }"
    >
      <!-- Glowing Tip -->
      <div class="absolute right-0 top-0 w-8 h-[3px] bg-white/70 blur-xs"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLoadingStore } from '@/stores/loading.store'

const loadingStore = useLoadingStore()
</script>
