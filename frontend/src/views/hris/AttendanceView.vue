<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Presensi & Absensi</h2>
        <p class="text-xs text-slate-500 mt-1">Pencatatan kehadiran staf barista, kasir dan kitchen real-time</p>
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
            @click="clockIn"
            :disabled="hasClockedIn"
            class="p-4 rounded-2xl bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-left shadow-lg shadow-emerald-600/30 transition-all cursor-pointer flex items-center justify-between disabled:opacity-80"
          >
            <div>
              <div class="text-base font-extrabold flex items-center gap-1.5">
                <CheckCircle2 class="w-4 h-4" /> Clock In
              </div>
              <div class="text-xs text-emerald-100 font-medium mt-0.5">
                Masuk: {{ clockInTime || '07:00' }}
              </div>
            </div>
            <Clock class="w-8 h-8 text-emerald-200" />
          </button>

          <button
            @click="clockOut"
            class="p-4 rounded-2xl bg-amber-600 hover:bg-amber-700 text-white font-bold text-left shadow-lg shadow-amber-600/30 transition-all cursor-pointer flex items-center justify-between"
          >
            <div>
              <div class="text-base font-extrabold flex items-center gap-1.5">
                <LogOut class="w-4 h-4" /> Clock Out
              </div>
              <div class="text-xs text-amber-100 font-medium mt-0.5">
                {{ clockOutTime || 'Belum ditekan' }}
              </div>
            </div>
            <LogOut class="w-8 h-8 text-amber-200" />
          </button>
        </div>
      </div>

      <!-- Today's Status Card (1 col) -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm flex flex-col justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-400 mb-1">Today's Status</div>
          <div class="text-3xl font-black text-emerald-600">{{ todaySummary.status }}</div>
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
            <div class="text-slate-400">Cuti / Sakit</div>
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
      <div class="flex items-center gap-6 border-b border-slate-200 pb-3 mb-5 text-xs font-bold">
        <button
          @click="activeSubTab = 'today'"
          class="cursor-pointer"
          :class="activeSubTab === 'today' ? 'text-blue-600 border-b-2 border-blue-600 pb-3 -mb-3' : 'text-slate-500'"
        >
          Daftar Presensi ({{ attendanceList.length }})
        </button>
      </div>

      <div class="overflow-x-auto">
        <div v-if="loading" class="py-12 text-center text-slate-400 text-xs">
          Memuat log presensi dari database...
        </div>
        <table v-else class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4">Nama Karyawan</th>
              <th class="py-3 px-4">Shift</th>
              <th class="py-3 px-4 text-center">Tanggal</th>
              <th class="py-3 px-4 text-center">Masuk</th>
              <th class="py-3 px-4 text-center">Keluar</th>
              <th class="py-3 px-4 text-center">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="att in attendanceList" :key="att.id || att.name" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3.5 px-4 font-bold text-slate-900 flex items-center gap-2.5">
                <div class="w-7 h-7 rounded-full bg-slate-200 text-slate-600 flex items-center justify-center font-bold text-[10px]">
                  {{ (att.name || 'U').charAt(0) }}
                </div>
                <div>
                  <div>{{ att.name }}</div>
                  <div class="text-[10px] text-slate-400 font-normal font-mono">{{ att.nik }}</div>
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { CheckCircle2, Clock, LogOut } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const currentTime = ref(new Date().toLocaleTimeString('id-ID'))
const currentDateStr = ref(new Date().toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }))
const activeSubTab = ref('today')
const hasClockedIn = ref(true)
const clockInTime = ref('07:12')
const clockOutTime = ref('')
const loading = ref(false)

const attendanceList = ref<any[]>([])

const fetchAttendances = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/hris/attendances')
    if (res.data?.data) {
      attendanceList.value = res.data.data.map((a: any) => ({
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
  } catch (err: any) {
    console.error('Failed to load attendances:', err)
    notifyStore.error('Gagal memuat catatan presensi dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAttendances()
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

const clockIn = () => {
  hasClockedIn.value = true
  clockInTime.value = new Date().toLocaleTimeString('id-ID')
  notifyStore.success(`Presensi masuk berhasil dicatat pada ${clockInTime.value}!`, 'Clock In Sukses')
}

const clockOut = () => {
  clockOutTime.value = new Date().toLocaleTimeString('id-ID')
  notifyStore.success(`Presensi pulang berhasil dicatat pada ${clockOutTime.value}! Terima kasih atas kerja keras Anda.`, 'Clock Out Sukses')
}
</script>
