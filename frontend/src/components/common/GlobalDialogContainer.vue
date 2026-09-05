<template>
  <div
    v-if="dialogStore.isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-xs p-4 animate-in fade-in duration-150"
    @click.self="dialogStore.handleCancel"
    @keydown.esc="dialogStore.handleCancel"
  >
    <div
      class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl border border-slate-100 animate-in zoom-in-95 duration-150"
      role="dialog"
      aria-modal="true"
    >
      <div class="flex items-start gap-4">
        <!-- Icon Badge -->
        <div
          class="w-11 h-11 rounded-2xl flex items-center justify-center shrink-0 border"
          :class="iconContainerClass"
        >
          <component :is="dialogIcon" class="w-5 h-5" />
        </div>

        <!-- Content Area -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between">
            <h3 class="text-base font-black text-slate-900 tracking-tight leading-snug">
              {{ dialogStore.options.title }}
            </h3>
            <button
              @click="dialogStore.handleCancel"
              class="text-slate-400 hover:text-slate-600 p-1 rounded-lg hover:bg-slate-100 transition-colors -mt-1 -mr-1"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <p v-if="dialogStore.options.message" class="text-xs text-slate-500 mt-1.5 leading-relaxed">
            {{ dialogStore.options.message }}
          </p>

          <!-- Input field for prompt dialogs -->
          <form v-if="dialogStore.options.type === 'prompt'" @submit.prevent="dialogStore.handleConfirm" class="mt-4">
            <input
              ref="inputRef"
              v-model="dialogStore.inputValue"
              :type="dialogStore.options.inputType || 'text'"
              :placeholder="dialogStore.options.placeholder"
              class="w-full px-3.5 py-2.5 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-2 focus:ring-blue-600/20 outline-none font-semibold text-slate-900 transition-all"
              autofocus
            />
          </form>
        </div>
      </div>

      <!-- Action Footer -->
      <div class="flex items-center justify-end gap-2.5 pt-4 mt-6 border-t border-slate-100">
        <button
          type="button"
          @click="dialogStore.handleCancel"
          class="px-4 py-2.5 text-xs font-bold text-slate-600 bg-slate-100 hover:bg-slate-200 rounded-xl transition-colors cursor-pointer"
        >
          {{ dialogStore.options.cancelText || 'Batal' }}
        </button>

        <button
          type="button"
          @click="dialogStore.handleConfirm"
          class="px-5 py-2.5 text-xs font-bold text-white rounded-xl shadow-md transition-all cursor-pointer"
          :class="confirmBtnClass"
        >
          {{ dialogStore.options.confirmText || 'Konfirmasi' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'
import {
  AlertTriangle,
  AlertCircle,
  HelpCircle,
  CheckCircle2,
  X
} from 'lucide-vue-next'
import { useDialogStore } from '@/stores/dialog.store'

const dialogStore = useDialogStore()
const inputRef = ref<HTMLInputElement | null>(null)

watch(() => dialogStore.isOpen, (open) => {
  if (open && dialogStore.options.type === 'prompt') {
    nextTick(() => {
      inputRef.value?.focus()
      inputRef.value?.select()
    })
  }
})

const dialogIcon = computed(() => {
  switch (dialogStore.options.type) {
    case 'danger':
      return AlertTriangle
    case 'warning':
      return AlertCircle
    case 'prompt':
      return HelpCircle
    default:
      return CheckCircle2
  }
})

const iconContainerClass = computed(() => {
  switch (dialogStore.options.type) {
    case 'danger':
      return 'bg-rose-50 text-rose-600 border-rose-100'
    case 'warning':
      return 'bg-amber-50 text-amber-600 border-amber-100'
    case 'prompt':
      return 'bg-blue-50 text-blue-600 border-blue-100'
    default:
      return 'bg-emerald-50 text-emerald-600 border-emerald-100'
  }
})

const confirmBtnClass = computed(() => {
  switch (dialogStore.options.type) {
    case 'danger':
      return 'bg-rose-600 hover:bg-rose-700 shadow-rose-600/30'
    case 'warning':
      return 'bg-amber-600 hover:bg-amber-700 shadow-amber-600/30'
    default:
      return 'bg-blue-600 hover:bg-blue-700 shadow-blue-600/30'
  }
})
</script>

