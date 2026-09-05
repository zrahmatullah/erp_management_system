<template>
  <div>
    <PageHeader title="Buat Purchase Order" show-back />
    <n-card>
      <n-form>
        <div class="grid grid-cols-2 gap-4">
          <n-form-item label="Supplier"><n-select /></n-form-item>
          <n-form-item label="Tgl Pengiriman"><n-date-picker class="w-full" /></n-form-item>
        </div>
      </n-form>
      <!-- Items table would go here -->
      <div class="flex justify-end gap-2 mt-6">
        <n-button>Simpan Draft</n-button>
        <n-button type="primary">Ajukan PO</n-button>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { NCard, NForm, NFormItem, NSelect, NDatePicker, NButton } from 'naive-ui'
import PageHeader from '@/components/common/PageHeader.vue'
</script>
