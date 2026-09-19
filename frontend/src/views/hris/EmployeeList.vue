<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2.5">
          <h2 class="text-2xl font-black text-slate-900 tracking-tight">Manajemen Karyawan & Staf</h2>
          <span class="px-2.5 py-0.5 rounded-full text-xs font-black bg-blue-100 text-blue-800 border border-blue-200 shadow-xs">
            HRIS Database
          </span>
        </div>
        <p class="text-xs text-slate-500 mt-1">Kelola data seluruh staf barista, kitchen, kasir, dan manajerial secara terpusat.</p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="fetchEmployees"
          class="p-2.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
          title="Segarkan Data"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>

        <button
          @click="openAddEmployee"
          class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah Karyawan Baru</span>
        </button>
      </div>
    </div>

    <!-- 4 Stat KPI Cards (Dynamically computed from real database) -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center shrink-0">
          <Users class="w-5 h-5" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Total Karyawan</span>
          <span class="text-xl font-black text-slate-900">{{ employees.length }}</span>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0">
          <UserCheck class="w-5 h-5" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Karyawan Aktif</span>
          <span class="text-xl font-black text-emerald-600">{{ countActive }}</span>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shrink-0">
          <CalendarOff class="w-5 h-5" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Sedang Cuti / Izin</span>
          <span class="text-xl font-black text-amber-600">{{ countOnLeave }}</span>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center gap-3.5">
        <div class="w-11 h-11 rounded-xl bg-yellow-50 text-yellow-600 flex items-center justify-center shrink-0">
          <AlertCircle class="w-5 h-5" />
        </div>
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Masa Percobaan</span>
          <span class="text-xl font-black text-yellow-600">{{ countProbation }}</span>
        </div>
      </div>
    </div>

    <!-- Subnav Tabs & Filters -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs p-5 space-y-4">
      <!-- HRIS Navigation Sub-tabs -->
      <div class="flex items-center gap-6 border-b border-slate-100 text-xs font-bold pb-3 overflow-x-auto">
        <button class="text-blue-600 border-b-2 border-blue-600 pb-3 -mb-3 cursor-pointer whitespace-nowrap">Semua Karyawan</button>
        <router-link to="/hris/attendance" class="text-slate-500 hover:text-slate-800 whitespace-nowrap">Presensi & Absensi</router-link>
        <router-link to="/hris/leaves" class="text-slate-500 hover:text-slate-800 whitespace-nowrap">Pengajuan & Approval Cuti</router-link>
        <router-link to="/hris/shifts" class="text-slate-500 hover:text-slate-800 whitespace-nowrap">Jadwal Shift Kerja</router-link>
        <router-link to="/hris/payroll" class="text-slate-500 hover:text-slate-800 whitespace-nowrap">Kalkulasi Penggajian</router-link>
      </div>

      <!-- Filters Bar -->
      <div class="flex flex-col sm:flex-row items-center justify-between gap-3">
        <div class="relative w-full sm:w-80">
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari nama, NIK, jabatan..."
            class="w-full pl-9 pr-4 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden"
          />
        </div>

        <div class="flex items-center gap-2 w-full sm:w-auto overflow-x-auto">
          <!-- Dynamic Departments from database -->
          <select
            v-model="filterDept"
            class="px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl font-medium outline-hidden"
          >
            <option value="">Semua Departemen</option>
            <option v-for="d in departmentList" :key="d.id" :value="d.name">
              {{ d.name }}
            </option>
          </select>

          <!-- Status Filter -->
          <select
            v-model="filterStatus"
            class="px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl font-medium outline-hidden"
          >
            <option value="">Semua Status</option>
            <option value="active">Aktif</option>
            <option value="probation">Masa Percobaan</option>
            <option value="resigned">Resigned / Keluar</option>
          </select>
        </div>
      </div>

      <!-- Employees Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 border-b border-slate-200/80 text-[11px] uppercase font-bold text-slate-500 tracking-wider">
            <tr>
              <th class="py-3 px-4">Karyawan & NIK</th>
              <th class="py-3 px-4">Departemen</th>
              <th class="py-3 px-4">Jabatan</th>
              <th class="py-3 px-4">Gaji Pokok</th>
              <th class="py-3 px-4">Tgl Masuk</th>
              <th class="py-3 px-4 text-center">Sisa Cuti</th>
              <th class="py-3 px-4 text-center">Status</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="8" class="py-12 text-center text-slate-400">
                <RotateCw class="w-5 h-5 animate-spin mx-auto mb-2 text-blue-600" />
                Memuat data staf dari database...
              </td>
            </tr>
            <tr v-else-if="filteredEmployees.length === 0">
              <td colspan="8" class="py-12 text-center text-slate-400">
                Tidak ada data karyawan yang sesuai dengan kriteria pencarian.
              </td>
            </tr>
            <tr
              v-for="emp in filteredEmployees"
              :key="emp.id"
              class="hover:bg-slate-50/70 transition-colors"
            >
              <!-- Name & NIK -->
              <td class="py-3 px-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-full bg-blue-100 text-blue-700 font-black flex items-center justify-center text-xs shrink-0 border border-blue-200">
                    {{ emp.first_name?.charAt(0) || 'K' }}
                  </div>
                  <div>
                    <div class="font-bold text-slate-900 text-xs">{{ emp.name }}</div>
                    <div class="text-[10px] font-mono text-slate-400">{{ emp.nik }}</div>
                  </div>
                </div>
              </td>

              <!-- Department -->
              <td class="py-3 px-4 font-semibold text-slate-800">{{ emp.department }}</td>

              <!-- Position -->
              <td class="py-3 px-4 font-medium text-slate-700">{{ emp.position }}</td>

              <!-- Basic Salary -->
              <td class="py-3 px-4 font-mono font-bold text-slate-900">
                Rp {{ formatNum(emp.basic_salary) }}
              </td>

              <!-- Join Date -->
              <td class="py-3 px-4 text-slate-500 font-mono">{{ emp.join_date }}</td>

              <!-- Remaining Leave -->
              <td class="py-3 px-4 text-center">
                <span class="px-2 py-0.5 rounded-md font-bold text-[11px] bg-slate-100 text-slate-800">
                  {{ emp.remaining_leave ?? 12 }} Hari
                </span>
              </td>

              <!-- Status -->
              <td class="py-3 px-4 text-center">
                <span
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider"
                  :class="getStatusBadgeClass(emp.status)"
                >
                  {{ emp.status }}
                </span>
              </td>

              <!-- Actions -->
              <td class="py-3 px-4 text-right">
                <div class="inline-flex items-center gap-1.5">
                  <button
                    @click="viewDetail(emp)"
                    class="p-1.5 rounded-lg border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 hover:text-blue-600 transition-colors cursor-pointer"
                    title="Lihat Detail Profil Karyawan"
                  >
                    <Eye class="w-4 h-4" />
                  </button>
                  <button
                    @click="openEditModal(emp)"
                    class="p-1.5 rounded-lg border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 hover:text-amber-600 transition-colors cursor-pointer"
                    title="Ubah Data Karyawan"
                  >
                    <Edit class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- MODAL: TAMBAH / UBAH KARYAWAN                                            -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <div
        v-if="showFormModal"
        class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 bg-slate-950/45 backdrop-blur-sm flex items-center justify-center p-4 overflow-y-auto"
      >
        <div class="bg-white rounded-3xl w-full max-w-2xl shadow-2xl border border-slate-200 overflow-hidden my-6">
          <div class="p-5 border-b border-slate-100 flex items-center justify-between bg-blue-50/50">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-2xl bg-blue-600 text-white flex items-center justify-center font-black shrink-0">
                <Users class="w-5 h-5" />
              </div>
              <div>
                <h3 class="font-black text-slate-900 text-sm">
                  {{ isEditing ? 'Ubah Data Karyawan' : 'Registrasi Karyawan Baru' }}
                </h3>
                <p class="text-[11px] text-slate-500 mt-0.5">
                  Lengkapi data profil kepegawaian & penggajian staf langsung ke database.
                </p>
              </div>
            </div>
            <button
              @click="showFormModal = false"
              class="p-1.5 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors cursor-pointer"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <form @submit.prevent="saveEmployee" class="p-6 space-y-4 text-xs max-h-[75vh] overflow-y-auto">
            <!-- Names Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-slate-700 mb-1">Nama Depan <span class="text-rose-500">*</span></label>
                <input
                  v-model="empForm.first_name"
                  type="text"
                  required
                  placeholder="Contoh: Sarah"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
                />
              </div>
              <div>
                <label class="block font-bold text-slate-700 mb-1">Nama Belakang</label>
                <input
                  v-model="empForm.last_name"
                  type="text"
                  placeholder="Contoh: Johnson"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
                />
              </div>
            </div>

            <!-- NIK & Email -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-slate-700 mb-1">NIK (Nomor Induk Karyawan)</label>
                <input
                  v-model="empForm.nik"
                  type="text"
                  placeholder="EMP-00X (Otomatis jika kosong)"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 font-mono outline-none"
                />
              </div>
              <div>
                <label class="block font-bold text-slate-700 mb-1">Alamat Email</label>
                <input
                  v-model="empForm.email"
                  type="email"
                  placeholder="nama@cafe-erp.com"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
                />
              </div>
            </div>

            <!-- Department & Position -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-slate-700 mb-1">Departemen <span class="text-rose-500">*</span></label>
                <select
                  v-model="empForm.department_id"
                  required
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium"
                >
                  <option value="" disabled>Pilih Departemen...</option>
                  <option v-for="d in departmentList" :key="d.id" :value="d.id">
                    {{ d.name }}
                  </option>
                </select>
              </div>
              <div>
                <label class="block font-bold text-slate-700 mb-1">Posisi / Jabatan <span class="text-rose-500">*</span></label>
                <select
                  v-model="empForm.position_id"
                  required
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium"
                >
                  <option value="" disabled>Pilih Posisi...</option>
                  <option v-for="pos in positionList" :key="pos.id" :value="pos.id">
                    {{ pos.title }} (Rp {{ formatNum(pos.base_salary) }})
                  </option>
                </select>
              </div>
            </div>

            <!-- Salary & Status -->
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div>
                <label class="block font-bold text-slate-700 mb-1">Gaji Pokok (Rp)</label>
                <input
                  v-model.number="empForm.basic_salary"
                  type="number"
                  step="50000"
                  placeholder="Contoh: 5500000"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 font-bold outline-none"
                />
              </div>
              <div>
                <label class="block font-bold text-slate-700 mb-1">Status Kepegawaian</label>
                <select
                  v-model="empForm.status"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
                >
                  <option value="active">Active (Aktif)</option>
                  <option value="probation">Probation (Percobaan)</option>
                  <option value="resigned">Resigned</option>
                </select>
              </div>
              <div>
                <label class="block font-bold text-slate-700 mb-1">Tipe Kontrak</label>
                <select
                  v-model="empForm.employment_type"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
                >
                  <option value="full_time">Full-time</option>
                  <option value="part_time">Part-time</option>
                  <option value="contract">Contract</option>
                </select>
              </div>
            </div>

            <!-- Bank Info -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-slate-700 mb-1">Nama Bank</label>
                <input
                  v-model="empForm.bank_name"
                  type="text"
                  placeholder="Contoh: Bank BCA"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-none"
                />
              </div>
              <div>
                <label class="block font-bold text-slate-700 mb-1">Nomor Rekening Bank</label>
                <input
                  v-model="empForm.bank_account"
                  type="text"
                  placeholder="Contoh: 1234567890"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white font-mono outline-none"
                />
              </div>
            </div>

            <!-- Phone & Join Date -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-slate-700 mb-1">Nomor Telepon / WhatsApp</label>
                <input
                  v-model="empForm.phone"
                  type="text"
                  placeholder="08123456789"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-none"
                />
              </div>
              <div>
                <label class="block font-bold text-slate-700 mb-1">Tanggal Mulai Bekerja</label>
                <input
                  v-model="empForm.join_date"
                  type="date"
                  class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-none"
                />
              </div>
            </div>

            <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
              <button
                type="button"
                @click="showFormModal = false"
                class="px-4 py-2 rounded-xl border border-slate-200 bg-white text-slate-700 text-xs font-bold hover:bg-slate-100 cursor-pointer"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-5 py-2 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
              >
                <Check class="w-4 h-4" />
                <span>{{ submitting ? 'Menyimpan...' : (isEditing ? 'Simpan Perubahan' : 'Simpan Karyawan') }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Detail Modal -->
    <EmployeeDetailModal
      v-if="showDetailModal"
      :show="showDetailModal"
      :employee="selectedEmployee"
      @close="showDetailModal = false"
      @refresh="fetchEmployees"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import EmployeeDetailModal from './EmployeeDetailModal.vue'
import {
  Users,
  UserCheck,
  CalendarOff,
  AlertCircle,
  Search,
  Eye,
  Edit,
  Plus,
  RotateCw,
  X,
  Check
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

// State
const loading = ref(false)
const submitting = ref(false)
const searchQuery = ref('')
const filterDept = ref('')
const filterStatus = ref('')

const employees = ref<any[]>([])
const departmentList = ref<any[]>([])
const positionList = ref<any[]>([])
const countOnLeave = ref(0)

// Modal states
const showFormModal = ref(false)
const isEditing = ref(false)
const editingEmployeeID = ref('')
const showDetailModal = ref(false)
const selectedEmployee = ref<any>(null)

const defaultForm = () => ({
  nik: '',
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  department_id: '',
  position_id: '',
  basic_salary: 5000000,
  employment_type: 'full_time',
  status: 'active',
  bank_name: 'Bank BCA',
  bank_account: '',
  bank_account_name: '',
  join_date: new Date().toISOString().split('T')[0]
})

const empForm = ref(defaultForm())

// Metrics
const countActive = computed(() => employees.value.filter(e => e.status === 'active').length)
const countProbation = computed(() => employees.value.filter(e => e.status === 'probation').length)

// Filtered Employees
const filteredEmployees = computed(() => {
  return employees.value.filter(e => {
    const q = searchQuery.value.toLowerCase()
    const matchSearch = !q ||
      e.name?.toLowerCase().includes(q) ||
      e.nik?.toLowerCase().includes(q) ||
      e.department?.toLowerCase().includes(q) ||
      e.position?.toLowerCase().includes(q)
    const matchDept = !filterDept.value || e.department === filterDept.value
    const matchStatus = !filterStatus.value || e.status === filterStatus.value
    return matchSearch && matchDept && matchStatus
  })
})

const formatNum = (val: number) => {
  return Number(val || 0).toLocaleString('id-ID')
}

const getStatusBadgeClass = (status: string) => {
  switch (status?.toLowerCase()) {
    case 'active':
      return 'bg-emerald-100 text-emerald-800'
    case 'probation':
      return 'bg-amber-100 text-amber-800'
    case 'resigned':
      return 'bg-rose-100 text-rose-800'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

// Data fetching
const fetchEmployees = async () => {
  loading.value = true
  try {
    const [empRes, attRes] = await Promise.all([
      axios.get('/api/v1/hris/employees'),
      axios.get('/api/v1/hris/attendances/today-status')
    ])
    if (empRes.data?.data) {
      employees.value = empRes.data.data
    }
    if (attRes.data?.on_leave !== undefined) {
      countOnLeave.value = attRes.data.on_leave
    }
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal mengambil daftar karyawan dari database', 'Error')
  } finally {
    loading.value = false
  }
}

const fetchMasters = async () => {
  try {
    const [dRes, pRes] = await Promise.all([
      axios.get('/api/v1/master/departments'),
      axios.get('/api/v1/master/positions')
    ])
    departmentList.value = Array.isArray(dRes.data) ? dRes.data : (dRes.data?.data || [])
    positionList.value = Array.isArray(pRes.data) ? pRes.data : (pRes.data?.data || [])
  } catch (err) {
    console.error('Failed to load master depts/positions:', err)
  }
}

// Modal actions
const openAddEmployee = () => {
  isEditing.value = false
  editingEmployeeID.value = ''
  empForm.value = defaultForm()
  if (departmentList.value.length > 0) empForm.value.department_id = departmentList.value[0].id
  if (positionList.value.length > 0) empForm.value.position_id = positionList.value[0].id
  showFormModal.value = true
}

const openEditModal = (emp: any) => {
  isEditing.value = true
  editingEmployeeID.value = emp.id
  empForm.value = {
    nik: emp.nik,
    first_name: emp.first_name || emp.name.split(' ')[0],
    last_name: emp.last_name || emp.name.split(' ').slice(1).join(' '),
    email: emp.email,
    phone: emp.phone,
    department_id: emp.department_id,
    position_id: emp.position_id,
    basic_salary: emp.basic_salary,
    employment_type: emp.employment_type || 'full_time',
    status: emp.status || 'active',
    bank_name: emp.bank_name || 'Bank BCA',
    bank_account: emp.bank_account || '',
    bank_account_name: emp.bank_account_name || emp.name,
    join_date: emp.join_date
  }
  showFormModal.value = true
}

const viewDetail = (emp: any) => {
  selectedEmployee.value = emp
  showDetailModal.value = true
}

const saveEmployee = async () => {
  submitting.value = true
  try {
    if (isEditing.value) {
      const res = await axios.put(`/api/v1/hris/employees/${editingEmployeeID.value}`, empForm.value)
      notifyStore.success(res.data?.message || 'Data karyawan berhasil diperbarui!', 'Sukses')
    } else {
      const res = await axios.post('/api/v1/hris/employees', empForm.value)
      notifyStore.success(res.data?.message || 'Karyawan baru berhasil didaftarkan!', 'Sukses')
    }
    showFormModal.value = false
    await fetchEmployees()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan data karyawan', 'Error')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchEmployees()
  fetchMasters()
})
</script>
