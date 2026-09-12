<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Trash2, Calendar, CheckCircle2, ArrowLeft } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const router = useRouter()
const notifyStore = useNotificationStore()

const entryDate = ref('2026-09-04')
const referenceNo = ref('JV-2026-0234')
const description = ref('Pembayaran sewa outlet bulan September')

interface JournalLine {
  id: string
  accountId: string
  accountName: string
  memo: string
  debit: number | null
  credit: number | null
}

import axios from 'axios'
import { onMounted } from 'vue'

const accounts = ref<Array<{ id: string, label: string }>>([])

const fetchAccounts = async () => {
  try {
    const res = await axios.get('/api/v1/master/accounts')
    const list = Array.isArray(res.data) ? res.data : (res.data?.data || [])
    accounts.value = list.map((a: any) => ({
      id: a.id,
      code: a.code,
      label: `${a.code} - ${a.name}`
    }))
    if (accounts.value.length > 0 && rows.value.length > 0 && !rows.value[0].accountId) {
      rows.value[0].accountId = accounts.value[0].id
      rows.value[0].accountName = accounts.value[0].label
      if (rows.value.length > 1 && accounts.value.length > 1) {
        rows.value[1].accountId = accounts.value[1].id
        rows.value[1].accountName = accounts.value[1].label
      }
    }
  } catch (err) {
    console.error('Failed to load accounts for journal:', err)
  }
}

onMounted(() => {
  fetchAccounts()
})

const rows = ref<JournalLine[]>([
  { id: '1', accountId: '', accountName: '', memo: '', debit: null, credit: null },
  { id: '2', accountId: '', accountName: '', memo: '', debit: null, credit: null }
])

const addRow = () => {
  const defaultAcc = accounts.value[0]?.id || ''
  const defaultLabel = accounts.value[0]?.label || ''
  rows.value.push({
    id: Date.now().toString(),
    accountId: defaultAcc,
    accountName: defaultLabel,
    memo: '',
    debit: null,
    credit: null
  })
}

const removeRow = (idx: number) => {
  if (rows.value.length > 2) {
    rows.value.splice(idx, 1)
  }
}

const totalDebit = computed(() => {
  return rows.value.reduce((acc, row) => acc + (Number(row.debit) || 0), 0)
})

const totalCredit = computed(() => {
  return rows.value.reduce((acc, row) => acc + (Number(row.credit) || 0), 0)
})

const diff = computed(() => {
  return Math.abs(totalDebit.value - totalCredit.value)
})

const isBalanced = computed(() => {
  return totalDebit.value > 0 && totalDebit.value === totalCredit.value
})

const formatNum = (val: number) => {
  return Number(val || 0).toLocaleString('id-ID')
}

const submitting = ref(false)

const handlePosting = async () => {
  if (!isBalanced.value) {
    notifyStore.warning('Jurnal belum seimbang! Total Debit harus sama dengan Total Kredit.', 'Validasi Keseimbangan')
    return
  }

  submitting.value = true
  try {
    const payload = {
      reference_no: referenceNo.value,
      entry_date: entryDate.value,
      description: description.value,
      lines: rows.value.map(r => ({
        account_id: r.accountId,
        description: r.memo || description.value,
        debit: Number(r.debit || 0),
        credit: Number(r.credit || 0)
      }))
    }

    const res = await axios.post('/api/v1/finance/journals', payload)
    notifyStore.success(res.data?.message || `Jurnal ${referenceNo.value} berhasil di-posting ke Buku Besar!`, 'Jurnal Diposting')
    router.push('/finance')
  } catch (err: any) {
    console.error('Failed to post journal:', err)
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan jurnal ke server', 'Error')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="max-w-6xl mx-auto space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button 
          @click="router.push('/finance')"
          class="p-2 border border-slate-200 rounded-xl hover:bg-slate-100 text-slate-600 transition-colors"
        >
          <ArrowLeft class="w-5 h-5" />
        </button>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight">Buat Jurnal Baru</h1>
      </div>
    </div>

    <!-- Top Card: Header Inputs -->
    <div class="bg-white rounded-2xl border border-slate-200 p-6 shadow-xs space-y-5">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <!-- Tanggal -->
        <div>
          <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-2">Tanggal</label>
          <div class="relative">
            <input 
              v-model="entryDate" 
              type="date" 
              class="w-full pl-10 pr-4 py-2.5 bg-slate-50/50 border border-slate-300 rounded-xl text-sm text-slate-800 font-medium focus:outline-hidden focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
            />
            <Calendar class="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
          </div>
        </div>

        <!-- No Referensi -->
        <div>
          <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-2">No. Referensi</label>
          <input 
            v-model="referenceNo" 
            type="text" 
            placeholder="JV-2026-XXXX" 
            class="w-full px-4 py-2.5 bg-slate-50/50 border border-slate-300 rounded-xl text-sm font-semibold text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
          />
        </div>
      </div>

      <!-- Deskripsi -->
      <div>
        <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-2">Deskripsi</label>
        <input 
          v-model="description" 
          type="text" 
          placeholder="Tuliskan keterangan ringkas transaksi..." 
          class="w-full px-4 py-2.5 bg-slate-50/50 border border-slate-300 rounded-xl text-sm text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
        />
      </div>
    </div>

    <!-- Journal Lines Table -->
    <div class="bg-white rounded-2xl border border-slate-200 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
              <th class="py-3 px-4 w-72">Akun</th>
              <th class="py-3 px-4">Keterangan</th>
              <th class="py-3 px-4 text-right w-44">Debit (Rp)</th>
              <th class="py-3 px-4 text-right w-44">Kredit (Rp)</th>
              <th class="py-3 px-2 w-12"></th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm">
            <tr v-for="(row, idx) in rows" :key="row.id" class="hover:bg-slate-50/50 transition-colors">
              <td class="p-3">
                <select 
                  v-model="row.accountId" 
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-sm text-slate-800 font-medium focus:outline-hidden focus:ring-2 focus:ring-blue-500"
                >
                  <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                    {{ acc.label }}
                  </option>
                </select>
              </td>
              <td class="p-3">
                <input 
                  v-model="row.memo" 
                  type="text" 
                  placeholder="Keterangan baris (opsional)" 
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-sm text-slate-700 focus:outline-hidden focus:ring-2 focus:ring-blue-500"
                />
              </td>
              <td class="p-3">
                <input 
                  v-model.number="row.debit" 
                  type="number" 
                  placeholder="0" 
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-sm font-semibold text-right text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-blue-500"
                />
              </td>
              <td class="p-3">
                <input 
                  v-model.number="row.credit" 
                  type="number" 
                  placeholder="0" 
                  class="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-sm font-semibold text-right text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-blue-500"
                />
              </td>
              <td class="p-3 text-center">
                <button 
                  v-if="rows.length > 2"
                  @click="removeRow(idx)" 
                  class="text-slate-400 hover:text-rose-600 transition-colors p-1"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Add Row Button -->
      <div class="p-4 border-t border-slate-100 bg-slate-50/50">
        <button 
          @click="addRow" 
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 border border-slate-300 rounded-lg text-sm font-semibold text-slate-700 bg-white hover:bg-slate-50 shadow-xs transition-colors"
        >
          <Plus class="w-4 h-4" />
          Add Row
        </button>
      </div>
    </div>

    <!-- Bottom Balance Summary & Action Bar -->
    <div class="bg-white rounded-2xl border border-slate-200 p-6 shadow-xs flex flex-col md:flex-row items-center justify-between gap-4">
      <div class="flex items-center gap-4 text-sm font-medium">
        <div class="text-slate-700">
          Total Debit: <span class="font-bold text-slate-900">Rp {{ formatNum(totalDebit) }}</span>
        </div>
        <span class="text-slate-300">|</span>
        <div class="text-slate-700">
          Total Kredit: <span class="font-bold text-slate-900">Rp {{ formatNum(totalCredit) }}</span>
        </div>
        <span class="text-slate-300">|</span>
        <div class="flex items-center gap-2">
          <span class="text-slate-700">Selisih: </span>
          <span class="font-bold text-slate-900">Rp {{ formatNum(diff) }}</span>
          <span 
            class="w-3 h-3 rounded-full inline-block"
            :class="isBalanced ? 'bg-emerald-500' : 'bg-rose-500 animate-pulse'"
            :title="isBalanced ? 'Seimbang' : 'Belum Seimbang'"
          ></span>
        </div>
      </div>

      <div class="flex items-center gap-3 w-full md:w-auto justify-end">
        <button 
          @click="router.push('/finance')"
          class="px-5 py-2.5 border border-slate-300 rounded-xl text-sm font-semibold text-slate-700 hover:bg-slate-50 transition-colors shadow-xs"
        >
          Simpan Draft
        </button>
        <button 
          @click="handlePosting"
          class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-sm font-semibold shadow-xs transition-colors"
        >
          Posting
        </button>
      </div>
    </div>
  </div>
</template>
