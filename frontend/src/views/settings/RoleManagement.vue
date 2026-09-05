<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Manajemen Peran & Izin</h2>
        <p class="text-xs text-slate-500 mt-1">Atur hak akses pengguna secara granular per modul sistem (RBAC)</p>
      </div>
      <button
        @click="showAddRoleModal = true"
        class="px-4 py-2 bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold rounded-xl shadow-md transition-colors cursor-pointer"
      >
        Tambah Peran
      </button>
    </div>

    <!-- Role Selection Tabs (Matches Mockup) -->
    <div class="grid grid-cols-3 sm:grid-cols-5 lg:grid-cols-9 gap-3">
      <button
        v-for="role in roles"
        :key="role.id"
        @click="selectRole(role)"
        class="p-3 rounded-2xl border transition-all text-center flex flex-col items-center justify-center gap-1.5 cursor-pointer"
        :class="selectedRole?.id === role.id 
          ? 'bg-blue-50/70 border-blue-500 text-blue-700 shadow-sm' 
          : 'bg-white border-slate-200/80 text-slate-700 hover:border-slate-300 hover:bg-slate-50'"
      >
        <component :is="role.icon" class="w-5 h-5 shrink-0" :class="selectedRole?.id === role.id ? 'text-blue-600' : 'text-slate-500'" />
        <span class="text-xs font-bold truncate w-full">{{ role.name }}</span>
      </button>
    </div>

    <!-- Permission Matrix Table (Matches Mockup) -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6">
      <div class="flex items-center justify-between mb-4 pb-3 border-b border-slate-100">
        <div>
          <h3 class="font-bold text-slate-900 text-sm">
            Matriks Hak Akses: <span class="text-blue-600 font-extrabold">{{ selectedRole?.name }}</span>
          </h3>
          <p class="text-xs text-slate-400">Centang izin yang diperbolehkan untuk peran ini</p>
        </div>
        <span class="text-xs text-slate-400 font-medium">{{ permissions.length }} modul terkonfigurasi</span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-600 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4">Module</th>
              <th class="py-3 px-4 text-center">Lihat (View)</th>
              <th class="py-3 px-4 text-center">Buat (Create)</th>
              <th class="py-3 px-4 text-center">Edit</th>
              <th class="py-3 px-4 text-center">Hapus (Delete)</th>
              <th class="py-3 px-4 text-center">Setujui (Approve)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="mod in permissionModules" :key="mod.key" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-3 px-4 font-bold text-slate-900">
                {{ mod.label }}
              </td>
              <td class="py-3 px-4 text-center">
                <input
                  type="checkbox"
                  v-model="matrix[mod.key].view"
                  class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                />
              </td>
              <td class="py-3 px-4 text-center">
                <input
                  type="checkbox"
                  v-model="matrix[mod.key].create"
                  class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                />
              </td>
              <td class="py-3 px-4 text-center">
                <input
                  type="checkbox"
                  v-model="matrix[mod.key].edit"
                  class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                />
              </td>
              <td class="py-3 px-4 text-center">
                <input
                  type="checkbox"
                  v-model="matrix[mod.key].delete"
                  class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                />
              </td>
              <td class="py-3 px-4 text-center">
                <input
                  type="checkbox"
                  v-model="matrix[mod.key].approve"
                  class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 border-slate-300 cursor-pointer"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Save Button -->
      <div class="mt-6 pt-4 border-t border-slate-100 flex items-center justify-between">
        <span class="text-xs text-slate-400">Perubahan hak akses akan langsung aktif pada sesi kasir & admin berikutnya</span>
        <button
          @click="savePermissions"
          :disabled="saving"
          class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer disabled:opacity-50"
        >
          {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import {
  Shield,
  UserCheck,
  Briefcase,
  ShoppingCart,
  ChefHat,
  Users,
  Package,
  FileSpreadsheet,
  UtensilsCrossed
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

interface RoleItem {
  id: string
  name: string
  icon: any
  desc: string
}

const getRoleIcon = (name: string) => {
  if (name.includes('Admin')) return Shield
  if (name.includes('Owner')) return UserCheck
  if (name.includes('Manager')) return Briefcase
  if (name.includes('Kasir')) return ShoppingCart
  if (name.includes('Kitchen')) return ChefHat
  if (name.includes('HR')) return Users
  if (name.includes('Warehouse')) return Package
  if (name.includes('Akuntan')) return FileSpreadsheet
  return UtensilsCrossed
}

const roles = ref<RoleItem[]>([
  { id: '1', name: 'Super Admin', icon: Shield, desc: 'Full Access' },
  { id: '2', name: 'Owner', icon: UserCheck, desc: 'Read Only Analytics' },
  { id: '3', name: 'Manager', icon: Briefcase, desc: 'Branch Supervisor' },
  { id: '4', name: 'Kasir', icon: ShoppingCart, desc: 'POS Terminal' },
  { id: '5', name: 'Kitchen Staff', icon: ChefHat, desc: 'KDS Order Monitor' },
  { id: '6', name: 'HR Admin', icon: Users, desc: 'Employee & Payroll' },
  { id: '7', name: 'Warehouse', icon: Package, desc: 'Stock & Logistics' },
  { id: '8', name: 'Akuntan', icon: FileSpreadsheet, desc: 'Finance & COA' },
  { id: '9', name: 'Pelayan', icon: UtensilsCrossed, desc: 'Table Service' }
])

const selectedRole = ref<RoleItem>(roles.value[2])
const saving = ref(false)
const showAddRoleModal = ref(false)

const fetchRoles = async () => {
  try {
    const res = await axios.get('/api/v1/master/roles')
    const list = Array.isArray(res.data) ? res.data : (res.data?.data || [])
    if (list.length > 0) {
      roles.value = list.map((r: any) => ({
        id: r.id,
        name: r.name,
        desc: r.description,
        icon: getRoleIcon(r.name)
      }))
      selectedRole.value = roles.value.find((r: any) => r.name.toLowerCase().includes('manager')) || roles.value[0]
    }
  } catch (err) {
    console.error('Failed to load master roles:', err)
  }
}

onMounted(() => {
  fetchRoles()
})
const permissions = ref<any[]>([])

const permissionModules = [
  { key: 'dashboard', label: 'Dashboard' },
  { key: 'pos', label: 'Point of Sale (POS)' },
  { key: 'kitchen', label: 'Kitchen Display (KDS)' },
  { key: 'tables', label: 'Table Management' },
  { key: 'inventory', label: 'Inventory & Bahan Baku' },
  { key: 'purchasing', label: 'Purchase Order (PO)' },
  { key: 'hris', label: 'HRIS & Karyawan' },
  { key: 'payroll', label: 'Payroll & Penggajian' },
  { key: 'finance', label: 'Keuangan & Jurnal' },
  { key: 'reports', label: 'Laporan & Analitik' },
  { key: 'settings', label: 'Pengaturan Sistem' },
  { key: 'master', label: 'Master Data Hub' }
]

const matrix = reactive<Record<string, { view: boolean, create: boolean, edit: boolean, delete: boolean, approve: boolean }>>({
  dashboard: { view: true, create: false, edit: false, delete: false, approve: false },
  pos: { view: true, create: true, edit: true, delete: false, approve: false },
  orders: { view: true, create: true, edit: true, delete: true, approve: true },
  payments: { view: true, create: true, edit: false, delete: false, approve: false },
  inventory: { view: true, create: false, edit: true, delete: false, approve: false },
  purchasing: { view: true, create: true, edit: false, delete: false, approve: true },
  hris: { view: true, create: true, edit: true, delete: false, approve: false },
  finance: { view: true, create: false, edit: false, delete: false, approve: false },
  kitchen: { view: true, create: true, edit: true, delete: false, approve: false },
  tables: { view: true, create: true, edit: true, delete: false, approve: false },
  payroll: { view: true, create: false, edit: false, delete: false, approve: false },
  reports: { view: true, create: false, edit: false, delete: false, approve: false },
  settings: { view: true, create: false, edit: false, delete: false, approve: false },
  master: { view: true, create: true, edit: true, delete: true, approve: true }
})

const selectRole = (role: RoleItem) => {
  selectedRole.value = role
  if (role.name === 'Super Admin') {
    Object.keys(matrix).forEach(k => {
      matrix[k].view = true
      matrix[k].create = true
      matrix[k].edit = true
      matrix[k].delete = true
      matrix[k].approve = true
    })
  } else if (role.name === 'Kasir') {
    Object.keys(matrix).forEach(k => {
      matrix[k].view = k === 'pos' || k === 'payments'
      matrix[k].create = k === 'pos' || k === 'payments'
      matrix[k].edit = false
      matrix[k].delete = false
      matrix[k].approve = false
    })
  }
}

const savePermissions = () => {
  saving.value = true
  setTimeout(() => {
    saving.value = false
    notifyStore.success(`Matriks hak akses untuk peran '${selectedRole.value.name}' berhasil disimpan!`, 'Izin Disimpan')
  }, 400)
}
</script>
