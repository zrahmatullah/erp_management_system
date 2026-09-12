<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Purchase Order Management</h2>
        <p class="text-xs text-slate-500 mt-1">Pengadaan bahan baku ke supplier, persetujuan PO, dan penerimaan stok gudang (GRN)</p>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="fetchPOList"
          class="p-2.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
          title="Segarkan Data"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>
        <button
          @click="openCreateModal"
          class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>Buat PO Baru</span>
        </button>
      </div>
    </div>

    <!-- 4 Top KPI Stat Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-400 mb-1">Total PO</div>
          <div class="text-2xl font-black text-slate-900">{{ poList.length }}</div>
          <div class="text-[11px] font-bold text-blue-600 mt-1">Rp {{ formatNum(totalAllAmount) }}</div>
        </div>
        <div class="w-10 h-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center">
          <Truck class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-400 mb-1">Menunggu Review</div>
          <div class="text-2xl font-black text-amber-600">{{ countDraft }}</div>
          <div class="text-[11px] font-bold text-amber-600 mt-1">Perlu Persetujuan</div>
        </div>
        <div class="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center">
          <Clock class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-400 mb-1">Dikirim ke Supplier</div>
          <div class="text-2xl font-black text-blue-600">{{ countSent }}</div>
          <div class="text-[11px] font-bold text-blue-600 mt-1">Menunggu Pengiriman</div>
        </div>
        <div class="w-10 h-10 rounded-xl bg-sky-50 text-sky-600 flex items-center justify-center">
          <Package class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-slate-200/80 shadow-sm flex items-start justify-between">
        <div>
          <div class="text-xs font-semibold text-slate-400 mb-1">Barang Diterima</div>
          <div class="text-2xl font-black text-emerald-600">{{ countReceived }}</div>
          <div class="text-[11px] font-bold text-emerald-600 mt-1">Stok Masuk Gudang</div>
        </div>
        <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
          <CheckCircle2 class="w-5 h-5" />
        </div>
      </div>
    </div>

    <!-- Main PO Table Section -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm p-6 space-y-4">
      <!-- Search and Status Filter Tabs -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-100 pb-4">
        <!-- Status Tabs -->
        <div class="flex items-center gap-2 overflow-x-auto pb-1 sm:pb-0">
          <button
            v-for="st in ['all', 'draft', 'sent', 'received', 'rejected']"
            :key="st"
            @click="activeStatusFilter = st"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer capitalize"
            :class="activeStatusFilter === st ? 'bg-slate-900 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'"
          >
            {{ st === 'all' ? 'Semua Status' : st }}
          </button>
        </div>

        <!-- Search Input -->
        <div class="relative max-w-xs w-full">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400">
            <Search class="w-3.5 h-3.5" />
          </span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari PO / Supplier / Catatan..."
            class="w-full pl-8 pr-3 py-1.5 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-none focus:bg-white focus:border-blue-600 transition-all"
          />
        </div>
      </div>

      <!-- Table View -->
      <div class="overflow-x-auto">
        <div v-if="loading" class="py-16 text-center text-slate-400 text-xs flex flex-col items-center justify-center gap-2">
          <RotateCw class="w-6 h-6 animate-spin text-blue-600" />
          <span>Memuat daftar Purchase Order dari database...</span>
        </div>
        <table v-else class="w-full text-left text-xs">
          <thead>
            <tr class="text-slate-500 border-b border-slate-200 uppercase font-bold text-[11px] bg-slate-50/50">
              <th class="py-3 px-4">No. PO</th>
              <th class="py-3 px-4">Supplier</th>
              <th class="py-3 px-4">Tgl Order</th>
              <th class="py-3 px-4 text-center">Items</th>
              <th class="py-3 px-4 text-right">Total Nominal</th>
              <th class="py-3 px-4 text-center">Status</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
            <tr v-for="po in filteredPOList" :key="po.id" class="hover:bg-slate-50/80 transition-colors">
              <td class="py-3.5 px-4 font-mono font-bold text-blue-600">
                {{ po.po_number }}
              </td>
              <td class="py-3.5 px-4">
                <div class="font-bold text-slate-900">{{ po.supplier }}</div>
                <div class="text-[10px] text-slate-400 truncate max-w-xs">{{ po.notes || 'Pengadaan reguler' }}</div>
              </td>
              <td class="py-3.5 px-4 text-slate-500 font-mono">
                {{ po.order_date }}
              </td>
              <td class="py-3.5 px-4 text-center">
                <span class="px-2 py-0.5 bg-slate-100 rounded-md font-bold text-slate-700 text-[11px]">
                  {{ po.total_items }} item
                </span>
              </td>
              <td class="py-3.5 px-4 text-right font-bold text-slate-900">
                Rp {{ formatNum(po.total_amount) }}
              </td>
              <td class="py-3.5 px-4 text-center">
                <span
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                  :class="getStatusBadge(po.status)"
                >
                  {{ po.status }}
                </span>
              </td>
              <td class="py-3.5 px-4 text-right">
                <div class="inline-flex items-center gap-1.5">
                  <button
                    @click="openDetailModal(po.id)"
                    class="px-2.5 py-1 rounded-lg bg-blue-50 hover:bg-blue-100 text-blue-600 text-[11px] font-bold transition-colors cursor-pointer flex items-center gap-1"
                    title="Lihat Detail & Barang"
                  >
                    <Eye class="w-3.5 h-3.5" />
                    <span>Detail</span>
                  </button>

                  <button
                    v-if="po.status === 'sent' || po.status === 'manager_approved' || po.status === 'approved'"
                    @click="quickReceivePO(po)"
                    class="px-2.5 py-1 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white text-[11px] font-bold transition-colors shadow-xs cursor-pointer flex items-center gap-1"
                    title="Terima Barang di Gudang & Tambah Stok"
                  >
                    <Check class="w-3.5 h-3.5" />
                    <span>Terima GRN</span>
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredPOList.length === 0">
              <td colspan="7" class="py-12 text-center text-slate-400 text-xs">
                Tidak ada data Purchase Order sesuai filter atau pencarian.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL: PO DETAIL & ITEMS -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-2xl max-w-3xl w-full p-6 shadow-2xl border border-slate-200 space-y-5 max-h-[90vh] overflow-y-auto">
          <!-- Modal Header -->
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div>
              <div class="flex items-center gap-2">
                <h3 class="text-lg font-black text-slate-900 tracking-tight">Detail Purchase Order {{ currentPODetail.po_number }}</h3>
                <span
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                  :class="getStatusBadge(currentPODetail.status)"
                >
                  {{ currentPODetail.status }}
                </span>
              </div>
              <p class="text-xs text-slate-500 mt-0.5">Dokumen pengadaan resmi inventori bahan baku</p>
            </div>
            <button @click="showDetailModal = false" class="p-1.5 rounded-xl hover:bg-slate-100 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- PO Overview Info -->
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 bg-slate-50/70 p-4 rounded-xl text-xs">
            <div>
              <span class="text-slate-400 block text-[10px] uppercase font-bold">Supplier</span>
              <strong class="text-slate-900">{{ currentPODetail.supplier }}</strong>
            </div>
            <div>
              <span class="text-slate-400 block text-[10px] uppercase font-bold">Tanggal Order</span>
              <strong class="text-slate-800 font-mono">{{ currentPODetail.order_date }}</strong>
            </div>
            <div>
              <span class="text-slate-400 block text-[10px] uppercase font-bold">Estimasi Datang</span>
              <strong class="text-slate-800 font-mono">{{ currentPODetail.expected_date || '-' }}</strong>
            </div>
            <div>
              <span class="text-slate-400 block text-[10px] uppercase font-bold">Total Nominal</span>
              <strong class="text-blue-600 font-black">Rp {{ formatNum(currentPODetail.total_amount) }}</strong>
            </div>
          </div>

          <div v-if="currentPODetail.notes" class="text-xs text-slate-600 bg-blue-50/50 p-3 rounded-xl border border-blue-100">
            <strong>Catatan:</strong> {{ currentPODetail.notes }}
          </div>

          <!-- Items Table -->
          <div class="space-y-2">
            <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">Item Barang yang Dipesan</h4>
            <div class="border border-slate-200 rounded-xl overflow-hidden">
              <table class="w-full text-left text-xs">
                <thead>
                  <tr class="bg-slate-50 text-slate-500 font-bold border-b border-slate-200 text-[11px]">
                    <th class="py-2.5 px-3">Barang</th>
                    <th class="py-2.5 px-3 text-center">Qty Pesan</th>
                    <th class="py-2.5 px-3 text-center">Qty Diterima</th>
                    <th class="py-2.5 px-3 text-right">Harga Satuan</th>
                    <th class="py-2.5 px-3 text-right">Subtotal</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100 text-slate-700">
                  <tr v-for="it in currentPODetail.items" :key="it.id">
                    <td class="py-2.5 px-3 font-semibold text-slate-900">
                      <div>{{ it.name }}</div>
                      <div class="text-[10px] text-slate-400 font-mono">{{ it.sku }} ({{ it.uom }})</div>
                    </td>
                    <td class="py-2.5 px-3 text-center font-bold">{{ it.quantity }} {{ it.uom }}</td>
                    <td class="py-2.5 px-3 text-center font-bold" :class="it.quantity_received >= it.quantity ? 'text-emerald-600' : 'text-slate-500'">
                      {{ it.quantity_received }} {{ it.uom }}
                    </td>
                    <td class="py-2.5 px-3 text-right">Rp {{ formatNum(it.unit_price) }}</td>
                    <td class="py-2.5 px-3 text-right font-bold text-slate-900">Rp {{ formatNum(it.total_price) }}</td>
                  </tr>
                  <tr v-if="!currentPODetail.items || currentPODetail.items.length === 0">
                    <td colspan="5" class="py-6 text-center text-slate-400 text-xs">Tidak ada item dalam PO ini.</td>
                  </tr>
                </tbody>
                <tfoot>
                  <tr class="bg-slate-50/80 border-t border-slate-200 font-black text-slate-900">
                    <td colspan="4" class="py-3 px-3 text-right">Total Keseluruhan:</td>
                    <td class="py-3 px-3 text-right text-blue-600 text-sm">Rp {{ formatNum(currentPODetail.total_amount) }}</td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>

          <!-- Actions Footer -->
          <div class="flex items-center justify-between pt-3 border-t border-slate-100">
            <button
              v-if="currentPODetail.status !== 'rejected' && currentPODetail.status !== 'received'"
              @click="updatePOStatus(currentPODetail.id, 'rejected')"
              class="px-3.5 py-2 border border-rose-300 text-rose-600 hover:bg-rose-50 text-xs font-bold rounded-xl transition-colors cursor-pointer"
            >
              Tolak PO
            </button>
            <div v-else></div>

            <div class="flex items-center gap-2">
              <button
                v-if="currentPODetail.status === 'draft' || currentPODetail.status === 'submitted'"
                @click="updatePOStatus(currentPODetail.id, 'sent')"
                class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer"
              >
                Setujui & Kirim ke Supplier
              </button>

              <button
                v-if="currentPODetail.status === 'sent' || currentPODetail.status === 'approved' || currentPODetail.status === 'manager_approved'"
                @click="updatePOStatus(currentPODetail.id, 'received')"
                class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold rounded-xl shadow-md shadow-emerald-600/30 transition-all cursor-pointer flex items-center gap-1.5"
              >
                <CheckCircle2 class="w-4 h-4" />
                <span>Terima Barang (GRN) & Tambah Stok</span>
              </button>

              <button
                @click="showDetailModal = false"
                class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-bold rounded-xl transition-colors cursor-pointer"
              >
                Tutup
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- MODAL: BUAT PURCHASE ORDER BARU -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-2xl max-w-3xl w-full p-6 shadow-2xl border border-slate-200 space-y-5 max-h-[90vh] overflow-y-auto">
          <!-- Header -->
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div>
              <h3 class="text-lg font-black text-slate-900 tracking-tight">Buat Purchase Order Baru</h3>
              <p class="text-xs text-slate-500 mt-0.5">Ajukan permintaan pengadaan bahan baku ke vendor supplier</p>
            </div>
            <button @click="showCreateModal = false" class="p-1.5 rounded-xl hover:bg-slate-100 text-slate-400 hover:text-slate-600 transition-colors cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Form Fields -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Pilih Supplier <span class="text-rose-500">*</span></label>
              <select
                v-model="newPO.supplier_id"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium"
              >
                <option value="">-- Pilih Vendor Supplier --</option>
                <option v-for="sup in supplierOptions" :key="sup.id" :value="sup.id">
                  {{ sup.name }}
                </option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Estimasi Tanggal Pengiriman</label>
              <input
                v-model="newPO.expected_date"
                type="date"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium"
              />
            </div>

            <div class="sm:col-span-2">
              <label class="block text-xs font-bold text-slate-700 mb-1">Catatan Pengadaan</label>
              <input
                v-model="newPO.notes"
                type="text"
                placeholder="Contoh: Pengadaan bahan baku operasional mingguan cabang Senopati"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none"
              />
            </div>
          </div>

          <!-- Item Selection Section -->
          <div class="space-y-3 pt-2">
            <div class="flex items-center justify-between">
              <label class="text-xs font-bold text-slate-800 uppercase tracking-wider">Daftar Barang Pesanan</label>
              <button
                type="button"
                @click="addPOItem"
                class="px-3 py-1 bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-bold rounded-xl transition-colors cursor-pointer flex items-center gap-1"
              >
                <Plus class="w-3.5 h-3.5" />
                <span>Tambah Baris Barang</span>
              </button>
            </div>

            <div class="border border-slate-200 rounded-xl overflow-hidden">
              <table class="w-full text-left text-xs">
                <thead>
                  <tr class="bg-slate-50 text-slate-500 font-bold border-b border-slate-200 text-[11px]">
                    <th class="py-2.5 px-3">Bahan Baku</th>
                    <th class="py-2.5 px-3 text-center w-28">Jumlah (Qty)</th>
                    <th class="py-2.5 px-3 text-right w-36">Harga Satuan (Rp)</th>
                    <th class="py-2.5 px-3 text-right w-36">Subtotal</th>
                    <th class="py-2.5 px-3 text-center w-12">Hapus</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100 text-slate-700 font-medium">
                  <tr v-for="(it, idx) in newPO.items" :key="idx">
                    <td class="py-2.5 px-3">
                      <select
                        v-model="it.inventory_item_id"
                        @change="onSelectItem(it)"
                        class="w-full px-2 py-1.5 text-xs bg-slate-50 border border-slate-200 rounded-lg outline-none"
                      >
                        <option value="">-- Pilih Bahan Baku --</option>
                        <option v-for="stock in availableStocks" :key="stock.id" :value="stock.id">
                          {{ stock.name }} (Stok saat ini: {{ stock.stock }} {{ stock.uom }})
                        </option>
                      </select>
                    </td>
                    <td class="py-2.5 px-3 text-center">
                      <input
                        v-model.number="it.quantity"
                        type="number"
                        min="1"
                        class="w-20 px-2 py-1 text-center font-bold bg-slate-50 border border-slate-200 rounded-lg outline-none focus:border-blue-600"
                      />
                    </td>
                    <td class="py-2.5 px-3 text-right">
                      <input
                        v-model.number="it.unit_price"
                        type="number"
                        min="0"
                        class="w-28 px-2 py-1 text-right font-bold bg-slate-50 border border-slate-200 rounded-lg outline-none focus:border-blue-600"
                      />
                    </td>
                    <td class="py-2.5 px-3 text-right font-bold text-slate-900">
                      Rp {{ formatNum(it.quantity * it.unit_price) }}
                    </td>
                    <td class="py-2.5 px-3 text-center">
                      <button
                        v-if="newPO.items.length > 1"
                        type="button"
                        @click="removePOItem(idx)"
                        class="p-1 text-rose-500 hover:text-rose-700 hover:bg-rose-50 rounded-md cursor-pointer transition-colors"
                      >
                        <Trash2 class="w-4 h-4" />
                      </button>
                    </td>
                  </tr>
                </tbody>
                <tfoot>
                  <tr class="bg-slate-50/80 border-t border-slate-200 font-black text-slate-900">
                    <td colspan="3" class="py-3 px-3 text-right">Total Perkiraan PO:</td>
                    <td class="py-3 px-3 text-right text-blue-600 text-sm">
                      Rp {{ formatNum(calculatedNewPOTotal) }}
                    </td>
                    <td></td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>

          <!-- Modal Actions -->
          <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer"
            >
              Batal
            </button>
            <button
              type="button"
              :disabled="submitting"
              @click="submitNewPO"
              class="px-5 py-2 text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 active:bg-blue-800 rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
            >
              <span>{{ submitting ? 'Menyimpan...' : 'Simpan & Ajukan PO' }}</span>
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
  Plus,
  Search,
  RotateCw,
  Eye,
  Check,
  CheckCircle2,
  Clock,
  Package,
  Truck,
  Trash2,
  X
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const loading = ref(false)
const submitting = ref(false)
const searchQuery = ref('')
const activeStatusFilter = ref('all')

const poList = ref<any[]>([])
const supplierOptions = ref<any[]>([])
const availableStocks = ref<any[]>([])

const showDetailModal = ref(false)
const currentPODetail = ref<any>({})

const showCreateModal = ref(false)
const newPO = ref({
  supplier_id: '',
  notes: '',
  expected_date: '',
  items: [
    { inventory_item_id: '', quantity: 10, unit_price: 50000 }
  ]
})

const formatNum = (num: number) => {
  return Number(num || 0).toLocaleString('id-ID')
}

const fetchPOList = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/inventory/purchase-orders')
    poList.value = Array.isArray(res.data) ? res.data : (res.data?.data || [])
  } catch (err: any) {
    console.error('Failed to load purchase orders:', err)
    notifyStore.error('Gagal mengambil daftar Purchase Order dari database', 'Koneksi Gagal')
  } finally {
    loading.value = false
  }
}

const fetchMasterData = async () => {
  try {
    const [supRes, stockRes] = await Promise.all([
      axios.get('/api/v1/master/suppliers'),
      axios.get('/api/v1/inventory/stocks')
    ])
    supplierOptions.value = Array.isArray(supRes.data) ? supRes.data : (supRes.data?.data || [])
    availableStocks.value = Array.isArray(stockRes.data) ? stockRes.data : (stockRes.data?.data || [])
  } catch (err) {
    console.error('Failed to load suppliers/stocks:', err)
  }
}

onMounted(() => {
  fetchPOList()
  fetchMasterData()
})

const filteredPOList = computed(() => {
  return poList.value.filter(po => {
    const matchStatus = activeStatusFilter.value === 'all' || po.status === activeStatusFilter.value
    const q = searchQuery.value.toLowerCase()
    const matchQuery = !q ||
      (po.po_number && po.po_number.toLowerCase().includes(q)) ||
      (po.supplier && po.supplier.toLowerCase().includes(q)) ||
      (po.notes && po.notes.toLowerCase().includes(q))
    return matchStatus && matchQuery
  })
})

const totalAllAmount = computed(() => {
  return poList.value.reduce((acc, po) => acc + (Number(po.total_amount) || 0), 0)
})

const countDraft = computed(() => poList.value.filter(po => po.status === 'draft' || po.status === 'submitted').length)
const countSent = computed(() => poList.value.filter(po => po.status === 'sent' || po.status === 'approved' || po.status === 'manager_approved').length)
const countReceived = computed(() => poList.value.filter(po => po.status === 'received').length)

const getStatusBadge = (status: string) => {
  switch (status) {
    case 'received':
      return 'bg-emerald-100 text-emerald-800'
    case 'sent':
    case 'approved':
    case 'manager_approved':
      return 'bg-blue-100 text-blue-800'
    case 'draft':
    case 'submitted':
      return 'bg-amber-100 text-amber-800'
    case 'rejected':
    case 'cancelled':
      return 'bg-rose-100 text-rose-800'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

const openDetailModal = async (poId: string) => {
  try {
    const res = await axios.get(`/api/v1/inventory/purchase-orders/${poId}`)
    if (res.data?.data) {
      currentPODetail.value = res.data.data
      showDetailModal.value = true
    }
  } catch (err: any) {
    console.error('Failed to load PO detail:', err)
    notifyStore.error('Gagal mengambil detail Purchase Order', 'Gagal Memuat')
  }
}

const updatePOStatus = async (poId: string, newStatus: string) => {
  try {
    const res = await axios.put(`/api/v1/inventory/purchase-orders/${poId}/status`, {
      status: newStatus
    })
    notifyStore.success(res.data?.message || `Status PO berhasil diubah menjadi ${newStatus}`, 'Sukses')
    showDetailModal.value = false
    await fetchPOList()
  } catch (err: any) {
    console.error('Failed to update PO status:', err)
    notifyStore.error(err.response?.data?.message || 'Gagal mengubah status PO', 'Error')
  }
}

const quickReceivePO = async (po: any) => {
  await updatePOStatus(po.id, 'received')
}

const openCreateModal = () => {
  newPO.value = {
    supplier_id: supplierOptions.value[0]?.id || '',
    notes: '',
    expected_date: new Date(Date.now() + 7 * 86400000).toISOString().split('T')[0],
    items: [
      {
        inventory_item_id: availableStocks.value[0]?.id || '',
        quantity: 10,
        unit_price: Number(availableStocks.value[0]?.cost || 25000)
      }
    ]
  }
  showCreateModal.value = true
}

const onSelectItem = (item: any) => {
  const stock = availableStocks.value.find(s => s.id === item.inventory_item_id)
  if (stock) {
    item.unit_price = Number(stock.cost || 25000)
  }
}

const addPOItem = () => {
  const defaultStock = availableStocks.value[0]
  newPO.value.items.push({
    inventory_item_id: defaultStock?.id || '',
    quantity: 10,
    unit_price: Number(defaultStock?.cost || 25000)
  })
}

const removePOItem = (idx: number) => {
  if (newPO.value.items.length > 1) {
    newPO.value.items.splice(idx, 1)
  }
}

const calculatedNewPOTotal = computed(() => {
  return newPO.value.items.reduce((acc, it) => acc + ((Number(it.quantity) || 0) * (Number(it.unit_price) || 0)), 0)
})

const submitNewPO = async () => {
  if (!newPO.value.supplier_id) {
    notifyStore.warning('Pilih supplier terlebih dahulu!', 'Validasi Input')
    return
  }
  if (newPO.value.items.some(i => !i.inventory_item_id || i.quantity <= 0)) {
    notifyStore.warning('Pastikan semua baris barang telah dipilih dan kuantitas valid!', 'Validasi Input')
    return
  }

  submitting.value = true
  try {
    const res = await axios.post('/api/v1/inventory/purchase-orders', newPO.value)
    notifyStore.success(res.data?.message || 'Purchase Order baru berhasil dibuat!', 'PO Dibuat')
    showCreateModal.value = false
    await fetchPOList()
  } catch (err: any) {
    console.error('Failed to create PO:', err)
    notifyStore.error(err.response?.data?.message || 'Gagal membuat Purchase Order', 'Gagal Simpan')
  } finally {
    submitting.value = false
  }
}
</script>
