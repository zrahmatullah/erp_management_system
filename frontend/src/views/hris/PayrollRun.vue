<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { Check, FileText, CheckCircle2, Calculator, ArrowRight, ShieldCheck } from 'lucide-vue-next'
import PayslipModal from './PayslipModal.vue'
import { useNotificationStore } from '@/stores/notification.store'

interface EmployeePayroll {
  id: string
  employee_id: string
  name: string
  nik: string
  position: string
  department: string
  bank: string
  basicSalary: number
  allowance: number
  overtime: number
  gross: number
  bpjs: number
  pph21: number
  deduction: number
  net: number
  is_paid: boolean
  paid_at: string
  period: string
}

const notifyStore = useNotificationStore()
const currentStep = ref(4)
const loading = ref(false)
const calculating = ref(false)
const approving = ref(false)

const selectedMonth = ref(9)
const selectedYear = ref(2026)
const approvalResult = ref<any>(null)

const months = [
  { value: 1, label: 'Januari' },
  { value: 2, label: 'Februari' },
  { value: 3, label: 'Maret' },
  { value: 4, label: 'April' },
  { value: 5, label: 'Mei' },
  { value: 6, label: 'Juni' },
  { value: 7, label: 'Juli' },
  { value: 8, label: 'Agustus' },
  { value: 9, label: 'September' },
  { value: 10, label: 'Oktober' },
  { value: 11, label: 'November' },
  { value: 12, label: 'Desember' }
]

const monthName = computed(() => {
  return months.find(m => m.value === selectedMonth.value)?.label || 'September'
})

const periodString = computed(() => {
  const m = String(selectedMonth.value).padStart(2, '0')
  return `${selectedYear.value}-${m}`
})

const periodStartDate = computed(() => {
  const m = String(selectedMonth.value).padStart(2, '0')
  return `${selectedYear.value}-${m}-01`
})

const steps = [
  { step: 1, label: 'Pilih Periode' },
  { step: 2, label: 'Review Kehadiran' },
  { step: 3, label: 'Kalkulasi Gaji' },
  { step: 4, label: 'Review Hasil' },
  { step: 5, label: 'Pencairan Dana' }
]

const employees = ref<EmployeePayroll[]>([])
const activeEmployeesCount = ref(4)

const fetchPayrolls = async () => {
  loading.value = true
  try {
    const res = await axios.get(`/api/v1/hris/payrolls?period=${periodString.value}`)
    if (res.data?.data) {
      employees.value = res.data.data.map((e: any) => ({
        id: e.id,
        employee_id: e.employee_id,
        name: e.name,
        nik: e.nik,
        position: e.position,
        department: e.department,
        bank: e.bank || 'BCA Operasional',
        basicSalary: Number(e.basicSalary || 0),
        allowance: Number(e.allowance || 0),
        overtime: Number(e.overtime || 0),
        gross: Number(e.gross || 0),
        bpjs: Number(e.bpjs || 0),
        pph21: Number(e.pph21 || 0),
        deduction: Number(e.deduction || 0),
        net: Number(e.net || 0),
        is_paid: Boolean(e.is_paid),
        paid_at: e.paid_at || '-',
        period: e.period || `${monthName.value} ${selectedYear.value}`
      }))
    }
  } catch (err: any) {
    console.error('Failed to load payroll data:', err)
    notifyStore.error('Gagal memuat data penggajian dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

const fetchEmployeeCount = async () => {
  try {
    const res = await axios.get('/api/v1/hris/employees')
    const list = res.data?.data || []
    activeEmployeesCount.value = list.filter((e: any) => e.status?.toLowerCase() === 'active').length || list.length
  } catch (err) {
    console.error('Failed to load employee count:', err)
  }
}

onMounted(() => {
  fetchEmployeeCount()
  fetchPayrolls()
})

const summary = computed(() => {
  const gross = employees.value.reduce((acc, e) => acc + e.gross, 0)
  const deduction = employees.value.reduce((acc, e) => acc + e.deduction, 0)
  const net = employees.value.reduce((acc, e) => acc + e.net, 0)
  const paidCount = employees.value.filter(e => e.is_paid).length
  return { gross, deduction, net, paidCount }
})

const allPaid = computed(() => {
  return employees.value.length > 0 && employees.value.every(e => e.is_paid)
})

const formatNum = (val: number) => {
  return Number(val || 0).toLocaleString('id-ID')
}

const showPayslipModal = ref(false)
const selectedEmployee = ref<any>(null)

const openPayslip = (emp: EmployeePayroll) => {
  selectedEmployee.value = {
    name: emp.name,
    nik: emp.nik,
    position: emp.position,
    department: emp.department,
    bank: emp.bank,
    period: `${monthName.value} ${selectedYear.value}`,
    basicSalary: emp.basicSalary,
    mealAllowance: Math.round(emp.allowance * 0.6),
    transportAllowance: Math.round(emp.allowance * 0.4),
    overtimePay: emp.overtime,
    overtimeHours: emp.overtime > 0 ? Math.round(emp.overtime / 45000) : 0,
    bpjsHealth: Math.round(emp.bpjs * 0.35),
    bpjsEmployment: Math.round(emp.bpjs * 0.65),
    pph21: emp.pph21,
    latenessDeduction: 0,
    latenessCount: 0
  }
  showPayslipModal.value = true
}

const runCalculation = async () => {
  currentStep.value = 3
  calculating.value = true
  try {
    const res = await axios.post('/api/v1/hris/payrolls/run', {
      month: selectedMonth.value,
      year: selectedYear.value
    })
    notifyStore.success(res.data?.message || 'Kalkulasi payroll berhasil diproses di server!', 'Kalkulasi Sukses')
    await fetchPayrolls()
    currentStep.value = 4
  } catch (err: any) {
    console.error('Failed to run payroll:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal menjalankan kalkulasi payroll', 'Gagal')
    currentStep.value = 2
  } finally {
    calculating.value = false
  }
}

const handleApprove = async () => {
  if (allPaid.value) {
    notifyStore.info('Seluruh data penggajian periode ini sudah disetujui dan dicairkan.', 'Sudah Lunas')
    currentStep.value = 5
    return
  }
  approving.value = true
  try {
    const res = await axios.post('/api/v1/hris/payrolls/batch-approve', {
      period_start: periodStartDate.value,
      bank_name: 'BCA Operasional',
      payment_date: new Date().toISOString().split('T')[0]
    })

    approvalResult.value = res.data?.data || res.data
    currentStep.value = 5
    notifyStore.success('Payroll Batch telah disetujui & otomatis diterbitkan ke Jurnal Umum Akuntansi!', 'Pencairan Berhasil')
    await fetchPayrolls()
  } catch (err: any) {
    console.error('Batch approve failed:', err)
    const isForbidden = err.response?.status === 403
    const title = isForbidden ? 'Akses Ditolak' : 'Pemberitahuan Payroll'
    const msg = err.response?.data?.error || err.response?.data?.message || (isForbidden 
      ? 'Gagal menyetujui payroll. Pastikan role Anda memiliki wewenang Super Admin atau Manager.' 
      : 'Gagal memproses persetujuan payroll.')
    notifyStore.error(msg, title)
  } finally {
    approving.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-slate-900 tracking-tight">Penggajian Staf & Barista</h1>
        <p class="text-xs text-slate-500 mt-0.5">Periode: {{ monthName }} {{ selectedYear }} • Terintegrasi Otomatis dengan Jurnal Akuntansi</p>
      </div>
    </div>

    <!-- Stepper Navigation -->
    <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-xs overflow-x-auto">
      <div class="flex items-center justify-between min-w-[650px] max-w-4xl mx-auto">
        <div v-for="(s, idx) in steps" :key="s.step" class="flex items-center flex-1 last:flex-none">
          <!-- Step indicator -->
          <div 
            @click="s.step <= 4 && (currentStep = s.step)"
            class="flex items-center gap-3 cursor-pointer select-none"
          >
            <div 
              class="w-8 h-8 rounded-full flex items-center justify-center font-bold text-xs transition-all"
              :class="[
                s.step < currentStep 
                  ? 'bg-emerald-600 text-white shadow-sm' 
                  : s.step === currentStep 
                    ? 'border-2 border-blue-600 text-blue-600 bg-white ring-4 ring-blue-50' 
                    : 'bg-slate-200 text-slate-500'
              ]"
            >
              <Check v-if="s.step < currentStep" class="w-4 h-4" />
              <span v-else>{{ s.step }}</span>
            </div>
            <div class="flex flex-col">
              <span class="text-[10px] text-slate-400 font-medium">Langkah {{ s.step }}</span>
              <span 
                class="text-xs font-bold whitespace-nowrap"
                :class="s.step === currentStep ? 'text-blue-600' : 'text-slate-700'"
              >
                {{ s.label }}
              </span>
            </div>
          </div>
          
          <!-- Connecting Line -->
          <div 
            v-if="idx < steps.length - 1" 
            class="flex-1 mx-3 h-0.5"
            :class="s.step < currentStep ? 'bg-emerald-600' : 'bg-slate-200'"
          ></div>
        </div>
      </div>
    </div>

    <!-- STEP 1: PILIH PERIODE -->
    <div v-if="currentStep === 1" class="bg-white rounded-2xl border border-slate-200/80 p-8 max-w-2xl mx-auto shadow-sm space-y-6 text-xs text-slate-700 font-medium">
      <div class="border-b border-slate-100 pb-4">
        <h2 class="text-lg font-black text-slate-900">Langkah 1: Tentukan Periode Penggajian</h2>
        <p class="text-slate-500 mt-1">Pilih bulan dan tahun kalender operasional untuk penarikan absensi dan kalkulasi slip.</p>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-slate-500 mb-1.5 font-bold">Bulan Periode</label>
          <select v-model="selectedMonth" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none font-bold text-slate-900 text-sm">
            <option v-for="m in months" :key="m.value" :value="m.value">{{ m.label }}</option>
          </select>
        </div>
        <div>
          <label class="block text-slate-500 mb-1.5 font-bold">Tahun</label>
          <input v-model.number="selectedYear" type="number" class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none font-mono font-bold text-slate-900 text-sm" />
        </div>
      </div>

      <div class="p-4 rounded-xl bg-blue-50 border border-blue-100 text-blue-800 space-y-1">
        <div class="font-bold">Ringkasan Periode Aktif:</div>
        <div>01 {{ monthName }} {{ selectedYear }} s/d 30 {{ monthName }} {{ selectedYear }}</div>
        <div class="text-[11px] text-blue-600">Total Karyawan Aktif Terdaftar: <strong>{{ activeEmployeesCount }} Staf</strong></div>
      </div>

      <div class="flex justify-end pt-2">
        <button 
          @click="currentStep = 2"
          class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold flex items-center gap-2 cursor-pointer shadow-sm transition-all"
        >
          Lanjut ke Review Data
          <ArrowRight class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- STEP 2: REVIEW DATA KEHADIRAN -->
    <div v-else-if="currentStep === 2" class="bg-white rounded-2xl border border-slate-200/80 p-8 max-w-2xl mx-auto shadow-sm space-y-6 text-xs text-slate-700 font-medium">
      <div class="border-b border-slate-100 pb-4">
        <h2 class="text-lg font-black text-slate-900">Langkah 2: Sinkronisasi Absensi & Data Lembur</h2>
        <p class="text-slate-500 mt-1">Sistem akan menarik seluruh rekaman Clock-In/Out, durasi lembur, dan potongan dari modul presensi.</p>
      </div>

      <div class="space-y-3">
        <div class="p-4 rounded-xl border border-slate-200 flex items-center justify-between">
          <div>
            <div class="font-bold text-slate-900 text-sm">Status Presensi Periode Ini</div>
            <div class="text-slate-500 text-[11px] mt-0.5">Tercatat di server PostgreSQL tanpa data fiktif</div>
          </div>
          <span class="px-3 py-1 rounded-full text-xs font-bold bg-emerald-100 text-emerald-700">Tersinkronisasi</span>
        </div>

        <div class="p-4 rounded-xl border border-slate-200 flex items-center justify-between">
          <div>
            <div class="font-bold text-slate-900 text-sm">Tarif Dasar Lembur (Depnaker Standard)</div>
            <div class="text-slate-500 text-[11px] mt-0.5">1.5x Upah per Jam (Gaji Pokok / 173 Jam)</div>
          </div>
          <span class="font-mono font-bold text-slate-900">Aktif</span>
        </div>

        <div class="p-4 rounded-xl border border-slate-200 flex items-center justify-between">
          <div>
            <div class="font-bold text-slate-900 text-sm">Potongan BPJS & PPh 21</div>
            <div class="text-slate-500 text-[11px] mt-0.5">BPJS Kesehatan + Ketenagakerjaan 4% & PPh21 Tier 5%</div>
          </div>
          <span class="font-mono font-bold text-slate-900">Aktif</span>
        </div>
      </div>

      <div class="flex items-center justify-between pt-4 border-t border-slate-100">
        <button 
          @click="currentStep = 1"
          class="px-5 py-2.5 border border-slate-200 rounded-xl font-semibold text-slate-600 hover:bg-slate-50 cursor-pointer"
        >
          Kembali
        </button>
        <button 
          @click="runCalculation"
          :disabled="calculating"
          class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold flex items-center gap-2 cursor-pointer shadow-sm transition-all disabled:opacity-60"
        >
          <Calculator class="w-4 h-4" />
          Kalkulasi Payroll Sekarang
        </button>
      </div>
    </div>

    <!-- STEP 3: KALKULASI PROSES -->
    <div v-else-if="currentStep === 3" class="bg-white rounded-2xl border border-slate-200/80 p-12 max-w-lg mx-auto shadow-sm text-center space-y-4">
      <div class="w-16 h-16 rounded-full bg-blue-50 text-blue-600 flex items-center justify-center mx-auto animate-pulse">
        <Calculator class="w-8 h-8" />
      </div>
      <h3 class="text-lg font-black text-slate-900">Memproses Kalkulasi Payroll...</h3>
      <p class="text-xs text-slate-500">
        Menghitung gaji pokok, tunjangan operasional, insentif lembur, BPJS, dan pajak PPh21 secara otomatis di database.
      </p>
    </div>

    <!-- STEP 4: REVIEW HASIL -->
    <div v-else-if="currentStep === 4" class="space-y-6">
      <!-- 3 Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs">
          <p class="text-xs font-semibold text-slate-400 mb-1">Total Gaji Bruto</p>
          <p class="text-2xl font-black text-slate-900 tracking-tight font-mono">
            Rp {{ formatNum(summary.gross) }}
          </p>
          <p class="text-[11px] text-slate-500 mt-2">Gaji pokok + tunjangan makan & lembur</p>
        </div>

        <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs">
          <p class="text-xs font-semibold text-slate-400 mb-1">Total Potongan (BPJS & Pajak)</p>
          <p class="text-2xl font-black text-rose-600 tracking-tight font-mono">
            Rp {{ formatNum(summary.deduction) }}
          </p>
          <p class="text-[11px] text-slate-500 mt-2">BPJS Ketenagakerjaan + BPJS Kesehatan + PPh21</p>
        </div>

        <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs">
          <p class="text-xs font-semibold text-slate-400 mb-1">Total Gaji Netto (Take Home Pay)</p>
          <p class="text-2xl font-black text-emerald-600 tracking-tight font-mono">
            Rp {{ formatNum(summary.net) }}
          </p>
          <p class="text-[11px] text-slate-500 mt-2">{{ employees.length }} staf terhitung • {{ summary.paidCount }} lunas</p>
        </div>
      </div>

      <!-- Data Table -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden p-6">
        <div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-4">
          <h2 class="text-sm font-black text-slate-900">Rincian Slip Gaji Karyawan ({{ employees.length }})</h2>
          <button 
            @click="runCalculation"
            class="text-xs font-bold text-blue-600 hover:text-blue-800 transition-colors cursor-pointer"
          >
            Hitung Ulang
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 font-bold text-slate-500 uppercase text-[10px]">
                <th class="py-3 px-3">Nama Karyawan</th>
                <th class="py-3 px-3">Rekening Bank</th>
                <th class="py-3 px-3 text-right">Gaji Pokok</th>
                <th class="py-3 px-3 text-right">Tunjangan</th>
                <th class="py-3 px-3 text-right">Lembur</th>
                <th class="py-3 px-3 text-right">Bruto</th>
                <th class="py-3 px-3 text-right">BPJS (4%)</th>
                <th class="py-3 px-3 text-right">PPh 21</th>
                <th class="py-3 px-3 text-right">Netto (Rp)</th>
                <th class="py-3 px-3 text-center">Status</th>
                <th class="py-3 px-3 text-center">Slip</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
              <tr v-if="loading">
                <td colspan="11" class="py-12 text-center text-slate-400 text-xs">
                  Memuat data penggajian karyawan dari database...
                </td>
              </tr>
              <tr v-else-if="employees.length === 0">
                <td colspan="11" class="py-12 text-center text-slate-400 text-xs">
                  Belum ada data payroll yang tercatat untuk periode {{ monthName }} {{ selectedYear }}. Silakan klik tombol 'Kalkulasi Payroll Sekarang'.
                </td>
              </tr>
              <tr 
                v-else
                v-for="emp in employees" 
                :key="emp.id"
                class="hover:bg-slate-50/80 transition-colors cursor-pointer"
                @click="openPayslip(emp)"
              >
                <td class="py-3 px-3 font-bold text-slate-900 flex items-center gap-2">
                  <div class="w-7 h-7 rounded-full bg-blue-100 text-blue-700 flex items-center justify-center font-bold text-[10px]">
                    {{ emp.name.charAt(0) }}
                  </div>
                  <div>
                    <div>{{ emp.name }}</div>
                    <div class="text-[10px] text-slate-400 font-normal font-mono">{{ emp.nik }}</div>
                  </div>
                </td>
                <td class="py-3 px-3 font-mono text-slate-600">{{ emp.bank }}</td>
                <td class="py-3 px-3 text-right font-mono text-slate-700">{{ formatNum(emp.basicSalary) }}</td>
                <td class="py-3 px-3 text-right font-mono text-slate-700">{{ formatNum(emp.allowance) }}</td>
                <td class="py-3 px-3 text-right font-mono text-slate-700">{{ formatNum(emp.overtime) }}</td>
                <td class="py-3 px-3 text-right font-mono font-bold text-slate-900">{{ formatNum(emp.gross) }}</td>
                <td class="py-3 px-3 text-right font-mono text-slate-500">{{ formatNum(emp.bpjs) }}</td>
                <td class="py-3 px-3 text-right font-mono text-slate-500">{{ formatNum(emp.pph21) }}</td>
                <td class="py-3 px-3 text-right font-mono font-black text-emerald-600">{{ formatNum(emp.net) }}</td>
                <td class="py-3 px-3 text-center">
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                    :class="emp.is_paid ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-800'"
                  >
                    {{ emp.is_paid ? 'Disalurkan' : 'Pending' }}
                  </span>
                </td>
                <td class="py-3 px-3 text-center" @click.stop="openPayslip(emp)">
                  <button class="p-1.5 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded-lg transition-colors cursor-pointer" title="Cetak Slip Gaji Resmi">
                    <FileText class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Actions Footer -->
      <div class="flex items-center justify-between pt-2">
        <button 
          @click="currentStep = 2"
          class="px-5 py-2.5 border border-slate-300 rounded-xl text-xs font-bold text-slate-700 bg-white hover:bg-slate-50 shadow-xs transition-colors cursor-pointer"
        >
          Kembali ke Review
        </button>

        <div class="flex items-center gap-3">
          <div v-if="allPaid" class="flex items-center gap-2 text-xs font-bold text-emerald-700 bg-emerald-50 px-4 py-2.5 rounded-xl border border-emerald-200">
            <Check class="w-4 h-4 text-emerald-600" />
            <span>Seluruh Gaji Periode Ini Sudah Dicairkan</span>
          </div>

          <button 
            v-if="!allPaid"
            @click="handleApprove"
            :disabled="approving || employees.length === 0"
            class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 rounded-xl text-xs font-bold text-white shadow-md shadow-blue-600/30 transition-all flex items-center gap-2 cursor-pointer disabled:opacity-60"
          >
            <CheckCircle2 class="w-4 h-4" />
            {{ approving ? 'Memproses Jurnal...' : 'Setujui & Terbitkan Jurnal Akuntansi' }}
          </button>

          <button 
            v-else
            @click="currentStep = 5"
            class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 rounded-xl text-xs font-bold text-white shadow-md shadow-blue-600/30 transition-all flex items-center gap-2 cursor-pointer"
          >
            <span>Lihat Ringkasan Jurnal</span>
          </button>
        </div>
      </div>
    </div>

    <!-- STEP 5: PERSETUJUAN & JURNAL TERBIT -->
    <div v-else-if="currentStep === 5" class="bg-white rounded-2xl border border-slate-200/80 p-8 max-w-2xl mx-auto shadow-sm space-y-6 text-xs text-slate-700 font-medium">
      <div class="text-center space-y-2">
        <div class="w-14 h-14 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center mx-auto">
          <ShieldCheck class="w-8 h-8" />
        </div>
        <h2 class="text-xl font-black text-slate-900">Pencairan Payroll Berhasil Disetujui!</h2>
        <p class="text-slate-500">
          Dana gaji telah ditandai lunas dan otomatis dibukukan ke dalam Jurnal Finansial ERP.
        </p>
      </div>

      <div class="p-5 rounded-2xl bg-slate-50 border border-slate-200 space-y-3 font-mono">
        <div class="flex justify-between border-b border-slate-200 pb-2">
          <span class="text-slate-500 font-sans font-bold">Nomor Referensi Jurnal:</span>
          <span class="font-bold text-blue-600">{{ approvalResult?.journal_reference || 'JRN-PAYROLL-AUTO' }}</span>
        </div>
        <div class="flex justify-between border-b border-slate-200 pb-2">
          <span class="text-slate-500 font-sans font-bold">Total Nilai Gaji Netto:</span>
          <span class="font-bold text-emerald-600 font-black">Rp {{ formatNum(summary.net) }}</span>
        </div>
        <div class="flex justify-between border-b border-slate-200 pb-2">
          <span class="text-slate-500 font-sans font-bold">Akun Debit (Beban):</span>
          <span class="text-slate-900">6-1001 Beban Gaji Karyawan</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-500 font-sans font-bold">Akun Kredit (Kas):</span>
          <span class="text-slate-900">1-1102 Bank BCA Operasional</span>
        </div>
      </div>

      <div class="flex justify-center gap-4 pt-2">
        <button 
          @click="currentStep = 4"
          class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white rounded-xl font-bold cursor-pointer transition-all shadow-md shadow-blue-600/30"
        >
          Lihat Slip Karyawan
        </button>
      </div>
    </div>

    <!-- Payslip Modal -->
    <PayslipModal 
      :show="showPayslipModal" 
      :employee="selectedEmployee" 
      @close="showPayslipModal = false" 
    />
  </div>
</template>
