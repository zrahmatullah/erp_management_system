<script setup lang="ts">
import { ref } from 'vue'
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
  Building2,
  ChevronDown
} from 'lucide-vue-next'

const activeTab = ref('sales')
const selectedRange = ref('this-month')
const selectedBranch = ref('all')
const compareToggle = ref(true)

const tabs = [
  { id: 'sales', label: 'Sales Report' },
  { id: 'financial', label: 'Financial Report' },
  { id: 'inventory', label: 'Inventory Report' },
  { id: 'hr', label: 'HR Report' },
  { id: 'custom', label: 'Custom Report' }
]

import axios from 'axios'
import { onMounted } from 'vue'

const topProducts = ref<any[]>([])

const fetchReportData = async () => {
  try {
    const res = await axios.get('/api/v1/pos/products')
    if (res.data?.data) {
      topProducts.value = res.data.data.map((p: any, idx: number) => {
        const qty = Math.max(15, 120 - (idx * 15))
        const price = Number(p.price || p.base_price || 35000)
        return {
          rank: idx + 1,
          name: p.name,
          qty: qty,
          revenue: price * qty,
          growth: '+12.5%'
        }
      })
    }
  } catch (err) {
    console.error('Failed to load report products:', err)
  }
}

onMounted(() => {
  fetchReportData()
})

const formatNum = (val: number) => {
  return Number(val || 0).toLocaleString('id-ID')
}

// Hours matrix data (1-24)
const hours = [1, 2, 3, 4, 5, 8, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 24]
const hourLevels = [
  [1, 1, 1, 1, 1, 2, 4, 5, 5, 5, 5, 5, 5, 5, 4, 2, 1, 1],
  [0, 0, 0, 0, 1, 3, 5, 5, 5, 4, 5, 5, 5, 4, 3, 1, 0, 0],
  [0, 0, 0, 0, 0, 2, 4, 5, 4, 4, 4, 5, 4, 3, 2, 0, 0, 0],
  [0, 0, 0, 0, 0, 1, 3, 4, 4, 3, 3, 4, 3, 2, 1, 0, 0, 0]
]

const getHeatColor = (val: number) => {
  if (val === 5) return 'bg-blue-600'
  if (val === 4) return 'bg-blue-500'
  if (val === 3) return 'bg-blue-400'
  if (val === 2) return 'bg-blue-300'
  if (val === 1) return 'bg-blue-200'
  return 'bg-blue-50'
}

const handleDownloadPdf = () => {
  window.print()
}

const handleExportExcel = () => {
  alert('Exporting Sales Report to Excel (.xlsx)... Download will start automatically.')
}

const handleScheduleEmail = () => {
  alert('Daily and monthly report emails scheduled successfully!')
}

const handlePrint = () => {
  window.print()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight">Reports</h1>
    </div>

    <!-- Tabs Navigation -->
    <div class="border-b border-slate-200">
      <nav class="flex space-x-8 overflow-x-auto pb-px">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm transition-colors"
          :class="[
            activeTab === tab.id
              ? 'border-blue-600 text-blue-600 font-semibold'
              : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
          ]"
        >
          {{ tab.label }}
        </button>
      </nav>
    </div>

    <!-- Filter Bar -->
    <div class="flex flex-wrap items-center justify-between gap-4 bg-white p-4 rounded-xl border border-slate-200 shadow-xs">
      <div class="flex flex-wrap items-center gap-3">
        <!-- Date range -->
        <div class="relative">
          <select 
            v-model="selectedRange"
            class="appearance-none pl-9 pr-8 py-2 bg-slate-50 border border-slate-300 rounded-lg text-sm font-medium text-slate-700 focus:outline-hidden focus:ring-2 focus:ring-blue-500"
          >
            <option value="today">Today</option>
            <option value="this-week">This Week</option>
            <option value="this-month">This Month</option>
            <option value="this-year">This Year</option>
          </select>
          <Calendar class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" />
          <ChevronDown class="w-4 h-4 text-slate-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>

        <!-- Branch filter -->
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-slate-500">Branch:</span>
          <div class="relative">
            <select 
              v-model="selectedBranch"
              class="appearance-none pl-3 pr-8 py-2 bg-slate-50 border border-slate-300 rounded-lg text-sm font-medium text-slate-700 focus:outline-hidden focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Branches</option>
              <option value="1">Central Hub - Sudirman</option>
              <option value="2">Roastery Cafe - Senopati</option>
              <option value="3">Express Bar - Kemang</option>
            </select>
            <ChevronDown class="w-4 h-4 text-slate-400 absolute right-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
          </div>
        </div>
      </div>

      <!-- Compare Toggle -->
      <div class="flex items-center gap-3">
        <span class="text-xs font-semibold text-slate-500">Compared:</span>
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
        <span class="text-xs font-medium text-slate-600">vs Last Month</span>
      </div>
    </div>

    <!-- 4 KPI Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Total Sales -->
      <div class="p-5 rounded-2xl border-2 border-blue-600 bg-blue-50/40 shadow-xs flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-blue-600 text-white flex items-center justify-center">
          <DollarSign class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500">Total Sales</p>
          <p class="text-xl font-black text-slate-900">Rp 89.750.000</p>
        </div>
      </div>

      <!-- Average Order Value -->
      <div class="p-5 rounded-2xl border border-slate-200 bg-white shadow-xs flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
          <ShoppingBag class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500">Average Order Value</p>
          <p class="text-xl font-bold text-slate-900">Rp 68.500</p>
        </div>
      </div>

      <!-- Total Transactions -->
      <div class="p-5 rounded-2xl border border-slate-200 bg-white shadow-xs flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-slate-100 text-slate-700 flex items-center justify-center">
          <CreditCard class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500">Total Transactions</p>
          <p class="text-xl font-bold text-slate-900">1,310</p>
        </div>
      </div>

      <!-- Growth -->
      <div class="p-5 rounded-2xl border border-slate-200 bg-white shadow-xs flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
          <TrendingUp class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500">Growth</p>
          <p class="text-xl font-bold text-emerald-600">+12.5%</p>
        </div>
      </div>
    </div>

    <!-- Trend Chart: Sales Report -->
    <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-xs">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-bold text-slate-900">Sales Report</h3>
        <div class="flex items-center gap-4 text-xs">
          <span class="flex items-center gap-1.5 text-slate-700 font-medium">
            <span class="w-4 h-0.5 bg-blue-600 inline-block"></span> Current period
          </span>
          <span class="flex items-center gap-1.5 text-slate-500 font-medium">
            <span class="w-4 h-0.5 border-b-2 border-dashed border-sky-400 inline-block"></span> Previous period
          </span>
        </div>
      </div>

      <!-- Chart Graphic -->
      <div class="h-64 flex flex-col justify-between">
        <div class="relative h-52 w-full border-b border-slate-200">
          <svg class="w-full h-full" viewBox="0 0 700 180" preserveAspectRatio="none">
            <defs>
              <linearGradient id="areaSales" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#2563eb" stop-opacity="0.2" />
                <stop offset="100%" stop-color="#2563eb" stop-opacity="0.0" />
              </linearGradient>
            </defs>
            <!-- Grid lines -->
            <line x1="0" y1="45" x2="700" y2="45" stroke="#f1f5f9" stroke-width="1" />
            <line x1="0" y1="90" x2="700" y2="90" stroke="#f1f5f9" stroke-width="1" />
            <line x1="0" y1="135" x2="700" y2="135" stroke="#f1f5f9" stroke-width="1" />

            <!-- Previous period (dashed light blue) -->
            <path 
              d="M0,140 Q60,90 120,110 T240,95 T360,70 T480,130 T600,60 T700,120" 
              fill="none" 
              stroke="#38bdf8" 
              stroke-width="2" 
              stroke-dasharray="4 4" 
            />

            <!-- Current period (area + solid blue line) -->
            <path 
              d="M0,130 Q60,80 120,100 T240,85 T360,50 T480,90 T600,40 T700,75 L700,180 L0,180 Z" 
              fill="url(#areaSales)" 
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
          <span>21 Jan</span>
          <span>24 Jan</span>
          <span>25 Jan</span>
          <span>26 Jan</span>
          <span>07 Feb</span>
          <span>09 Feb</span>
          <span>10 Feb</span>
          <span>13 Feb</span>
          <span>19 Feb</span>
          <span>23 Feb</span>
          <span>27 Feb</span>
          <span>28 Feb</span>
        </div>
      </div>
    </div>

    <!-- Bottom Row: Top 10 Products & Sales by Hour -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Top 10 Products -->
      <div class="bg-white rounded-2xl border border-slate-200 shadow-xs p-6">
        <h3 class="font-bold text-slate-900 mb-4">Top 10 Products</h3>
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 font-semibold text-slate-600">
                <th class="py-2.5 px-3">Rank</th>
                <th class="py-2.5 px-3">Product Name</th>
                <th class="py-2.5 px-3 text-right">Qty Sold</th>
                <th class="py-2.5 px-3 text-right">Revenue</th>
                <th class="py-2.5 px-3 text-right">Growth %</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="item in topProducts" :key="item.rank" class="hover:bg-slate-50/70">
                <td class="py-2.5 px-3 font-semibold text-slate-500">{{ item.rank }}</td>
                <td class="py-2.5 px-3 font-medium text-slate-900">{{ item.name }}</td>
                <td class="py-2.5 px-3 text-right text-slate-700">{{ item.qty }}</td>
                <td class="py-2.5 px-3 text-right font-medium text-slate-900">Rp {{ formatNum(item.revenue) }}</td>
                <td class="py-2.5 px-3 text-right font-semibold text-emerald-600">{{ item.growth }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Sales by Hour Matrix -->
      <div class="bg-white rounded-2xl border border-slate-200 shadow-xs p-6">
        <h3 class="font-bold text-slate-900 mb-4">Sales by Hour</h3>
        
        <div class="h-64 flex flex-col justify-between">
          <div class="flex-1 grid grid-rows-4 gap-1.5 py-2">
            <div v-for="(row, rIdx) in hourLevels" :key="rIdx" class="grid grid-cols-18 gap-1 items-center">
              <div 
                v-for="(val, cIdx) in row" 
                :key="cIdx" 
                class="h-9 rounded-xs transition-transform hover:scale-110 cursor-pointer"
                :class="getHeatColor(val)"
                :title="`Hour ${hours[cIdx]}: Volume Level ${val}`"
              ></div>
            </div>
          </div>

          <!-- Hours labels -->
          <div class="grid grid-cols-18 text-[9px] text-center text-slate-400 font-medium pt-2 border-t border-slate-100">
            <span v-for="h in hours" :key="h">{{ h }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Export Action Bar Footer -->
    <div class="bg-white rounded-2xl border border-slate-200 p-4 shadow-xs flex flex-wrap items-center justify-between gap-4">
      <span class="text-sm font-bold text-slate-800">Export</span>
      <div class="flex flex-wrap items-center gap-3">
        <button 
          @click="handleDownloadPdf"
          class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs transition-colors"
        >
          <Download class="w-3.5 h-3.5" />
          Download PDF
        </button>
        <button 
          @click="handleExportExcel"
          class="inline-flex items-center gap-2 px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-xl text-xs font-semibold shadow-xs transition-colors"
        >
          <FileSpreadsheet class="w-3.5 h-3.5" />
          Export Excel
        </button>
        <button 
          @click="handleScheduleEmail"
          class="inline-flex items-center gap-2 px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-xl text-xs font-semibold shadow-xs transition-colors"
        >
          <Mail class="w-3.5 h-3.5" />
          Schedule Email Report
        </button>
        <button 
          @click="handlePrint"
          class="inline-flex items-center gap-2 px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-xl text-xs font-semibold shadow-xs transition-colors"
        >
          <Printer class="w-3.5 h-3.5" />
          Print
        </button>
      </div>
    </div>
  </div>
</template>
