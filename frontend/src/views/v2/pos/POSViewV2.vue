<template>
  <div class="pos-v2 space-y-4 animate-in fade-in duration-200">
    <!-- Top Mode Switcher Bar -->
    <div 
      class="p-2.5 rounded-2xl border flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-xs"
      style="background-color: var(--bg-primary); border-color: var(--border-color);"
    >
      <div class="flex items-center gap-1.5 overflow-x-auto">
        <button
          type="button"
          @click="activePosTab = 'order'"
          class="flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer select-none"
          :class="activePosTab === 'order' 
            ? 'bg-blue-600 text-white shadow-xs' 
            : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100 hover:bg-slate-100 dark:hover:bg-slate-800'"
        >
          <Utensils class="w-4 h-4" />
          <span>Katalog & Pesanan Baru</span>
        </button>

        <button
          type="button"
          @click="activePosTab = 'billing'"
          class="flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer select-none relative"
          :class="activePosTab === 'billing' 
            ? 'bg-blue-600 text-white shadow-xs' 
            : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100 hover:bg-slate-100 dark:hover:bg-slate-800'"
        >
          <CreditCard class="w-4 h-4" />
          <span>Pembayaran Tagihan Meja</span>
          <span
            v-if="activeOrders.length > 0"
            class="px-2 py-0.2 rounded-full text-[10px] font-bold"
            :class="activePosTab === 'billing' ? 'bg-white text-blue-700' : 'bg-rose-500 text-white'"
          >
            {{ activeOrders.length }}
          </span>
        </button>

        <button
          type="button"
          @click="activePosTab = 'history'"
          class="flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-semibold transition-all cursor-pointer select-none"
          :class="activePosTab === 'history' 
            ? 'bg-blue-600 text-white shadow-xs' 
            : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-100 hover:bg-slate-100 dark:hover:bg-slate-800'"
        >
          <ReceiptText class="w-4 h-4" />
          <span>Riwayat Transaksi & Struk</span>
        </button>
      </div>

      <!-- Cashier Meta Pill -->
      <div class="flex items-center gap-3 text-xs text-slate-500 dark:text-slate-400 px-2">
        <div class="flex items-center gap-1.5">
          <User class="w-3.5 h-3.5 text-blue-600 dark:text-blue-400" />
          <span>Kasir: <strong class="text-slate-900 dark:text-slate-100">{{ cashierName }}</strong></span>
        </div>
        <span class="text-slate-300 dark:text-slate-700">|</span>
        <div class="flex items-center gap-1.5">
          <Clock class="w-3.5 h-3.5" />
          <span>Shift Pagi</span>
        </div>
      </div>
    </div>

    <!-- TAB 1: ORDER BARU & KATALOG -->
    <div v-if="activePosTab === 'order'" class="flex flex-col lg:flex-row gap-5 h-[calc(100vh-10.5rem)]">
      <!-- Left: Products Catalog Grid -->
      <div class="flex-1 flex flex-col min-w-0 space-y-3">
        <!-- Search & Filter Controls -->
        <div 
          class="p-4 rounded-2xl border space-y-3 shrink-0"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">
              Pilih menu kafe untuk ditambahkan ke keranjang pesanan
            </span>

            <div class="relative w-full sm:w-64">
              <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5" />
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Cari menu kopi, makanan..."
                class="w-full pl-8.5 pr-3 py-1.5 text-xs rounded-xl border outline-none transition-all focus:ring-1 focus:ring-blue-500"
                style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
              />
            </div>
          </div>

          <!-- Category Pills -->
          <div class="flex items-center gap-2 overflow-x-auto pb-1 text-xs select-none">
            <button
              v-for="cat in availableCategories"
              :key="cat"
              type="button"
              @click="selectedCategory = cat"
              class="px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all cursor-pointer shrink-0 border"
              :class="selectedCategory === cat 
                ? 'bg-blue-600 text-white border-blue-600 shadow-xs' 
                : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200 border-slate-200/80 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800'"
            >
              {{ cat }}
            </button>
          </div>
        </div>

        <!-- Products Scrollable Grid -->
        <div class="flex-1 overflow-y-auto pr-1">
          <div v-if="filteredProducts.length === 0" class="p-12 text-center text-xs text-slate-400 border border-dashed rounded-2xl" style="border-color: var(--border-color);">
            Menu dengan kata kunci "{{ searchQuery }}" tidak ditemukan.
          </div>

          <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-3.5">
            <div
              v-for="prod in filteredProducts"
              :key="prod.id"
              class="p-3.5 rounded-2xl border transition-all duration-200 hover:shadow-xs group flex flex-col justify-between"
              style="background-color: var(--bg-primary); border-color: var(--border-color);"
            >
              <div>
                <!-- Image / Icon Header -->
                <div class="h-28 rounded-xl bg-slate-100 dark:bg-slate-800/80 flex items-center justify-center text-slate-400 mb-3 relative overflow-hidden">
                  <img v-if="prod.image_url" :src="prod.image_url" :alt="prod.name" class="w-full h-full object-cover" />
                  <Utensils v-else class="w-8 h-8 text-slate-300 dark:text-slate-600" />

                  <!-- Status Stock Badge -->
                  <span 
                    class="absolute top-2 left-2 text-[10px] font-bold px-2 py-0.5 rounded-md"
                    :class="prod.stock_status === 'in_stock' 
                      ? 'bg-emerald-500 text-white' 
                      : prod.stock_status === 'low_stock' 
                      ? 'bg-amber-500 text-white' 
                      : 'bg-rose-500 text-white'"
                  >
                    {{ prod.stock_status === 'in_stock' ? `${prod.calculated_portions} Porsi` : prod.stock_status === 'low_stock' ? 'Stok Tipis' : 'Habis' }}
                  </span>
                </div>

                <div class="font-bold text-xs text-slate-900 dark:text-slate-100 truncate" :title="prod.name">
                  {{ prod.name }}
                </div>
                <div class="text-[11px] text-slate-400 mt-0.5 truncate">{{ prod.category }}</div>
              </div>

              <div class="mt-3 pt-2.5 border-t flex items-center justify-between" style="border-color: var(--border-color);">
                <div class="font-mono font-bold text-xs text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(prod.price) }}
                </div>

                <button
                  type="button"
                  v-if="prod.stock_status !== 'out_of_stock'"
                  @click="addToCart(prod)"
                  class="px-2.5 py-1 rounded-lg text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white transition-all cursor-pointer shadow-xs flex items-center gap-1"
                >
                  <Plus class="w-3.5 h-3.5" />
                  <span>Tambah</span>
                </button>
                <button
                  type="button"
                  v-else
                  @click="openDraftPOModal(prod)"
                  class="px-2 py-1 rounded-lg text-[10px] font-bold bg-amber-50 text-amber-700 border border-amber-300 hover:bg-amber-100 dark:bg-amber-950/40 dark:text-amber-300 dark:border-amber-800 transition-all cursor-pointer"
                >
                  Minta PO
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Ergonomic Cart & Checkout Panel -->
      <div 
        class="w-full lg:w-96 rounded-2xl border shadow-xs flex flex-col shrink-0 overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <!-- Cart Header -->
        <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
          <div class="flex items-center gap-2">
            <ShoppingCart class="w-4 h-4 text-blue-600 dark:text-blue-400" />
            <h3 class="font-bold text-sm text-slate-900 dark:text-slate-100">Keranjang Kasir</h3>
            <span class="text-xs font-mono px-2 py-0.2 rounded-full bg-blue-100 dark:bg-blue-900/60 text-blue-700 dark:text-blue-300 font-bold">
              {{ cartItems.length }} Item
            </span>
          </div>

          <button 
            v-if="cartItems.length > 0"
            type="button"
            @click="clearCart"
            class="text-[11px] font-semibold text-rose-600 hover:underline cursor-pointer"
          >
            Kosongkan
          </button>
        </div>

        <!-- Dine-in / Takeaway Switch & Table Selector -->
        <div class="p-3 border-b space-y-2.5" style="border-color: var(--border-color); background-color: var(--bg-content);">
          <div class="grid grid-cols-2 gap-2 text-xs">
            <button
              type="button"
              @click="orderType = 'dine_in'"
              class="py-1.5 rounded-xl font-semibold border transition-all cursor-pointer text-center"
              :class="orderType === 'dine_in' 
                ? 'bg-blue-600 text-white border-blue-600 shadow-xs' 
                : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 border-slate-200 dark:border-slate-800'"
            >
              🍽️ Dine-In
            </button>
            <button
              type="button"
              @click="orderType = 'takeaway'"
              class="py-1.5 rounded-xl font-semibold border transition-all cursor-pointer text-center"
              :class="orderType === 'takeaway' 
                ? 'bg-blue-600 text-white border-blue-600 shadow-xs' 
                : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 border-slate-200 dark:border-slate-800'"
            >
              🛍️ Takeaway
            </button>
          </div>

          <div v-if="orderType === 'dine_in'" class="flex items-center gap-2">
            <label class="text-xs font-semibold text-slate-600 dark:text-slate-400 shrink-0">Nomor Meja:</label>
            <select
              v-model="selectedTableNumber"
              class="flex-1 px-3 py-1.5 text-xs rounded-xl border outline-none"
              style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
            >
              <option value="T-01">Meja T-01 (2 Kursi)</option>
              <option value="T-02">Meja T-02 (2 Kursi)</option>
              <option value="T-03">Meja T-03 (2 Kursi)</option>
              <option value="T-04">Meja T-04 (4 Kursi)</option>
              <option value="T-05">Meja T-05 (4 Kursi)</option>
              <option value="T-06">Meja T-06 (6 Kursi)</option>
              <option value="T-07">Meja T-07 (8 Kursi)</option>
            </select>
          </div>
        </div>

        <!-- Scrollable Cart Item Rows -->
        <div class="flex-1 overflow-y-auto p-3 space-y-2.5">
          <div v-if="cartItems.length === 0" class="h-full flex flex-col items-center justify-center text-center p-6 text-slate-400">
            <ShoppingCart class="w-10 h-10 mb-2 opacity-30" />
            <p class="text-xs font-medium">Keranjang masih kosong</p>
            <p class="text-[11px] text-slate-400 mt-0.5">Pilih menu dari katalog di sebelah kiri</p>
          </div>

          <div 
            v-for="(item, idx) in cartItems" 
            :key="idx"
            class="p-3 rounded-xl border space-y-2 text-xs"
            style="background-color: var(--bg-content); border-color: var(--border-color);"
          >
            <div class="flex items-start justify-between gap-2">
              <div class="min-w-0">
                <div class="font-bold text-slate-900 dark:text-slate-100 truncate">{{ item.product.name }}</div>
                <div class="text-[11px] font-mono text-slate-400">Rp {{ formatNum(item.product.price) }} / porsi</div>
              </div>
              <div class="font-mono font-bold text-slate-900 dark:text-slate-100">
                Rp {{ formatNum(item.product.price * item.quantity) }}
              </div>
            </div>

            <!-- Notes & Quantity Adjuster -->
            <div class="flex items-center justify-between gap-2 pt-1">
              <input
                v-model="item.notes"
                type="text"
                placeholder="Catatan (less ice, less sugar)..."
                class="flex-1 px-2.5 py-1 text-[11px] rounded-lg border outline-none"
                style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
              />

              <div class="flex items-center gap-1.5 shrink-0">
                <button
                  type="button"
                  @click="decreaseQty(idx)"
                  class="w-6 h-6 rounded-md border flex items-center justify-center text-slate-600 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700 cursor-pointer"
                >
                  <Minus class="w-3 h-3" />
                </button>
                <span class="w-6 text-center font-bold font-mono">{{ item.quantity }}</span>
                <button
                  type="button"
                  @click="increaseQty(idx)"
                  class="w-6 h-6 rounded-md border flex items-center justify-center text-slate-600 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700 cursor-pointer"
                >
                  <Plus class="w-3 h-3" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Checkout Summary & Action Buttons -->
        <div class="p-4 border-t space-y-3" style="border-color: var(--border-color); background-color: var(--bg-primary);">
          <div class="space-y-1.5 text-xs">
            <div class="flex justify-between text-slate-500">
              <span>Subtotal:</span>
              <span class="font-mono">Rp {{ formatNum(cartSubtotal) }}</span>
            </div>
            <div class="flex justify-between text-slate-500">
              <span>Pajak Resto (PB1 10%):</span>
              <span class="font-mono">Rp {{ formatNum(cartTax) }}</span>
            </div>
            <div class="flex justify-between font-bold text-sm text-slate-900 dark:text-slate-50 pt-1 border-t" style="border-color: var(--border-color);">
              <span>Total Tagihan:</span>
              <span class="font-mono text-blue-600 dark:text-blue-400">Rp {{ formatNum(cartTotal) }}</span>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2 pt-1">
            <button
              type="button"
              @click="sendToKitchen"
              :disabled="cartItems.length === 0"
              class="py-2.5 px-3 rounded-xl text-xs font-semibold border flex items-center justify-center gap-1.5 transition-all cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-100 dark:hover:bg-slate-800"
              style="border-color: var(--border-color); color: var(--text-primary);"
            >
              <Send class="w-3.5 h-3.5 text-blue-600" />
              <span>Kirim KDS</span>
            </button>

            <button
              type="button"
              @click="openCheckoutPayment"
              :disabled="cartItems.length === 0"
              class="py-2.5 px-3 rounded-xl text-xs font-bold bg-blue-600 hover:bg-blue-700 text-white flex items-center justify-center gap-1.5 transition-all cursor-pointer shadow-xs disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <CreditCard class="w-3.5 h-3.5" />
              <span>Bayar Sekarang</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: PEMBAYARAN TAGIHAN MEJA (BILLING) -->
    <div v-else-if="activePosTab === 'billing'" class="space-y-4">
      <div 
        class="p-6 rounded-2xl border shadow-xs"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Daftar Tagihan Meja Aktif</h3>
            <p class="text-xs text-slate-400">Pilih meja yang siap melakukan pembayaran atau cetak pre-bill struk</p>
          </div>
          <button
            type="button"
            @click="fetchActiveOrders"
            class="px-3 py-1.5 rounded-xl border text-xs font-semibold flex items-center gap-1.5 hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer"
            style="border-color: var(--border-color);"
          >
            <RotateCw class="w-3.5 h-3.5" />
            <span>Perbarui Tagihan</span>
          </button>
        </div>

        <div v-if="activeOrders.length === 0" class="p-12 text-center text-xs text-slate-400 border border-dashed rounded-2xl" style="border-color: var(--border-color);">
          Tidak ada tagihan meja yang pending saat ini. Semua pesanan telah diselesaikan.
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="ord in activeOrders"
            :key="ord.id"
            class="p-4 rounded-xl border flex flex-col justify-between space-y-3"
            style="background-color: var(--bg-content); border-color: var(--border-color);"
          >
            <div>
              <div class="flex items-center justify-between mb-1">
                <span class="font-bold text-sm text-slate-900 dark:text-slate-100">{{ ord.table_number }}</span>
                <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-amber-100 text-amber-700">Belum Bayar</span>
              </div>
              <div class="text-[11px] font-mono text-slate-400">{{ ord.order_number }} • {{ ord.customer_name || 'Pelanggan Walk-in' }}</div>
              <div class="text-xs text-slate-600 dark:text-slate-300 mt-2 font-medium">
                {{ ord.items_count }} Item pesanan
              </div>
            </div>

            <div class="pt-3 border-t flex items-center justify-between" style="border-color: var(--border-color);">
              <div>
                <span class="text-[10px] text-slate-400 block">Total:</span>
                <span class="font-mono font-bold text-sm text-blue-600 dark:text-blue-400">
                  Rp {{ formatNum(ord.total_amount) }}
                </span>
              </div>

              <button
                type="button"
                @click="payActiveOrder(ord)"
                class="px-4 py-2 rounded-xl text-xs font-bold bg-emerald-600 hover:bg-emerald-700 text-white cursor-pointer shadow-xs flex items-center gap-1.5"
              >
                <CreditCard class="w-3.5 h-3.5" />
                <span>Proses Bayar</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 3: RIWAYAT TRANSAKSI (REUSED V1 COMPONENT) -->
    <div v-else class="space-y-4">
      <TransactionHistoryView />
    </div>

    <!-- Payment Modal (Reused) -->
    <PaymentModal
      v-if="showPaymentModal"
      :show="showPaymentModal"
      :order-data="paymentOrderPayload"
      @close="showPaymentModal = false"
      @payment-completed="handlePaymentCompleted"
    />

    <!-- Draft PO Modal for Out of Stock -->
    <Teleport to="body">
      <div v-if="showDraftPOModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/50 backdrop-blur-sm animate-in fade-in">
        <div 
          class="w-full max-w-md p-6 rounded-3xl border shadow-xl space-y-4 animate-in zoom-in-95 duration-150"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between pb-3 border-b" style="border-color: var(--border-color);">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Pengajuan Restok Menu</h3>
              <p class="text-xs text-slate-400">Kirim permintaan restok bahan ke Tim Gudang</p>
            </div>
            <button @click="showDraftPOModal = false" class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="space-y-3 text-xs">
            <div class="p-3.5 rounded-xl border" style="background-color: var(--bg-content); border-color: var(--border-color);">
              <div class="font-bold text-sm text-slate-900 dark:text-slate-100">{{ selectedProductForPO?.name }}</div>
              <div class="text-[11px] text-rose-600 font-semibold mt-0.5">Stok Saat Ini: 0 Porsi (Habis)</div>
            </div>

            <div>
              <label class="block font-semibold mb-1 text-slate-700 dark:text-slate-300">Target Porsi yang Dibutuhkan</label>
              <input
                v-model.number="draftPOTargetPortions"
                type="number"
                min="5"
                class="w-full px-3 py-2 rounded-xl border outline-none font-bold"
                style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
              />
            </div>
          </div>

          <div class="pt-3 border-t flex justify-end gap-2.5" style="border-color: var(--border-color);">
            <button 
              type="button" 
              @click="showDraftPOModal = false" 
              class="px-4 py-2 rounded-xl border text-xs font-semibold hover:bg-slate-100 dark:hover:bg-slate-800 cursor-pointer"
              style="border-color: var(--border-color);"
            >
              Batal
            </button>
            <button 
              type="button" 
              @click="submitDraftPO"
              class="px-5 py-2 rounded-xl text-xs font-bold bg-amber-500 hover:bg-amber-600 text-white cursor-pointer shadow-xs"
            >
              Kirim ke Gudang
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import PaymentModal from '@/views/pos/PaymentModal.vue'
import TransactionHistoryView from '@/views/pos/TransactionHistoryView.vue'
import {
  Utensils,
  CreditCard,
  ReceiptText,
  User,
  Clock,
  Search,
  ShoppingCart,
  Plus,
  Minus,
  Send,
  RotateCw,
  X
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth.store'
import { useNotificationStore } from '@/stores/notification.store'

const authStore = useAuthStore()
const notifyStore = useNotificationStore()

const activePosTab = ref<'order' | 'billing' | 'history'>('order')
const searchQuery = ref('')
const selectedCategory = ref('Semua')
const orderType = ref<'dine_in' | 'takeaway'>('dine_in')
const selectedTableNumber = ref('T-01')

const cashierName = computed(() => authStore.fullName || authStore.user?.firstName || 'Jane Cashier')

const availableCategories = ['Semua', 'Espresso Based', 'Manual Brew', 'Non-Coffee', 'Pastry & Bakery', 'Main Course']

interface ProductItem {
  id: string
  name: string
  category: string
  price: number
  stock_status: 'in_stock' | 'low_stock' | 'out_of_stock'
  calculated_portions: number
  image_url?: string
}

const mockProducts = ref<ProductItem[]>([
  { id: '1', name: 'Espresso Single Origin', category: 'Espresso Based', price: 25000, stock_status: 'in_stock', calculated_portions: 45 },
  { id: '2', name: 'Americano Ice / Hot', category: 'Espresso Based', price: 28000, stock_status: 'in_stock', calculated_portions: 38 },
  { id: '3', name: 'Caffe Latte Oatmilk', category: 'Espresso Based', price: 38000, stock_status: 'in_stock', calculated_portions: 22 },
  { id: '4', name: 'Cappuccino Cinnamon', category: 'Espresso Based', price: 35000, stock_status: 'in_stock', calculated_portions: 18 },
  { id: '5', name: 'V60 Gayo Natural', category: 'Manual Brew', price: 35000, stock_status: 'in_stock', calculated_portions: 12 },
  { id: '6', name: 'Japanese Drip Iced Coffee', category: 'Manual Brew', price: 38000, stock_status: 'low_stock', calculated_portions: 4 },
  { id: '7', name: 'Matcha Latte Uji Kyoto', category: 'Non-Coffee', price: 40000, stock_status: 'in_stock', calculated_portions: 15 },
  { id: '8', name: 'Croissant Butter French', category: 'Pastry & Bakery', price: 32000, stock_status: 'low_stock', calculated_portions: 3 },
  { id: '9', name: 'Truffle Fries Parmesan', category: 'Main Course', price: 42000, stock_status: 'in_stock', calculated_portions: 16 },
  { id: '10', name: 'Avocado Coffee Special', category: 'Espresso Based', price: 42000, stock_status: 'out_of_stock', calculated_portions: 0 }
])

const filteredProducts = computed(() => {
  return mockProducts.value.filter(p => {
    const matchCat = selectedCategory.value === 'Semua' || p.category === selectedCategory.value
    const matchQuery = !searchQuery.value.trim() || p.name.toLowerCase().includes(searchQuery.value.toLowerCase().trim())
    return matchCat && matchQuery
  })
})

interface CartItem {
  product: ProductItem
  quantity: number
  notes: string
}

const cartItems = ref<CartItem[]>([])

const addToCart = (product: ProductItem) => {
  const existing = cartItems.value.find(c => c.product.id === product.id)
  if (existing) {
    existing.quantity++
  } else {
    cartItems.value.push({ product, quantity: 1, notes: '' })
  }
}

const increaseQty = (idx: number) => {
  cartItems.value[idx].quantity++
}

const decreaseQty = (idx: number) => {
  if (cartItems.value[idx].quantity > 1) {
    cartItems.value[idx].quantity--
  } else {
    cartItems.value.splice(idx, 1)
  }
}

const clearCart = () => {
  cartItems.value = []
}

const cartSubtotal = computed(() => {
  return cartItems.value.reduce((acc, item) => acc + (item.product.price * item.quantity), 0)
})

const cartTax = computed(() => Math.round(cartSubtotal.value * 0.1))
const cartTotal = computed(() => cartSubtotal.value + cartTax.value)

// Billing state
const activeOrders = ref([
  { id: 'ord-01', order_number: 'ORD-2026-0928', table_number: 'Meja T-04', customer_name: 'Pak Hendra', items_count: 3, total_amount: 114400 },
  { id: 'ord-02', order_number: 'ORD-2026-0929', table_number: 'Meja T-05', customer_name: 'Ibu Ratna', items_count: 5, total_amount: 185900 }
])

const fetchActiveOrders = () => {
  notifyStore.info('Daftar tagihan meja aktif diperbarui', 'POS Kasir')
}

// Payment modal trigger
const showPaymentModal = ref(false)
const paymentOrderPayload = ref<any>(null)

const openCheckoutPayment = () => {
  paymentOrderPayload.value = {
    order_id: 'ORD-NEW-' + Date.now(),
    total_amount: cartTotal.value,
    items: cartItems.value,
    table_number: orderType.value === 'dine_in' ? selectedTableNumber.value : 'Takeaway'
  }
  showPaymentModal.value = true
}

const payActiveOrder = (ord: any) => {
  paymentOrderPayload.value = ord
  showPaymentModal.value = true
}

const handlePaymentCompleted = () => {
  showPaymentModal.value = false
  clearCart()
  notifyStore.success('Pembayaran kasir selesai! Struk tercetak.', 'Transaksi Sukses')
}

const sendToKitchen = () => {
  notifyStore.success(`Pesanan ${selectedTableNumber.value} berhasil dikirim ke antrean KDS Barista!`, 'Tiket KDS')
}

// Draft PO Modal
const showDraftPOModal = ref(false)
const selectedProductForPO = ref<ProductItem | null>(null)
const draftPOTargetPortions = ref(50)

const openDraftPOModal = (prod: ProductItem) => {
  selectedProductForPO.value = prod
  draftPOTargetPortions.value = 50
  showDraftPOModal.value = true
}

const submitDraftPO = () => {
  showDraftPOModal.value = false
  notifyStore.success(`Permintaan restok untuk "${selectedProductForPO.value?.name}" sebanyak ${draftPOTargetPortions.value} porsi telah diteruskan ke Tim Gudang.`, 'Draft PO Gudang')
}

const formatNum = (val: number | string) => Number(val || 0).toLocaleString('id-ID')
</script>

