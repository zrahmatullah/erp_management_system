<template>
  <div class="h-[calc(100vh-5.5rem)] flex flex-col space-y-4">
    <!-- Top Header Bar -->
    <div class="bg-slate-900 text-white rounded-2xl p-4 flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-lg shrink-0 border border-slate-800">
      <div class="flex items-center gap-4">
        <div class="w-11 h-11 rounded-xl bg-blue-600 flex items-center justify-center font-bold shadow-md shadow-blue-600/30">
          <ChefHat class="w-6 h-6 text-white" />
        </div>
        <div>
          <h2 class="text-lg font-black tracking-tight flex items-center gap-2">
            Kitchen Display System (KDS)
            <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
          </h2>
          <p class="text-xs text-slate-400">Antrian Pesanan Dapur & Barista Real-time</p>
        </div>
      </div>

      <!-- Station Filter Pills & Clock -->
      <div class="flex items-center gap-3 flex-wrap">
        <div class="flex items-center gap-1.5 bg-slate-800 p-1 rounded-xl border border-slate-700 text-xs">
          <button
            v-for="st in ['Semua Station', 'barista', 'kitchen']"
            :key="st"
            @click="stationFilter = st"
            class="px-3 py-1.5 rounded-lg font-bold transition-all cursor-pointer capitalize"
            :class="stationFilter === st ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-400 hover:text-white'"
          >
            {{ st === 'Semua Station' ? 'Semua Station' : (st === 'barista' ? '☕ Barista' : '🍳 Kitchen') }}
          </button>
        </div>

        <button
          @click="fetchKDSTickets"
          class="p-2.5 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-xl border border-slate-700 transition-colors cursor-pointer"
          title="Refresh Antrian Dapur"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>

        <!-- Real-time Clock -->
        <div class="flex items-center gap-2 font-mono text-lg font-black bg-slate-800 px-4 py-1.5 rounded-xl border border-slate-700 text-emerald-400 shadow-inner">
          <Clock class="w-4 h-4 text-emerald-400" />
          <span>{{ currentTime }}</span>
        </div>
      </div>
    </div>

    <!-- Sub-tab Filter: Dine-In vs Takeaway vs Delivery vs All -->
    <div class="bg-slate-900/80 backdrop-blur-md rounded-2xl p-2 border border-slate-800 flex items-center justify-between gap-2 overflow-x-auto shrink-0">
      <div class="flex items-center gap-2">
        <button
          v-for="tab in orderTypeTabs"
          :key="tab.id"
          @click="activeTypeFilter = tab.id"
          class="flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-black transition-all cursor-pointer whitespace-nowrap"
          :class="activeTypeFilter === tab.id 
            ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' 
            : 'text-slate-400 hover:text-white hover:bg-slate-800/80'"
        >
          <component :is="tab.icon" class="w-4 h-4" />
          <span>{{ tab.label }}</span>
          <span
            class="px-2 py-0.5 rounded-full text-[10px] font-extrabold"
            :class="activeTypeFilter === tab.id ? 'bg-white text-blue-700' : 'bg-slate-800 text-slate-300'"
          >
            {{ getCountByType(tab.id) }}
          </span>
        </button>
      </div>

      <!-- Quick Legend Badge -->
      <div class="hidden lg:flex items-center gap-3 pr-3 text-[11px] font-bold text-slate-400">
        <span class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-blue-500"></span> Dine-in (#D)
        </span>
        <span class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span> Takeaway (#TA)
        </span>
        <span class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-purple-500"></span> Delivery (#DL)
        </span>
      </div>
    </div>

    <!-- 4 Kanban Columns: Pesanan Masuk, Sedang Dimasak, Siap Saji, Disajikan -->
    <div class="flex-1 grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 min-h-0 overflow-hidden">
      <!-- 1. Pesanan Masuk (Pending) -->
      <div class="flex flex-col h-full bg-slate-900/70 rounded-2xl border border-slate-800 overflow-hidden">
        <div class="px-4 py-3 bg-amber-500 text-slate-950 font-black text-xs uppercase tracking-wider flex items-center justify-between shrink-0 shadow-sm">
          <div class="flex items-center gap-2">
            <Clock class="w-4 h-4" />
            <span>Pesanan Masuk (Pending)</span>
          </div>
          <span class="px-2 py-0.5 rounded-full bg-black/20 text-xs font-black">{{ pendingTickets.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in pendingTickets"
            :key="t.id"
            class="bg-slate-900 p-4 rounded-xl border border-amber-500/40 shadow-md space-y-3 hover:border-amber-500 transition-all"
          >
            <!-- Header: Queue Badge & Timer -->
            <div class="flex items-start justify-between gap-2">
              <div>
                <span
                  class="px-2.5 py-1 rounded-lg text-xs font-black shadow-xs inline-flex items-center gap-1"
                  :class="getQueueBadgeClass(t.orderType)"
                >
                  Antrian #{{ t.queueNumber || '-' }}
                </span>
                <div class="text-[11px] font-bold text-slate-300 mt-1 flex items-center gap-1.5">
                  <span v-if="t.orderType === 'dine_in'">Meja {{ t.table }}</span>
                  <span v-else-if="t.orderType === 'takeaway'">Takeaway (Bungkus)</span>
                  <span v-else>Delivery Kurir</span>
                  <span class="text-slate-600">•</span>
                  <span class="text-slate-400 font-normal truncate max-w-[100px]">{{ t.customerName }}</span>
                </div>
              </div>

              <!-- Waiting Time -->
              <span
                class="font-mono text-xs font-black px-2 py-0.5 rounded-md"
                :class="getTimerClass(t.elapsedMin)"
              >
                ⏱ {{ t.elapsedMin }}m
              </span>
            </div>

            <!-- Item Detail Box -->
            <div class="bg-slate-950/60 p-3 rounded-lg border border-slate-800 space-y-1.5">
              <div class="flex items-start justify-between gap-2">
                <div class="flex items-center gap-2">
                  <span class="w-6 h-6 rounded-md bg-amber-500/20 text-amber-400 font-black text-xs flex items-center justify-center">
                    {{ t.quantity }}x
                  </span>
                  <span class="font-bold text-white text-xs">{{ t.productName }}</span>
                </div>
                <span class="px-2 py-0.5 rounded text-[9px] font-bold uppercase bg-slate-800 text-slate-400">
                  {{ t.station }}
                </span>
              </div>

              <div v-if="t.notes" class="text-[11px] text-amber-300 italic bg-amber-950/30 p-1.5 rounded border border-amber-900/50">
                Catatan: "{{ t.notes }}"
              </div>
            </div>

            <div class="flex items-center justify-between text-[10px] text-slate-500">
              <span>{{ t.orderNumber }}</span>
              <span>Masuk: {{ t.time }}</span>
            </div>

            <button
              @click="advanceTicket(t, 'cooking')"
              class="w-full py-2 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white font-black text-xs rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center justify-center gap-1.5"
            >
              <span>Mulai Masak</span>
              <span>→</span>
            </button>
          </div>

          <div v-if="pendingTickets.length === 0" class="h-40 flex flex-col items-center justify-center text-slate-600 text-xs space-y-1">
            <CheckCircle2 class="w-8 h-8 text-slate-700" />
            <span>Tidak ada antrian pesanan masuk</span>
          </div>
        </div>
      </div>

      <!-- 2. Sedang Dimasak (Cooking) -->
      <div class="flex flex-col h-full bg-slate-900/70 rounded-2xl border border-slate-800 overflow-hidden">
        <div class="px-4 py-3 bg-blue-600 text-white font-black text-xs uppercase tracking-wider flex items-center justify-between shrink-0 shadow-sm">
          <div class="flex items-center gap-2">
            <Flame class="w-4 h-4 text-orange-300" />
            <span>Sedang Dimasak (Cooking)</span>
          </div>
          <span class="px-2 py-0.5 rounded-full bg-black/20 text-xs font-black">{{ inProgressTickets.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in inProgressTickets"
            :key="t.id"
            class="bg-slate-900 p-4 rounded-xl border border-blue-500/50 shadow-md space-y-3 hover:border-blue-400 transition-all"
          >
            <!-- Header -->
            <div class="flex items-start justify-between gap-2">
              <div>
                <span
                  class="px-2.5 py-1 rounded-lg text-xs font-black shadow-xs inline-flex items-center gap-1"
                  :class="getQueueBadgeClass(t.orderType)"
                >
                  Antrian #{{ t.queueNumber || '-' }}
                </span>
                <div class="text-[11px] font-bold text-slate-300 mt-1 flex items-center gap-1.5">
                  <span v-if="t.orderType === 'dine_in'">Meja {{ t.table }}</span>
                  <span v-else-if="t.orderType === 'takeaway'">Takeaway</span>
                  <span v-else>Delivery</span>
                  <span class="text-slate-600">•</span>
                  <span class="text-slate-400 font-normal truncate max-w-[100px]">{{ t.customerName }}</span>
                </div>
              </div>

              <!-- Waiting Time -->
              <span
                class="font-mono text-xs font-black px-2 py-0.5 rounded-md"
                :class="getTimerClass(t.elapsedMin)"
              >
                ⏱ {{ t.elapsedMin }}m
              </span>
            </div>

            <!-- Item Detail Box -->
            <div class="bg-slate-950/60 p-3 rounded-lg border border-slate-800 space-y-1.5">
              <div class="flex items-start justify-between gap-2">
                <div class="flex items-center gap-2">
                  <span class="w-6 h-6 rounded-md bg-blue-500/20 text-blue-400 font-black text-xs flex items-center justify-center">
                    {{ t.quantity }}x
                  </span>
                  <span class="font-bold text-white text-xs">{{ t.productName }}</span>
                </div>
                <span class="px-2 py-0.5 rounded text-[9px] font-bold uppercase bg-slate-800 text-slate-400">
                  {{ t.station }}
                </span>
              </div>

              <div v-if="t.notes" class="text-[11px] text-amber-300 italic bg-amber-950/30 p-1.5 rounded border border-amber-900/50">
                Catatan: "{{ t.notes }}"
              </div>
            </div>

            <div class="flex items-center justify-between text-[10px] text-slate-500">
              <span>{{ t.orderNumber }}</span>
              <span>Masuk: {{ t.time }}</span>
            </div>

            <button
              @click="advanceTicket(t, 'ready')"
              class="w-full py-2 bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white font-black text-xs rounded-xl shadow-md shadow-emerald-600/30 transition-all cursor-pointer flex items-center justify-center gap-1.5"
            >
              <span>Selesai Masak</span>
              <span>→</span>
            </button>
          </div>

          <div v-if="inProgressTickets.length === 0" class="h-40 flex flex-col items-center justify-center text-slate-600 text-xs space-y-1">
            <Flame class="w-8 h-8 text-slate-700" />
            <span>Dapur sedang santai, tidak ada masakan</span>
          </div>
        </div>
      </div>

      <!-- 3. Siap Saji (Ready) -->
      <div class="flex flex-col h-full bg-slate-900/70 rounded-2xl border border-slate-800 overflow-hidden">
        <div class="px-4 py-3 bg-emerald-600 text-white font-black text-xs uppercase tracking-wider flex items-center justify-between shrink-0 shadow-sm">
          <div class="flex items-center gap-2">
            <BellRing class="w-4 h-4 text-emerald-200" />
            <span>Siap Saji / Pick-Up</span>
          </div>
          <span class="px-2 py-0.5 rounded-full bg-black/20 text-xs font-black">{{ readyTickets.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in readyTickets"
            :key="t.id"
            class="bg-slate-900 p-4 rounded-xl border border-emerald-500/50 shadow-md space-y-3 hover:border-emerald-400 transition-all"
          >
            <!-- Header -->
            <div class="flex items-start justify-between gap-2">
              <div>
                <span
                  class="px-2.5 py-1 rounded-lg text-xs font-black shadow-xs inline-flex items-center gap-1"
                  :class="getQueueBadgeClass(t.orderType)"
                >
                  Antrian #{{ t.queueNumber || '-' }}
                </span>
                <div class="text-[11px] font-bold text-slate-300 mt-1 flex items-center gap-1.5">
                  <span v-if="t.orderType === 'dine_in'">Meja {{ t.table }}</span>
                  <span v-else-if="t.orderType === 'takeaway'">Takeaway</span>
                  <span v-else>Delivery</span>
                  <span class="text-slate-600">•</span>
                  <span class="text-slate-400 font-normal truncate max-w-[100px]">{{ t.customerName }}</span>
                </div>
              </div>

              <span class="px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-400 text-[10px] font-extrabold flex items-center gap-1">
                <Check class="w-3 h-3" /> Ready
              </span>
            </div>

            <!-- Item Detail Box -->
            <div class="bg-slate-950/60 p-3 rounded-lg border border-slate-800 space-y-1.5">
              <div class="flex items-start justify-between gap-2">
                <div class="flex items-center gap-2">
                  <span class="w-6 h-6 rounded-md bg-emerald-500/20 text-emerald-400 font-black text-xs flex items-center justify-center">
                    {{ t.quantity }}x
                  </span>
                  <span class="font-bold text-white text-xs">{{ t.productName }}</span>
                </div>
                <span class="px-2 py-0.5 rounded text-[9px] font-bold uppercase bg-slate-800 text-slate-400">
                  {{ t.station }}
                </span>
              </div>
            </div>

            <button
              @click="advanceTicket(t, 'served')"
              class="w-full py-2 bg-slate-800 hover:bg-slate-700 text-emerald-400 border border-emerald-500/30 font-black text-xs rounded-xl transition-all cursor-pointer flex items-center justify-center gap-1.5"
            >
              <span>{{ t.orderType === 'dine_in' ? 'Antar ke Meja' : 'Serahkan ke Pelanggan' }}</span>
              <span>✓</span>
            </button>
          </div>

          <div v-if="readyTickets.length === 0" class="h-40 flex flex-col items-center justify-center text-slate-600 text-xs space-y-1">
            <BellRing class="w-8 h-8 text-slate-700" />
            <span>Belum ada sajian yang siap di-pickup</span>
          </div>
        </div>
      </div>

      <!-- 4. Disajikan (Served) -->
      <div class="flex flex-col h-full bg-slate-900/70 rounded-2xl border border-slate-800 overflow-hidden">
        <div class="px-4 py-3 bg-slate-800 text-slate-300 font-black text-xs uppercase tracking-wider flex items-center justify-between shrink-0 shadow-sm">
          <div class="flex items-center gap-2">
            <CheckCircle2 class="w-4 h-4 text-slate-400" />
            <span>Selesai (Served)</span>
          </div>
          <span class="px-2 py-0.5 rounded-full bg-black/20 text-xs font-black">{{ servedTickets.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in servedTickets"
            :key="t.id"
            class="bg-slate-900/50 p-4 rounded-xl border border-slate-800 space-y-2.5 opacity-70"
          >
            <div class="flex items-center justify-between">
              <span class="font-black text-xs text-slate-300">#{{ t.queueNumber }} • {{ t.productName }}</span>
              <span class="text-[10px] text-slate-500 font-mono">{{ t.time }}</span>
            </div>
            <div class="text-[11px] text-slate-400 flex items-center justify-between">
              <span>{{ t.orderType === 'dine_in' ? `Meja ${t.table}` : 'Takeaway/Delivery' }}</span>
              <span class="text-emerald-400 font-bold flex items-center gap-1">
                <Check class="w-3 h-3" /> Disajikan
              </span>
            </div>
          </div>

          <div v-if="servedTickets.length === 0" class="h-40 flex flex-col items-center justify-center text-slate-600 text-xs space-y-1">
            <span>Belum ada riwayat sajian baru</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import {
  ChefHat,
  Clock,
  Check,
  RotateCw,
  Flame,
  BellRing,
  CheckCircle2,
  Armchair,
  ShoppingBag,
  Truck,
  Layers
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const currentTime = ref(new Date().toLocaleTimeString('id-ID'))
const loading = ref(false)
const stationFilter = ref('Semua Station')
const activeTypeFilter = ref('all')

let clockInterval: any = null
let autoRefreshInterval: any = null

interface Ticket {
  id: string
  orderId: string
  orderNumber: string
  queueNumber: string
  orderType: string
  customerName: string
  table: string
  productName: string
  quantity: number
  status: 'pending' | 'cooking' | 'ready' | 'served'
  station: string
  notes: string
  time: string
  elapsedMin: number
}

const tickets = ref<Ticket[]>([])

const orderTypeTabs = [
  { id: 'all', label: 'Semua Tiket', icon: Layers },
  { id: 'dine_in', label: 'Dine-In (Meja)', icon: Armchair },
  { id: 'takeaway', label: 'Takeaway (Bungkus)', icon: ShoppingBag },
  { id: 'delivery', label: 'Delivery (Kurir)', icon: Truck }
]

const fetchKDSTickets = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/kds/tickets')
    if (res.data?.data) {
      tickets.value = res.data.data.map((t: any) => ({
        id: t.id,
        orderId: t.order_id,
        orderNumber: t.order_number,
        queueNumber: t.queue_number || '-',
        orderType: (t.order_type || 'dine_in').toLowerCase(),
        customerName: t.customer_name || 'Guest',
        table: t.table_number || '-',
        productName: t.product_name,
        quantity: t.quantity || 1,
        status: t.status || 'pending',
        station: t.station || 'barista',
        notes: t.notes || '',
        time: t.time || '12:00',
        elapsedMin: t.elapsed_min || 0
      }))
    }
  } catch (err: any) {
    console.error('Failed to load KDS tickets:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchKDSTickets()
  clockInterval = setInterval(() => {
    currentTime.value = new Date().toLocaleTimeString('id-ID')
  }, 1000)

  // Auto refresh tickets every 12 seconds
  autoRefreshInterval = setInterval(() => {
    fetchKDSTickets()
  }, 12000)
})

onUnmounted(() => {
  if (clockInterval) clearInterval(clockInterval)
  if (autoRefreshInterval) clearInterval(autoRefreshInterval)
})

// Filter by Station and Order Type
const filteredTickets = computed(() => {
  return tickets.value.filter(t => {
    const matchType = activeTypeFilter.value === 'all' || t.orderType === activeTypeFilter.value
    const matchStation = stationFilter.value === 'Semua Station' || t.station.toLowerCase() === stationFilter.value.toLowerCase()
    return matchType && matchStation
  })
})

const pendingTickets = computed(() => filteredTickets.value.filter(t => t.status === 'pending'))
const inProgressTickets = computed(() => filteredTickets.value.filter(t => t.status === 'cooking'))
const readyTickets = computed(() => filteredTickets.value.filter(t => t.status === 'ready'))
const servedTickets = computed(() => filteredTickets.value.filter(t => t.status === 'served'))

const getCountByType = (type: string) => {
  if (type === 'all') return tickets.value.filter(t => t.status !== 'served').length
  return tickets.value.filter(t => t.orderType === type && t.status !== 'served').length
}

const getQueueBadgeClass = (orderType: string) => {
  switch (orderType) {
    case 'dine_in':
      return 'bg-blue-600 text-white'
    case 'takeaway':
      return 'bg-amber-500 text-white'
    case 'delivery':
      return 'bg-purple-600 text-white'
    default:
      return 'bg-slate-700 text-white'
  }
}

const getTimerClass = (minutes: number) => {
  if (minutes < 10) return 'bg-emerald-950/60 text-emerald-400 border border-emerald-800/40'
  if (minutes < 20) return 'bg-amber-950/60 text-amber-400 border border-amber-800/40'
  return 'bg-red-950/80 text-red-400 border border-red-800/60 animate-pulse'
}

const advanceTicket = async (ticket: Ticket, nextStatus: any) => {
  ticket.status = nextStatus
  try {
    await axios.put(`/api/v1/kds/items/${ticket.id}/status`, { status: nextStatus })
    notifyStore.success(
      `Pesanan ${ticket.productName} (Antrian #${ticket.queueNumber}) dipindahkan ke '${nextStatus}'!`,
      'KDS Status Diperbarui'
    )
  } catch (err) {
    console.warn('Could not persist status to backend:', err)
  }
}
</script>
