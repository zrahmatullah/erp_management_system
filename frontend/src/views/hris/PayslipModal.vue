<script setup lang="ts">
import { Printer, X, Coffee } from 'lucide-vue-next'

const props = defineProps<{
  show: boolean
  employee?: {
    name: string
    nik: string
    position: string
    department: string
    bank: string
    period: string
    basicSalary: number
    mealAllowance: number
    transportAllowance: number
    overtimePay: number
    overtimeHours: number
    bpjsHealth: number
    bpjsEmployment: number
    pph21: number
    latenessDeduction: number
    latenessCount: number
  } | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const formatRupiah = (val: number) => {
  return 'Rp ' + Number(val || 0).toLocaleString('id-ID')
}

const printSlip = () => {
  window.print()
}

// Defaults for preview if no employee passed
const emp = props.employee || {
  name: 'Sarah Andini',
  nik: 'EMP-003',
  position: 'Barista',
  department: 'Front of House',
  bank: 'BCA - 1234567890',
  period: 'September 2026',
  basicSalary: 5500000,
  mealAllowance: 600000,
  transportAllowance: 400000,
  overtimePay: 570000,
  overtimeHours: 12,
  bpjsHealth: 55000,
  bpjsEmployment: 110000,
  pph21: 145000,
  latenessDeduction: 50000,
  latenessCount: 2
}

const totalPendapatan = () => {
  return (emp.basicSalary || 0) + (emp.mealAllowance || 0) + (emp.transportAllowance || 0) + (emp.overtimePay || 0)
}

const totalPotongan = () => {
  return (emp.bpjsHealth || 0) + (emp.bpjsEmployment || 0) + (emp.pph21 || 0) + (emp.latenessDeduction || 0)
}

const takeHomePay = () => {
  return totalPendapatan() - totalPotongan()
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-xs p-4 overflow-y-auto">
    <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 w-full max-w-3xl overflow-hidden my-8 animate-in fade-in zoom-in-95 duration-200">
      
      <!-- Top Action Bar -->
      <div class="px-6 py-4 bg-slate-50 border-b border-slate-200 flex items-center justify-between">
        <span class="text-xs font-semibold uppercase tracking-wider text-slate-500">Official Payslip Document</span>
        <div class="flex items-center gap-3">
          <button @click="printSlip" class="inline-flex items-center gap-2 px-3 py-1.5 border border-slate-300 rounded-lg text-sm font-medium text-slate-700 bg-white hover:bg-slate-50 transition-colors shadow-xs">
            <Printer class="w-4 h-4 text-slate-600" />
            Print
          </button>
          <button @click="emit('close')" class="p-1.5 text-slate-400 hover:text-slate-700 rounded-lg hover:bg-slate-200/60 transition-colors">
            <X class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Payslip Body (Printable Paper Layout) -->
      <div class="p-8 md:p-12" id="printable-payslip">
        <!-- Header -->
        <div class="flex items-center justify-between pb-6 border-b border-slate-200">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 rounded-xl bg-amber-900/10 flex items-center justify-center text-amber-900 border border-amber-900/20">
              <Coffee class="w-6 h-6" />
            </div>
            <div>
              <h2 class="text-xl font-bold text-slate-900 tracking-tight">Cafe Harmony</h2>
              <p class="text-xs text-slate-500">Enterprise Food & Beverage ERP</p>
            </div>
          </div>
          <div class="text-right">
            <span class="text-sm font-medium text-slate-600">Periode: {{ emp.period || 'September 2026' }}</span>
          </div>
        </div>

        <!-- Title -->
        <div class="py-6 text-center">
          <h1 class="text-lg font-extrabold text-slate-900 tracking-wider">SLIP GAJI / PAYSLIP</h1>
        </div>

        <!-- Employee Info -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-y-1.5 text-sm mb-6 text-slate-700 font-medium">
          <div class="flex gap-2">
            <span class="w-28 text-slate-500">Nama</span>
            <span>: {{ emp.name }}</span>
          </div>
          <div class="flex gap-2">
            <span class="w-28 text-slate-500">Jabatan</span>
            <span>: {{ emp.position }}</span>
          </div>
          <div class="flex gap-2">
            <span class="w-28 text-slate-500">NIK</span>
            <span>: {{ emp.nik }}</span>
          </div>
          <div class="flex gap-2">
            <span class="w-28 text-slate-500">Departemen</span>
            <span>: {{ emp.department }}</span>
          </div>
          <div class="flex gap-2 sm:col-span-2">
            <span class="w-28 text-slate-500">Bank</span>
            <span>: {{ emp.bank }}</span>
          </div>
        </div>

        <!-- Tables: Pendapatan vs Potongan -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <!-- Pendapatan -->
          <div class="border border-slate-200 rounded-lg overflow-hidden">
            <div class="bg-slate-100 px-4 py-2 border-b border-slate-200 text-center font-bold text-xs uppercase tracking-wider text-slate-800">
              PENDAPATAN
            </div>
            <div class="divide-y divide-slate-100 text-sm">
              <div class="flex justify-between px-4 py-2.5">
                <span class="text-slate-600">Gaji Pokok</span>
                <span class="font-medium text-slate-900">{{ formatRupiah(emp.basicSalary) }}</span>
              </div>
              <div class="flex justify-between px-4 py-2.5">
                <span class="text-slate-600">Tunjangan Makan</span>
                <span class="font-medium text-slate-900">{{ formatRupiah(emp.mealAllowance) }}</span>
              </div>
              <div class="flex justify-between px-4 py-2.5">
                <span class="text-slate-600">Tunjangan Transport</span>
                <span class="font-medium text-slate-900">{{ formatRupiah(emp.transportAllowance) }}</span>
              </div>
              <div class="flex justify-between px-4 py-2.5">
                <span class="text-slate-600">Lembur ({{ emp.overtimeHours }} jam)</span>
                <span class="font-medium text-slate-900">{{ formatRupiah(emp.overtimePay) }}</span>
              </div>
              <div class="flex justify-between px-4 py-2.5 bg-slate-50/80 font-bold border-t border-slate-200">
                <span class="text-slate-900">Total Pendapatan</span>
                <span class="text-slate-900">{{ formatRupiah(totalPendapatan()) }}</span>
              </div>
            </div>
          </div>

          <!-- Potongan -->
          <div class="border border-slate-200 rounded-lg overflow-hidden">
            <div class="bg-slate-100 px-4 py-2 border-b border-slate-200 text-center font-bold text-xs uppercase tracking-wider text-slate-800">
              POTONGAN
            </div>
            <div class="divide-y divide-slate-100 text-sm">
              <div class="flex justify-between px-4 py-2.5">
                <span class="text-slate-600">BPJS Kesehatan</span>
                <span class="font-medium text-slate-900">{{ formatRupiah(emp.bpjsHealth) }}</span>
              </div>
              <div class="flex justify-between px-4 py-2.5">
                <span class="text-slate-600">BPJS Ketenagakerjaan</span>
                <span class="font-medium text-slate-900">{{ formatRupiah(emp.bpjsEmployment) }}</span>
              </div>
              <div class="flex justify-between px-4 py-2.5">
                <span class="text-slate-600">PPh 21</span>
                <span class="font-medium text-slate-900">{{ formatRupiah(emp.pph21) }}</span>
              </div>
              <div class="flex justify-between px-4 py-2.5">
                <span class="text-slate-600">Keterlambatan ({{ emp.latenessCount }}x)</span>
                <span class="font-medium text-slate-900">{{ formatRupiah(emp.latenessDeduction) }}</span>
              </div>
              <div class="flex justify-between px-4 py-2.5 bg-slate-50/80 font-bold border-t border-slate-200">
                <span class="text-slate-900">Total Potongan</span>
                <span class="text-slate-900">{{ formatRupiah(totalPotongan()) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Take Home Pay Box -->
        <div class="bg-slate-100 border border-slate-300 rounded-xl p-4 text-center">
          <div class="text-base sm:text-lg font-black text-slate-900 tracking-wide">
            GAJI BERSIH (TAKE HOME PAY): {{ formatRupiah(takeHomePay()) }}
          </div>
        </div>

        <!-- Signatures (Hidden on small screens, shown when printing) -->
        <div class="hidden print:grid grid-cols-2 gap-8 mt-12 pt-8 text-center text-xs text-slate-600">
          <div>
            <p>Diterima Oleh,</p>
            <div class="h-16"></div>
            <p class="font-semibold underline">{{ emp.name }}</p>
          </div>
          <div>
            <p>Finance / HR Manager,</p>
            <div class="h-16"></div>
            <p class="font-semibold underline">Agus Hermanto, SE</p>
          </div>
        </div>

      </div>

    </div>
  </div>
</template>

