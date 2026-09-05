<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-xs flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2">
          <div class="p-2 rounded-xl bg-blue-50 text-blue-600">
            <ReceiptText class="w-5 h-5" />
          </div>
          <h1 class="text-xl font-black text-slate-900 tracking-tight">Riwayat Transaksi & Struk Kasir</h1>
        </div>
        <p class="text-xs text-slate-500 mt-1">
          Daftar seluruh transaksi yang telah berhasil dibayar (Completed) dengan fitur cetak ulang billing struk kasir.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="fetchTransactions"
          :disabled="loading"
          class="flex items-center gap-2 px-4 py-2.5 bg-slate-100 hover:bg-slate-200 active:bg-slate-300 text-slate-700 font-bold text-xs rounded-xl transition-all cursor-pointer disabled:opacity-50"
        >
          <RotateCw class="w-3.5 h-3.5" :class="loading ? 'animate-spin' : ''" />
          <span>Segarkan Data</span>
        </button>
        <router-link
          to="/pos"
          class="flex items-center gap-2 px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold text-xs rounded-xl shadow-md shadow-blue-600/25 transition-all cursor-pointer"
        >
          <Utensils class="w-3.5 h-3.5" />
          <span>Kembali ke POS</span>
        </router-link>
      </div>
    </div>

    <!-- Revenue & Summary Metrics Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Total Omzet Lunas</span>
          <div class="text-xl font-black text-slate-900 mt-1">
            Rp {{ totalRevenue.toLocaleString('id-ID') }}
          </div>
          <span class="text-[10px] text-emerald-600 font-bold flex items-center gap-1 mt-1">
            <CheckCircle2 class="w-3 h-3 inline" /> {{ transactions.length }} Transaksi Berhasil
          </span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
          <Banknote class="w-6 h-6" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Rata-rata Transaksi</span>
          <div class="text-xl font-black text-slate-900 mt-1">
            Rp {{ averageOrderValue.toLocaleString('id-ID') }}
          </div>
          <span class="text-[10px] text-slate-500 font-medium mt-1">Nilai per struk pembayaran</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-blue-50 text-blue-600 flex items-center justify-center">
          <TrendingUp class="w-6 h-6" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Metode Tunai (Cash)</span>
          <div class="text-xl font-black text-slate-900 mt-1">
            Rp {{ cashRevenue.toLocaleString('id-ID') }}
          </div>
          <span class="text-[10px] text-slate-500 font-medium mt-1">{{ cashCount }} transaksi tunai</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-amber-50 text-amber-600 flex items-center justify-center">
          <Wallet class="w-6 h-6" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs flex items-center justify-between">
        <div>
          <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Non-Tunai (QRIS/Transfer)</span>
          <div class="text-xl font-black text-slate-900 mt-1">
            Rp {{ nonCashRevenue.toLocaleString('id-ID') }}
          </div>
          <span class="text-[10px] text-slate-500 font-medium mt-1">{{ nonCashCount }} transaksi non-tunai</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-indigo-50 text-indigo-600 flex items-center justify-center">
          <CreditCard class="w-6 h-6" />
        </div>
      </div>
    </div>

    <!-- Filters & Search Toolbar -->
    <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-xs space-y-3">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <!-- Live Search -->
        <div class="relative w-full md:w-80">
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari order #, antrian, tamu, meja..."
            class="w-full pl-9 pr-3 py-2 bg-slate-50 border border-slate-200 rounded-xl text-xs font-medium focus:bg-white focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none transition-all"
          />
        </div>

        <!-- Filter Payment Method -->
        <div class="flex items-center gap-1.5 overflow-x-auto text-xs pb-1">
          <span class="text-slate-400 text-[11px] font-bold mr-1">Metode:</span>
          <button
            v-for="m in ['Semua', 'Tunai', 'QRIS', 'Transfer', 'E-Wallet']"
            :key="m"
            @click="selectedMethod = m"
            class="px-3 py-1.5 rounded-lg font-bold transition-all cursor-pointer whitespace-nowrap"
            :class="selectedMethod === m ? 'bg-blue-600 text-white shadow-xs' : 'bg-slate-50 text-slate-600 hover:bg-slate-100'"
          >
            {{ m }}
          </button>
        </div>

        <!-- Filter Order Type -->
        <div class="flex items-center gap-1.5 overflow-x-auto text-xs pb-1">
          <span class="text-slate-400 text-[11px] font-bold mr-1">Tipe:</span>
          <button
            v-for="t in ['Semua', 'Dine-in', 'Takeaway', 'Delivery']"
            :key="t"
            @click="selectedOrderType = t"
            class="px-3 py-1.5 rounded-lg font-bold transition-all cursor-pointer whitespace-nowrap"
            :class="selectedOrderType === t ? 'bg-indigo-600 text-white shadow-xs' : 'bg-slate-50 text-slate-600 hover:bg-slate-100'"
          >
            {{ t }}
          </button>
        </div>
      </div>
    </div>

    <!-- Table of Completed Transactions -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <!-- Loading State -->
      <div v-if="loading" class="py-24 flex flex-col items-center justify-center text-center space-y-3">
        <div class="w-10 h-10 border-4 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
        <p class="text-xs font-semibold text-slate-500">Mengambil data transaksi selesai dari server...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredTransactions.length === 0" class="py-24 flex flex-col items-center justify-center text-center space-y-3">
        <div class="w-16 h-16 rounded-full bg-slate-50 border border-slate-200 flex items-center justify-center text-slate-400">
          <ReceiptText class="w-8 h-8" />
        </div>
        <div>
          <h4 class="text-sm font-bold text-slate-800">Tidak ada riwayat transaksi ditemukan</h4>
          <p class="text-xs text-slate-400 mt-1">Belum ada transaksi pembayaran yang cocok dengan filter pencarian Anda.</p>
        </div>
      </div>

      <!-- Transactions Table -->
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 text-[11px] uppercase font-bold text-slate-400 border-b border-slate-200 tracking-wider">
            <tr>
              <th class="py-3.5 px-4">No. Order & Antrian</th>
              <th class="py-3.5 px-4">Waktu Selesai</th>
              <th class="py-3.5 px-4">Meja / Tipe</th>
              <th class="py-3.5 px-4">Pelanggan</th>
              <th class="py-3.5 px-4">Ringkasan Menu</th>
              <th class="py-3.5 px-4">Metode Bayar</th>
              <th class="py-3.5 px-4 text-right">Total Tagihan</th>
              <th class="py-3.5 px-4 text-center">Aksi Struk</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 font-medium">
            <tr
              v-for="tx in filteredTransactions"
              :key="tx.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <!-- Order & Queue -->
              <td class="py-3.5 px-4">
                <div class="flex items-center gap-2 flex-wrap">
                  <span class="font-mono font-black text-slate-900">{{ tx.order_number }}</span>
                  <span
                    v-if="tx.queue_number && tx.queue_number !== '-'"
                    class="px-2 py-0.5 rounded-md text-[10px] font-black bg-indigo-600 text-white shadow-xs"
                  >
                    #{{ tx.queue_number }}
                  </span>
                </div>
                <span class="text-[10px] text-slate-400 block mt-0.5">ID: {{ tx.id.substring(0, 8) }}</span>
              </td>

              <!-- Timestamp -->
              <td class="py-3.5 px-4 whitespace-nowrap text-slate-500">
                <div class="flex items-center gap-1.5">
                  <Clock class="w-3.5 h-3.5 text-slate-400" />
                  <span>{{ formatTimestamp(tx.paid_at || tx.created_at) }}</span>
                </div>
              </td>

              <!-- Table / Type -->
              <td class="py-3.5 px-4 whitespace-nowrap">
                <div v-if="tx.order_type === 'dine_in' || tx.table_number !== '-'" class="flex items-center gap-1.5">
                  <span class="px-2.5 py-1 rounded-lg text-xs font-bold bg-blue-50 text-blue-700 border border-blue-100 flex items-center gap-1">
                    <Armchair class="w-3 h-3" />
                    Meja {{ tx.table_number || '-' }}
                  </span>
                </div>
                <div v-else-if="tx.order_type === 'takeaway'" class="flex items-center gap-1.5">
                  <span class="px-2.5 py-1 rounded-lg text-xs font-bold bg-amber-50 text-amber-700 border border-amber-200/80 flex items-center gap-1">
                    <ShoppingBag class="w-3 h-3" />
                    Takeaway
                  </span>
                </div>
                <div v-else class="flex items-center gap-1.5">
                  <span class="px-2.5 py-1 rounded-lg text-xs font-bold bg-purple-50 text-purple-700 border border-purple-200/80 flex items-center gap-1">
                    <MapPin class="w-3 h-3" />
                    Delivery
                  </span>
                </div>
              </td>

              <!-- Customer -->
              <td class="py-3.5 px-4">
                <div class="font-bold text-slate-900">{{ tx.customer_name || 'Pelanggan' }}</div>
                <span class="text-[10px] text-slate-400">Kasir: Jane Cashier</span>
              </td>

              <!-- Items Summary -->
              <td class="py-3.5 px-4">
                <div class="max-w-xs space-y-1">
                  <div
                    v-for="item in (tx.items || []).slice(0, 2)"
                    :key="item.id"
                    class="text-xs text-slate-700 truncate"
                  >
                    <span class="font-bold text-slate-900">{{ item.quantity }}x</span> {{ item.name }}
                  </div>
                  <div v-if="(tx.items || []).length > 2" class="text-[10px] text-blue-600 font-bold">
                    + {{ (tx.items || []).length - 2 }} item lainnya
                  </div>
                </div>
              </td>

              <!-- Payment Method -->
              <td class="py-3.5 px-4 whitespace-nowrap">
                <div class="flex items-center gap-1.5">
                  <span class="px-2.5 py-0.5 rounded-md text-[10px] font-black uppercase tracking-wider bg-slate-100 text-slate-800 border border-slate-200">
                    {{ tx.payment_method || 'TUNAI' }}
                  </span>
                </div>
                <div v-if="tx.change_due > 0" class="text-[10px] text-emerald-600 font-semibold mt-0.5">
                  Kembali: Rp {{ Number(tx.change_due).toLocaleString('id-ID') }}
                </div>
              </td>

              <!-- Total Amount -->
              <td class="py-3.5 px-4 text-right whitespace-nowrap">
                <div class="text-sm font-black text-slate-900">
                  Rp {{ Number(tx.total_amount || 0).toLocaleString('id-ID') }}
                </div>
                <span class="text-[10px] font-bold text-emerald-600">LUNAS</span>
              </td>

              <!-- Action: Print Billing Receipt -->
              <td class="py-3.5 px-4 text-center whitespace-nowrap">
                <div class="flex items-center justify-center gap-1.5">
                  <button
                    @click="handlePrintReceipt(tx)"
                    class="p-2 rounded-xl bg-blue-50 hover:bg-blue-600 text-blue-600 hover:text-white transition-all duration-150 cursor-pointer shadow-xs group"
                    title="Cetak Struk Billing (Blob Thermal)"
                  >
                    <Printer class="w-4 h-4 group-hover:scale-110 transition-transform" />
                  </button>
                  <button
                    @click="previewReceipt(tx)"
                    class="p-2 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-600 transition-colors cursor-pointer"
                    title="Lihat Detail Struk"
                  >
                    <Eye class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Preview Struk Billing -->
    <Transition name="fade">
      <div
        v-if="selectedPreviewReceipt"
        class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/40 backdrop-blur-sm p-4 transition-all"
      >
        <div class="bg-white rounded-3xl max-w-sm w-full p-6 shadow-2xl border border-slate-100 flex flex-col max-h-[90vh]">
          <!-- Modal Header -->
          <div class="flex items-center justify-between pb-3 border-b border-slate-100 mb-4">
            <div class="flex items-center gap-2">
              <Receipt class="w-4 h-4 text-blue-600" />
              <h3 class="text-sm font-black text-slate-900">Preview Struk Kasir</h3>
            </div>
            <button
              @click="selectedPreviewReceipt = null"
              class="text-slate-400 hover:text-slate-600 p-1 rounded-lg hover:bg-slate-100 transition-colors cursor-pointer"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <!-- Thermal Receipt Body Simulation -->
          <div class="flex-1 overflow-y-auto pr-1 bg-slate-50 p-4 rounded-2xl border border-slate-200 font-mono text-xs text-slate-800 space-y-3">
            <div class="text-center space-y-0.5">
              <div class="font-black text-sm tracking-wider">CAFE ERP SYSTEM</div>
              <div class="text-[10px] text-slate-500">Kopi Premium & Kitchen</div>
              <div class="text-[9px] text-slate-400">Jl. Senopati Raya No. 45, Jakarta Selatan</div>
              <div
                v-if="selectedPreviewReceipt.queue_number && selectedPreviewReceipt.queue_number !== '-'"
                class="inline-block mt-2 px-3 py-1 rounded-md bg-slate-900 text-white text-xs font-black"
              >
                NO. ANTRIAN: #{{ selectedPreviewReceipt.queue_number }}
              </div>
            </div>

            <div class="border-t border-b border-dashed border-slate-300 py-2 text-[10px] space-y-0.5">
              <div class="flex justify-between">
                <span>No. Order:</span>
                <span class="font-bold">{{ selectedPreviewReceipt.order_number }}</span>
              </div>
              <div class="flex justify-between">
                <span>Meja / Tipe:</span>
                <span>{{ selectedPreviewReceipt.order_type === 'dine_in' ? `Meja ${selectedPreviewReceipt.table_number}` : selectedPreviewReceipt.order_type }}</span>
              </div>
              <div class="flex justify-between">
                <span>Tamu:</span>
                <span>{{ selectedPreviewReceipt.customer_name }}</span>
              </div>
              <div class="flex justify-between">
                <span>Waktu:</span>
                <span>{{ selectedPreviewReceipt.paid_at || selectedPreviewReceipt.created_at }}</span>
              </div>
            </div>

            <!-- Items -->
            <div class="space-y-1.5 py-1">
              <div
                v-for="item in selectedPreviewReceipt.items"
                :key="item.id"
                class="flex justify-between items-start text-[11px]"
              >
                <div>
                  <div class="font-bold">{{ item.name }}</div>
                  <div class="text-[9px] text-slate-500">{{ item.quantity }} x Rp {{ Number(item.unit_price).toLocaleString('id-ID') }}</div>
                </div>
                <div class="font-bold">Rp {{ Number(item.total_price || item.unit_price * item.quantity).toLocaleString('id-ID') }}</div>
              </div>
            </div>

            <!-- Totals -->
            <div class="border-t border-slate-300 pt-2 text-[10px] space-y-1">
              <div class="flex justify-between">
                <span>Subtotal</span>
                <span>Rp {{ Number(selectedPreviewReceipt.subtotal || 0).toLocaleString('id-ID') }}</span>
              </div>
              <div class="flex justify-between">
                <span>PPN 10%</span>
                <span>Rp {{ Number(selectedPreviewReceipt.tax_amount || 0).toLocaleString('id-ID') }}</span>
              </div>
              <div class="flex justify-between text-xs font-black border-t border-b border-slate-900 py-1.5">
                <span>TOTAL</span>
                <span>Rp {{ Number(selectedPreviewReceipt.total_amount || 0).toLocaleString('id-ID') }}</span>
              </div>
              <div class="flex justify-between pt-1">
                <span>Metode:</span>
                <span class="font-bold uppercase">{{ selectedPreviewReceipt.payment_method }}</span>
              </div>
              <div class="flex justify-between">
                <span>Bayar:</span>
                <span>Rp {{ Number(selectedPreviewReceipt.amount_paid || selectedPreviewReceipt.total_amount).toLocaleString('id-ID') }}</span>
              </div>
              <div class="flex justify-between font-bold text-emerald-600">
                <span>Kembalian:</span>
                <span>Rp {{ Number(selectedPreviewReceipt.change_due || 0).toLocaleString('id-ID') }}</span>
              </div>
            </div>

            <div class="text-center text-[9px] text-slate-400 pt-2 border-t border-dashed border-slate-300">
              <div>TERIMA KASIH ATAS KUNJUNGAN ANDA</div>
              <div>Simpan struk ini sebagai bukti transaksi yang sah.</div>
            </div>
          </div>

          <!-- Actions -->
          <div class="pt-4 flex items-center gap-2">
            <button
              @click="selectedPreviewReceipt = null"
              class="flex-1 py-2.5 rounded-xl border border-slate-200 text-xs font-bold text-slate-600 hover:bg-slate-50 transition-colors cursor-pointer"
            >
              Tutup
            </button>
            <button
              @click="handlePrintReceipt(selectedPreviewReceipt)"
              class="flex-1 py-2.5 rounded-xl bg-blue-600 hover:bg-blue-700 text-xs font-bold text-white shadow-md shadow-blue-600/30 transition-all flex items-center justify-center gap-1.5 cursor-pointer"
            >
              <Printer class="w-3.5 h-3.5" />
              <span>Cetak Struk</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import {
  ReceiptText,
  RotateCw,
  Utensils,
  Banknote,
  TrendingUp,
  Wallet,
  CreditCard,
  Search,
  Clock,
  Armchair,
  ShoppingBag,
  MapPin,
  Printer,
  Eye,
  X,
  Receipt,
  CheckCircle2
} from 'lucide-vue-next'
import { printReceiptBlob, type ReceiptData } from '@/utils/receiptPrinter'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const loading = ref(false)
const transactions = ref<any[]>([])
const searchQuery = ref('')
const selectedMethod = ref('Semua')
const selectedOrderType = ref('Semua')
const selectedPreviewReceipt = ref<any | null>(null)

const fetchTransactions = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/pos/transactions')
    transactions.value = res.data?.data || []
  } catch (err: any) {
    console.error('Error fetching completed transactions:', err)
    notifyStore.error('Gagal mengambil data riwayat transaksi')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchTransactions()
})

const filteredTransactions = computed(() => {
  return transactions.value.filter(tx => {
    // Search query matching
    const q = searchQuery.value.toLowerCase().trim()
    const matchesQuery = !q ||
      tx.order_number?.toLowerCase().includes(q) ||
      tx.queue_number?.toLowerCase().includes(q) ||
      tx.customer_name?.toLowerCase().includes(q) ||
      tx.table_number?.toLowerCase().includes(q)

    // Method matching
    const matchesMethod = selectedMethod.value === 'Semua' ||
      tx.payment_method?.toLowerCase() === selectedMethod.value.toLowerCase()

    // Order type matching
    let txType = 'dine-in'
    if (tx.order_type === 'takeaway') txType = 'takeaway'
    if (tx.order_type === 'delivery') txType = 'delivery'
    const matchesType = selectedOrderType.value === 'Semua' ||
      txType === selectedOrderType.value.toLowerCase()

    return matchesQuery && matchesMethod && matchesType
  })
})

// Metrics Computations
const totalRevenue = computed(() => {
  return transactions.value.reduce((sum, tx) => sum + (Number(tx.total_amount) || 0), 0)
})

const averageOrderValue = computed(() => {
  if (transactions.value.length === 0) return 0
  return Math.round(totalRevenue.value / transactions.value.length)
})

const cashRevenue = computed(() => {
  return transactions.value
    .filter(tx => (tx.payment_method || '').toLowerCase() === 'tunai' || (tx.payment_method || '').toLowerCase() === 'cash')
    .reduce((sum, tx) => sum + (Number(tx.total_amount) || 0), 0)
})

const cashCount = computed(() => {
  return transactions.value
    .filter(tx => (tx.payment_method || '').toLowerCase() === 'tunai' || (tx.payment_method || '').toLowerCase() === 'cash')
    .length
})

const nonCashRevenue = computed(() => {
  return totalRevenue.value - cashRevenue.value
})

const nonCashCount = computed(() => {
  return transactions.value.length - cashCount.value
})

const formatTimestamp = (dateStr?: string) => {
  if (!dateStr) return '-'
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleString('id-ID', {
      day: '2-digit',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

const buildReceiptData = (tx: any): ReceiptData => {
  return {
    order_number: tx.order_number,
    queue_number: tx.queue_number,
    table_number: tx.table_number,
    customer_name: tx.customer_name,
    order_type: tx.order_type === 'dine_in' ? 'Dine-in' : (tx.order_type === 'takeaway' ? 'Takeaway' : 'Delivery'),
    paid_at: tx.paid_at || tx.created_at,
    cashier_name: 'Jane Cashier',
    items: (tx.items || []).map((i: any) => ({
      name: i.name,
      quantity: i.quantity,
      unit_price: i.unit_price,
      total_price: i.total_price || (i.unit_price * i.quantity)
    })),
    subtotal: tx.subtotal || 0,
    tax: tx.tax_amount || 0,
    total: tx.total_amount || 0,
    payment_method: tx.payment_method || 'Tunai',
    amount_paid: tx.amount_paid || tx.total_amount || 0,
    change_due: tx.change_due || 0
  }
}

const handlePrintReceipt = async (tx: any) => {
  const receiptData = buildReceiptData(tx)
  notifyStore.info(`Mencetak struk untuk order ${tx.order_number}...`, 'Cetak Struk')
  await printReceiptBlob(receiptData)
}

const previewReceipt = (tx: any) => {
  selectedPreviewReceipt.value = tx
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

