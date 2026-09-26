<template>
  <aside 
    class="sidebar-container flex flex-col h-screen select-none shrink-0 transition-all duration-300 relative z-30 border-r"
    :class="isCollapsed ? 'w-20' : 'w-64'"
    style="background-color: var(--bg-sidebar); border-color: var(--border-color); color: var(--text-secondary);"
  >
    <!-- Brand Header -->
    <div class="h-16 flex items-center px-4.5 gap-3 border-b shrink-0" style="border-color: var(--border-color);">
      <!-- Circular Logo matching Gridlines UI -->
      <div class="w-9 h-9 rounded-full bg-gradient-to-tr from-blue-600 via-blue-500 to-indigo-500 flex items-center justify-center text-white shadow-md shadow-blue-500/25 shrink-0">
        <Coffee class="w-4.5 h-4.5" />
      </div>
      <div v-if="!isCollapsed" class="min-w-0 transition-opacity duration-200">
        <h1 class="font-black text-[15px] tracking-tight leading-none flex items-center gap-1.5" style="color: var(--text-primary);">
          Cafe ERP
          <span class="text-[9px] font-bold uppercase tracking-wider bg-blue-50 dark:bg-blue-950/70 text-blue-600 dark:text-blue-400 border border-blue-200/60 dark:border-blue-800/60 px-1.5 py-0.5 rounded-full">Pro</span>
        </h1>
        <p class="text-[11px] font-medium truncate mt-1" style="color: var(--text-muted);">Management System</p>
      </div>
    </div>

    <!-- Navigation Menu with Custom Sleek Scrollbar -->
    <div class="flex-1 overflow-y-auto px-3 py-3 space-y-4 sidebar-scroll">
      <div v-for="section in filteredSections" :key="section.title" class="space-y-1">
        <!-- Category Section Header -->
        <div 
          v-if="!isCollapsed" 
          class="px-3 pt-1 pb-1 text-[11px] font-semibold tracking-wide uppercase"
          style="color: var(--text-muted);"
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
              :class="isRouteActive(item.path!) ? 'sidebar-item-active font-bold shadow-xs' : 'sidebar-item-inactive'"
              :title="isCollapsed ? item.label : undefined"
            >
              <component 
                :is="item.icon" 
                class="w-4.5 h-4.5 shrink-0 transition-colors" 
                :class="isRouteActive(item.path!) ? 'text-[var(--sidebar-item-active-text)]' : 'text-[var(--text-muted)] group-hover:text-[var(--sidebar-item-hover-text)]'"
              />
              <span v-if="!isCollapsed" class="truncate">{{ item.label }}</span>
            </router-link>

            <!-- Item with children / Collapsible -->
            <div v-else class="space-y-1">
              <button
                type="button"
                @click="toggleGroup(item.label)"
                class="w-full group flex items-center justify-between px-3 py-2.5 rounded-xl text-[13px] font-medium transition-all duration-150 cursor-pointer"
                :class="isGroupActive(item) ? 'sidebar-item-active font-bold' : 'sidebar-item-inactive'"
                :title="isCollapsed ? item.label : undefined"
              >
                <div class="flex items-center gap-3 min-w-0">
                  <component 
                    :is="item.icon" 
                    class="w-4.5 h-4.5 shrink-0 transition-colors" 
                    :class="isGroupActive(item) ? 'text-[var(--sidebar-item-active-text)]' : 'text-[var(--text-muted)] group-hover:text-[var(--sidebar-item-hover-text)]'"
                  />
                  <span v-if="!isCollapsed" class="truncate">{{ item.label }}</span>
                </div>
                <ChevronDown 
                  v-if="!isCollapsed"
                  class="w-4 h-4 transition-transform duration-200 shrink-0" 
                  :class="openGroups[item.label] ? 'rotate-180 text-[var(--sidebar-item-hover-text)]' : 'text-[var(--text-muted)]'" 
                />
              </button>

              <!-- Submenu Items (Indented with Guideline) -->
              <div 
                v-if="!isCollapsed && (openGroups[item.label] || isGroupActive(item))" 
                class="pl-4 pr-1 space-y-1 border-l ml-5 py-1 transition-all"
                style="border-color: var(--border-color);"
              >
                <router-link
                  v-for="sub in item.children"
                  :key="sub.path"
                  :to="sub.path"
                  class="flex items-center justify-between px-3 py-1.5 rounded-lg text-xs font-medium transition-all"
                  :class="isRouteActive(sub.path) ? 'sidebar-subitem-active font-bold' : 'sidebar-subitem-inactive'"
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
    <div 
      class="p-3 border-t space-y-2.5 shrink-0"
      style="background-color: var(--bg-sidebar); border-color: var(--border-color);"
    >
      <!-- Collapse Menu Toggle -->
      <button
        type="button"
        @click="isCollapsed = !isCollapsed"
        class="w-full flex items-center gap-3 px-3 py-2 text-xs font-semibold rounded-xl transition-all cursor-pointer sidebar-collapse-btn"
        :title="isCollapsed ? 'Perluas Sidebar' : 'Ciutkan Sidebar'"
      >
        <Columns2 class="w-4 h-4 shrink-0 text-[var(--text-muted)]" />
        <span v-if="!isCollapsed" class="truncate">Collapse menu</span>
      </button>

      <!-- Theme Switcher Segmented Control (Light / Dark / System) -->
      <div 
        v-if="!isCollapsed" 
        role="radiogroup" 
        aria-label="Pilih tema tampilan antarmuka"
        class="p-1 rounded-xl border flex items-center text-xs font-semibold gap-1 theme-switcher-group"
      >
        <!-- Light (Matahari) -->
        <button
          type="button"
          role="radio"
          :aria-checked="appStore.themeMode === 'light'"
          aria-label="Mode Terang (Matahari)"
          @click="appStore.setTheme('light')"
          class="flex-1 py-1.5 px-2 rounded-lg flex items-center justify-center gap-1.5 transition-all cursor-pointer select-none focus-visible:ring-2 focus-visible:ring-blue-500"
          :class="appStore.themeMode === 'light' ? 'theme-pill-active' : 'theme-pill-inactive'"
          title="Mode Terang (Matahari: Putih & Biru)"
        >
          <Sun class="w-3.5 h-3.5 text-amber-500 shrink-0" />
          <span class="text-[11px]">Light</span>
        </button>

        <!-- Dark (Bulan) -->
        <button
          type="button"
          role="radio"
          :aria-checked="appStore.themeMode === 'dark'"
          aria-label="Mode Gelap (Bulan)"
          @click="appStore.setTheme('dark')"
          class="flex-1 py-1.5 px-2 rounded-lg flex items-center justify-center gap-1.5 transition-all cursor-pointer select-none focus-visible:ring-2 focus-visible:ring-blue-500"
          :class="appStore.themeMode === 'dark' ? 'theme-pill-active' : 'theme-pill-inactive'"
          title="Mode Gelap (Bulan: Hitam & Biru)"
        >
          <Moon class="w-3.5 h-3.5 text-blue-400 shrink-0" />
          <span class="text-[11px]">Dark</span>
        </button>

        <!-- System (Sistem) -->
        <button
          type="button"
          role="radio"
          :aria-checked="appStore.themeMode === 'system'"
          aria-label="Ikuti Tema Sistem OS"
          @click="appStore.setTheme('system')"
          class="flex-1 py-1.5 px-2 rounded-lg flex items-center justify-center gap-1.5 transition-all cursor-pointer select-none focus-visible:ring-2 focus-visible:ring-blue-500"
          :class="appStore.themeMode === 'system' ? 'theme-pill-active' : 'theme-pill-inactive'"
          title="Otomatis mengikuti preferensi OS Sistem"
        >
          <Monitor class="w-3.5 h-3.5 shrink-0" :class="appStore.themeMode === 'system' ? 'text-blue-500 dark:text-blue-400' : 'text-slate-400'" />
          <span class="text-[11px]">Auto</span>
        </button>
      </div>

      <!-- Theme Switcher for Collapsed Sidebar (Cycles Light -> Dark -> System) -->
      <div v-else class="flex justify-center">
        <button
          type="button"
          role="button"
          :aria-label="`Tema saat ini: ${appStore.themeMode}. Klik untuk ganti tema.`"
          @click="cycleTheme"
          class="w-10 h-10 rounded-xl flex items-center justify-center transition-all cursor-pointer border focus-visible:ring-2 focus-visible:ring-blue-500"
          style="background-color: var(--bg-content); border-color: var(--border-color); color: var(--text-primary);"
          :title="collapsedTooltip"
        >
          <Sun v-if="appStore.themeMode === 'light'" class="w-4.5 h-4.5 text-amber-500" />
          <Moon v-else-if="appStore.themeMode === 'dark'" class="w-4.5 h-4.5 text-blue-400" />
          <Monitor v-else class="w-4.5 h-4.5 text-blue-500" />
        </button>
      </div>

      <!-- User Card & Logout -->
      <div class="pt-2 border-t flex items-center justify-between gap-2 px-1" style="border-color: var(--border-color);">
        <div class="flex items-center gap-2.5 min-w-0">
          <div class="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800 flex items-center justify-center text-xs font-black shrink-0 shadow-2xs">
            {{ avatarInitials }}
          </div>
          <div v-if="!isCollapsed" class="min-w-0">
            <div class="text-xs font-bold truncate leading-tight" style="color: var(--text-primary);">{{ userName }}</div>
            <div class="text-[10px] font-medium truncate flex items-center gap-1 mt-0.5" style="color: var(--text-muted);">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
              {{ userRoleBadge }}
            </div>
          </div>
        </div>
        <button 
          @click="handleLogout" 
          title="Keluar dari Sistem" 
          class="p-1.5 rounded-lg hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/40 transition-colors cursor-pointer shrink-0"
          style="color: var(--text-muted);"
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
  Moon,
  Monitor
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

const cycleTheme = () => {
  if (appStore.themeMode === 'light') {
    appStore.setTheme('dark')
  } else if (appStore.themeMode === 'dark') {
    appStore.setTheme('system')
  } else {
    appStore.setTheme('light')
  }
}

const collapsedTooltip = computed(() => {
  if (appStore.themeMode === 'light') return 'Mode Terang (Matahari) - Klik untuk beralih ke Mode Gelap'
  if (appStore.themeMode === 'dark') return 'Mode Gelap (Bulan) - Klik untuk beralih ke Mode Sistem'
  return `Mode Sistem Auto (${appStore.resolvedTheme === 'dark' ? 'Gelap' : 'Terang'}) - Klik untuk beralih ke Mode Terang`
})

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
/* Sidebar menu items driven by CSS tokens */
.sidebar-item-active {
  background-color: var(--sidebar-item-active-bg);
  color: var(--sidebar-item-active-text);
  border: 1px solid var(--sidebar-item-active-border);
}

.sidebar-item-inactive {
  color: var(--text-secondary);
}
.sidebar-item-inactive:hover {
  background-color: var(--sidebar-item-hover-bg);
  color: var(--sidebar-item-hover-text);
}

.sidebar-subitem-active {
  background-color: var(--sidebar-item-active-bg);
  color: var(--sidebar-item-active-text);
}

.sidebar-subitem-inactive {
  color: var(--text-muted);
}
.sidebar-subitem-inactive:hover {
  background-color: var(--sidebar-item-hover-bg);
  color: var(--sidebar-item-hover-text);
}

.sidebar-collapse-btn {
  color: var(--text-muted);
}
.sidebar-collapse-btn:hover {
  background-color: var(--sidebar-item-hover-bg);
  color: var(--sidebar-item-hover-text);
}

.theme-switcher-group {
  background-color: var(--bg-content);
  border-color: var(--border-color);
}

.theme-pill-active {
  background-color: var(--bg-primary);
  color: var(--accent-primary);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
  font-weight: 700;
}

.theme-pill-inactive {
  color: var(--text-muted);
}
.theme-pill-inactive:hover {
  color: var(--text-primary);
}

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
