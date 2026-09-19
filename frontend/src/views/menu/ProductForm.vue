<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <button
          @click="$router.push('/menu/products')"
          class="p-2 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
        >
          <ArrowLeft class="w-4 h-4" />
        </button>
        <div>
          <h2 class="text-2xl font-black text-slate-900 tracking-tight">Tambah Menu & Produk Baru</h2>
          <p class="text-xs text-slate-500 mt-0.5">Daftarkan item menu POS & KDS langsung ke database sistem</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="$router.push('/menu/products')"
          class="px-4 py-2.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-700 text-xs font-bold transition-colors shadow-xs cursor-pointer"
        >
          Batal
        </button>
        <button
          @click="submitProduct"
          :disabled="submitting"
          class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>{{ submitting ? 'Menyimpan...' : 'Simpan Produk' }}</span>
        </button>
      </div>
    </div>

    <!-- Main Form Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Primary Details -->
      <div class="lg:col-span-2 bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-5">
        <h4 class="text-sm font-black text-slate-900 uppercase tracking-wider pb-3 border-b border-slate-100">
          Informasi Utama Menu
        </h4>

        <div class="space-y-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">Nama Produk / Menu <span class="text-rose-500">*</span></label>
            <input
              v-model="form.name"
              type="text"
              placeholder="Contoh: Iced Caramel Macchiato"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-semibold text-slate-900"
            />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Kategori Menu <span class="text-rose-500">*</span></label>
              <select
                v-model="form.category_id"
                class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium text-slate-800"
              >
                <option value="">-- Pilih Kategori --</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                  {{ cat.name }}
                </option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">SKU / Kode Barang <span class="text-rose-500">*</span></label>
              <input
                v-model="form.sku"
                type="text"
                placeholder="BEV-ICM-01"
                class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-mono text-slate-800 uppercase"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Harga Dasar Jual (IDR) <span class="text-rose-500">*</span></label>
              <input
                v-model.number="form.base_price"
                type="number"
                min="0"
                step="500"
                class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-bold text-blue-700"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Target Station KDS <span class="text-rose-500">*</span></label>
              <select
                v-model="form.target_station"
                class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-medium text-slate-800"
              >
                <option value="bar">Bar (Minuman / Coffee)</option>
                <option value="kitchen">Kitchen (Makanan Berat / Snack)</option>
                <option value="dessert">Dessert & Pastry</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">Deskripsi Singkat / Resep Singkat</label>
            <textarea
              v-model="form.description"
              rows="3"
              placeholder="Perpaduan espresso arabica, fresh milk, sirup vanilla, dan drizzle karamel gurih manis."
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none text-slate-800 resize-none"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- Right Column: Settings & Preview -->
      <div class="space-y-6">
        <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-4">
          <h4 class="text-sm font-black text-slate-900 uppercase tracking-wider pb-3 border-b border-slate-100">
            Pengaturan Visibilitas
          </h4>

          <div class="flex items-center justify-between p-3.5 rounded-xl bg-slate-50 border border-slate-200">
            <div>
              <span class="text-xs font-bold text-slate-800 block">Tersedia untuk Dipesan</span>
              <span class="text-[11px] text-slate-500">Muncul di katalog POS kasir</span>
            </div>
            <input
              type="checkbox"
              v-model="form.is_active"
              class="w-5 h-5 accent-blue-600 rounded cursor-pointer"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">URL Gambar Produk</label>
            <input
              v-model="form.image_url"
              type="text"
              placeholder="https://images.unsplash.com/..."
              class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none text-slate-700"
            />
          </div>

          <!-- Preview Card -->
          <div class="pt-3 border-t border-slate-100">
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-2">Pratinjau POS Card</span>
            <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/60 flex items-center gap-3">
              <div class="w-12 h-12 rounded-xl bg-blue-100 text-blue-600 flex items-center justify-center font-black text-base shrink-0">
                ☕
              </div>
              <div class="min-w-0 flex-1">
                <div class="text-xs font-bold text-slate-900 truncate">{{ form.name || 'Nama Menu' }}</div>
                <div class="text-[11px] font-extrabold text-blue-600">Rp {{ (form.base_price || 0).toLocaleString('id-ID') }}</div>
                <span class="inline-block mt-0.5 px-2 py-0.5 rounded text-[9px] font-bold bg-slate-200 text-slate-700 uppercase">
                  {{ form.target_station }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ArrowLeft, Plus } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const router = useRouter()
const notifyStore = useNotificationStore()

const submitting = ref(false)
const categories = ref<any[]>([])

const form = ref({
  name: '',
  category_id: '',
  sku: '',
  base_price: 35000,
  target_station: 'bar',
  description: '',
  image_url: '',
  is_active: true
})

const fetchCategories = async () => {
  try {
    const res = await axios.get('/api/v1/master/categories')
    const list = Array.isArray(res.data) ? res.data : (res.data?.data || [])
    categories.value = list
    if (list.length > 0 && !form.value.category_id) {
      form.value.category_id = list[0].id
    }
  } catch (err) {
    console.error('Failed to load categories:', err)
  }
}

const submitProduct = async () => {
  if (!form.value.name.trim()) {
    notifyStore.warning('Nama menu wajib diisi!', 'Validasi')
    return
  }
  if (!form.value.sku.trim()) {
    notifyStore.warning('SKU wajib diisi!', 'Validasi')
    return
  }
  if (!form.value.category_id) {
    notifyStore.warning('Pilih kategori menu terlebih dahulu!', 'Validasi')
    return
  }

  submitting.value = true
  try {
    const res = await axios.post('/api/v1/master/products', form.value)
    notifyStore.success(res.data?.message || 'Produk baru berhasil disimpan ke database!', 'Sukses')
    router.push('/menu/products')
  } catch (err: any) {
    console.error('Failed to create product:', err)
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan produk ke database', 'Error')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchCategories()
})
</script>
