<template>
  <header 
    class="h-16 border-b flex items-center justify-between px-4 sm:px-6 z-20 shrink-0 sticky top-0 transition-colors backdrop-blur-md"
    style="background-color: var(--bg-topbar); border-color: var(--border-color);"
  >
    <!-- Left: Sidebar Toggle (Mobile / Desktop) & Modern Breadcrumb -->
    <div class="flex items-center gap-3">
      <!-- Precision h-9 w-9 Sidebar Toggle Button -->
      <button
        type="button"
        @click="appStore.toggleSidebar"
        class="h-9 w-9 rounded-xl flex items-center justify-center transition-all duration-150 cursor-pointer border text-slate-600 dark:text-slate-300 hover:text-[#4361EE] dark:hover:text-blue-400 bg-white hover:bg-slate-50 dark:bg-slate-800/90 dark:hover:bg-slate-700/90 border-slate-200/80 dark:border-slate-700/80 shadow-xs hover:border-blue-300 dark:hover:border-blue-600 hover:shadow-sm"
        :title="appStore.sidebarCollapsed ? 'Perluas Sidebar (Ctrl+B)' : 'Ciutkan Sidebar (Ctrl+B)'"
        aria-label="Toggle Sidebar"
      >
        <Menu v-if="appStore.sidebarCollapsed" class="w-4 h-4 shrink-0" />
        <PanelLeftClose v-else class="w-4 h-4 shrink-0" />
      </button>

      <div class="h-5 w-px bg-slate-200 dark:bg-slate-700/80 mx-0.5 hidden sm:block"></div>

      <BreadcrumbV2 class="hidden md:flex" />
    </div>

    <!-- Right: Quick Search, UI Version Switcher, Theme Toggle, Notifications, Profile -->
    <div class="flex items-center gap-2 sm:gap-2.5">
      <!-- Minimalist Command Search Bar with unified h-9 height -->
      <div class="relative hidden lg:flex items-center w-64 xl:w-72">
        <span class="absolute left-3 flex items-center pointer-events-none text-slate-400 dark:text-slate-500">
          <Search class="w-3.5 h-3.5" />
        </span>
        <input
          v-model="searchQuery"
          @keydown.enter="handleSearch"
          type="text"
          placeholder="Cari transaksi, menu, stok..."
          class="h-9 w-full pl-9 pr-11 text-xs rounded-xl border outline-none transition-all duration-150 border-slate-200/80 dark:border-slate-700/80 bg-white dark:bg-slate-800/90 text-slate-800 dark:text-slate-200 placeholder:text-slate-400 dark:placeholder:text-slate-500 focus:border-[#4361EE] focus:ring-2 focus:ring-blue-500/15 shadow-xs"
        />
        <span class="absolute right-2.5 top-1/2 -translate-y-1/2 text-[10px] font-mono px-1.5 py-0.5 rounded-md border border-slate-200 dark:border-slate-700 text-slate-400 bg-slate-50 dark:bg-slate-800 pointer-events-none">
          ⌘K
        </span>
      </div>

      <!-- UI Version Switcher Pill (Instant Rollback to V1 / Switch to V2) with unified h-9 height -->
      <UiVersionSwitcherPill />

      <!-- Minimalist Theme Toggle (Light / Dark) with precision h-9 w-9 -->
      <button
        type="button"
        @click="appStore.toggleTheme"
        class="h-9 w-9 rounded-xl flex items-center justify-center transition-all duration-150 cursor-pointer border text-slate-600 dark:text-slate-300 hover:text-amber-500 dark:hover:text-amber-400 bg-white hover:bg-slate-50 dark:bg-slate-800/90 dark:hover:bg-slate-700/90 border-slate-200/80 dark:border-slate-700/80 shadow-xs hover:border-amber-300 dark:hover:border-amber-600 hover:shadow-sm"
        :title="appStore.resolvedTheme === 'dark' ? 'Ganti ke Mode Terang' : 'Ganti ke Mode Gelap'"
        aria-label="Toggle Theme"
      >
        <Sun v-if="appStore.resolvedTheme === 'dark'" class="w-4 h-4 text-amber-400 shrink-0" />
        <Moon v-else class="w-4 h-4 text-slate-600 shrink-0" />
      </button>

      <!-- Notification Bell with Ping and precision h-9 w-9 -->
      <div class="relative">
        <button
          type="button"
          @click="toggleNotifications"
          class="h-9 w-9 rounded-xl flex items-center justify-center transition-all duration-150 cursor-pointer border text-slate-600 dark:text-slate-300 hover:text-[#4361EE] dark:hover:text-blue-400 bg-white hover:bg-slate-50 dark:bg-slate-800/90 dark:hover:bg-slate-700/90 border-slate-200/80 dark:border-slate-700/80 shadow-xs hover:border-blue-300 dark:hover:border-blue-600 hover:shadow-sm"
          title="Notifikasi & Peringatan"
          aria-label="Lihat Notifikasi"
        >
          <Bell class="w-4 h-4 shrink-0" />
        </button>
        <span class="absolute top-2 right-2 w-2 h-2 rounded-full bg-rose-500 ring-2 ring-white dark:ring-slate-800 pointer-events-none"></span>
      </div>

      <div class="h-5 w-px bg-slate-200 dark:bg-slate-700/80 mx-0.5 hidden sm:block"></div>

      <!-- User Profile Dropdown Pill with precision h-9 height -->
      <n-dropdown :options="userMenuOptions" @select="handleUserMenu">
        <button
          type="button"
          class="h-9 inline-flex items-center gap-2.5 px-2 rounded-xl transition-all duration-150 cursor-pointer border border-slate-200/80 dark:border-slate-700/80 bg-white hover:bg-slate-50 dark:bg-slate-800/90 dark:hover:bg-slate-800 text-slate-800 dark:text-slate-200 hover:border-blue-300 dark:hover:border-blue-600 shadow-xs hover:shadow-sm group select-none"
        >
          <div class="w-6.5 h-6.5 rounded-lg bg-gradient-to-tr from-[#3B54DE] to-[#5572FA] text-white font-bold flex items-center justify-center shadow-xs text-xs tracking-tight shrink-0">
            {{ userInitials }}
          </div>
          <div class="text-left hidden sm:block">
            <div class="text-xs font-bold leading-tight group-hover:text-[#4361EE] dark:group-hover:text-blue-400 transition-colors truncate max-w-[130px]">
              {{ userDisplayName }}
            </div>
            <div class="text-[10px] font-medium capitalize text-slate-400 dark:text-slate-500 leading-none truncate max-w-[130px] mt-0.5">
              {{ userRoleName }}
            </div>
          </div>
          <ChevronDown class="w-3.5 h-3.5 text-slate-400 group-hover:text-slate-600 dark:group-hover:text-slate-200 transition-transform shrink-0" />
        </button>
      </n-dropdown>
    </div>

    <!-- Logout Loading Overlay -->
    <Teleport to="body">
      <transition name="fade">
        <div v-if="isLoggingOut" class="fixed inset-0 z-99999 flex flex-col items-center justify-center bg-slate-950/50 backdrop-blur-md">
          <div 
            class="rounded-3xl p-8 max-w-xs w-full shadow-2xl flex flex-col items-center text-center animate-in zoom-in-95 duration-200 border"
            style="background-color: var(--bg-primary); border-color: var(--border-color);"
          >
            <div class="relative w-14 h-14 mb-4 flex items-center justify-center">
              <div class="absolute inset-0 rounded-full border-3 border-rose-100 dark:border-rose-950 border-t-rose-600 animate-spin"></div>
              <div class="w-9 h-9 rounded-xl bg-rose-50 dark:bg-rose-950/60 text-rose-600 flex items-center justify-center">
                <LogOut class="w-4.5 h-4.5" />
              </div>
            </div>
            <h3 class="text-sm font-bold tracking-tight" style="color: var(--text-primary);">Mengakhiri Sesi...</h3>
            <p class="text-xs mt-1 text-slate-400">Menutup akses akun dengan aman</p>
          </div>
        </div>
      </transition>
    </Teleport>
  </header>
</template>

<script setup lang="ts">
import { h, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { NDropdown } from 'naive-ui'
import {
  Menu,
  PanelLeftClose,
  Search,
  Sun,
  Moon,
  Bell,
  ChevronDown,
  User,
  Database,
  Settings,
  LogOut
} from 'lucide-vue-next'
import BreadcrumbV2 from './BreadcrumbV2.vue'
import UiVersionSwitcherPill from '@/components/v2/common/UiVersionSwitcherPill.vue'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'
import { useNotificationStore } from '@/stores/notification.store'

const authStore = useAuthStore()
const appStore = useAppStore()
const notifyStore = useNotificationStore()
const router = useRouter()

const searchQuery = ref('')
const isLoggingOut = ref(false)

const userDisplayName = computed(() => authStore.fullName?.trim() || authStore.user?.firstName || 'Super Administrator')
const userRoleName = computed(() => authStore.currentRole || 'Administrator')
const userInitials = computed(() => {
  const name = userDisplayName.value.trim()
  const parts = name.split(' ')
  if (parts.length >= 2) return `${parts[0][0]}${parts[1][0]}`.toUpperCase()
  return name.slice(0, 2).toUpperCase() || 'SA'
})

const renderIcon = (icon: any) => () => h(icon, { class: 'w-4 h-4' })

const userMenuOptions = [
  { label: 'Profil Pengguna', key: 'profile', icon: renderIcon(User) },
  { label: 'Master Data Hub', key: 'master', icon: renderIcon(Database) },
  { label: 'Pengaturan & Peran', key: 'settings', icon: renderIcon(Settings) },
  { type: 'divider', key: 'd1' },
  { label: 'Keluar (Logout)', key: 'logout', icon: renderIcon(LogOut) }
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
  } else if (key === 'profile') {
    router.push('/dashboard')
  }
}

const toggleNotifications = () => {
  notifyStore.info('Tidak ada notifikasi penting yang belum dibaca.', 'Notifikasi Sistem')
}

const handleSearch = () => {
  if (!searchQuery.value.trim()) return
  const query = searchQuery.value.toLowerCase().trim()
  if (query.includes('pos') || query.includes('kasir') || query.includes('order')) {
    router.push('/pos')
  } else if (query.includes('stok') || query.includes('barang') || query.includes('gudang')) {
    router.push('/inventory')
  } else if (query.includes('karyawan') || query.includes('gaji') || query.includes('hris')) {
    router.push('/hris')
  } else if (query.includes('uang') || query.includes('jurnal') || query.includes('keuangan')) {
    router.push('/finance')
  } else if (query.includes('laporan') || query.includes('analitik')) {
    router.push('/reports')
  } else {
    notifyStore.info(`Hasil pencarian untuk "${searchQuery.value}"`, 'Pencarian Global')
  }
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
