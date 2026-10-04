<template>
  <aside 
    class="sidebar-container flex flex-col h-screen select-none shrink-0 transition-all duration-300 relative z-30 border-r"
    :class="isCollapsed ? 'w-16' : 'w-64'"
    style="background-color: var(--bg-sidebar); border-color: var(--border-color); color: var(--text-secondary);"
  >
    <!-- Modern Workspace Header (Aligned with AppHeader h-16) -->
    <div 
      class="h-16 flex items-center border-b shrink-0 transition-all duration-200"
      :class="isCollapsed ? 'justify-center px-2' : 'justify-between px-3.5'"
      style="border-color: var(--border-color);"
    >
      <!-- Brand Workspace Pill -->
      <button 
        type="button"
        class="flex items-center gap-2.5 min-w-0 text-left cursor-pointer bg-transparent border-0 p-0"
        :class="isCollapsed ? 'justify-center w-full' : ''"
        @click="isCollapsed && toggleSidebar()"
        :title="isCollapsed ? 'Klik untuk memperluas sidebar (Ctrl+B)' : undefined"
        :aria-label="isCollapsed ? 'Perluas sidebar' : 'Workspace Cafe ERP'"
      >
        <div class="w-8.5 h-8.5 rounded-xl bg-gradient-to-tr from-blue-600 via-blue-500 to-indigo-500 flex items-center justify-center text-white shadow-md shadow-blue-500/20 shrink-0">
          <Store class="w-4.5 h-4.5" />
        </div>
        <div v-if="!isCollapsed" class="min-w-0 transition-opacity duration-200">
          <div class="font-bold text-[13px] tracking-tight truncate leading-tight" style="color: var(--text-primary);">
            Cafe ERP
          </div>
          <div class="text-[10px] font-medium flex items-center gap-1.5 mt-0.5 truncate" style="color: var(--text-muted);">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0"></span>
            <span class="truncate">Senopati Central</span>
          </div>
        </div>
      </button>

      <!-- Collapse Toggle Button (Top Right when Expanded) -->
      <button
        v-if="!isCollapsed"
        type="button"
        @click="toggleSidebar"
        class="w-7 h-7 rounded-lg flex items-center justify-center transition-colors cursor-pointer sidebar-action-btn"
        title="Ciutkan Sidebar (Ctrl+B)"
        aria-label="Ciutkan Sidebar"
      >
        <PanelLeftClose class="w-4 h-4 shrink-0" />
      </button>
    </div>

    <!-- Quick Command / Filter Bar (Visible when Expanded) -->
    <div v-if="!isCollapsed" class="px-3 pt-3 pb-1 shrink-0">
      <div class="relative flex items-center">
        <span class="absolute left-2.5 flex items-center pointer-events-none text-slate-400">
          <Search class="w-3.5 h-3.5" />
        </span>
        <input
          ref="searchInputRef"
          v-model="searchQuery"
          type="text"
          placeholder="Cari modul..."
          aria-label="Cari menu atau fitur"
          class="w-full pl-8 pr-12 py-1.5 text-xs rounded-xl border outline-none transition-all sidebar-search-input"
        />
        <div class="absolute right-2 flex items-center gap-1">
          <button
            v-if="searchQuery"
            type="button"
            @click="searchQuery = ''"
            class="p-0.5 rounded hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors cursor-pointer"
            aria-label="Bersihkan pencarian"
          >
            <X class="w-3 h-3" />
          </button>
          <kbd v-else class="text-[9px] font-mono px-1 py-0.5 rounded border border-slate-200 dark:border-slate-700 text-slate-400 bg-slate-100/80 dark:bg-slate-800/80 pointer-events-none">
            ⌘K
          </kbd>
        </div>
      </div>
    </div>

    <!-- Navigation Menu with Custom Sleek Scrollbar -->
    <div class="flex-1 overflow-y-auto px-2.5 py-3 space-y-3.5 sidebar-scroll">
      <!-- Empty state when search produces no results -->
      <div v-if="filteredSections.length === 0" class="py-8 px-3 text-center">
        <p class="text-xs font-medium" style="color: var(--text-muted);">Tidak ada menu ditemukan</p>
        <button
          type="button"
          @click="searchQuery = ''"
          class="text-[11px] font-semibold text-blue-600 dark:text-blue-400 hover:underline mt-1.5 cursor-pointer"
        >
          Reset pencarian
        </button>
      </div>

      <!-- Navigation Sections -->
      <div v-for="(section, sIdx) in filteredSections" :key="section.title" class="space-y-1">
        <!-- Category Section Header -->
        <div 
          v-if="!isCollapsed" 
          class="px-2.5 pt-1.5 pb-1 text-[10px] font-bold tracking-widest uppercase select-none"
          style="color: var(--text-muted);"
        >
          {{ section.title }}
        </div>
        <div 
          v-else-if="sIdx > 0" 
          class="my-2 mx-auto w-6 border-t" 
          style="border-color: var(--border-color);"
        ></div>

        <!-- Section Menu Items -->
        <div class="space-y-0.5">
          <template v-for="item in section.items" :key="item.path || item.label">
            <!-- Item without children -->
            <router-link
              v-if="!item.children"
              :to="item.path!"
              class="group flex items-center rounded-xl text-[13px] font-medium transition-all duration-150 relative"
              :class="[
                isRouteActive(item.path!) ? 'sidebar-item-active font-semibold shadow-xs' : 'sidebar-item-inactive',
                isCollapsed ? 'justify-center p-2.5 mx-auto' : 'gap-3 px-3 py-2'
              ]"
              :title="isCollapsed ? item.label : undefined"
              :aria-label="item.label"
            >
              <!-- Left Accent Pill when Active -->
              <div 
                v-if="isRouteActive(item.path!) && !isCollapsed" 
                class="absolute left-0 top-1.5 bottom-1.5 w-1 rounded-r-full bg-blue-600 dark:bg-blue-400 shadow-sm shadow-blue-500/50"
              ></div>

              <component 
                :is="item.icon" 
                class="w-4.5 h-4.5 shrink-0 transition-colors" 
                :class="isRouteActive(item.path!) ? 'text-[var(--sidebar-item-active-text)]' : 'text-[var(--text-muted)] group-hover:text-[var(--sidebar-item-hover-text)]'"
              />
              <span v-if="!isCollapsed" class="truncate">{{ item.label }}</span>
            </router-link>

            <!-- Item with children / Collapsible -->
            <div v-else class="space-y-0.5">
              <button
                type="button"
                @click="isCollapsed ? (isCollapsed = false) : toggleGroup(item.label)"
                class="w-full group flex items-center rounded-xl text-[13px] font-medium transition-all duration-150 cursor-pointer relative"
                :class="[
                  isGroupActive(item) ? 'sidebar-item-active font-semibold' : 'sidebar-item-inactive',
                  isCollapsed ? 'justify-center p-2.5 mx-auto' : 'justify-between px-3 py-2'
                ]"
                :title="isCollapsed ? item.label : undefined"
                :aria-label="item.label"
              >
                <!-- Left Accent Pill when Group Active -->
                <div 
                  v-if="isGroupActive(item) && !isCollapsed" 
                  class="absolute left-0 top-1.5 bottom-1.5 w-1 rounded-r-full bg-blue-600 dark:bg-blue-400 shadow-sm shadow-blue-500/50"
                ></div>

                <div class="flex items-center gap-3 min-w-0" :class="isCollapsed ? 'justify-center' : ''">
                  <component 
                    :is="item.icon" 
                    class="w-4.5 h-4.5 shrink-0 transition-colors" 
                    :class="isGroupActive(item) ? 'text-[var(--sidebar-item-active-text)]' : 'text-[var(--text-muted)] group-hover:text-[var(--sidebar-item-hover-text)]'"
                  />
                  <span v-if="!isCollapsed" class="truncate">{{ item.label }}</span>
                </div>
                <ChevronDown 
                  v-if="!isCollapsed"
                  class="w-3.5 h-3.5 transition-transform duration-200 shrink-0" 
                  :class="openGroups[item.label] || searchQuery.trim().length > 0 ? 'rotate-180 text-[var(--sidebar-item-hover-text)]' : 'text-[var(--text-muted)]'" 
                />
              </button>

              <!-- Submenu Items (Indented with Guideline Tree) -->
              <div 
                v-if="!isCollapsed && (openGroups[item.label] || isGroupActive(item) || searchQuery.trim().length > 0)" 
                class="pl-3.5 ml-5 py-1 space-y-0.5 border-l transition-all"
                style="border-color: var(--border-color);"
              >
                <router-link
                  v-for="sub in item.children"
                  :key="sub.path"
                  :to="sub.path"
                  class="flex items-center justify-between px-2.5 py-1.5 rounded-lg text-xs font-medium transition-all group"
                  :class="isRouteActive(sub.path) ? 'sidebar-subitem-active font-semibold' : 'sidebar-subitem-inactive'"
                >
                  <div class="flex items-center gap-2 min-w-0">
                    <span 
                      class="w-1.5 h-1.5 rounded-full shrink-0 transition-colors"
                      :class="isRouteActive(sub.path) ? 'bg-blue-600 dark:bg-blue-400' : 'bg-slate-300 dark:bg-slate-600 group-hover:bg-slate-400'"
                    ></span>
                    <span class="truncate">{{ sub.label }}</span>
                  </div>
                  <span v-if="sub.badge" class="px-1.5 py-0.2 rounded text-[10px] bg-blue-100/80 dark:bg-blue-900/60 text-blue-700 dark:text-blue-300 font-bold shrink-0">
                    {{ sub.badge }}
                  </span>
                </router-link>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- Bottom Actions Section (Status + Theme Switcher) -->
    <div 
      class="p-2.5 border-t space-y-2 shrink-0 transition-all"
      style="background-color: var(--bg-sidebar); border-color: var(--border-color);"
    >
      <!-- Outlet Status Indicator (Expanded only) -->
      <div 
        v-if="!isCollapsed" 
        class="px-2.5 py-1.5 rounded-xl border flex items-center justify-between text-xs sidebar-status-card"
      >
        <div class="flex items-center gap-2 min-w-0">
          <span class="relative flex h-2 w-2 shrink-0">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
          <span class="text-[11px] font-medium truncate" style="color: var(--text-secondary);">Senopati Branch</span>
        </div>
        <span class="text-[10px] font-bold text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 dark:bg-emerald-500/15 border border-emerald-500/20 px-1.5 py-0.5 rounded-md uppercase tracking-wider">
          Online
        </span>
      </div>

      <!-- Theme Switcher Segmented Control (Light / Dark / System) -->
      <div 
        v-if="!isCollapsed" 
        role="radiogroup" 
        aria-label="Pilih tema tampilan antarmuka"
        class="p-1 rounded-xl border flex items-center text-xs font-semibold gap-1 theme-switcher-group"
      >
        <!-- Light -->
        <button
          type="button"
          role="radio"
          :aria-checked="appStore.themeMode === 'light'"
          aria-label="Mode Terang"
          @click="appStore.setTheme('light')"
          class="flex-1 py-1.5 px-2 rounded-lg flex items-center justify-center gap-1.5 transition-all cursor-pointer select-none focus-visible:ring-2 focus-visible:ring-blue-500"
          :class="appStore.themeMode === 'light' ? 'theme-pill-active' : 'theme-pill-inactive'"
          title="Mode Terang"
        >
          <Sun class="w-3.5 h-3.5 text-amber-500 shrink-0" />
          <span class="text-[11px]">Light</span>
        </button>

        <!-- Dark -->
        <button
          type="button"
          role="radio"
          :aria-checked="appStore.themeMode === 'dark'"
          aria-label="Mode Gelap"
          @click="appStore.setTheme('dark')"
          class="flex-1 py-1.5 px-2 rounded-lg flex items-center justify-center gap-1.5 transition-all cursor-pointer select-none focus-visible:ring-2 focus-visible:ring-blue-500"
          :class="appStore.themeMode === 'dark' ? 'theme-pill-active' : 'theme-pill-inactive'"
          title="Mode Gelap"
        >
          <Moon class="w-3.5 h-3.5 text-blue-400 shrink-0" />
          <span class="text-[11px]">Dark</span>
        </button>

        <!-- System -->
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

      <!-- Theme Switcher & Expand Toggle for Collapsed Sidebar -->
      <div v-else class="flex flex-col items-center gap-2">
        <button
          type="button"
          :aria-label="`Tema saat ini: ${appStore.themeMode}. Klik untuk ganti tema.`"
          @click="cycleTheme"
          class="w-10 h-10 rounded-xl flex items-center justify-center transition-all cursor-pointer border focus-visible:ring-2 focus-visible:ring-blue-500"
          style="background-color: var(--bg-content); border-color: var(--border-color); color: var(--text-primary);"
          :title="collapsedTooltip"
        >
          <Sun v-if="appStore.themeMode === 'light'" class="w-4 h-4 text-amber-500" />
          <Moon v-else-if="appStore.themeMode === 'dark'" class="w-4 h-4 text-blue-400" />
          <Monitor v-else class="w-4 h-4 text-blue-500" />
        </button>

        <button
          type="button"
          @click="toggleSidebar"
          class="w-10 h-10 rounded-xl flex items-center justify-center transition-all cursor-pointer sidebar-action-btn"
          title="Perluas Sidebar (Ctrl+B)"
          aria-label="Perluas Sidebar"
        >
          <PanelLeftOpen class="w-4 h-4 text-[var(--text-muted)]" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  LayoutDashboard,
  UtensilsCrossed,
  Package,
  Users,
  DollarSign,
  FileBarChart,
  Settings,
  ChevronDown,
  PanelLeftClose,
  PanelLeftOpen,
  Sun,
  Moon,
  Monitor,
  Store,
  Search,
  X
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
const authStore = useAuthStore()
const appStore = useAppStore()

const isCollapsed = computed({
  get: () => appStore.sidebarCollapsed,
  set: (val: boolean) => {
    appStore.sidebarCollapsed = val
  }
})

const toggleSidebar = () => {
  appStore.toggleSidebar()
}

const cycleTheme = () => {
  appStore.cycleTheme()
}

const collapsedTooltip = computed(() => {
  if (appStore.themeMode === 'light') return 'Mode Terang (Matahari) - Klik untuk beralih ke Mode Gelap'
  if (appStore.themeMode === 'dark') return 'Mode Gelap (Bulan) - Klik untuk beralih ke Mode Sistem'
  return `Mode Sistem Auto (${appStore.resolvedTheme === 'dark' ? 'Gelap' : 'Terang'}) - Klik untuk beralih ke Mode Terang`
})

const searchQuery = ref('')
const searchInputRef = ref<HTMLInputElement | null>(null)

const handleKeydown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    if (isCollapsed.value) {
      isCollapsed.value = false
    }
    setTimeout(() => {
      searchInputRef.value?.focus()
    }, 50)
  }
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'b') {
    e.preventDefault()
    toggleSidebar()
  }
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('keydown', handleKeydown)
  }
})

onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('keydown', handleKeydown)
  }
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
          { label: 'Shift & Kasir Balancing', path: '/pos/shifts' },
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
          { label: 'Data Master CRUD', path: '/settings/master', badge: 'Super Admin' },
          { label: 'Peran & Matriks Izin', path: '/settings/roles' }
        ]
      }
    ]
  }
]

const filteredSections = computed(() => {
  const role = (authStore.currentRole || '').toLowerCase()
  const q = searchQuery.value.trim().toLowerCase()

  return menuSections
    .map(section => {
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

      if (q) {
        items = items
          .map(item => {
            const matchesItem = item.label.toLowerCase().includes(q)
            if (item.children) {
              const matchedChildren = item.children.filter(c => c.label.toLowerCase().includes(q))
              if (matchedChildren.length > 0) {
                return {
                  ...item,
                  children: matchedChildren
                }
              }
            }
            return matchesItem ? item : null
          })
          .filter(Boolean) as MenuItem[]
      }

      return {
        ...section,
        items
      }
    })
    .filter(section => section.items.length > 0)
})
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

.sidebar-action-btn {
  color: var(--text-muted);
}
.sidebar-action-btn:hover {
  background-color: var(--sidebar-item-hover-bg);
  color: var(--sidebar-item-hover-text);
}

.sidebar-search-input {
  background-color: var(--bg-content);
  border-color: var(--border-color);
  color: var(--text-primary);
}
.sidebar-search-input:focus {
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.15);
}

.sidebar-status-card {
  background-color: var(--bg-content);
  border-color: var(--border-color);
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
  scrollbar-color: rgba(203, 213, 225, 0.4) transparent;
}

.sidebar-scroll::-webkit-scrollbar {
  width: 4px;
}

.sidebar-scroll::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-scroll::-webkit-scrollbar-thumb {
  background: rgba(203, 213, 225, 0.3);
  border-radius: 9999px;
  transition: background-color 0.2s ease;
}

.sidebar-scroll:hover::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.6);
}
</style>
