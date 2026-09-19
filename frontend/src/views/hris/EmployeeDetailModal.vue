<template>
  <Teleport to="body">
    <Transition name="emp-modal-fade">
      <div v-if="show && employeeData" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center bg-slate-950/45 backdrop-blur-sm p-4 overflow-y-auto transition-all duration-200">
        <div class="bg-white rounded-3xl max-w-4xl w-full p-6 shadow-2xl shadow-slate-950/25 border border-slate-100/90 my-8 transform transition-all duration-300 ease-out">
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

          <!-- Top Summary Grid -->
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
            <!-- Left: Profile Header (2 cols) -->
            <div class="lg:col-span-2 p-5 rounded-2xl border border-slate-200/80 bg-slate-50/50 flex items-center gap-5">
              <div class="w-20 h-20 rounded-full overflow-hidden bg-blue-100 border-2 border-blue-500 shrink-0 flex items-center justify-center text-blue-700 text-2xl font-black">
                <img v-if="employeeData.avatar" :src="employeeData.avatar" class="w-full h-full object-cover" />
                <span v-else>{{ (employeeData.full_name || employeeData.name || 'K').charAt(0) }}</span>
              </div>
              <div>
                <div class="flex items-center gap-3">
                  <h3 class="text-xl font-black text-slate-900">{{ employeeData.full_name || employeeData.name }}</h3>
                  <span
                    class="px-2.5 py-0.5 rounded-full text-xs font-bold capitalize"
                    :class="isActiveStatus(employeeData.status) ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-800'"
                  >
                    {{ employeeData.status }}
                  </span>
                </div>
                <div class="text-xs text-slate-600 font-medium mt-1 space-y-0.5">
                  <div>Posisi: <strong class="text-slate-900">{{ employeeData.position_name || employeeData.position || '-' }}</strong></div>
                  <div>Department: <strong class="text-slate-900">{{ employeeData.department_name || employeeData.department || '-' }}</strong></div>
                  <div>ID Pegawai: <strong class="text-slate-900 font-mono">{{ employeeData.employee_number || employeeData.id }}</strong></div>
                </div>
              </div>
            </div>

            <!-- Right: Metrics Cards (1 col) -->
            <div class="grid grid-cols-2 lg:grid-cols-1 gap-3">
              <div class="p-3.5 rounded-xl border border-slate-200 bg-white flex items-center justify-between">
                <div>
                  <div class="text-[11px] text-slate-400 font-semibold">Sisa Cuti Tahunan</div>
                  <div class="text-xl font-black text-slate-900">{{ employeeData.remaining_leave ?? employeeData.remainingLeave ?? 12 }} hari</div>
                </div>
                <CalendarOff class="w-6 h-6 text-amber-500" />
              </div>
              <div class="p-3.5 rounded-xl border border-slate-200 bg-white flex items-center justify-between">
                <div>
                  <div class="text-[11px] text-slate-400 font-semibold">Tipe Pekerjaan</div>
                  <div class="text-base font-black text-blue-600">{{ employeeData.employment_type || 'Full-time' }}</div>
                </div>
                <Calendar class="w-6 h-6 text-blue-500" />
              </div>
            </div>
          </div>

          <!-- Navigation Tabs -->
          <div class="flex items-center gap-6 border-b border-slate-200 text-xs font-bold mb-5 pb-1">
            <button class="pb-2 text-blue-600 border-b-2 border-blue-600 cursor-pointer">Informasi Personal & Payroll</button>
          </div>

          <!-- Informasi Details Grid -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-y-4 gap-x-8 text-xs font-medium text-slate-700">
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-400">Nomor Induk Kependudukan (NIK):</span>
              <span class="font-bold text-slate-900 font-mono">{{ employeeData.nik || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-400">Cabang Penempatan:</span>
              <span class="font-bold text-slate-900">{{ employeeData.branch_name || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-400">Email Kerja:</span>
              <span class="font-bold text-slate-900">{{ employeeData.email || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-400">No. Telepon / WhatsApp:</span>
              <span class="font-bold text-slate-900">{{ employeeData.phone || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-400">Tanggal Mulai Bekerja:</span>
              <span class="font-bold text-slate-900">{{ employeeData.join_date ? formatDate(employeeData.join_date) : '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-400">Gaji Pokok:</span>
              <span class="font-bold text-emerald-600 font-mono">{{ employeeData.basic_salary ? 'Rp ' + Number(employeeData.basic_salary).toLocaleString('id-ID') : '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-400">Nama Bank:</span>
              <span class="font-bold text-slate-900">{{ employeeData.bank_name || employeeData.bank || '-' }}</span>
            </div>
            <div class="flex justify-between py-2 border-b border-slate-100">
              <span class="text-slate-400">No. Rekening Bank:</span>
              <span class="font-mono font-bold text-slate-900">{{ employeeData.bank_account || employeeData.bankAccount || '-' }}</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
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
  return props.employee || null
})

const isActiveStatus = (status?: string) => {
  if (!status) return false
  const s = status.toLowerCase()
  return s === 'active' || s === 'aktif'
}

const formatDate = (dateStr: string) => {
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
  } catch {
    return dateStr
  }
}
</script>

<style scoped>
.emp-modal-fade-enter-active,
.emp-modal-fade-leave-active {
  transition: opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.emp-modal-fade-enter-active > div,
.emp-modal-fade-leave-active > div {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.emp-modal-fade-enter-from,
.emp-modal-fade-leave-to {
  opacity: 0;
}

.emp-modal-fade-enter-from > div,
.emp-modal-fade-leave-to > div {
  opacity: 0;
  transform: scale(0.94) translateY(8px);
}
</style>

