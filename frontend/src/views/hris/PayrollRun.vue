<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { Check, ChevronLeft, ChevronRight, FileText, CheckCircle2 } from 'lucide-vue-next'
import PayslipModal from './PayslipModal.vue'
import { useNotificationStore } from '@/stores/notification.store'

interface EmployeePayroll {
  id: string
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
}

const notifyStore = useNotificationStore()
const currentStep = ref(4)
const loading = ref(false)

const steps = [
  { step: 1, label: 'Pilih Periode' },
  { step: 2, label: 'Review Data' },
  { step: 3, label: 'Kalkulasi' },
  { step: 4, label: 'Review Hasil' },
  { step: 5, label: 'Persetujuan' }
]

const employees = ref<EmployeePayroll[]>([])

const fetchPayrolls = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/hris/payrolls')
    if (res.data?.data) {
      employees.value = res.data.data.map((e: any, idx: number) => {
        const banks = ['BCA', 'Mandiri', 'BNI', 'BRI']
        const bankName = banks[idx % banks.length]
        return {
          id: e.id,
          name: e.name,
          nik: e.nik,
          position: e.position,
          department: e.department,
          bank: `${bankName} - 88123${idx + 100}`,
          basicSalary: Number(e.basicSalary || 0),
          allowance: Number(e.allowance || 0),
          overtime: Number(e.overtime || 0),
          gross: Number(e.gross || 0),
          bpjs: Number(e.bpjs || 0),
          pph21: Number(e.pph21 || 0),
          deduction: Number(e.deduction || 0),
          net: Number(e.net || 0)
        }
      })
    }
  } catch (err: any) {
    console.error('Failed to load payroll data:', err)
    notifyStore.error('Gagal memuat data penggajian dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPayrolls()
})

const summary = computed(() => {
  const gross = employees.value.reduce((acc, e) => acc + e.gross, 0)
  const deduction = employees.value.reduce((acc, e) => acc + e.deduction, 0)
  const net = employees.value.reduce((acc, e) => acc + e.net, 0)
  return { gross, deduction, net }
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
    period: 'September 2026',
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

const handleApprove = () => {
  currentStep.value = 5
  notifyStore.success('Payroll Batch September 2026 telah berhasil disetujui dan dijadwalkan untuk transfer bank!', 'Persetujuan Selesai')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight">Penggajian - September 2026</h1>
    </div>

    <!-- Stepper Navigation -->
    <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
      <div class="flex items-center justify-between max-w-4xl mx-auto">
        <div v-for="(s, idx) in steps" :key="s.step" class="flex items-center flex-1 last:flex-none">
          <!-- Step indicator -->
          <div class="flex items-center gap-3">
            <div 
              class="w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm transition-all"
              :class="[
                s.step < currentStep 
                  ? 'bg-emerald-600 text-white' 
                  : s.step === currentStep 
                    ? 'border-2 border-blue-600 text-blue-600 bg-white ring-4 ring-blue-50' 
                    : 'bg-slate-200 text-slate-500'
              ]"
            >
              <Check v-if="s.step < currentStep" class="w-4 h-4" />
              <span v-else>{{ s.step }}</span>
            </div>
            <div class="flex flex-col">
              <span class="text-xs text-slate-400 font-medium">Step {{ s.step }}</span>
              <span 
                class="text-sm font-semibold whitespace-nowrap"
                :class="s.step === currentStep ? 'text-blue-600' : 'text-slate-700'"
              >
                {{ s.label }}
              </span>
            </div>
          </div>
          
          <!-- Connecting Line -->
          <div 
            v-if="idx < steps.length - 1" 
            class="flex-1 mx-4 h-0.5"
            :class="s.step < currentStep ? 'bg-emerald-600' : 'bg-slate-200'"
          ></div>
        </div>
      </div>
    </div>

    <!-- Section Title -->
    <div>
      <h2 class="text-xl font-bold text-slate-900">Review Hasil</h2>
    </div>

    <!-- 3 Summary Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
        <p class="text-sm font-medium text-slate-500 mb-2">Gaji Bruto Total</p>
        <p class="text-3xl font-extrabold text-slate-900 tracking-tight">
          Rp {{ formatNum(summary.gross) }}
        </p>
      </div>

      <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
        <p class="text-sm font-medium text-slate-500 mb-2">Total Potongan</p>
        <p class="text-3xl font-extrabold text-slate-900 tracking-tight">
          Rp {{ formatNum(summary.deduction) }}
        </p>
      </div>

      <div class="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
        <p class="text-sm font-medium text-slate-500 mb-2">Gaji Netto Total</p>
        <p class="text-3xl font-extrabold text-slate-900 tracking-tight">
          Rp {{ formatNum(summary.net) }}
        </p>
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
              <th class="py-3.5 px-4">Nama Karyawan</th>
              <th class="py-3.5 px-4 text-right">Gaji Pokok (Rp)</th>
              <th class="py-3.5 px-4 text-right">Tunjangan (Rp)</th>
              <th class="py-3.5 px-4 text-right">Lembur (Rp)</th>
              <th class="py-3.5 px-4 text-right">Bruto (Rp)</th>
              <th class="py-3.5 px-4 text-right">BPJS (Rp)</th>
              <th class="py-3.5 px-4 text-right">PPh21 (Rp)</th>
              <th class="py-3.5 px-4 text-right">Total Potongan (Rp)</th>
              <th class="py-3.5 px-4 text-right">Netto (Rp)</th>
              <th class="py-3.5 px-4 text-center">Slip</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm">
            <tr v-if="loading">
              <td colspan="10" class="py-12 text-center text-slate-400 text-xs">
                Memuat data penggajian karyawan dari database...
              </td>
            </tr>
            <tr v-else-if="employees.length === 0">
              <td colspan="10" class="py-12 text-center text-slate-400 text-xs">
                Belum ada data payroll yang tercatat untuk periode ini.
              </td>
            </tr>
            <tr 
              v-else
              v-for="emp in employees" 
              :key="emp.id"
              class="hover:bg-slate-50/80 transition-colors cursor-pointer"
              @click="openPayslip(emp)"
            >
              <td class="py-3.5 px-4 font-medium text-slate-900">
                {{ emp.name }}
              </td>
              <td class="py-3.5 px-4 text-right text-slate-700">
                {{ formatNum(emp.basicSalary) }}
              </td>
              <td class="py-3.5 px-4 text-right text-slate-700">
                {{ formatNum(emp.allowance) }}
              </td>
              <td class="py-3.5 px-4 text-right text-slate-700">
                {{ formatNum(emp.overtime) }}
              </td>
              <td class="py-3.5 px-4 text-right font-semibold text-slate-900">
                {{ formatNum(emp.gross) }}
              </td>
              <td class="py-3.5 px-4 text-right text-slate-700">
                {{ formatNum(emp.bpjs) }}
              </td>
              <td class="py-3.5 px-4 text-right text-slate-700">
                {{ formatNum(emp.pph21) }}
              </td>
              <td class="py-3.5 px-4 text-right text-red-600 font-medium">
                {{ formatNum(emp.deduction) }}
              </td>
              <td class="py-3.5 px-4 text-right font-bold text-slate-900">
                {{ formatNum(emp.net) }}
              </td>
              <td class="py-3.5 px-4 text-center" @click.stop="openPayslip(emp)">
                <button class="p-1 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded-md transition-colors" title="Lihat Slip Gaji">
                  <FileText class="w-4 h-4" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="px-6 py-4 border-t border-slate-200 flex items-center justify-center gap-1">
        <button class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-400 hover:bg-slate-50">
          <ChevronLeft class="w-4 h-4" />
        </button>
        <button class="w-8 h-8 flex items-center justify-center rounded-lg border border-blue-600 bg-blue-50 text-blue-600 font-medium text-sm">
          1
        </button>
        <button class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-sm">
          2
        </button>
        <button class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-sm">
          3
        </button>
        <button class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50">
          <ChevronRight class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Actions Footer -->
    <div class="flex items-center justify-end gap-3 pt-2">
      <button 
        @click="currentStep = Math.max(1, currentStep - 1)"
        class="px-5 py-2.5 border border-slate-300 rounded-xl text-sm font-semibold text-slate-700 bg-white hover:bg-slate-50 shadow-xs transition-colors"
      >
        Kembali
      </button>
      <button 
        @click="handleApprove"
        class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 rounded-xl text-sm font-semibold text-white shadow-xs transition-colors flex items-center gap-2"
      >
        <CheckCircle2 class="w-4 h-4" />
        Setujui & Lanjutkan
      </button>
    </div>

    <!-- Payslip Modal -->
    <PayslipModal 
      :show="showPayslipModal" 
      :employee="selectedEmployee" 
      @close="showPayslipModal = false" 
    />
  </div>
</template>
