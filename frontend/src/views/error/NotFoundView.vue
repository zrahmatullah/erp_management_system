<template>
  <div class="min-h-screen w-full flex flex-col justify-between bg-slate-50 dark:bg-slate-950 text-slate-800 dark:text-slate-100 transition-colors duration-200 select-none p-4 sm:p-8">
    <!-- Top Header -->
    <header class="w-full max-w-6xl mx-auto flex items-center justify-between py-2">
      <router-link to="/dashboard" class="flex items-center gap-2.5 group">
        <div class="w-10 h-10 rounded-2xl bg-blue-600 flex items-center justify-center text-white shadow-md shadow-blue-500/25 group-hover:scale-105 transition-transform">
          <Coffee class="w-5 h-5" />
        </div>
        <div>
          <span class="text-base font-black text-slate-900 dark:text-white tracking-tight">Cafe ERP</span>
          <span class="text-[10px] font-bold text-blue-600 dark:text-blue-400 uppercase tracking-wider block">Enterprise Suite</span>
        </div>
      </router-link>

      <!-- Theme Switcher Capsule -->
      <button 
        @click="appStore.toggleTheme"
        class="flex items-center gap-2 px-3 py-1.5 rounded-full border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 text-xs font-bold text-slate-700 dark:text-slate-200 hover:border-blue-500/40 shadow-xs transition-all cursor-pointer"
        :title="appStore.theme === 'dark' ? 'Ganti ke Mode Terang' : 'Ganti ke Mode Gelap'"
      >
        <Sun v-if="appStore.theme === 'dark'" class="w-4 h-4 text-amber-400" />
        <Moon v-else class="w-4 h-4 text-indigo-500" />
        <span class="hidden sm:inline">{{ appStore.theme === 'dark' ? 'Mode Gelap' : 'Mode Terang' }}</span>
      </button>
    </header>

    <!-- Main 404 Hero Container -->
    <main class="w-full max-w-xl mx-auto my-auto text-center px-4 py-8">
      <!-- Visual Illustration / Graphic -->
      <div class="relative w-36 h-36 sm:w-44 sm:h-44 mx-auto mb-6 flex items-center justify-center">
        <!-- Glowing background ring -->
        <div class="absolute inset-0 rounded-full bg-blue-500/10 dark:bg-blue-500/20 blur-2xl animate-pulse"></div>
        <div class="absolute -inset-4 rounded-full bg-indigo-500/10 dark:bg-indigo-500/15 blur-xl"></div>
        
        <!-- Big Card Container -->
        <div class="w-full h-full rounded-3xl bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800 shadow-xl flex flex-col items-center justify-center relative overflow-hidden">
          <div class="text-5xl sm:text-6xl font-black bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 bg-clip-text text-transparent tracking-tighter">
            404
          </div>
          <div class="flex items-center gap-1.5 mt-1 text-[11px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-widest">
            <Compass class="w-3.5 h-3.5 animate-spin" style="animation-duration: 12s;" />
            <span>Not Found</span>
          </div>
        </div>
      </div>

      <!-- Main Titles -->
      <h1 class="text-2xl sm:text-3xl font-black text-slate-900 dark:text-white tracking-tight mb-2">
        Halaman Tidak Ditemukan
      </h1>
      <p class="text-xs sm:text-sm text-slate-500 dark:text-slate-400 max-w-md mx-auto leading-relaxed mb-4">
        Maaf, tautan atau menu yang Anda akses tidak tersedia di sistem Cafe ERP, telah dipindahkan, atau akun Anda belum memiliki izin akses yang sesuai.
      </p>

      <!-- URL Badge -->
      <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-xl bg-slate-100 dark:bg-slate-900/80 border border-slate-200/80 dark:border-slate-800 text-xs font-mono text-slate-600 dark:text-slate-300 mb-8 max-w-full truncate">
        <span class="w-2 h-2 rounded-full bg-rose-500 shrink-0"></span>
        <span class="truncate">{{ $route.fullPath }}</span>
      </div>

      <!-- Action Buttons -->
      <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
        <button
          type="button"
          @click="goBack"
          class="w-full sm:w-auto px-5 py-2.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 hover:bg-slate-50 dark:hover:bg-slate-800/80 text-xs font-bold text-slate-700 dark:text-slate-200 shadow-xs transition-all flex items-center justify-center gap-2 cursor-pointer"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>Kembali ke Sebelumnya</span>
        </button>

        <router-link
          to="/dashboard"
          class="w-full sm:w-auto px-6 py-2.5 rounded-xl bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold shadow-md shadow-blue-600/25 transition-all flex items-center justify-center gap-2 cursor-pointer"
        >
          <LayoutDashboard class="w-4 h-4" />
          <span>Ke Dashboard Utama</span>
        </router-link>
      </div>
    </main>

    <!-- Footer Copyright -->
    <footer class="w-full max-w-6xl mx-auto py-3 text-center text-xs text-slate-400 dark:text-slate-600 font-medium">
      &copy; {{ new Date().getFullYear() }} Cafe Monitoring & ERP System. Hak Cipta Dilindungi.
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app.store'
import {
  Coffee,
  Sun,
  Moon,
  Compass,
  ArrowLeft,
  LayoutDashboard
} from 'lucide-vue-next'

const router = useRouter()
const appStore = useAppStore()

const goBack = () => {
  if (window.history.length > 2) {
    router.back()
  } else {
    router.push('/dashboard')
  }
}
</script>

