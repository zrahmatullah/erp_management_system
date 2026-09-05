<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-xs flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2">
          <div class="p-2 rounded-xl bg-blue-50 text-blue-600">
            <Layers class="w-5 h-5" />
          </div>
          <h1 class="text-xl font-black text-slate-900 tracking-tight">Kartu Stok (Stock Card Ledger)</h1>
        </div>
        <p class="text-xs text-slate-500 mt-1">
          Buku besar histori pergerakan stok keluar-masuk barang, konsumsi bahan baku dari transaksi POS, dan saldo real-time.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="fetchStockMovements"
          :disabled="loading"
          class="flex items-center gap-2 px-4 py-2.5 bg-slate-100 hover:bg-slate-200 active:bg-slate-300 text-slate-700 font-bold text-xs rounded-xl transition-all cursor-pointer disabled:opacity-50"
        >
          <RotateCw class="w-3.5 h-3.5" :class="loading ? 'animate-spin' : ''" />
          <span>Segarkan Data</span>
        </button>
        <router-link
          to="/inventory"
          class="flex items-center gap-2 px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold text-xs rounded-xl shadow-md shadow-blue-600/25 transition-all cursor-pointer"
        >
          <Package class="w-3.5 h-3.5" />
          <span>Daftar Stok</span>
        </router-link>
      </div>
    </div>

    <!-- Metric Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Total Masuk (IN)</span>
          <div class="text-2xl font-black text-emerald-600 mt-1">
            +{{ formatNumber(summary.total_in) }}
          </div>
          <span class="text-[10px] text-slate-500 font-medium mt-0.5 block">Pembelian & Saldo Awal</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
          <ArrowDownLeft class="w-6 h-6" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Total Keluar (OUT)</span>
          <div class="text-2xl font-black text-rose-600 mt-1">
            -{{ formatNumber(summary.total_out) }}
          </div>
          <span class="text-[10px] text-slate-500 font-medium mt-0.5 block">Penjualan Kasir POS & Dapur</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-rose-50 text-rose-600 flex items-center justify-center">
          <ArrowUpRight class="w-6 h-6" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Total Transaksi Mutasi</span>
          <div class="text-2xl font-black text-slate-900 mt-1">
            {{ summary.total_movements }}
          </div>
          <span class="text-[10px] text-slate-500 font-medium mt-0.5 block">Log mutasi tercatat</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-blue-50 text-blue-600 flex items-center justify-center">
          <Activity class="w-6 h-6" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Status Integrasi</span>
          <div class="text-base font-black text-slate-900 mt-1 flex items-center gap-1.5">
            <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
            <span>Real-Time POS</span>
          </div>
          <span class="text-[10px] text-emerald-600 font-semibold mt-0.5 block">Auto potong saat payment</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-indigo-50 text-indigo-600 flex items-center justify-center">
          <CheckCircle2 class="w-6 h-6" />
        </div>
      </div>
    </div>

    <!-- Filter & Search Toolbar -->
    <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs space-y-3">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <!-- Live Search -->
        <div class="relative w-full md:w-80">
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari barang, SKU, no. order #ORD..."
            class="w-full pl-9 pr-3 py-2 bg-slate-50 border border-slate-200 rounded-xl text-xs font-medium focus:bg-white focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none transition-all"
          />
        </div>

        <!-- Direction / Type Filter Pills -->
        <div class="flex items-center gap-1.5 overflow-x-auto text-xs pb-1">
          <span class="text-slate-400 text-[11px] font-bold mr-1">Tipe Mutasi:</span>
          <button
            v-for="t in [
              { label: 'Semua', val: 'all' },
              { label: 'Masuk (IN)', val: 'in' },
              { label: 'Keluar (OUT)', val: 'out' },
              { label: 'Penyesuaian', val: 'adjustment' }
            ]"
            :key="t.val"
            @click="selectedType = t.val; fetchStockMovements()"
            class="px-3 py-1.5 rounded-lg font-bold transition-all cursor-pointer whitespace-nowrap"
            :class="selectedType === t.val ? 'bg-blue-600 text-white shadow-xs' : 'bg-slate-50 text-slate-600 hover:bg-slate-100'"
          >
            {{ t.label }}
          </button>
        </div>

        <!-- Filter Item Type (Bahan Baku / Produk Jadi) -->
        <div class="flex items-center gap-1.5 overflow-x-auto text-xs pb-1">
          <span class="text-slate-400 text-[11px] font-bold mr-1">Klasifikasi:</span>
          <button
            v-for="c in ['Semua', 'Bahan Baku', 'Produk Jadi']"
            :key="c"
            @click="filterClassification = c"
            class="px-3 py-1.5 rounded-lg font-bold transition-all cursor-pointer whitespace-nowrap"
            :class="filterClassification === c ? 'bg-indigo-600 text-white shadow-xs' : 'bg-slate-50 text-slate-600 hover:bg-slate-100'"
          >
            {{ c }}
          </button>
        </div>
      </div>
    </div>

    <!-- Stock Card Ledger Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <!-- Loading State -->
      <div v-if="loading" class="py-24 flex flex-col items-center justify-center text-center space-y-3">
        <div class="w-10 h-10 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-xs font-semibold text-slate-500">Mengambil buku besar kartu stok dari server...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredMovements.length === 0" class="py-24 flex flex-col items-center justify-center text-center space-y-3">
        <div class="w-16 h-16 rounded-full bg-slate-50 border border-slate-200 flex items-center justify-center text-slate-400">
          <Layers class="w-8 h-8" />
        </div>
        <div>
          <h4 class="text-sm font-bold text-slate-800">Tidak ada data mutasi stok</h4>
          <p class="text-xs text-slate-400 mt-1">Belum ada mutasi masuk atau keluar yang sesuai dengan filter Anda.</p>
        </div>
      </div>

      <!-- Table of Movements -->
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 text-[11px] uppercase font-bold text-slate-400 border-b border-slate-200 tracking-wider">
            <tr>
              <th class="py-3.5 px-4">Waktu & Tanggal</th>
              <th class="py-3.5 px-4">Nama Barang & SKU</th>
              <th class="py-3.5 px-4">Gudang / Lokasi</th>
              <th class="py-3.5 px-4">Jenis Transaksi</th>
              <th class="py-3.5 px-4 text-right text-emerald-600">Masuk (IN)</th>
              <th class="py-3.5 px-4 text-right text-rose-600">Keluar (OUT)</th>
              <th class="py-3.5 px-4 text-right font-black text-slate-900">Saldo Akhir</th>
              <th class="py-3.5 px-4">No. Referensi / Order</th>
              <th class="py-3.5 px-4">Keterangan & Operator</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 font-medium">
            <tr
              v-for="m in filteredMovements"
              :key="m.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <!-- Timestamp -->
              <td class="py-3.5 px-4 whitespace-nowrap text-slate-500">
                <div class="flex items-center gap-1.5">
                  <Clock class="w-3.5 h-3.5 text-slate-400" />
                  <span>{{ formatTimestamp(m.created_at) }}</span>
                </div>
              </td>

              <!-- Item Details -->
              <td class="py-3.5 px-4">
                <div class="flex items-center gap-2">
                  <div class="font-bold text-slate-900">{{ m.item_name }}</div>
                  <span
                    class="px-1.5 py-0.5 rounded text-[9px] font-black uppercase"
                    :class="m.item_type === 'product' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'"
                  >
                    {{ m.item_type === 'product' ? 'Produk' : 'Bahan' }}
                  </span>
                </div>
                <div class="text-[10px] text-slate-400 font-mono mt-0.5">{{ m.sku }} • {{ m.category }}</div>
              </td>

              <!-- Warehouse -->
              <td class="py-3.5 px-4 whitespace-nowrap text-slate-600">
                {{ m.warehouse }}
              </td>

              <!-- Movement Type Badge -->
              <td class="py-3.5 px-4 whitespace-nowrap">
                <span
                  class="px-2.5 py-1 rounded-lg text-[10px] font-bold inline-flex items-center gap-1"
                  :class="getMovementTypeBadge(m.type)"
                >
                  <ArrowDownLeft v-if="m.direction === 'IN'" class="w-3 h-3 text-emerald-600" />
                  <ArrowUpRight v-else-if="m.direction === 'OUT'" class="w-3 h-3 text-rose-600" />
                  <span>{{ formatMovementType(m.type) }}</span>
                </span>
              </td>

              <!-- Masuk (IN) -->
              <td class="py-3.5 px-4 text-right whitespace-nowrap font-bold text-emerald-600">
                <span v-if="m.direction === 'IN'">+{{ formatNumber(m.quantity) }} {{ m.uom }}</span>
                <span v-else class="text-slate-300">-</span>
              </td>

              <!-- Keluar (OUT) -->
              <td class="py-3.5 px-4 text-right whitespace-nowrap font-bold text-rose-600">
                <span v-if="m.direction === 'OUT'">-{{ formatNumber(m.quantity) }} {{ m.uom }}</span>
                <span v-else class="text-slate-300">-</span>
              </td>

              <!-- Saldo Akhir -->
              <td class="py-3.5 px-4 text-right whitespace-nowrap font-black text-slate-900 bg-slate-50/50">
                {{ formatNumber(m.balance_after) }} {{ m.uom }}
              </td>

              <!-- Reference -->
              <td class="py-3.5 px-4 whitespace-nowrap">
                <span
                  v-if="m.reference_no && m.reference_no !== '-'"
                  class="px-2 py-0.5 rounded-md text-[10px] font-mono font-bold bg-slate-100 text-slate-800 border border-slate-200"
                >
                  {{ m.reference_no }}
                </span>
                <span v-else class="text-slate-400 text-[10px]">-</span>
              </td>

              <!-- Remarks & Operator -->
              <td class="py-3.5 px-4 max-w-xs">
                <div class="text-xs text-slate-700 truncate" :title="m.remarks">{{ m.remarks }}</div>
                <div class="text-[10px] text-slate-400 mt-0.5">Oleh: {{ m.operator || 'Kasir POS' }}</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import {
  Layers,
  RotateCw,
  Package,
  ArrowDownLeft,
  ArrowUpRight,
  Activity,
  CheckCircle2,
  Search,
  Clock
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const loading = ref(false)
const movements = ref<any[]>([])
const summary = ref({
  total_in: 0,
  total_out: 0,
  total_movements: 0
})

const searchQuery = ref('')
const selectedType = ref('all')
const filterClassification = ref('Semua')

const fetchStockMovements = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (selectedType.value !== 'all') {
      params.type = selectedType.value
    }
    const res = await axios.get('/api/v1/inventory/stock-movements', { params })
    movements.value = res.data?.data || []
    if (res.data?.summary) {
      summary.value = res.data.summary
    }
  } catch (err: any) {
    console.error('Failed to fetch stock movements:', err)
    notifyStore.error('Gagal memuat kartu stok dari server', 'Kesalahan')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStockMovements()
})

const filteredMovements = computed(() => {
  return movements.value.filter(m => {
    // Search query
    const q = searchQuery.value.toLowerCase().trim()
    const matchesQuery = !q ||
      m.item_name?.toLowerCase().includes(q) ||
      m.sku?.toLowerCase().includes(q) ||
      m.remarks?.toLowerCase().includes(q) ||
      m.reference_no?.toLowerCase().includes(q)

    // Classification filter (Produk Jadi vs Bahan Baku)
    let matchesClass = true
    if (filterClassification.value === 'Produk Jadi') {
      matchesClass = m.item_type === 'product'
    } else if (filterClassification.value === 'Bahan Baku') {
      matchesClass = m.item_type === 'raw_material'
    }

    return matchesQuery && matchesClass
  })
})

const formatNumber = (val: number | string) => {
  const num = Number(val)
  if (isNaN(num)) return '0'
  return num.toLocaleString('id-ID', { maximumFractionDigits: 3 })
}

const formatTimestamp = (dateStr?: string) => {
  if (!dateStr) return '-'
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleString('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

const formatMovementType = (type: string) => {
  switch (type) {
    case 'in_purchase':
      return 'Pembelian PO'
    case 'initial_stock':
      return 'Saldo Awal'
    case 'out_pos_sales':
      return 'Penjualan POS'
    case 'waste':
      return 'Waste / Basi'
    case 'adjustment':
      return 'Penyesuaian Opname'
    default:
      return type
  }
}

const getMovementTypeBadge = (type: string) => {
  switch (type) {
    case 'in_purchase':
    case 'initial_stock':
      return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
    case 'out_pos_sales':
      return 'bg-rose-50 text-rose-700 border border-rose-200'
    case 'waste':
      return 'bg-amber-50 text-amber-700 border border-amber-200'
    case 'adjustment':
      return 'bg-purple-50 text-purple-700 border border-purple-200'
    default:
      return 'bg-slate-50 text-slate-700 border border-slate-200'
  }
}
</script>

