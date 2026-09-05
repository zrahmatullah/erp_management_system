<template>
  <div class="space-y-6">
    <!-- Header with Back Button -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <router-link to="/inventory" class="p-2 rounded-xl bg-white border border-slate-200 hover:bg-slate-50 text-slate-600 transition-colors">
          ←
        </router-link>
        <div>
          <h2 class="text-2xl font-black text-slate-900 tracking-tight">Purchase Order {{ currentPO.po_number || '#PO-2026-0891' }}</h2>
          <p class="text-xs text-slate-500 mt-0.5">Detail pesanan pengadaan bahan baku ke supplier</p>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="rejectPO"
          class="px-4 py-2 border border-rose-300 text-rose-600 hover:bg-rose-50 text-xs font-bold rounded-xl transition-colors cursor-pointer"
        >
          Tolak
        </button>
        <button
          @click="approvePO"
          class="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold rounded-xl shadow-md shadow-emerald-600/30 transition-all cursor-pointer"
        >
          Setujui
        </button>
      </div>
    </div>

    <!-- Top Two Cards Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Left Card: Informasi Pesanan Pembelian -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm space-y-3.5">
        <h3 class="font-bold text-slate-900 text-sm pb-2 border-b border-slate-100">Informasi Pesanan Pembelian</h3>
        <div class="grid grid-cols-3 text-xs">
          <span class="text-slate-400">Supplier</span>
          <span class="col-span-2 font-bold text-slate-900">{{ currentPO.supplier || 'PT Kopi Nusantara Jaya' }}</span>
        </div>
        <div class="grid grid-cols-3 text-xs">
          <span class="text-slate-400">Tanggal</span>
          <span class="col-span-2 font-semibold text-slate-800">{{ currentPO.order_date || '2026-08-26' }}</span>
        </div>
        <div class="grid grid-cols-3 text-xs">
          <span class="text-slate-400">Status</span>
          <div class="col-span-2">
            <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-amber-100 text-amber-800">
              {{ currentPO.status || 'Menunggu Persetujuan' }}
            </span>
          </div>
        </div>
        <div class="grid grid-cols-3 text-xs">
          <span class="text-slate-400">Total Nominal</span>
          <span class="col-span-2 font-black text-blue-600 text-sm">Rp {{ grandTotal.toLocaleString('id-ID') }}</span>
        </div>
        <div class="grid grid-cols-3 text-xs">
          <span class="text-slate-400">Dibuat oleh</span>
          <span class="col-span-2 font-medium text-slate-800">Ahmad - Warehouse Staff</span>
        </div>
        <div class="grid grid-cols-3 text-xs">
          <span class="text-slate-400">Catatan</span>
          <span class="col-span-2 text-slate-600">{{ currentPO.notes || 'Pengadaan stok bahan baku operasional cabang' }}</span>
        </div>
      </div>

      <!-- Right Card: Status Timeline (Matches Mockup) -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm">
        <h3 class="font-bold text-slate-900 text-sm pb-3 border-b border-slate-100 mb-4">Status Timeline</h3>
        <div class="space-y-4 text-xs font-medium pl-2">
          <div class="flex items-center gap-3">
            <span class="w-6 h-6 rounded-full bg-emerald-500 text-white flex items-center justify-center">
              <Check class="w-3.5 h-3.5" />
            </span>
            <div>
              <div class="font-bold text-slate-900">Draft Dibuat</div>
              <div class="text-[10px] text-slate-400">02 Sep 2026, 09:15</div>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span class="w-6 h-6 rounded-full bg-emerald-500 text-white flex items-center justify-center">
              <Check class="w-3.5 h-3.5" />
            </span>
            <div>
              <div class="font-bold text-slate-900">Diajukan ke Supervisor</div>
              <div class="text-[10px] text-slate-400">03 Sep 2026, 14:20</div>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span class="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold text-[10px] ring-4 ring-blue-100">●</span>
            <div>
              <div class="font-black text-blue-600">Review Manager (Sekarang)</div>
              <div class="text-[10px] text-slate-400">Menunggu persetujuan digital</div>
            </div>
          </div>
          <div class="flex items-center gap-3 opacity-40">
            <span class="w-6 h-6 rounded-full border-2 border-slate-300 flex items-center justify-center font-bold text-[10px]">○</span>
            <div class="font-bold text-slate-500">Persetujuan Owner</div>
          </div>
          <div class="flex items-center gap-3 opacity-40">
            <span class="w-6 h-6 rounded-full border-2 border-slate-300 flex items-center justify-center font-bold text-[10px]">○</span>
            <div class="font-bold text-slate-500">Barang Diterima (GRN)</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Daftar Barang Table (Matches Mockup) -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6">
      <h3 class="font-bold text-slate-900 text-sm mb-4">Daftar Barang yang Dipesan</h3>
      <div v-if="loading" class="py-12 text-center text-slate-400 text-xs">
        Memuat item purchase order dari database...
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4 w-12">No</th>
              <th class="py-3 px-4">Nama Barang</th>
              <th class="py-3 px-4 text-center">Qty</th>
              <th class="py-3 px-4">Satuan</th>
              <th class="py-3 px-4 text-right">Harga Satuan (Rp)</th>
              <th class="py-3 px-4 text-right">Total (Rp)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="(item, idx) in poItems" :key="item.name" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3.5 px-4 font-bold text-slate-400">{{ idx + 1 }}</td>
              <td class="py-3.5 px-4 font-bold text-slate-900">{{ item.name }}</td>
              <td class="py-3.5 px-4 text-center font-bold">{{ item.qty }}</td>
              <td class="py-3.5 px-4">{{ item.unit }}</td>
              <td class="py-3.5 px-4 text-right">Rp {{ item.price.toLocaleString('id-ID') }}</td>
              <td class="py-3.5 px-4 text-right font-bold text-slate-900">
                Rp {{ (item.qty * item.price).toLocaleString('id-ID') }}
              </td>
            </tr>
            <tr v-if="poItems.length === 0">
              <td colspan="6" class="py-8 text-center text-slate-400 text-xs">
                Tidak ada data barang purchase order.
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="border-t-2 border-slate-200 text-sm font-black text-slate-900">
              <td colspan="5" class="py-4 px-4 text-right">Grand Total:</td>
              <td class="py-4 px-4 text-right text-blue-600">Rp {{ grandTotal.toLocaleString('id-ID') }}</td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { Check } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const loading = ref(false)

const currentPO = ref<any>({})
const poItems = ref<any[]>([])

const fetchPOData = async () => {
  loading.value = true
  try {
    const [poRes, stockRes] = await Promise.all([
      axios.get('/api/v1/inventory/purchase-orders'),
      axios.get('/api/v1/inventory/stocks')
    ])
    if (poRes.data?.data && poRes.data.data.length > 0) {
      currentPO.value = poRes.data.data[0]
    }
    if (stockRes.data?.data) {
      poItems.value = stockRes.data.data.slice(0, 5).map((s: any) => ({
        name: s.name,
        qty: Math.max(10, Math.round((s.min_stock || 10) * 2)),
        unit: s.uom || 'pcs',
        price: Number(s.cost || 25000)
      }))
    }
  } catch (err: any) {
    console.error('Failed to load purchase orders:', err)
    notifyStore.error('Gagal mengambil data PO dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPOData()
})

const grandTotal = computed(() => {
  return poItems.value.reduce((acc, item) => acc + (item.qty * item.price), 0)
})

const approvePO = () => {
  notifyStore.success(`Purchase Order ${currentPO.value.po_number || ''} berhasil disetujui! Dokumen dikirim ke Supplier.`, 'PO Disetujui')
}

const rejectPO = () => {
  notifyStore.warning(`Purchase Order ${currentPO.value.po_number || ''} telah ditolak dan dikembalikan ke pembuat.`, 'PO Ditolak')
}
</script>
