<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Jadwal Shift Kerja</h2>
        <!-- Date Navigator (Matches Mockup) -->
        <div class="inline-flex items-center gap-2 bg-white px-3 py-1.5 rounded-xl border border-slate-200/80 shadow-xs text-xs font-bold text-slate-700">
          <button class="hover:text-blue-600 cursor-pointer">&lt;</button>
          <span>01 - 07 Sep 2026</span>
          <button class="hover:text-blue-600 cursor-pointer">&gt;</button>
        </div>
      </div>

      <button
        @click="manageShifts"
        class="px-4 py-2 bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold rounded-xl shadow-xs transition-colors cursor-pointer"
      >
        Kelola Shift
      </button>
    </div>

    <!-- Main Grid: Weekly Schedule + Shift Definitions -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Weekly Schedule Table (Left 3 cols) -->
      <div class="lg:col-span-3 bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6 overflow-x-auto">
        <table class="w-full text-left text-xs min-w-[650px]">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-3 w-40">Karyawan</th>
              <th v-for="day in ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu']" :key="day" class="py-3 px-2 text-center">
                {{ day }}
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="staff in staffRoster" :key="staff.name" class="hover:bg-slate-50/60 transition-colors">
              <td class="py-3 px-3 font-bold text-slate-900 flex items-center gap-2.5">
                <div class="w-7 h-7 rounded-full bg-slate-200 text-slate-700 flex items-center justify-center font-bold text-[10px]">
                  {{ staff.name.charAt(0) }}
                </div>
                <span class="truncate">{{ staff.name }}</span>
              </td>
              <td v-for="(shift, idx) in staff.shifts" :key="idx" class="py-3 px-1.5 text-center">
                <span
                  class="inline-block px-2 py-1 rounded-lg text-[10px] font-extrabold shadow-2xs"
                  :class="getShiftStyle(shift)"
                >
                  {{ shift }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Definisi Shift (Right 1 col - Matches Mockup) -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-5 space-y-4">
        <div class="flex items-center justify-between pb-2 border-b border-slate-100">
          <h3 class="font-bold text-slate-900 text-sm">Definisi Shift</h3>
          <span class="text-[10px] text-blue-600 font-bold cursor-pointer">+ Baru</span>
        </div>

        <div class="space-y-3">
          <!-- Pagi -->
          <div class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between hover:border-blue-300 transition-colors">
            <div>
              <div class="font-bold text-slate-900 text-xs flex items-center gap-1.5">
                <span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span>
                Pagi
              </div>
              <div class="text-[11px] text-slate-500 mt-0.5 font-mono">07:00 - 15:00</div>
            </div>
            <button class="p-1.5 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer" title="Edit Shift">
              <Edit class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Siang -->
          <div class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between hover:border-emerald-300 transition-colors">
            <div>
              <div class="font-bold text-slate-900 text-xs flex items-center gap-1.5">
                <span class="w-2.5 h-2.5 rounded-full bg-emerald-600"></span>
                Siang
              </div>
              <div class="text-[11px] text-slate-500 mt-0.5 font-mono">15:00 - 23:00</div>
            </div>
            <button class="p-1.5 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer" title="Edit Shift">
              <Edit class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Malam -->
          <div class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between hover:border-purple-300 transition-colors">
            <div>
              <div class="font-bold text-slate-900 text-xs flex items-center gap-1.5">
                <span class="w-2.5 h-2.5 rounded-full bg-purple-600"></span>
                Malam
              </div>
              <div class="text-[11px] text-slate-500 mt-0.5 font-mono">23:00 - 07:00</div>
            </div>
            <button class="p-1.5 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer" title="Edit Shift">
              <Edit class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { Edit } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const loading = ref(false)
const staffRoster = ref<any[]>([])

const fetchRoster = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/hris/schedules')
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

const manageShifts = () => {
  notifyStore.info('Modal kelola master template shift staf dibuka.', 'Master Shift')
}
</script>

