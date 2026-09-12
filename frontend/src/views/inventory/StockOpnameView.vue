<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Stock Opname & Audit Fisik</h2>
        <p class="text-xs text-slate-500 mt-1">Audit fisik persediaan gudang dan penyesuaian otomatis ke kartu stok & buku besar</p>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="initOpnameSession"
          class="p-2.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
          title="Segarkan Data"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="flex items-center gap-4 border-b border-slate-200 text-xs font-bold pb-2">
      <button
        @click="activeMainTab = 'session'"
        class="pb-2 relative transition-colors cursor-pointer"
        :class="activeMainTab === 'session' ? 'text-blue-600 border-b-2 border-blue-600 -mb-2.5' : 'text-slate-400 hover:text-slate-700'"
      >
        <span>Sesi Hitung Fisik Aktif</span>
      </button>
      <button
        @click="activeMainTab = 'history'"
        class="pb-2 relative transition-colors cursor-pointer"
        :class="activeMainTab === 'history' ? 'text-blue-600 border-b-2 border-blue-600 -mb-2.5' : 'text-slate-400 hover:text-slate-700'"
      >
        <span>Riwayat Audit Opname ({{ opnameHistory.length }})</span>
      </button>
    </div>

    <!-- TAB 1: SESI HITUNG FISIK -->
    <div v-if="activeMainTab === 'session'" class="space-y-6">
      <!-- Info Bar -->
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-sm flex flex-wrap items-center justify-between gap-4 text-xs font-medium text-slate-600">
        <div class="flex items-center gap-2">
          <Warehouse class="w-4 h-4 text-blue-600" />
          <span>Gudang: <strong class="text-slate-900">Main Warehouse (Pusat)</strong></span>
        </div>
        <div class="flex items-center gap-2">
          <Calendar class="w-4 h-4 text-blue-600" />
          <span>Tanggal: <strong class="text-slate-900">{{ opnameDate }}</strong></span>
        </div>
        <div class="flex items-center gap-2">
          <User class="w-4 h-4 text-blue-600" />
          <span>Petugas: <strong class="text-slate-900">Staff Audit Inventory</strong></span>
        </div>
      </div>

      <!-- Table Section -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6">
        <div v-if="loading" class="py-16 text-center text-slate-400 text-xs flex flex-col items-center justify-center gap-2">
          <RotateCw class="w-6 h-6 animate-spin text-blue-600" />
          <span>Mengambil data stok bahan baku dari database server...</span>
        </div>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
                <th class="py-3 px-4 w-12">No</th>
                <th class="py-3 px-4">Kode SKU</th>
                <th class="py-3 px-4">Nama Bahan Baku</th>
                <th class="py-3 px-4">Satuan</th>
                <th class="py-3 px-4 text-center">Stok Sistem</th>
                <th class="py-3 px-4 text-center w-32">Stok Fisik</th>
                <th class="py-3 px-4 text-center">Selisih</th>
                <th class="py-3 px-4">Keterangan / Alasan</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
              <tr v-for="(item, idx) in opnameItems" :key="item.id || item.code" class="hover:bg-slate-50/80 transition-colors">
                <td class="py-3.5 px-4 font-bold text-slate-400">{{ idx + 1 }}</td>
                <td class="py-3.5 px-4 font-mono font-bold text-blue-600">{{ item.code }}</td>
                <td class="py-3.5 px-4 font-semibold text-slate-800">{{ item.name }}</td>
                <td class="py-3.5 px-4">{{ item.unit }}</td>
                <td class="py-3.5 px-4 text-center font-bold text-slate-700">{{ item.systemStock }}</td>
                <td class="py-3.5 px-4 text-center">
                  <input
                    v-model.number="item.physicalStock"
                    type="number"
                    step="0.1"
                    class="w-24 text-center py-1 font-bold text-slate-900 bg-slate-50 border border-slate-300 rounded-lg outline-none focus:border-blue-500 focus:bg-white"
                  />
                </td>
                <td class="py-3.5 px-4 text-center font-extrabold text-sm"
                  :class="diff(item) < 0 ? 'text-rose-600' : diff(item) > 0 ? 'text-emerald-600' : 'text-slate-400'">
                  {{ diff(item) > 0 ? `+${diff(item)}` : diff(item) }}
                </td>
                <td class="py-3.5 px-4">
                  <input
                    v-model="item.notes"
                    type="text"
                    placeholder="Catatan perbedaan jika ada..."
                    class="w-full px-2.5 py-1 text-xs bg-slate-50 border border-slate-200 rounded-lg outline-none focus:bg-white focus:border-blue-500"
                  />
                </td>
              </tr>
              <tr v-if="opnameItems.length === 0">
                <td colspan="8" class="py-8 text-center text-slate-400 text-xs">
                  Tidak ada data barang inventori ditemukan di database server.
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Bottom Summary Bar -->
        <div class="mt-6 pt-4 border-t border-slate-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4 text-xs font-medium">
          <div class="flex items-center gap-6">
            <div>Total Item: <strong class="text-slate-900 text-sm ml-1">{{ opnameItems.length }}</strong></div>
            <div>Sesuai: <strong class="text-slate-900 text-sm ml-1">{{ countMatch }}</strong></div>
            <div>Selisih Kurang: <strong class="text-rose-600 text-sm ml-1">{{ countUnder }}</strong></div>
            <div>Selisih Lebih: <strong class="text-emerald-600 text-sm ml-1">{{ countOver }}</strong></div>
          </div>

          <div class="flex items-center gap-3">
            <button
              @click="finishOpname"
              :disabled="submitting"
              class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
            >
              <CheckCircle2 class="w-4 h-4" />
              <span>{{ submitting ? 'Memproses Penyesuaian...' : 'Selesai & Terapkan Penyesuaian Stok' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: RIWAYAT AUDIT OPNAMES -->
    <div v-else class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6 space-y-4">
      <div class="flex items-center justify-between border-b border-slate-100 pb-3">
        <h3 class="text-sm font-bold text-slate-900">Daftar Dokumen Audit Stock Opname Sebelumnya</h3>
        <span class="text-xs text-slate-400">Total: {{ opnameHistory.length }} Audit</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4">No. Opname</th>
              <th class="py-3 px-4">Tanggal</th>
              <th class="py-3 px-4">Gudang</th>
              <th class="py-3 px-4 text-center">Item Tercatat</th>
              <th class="py-3 px-4">Catatan</th>
              <th class="py-3 px-4 text-center">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="op in opnameHistory" :key="op.id" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3.5 px-4 font-mono font-bold text-blue-600">{{ op.opname_number }}</td>
              <td class="py-3.5 px-4 font-mono text-slate-500">{{ op.opname_date }}</td>
              <td class="py-3.5 px-4 font-bold text-slate-800">{{ op.warehouse }}</td>
              <td class="py-3.5 px-4 text-center">
                <span class="px-2 py-0.5 bg-slate-100 rounded-md font-bold text-slate-700 text-[11px]">
                  {{ op.total_items }} item
                </span>
              </td>
              <td class="py-3.5 px-4 text-slate-500">{{ op.notes || '-' }}</td>
              <td class="py-3.5 px-4 text-center">
                <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
                  {{ op.status }}
                </span>
              </td>
            </tr>
            <tr v-if="opnameHistory.length === 0">
              <td colspan="6" class="py-12 text-center text-slate-400 text-xs">
                Belum ada data riwayat stock opname.
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
import { Warehouse, Calendar, User, RotateCw, CheckCircle2 } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const loading = ref(false)
const submitting = ref(false)
const activeMainTab = ref('session')
const opnameDate = ref(new Date().toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' }))

const opnameItems = ref<any[]>([])
const opnameHistory = ref<any[]>([])

const initOpnameSession = async () => {
  loading.value = true
  try {
    const [stocksRes, histRes] = await Promise.all([
      axios.get('/api/v1/inventory/stocks'),
      axios.get('/api/v1/inventory/opnames')
    ])

    if (stocksRes.data?.data) {
      opnameItems.value = stocksRes.data.data.map((s: any) => ({
        id: s.id,
        code: s.sku || 'INV-000',
        name: s.name,
        unit: s.uom || 'pcs',
        systemStock: Number(s.stock || 0),
        physicalStock: Number(s.stock || 0),
        notes: ''
      }))
    }

    if (histRes.data?.data) {
      opnameHistory.value = histRes.data.data
    }
  } catch (err: any) {
    console.error('Failed to load stock opname data:', err)
    notifyStore.error('Gagal mengambil data inventori dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initOpnameSession()
})

const diff = (item: any) => {
  return Math.round(((item.physicalStock || 0) - item.systemStock) * 100) / 100
}

const countMatch = computed(() => opnameItems.value.filter(i => diff(i) === 0).length)
const countUnder = computed(() => opnameItems.value.filter(i => diff(i) < 0).length)
const countOver = computed(() => opnameItems.value.filter(i => diff(i) > 0).length)

const finishOpname = async () => {
  if (opnameItems.value.length === 0) {
    notifyStore.warning('Tidak ada item inventori untuk di-opname!', 'Validasi')
    return
  }

  submitting.value = true
  try {
    const payload = {
      notes: `Stock Opname Rutin - ${opnameDate.value}`,
      items: opnameItems.value.map(i => ({
        inventory_item_id: i.id,
        system_stock: i.systemStock,
        physical_stock: i.physicalStock,
        notes: i.notes || (diff(i) !== 0 ? `Penyesuaian stok ${diff(i) > 0 ? '+' : ''}${diff(i)} ${i.unit}` : 'Sesuai audit fisik')
      }))
    }

    const res = await axios.post('/api/v1/inventory/opnames', payload)
    notifyStore.success(res.data?.message || 'Stock opname berhasil disimpan dan penyesuaian telah dicatat ke kartu stok!', 'Opname Berhasil')
    await initOpnameSession()
    activeMainTab.value = 'history'
  } catch (err: any) {
    console.error('Failed to finish stock opname:', err)
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan hasil stock opname ke database', 'Error')
  } finally {
    submitting.value = false
  }
}
</script>
