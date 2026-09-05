<template>
  <div class="flex h-screen w-screen overflow-hidden bg-[#F8FAFC]">
    <Sidebar />
    <div class="flex-1 flex flex-col min-w-0 h-full overflow-hidden">
      <AppHeader />
      <main class="flex-1 overflow-y-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import AppHeader from './AppHeader.vue'
</script>
