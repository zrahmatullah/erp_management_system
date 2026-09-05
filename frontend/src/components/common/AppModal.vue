<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/40 backdrop-blur-md p-4 transition-all duration-300"
        @click.self="$emit('close')"
        @keydown.esc="$emit('close')"
      >
        <div
          class="bg-white rounded-3xl w-full shadow-2xl shadow-slate-950/25 border border-slate-100/90 flex flex-col overflow-hidden transform transition-all duration-300 ease-out"
          :class="maxWidth || 'max-w-md'"
          role="dialog"
          aria-modal="true"
        >
          <!-- Modal Header -->
          <div class="flex items-center justify-between p-5 border-b border-slate-100 shrink-0">
            <div class="flex items-center gap-3">
              <div
                v-if="icon"
                class="w-10 h-10 rounded-2xl bg-blue-50 text-blue-600 border border-blue-100 flex items-center justify-center shrink-0 shadow-xs"
              >
                <component :is="icon" class="w-5 h-5" />
              </div>
              <div>
                <h3 class="text-base font-black text-slate-900 tracking-tight leading-snug">
                  {{ title }}
                </h3>
                <p v-if="subtitle" class="text-xs text-slate-400 mt-0.5">
                  {{ subtitle }}
                </p>
              </div>
            </div>

            <button
              @click="$emit('close')"
              class="text-slate-400 hover:text-slate-600 p-2 rounded-xl hover:bg-slate-100 transition-colors cursor-pointer"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <!-- Modal Body -->
          <div class="p-5 overflow-y-auto max-h-[75vh]">
            <slot />
          </div>

          <!-- Modal Footer -->
          <div v-if="$slots.footer" class="p-4 border-t border-slate-100 bg-slate-50/50 flex items-center justify-end gap-2.5 shrink-0">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { X } from 'lucide-vue-next'

defineProps<{
  show: boolean
  title: string
  subtitle?: string
  icon?: any
  maxWidth?: string
}>()

defineEmits(['close'])
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-fade-enter-active > div,
.modal-fade-leave-active > div {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-from > div,
.modal-fade-leave-to > div {
  opacity: 0;
  transform: scale(0.94) translateY(8px);
}
</style>
