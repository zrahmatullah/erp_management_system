<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Inventory Management</h2>
        <p class="text-xs text-slate-500 mt-1">Pantau pergerakan stok bahan baku, resep, dan logistik multi-cabang</p>
      </div>
    </div>

    <!-- Top Metrics & Quick Actions Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- 4 Metric Cards (Left 3 cols) -->
      <div class="lg:col-span-3 grid grid-cols-2 sm:grid-cols-4 gap-4">
        <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm">
          <div class="text-xs font-semibold text-slate-500 mb-1">Total Items</div>
          <div class="text-3xl font-black text-slate-900">{{ items.length }}</div>
        </div>
        <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm">
          <div class="text-xs font-semibold text-slate-500 mb-1">Low Stock Alerts</div>
          <div class="text-3xl font-black text-rose-600">{{ countLowStock }}</div>
        </div>
        <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm">
          <div class="text-xs font-semibold text-slate-500 mb-1">Pending PO</div>
          <div class="text-3xl font-black text-slate-900">3</div>
        </div>
        <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm">
          <div class="text-xs font-semibold text-slate-500 mb-1">Gudang Aktif</div>
          <div class="text-3xl font-black text-emerald-600">2</div>
        </div>
      </div>

      <!-- Quick Actions (Right 1 col) -->
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-sm flex flex-col justify-between gap-2">
        <div class="text-xs font-bold text-slate-700 px-1">Quick Actions</div>
        <div class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-2 gap-2">
          <router-link to="/settings/master" class="py-2 px-3 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl text-center shadow-xs transition-colors">
            New Item
          </router-link>
          <router-link to="/inventory/po" class="py-2 px-3 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl text-center shadow-xs transition-colors">
            Create PO
          </router-link>
          <router-link to="/inventory/opname" class="py-2 px-3 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl text-center shadow-xs transition-colors">
            Stock Opname
          </router-link>
          <button @click="alertTransfer" class="py-2 px-3 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl text-center shadow-xs transition-colors cursor-pointer">
            Transfer Stock
          </button>
        </div>
      </div>
    </div>

    <!-- Quick Filter Panel (Matches Mockup) -->
    <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm">
      <div class="text-xs font-bold text-slate-800 mb-3">Quick Filter</div>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
        <!-- Search -->
        <div class="lg:col-span-2 relative">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400">
            <Search class="w-3.5 h-3.5" />
          </span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="e.g., 'Coffee Beans'"
            class="w-full pl-8 pr-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none focus:bg-white focus:border-blue-500"
          />
        </div>

        <!-- Category -->
        <div>
          <select v-model="filterCategory" class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none">
            <option value="">Semua Kategori</option>
            <option v-for="cat in availableCategories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>

        <!-- Warehouse -->
        <div>
          <select v-model="filterWarehouse" class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none">
            <option value="">Semua Gudang</option>
            <option value="Main Store Senopati">Main Store Senopati</option>
            <option value="Central Kitchen">Central Kitchen</option>
          </select>
        </div>

        <!-- Status -->
        <div>
          <select v-model="filterStatus" class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none">
            <option value="">Semua Status</option>
            <option value="Normal">Normal</option>
            <option value="Low Stock">Low Stock</option>
            <option value="Out of Stock">Out of Stock</option>
          </select>
        </div>

        <!-- Reset Button -->
        <div>
          <button
            @click="resetFilters"
            class="w-full py-2 px-3 border border-slate-200 hover:bg-slate-50 text-slate-700 text-xs font-bold rounded-xl transition-colors flex items-center justify-center gap-1.5 cursor-pointer"
          >
            <RotateCw class="w-3.5 h-3.5" />
            <span>Reset</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Inventory Table Section -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6">
      <div v-if="loading" class="py-12 text-center text-slate-400 text-xs">
        Memuat data inventori bahan baku dari server...
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4">Item Details</th>
              <th class="py-3 px-4">Category</th>
              <th class="py-3 px-4 text-center">Unit</th>
              <th class="py-3 px-4 text-center">Current Stock</th>
              <th class="py-3 px-4 text-center">Min / Max</th>
              <th class="py-3 px-4">Warehouse</th>
              <th class="py-3 px-4 text-center">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="item in filteredItems" :key="item.code" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3.5 px-4 font-bold text-slate-900">
                <div class="text-xs">{{ item.name }}</div>
                <div class="text-[10px] text-slate-400 font-mono font-normal">{{ item.code }}</div>
              </td>
              <td class="py-3.5 px-4">{{ item.category }}</td>
              <td class="py-3.5 px-4 text-center">{{ item.unit }}</td>
              <td class="py-3.5 px-4 text-center font-black text-sm text-slate-900">{{ item.stock }}</td>
              <td class="py-3.5 px-4 text-center text-slate-400 font-mono text-[11px]">
                {{ item.min }} / {{ item.max }}
              </td>
              <td class="py-3.5 px-4">{{ item.warehouse }}</td>
              <td class="py-3.5 px-4 text-center">
                <span
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                  :class="getStatusClass(item.status)"
                >
                  {{ item.status }}
                </span>
              </td>
            </tr>
            <tr v-if="filteredItems.length === 0">
              <td colspan="7" class="py-8 text-center text-slate-400 text-xs">
                Tidak ada data barang inventori sesuai filter.
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
import { Search, RotateCw } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const searchQuery = ref('')
const filterCategory = ref('')
const filterWarehouse = ref('')
const filterStatus = ref('')
const loading = ref(false)

const items = ref<any[]>([])

const fetchStockData = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/inventory/stocks')
    if (res.data?.data) {
      items.value = res.data.data.map((s: any) => ({
        code: s.sku || 'INV-000',
        name: s.name,
        category: s.category || 'General',
        unit: s.uom || 'pcs',
        stock: Number(s.stock || 0),
        min: Number(s.min_stock || 10),
        max: Number(s.min_stock || 10) * 10,
        warehouse: s.warehouse || 'Main Store Senopati',
        status: Number(s.stock || 0) <= 0 ? 'Out of Stock' : (Number(s.stock || 0) <= Number(s.min_stock || 10) ? 'Low Stock' : 'Normal')
      }))
    }
  } catch (err: any) {
    console.error('Failed to load stock list:', err)
    notifyStore.error('Gagal mengambil daftar inventori dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStockData()
})

const availableCategories = computed(() => {
  const cats = new Set<string>()
  items.value.forEach(i => {
    if (i.category) cats.add(i.category)
  })
  return Array.from(cats)
})

const countLowStock = computed(() => items.value.filter(i => i.status === 'Low Stock' || i.status === 'Out of Stock').length)

const filteredItems = computed(() => {
  return items.value.filter(i => {
    const mSearch = !searchQuery.value || i.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || i.code.toLowerCase().includes(searchQuery.value.toLowerCase())
    const mCat = !filterCategory.value || i.category === filterCategory.value
    const mWar = !filterWarehouse.value || i.warehouse === filterWarehouse.value
    const mStat = !filterStatus.value || i.status === filterStatus.value
    return mSearch && mCat && mWar && mStat
  })
})

const getStatusClass = (status: string) => {
  switch (status) {
    case 'Normal':
      return 'bg-emerald-100 text-emerald-700'
    case 'Low Stock':
      return 'bg-rose-100 text-rose-700'
    case 'Out of Stock':
      return 'bg-red-200 text-red-800'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

const resetFilters = () => {
  searchQuery.value = ''
  filterCategory.value = ''
  filterWarehouse.value = ''
  filterStatus.value = ''
  notifyStore.info('Filter inventori berhasil direset.', 'Filter')
}

const alertTransfer = () => {
  notifyStore.info('Modul Transfer Stok multi-outlet siap digunakan.', 'Transfer Stok')
}
</script>
