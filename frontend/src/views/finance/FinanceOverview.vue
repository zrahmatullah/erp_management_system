<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { 
  Receipt, 
  CreditCard, 
  TrendingUp, 
  Wallet, 
  Plus
} from 'lucide-vue-next'

const router = useRouter()

const activeTab = ref('overview')
const loading = ref(false)

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'journal', label: 'Journal Entries' },
  { id: 'coa', label: 'Chart of Accounts' },
  { id: 'expenses', label: 'Expenses' },
  { id: 'tax', label: 'Tax' },
  { id: 'budget', label: 'Budget' },
  { id: 'reconciliation', label: 'Bank Reconciliation' }
]

const overviewData = ref({
  total_revenue: 0,
  total_expenses: 0,
  net_profit: 0,
  cash_on_hand: 0
})

const recentJournals = ref<any[]>([])

const formatNum = (num: number) => {
  return Number(num || 0).toLocaleString('id-ID')
}

const fetchFinanceData = async () => {
  loading.value = true
  try {
    const [ovRes, jRes] = await Promise.all([
      axios.get('/api/v1/finance/overview'),
      axios.get('/api/v1/finance/journals')
    ])
    if (ovRes.data?.data) {
      overviewData.value = ovRes.data.data
    }
    if (jRes.data?.data) {
      recentJournals.value = jRes.data.data.map((j: any) => ({
        date: j.entry_date || '2026-09-04',
        ref: j.reference_no || j.entry_number || 'JV-001',
        desc: j.description || 'Jurnal Transaksi',
        debit: Number(j.total_debit ?? j.debit ?? 0),
        credit: Number(j.total_credit ?? j.credit ?? 0),
        status: j.status || 'Posted'
      }))
    }
  } catch (err) {
    console.error('Failed to load finance data:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchFinanceData()
})

const monthlyData = [
  { m: 'Jan', rev: 54, exp: 32 },
  { m: 'Feb', rev: 48, exp: 30 },
  { m: 'Mar', rev: 55, exp: 37 },
  { m: 'Apr', rev: 61, exp: 42 },
  { m: 'May', rev: 63, exp: 43 },
  { m: 'Jun', rev: 68, exp: 41 },
  { m: 'Jul', rev: 62, exp: 39 },
  { m: 'Aug', rev: 64, exp: 37 },
  { m: 'Sep', rev: 60, exp: 36 },
  { m: 'Oct', rev: 74, exp: 44 },
  { m: 'Nov', rev: 69, exp: 45 },
  { m: 'Dec', rev: 73, exp: 47 }
]
</script>

<template>
  <div class="space-y-6">
    <!-- Breadcrumb & Top Bar -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <p class="text-xs font-semibold text-slate-400 mb-1">/ Finance and Accounting</p>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight">Financial Overview</h1>
      </div>
      <div class="flex items-center gap-3">
        <button 
          @click="router.push('/finance/journal')"
          class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-sm font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          Buat Jurnal Baru
        </button>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="border-b border-slate-200">
      <nav class="flex space-x-8 overflow-x-auto pb-px">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm transition-colors cursor-pointer"
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

    <!-- 4 Key Metrics Cards (From real database) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Total Revenue -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-blue-50 flex items-center justify-center text-blue-600">
          <Receipt class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500">Total Revenue</p>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(overviewData.total_revenue) }}</p>
        </div>
      </div>

      <!-- Total Expenses -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-rose-50 flex items-center justify-center text-rose-600">
          <CreditCard class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500">Total Expenses</p>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(overviewData.total_expenses) }}</p>
        </div>
      </div>

      <!-- Net Profit -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-emerald-50 flex items-center justify-center text-emerald-600">
          <TrendingUp class="w-6 h-6" />
        </div>
        <div>
          <div class="flex items-center gap-1.5">
            <p class="text-xs font-medium text-slate-500">Net Profit</p>
          </div>
          <p class="text-xl font-bold" :class="overviewData.net_profit >= 0 ? 'text-slate-900' : 'text-rose-600'">
            Rp {{ formatNum(overviewData.net_profit) }}
          </p>
        </div>
      </div>

      <!-- Cash on Hand -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-slate-100 flex items-center justify-center text-slate-600">
          <Wallet class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs font-medium text-slate-500">Cash on Hand</p>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(overviewData.cash_on_hand) }}</p>
        </div>
      </div>
    </div>

    <!-- Middle Row: Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Revenue vs Expenses Bar Chart -->
      <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-xs">
        <div class="flex items-center justify-between mb-6">
          <h3 class="font-bold text-slate-900">Revenue vs Expenses</h3>
          <div class="flex items-center gap-4 text-xs">
            <span class="flex items-center gap-1.5 text-slate-600">
              <span class="w-3 h-3 rounded-xs bg-blue-600"></span> Revenue
            </span>
            <span class="flex items-center gap-1.5 text-slate-600">
              <span class="w-3 h-3 rounded-xs bg-slate-800"></span> Expenses
            </span>
          </div>
        </div>

        <!-- Custom SVG Bar Chart -->
        <div class="h-64 flex flex-col justify-end">
          <div class="h-52 flex items-end justify-between gap-2 pt-4 px-2 border-b border-slate-200">
            <div v-for="item in monthlyData" :key="item.m" class="flex flex-col items-center flex-1 h-full justify-end group relative">
              <div class="absolute -top-10 opacity-0 group-hover:opacity-100 transition-opacity bg-slate-900 text-white text-[10px] py-1 px-2 rounded pointer-events-none whitespace-nowrap z-10">
                Rev: {{ item.rev }}M | Exp: {{ item.exp }}M
              </div>
              <div class="w-full flex items-end justify-center gap-1">
                <div class="w-2 sm:w-3 bg-blue-600 rounded-t-xs transition-all duration-300 group-hover:bg-blue-500" :style="{ height: `${item.rev * 2.2}px` }"></div>
                <div class="w-2 sm:w-3 bg-slate-800 rounded-t-xs transition-all duration-300 group-hover:bg-slate-700" :style="{ height: `${item.exp * 2.2}px` }"></div>
              </div>
              <span class="text-[10px] text-slate-400 mt-2 font-medium">{{ item.m }}</span>
            </div>
          </div>
          <div class="flex justify-between text-[10px] text-slate-400 pt-1">
            <span>0</span>
            <span>20M</span>
            <span>40M</span>
            <span>60M</span>
            <span>80M</span>
          </div>
        </div>
      </div>

      <!-- Cash Flow Trend -->
      <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-xs flex flex-col justify-between">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-bold text-slate-900">Cash Flow Trend</h3>
          <div class="flex items-center gap-3 text-xs">
            <span class="flex items-center gap-1.5 text-slate-600"><span class="w-3 h-0.5 bg-sky-400"></span> Inflow</span>
            <span class="flex items-center gap-1.5 text-slate-600"><span class="w-3 h-0.5 bg-blue-600"></span> Net Flow</span>
          </div>
        </div>

        <!-- SVG Line Trend -->
        <div class="relative h-64 flex flex-col justify-end">
          <div class="h-52 w-full pt-4">
            <svg class="w-full h-full" viewBox="0 0 500 200" preserveAspectRatio="none">
              <defs>
                <linearGradient id="gradFlow" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#38bdf8" stop-opacity="0.3" />
                  <stop offset="100%" stop-color="#38bdf8" stop-opacity="0.0" />
                </linearGradient>
                <linearGradient id="gradNet" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#2563eb" stop-opacity="0.4" />
                  <stop offset="100%" stop-color="#2563eb" stop-opacity="0.0" />
                </linearGradient>
              </defs>
              <line x1="0" y1="50" x2="500" y2="50" stroke="#f1f5f9" stroke-width="1" />
              <line x1="0" y1="100" x2="500" y2="100" stroke="#f1f5f9" stroke-width="1" />
              <line x1="0" y1="150" x2="500" y2="150" stroke="#f1f5f9" stroke-width="1" />

              <path d="M0,160 Q50,110 100,125 T200,80 T300,60 T400,90 T500,50 L500,200 L0,200 Z" fill="url(#gradFlow)" />
              <path d="M0,160 Q50,110 100,125 T200,80 T300,60 T400,90 T500,50" fill="none" stroke="#38bdf8" stroke-width="2.5" />

              <path d="M0,180 Q50,140 100,150 T200,110 T300,95 T400,120 T500,80 L500,200 L0,200 Z" fill="url(#gradNet)" />
              <path d="M0,180 Q50,140 100,150 T200,110 T300,95 T400,120 T500,80" fill="none" stroke="#2563eb" stroke-width="2.5" />
            </svg>
          </div>
          <div class="flex justify-between text-[10px] text-slate-400 pt-1">
            <span>Jan</span>
            <span>Mar</span>
            <span>May</span>
            <span>Jul</span>
            <span>Sep</span>
            <span>Nov</span>
            <span>Dec</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Row: Recent Journal Entries & Expense Breakdown -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Journal Entries (2 cols) -->
      <div class="lg:col-span-2 bg-white rounded-2xl border border-slate-200 shadow-xs p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-bold text-slate-900">Recent Journal Entries</h3>
          <button @click="router.push('/finance/journal')" class="text-xs text-blue-600 font-semibold hover:underline cursor-pointer">
            View All Entries →
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
                <th class="py-3 px-3">Date</th>
                <th class="py-3 px-3">Reference</th>
                <th class="py-3 px-3">Description</th>
                <th class="py-3 px-3 text-right">Debit (Rp)</th>
                <th class="py-3 px-3 text-right">Credit (Rp)</th>
                <th class="py-3 px-3 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-xs">
              <tr v-for="(entry, idx) in recentJournals" :key="idx" class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3 px-3 text-slate-600">{{ entry.date }}</td>
                <td class="py-3 px-3 font-semibold text-slate-900">{{ entry.ref }}</td>
                <td class="py-3 px-3 text-slate-700 max-w-xs truncate">{{ entry.desc }}</td>
                <td class="py-3 px-3 text-right font-medium text-slate-900">
                  {{ entry.debit > 0 ? formatNum(entry.debit) : '-' }}
                </td>
                <td class="py-3 px-3 text-right font-medium text-slate-900">
                  {{ entry.credit > 0 ? formatNum(entry.credit) : '-' }}
                </td>
                <td class="py-3 px-3 text-center">
                  <span class="inline-block px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                    {{ entry.status }}
                  </span>
                </td>
              </tr>
              <tr v-if="recentJournals.length === 0">
                <td colspan="6" class="py-8 text-center text-slate-400 text-xs">
                  Tidak ada catatan jurnal di database server.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Expense Breakdown Donut -->
      <div class="bg-white rounded-2xl border border-slate-200 shadow-xs p-6 flex flex-col justify-between">
        <h3 class="font-bold text-slate-900 mb-2">Expense Breakdown</h3>

        <!-- Donut Visual SVG -->
        <div class="relative flex items-center justify-center my-4">
          <svg class="w-48 h-48 transform -rotate-90" viewBox="0 0 100 100">
            <circle cx="50" cy="50" r="38" fill="transparent" stroke="#1e3a8a" stroke-width="16" stroke-dasharray="107.4 238.8" stroke-dashoffset="0" />
            <circle cx="50" cy="50" r="38" fill="transparent" stroke="#2563eb" stroke-width="16" stroke-dasharray="71.6 238.8" stroke-dashoffset="-107.4" />
            <circle cx="50" cy="50" r="38" fill="transparent" stroke="#38bdf8" stroke-width="16" stroke-dasharray="23.9 238.8" stroke-dashoffset="-179" />
            <circle cx="50" cy="50" r="38" fill="transparent" stroke="#34d399" stroke-width="16" stroke-dasharray="19.1 238.8" stroke-dashoffset="-202.9" />
            <circle cx="50" cy="50" r="38" fill="transparent" stroke="#ec4899" stroke-width="16" stroke-dasharray="16.7 238.8" stroke-dashoffset="-222" />
          </svg>
          <div class="absolute text-center">
            <span class="text-xs text-slate-400 font-medium">Top Cost</span>
            <p class="text-lg font-black text-slate-900">45%</p>
            <span class="text-[10px] text-slate-500">Ingredients</span>
          </div>
        </div>

        <!-- Legend -->
        <div class="flex flex-wrap items-center justify-center gap-3 text-xs pt-2">
          <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-blue-900"></span> Ingredients 45%</span>
          <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span> Salary 30%</span>
          <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-sky-400"></span> Rent 10%</span>
          <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-emerald-400"></span> Utilities 8%</span>
          <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-pink-500"></span> Others 7%</span>
        </div>
      </div>
    </div>
  </div>
</template>
