<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Presensi & Absensi</h2>
        <p class="text-xs text-slate-500 mt-1">Pencatatan kehadiran staf barista, kasir dan kitchen real-time terhubung ke database</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="openManualAttendanceModal"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <span>+</span> Catat Presensi Staf
        </button>
      </div>
    </div>

    <!-- Top Grid: Big Clock & Action Buttons + Today's Status -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Big Clock & Clock-In/Out (2 cols) -->
      <div class="lg:col-span-2 bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm flex flex-col justify-between">
        <div>
          <div class="text-4xl font-mono font-black text-slate-900 tracking-tight">
            {{ currentTime }}
          </div>
          <div class="text-xs font-semibold text-slate-400 mt-1">
            {{ currentDateStr }} • Shift Operasional
          </div>
        </div>

        <!-- Clock in / Clock out buttons -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-6">
          <button
            @click="handleClockIn"
            :disabled="hasClockedIn || actionLoading"
            class="p-4 rounded-2xl bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-left shadow-lg shadow-emerald-600/30 transition-all cursor-pointer flex items-center justify-between disabled:opacity-60 disabled:cursor-not-allowed"
          >
            <div>
              <div class="text-base font-extrabold flex items-center gap-1.5">
                <CheckCircle2 class="w-4 h-4" /> Clock In
              </div>
              <div class="text-xs text-emerald-100 font-medium mt-0.5">
                {{ hasClockedIn ? `Masuk: ${clockInTime}` : 'Klik untuk Mulai Shift' }}
              </div>
            </div>
            <Clock class="w-8 h-8 text-emerald-200" />
          </button>

          <button
            @click="handleClockOut"
            :disabled="!hasClockedIn || hasClockedOut || actionLoading"
            class="p-4 rounded-2xl bg-amber-600 hover:bg-amber-700 text-white font-bold text-left shadow-lg shadow-amber-600/30 transition-all cursor-pointer flex items-center justify-between disabled:opacity-60 disabled:cursor-not-allowed"
          >
            <div>
              <div class="text-base font-extrabold flex items-center gap-1.5">
                <LogOut class="w-4 h-4" /> Clock Out
              </div>
              <div class="text-xs text-amber-100 font-medium mt-0.5">
                {{ hasClockedOut ? `Pulang: ${clockOutTime}` : (hasClockedIn ? 'Klik untuk Selesai Shift' : 'Belum Clock In') }}
              </div>
            </div>
            <LogOut class="w-8 h-8 text-amber-200" />
          </button>
        </div>
      </div>

      <!-- Today's Status Card (1 col) -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm flex flex-col justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-400 mb-1">Status Kehadiran Hari Ini</div>
          <div class="text-3xl font-black text-emerald-600">
            {{ todaySummary.present > 0 ? 'Aktif Berjalan' : 'Belum Ada Hadir' }}
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3 pt-4 border-t border-slate-100 text-xs font-medium">
          <div>
            <div class="text-slate-400">Total Hadir</div>
            <div class="font-bold text-emerald-600 text-sm">{{ todaySummary.present }} Staf</div>
          </div>
          <div>
            <div class="text-slate-400">Terlambat</div>
            <div class="font-bold text-amber-600 text-sm">{{ todaySummary.late }} Staf</div>
          </div>
          <div>
            <div class="text-slate-400">Cuti / Izin</div>
            <div class="font-bold text-blue-600 text-sm">{{ todaySummary.leave }} Staf</div>
          </div>
          <div>
            <div class="text-slate-400">Total Tercatat</div>
            <div class="font-bold text-slate-700 text-sm">{{ attendanceList.length }} Orang</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Attendance Table Section -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6">
      <!-- Tabs Hari Ini / Rekap -->
      <div class="flex items-center justify-between border-b border-slate-200 pb-3 mb-5">
        <div class="text-xs font-bold text-blue-600 border-b-2 border-blue-600 pb-3 -mb-3">
          Log Presensi Staf ({{ attendanceList.length }})
        </div>
        <button
          @click="fetchAttendanceData"
          class="text-xs font-semibold text-slate-500 hover:text-slate-800 transition-colors cursor-pointer"
        >
          Perbarui Data
        </button>
      </div>

      <div class="overflow-x-auto">
        <div v-if="loading" class="py-12 text-center text-slate-400 text-xs">
          Memuat log presensi dari database server...
        </div>
        <table v-else class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4">Nama Karyawan</th>
              <th class="py-3 px-4">Shift</th>
              <th class="py-3 px-4 text-center">Tanggal</th>
              <th class="py-3 px-4 text-center">Jam Masuk</th>
              <th class="py-3 px-4 text-center">Jam Keluar</th>
              <th class="py-3 px-4 text-center">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="att in attendanceList" :key="att.id" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3.5 px-4 font-bold text-slate-900 flex items-center gap-2.5">
                <div class="w-7 h-7 rounded-full bg-blue-100 text-blue-700 flex items-center justify-center font-bold text-[10px]">
                  {{ (att.name || 'K').charAt(0) }}
                </div>
                <div>
                  <div>{{ att.name }}</div>
                  <div class="text-[10px] text-slate-400 font-normal font-mono">{{ att.nik || '-' }}</div>
                </div>
              </td>
              <td class="py-3.5 px-4">{{ att.shift }}</td>
              <td class="py-3.5 px-4 text-center font-mono text-slate-500">{{ att.date || '-' }}</td>
              <td class="py-3.5 px-4 text-center font-mono font-bold text-slate-800">{{ att.in }}</td>
              <td class="py-3.5 px-4 text-center font-mono text-slate-400">{{ att.out }}</td>
              <td class="py-3.5 px-4 text-center">
                <span
                  class="px-2.5 py-0.5 rounded-full text-[11px] font-bold"
                  :class="getAttendanceBadge(att.status)"
                >
                  {{ att.status }}
                </span>
              </td>
            </tr>
            <tr v-if="attendanceList.length === 0">
              <td colspan="6" class="py-8 text-center text-slate-400 text-xs">
                Belum ada catatan presensi aktif di database server.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Catat Presensi Manual -->
    <Teleport to="body">
      <Transition name="emp-modal-fade">
        <div v-if="showManualModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center bg-slate-950/45 backdrop-blur-sm p-4 overflow-y-auto">
          <div class="bg-white rounded-3xl max-w-lg w-full p-6 shadow-2xl border border-slate-100/90 my-8">
            <div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-5">
              <h3 class="text-lg font-black text-slate-900 tracking-tight">Catat Presensi Staf</h3>
              <button @click="showManualModal = false" class="text-slate-400 hover:text-slate-600 p-1.5 rounded-lg cursor-pointer">
                <X class="w-4 h-4" />
              </button>
            </div>

            <form @submit.prevent="submitManualAttendance" class="space-y-4 text-xs font-medium text-slate-700">
              <div>
                <label class="block text-slate-500 mb-1">Pilih Karyawan</label>
                <select v-model="manualForm.employee_id" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 font-semibold">
                  <option value="" disabled>Pilih Staf...</option>
                  <option v-for="emp in activeEmployees" :key="emp.id" :value="emp.id">
                    {{ emp.full_name }} ({{ emp.position_name || 'Staff' }})
                  </option>
                </select>
              </div>

              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-slate-500 mb-1">Tipe Aksi</label>
                  <select v-model="manualForm.action_type" class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900">
                    <option value="clock-in">Clock In (Masuk)</option>
                    <option value="clock-out">Clock Out (Pulang)</option>
                  </select>
                </div>
                <div>
                  <label class="block text-slate-500 mb-1">Catatan / Keterangan</label>
                  <input v-model="manualForm.notes" type="text" placeholder="Misal: Hadir tepat waktu" class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none" />
                </div>
              </div>

              <div class="flex justify-end gap-3 pt-4 border-t border-slate-100">
                <button type="button" @click="showManualModal = false" class="px-4 py-2 border border-slate-200 rounded-xl text-slate-600 hover:bg-slate-50 cursor-pointer">
                  Batal
                </button>
                <button type="submit" :disabled="actionLoading" class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold cursor-pointer disabled:opacity-60">
                  Simpan Presensi
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
import { CheckCircle2, Clock, LogOut, X } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const currentTime = ref(new Date().toLocaleTimeString('id-ID'))
const currentDateStr = ref(new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }))

const hasClockedIn = ref(false)
const hasClockedOut = ref(false)
const clockInTime = ref('')
const clockOutTime = ref('')
const loading = ref(false)
const actionLoading = ref(false)

const attendanceList = ref<any[]>([])
const activeEmployees = ref<any[]>([])
const showManualModal = ref(false)
const manualForm = ref({
  employee_id: '',
  action_type: 'clock-in',
  notes: ''
})

const fetchAttendanceData = async () => {
  loading.value = true
  try {
    const [attRes, statusRes] = await Promise.all([
      axios.get('/api/v1/hris/attendances'),
      axios.get('/api/v1/hris/attendances/today-status')
    ])

    if (attRes.data?.data) {
      attendanceList.value = attRes.data.data.map((a: any) => ({
        id: a.id,
        nik: a.nik,
        name: a.name,
        shift: a.shift || 'Pagi',
        in: a.clock_in || '-',
        out: a.clock_out || '-',
        date: a.date || '',
        status: a.status === 'late' ? 'Terlambat' : (a.status === 'present' ? 'Hadir' : (a.status === 'leave' ? 'Cuti' : 'Belum Masuk'))
      }))
    }

    if (statusRes.data?.data) {
      hasClockedIn.value = statusRes.data.data.has_clocked_in
      hasClockedOut.value = statusRes.data.data.has_clocked_out
      clockInTime.value = statusRes.data.data.clock_in_time || ''
      clockOutTime.value = statusRes.data.data.clock_out_time || ''
    }
  } catch (err: any) {
    console.error('Failed to load attendances:', err)
    notifyStore.error('Gagal memuat catatan presensi dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

const fetchEmployees = async () => {
  try {
    const res = await axios.get('/api/v1/hris/employees')
    activeEmployees.value = (res.data?.data || []).filter((e: any) => e.status?.toLowerCase() === 'active')
    if (activeEmployees.value.length > 0 && !manualForm.value.employee_id) {
      manualForm.value.employee_id = activeEmployees.value[0].id
    }
  } catch (err) {
    console.error('Failed to load employees:', err)
  }
}

onMounted(() => {
  fetchAttendanceData()
  fetchEmployees()
  setInterval(() => {
    currentTime.value = new Date().toLocaleTimeString('id-ID')
  }, 1000)
})

const todaySummary = computed(() => {
  const present = attendanceList.value.filter(a => a.status === 'Hadir').length
  const late = attendanceList.value.filter(a => a.status === 'Terlambat').length
  const leave = attendanceList.value.filter(a => a.status === 'Cuti').length
  return {
    status: present > 0 ? 'Aktif' : 'Menunggu',
    present,
    late,
    leave
  }
})

const getAttendanceBadge = (status: string) => {
  switch (status) {
    case 'Hadir':
      return 'bg-emerald-100 text-emerald-700'
    case 'Terlambat':
      return 'bg-amber-100 text-amber-800'
    case 'Cuti':
      return 'bg-blue-100 text-blue-700'
    case 'Belum Masuk':
      return 'bg-slate-100 text-slate-600'
    default:
      return 'bg-rose-100 text-rose-700'
  }
}

const handleClockIn = async () => {
  actionLoading.value = true
  try {
    const res = await axios.post('/api/v1/hris/attendances/clock-in', {
      notes: 'Presensi masuk via Web App'
    })
    hasClockedIn.value = true
    clockInTime.value = res.data?.data?.clock_in || new Date().toLocaleTimeString('id-ID')
    notifyStore.success(`Presensi masuk berhasil dicatat pada ${clockInTime.value}!`, 'Clock In Sukses')
    fetchAttendanceData()
  } catch (err: any) {
    console.error('Clock-in failed:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal melakukan Clock In', 'Gagal')
  } finally {
    actionLoading.value = false
  }
}

const handleClockOut = async () => {
  actionLoading.value = true
  try {
    const res = await axios.post('/api/v1/hris/attendances/clock-out', {
      notes: 'Presensi keluar via Web App'
    })
    hasClockedOut.value = true
    clockOutTime.value = res.data?.data?.clock_out || new Date().toLocaleTimeString('id-ID')
    notifyStore.success(`Presensi pulang berhasil dicatat pada ${clockOutTime.value}! Terima kasih atas kerja keras Anda.`, 'Clock Out Sukses')
    fetchAttendanceData()
  } catch (err: any) {
    console.error('Clock-out failed:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal melakukan Clock Out', 'Gagal')
  } finally {
    actionLoading.value = false
  }
}

const openManualAttendanceModal = () => {
  showManualModal.value = true
}

const submitManualAttendance = async () => {
  if (!manualForm.value.employee_id) {
    notifyStore.warning('Silakan pilih karyawan terlebih dahulu', 'Validasi')
    return
  }
  actionLoading.value = true
  try {
    const endpoint = manualForm.value.action_type === 'clock-in'
      ? '/api/v1/hris/attendances/clock-in'
      : '/api/v1/hris/attendances/clock-out'
    
    await axios.post(endpoint, {
      employee_id: manualForm.value.employee_id,
      notes: manualForm.value.notes || 'Catatan manual HR'
    })

    notifyStore.success('Presensi staf berhasil dicatat ke database!', 'Presensi Berhasil')
    showManualModal.value = false
    manualForm.value.notes = ''
    fetchAttendanceData()
  } catch (err: any) {
    console.error('Manual attendance failed:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal menyimpan presensi staf', 'Gagal')
  } finally {
    actionLoading.value = false
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
