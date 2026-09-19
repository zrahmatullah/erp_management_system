<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div class="flex flex-wrap items-center gap-4">
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Jadwal Shift Kerja</h2>
        <!-- Date Navigator -->
        <div class="inline-flex items-center gap-3 bg-white px-3 py-1.5 rounded-xl border border-slate-200/80 shadow-xs text-xs font-bold text-slate-700">
          <button @click="previousWeek" class="hover:text-blue-600 cursor-pointer p-1 rounded hover:bg-slate-100 transition-colors" title="Minggu Sebelumnya">
            &lt;
          </button>
          <span>{{ dateRangeLabel }}</span>
          <button @click="nextWeek" class="hover:text-blue-600 cursor-pointer p-1 rounded hover:bg-slate-100 transition-colors" title="Minggu Berikutnya">
            &gt;
          </button>
          <button @click="resetToCurrentWeek" class="ml-1 text-[10px] text-blue-600 font-extrabold hover:underline cursor-pointer">
            Minggu Ini
          </button>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="openAssignModal()"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl shadow-xs transition-colors cursor-pointer flex items-center gap-1.5"
        >
          <span>+</span> Atur Shift Staf
        </button>
      </div>
    </div>

    <!-- Main Grid: Weekly Schedule + Shift Definitions -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Weekly Schedule Table (Left 3 cols) -->
      <div class="lg:col-span-3 bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6 overflow-x-auto">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100 mb-4">
          <div class="text-xs font-bold text-slate-700">
            Roster Jadwal Tim ({{ staffRoster.length }} Staf)
          </div>
          <div class="text-[11px] text-slate-400">
            *Klik pada label shift karyawan untuk mengubah jadwal langsung ke database.
          </div>
        </div>

        <div v-if="loading" class="py-12 text-center text-slate-400 text-xs">
          Memuat jadwal shift dari database server...
        </div>
        <table v-else class="w-full text-left text-xs min-w-[650px]">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-3 w-44">Karyawan</th>
              <th v-for="day in weekDays" :key="day.iso" class="py-3 px-2 text-center">
                <div>{{ day.dayName }}</div>
                <div class="text-[10px] text-slate-400 font-normal font-mono">{{ day.formatted }}</div>
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="staff in staffRoster" :key="staff.employee_id || staff.name" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-3 px-3 font-bold text-slate-900 flex items-center gap-2.5">
                <div class="w-7 h-7 rounded-full bg-blue-100 text-blue-700 flex items-center justify-center font-bold text-[10px]">
                  {{ (staff.name || 'K').charAt(0) }}
                </div>
                <div class="truncate">
                  <div class="font-bold text-slate-900">{{ staff.name }}</div>
                  <div class="text-[10px] text-slate-400 font-normal">{{ staff.position || '-' }}</div>
                </div>
              </td>
              <td 
                v-for="(shift, idx) in staff.shifts" 
                :key="idx" 
                class="py-3 px-1.5 text-center cursor-pointer hover:bg-blue-50/50 transition-colors rounded-lg"
                @click="openAssignModal(staff, weekDays[idx].iso, shift)"
                title="Klik untuk mengubah shift ini"
              >
                <span
                  class="inline-block px-2.5 py-1 rounded-lg text-[10px] font-extrabold shadow-2xs transition-transform hover:scale-105"
                  :class="getShiftStyle(shift)"
                >
                  {{ shift }}
                </span>
              </td>
            </tr>
            <tr v-if="staffRoster.length === 0">
              <td colspan="8" class="py-8 text-center text-slate-400 text-xs">
                Tidak ada data roster staf ditemukan.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Definisi Shift (Right 1 col) -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-5 space-y-4">
        <div class="flex items-center justify-between pb-2 border-b border-slate-100">
          <h3 class="font-bold text-slate-900 text-sm">Definisi Master Shift</h3>
        </div>

        <div class="space-y-3">
          <div 
            v-for="shift in shiftDefinitions" 
            :key="shift.id"
            class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between hover:border-blue-300 transition-colors"
          >
            <div>
              <div class="font-bold text-slate-900 text-xs flex items-center gap-1.5">
                <span 
                  class="w-2.5 h-2.5 rounded-full" 
                  :style="{ backgroundColor: shift.color || '#2563eb' }"
                ></span>
                {{ shift.name }}
              </div>
              <div class="text-[11px] text-slate-500 mt-0.5 font-mono">
                {{ (shift.start_time || '07:00').substring(0, 5) }} - {{ (shift.end_time || '15:00').substring(0, 5) }}
              </div>
            </div>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-slate-100 text-slate-600">
              Operasional
            </span>
          </div>
          <div v-if="shiftDefinitions.length === 0" class="py-4 text-center text-slate-400 text-xs">
            Memuat data master shift...
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 text-[11px] text-slate-500 space-y-1">
          <div class="font-bold text-slate-700">Keterangan Label:</div>
          <div class="flex items-center gap-2"><span class="w-2.5 h-2.5 rounded bg-blue-600"></span> Shift Pagi</div>
          <div class="flex items-center gap-2"><span class="w-2.5 h-2.5 rounded bg-emerald-600"></span> Shift Siang</div>
          <div class="flex items-center gap-2"><span class="w-2.5 h-2.5 rounded bg-purple-600"></span> Shift Malam</div>
          <div class="flex items-center gap-2"><span class="w-2.5 h-2.5 rounded bg-slate-200"></span> Hari Libur (OFF)</div>
        </div>
      </div>
    </div>

    <!-- Modal Atur Shift Karyawan -->
    <Teleport to="body">
      <Transition name="emp-modal-fade">
        <div v-if="showAssignModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center bg-slate-950/45 backdrop-blur-sm p-4 overflow-y-auto">
          <div class="bg-white rounded-3xl max-w-md w-full p-6 shadow-2xl border border-slate-100/90 my-8">
            <div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-5">
              <h3 class="text-lg font-black text-slate-900 tracking-tight">Atur Jadwal Shift Staf</h3>
              <button @click="showAssignModal = false" class="text-slate-400 hover:text-slate-600 p-1.5 rounded-lg cursor-pointer">
                <X class="w-4 h-4" />
              </button>
            </div>

            <form @submit.prevent="submitAssignShift" class="space-y-4 text-xs font-medium text-slate-700">
              <div>
                <label class="block text-slate-500 mb-1">Pilih Karyawan</label>
                <select v-model="assignForm.employee_id" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 font-semibold">
                  <option value="" disabled>Pilih Staf...</option>
                  <option v-for="emp in activeEmployees" :key="emp.id" :value="emp.id">
                    {{ emp.full_name }} ({{ emp.position_name || 'Staff' }})
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-slate-500 mb-1">Tanggal Shift</label>
                <input v-model="assignForm.date" type="date" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none font-mono" />
              </div>

              <div>
                <label class="block text-slate-500 mb-1">Pilih Shift</label>
                <div class="grid grid-cols-2 gap-2">
                  <button
                    type="button"
                    v-for="s in ['Pagi', 'Siang', 'Malam', 'OFF']"
                    :key="s"
                    @click="assignForm.shift_name = s"
                    class="py-2.5 px-3 rounded-xl border text-xs font-bold transition-all cursor-pointer flex items-center justify-center gap-1.5"
                    :class="assignForm.shift_name === s ? 'border-blue-600 bg-blue-50 text-blue-700 ring-2 ring-blue-500/20' : 'border-slate-200 hover:bg-slate-50 text-slate-700'"
                  >
                    {{ s }}
                  </button>
                </div>
              </div>

              <div>
                <label class="block text-slate-500 mb-1">Catatan Tambahan (Opsional)</label>
                <input v-model="assignForm.notes" type="text" placeholder="Misal: Tukar shift dengan Sarah" class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none" />
              </div>

              <div class="flex justify-end gap-3 pt-4 border-t border-slate-100">
                <button type="button" @click="showAssignModal = false" class="px-4 py-2 border border-slate-200 rounded-xl text-slate-600 hover:bg-slate-50 cursor-pointer">
                  Batal
                </button>
                <button type="submit" :disabled="saving" class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold cursor-pointer disabled:opacity-60">
                  Simpan Jadwal
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { X } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const loading = ref(false)
const saving = ref(false)
const staffRoster = ref<any[]>([])
const shiftDefinitions = ref<any[]>([])
const activeEmployees = ref<any[]>([])

// Date calculation: Start of current week (Monday)
const getMonday = (d: Date) => {
  const date = new Date(d)
  const day = date.getDay()
  const diff = date.getDate() - day + (day === 0 ? -6 : 1)
  date.setDate(diff)
  date.setHours(0, 0, 0, 0)
  return date
}

const currentWeekStart = ref<Date>(getMonday(new Date()))

const weekDays = computed(() => {
  const days = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu']
  return days.map((dayName, idx) => {
    const d = new Date(currentWeekStart.value)
    d.setDate(d.getDate() + idx)
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const dayNum = String(d.getDate()).padStart(2, '0')
    return {
      dayName,
      formatted: `${dayNum}/${m}`,
      iso: `${y}-${m}-${dayNum}`
    }
  })
})

const dateRangeLabel = computed(() => {
  if (weekDays.value.length === 0) return ''
  const first = weekDays.value[0]
  const last = weekDays.value[6]
  const dStart = new Date(first.iso)
  const dEnd = new Date(last.iso)
  const startStr = dStart.toLocaleDateString('id-ID', { day: '2-digit', month: 'short' })
  const endStr = dEnd.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
  return `${startStr} - ${endStr}`
})

const previousWeek = () => {
  const d = new Date(currentWeekStart.value)
  d.setDate(d.getDate() - 7)
  currentWeekStart.value = d
  fetchRoster()
}

const nextWeek = () => {
  const d = new Date(currentWeekStart.value)
  d.setDate(d.getDate() + 7)
  currentWeekStart.value = d
  fetchRoster()
}

const resetToCurrentWeek = () => {
  currentWeekStart.value = getMonday(new Date())
  fetchRoster()
}

// Form state
const showAssignModal = ref(false)
const assignForm = ref({
  employee_id: '',
  shift_name: 'Pagi',
  date: '',
  notes: ''
})

const openAssignModal = (staff?: any, dateIso?: string, currentShift?: string) => {
  if (staff && staff.employee_id) {
    assignForm.value.employee_id = staff.employee_id
  } else if (activeEmployees.value.length > 0 && !assignForm.value.employee_id) {
    assignForm.value.employee_id = activeEmployees.value[0].id
  }

  assignForm.value.date = dateIso || weekDays.value[0].iso
  assignForm.value.shift_name = currentShift || 'Pagi'
  assignForm.value.notes = ''
  showAssignModal.value = true
}

const submitAssignShift = async () => {
  if (!assignForm.value.employee_id || !assignForm.value.date) {
    notifyStore.warning('Karyawan dan tanggal harus dipilih', 'Validasi')
    return
  }

  saving.value = true
  try {
    await axios.post('/api/v1/hris/schedules', {
      employee_id: assignForm.value.employee_id,
      date: assignForm.value.date,
      shift_name: assignForm.value.shift_name,
      notes: assignForm.value.notes
    })

    notifyStore.success(`Shift berhasil disimpan untuk tanggal ${assignForm.value.date}!`, 'Jadwal Diperbarui')
    showAssignModal.value = false
    fetchRoster()
  } catch (err: any) {
    console.error('Failed to assign shift:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal menyimpan jadwal shift', 'Gagal')
  } finally {
    saving.value = false
  }
}

const fetchShifts = async () => {
  try {
    const res = await axios.get('/api/v1/master/shifts')
    shiftDefinitions.value = Array.isArray(res.data) ? res.data : (res.data?.data || [])
  } catch (err) {
    console.error('Failed to load master shifts:', err)
  }
}

const fetchEmployees = async () => {
  try {
    const res = await axios.get('/api/v1/hris/employees')
    activeEmployees.value = (res.data?.data || []).filter((e: any) => e.status?.toLowerCase() === 'active')
  } catch (err) {
    console.error('Failed to load employees:', err)
  }
}

const fetchRoster = async () => {
  loading.value = true
  try {
    const startDate = weekDays.value[0]?.iso
    const endDate = weekDays.value[6]?.iso
    const res = await axios.get(`/api/v1/hris/schedules?start_date=${startDate}&end_date=${endDate}`)
    if (res.data?.data) {
      staffRoster.value = res.data.data
    }
  } catch (err: any) {
    console.error('Failed to load shift schedule:', err)
    notifyStore.error('Gagal memuat jadwal shift dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchEmployees()
  fetchShifts()
  fetchRoster()
})

const getShiftStyle = (shift: string) => {
  switch (shift) {
    case 'Pagi':
      return 'bg-blue-600 text-white'
    case 'Siang':
      return 'bg-emerald-600 text-white'
    case 'Malam':
      return 'bg-purple-600 text-white'
    default:
      return 'bg-slate-100 text-slate-500 font-medium'
  }
}
</script>

<style scoped>
.emp-modal-fade-enter-active,
.emp-modal-fade-leave-active {
  transition: opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.emp-modal-fade-enter-from,
.emp-modal-fade-leave-to {
  opacity: 0;
}
</style>

