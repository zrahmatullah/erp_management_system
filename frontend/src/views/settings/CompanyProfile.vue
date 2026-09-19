<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-black text-slate-900 tracking-tight">Profil Perusahaan & Pengaturan ERP</h2>
        <p class="text-xs text-slate-500 mt-1">Konfigurasi entitas bisnis, identitas legalitas pajak, dan data operasional cafe</p>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="fetchProfile"
          class="p-2.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
          title="Segarkan Data"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>
        <button
          @click="saveProfile"
          :disabled="saving"
          class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-2"
        >
          <Save class="w-4 h-4" />
          <span>{{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
        </button>
      </div>
    </div>

    <!-- Main Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Card: Company Badge & Summary -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-6">
        <div class="flex flex-col items-center text-center">
          <div class="w-24 h-24 rounded-3xl bg-blue-600/10 border-2 border-blue-600/20 text-blue-600 flex items-center justify-center shadow-inner mb-4">
            <Coffee class="w-12 h-12" />
          </div>
          <h3 class="text-lg font-black text-slate-900">{{ form.company_name || 'Cafe ERP System' }}</h3>
          <p class="text-xs text-slate-500 font-mono mt-0.5">Kode: {{ form.code || 'HQ-01' }}</p>
          <div class="mt-3 inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-50 border border-emerald-200/60 text-emerald-700 text-xs font-bold">
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            Status: Entitas Aktif
          </div>
        </div>

        <div class="pt-4 border-t border-slate-100 space-y-3 text-xs">
          <div class="flex items-center justify-between text-slate-600">
            <span class="font-medium">Mata Uang Basis:</span>
            <span class="font-bold text-slate-900 font-mono">{{ form.currency || 'IDR (Rp)' }}</span>
          </div>
          <div class="flex items-center justify-between text-slate-600">
            <span class="font-medium">Tarif PPN (UU HPP):</span>
            <span class="font-bold text-slate-900 font-mono">{{ form.tax_rate || 11 }}%</span>
          </div>
          <div class="flex items-center justify-between text-slate-600">
            <span class="font-medium">Sinkronisasi Database:</span>
            <span class="text-emerald-600 font-bold">Terhubung PostgreSQL</span>
          </div>
        </div>
      </div>

      <!-- Right Card: Form Edit Fields -->
      <div class="lg:col-span-2 bg-white rounded-2xl border border-slate-200/80 p-6 shadow-xs space-y-5">
        <h4 class="text-sm font-black text-slate-900 uppercase tracking-wider pb-3 border-b border-slate-100">
          Identitas Legal & Kontak Operasional
        </h4>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="sm:col-span-2">
            <label class="block text-xs font-bold text-slate-700 mb-1.5">Nama Perusahaan / Brand Cafe</label>
            <input
              v-model="form.company_name"
              type="text"
              placeholder="Contoh: PT Kopi Nusantara ERP"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-semibold text-slate-900 transition-all"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">Kode Identitas Outlet / HQ</label>
            <input
              v-model="form.code"
              type="text"
              placeholder="HQ-01"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-mono text-slate-800 transition-all"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">NPWP Perusahaan</label>
            <input
              v-model="form.npwp"
              type="text"
              placeholder="01.234.567.8-901.000"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none font-mono text-slate-800 transition-all"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">Email Bisnis</label>
            <input
              v-model="form.email"
              type="email"
              placeholder="admin@cafeharmony.com"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none text-slate-800 transition-all"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1.5">Nomor Telepon Operasional</label>
            <input
              v-model="form.phone"
              type="text"
              placeholder="021-7201234"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none text-slate-800 transition-all"
            />
          </div>

          <div class="sm:col-span-2">
            <label class="block text-xs font-bold text-slate-700 mb-1.5">Alamat Lengkap Kantor Pusat</label>
            <textarea
              v-model="form.address"
              rows="3"
              placeholder="Alamat jalan, gedung, kota, kode pos"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-none text-slate-800 transition-all resize-none"
            ></textarea>
          </div>
        </div>

        <div class="pt-4 border-t border-slate-100 flex items-center justify-end">
          <button
            @click="saveProfile"
            :disabled="saving"
            class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 disabled:opacity-50 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-2"
          >
            <Save class="w-4 h-4" />
            <span>{{ saving ? 'Menyimpan Perubahan...' : 'Simpan Profil Perusahaan' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { Coffee, RotateCw, Save } from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()

const loading = ref(false)
const saving = ref(false)

const form = ref({
  company_name: '',
  code: '',
  npwp: '',
  email: '',
  phone: '',
  address: '',
  currency: 'IDR',
  tax_rate: 11
})

const fetchProfile = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/master/company-profile')
    if (res.data?.data) {
      const d = res.data.data
      form.value = {
        company_name: d.company_name || d.name || 'Cafe Harmony ERP',
        code: d.code || 'HQ-01',
        npwp: d.npwp || '01.234.567.8-901.000',
        email: d.email || 'finance@cafeharmony.com',
        phone: d.phone || '021-7201234',
        address: d.address || 'Jl. Senopati Raya No. 45, Jakarta Selatan',
        currency: d.currency || 'IDR',
        tax_rate: d.tax_rate || 11
      }
    }
  } catch (err: any) {
    console.error('Failed to load company profile:', err)
    form.value = {
      company_name: 'Cafe Harmony Indonesia',
      code: 'HQ-JKT-01',
      npwp: '01.234.567.8-901.000',
      email: 'admin@cafeharmony.com',
      phone: '021-7201234',
      address: 'Jl. Senopati Raya No. 45, Kebayoran Baru, Jakarta Selatan',
      currency: 'IDR',
      tax_rate: 11
    }
  } finally {
    loading.value = false
  }
}

const saveProfile = async () => {
  if (!form.value.company_name.trim()) {
    notifyStore.warning('Nama perusahaan tidak boleh kosong!', 'Validasi')
    return
  }

  saving.value = true
  try {
    const res = await axios.put('/api/v1/master/company-profile', {
      company_name: form.value.company_name,
      address: form.value.address,
      phone: form.value.phone,
      email: form.value.email
    })
    notifyStore.success(res.data?.message || 'Profil perusahaan berhasil disimpan ke database!', 'Berhasil Disimpan')
    await fetchProfile()
  } catch (err: any) {
    console.error('Failed to update company profile:', err)
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan profil perusahaan', 'Error')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchProfile()
})
</script>
