<template>
  <div class="fixed top-5 right-5 z-[9999] pointer-events-none flex flex-col gap-3 w-96 max-w-[calc(100vw-2.5rem)]">
    <transition-group
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-4"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-for="toast in notifyStore.toasts"
        :key="toast.id"
        class="pointer-events-auto flex items-start gap-3 p-4 rounded-2xl border shadow-xl bg-white/95 backdrop-blur-md transition-all select-none"
        :class="getToastClasses(toast.type)"
      >
        <!-- Icon -->
        <div
          class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0"
          :class="getIconWrapperClasses(toast.type)"
        >
          <component :is="getIconComponent(toast.type)" class="w-5 h-5" />
        </div>

        <!-- Content -->
        <div class="flex-1 min-w-0 pt-0.5">
          <h4 v-if="toast.title" class="text-xs font-bold text-slate-900 leading-tight">
            {{ toast.title }}
          </h4>
          <p class="text-xs text-slate-600 font-medium leading-relaxed mt-0.5 break-words">
            {{ toast.message }}
          </p>
        </div>

        <!-- Close Button -->
        <button
          @click="notifyStore.removeToast(toast.id)"
          class="text-slate-400 hover:text-slate-700 p-1 rounded-lg hover:bg-slate-100 transition-colors shrink-0 cursor-pointer"
          title="Tutup"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore, type Toast } from '@/stores/notification.store'
import {
  CheckCircle2,
  AlertCircle,
  AlertTriangle,
  Info,
  X
} from 'lucide-vue-next'

const notifyStore = useNotificationStore()

const getToastClasses = (type: Toast['type']) => {
  switch (type) {
    case 'success':
      return 'border-emerald-200/90 shadow-emerald-500/10'
    case 'error':
      return 'border-rose-200/90 shadow-rose-500/10'
    case 'warning':
      return 'border-amber-200/90 shadow-amber-500/10'
    case 'info':
    default:
      return 'border-blue-200/90 shadow-blue-500/10'
  }
}

const getIconWrapperClasses = (type: Toast['type']) => {
  switch (type) {
    case 'success':
      return 'bg-emerald-50 text-emerald-600'
    case 'error':
      return 'bg-rose-50 text-rose-600'
    case 'warning':
      return 'bg-amber-50 text-amber-600'
    case 'info':
    default:
      return 'bg-blue-50 text-blue-600'
  }
}

const getIconComponent = (type: Toast['type']) => {
  switch (type) {
    case 'success':
      return CheckCircle2
    case 'error':
      return AlertCircle
    case 'warning':
      return AlertTriangle
    case 'info':
    default:
      return Info
  }
}
</script>

