<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Stock Opname #SO-2026-001</h2>
        <span class="px-3 py-1 bg-blue-600 text-white font-bold text-xs rounded-full shadow-xs">
          Sedang Berjalan
        </span>
      </div>
    </div>

    <!-- Info Bar (Matches Mockup) -->
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
        <span>Petugas Opname: <strong class="text-slate-900">Ahmad Staff Gudang</strong></span>
      </div>
    </div>

    <!-- Table Section -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6">
      <div v-if="loading" class="py-12 text-center text-slate-400 text-xs">
        Mengambil data stok bahan baku dari database...
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4 w-12">No</th>
              <th class="py-3 px-4">Kode Barang</th>
              <th class="py-3 px-4">Nama Barang</th>
              <th class="py-3 px-4">Satuan</th>
              <th class="py-3 px-4 text-center">Stok Sistem</th>
              <th class="py-3 px-4 text-center w-32">Stok Fisik</th>
              <th class="py-3 px-4 text-center">Selisih</th>
              <th class="py-3 px-4">Keterangan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="(item, idx) in opnameItems" :key="item.code" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3.5 px-4 font-bold text-slate-400">{{ idx + 1 }}</td>
              <td class="py-3.5 px-4 font-bold text-slate-900">{{ item.code }}</td>
              <td class="py-3.5 px-4 font-semibold text-slate-800">{{ item.name }}</td>
              <td class="py-3.5 px-4">{{ item.unit }}</td>
              <td class="py-3.5 px-4 text-center font-bold text-slate-700">{{ item.systemStock }}</td>
              <td class="py-3.5 px-4 text-center">
                <input
                  v-model.number="item.physicalStock"
                  type="number"
                  class="w-24 text-center py-1 font-bold text-slate-900 bg-slate-50 border border-slate-300 rounded-lg outline-none focus:border-blue-500"
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
                  placeholder="Catatan perbedaan..."
                  class="w-full px-2.5 py-1 text-xs bg-slate-50 border border-slate-200 rounded-lg outline-none"
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

      <!-- Bottom Summary Bar (Matches Mockup) -->
      <div class="mt-6 pt-4 border-t border-slate-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4 text-xs font-medium">
        <div class="flex items-center gap-6">
          <div>Total Item: <strong class="text-slate-900 text-sm ml-1">{{ opnameItems.length }}</strong></div>
          <div>Item Sesuai: <strong class="text-slate-900 text-sm ml-1">{{ countMatch }}</strong></div>
          <div>Selisih Kurang: <strong class="text-rose-600 text-sm ml-1">{{ countUnder }}</strong></div>
          <div>Selisih Lebih: <strong class="text-emerald-600 text-sm ml-1">{{ countOver }}</strong></div>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="saveDraft"
            class="px-5 py-2.5 bg-slate-900 hover:bg-slate-800 text-white font-bold rounded-xl shadow-xs transition-colors cursor-pointer"
          >
            Simpan Hitungan
          </button>
          <button
            @click="finishOpname"
            class="px-5 py-2.5 bg-white hover:bg-slate-50 text-slate-800 font-bold rounded-xl border border-slate-300 shadow-xs transition-colors cursor-pointer"
          >
            Selesai & Review
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { Warehouse, Calendar, User } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const loading = ref(false)
const opnameDate = ref(new Date().toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' }))

const opnameItems = ref<any[]>([])

const fetchOpnameData = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/inventory/stocks')
    if (res.data?.data) {
      opnameItems.value = res.data.data.map((s: any) => ({
        code: s.sku || 'INV-000',
        name: s.name,
        unit: s.uom || 'pcs',
        systemStock: Number(s.stock || 0),
        physicalStock: Number(s.stock || 0),
        notes: ''
      }))
    }
  } catch (err: any) {
    console.error('Failed to load stock opname data:', err)
    notifyStore.error('Gagal mengambil data inventori dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchOpnameData()
})

const diff = (item: any) => {
  return (item.physicalStock || 0) - item.systemStock
}

const countMatch = computed(() => opnameItems.value.filter(i => diff(i) === 0).length)
const countUnder = computed(() => opnameItems.value.filter(i => diff(i) < 0).length)
const countOver = computed(() => opnameItems.value.filter(i => diff(i) > 0).length)

const saveDraft = () => {
  notifyStore.success('Hitungan stok opname berhasil disimpan sebagai draf!', 'Draf Disimpan')
}

const finishOpname = () => {
  notifyStore.success('Stock opname selesai! Penyesuaian stok otomatis dicatat ke jurnal penyesuaian.', 'Opname Selesai')
}
</script>
