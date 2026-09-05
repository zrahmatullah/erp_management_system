import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface DialogOptions {
  title: string
  message?: string
  type?: 'confirm' | 'danger' | 'warning' | 'prompt' | 'info'
  confirmText?: string
  cancelText?: string
  placeholder?: string
  inputType?: string
  defaultValue?: string
}

export const useDialogStore = defineStore('dialog', () => {
  const isOpen = ref(false)
  const options = ref<DialogOptions>({
    title: '',
    message: '',
    type: 'confirm',
    confirmText: 'Konfirmasi',
    cancelText: 'Batal'
  })
  const inputValue = ref('')

  let resolvePromise: ((value: any) => void) | null = null

  const confirm = (opts: DialogOptions): Promise<boolean> => {
    options.value = {
      confirmText: 'Konfirmasi',
      cancelText: 'Batal',
      type: 'confirm',
      ...opts
    }
    inputValue.value = ''
    isOpen.value = true

    return new Promise((resolve) => {
      resolvePromise = resolve
    })
  }

  const prompt = (opts: DialogOptions): Promise<string | null> => {
    options.value = {
      confirmText: 'Simpan',
      cancelText: 'Batal',
      type: 'prompt',
      ...opts
    }
    inputValue.value = opts.defaultValue || ''
    isOpen.value = true

    return new Promise((resolve) => {
      resolvePromise = resolve
    })
  }

  const handleConfirm = () => {
    isOpen.value = false
    if (resolvePromise) {
      if (options.value.type === 'prompt') {
        resolvePromise(inputValue.value)
      } else {
        resolvePromise(true)
      }
      resolvePromise = null
    }
  }

  const handleCancel = () => {
    isOpen.value = false
    if (resolvePromise) {
      if (options.value.type === 'prompt') {
        resolvePromise(null)
      } else {
        resolvePromise(false)
      }
      resolvePromise = null
    }
  }

  return {
    isOpen,
    options,
    inputValue,
    confirm,
    prompt,
    handleConfirm,
    handleCancel
  }
})
