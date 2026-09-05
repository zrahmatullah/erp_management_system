<template>
  <div class="space-y-6">
    <!-- Stat Cards 4 Kolom -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
      <!-- Card 1: Revenue -->
      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-500 mb-1">1. Total Revenue</div>
          <div class="text-2xl font-black text-slate-900 tracking-tight">
            Rp {{ formatNum(stats.today_sales > 0 ? stats.today_sales : 45890000) }}
          </div>
          <div class="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600 mt-2 bg-emerald-50 px-2 py-0.5 rounded-full">
            <TrendingUp class="w-3.5 h-3.5" /> 12.5% vs yesterday
          </div>
        </div>
        <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
          <DollarSign class="w-6 h-6" />
        </div>
      </div>

      <!-- Card 2: Orders Today -->
      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-500 mb-1">2. Orders Today</div>
          <div class="text-2xl font-black text-slate-900 tracking-tight">
            {{ stats.total_orders > 0 ? stats.total_orders : 127 }}
          </div>
          <div class="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600 mt-2 bg-emerald-50 px-2 py-0.5 rounded-full">
            <TrendingUp class="w-3.5 h-3.5" /> 8.2%
          </div>
        </div>
        <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
          <ClipboardList class="w-6 h-6" />
        </div>
      </div>

      <!-- Card 3: Low Stock Warnings -->
      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-500 mb-1">3. Low Stock Items</div>
          <div class="text-2xl font-black text-slate-900 tracking-tight">
            {{ stats.low_stock_count }} items
          </div>
          <div class="inline-flex items-center gap-1 text-xs font-semibold text-amber-600 mt-2 bg-amber-50 px-2 py-0.5 rounded-full">
            <AlertTriangle class="w-3.5 h-3.5" /> Perlu Restock
          </div>
        </div>
        <div class="w-11 h-11 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center">
          <Package class="w-6 h-6" />
        </div>
      </div>

      <!-- Card 4: Active Occupied Tables -->
      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-500 mb-1">4. Occupied Tables</div>
          <div class="text-2xl font-black text-slate-900 tracking-tight">
            {{ stats.active_tables }} Meja
          </div>
          <div class="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600 mt-2 bg-emerald-50 px-2 py-0.5 rounded-full">
            <TrendingUp class="w-3.5 h-3.5" /> Terisi Aktif
          </div>
        </div>
        <div class="w-11 h-11 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
          <Armchair class="w-6 h-6" />
        </div>
      </div>
    </div>

    <!-- Middle Section: Revenue Trend Chart & Sales by Category -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Weekly Revenue Trend (2 cols) -->
      <div class="lg:col-span-2 bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="font-bold text-slate-900 text-base">Weekly Revenue Trend</h3>
            <p class="text-xs text-slate-400">Pola pergerakan omset mingguan seluruh outlet</p>
          </div>
          <div class="flex items-center gap-2 text-xs font-medium">
            <span class="inline-flex items-center gap-1 text-blue-600">
              <span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span> Minggu Ini
            </span>
            <span class="inline-flex items-center gap-1 text-slate-400">
              <span class="w-2.5 h-2.5 rounded-full bg-slate-300"></span> Minggu Lalu
            </span>
          </div>
        </div>

        <!-- Custom SVG Curve Visual -->
        <div class="h-64 flex flex-col justify-between">
          <div class="relative h-52 w-full border-b border-slate-100">
            <svg class="w-full h-full" viewBox="0 0 600 200" preserveAspectRatio="none">
              <defs>
                <linearGradient id="gradientArea" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#2563EB" stop-opacity="0.25" />
                  <stop offset="100%" stop-color="#2563EB" stop-opacity="0.0" />
                </linearGradient>
              </defs>
              <line x1="0" y1="50" x2="600" y2="50" stroke="#f1f5f9" stroke-width="1" />
              <line x1="0" y1="100" x2="600" y2="100" stroke="#f1f5f9" stroke-width="1" />
              <line x1="0" y1="150" x2="600" y2="150" stroke="#f1f5f9" stroke-width="1" />

              <path
                d="M0,130 Q100,160 200,90 T400,120 T600,60"
                fill="none"
                stroke="#CBD5E1"
                stroke-width="2"
                stroke-dasharray="4 4"
              />
              <path
                d="M0,120 Q100,70 200,80 T400,50 T600,30 L600,200 L0,200 Z"
                fill="url(#gradientArea)"
              />
              <path
                d="M0,120 Q100,70 200,80 T400,50 T600,30"
                fill="none"
                stroke="#2563EB"
                stroke-width="3"
              />
            </svg>
          </div>
          <div class="flex justify-between text-[11px] text-slate-400 pt-2 font-medium">
            <span>Sen</span>
            <span>Sel</span>
            <span>Rab</span>
            <span>Kam</span>
            <span>Jum</span>
            <span>Sab</span>
            <span>Min</span>
          </div>
        </div>
      </div>

      <!-- Sales by Category (1 col) -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm flex flex-col justify-between">
        <div>
          <h3 class="font-bold text-slate-900 text-base">Sales by Category</h3>
          <p class="text-xs text-slate-400 mb-6">Kontribusi omset per kelompok menu</p>

          <div class="relative flex items-center justify-center my-4">
            <svg class="w-44 h-44 transform -rotate-90" viewBox="0 0 100 100">
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#2563EB" stroke-width="14" stroke-dasharray="100 238" stroke-dashoffset="0" />
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#38BDF8" stroke-width="14" stroke-dasharray="65 238" stroke-dashoffset="-100" />
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#F59E0B" stroke-width="14" stroke-dasharray="45 238" stroke-dashoffset="-165" />
              <circle cx="50" cy="50" r="38" fill="transparent" stroke="#10B981" stroke-width="14" stroke-dasharray="28 238" stroke-dashoffset="-210" />
            </svg>
            <div class="absolute text-center">
              <div class="text-xl font-black text-slate-900">100%</div>
              <div class="text-[10px] text-slate-400 uppercase font-bold tracking-wider">Katalog</div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-2 text-xs pt-2 border-t border-slate-100">
          <div class="flex items-center gap-2 text-slate-600">
            <span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span> Coffee (42%)
          </div>
          <div class="flex items-center gap-2 text-slate-600">
            <span class="w-2.5 h-2.5 rounded-full bg-sky-400"></span> Non-Coffee (27%)
          </div>
          <div class="flex items-center gap-2 text-slate-600">
            <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span> Food (19%)
          </div>
          <div class="flex items-center gap-2 text-slate-600">
            <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span> Pastry (12%)
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Section: Recent Orders & Table Layout Preview -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Orders Table (2 cols) -->
      <div class="lg:col-span-2 bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="font-bold text-slate-900 text-base">Recent Live Orders</h3>
            <p class="text-xs text-slate-400">Pesanan kasir dan pesanan meja terkini</p>
          </div>
          <router-link to="/pos" class="text-xs text-blue-600 font-semibold hover:underline">
            Buka Kasir POS →
          </router-link>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr class="text-slate-400 border-b border-slate-100 font-semibold">
                <th class="pb-3">Order ID</th>
                <th class="pb-3">Customer</th>
                <th class="pb-3">Meja / Tipe</th>
                <th class="pb-3">Total</th>
                <th class="pb-3">Status</th>
                <th class="pb-3 text-right">Time</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 font-medium text-slate-700">
              <tr v-for="ord in recentOrders" :key="ord.id" class="hover:bg-slate-50/80 transition-colors">
                <td class="py-3 font-semibold text-slate-900">{{ ord.order_number }}</td>
                <td class="py-3">{{ ord.customer }}</td>
                <td class="py-3 text-slate-500">Meja {{ ord.table_number }} ({{ ord.order_type }})</td>
                <td class="py-3 font-bold text-slate-900">Rp {{ formatNum(ord.total) }}</td>
                <td class="py-3">
                  <span 
                    class="px-2.5 py-0.5 rounded-full text-[11px] font-semibold"
                    :class="ord.status === 'completed' ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'"
                  >
                    {{ ord.status }}
                  </span>
                </td>
                <td class="py-3 text-right text-slate-400">{{ ord.created_at }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Table Layout Mini Preview (1 col) -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/80 shadow-sm flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between mb-2">
            <h3 class="font-bold text-slate-900 text-base">Table Layout</h3>
            <router-link to="/pos/tables" class="text-xs text-blue-600 font-semibold hover:underline">
              Buka Denah →
            </router-link>
          </div>
          <p class="text-xs text-slate-400">Visual denah meja operasional cabang</p>
        </div>

        <!-- Mini Grid Floor Plan -->
        <div class="bg-slate-50 rounded-xl p-4 border border-slate-100 my-3 grid grid-cols-4 gap-2.5">
          <div 
            v-for="t in floorTables.slice(0, 8)" 
            :key="t.id"
            class="p-2.5 rounded-lg border text-center text-xs font-bold transition-all shadow-xs"
            :class="t.status === 'occupied' 
              ? 'bg-red-500 text-white border-red-600' 
              : t.status === 'reserved' 
                ? 'bg-amber-400 text-slate-900 border-amber-500' 
                : 'bg-white border-emerald-500 text-emerald-700'"
          >
            <div>{{ t.table_number }}</div>
            <div class="text-[10px] font-normal" :class="t.status === 'occupied' ? 'text-red-100' : 'text-slate-400'">
              {{ t.status === 'occupied' ? 'Terisi' : `${t.capacity} Seat` }}
            </div>
          </div>
        </div>

        <div class="flex justify-around text-[11px] text-slate-500 pt-2 border-t border-slate-100">
          <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded bg-emerald-500"></span> Tersedia</div>
          <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded bg-red-500"></span> Terisi</div>
          <div class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded bg-amber-400"></span> Reservasi</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import {
  DollarSign,
  ClipboardList,
  Package,
  Armchair,
  TrendingUp,
  AlertTriangle
} from 'lucide-vue-next'

const stats = ref({
  today_sales: 0,
  total_orders: 0,
  low_stock_count: 0,
  active_tables: 0,
  avg_order_value: 0
})

const recentOrders = ref<any[]>([])
const floorTables = ref<any[]>([])

const formatNum = (val: number) => {
  return Number(val || 0).toLocaleString('id-ID')
}

const fetchDashboardData = async () => {
  try {
    const [sRes, oRes, tRes] = await Promise.all([
      axios.get('/api/v1/dashboard/stats'),
      axios.get('/api/v1/pos/orders'),
      axios.get('/api/v1/pos/tables')
    ])
    if (sRes.data?.data) {
      stats.value = sRes.data.data
    }
    if (oRes.data?.data) {
      recentOrders.value = oRes.data.data
    }
    if (tRes.data?.data) {
      floorTables.value = tRes.data.data
    }
  } catch (err) {
    console.error('Failed to load dashboard stats', err)
  }
}

onMounted(() => {
  fetchDashboardData()
})
</script>
