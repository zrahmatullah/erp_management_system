<template>
  <div>
    <PageHeader title="Tambah Produk" show-back />
    <n-card>
      <n-tabs type="line">
        <n-tab-pane name="info" tab="Info Dasar">
          <n-form class="max-w-xl">
            <n-form-item label="Nama Produk"><n-input /></n-form-item>
            <n-form-item label="Kategori"><n-select /></n-form-item>
            <n-form-item label="Harga Dasar"><n-input-number /></n-form-item>
            <n-form-item label="Deskripsi"><n-input type="textarea" /></n-form-item>
          </n-form>
        </n-tab-pane>
        <n-tab-pane name="variants" tab="Varian">
          <!-- Varian form -->
        </n-tab-pane>
        <n-tab-pane name="images" tab="Gambar">
          <!-- Images upload -->
        </n-tab-pane>
        <n-tab-pane name="recipe" tab="Resep">
          <!-- Recipe builder -->
        </n-tab-pane>
      </n-tabs>
      <div class="flex justify-end mt-6 gap-2">
        <n-button>Batal</n-button>
        <n-button type="primary">Simpan</n-button>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { NCard, NTabs, NTabPane, NForm, NFormItem, NInput, NSelect, NInputNumber, NButton } from 'naive-ui'
import PageHeader from '@/components/common/PageHeader.vue'
</script>
