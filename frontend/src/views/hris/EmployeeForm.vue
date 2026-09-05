<template>
  <div>
    <PageHeader title="Tambah Karyawan" show-back />
    <n-card>
      <n-tabs>
        <n-tab-pane name="personal" tab="Data Pribadi">
          <!-- form items -->
          <n-form class="max-w-xl">
            <n-form-item label="Nama Depan"><n-input /></n-form-item>
            <n-form-item label="Nama Belakang"><n-input /></n-form-item>
            <n-form-item label="Email"><n-input /></n-form-item>
          </n-form>
        </n-tab-pane>
        <n-tab-pane name="employment" tab="Kepegawaian"></n-tab-pane>
        <n-tab-pane name="bank" tab="Rekening Bank"></n-tab-pane>
      </n-tabs>
      <div class="flex justify-end mt-4">
        <n-button type="primary">Simpan</n-button>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { NCard, NTabs, NTabPane, NForm, NFormItem, NInput, NButton } from 'naive-ui'
import PageHeader from '@/components/common/PageHeader.vue'
</script>
