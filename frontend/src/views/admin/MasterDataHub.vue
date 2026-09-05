<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-sm">
      <div>
        <div class="flex items-center gap-2 text-xs font-semibold text-blue-600 uppercase tracking-wider mb-1">
          <span>⭐ Super Admin Exclusive</span>
        </div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Data Master Management Hub</h2>
        <p class="text-xs text-slate-500 mt-0.5">Kelola seluruh entitas data master sistem secara terpusat (CRUD Lengkap)</p>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="openCreateModal"
          class="inline-flex items-center gap-2 px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah {{ activeTabLabel }}</span>
        </button>
      </div>
    </div>

    <!-- Master Navigation Tabs -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-2 overflow-x-auto">
      <div class="flex gap-1.5 min-w-max">
        <button
          v-for="tab in masterTabs"
          :key="tab.key"
          @click="activeTab = tab.key"
          class="px-3.5 py-2 rounded-xl text-xs font-bold transition-all flex items-center gap-2 cursor-pointer"
          :class="activeTab === tab.key 
            ? 'bg-blue-600 text-white shadow-sm shadow-blue-500/20' 
            : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'"
        >
          <component :is="tab.icon" class="w-4 h-4 shrink-0" />
          <span>{{ tab.label }}</span>
          <span
            v-if="tabCounts[tab.key] !== undefined"
            class="px-1.5 py-0.2 rounded-md text-[10px]"
            :class="activeTab === tab.key ? 'bg-blue-700/80 text-white' : 'bg-slate-200 text-slate-700'"
          >
            {{ tabCounts[tab.key] }}
          </span>
        </button>
      </div>
    </div>

    <!-- Main Content Area based on Active Tab -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6">
      <!-- Search & Filters -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div class="relative flex-1 max-w-md">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400">
            <Search class="w-4 h-4" />
          </span>
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="`Cari ${activeTabLabel.toLowerCase()}...`"
            class="w-full pl-9 pr-4 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-1 focus:ring-blue-600 outline-none transition-all"
          />
        </div>
        <div class="text-xs text-slate-400 font-medium">
          Menampilkan <span class="font-bold text-slate-700">{{ filteredData.length }}</span> data
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="py-16 text-center text-slate-400 text-sm flex flex-col items-center justify-center gap-2">
        <RotateCw class="w-6 h-6 animate-spin text-blue-600" />
        <span>Memuat data master...</span>
      </div>

      <!-- Table Section -->
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th v-for="col in currentColumns" :key="col.key" class="py-3 px-4">
                {{ col.label }}
              </th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-if="filteredData.length === 0">
              <td :colspan="currentColumns.length + 1" class="py-8 text-center text-slate-400">
                Belum ada data {{ activeTabLabel.toLowerCase() }}
              </td>
            </tr>
            <tr v-for="row in filteredData" :key="row.id" class="hover:bg-slate-50/80 transition-colors">
              <td v-for="col in currentColumns" :key="col.key" class="py-3 px-4">
                <!-- Status Boolean -->
                <span
                  v-if="col.key === 'is_active'"
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                  :class="row[col.key] ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-100 text-slate-500'"
                >
                  {{ row[col.key] ? 'Aktif' : 'Non-Aktif' }}
                </span>
                <!-- Currency Formatting -->
                <span v-else-if="col.key.includes('price') || col.key.includes('cost') || col.key === 'balance'" class="font-semibold text-slate-900">
                  Rp {{ formatNum(row[col.key]) }}
                </span>
                <!-- Default Display -->
                <span v-else class="text-slate-800">
                  {{ row[col.key] !== null && row[col.key] !== undefined && row[col.key] !== '' ? row[col.key] : '-' }}
                </span>
              </td>
              <td class="py-3 px-4 text-right">
                <div class="inline-flex items-center gap-1.5">
                  <button
                    @click="openEditModal(row)"
                    class="px-2.5 py-1 text-xs font-semibold text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors cursor-pointer inline-flex items-center gap-1"
                  >
                    <Edit class="w-3.5 h-3.5" /> Edit
                  </button>
                  <button
                    @click="deleteItem(row.id)"
                    class="px-2.5 py-1 text-xs font-semibold text-rose-600 bg-rose-50 hover:bg-rose-100 rounded-lg transition-colors cursor-pointer inline-flex items-center gap-1"
                  >
                    <Trash2 class="w-3.5 h-3.5" /> Hapus
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form Tambah / Edit -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-xs p-4">
      <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl border border-slate-100 animate-in fade-in zoom-in-95 duration-150">
        <div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-5">
          <h3 class="text-base font-bold text-slate-900">
            {{ isEditing ? 'Edit' : 'Tambah' }} {{ activeTabLabel }}
          </h3>
          <button @click="showModal = false" class="text-slate-400 hover:text-slate-600 p-1 rounded-lg hover:bg-slate-100 transition-colors">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="saveForm" class="space-y-4">
          <!-- Dynamic Form Fields based on Active Tab -->
          <div v-for="field in currentFormFields" :key="field.key">
            <label class="block text-xs font-semibold text-slate-700 mb-1.5">{{ field.label }}</label>

            <!-- Select Input -->
            <select
              v-if="field.type === 'select'"
              v-model="formData[field.key]"
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
            >
              <option v-for="opt in field.options" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
            </select>

            <!-- Checkbox Input -->
            <label v-else-if="field.type === 'checkbox'" class="flex items-center gap-2 cursor-pointer mt-1">
              <input v-model="formData[field.key]" type="checkbox" class="w-4 h-4 rounded text-blue-600" />
              <span class="text-xs text-slate-700">{{ field.checkboxLabel || 'Status Aktif' }}</span>
            </label>

            <!-- Number Input -->
            <input
              v-else-if="field.type === 'number'"
              v-model.number="formData[field.key]"
              type="number"
              :placeholder="field.placeholder"
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
            />

            <!-- Text Input -->
            <input
              v-else
              v-model="formData[field.key]"
              type="text"
              :required="field.required !== false"
              :placeholder="field.placeholder"
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
            />
          </div>

          <!-- Buttons -->
          <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 mt-6">
            <button
              type="button"
              @click="showModal = false"
              class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="px-5 py-2 text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 active:bg-blue-800 rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer"
            >
              {{ saving ? 'Menyimpan...' : 'Simpan Data' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import axios from 'axios'
import {
  Plus,
  Search,
  RotateCw,
  Edit,
  Trash2,
  X,
  Building2,
  User,
  Shield,
  FolderTree,
  Coffee,
  Package,
  Armchair,
  Receipt,
  Warehouse,
  Truck,
  Users,
  Briefcase,
  Clock
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const formatNum = (v: any) => Number(v || 0).toLocaleString('id-ID')

// 13 Master Tabs with Lucide Icons
const masterTabs = [
  { key: 'branches', label: 'Cabang / Outlet', icon: Building2 },
  { key: 'users', label: 'Pengguna & Akun', icon: User },
  { key: 'roles', label: 'Peran & Izin', icon: Shield },
  { key: 'categories', label: 'Kategori Menu', icon: FolderTree },
  { key: 'products', label: 'Menu & Varian', icon: Coffee },
  { key: 'inventory-items', label: 'Bahan Baku', icon: Package },
  { key: 'tables', label: 'Meja & Denah', icon: Armchair },
  { key: 'accounts', label: 'Bagan Akun (COA)', icon: Receipt },
  { key: 'warehouses', label: 'Gudang', icon: Warehouse },
  { key: 'suppliers', label: 'Supplier / Vendor', icon: Truck },
  { key: 'departments', label: 'HR Departemen', icon: Users },
  { key: 'positions', label: 'HR Jabatan', icon: Briefcase },
  { key: 'shifts', label: 'Shift Kerja', icon: Clock }
]

const activeTab = ref('branches')
const loading = ref(false)
const saving = ref(false)
const searchQuery = ref('')
const tableData = ref<any[]>([])
const tabCounts = ref<Record<string, number>>({})

const showModal = ref(false)
const isEditing = ref(false)
const formData = ref<Record<string, any>>({})
const currentEditId = ref<string | null>(null)

// Extra dropdown options
const branchOptions = ref<{label: string, value: string}[]>([])
const categoryOptions = ref<{label: string, value: string}[]>([])

const activeTabLabel = computed(() => {
  return masterTabs.find(t => t.key === activeTab.value)?.label || 'Data'
})

// Column definitions for each master table
const columnsConfig: Record<string, { label: string, key: string, type?: string, bold?: boolean }[]> = {
  branches: [
    { label: 'Kode', key: 'code', bold: true },
    { label: 'Nama Cabang', key: 'name' },
    { label: 'Telepon', key: 'phone' },
    { label: 'Alamat', key: 'address' },
    { label: 'Status', key: 'is_active' }
  ],
  users: [
    { label: 'Username', key: 'username', bold: true },
    { label: 'Nama Lengkap', key: 'full_name' },
    { label: 'Email', key: 'email' },
    { label: 'Peran (Role)', key: 'role_name' },
    { label: 'Cabang', key: 'branch_name' },
    { label: 'Status', key: 'is_active' }
  ],
  roles: [
    { label: 'Nama Peran', key: 'name', bold: true },
    { label: 'Deskripsi', key: 'description' },
    { label: 'Jumlah User', key: 'user_count' }
  ],
  categories: [
    { label: 'Nama Kategori', key: 'name', bold: true },
    { label: 'Slug', key: 'slug' },
    { label: 'Urutan', key: 'sort_order' },
    { label: 'Status', key: 'is_active' }
  ],
  products: [
    { label: 'SKU', key: 'sku', bold: true },
    { label: 'Nama Menu', key: 'name' },
    { label: 'Kategori', key: 'category_name' },
    { label: 'Stasiun', key: 'target_station' },
    { label: 'Harga Dasar', key: 'base_price', type: 'currency' },
    { label: 'Status', key: 'is_active' }
  ],
  'inventory-items': [
    { label: 'SKU', key: 'sku', bold: true },
    { label: 'Nama Bahan', key: 'name' },
    { label: 'Kategori', key: 'category' },
    { label: 'Satuan', key: 'uom' },
    { label: 'Stok Saat Ini', key: 'current_stock' },
    { label: 'Min Stok', key: 'min_stock' },
    { label: 'HPP Rata-rata', key: 'average_cost', type: 'currency' },
    { label: 'Status Stok', key: 'status' }
  ],
  tables: [
    { label: 'Nomor Meja', key: 'table_number', bold: true },
    { label: 'Cabang', key: 'branch_name' },
    { label: 'Zonasi', key: 'zone_name' },
    { label: 'Kapasitas Kursi', key: 'capacity' },
    { label: 'Status Meja', key: 'status' }
  ],
  accounts: [
    { label: 'Kode Akun', key: 'code', bold: true },
    { label: 'Nama Akun', key: 'name' },
    { label: 'Klasifikasi', key: 'account_type' },
    { label: 'Deskripsi', key: 'description' },
    { label: 'Status', key: 'is_active' }
  ],
  warehouses: [
    { label: 'Nama Gudang', key: 'name', bold: true },
    { label: 'Cabang', key: 'branch_name' },
    { label: 'Tipe', key: 'type' },
    { label: 'Alamat', key: 'address' }
  ],
  suppliers: [
    { label: 'Kode', key: 'code', bold: true },
    { label: 'Nama Supplier', key: 'name' },
    { label: 'Kontak Person', key: 'contact_person' },
    { label: 'Telepon', key: 'phone' },
    { label: 'Email', key: 'email' }
  ],
  departments: [
    { label: 'Nama Departemen', key: 'name', bold: true },
    { label: 'Deskripsi', key: 'description' }
  ],
  positions: [
    { label: 'Nama Jabatan', key: 'title', bold: true },
    { label: 'Departemen', key: 'department_name' },
    { label: 'Level', key: 'level' },
    { label: 'Gaji Pokok Acuan', key: 'base_salary', type: 'currency' }
  ],
  shifts: [
    { label: 'Nama Shift', key: 'name', bold: true },
    { label: 'Jam Masuk', key: 'start_time' },
    { label: 'Jam Pulang', key: 'end_time' },
    { label: 'Warna Tag', key: 'color' }
  ]
}

const currentColumns = computed(() => {
  return columnsConfig[activeTab.value] || []
})

// Form field definitions for Create / Edit modal
const formFieldsConfig: Record<string, any[]> = {
  branches: [
    { label: 'Kode Cabang', key: 'code', placeholder: 'e.g. B-SENO-01' },
    { label: 'Nama Cabang', key: 'name', placeholder: 'e.g. Kopi Kenangan Senopati' },
    { label: 'Alamat', key: 'address', placeholder: 'Jl. Senopati No. 45' },
    { label: 'Telepon', key: 'phone', placeholder: '021-7201234' },
    { label: 'Email', key: 'email', placeholder: 'outlet@cafe-erp.com' },
    { label: 'Status', key: 'is_active', type: 'checkbox', checkboxLabel: 'Cabang Beroperasi Aktif' }
  ],
  users: [
    { label: 'Username', key: 'username', placeholder: 'e.g. barista_budi' },
    { label: 'Email', key: 'email', placeholder: 'budi@cafe-erp.com' },
    { label: 'Nama Lengkap', key: 'full_name', placeholder: 'Budi Santoso' },
    { label: 'Nomor Telepon', key: 'phone', placeholder: '08123456789' },
    { label: 'PIN Kasir (4-6 Digit)', key: 'pin_code', placeholder: '1234' },
    { label: 'Status', key: 'is_active', type: 'checkbox', checkboxLabel: 'Pengguna Aktif' }
  ],
  categories: [
    { label: 'Nama Kategori', key: 'name', placeholder: 'e.g. Coffee, Pastry, Beverages' },
    { label: 'Slug / URL Code', key: 'slug', placeholder: 'e.g. coffee, pastry' },
    { label: 'Urutan Tampilan', key: 'sort_order', type: 'number', placeholder: '1' },
    { label: 'Status', key: 'is_active', type: 'checkbox', checkboxLabel: 'Kategori Ditampilkan di POS' }
  ],
  products: [
    { label: 'SKU Produk', key: 'sku', placeholder: 'PRD-LATTE-01' },
    { label: 'Nama Menu', key: 'name', placeholder: 'Iced Caramel Macchiato' },
    { label: 'Harga Dasar (Rp)', key: 'base_price', type: 'number', placeholder: '35000' },
    { label: 'Stasiun Pemrosesan', key: 'target_station', type: 'select', options: [
      { label: 'Barista (Minuman)', value: 'barista' },
      { label: 'Kitchen (Makanan)', value: 'kitchen' },
      { label: 'Kasir Langsung', value: 'cashier' }
    ]},
    { label: 'Status', key: 'is_active', type: 'checkbox', checkboxLabel: 'Menu Siap Dijual' }
  ],
  'inventory-items': [
    { label: 'SKU Bahan Baku', key: 'sku', placeholder: 'INV-CB-001' },
    { label: 'Nama Bahan', key: 'name', placeholder: 'Biji Kopi Arabica Gayo 1kg' },
    { label: 'Kategori Bahan', key: 'category', type: 'select', options: [
      { label: 'Biji Kopi (Coffee Beans)', value: 'coffee_beans' },
      { label: 'Susu & Dairy', value: 'dairy' },
      { label: 'Sirup & Pemanis', value: 'syrup' },
      { label: 'Packaging / Cup / Sedotan', value: 'packaging' },
      { label: 'Dry Goods / Makanan Kering', value: 'dry_goods' }
    ]},
    { label: 'Satuan Pakai (UOM)', key: 'uom', type: 'select', options: [
      { label: 'Kilogram (kg)', value: 'kg' },
      { label: 'Gram (g)', value: 'gram' },
      { label: 'Liter (L)', value: 'liter' },
      { label: 'Mililiter (ml)', value: 'ml' },
      { label: 'Pieces (pcs)', value: 'pcs' },
      { label: 'Pack / Kotak', value: 'pack' }
    ]},
    { label: 'Minimum Reorder Stock', key: 'min_stock', type: 'number', placeholder: '10' },
    { label: 'Rata-rata Harga Beli / HPP (Rp)', key: 'average_cost', type: 'number', placeholder: '120000' },
    { label: 'Status', key: 'is_active', type: 'checkbox', checkboxLabel: 'Bahan Aktif Digunakan' }
  ],
  tables: [
    { label: 'Nomor Meja', key: 'table_number', placeholder: 'T-01, VIP-1' },
    { label: 'Kapasitas Kursi', key: 'capacity', type: 'number', placeholder: '4' },
    { label: 'Status Meja', key: 'status', type: 'select', options: [
      { label: 'Tersedia (Available)', value: 'available' },
      { label: 'Terisi (Occupied)', value: 'occupied' },
      { label: 'Reservasi (Reserved)', value: 'reserved' }
    ]}
  ],
  accounts: [
    { label: 'Kode Akun', key: 'code', placeholder: '1101' },
    { label: 'Nama Akun', key: 'name', placeholder: 'Kas Kasir Senopati' },
    { label: 'Klasifikasi Akun', key: 'account_type', type: 'select', options: [
      { label: 'Asset (Aset / Harta)', value: 'asset' },
      { label: 'Liability (Kewajiban / Hutang)', value: 'liability' },
      { label: 'Equity (Ekuitas / Modal)', value: 'equity' },
      { label: 'Revenue (Pendapatan)', value: 'revenue' },
      { label: 'Expense (Beban / Biaya)', value: 'expense' }
    ]},
    { label: 'Deskripsi Akun', key: 'description', placeholder: 'Keterangan fungsi akun' },
    { label: 'Status', key: 'is_active', type: 'checkbox', checkboxLabel: 'Akun Aktif' }
  ]
}

const currentFormFields = computed(() => {
  return formFieldsConfig[activeTab.value] || [
    { label: 'Nama Data', key: 'name', placeholder: 'Masukkan nama' }
  ]
})

// Filtered data by search
const filteredData = computed(() => {
  if (!searchQuery.value) return tableData.value
  const q = searchQuery.value.toLowerCase()
  return tableData.value.filter(item => {
    return Object.values(item).some(val => String(val).toLowerCase().includes(q))
  })
})

// Fetch data from API
const fetchData = async () => {
  loading.value = true
  try {
    const res = await axios.get(`/api/v1/master/${activeTab.value}`)
    tableData.value = Array.isArray(res.data) ? res.data : []
    tabCounts.value[activeTab.value] = tableData.value.length
  } catch (err) {
    console.error('Failed to load master data:', err)
    tableData.value = []
  } finally {
    loading.value = false
  }
}

watch(activeTab, () => {
  searchQuery.value = ''
  fetchData()
})

onMounted(() => {
  fetchData()
})

// Open create modal
const openCreateModal = () => {
  isEditing.value = false
  currentEditId.value = null
  const initial: Record<string, any> = { is_active: true }
  currentFormFields.value.forEach(f => {
    if (f.key !== 'is_active') {
      initial[f.key] = f.type === 'number' ? 0 : f.type === 'select' ? f.options?.[0]?.value : ''
    }
  })
  formData.value = initial
  showModal.value = true
}

// Open edit modal
const openEditModal = (row: any) => {
  isEditing.value = true
  currentEditId.value = row.id
  formData.value = { ...row }
  showModal.value = true
}

// Save form
const saveForm = async () => {
  saving.value = true
  try {
    if (isEditing.value && currentEditId.value) {
      await axios.put(`/api/v1/master/${activeTab.value}/${currentEditId.value}`, formData.value)
    } else {
      await axios.post(`/api/v1/master/${activeTab.value}`, formData.value)
    }
    showModal.value = false
    notifyStore.success(`Data ${activeTabLabel.value} berhasil disimpan!`, 'Master Data Disimpan')
    await fetchData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.error || err.message || 'Gagal menyimpan data master', 'Kesalahan Simpan')
  } finally {
    saving.value = false
  }
}

// Delete item
const deleteItem = async (id: string) => {
  if (!confirm('Apakah Anda yakin ingin menghapus data master ini?')) return
  try {
    await axios.delete(`/api/v1/master/${activeTab.value}/${id}`)
    notifyStore.success(`Data ${activeTabLabel.value} berhasil dihapus!`, 'Master Data Dihapus')
    await fetchData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.error || err.message || 'Gagal menghapus data master', 'Kesalahan Hapus')
  }
}
</script>

