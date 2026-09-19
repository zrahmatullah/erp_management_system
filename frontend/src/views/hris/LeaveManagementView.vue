<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Pengajuan Cuti & Izin</h2>
        <p class="text-xs text-slate-500 mt-1">Kelola permohonan libur kerja staf dan persetujuan supervisor berbasis database</p>
      </div>
      <button
        @click="openRequestLeave"
        class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
      >
        <span>+</span> Ajukan Cuti
      </button>
    </div>

    <!-- Tabs -->
    <div class="flex items-center gap-6 border-b border-slate-200 text-sm font-bold">
      <button
        v-for="t in ['Cuti Tim', 'Semua Cuti']"
        :key="t"
        @click="activeTab = t"
        class="pb-3 transition-colors relative cursor-pointer"
        :class="activeTab === t ? 'text-blue-600' : 'text-slate-500 hover:text-slate-800'"
      >
        {{ t }}
        <span v-if="activeTab === t" class="absolute bottom-0 inset-x-0 h-0.5 bg-blue-600 rounded-full"></span>
      </button>
    </div>

    <!-- Main Grid: Table & Saldo Cuti Tim -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Table Section (3 cols) -->
      <div class="lg:col-span-3 bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6 space-y-4">
        <!-- Filters -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <input
            v-model="search"
            type="text"
            placeholder="Cari Staf / NIK..."
            class="px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none"
          />
          <select v-model="filterType" class="px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none">
            <option value="">Semua Tipe Cuti</option>
            <option value="Cuti Tahunan">Cuti Tahunan</option>
            <option value="Cuti Sakit">Cuti Sakit</option>
            <option value="Cuti Alasan Penting">Cuti Alasan Penting</option>
          </select>
          <select v-model="filterStatus" class="px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none">
            <option value="">Semua Status</option>
            <option value="Menunggu">Menunggu</option>
            <option value="Disetujui">Disetujui</option>
            <option value="Ditolak">Ditolak</option>
          </select>
        </div>

        <div class="overflow-x-auto">
          <div v-if="loading" class="py-12 text-center text-slate-400 text-xs">
            Memuat permohonan cuti dari database...
          </div>
          <table v-else class="w-full text-left text-xs">
            <thead>
              <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
                <th class="py-3 px-3">Karyawan</th>
                <th class="py-3 px-3">Tipe Cuti</th>
                <th class="py-3 px-3">Dari</th>
                <th class="py-3 px-3">Sampai</th>
                <th class="py-3 px-3 text-center">Jumlah Hari</th>
                <th class="py-3 px-3">Alasan</th>
                <th class="py-3 px-3 text-center">Status</th>
                <th class="py-3 px-3 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
              <tr v-for="leave in filteredLeaves" :key="leave.id" class="hover:bg-slate-50/60 transition-colors">
                <td class="py-3 px-3 font-bold text-slate-900 flex items-center gap-2">
                  <div class="w-7 h-7 rounded-full bg-blue-100 text-blue-700 flex items-center justify-center text-[10px] font-bold">
                    {{ (leave.name || 'K').charAt(0) }}
                  </div>
                  <div>
                    <div>{{ leave.name }}</div>
                    <div class="text-[10px] text-slate-400 font-normal font-mono">{{ leave.nik || '-' }}</div>
                  </div>
                </td>
                <td class="py-3 px-3">{{ leave.type }}</td>
                <td class="py-3 px-3 text-slate-500 font-mono">{{ leave.from }}</td>
                <td class="py-3 px-3 text-slate-500 font-mono">{{ leave.to }}</td>
                <td class="py-3 px-3 text-center font-bold">{{ leave.days }} Hari</td>
                <td class="py-3 px-3 text-slate-600 truncate max-w-xs">{{ leave.reason }}</td>
                <td class="py-3 px-3 text-center">
                  <span
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                    :class="getStatusBadge(leave.status)"
                  >
                    {{ leave.status }}
                  </span>
                </td>
                <td class="py-3 px-3 text-right">
                  <div v-if="leave.status === 'Menunggu'" class="inline-flex items-center gap-1.5">
                    <button
                      @click="updateStatus(leave, 'approved')"
                      :disabled="actionLoading"
                      class="px-2.5 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-[11px] font-bold transition-colors cursor-pointer disabled:opacity-60"
                    >
                      Setujui
                    </button>
                    <button
                      @click="updateStatus(leave, 'rejected')"
                      :disabled="actionLoading"
                      class="px-2.5 py-1 border border-rose-300 text-rose-600 hover:bg-rose-50 rounded-lg text-[11px] font-bold transition-colors cursor-pointer disabled:opacity-60"
                    >
                      Tolak
                    </button>
                  </div>
                  <span v-else class="text-[11px] text-slate-400 italic">Selesai</span>
                </td>
              </tr>
              <tr v-if="filteredLeaves.length === 0">
                <td colspan="8" class="py-8 text-center text-slate-400 text-xs">
                  Tidak ada data permohonan cuti ditemukan.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Saldo Cuti Tim (Right 1 col) -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-5 space-y-4">
        <h3 class="font-bold text-slate-900 text-sm pb-2 border-b border-slate-100">Pemakaian Cuti Tim</h3>

        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-1">
          <div class="text-xs font-bold text-slate-800">Cuti Tahunan</div>
          <div class="text-[11px] text-slate-500">{{ countApprovedAnnual }} Hari Terpakai</div>
          <div class="w-full bg-slate-200 h-1.5 rounded-full mt-2 overflow-hidden">
            <div class="bg-blue-600 h-full rounded-full" :style="{ width: `${Math.min(100, countApprovedAnnual * 8.3)}%` }"></div>
          </div>
        </div>

        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-1">
          <div class="text-xs font-bold text-slate-800">Cuti Sakit</div>
          <div class="text-[11px] text-slate-500">{{ countApprovedSick }} Hari Terpakai</div>
          <div class="w-full bg-slate-200 h-1.5 rounded-full mt-2 overflow-hidden">
            <div class="bg-emerald-600 h-full rounded-full" :style="{ width: `${Math.min(100, countApprovedSick * 20)}%` }"></div>
          </div>
        </div>

        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-1">
          <div class="text-xs font-bold text-slate-800">Cuti Alasan Penting</div>
          <div class="text-[11px] text-slate-500">{{ countApprovedSpecial }} Hari Terpakai</div>
          <div class="w-full bg-slate-200 h-1.5 rounded-full mt-2 overflow-hidden">
            <div class="bg-purple-600 h-full rounded-full" :style="{ width: `${Math.min(100, countApprovedSpecial * 33.3)}%` }"></div>
          </div>
        </div>

        <div class="pt-2 text-[11px] text-slate-400 border-t border-slate-100">
          *Persetujuan cuti otomatis mengurangi sisa kuota cuti tahunan pegawai di sistem.
        </div>
      </div>
    </div>

    <!-- Modal Form Pengajuan Cuti -->
    <Teleport to="body">
      <Transition name="emp-modal-fade">
        <div v-if="showLeaveModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center bg-slate-950/45 backdrop-blur-sm p-4 overflow-y-auto">
          <div class="bg-white rounded-3xl max-w-lg w-full p-6 shadow-2xl border border-slate-100/90 my-8">
            <div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-5">
              <h3 class="text-lg font-black text-slate-900 tracking-tight">Formulir Pengajuan Cuti</h3>
              <button @click="showLeaveModal = false" class="text-slate-400 hover:text-slate-600 p-1.5 rounded-lg cursor-pointer">
                <X class="w-4 h-4" />
              </button>
            </div>

            <form @submit.prevent="submitLeaveRequest" class="space-y-4 text-xs font-medium text-slate-700">
              <div>
                <label class="block text-slate-500 mb-1">Pilih Karyawan</label>
                <select v-model="leaveForm.employee_id" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900 font-semibold">
                  <option value="" disabled>Pilih Staf...</option>
                  <option v-for="emp in activeEmployees" :key="emp.id" :value="emp.id">
                    {{ emp.full_name }} (Sisa Cuti: {{ emp.remaining_leave ?? 12 }} hari)
                  </option>
                </select>
              </div>

              <div>
                <label class="block text-slate-500 mb-1">Tipe Cuti</label>
                <select v-model="leaveForm.leave_type" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900">
                  <option value="Cuti Tahunan">Cuti Tahunan</option>
                  <option value="Cuti Sakit">Cuti Sakit</option>
                  <option value="Cuti Alasan Penting">Cuti Alasan Penting</option>
                </select>
              </div>

              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-slate-500 mb-1">Tanggal Mulai</label>
                  <input v-model="leaveForm.start_date" type="date" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none font-mono" />
                </div>
                <div>
                  <label class="block text-slate-500 mb-1">Tanggal Selesai</label>
                  <input v-model="leaveForm.end_date" type="date" required class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none font-mono" />
                </div>
              </div>

              <div>
                <label class="block text-slate-500 mb-1">Alasan Cuti / Keterangan</label>
                <textarea v-model="leaveForm.reason" rows="3" required placeholder="Contoh: Keperluan keluarga mendesak di luar kota" class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl outline-none text-slate-900"></textarea>
              </div>

              <div class="flex justify-end gap-3 pt-4 border-t border-slate-100">
                <button type="button" @click="showLeaveModal = false" class="px-4 py-2 border border-slate-200 rounded-xl text-slate-600 hover:bg-slate-50 cursor-pointer">
                  Batal
                </button>
                <button type="submit" :disabled="actionLoading" class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold cursor-pointer disabled:opacity-60">
                  Kirim Permohonan
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
const activeTab = ref('Cuti Tim')
const search = ref('')
const filterType = ref('')
const filterStatus = ref('')
const loading = ref(false)
const actionLoading = ref(false)

const leaves = ref<any[]>([])
const activeEmployees = ref<any[]>([])
const showLeaveModal = ref(false)

const today = new Date().toISOString().split('T')[0]
const leaveForm = ref({
  employee_id: '',
  leave_type: 'Cuti Tahunan',
  start_date: today,
  end_date: today,
  reason: ''
})

const fetchLeaves = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/hris/leaves')
    if (res.data?.data) {
      leaves.value = res.data.data.map((l: any) => ({
        id: l.id,
        nik: l.nik,
        name: l.name,
        type: l.leave_type,
        from: l.start_date,
        to: l.end_date,
        days: l.total_days,
        reason: l.reason,
        status: l.status === 'approved' ? 'Disetujui' : (l.status === 'rejected' ? 'Ditolak' : 'Menunggu')
      }))
    }
  } catch (err: any) {
    console.error('Failed to load leaves:', err)
    notifyStore.error('Gagal mengambil daftar cuti dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

const fetchEmployees = async () => {
  try {
    const res = await axios.get('/api/v1/hris/employees')
    activeEmployees.value = (res.data?.data || []).filter((e: any) => e.status?.toLowerCase() === 'active')
    if (activeEmployees.value.length > 0 && !leaveForm.value.employee_id) {
      leaveForm.value.employee_id = activeEmployees.value[0].id
    }
  } catch (err) {
    console.error('Failed to load employees:', err)
  }
}

onMounted(() => {
  fetchLeaves()
  fetchEmployees()
})

const filteredLeaves = computed(() => {
  return leaves.value.filter(l => {
    const matchSearch = !search.value || l.name.toLowerCase().includes(search.value.toLowerCase()) || (l.nik && l.nik.toLowerCase().includes(search.value.toLowerCase()))
    const matchType = !filterType.value || l.type === filterType.value
    const matchStatus = !filterStatus.value || l.status === filterStatus.value
    return matchSearch && matchType && matchStatus
  })
})

const countApprovedAnnual = computed(() => leaves.value.filter(l => l.type === 'Cuti Tahunan' && l.status === 'Disetujui').reduce((acc, l) => acc + l.days, 0))
const countApprovedSick = computed(() => leaves.value.filter(l => l.type === 'Cuti Sakit' && l.status === 'Disetujui').reduce((acc, l) => acc + l.days, 0))
const countApprovedSpecial = computed(() => leaves.value.filter(l => l.type === 'Cuti Alasan Penting' && l.status === 'Disetujui').reduce((acc, l) => acc + l.days, 0))

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'Menunggu':
      return 'bg-amber-100 text-amber-800'
    case 'Disetujui':
      return 'bg-emerald-100 text-emerald-700'
    case 'Ditolak':
      return 'bg-rose-100 text-rose-700'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

const updateStatus = async (leave: any, targetStatus: 'approved' | 'rejected') => {
  actionLoading.value = true
  try {
    await axios.put(`/api/v1/hris/leaves/${leave.id}/status`, {
      status: targetStatus
    })

    const label = targetStatus === 'approved' ? 'Disetujui' : 'Ditolak'
    leave.status = label
    notifyStore.success(`Permohonan cuti ${leave.name} berhasil diubah statusnya menjadi '${label}'.`, 'Status Cuti Diperbarui')
    fetchLeaves()
    fetchEmployees() // Refresh employee quota
  } catch (err: any) {
    console.error('Update leave status failed:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal mengubah status permohonan cuti', 'Akses Ditolak')
  } finally {
    actionLoading.value = false
  }
}

const openRequestLeave = () => {
  showLeaveModal.value = true
}

const submitLeaveRequest = async () => {
  if (!leaveForm.value.employee_id) {
    notifyStore.warning('Silakan pilih karyawan yang mengajukan cuti', 'Validasi')
    return
  }
  if (!leaveForm.value.start_date || !leaveForm.value.end_date) {
    notifyStore.warning('Tanggal cuti harus diisi lengkap', 'Validasi')
    return
  }

  actionLoading.value = true
  try {
    await axios.post('/api/v1/hris/leaves', {
      employee_id: leaveForm.value.employee_id,
      leave_type: leaveForm.value.leave_type,
      start_date: leaveForm.value.start_date,
      end_date: leaveForm.value.end_date,
      reason: leaveForm.value.reason
    })

    notifyStore.success('Permohonan cuti berhasil diajukan dan masuk ke antrean persetujuan supervisor!', 'Pengajuan Berhasil')
    showLeaveModal.value = false
    leaveForm.value.reason = ''
    fetchLeaves()
  } catch (err: any) {
    console.error('Submit leave failed:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal mengajukan cuti', 'Gagal')
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
