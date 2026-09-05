<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-xs p-4 overflow-y-auto">
    <div class="bg-white rounded-2xl max-w-4xl w-full p-6 shadow-2xl border border-slate-100 my-8 animate-in fade-in zoom-in-95 duration-150">
      <!-- Header -->
      <div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-6">
        <div class="flex items-center gap-3">
          <button @click="$emit('close')" class="p-2 rounded-xl hover:bg-slate-100 text-slate-500 transition-colors cursor-pointer">
            <ArrowLeft class="w-4 h-4" />
          </button>
          <h2 class="text-xl font-black text-slate-900 tracking-tight">Detail Profil Karyawan</h2>
        </div>
        <button @click="$emit('close')" class="text-slate-400 hover:text-slate-600 p-1.5 rounded-lg hover:bg-slate-100 transition-colors cursor-pointer">
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Top Summary Grid (Matches Mockup) -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
        <!-- Left: Profile Header (2 cols) -->
        <div class="lg:col-span-2 p-5 rounded-2xl border border-slate-200/80 bg-slate-50/50 flex items-center gap-5">
          <div class="w-20 h-20 rounded-full overflow-hidden bg-blue-100 border-2 border-blue-500 shrink-0">
            <img :src="employeeData.avatar" class="w-full h-full object-cover" />
          </div>
          <div>
            <div class="flex items-center gap-3">
              <h3 class="text-xl font-black text-slate-900">{{ employeeData.name }}</h3>
              <span
                class="px-2.5 py-0.5 rounded-full text-xs font-bold"
                :class="employeeData.status === 'Active' ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-800'"
              >
                {{ employeeData.status }}
              </span>
            </div>
            <div class="text-xs text-slate-600 font-medium mt-1 space-y-0.5">
              <div>Posisi: <strong class="text-slate-900">{{ employeeData.position }}</strong></div>
              <div>Department: <strong class="text-slate-900">{{ employeeData.department }}</strong></div>
              <div>Employee ID: <strong class="text-slate-900 font-mono">{{ employeeData.id }}</strong></div>
            </div>
          </div>
        </div>

        <!-- Right: Metrics Cards (1 col) -->
        <div class="grid grid-cols-2 lg:grid-cols-1 gap-3">
          <div class="p-3.5 rounded-xl border border-slate-200 bg-white flex items-center justify-between">
            <div>
              <div class="text-[11px] text-slate-400 font-semibold">Sisa Cuti Tahunan</div>
              <div class="text-xl font-black text-slate-900">{{ employeeData.remainingLeave || 8 }} hari</div>
            </div>
            <CalendarOff class="w-6 h-6 text-amber-500" />
          </div>
          <div class="p-3.5 rounded-xl border border-slate-200 bg-white flex items-center justify-between">
            <div>
              <div class="text-[11px] text-slate-400 font-semibold">Kehadiran Bulan Ini</div>
              <div class="text-xl font-black text-emerald-600">95%</div>
            </div>
            <Calendar class="w-6 h-6 text-blue-500" />
          </div>
        </div>
      </div>

      <!-- Navigation Tabs (Matches Mockup) -->
      <div class="flex items-center gap-6 border-b border-slate-200 text-xs font-bold mb-5 pb-1">
        <button class="pb-2 text-blue-600 border-b-2 border-blue-600 cursor-pointer">Informasi</button>
        <button class="pb-2 text-slate-400 hover:text-slate-700 cursor-pointer">Riwayat Kepegawaian</button>
        <button class="pb-2 text-slate-400 hover:text-slate-700 cursor-pointer">Absensi</button>
        <button class="pb-2 text-slate-400 hover:text-slate-700 cursor-pointer">Cuti</button>
        <button class="pb-2 text-slate-400 hover:text-slate-700 cursor-pointer">Slip Gaji</button>
        <button class="pb-2 text-slate-400 hover:text-slate-700 cursor-pointer">Performa</button>
      </div>

      <!-- Informasi Details Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-y-4 gap-x-8 text-xs font-medium text-slate-700">
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">NIK:</span>
          <span class="font-bold text-slate-900 font-mono">{{ employeeData.nik || '3275012345678901' }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">NPWP:</span>
          <span class="font-bold text-slate-900 font-mono">12.345.678.9-123.000</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">Tanggal Lahir:</span>
          <span class="font-bold text-slate-900">15 Mei 1995</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">Jenis Kelamin:</span>
          <span class="font-bold text-slate-900">Perempuan</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">Status Pernikahan:</span>
          <span class="font-bold text-slate-900">Belum Kawin</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">No. Telepon:</span>
          <span class="font-bold text-slate-900">{{ employeeData.phone || '0812-3456-7890' }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">Email:</span>
          <span class="font-bold text-slate-900">{{ employeeData.email || 'employee@cafeharmony.com' }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">Nama Bank:</span>
          <span class="font-bold text-slate-900">{{ employeeData.bank || 'BCA (Bank Central Asia)' }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">No. Rekening:</span>
          <span class="font-mono font-bold text-slate-900">{{ employeeData.bankAccount || '1234567890' }}</span>
        </div>
        <div class="flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">Tanggal Bergabung:</span>
          <span class="font-bold text-slate-900">{{ employeeData.joinDate || '10 Januari 2022' }}</span>
        </div>
        <div class="sm:col-span-2 flex justify-between py-2 border-b border-slate-100">
          <span class="text-slate-400">Alamat Lengkap:</span>
          <span class="font-bold text-slate-900">Jl. Sudirman No. 45, Jakarta Selatan</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ArrowLeft, X, CalendarOff, Calendar } from 'lucide-vue-next'

const props = defineProps<{
  show: boolean
  employee?: any
}>()

defineEmits(['close'])

const employeeData = computed(() => {
  return props.employee || {
    name: 'Sarah Andini',
    id: 'EMP-003',
    department: 'Front of House (Minuman)',
    position: 'Barista',
    status: 'Active',
    avatar: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?w=300',
    remainingLeave: 8,
    nik: '3275012345678901',
    phone: '0812-3456-7890',
    email: 'sarah.andini@cafeharmony.com',
    bank: 'BCA (Bank Central Asia)',
    bankAccount: '1234567890',
    joinDate: '10 Januari 2022'
  }
})
</script>

