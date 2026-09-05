<template>
  <Transition name="payment-fade">
    <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/40 backdrop-blur-md p-4 transition-all duration-300">
      <div class="bg-white rounded-3xl max-w-2xl w-full p-6 shadow-2xl shadow-slate-950/25 border border-slate-100/90 transform transition-all duration-300 ease-out overflow-hidden relative">

        <!-- ======================================================== -->
        <!-- STATE 1: FORM PEMBAYARAN -->
        <!-- ======================================================== -->
        <div v-if="paymentState === 'form'" class="space-y-5">
          <!-- Header -->
          <div class="flex items-center justify-between pb-4 border-b border-slate-100">
            <div class="flex items-center gap-2.5 flex-wrap">
              <h3 class="text-lg font-black text-slate-900 tracking-tight">{{ orderNumber }}</h3>
              <span v-if="queueNumber && queueNumber !== '-'" class="px-2.5 py-0.5 rounded-lg text-xs font-black bg-indigo-600 text-white shadow-xs">
                Antrian #{{ queueNumber }}
              </span>
              <span v-if="table && table !== '-'" class="px-2.5 py-0.5 rounded-lg text-xs font-bold bg-blue-50 text-blue-600 border border-blue-100 flex items-center gap-1.5">
                <Armchair class="w-3.5 h-3.5" /> Meja {{ table }}
              </span>
              <span class="px-2.5 py-0.5 rounded-lg text-xs font-bold bg-amber-50 text-amber-600 border border-amber-100 flex items-center gap-1.5">
                <MapPin class="w-3.5 h-3.5" /> {{ orderType }}
              </span>
            </div>
            <button @click="$emit('close')" class="text-slate-400 hover:text-slate-600 p-1 rounded-lg hover:bg-slate-100 transition-colors cursor-pointer">
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Left: Order Summary -->
            <div class="space-y-4 border-r border-slate-100 pr-0 md:pr-4">
              <div class="space-y-2.5 max-h-56 overflow-y-auto pr-1">
                <div v-for="item in displayItems" :key="item.name" class="flex items-center justify-between text-xs font-medium text-slate-700">
                  <div class="flex items-center gap-2">
                    <span class="font-bold text-slate-900">{{ item.name }}</span>
                    <span class="text-slate-400">x{{ item.qty }}</span>
                  </div>
                  <span class="font-semibold text-slate-900">Rp {{ ((item.price || 0) * (item.qty || 1)).toLocaleString('id-ID') }}</span>
                </div>
              </div>

              <div class="pt-4 border-t border-slate-100 space-y-2 text-xs">
                <div class="flex justify-between text-slate-500">
                  <span>Subtotal</span>
                  <span class="font-semibold text-slate-800">Rp {{ subtotal.toLocaleString('id-ID') }}</span>
                </div>
                <div class="flex justify-between text-slate-500">
                  <span>Diskon</span>
                  <span class="font-semibold text-slate-800">Rp 0</span>
                </div>
                <div class="flex justify-between text-slate-500">
                  <span>PPN 10%</span>
                  <span class="font-semibold text-slate-800">Rp {{ tax.toLocaleString('id-ID') }}</span>
                </div>
                <div class="flex justify-between text-base font-black text-slate-900 pt-2 border-t border-slate-100">
                  <span>Total</span>
                  <span class="text-blue-600">Rp {{ total.toLocaleString('id-ID') }}</span>
                </div>
              </div>
            </div>

            <!-- Right: Payment Methods & Cash Processing -->
            <div class="space-y-4">
              <!-- Method Tabs -->
              <div class="grid grid-cols-4 gap-1 bg-slate-100 p-1 rounded-xl">
                <button
                  v-for="m in ['Tunai', 'QRIS', 'Transfer', 'E-Wallet']"
                  :key="m"
                  @click="method = m"
                  class="py-1.5 text-xs font-bold rounded-lg transition-all cursor-pointer"
                  :class="method === m ? 'bg-white text-blue-600 shadow-xs' : 'text-slate-500 hover:text-slate-800'"
                >
                  {{ m }}
                </button>
              </div>

              <!-- Total Display Box -->
              <div class="bg-slate-50 p-4 rounded-xl border border-slate-100 text-center">
                <div class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Total yang harus dibayar:</div>
                <div class="text-2xl font-black text-slate-900 tracking-tight mt-1">
                  Rp {{ total.toLocaleString('id-ID') }}
                </div>
              </div>

              <!-- Cash input -->
              <div v-if="method === 'Tunai'" class="space-y-2">
                <label class="block text-xs font-semibold text-slate-700">Jumlah Dibayar</label>
                <input
                  v-model.number="paidAmount"
                  type="number"
                  class="w-full px-3 py-2 text-sm font-bold text-slate-900 bg-white border border-blue-400 rounded-xl outline-none focus:ring-2 focus:ring-blue-500/20"
                />

                <!-- Kembalian Highlight -->
                <div class="p-2.5 bg-emerald-50 border border-emerald-200 rounded-xl flex items-center justify-between text-xs font-bold text-emerald-800">
                  <span>Kembalian:</span>
                  <span class="text-sm font-extrabold text-emerald-600">
                    Rp {{ Math.max(0, paidAmount - total).toLocaleString('id-ID') }}
                  </span>
                </div>

                <!-- Quick Cash Pills -->
                <div class="grid grid-cols-4 gap-1.5 pt-1">
                  <button
                    type="button"
                    v-for="pill in [50000, 100000, 200000, total]"
                    :key="pill"
                    @click="paidAmount = pill"
                    class="py-1.5 px-2 text-[11px] font-semibold rounded-lg border border-slate-200 bg-white hover:bg-slate-50 text-slate-700 transition-colors cursor-pointer"
                  >
                    {{ pill === total ? 'Uang Pas' : `Rp ${(pill / 1000)}k` }}
                  </button>
                </div>
              </div>

              <!-- QRIS display -->
              <div v-else-if="method === 'QRIS'" class="p-4 bg-slate-50 rounded-xl border border-slate-200 text-center space-y-2">
                <QrCode class="w-8 h-8 text-blue-600 mx-auto" />
                <div class="text-xs font-bold text-slate-800">Scan QRIS Dinamis</div>
                <div class="w-32 h-32 mx-auto bg-white p-2 rounded-lg border border-slate-300 flex items-center justify-center font-mono text-[10px] text-slate-400">
                  [QRIS CODE HERE]
                </div>
                <p class="text-[10px] text-slate-500">Mendukung BCA, Mandiri, GoPay, OVO, ShopeePay</p>
              </div>

              <!-- Transfer / E-Wallet fallback -->
              <div v-else class="p-4 bg-slate-50 rounded-xl border border-slate-200 text-center text-xs text-slate-600 space-y-1">
                <div class="font-bold text-slate-800">Konfirmasi Pembayaran {{ method }}</div>
                <p class="text-[11px] text-slate-500">Pastikan mutasi dana atau transfer pelanggan telah berhasil masuk.</p>
              </div>

              <!-- Process button -->
              <button
                @click="processPayment"
                class="w-full py-3 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-bold rounded-xl shadow-lg shadow-blue-600/30 transition-all cursor-pointer mt-2 flex items-center justify-center gap-2"
              >
                <CreditCard class="w-4 h-4" />
                <span>Bayar Sekarang (Rp {{ total.toLocaleString('id-ID') }})</span>
              </button>
            </div>
          </div>
        </div>

        <!-- ======================================================== -->
        <!-- STATE 2: LOADING ANIMASI PEMROSESAN -->
        <!-- ======================================================== -->
        <div v-else-if="paymentState === 'loading'" class="py-16 flex flex-col items-center justify-center text-center space-y-4">
          <div class="relative w-20 h-20 flex items-center justify-center">
            <div class="absolute inset-0 rounded-full border-4 border-blue-100 animate-ping opacity-40"></div>
            <div class="w-16 h-16 rounded-full border-4 border-blue-600 border-t-transparent animate-spin"></div>
            <CreditCard class="w-7 h-7 text-blue-600 absolute" />
          </div>
          <div>
            <h4 class="text-base font-black text-slate-900">Memproses Transaksi Pembayaran...</h4>
            <p class="text-xs text-slate-500 mt-1">Menghubungkan ke sistem kasir & mencatat transaksi</p>
          </div>
        </div>

        <!-- ======================================================== -->
        <!-- STATE 3: PAYMENT DONE & CETAK STRUK OTOMATIS -->
        <!-- ======================================================== -->
        <div v-else-if="paymentState === 'done'" class="py-8 px-4 flex flex-col items-center justify-center text-center space-y-5 animate-in fade-in zoom-in duration-300">
          <!-- Big Success Badge -->
          <div class="relative">
            <div class="w-20 h-20 rounded-full bg-emerald-50 border-4 border-emerald-500 flex items-center justify-center shadow-lg shadow-emerald-500/20">
              <CheckCircle2 class="w-10 h-10 text-emerald-600" />
            </div>
            <span class="absolute -top-1 -right-1 px-2.5 py-0.5 rounded-full text-[10px] font-black bg-emerald-600 text-white shadow-xs">
              LUNAS
            </span>
          </div>

          <div>
            <h3 class="text-xl font-black text-slate-900 tracking-tight">Pembayaran Sukses!</h3>
            <p class="text-xs text-slate-500 mt-1">Transaksi telah selesai dicatat dan status meja telah tersedia kembali</p>
          </div>

          <!-- Transaction Summary Card -->
          <div class="w-full max-w-md bg-slate-50 rounded-2xl p-4 border border-slate-200/80 space-y-2.5 text-xs text-left">
            <div class="flex justify-between items-center pb-2 border-b border-slate-200">
              <span class="text-slate-500 font-medium">No. Order</span>
              <span class="font-mono font-bold text-slate-900">{{ completedReceipt?.order_number || orderNumber }}</span>
            </div>
            <div v-if="completedReceipt?.queue_number && completedReceipt.queue_number !== '-'" class="flex justify-between items-center pb-2 border-b border-slate-200">
              <span class="text-slate-500 font-medium">Nomor Antrian</span>
              <span class="px-2.5 py-0.5 rounded-lg bg-indigo-600 text-white font-black text-xs">
                Antrian #{{ completedReceipt.queue_number }}
              </span>
            </div>
            <div class="flex justify-between items-center pb-2 border-b border-slate-200">
              <span class="text-slate-500 font-medium">Total Tagihan</span>
              <span class="font-black text-blue-600 text-sm">Rp {{ total.toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between items-center pb-2 border-b border-slate-200">
              <span class="text-slate-500 font-medium">Metode & Bayar</span>
              <span class="font-bold text-slate-800 uppercase">{{ method }} (Rp {{ (completedReceipt?.amount_paid || paidAmount).toLocaleString('id-ID') }})</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-slate-500 font-medium">Kembalian</span>
              <span class="font-black text-emerald-600 text-sm">
                Rp {{ ((completedReceipt?.change_due !== undefined ? completedReceipt.change_due : Math.max(0, paidAmount - total))).toLocaleString('id-ID') }}
              </span>
            </div>
          </div>

          <!-- Auto Print Status Notice -->
          <div class="flex items-center gap-2 text-xs text-emerald-700 bg-emerald-50 px-4 py-2 rounded-xl border border-emerald-200">
            <Printer class="w-4 h-4 text-emerald-600" />
            <span>Struk billing otomatis dikirim ke printer kasir (Blob thermal receipt)</span>
          </div>

          <!-- Action Buttons: Print Again & Close -->
          <div class="flex items-center gap-3 w-full max-w-md pt-2">
            <button
              @click="reprintReceipt"
              class="flex-1 py-3 bg-slate-100 hover:bg-slate-200 active:bg-slate-300 text-slate-800 font-bold text-xs rounded-xl transition-all flex items-center justify-center gap-2 cursor-pointer border border-slate-200"
            >
              <Printer class="w-4 h-4 text-slate-600" />
              <span>Cetak Ulang Struk</span>
            </button>
            <button
              @click="finishAndClose"
              class="flex-1 py-3 bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white font-bold text-xs rounded-xl shadow-lg shadow-emerald-600/30 transition-all flex items-center justify-center gap-2 cursor-pointer"
            >
              <Check class="w-4 h-4" />
              <span>Selesai & Tutup</span>
            </button>
          </div>
        </div>

      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import axios from 'axios'
import {
  Armchair,
  MapPin,
  X,
  QrCode,
  CheckCircle2,
  Printer,
  CreditCard,
  Check
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'
import { printReceiptBlob, type ReceiptData } from '@/utils/receiptPrinter'

const props = defineProps({
  show: Boolean,
  orderId: { type: String, default: '' },
  orderNumber: { type: String, default: '#ORD-1209' },
  queueNumber: { type: String, default: '' },
  table: { type: String, default: 'T-05' },
  orderType: { type: String, default: 'Dine-in' },
  items: { type: Array as () => Array<{ name: string, qty: number, price: number }>, default: () => [] },
  rawItems: { type: Array as () => Array<{ id: string, name: string, price: number, qty: number }>, default: () => [] },
  subtotal: { type: Number, default: 0 },
  tax: { type: Number, default: 0 },
  total: { type: Number, default: 0 }
})

const emit = defineEmits(['close', 'success'])

const notifyStore = useNotificationStore()
const paymentState = ref<'form' | 'loading' | 'done'>('form')
const method = ref('Tunai')
const paidAmount = ref(props.total || 0)
const completedReceipt = ref<ReceiptData | null>(null)

watch(() => props.total, (val) => {
  paidAmount.value = val || 0
})

watch(() => props.show, (newVal) => {
  if (newVal) {
    paymentState.value = 'form'
    paidAmount.value = props.total || 0
    completedReceipt.value = null
  }
})

const displayItems = computed(() => {
  if (props.items && props.items.length > 0) return props.items
  if (props.rawItems && props.rawItems.length > 0) {
    return props.rawItems.map(i => ({ name: i.name, qty: i.qty, price: i.price }))
  }
  return []
})

const processPayment = async () => {
  if (method.value === 'Tunai' && paidAmount.value < props.total) {
    notifyStore.error('Nominal pembayaran tunai kurang dari total tagihan')
    return
  }

  paymentState.value = 'loading'
  try {
    let activeId = props.orderId
    let generatedOrderNum = props.orderNumber
    let generatedQueueNum = props.queueNumber

    // If order has not been created yet (direct payment from cart in POS)
    if (!activeId && props.rawItems && props.rawItems.length > 0) {
      const orderRes = await axios.post('/api/v1/pos/orders', {
        customer_name: props.orderType === 'Dine-in' ? `Tamu ${props.table}` : 'Pelanggan Kasir',
        table_number: props.orderType === 'Dine-in' ? props.table : '',
        order_type: props.orderType.toLowerCase().replace('-', '_'),
        items: props.rawItems.map(i => ({
          product_id: i.id,
          quantity: i.qty,
          unit_price: i.price
        }))
      })
      activeId = orderRes.data?.order_id
      generatedOrderNum = orderRes.data?.order_number || props.orderNumber
      generatedQueueNum = orderRes.data?.queue_number || props.queueNumber
    }

    let payRes: any = null
    if (activeId) {
      payRes = await axios.post(`/api/v1/pos/orders/${activeId}/pay`, {
        payment_method: method.value,
        amount_paid: paidAmount.value,
        total_amount: props.total
      })
    }

    // Build complete receipt data for printing
    const receiptData: ReceiptData = {
      order_number: payRes?.data?.order_number || generatedOrderNum,
      queue_number: payRes?.data?.queue_number || generatedQueueNum,
      table_number: props.table,
      customer_name: payRes?.data?.customer_name || (props.orderType === 'Dine-in' ? `Tamu ${props.table}` : 'Pelanggan'),
      order_type: props.orderType,
      paid_at: payRes?.data?.paid_at,
      cashier_name: 'Jane Cashier',
      items: displayItems.value.map(i => ({
        name: i.name,
        quantity: i.qty,
        unit_price: i.price,
        total_price: i.price * i.qty
      })),
      subtotal: props.subtotal,
      tax: props.tax,
      total: props.total,
      payment_method: method.value,
      amount_paid: paidAmount.value,
      change_due: payRes?.data?.change_due ?? Math.max(0, paidAmount.value - props.total)
    }

    completedReceipt.value = receiptData

    // Transition to 'done' state
    paymentState.value = 'done'
    notifyStore.success('Pembayaran sukses! Status meja diperbarui.', 'Lunas')

    // Automatically trigger receipt print via hidden iframe Blob URL
    setTimeout(() => {
      printReceiptBlob(receiptData)
    }, 400)

  } catch (err: any) {
    console.error('Payment error:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal memproses pembayaran')
    paymentState.value = 'form'
  }
}

const reprintReceipt = () => {
  if (completedReceipt.value) {
    printReceiptBlob(completedReceipt.value)
    notifyStore.info('Mencetak ulang struk pembayaran...', 'Cetak Struk')
  }
}

const finishAndClose = () => {
  emit('success')
  emit('close')
}
</script>

<style scoped>
.payment-fade-enter-active,
.payment-fade-leave-active {
  transition: opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.payment-fade-enter-active > div,
.payment-fade-leave-active > div {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.payment-fade-enter-from,
.payment-fade-leave-to {
  opacity: 0;
}

.payment-fade-enter-from > div,
.payment-fade-leave-to > div {
  opacity: 0;
  transform: scale(0.94) translateY(8px);
}
</style>

