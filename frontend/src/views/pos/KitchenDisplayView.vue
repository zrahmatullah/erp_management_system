<template>
  <div class="h-[calc(100vh-5.5rem)] flex flex-col space-y-4">
    <!-- Top Header Bar -->
    <div class="bg-slate-900 text-white rounded-2xl p-4 flex items-center justify-between shadow-lg shrink-0">
      <div class="flex items-center gap-4">
        <div class="w-10 h-10 rounded-xl bg-blue-600 flex items-center justify-center font-bold">
          <ChefHat class="w-6 h-6 text-white" />
        </div>
        <div>
          <h2 class="text-lg font-black tracking-tight flex items-center gap-2">
            Kitchen Display System (KDS)
            <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
          </h2>
          <p class="text-xs text-slate-400">Station: <strong>Barista & Central Kitchen</strong></p>
        </div>
      </div>

      <!-- Real-time Clock -->
      <div class="flex items-center gap-2 font-mono text-xl font-black bg-slate-800 px-4 py-2 rounded-xl border border-slate-700 text-emerald-400">
        <Clock class="w-5 h-5 text-emerald-400" />
        <span>{{ currentTime }}</span>
      </div>
    </div>

    <!-- 4 Kanban Columns: Pesanan Masuk, Sedang Dimasak, Siap Saji, Disajikan -->
    <div class="flex-1 grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4 min-h-0 overflow-hidden">
      <!-- 1. Pesanan Masuk (Pending) -->
      <div class="flex flex-col h-full bg-slate-900/60 rounded-2xl border border-slate-800/80 overflow-hidden">
        <div class="px-4 py-3 bg-amber-500 text-slate-950 font-black text-sm uppercase tracking-wider flex items-center justify-between">
          <span>Pesanan Masuk</span>
          <span class="px-2 py-0.5 rounded-full bg-black/20 text-xs">{{ pendingTickets.length }}</span>
        </div>
        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in pendingTickets"
            :key="t.id"
            class="bg-slate-900 p-4 rounded-xl border border-amber-500/50 shadow-md space-y-3"
          >
            <div class="flex items-center justify-between">
              <span class="font-black text-sm text-white">{{ t.orderNumber }}</span>
              <span class="px-2 py-0.5 rounded bg-amber-500/20 text-amber-300 text-[10px] font-bold">Dine-in</span>
            </div>
            <div class="flex items-center justify-between text-xs text-slate-400">
              <span class="font-bold text-slate-200">Meja {{ t.table }}</span>
              <span class="font-mono text-amber-400">⏱ {{ t.timer }}</span>
            </div>
            <div class="space-y-1.5 text-xs text-slate-300 border-t border-slate-800 pt-2">
              <div v-for="(item, idx) in t.items" :key="idx" class="flex items-center gap-2">
                <span class="w-1.5 h-1.5 rounded-full bg-amber-400"></span>
                <span>{{ item }}</span>
              </div>
            </div>
            <button
              @click="advanceTicket(t, 'cooking')"
              class="w-full py-1.5 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs rounded-lg transition-colors cursor-pointer"
            >
              Mulai Masak →
            </button>
          </div>
          <div v-if="pendingTickets.length === 0" class="h-32 flex items-center justify-center text-slate-600 text-xs">
            Tidak ada tiket antrian baru
          </div>
        </div>
      </div>

      <!-- 2. Sedang Dimasak (Cooking) -->
      <div class="flex flex-col h-full bg-slate-900/60 rounded-2xl border border-slate-800/80 overflow-hidden">
        <div class="px-4 py-3 bg-blue-600 text-white font-black text-sm uppercase tracking-wider flex items-center justify-between">
          <span>Sedang Dimasak</span>
          <span class="px-2 py-0.5 rounded-full bg-black/20 text-xs">{{ inProgressTickets.length }}</span>
        </div>
        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in inProgressTickets"
            :key="t.id"
            class="bg-slate-900 p-4 rounded-xl border border-blue-500/50 shadow-md space-y-3"
          >
            <div class="flex items-center justify-between">
              <span class="font-black text-sm text-white">{{ t.orderNumber }}</span>
              <span class="px-2 py-0.5 rounded bg-blue-500/20 text-blue-300 text-[10px] font-bold">Dine-in</span>
            </div>
            <div class="flex items-center justify-between text-xs text-slate-400">
              <span class="font-bold text-slate-200">Meja {{ t.table }}</span>
              <span class="font-mono text-blue-400">⏱ {{ t.timer }}</span>
            </div>
            <div class="space-y-1.5 text-xs text-slate-300 border-t border-slate-800 pt-2">
              <div v-for="(item, idx) in t.items" :key="idx" class="flex items-center gap-2">
                <span class="w-1.5 h-1.5 rounded-full bg-blue-400"></span>
                <span>{{ item }}</span>
              </div>
            </div>
            <button
              @click="advanceTicket(t, 'ready')"
              class="w-full py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-lg transition-colors cursor-pointer"
            >
              Selesai Masak →
            </button>
          </div>
          <div v-if="inProgressTickets.length === 0" class="h-32 flex items-center justify-center text-slate-600 text-xs">
            Dapur sedang tidak memasak pesanan
          </div>
        </div>
      </div>

      <!-- 3. Siap Saji (Ready) -->
      <div class="flex flex-col h-full bg-slate-900/60 rounded-2xl border border-slate-800/80 overflow-hidden">
        <div class="px-4 py-3 bg-emerald-600 text-white font-black text-sm uppercase tracking-wider flex items-center justify-between">
          <span>Siap Saji</span>
          <span class="px-2 py-0.5 rounded-full bg-black/20 text-xs">{{ readyTickets.length }}</span>
        </div>
        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in readyTickets"
            :key="t.id"
            class="bg-slate-900 p-4 rounded-xl border border-emerald-500/50 shadow-md space-y-3"
          >
            <div class="flex items-center justify-between">
              <span class="font-black text-sm text-white">{{ t.orderNumber }}</span>
              <span class="px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-300 text-[10px] font-bold">Dine-in</span>
            </div>
            <div class="flex items-center justify-between text-xs text-slate-400">
              <span class="font-bold text-slate-200">Meja {{ t.table }}</span>
              <span class="font-mono text-emerald-400">⏱ {{ t.timer }}</span>
            </div>
            <div class="space-y-1.5 text-xs text-slate-300 border-t border-slate-800 pt-2">
              <div v-for="(item, idx) in t.items" :key="idx" class="flex items-center gap-2">
                <Check class="w-3.5 h-3.5 text-emerald-400" />
                <span>{{ item }}</span>
              </div>
            </div>
            <button
              @click="advanceTicket(t, 'served')"
              class="w-full py-1.5 bg-slate-800 hover:bg-slate-700 text-emerald-400 font-bold text-xs rounded-lg border border-emerald-500/30 transition-colors cursor-pointer"
            >
              Next Status (Disajikan) →
            </button>
          </div>
          <div v-if="readyTickets.length === 0" class="h-32 flex items-center justify-center text-slate-600 text-xs">
            Belum ada pesanan siap di pick-up
          </div>
        </div>
      </div>

      <!-- 4. Disajikan (Served) -->
      <div class="flex flex-col h-full bg-slate-900/60 rounded-2xl border border-slate-800/80 overflow-hidden">
        <div class="px-4 py-3 bg-slate-700 text-white font-black text-sm uppercase tracking-wider flex items-center justify-between">
          <span>Disajikan</span>
          <span class="px-2 py-0.5 rounded-full bg-black/20 text-xs">{{ servedTickets.length }}</span>
        </div>
        <div class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-for="t in servedTickets"
            :key="t.id"
            class="bg-slate-900/50 p-4 rounded-xl border border-slate-800 space-y-3 opacity-70"
          >
            <div class="flex items-center justify-between">
              <span class="font-black text-sm text-slate-300">{{ t.orderNumber }}</span>
              <span class="px-2 py-0.5 rounded bg-slate-800 text-slate-400 text-[10px] font-bold">Dine-in</span>
            </div>
            <div class="flex items-center justify-between text-xs text-slate-500">
              <span>Meja {{ t.table }}</span>
              <span class="font-mono">⏱ {{ t.timer }}</span>
            </div>
            <div class="space-y-1 text-xs text-slate-400 border-t border-slate-800 pt-2">
              <div v-for="(item, idx) in t.items" :key="idx" class="flex items-center gap-1.5">
                <Check class="w-3.5 h-3.5 text-emerald-400" />
                <span>{{ item }}</span>
              </div>
            </div>
            <div class="text-[11px] text-center text-emerald-400 font-bold py-1">
              Selesai disajikan ke meja
            </div>
          </div>
          <div v-if="servedTickets.length === 0" class="h-32 flex items-center justify-center text-slate-600 text-xs">
            Tidak ada riwayat sajian baru
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { ChefHat, Clock, Check } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const currentTime = ref(new Date().toLocaleTimeString('id-ID'))

interface Ticket {
  id: string
  orderNumber: string
  table: string
  timer: string
  status: 'pending' | 'cooking' | 'ready' | 'served'
  items: string[]
}

const tickets = ref<Ticket[]>([])

const fetchKDSTickets = async () => {
  try {
    const res = await axios.get('/api/v1/kds/tickets')
    if (res.data?.data) {
      tickets.value = res.data.data.map((t: any) => ({
        id: t.id,
        orderNumber: t.order_number,
        table: t.table_number || 'T-01',
        timer: t.time || `${t.elapsed_min || 0}m`,
        status: t.status || 'cooking',
        items: [`${t.product_name} x${t.quantity}`]
      }))
    }
  } catch (err: any) {
    console.error('Failed to load KDS tickets:', err)
  }
}

onMounted(() => {
  fetchKDSTickets()
  setInterval(() => {
    currentTime.value = new Date().toLocaleTimeString('id-ID')
  }, 1000)
})

const pendingTickets = computed(() => tickets.value.filter(t => t.status === 'pending'))
const inProgressTickets = computed(() => tickets.value.filter(t => t.status === 'cooking'))
const readyTickets = computed(() => tickets.value.filter(t => t.status === 'ready'))
const servedTickets = computed(() => tickets.value.filter(t => t.status === 'served'))

const advanceTicket = async (ticket: Ticket, nextStatus: any) => {
  ticket.status = nextStatus
  try {
    await axios.put(`/api/v1/kds/items/${ticket.id}/status`, { status: nextStatus })
  } catch (err) {
    console.warn('Could not persist status to backend, updated locally:', err)
  }
  notifyStore.success(`Pesanan ${ticket.orderNumber} (Meja ${ticket.table}) dipindahkan ke '${nextStatus}'!`, 'KDS Status')
}
</script>
