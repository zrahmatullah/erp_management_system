<template>
  <div 
    class="layout-main-wrapper flex h-screen w-screen overflow-hidden relative transition-colors duration-200"
    style="background-color: var(--bg-content); color: var(--text-primary);"
  >
    <TopProgressBar />
    <Sidebar />
    <div class="flex-1 flex flex-col min-w-0 h-full overflow-hidden">
      <AppHeader />
      <main 
        class="content-container flex-1 overflow-y-auto p-6 relative transition-colors"
        style="background-color: var(--bg-content);"
      >
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" :key="$route.fullPath" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import AppHeader from './AppHeader.vue'
import TopProgressBar from '@/components/common/TopProgressBar.vue'
</script>

<style scoped>
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.22s cubic-bezier(0.16, 1, 0.3, 1), transform 0.22s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
