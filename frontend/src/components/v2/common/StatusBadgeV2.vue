<template>
  <span 
    class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold border select-none transition-colors"
    :class="badgeStyle"
  >
    <span class="w-1.5 h-1.5 rounded-full shrink-0" :class="dotStyle"></span>
    <span>{{ formattedLabel }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  status: string
  type?: 'order' | 'payment' | 'po' | 'leave' | 'expense' | 'employee' | 'general'
  label?: string
}>(), {
  type: 'general'
})

const normStatus = computed(() => props.status?.toLowerCase().trim() || 'default')

const formattedLabel = computed(() => {
  if (props.label) return props.label
  const s = normStatus.value
  const dict: Record<string, string> = {
    active: 'Aktif',
    completed: 'Selesai',
    paid: 'Lunas',
    pending: 'Menunggu',
    cancelled: 'Dibatalkan',
    draft: 'Draft',
    approved: 'Disetujui',
    rejected: 'Ditolak',
    unpaid: 'Belum Bayar',
    open: 'Buka',
    closed: 'Tutup',
    on_leave: 'Sedang Cuti',
    inactive: 'Nonaktif'
  }
  return dict[s] || props.status.toUpperCase()
})

const badgeStyle = computed(() => {
  const s = normStatus.value
  switch (s) {
    case 'active':
    case 'completed':
    case 'paid':
    case 'approved':
    case 'open':
      return 'bg-emerald-50 text-emerald-700 border-emerald-200/80 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-800/60'
    case 'pending':
    case 'unpaid':
    case 'draft':
      return 'bg-amber-50 text-amber-700 border-amber-200/80 dark:bg-amber-950/40 dark:text-amber-300 dark:border-amber-800/60'
    case 'cancelled':
    case 'rejected':
    case 'inactive':
    case 'closed':
      return 'bg-rose-50 text-rose-700 border-rose-200/80 dark:bg-rose-950/40 dark:text-rose-300 dark:border-rose-800/60'
    case 'on_leave':
      return 'bg-purple-50 text-purple-700 border-purple-200/80 dark:bg-purple-950/40 dark:text-purple-300 dark:border-purple-800/60'
    default:
      return 'bg-slate-50 text-slate-700 border-slate-200/80 dark:bg-slate-800/40 dark:text-slate-300 dark:border-slate-700/60'
  }
})

const dotStyle = computed(() => {
  const s = normStatus.value
  switch (s) {
    case 'active':
    case 'completed':
    case 'paid':
    case 'approved':
    case 'open':
      return 'bg-emerald-500'
    case 'pending':
    case 'unpaid':
    case 'draft':
      return 'bg-amber-500'
    case 'cancelled':
    case 'rejected':
    case 'inactive':
    case 'closed':
      return 'bg-rose-500'
    case 'on_leave':
      return 'bg-purple-500'
    default:
      return 'bg-slate-400'
  }
})
</script>

