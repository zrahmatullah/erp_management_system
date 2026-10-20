<template>
  <header 
    class="h-16 border-b flex items-center justify-between px-4 sm:px-6 z-20 shrink-0 sticky top-0 transition-colors backdrop-blur-md"
    style="background-color: var(--bg-topbar); border-color: var(--border-color);"
  >
    <!-- Left: Sidebar Toggle (Mobile / Desktop) & Modern Breadcrumb -->
    <div class="flex items-center gap-3">
      <button
        type="button"
        @click="appStore.toggleSidebar"
        class="w-9 h-9 rounded-xl flex items-center justify-center transition-all cursor-pointer border text-slate-600 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 bg-slate-100/80 hover:bg-slate-200/80 dark:bg-slate-800/80 dark:hover:bg-slate-700/80 border-slate-200/70 dark:border-slate-700/60"
        :title="appStore.sidebarCollapsed ? 'Perluas Sidebar' : 'Ciutkan Sidebar'"
        aria-label="Toggle Sidebar"
      >
        <Menu v-if="appStore.sidebarCollapsed" class="w-4 h-4" />
        <PanelLeftClose v-else class="w-4 h-4" />
      </button>

      <div class="h-5 w-px bg-slate-200 dark:bg-slate-800 mx-1 hidden sm:block"></div>

      <BreadcrumbV2 class="hidden md:flex" />
    </div>

    <!-- Right: Quick Search, UI Version Switcher, Theme Toggle, Notifications, Profile -->
    <div class="flex items-center gap-2 sm:gap-3">
      <!-- Minimalist Command Search Bar -->
      <div class="relative hidden lg:block w-64 xl:w-72">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
          <Search class="w-3.5 h-3.5" />
        </span>
        <input
          v-model="searchQuery"
          @keydown.enter="handleSearch"
          type="text"
          placeholder="Cari transaksi, menu, stok..."
          class="w-full pl-9 pr-14 py-1.5 text-xs rounded-xl border outline-none transition-all focus:ring-2 focus:ring-blue-500/20"
          style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
        />
        <span class="absolute right-2.5 top-1/2 -translate-y-1/2 text-[10px] font-mono px-1.5 py-0.5 rounded border border-slate-200 dark:border-slate-700 text-slate-400 bg-slate-100 dark:bg-slate-800 pointer-events-none">
          ⌘K
        </span>
      </div>

      <!-- UI Version Switcher Pill (Instant Rollback to V1 / Switch to V2) -->
      <UiVersionSwitcherPill />

      <!-- Minimalist Theme Toggle (Light / Dark) -->
      <button
        type="button"
        @click="appStore.toggleTheme"
        class="w-8.5 h-8.5 rounded-xl flex items-center justify-center transition-colors cursor-pointer border text-slate-600 dark:text-slate-300 hover:text-amber-500 dark:hover:text-blue-400 bg-slate-100/80 hover:bg-slate-200/80 dark:bg-slate-800/80 dark:hover:bg-slate-700/80 border-slate-200/70 dark:border-slate-700/60"
        :title="appStore.resolvedTheme === 'dark' ? 'Ganti ke Mode Terang' : 'Ganti ke Mode Gelap'"
        aria-label="Toggle Theme"
      >
        <Sun v-if="appStore.resolvedTheme === 'dark'" class="w-4 h-4 text-amber-400" />
        <Moon v-else class="w-4 h-4 text-slate-600" />
      </button>

      <!-- Notification Bell with Ping -->
      <div class="relative">
        <button
          type="button"
          @click="toggleNotifications"
          class="w-8.5 h-8.5 rounded-xl flex items-center justify-center transition-colors cursor-pointer border text-slate-600 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 bg-slate-100/80 hover:bg-slate-200/80 dark:bg-slate-800/80 dark:hover:bg-slate-700/80 border-slate-200/70 dark:border-slate-700/60"
          title="Notifikasi & Peringatan"
          aria-label="Lihat Notifikasi"
        >
          <Bell class="w-4 h-4" />
        </button>
        <span class="absolute top-1 right-1 w-2 h-2 rounded-full bg-rose-500 ring-2 ring-white dark:ring-slate-900 pointer-events-none"></span>
      </div>

      <div class="h-6 w-px bg-slate-200 dark:bg-slate-800 mx-0.5"></div>

      <!-- User Profile Dropdown Pill -->
      <n-dropdown :options="userMenuOptions" @select="handleUserMenu">
        <button
          type="button"
          class="flex items-center gap-2.5 py-1 px-2 rounded-xl transition-all cursor-pointer border hover:shadow-xs group"
          style="border-color: var(--border-color); background-color: var(--bg-primary);"
        >
          <div class="w-8 h-8 rounded-lg bg-gradient-to-tr from-blue-600 to-indigo-600 text-white font-bold flex items-center justify-center shadow-xs text-xs tracking-tight">
            {{ userInitials }}
          </div>
          <div class="text-left hidden sm:block">
            <div class="text-xs font-semibold leading-tight group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors" style="color: var(--text-primary);">
              {{ userDisplayName }}
            </div>
            <div class="text-[10px] font-medium capitalize text-slate-400">
              {{ userRoleName }}
            </div>
          </div>
          <ChevronDown class="w-3.5 h-3.5 text-slate-400 group-hover:text-slate-600 transition-transform" />
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
  LogOut,
  Sparkles
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

