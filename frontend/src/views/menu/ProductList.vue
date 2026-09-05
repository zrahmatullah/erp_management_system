<template>
  <div>
    <PageHeader title="Daftar Produk">
      <n-button type="primary" @click="$router.push('/menu/products/create')">Tambah Produk</n-button>
    </PageHeader>
    <n-card>
      <div class="flex gap-4 mb-4">
        <n-input v-model:value="searchQuery" placeholder="Cari menu..." class="w-64" />
      </div>
      <n-data-table :columns="columns" :data="filteredData" :pagination="pagination" :loading="loading" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { NCard, NDataTable, NButton, NInput, NSwitch } from 'naive-ui'
import PageHeader from '@/components/common/PageHeader.vue'

const pagination = { pageSize: 10 }
const searchQuery = ref('')
const loading = ref(false)
const data = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/master/products')
    const list = Array.isArray(res.data) ? res.data : (res.data?.data || [])
    data.value = list.map((p: any) => ({
      id: p.id,
      name: p.name,
      category: p.category_name || p.target_station || 'Menu',
      price: p.base_price || p.price || 0,
      status: p.is_active !== false
    }))
  } catch (err) {
    console.error('Failed to load products:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const filteredData = computed(() => {
  if (!searchQuery.value) return data.value
  return data.value.filter(p => p.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
})

const columns = [
  { title: 'Nama', key: 'name' },
  { title: 'Kategori / Station', key: 'category' },
  { title: 'Harga', key: 'price', render(row: any) { return `Rp ${Number(row.price || 0).toLocaleString('id-ID')}` } },
  { title: 'Status', key: 'status', render(row: any) { return h(NSwitch, { value: row.status }) } },
  { title: 'Aksi', key: 'actions', render(row: any) { return h(NButton, { size: 'small' }, { default: () => 'Edit' }) } }
]
</script>
