<template>
  <aside 
    class="bg-white dark:bg-slate-900 text-slate-700 dark:text-slate-200 flex flex-col h-screen select-none border-r border-slate-200/80 dark:border-slate-800 shrink-0 transition-all duration-300 relative z-30"
    :class="isCollapsed ? 'w-20' : 'w-64'"
  >
    <!-- Brand Header -->
    <div class="h-16 flex items-center px-4.5 gap-3 border-b border-slate-100 dark:border-slate-800 shrink-0">
      <!-- Circular Logo matching Gridlines UI -->
      <div class="w-9 h-9 rounded-full bg-gradient-to-tr from-blue-600 via-blue-500 to-indigo-500 flex items-center justify-center text-white shadow-md shadow-blue-500/25 shrink-0">
        <Coffee class="w-4.5 h-4.5" />
      </div>
      <div v-if="!isCollapsed" class="min-w-0 transition-opacity duration-200">
        <h1 class="text-slate-900 dark:text-white font-black text-[15px] tracking-tight leading-none flex items-center gap-1.5">
          Cafe ERP
          <span class="text-[9px] font-bold uppercase tracking-wider bg-blue-50 dark:bg-blue-950/70 text-blue-600 dark:text-blue-400 border border-blue-200/60 dark:border-blue-800/60 px-1.5 py-0.5 rounded-full">Pro</span>
        </h1>
        <p class="text-[11px] text-slate-400 dark:text-slate-500 font-medium truncate mt-1">Management System</p>
      </div>
    </div>

    <!-- Navigation Menu with Custom Sleek Scrollbar -->
    <div class="flex-1 overflow-y-auto px-3 py-3 space-y-4 sidebar-scroll">
      <div v-for="section in filteredSections" :key="section.title" class="space-y-1">
        <!-- Category Section Header -->
        <div 
          v-if="!isCollapsed" 
          class="px-3 pt-1 pb-1 text-[11px] font-semibold text-slate-400 dark:text-slate-500 tracking-wide"
        >
          {{ section.title }}
        </div>

        <!-- Section Menu Items -->
        <div class="space-y-1">
          <template v-for="item in section.items" :key="item.path || item.label">
            <!-- Item without children -->
            <router-link
              v-if="!item.children"
              :to="item.path!"
              class="group flex items-center gap-3 px-3 py-2.5 rounded-xl text-[13px] font-medium transition-all duration-150 relative"
              :class="isRouteActive(item.path!) 
                ? 'bg-blue-50/90 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 font-bold border border-blue-200/50 dark:border-blue-800/50 shadow-xs' 
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100/70 dark:hover:bg-slate-800/70'"
              :title="isCollapsed ? item.label : undefined"
            >
              <component 
                :is="item.icon" 
                class="w-4.5 h-4.5 shrink-0 transition-colors" 
                :class="isRouteActive(item.path!) ? 'text-blue-600 dark:text-blue-400' : 'text-slate-400 dark:text-slate-500 group-hover:text-slate-700 dark:group-hover:text-slate-200'"
              />
              <span v-if="!isCollapsed" class="truncate">{{ item.label }}</span>
            </router-link>

            <!-- Item with children / Collapsible -->
            <div v-else class="space-y-1">
              <button
                type="button"
                @click="toggleGroup(item.label)"
                class="w-full group flex items-center justify-between px-3 py-2.5 rounded-xl text-[13px] font-medium transition-all duration-150 cursor-pointer"
                :class="isGroupActive(item) 
                  ? 'text-blue-600 dark:text-blue-400 font-bold bg-blue-50/40 dark:bg-blue-900/20' 
                  : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100/70 dark:hover:bg-slate-800/70'"
                :title="isCollapsed ? item.label : undefined"
              >
                <div class="flex items-center gap-3 min-w-0">
                  <component 
                    :is="item.icon" 
                    class="w-4.5 h-4.5 shrink-0 transition-colors" 
                    :class="isGroupActive(item) ? 'text-blue-600 dark:text-blue-400' : 'text-slate-400 dark:text-slate-500 group-hover:text-slate-700 dark:group-hover:text-slate-200'"
                  />
                  <span v-if="!isCollapsed" class="truncate">{{ item.label }}</span>
                </div>
                <ChevronDown 
                  v-if="!isCollapsed"
                  class="w-4 h-4 text-slate-400 dark:text-slate-500 transition-transform duration-200 shrink-0" 
                  :class="openGroups[item.label] ? 'rotate-180 text-slate-600 dark:text-slate-300' : ''" 
                />
              </button>

              <!-- Submenu Items (Indented with Guideline) -->
              <div 
                v-if="!isCollapsed && (openGroups[item.label] || isGroupActive(item))" 
                class="pl-4 pr-1 space-y-1 border-l border-slate-200/80 dark:border-slate-800 ml-5 py-1 transition-all"
              >
                <router-link
                  v-for="sub in item.children"
                  :key="sub.path"
                  :to="sub.path"
                  class="flex items-center justify-between px-3 py-1.5 rounded-lg text-xs font-medium transition-all"
                  :class="isRouteActive(sub.path)
                    ? 'bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 font-bold'
                    : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100/60 dark:hover:bg-slate-800/60'"
                >
                  <span class="truncate">{{ sub.label }}</span>
                  <span v-if="sub.badge" class="px-1.5 py-0.5 rounded text-[10px] bg-blue-100/80 dark:bg-blue-900/60 text-blue-700 dark:text-blue-300 font-bold shrink-0">
                    {{ sub.badge }}
                  </span>
                </router-link>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- Bottom Actions Section (Collapse + Theme Switcher + User Profile) -->
    <div class="p-3 border-t border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/90 space-y-2.5 shrink-0">
      <!-- Collapse Menu Toggle -->
      <button
        type="button"
        @click="isCollapsed = !isCollapsed"
        class="w-full flex items-center gap-3 px-3 py-2 text-xs font-semibold text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100/80 dark:hover:bg-slate-800/80 rounded-xl transition-all cursor-pointer"
        :title="isCollapsed ? 'Perluas Sidebar' : 'Ciutkan Sidebar'"
      >
        <Columns2 class="w-4 h-4 text-slate-400 dark:text-slate-500 shrink-0" />
        <span v-if="!isCollapsed" class="truncate">Collapse menu</span>
      </button>

      <!-- Theme Switcher Segmented Capsule (Light / Dark) -->
      <div class="bg-slate-200/60 dark:bg-slate-800 p-1 rounded-2xl flex items-center text-xs font-medium gap-1">
        <button
          type="button"
          @click="appStore.setTheme('light')"
          class="flex-1 py-1.5 px-2.5 rounded-xl flex items-center justify-center gap-1.5 transition-all cursor-pointer"
          :class="appStore.theme === 'light' ? 'bg-white bg-keep-white text-slate-900 shadow-xs font-bold' : 'text-slate-400 dark:text-slate-500 hover:text-slate-600 dark:hover:text-slate-300'"
          title="Mode Terang (Light Mode)"
        >
          <Sun class="w-3.5 h-3.5" :class="appStore.theme === 'light' ? 'text-amber-500' : 'text-slate-400 dark:text-slate-500'" />
          <span v-if="!isCollapsed">Light</span>
        </button>
        <button
          type="button"
          @click="appStore.setTheme('dark')"
          class="flex-1 py-1.5 px-2.5 rounded-xl flex items-center justify-center gap-1.5 transition-all cursor-pointer"
          :class="appStore.theme === 'dark' ? 'bg-slate-900 dark:bg-slate-700 text-white shadow-xs font-bold' : 'text-slate-400 dark:text-slate-500 hover:text-slate-600 dark:hover:text-slate-300'"
          title="Mode Gelap (Dark Mode)"
        >
          <Moon class="w-3.5 h-3.5" :class="appStore.theme === 'dark' ? 'text-indigo-400' : 'text-slate-400 dark:text-slate-500'" />
          <span v-if="!isCollapsed">Dark</span>
        </button>
      </div>

      <!-- User Card & Logout -->
      <div class="pt-2 border-t border-slate-200/60 dark:border-slate-800 flex items-center justify-between gap-2 px-1">
        <div class="flex items-center gap-2.5 min-w-0">
          <div class="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800 flex items-center justify-center text-xs font-black shrink-0 shadow-2xs">
            {{ avatarInitials }}
          </div>
          <div v-if="!isCollapsed" class="min-w-0">
            <div class="text-xs font-bold text-slate-800 dark:text-slate-200 truncate leading-tight">{{ userName }}</div>
            <div class="text-[10px] text-slate-400 dark:text-slate-500 font-medium truncate flex items-center gap-1 mt-0.5">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
              {{ userRoleBadge }}
            </div>
          </div>
        </div>
        <button 
          @click="handleLogout" 
          title="Keluar dari Sistem" 
          class="p-1.5 rounded-lg text-slate-400 dark:text-slate-500 hover:text-rose-600 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors cursor-pointer shrink-0"
        >
          <LogOut class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  LayoutDashboard,
  UtensilsCrossed,
  Package,
  Users,
  DollarSign,
  FileBarChart,
  Settings,
  Coffee,
  ChevronDown,
  LogOut,
  Columns2,
  Sun,
  Moon
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth.store'
import { useAppStore } from '@/stores/app.store'

interface MenuItem {
  label: string
  path?: string
  icon: any
  badge?: string
  children?: {
    label: string
    path: string
    badge?: string
  }[]
}

interface MenuSection {
  title: string
  items: MenuItem[]
}

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const isCollapsed = ref(false)
const currentTheme = computed<'light' | 'dark'>({
  get: () => appStore.theme,
  set: (val) => appStore.setTheme(val)
})

const userName = computed(() => {
  return authStore.fullName || authStore.user?.firstName || 'Pengguna'
})

const userRoleBadge = computed(() => {
  return authStore.currentRole || 'Staff'
})

const avatarInitials = computed(() => {
  const name = userName.value.trim()
  if (!name) return 'U'
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.slice(0, 2).toUpperCase()
})

const openGroups = ref<Record<string, boolean>>({
  'POS & Pesanan': true,
  'Inventori': true,
  'HRIS & Staf': false,
  'Keuangan': false,
  'Pengaturan': false
})

const toggleGroup = (group: string) => {
  openGroups.value[group] = !openGroups.value[group]
}

const isRouteActive = (path: string) => {
  return route.path === path
}

const isGroupActive = (item: any) => {
  if (!item.children) return false
  return item.children.some((sub: any) => route.path.startsWith(sub.path))
}

const menuSections: MenuSection[] = [
  {
    title: 'Home',
    items: [
      {
        label: 'Dashboard',
        path: '/dashboard',
        icon: LayoutDashboard
      },
      {
        label: 'POS & Pesanan',
        icon: UtensilsCrossed,
        children: [
          { label: 'Kasir POS', path: '/pos' },
          { label: 'Denah Meja', path: '/pos/tables' },
          { label: 'Kitchen KDS', path: '/pos/kds' },
          { label: 'Riwayat Transaksi', path: '/pos/transactions' },
          { label: 'Katalog Menu & Produk', path: '/menu/products' }
        ]
      }
    ]
  },
  {
    title: 'Gudang & Inventori',
    items: [
      {
        label: 'Inventori',
        icon: Package,
        children: [
          { label: 'Daftar Stok', path: '/inventory' },
          { label: 'Kartu Stok', path: '/inventory/stock-card' },
          { label: 'Stock Opname', path: '/inventory/opname' },
          { label: 'P2P Procurement Hub', path: '/inventory/po', badge: 'P2P' }
        ]
      }
    ]
  },
  {
    title: 'HRIS & Karyawan',
    items: [
      {
        label: 'HRIS & Staf',
        icon: Users,
        children: [
          { label: 'Daftar Karyawan', path: '/hris' },
          { label: 'Absensi Presensi', path: '/hris/attendance' },
          { label: 'Jadwal Shift', path: '/hris/shifts' },
          { label: 'Pengajuan Cuti', path: '/hris/leaves' },
          { label: 'Penggajian / Payroll', path: '/hris/payroll' }
        ]
      }
    ]
  },
  {
    title: 'Keuangan & Laporan',
    items: [
      {
        label: 'Keuangan',
        icon: DollarSign,
        children: [
          { label: 'Ikhtisar Keuangan', path: '/finance' },
          { label: 'Jurnal Umum', path: '/finance/journal' }
        ]
      },
      {
        label: 'Laporan & Analitik',
        path: '/reports',
        icon: FileBarChart
      }
    ]
  },
  {
    title: 'Pengaturan',
    items: [
      {
        label: 'Pengaturan',
        icon: Settings,
        children: [
          { label: 'Profil Perusahaan', path: '/settings/company' },
          { label: '⭐ Data Master CRUD', path: '/settings/master', badge: 'Super Admin' },
          { label: 'Peran & Matriks Izin', path: '/settings/roles' }
        ]
      }
    ]
  }
]

const filteredSections = computed(() => {
  const role = (authStore.currentRole || '').toLowerCase()

  return menuSections.map(section => {
    let items = section.items

    if (!role || role.includes('super admin') || role.includes('owner') || role.includes('admin')) {
      // Super admin has full access
    } else if (role.includes('manager')) {
      items = items.filter(m => m.label !== 'Pengaturan')
    } else if (role.includes('kasir') || role.includes('cashier')) {
      items = items
        .filter(m => m.label === 'POS & Pesanan' || m.label === 'Dashboard' || m.label === 'HRIS & Staf')
        .map(m => {
          if (m.label === 'HRIS & Staf' && m.children) {
            return {
              ...m,
              children: m.children.filter(c => c.path === '/hris/attendance' || c.path === '/hris/leaves')
            }
          }
          return m
        })
    } else if (role.includes('gudang') || role.includes('warehouse') || role.includes('inventory')) {
      items = items
        .filter(m => m.label === 'Inventori' || m.label === 'Dashboard' || m.label === 'HRIS & Staf')
        .map(m => {
          if (m.label === 'HRIS & Staf' && m.children) {
            return {
              ...m,
              children: m.children.filter(c => c.path === '/hris/attendance' || c.path === '/hris/leaves')
            }
          }
          return m
        })
    } else if (role.includes('akuntan') || role.includes('finance')) {
      items = items
        .filter(m => m.label === 'Keuangan' || m.label === 'Laporan & Analitik' || m.label === 'Dashboard' || m.label === 'Inventori' || m.label === 'HRIS & Staf')
        .map(m => {
          if (m.label === 'Inventori' && m.children) {
            return {
              ...m,
              children: m.children.filter(c => c.path === '/inventory/po')
            }
          }
          if (m.label === 'HRIS & Staf' && m.children) {
            return {
              ...m,
              children: m.children.filter(c => c.path === '/hris/payroll')
            }
          }
          return m
        })
    }

    return {
      ...section,
      items
    }
  }).filter(section => section.items.length > 0)
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
/* Ultra-sleek, modern scrollbar */
.sidebar-scroll {
  scrollbar-width: thin;
  scrollbar-color: rgba(203, 213, 225, 0.5) transparent;
}

.sidebar-scroll::-webkit-scrollbar {
  width: 4px;
}

.sidebar-scroll::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-scroll::-webkit-scrollbar-thumb {
  background: rgba(203, 213, 225, 0.4);
  border-radius: 9999px;
  transition: background-color 0.2s ease;
}

.sidebar-scroll:hover::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.7);
}
</style>
