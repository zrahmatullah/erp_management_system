<template>
  <aside class="w-64 bg-[#0F172A] text-slate-300 flex flex-col h-screen select-none border-r border-slate-800 shrink-0">
    <!-- Brand Header -->
    <div class="h-16 flex items-center px-6 gap-3 border-b border-slate-800/80">
      <div class="w-9 h-9 rounded-xl bg-blue-600 flex items-center justify-center text-white shadow-lg shadow-blue-600/30">
        <Coffee class="w-5 h-5" />
      </div>
      <div>
        <h1 class="text-white font-bold text-base tracking-wide flex items-center gap-1.5">
          Cafe ERP
          <span class="text-[10px] uppercase font-semibold tracking-wider bg-blue-500/20 text-blue-400 px-1.5 py-0.5 rounded">Pro</span>
        </h1>
        <p class="text-xs text-slate-400">All-in-One Management</p>
      </div>
    </div>

    <!-- Navigation Menu -->
    <div class="flex-1 overflow-y-auto px-3 py-4 space-y-1.5 custom-scrollbar">
      <template v-for="item in menuItems" :key="item.path || item.label">
        <!-- Item without children -->
        <router-link
          v-if="!item.children"
          :to="item.path"
          class="flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-150"
          :class="isRouteActive(item.path) 
            ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' 
            : 'text-slate-300 hover:text-white hover:bg-slate-800/60'"
        >
          <component :is="item.icon" class="w-4 h-4 shrink-0" />
          <span>{{ item.label }}</span>
        </router-link>

        <!-- Item with children / Collapsible -->
        <div v-else class="space-y-1">
          <button
            @click="toggleGroup(item.label)"
            class="w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl text-sm font-medium text-slate-300 hover:text-white hover:bg-slate-800/60 transition-all duration-150"
            :class="isGroupActive(item) ? 'text-white' : ''"
          >
            <div class="flex items-center gap-3">
              <component :is="item.icon" class="w-4 h-4 shrink-0" :class="isGroupActive(item) ? 'text-blue-400' : ''" />
              <span>{{ item.label }}</span>
            </div>
            <ChevronDown class="w-4 h-4 transition-transform duration-200" :class="openGroups[item.label] ? 'rotate-180' : ''" />
          </button>

          <!-- Submenu Items -->
          <div v-show="openGroups[item.label] || isGroupActive(item)" class="pl-7 pr-1 space-y-1 border-l border-slate-800 ml-4 py-1">
            <router-link
              v-for="sub in item.children"
              :key="sub.path"
              :to="sub.path"
              class="flex items-center justify-between px-3 py-1.5 rounded-lg text-xs font-medium transition-all"
              :class="isRouteActive(sub.path)
                ? 'bg-blue-600/90 text-white font-semibold'
                : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/40'"
            >
              <span>{{ sub.label }}</span>
              <span v-if="sub.badge" class="px-1.5 py-0.5 rounded text-[10px] bg-blue-500/30 text-blue-300 font-semibold">{{ sub.badge }}</span>
            </router-link>
          </div>
        </div>
      </template>
    </div>

    <!-- Bottom User Section & Logout -->
    <div class="p-3 border-t border-slate-800/80 bg-slate-950/40">
      <div class="flex items-center justify-between px-2 py-2">
        <div class="flex items-center gap-2.5 overflow-hidden">
          <div class="w-8 h-8 rounded-full bg-slate-700 flex items-center justify-center text-sm font-bold text-slate-200 shrink-0">
            SA
          </div>
          <div class="overflow-hidden">
            <div class="text-xs font-semibold text-white truncate">Super Admin</div>
            <div class="text-[10px] text-slate-400 truncate">admin@cafe-erp.com</div>
          </div>
        </div>
        <button @click="handleLogout" title="Logout" class="p-1.5 rounded-lg text-slate-400 hover:text-red-400 hover:bg-slate-800 transition-colors">
          <LogOut class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  LayoutDashboard,
  UtensilsCrossed,
  Package,
  Users,
  DollarSign,
  FileBarChart,
  Settings,
  Database,
  Coffee,
  ChevronDown,
  LogOut
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth.store'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const openGroups = ref<Record<string, boolean>>({
  'POS & Pesanan': true,
  'Inventori': true,
  'HRIS & Staf': false,
  'Keuangan': false,
  'Pengaturan': true
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

const menuItems = [
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
      { label: 'Kitchen KDS', path: '/pos/kds' }
    ]
  },
  {
    label: 'Inventori',
    icon: Package,
    children: [
      { label: 'Daftar Stok', path: '/inventory' },
      { label: 'Stock Opname', path: '/inventory/opname' },
      { label: 'Purchase Order', path: '/inventory/po' }
    ]
  },
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
  },
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
  },
  {
    label: 'Pengaturan',
    icon: Settings,
    children: [
      { label: '⭐ Data Master CRUD', path: '/settings/master', badge: 'Super Admin' },
      { label: 'Peran & Matriks Izin', path: '/settings/roles' }
    ]
  }
]

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 4px;
}
</style>
