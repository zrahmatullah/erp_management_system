<template>
  <div class="hris-v2 space-y-6 animate-in fade-in duration-200">
    <!-- Page Header V2 -->
    <PageHeaderV2
      title="Manajemen Karyawan & HRIS"
      description="Pusat administrasi kepegawaian, monitoring presensi harian, permohonan cuti, dan penggajian staf."
      badge="Database SDM"
    >
      <template #actions>
        <button
          type="button"
          @click="fetchEmployees"
          class="p-2 rounded-xl border transition-colors cursor-pointer text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800"
          style="border-color: var(--border-color); background-color: var(--bg-primary);"
          title="Perbarui Data"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>

        <button
          type="button"
          @click="openAddEmployee"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Tambah Karyawan Baru</span>
        </button>
      </template>

      <!-- HRIS Navigation Tabs Slot -->
      <template #tabs>
        <div class="flex items-center gap-6 text-xs font-semibold overflow-x-auto pb-2 select-none">
          <button type="button" class="text-blue-600 dark:text-blue-400 border-b-2 border-blue-600 pb-2 -mb-2 cursor-pointer whitespace-nowrap">
            Daftar Karyawan
          </button>
          <router-link to="/hris/attendance" class="text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200 pb-2 -mb-2 whitespace-nowrap">
            Presensi & Absensi
          </router-link>
          <router-link to="/hris/leaves" class="text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200 pb-2 -mb-2 whitespace-nowrap">
            Persetujuan Cuti
          </router-link>
          <router-link to="/hris/shifts" class="text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200 pb-2 -mb-2 whitespace-nowrap">
            Jadwal Roster Shift
          </router-link>
          <router-link to="/hris/payroll" class="text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200 pb-2 -mb-2 whitespace-nowrap">
            Penggajian (Payroll)
          </router-link>
        </div>
      </template>
    </PageHeaderV2>

    <!-- 4 HRIS Stat Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
      <StatCardV2
        title="Total Karyawan"
        :value="employees.length"
        subtitle="Staf terdaftar aktif"
        color="blue"
        :icon="Users"
      />
      <StatCardV2
        title="Karyawan Aktif"
        :value="countActive"
        trend="up"
        trendValue="100%"
        subtitle="Siap bertugas"
        color="emerald"
        :icon="UserCheck"
      />
      <StatCardV2
        title="Sedang Cuti / Izin"
        :value="countOnLeave"
        subtitle="Izin cuti disetujui"
        color="amber"
        :icon="CalendarOff"
      />
      <StatCardV2
        title="Masa Percobaan"
        :value="countProbation"
        subtitle="Staf orientasi baru"
        color="indigo"
        :icon="AlertCircle"
      />
    </div>

    <!-- Filter & Search Toolbar -->
    <div 
      class="p-4 rounded-2xl border space-y-3 shadow-xs"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 text-xs">
        <div class="relative lg:col-span-2">
          <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5 pointer-events-none" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari nama karyawan, NIK, atau posisi..."
            class="w-full pl-8.5 pr-3 py-1.5 text-xs rounded-xl border outline-none transition-all focus:ring-1 focus:ring-blue-500"
            style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
          />
        </div>

        <div>
          <select
            v-model="filterDepartment"
            class="w-full px-3 py-1.5 text-xs rounded-xl border outline-none"
            style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
          >
            <option value="">Semua Departemen</option>
            <option value="Operasional Barista">Operasional Barista</option>
            <option value="Kitchen & Cook">Kitchen & Cook</option>
            <option value="Kasir & Layanan">Kasir & Layanan</option>
            <option value="Logistik Gudang">Logistik Gudang</option>
            <option value="Manajemen & Admin">Manajemen & Admin</option>
          </select>
        </div>

        <div>
          <select
            v-model="filterStatus"
            class="w-full px-3 py-1.5 text-xs rounded-xl border outline-none"
            style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
          >
            <option value="">Semua Status</option>
            <option value="Aktif">Aktif</option>
            <option value="Sedang Cuti">Sedang Cuti</option>
            <option value="Percobaan">Masa Percobaan</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Modern Clean Employee Table -->
    <div 
      class="rounded-2xl border shadow-xs overflow-hidden"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
              <th class="py-3 px-4">Nama Lengkap & NIK</th>
              <th class="py-3 px-4">Posisi & Jabatan</th>
              <th class="py-3 px-4">Departemen</th>
              <th class="py-3 px-4">Tanggal Masuk</th>
              <th class="py-3 px-4 text-center">Status</th>
              <th class="py-3 px-4 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y" style="border-color: var(--border-color);">
            <tr v-if="filteredEmployees.length === 0">
              <td colspan="6" class="py-12 text-center text-slate-400">
                Data karyawan tidak ditemukan.
              </td>
            </tr>

            <tr 
              v-for="emp in filteredEmployees" 
              :key="emp.id"
              class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
            >
              <!-- Name & Initials -->
              <td class="py-3.5 px-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-600 text-white font-bold flex items-center justify-center text-xs shrink-0 shadow-xs">
                    {{ emp.full_name.slice(0, 2).toUpperCase() }}
                  </div>
                  <div>
                    <div class="font-bold text-slate-900 dark:text-slate-100">{{ emp.full_name }}</div>
                    <div class="font-mono text-[10.5px] text-slate-400">{{ emp.nik }} • {{ emp.email }}</div>
                  </div>
                </div>
              </td>

              <!-- Position -->
              <td class="py-3.5 px-4 font-semibold text-slate-800 dark:text-slate-200">
                {{ emp.job_title }}
              </td>

              <!-- Department -->
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-md text-[11px] font-medium bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300">
                  {{ emp.department }}
                </span>
              </td>

              <!-- Join Date -->
              <td class="py-3.5 px-4 text-slate-500 font-mono text-[11px]">
                {{ emp.hire_date }}
              </td>

              <!-- Status Badge -->
              <td class="py-3.5 px-4 text-center">
                <StatusBadgeV2
                  :status="emp.status === 'Aktif' ? 'active' : emp.status === 'Sedang Cuti' ? 'on_leave' : 'pending'"
                  :label="emp.status"
                />
              </td>

              <!-- Actions -->
              <td class="py-3.5 px-4 text-center">
                <div class="flex items-center justify-center gap-1.5">
                  <button
                    type="button"
                    @click="viewEmployeeDetail(emp)"
                    class="p-1.5 rounded-lg border text-slate-500 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-950/40 transition-colors cursor-pointer"
                    title="Detail Profil & Dokumen"
                  >
                    <Eye class="w-3.5 h-3.5" />
                  </button>

                  <button
                    type="button"
                    @click="viewPayslip(emp)"
                    class="p-1.5 rounded-lg border text-slate-500 hover:text-purple-600 hover:bg-purple-50 dark:hover:bg-purple-950/40 transition-colors cursor-pointer"
                    title="Slip Gaji Payroll"
                  >
                    <Banknote class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Reused Modals -->
    <EmployeeDetailModal
      v-if="showDetailModal"
      :show="showDetailModal"
      :employee="selectedEmployee"
      @close="showDetailModal = false"
    />

    <PayslipModal
      v-if="showPayslipModal"
      :show="showPayslipModal"
      :employee="selectedEmployee"
      @close="showPayslipModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import PageHeaderV2 from '@/components/v2/common/PageHeaderV2.vue'
import StatCardV2 from '@/components/v2/common/StatCardV2.vue'
import StatusBadgeV2 from '@/components/v2/common/StatusBadgeV2.vue'
import EmployeeDetailModal from '@/views/hris/EmployeeDetailModal.vue'
import PayslipModal from '@/views/hris/PayslipModal.vue'
import {
  RotateCw,
  Plus,
  Users,
  UserCheck,
  CalendarOff,
  AlertCircle,
  Search,
  Eye,
  Banknote
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const loading = ref(false)
const searchQuery = ref('')
const filterDepartment = ref('')
const filterStatus = ref('')

const showDetailModal = ref(false)
const showPayslipModal = ref(false)
const selectedEmployee = ref<any>(null)

interface Employee {
  id: string
  nik: string
  full_name: string
  email: string
  department: string
  job_title: string
  hire_date: string
  status: 'Aktif' | 'Sedang Cuti' | 'Percobaan'
}

const employees = ref<Employee[]>([
  { id: '1', nik: 'EMP-001', full_name: 'Budi Pratama', email: 'budi.barista@cafe-erp.local', department: 'Operasional Barista', job_title: 'Head Barista', hire_date: '15 Jan 2025', status: 'Aktif' },
  { id: '2', nik: 'EMP-002', full_name: 'Jane Doe', email: 'jane.cashier@cafe-erp.local', department: 'Kasir & Layanan', job_title: 'Kasir POS Senior', hire_date: '01 Mar 2025', status: 'Aktif' },
  { id: '3', nik: 'EMP-003', full_name: 'Ahmad Gudang', email: 'ahmad.inventory@cafe-erp.local', department: 'Logistik Gudang', job_title: 'Warehouse & Logistik', hire_date: '10 Feb 2025', status: 'Aktif' },
  { id: '4', nik: 'EMP-004', full_name: 'Dewi Finance', email: 'dewi.finance@cafe-erp.local', department: 'Manajemen & Admin', job_title: 'Akuntan Keuangan', hire_date: '05 Jan 2025', status: 'Aktif' },
  { id: '5', nik: 'EMP-005', full_name: 'Rian Barista', email: 'rian@cafe-erp.local', department: 'Operasional Barista', job_title: 'Junior Barista', hire_date: '01 Sep 2026', status: 'Percobaan' },
  { id: '6', nik: 'EMP-006', full_name: 'Siti Kitchen', email: 'siti@cafe-erp.local', department: 'Kitchen & Cook', job_title: 'Pastry Chef', hire_date: '12 Apr 2025', status: 'Sedang Cuti' },
  { id: '7', nik: 'EMP-007', full_name: 'Sarah Store Manager', email: 'sarah.manager@cafe-erp.local', department: 'Manajemen & Admin', job_title: 'Store Manager', hire_date: '01 Jan 2025', status: 'Aktif' }
])

const countActive = computed(() => employees.value.filter(e => e.status === 'Aktif').length)
const countOnLeave = computed(() => employees.value.filter(e => e.status === 'Sedang Cuti').length)
const countProbation = computed(() => employees.value.filter(e => e.status === 'Percobaan').length)

const filteredEmployees = computed(() => {
  return employees.value.filter(e => {
    const matchQuery = !searchQuery.value.trim() ||
      e.full_name.toLowerCase().includes(searchQuery.value.toLowerCase().trim()) ||
      e.nik.toLowerCase().includes(searchQuery.value.toLowerCase().trim()) ||
      e.job_title.toLowerCase().includes(searchQuery.value.toLowerCase().trim())
    const matchDept = !filterDepartment.value || e.department === filterDepartment.value
    const matchStatus = !filterStatus.value || e.status === filterStatus.value
    return matchQuery && matchDept && matchStatus
  })
})

const fetchEmployees = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    notifyStore.success('Database kepegawaian berhasil diperbarui', 'HRIS Sinkron')
  }, 350)
}

const openAddEmployee = () => {
  notifyStore.info('Form pendaftaran karyawan baru dibuka', 'HRIS')
}

const viewEmployeeDetail = (emp: any) => {
  selectedEmployee.value = emp
  showDetailModal.value = true
}

const viewPayslip = (emp: any) => {
  selectedEmployee.value = emp
  showPayslipModal.value = true
}

onMounted(() => {
  fetchEmployees()
})
</script>

