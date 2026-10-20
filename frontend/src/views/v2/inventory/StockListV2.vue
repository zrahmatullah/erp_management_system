<template>
  <div class="inventory-v2 space-y-6 animate-in fade-in duration-200">
    <!-- Page Header V2 -->
    <PageHeaderV2
      title="Inventaris & Manajemen Stok"
      description="Pantau mutasi bahan baku, ketersediaan stok fisik outlet, dan integrasi rantai pasok multi-gudang."
      badge="Logistik Aktif"
    >
      <template #actions>
        <router-link
          to="/settings/master"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Tambah Bahan Baru</span>
        </router-link>

        <router-link
          to="/inventory/po"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold border transition-all flex items-center gap-1.5 cursor-pointer text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800"
          style="border-color: var(--border-color); background-color: var(--bg-primary);"
        >
          <FileText class="w-3.5 h-3.5 text-blue-600" />
          <span>Buat PO Restok</span>
        </router-link>

        <router-link
          to="/inventory/opname"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold border transition-all flex items-center gap-1.5 cursor-pointer text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800"
          style="border-color: var(--border-color); background-color: var(--bg-primary);"
        >
          <ClipboardCheck class="w-3.5 h-3.5 text-emerald-600" />
          <span>Stock Opname</span>
        </router-link>
      </template>
    </PageHeaderV2>

    <!-- 4 Stat Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
      <StatCardV2
        title="Total SKU Bahan"
        :value="items.length"
        subtitle="Terdaftar di katalog"
        color="blue"
        :icon="Package"
      />
      <StatCardV2
        title="Stok Menipis / Kritis"
        :value="countLowStock"
        trend="down"
        trendValue="Perlu Restok"
        subtitle="Di bawah safety stock"
        color="rose"
        :icon="AlertTriangle"
      />
      <StatCardV2
        title="Purchase Orders Aktif"
        :value="3"
        subtitle="Sedang diproses supplier"
        color="amber"
        :icon="FileSpreadsheet"
      />
      <StatCardV2
        title="Total Nilai Valuasi"
        :value="`Rp ${formatNum(totalValuation)}`"
        subtitle="Metode FIFO Aktif"
        color="emerald"
        :icon="Boxes"
      />
    </div>

    <!-- Filter Toolbar -->
    <div 
      class="p-4 rounded-2xl border space-y-3 shadow-xs"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3 text-xs">
        <!-- Search -->
        <div class="lg:col-span-2 relative">
          <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5 pointer-events-none" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari SKU atau nama bahan (e.g. 'Coffee Beans', 'Susu')..."
            class="w-full pl-8.5 pr-3 py-1.5 text-xs rounded-xl border outline-none transition-all focus:ring-1 focus:ring-blue-500"
            style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
          />
        </div>

        <!-- Category -->
        <div>
          <select
            v-model="filterCategory"
            class="w-full px-3 py-1.5 text-xs rounded-xl border outline-none"
            style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
          >
            <option value="">Semua Kategori</option>
            <option v-for="cat in availableCategories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>

        <!-- Warehouse -->
        <div>
          <select
            v-model="filterWarehouse"
            class="w-full px-3 py-1.5 text-xs rounded-xl border outline-none"
            style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
          >
            <option value="">Semua Gudang</option>
            <option value="Main Store Senopati">Main Store Senopati</option>
            <option value="Central Kitchen">Central Kitchen</option>
          </select>
        </div>

        <!-- Reset Button -->
        <div>
          <button
            type="button"
            @click="resetFilters"
            class="w-full py-1.5 px-3 rounded-xl border font-semibold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer"
            style="border-color: var(--border-color);"
          >
            Reset Filter
          </button>
        </div>
      </div>
    </div>

    <!-- Clean Inventory Table Card -->
    <div 
      class="rounded-2xl border shadow-xs overflow-hidden"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
              <th class="py-3 px-4">SKU / Nama Bahan</th>
              <th class="py-3 px-4">Kategori</th>
              <th class="py-3 px-4">Gudang</th>
              <th class="py-3 px-4 text-right">Stok Fisik</th>
              <th class="py-3 px-4 text-right">Biaya Satuan</th>
              <th class="py-3 px-4 text-right">Total Nilai</th>
              <th class="py-3 px-4 text-center">Status</th>
              <th class="py-3 px-4 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y" style="border-color: var(--border-color);">
            <tr v-if="filteredItems.length === 0">
              <td colspan="8" class="py-12 text-center text-slate-400">
                Tidak ada data bahan inventaris yang cocok dengan filter pencarian.
              </td>
            </tr>

            <tr 
              v-for="item in filteredItems" 
              :key="item.id"
              class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
            >
              <!-- Name & SKU -->
              <td class="py-3.5 px-4">
                <div class="font-bold text-slate-900 dark:text-slate-100">{{ item.name }}</div>
                <div class="font-mono text-[10.5px] text-slate-400 mt-0.5">{{ item.sku }}</div>
              </td>

              <!-- Category -->
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-md text-[11px] font-medium bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300">
                  {{ item.category }}
                </span>
              </td>

              <!-- Warehouse -->
              <td class="py-3.5 px-4 text-slate-600 dark:text-slate-300">
                {{ item.warehouse }}
              </td>

              <!-- Stock Level -->
              <td class="py-3.5 px-4 text-right">
                <div class="font-mono font-bold text-slate-900 dark:text-slate-100">
                  {{ item.current_stock }} {{ item.unit }}
                </div>
                <div class="text-[10px] text-slate-400 mt-0.5">Min: {{ item.min_stock }} {{ item.unit }}</div>
              </td>

              <!-- Cost -->
              <td class="py-3.5 px-4 text-right font-mono text-slate-600 dark:text-slate-300">
                Rp {{ formatNum(item.cost_per_unit) }}
              </td>

              <!-- Total Value -->
              <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                Rp {{ formatNum(item.current_stock * item.cost_per_unit) }}
              </td>

              <!-- Status Badge -->
              <td class="py-3.5 px-4 text-center">
                <StatusBadgeV2 
                  :status="item.status === 'Normal' ? 'active' : item.status === 'Low Stock' ? 'pending' : 'cancelled'"
                  :label="item.status"
                />
              </td>

              <!-- Actions -->
              <td class="py-3.5 px-4 text-center">
                <div class="flex items-center justify-center gap-1.5">
                  <router-link
                    :to="`/inventory/stock-card?sku=${item.sku}`"
                    class="p-1.5 rounded-lg border text-slate-500 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-950/40 transition-colors"
                    title="Lihat Mutasi Kartu Stok"
                  >
                    <FileSpreadsheet class="w-3.5 h-3.5" />
                  </router-link>

                  <router-link
                    to="/inventory/po"
                    class="p-1.5 rounded-lg border text-slate-500 hover:text-amber-600 hover:bg-amber-50 dark:hover:bg-amber-950/40 transition-colors"
                    title="Buat Purchase Order"
                  >
                    <ShoppingCart class="w-3.5 h-3.5" />
                  </router-link>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import PageHeaderV2 from '@/components/v2/common/PageHeaderV2.vue'
import StatCardV2 from '@/components/v2/common/StatCardV2.vue'
import StatusBadgeV2 from '@/components/v2/common/StatusBadgeV2.vue'
import {
  Plus,
  FileText,
  ClipboardCheck,
  Package,
  AlertTriangle,
  FileSpreadsheet,
  Boxes,
  Search,
  ShoppingCart
} from 'lucide-vue-next'

const searchQuery = ref('')
const filterCategory = ref('')
const filterWarehouse = ref('')
const filterStatus = ref('')

const availableCategories = ['Biji Kopi (Beans)', 'Dairy & Susu', 'Sirup & Perisa', 'Pastry & Bakery', 'Kemasan & Cup']

interface InventoryItem {
  id: string
  sku: string
  name: string
  category: string
  warehouse: string
  current_stock: number
  min_stock: number
  unit: string
  cost_per_unit: number
  status: 'Normal' | 'Low Stock' | 'Out of Stock'
}

const items = ref<InventoryItem[]>([
  { id: '1', sku: 'RM-COF-001', name: 'Arabica Gayo Natural Roast', category: 'Biji Kopi (Beans)', warehouse: 'Main Store Senopati', current_stock: 18.5, min_stock: 5, unit: 'kg', cost_per_unit: 185000, status: 'Normal' },
  { id: '2', sku: 'RM-COF-002', name: 'Robusta Temanggung Fine', category: 'Biji Kopi (Beans)', warehouse: 'Main Store Senopati', current_stock: 12.0, min_stock: 5, unit: 'kg', cost_per_unit: 95000, status: 'Normal' },
  { id: '3', sku: 'RM-DRY-001', name: 'Fresh Milk Greenfields 1L', category: 'Dairy & Susu', warehouse: 'Main Store Senopati', current_stock: 4, min_stock: 12, unit: 'karton', cost_per_unit: 215000, status: 'Low Stock' },
  { id: '4', sku: 'RM-DRY-002', name: 'Oatside Barista Blend 1L', category: 'Dairy & Susu', warehouse: 'Main Store Senopati', current_stock: 2, min_stock: 8, unit: 'karton', cost_per_unit: 345000, status: 'Low Stock' },
  { id: '5', sku: 'RM-SRP-001', name: 'Monin Salted Caramel 700ml', category: 'Sirup & Perisa', warehouse: 'Central Kitchen', current_stock: 6, min_stock: 2, unit: 'botol', cost_per_unit: 165000, status: 'Normal' },
  { id: '6', sku: 'RM-SRP-002', name: 'Monin Vanilla Syrup 700ml', category: 'Sirup & Perisa', warehouse: 'Central Kitchen', current_stock: 0, min_stock: 3, unit: 'botol', cost_per_unit: 165000, status: 'Out of Stock' },
  { id: '7', sku: 'RM-PKG-001', name: 'Paper Cup Hot 8oz Cafe ERP', category: 'Kemasan & Cup', warehouse: 'Main Store Senopati', current_stock: 450, min_stock: 100, unit: 'pcs', cost_per_unit: 850, status: 'Normal' },
  { id: '8', sku: 'RM-PST-001', name: 'Frozen Croissant Dough 60g', category: 'Pastry & Bakery', warehouse: 'Central Kitchen', current_stock: 85, min_stock: 30, unit: 'pcs', cost_per_unit: 12500, status: 'Normal' }
])

const countLowStock = computed(() => items.value.filter(i => i.status === 'Low Stock' || i.status === 'Out of Stock').length)

const totalValuation = computed(() => {
  return items.value.reduce((acc, item) => acc + (item.current_stock * item.cost_per_unit), 0)
})

const filteredItems = computed(() => {
  return items.value.filter(item => {
    const matchQuery = !searchQuery.value.trim() ||
      item.name.toLowerCase().includes(searchQuery.value.toLowerCase().trim()) ||
      item.sku.toLowerCase().includes(searchQuery.value.toLowerCase().trim())
    const matchCat = !filterCategory.value || item.category === filterCategory.value
    const matchWh = !filterWarehouse.value || item.warehouse === filterWarehouse.value
    const matchStatus = !filterStatus.value || item.status === filterStatus.value
    return matchQuery && matchCat && matchWh && matchStatus
  })
})

const resetFilters = () => {
  searchQuery.value = ''
  filterCategory.value = ''
  filterWarehouse.value = ''
  filterStatus.value = ''
}

const formatNum = (val: number | string) => Number(val || 0).toLocaleString('id-ID')
</script>

