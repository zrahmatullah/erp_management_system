<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Manajemen Meja & Denah</h2>
        <p class="text-xs text-slate-500 mt-1">Pantau okupansi dan tata letak meja operasional secara real-time</p>
      </div>
      <button
        @click="openAddTableModal"
        class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
      >
        <span>+</span> Tambah Meja
      </button>
    </div>

    <!-- Zone Tabs -->
    <div class="flex items-center gap-6 border-b border-slate-200 text-sm font-bold">
      <button
        v-for="z in ['Semua Meja', 'Lantai 1', 'Lantai 2', 'Outdoor']"
        :key="z"
        @click="activeZone = z"
        class="pb-3 transition-colors relative cursor-pointer"
        :class="activeZone === z ? 'text-blue-600' : 'text-slate-500 hover:text-slate-800'"
      >
        {{ z }}
        <span v-if="activeZone === z" class="absolute bottom-0 inset-x-0 h-0.5 bg-blue-600 rounded-full"></span>
      </button>
    </div>

    <!-- Interactive Floor Plan Grid -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-8 min-h-[460px] relative overflow-hidden">
      <div v-if="loading" class="py-20 text-center text-slate-400 text-xs">
        Memuat data meja dari database server...
      </div>
      <div v-else-if="filteredTables.length === 0" class="py-20 text-center text-slate-400 text-xs">
        Belum ada data meja terdaftar di zona ini.
      </div>
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-6">
        <div
          v-for="table in filteredTables"
          :key="table.id || table.code"
          @click="selectTable(table)"
          class="p-4 rounded-2xl border-2 transition-all cursor-pointer flex flex-col items-center justify-center relative select-none hover:shadow-lg hover:scale-102"
          :class="getTableClass(table.status)"
        >
          <!-- Table Code -->
          <div class="text-sm font-extrabold">{{ table.code }}</div>

          <!-- Order Tag if Occupied -->
          <div v-if="table.status === 'occupied'" class="flex flex-col items-center mt-1 text-center">
            <span class="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-red-700 text-white shadow-xs">
              {{ table.orderNumber || 'Terisi' }}
            </span>
            <span v-if="table.customer" class="text-[10px] font-bold text-white/95 truncate max-w-[120px] mt-0.5">
              {{ table.customer }}
            </span>
            <span v-if="table.total" class="text-[10px] font-extrabold text-white/90">
              Rp {{ Number(table.total).toLocaleString('id-ID') }}
            </span>
          </div>

          <!-- Order Tag if Billing -->
          <div v-else-if="table.status === 'billing'" class="flex flex-col items-center mt-1 text-center">
            <span class="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-blue-600 text-white shadow-xs">
              {{ table.orderNumber || 'Billing' }}
            </span>
            <span v-if="table.customer" class="text-[10px] font-bold text-blue-900 truncate max-w-[120px] mt-0.5">
              {{ table.customer }}
            </span>
            <span v-if="table.total" class="text-[10px] font-extrabold text-blue-700">
              Rp {{ Number(table.total).toLocaleString('id-ID') }}
            </span>
          </div>

          <!-- Time if Reserved -->
          <div v-else-if="table.status === 'reserved'" class="mt-1 flex flex-col items-center">
            <span class="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-amber-200 text-amber-900">
              Reservasi 19:00
            </span>
          </div>

          <!-- Seats Count -->
          <div class="flex items-center gap-1.5 mt-2 text-[11px] font-semibold opacity-90">
            <Users class="w-3.5 h-3.5" />
            <span>{{ table.capacity }} Kursi</span>
          </div>

          <!-- Hover Action Hint -->
          <div class="mt-1 text-[9px] font-bold uppercase tracking-wider opacity-75">
            Klik untuk Detail / Pesan
          </div>
        </div>
      </div>

      <!-- Legend floating box -->
      <div class="absolute bottom-6 right-6 bg-white/95 backdrop-blur-xs p-4 rounded-2xl border border-slate-200 shadow-lg text-xs space-y-2">
        <div class="font-bold text-slate-800 mb-1">Status Meja</div>
        <div class="flex items-center gap-2">
          <span class="w-3.5 h-3.5 rounded border-2 border-emerald-500 bg-emerald-50"></span>
          <span class="text-slate-600 font-medium">Available (Tersedia)</span>
        </div>
        <div class="flex items-center gap-2">
          <span class="w-3.5 h-3.5 rounded bg-red-500 border-2 border-red-600"></span>
          <span class="text-slate-600 font-medium">Occupied (Terisi)</span>
        </div>
        <div class="flex items-center gap-2">
          <span class="w-3.5 h-3.5 rounded bg-amber-200 border-2 border-amber-400"></span>
          <span class="text-slate-600 font-medium">Reserved (Reservasi)</span>
        </div>
        <div class="flex items-center gap-2">
          <span class="w-3.5 h-3.5 rounded border-2 border-blue-500 bg-blue-50"></span>
          <span class="text-slate-600 font-medium">Billing</span>
        </div>
      </div>
    </div>

    <!-- Bottom Status Counter Bar -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-4 flex flex-wrap items-center justify-around text-xs font-medium">
      <div class="text-slate-700">
        Total <span class="font-black text-slate-900 text-sm ml-1">{{ tables.length }} meja</span>
      </div>
      <div class="text-slate-700">
        Tersedia <span class="font-black text-emerald-600 text-sm ml-1">{{ countAvailable }}</span>
      </div>
      <div class="text-slate-700">
        Terisi <span class="font-black text-red-600 text-sm ml-1">{{ countOccupied }}</span>
      </div>
      <div class="text-slate-700">
        Reservasi <span class="font-black text-amber-500 text-sm ml-1">{{ countReserved }}</span>
      </div>
      <div class="text-slate-700">
        Billing <span class="font-black text-blue-600 text-sm ml-1">{{ countBilling }}</span>
      </div>
    </div>

    <!-- Standardized Add Table Modal -->
    <AppModal
      :show="showAddModal"
      title="Tambah Meja Baru"
      subtitle="Daftarkan meja operasional baru ke sistem cafe"
      :icon="Armchair"
      max-width="max-w-md"
      @close="showAddModal = false"
    >
      <form @submit.prevent="saveNewTable" class="space-y-4">
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1.5">Nomor / Kode Meja</label>
          <input
            v-model="newTable.number"
            type="text"
            required
            placeholder="e.g. T-13, VIP-02"
            class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-2 focus:ring-blue-600/20 outline-none font-bold text-slate-900 transition-all"
            autofocus
          />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1.5">Kapasitas Kursi</label>
            <select
              v-model.number="newTable.capacity"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium text-slate-800"
            >
              <option :value="2">2 Kursi</option>
              <option :value="4">4 Kursi</option>
              <option :value="6">6 Kursi</option>
              <option :value="8">8 Kursi</option>
              <option :value="10">10 Kursi</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-700 mb-1.5">Zona Cafe</label>
            <select
              v-model="newTable.zone"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium text-slate-800"
            >
              <option value="Lantai 1">Lantai 1</option>
              <option value="Lantai 2">Lantai 2</option>
              <option value="Outdoor">Outdoor</option>
            </select>
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1.5">Status Awal</label>
          <select
            v-model="newTable.status"
            class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium text-slate-800"
          >
            <option value="available">Tersedia (Available)</option>
            <option value="reserved">Reservasi (Reserved)</option>
          </select>
        </div>

        <div class="flex items-center justify-end gap-2.5 pt-4 mt-6 border-t border-slate-100">
          <button
            type="button"
            @click="showAddModal = false"
            class="px-4 py-2.5 text-xs font-bold text-slate-600 bg-slate-100 hover:bg-slate-200 rounded-xl transition-colors cursor-pointer"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="saving"
            class="px-5 py-2.5 text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer disabled:opacity-50"
          >
            {{ saving ? 'Menyimpan...' : 'Simpan Meja' }}
          </button>
        </div>
      </form>
    </AppModal>

    <!-- Interactive Table Detail & Order Modal -->
    <TableDetailModal
      v-if="showTableDetailModal && selectedTableForModal"
      :show="showTableDetailModal"
      :table="selectedTableForModal"
      @close="showTableDetailModal = false"
      @updated="fetchTables"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { Users, Armchair } from 'lucide-vue-next'
import AppModal from '@/components/common/AppModal.vue'
import TableDetailModal from './TableDetailModal.vue'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const activeZone = ref('Semua Meja')
const loading = ref(false)
const saving = ref(false)
const tables = ref<any[]>([])

const showAddModal = ref(false)
const showTableDetailModal = ref(false)
const selectedTableForModal = ref<any>(null)

const newTable = ref({
  number: '',
  capacity: 4,
  zone: 'Lantai 1',
  status: 'available'
})

const fetchTables = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/pos/tables')
    if (res.data?.data) {
      tables.value = res.data.data.map((t: any) => ({
        id: t.id,
        code: t.table_number || t.code,
        capacity: t.capacity || 4,
        status: t.status || 'available',
        zone: t.zone_name || 'Lantai 1',
        orderNumber: t.active_order?.order_number,
        customer: t.active_order?.customer,
        total: t.active_order?.total,
        orderTime: t.active_order?.time
      }))
    }
  } catch (err: any) {
    console.error('Failed to load tables:', err)
    notifyStore.error('Gagal mengambil data meja dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchTables()
})

const filteredTables = computed(() => {
  if (activeZone.value === 'Semua Meja') return tables.value
  return tables.value.filter(t => (t.zone || '').toLowerCase().includes(activeZone.value.toLowerCase()))
})

const countAvailable = computed(() => tables.value.filter(t => t.status === 'available').length)
const countOccupied = computed(() => tables.value.filter(t => t.status === 'occupied').length)
const countReserved = computed(() => tables.value.filter(t => t.status === 'reserved').length)
const countBilling = computed(() => tables.value.filter(t => t.status === 'billing').length)

const getTableClass = (status: string) => {
  switch (status) {
    case 'available':
      return 'border-emerald-500 bg-white text-emerald-800'
    case 'occupied':
      return 'border-red-600 bg-red-500 text-white shadow-sm'
    case 'reserved':
      return 'border-amber-400 bg-amber-100 text-amber-900'
    case 'billing':
      return 'border-blue-500 bg-blue-50 text-blue-800'
    default:
      return 'border-slate-300 bg-white'
  }
}

const selectTable = (table: any) => {
  selectedTableForModal.value = table
  showTableDetailModal.value = true
}

const openAddTableModal = () => {
  newTable.value = {
    number: `T-${tables.value.length + 1}`,
    capacity: 4,
    zone: 'Lantai 1',
    status: 'available'
  }
  showAddModal.value = true
}

const saveNewTable = async () => {
  if (!newTable.value.number) return
  saving.value = true
  try {
    await axios.post('/api/v1/master/tables', {
      table_number: newTable.value.number,
      capacity: newTable.value.capacity,
      status: newTable.value.status
    })
    notifyStore.success(`Meja ${newTable.value.number} berhasil ditambahkan ke database!`, 'Meja Ditambahkan')
    showAddModal.value = false
    await fetchTables()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.error || 'Gagal menambahkan meja', 'Kesalahan API')
  } finally {
    saving.value = false
  }
}
</script>
