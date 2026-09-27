<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { 
  Receipt, 
  CreditCard, 
  TrendingUp, 
  Wallet, 
  Plus,
  Search,
  Filter,
  Eye,
  CheckCircle2,
  AlertTriangle,
  AlertCircle,
  FileText,
  DollarSign,
  ArrowUpRight,
  ArrowDownLeft,
  Calendar,
  Building,
  Check,
  X,
  RefreshCw,
  ExternalLink,
  Scale,
  ShieldCheck,
  Percent
} from 'lucide-vue-next'

const router = useRouter()

// -----------------------------------------------------------------------------
// TABS CONFIGURATION
// -----------------------------------------------------------------------------
const activeTab = ref('overview')

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'journal', label: 'Journal Entries' },
  { id: 'coa', label: 'Chart of Accounts' },
  { id: 'expenses', label: 'Expenses' },
  { id: 'tax', label: 'Tax' },
  { id: 'budget', label: 'Budget' },
  { id: 'reconciliation', label: 'Bank Reconciliation' }
]

// -----------------------------------------------------------------------------
// STATE: COMMON & OVERVIEW
// -----------------------------------------------------------------------------
const loading = ref(false)
const notification = ref<{ message: string; type: 'success' | 'error' } | null>(null)

const showNotify = (message: string, type: 'success' | 'error' = 'success') => {
  notification.value = { message, type }
  setTimeout(() => {
    notification.value = null
  }, 4000)
}

const formatNum = (num: number) => {
  return Number(num || 0).toLocaleString('id-ID')
}

const overviewData = ref({
  total_revenue: 0,
  total_expenses: 0,
  net_profit: 0,
  cash_on_hand: 0
})

const recentJournals = ref<any[]>([])

const monthlyData = [
  { m: 'Jan', rev: 54, exp: 32 },
  { m: 'Feb', rev: 48, exp: 30 },
  { m: 'Mar', rev: 55, exp: 37 },
  { m: 'Apr', rev: 61, exp: 42 },
  { m: 'May', rev: 63, exp: 43 },
  { m: 'Jun', rev: 68, exp: 41 },
  { m: 'Jul', rev: 62, exp: 39 },
  { m: 'Aug', rev: 64, exp: 37 },
  { m: 'Sep', rev: 60, exp: 36 },
  { m: 'Oct', rev: 74, exp: 44 },
  { m: 'Nov', rev: 69, exp: 45 },
  { m: 'Dec', rev: 73, exp: 47 }
]

// -----------------------------------------------------------------------------
// STATE: JOURNAL ENTRIES
// -----------------------------------------------------------------------------
const journals = ref<any[]>([])
const journalSearch = ref('')
const journalStatusFilter = ref('all')
const selectedJournalLines = ref<any[]>([])
const selectedJournalRef = ref('')
const showLinesModal = ref(false)
const showNewJournalModal = ref(false)

const newJournalForm = ref({
  reference_no: '',
  entry_date: new Date().toISOString().split('T')[0],
  description: '',
  lines: [
    { account_id: '', description: '', debit: 0, credit: 0 },
    { account_id: '', description: '', debit: 0, credit: 0 }
  ]
})

const newJournalTotalDebit = computed(() => {
  return newJournalForm.value.lines.reduce((sum, l) => sum + (Number(l.debit) || 0), 0)
})

const newJournalTotalCredit = computed(() => {
  return newJournalForm.value.lines.reduce((sum, l) => sum + (Number(l.credit) || 0), 0)
})

const isNewJournalBalanced = computed(() => {
  return (
    newJournalTotalDebit.value > 0 &&
    newJournalTotalDebit.value === newJournalTotalCredit.value
  )
})

const filteredJournals = computed(() => {
  return journals.value.filter(j => {
    const matchSearch = 
      j.reference_no.toLowerCase().includes(journalSearch.value.toLowerCase()) ||
      j.description.toLowerCase().includes(journalSearch.value.toLowerCase())
    const matchStatus = journalStatusFilter.value === 'all' || j.status === journalStatusFilter.value
    return matchSearch && matchStatus
  })
})

const fetchJournals = async () => {
  try {
    const res = await axios.get('/api/v1/finance/journals')
    if (res.data?.data) {
      journals.value = res.data.data
      recentJournals.value = res.data.data.slice(0, 5)
    }
  } catch (err) {
    console.error('Failed to load journals:', err)
  }
}

const viewJournalLines = async (journal: any) => {
  selectedJournalRef.value = journal.reference_no
  try {
    const res = await axios.get(`/api/v1/finance/journals/${journal.id}/lines`)
    selectedJournalLines.value = res.data?.data || []
    showLinesModal.value = true
  } catch (err) {
    console.error('Failed to load journal lines:', err)
    showNotify('Gagal memuat rincian jurnal baris', 'error')
  }
}

const addJournalLine = () => {
  newJournalForm.value.lines.push({
    account_id: accounts.value[0]?.id || '',
    description: '',
    debit: 0,
    credit: 0
  })
}

const removeJournalLine = (idx: number) => {
  if (newJournalForm.value.lines.length > 2) {
    newJournalForm.value.lines.splice(idx, 1)
  }
}

const submitNewJournal = async () => {
  if (!isNewJournalBalanced.value) {
    showNotify('Jurnal harus seimbang: Total Debit harus sama dengan Total Credit!', 'error')
    return
  }
  try {
    loading.value = true
    await axios.post('/api/v1/finance/journals', newJournalForm.value)
    showNotify('Jurnal umum berhasil diposting ke Buku Besar!')
    showNewJournalModal.value = false
    newJournalForm.value = {
      reference_no: '',
      entry_date: new Date().toISOString().split('T')[0],
      description: '',
      lines: [
        { account_id: '', description: '', debit: 0, credit: 0 },
        { account_id: '', description: '', debit: 0, credit: 0 }
      ]
    }
    await fetchJournals()
    await fetchOverview()
  } catch (err: any) {
    console.error(err)
    showNotify(err.response?.data?.message || 'Gagal memposting jurnal', 'error')
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// STATE: CHART OF ACCOUNTS (COA)
// -----------------------------------------------------------------------------
const accounts = ref<any[]>([])
const coaSearch = ref('')
const coaTypeFilter = ref('all')
const showNewAccountModal = ref(false)

const newAccountForm = ref({
  code: '',
  name: '',
  account_type: 'asset',
  description: '',
  is_active: true
})

const filteredAccounts = computed(() => {
  return accounts.value.filter(a => {
    const matchSearch = 
      a.code.toLowerCase().includes(coaSearch.value.toLowerCase()) ||
      a.name.toLowerCase().includes(coaSearch.value.toLowerCase())
    const matchType = coaTypeFilter.value === 'all' || a.account_type === coaTypeFilter.value
    return matchSearch && matchType
  })
})

const fetchAccounts = async () => {
  try {
    const res = await axios.get('/api/v1/finance/coa')
    if (res.data?.data) {
      accounts.value = res.data.data
      if (newJournalForm.value.lines[0].account_id === '' && accounts.value.length > 0) {
        newJournalForm.value.lines[0].account_id = accounts.value[0].id
        newJournalForm.value.lines[1].account_id = accounts.value[1]?.id || accounts.value[0].id
      }
    }
  } catch (err) {
    console.error('Failed to load COA:', err)
  }
}

const submitNewAccount = async () => {
  if (!newAccountForm.value.code || !newAccountForm.value.name) {
    showNotify('Kode dan Nama Akun wajib diisi!', 'error')
    return
  }
  try {
    loading.value = true
    await axios.post('/api/v1/finance/coa', newAccountForm.value)
    showNotify('Bagan Akun (COA) berhasil ditambahkan!')
    showNewAccountModal.value = false
    newAccountForm.value = {
      code: '',
      name: '',
      account_type: 'asset',
      description: '',
      is_active: true
    }
    await fetchAccounts()
  } catch (err: any) {
    showNotify(err.response?.data?.message || 'Gagal menambah akun', 'error')
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// STATE: EXPENSES & PETTY CASH
// -----------------------------------------------------------------------------
const expenses = ref<any[]>([])
const expenseCategoryFilter = ref('all')
const expenseStatusFilter = ref('all')
const showNewExpenseModal = ref(false)

const newExpenseForm = ref({
  category: 'Operasional & Kas Kecil',
  account_id: '',
  amount: 0,
  expense_date: new Date().toISOString().split('T')[0],
  submitted_by: 'Kasir / Staff',
  description: '',
  receipt_url: ''
})

const filteredExpenses = computed(() => {
  return expenses.value.filter(e => {
    const matchCat = expenseCategoryFilter.value === 'all' || e.category === expenseCategoryFilter.value
    const matchStatus = expenseStatusFilter.value === 'all' || e.status === expenseStatusFilter.value
    return matchCat && matchStatus
  })
})

const totalExpenseMonth = computed(() => {
  return expenses.value.reduce((sum, e) => sum + (Number(e.amount) || 0), 0)
})

const pendingExpenseCount = computed(() => {
  return expenses.value.filter(e => e.status === 'submitted').length
})

const fetchExpenses = async () => {
  try {
    const res = await axios.get('/api/v1/finance/expenses')
    if (res.data?.data) {
      expenses.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load expenses:', err)
  }
}

const submitNewExpense = async () => {
  if (newExpenseForm.value.amount <= 0 || !newExpenseForm.value.description) {
    showNotify('Nominal harus lebih dari 0 dan Keterangan wajib diisi!', 'error')
    return
  }
  try {
    loading.value = true
    await axios.post('/api/v1/finance/expenses', newExpenseForm.value)
    showNotify('Pengeluaran berhasil diajukan!')
    showNewExpenseModal.value = false
    newExpenseForm.value = {
      category: 'Operasional & Kas Kecil',
      account_id: '',
      amount: 0,
      expense_date: new Date().toISOString().split('T')[0],
      submitted_by: 'Kasir / Staff',
      description: '',
      receipt_url: ''
    }
    await fetchExpenses()
  } catch (err: any) {
    showNotify(err.response?.data?.message || 'Gagal mengajukan pengeluaran', 'error')
  } finally {
    loading.value = false
  }
}

const updateExpenseStatus = async (id: string, status: 'approved' | 'rejected' | 'reimbursed') => {
  try {
    loading.value = true
    await axios.put(`/api/v1/finance/expenses/${id}/status`, {
      status,
      approved_by: 'Finance Manager'
    })
    const actionLabel = status === 'reimbursed' ? 'dicairkan & dijurnal' : status === 'approved' ? 'disetujui' : 'ditolak'
    showNotify(`Pengeluaran berhasil ${actionLabel}!`)
    await fetchExpenses()
    await fetchJournals()
    await fetchOverview()
  } catch (err: any) {
    showNotify(err.response?.data?.message || 'Gagal memperbarui status pengeluaran', 'error')
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// STATE: TAX MANAGEMENT
// -----------------------------------------------------------------------------
const taxData = ref({
  gross_sales: 0,
  pb1_tax: 0,
  ppn_masukan: 0,
  pph21_tax: 0,
  total_tax_liability: 0,
  period: 'September 2026',
  tax_transactions: [] as any[]
})

const showTaxSettleModal = ref(false)
const selectedTaxToSettle = ref<any>(null)
const taxNTPN = ref('')

const fetchTaxSummary = async () => {
  try {
    const res = await axios.get('/api/v1/finance/tax/summary')
    if (res.data?.data) {
      taxData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load tax summary:', err)
  }
}

const openTaxSettle = (item: any) => {
  selectedTaxToSettle.value = item
  taxNTPN.value = `NTPN-${Math.floor(10000000 + Math.random() * 90000000)}`
  showTaxSettleModal.value = true
}

const confirmTaxSettle = () => {
  if (selectedTaxToSettle.value) {
    selectedTaxToSettle.value.status = 'disetor'
    showNotify(`Pajak ${selectedTaxToSettle.value.tax_type} berhasil disetor dengan NTPN: ${taxNTPN.value}`)
    showTaxSettleModal.value = false
  }
}

// -----------------------------------------------------------------------------
// STATE: BUDGETING & COST CONTROL
// -----------------------------------------------------------------------------
const budgetData = ref({
  period: '2026-09',
  total_budget: 0,
  total_actual: 0,
  remaining: 0,
  utilization_pct: 0,
  items: [] as any[]
})

const showNewBudgetModal = ref(false)
const newBudgetForm = ref({
  period: '2026-09',
  category: 'Bahan Baku & Minuman',
  budget_amount: 0,
  notes: ''
})

const fetchBudgets = async () => {
  try {
    const res = await axios.get('/api/v1/finance/budgets')
    if (res.data?.data) {
      budgetData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load budgets:', err)
  }
}

const submitNewBudget = async () => {
  if (newBudgetForm.value.budget_amount <= 0) {
    showNotify('Nominal anggaran harus lebih dari 0!', 'error')
    return
  }
  try {
    loading.value = true
    await axios.post('/api/v1/finance/budgets', newBudgetForm.value)
    showNotify('Anggaran berhasil diperbarui!')
    showNewBudgetModal.value = false
    await fetchBudgets()
  } catch (err: any) {
    showNotify(err.response?.data?.message || 'Gagal menyimpan anggaran', 'error')
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// STATE: BANK RECONCILIATION
// -----------------------------------------------------------------------------
const reconData = ref({
  bank_name: 'Bank BCA Operasional (1102)',
  period: '2026-09',
  statement_balance: 0,
  erp_book_balance: 0,
  difference: 0,
  matched_count: 0,
  unmatched_count: 0,
  statements: [] as any[]
})

const selectedBank = ref('Bank BCA Operasional (1102)')

const fetchReconciliation = async () => {
  try {
    const res = await axios.get(`/api/v1/finance/reconciliation?bank_name=${encodeURIComponent(selectedBank.value)}`)
    if (res.data?.data) {
      reconData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load bank reconciliation:', err)
  }
}

const toggleMatchReconciliation = async (item: any) => {
  const newStatus = item.match_status === 'matched' ? 'unmatched' : 'matched'
  try {
    await axios.post('/api/v1/finance/reconciliation/match', {
      id: item.id,
      match_status: newStatus
    })
    item.match_status = newStatus
    if (newStatus === 'matched') {
      reconData.value.matched_count++
      reconData.value.unmatched_count--
      showNotify(`Transaksi ${item.ref_number} ditandai Cocok (Matched)!`)
    } else {
      reconData.value.matched_count--
      reconData.value.unmatched_count++
      showNotify(`Status transaksi ${item.ref_number} dikembalikan ke Unmatched.`)
    }
  } catch (err) {
    console.error(err)
    showNotify('Gagal memperbarui status rekonsiliasi', 'error')
  }
}

const autoMatchAll = () => {
  reconData.value.statements.forEach(s => {
    if (s.match_status === 'unmatched') {
      s.match_status = 'matched'
      reconData.value.matched_count++
      reconData.value.unmatched_count--
    }
  })
  showNotify('Seluruh mutasi berhasil dicocokkan otomatis dengan buku besar!')
}

// -----------------------------------------------------------------------------
// OVERVIEW API
// -----------------------------------------------------------------------------
const fetchOverview = async () => {
  try {
    const res = await axios.get('/api/v1/finance/overview')
    if (res.data?.data) {
      overviewData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load overview:', err)
  }
}

// -----------------------------------------------------------------------------
// INITIAL MOUNT
// -----------------------------------------------------------------------------
onMounted(async () => {
  loading.value = true
  await Promise.all([
    fetchOverview(),
    fetchJournals(),
    fetchAccounts(),
    fetchExpenses(),
    fetchTaxSummary(),
    fetchBudgets(),
    fetchReconciliation()
  ])
  loading.value = false
})
</script>

<template>
  <div class="space-y-6">
    <!-- Notification Toast -->
    <div 
      v-if="notification"
      class="fixed top-5 right-5 z-50 flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg border text-sm font-semibold transition-all transform duration-300"
      :class="notification.type === 'success' ? 'bg-emerald-50 border-emerald-300 text-emerald-800' : 'bg-rose-50 border-rose-300 text-rose-800'"
    >
      <CheckCircle2 v-if="notification.type === 'success'" class="w-5 h-5 text-emerald-600" />
      <AlertCircle v-else class="w-5 h-5 text-rose-600" />
      <span>{{ notification.message }}</span>
    </div>

    <!-- Breadcrumb & Top Bar -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <p class="text-xs font-semibold text-slate-400 mb-1">/ Finance and Accounting</p>
        <h1 class="text-2xl font-bold text-slate-900 tracking-tight">Financial Overview</h1>
      </div>
      <div class="flex items-center gap-3">
        <button 
          @click="showNewExpenseModal = true"
          class="inline-flex items-center gap-2 px-3.5 py-2 bg-white hover:bg-slate-50 text-slate-700 border border-slate-200 rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <CreditCard class="w-4 h-4 text-slate-500" />
          Ajukan Pengeluaran
        </button>
        <button 
          @click="showNewJournalModal = true"
          class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-sm font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          Buat Jurnal Baru
        </button>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="border-b border-slate-200">
      <nav class="flex space-x-8 overflow-x-auto pb-px">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm transition-colors cursor-pointer"
          :class="[
            activeTab === tab.id
              ? 'border-blue-600 text-blue-600 font-semibold'
              : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
          ]"
        >
          {{ tab.label }}
        </button>
      </nav>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 1: FINANCIAL OVERVIEW                                               -->
    <!-- ======================================================================= -->
    <div v-if="activeTab === 'overview'" class="space-y-6">
      <!-- 4 Key Metrics Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Total Revenue -->
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-blue-50 flex items-center justify-center text-blue-600">
            <Receipt class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-medium text-slate-500">Total Revenue</p>
            <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(overviewData.total_revenue) }}</p>
          </div>
        </div>

        <!-- Total Expenses -->
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-rose-50 flex items-center justify-center text-rose-600">
            <CreditCard class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-medium text-slate-500">Total Expenses</p>
            <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(overviewData.total_expenses) }}</p>
          </div>
        </div>

        <!-- Net Profit -->
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-emerald-50 flex items-center justify-center text-emerald-600">
            <TrendingUp class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-medium text-slate-500">Net Profit</p>
            <p class="text-xl font-bold" :class="overviewData.net_profit >= 0 ? 'text-slate-900' : 'text-rose-600'">
              Rp {{ formatNum(overviewData.net_profit) }}
            </p>
          </div>
        </div>

        <!-- Cash on Hand -->
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-slate-100 flex items-center justify-center text-slate-600">
            <Wallet class="w-6 h-6" />
          </div>
          <div>
            <p class="text-xs font-medium text-slate-500">Cash on Hand</p>
            <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(overviewData.cash_on_hand) }}</p>
          </div>
        </div>
      </div>

      <!-- Middle Row: Charts -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Revenue vs Expenses Bar Chart -->
        <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-xs">
          <div class="flex items-center justify-between mb-6">
            <h3 class="font-bold text-slate-900">Revenue vs Expenses</h3>
            <div class="flex items-center gap-4 text-xs">
              <span class="flex items-center gap-1.5 text-slate-600">
                <span class="w-3 h-3 rounded-xs bg-blue-600"></span> Revenue
              </span>
              <span class="flex items-center gap-1.5 text-slate-600">
                <span class="w-3 h-3 rounded-xs bg-slate-800"></span> Expenses
              </span>
            </div>
          </div>

          <div class="h-64 flex flex-col justify-end">
            <div class="h-52 flex items-end justify-between gap-2 pt-4 px-2 border-b border-slate-200">
              <div v-for="item in monthlyData" :key="item.m" class="flex flex-col items-center flex-1 h-full justify-end group relative">
                <div class="absolute -top-10 opacity-0 group-hover:opacity-100 transition-opacity bg-slate-900 text-white text-[10px] py-1 px-2 rounded pointer-events-none whitespace-nowrap z-10">
                  Rev: {{ item.rev }}M | Exp: {{ item.exp }}M
                </div>
                <div class="w-full flex items-end justify-center gap-1">
                  <div class="w-2 sm:w-3 bg-blue-600 rounded-t-xs transition-all duration-300 group-hover:bg-blue-500" :style="{ height: `${item.rev * 2.2}px` }"></div>
                  <div class="w-2 sm:w-3 bg-slate-800 rounded-t-xs transition-all duration-300 group-hover:bg-slate-700" :style="{ height: `${item.exp * 2.2}px` }"></div>
                </div>
                <span class="text-[10px] text-slate-400 mt-2 font-medium">{{ item.m }}</span>
              </div>
            </div>
            <div class="flex justify-between text-[10px] text-slate-400 pt-1">
              <span>0</span>
              <span>20M</span>
              <span>40M</span>
              <span>60M</span>
              <span>80M</span>
            </div>
          </div>
        </div>

        <!-- Cash Flow Trend -->
        <div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-xs flex flex-col justify-between">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-bold text-slate-900">Cash Flow Trend</h3>
            <div class="flex items-center gap-3 text-xs">
              <span class="flex items-center gap-1.5 text-slate-600"><span class="w-3 h-0.5 bg-sky-400"></span> Inflow</span>
              <span class="flex items-center gap-1.5 text-slate-600"><span class="w-3 h-0.5 bg-blue-600"></span> Net Flow</span>
            </div>
          </div>

          <div class="relative h-64 flex flex-col justify-end">
            <div class="h-52 w-full pt-4">
              <svg class="w-full h-full" viewBox="0 0 500 200" preserveAspectRatio="none">
                <defs>
                  <linearGradient id="gradFlow" x1="0%" y1="0%" x2="0%" y2="100%">
                    <stop offset="0%" stop-color="#38bdf8" stop-opacity="0.3" />
                    <stop offset="100%" stop-color="#38bdf8" stop-opacity="0.0" />
                  </linearGradient>
                  <linearGradient id="gradNet" x1="0%" y1="0%" x2="0%" y2="100%">
                    <stop offset="0%" stop-color="#2563eb" stop-opacity="0.4" />
                    <stop offset="100%" stop-color="#2563eb" stop-opacity="0.0" />
                  </linearGradient>
                </defs>
                <line x1="0" y1="50" x2="500" y2="50" stroke="#f1f5f9" stroke-width="1" />
                <line x1="0" y1="100" x2="500" y2="100" stroke="#f1f5f9" stroke-width="1" />
                <line x1="0" y1="150" x2="500" y2="150" stroke="#f1f5f9" stroke-width="1" />
                <path d="M0,160 Q50,110 100,125 T200,80 T300,60 T400,90 T500,50 L500,200 L0,200 Z" fill="url(#gradFlow)" />
                <path d="M0,160 Q50,110 100,125 T200,80 T300,60 T400,90 T500,50" fill="none" stroke="#38bdf8" stroke-width="2.5" />
                <path d="M0,180 Q50,140 100,150 T200,110 T300,95 T400,120 T500,80 L500,200 L0,200 Z" fill="url(#gradNet)" />
                <path d="M0,180 Q50,140 100,150 T200,110 T300,95 T400,120 T500,80" fill="none" stroke="#2563eb" stroke-width="2.5" />
              </svg>
            </div>
            <div class="flex justify-between text-[10px] text-slate-400 pt-1">
              <span>Jan</span>
              <span>Mar</span>
              <span>May</span>
              <span>Jul</span>
              <span>Sep</span>
              <span>Nov</span>
              <span>Dec</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Bottom Row: Recent Journal Entries & Expense Breakdown -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Recent Journal Entries -->
        <div class="lg:col-span-2 bg-white rounded-2xl border border-slate-200 shadow-xs p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-bold text-slate-900">Recent Journal Entries</h3>
            <button @click="activeTab = 'journal'" class="text-xs text-blue-600 font-semibold hover:underline cursor-pointer">
              View All Entries →
            </button>
          </div>

          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
                  <th class="py-3 px-3">Date</th>
                  <th class="py-3 px-3">Reference</th>
                  <th class="py-3 px-3">Description</th>
                  <th class="py-3 px-3 text-right">Debit (Rp)</th>
                  <th class="py-3 px-3 text-right">Credit (Rp)</th>
                  <th class="py-3 px-3 text-center">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 text-xs">
                <tr v-for="(entry, idx) in recentJournals" :key="idx" class="hover:bg-slate-50/70 transition-colors">
                  <td class="py-3 px-3 text-slate-600">{{ entry.entry_date || entry.date }}</td>
                  <td class="py-3 px-3 font-semibold text-slate-900">{{ entry.reference_no || entry.ref }}</td>
                  <td class="py-3 px-3 text-slate-700 max-w-xs truncate">{{ entry.description || entry.desc }}</td>
                  <td class="py-3 px-3 text-right font-medium text-slate-900">
                    {{ entry.total_debit || entry.debit ? formatNum(entry.total_debit || entry.debit) : '-' }}
                  </td>
                  <td class="py-3 px-3 text-right font-medium text-slate-900">
                    {{ entry.total_credit || entry.credit ? formatNum(entry.total_credit || entry.credit) : '-' }}
                  </td>
                  <td class="py-3 px-3 text-center">
                    <span class="inline-block px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                      {{ entry.status }}
                    </span>
                  </td>
                </tr>
                <tr v-if="recentJournals.length === 0">
                  <td colspan="6" class="py-8 text-center text-slate-400 text-xs">
                    Tidak ada catatan jurnal transaksi.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Expense Breakdown Donut -->
        <div class="bg-white rounded-2xl border border-slate-200 shadow-xs p-6 flex flex-col justify-between">
          <h3 class="font-bold text-slate-900 mb-2">Expense Breakdown</h3>
          <div class="relative flex items-center justify-center my-4">
            <svg class="w-48 h-48 transform -rotate-90" viewBox="0 0 100 100">
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#1e3a8a" stroke-width="16" stroke-dasharray="107.4 238.8" stroke-dashoffset="0" />
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#2563eb" stroke-width="16" stroke-dasharray="71.6 238.8" stroke-dashoffset="-107.4" />
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#38bdf8" stroke-width="16" stroke-dasharray="23.9 238.8" stroke-dashoffset="-179" />
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#34d399" stroke-width="16" stroke-dasharray="19.1 238.8" stroke-dashoffset="-202.9" />
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#ec4899" stroke-width="16" stroke-dasharray="16.7 238.8" stroke-dashoffset="-222" />
            </svg>
            <div class="absolute text-center">
              <span class="text-xs text-slate-400 font-medium">Top Cost</span>
              <p class="text-lg font-black text-slate-900">45%</p>
              <span class="text-[10px] text-slate-500">Ingredients</span>
            </div>
          </div>
          <div class="flex flex-wrap items-center justify-center gap-3 text-xs pt-2">
            <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-blue-900"></span> Ingredients 45%</span>
            <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span> Salary 30%</span>
            <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-sky-400"></span> Rent 10%</span>
            <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-emerald-400"></span> Utilities 8%</span>
            <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-full bg-pink-500"></span> Others 7%</span>
          </div>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 2: JOURNAL ENTRIES (BUKU JURNAL UMUM)                               -->
    <!-- ======================================================================= -->
    <div v-if="activeTab === 'journal'" class="space-y-6">
      <!-- Search & Action Bar -->
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white p-4 rounded-2xl border border-slate-200 shadow-xs">
        <div class="flex flex-1 items-center gap-3 w-full sm:w-auto">
          <div class="relative flex-1 max-w-md">
            <Search class="w-4 h-4 absolute left-3.5 top-3 text-slate-400" />
            <input 
              v-model="journalSearch"
              type="text"
              placeholder="Cari nomor referensi, deskripsi transaksi..."
              class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-xl text-xs focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <select 
            v-model="journalStatusFilter"
            class="px-3 py-2 border border-slate-200 rounded-xl text-xs bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">Semua Status</option>
            <option value="posted">Posted</option>
            <option value="draft">Draft</option>
          </select>
        </div>
        <div class="flex items-center gap-3">
          <button 
            @click="showNewJournalModal = true"
            class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
          >
            <Plus class="w-4 h-4" />
            Jurnal Baru
          </button>
        </div>
      </div>

      <!-- Journal Table -->
      <div class="bg-white rounded-2xl border border-slate-200 shadow-xs overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
                <th class="py-3 px-4">Tanggal</th>
                <th class="py-3 px-4">No. Referensi</th>
                <th class="py-3 px-4">Keterangan Jurnal</th>
                <th class="py-3 px-4 text-right">Total Debit (Rp)</th>
                <th class="py-3 px-4 text-right">Total Credit (Rp)</th>
                <th class="py-3 px-4 text-center">Status</th>
                <th class="py-3 px-4 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-xs">
              <tr v-for="j in filteredJournals" :key="j.id" class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3.5 px-4 text-slate-600 whitespace-nowrap">{{ j.entry_date }}</td>
                <td class="py-3.5 px-4 font-bold text-slate-900 whitespace-nowrap">{{ j.reference_no }}</td>
                <td class="py-3.5 px-4 text-slate-700 max-w-sm">{{ j.description }}</td>
                <td class="py-3.5 px-4 text-right font-semibold text-emerald-600 whitespace-nowrap">
                  {{ formatNum(j.total_debit) }}
                </td>
                <td class="py-3.5 px-4 text-right font-semibold text-blue-600 whitespace-nowrap">
                  {{ formatNum(j.total_credit) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                    {{ j.status }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-center">
                  <button 
                    @click="viewJournalLines(j)"
                    class="inline-flex items-center gap-1.5 px-2.5 py-1 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg text-[11px] font-medium transition-colors cursor-pointer"
                  >
                    <Eye class="w-3.5 h-3.5" />
                    Rincian
                  </button>
                </td>
              </tr>
              <tr v-if="filteredJournals.length === 0">
                <td colspan="7" class="py-12 text-center text-slate-400">
                  Tidak ditemukan jurnal yang cocok dengan kriteria pencarian.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 3: CHART OF ACCOUNTS (BAGAN AKUN)                                   -->
    <!-- ======================================================================= -->
    <div v-if="activeTab === 'coa'" class="space-y-6">
      <!-- Search & Filters -->
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white p-4 rounded-2xl border border-slate-200 shadow-xs">
        <div class="flex flex-1 items-center gap-3 w-full sm:w-auto">
          <div class="relative flex-1 max-w-md">
            <Search class="w-4 h-4 absolute left-3.5 top-3 text-slate-400" />
            <input 
              v-model="coaSearch"
              type="text"
              placeholder="Cari kode akun, nama bagan akun..."
              class="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-xl text-xs focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <select 
            v-model="coaTypeFilter"
            class="px-3 py-2 border border-slate-200 rounded-xl text-xs bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">Semua Tipe Akun</option>
            <option value="asset">Aset (1xxx)</option>
            <option value="liability">Kewajiban / Hutang (2xxx)</option>
            <option value="equity">Ekuitas / Modal (3xxx)</option>
            <option value="revenue">Pendapatan (4xxx)</option>
            <option value="expense">Beban & HPP (5xxx/6xxx)</option>
          </select>
        </div>
        <div>
          <button 
            @click="showNewAccountModal = true"
            class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
          >
            <Plus class="w-4 h-4" />
            Tambah Akun
          </button>
        </div>
      </div>

      <!-- COA Table -->
      <div class="bg-white rounded-2xl border border-slate-200 shadow-xs overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
                <th class="py-3 px-4">Kode Akun</th>
                <th class="py-3 px-4">Nama Akun</th>
                <th class="py-3 px-4">Klasifikasi Tipe</th>
                <th class="py-3 px-4">Keterangan</th>
                <th class="py-3 px-4 text-center">Saldo Normal</th>
                <th class="py-3 px-4 text-right">Saldo Berjalan (Rp)</th>
                <th class="py-3 px-4 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-xs">
              <tr v-for="a in filteredAccounts" :key="a.id" class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3.5 px-4 font-mono font-bold text-slate-900">{{ a.code }}</td>
                <td class="py-3.5 px-4 font-semibold text-slate-800">{{ a.name }}</td>
                <td class="py-3.5 px-4">
                  <span 
                    class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-semibold uppercase"
                    :class="{
                      'bg-blue-50 text-blue-700 border border-blue-200': a.account_type === 'asset',
                      'bg-amber-50 text-amber-700 border border-amber-200': a.account_type === 'liability',
                      'bg-purple-50 text-purple-700 border border-purple-200': a.account_type === 'equity',
                      'bg-emerald-50 text-emerald-700 border border-emerald-200': a.account_type === 'revenue',
                      'bg-rose-50 text-rose-700 border border-rose-200': a.account_type === 'expense'
                    }"
                  >
                    {{ a.account_type }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-slate-500 max-w-xs truncate">{{ a.description || '-' }}</td>
                <td class="py-3.5 px-4 text-center uppercase font-mono text-[11px] text-slate-600">
                  {{ a.normal_balance }}
                </td>
                <td class="py-3.5 px-4 text-right font-bold" :class="a.balance >= 0 ? 'text-slate-900' : 'text-rose-600'">
                  Rp {{ formatNum(a.balance) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <span class="inline-block px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                    Aktif
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 4: EXPENSES & PETTY CASH (PENGELUARAN KAS KECIL)                     -->
    <!-- ======================================================================= -->
    <div v-if="activeTab === 'expenses'" class="space-y-6">
      <!-- 3 Summary Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-slate-500">Total Pengeluaran Bulan Ini</p>
            <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(totalExpenseMonth) }}</p>
          </div>
          <div class="w-10 h-10 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center">
            <CreditCard class="w-5 h-5" />
          </div>
        </div>
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-slate-500">Menunggu Persetujuan</p>
            <p class="text-xl font-bold text-amber-600">{{ pendingExpenseCount }} Pengajuan</p>
          </div>
          <div class="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center">
            <AlertTriangle class="w-5 h-5" />
          </div>
        </div>
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-slate-500">Disetujui / Dicairkan</p>
            <p class="text-xl font-bold text-emerald-600">
              {{ expenses.filter(e => e.status === 'approved' || e.status === 'reimbursed').length }} Transaksi
            </p>
          </div>
          <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
            <CheckCircle2 class="w-5 h-5" />
          </div>
        </div>
      </div>

      <!-- Filters & Action -->
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white p-4 rounded-2xl border border-slate-200 shadow-xs">
        <div class="flex flex-wrap items-center gap-3">
          <select 
            v-model="expenseCategoryFilter"
            class="px-3 py-2 border border-slate-200 rounded-xl text-xs bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">Semua Kategori</option>
            <option value="Operasional & Kas Kecil">Operasional & Kas Kecil</option>
            <option value="Bahan Baku">Bahan Baku</option>
            <option value="Pemeliharaan Peralatan">Pemeliharaan Peralatan</option>
            <option value="Promosi & Pemasaran">Promosi & Pemasaran</option>
          </select>
          <select 
            v-model="expenseStatusFilter"
            class="px-3 py-2 border border-slate-200 rounded-xl text-xs bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">Semua Status</option>
            <option value="submitted">Submitted</option>
            <option value="approved">Approved</option>
            <option value="reimbursed">Reimbursed</option>
            <option value="rejected">Rejected</option>
          </select>
        </div>
        <button 
          @click="showNewExpenseModal = true"
          class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          Ajukan Pengeluaran
        </button>
      </div>

      <!-- Expenses Table -->
      <div class="bg-white rounded-2xl border border-slate-200 shadow-xs overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
                <th class="py-3 px-4">Tanggal</th>
                <th class="py-3 px-4">No. Pengeluaran</th>
                <th class="py-3 px-4">Kategori</th>
                <th class="py-3 px-4">Keterangan</th>
                <th class="py-3 px-4">Diajukan Oleh</th>
                <th class="py-3 px-4 text-right">Nominal (Rp)</th>
                <th class="py-3 px-4 text-center">Status</th>
                <th class="py-3 px-4 text-center">Aksi Manajemen</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-xs">
              <tr v-for="e in filteredExpenses" :key="e.id" class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3.5 px-4 text-slate-600 whitespace-nowrap">{{ e.expense_date }}</td>
                <td class="py-3.5 px-4 font-mono font-bold text-slate-900 whitespace-nowrap">{{ e.expense_number }}</td>
                <td class="py-3.5 px-4">
                  <span class="inline-block px-2 py-0.5 rounded-md text-[10px] font-semibold bg-slate-100 text-slate-800">
                    {{ e.category }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-slate-700 max-w-xs">{{ e.description }}</td>
                <td class="py-3.5 px-4 text-slate-600">{{ e.submitted_by }}</td>
                <td class="py-3.5 px-4 text-right font-bold text-slate-900 whitespace-nowrap">
                  Rp {{ formatNum(e.amount) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <span 
                    class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-semibold uppercase"
                    :class="{
                      'bg-amber-50 text-amber-700 border border-amber-200': e.status === 'submitted',
                      'bg-blue-50 text-blue-700 border border-blue-200': e.status === 'approved',
                      'bg-emerald-50 text-emerald-700 border border-emerald-200': e.status === 'reimbursed',
                      'bg-rose-50 text-rose-700 border border-rose-200': e.status === 'rejected'
                    }"
                  >
                    {{ e.status }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-center whitespace-nowrap">
                  <div class="flex items-center justify-center gap-1.5">
                    <button 
                      v-if="e.status === 'submitted'"
                      @click="updateExpenseStatus(e.id, 'approved')"
                      class="px-2 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded text-[10px] font-semibold transition-colors cursor-pointer"
                    >
                      Setujui
                    </button>
                    <button 
                      v-if="e.status === 'approved'"
                      @click="updateExpenseStatus(e.id, 'reimbursed')"
                      class="px-2 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded text-[10px] font-semibold transition-colors cursor-pointer"
                    >
                      Cairkan Kas
                    </button>
                    <button 
                      v-if="e.status === 'submitted'"
                      @click="updateExpenseStatus(e.id, 'rejected')"
                      class="px-2 py-1 bg-rose-50 hover:bg-rose-100 text-rose-700 border border-rose-200 rounded text-[10px] font-semibold transition-colors cursor-pointer"
                    >
                      Tolak
                    </button>
                    <span v-if="e.status === 'reimbursed'" class="text-[11px] text-emerald-600 font-medium">
                      ✓ Terjurnal
                    </span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 5: TAX MANAGEMENT (MANAJEMEN PAJAK KAFE)                            -->
    <!-- ======================================================================= -->
    <div v-if="activeTab === 'tax'" class="space-y-6">
      <!-- 4 Tax Summary Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold text-slate-500">Pajak Restoran PB1 (10%)</span>
            <span class="p-2 bg-blue-50 text-blue-600 rounded-xl"><Percent class="w-4 h-4" /></span>
          </div>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(taxData.pb1_tax) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Dasar Pengenaan: Penjualan F&B Kasir POS</p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold text-slate-500">PPN Masukan (11%)</span>
            <span class="p-2 bg-emerald-50 text-emerald-600 rounded-xl"><ArrowDownLeft class="w-4 h-4" /></span>
          </div>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(taxData.ppn_masukan) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Faktur Pajak Pembelian Bahan P2P</p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold text-slate-500">PPh 21 Karyawan</span>
            <span class="p-2 bg-purple-50 text-purple-600 rounded-xl"><Building class="w-4 h-4" /></span>
          </div>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(taxData.pph21_tax) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Pemotongan PPh atas Penggajian HRIS</p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold text-slate-500">Total Pajak Terutang</span>
            <span class="p-2 bg-rose-50 text-rose-600 rounded-xl"><Scale class="w-4 h-4" /></span>
          </div>
          <p class="text-xl font-bold text-rose-600">Rp {{ formatNum(taxData.total_tax_liability) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Wajib disetor paling lambat tgl 15 bulan depan</p>
        </div>
      </div>

      <!-- Info Banner Peraturan Pajak -->
      <div class="bg-blue-50/60 border border-blue-200 rounded-2xl p-5 flex items-start gap-4">
        <ShieldCheck class="w-6 h-6 text-blue-600 shrink-0 mt-0.5" />
        <div class="text-xs text-blue-900 leading-relaxed">
          <span class="font-bold text-sm block mb-1">Ketentuan Kepatuhan Pajak Usaha Restoran & Kafe:</span>
          Sesuai UU Hubungan Keuangan Pemerintah Pusat dan Daerah (HKPD No. 1/2022), Kafe memungut **Pajak Barang dan Jasa Tertentu (PB1)** sebesar 10% dari konsumen akhir untuk disetorkan ke Bapenda Daerah. Faktur Pembelian Supplier resmi diakui sebagai **PPN Masukan (11%)**, dan Pajak Penghasilan Karyawan disetorkan via **e-Billing DJP PPh Pasal 21**.
        </div>
      </div>

      <!-- Tax Transactions Table -->
      <div class="bg-white rounded-2xl border border-slate-200 shadow-xs overflow-hidden">
        <div class="p-4 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-bold text-slate-900 text-sm">Rekapitulasi Transaksi Pajak (Periode: {{ taxData.period }})</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
                <th class="py-3 px-4">Jenis Pajak</th>
                <th class="py-3 px-4">Referensi Dokumen</th>
                <th class="py-3 px-4 text-right">Dasar Pengenaan (DPP)</th>
                <th class="py-3 px-4 text-center">Tarif</th>
                <th class="py-3 px-4 text-right">Pajak Terhitung (Rp)</th>
                <th class="py-3 px-4">Jatuh Tempo</th>
                <th class="py-3 px-4 text-center">Status</th>
                <th class="py-3 px-4 text-center">Aksi Setor</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-xs">
              <tr v-for="t in taxData.tax_transactions" :key="t.id" class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3.5 px-4 font-bold text-slate-900">{{ t.tax_type }}</td>
                <td class="py-3.5 px-4 font-mono text-slate-600">{{ t.reference }}</td>
                <td class="py-3.5 px-4 text-right text-slate-700">Rp {{ formatNum(t.taxable_amount) }}</td>
                <td class="py-3.5 px-4 text-center font-semibold text-slate-800">{{ t.rate }}%</td>
                <td class="py-3.5 px-4 text-right font-bold text-slate-900">Rp {{ formatNum(t.tax_amount) }}</td>
                <td class="py-3.5 px-4 text-slate-600">{{ t.due_date }}</td>
                <td class="py-3.5 px-4 text-center">
                  <span 
                    class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-semibold uppercase"
                    :class="{
                      'bg-rose-50 text-rose-700 border border-rose-200': t.status === 'terutang',
                      'bg-blue-50 text-blue-700 border border-blue-200': t.status === 'dapat_dikreditkan',
                      'bg-emerald-50 text-emerald-700 border border-emerald-200': t.status === 'disetor'
                    }"
                  >
                    {{ t.status }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-center">
                  <button 
                    v-if="t.status === 'terutang'"
                    @click="openTaxSettle(t)"
                    class="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-[11px] font-semibold transition-colors cursor-pointer"
                  >
                    Setor Pajak
                  </button>
                  <span v-else class="text-[11px] text-emerald-600 font-medium">✓ Selesai</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 6: BUDGETING & COST CONTROL (ANGGARAN & REALISASI)                  -->
    <!-- ======================================================================= -->
    <div v-if="activeTab === 'budget'" class="space-y-6">
      <!-- 4 Budget KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <p class="text-xs font-medium text-slate-500 mb-1">Total Target Anggaran</p>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(budgetData.total_budget) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Periode: {{ budgetData.period }}</p>
        </div>
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <p class="text-xs font-medium text-slate-500 mb-1">Realisasi Pengeluaran</p>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(budgetData.total_actual) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Beban Buku Besar & Kas</p>
        </div>
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <p class="text-xs font-medium text-slate-500 mb-1">Sisa Anggaran Tersedia</p>
          <p class="text-xl font-bold" :class="budgetData.remaining >= 0 ? 'text-emerald-600' : 'text-rose-600'">
            Rp {{ formatNum(budgetData.remaining) }}
          </p>
          <p class="text-[11px] text-slate-400 mt-1">Sisa limit belanja operasional</p>
        </div>
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <p class="text-xs font-medium text-slate-500 mb-1">Persentase Utilisasi</p>
          <p class="text-xl font-bold" :class="budgetData.utilization_pct > 100 ? 'text-rose-600' : budgetData.utilization_pct > 80 ? 'text-amber-600' : 'text-blue-600'">
            {{ budgetData.utilization_pct.toFixed(1) }}%
          </p>
          <p class="text-[11px] text-slate-400 mt-1">Target utilisasi aman: &lt; 90%</p>
        </div>
      </div>

      <!-- Action Button -->
      <div class="flex justify-end">
        <button 
          @click="showNewBudgetModal = true"
          class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          Atur / Tambah Anggaran
        </button>
      </div>

      <!-- Budget Breakdown Progress Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div 
          v-for="b in budgetData.items" 
          :key="b.id"
          class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs space-y-3"
        >
          <div class="flex items-center justify-between">
            <div>
              <h4 class="font-bold text-slate-900 text-sm">{{ b.category }}</h4>
              <p class="text-[11px] text-slate-400">{{ b.notes }}</p>
            </div>
            <span 
              class="px-2.5 py-0.5 rounded-full text-[10px] font-semibold uppercase"
              :class="{
                'bg-emerald-50 text-emerald-700 border border-emerald-200': b.status === 'on_track',
                'bg-amber-50 text-amber-700 border border-amber-200': b.status === 'warning',
                'bg-rose-50 text-rose-700 border border-rose-200': b.status === 'over_budget'
              }"
            >
              {{ b.status === 'over_budget' ? 'Over Budget' : b.status === 'warning' ? 'Waspada' : 'Aman' }}
            </span>
          </div>

          <div class="space-y-1">
            <div class="flex justify-between text-xs font-semibold">
              <span class="text-slate-600">Realisasi: Rp {{ formatNum(b.actual_amount) }}</span>
              <span class="text-slate-900 font-bold">{{ b.utilization.toFixed(1) }}%</span>
            </div>
            <!-- Progress bar -->
            <div class="w-full bg-slate-100 rounded-full h-2.5 overflow-hidden">
              <div 
                class="h-2.5 rounded-full transition-all duration-500"
                :class="{
                  'bg-emerald-500': b.utilization <= 80,
                  'bg-amber-500': b.utilization > 80 && b.utilization <= 100,
                  'bg-rose-500': b.utilization > 100
                }"
                :style="{ width: `${Math.min(b.utilization, 100)}%` }"
              ></div>
            </div>
            <div class="flex justify-between text-[11px] text-slate-400 pt-0.5">
              <span>Target: Rp {{ formatNum(b.budget_amount) }}</span>
              <span :class="b.remaining >= 0 ? 'text-slate-500' : 'text-rose-500 font-bold'">
                {{ b.remaining >= 0 ? 'Sisa: Rp ' + formatNum(b.remaining) : 'Defisit: Rp ' + formatNum(Math.abs(b.remaining)) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 7: BANK RECONCILIATION (REKONSILIASI BANK)                           -->
    <!-- ======================================================================= -->
    <div v-if="activeTab === 'reconciliation'" class="space-y-6">
      <!-- 3 Balance Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <p class="text-xs font-medium text-slate-500 mb-1">Saldo Rekening Koran Bank</p>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(reconData.statement_balance) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">{{ reconData.bank_name }}</p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <p class="text-xs font-medium text-slate-500 mb-1">Saldo Buku Kas ERP (Akun 1102)</p>
          <p class="text-xl font-bold text-slate-900">Rp {{ formatNum(reconData.erp_book_balance) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">General Ledger Buku Besar</p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200 shadow-xs">
          <p class="text-xs font-medium text-slate-500 mb-1">Selisih Rekonsiliasi</p>
          <p class="text-xl font-bold text-blue-600">Rp {{ formatNum(reconData.difference) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Cocok: {{ reconData.matched_count }} | Belum: {{ reconData.unmatched_count }}</p>
        </div>
      </div>

      <!-- Action & Auto Match Toolbar -->
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white p-4 rounded-2xl border border-slate-200 shadow-xs">
        <div class="flex items-center gap-3">
          <select 
            v-model="selectedBank"
            @change="fetchReconciliation"
            class="px-3 py-2 border border-slate-200 rounded-xl text-xs bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="Bank BCA Operasional (1102)">Bank BCA Operasional (1102)</option>
            <option value="Bank Mandiri Operasional (1103)">Bank Mandiri Operasional (1103)</option>
          </select>
        </div>
        <div class="flex items-center gap-3">
          <button 
            @click="autoMatchAll"
            class="inline-flex items-center gap-2 px-3.5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
          >
            <CheckCircle2 class="w-4 h-4" />
            Auto-Match Sesuai Mutasi
          </button>
        </div>
      </div>

      <!-- Bank Statement Mutasi Table -->
      <div class="bg-white rounded-2xl border border-slate-200 shadow-xs overflow-hidden">
        <div class="p-4 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-bold text-slate-900 text-sm">Mutasi Rekening Koran vs Pencatatan Kas ERP</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-600">
                <th class="py-3 px-4">Tanggal Mutasi</th>
                <th class="py-3 px-4">No. Referensi Bank</th>
                <th class="py-3 px-4">Keterangan Transaksi</th>
                <th class="py-3 px-4 text-center">Jenis</th>
                <th class="py-3 px-4 text-right">Nominal (Rp)</th>
                <th class="py-3 px-4 text-center">Status Cocok</th>
                <th class="py-3 px-4 text-center">Aksi Rekonsiliasi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 text-xs">
              <tr v-for="s in reconData.statements" :key="s.id" class="hover:bg-slate-50/70 transition-colors">
                <td class="py-3.5 px-4 text-slate-600 whitespace-nowrap">{{ s.statement_date }}</td>
                <td class="py-3.5 px-4 font-mono font-bold text-slate-900 whitespace-nowrap">{{ s.ref_number }}</td>
                <td class="py-3.5 px-4 text-slate-700 max-w-sm">{{ s.description }}</td>
                <td class="py-3.5 px-4 text-center">
                  <span 
                    class="inline-block px-2 py-0.5 rounded text-[10px] font-bold uppercase"
                    :class="s.type === 'credit' ? 'bg-emerald-50 text-emerald-700' : 'bg-rose-50 text-rose-700'"
                  >
                    {{ s.type === 'credit' ? 'Kredit (Masuk)' : 'Debet (Keluar)' }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-right font-bold text-slate-900 whitespace-nowrap">
                  Rp {{ formatNum(s.amount) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <span 
                    class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-semibold uppercase"
                    :class="s.match_status === 'matched' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-amber-50 text-amber-700 border border-amber-200'"
                  >
                    {{ s.match_status === 'matched' ? 'Matched' : 'Unmatched' }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-center">
                  <button 
                    @click="toggleMatchReconciliation(s)"
                    class="px-2.5 py-1 rounded-lg text-[11px] font-semibold transition-colors cursor-pointer"
                    :class="s.match_status === 'matched' ? 'bg-slate-100 hover:bg-slate-200 text-slate-700' : 'bg-blue-600 hover:bg-blue-700 text-white'"
                  >
                    {{ s.match_status === 'matched' ? 'Batalkan Cocok' : 'Cocokkan Transaksi' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- MODALS SECTION                                                          -->
    <!-- ======================================================================= -->

    <!-- 1. Modal View Journal Lines -->
    <div v-if="showLinesModal" class="fixed inset-0 z-50 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-2xl w-full shadow-2xl border border-slate-200 overflow-hidden">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-bold text-slate-900 text-base">Rincian Baris Jurnal: {{ selectedJournalRef }}</h3>
          <button @click="showLinesModal = false" class="text-slate-400 hover:text-slate-600 cursor-pointer">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="p-5">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 font-semibold text-slate-600">
                <th class="py-2.5 px-3">Kode & Nama Akun</th>
                <th class="py-2.5 px-3">Keterangan Memo</th>
                <th class="py-2.5 px-3 text-right">Debit (Rp)</th>
                <th class="py-2.5 px-3 text-right">Credit (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="l in selectedJournalLines" :key="l.id">
                <td class="py-2.5 px-3 font-semibold text-slate-800">{{ l.account_code }} - {{ l.account_name }}</td>
                <td class="py-2.5 px-3 text-slate-500">{{ l.description || '-' }}</td>
                <td class="py-2.5 px-3 text-right font-medium text-emerald-600">
                  {{ l.debit > 0 ? formatNum(l.debit) : '-' }}
                </td>
                <td class="py-2.5 px-3 text-right font-medium text-blue-600">
                  {{ l.credit > 0 ? formatNum(l.credit) : '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="p-4 bg-slate-50 border-t border-slate-200 flex justify-end">
          <button @click="showLinesModal = false" class="px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white rounded-xl text-xs font-semibold cursor-pointer">
            Tutup
          </button>
        </div>
      </div>
    </div>

    <!-- 2. Modal Buat Jurnal Baru -->
    <div v-if="showNewJournalModal" class="fixed inset-0 z-50 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-3xl w-full shadow-2xl border border-slate-200 overflow-hidden max-h-[90vh] flex flex-col">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-bold text-slate-900 text-base">Buat Jurnal Umum Baru</h3>
          <button @click="showNewJournalModal = false" class="text-slate-400 hover:text-slate-600 cursor-pointer">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="p-6 space-y-4 overflow-y-auto flex-1 text-xs">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block font-semibold text-slate-700 mb-1">Tanggal Transaksi</label>
              <input v-model="newJournalForm.entry_date" type="date" class="w-full px-3 py-2 border border-slate-200 rounded-xl" />
            </div>
            <div>
              <label class="block font-semibold text-slate-700 mb-1">Keterangan / Memo Jurnal</label>
              <input v-model="newJournalForm.description" type="text" placeholder="e.g. Pembayaran listrik outlet" class="w-full px-3 py-2 border border-slate-200 rounded-xl" />
            </div>
          </div>

          <div>
            <div class="flex items-center justify-between mb-2">
              <span class="font-bold text-slate-800">Baris Akun (Debit & Credit)</span>
              <button @click="addJournalLine" class="text-xs text-blue-600 font-semibold hover:underline cursor-pointer">
                + Tambah Baris
              </button>
            </div>
            <div class="space-y-2">
              <div v-for="(line, idx) in newJournalForm.lines" :key="idx" class="flex items-center gap-2">
                <select v-model="line.account_id" class="flex-2 px-3 py-2 border border-slate-200 rounded-xl text-xs bg-white">
                  <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </select>
                <input v-model.number="line.debit" type="number" placeholder="Debit" class="flex-1 px-3 py-2 border border-slate-200 rounded-xl text-xs text-right font-medium" />
                <input v-model.number="line.credit" type="number" placeholder="Credit" class="flex-1 px-3 py-2 border border-slate-200 rounded-xl text-xs text-right font-medium" />
                <button @click="removeJournalLine(idx)" class="p-2 text-rose-500 hover:bg-rose-50 rounded-lg cursor-pointer">
                  <X class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>

          <!-- Total Balance Check -->
          <div class="bg-slate-50 p-4 rounded-xl flex items-center justify-between font-bold">
            <span>Total:</span>
            <div class="flex gap-6">
              <span class="text-emerald-600">Debit: Rp {{ formatNum(newJournalTotalDebit) }}</span>
              <span class="text-blue-600">Credit: Rp {{ formatNum(newJournalTotalCredit) }}</span>
              <span :class="isNewJournalBalanced ? 'text-emerald-700' : 'text-rose-600'">
                {{ isNewJournalBalanced ? '✓ Seimbang' : '✗ Tidak Seimbang' }}
              </span>
            </div>
          </div>
        </div>
        <div class="p-4 bg-slate-50 border-t border-slate-200 flex justify-end gap-3">
          <button @click="showNewJournalModal = false" class="px-4 py-2 border border-slate-200 text-slate-700 rounded-xl text-xs font-semibold cursor-pointer">
            Batal
          </button>
          <button 
            @click="submitNewJournal"
            :disabled="!isNewJournalBalanced || loading"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded-xl text-xs font-semibold shadow-xs cursor-pointer"
          >
            Posting Jurnal
          </button>
        </div>
      </div>
    </div>

    <!-- 3. Modal Tambah Akun COA -->
    <div v-if="showNewAccountModal" class="fixed inset-0 z-50 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-md w-full shadow-2xl border border-slate-200 overflow-hidden">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-bold text-slate-900 text-base">Tambah Bagan Akun (COA) Baru</h3>
          <button @click="showNewAccountModal = false" class="text-slate-400 hover:text-slate-600 cursor-pointer">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="p-5 space-y-3 text-xs">
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Kode Akun (e.g. 1103, 6301)</label>
            <input v-model="newAccountForm.code" type="text" placeholder="1103" class="w-full px-3 py-2 border border-slate-200 rounded-xl font-mono" />
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Nama Akun</label>
            <input v-model="newAccountForm.name" type="text" placeholder="Bank Mandiri Operasional" class="w-full px-3 py-2 border border-slate-200 rounded-xl" />
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Klasifikasi Akun</label>
            <select v-model="newAccountForm.account_type" class="w-full px-3 py-2 border border-slate-200 rounded-xl bg-white">
              <option value="asset">Asset (Aktiva / Harta)</option>
              <option value="liability">Liability (Kewajiban / Hutang)</option>
              <option value="equity">Equity (Modal / Ekuitas)</option>
              <option value="revenue">Revenue (Pendapatan)</option>
              <option value="expense">Expense (Beban Operasional / HPP)</option>
            </select>
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Deskripsi Tambahan</label>
            <textarea v-model="newAccountForm.description" rows="2" placeholder="Catatan kegunaan akun..." class="w-full px-3 py-2 border border-slate-200 rounded-xl"></textarea>
          </div>
        </div>
        <div class="p-4 bg-slate-50 border-t border-slate-200 flex justify-end gap-3">
          <button @click="showNewAccountModal = false" class="px-4 py-2 border border-slate-200 text-slate-700 rounded-xl text-xs font-semibold cursor-pointer">
            Batal
          </button>
          <button @click="submitNewAccount" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs cursor-pointer">
            Simpan Akun
          </button>
        </div>
      </div>
    </div>

    <!-- 4. Modal Ajukan Pengeluaran -->
    <div v-if="showNewExpenseModal" class="fixed inset-0 z-50 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-md w-full shadow-2xl border border-slate-200 overflow-hidden">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-bold text-slate-900 text-base">Ajukan Pengeluaran Baru</h3>
          <button @click="showNewExpenseModal = false" class="text-slate-400 hover:text-slate-600 cursor-pointer">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="p-5 space-y-3 text-xs">
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Kategori Pengeluaran</label>
            <select v-model="newExpenseForm.category" class="w-full px-3 py-2 border border-slate-200 rounded-xl bg-white">
              <option value="Operasional & Kas Kecil">Operasional & Kas Kecil</option>
              <option value="Bahan Baku">Bahan Baku</option>
              <option value="Pemeliharaan Peralatan">Pemeliharaan Peralatan</option>
              <option value="Promosi & Pemasaran">Promosi & Pemasaran</option>
              <option value="Utilitas Listrik, Air & Gas">Utilitas Listrik, Air & Gas</option>
            </select>
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Nominal (Rp)</label>
            <input v-model.number="newExpenseForm.amount" type="number" placeholder="500000" class="w-full px-3 py-2 border border-slate-200 rounded-xl font-bold" />
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Keterangan Pengeluaran</label>
            <textarea v-model="newExpenseForm.description" rows="2" placeholder="e.g. Pembelian bahan darurat di supermarket..." class="w-full px-3 py-2 border border-slate-200 rounded-xl"></textarea>
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Diajukan Oleh</label>
            <input v-model="newExpenseForm.submitted_by" type="text" placeholder="Sarah (Manager)" class="w-full px-3 py-2 border border-slate-200 rounded-xl" />
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">URL Bukti / Foto Nota Kuitansi</label>
            <input v-model="newExpenseForm.receipt_url" type="text" placeholder="https://..." class="w-full px-3 py-2 border border-slate-200 rounded-xl" />
          </div>
        </div>
        <div class="p-4 bg-slate-50 border-t border-slate-200 flex justify-end gap-3">
          <button @click="showNewExpenseModal = false" class="px-4 py-2 border border-slate-200 text-slate-700 rounded-xl text-xs font-semibold cursor-pointer">
            Batal
          </button>
          <button @click="submitNewExpense" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs cursor-pointer">
            Kirim Pengajuan
          </button>
        </div>
      </div>
    </div>

    <!-- 5. Modal Setor Pajak -->
    <div v-if="showTaxSettleModal" class="fixed inset-0 z-50 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-md w-full shadow-2xl border border-slate-200 overflow-hidden">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-bold text-slate-900 text-base">Konfirmasi Penyetoran Pajak</h3>
          <button @click="showTaxSettleModal = false" class="text-slate-400 hover:text-slate-600 cursor-pointer">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="p-5 space-y-3 text-xs">
          <div class="bg-slate-50 p-3 rounded-xl space-y-1">
            <p class="font-bold text-slate-900">{{ selectedTaxToSettle?.tax_type }}</p>
            <p class="text-slate-500">Jumlah Setor: <span class="font-bold text-slate-900">Rp {{ formatNum(selectedTaxToSettle?.tax_amount) }}</span></p>
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Nomor Transaksi Penerimaan Negara (NTPN)</label>
            <input v-model="taxNTPN" type="text" class="w-full px-3 py-2 border border-slate-200 rounded-xl font-mono font-bold text-blue-600" />
          </div>
        </div>
        <div class="p-4 bg-slate-50 border-t border-slate-200 flex justify-end gap-3">
          <button @click="showTaxSettleModal = false" class="px-4 py-2 border border-slate-200 text-slate-700 rounded-xl text-xs font-semibold cursor-pointer">
            Batal
          </button>
          <button @click="confirmTaxSettle" class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-semibold shadow-xs cursor-pointer">
            Konfirmasi Setor Pajak
          </button>
        </div>
      </div>
    </div>

    <!-- 6. Modal Set Budget -->
    <div v-if="showNewBudgetModal" class="fixed inset-0 z-50 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-md w-full shadow-2xl border border-slate-200 overflow-hidden">
        <div class="p-5 border-b border-slate-200 flex items-center justify-between">
          <h3 class="font-bold text-slate-900 text-base">Atur Anggaran Kategori</h3>
          <button @click="showNewBudgetModal = false" class="text-slate-400 hover:text-slate-600 cursor-pointer">
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="p-5 space-y-3 text-xs">
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Kategori Pengeluaran</label>
            <select v-model="newBudgetForm.category" class="w-full px-3 py-2 border border-slate-200 rounded-xl bg-white">
              <option value="Bahan Baku & Minuman">Bahan Baku & Minuman</option>
              <option value="Gaji & Upah Karyawan">Gaji & Upah Karyawan</option>
              <option value="Utilitas Listrik, Air & Gas">Utilitas Listrik, Air & Gas</option>
              <option value="Sewa Tempat & Gedung">Sewa Tempat & Gedung</option>
              <option value="Promosi & Pemasaran">Promosi & Pemasaran</option>
              <option value="Operasional & Kas Kecil">Operasional & Kas Kecil</option>
            </select>
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Nominal Anggaran Target (Rp)</label>
            <input v-model.number="newBudgetForm.budget_amount" type="number" placeholder="15000000" class="w-full px-3 py-2 border border-slate-200 rounded-xl font-bold" />
          </div>
          <div>
            <label class="block font-semibold text-slate-700 mb-1">Catatan Anggaran</label>
            <textarea v-model="newBudgetForm.notes" rows="2" placeholder="Catatan alokasi dana..." class="w-full px-3 py-2 border border-slate-200 rounded-xl"></textarea>
          </div>
        </div>
        <div class="p-4 bg-slate-50 border-t border-slate-200 flex justify-end gap-3">
          <button @click="showNewBudgetModal = false" class="px-4 py-2 border border-slate-200 text-slate-700 rounded-xl text-xs font-semibold cursor-pointer">
            Batal
          </button>
          <button @click="submitNewBudget" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold shadow-xs cursor-pointer">
            Simpan Anggaran
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
