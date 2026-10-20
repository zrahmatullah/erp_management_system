<template>
  <div class="kds-v2 h-[calc(100vh-5.5rem)] flex flex-col space-y-4 animate-in fade-in duration-200">
    <!-- Top Header Bar -->
    <div 
      class="p-4 rounded-2xl border flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-xs shrink-0"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      <div class="flex items-center gap-3.5">
        <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-amber-500 to-orange-500 flex items-center justify-center text-white shadow-xs">
          <ChefHat class="w-5 h-5" />
        </div>
        <div>
          <h2 class="text-base font-bold tracking-tight text-slate-900 dark:text-slate-50 flex items-center gap-2">
            Kitchen Display System (KDS)
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          </h2>
          <p class="text-xs text-slate-500 dark:text-slate-400">Antrean Pesanan Barista & Dapur Real-time</p>
        </div>
      </div>

      <!-- Station Filters, Refresh, Clock -->
      <div class="flex items-center gap-3 flex-wrap">
        <div class="flex items-center gap-1 p-1 rounded-xl border text-xs" style="border-color: var(--border-color); background-color: var(--bg-content);">
          <button
            v-for="st in ['Semua Station', 'barista', 'kitchen']"
            :key="st"
            type="button"
            @click="stationFilter = st"
            class="px-3 py-1.5 rounded-lg font-semibold transition-all cursor-pointer capitalize"
            :class="stationFilter === st 
              ? 'bg-blue-600 text-white shadow-xs' 
              : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100'"
          >
            {{ st === 'Semua Station' ? 'Semua' : (st === 'barista' ? '☕ Barista' : '🍳 Kitchen') }}
          </button>
        </div>

        <button
          type="button"
          @click="fetchKDSTickets"
          class="p-2 rounded-xl border transition-colors cursor-pointer text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800"
          style="border-color: var(--border-color); background-color: var(--bg-content);"
          title="Perbarui Tiket"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>

        <div class="flex items-center gap-2 font-mono text-xs font-bold px-3 py-1.5 rounded-xl border text-emerald-600 dark:text-emerald-400"
             style="border-color: var(--border-color); background-color: var(--bg-content);">
          <Clock class="w-3.5 h-3.5" />
          <span>{{ currentTime }}</span>
        </div>
      </div>
    </div>

    <!-- Sub-tab Filter: Dine-In vs Takeaway vs Delivery -->
    <div 
      class="p-2 rounded-2xl border flex items-center justify-between gap-2 overflow-x-auto shrink-0 shadow-xs"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      <div class="flex items-center gap-1.5">
        <button
          v-for="tab in orderTypeTabs"
          :key="tab.id"
          type="button"
          @click="activeTypeFilter = tab.id"
          class="flex items-center gap-2 px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all cursor-pointer whitespace-nowrap border"
          :class="activeTypeFilter === tab.id 
            ? 'bg-blue-600 text-white border-blue-600 shadow-xs' 
            : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 border-transparent hover:bg-slate-100 dark:hover:bg-slate-800'"
        >
          <component :is="tab.icon" class="w-3.5 h-3.5" />
          <span>{{ tab.label }}</span>
          <span 
            class="px-1.5 py-0.2 rounded-full text-[10px] font-bold"
            :class="activeTypeFilter === tab.id ? 'bg-white text-blue-700' : 'bg-slate-200 dark:bg-slate-800 text-slate-700 dark:text-slate-300'"
          >
            {{ getCountByType(tab.id) }}
          </span>
        </button>
      </div>

      <div class="hidden lg:flex items-center gap-4 pr-3 text-[11px] text-slate-400">
        <span class="flex items-center gap-1.5"><span class="w-2 h-2 rounded-full bg-blue-500"></span> Dine-in</span>
        <span class="flex items-center gap-1.5"><span class="w-2 h-2 rounded-full bg-amber-500"></span> Takeaway</span>
      </div>
    </div>

    <!-- 4 Kanban Columns: Pesanan Masuk, Sedang Dimasak, Siap Saji, Disajikan -->
    <div class="flex-1 grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 min-h-0 overflow-hidden">
      <!-- 1. Pesanan Masuk -->
      <div class="flex flex-col h-full rounded-2xl border overflow-hidden"
           style="background-color: var(--bg-primary); border-color: var(--border-color);">
        <div class="px-4 py-2.5 bg-amber-500 text-slate-950 font-bold text-xs uppercase tracking-wider flex items-center justify-between shrink-0">
          <div class="flex items-center gap-2">
            <Clock class="w-3.5 h-3.5" />
            <span>1. Pesanan Masuk</span>
          </div>
          <span class="px-2 py-0.2 rounded-full bg-black/20 text-xs font-mono font-bold">{{ pendingTickets.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in pendingTickets"
            :key="t.id"
            class="p-3.5 rounded-xl border space-y-2 text-xs"
            style="background-color: var(--bg-content); border-color: var(--border-color);"
          >
            <div class="flex items-center justify-between font-bold">
              <span class="text-blue-600 dark:text-blue-400 font-mono">#{{ t.queueNumber }}</span>
              <span class="text-slate-400 font-normal">{{ t.time }}</span>
            </div>
            <div class="font-bold text-sm text-slate-900 dark:text-slate-100">
              {{ t.quantity }}x {{ t.productName }}
            </div>
            <div class="text-[11px] text-slate-500">
              Meja: <strong>{{ t.table }}</strong> • {{ t.customerName }}
            </div>
            <div v-if="t.notes" class="text-[11px] italic text-amber-600 dark:text-amber-400">
              "{{ t.notes }}"
            </div>
            <button
              type="button"
              @click="updateTicketStatus(t.id, 'cooking')"
              class="w-full py-1.5 rounded-lg text-xs font-semibold bg-amber-500 hover:bg-amber-600 text-slate-950 cursor-pointer transition-colors shadow-xs"
            >
              Mulai Proses
            </button>
          </div>
        </div>
      </div>

      <!-- 2. Sedang Dimasak -->
      <div class="flex flex-col h-full rounded-2xl border overflow-hidden"
           style="background-color: var(--bg-primary); border-color: var(--border-color);">
        <div class="px-4 py-2.5 bg-blue-600 text-white font-bold text-xs uppercase tracking-wider flex items-center justify-between shrink-0">
          <div class="flex items-center gap-2">
            <ChefHat class="w-3.5 h-3.5" />
            <span>2. Sedang Diproses</span>
          </div>
          <span class="px-2 py-0.2 rounded-full bg-white/20 text-xs font-mono font-bold">{{ cookingTickets.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in cookingTickets"
            :key="t.id"
            class="p-3.5 rounded-xl border space-y-2 text-xs"
            style="background-color: var(--bg-content); border-color: var(--border-color);"
          >
            <div class="flex items-center justify-between font-bold">
              <span class="text-blue-600 dark:text-blue-400 font-mono">#{{ t.queueNumber }}</span>
              <span class="text-slate-400 font-normal">{{ t.time }}</span>
            </div>
            <div class="font-bold text-sm text-slate-900 dark:text-slate-100">
              {{ t.quantity }}x {{ t.productName }}
            </div>
            <div class="text-[11px] text-slate-500">
              Meja: <strong>{{ t.table }}</strong> • {{ t.customerName }}
            </div>
            <button
              type="button"
              @click="updateTicketStatus(t.id, 'ready')"
              class="w-full py-1.5 rounded-lg text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white cursor-pointer transition-colors shadow-xs"
            >
              Tandai Siap Saji
            </button>
          </div>
        </div>
      </div>

      <!-- 3. Siap Saji -->
      <div class="flex flex-col h-full rounded-2xl border overflow-hidden"
           style="background-color: var(--bg-primary); border-color: var(--border-color);">
        <div class="px-4 py-2.5 bg-emerald-600 text-white font-bold text-xs uppercase tracking-wider flex items-center justify-between shrink-0">
          <div class="flex items-center gap-2">
            <CheckCircle2 class="w-3.5 h-3.5" />
            <span>3. Siap Disajikan</span>
          </div>
          <span class="px-2 py-0.2 rounded-full bg-white/20 text-xs font-mono font-bold">{{ readyTickets.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in readyTickets"
            :key="t.id"
            class="p-3.5 rounded-xl border space-y-2 text-xs"
            style="background-color: var(--bg-content); border-color: var(--border-color);"
          >
            <div class="flex items-center justify-between font-bold">
              <span class="text-emerald-600 font-mono">#{{ t.queueNumber }}</span>
              <span class="text-slate-400 font-normal">{{ t.time }}</span>
            </div>
            <div class="font-bold text-sm text-slate-900 dark:text-slate-100">
              {{ t.quantity }}x {{ t.productName }}
            </div>
            <div class="text-[11px] text-slate-500">
              Antarkan ke: <strong>{{ t.table }}</strong>
            </div>
            <button
              type="button"
              @click="updateTicketStatus(t.id, 'served')"
              class="w-full py-1.5 rounded-lg text-xs font-semibold bg-emerald-600 hover:bg-emerald-700 text-white cursor-pointer transition-colors shadow-xs"
            >
              Selesai Saji
            </button>
          </div>
        </div>
      </div>

      <!-- 4. Sudah Disajikan -->
      <div class="flex flex-col h-full rounded-2xl border overflow-hidden"
           style="background-color: var(--bg-primary); border-color: var(--border-color);">
        <div class="px-4 py-2.5 bg-slate-700 text-white font-bold text-xs uppercase tracking-wider flex items-center justify-between shrink-0">
          <div class="flex items-center gap-2">
            <CheckCircle2 class="w-3.5 h-3.5" />
            <span>4. Selesai</span>
          </div>
          <span class="px-2 py-0.2 rounded-full bg-white/20 text-xs font-mono font-bold">{{ servedTickets.length }}</span>
        </div>

        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in servedTickets"
            :key="t.id"
            class="p-3.5 rounded-xl border text-xs opacity-75"
            style="background-color: var(--bg-content); border-color: var(--border-color);"
          >
            <div class="flex items-center justify-between font-bold">
              <span class="text-slate-500 font-mono">#{{ t.queueNumber }}</span>
              <span class="text-slate-400 font-normal">{{ t.time }}</span>
            </div>
            <div class="font-medium text-slate-900 dark:text-slate-100 mt-1">
              {{ t.quantity }}x {{ t.productName }}
            </div>
            <div class="text-[11px] text-slate-400">Meja: {{ t.table }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  ChefHat,
  RotateCw,
  Clock,
  CheckCircle2,
  Layers,
  Armchair,
  ShoppingBag,
  Truck
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const loading = ref(false)
const stationFilter = ref('Semua Station')
const activeTypeFilter = ref('all')
const currentTime = ref('')
let clockInterval: any = null

interface KDSTicket {
  id: string
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
}

const tickets = ref<KDSTicket[]>([
  { id: 't1', queueNumber: 'D-01', orderType: 'dine_in', customerName: 'Budi Santoso', table: 'Meja T-04', productName: 'Caffe Latte Oatmilk', quantity: 2, status: 'pending', station: 'barista', notes: 'Less ice, Oatly milk', time: '12:04' },
  { id: 't2', queueNumber: 'D-02', orderType: 'dine_in', customerName: 'Siti Rahma', table: 'Meja T-02', productName: 'Croissant Butter French', quantity: 1, status: 'cooking', station: 'kitchen', notes: 'Hangatkan 2 menit', time: '12:01' },
  { id: 't3', queueNumber: 'TA-05', orderType: 'takeaway', customerName: 'Fajar', table: 'Takeaway', productName: 'V60 Gayo Natural', quantity: 1, status: 'ready', station: 'barista', notes: 'Ice cup takeaway', time: '11:55' },
  { id: 't4', queueNumber: 'D-03', orderType: 'dine_in', customerName: 'Hendra', table: 'Meja T-05', productName: 'Americano Ice', quantity: 2, status: 'served', station: 'barista', notes: '', time: '11:42' }
])

const orderTypeTabs = [
  { id: 'all', label: 'Semua Tiket', icon: Layers },
  { id: 'dine_in', label: 'Dine-In', icon: Armchair },
  { id: 'takeaway', label: 'Takeaway', icon: ShoppingBag }
]

const getCountByType = (type: string) => {
  if (type === 'all') return tickets.value.length
  return tickets.value.filter(t => t.orderType === type).length
}

const filteredTickets = computed(() => {
  return tickets.value.filter(t => {
    const matchType = activeTypeFilter.value === 'all' || t.orderType === activeTypeFilter.value
    const matchStation = stationFilter.value === 'Semua Station' || t.station === stationFilter.value
    return matchType && matchStation
  })
})

const pendingTickets = computed(() => filteredTickets.value.filter(t => t.status === 'pending'))
const cookingTickets = computed(() => filteredTickets.value.filter(t => t.status === 'cooking'))
const readyTickets = computed(() => filteredTickets.value.filter(t => t.status === 'ready'))
const servedTickets = computed(() => filteredTickets.value.filter(t => t.status === 'served'))

const updateTicketStatus = (id: string, newStatus: 'cooking' | 'ready' | 'served') => {
  const target = tickets.value.find(t => t.id === id)
  if (target) {
    target.status = newStatus
    notifyStore.success(`Status tiket #${target.queueNumber} diperbarui menjadi ${newStatus.toUpperCase()}`, 'KDS Update')
  }
}

const fetchKDSTickets = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    notifyStore.info('Antrean pesanan KDS sinkron dengan kasir POS', 'KDS Refresh')
  }, 300)
}

onMounted(() => {
  currentTime.value = new Date().toLocaleTimeString('id-ID')
  clockInterval = setInterval(() => {
    currentTime.value = new Date().toLocaleTimeString('id-ID')
  }, 1000)
})

onUnmounted(() => {
  if (clockInterval) clearInterval(clockInterval)
})
</script>

