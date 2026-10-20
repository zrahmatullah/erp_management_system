<template>
  <aside 
    class="sidebar-v2 flex flex-col h-screen select-none shrink-0 transition-all duration-300 relative z-30 border-r backdrop-blur-md"
    :class="appStore.sidebarCollapsed ? 'w-18' : 'w-64'"
    style="background-color: var(--bg-sidebar); border-color: var(--border-color); color: var(--text-secondary);"
  >
    <!-- Brand / Workspace Header -->
    <div 
      class="h-16 flex items-center border-b shrink-0 transition-all duration-200"
      :class="appStore.sidebarCollapsed ? 'justify-center px-2' : 'justify-between px-4'"
      style="border-color: var(--border-color);"
    >
      <button 
        type="button"
        class="flex items-center gap-3 min-w-0 text-left cursor-pointer bg-transparent border-0 p-0"
        :class="appStore.sidebarCollapsed ? 'justify-center w-full' : ''"
        @click="appStore.sidebarCollapsed && appStore.toggleSidebar()"
        :title="appStore.sidebarCollapsed ? 'Klik untuk memperluas sidebar (Ctrl+B)' : undefined"
        aria-label="Cafe ERP Workspace"
      >
        <div class="w-9 h-9 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white shadow-sm shadow-blue-500/20 shrink-0">
          <Coffee class="w-5 h-5" />
        </div>
        <div v-if="!appStore.sidebarCollapsed" class="min-w-0 transition-opacity duration-200">
          <div class="font-bold text-sm tracking-tight truncate leading-tight" style="color: var(--text-primary);">
            Cafe ERP <span class="text-[10px] font-mono px-1 py-0.2 ml-1 rounded bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300">V2</span>
          </div>
          <div class="text-[11px] font-medium flex items-center gap-1.5 mt-0.5 truncate text-slate-400">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0"></span>
            <span class="truncate">Senopati Central</span>
          </div>
        </div>
      </button>

      <!-- Collapse Toggle Button -->
      <button
        v-if="!appStore.sidebarCollapsed"
        type="button"
        @click="appStore.toggleSidebar"
        class="w-7 h-7 rounded-lg flex items-center justify-center transition-colors cursor-pointer text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800"
        title="Ciutkan Sidebar (Ctrl+B)"
        aria-label="Ciutkan Sidebar"
      >
        <PanelLeftClose class="w-4 h-4 shrink-0" />
      </button>
    </div>

    <!-- Quick Search Filter Bar (Expanded) -->
    <div v-if="!appStore.sidebarCollapsed" class="px-3.5 pt-3 pb-1 shrink-0">
      <div class="relative flex items-center">
        <span class="absolute left-2.5 flex items-center pointer-events-none text-slate-400">
          <Search class="w-3.5 h-3.5" />
        </span>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Filter menu..."
          aria-label="Filter menu"
          class="w-full pl-8 pr-8 py-1.5 text-xs rounded-xl border outline-none transition-all focus:ring-1 focus:ring-blue-500"
          style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
        />
        <button
          v-if="searchQuery"
          type="button"
          @click="searchQuery = ''"
          class="absolute right-2 p-0.5 rounded text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 cursor-pointer"
        >
          <X class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>

    <!-- Navigation Menu Items -->
    <div class="flex-1 overflow-y-auto px-3 py-3 space-y-4 sidebar-v2-scroll">
      <!-- Empty Search Result -->
      <div v-if="filteredSections.length === 0" class="py-8 px-2 text-center">
        <p class="text-xs text-slate-400">Menu tidak ditemukan</p>
        <button
          type="button"
          @click="searchQuery = ''"
          class="text-xs text-blue-600 dark:text-blue-400 hover:underline mt-1 cursor-pointer font-medium"
        >
          Reset Filter
        </button>
      </div>

      <!-- Sections Loop -->
      <div v-for="(section, sIdx) in filteredSections" :key="section.title" class="space-y-1">
        <!-- Section Header -->
        <div 
          v-if="!appStore.sidebarCollapsed" 
          class="px-2 pt-2 pb-1 text-[10px] font-bold tracking-wider uppercase text-slate-400 select-none"
        >
          {{ section.title }}
        </div>
        <div 
          v-else-if="sIdx > 0" 
          class="my-2 mx-auto w-6 border-t border-slate-200 dark:border-slate-800"
        ></div>

        <!-- Section Items -->
        <div class="space-y-0.5">
          <template v-for="item in section.items" :key="item.path || item.label">
            <!-- Non-expandable Menu Item -->
            <router-link
              v-if="!item.children"
              :to="item.path!"
              class="group flex items-center rounded-xl text-xs font-medium transition-all duration-150 relative cursor-pointer"
              :class="[
                isRouteActive(item.path!) 
                  ? 'bg-blue-50/90 text-blue-700 dark:bg-blue-950/40 dark:text-blue-300 font-semibold shadow-xs' 
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100/80 dark:text-slate-400 dark:hover:text-slate-100 dark:hover:bg-slate-800/60',
                appStore.sidebarCollapsed ? 'justify-center p-2.5 mx-auto w-10 h-10' : 'gap-3 px-3 py-2'
              ]"
              :title="appStore.sidebarCollapsed ? item.label : undefined"
            >
              <component 
                :is="item.icon" 
                class="w-4.5 h-4.5 shrink-0 transition-colors" 
                :class="isRouteActive(item.path!) ? 'text-blue-600 dark:text-blue-400' : 'text-slate-400 group-hover:text-slate-600 dark:group-hover:text-slate-200'"
              />
              <span v-if="!appStore.sidebarCollapsed" class="truncate">{{ item.label }}</span>
              <span 
                v-if="!appStore.sidebarCollapsed && item.badge" 
                class="ml-auto text-[10px] font-semibold px-1.5 py-0.2 rounded-full"
                :class="item.badgeClass || 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300'"
              >
                {{ item.badge }}
              </span>
            </router-link>

            <!-- Expandable Submenu Item -->
            <div v-else class="space-y-0.5">
              <button
                type="button"
                @click="toggleSubmenu(item.label)"
                class="w-full group flex items-center rounded-xl text-xs font-medium transition-all duration-150 relative cursor-pointer text-left"
                :class="[
                  isSubmenuActive(item) 
                    ? 'text-blue-700 dark:text-blue-300 font-semibold' 
                    : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100/80 dark:text-slate-400 dark:hover:text-slate-100 dark:hover:bg-slate-800/60',
                  appStore.sidebarCollapsed ? 'justify-center p-2.5 mx-auto w-10 h-10' : 'gap-3 px-3 py-2 justify-between'
                ]"
                :title="appStore.sidebarCollapsed ? item.label : undefined"
              >
                <div class="flex items-center gap-3 min-w-0">
                  <component 
                    :is="item.icon" 
                    class="w-4.5 h-4.5 shrink-0 text-slate-400 group-hover:text-slate-600 dark:group-hover:text-slate-200" 
                  />
                  <span v-if="!appStore.sidebarCollapsed" class="truncate">{{ item.label }}</span>
                </div>
                <ChevronDown 
                  v-if="!appStore.sidebarCollapsed" 
                  class="w-3.5 h-3.5 text-slate-400 transition-transform duration-200"
                  :class="expandedSubmenus.has(item.label) ? 'rotate-180' : ''"
                />
              </button>

              <!-- Submenu Children (when Expanded) -->
              <div 
                v-if="!appStore.sidebarCollapsed && expandedSubmenus.has(item.label)" 
                class="pl-7 pr-1 space-y-0.5 border-l border-slate-200/80 dark:border-slate-800 ml-4 py-1"
              >
                <router-link
                  v-for="sub in item.children"
                  :key="sub.path"
                  :to="sub.path"
                  class="flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-[11px] font-medium transition-colors"
                  :class="isRouteActive(sub.path)
                    ? 'text-blue-600 dark:text-blue-400 font-semibold bg-blue-50/60 dark:bg-blue-950/30'
                    : 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-slate-100 hover:bg-slate-100/60 dark:hover:bg-slate-800/40'"
                >
                  <span class="w-1 h-1 rounded-full" :class="isRouteActive(sub.path) ? 'bg-blue-600 dark:bg-blue-400' : 'bg-slate-300 dark:bg-slate-700'"></span>
                  <span class="truncate">{{ sub.label }}</span>
                </router-link>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- Bottom Actions: Branch Status & Theme Quick Toggle -->
    <div 
      class="p-3 border-t space-y-2 shrink-0 transition-all"
      style="border-color: var(--border-color); background-color: var(--bg-sidebar);"
    >
      <!-- Outlet Status Card -->
      <div 
        v-if="!appStore.sidebarCollapsed"
        class="px-2.5 py-1.5 rounded-xl border flex items-center justify-between text-xs"
        style="border-color: var(--border-color); background-color: var(--bg-content);"
      >
        <div class="flex items-center gap-2 min-w-0">
          <span class="relative flex h-2 w-2 shrink-0">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
          <span class="text-[11px] font-medium truncate text-slate-600 dark:text-slate-300">Senopati Branch</span>
        </div>
        <span class="text-[10px] font-bold text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 px-1.5 py-0.5 rounded uppercase">
          Online
        </span>
      </div>

      <!-- Quick Expand / Collapse hint when collapsed -->
      <div v-else class="flex justify-center">
        <button
          type="button"
          @click="appStore.toggleSidebar"
          class="w-9 h-9 rounded-xl flex items-center justify-center transition-colors cursor-pointer border text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 border-slate-200 dark:border-slate-800"
          title="Perluas Sidebar (Ctrl+B)"
        >
          <PanelLeftOpen class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  Coffee,
  PanelLeftClose,
  PanelLeftOpen,
  Search,
  X,
  ChevronDown,
  LayoutDashboard,
  ShoppingCart,
  Layers,
  MonitorDot,
  History,
  Clock,
  Package,
  FileSpreadsheet,
  ClipboardCheck,
  FileText,
  Users,
  CalendarCheck,
  CalendarDays,
  FileCheck2,
  Banknote,
  DollarSign,
  BookOpen,
  BarChart3,
  SlidersHorizontal,
  ShieldCheck,
  Building2,
  UtensilsCrossed
} from 'lucide-vue-next'
import { useAppStore } from '@/stores/app.store'

const appStore = useAppStore()
const route = useRoute()

const searchQuery = ref('')
const expandedSubmenus = ref<Set<string>>(new Set(['Operasional POS', 'Logistik & Stok', 'HR & Kepegawaian']))

interface NavItem {
  label: string
  path?: string
  icon: any
  badge?: string
  badgeClass?: string
  children?: { label: string; path: string }[]
}

interface NavSection {
  title: string
  items: NavItem[]
}

const navSections: NavSection[] = [
  {
    title: 'Utama',
    items: [
      { label: 'Overview Dashboard', path: '/dashboard', icon: LayoutDashboard }
    ]
  },
  {
    title: 'Operasional POS',
    items: [
      { label: 'Point of Sale (POS)', path: '/pos', icon: ShoppingCart, badge: 'Live', badgeClass: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300' },
      { label: 'Manajemen Meja', path: '/pos/tables', icon: Layers },
      { label: 'Kitchen Display (KDS)', path: '/pos/kds', icon: MonitorDot },
      { label: 'Riwayat Transaksi', path: '/pos/transactions', icon: History },
      { label: 'Sesi Kasir & Shift', path: '/pos/shifts', icon: Clock }
    ]
  },
  {
    title: 'Logistik & Stok',
    items: [
      { label: 'Inventaris Stok', path: '/inventory', icon: Package },
      { label: 'Kartu Stok Mutasi', path: '/inventory/stock-card', icon: FileSpreadsheet },
      { label: 'Stock Opname', path: '/inventory/opname', icon: ClipboardCheck },
      { label: 'Purchase Orders (PO)', path: '/inventory/po', icon: FileText }
    ]
  },
  {
    title: 'HR & Kepegawaian',
    items: [
      { label: 'Daftar Karyawan', path: '/hris', icon: Users },
      { label: 'Presensi Absensi', path: '/hris/attendance', icon: CalendarCheck },
      { label: 'Jadwal Roster Shift', path: '/hris/shifts', icon: CalendarDays },
      { label: 'Manajemen Cuti', path: '/hris/leaves', icon: FileCheck2 },
      { label: 'Penggajian (Payroll)', path: '/hris/payroll', icon: Banknote }
    ]
  },
  {
    title: 'Keuangan & Akuntansi',
    items: [
      { label: 'Ikhtisar Keuangan', path: '/finance', icon: DollarSign },
      { label: 'Buku Jurnal Umum', path: '/finance/journal', icon: BookOpen }
    ]
  },
  {
    title: 'Analitik & Master',
    items: [
      { label: 'Laporan & Analitik', path: '/reports', icon: BarChart3 },
      { label: 'Master Data Hub', path: '/settings/master', icon: SlidersHorizontal },
      { label: 'Akses & Peran', path: '/settings/roles', icon: ShieldCheck },
      { label: 'Profil Usaha', path: '/settings/company', icon: Building2 },
      { label: 'Katalog Menu & Resep', path: '/menu/products', icon: UtensilsCrossed }
    ]
  }
]

const filteredSections = computed(() => {
  if (!searchQuery.value.trim()) return navSections
  const query = searchQuery.value.toLowerCase().trim()

  return navSections
    .map(section => {
      const matchedItems = section.items.filter(item => {
        if (item.label.toLowerCase().includes(query)) return true
        if (item.children?.some(c => c.label.toLowerCase().includes(query))) return true
        return false
      })
      return {
        ...section,
        items: matchedItems
      }
    })
    .filter(section => section.items.length > 0)
})

const isRouteActive = (path: string) => {
  if (path === '/dashboard') return route.path === '/dashboard' || route.path === '/'
  return route.path === path || route.path.startsWith(`${path}/`)
}

const isSubmenuActive = (item: NavItem) => {
  if (!item.children) return false
  return item.children.some(child => isRouteActive(child.path))
}

const toggleSubmenu = (label: string) => {
  if (expandedSubmenus.value.has(label)) {
    expandedSubmenus.value.delete(label)
  } else {
    expandedSubmenus.value.add(label)
  }
}

// Global shortcut handler (Ctrl+B)
const handleKeyDown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'b') {
    e.preventDefault()
    appStore.toggleSidebar()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
.sidebar-v2-scroll::-webkit-scrollbar {
  width: 4px;
}
.sidebar-v2-scroll::-webkit-scrollbar-track {
  background: transparent;
}
.sidebar-v2-scroll::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.3);
  border-radius: 4px;
}
.sidebar-v2-scroll::-webkit-scrollbar-thumb:hover {
  background: rgba(148, 163, 184, 0.5);
}
</style>

