<template>
  <header class="h-16 bg-white border-b border-slate-200/80 flex items-center justify-between px-6 z-10">
    <!-- Left: Breadcrumb / Page context -->
    <div class="flex items-center gap-3">
      <Breadcrumb />
    </div>

    <!-- Right: Search, Notifications, User Profile -->
    <div class="flex items-center gap-4">
      <!-- Search input -->
      <div class="relative hidden sm:block w-72">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
          <Search class="w-4 h-4" />
        </span>
        <input
          type="text"
          placeholder="Search orders, reports..."
          class="w-full pl-9 pr-4 py-1.5 text-xs bg-slate-100/80 hover:bg-slate-100 focus:bg-white text-slate-700 rounded-lg border border-slate-200 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all"
        />
      </div>

      <!-- Help Button -->
      <button class="w-8 h-8 rounded-full flex items-center justify-center text-slate-500 hover:bg-slate-100 transition-colors" title="Help & Support">
        <HelpCircle class="w-4 h-4" />
      </button>

      <!-- Notification Bell -->
      <div class="relative cursor-pointer">
        <div class="w-8 h-8 rounded-full flex items-center justify-center text-slate-500 hover:bg-slate-100 transition-colors">
          <Bell class="w-4 h-4" />
        </div>
        <span class="absolute top-0 right-0 w-4 h-4 bg-red-500 text-white text-[10px] font-bold rounded-full flex items-center justify-center shadow-sm">
          5
        </span>
      </div>

      <div class="h-6 w-px bg-slate-200 mx-1"></div>

      <!-- User Profile Dropdown -->
      <n-dropdown :options="userOptions" @select="handleUserMenu">
        <div class="flex items-center gap-2.5 cursor-pointer py-1 px-2 rounded-lg hover:bg-slate-50 transition-colors">
          <div class="w-9 h-9 rounded-full bg-gradient-to-tr from-blue-600 to-indigo-600 text-white font-semibold flex items-center justify-center shadow-sm text-xs">
            SA
          </div>
          <div class="text-left hidden md:block">
            <div class="text-xs font-semibold text-slate-800 leading-tight">Super Administrator</div>
            <div class="text-[10px] text-slate-500 font-medium">Headquarters</div>
          </div>
          <ChevronDown class="w-3.5 h-3.5 text-slate-400 ml-1" />
        </div>
      </n-dropdown>
    </div>
  </header>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { NDropdown } from 'naive-ui'
import Breadcrumb from './Breadcrumb.vue'
import { useAuthStore } from '@/stores/auth.store'
import { useRouter } from 'vue-router'
import {
  Search,
  HelpCircle,
  Bell,
  ChevronDown,
  User,
  Database,
  Settings,
  LogOut
} from 'lucide-vue-next'

const authStore = useAuthStore()
const router = useRouter()

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
    authStore.logout()
    router.push('/login')
  } else if (key === 'master') {
    router.push('/settings/master')
  } else if (key === 'settings') {
    router.push('/settings/roles')
  }
}
</script>
