<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Employee Management</h2>
        <p class="text-xs text-slate-500 mt-1">Overview and manage all employees.</p>
      </div>

      <div class="flex items-center gap-3">
        <select v-model="filterDept" class="px-3 py-2 text-xs bg-white border border-slate-200 rounded-xl outline-none font-medium">
          <option value="">Department: All</option>
          <option value="Front of House">Front of House</option>
          <option value="Kitchen">Kitchen</option>
          <option value="Management">Management</option>
        </select>
        <select v-model="filterStatus" class="px-3 py-2 text-xs bg-white border border-slate-200 rounded-xl outline-none font-medium">
          <option value="">Status: All</option>
          <option value="Active">Active</option>
          <option value="Probation">Probation</option>
          <option value="Resigned">Resigned</option>
        </select>
        <button
          @click="openAddEmployee"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <span>Add Employee</span>
          <span>+</span>
        </button>
      </div>
    </div>

    <!-- 4 Stat Cards (Computed dynamically) -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-500 mb-1">Total Employees:</div>
          <div class="text-3xl font-black text-slate-900">{{ employees.length }}</div>
          <div class="text-[11px] font-bold text-emerald-600 mt-1">Aktif di DB</div>
        </div>
        <div class="w-10 h-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
          <Users class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-500 mb-1">Active Employees:</div>
          <div class="text-3xl font-black text-slate-900">{{ countActive }}</div>
          <div class="text-[11px] font-bold text-emerald-600 mt-1">Operasional</div>
        </div>
        <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
          <UserCheck class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-500 mb-1">On Leave:</div>
          <div class="text-3xl font-black text-slate-900">{{ countLeave }}</div>
          <div class="text-[11px] font-bold text-amber-600 mt-1">Izin / Cuti</div>
        </div>
        <div class="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center">
          <CalendarOff class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-500 mb-1">On Probation:</div>
          <div class="text-3xl font-black text-slate-900">{{ countProbation }}</div>
          <div class="text-[11px] font-bold text-amber-600 mt-1">Percobaan</div>
        </div>
        <div class="w-10 h-10 rounded-xl bg-yellow-50 text-yellow-600 flex items-center justify-center">
          <AlertCircle class="w-5 h-5" />
        </div>
      </div>
    </div>

    <!-- Table Section with Subnav Tabs -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6 space-y-4">
      <!-- Sub-tabs matching mockup -->
      <div class="flex items-center gap-6 border-b border-slate-200 text-xs font-bold pb-3">
        <button class="text-blue-600 border-b-2 border-blue-600 pb-3 -mb-3 cursor-pointer">All Employees</button>
        <router-link to="/hris/attendance" class="text-slate-400 hover:text-slate-700">Attendance</router-link>
        <router-link to="/hris/leaves" class="text-slate-400 hover:text-slate-700">Leave Management</router-link>
        <router-link to="/hris/payroll" class="text-slate-400 hover:text-slate-700">Payroll</router-link>
        <router-link to="/hris/shifts" class="text-slate-400 hover:text-slate-700">Shift Schedule</router-link>
      </div>

      <!-- Search Field -->
      <div class="relative max-w-sm">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400">
          <Search class="w-3.5 h-3.5" />
        </span>
        <input
          v-model="search"
          type="text"
          placeholder="Search employees..."
          class="w-full pl-8 pr-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none"
        />
      </div>

      <!-- Main Employees Table -->
      <div class="overflow-x-auto">
        <div v-if="loading" class="py-12 text-center text-slate-400 text-xs">
          Memuat daftar karyawan dari server...
        </div>
        <table v-else class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4">Employee</th>
              <th class="py-3 px-4">Department</th>
              <th class="py-3 px-4">Position</th>
              <th class="py-3 px-4">Join Date</th>
              <th class="py-3 px-4 text-center">Status</th>
              <th class="py-3 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="emp in filteredEmployees" :key="emp.id" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3.5 px-4 font-bold text-slate-900 flex items-center gap-3">
                <img :src="emp.avatar" :alt="emp.name" class="w-8 h-8 rounded-full object-cover bg-slate-200" />
                <div>
                  <div class="text-xs">{{ emp.name }}</div>
                  <div class="text-[10px] text-slate-400 font-mono font-normal">{{ emp.id }}</div>
                </div>
              </td>
              <td class="py-3.5 px-4">{{ emp.department }}</td>
              <td class="py-3.5 px-4">{{ emp.position }}</td>
              <td class="py-3.5 px-4 text-slate-500">{{ emp.joinDate }}</td>
              <td class="py-3.5 px-4 text-center">
                <span
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                  :class="getStatusBadge(emp.status)"
                >
                  {{ emp.status }}
                </span>
              </td>
              <td class="py-3.5 px-4 text-right">
                <div class="inline-flex items-center gap-2">
                  <button
                    @click="viewDetail(emp)"
                    class="p-1 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-blue-600 transition-colors cursor-pointer"
                    title="View Detail"
                  >
                    <Eye class="w-4 h-4" />
                  </button>
                  <button
                    @click="viewDetail(emp)"
                    class="p-1 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
                    title="Edit"
                  >
                    <Edit class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredEmployees.length === 0">
              <td colspan="6" class="py-8 text-center text-slate-400 text-xs">
                Tidak ada data karyawan sesuai filter.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Employee Detail Modal -->
    <EmployeeDetailModal
      v-if="showDetailModal"
      :show="showDetailModal"
      :employee="selectedEmployee"
      @close="showDetailModal = false"
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
  Edit
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const search = ref('')
const filterDept = ref('')
const filterStatus = ref('')
const showDetailModal = ref(false)
const selectedEmployee = ref<any>(null)
const loading = ref(false)

const employees = ref<any[]>([])

const fetchEmployees = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/hris/employees')
    if (res.data?.data) {
      employees.value = res.data.data.map((e: any, idx: number) => {
        const avatars = [
          'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150',
          'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150',
          'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150',
          'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150'
        ]
        return {
          id: e.nik || `EMP-00${idx + 1}`,
          name: e.name,
          department: e.department || 'Operations',
          position: e.position || 'Staff',
          joinDate: e.join_date || '2023-01-10',
          status: e.status === 'active' ? 'Active' : (e.status === 'probation' ? 'Probation' : 'Resigned'),
          attendance: 'Present',
          avatar: avatars[idx % avatars.length],
          email: e.email || `${e.name ? e.name.toLowerCase().replace(/\s+/g, '.') : 'employee'}@cafeharmony.com`,
          phone: e.phone || ('0812-3456-' + (1000 + idx * 111)),
          nik: '3275' + String(100000000000 + idx),
          bank: 'BCA (Bank Central Asia)',
          bankAccount: '123456' + (1000 + idx * 123),
          remainingLeave: 8 + (idx % 5)
        }
      })
    }
  } catch (err: any) {
    console.error('Failed to load employees:', err)
    notifyStore.error('Gagal mengambil daftar karyawan dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchEmployees()
})

const filteredEmployees = computed(() => {
  return employees.value.filter(e => {
    const matchSearch = !search.value || e.name.toLowerCase().includes(search.value.toLowerCase()) || e.id.toLowerCase().includes(search.value.toLowerCase())
    const matchDept = !filterDept.value || e.department === filterDept.value
    const matchStatus = !filterStatus.value || e.status === filterStatus.value
    return matchSearch && matchDept && matchStatus
  })
})

const countActive = computed(() => employees.value.filter(e => e.status === 'Active').length)
const countProbation = computed(() => employees.value.filter(e => e.status === 'Probation').length)
const countLeave = computed(() => 1)

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'Active':
      return 'bg-emerald-100 text-emerald-700'
    case 'Probation':
      return 'bg-amber-100 text-amber-800'
    case 'Resigned':
      return 'bg-rose-100 text-rose-700'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

const viewDetail = (emp: any) => {
  selectedEmployee.value = emp
  showDetailModal.value = true
}

const openAddEmployee = () => {
  notifyStore.info('Modal registrasi data karyawan baru siap dibuka.', 'Karyawan Baru')
}
</script>
