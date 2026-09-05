<template>
  <n-tag :type="tagType" round size="small">
    {{ label }}
  </n-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NTag } from 'naive-ui'

const props = defineProps<{
  status: string
  type: 'order'|'payment'|'po'|'leave'|'expense'|'employee'
}>()

const tagType = computed(() => {
  const map: Record<string, any> = {
    pending: 'warning',
    completed: 'success',
    cancelled: 'error',
    active: 'success',
    draft: 'default',
    approved: 'success',
    rejected: 'error'
  }
  return map[props.status?.toLowerCase()] || 'default'
})

const label = computed(() => {
  // basic translation map
  return props.status.toUpperCase()
})
</script>
