<template>
  <Transition name="payment-fade">
    <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/40 backdrop-blur-md p-4 transition-all duration-300">
      <div class="bg-white rounded-3xl max-w-2xl w-full p-6 shadow-2xl shadow-slate-950/25 border border-slate-100/90 transform transition-all duration-300 ease-out">
      <div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-5">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-black text-slate-900 tracking-tight">{{ orderNumber }}</h3>
          <span class="px-2.5 py-0.5 rounded-lg text-xs font-bold bg-blue-50 text-blue-600 border border-blue-100 flex items-center gap-1.5">
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
            <div v-for="item in items" :key="item.name" class="flex items-center justify-between text-xs font-medium text-slate-700">
              <div class="flex items-center gap-2">
                <span class="font-bold text-slate-900">{{ item.name }}</span>
                <span class="text-slate-400">x{{ item.qty }}</span>
              </div>
              <span class="font-semibold text-slate-900">Rp {{ (item.price * item.qty).toLocaleString('id-ID') }}</span>
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

          <!-- Process button -->
          <button
            @click="processPayment"
            :disabled="processing"
            class="w-full py-3 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 text-white font-bold rounded-xl shadow-lg shadow-blue-600/30 transition-all cursor-pointer mt-2"
          >
            {{ processing ? 'Memproses...' : 'Bayar Sekarang' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</Transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import axios from 'axios'
import { Armchair, MapPin, X, QrCode } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const props = defineProps({
  show: Boolean,
  orderId: { type: String, default: '' },
  orderNumber: { type: String, default: '#ORD-1209' },
  table: { type: String, default: 'T-05' },
  orderType: { type: String, default: 'Dine-in' },
  items: { type: Array as () => Array<{ name: string, qty: number, price: number }>, default: () => [] },
  subtotal: { type: Number, default: 0 },
  tax: { type: Number, default: 0 },
  total: { type: Number, default: 0 }
})

const emit = defineEmits(['close', 'success'])

const notifyStore = useNotificationStore()
const method = ref('Tunai')
const paidAmount = ref(props.total || 0)
const processing = ref(false)

watch(() => props.total, (val) => {
  paidAmount.value = val || 0
})

const processPayment = async () => {
  if (method.value === 'Tunai' && paidAmount.value < props.total) {
    notifyStore.error('Nominal pembayaran tunai kurang dari total tagihan')
    return
  }

  processing.value = true
  try {
    if (props.orderId) {
      await axios.post(`/api/v1/pos/orders/${props.orderId}/pay`, {
        payment_method: method.value,
        amount_paid: paidAmount.value,
        total_amount: props.total
      })
    }
    notifyStore.success('Pembayaran berhasil diproses!')
    emit('success')
    emit('close')
  } catch (err: any) {
    console.error('Payment error:', err)
    notifyStore.error(err.response?.data?.error || 'Gagal memproses pembayaran')
  } finally {
    processing.value = false
  }
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
