<template>
  <div class="reports-v2 space-y-6 animate-in fade-in duration-200">
    <!-- Page Header V2 -->
    <PageHeaderV2
      title="Laporan Bisnis & Analitik Eksekutif"
      description="Analisis performa penjualan harian, profitabilitas cabang, audit inventaris, dan efisiensi tenaga kerja."
      badge="Business Intelligence"
    >
      <template #actions>
        <select 
          v-model="selectedRange"
          class="px-3 py-2 rounded-xl border text-xs font-semibold outline-none"
          style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
        >
          <option value="today">Hari Ini</option>
          <option value="this-week">Minggu Ini</option>
          <option value="this-month">Bulan Ini (September - Oktober)</option>
          <option value="this-year">Tahun 2026</option>
        </select>

        <button
          type="button"
          @click="exportExcel"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold border transition-all flex items-center gap-1.5 cursor-pointer text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800"
          style="border-color: var(--border-color); background-color: var(--bg-primary);"
        >
          <FileSpreadsheet class="w-3.5 h-3.5 text-emerald-600" />
          <span>Export Excel</span>
        </button>

        <button
          type="button"
          @click="exportPDF"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Unduh PDF</span>
        </button>
      </template>

      <!-- Sub-tabs Slot -->
      <template #tabs>
        <div class="flex items-center gap-6 text-xs font-semibold overflow-x-auto pb-2 select-none">
          <button
            v-for="t in reportTabs"
            :key="t.id"
            type="button"
            @click="activeTab = t.id"
            class="flex items-center gap-2 pb-2 -mb-2 cursor-pointer transition-colors whitespace-nowrap"
            :class="activeTab === t.id 
              ? 'text-blue-600 dark:text-blue-400 border-b-2 border-blue-600' 
              : 'text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200'"
          >
            <component :is="t.icon" class="w-3.5 h-3.5" />
            <span>{{ t.label }}</span>
          </button>
        </div>
      </template>
    </PageHeaderV2>

    <!-- TAB 1: SALES REPORT -->
    <div v-if="activeTab === 'sales'" class="space-y-6">
      <!-- 4 Stat Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Omzet Penjualan"
          value="Rp 148.520.000"
          trend="up"
          trendValue="+14.5%"
          subtitle="Dibandingkan bulan lalu"
          color="emerald"
          :icon="DollarSign"
        />
        <StatCardV2
          title="Rata-rata Order (AOV)"
          value="Rp 54.200"
          trend="up"
          trendValue="+4.2%"
          subtitle="Per transaksi pelanggan"
          color="blue"
          :icon="TrendingUp"
        />
        <StatCardV2
          title="Total Transaksi Sukses"
          value="2.740"
          subtitle="Kasir POS & Dine-in"
          color="indigo"
          :icon="ShoppingBag"
        />
        <StatCardV2
          title="Metode Bayar Favorit"
          value="QRIS (68%)"
          subtitle="Non-tunai mendominasi"
          color="amber"
          :icon="CreditCard"
        />
      </div>

      <!-- Trend Sales Bar Chart & Best Sellers -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Sales Trend Bar Chart (8 cols) -->
        <div 
          class="lg:col-span-8 p-6 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between mb-6">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Tren Penjualan Harian</h3>
              <p class="text-xs text-slate-400 mt-0.5">Grafik dinamika omzet harian 12 periode terakhir</p>
            </div>
            <span class="text-xs font-bold text-emerald-600 bg-emerald-50 dark:bg-emerald-950/40 px-2.5 py-0.5 rounded-full border border-emerald-200">
              Tren Positif +14.5%
            </span>
          </div>

          <div class="h-56 flex items-end justify-between gap-2 pt-6 pb-2 border-b" style="border-color: var(--border-color);">
            <div 
              v-for="(day, idx) in dailyTrend" 
              :key="idx"
              class="flex-1 flex flex-col items-center gap-1.5 h-full justify-end group cursor-pointer"
            >
              <div 
                class="w-full max-w-[28px] bg-blue-600 hover:bg-blue-700 rounded-t-lg transition-all relative"
                :style="{ height: `${(day.amount / 6000000) * 100}%` }"
                :title="`${day.date}: Rp ${formatNum(day.amount)}`"
              ></div>
              <span class="text-[10px] font-mono text-slate-400 group-hover:text-slate-900 dark:group-hover:text-slate-100">
                {{ day.date }}
              </span>
            </div>
          </div>
        </div>

        <!-- Best Sellers Leaderboard (4 cols) -->
        <div 
          class="lg:col-span-4 p-6 rounded-2xl border shadow-xs space-y-4"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Menu Terlaris (Top 5)</h3>
            <p class="text-xs text-slate-400 mt-0.5">Produk dengan kontribusi omzet tertinggi</p>
          </div>

          <div class="space-y-3">
            <div 
              v-for="(prod, pIdx) in topProducts" 
              :key="prod.name"
              class="p-3 rounded-xl border flex items-center justify-between text-xs"
              style="background-color: var(--bg-content); border-color: var(--border-color);"
            >
              <div class="flex items-center gap-3">
                <span class="w-6 h-6 rounded-md font-bold font-mono text-xs flex items-center justify-center"
                      :class="pIdx === 0 ? 'bg-amber-100 text-amber-800' : pIdx === 1 ? 'bg-slate-200 text-slate-700' : 'bg-slate-100 text-slate-600'">
                  #{{ pIdx + 1 }}
                </span>
                <div>
                  <div class="font-bold text-slate-900 dark:text-slate-100">{{ prod.name }}</div>
                  <div class="text-[10px] text-slate-400">{{ prod.sold }} Porsi Terjual</div>
                </div>
              </div>
              <div class="text-right font-mono font-bold text-blue-600 dark:text-blue-400">
                Rp {{ formatNum(prod.revenue) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB LAINNYA -->
    <div v-else class="p-12 text-center text-xs text-slate-400 border border-dashed rounded-2xl" style="border-color: var(--border-color);">
      Laporan terintegrasi untuk "{{ activeTab.toUpperCase() }}" siap diexport ke format PDF dan Excel spreadsheet.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import PageHeaderV2 from '@/components/v2/common/PageHeaderV2.vue'
import StatCardV2 from '@/components/v2/common/StatCardV2.vue'
import {
  FileSpreadsheet,
  Download,
  ShoppingBag,
  DollarSign,
  Package,
  Users,
  Sliders,
  TrendingUp,
  CreditCard
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const activeTab = ref('sales')
const selectedRange = ref('this-month')

const reportTabs = [
  { id: 'sales', label: 'Penjualan (Sales)', icon: ShoppingBag },
  { id: 'financial', label: 'Laba Rugi (Finance)', icon: DollarSign },
  { id: 'inventory', label: 'Valuasi Stok (Inventory)', icon: Package },
  { id: 'hr', label: 'Produktivitas SDM (HR)', icon: Users },
  { id: 'custom', label: 'Laporan Kustom', icon: Sliders }
]

const dailyTrend = [
  { date: '01 Okt', amount: 3850000 },
  { date: '02 Okt', amount: 4200000 },
  { date: '03 Okt', amount: 4950000 },
  { date: '04 Okt', amount: 5600000 },
  { date: '05 Okt', amount: 3900000 },
  { date: '06 Okt', amount: 4100000 },
  { date: '07 Okt', amount: 4800000 },
  { date: '08 Okt', amount: 5200000 }
]

const topProducts = [
  { name: 'Caffe Latte Oatmilk', sold: 480, revenue: 18240000 },
  { name: 'Americano Ice Single', sold: 420, revenue: 11760000 },
  { name: 'V60 Gayo Natural', sold: 290, revenue: 10150000 },
  { name: 'Croissant Butter French', sold: 210, revenue: 6720000 },
  { name: 'Truffle Fries Parmesan', sold: 160, revenue: 6720000 }
]

const exportExcel = () => {
  notifyStore.success('File spreadsheet laporan (.xlsx) berhasil diunduh', 'Export Excel')
}

const exportPDF = () => {
  notifyStore.success('Dokumen ringkasan eksekutif (.pdf) berhasil diekspor', 'Export PDF')
}

const formatNum = (val: number | string) => Number(val || 0).toLocaleString('id-ID')
</script>

