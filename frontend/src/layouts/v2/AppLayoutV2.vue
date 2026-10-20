<template>
  <div 
    class="layout-v2-wrapper flex h-screen w-screen overflow-hidden relative transition-colors duration-200"
    style="background-color: var(--bg-content); color: var(--text-primary);"
  >
    <TopProgressBar />
    <SidebarV2 />
    <div class="flex-1 flex flex-col min-w-0 h-full overflow-hidden">
      <AppHeaderV2 />
      <main 
        class="content-v2-container flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8 relative transition-colors"
        style="background-color: var(--bg-content);"
      >
        <div class="w-full max-w-[1600px] mx-auto">
          <router-view v-slot="{ Component }">
            <transition name="page-fade-v2" mode="out-in">
              <component :is="Component" :key="$route.fullPath" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import SidebarV2 from '@/components/v2/layout/SidebarV2.vue'
import AppHeaderV2 from '@/components/v2/layout/AppHeaderV2.vue'
import TopProgressBar from '@/components/common/TopProgressBar.vue'
</script>

<style scoped>
.page-fade-v2-enter-active,
.page-fade-v2-leave-active {
  transition: opacity 0.2s cubic-bezier(0.16, 1, 0.3, 1), transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-fade-v2-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.page-fade-v2-leave-to {
  opacity: 0;
  transform: translateY(-2px);
}
</style>

