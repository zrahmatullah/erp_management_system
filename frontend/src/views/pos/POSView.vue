<template>
  <div class="space-y-4">
    <!-- Top Mode Switcher: 2 Tabs di Kasir (Order Baru & Pembayaran) -->
    <div class="bg-white rounded-2xl p-2.5 border border-slate-200/80 shadow-xs flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="flex items-center gap-2">
        <button
          @click="activePosTab = 'order'"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-black transition-all cursor-pointer"
          :class="activePosTab === 'order' 
            ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' 
            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'"
        >
          <Utensils class="w-4 h-4" />
          <span>Katalog & Pesanan Baru</span>
        </button>

        <button
          @click="activePosTab = 'billing'"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-black transition-all cursor-pointer relative"
          :class="activePosTab === 'billing' 
            ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' 
            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'"
        >
          <CreditCard class="w-4 h-4" />
          <span>Kasir Pembayaran Meja</span>
          <span
            v-if="activeOrders.length > 0"
            class="px-2 py-0.5 rounded-full text-[10px] font-extrabold transition-all"
            :class="activePosTab === 'billing' ? 'bg-white text-blue-700' : 'bg-red-500 text-white'"
          >
            {{ activeOrders.length }} Tagihan
          </span>
        </button>

        <button
          @click="activePosTab = 'history'"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-black transition-all cursor-pointer"
          :class="activePosTab === 'history' 
            ? 'bg-blue-600 text-white shadow-md shadow-blue-600/30' 
            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'"
        >
          <ReceiptText class="w-4 h-4" />
          <span>Riwayat Transaksi & Struk</span>
        </button>
      </div>

      <!-- Active Cashier Bar -->
      <div class="flex items-center gap-4 text-xs font-semibold text-slate-700 pr-2">
        <div class="flex items-center gap-1.5">
          <User class="w-4 h-4 text-blue-600" />
          <span>Kasir: <strong class="text-slate-900">Jane Cashier</strong></span>
        </div>
        <span class="text-slate-300">|</span>
        <div class="flex items-center gap-1.5 text-slate-500">
          <Clock class="w-4 h-4 text-slate-400" />
          <span>Shift Pagi</span>
        </div>
      </div>
    </div>

    <!-- TAB 1: ORDER BARU & KATALOG MENU -->
    <div v-if="activePosTab === 'order'" class="h-[calc(100vh-10rem)] flex flex-col lg:flex-row gap-6">
      <!-- Left Section: Catalog & Filter -->
      <div class="flex-1 flex flex-col min-w-0 space-y-4">
        <!-- Search & Category Header -->
        <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-sm space-y-4 shrink-0">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div class="text-xs font-bold text-slate-700">
              Pilih menu untuk ditambahkan ke keranjang pesanan
            </div>

            <!-- Live Search Bar -->
            <div class="relative w-full sm:w-72">
              <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Cari menu kopi, makanan..."
                class="w-full pl-9 pr-3 py-1.5 bg-slate-50 border border-slate-200 rounded-xl text-xs focus:bg-white focus:border-blue-500 transition-all outline-none"
              />
            </div>
          </div>

          <!-- Horizontal Category Pills -->
          <div class="flex items-center gap-2 overflow-x-auto pb-1 scrollbar-none text-xs">
            <button
              v-for="cat in availableCategories"
              :key="cat"
              @click="selectedCategory = cat"
              class="px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer shrink-0"
              :class="selectedCategory === cat 
                ? 'bg-blue-600 text-white shadow-sm shadow-blue-500/20' 
                : 'bg-slate-50 text-slate-600 hover:bg-slate-100 hover:text-slate-900 border border-slate-200/60'"
            >
              {{ cat }}
            </button>
          </div>
        </div>

        <!-- Product Grid (4 columns) -->
        <div class="flex-1 overflow-y-auto pr-1">
          <div v-if="loading" class="h-64 flex items-center justify-center text-slate-400 text-xs">
            Memuat menu dari database server...
          </div>
          <div v-else-if="filteredProducts.length === 0" class="h-64 flex flex-col items-center justify-center text-slate-400 text-xs">
            <span>Tidak ada produk yang sesuai kriteria pencarian</span>
          </div>
          <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4">
            <div
              v-for="product in filteredProducts"
              :key="product.id"
              class="bg-white rounded-2xl p-3 border border-slate-200/80 shadow-sm flex flex-col justify-between hover:border-blue-300 hover:shadow-md transition-all group"
            >
              <div>
                <div class="w-full h-28 rounded-xl overflow-hidden mb-2.5 bg-slate-100">
                  <img 
                    :src="product.image || product.image_url || 'https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?w=300'" 
                    :alt="product.name" 
                    class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200" 
                  />
                </div>
                <h4 class="font-bold text-slate-900 text-xs truncate">{{ product.name }}</h4>
                <p class="text-xs font-extrabold text-blue-600 mt-1">Rp {{ (product.price || product.base_price || 0).toLocaleString('id-ID') }}</p>
              </div>

              <button
                @click="addToCart(product)"
                class="mt-3 w-full py-1.5 bg-blue-50 hover:bg-blue-600 text-blue-600 hover:text-white text-xs font-bold rounded-xl transition-colors flex items-center justify-center gap-1 cursor-pointer"
              >
                <span>Add</span>
                <span>+</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Section: Order Cart Panel -->
      <div class="w-full lg:w-96 bg-white rounded-2xl border border-slate-200/80 shadow-sm flex flex-col h-full overflow-hidden shrink-0">
        <!-- Order Type Tabs -->
        <div class="p-4 border-b border-slate-100 space-y-3 shrink-0">
          <div class="grid grid-cols-3 gap-1 bg-slate-100 p-1 rounded-xl">
            <button
              v-for="ot in ['Dine-in', 'Takeaway', 'Delivery']"
              :key="ot"
              @click="orderType = ot"
              class="py-1.5 text-xs font-bold rounded-lg transition-all cursor-pointer text-center"
              :class="orderType === ot ? 'bg-white text-blue-600 shadow-xs' : 'text-slate-500 hover:text-slate-800'"
            >
              {{ ot }}
            </button>
          </div>

          <!-- Table Selector -->
          <div v-if="orderType === 'Dine-in'" class="flex items-center justify-between bg-slate-50 px-3 py-2 rounded-xl border border-slate-200/70">
            <span class="text-xs font-medium text-slate-600">Table:</span>
            <select v-model="selectedTable" class="text-xs font-bold text-slate-800 bg-transparent outline-none cursor-pointer">
              <option v-for="t in tablesList" :key="t.id" :value="t.code">
                {{ t.code }} ({{ t.capacity }} Kursi)
              </option>
            </select>
          </div>
        </div>

        <!-- Cart Items List -->
        <div class="flex-1 overflow-y-auto p-4 space-y-3">
          <div v-if="cart.length === 0" class="h-full flex flex-col items-center justify-center text-slate-400 text-xs py-8">
            <ShoppingCart class="w-10 h-10 text-slate-300 mb-2" />
            <span class="font-semibold text-slate-500">Pesanan masih kosong</span>
            <span class="text-[11px] text-slate-400 mt-0.5">Pilih menu dari katalog di sebelah kiri</span>
          </div>

          <div
            v-for="item in cart"
            :key="item.id"
            class="flex items-center justify-between pb-3 border-b border-slate-100"
          >
            <div class="flex-1 pr-2">
              <h5 class="text-xs font-bold text-slate-800 leading-tight">{{ item.name }}</h5>
              <div class="text-[11px] text-slate-400 font-medium">Rp {{ item.price.toLocaleString('id-ID') }}</div>
            </div>

            <!-- Quantity Stepper -->
            <div class="flex items-center gap-2">
              <div class="flex items-center border border-slate-200 rounded-lg overflow-hidden bg-slate-50">
                <button
                  @click="updateQty(item.id, -1)"
                  class="w-6 h-6 flex items-center justify-center text-xs font-bold text-slate-600 hover:bg-slate-200 cursor-pointer"
                >
                  <Minus class="w-3 h-3" />
                </button>
                <span class="w-7 text-center text-xs font-bold text-slate-800">{{ item.qty }}</span>
                <button
                  @click="updateQty(item.id, 1)"
                  class="w-6 h-6 flex items-center justify-center text-xs font-bold text-slate-600 hover:bg-slate-200 cursor-pointer"
                >
                  <Plus class="w-3 h-3" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Bill Breakdown & Action Buttons -->
        <div class="p-4 border-t border-slate-100 bg-slate-50/50 space-y-3 shrink-0">
          <div class="space-y-1.5 text-xs">
            <div class="flex justify-between text-slate-500">
              <span>Subtotal</span>
              <span class="font-semibold text-slate-800">Rp {{ subtotal.toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-slate-500">
              <span>PPN 10%</span>
              <span class="font-semibold text-slate-800">Rp {{ taxAmount.toLocaleString('id-ID') }}</span>
            </div>
            <div v-if="discount > 0" class="flex justify-between text-emerald-600">
              <span>Diskon</span>
              <span>-Rp {{ discount.toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-slate-900 font-black text-sm pt-2 border-t border-slate-200">
              <span>Total Akhir</span>
              <span class="text-blue-600">Rp {{ grandTotal.toLocaleString('id-ID') }}</span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2 pt-2">
            <button
              @click="clearCart"
              :disabled="cart.length === 0"
              class="py-2 px-3 border border-slate-200 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 transition-colors disabled:opacity-50 cursor-pointer"
            >
              Kosongkan
            </button>
            <button
              @click="saveOrderToTable"
              :disabled="cart.length === 0 || savingOrder"
              class="py-2 px-3 bg-amber-500 hover:bg-amber-600 active:bg-amber-700 text-white rounded-xl text-xs font-bold transition-all disabled:opacity-50 cursor-pointer flex items-center justify-center gap-1.5 shadow-sm shadow-amber-500/20"
            >
              <Send class="w-3.5 h-3.5" />
              <span>
                {{ savingOrder ? 'Mengirim...' : (orderType === 'Dine-in' ? 'Kirim ke Meja' : (orderType === 'Takeaway' ? 'Pesan Takeaway' : 'Pesan Delivery')) }}
              </span>
            </button>
          </div>

          <button
            @click="openPayment"
            :disabled="cart.length === 0"
            class="w-full py-3 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white rounded-xl font-bold text-xs shadow-lg shadow-blue-600/30 transition-all disabled:opacity-50 cursor-pointer flex items-center justify-center gap-2"
          >
            <CreditCard class="w-4 h-4" />
            <span>Bayar Sekarang (Rp {{ grandTotal.toLocaleString('id-ID') }})</span>
          </button>
        </div>
      </div>
    </div>

    <!-- TAB 2: KASIR PEMBAYARAN MEJA (MENUNGGU PEMBAYARAN / BILLING) -->
    <div v-else-if="activePosTab === 'billing'" class="space-y-4">
      <!-- Filter and Metrics Bar -->
      <div class="bg-white rounded-2xl p-4 border border-slate-200/80 shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-4">
        <!-- Zone Filter Pills -->
        <div class="flex items-center gap-2 overflow-x-auto text-xs font-bold">
          <button
            v-for="z in ['Semua Tagihan', 'Lantai 1', 'Lantai 2', 'Outdoor', 'Takeaway']"
            :key="z"
            @click="billingFilter = z"
            class="px-4 py-2 rounded-xl transition-all cursor-pointer whitespace-nowrap"
            :class="billingFilter === z ? 'bg-blue-600 text-white shadow-sm shadow-blue-600/20' : 'bg-slate-50 text-slate-600 hover:bg-slate-100'"
          >
            {{ z }}
          </button>
        </div>

        <!-- Search input & summary metrics -->
        <div class="flex items-center gap-3">
          <div class="relative w-64">
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
            <input
              v-model="billingSearch"
              type="text"
              placeholder="Cari meja / pelanggan / order..."
              class="w-full pl-9 pr-3 py-1.5 bg-slate-50 border border-slate-200 rounded-xl text-xs focus:bg-white focus:border-blue-500 outline-none"
            />
          </div>

          <button
            @click="fetchActiveOrders"
            class="p-2 text-slate-500 hover:text-blue-600 rounded-xl hover:bg-slate-100 transition-colors cursor-pointer"
            title="Refresh Tagihan"
          >
            <RotateCw class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Active Orders Cards Grid -->
      <div>
        <div v-if="loadingOrders" class="py-20 text-center text-slate-400 text-xs bg-white rounded-2xl border border-slate-200/80 shadow-sm">
          Memuat daftar tagihan aktif dari server...
        </div>
        <div v-else-if="filteredActiveOrders.length === 0" class="py-20 flex flex-col items-center justify-center text-slate-400 text-xs bg-white rounded-2xl border border-slate-200/80 shadow-sm space-y-2">
          <CheckCircle2 class="w-10 h-10 text-emerald-500" />
          <span class="font-bold text-slate-700 text-sm">Semua Meja Telah Lunas</span>
          <span class="text-slate-400">Tidak ada pesanan yang sedang menunggu pembayaran di kasir saat ini.</span>
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          <div
            v-for="order in filteredActiveOrders"
            :key="order.id"
            class="bg-white rounded-2xl border border-slate-200/90 shadow-sm hover:shadow-md transition-all p-5 flex flex-col justify-between space-y-4"
          >
            <!-- Card Header -->
            <div>
              <div class="flex items-start justify-between">
                <div class="flex items-center gap-2.5">
                  <div class="px-3 py-1.5 rounded-xl bg-blue-50 border border-blue-200 text-blue-700 font-black text-sm">
                    {{ order.table_number && order.table_number !== '-' ? `Meja ${order.table_number}` : (order.order_type === 'takeaway' ? 'Takeaway' : 'Delivery') }}
                  </div>
                  <div>
                    <div class="flex items-center gap-1.5">
                      <h4 class="font-bold text-slate-900 text-xs">{{ order.customer }}</h4>
                      <span v-if="order.queue_number && order.queue_number !== '-'" class="px-2 py-0.5 rounded-md bg-indigo-600 text-white font-black text-[10px] shadow-xs">
                        #{{ order.queue_number }}
                      </span>
                    </div>
                    <p class="text-[10px] text-slate-400 font-mono">{{ order.order_number }} • {{ order.zone_name }}</p>
                  </div>
                </div>

                <span class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase bg-amber-100 text-amber-800">
                  {{ order.status }}
                </span>
              </div>

              <!-- Time & Duration -->
              <div class="flex items-center gap-2 mt-3 text-[11px] text-slate-500 font-medium">
                <Clock class="w-3.5 h-3.5 text-slate-400" />
                <span>Masuk: {{ order.created_at }}</span>
                <span>•</span>
                <span class="text-blue-600 font-semibold">{{ order.duration_minutes || 1 }} menit lalu</span>
              </div>

              <!-- Items preview -->
              <div class="mt-3.5 pt-3 border-t border-slate-100 space-y-1.5">
                <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1">Menu Dipesan:</div>
                <div
                  v-for="item in (order.items || []).slice(0, 3)"
                  :key="item.id"
                  class="flex items-center justify-between text-xs text-slate-700"
                >
                  <span class="truncate max-w-[70%] font-medium">
                    <strong class="text-slate-900">{{ item.quantity }}x</strong> {{ item.name }}
                  </span>
                  <span class="font-semibold text-slate-800">Rp {{ (item.total_price || item.unit_price * item.quantity).toLocaleString('id-ID') }}</span>
                </div>
                <div v-if="order.items && order.items.length > 3" class="text-[10px] text-slate-400 font-medium italic">
                  + {{ order.items.length - 3 }} menu lainnya
                </div>
              </div>
            </div>

            <!-- Card Footer: Total & Action -->
            <div class="pt-3 border-t border-slate-100 space-y-3">
              <div class="flex items-center justify-between">
                <div>
                  <div class="text-[10px] text-slate-400 uppercase font-semibold">Total Tagihan:</div>
                  <div class="text-base font-black text-blue-600">Rp {{ (order.total || 0).toLocaleString('id-ID') }}</div>
                </div>
                <button
                  @click="printTempBill(order)"
                  class="px-2.5 py-1.5 rounded-xl border border-slate-200 hover:bg-slate-100 text-slate-600 text-[11px] font-bold transition-colors cursor-pointer"
                  title="Cetak Bill Sementara"
                >
                  <Receipt class="w-3.5 h-3.5 inline mr-1" />
                  Bill
                </button>
              </div>

              <button
                @click="payActiveOrder(order)"
                class="w-full py-2.5 bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white text-xs font-black rounded-xl shadow-md shadow-emerald-600/30 transition-all flex items-center justify-center gap-2 cursor-pointer"
              >
                <CreditCard class="w-4 h-4" />
                <span>Bayar Meja Ini (Rp {{ (order.total || 0).toLocaleString('id-ID') }})</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 3: RIWAYAT TRANSAKSI & STRUK -->
    <div v-else-if="activePosTab === 'history'">
      <TransactionHistoryView />
    </div>

    <!-- Direct Cart Payment Modal -->
    <PaymentModal
      v-if="showPaymentModal"
      :show="showPaymentModal"
      :total="grandTotal"
      :order-type="orderType"
      :table="orderType === 'Dine-in' ? selectedTable : '-'"
      :raw-items="cart"
      :items="cart.map(i => ({ name: i.name, qty: i.qty, price: i.price }))"
      :subtotal="subtotal"
      :tax="taxAmount"
      @close="showPaymentModal = false"
      @success="handlePaymentSuccess"
    />

    <!-- Active Table Order Payment Modal (Tab 2) -->
    <PaymentModal
      v-if="showBillingPaymentModal && selectedOrderForPayment"
      :show="showBillingPaymentModal"
      :order-id="selectedOrderForPayment.id"
      :order-number="selectedOrderForPayment.order_number"
      :queue-number="selectedOrderForPayment.queue_number"
      :table="selectedOrderForPayment.table_number"
      :order-type="selectedOrderForPayment.order_type === 'dine_in' ? 'Dine-in' : (selectedOrderForPayment.order_type === 'takeaway' ? 'Takeaway' : 'Delivery')"
      :items="selectedOrderForPayment.items?.map((i: any) => ({ name: i.name, qty: i.quantity, price: i.unit_price }))"
      :subtotal="selectedOrderForPayment.subtotal"
      :tax="selectedOrderForPayment.tax"
      :total="selectedOrderForPayment.total"
      @close="showBillingPaymentModal = false"
      @success="handleBillingPaymentSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import PaymentModal from './PaymentModal.vue'
import TransactionHistoryView from './TransactionHistoryView.vue'
import {
  User,
  Clock,
  Search,
  ShoppingCart,
  Plus,
  Minus,
  Utensils,
  CreditCard,
  Receipt,
  ReceiptText,
  RotateCw,
  CheckCircle2,
  Send
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const route = useRoute()
const notifyStore = useNotificationStore()

const activePosTab = ref('order') // 'order', 'billing', or 'history'
const selectedCategory = ref('All')
const searchQuery = ref('')
const orderType = ref('Dine-in')
const selectedTable = ref('T-01')
const discount = ref(0)
const showPaymentModal = ref(false)
const loading = ref(false)
const savingOrder = ref(false)

const products = ref<any[]>([])
const tablesList = ref<any[]>([])
const cart = ref<Array<{ id: string, name: string, price: number, qty: number }>>([])

// Billing Tab 2 state
const activeOrders = ref<any[]>([])
const loadingOrders = ref(false)
const billingFilter = ref('Semua Tagihan')
const billingSearch = ref('')
const selectedOrderForPayment = ref<any>(null)
const showBillingPaymentModal = ref(false)

const fetchPOSData = async () => {
  loading.value = true
  try {
    const [pRes, tRes] = await Promise.all([
      axios.get('/api/v1/pos/products'),
      axios.get('/api/v1/pos/tables')
    ])
    if (pRes.data?.data) {
      products.value = pRes.data.data.map((p: any) => ({
        id: p.id,
        name: p.name,
        price: p.price || p.base_price || 0,
        category: p.category || 'Coffee',
        image: p.image_url || p.image || 'https://images.unsplash.com/photo-1570968915860-54d5c301fa9f?w=300'
      }))
    }
    if (tRes.data?.data) {
      tablesList.value = tRes.data.data.map((t: any) => ({
        id: t.id,
        code: t.table_number || t.code,
        capacity: t.capacity || 4,
        status: t.status || 'available'
      }))
      if (tablesList.value.length > 0 && !route.query.table) {
        selectedTable.value = tablesList.value[0].code
      }
    }
  } catch (err: any) {
    console.error('Failed to load POS data:', err)
    notifyStore.error('Gagal mengambil katalog menu kasir dari server', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

const fetchActiveOrders = async () => {
  loadingOrders.value = true
  try {
    const res = await axios.get('/api/v1/pos/orders/active')
    if (res.data?.data) {
      activeOrders.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load active orders:', err)
  } finally {
    loadingOrders.value = false
  }
}

onMounted(() => {
  fetchPOSData()
  fetchActiveOrders()

  if (route.query.tab === 'billing') {
    activePosTab.value = 'billing'
  } else if (route.query.tab === 'history') {
    activePosTab.value = 'history'
  }
  if (route.query.table) {
    selectedTable.value = String(route.query.table)
    orderType.value = 'Dine-in'
    billingSearch.value = String(route.query.table)
  }
})

const availableCategories = computed(() => {
  const cats = new Set<string>()
  cats.add('All')
  products.value.forEach(p => {
    if (p.category) cats.add(p.category)
  })
  return Array.from(cats)
})

const filteredProducts = computed(() => {
  return products.value.filter(p => {
    const matchCat = selectedCategory.value === 'All' || p.category === selectedCategory.value
    const matchSearch = p.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchCat && matchSearch
  })
})

const subtotal = computed(() => {
  return cart.value.reduce((acc, item) => acc + (item.price * item.qty), 0)
})

const taxAmount = computed(() => {
  return Math.round((subtotal.value - discount.value) * 0.10)
})

const grandTotal = computed(() => {
  return Math.max(0, subtotal.value - discount.value + taxAmount.value)
})

const addToCart = (product: any) => {
  const existing = cart.value.find(i => i.id === product.id)
  if (existing) {
    existing.qty++
  } else {
    cart.value.push({ id: product.id, name: product.name, price: product.price, qty: 1 })
  }
  notifyStore.info(`${product.name} ditambahkan ke pesanan`, 'Menu Ditambahkan')
}

const updateQty = (id: string, delta: number) => {
  const idx = cart.value.findIndex(i => i.id === id)
  if (idx > -1) {
    cart.value[idx].qty += delta
    if (cart.value[idx].qty <= 0) {
      const removed = cart.value.splice(idx, 1)
      notifyStore.info(`${removed[0].name} dihapus dari pesanan`, 'Item Dihapus')
    }
  }
}

const clearCart = () => {
  if (cart.value.length === 0) return
  cart.value = []
  notifyStore.info('Keranjang pesanan telah dikosongkan.', 'Keranjang Kosong')
}

const saveOrderToTable = async () => {
  if (cart.value.length === 0) return
  savingOrder.value = true
  try {
    const isDineIn = orderType.value === 'Dine-in'
    const res = await axios.post('/api/v1/pos/orders', {
      customer_name: isDineIn ? `Tamu Meja ${selectedTable.value}` : `Pelanggan ${orderType.value}`,
      table_number: isDineIn ? selectedTable.value : '',
      order_type: orderType.value.toLowerCase().replace('-', '_'),
      items: cart.value.map(i => ({
        product_id: i.id,
        quantity: i.qty,
        unit_price: i.price
      }))
    })

    const queueNum = res.data?.queue_number || '-'
    if (isDineIn) {
      notifyStore.success(`Pesanan Meja ${selectedTable.value} berhasil disimpan! No. Antrian: #${queueNum}`, 'Meja Occupied & Pesanan Terkirim')
    } else {
      notifyStore.success(`Pesanan ${orderType.value} berhasil diproses! No. Antrian: #${queueNum}`, 'Antrian Dapur Dibuat')
    }
    cart.value = []
    await fetchActiveOrders()
    await fetchPOSData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.error || 'Gagal menyimpan pesanan', 'Kesalahan')
  } finally {
    savingOrder.value = false
  }
}

const openPayment = () => {
  if (cart.value.length === 0) return
  showPaymentModal.value = true
}

const handlePaymentSuccess = async () => {
  cart.value = []
  showPaymentModal.value = false
  await fetchActiveOrders()
  await fetchPOSData()
}

// Tab 2 Billing Computed & Handlers
const filteredActiveOrders = computed(() => {
  return activeOrders.value.filter(o => {
    const matchZone = billingFilter.value === 'Semua Tagihan' || o.zone_name === billingFilter.value
    const matchSearch = (o.customer || '').toLowerCase().includes(billingSearch.value.toLowerCase()) ||
                        (o.table_number || '').toLowerCase().includes(billingSearch.value.toLowerCase()) ||
                        (o.order_number || '').toLowerCase().includes(billingSearch.value.toLowerCase())
    return matchZone && matchSearch
  })
})

const printTempBill = (order: any) => {
  notifyStore.info(`Mencetak bill sementara untuk meja ${order.table_number}...`, 'Cetak Bill')
  window.print()
}

const payActiveOrder = (order: any) => {
  selectedOrderForPayment.value = order
  showBillingPaymentModal.value = true
}

const handleBillingPaymentSuccess = async () => {
  showBillingPaymentModal.value = false
  selectedOrderForPayment.value = null
  await fetchActiveOrders()
  await fetchPOSData()
}
</script>
