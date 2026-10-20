<template>
  <nav aria-label="Breadcrumb" class="flex items-center text-xs select-none">
    <ol class="inline-flex items-center space-x-1.5 md:space-x-2">
      <!-- Home Root -->
      <li class="inline-flex items-center">
        <router-link
          to="/dashboard"
          class="inline-flex items-center gap-1.5 font-medium text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200 transition-colors"
        >
          <Home class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Workspace</span>
        </router-link>
      </li>

      <!-- Dynamic Crumbs from Store or Route -->
      <template v-if="appStore.breadcrumbs.length > 0">
        <li v-for="(crumb, idx) in appStore.breadcrumbs" :key="idx" class="flex items-center">
          <ChevronRight class="w-3.5 h-3.5 text-slate-400 dark:text-slate-600 shrink-0 mx-1" />
          <router-link
            v-if="crumb.path && idx < appStore.breadcrumbs.length - 1"
            :to="crumb.path"
            class="font-medium text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200 transition-colors truncate max-w-[120px]"
          >
            {{ crumb.label }}
          </router-link>
          <span
            v-else
            class="font-semibold text-slate-900 dark:text-slate-100 truncate max-w-[160px]"
          >
            {{ crumb.label }}
          </span>
        </li>
      </template>

      <!-- Fallback: Route-based crumb if store breadcrumbs is empty -->
      <template v-else-if="currentRouteTitle">
        <li class="flex items-center">
          <ChevronRight class="w-3.5 h-3.5 text-slate-400 dark:text-slate-600 shrink-0 mx-1" />
          <span class="font-semibold text-slate-900 dark:text-slate-100 truncate">
            {{ currentRouteTitle }}
          </span>
        </li>
      </template>
    </ol>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Home, ChevronRight } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app.store'

const appStore = useAppStore()
const route = useRoute()

const currentRouteTitle = computed(() => {
  if (route.name === 'Dashboard') return 'Overview'
  if (route.name === 'POS') return 'Point of Sale'
  if (route.name === 'TableManagement') return 'Manajemen Meja'
  if (route.name === 'KitchenDisplay') return 'Kitchen Display (KDS)'
  if (route.name === 'POSTransactions') return 'Riwayat Transaksi'
  if (route.name === 'CashierShifts') return 'Sesi Kasir & Shift'
  if (route.name === 'InventoryStock') return 'Inventaris Stok'
  if (route.name === 'StockCard') return 'Kartu Stok Mutasi'
  if (route.name === 'StockOpname') return 'Stock Opname Fisik'
  if (route.name === 'PurchaseOrders') return 'Purchase Orders (PO)'
  if (route.name === 'EmployeeList') return 'Daftar Karyawan'
  if (route.name === 'Attendance') return 'Presensi Kehadiran'
  if (route.name === 'ShiftSchedule') return 'Jadwal Roster Shift'
  if (route.name === 'LeaveManagement') return 'Manajemen Cuti'
  if (route.name === 'PayrollRun') return 'Penggajian (Payroll)'
  if (route.name === 'FinanceOverview') return 'Ikhtisar Keuangan'
  if (route.name === 'JournalEntry') return 'Buku Jurnal Umum'
  if (route.name === 'Reports') return 'Laporan & Analitik'
  if (route.name === 'MasterDataHub') return 'Master Data Hub'
  if (route.name === 'RoleManagement') return 'Akses & Peran'
  if (route.name === 'CompanyProfile') return 'Profil Usaha'
  return String(route.name || '')
})
</script>

