<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-xs p-4">
    <div class="bg-white rounded-2xl max-w-2xl w-full p-6 shadow-2xl border border-slate-100">
      <div class="flex items-center justify-between pb-4 border-b border-slate-100 mb-4">
        <div>
          <h3 class="text-lg font-black text-slate-900">Meja {{ table?.code }}</h3>
          <p class="text-xs text-slate-500 font-medium">{{ table?.zone || 'Lantai 1' }} • Kapasitas {{ table?.capacity || 4 }} Kursi</p>
        </div>
        <button @click="$emit('close')" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-100 cursor-pointer">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Occupied state -->
      <div v-if="table?.status === 'occupied' || table?.status === 'billing'" class="space-y-4">
        <div class="p-4 bg-slate-50 rounded-xl border border-slate-200/80 flex items-center justify-between">
          <div>
            <div class="text-xs font-bold text-slate-500 uppercase">Pesanan Aktif</div>
            <div class="text-base font-black text-slate-900 font-mono">{{ orderDetail?.order_number || table?.orderNumber || 'ORD-ACTIVE' }}</div>
            <div class="text-xs text-slate-600 font-medium mt-0.5">Tamu: <strong>{{ orderDetail?.customer || table?.customer || 'Pelanggan' }}</strong></div>
          </div>
          <div class="text-right">
            <span class="px-2.5 py-1 rounded-full text-xs font-extrabold uppercase bg-red-100 text-red-700">
              {{ table?.status }}
            </span>
          </div>
        </div>

        <div>
          <div class="text-xs font-bold text-slate-700 uppercase mb-2">Item Menu Dipesan</div>
          <div v-if="!orderDetail?.items || orderDetail.items.length === 0" class="py-6 text-center text-slate-400 text-xs border border-dashed rounded-xl">
            Belum ada item pesanan
          </div>
          <div v-else class="space-y-2 max-h-48 overflow-y-auto pr-1">
            <div
              v-for="item in orderDetail.items"
              :key="item.id"
              class="flex items-center justify-between p-2.5 rounded-xl bg-white border border-slate-150 text-xs"
            >
              <div>
                <span class="font-bold text-slate-900">{{ item.name }}</span>
                <span class="text-slate-400 ml-2">x{{ item.quantity }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="font-bold text-slate-800">Rp {{ Number(item.total_price || item.unit_price * item.quantity).toLocaleString('id-ID') }}</span>
                <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-amber-100 text-amber-800">{{ item.kitchen_status || 'Pending' }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-100 flex items-center justify-between">
          <div>
            <span class="text-xs text-slate-400 font-semibold">Total Tagihan:</span>
            <div class="text-lg font-black text-blue-600">Rp {{ Number(orderDetail?.total || table?.total || 0).toLocaleString('id-ID') }}</div>
          </div>
          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="goToBilling"
              class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-xs font-bold cursor-pointer"
            >
              Ke Kasir Pembayaran
            </button>
            <button
              type="button"
              @click="showPayment = true"
              class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold cursor-pointer"
            >
              Bayar Meja Ini
            </button>
          </div>
        </div>
      </div>

      <!-- Available state -->
      <div v-else class="space-y-4">
        <div class="p-4 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-800 text-xs">
          <strong>Meja Tersedia:</strong> Siap untuk menerima pesanan baru.
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-700 mb-1">Nama Tamu</label>
          <input
            v-model="customerName"
            type="text"
            placeholder="e.g. Tamu Meja 1"
            class="w-full px-3.5 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none focus:bg-white focus:border-blue-600 font-bold"
          />
        </div>
        <div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
          <button
            type="button"
            @click="setReserved"
            class="px-4 py-2 rounded-xl border border-amber-300 bg-amber-50 text-amber-800 text-xs font-bold cursor-pointer"
          >
            Set Reservasi
          </button>
          <button
            type="button"
            @click="createOrder"
            :disabled="creating"
            class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold cursor-pointer disabled:opacity-50"
          >
            {{ creating ? 'Membuka...' : 'Buka Pesanan Dine-in' }}
          </button>
        </div>
      </div>

      <PaymentModal
        v-if="showPayment && orderDetail"
        :show="showPayment"
        :order-id="orderDetail.id"
        :order-number="orderDetail.order_number"
        :table="table?.code"
        :order-type="'Dine-in'"
        :items="orderDetail.items?.map((i: any) => ({ name: i.name, qty: i.quantity, price: i.unit_price }))"
        :subtotal="orderDetail.subtotal"
        :tax="orderDetail.tax"
        :total="orderDetail.total"
        @close="showPayment = false"
        @success="handleSuccess"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { X } from 'lucide-vue-next'
import PaymentModal from './PaymentModal.vue'
import { useNotificationStore } from '@/stores/notification.store'

const props = defineProps({
  show: Boolean,
  table: { type: Object, default: () => ({}) }
})

const emit = defineEmits(['close', 'updated'])

const router = useRouter()
const notifyStore = useNotificationStore()

const orderDetail = ref<any>(null)
const customerName = ref('')
const creating = ref(false)
const showPayment = ref(false)

const loadOrderDetail = async () => {
  if (!props.table?.code && !props.table?.id) return
  try {
    const tableId = props.table.code || props.table.id
    const res = await axios.get(`/api/v1/pos/tables/${tableId}/order`)
    if (res.data?.data) {
      orderDetail.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load table order detail:', err)
  }
}

watch(() => props.table, () => {
  if (props.table && (props.table.status === 'occupied' || props.table.status === 'billing')) {
    loadOrderDetail()
  }
  customerName.value = `Tamu ${props.table?.code || ''}`
}, { immediate: true })

onMounted(() => {
  if (props.table && (props.table.status === 'occupied' || props.table.status === 'billing')) {
    loadOrderDetail()
  }
})

const createOrder = async () => {
  creating.value = true
  try {
    await axios.post('/api/v1/pos/orders', {
      customer_name: customerName.value || `Tamu ${props.table?.code}`,
      table_number: props.table?.code,
      order_type: 'dine_in',
      items: []
    })
    notifyStore.success(`Pesanan Meja ${props.table?.code} dibuka!`, 'Pesanan Baru')
    emit('updated')
    await loadOrderDetail()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.error || 'Gagal membuka pesanan', 'Kesalahan')
  } finally {
    creating.value = false
  }
}

const setReserved = async () => {
  try {
    const tableId = props.table?.code || props.table?.id
    await axios.put(`/api/v1/pos/tables/${tableId}/status`, { status: 'reserved' })
    notifyStore.success('Meja diset ke reservasi', 'Status Meja')
    emit('updated')
    emit('close')
  } catch (err: any) {
    notifyStore.error('Gagal mengubah status', 'Kesalahan')
  }
}

const goToBilling = () => {
  emit('close')
  router.push(`/pos?tab=billing&table=${props.table?.code}`)
}

const handleSuccess = () => {
  showPayment.value = false
  emit('updated')
  emit('close')
}
</script>

