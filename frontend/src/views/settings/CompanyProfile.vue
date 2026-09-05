<template>
  <div>
    <PageHeader title="Profil Perusahaan" />
    <n-card>
      <n-form label-placement="left" label-width="150" class="max-w-2xl">
        <n-form-item label="Nama Perusahaan"><n-input value="Cafe ERP Indonesia" /></n-form-item>
        <n-form-item label="NPWP"><n-input value="01.234.567.8-901.000" /></n-form-item>
        <n-form-item label="Email"><n-input value="hello@cafeeerp.com" /></n-form-item>
        <n-form-item label="Telepon"><n-input value="021-5551234" /></n-form-item>
        <n-form-item label="Alamat"><n-input type="textarea" value="Jl. Kopi Manis No 1, Jakarta" /></n-form-item>
        <n-form-item label="Logo"><n-upload /></n-form-item>
        <div class="flex justify-end mt-4">
          <n-button type="primary">Simpan Perubahan</n-button>
        </div>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { NCard, NForm, NFormItem, NInput, NUpload, NButton } from 'naive-ui'
import PageHeader from '@/components/common/PageHeader.vue'
</script>
