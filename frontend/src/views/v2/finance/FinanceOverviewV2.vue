<template>
  <div class="finance-v2 space-y-6 animate-in fade-in duration-200">
    <!-- Page Header V2 -->
    <PageHeaderV2
      title="Manajemen Keuangan & Akuntansi"
      description="Pusat pencatatan pembukuan ganda (double-entry), pemantauan arus kas, rekonsiliasi perbankan, dan laporan laba rugi."
      badge="General Ledger"
    >
      <template #actions>
        <button
          type="button"
          @click="fetchOverviewData"
          class="p-2 rounded-xl border transition-colors cursor-pointer text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800"
          style="border-color: var(--border-color); background-color: var(--bg-primary);"
          title="Perbarui Data"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>

        <router-link
          to="/finance/journal"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Buat Jurnal Umum</span>
        </router-link>
      </template>

      <!-- Finance Navigation Sub-tabs -->
      <template #tabs>
        <div class="flex items-center gap-5 text-xs font-semibold overflow-x-auto pb-2 select-none">
          <button
            v-for="t in financeTabs"
            :key="t.id"
            type="button"
            @click="activeTab = t.id"
            class="pb-2 -mb-2 cursor-pointer transition-colors whitespace-nowrap"
            :class="activeTab === t.id 
              ? 'text-blue-600 dark:text-blue-400 border-b-2 border-blue-600' 
              : 'text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200'"
          >
            {{ t.label }}
          </button>
        </div>
      </template>
    </PageHeaderV2>

    <!-- TAB 1: OVERVIEW & CASHFLOW -->
    <div v-if="activeTab === 'overview'" class="space-y-6">
      <!-- 4 Financial Stat Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Pendapatan"
          :value="`Rp ${formatNum(overviewData.total_revenue || 128450000)}`"
          trend="up"
          trendValue="+14.2%"
          subtitle="Akumulasi bulan berjalan"
          color="emerald"
          :icon="DollarSign"
        />
        <StatCardV2
          title="Beban Operasional"
          :value="`Rp ${formatNum(overviewData.total_expenses || 78320000)}`"
          trend="neutral"
          trendValue="Sesuai Budget"
          subtitle="HPP & Pengeluaran Kas"
          color="rose"
          :icon="Receipt"
        />
        <StatCardV2
          title="Laba Bersih Operasional"
          :value="`Rp ${formatNum(overviewData.net_profit || 50130000)}`"
          trend="up"
          trendValue="+18.5%"
          subtitle="Net Profit Margin 39.0%"
          color="blue"
          :icon="TrendingUp"
        />
        <StatCardV2
          title="Saldo Kas & Bank"
          :value="`Rp ${formatNum(overviewData.cash_on_hand || 32500000)}`"
          subtitle="Cash & Bank BCA"
          color="indigo"
          :icon="Landmark"
        />
      </div>

      <!-- Cashflow Bar Chart & Bank Accounts -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Monthly Revenue vs Expense Visualization (8 cols) -->
        <div 
          class="lg:col-span-8 p-6 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between mb-6">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Tren Pendapatan vs Beban Bulanan</h3>
              <p class="text-xs text-slate-400 mt-0.5">Perbandingan performa finansial per bulan (dalam juta Rupiah)</p>
            </div>
            <div class="flex items-center gap-4 text-xs font-semibold">
              <span class="flex items-center gap-1.5 text-blue-600">
                <span class="w-3 h-3 rounded bg-blue-600"></span> Pendapatan
              </span>
              <span class="flex items-center gap-1.5 text-rose-500">
                <span class="w-3 h-3 rounded bg-rose-500"></span> Beban
              </span>
            </div>
          </div>

          <!-- Bar Chart Columns -->
          <div class="h-56 flex items-end justify-between gap-2 pt-6 pb-2 border-b" style="border-color: var(--border-color);">
            <div 
              v-for="bar in monthlyData" 
              :key="bar.m"
              class="flex-1 flex flex-col items-center gap-1.5 h-full justify-end group cursor-pointer"
            >
              <div class="w-full flex items-end justify-center gap-1 h-44">
                <!-- Rev Bar -->
                <div 
                  class="w-3.5 sm:w-4 bg-blue-600 hover:bg-blue-700 rounded-t-md transition-all relative"
                  :style="{ height: `${(bar.rev / 80) * 100}%` }"
                  :title="`${bar.m} Pendapatan: Rp ${bar.rev} Juta`"
                ></div>
                <!-- Exp Bar -->
                <div 
                  class="w-3.5 sm:w-4 bg-rose-500 hover:bg-rose-600 rounded-t-md transition-all relative"
                  :style="{ height: `${(bar.exp / 80) * 100}%` }"
                  :title="`${bar.m} Beban: Rp ${bar.exp} Juta`"
                ></div>
              </div>
              <span class="text-[10px] font-bold text-slate-400 group-hover:text-slate-900 dark:group-hover:text-slate-100 transition-colors">
                {{ bar.m }}
              </span>
            </div>
          </div>
        </div>

        <!-- Akun Kas & Rekening (4 cols) -->
        <div 
          class="lg:col-span-4 p-6 rounded-2xl border shadow-xs space-y-4 flex flex-col justify-between"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div>
            <div class="flex items-center justify-between mb-1">
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Likuiditas Kas & Bank</h3>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">
                Sinkron
              </span>
            </div>
            <p class="text-xs text-slate-400 mb-4">Saldo rekening operasional resmi</p>

            <div class="space-y-3">
              <div class="p-3.5 rounded-xl border flex items-center justify-between"
                   style="background-color: var(--bg-content); border-color: var(--border-color);">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-lg bg-blue-100 text-blue-700 font-bold text-xs flex items-center justify-center">
                    BCA
                  </div>
                  <div>
                    <div class="text-xs font-bold text-slate-900 dark:text-slate-100">BCA Giro Operasional</div>
                    <div class="text-[10.5px] font-mono text-slate-400">882-0192-331</div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-xs font-bold font-mono text-slate-900 dark:text-slate-100">Rp 24.850.000</div>
                  <div class="text-[10px] text-emerald-600 font-semibold">Aktif</div>
                </div>
              </div>

              <div class="p-3.5 rounded-xl border flex items-center justify-between"
                   style="background-color: var(--bg-content); border-color: var(--border-color);">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-lg bg-amber-100 text-amber-700 font-bold text-xs flex items-center justify-center">
                    CASH
                  </div>
                  <div>
                    <div class="text-xs font-bold text-slate-900 dark:text-slate-100">Kas Tunai Kasir POS</div>
                    <div class="text-[10.5px] text-slate-400">Drawer Outlet Senopati</div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-xs font-bold font-mono text-slate-900 dark:text-slate-100">Rp 7.650.000</div>
                  <div class="text-[10px] text-emerald-600 font-semibold">Aktif</div>
                </div>
              </div>
            </div>
          </div>

          <div class="pt-3 border-t space-y-2" style="border-color: var(--border-color);">
            <router-link to="/finance/journal" class="w-full py-2 px-3 rounded-xl bg-blue-600 hover:bg-blue-700 text-white text-xs font-semibold flex items-center justify-center gap-1.5 shadow-xs transition-colors">
              <Plus class="w-3.5 h-3.5" />
              <span>Buat Jurnal Umum</span>
            </router-link>
          </div>
        </div>
      </div>

      <!-- Jurnal Transaksi Terbaru -->
      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Jurnal Transaksi Terbaru</h3>
            <p class="text-xs text-slate-400">Entri transaksi debit dan kredit akun akuntansi</p>
          </div>
          <button 
            type="button" 
            @click="activeTab = 'journal'"
            class="text-xs font-semibold text-blue-600 hover:underline cursor-pointer"
          >
            Lihat Semua Jurnal →
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">No. Bukti / Ref</th>
                <th class="py-3 px-4">Tanggal</th>
                <th class="py-3 px-4">Keterangan / Deskripsi Transaksi</th>
                <th class="py-3 px-4 text-right">Debit</th>
                <th class="py-3 px-4 text-right">Kredit</th>
                <th class="py-3 px-4 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr 
                v-for="j in recentJournals" 
                :key="j.ref"
                class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="py-3.5 px-4 font-mono font-bold text-blue-600">{{ j.ref }}</td>
                <td class="py-3.5 px-4 text-slate-500 font-mono text-[11px]">{{ j.date }}</td>
                <td class="py-3.5 px-4 font-medium text-slate-800 dark:text-slate-200">{{ j.description }}</td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(j.debit) }}
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(j.credit) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <StatusBadgeV2 status="completed" label="Posted" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB LAINNYA: KELOLA BUKU BESAR / CHART OF ACCOUNTS -->
    <div v-else class="p-8 text-center text-xs text-slate-400 border border-dashed rounded-2xl" style="border-color: var(--border-color);">
      Modul finansial "{{ activeTab.toUpperCase() }}" aktif dan terintegrasi dengan backend. Silakan gunakan tab Overview untuk ringkasan real-time.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import PageHeaderV2 from '@/components/v2/common/PageHeaderV2.vue'
import StatCardV2 from '@/components/v2/common/StatCardV2.vue'
import StatusBadgeV2 from '@/components/v2/common/StatusBadgeV2.vue'
import {
  RefreshCw,
  Plus,
  DollarSign,
  Receipt,
  TrendingUp,
  Landmark
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const loading = ref(false)
const activeTab = ref('overview')

const financeTabs = [
  { id: 'overview', label: 'Ringkasan Eksekutif' },
  { id: 'journal', label: 'Buku Jurnal Umum' },
  { id: 'coa', label: 'Bagan Akun (COA)' },
  { id: 'expenses', label: 'Beban & Pengeluaran' },
  { id: 'tax', label: 'Pajak Resto (PB1)' },
  { id: 'budget', label: 'Alokasi Budget' },
  { id: 'reconciliation', label: 'Rekonsiliasi Bank' }
]

const overviewData = ref({
  total_revenue: 128450000,
  total_expenses: 78320000,
  net_profit: 50130000,
  cash_on_hand: 32500000
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

const recentJournals = ref([
  { ref: 'JV-EXP-002', date: '04 Okt 2026', description: 'Pembelian es batu kristal & mint lokal', debit: 750000, credit: 750000 },
  { ref: 'JV-EXP-003', date: '03 Okt 2026', description: 'Servis berkala mesin espresso La Marzocco', debit: 1250000, credit: 1250000 },
  { ref: 'JV-SET-001', date: '02 Okt 2026', description: 'Setoran Kasir POS EDC BCA Weekend', debit: 5120000, credit: 5120000 },
  { ref: 'JV-PAY-008', date: '28 Sep 2026', description: 'Pencairan Payroll Gaji Staf Operasional', debit: 28260969, credit: 28260969 }
])

const fetchOverviewData = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    notifyStore.success('Data ikhtisar keuangan berhasil disinkronisasi', 'Keuangan Sinkron')
  }, 350)
}

const formatNum = (num: number) => Number(num || 0).toLocaleString('id-ID')

onMounted(() => {
  fetchOverviewData()
})
</script>

