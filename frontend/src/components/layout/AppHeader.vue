<template>
  <header class="h-16 bg-white dark:bg-slate-900 border-b border-slate-200/80 dark:border-slate-800 flex items-center justify-between px-6 z-10 transition-colors">
    <!-- Left: Breadcrumb / Page context -->
    <div class="flex items-center gap-3">
      <Breadcrumb />
    </div>

    <!-- Right: Search, Notifications, Theme Toggle, User Profile -->
    <div class="flex items-center gap-3 sm:gap-4">
      <!-- Search input -->
      <div class="relative hidden sm:block w-64 lg:w-72">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
          <Search class="w-4 h-4" />
        </span>
        <input
          type="text"
          placeholder="Cari transaksi, laporan..."
          class="w-full pl-9 pr-4 py-1.5 text-xs bg-slate-100/80 hover:bg-slate-100 dark:bg-slate-800 dark:hover:bg-slate-700/80 focus:bg-white dark:focus:bg-slate-800 text-slate-700 dark:text-slate-200 rounded-lg border border-slate-200 dark:border-slate-700 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all placeholder:text-slate-400 dark:placeholder:text-slate-500"
        />
      </div>

      <!-- Theme Switcher (Light / Dark) -->
      <button 
        @click="appStore.toggleTheme" 
        class="w-8 h-8 rounded-full flex items-center justify-center text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
        :title="appStore.theme === 'dark' ? 'Ganti ke Mode Terang (Light Mode)' : 'Ganti ke Mode Gelap (Dark Mode)'"
      >
        <Sun v-if="appStore.theme === 'dark'" class="w-4 h-4 text-amber-400 hover:rotate-45 transition-transform" />
        <Moon v-else class="w-4 h-4 text-slate-600 hover:-rotate-12 transition-transform" />
      </button>

      <!-- Help Button -->
      <button class="w-8 h-8 rounded-full flex items-center justify-center text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer" title="Bantuan & Dukungan">
        <HelpCircle class="w-4 h-4" />
      </button>

      <!-- Notification Bell -->
      <div class="relative cursor-pointer">
        <div class="w-8 h-8 rounded-full flex items-center justify-center text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors">
          <Bell class="w-4 h-4" />
        </div>
        <span class="absolute top-0 right-0 w-4 h-4 bg-red-500 text-white text-[10px] font-bold rounded-full flex items-center justify-center shadow-sm">
          5
        </span>
      </div>

      <div class="h-6 w-px bg-slate-200 dark:bg-slate-800 mx-1"></div>

      <!-- User Profile Dropdown -->
      <n-dropdown :options="userOptions" @select="handleUserMenu">
        <div class="flex items-center gap-2.5 cursor-pointer py-1 px-2 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-800/80 transition-colors">
          <div class="w-9 h-9 rounded-full bg-gradient-to-tr from-blue-600 to-indigo-600 text-white font-semibold flex items-center justify-center shadow-sm text-xs">
            {{ userInitials }}
          </div>
          <div class="text-left hidden md:block">
            <div class="text-xs font-semibold text-slate-800 dark:text-slate-200 leading-tight">{{ userDisplayName }}</div>
            <div class="text-[10px] text-slate-500 dark:text-slate-400 font-medium capitalize">{{ userRoleName }}</div>
          </div>
          <ChevronDown class="w-3.5 h-3.5 text-slate-400 ml-1" />
        </div>
      </n-dropdown>
    </div>

    <!-- Logout Loading Overlay -->
    <Teleport to="body">
      <transition name="fade">
        <div v-if="isLoggingOut" class="fixed inset-0 z-99999 flex flex-col items-center justify-center bg-slate-950/40 backdrop-blur-md">
          <div class="bg-white rounded-3xl p-8 max-w-xs w-full shadow-2xl shadow-slate-950/30 flex flex-col items-center text-center animate-in zoom-in-95 duration-200">
            <div class="relative w-16 h-16 mb-4 flex items-center justify-center">
              <div class="absolute inset-0 rounded-full border-3 border-rose-100 border-t-rose-600 animate-spin"></div>
              <div class="w-10 h-10 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center">
                <LogOut class="w-5 h-5" />
              </div>
            </div>
            <h3 class="text-base font-bold text-slate-900 tracking-tight">Mengakhiri Sesi...</h3>
            <p class="text-xs text-slate-500 mt-1">Menutup akses aman dan mengalihkan ke halaman login</p>
          </div>
        </div>
      </transition>
    </Teleport>
  </header>
</template>

<script setup lang="ts">
import { h, ref, computed } from 'vue'
import { NDropdown } from 'naive-ui'
import Breadcrumb from './Breadcrumb.vue'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { useNotificationStore } from '@/stores/notification.store'
import { useRouter } from 'vue-router'
import {
  Search,
  HelpCircle,
  Bell,
  ChevronDown,
  User,
  Database,
  Settings,
  LogOut,
  Sun,
  Moon
} from 'lucide-vue-next'

const authStore = useAuthStore()
const appStore = useAppStore()
const notifyStore = useNotificationStore()
const router = useRouter()
const isLoggingOut = ref(false)

const userDisplayName = computed(() => authStore.fullName?.trim() || authStore.user?.firstName || 'Super Administrator')
const userRoleName = computed(() => authStore.currentRole || 'Administrator')
const userInitials = computed(() => {
  const name = userDisplayName.value.trim()
  const parts = name.split(' ')
  if (parts.length >= 2) return `${parts[0][0]}${parts[1][0]}`.toUpperCase()
  return name.slice(0, 2).toUpperCase() || 'SA'
})

const renderIcon = (icon: any) => {
  return () => h(icon, { class: 'w-4 h-4' })
}

const userOptions = [
  { label: 'Profil Pengguna', key: 'profile', icon: renderIcon(User) },
  { label: 'Data Master Super Admin', key: 'master', icon: renderIcon(Database) },
  { label: 'Pengaturan Sistem', key: 'settings', icon: renderIcon(Settings) },
  { type: 'divider', key: 'd1' },
  { label: 'Logout', key: 'logout', icon: renderIcon(LogOut) }
]

const handleUserMenu = (key: string) => {
  if (key === 'logout') {
    isLoggingOut.value = true
    setTimeout(() => {
      authStore.logout()
      notifyStore.info('Anda telah berhasil keluar dari sistem.', 'Sesi Berakhir')
      router.push('/login')
      setTimeout(() => {
        isLoggingOut.value = false
      }, 300)
    }, 500)
  } else if (key === 'master') {
    router.push('/settings/master')
  } else if (key === 'settings') {
    router.push('/settings/roles')
  }
}
</script>
