<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { 
  DollarSign, 
  ShoppingBag, 
  CreditCard, 
  TrendingUp, 
  Download, 
  FileSpreadsheet, 
  Mail, 
  Printer, 
  Calendar, 
  ChevronDown, 
  FileBarChart, 
  Package, 
  Users, 
  AlertTriangle, 
  CheckCircle2, 
  Clock, 
  Briefcase, 
  Layers, 
  Filter, 
  RefreshCw, 
  Search, 
  Sliders, 
  Plus, 
  Trash2, 
  X, 
  ArrowUpRight, 
  ArrowDownRight,
  TrendingDown,
  Sparkles,
  Award
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const route = useRoute()
const router = useRouter()
const notifyStore = useNotificationStore()

// ----------------------------------------------------------------------------
// Active Tab & Filters
// ----------------------------------------------------------------------------
const activeTab = ref((route.query.tab as string) || 'sales')
const selectedRange = ref('this-month')
const selectedBranch = ref('all')
const compareToggle = ref(true)
const isLoading = ref(false)
const lastUpdated = ref(new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' }))

const tabs = [
  { id: 'sales', label: 'Sales Report', icon: ShoppingBag, desc: 'Tren omzet, produk terlaris & jam sibuk' },
  { id: 'financial', label: 'Financial Report', icon: DollarSign, desc: 'Laba rugi, HPP, beban operasional & arus kas' },
  { id: 'inventory', label: 'Inventory Report', icon: Package, desc: 'Valuasi stok, analisis Pareto & audit opname' },
  { id: 'hr', label: 'HR Report', icon: Users, desc: 'Presensi, produktivitas shift & beban payroll' },
  { id: 'custom', label: 'Custom Report', icon: Sliders, desc: 'Query generator & template laporan kustom' }
]

watch(activeTab, (newTab) => {
  router.replace({ query: { ...route.query, tab: newTab } })
  loadActiveTabData()
})

const branches = ref<any[]>([])

// ----------------------------------------------------------------------------
// Tab 1: Sales Report State
// ----------------------------------------------------------------------------
const salesData = ref({
  kpi: {
    total_sales: 2806100,
    avg_order_value: 51964.8,
    total_orders: 54,
    growth: '+14.5%'
  },
  trend: {
    labels: ['21 Jan', '24 Jan', '25 Jan', '26 Jan', '07 Feb', '09 Feb', '10 Feb', '13 Feb', '19 Feb', '23 Feb', '27 Feb', '28 Feb'],
    current: [140000, 195000, 250000, 168000, 310000, 225000, 335000, 280000, 390000, 255000, 420000, 510000],
    previous: [110000, 165000, 195000, 220000, 195000, 280000, 250000, 310000, 335000, 225000, 365000, 390000]
  },
  top_products: [] as any[],
  sales_by_hour: {
    hours: [1, 2, 3, 4, 5, 8, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 24],
    hour_levels: [
      [1, 1, 1, 1, 1, 2, 4, 5, 5, 5, 5, 5, 5, 5, 4, 2, 1, 1],
      [0, 0, 0, 0, 1, 3, 5, 5, 5, 4, 5, 5, 5, 4, 3, 1, 0, 0],
      [0, 0, 0, 0, 0, 2, 4, 5, 4, 4, 4, 5, 4, 3, 2, 0, 0, 0],
      [0, 0, 0, 0, 0, 1, 3, 4, 4, 3, 3, 4, 3, 2, 1, 0, 0, 0]
    ]
  }
})

// ----------------------------------------------------------------------------
// Tab 2: Financial Report State
// ----------------------------------------------------------------------------
const financialData = ref({
  kpi: {
    total_revenue: 2806100,
    total_expenses: 18360378,
    gross_profit: 1739782,
    gross_margin_pct: '62.0%',
    net_profit: -15554278,
    net_margin_pct: '-554.3%'
  },
  pnl_statement: [] as any[],
  expense_categories: [] as any[],
  cash_flow: {
    operating_cash_flow: 539782,
    investing_cash_flow: -1500000,
    financing_cash_flow: 0,
    cash_on_hand: 32500000,
    bank_balance: 68450000,
    total_liquid_cash: 100950000
  },
  tax_summary: {
    taxable_sales: 2806100,
    pb1_tax_rate: '10%',
    tax_collected: 280610,
    tax_payable: 280610,
    status: 'Siap Setor Kas Daerah (PB1 Restoran)'
  }
})

// ----------------------------------------------------------------------------
// Tab 3: Inventory Report State
// ----------------------------------------------------------------------------
const inventoryCategory = ref('all')
const inventorySearch = ref('')
const inventoryData = ref({
  kpi: {
    total_valuation: 16596950,
    total_items: 7,
    low_stock_items: 0,
    out_of_stock_items: 0,
    turnover_ratio: '8.6x / tahun'
  },
  stock_items: [] as any[],
  fast_moving: [] as any[],
  slow_moving: [] as any[],
  discrepancies: [] as any[]
})

// ----------------------------------------------------------------------------
// Tab 4: HR Report State
// ----------------------------------------------------------------------------
const hrData = ref({
  kpi: {
    total_employees: 4,
    attendance_rate_pct: '96.8%',
    total_payroll_cost: 25300000,
    total_overtime_hours: '24.5 Jam'
  },
  department_stats: [] as any[],
  attendance_breakdown: {
    hadir_ontime: 88,
    terlambat: 4,
    izin_cuti: 3,
    alpa: 0,
    total_shifts: 95
  },
  staff_performance: [] as any[]
})

// ----------------------------------------------------------------------------
// Tab 5: Custom Report State
// ----------------------------------------------------------------------------
const customReportForm = ref({
  source: 'sales',
  date_from: '2026-10-01',
  date_to: '2026-10-31',
  group_by: 'daily',
  branch_id: 'all'
})
const customReportResult = ref<any>(null)
const customTemplates = ref<any[]>([])
const isGeneratingCustom = ref(false)

// ----------------------------------------------------------------------------
// Schedule Email Modal State
// ----------------------------------------------------------------------------
const showScheduleModal = ref(false)
const scheduleList = ref<any[]>([])
const newScheduleForm = ref({
  report_type: 'sales',
  frequency: 'daily',
  time_of_day: '22:00',
  recipient_emails: 'owner@cafe-erp.com, manager@cafe-erp.com',
  format: 'pdf'
})

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------
const formatNum = (val: number | string) => {
  return Number(val || 0).toLocaleString('id-ID')
}

const getHeatColor = (val: number) => {
  if (val === 5) return 'bg-blue-600'
  if (val === 4) return 'bg-blue-500'
  if (val === 3) return 'bg-blue-400'
  if (val === 2) return 'bg-blue-300'
  if (val === 1) return 'bg-blue-200'
  return 'bg-blue-50'
}

// ----------------------------------------------------------------------------
// Data Loaders
// ----------------------------------------------------------------------------
const fetchBranches = async () => {
  try {
    const res = await axios.get('/api/v1/master/branches')
    if (res.data) {
      branches.value = Array.isArray(res.data) ? res.data : (res.data.data || [])
    }
  } catch (err) {
    console.warn('Could not load branches:', err)
  }
}

const fetchSalesReport = async () => {
  try {
    const res = await axios.get('/api/v1/reports/sales', {
      params: { range: selectedRange.value, branch_id: selectedBranch.value }
    })
    if (res.data?.data) {
      salesData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load sales report:', err)
  }
}

const fetchFinancialReport = async () => {
  try {
    const res = await axios.get('/api/v1/reports/financial', {
      params: { branch_id: selectedBranch.value }
    })
    if (res.data?.data) {
      financialData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load financial report:', err)
  }
}

const fetchInventoryReport = async () => {
  try {
    const res = await axios.get('/api/v1/reports/inventory', {
      params: { category: inventoryCategory.value, branch_id: selectedBranch.value }
    })
    if (res.data?.data) {
      inventoryData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load inventory report:', err)
  }
}

const fetchHRReport = async () => {
  try {
    const res = await axios.get('/api/v1/reports/hr', {
      params: { branch_id: selectedBranch.value }
    })
    if (res.data?.data) {
      hrData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load HR report:', err)
  }
}

const fetchCustomTemplates = async () => {
  try {
    const res = await axios.get('/api/v1/reports/templates')
    if (res.data?.data) {
      customTemplates.value = res.data.data
    }
  } catch (err) {
    console.warn('Failed to load templates:', err)
  }
}

const fetchSchedules = async () => {
  try {
    const res = await axios.get('/api/v1/reports/schedules')
    if (res.data?.data) {
      scheduleList.value = res.data.data
    }
  } catch (err) {
    console.warn('Failed to load schedules:', err)
  }
}

const generateCustomReport = async () => {
  isGeneratingCustom.value = true
  try {
    const res = await axios.post('/api/v1/reports/custom', customReportForm.value)
    if (res.data?.data) {
      customReportResult.value = res.data.data
      notifyStore.success('Laporan kustom berhasil digenerate!', 'Custom Report')
    }
  } catch (err) {
    notifyStore.error('Gagal men-generate laporan kustom', 'Error')
  } finally {
    isGeneratingCustom.value = false
  }
}

const loadPresetTemplate = (template: any) => {
  customReportForm.value.source = template.source
  if (template.config?.group_by) {
    customReportForm.value.group_by = template.config.group_by
  }
  generateCustomReport()
}

const loadActiveTabData = async () => {
  isLoading.value = true
  try {
    if (activeTab.value === 'sales') {
      await fetchSalesReport()
    } else if (activeTab.value === 'financial') {
      await fetchFinancialReport()
    } else if (activeTab.value === 'inventory') {
      await fetchInventoryReport()
    } else if (activeTab.value === 'hr') {
      await fetchHRReport()
    } else if (activeTab.value === 'custom') {
      await fetchCustomTemplates()
      if (!customReportResult.value) {
        await generateCustomReport()
      }
    }
    lastUpdated.value = new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  } finally {
    isLoading.value = false
  }
}

// ----------------------------------------------------------------------------
// Schedule Email Handlers
// ----------------------------------------------------------------------------
const handleOpenScheduleModal = async () => {
  showScheduleModal.value = true
  await fetchSchedules()
}

const handleCreateSchedule = async () => {
  try {
    const emails = newScheduleForm.value.recipient_emails
      .split(',')
      .map(e => e.trim())
      .filter(e => e.length > 0)

    const payload = {
      report_type: newScheduleForm.value.report_type,
      frequency: newScheduleForm.value.frequency,
      time_of_day: newScheduleForm.value.time_of_day,
      recipient_emails: emails,
      format: newScheduleForm.value.format
    }

    await axios.post('/api/v1/reports/schedules', payload)
    notifyStore.success('Jadwal pengiriman email otomatis berhasil dikonfigurasi!', 'Jadwal Aktif')
    await fetchSchedules()
  } catch (err) {
    notifyStore.error('Gagal menyimpan jadwal pengiriman email', 'Error')
  }
}

const handleDeleteSchedule = async (id: string) => {
  try {
    await axios.delete(`/api/v1/reports/schedules/${id}`)
    notifyStore.success('Jadwal laporan berhasil dihapus', 'Dihapus')
    await fetchSchedules()
  } catch (err) {
    notifyStore.error('Gagal menghapus jadwal', 'Error')
  }
}

// ----------------------------------------------------------------------------
// Export Actions
// ----------------------------------------------------------------------------
const handleDownloadPdf = () => {
  window.print()
}

const handleExportExcel = () => {
  let csvContent = '\uFEFF' // UTF-8 BOM for Excel
  let filename = `Cafe_ERP_${activeTab.value}_Report_${new Date().toISOString().slice(0, 10)}.csv`

  if (activeTab.value === 'sales') {
    csvContent += 'Rank,Nama Produk,Jumlah Terjual (Qty),Total Pendapatan (Rp),Pertumbuhan\n'
    salesData.value.top_products.forEach((p: any) => {
      csvContent += `"${p.rank}","${p.name}","${p.qty}","${p.revenue}","${p.growth}"\n`
    })
  } else if (activeTab.value === 'financial') {
    csvContent += 'Kategori,Deskripsi Item,Jumlah (Rp),Persentase\n'
    financialData.value.pnl_statement.forEach((cat: any) => {
      csvContent += `"[${cat.category}]","","${cat.subtotal}","100%"\n`
      cat.items.forEach((item: any) => {
        csvContent += `"","${item.name}","${item.amount}","${item.pct}"\n`
      })
    })
  } else if (activeTab.value === 'inventory') {
    csvContent += 'SKU,Nama Bahan Baku,Kategori,Stok Tersedia,Satuan,HPP Rata-rata (Rp),Nilai Aset (Rp),Status\n'
    inventoryData.value.stock_items.forEach((item: any) => {
      csvContent += `"${item.sku}","${item.name}","${item.category}","${item.quantity}","${item.uom}","${item.average_cost}","${item.valuation}","${item.status}"\n`
    })
  } else if (activeTab.value === 'hr') {
    csvContent += 'NIK,Nama Karyawan,Departemen,Posisi,Total Shift,Tingkat Hadir,Jam Lembur,Gaji Pokok (Rp),Status\n'
    hrData.value.staff_performance.forEach((s: any) => {
      csvContent += `"${s.nik}","${s.name}","${s.department}","${s.position}","${s.shifts}","${s.attendance_rate}","${s.overtime_hours}","${s.basic_salary}","${s.status}"\n`
    })
  } else if (activeTab.value === 'custom' && customReportResult.value) {
    const headers = customReportResult.value.headers.map((h: any) => h.label).join(',')
    csvContent += headers + '\n'
    customReportResult.value.rows.forEach((r: any) => {
      const row = customReportResult.value.headers.map((h: any) => `"${r[h.key] || ''}"`).join(',')
      csvContent += row + '\n'
    })
  }

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  notifyStore.success(`File ${filename} berhasil diekspor dan diunduh!`, 'Ekspor Berhasil')
}

const handlePrint = () => {
  window.print()
}

// Filtered stock list for Inventory Tab
const filteredStockItems = computed(() => {
  let list = inventoryData.value.stock_items || []
  if (inventoryCategory.value !== 'all') {
    list = list.filter((i: any) => i.category === inventoryCategory.value)
  }
  if (inventorySearch.value.trim()) {
    const q = inventorySearch.value.toLowerCase()
    list = list.filter((i: any) => i.name.toLowerCase().includes(q) || i.sku.toLowerCase().includes(q))
  }
  return list
})

onMounted(async () => {
  await fetchBranches()
  await loadActiveTabData()
})
</script>

<template>
  <div class="space-y-6 pb-12">
    <!-- Top Action Banner (No Print) -->
    <div class="flex flex-wrap items-center justify-between gap-4 bg-white p-4 rounded-2xl border border-slate-200/80 shadow-xs no-print">
      <div class="flex items-center gap-3">
        <div class="w-11 h-11 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-600 text-white flex items-center justify-center shadow-sm">
          <FileBarChart class="w-6 h-6" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h1 class="text-xl font-bold text-slate-900 tracking-tight">Reports & Business Intelligence</h1>
            <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-bold bg-blue-50 text-blue-700 border border-blue-200/60">
              <Sparkles class="w-3 h-3 text-blue-600" />
              Live Connected
            </span>
          </div>
          <p class="text-xs text-slate-500 font-medium">Sinkronisasi data POS penjualan, keuangan buku besar, inventori gudang, dan absensi HRIS</p>
        </div>
      </div>

      <div class="flex items-center gap-2 text-xs">
        <span class="text-slate-400 font-medium">Terakhir sinkron: {{ lastUpdated }}</span>
        <button
          @click="loadActiveTabData"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-slate-200 bg-slate-50 hover:bg-slate-100 text-slate-700 font-semibold transition-colors cursor-pointer"
          :disabled="isLoading"
        >
          <RefreshCw class="w-3.5 h-3.5 text-slate-500" :class="{ 'animate-spin': isLoading }" />
          <span>Refresh</span>
        </button>
      </div>
    </div>

    <!-- Print Only Header -->
    <div class="hidden print:block border-b-2 border-slate-800 pb-4 mb-6">
      <div class="flex justify-between items-start">
        <div>
          <h1 class="text-2xl font-black text-slate-900 uppercase tracking-wider">Cafe ERP System</h1>
          <p class="text-sm font-semibold text-slate-600 mt-1">Laporan Resmi Manajemen & Audit Operasional</p>
          <p class="text-xs text-slate-500">Modul: {{ tabs.find(t => t.id === activeTab)?.label }} | Dicetak: {{ new Date().toLocaleString('id-ID') }}</p>
        </div>
        <div class="text-right">
          <p class="text-xs font-bold text-slate-800">Cabang: {{ selectedBranch === 'all' ? 'Semua Cabang (All Branches)' : selectedBranch }}</p>
          <p class="text-xs text-slate-600">Periode: {{ selectedRange }}</p>
        </div>
      </div>
    </div>

    <!-- 5 Modern Tabs Navigation (No Print) -->
    <div class="bg-white p-2 rounded-2xl border border-slate-200/80 shadow-xs no-print">
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-2">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="flex items-center gap-3 p-3 rounded-xl transition-all cursor-pointer text-left group"
          :class="[
            activeTab === tab.id
              ? 'bg-blue-600 text-white shadow-sm ring-2 ring-blue-600/30'
              : 'hover:bg-slate-50 text-slate-700 border border-transparent hover:border-slate-200/60'
          ]"
        >
          <div 
            class="w-9 h-9 rounded-lg flex items-center justify-center shrink-0 transition-colors"
            :class="activeTab === tab.id ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-600 group-hover:bg-blue-50 group-hover:text-blue-600'"
          >
            <component :is="tab.icon" class="w-4.5 h-4.5" />
          </div>
          <div class="min-w-0">
            <p class="text-xs font-bold truncate leading-tight" :class="activeTab === tab.id ? 'text-white' : 'text-slate-900'">
              {{ tab.label }}
            </p>
            <p class="text-[10px] truncate mt-0.5" :class="activeTab === tab.id ? 'text-blue-100' : 'text-slate-400'">
              {{ tab.desc }}
            </p>
          </div>
        </button>
      </div>
    </div>

    <!-- Filter Control Bar (No Print) -->
    <div class="flex flex-wrap items-center justify-between gap-4 bg-white p-4 rounded-2xl border border-slate-200/80 shadow-xs no-print">
      <div class="flex flex-wrap items-center gap-3">
        <!-- Date range -->
        <div class="relative">
          <select 
            v-model="selectedRange"
            @change="loadActiveTabData"
            class="appearance-none pl-9 pr-8 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-semibold text-slate-700 focus:outline-hidden focus:ring-2 focus:ring-blue-500 cursor-pointer shadow-2xs"
          >
            <option value="today">Today (Hari Ini)</option>
            <option value="this-week">This Week (Minggu Ini)</option>
            <option value="this-month">This Month (Bulan Ini)</option>
            <option value="this-year">This Year (Tahun Ini)</option>
          </select>
          <Calendar class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" />
          <ChevronDown class="w-4 h-4 text-slate-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>

        <!-- Branch filter -->
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-slate-500">Cabang:</span>
          <div class="relative">
            <select 
              v-model="selectedBranch"
              @change="loadActiveTabData"
              class="appearance-none pl-3 pr-8 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-semibold text-slate-700 focus:outline-hidden focus:ring-2 focus:ring-blue-500 cursor-pointer shadow-2xs"
            >
              <option value="all">Semua Cabang (All Branches)</option>
              <option v-for="b in branches" :key="b.id" :value="b.id">{{ b.name }}</option>
            </select>
            <ChevronDown class="w-4 h-4 text-slate-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
          </div>
        </div>

        <!-- Extra Category filter for Inventory -->
        <div v-if="activeTab === 'inventory'" class="flex items-center gap-2">
          <span class="text-xs font-bold text-slate-500">Kategori:</span>
          <div class="relative">
            <select 
              v-model="inventoryCategory"
              @change="loadActiveTabData"
              class="appearance-none pl-3 pr-8 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-semibold text-slate-700 focus:outline-hidden focus:ring-2 focus:ring-blue-500 cursor-pointer shadow-2xs"
            >
              <option value="all">Semua Kategori Bahan</option>
              <option value="coffee_beans">Bahan Kopi (Beans)</option>
              <option value="dairy">Dairy & Susu Segar</option>
              <option value="syrup">Sirup & Flavour</option>
              <option value="packaging">Packaging & Cup</option>
              <option value="dry_goods">Dry Goods & Powder</option>
            </select>
            <ChevronDown class="w-4 h-4 text-slate-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
          </div>
        </div>
      </div>

      <!-- Compare Toggle (Sales Tab) -->
      <div v-if="activeTab === 'sales'" class="flex items-center gap-3">
        <span class="text-xs font-semibold text-slate-500">Bandingkan:</span>
        <button 
          @click="compareToggle = !compareToggle"
          class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-hidden"
          :class="compareToggle ? 'bg-blue-600' : 'bg-slate-200'"
        >
          <span 
            class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow-sm ring-0 transition duration-200 ease-in-out"
            :class="compareToggle ? 'translate-x-5' : 'translate-x-0'"
          />
        </button>
        <span class="text-xs font-medium text-slate-700">vs Periode Lalu</span>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- TAB 1: SALES REPORT                                                  -->
    <!-- ==================================================================== -->
    <div v-if="activeTab === 'sales'" class="space-y-6">
      <!-- 4 KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Total Sales -->
        <div class="p-5 rounded-2xl border-2 border-blue-600 bg-blue-50/50 shadow-xs flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-blue-600 text-white flex items-center justify-center shrink-0 shadow-xs">
            <DollarSign class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Sales</p>
            <p class="text-2xl font-black text-slate-900">Rp {{ formatNum(salesData.kpi.total_sales) }}</p>
          </div>
        </div>

        <!-- Average Order Value -->
        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center shrink-0">
            <ShoppingBag class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Average Order Value</p>
            <p class="text-2xl font-bold text-slate-900">Rp {{ formatNum(salesData.kpi.avg_order_value) }}</p>
          </div>
        </div>

        <!-- Total Transactions -->
        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-slate-100 text-slate-700 flex items-center justify-center shrink-0">
            <CreditCard class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Transactions</p>
            <p class="text-2xl font-bold text-slate-900">{{ salesData.kpi.total_orders }} Pesanan</p>
          </div>
        </div>

        <!-- Growth -->
        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0">
            <TrendingUp class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Growth</p>
            <p class="text-2xl font-bold text-emerald-600">{{ salesData.kpi.growth }}</p>
          </div>
        </div>
      </div>

      <!-- Trend Chart: Sales Report -->
      <div class="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div>
            <h3 class="font-bold text-slate-900 text-base">Sales Report Trend</h3>
            <p class="text-xs text-slate-500">Perbandingan pola penjualan terhadap baseline periode lalu</p>
          </div>
          <div class="flex items-center gap-4 text-xs font-semibold">
            <span class="flex items-center gap-1.5 text-blue-600">
              <span class="w-4 h-1 bg-blue-600 rounded-full inline-block"></span> Periode Berjalan
            </span>
            <span v-if="compareToggle" class="flex items-center gap-1.5 text-sky-500">
              <span class="w-4 h-0.5 border-b-2 border-dashed border-sky-400 inline-block"></span> Periode Sebelumnya
            </span>
          </div>
        </div>

        <!-- SVG Line Chart Graphic -->
        <div class="h-64 flex flex-col justify-between pt-2">
          <div class="relative h-52 w-full border-b border-slate-200">
            <svg class="w-full h-full" viewBox="0 0 700 180" preserveAspectRatio="none">
              <defs>
                <linearGradient id="areaSalesGrad2" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#2563eb" stop-opacity="0.28" />
                  <stop offset="100%" stop-color="#2563eb" stop-opacity="0.0" />
                </linearGradient>
              </defs>
              <line x1="0" y1="45" x2="700" y2="45" stroke="#f1f5f9" stroke-width="1" />
              <line x1="0" y1="90" x2="700" y2="90" stroke="#f1f5f9" stroke-width="1" />
              <line x1="0" y1="135" x2="700" y2="135" stroke="#f1f5f9" stroke-width="1" />

              <!-- Previous period (dashed) -->
              <path 
                v-if="compareToggle"
                d="M0,140 Q60,90 120,110 T240,95 T360,70 T480,130 T600,60 T700,120" 
                fill="none" 
                stroke="#38bdf8" 
                stroke-width="2.5" 
                stroke-dasharray="5 5" 
              />

              <!-- Current period (area + solid) -->
              <path 
                d="M0,130 Q60,80 120,100 T240,85 T360,50 T480,90 T600,40 T700,75 L700,180 L0,180 Z" 
                fill="url(#areaSalesGrad2)" 
              />
              <path 
                d="M0,130 Q60,80 120,100 T240,85 T360,50 T480,90 T600,40 T700,75" 
                fill="none" 
                stroke="#2563eb" 
                stroke-width="3" 
              />
            </svg>
          </div>
          <div class="flex justify-between text-[11px] text-slate-400 pt-2 font-medium">
            <span v-for="label in salesData.trend.labels" :key="label">{{ label }}</span>
          </div>
        </div>
      </div>

      <!-- Bottom Row: Top 10 Products & Sales by Hour -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Top 10 Products -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h3 class="font-bold text-slate-900 text-base">Top 10 Menu Terlaris</h3>
              <p class="text-xs text-slate-500">Menu kopi & makanan penyumbang omzet terbesar</p>
            </div>
            <span class="px-2.5 py-1 bg-blue-50 text-blue-700 text-xs font-bold rounded-lg border border-blue-200/60">
              Pareto Best Sellers
            </span>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse text-xs">
              <thead>
                <tr class="bg-slate-50 border-b border-slate-200 font-semibold text-slate-600">
                  <th class="py-2.5 px-3">Rank</th>
                  <th class="py-2.5 px-3">Nama Menu</th>
                  <th class="py-2.5 px-3 text-right">Terjual</th>
                  <th class="py-2.5 px-3 text-right">Omzet</th>
                  <th class="py-2.5 px-3 text-right">Tren</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="item in salesData.top_products" :key="item.rank" class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-2.5 px-3 font-bold">
                    <span 
                      v-if="item.rank <= 3" 
                      class="inline-flex items-center justify-center w-5 h-5 rounded-full text-[11px] font-black text-white"
                      :class="{
                        'bg-amber-500 shadow-2xs': item.rank === 1,
                        'bg-slate-400': item.rank === 2,
                        'bg-amber-700': item.rank === 3
                      }"
                    >
                      {{ item.rank }}
                    </span>
                    <span v-else class="text-slate-400 font-mono">{{ item.rank }}</span>
                  </td>
                  <td class="py-2.5 px-3 font-semibold text-slate-900">{{ item.name }}</td>
                  <td class="py-2.5 px-3 text-right font-medium text-slate-700">{{ item.qty }} pcs</td>
                  <td class="py-2.5 px-3 text-right font-bold text-slate-900">Rp {{ formatNum(item.revenue) }}</td>
                  <td class="py-2.5 px-3 text-right font-bold text-emerald-600">{{ item.growth }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Sales by Hour Matrix -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h3 class="font-bold text-slate-900 text-base">Sales by Hour (Matriks Jam Sibuk)</h3>
              <p class="text-xs text-slate-500">Heatmap volume pesanan untuk optimasi jadwal shift barista</p>
            </div>
            <span class="px-2.5 py-1 bg-amber-50 text-amber-700 text-xs font-bold rounded-lg border border-amber-200/60">
              Peak: 11.00 - 14.00 & 18.00 - 20.00
            </span>
          </div>
          
          <div class="h-64 flex flex-col justify-between pt-2">
            <div class="flex-1 grid grid-rows-4 gap-1.5 py-2">
              <div 
                v-for="(row, rIdx) in salesData.sales_by_hour.hour_levels" 
                :key="rIdx" 
                style="display: grid; grid-template-columns: repeat(18, minmax(0, 1fr)); gap: 4px; align-items: center;"
              >
                <div 
                  v-for="(val, cIdx) in row" 
                  :key="cIdx" 
                  class="h-9 rounded-sm transition-transform hover:scale-110 cursor-pointer shadow-2xs"
                  :class="getHeatColor(val)"
                  :title="`Jam ${salesData.sales_by_hour.hours[cIdx]}:00 - Level Aktivitas: ${val}/5`"
                ></div>
              </div>
            </div>

            <!-- Hours labels -->
            <div 
              style="display: grid; grid-template-columns: repeat(18, minmax(0, 1fr));"
              class="text-[10px] text-center text-slate-400 font-bold pt-2 border-t border-slate-100"
            >
              <span v-for="h in salesData.sales_by_hour.hours" :key="h">{{ h }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- TAB 2: FINANCIAL REPORT                                              -->
    <!-- ==================================================================== -->
    <div v-if="activeTab === 'financial'" class="space-y-6">
      <!-- 4 Financial KPIs -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Revenue</p>
          <p class="text-2xl font-black text-blue-600 mt-1">Rp {{ formatNum(financialData.kpi.total_revenue) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Pendapatan Kotor Kasir & Pesanan</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">HPP & Biaya Beban</p>
          <p class="text-2xl font-black text-rose-600 mt-1">Rp {{ formatNum(financialData.kpi.total_expenses) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">COGS Bahan Baku + Biaya Operasional</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Laba Kotor (Gross Profit)</p>
          <p class="text-2xl font-black text-emerald-600 mt-1">Rp {{ formatNum(financialData.kpi.gross_profit) }}</p>
          <p class="text-[11px] text-emerald-600 font-semibold mt-1">Margin Kotor: {{ financialData.kpi.gross_margin_pct }}</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Kas Likuid</p>
          <p class="text-2xl font-black text-slate-900 mt-1">Rp {{ formatNum(financialData.cash_flow.total_liquid_cash) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Kas Kasir + Rekening Bank Operasional</p>
        </div>
      </div>

      <!-- Income Statement (P&L) Table -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-5">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div>
            <h3 class="font-bold text-slate-900 text-lg">Laporan Laba Rugi (Income Statement / P&L)</h3>
            <p class="text-xs text-slate-500 font-medium">Rekonsiliasi akun pendapatan, HPP bahan baku, dan beban operasional</p>
          </div>
          <span class="px-3 py-1 bg-emerald-50 text-emerald-700 rounded-full text-xs font-bold border border-emerald-200">
            Jurnal Umum Terposting
          </span>
        </div>

        <div class="space-y-5">
          <div v-for="(cat, idx) in financialData.pnl_statement" :key="idx" class="space-y-2">
            <div class="flex justify-between items-center bg-slate-50/80 px-4 py-2.5 rounded-xl font-bold text-slate-800 text-sm border border-slate-100">
              <span>{{ cat.category }}</span>
              <span class="font-mono">Rp {{ formatNum(cat.subtotal) }}</span>
            </div>
            <div class="divide-y divide-slate-100 pl-4 pr-2 text-xs">
              <div v-for="(item, iIdx) in cat.items" :key="iIdx" class="py-2.5 flex justify-between items-center text-slate-600 hover:bg-slate-50/50 px-2 rounded-lg transition-colors">
                <span>{{ item.name }}</span>
                <div class="flex items-center gap-4">
                  <span class="text-slate-400 font-mono text-[11px]">{{ item.pct }}</span>
                  <span class="font-bold text-slate-900 font-mono">Rp {{ formatNum(item.amount) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Cash Flow & Tax Cards -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Arus Kas -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
          <h3 class="font-bold text-slate-900 text-base">Ringkasan Arus Kas (Cash Flow)</h3>
          <div class="space-y-3 text-xs">
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-600 font-medium">Arus Kas Operasional (Operating)</span>
              <span class="font-bold text-emerald-600">+ Rp {{ formatNum(financialData.cash_flow.operating_cash_flow) }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-600 font-medium">Arus Kas Investasi (Alat & Mesin)</span>
              <span class="font-bold text-rose-600">- Rp {{ formatNum(Math.abs(financialData.cash_flow.investing_cash_flow)) }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-600 font-medium">Saldo Kas Kasir (Petty Cash on Hand)</span>
              <span class="font-bold text-slate-900">Rp {{ formatNum(financialData.cash_flow.cash_on_hand) }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-600 font-medium">Saldo Rekening Bank Operasional</span>
              <span class="font-bold text-slate-900">Rp {{ formatNum(financialData.cash_flow.bank_balance) }}</span>
            </div>
          </div>
        </div>

        <!-- Rekapitulasi Pajak PB1 Restoran -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
          <h3 class="font-bold text-slate-900 text-base">Kepatuhan Pajak (Tax Compliance PB1)</h3>
          <div class="space-y-3 text-xs">
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-600 font-medium">Dasar Pengenaan Pajak (DPP Penjualan)</span>
              <span class="font-bold text-slate-900">Rp {{ formatNum(financialData.tax_summary.taxable_sales) }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-600 font-medium">Tarif Pajak Restoran Daerah (PB1)</span>
              <span class="font-bold text-blue-600">{{ financialData.tax_summary.pb1_tax_rate }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-600 font-medium">Total Pajak Terkumpul dari Konsumen</span>
              <span class="font-bold text-emerald-600">Rp {{ formatNum(financialData.tax_summary.tax_collected) }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100 items-center">
              <span class="text-slate-600 font-medium">Status Penyetoran Pajak</span>
              <span class="font-semibold text-amber-700 bg-amber-50 px-2 py-0.5 rounded-lg border border-amber-200">
                {{ financialData.tax_summary.status }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- TAB 3: INVENTORY REPORT                                              -->
    <!-- ==================================================================== -->
    <div v-if="activeTab === 'inventory'" class="space-y-6">
      <!-- 4 Inventory KPIs -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Nilai Aset Stok</p>
          <p class="text-2xl font-black text-blue-600 mt-1">Rp {{ formatNum(inventoryData.kpi.total_valuation) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Valuasi Seluruh Bahan Gudang</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Item Terdaftar</p>
          <p class="text-2xl font-black text-slate-900 mt-1">{{ inventoryData.kpi.total_items }} SKU</p>
          <p class="text-[11px] text-slate-400 mt-1">Terpantau Sistem Kartu Stok</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Peringatan Stok Menipis</p>
          <p class="text-2xl font-black text-amber-500 mt-1">{{ inventoryData.kpi.low_stock_items }} Item</p>
          <p class="text-[11px] text-slate-400 mt-1">Mendekati Batas Minimum</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Rasio Perputaran (Turnover)</p>
          <p class="text-2xl font-black text-emerald-600 mt-1">{{ inventoryData.kpi.turnover_ratio }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Efisiensi Siklus Perputaran Bahan</p>
        </div>
      </div>

      <!-- Inventory Stock Status Table -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
        <div class="flex flex-wrap items-center justify-between gap-4 pb-2 border-b border-slate-100">
          <div>
            <h3 class="font-bold text-slate-900 text-lg">Daftar Valuasi & Status Stok Bahan Baku</h3>
            <p class="text-xs text-slate-500 font-medium">Berdasarkan data kartu stok live dan gudang aktif</p>
          </div>
          <div class="relative w-72">
            <input 
              v-model="inventorySearch"
              type="text" 
              placeholder="Cari SKU atau nama item..."
              class="w-full pl-9 pr-3 py-2 text-xs bg-slate-50 border border-slate-300 rounded-xl focus:outline-hidden focus:ring-2 focus:ring-blue-500"
            />
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 font-semibold text-slate-600">
                <th class="py-2.5 px-3">SKU</th>
                <th class="py-2.5 px-3">Nama Bahan Baku</th>
                <th class="py-2.5 px-3">Kategori</th>
                <th class="py-2.5 px-3 text-right">Stok Fisik</th>
                <th class="py-2.5 px-3 text-right">Min. Stok</th>
                <th class="py-2.5 px-3 text-right">HPP Satuan</th>
                <th class="py-2.5 px-3 text-right">Total Valuasi</th>
                <th class="py-2.5 px-3 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="item in filteredStockItems" :key="item.sku" class="hover:bg-slate-50/70 transition-colors">
                <td class="py-2.5 px-3 font-mono font-medium text-slate-500">{{ item.sku }}</td>
                <td class="py-2.5 px-3 font-semibold text-slate-900">{{ item.name }}</td>
                <td class="py-2.5 px-3 text-slate-600">{{ item.category }}</td>
                <td class="py-2.5 px-3 text-right font-bold text-slate-800">{{ item.quantity }} {{ item.uom }}</td>
                <td class="py-2.5 px-3 text-right text-slate-500">{{ item.min_stock }} {{ item.uom }}</td>
                <td class="py-2.5 px-3 text-right text-slate-700">Rp {{ formatNum(item.average_cost) }}</td>
                <td class="py-2.5 px-3 text-right font-bold text-slate-900">Rp {{ formatNum(item.valuation) }}</td>
                <td class="py-2.5 px-3 text-center">
                  <span 
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                    :class="{
                      'bg-emerald-50 text-emerald-700 border border-emerald-200': item.status === 'Aman',
                      'bg-amber-50 text-amber-700 border border-amber-200': item.status === 'Menipis',
                      'bg-rose-50 text-rose-700 border border-rose-200': item.status === 'Habis'
                    }"
                  >
                    {{ item.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pareto 80/20 & Opname Discrepancy Rows -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Pareto Analysis Fast Moving -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h3 class="font-bold text-slate-900 text-base">Analisis Pareto 80/20 (Fast Moving)</h3>
              <p class="text-xs text-slate-500">Item bahan baku dengan konsumsi tertinggi</p>
            </div>
            <span class="text-xs text-blue-600 font-bold">Prioritas Pengadaan</span>
          </div>
          <div class="divide-y divide-slate-100 text-xs">
            <div v-for="(item, idx) in inventoryData.fast_moving" :key="idx" class="py-3 flex justify-between items-center">
              <div>
                <p class="font-semibold text-slate-900">{{ item.name }}</p>
                <p class="text-[11px] text-slate-400">{{ item.category }}</p>
              </div>
              <div class="text-right">
                <span class="px-2 py-0.5 rounded-md text-[10px] font-bold bg-blue-50 text-blue-700 border border-blue-200">
                  Turnover: {{ item.turnover }}
                </span>
                <p class="text-[11px] text-slate-500 font-medium mt-1">Porsi Penjualan: {{ item.sales_share }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Selisih Opname & Waste/Spoilage -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h3 class="font-bold text-slate-900 text-base">Selisih Opname & Kerugian Waste</h3>
              <p class="text-xs text-slate-500">Laporan deviasi stok fisik vs sistem</p>
            </div>
            <span class="text-xs text-rose-600 font-bold">Audit Kerugian</span>
          </div>
          <div class="divide-y divide-slate-100 text-xs">
            <div v-for="(disc, idx) in inventoryData.discrepancies" :key="idx" class="py-3 space-y-1.5">
              <div class="flex justify-between font-semibold">
                <span class="text-slate-900">{{ disc.item_name }}</span>
                <span class="text-rose-600 font-bold">- Rp {{ formatNum(disc.loss_value) }}</span>
              </div>
              <div class="flex justify-between text-slate-500 text-[11px]">
                <span>Sistem: {{ disc.system_stock }} | Fisik: {{ disc.physical_stock }} (Selisih {{ disc.difference }})</span>
                <span class="text-slate-400 font-mono">{{ disc.date }}</span>
              </div>
              <p class="text-[11px] text-amber-800 bg-amber-50 p-2 rounded-lg border border-amber-200/60">
                Penyebab: {{ disc.reason }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- TAB 4: HR REPORT                                                     -->
    <!-- ==================================================================== -->
    <div v-if="activeTab === 'hr'" class="space-y-6">
      <!-- 4 HR KPIs -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Karyawan Aktif</p>
          <p class="text-2xl font-black text-blue-600 mt-1">{{ hrData.kpi.total_employees }} Orang</p>
          <p class="text-[11px] text-slate-400 mt-1">Staf Operasional Terdaftar</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Rata-rata Kehadiran</p>
          <p class="text-2xl font-black text-emerald-600 mt-1">{{ hrData.kpi.attendance_rate_pct }}</p>
          <p class="text-[11px] text-emerald-600 font-semibold mt-1">95 Shift Sukses Dilaksanakan</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Biaya Payroll Gaji</p>
          <p class="text-2xl font-black text-slate-900 mt-1">Rp {{ formatNum(hrData.kpi.total_payroll_cost) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Gaji Pokok + Tunjangan Staf</p>
        </div>

        <div class="p-5 rounded-2xl border border-slate-200/80 bg-white shadow-xs">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Akumulasi Lembur</p>
          <p class="text-2xl font-black text-amber-500 mt-1">{{ hrData.kpi.total_overtime_hours }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Kompensasi Jam Sibuk & Event</p>
        </div>
      </div>

      <!-- Staff Performance Roster Table -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
        <div class="flex items-center justify-between pb-2 border-b border-slate-100">
          <div>
            <h3 class="font-bold text-slate-900 text-lg">Matriks Kinerja & Kehadiran Karyawan</h3>
            <p class="text-xs text-slate-500 font-medium">Rekapitulasi shift, rasio on-time, dan biaya gaji bulan berjalan</p>
          </div>
          <span class="px-3 py-1 bg-blue-50 text-blue-700 rounded-full text-xs font-bold border border-blue-200">
            HRIS Terintegrasi
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 font-semibold text-slate-600">
                <th class="py-2.5 px-3">NIK</th>
                <th class="py-2.5 px-3">Nama Karyawan</th>
                <th class="py-2.5 px-3">Departemen</th>
                <th class="py-2.5 px-3">Jabatan</th>
                <th class="py-2.5 px-3 text-right">Shift Diambil</th>
                <th class="py-2.5 px-3 text-right">Tingkat Hadir</th>
                <th class="py-2.5 px-3 text-right">Jam Lembur</th>
                <th class="py-2.5 px-3 text-right">Gaji Pokok</th>
                <th class="py-2.5 px-3 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="staff in hrData.staff_performance" :key="staff.nik" class="hover:bg-slate-50/70 transition-colors">
                <td class="py-2.5 px-3 font-mono font-medium text-slate-500">{{ staff.nik }}</td>
                <td class="py-2.5 px-3 font-semibold text-slate-900">{{ staff.name }}</td>
                <td class="py-2.5 px-3 text-slate-600">{{ staff.department }}</td>
                <td class="py-2.5 px-3 text-slate-800">{{ staff.position }}</td>
                <td class="py-2.5 px-3 text-right font-medium text-slate-900">{{ staff.shifts }} shift</td>
                <td class="py-2.5 px-3 text-right font-bold text-emerald-600">{{ staff.attendance_rate }}</td>
                <td class="py-2.5 px-3 text-right text-slate-700">{{ staff.overtime_hours }} jam</td>
                <td class="py-2.5 px-3 text-right font-medium text-slate-900 font-mono">Rp {{ formatNum(staff.basic_salary) }}</td>
                <td class="py-2.5 px-3 text-center">
                  <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-700 border border-emerald-200">
                    {{ staff.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Department Distribution Row -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Departemen Breakdown -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
          <h3 class="font-bold text-slate-900 text-base">Distribusi Beban Gaji per Departemen</h3>
          <div class="divide-y divide-slate-100 text-xs">
            <div v-for="(dept, idx) in hrData.department_stats" :key="idx" class="py-3 flex justify-between items-center">
              <div>
                <p class="font-semibold text-slate-900">{{ dept.department }}</p>
                <p class="text-[11px] text-slate-400">Headcount: {{ dept.headcount }} staf | Kehadiran: {{ dept.avg_attendance }}</p>
              </div>
              <div class="text-right">
                <p class="font-bold text-slate-900 font-mono">Rp {{ formatNum(dept.total_payroll) }}</p>
                <p class="text-[11px] text-slate-500 font-medium">Porsi: {{ dept.share }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Presensi Status Counts -->
        <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
          <h3 class="font-bold text-slate-900 text-base">Komposisi Presensi Bulan Ini</h3>
          <div class="grid grid-cols-2 gap-3 text-center">
            <div class="p-3.5 bg-emerald-50/80 border border-emerald-200/80 rounded-2xl">
              <p class="text-xs text-emerald-700 font-semibold">Tepat Waktu (On-Time)</p>
              <p class="text-2xl font-black text-emerald-800 mt-1">{{ hrData.attendance_breakdown.hadir_ontime }} Shift</p>
            </div>
            <div class="p-3.5 bg-amber-50/80 border border-amber-200/80 rounded-2xl">
              <p class="text-xs text-amber-700 font-semibold">Terlambat (&lt; 15 Mnt)</p>
              <p class="text-2xl font-black text-amber-800 mt-1">{{ hrData.attendance_breakdown.terlambat }} Shift</p>
            </div>
            <div class="p-3.5 bg-blue-50/80 border border-blue-200/80 rounded-2xl">
              <p class="text-xs text-blue-700 font-semibold">Izin / Cuti Sah</p>
              <p class="text-2xl font-black text-blue-800 mt-1">{{ hrData.attendance_breakdown.izin_cuti }} Hari</p>
            </div>
            <div class="p-3.5 bg-slate-50 border border-slate-200 rounded-2xl">
              <p class="text-xs text-slate-600 font-semibold">Alpa / Tanpa Keterangan</p>
              <p class="text-2xl font-black text-slate-800 mt-1">{{ hrData.attendance_breakdown.alpa }} Hari</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- TAB 5: CUSTOM REPORT BUILDER                                         -->
    <!-- ==================================================================== -->
    <div v-if="activeTab === 'custom'" class="space-y-6">
      <!-- Preset Templates Quick Bar -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-5 space-y-3">
        <h3 class="font-bold text-slate-900 text-sm">Preset Template Laporan Siap Pakai</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
          <button
            v-for="tmpl in customTemplates"
            :key="tmpl.id"
            @click="loadPresetTemplate(tmpl)"
            class="p-3.5 text-left border rounded-xl transition-all cursor-pointer hover:border-blue-500 hover:shadow-xs group"
            :class="customReportForm.source === tmpl.source ? 'border-blue-500 bg-blue-50/40 ring-1 ring-blue-500' : 'border-slate-200 bg-slate-50/50'"
          >
            <p class="text-xs font-bold text-slate-900 group-hover:text-blue-600 transition-colors">{{ tmpl.name }}</p>
            <p class="text-[11px] text-slate-500 mt-1 line-clamp-2">{{ tmpl.description }}</p>
          </button>
        </div>
      </div>

      <!-- Report Generator Form -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
        <h3 class="font-bold text-slate-900 text-lg">Custom Report Builder (Konfigurasi Laporan)</h3>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <!-- Modul Sumber Data -->
          <div class="space-y-1.5">
            <label class="text-xs font-semibold text-slate-700">Sumber Data:</label>
            <select 
              v-model="customReportForm.source"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-medium text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-blue-500 cursor-pointer"
            >
              <option value="sales">Penjualan (Sales & Kasir)</option>
              <option value="financial">Keuangan (Financial & Ledger)</option>
              <option value="inventory">Inventori & Bahan Baku</option>
              <option value="hr">HRIS & Kepegawaian</option>
            </select>
          </div>

          <!-- Dari Tanggal -->
          <div class="space-y-1.5">
            <label class="text-xs font-semibold text-slate-700">Mulai Tanggal:</label>
            <input 
              v-model="customReportForm.date_from"
              type="date" 
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-medium text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- Sampai Tanggal -->
          <div class="space-y-1.5">
            <label class="text-xs font-semibold text-slate-700">Sampai Tanggal:</label>
            <input 
              v-model="customReportForm.date_to"
              type="date" 
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-medium text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <!-- Pengelompokan Data -->
          <div class="space-y-1.5">
            <label class="text-xs font-semibold text-slate-700">Pengelompokan (Group by):</label>
            <select 
              v-model="customReportForm.group_by"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-medium text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-blue-500 cursor-pointer"
            >
              <option value="none">Tanpa Pengelompokan (Detail)</option>
              <option value="daily">Harian (Per Hari)</option>
              <option value="category">Berdasarkan Kategori</option>
              <option value="status">Berdasarkan Status</option>
            </select>
          </div>
        </div>

        <div class="flex justify-end pt-2">
          <button
            @click="generateCustomReport"
            class="inline-flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold shadow-xs transition-colors cursor-pointer"
            :disabled="isGeneratingCustom"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isGeneratingCustom }" />
            <span>Generate Laporan Kustom</span>
          </button>
        </div>
      </div>

      <!-- Generated Results Data Table -->
      <div v-if="customReportResult" class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-slate-100 pb-3">
          <div>
            <h3 class="font-bold text-slate-900 text-base">Hasil Generate Laporan Kustom</h3>
            <p class="text-xs text-slate-500">Waktu pembuatan: {{ customReportResult.generated_at }} | Sumber: {{ customReportResult.source }}</p>
          </div>
          <span class="px-2.5 py-1 bg-emerald-50 text-emerald-700 font-bold text-xs rounded-lg border border-emerald-200">
            {{ customReportResult.rows.length }} Baris Data
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 font-semibold text-slate-600">
                <th v-for="h in customReportResult.headers" :key="h.key" class="py-2.5 px-3">
                  {{ h.label }}
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="(row, rIdx) in customReportResult.rows" :key="rIdx" class="hover:bg-slate-50/70">
                <td v-for="h in customReportResult.headers" :key="h.key" class="py-2.5 px-3">
                  <span v-if="typeof row[h.key] === 'number' && (h.key.includes('total') || h.key.includes('subtotal') || h.key.includes('amount') || h.key.includes('salary') || h.key.includes('debit') || h.key.includes('credit') || h.key.includes('cost') || h.key.includes('valuation'))" class="font-mono font-medium">
                    Rp {{ formatNum(row[h.key]) }}
                  </span>
                  <span v-else>{{ row[h.key] }}</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- GLOBAL EXPORT BAR FOOTER (Present on all tabs!)                      -->
    <!-- ==================================================================== -->
    <div class="bg-white rounded-2xl border border-slate-200/80 p-4 shadow-xs flex flex-wrap items-center justify-between gap-4 no-print">
      <div class="flex items-center gap-2">
        <span class="text-sm font-bold text-slate-800">Export & Distribusi Laporan:</span>
        <span class="text-xs text-slate-400 font-medium">Pilih format dokumen resmi</span>
      </div>
      <div class="flex flex-wrap items-center gap-3">
        <button 
          @click="handleDownloadPdf"
          class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <Download class="w-3.5 h-3.5" />
          Download PDF
        </button>
        <button 
          @click="handleExportExcel"
          class="inline-flex items-center gap-2 px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <FileSpreadsheet class="w-3.5 h-3.5" />
          Export Excel (.csv)
        </button>
        <button 
          @click="handleOpenScheduleModal"
          class="inline-flex items-center gap-2 px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <Mail class="w-3.5 h-3.5" />
          Schedule Email Report
        </button>
        <button 
          @click="handlePrint"
          class="inline-flex items-center gap-2 px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5" />
          Print
        </button>
      </div>
    </div>

    <!-- ==================================================================== -->
    <!-- SCHEDULE EMAIL MODAL                                                 -->
    <!-- ==================================================================== -->
    <div v-if="showScheduleModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-xs p-4 no-print">
      <div class="bg-white rounded-2xl shadow-xl border border-slate-200 w-full max-w-xl overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="p-5 border-b border-slate-100 flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div class="w-9 h-9 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center">
              <Mail class="w-5 h-5" />
            </div>
            <div>
              <h3 class="font-bold text-slate-900 text-base">Jadwal Pengiriman Laporan via Email</h3>
              <p class="text-xs text-slate-500">Kirim laporan otomatis secara berkala ke pimpinan / akuntan</p>
            </div>
          </div>
          <button @click="showScheduleModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors cursor-pointer">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-5 space-y-4 max-h-[70vh] overflow-y-auto">
          <!-- Active Schedules List -->
          <div class="space-y-2">
            <p class="text-xs font-bold text-slate-700 uppercase tracking-wider">Jadwal Aktif Saat Ini</p>
            <div class="divide-y divide-slate-100 border border-slate-200 rounded-xl overflow-hidden text-xs">
              <div v-if="scheduleList.length === 0" class="p-4 text-center text-slate-400">
                Belum ada jadwal laporan otomatis aktif.
              </div>
              <div v-for="sch in scheduleList" :key="sch.id" class="p-3 flex items-center justify-between bg-slate-50/50">
                <div>
                  <div class="flex items-center gap-2">
                    <span class="font-bold text-slate-900 uppercase text-[11px]">{{ sch.report_type }} Report</span>
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-100 text-blue-700 capitalize">{{ sch.frequency }} @ {{ sch.time_of_day }}</span>
                    <span class="px-1.5 py-0.5 rounded text-[10px] font-mono bg-slate-200 text-slate-700 uppercase">{{ sch.format }}</span>
                  </div>
                  <p class="text-[11px] text-slate-500 mt-1">Penerima: {{ (sch.recipient_emails || []).join(', ') }}</p>
                </div>
                <button @click="handleDeleteSchedule(sch.id)" class="p-1.5 text-rose-500 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer" title="Hapus jadwal">
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>

          <!-- Add Schedule Form -->
          <div class="border-t border-slate-100 pt-4 space-y-3">
            <p class="text-xs font-bold text-slate-700 uppercase tracking-wider">Tambah Jadwal Baru</p>

            <div class="grid grid-cols-2 gap-3 text-xs">
              <div>
                <label class="font-semibold text-slate-700 block mb-1">Tipe Laporan:</label>
                <select v-model="newScheduleForm.report_type" class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl">
                  <option value="sales">Sales Report</option>
                  <option value="financial">Financial Report</option>
                  <option value="inventory">Inventory Report</option>
                  <option value="hr">HR Report</option>
                  <option value="all">Semua Laporan Lengkap</option>
                </select>
              </div>

              <div>
                <label class="font-semibold text-slate-700 block mb-1">Frekuensi:</label>
                <select v-model="newScheduleForm.frequency" class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl">
                  <option value="daily">Harian (Daily)</option>
                  <option value="weekly">Mingguan (Weekly - Senin)</option>
                  <option value="monthly">Bulanan (Monthly - Tgl 1)</option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3 text-xs">
              <div>
                <label class="font-semibold text-slate-700 block mb-1">Jam Kirim:</label>
                <input v-model="newScheduleForm.time_of_day" type="time" class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl" />
              </div>

              <div>
                <label class="font-semibold text-slate-700 block mb-1">Format Lampiran:</label>
                <select v-model="newScheduleForm.format" class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl">
                  <option value="pdf">Dokumen PDF (.pdf)</option>
                  <option value="excel">Spreadsheet Excel (.csv/.xlsx)</option>
                </select>
              </div>
            </div>

            <div class="text-xs">
              <label class="font-semibold text-slate-700 block mb-1">Daftar Email Penerima (Pisahkan dengan koma):</label>
              <input 
                v-model="newScheduleForm.recipient_emails" 
                type="text" 
                placeholder="owner@cafe.com, finance@cafe.com"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl"
              />
            </div>
          </div>
        </div>

        <div class="p-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3">
          <button @click="showScheduleModal = false" class="px-4 py-2 border border-slate-300 text-slate-700 rounded-xl text-xs font-semibold hover:bg-slate-100 cursor-pointer">
            Tutup
          </button>
          <button @click="handleCreateSchedule" class="px-4 py-2 bg-blue-600 text-white rounded-xl text-xs font-semibold hover:bg-blue-700 shadow-xs flex items-center gap-1.5 cursor-pointer">
            <Plus class="w-3.5 h-3.5" />
            <span>Simpan Jadwal</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media print {
  .no-print {
    display: none !important;
  }
  body {
    background-color: white !important;
    color: black !important;
  }
}
</style>
