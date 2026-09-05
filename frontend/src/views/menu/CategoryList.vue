<template>
  <div>
    <PageHeader title="Kategori Menu">
      <n-button type="primary" @click="showModal = true">Tambah Kategori</n-button>
    </PageHeader>
    <n-card>
      <n-data-table :columns="columns" :data="data" :pagination="pagination" :loading="loading" />
    </n-card>
    
    <n-modal v-model:show="showModal" preset="card" title="Tambah Kategori" class="w-[500px]">
      <n-form>
        <n-form-item label="Nama Kategori"><n-input v-model:value="newCategory.name" placeholder="e.g. Coffee" /></n-form-item>
        <n-form-item label="Slug / Code"><n-input v-model:value="newCategory.slug" placeholder="e.g. coffee" /></n-form-item>
        <n-button type="primary" block :loading="saving" @click="saveCategory">Simpan</n-button>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import axios from 'axios'
import { NCard, NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSwitch } from 'naive-ui'
import PageHeader from '@/components/common/PageHeader.vue'

const showModal = ref(false)
const loading = ref(false)
const saving = ref(false)
const pagination = { pageSize: 10 }
const data = ref<any[]>([])
const newCategory = ref({ name: '', slug: '' })

const fetchCategories = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/master/categories')
    const list = Array.isArray(res.data) ? res.data : (res.data?.data || [])
    data.value = list.map((c: any) => ({
      id: c.id,
      name: c.name,
      slug: c.slug,
      products: c.products_count || 5,
      status: c.is_active !== false
    }))
  } catch (err) {
    console.error('Failed to load categories:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCategories()
})

const saveCategory = async () => {
  if (!newCategory.value.name) return
  saving.value = true
  try {
    await axios.post('/api/v1/master/categories', {
      name: newCategory.value.name,
      slug: newCategory.value.slug || newCategory.value.name.toLowerCase().replace(/\s+/g, '-'),
      is_active: true
    })
    showModal.value = false
    newCategory.value = { name: '', slug: '' }
    await fetchCategories()
  } catch (err) {
    console.error('Failed to save category:', err)
  } finally {
    saving.value = false
  }
}

const columns = [
  { title: 'Nama Kategori', key: 'name' },
  { title: 'Slug / Code', key: 'slug' },
  { title: 'Status', key: 'status', render(row: any) { return h(NSwitch, { value: row.status }) } }
]
</script>
