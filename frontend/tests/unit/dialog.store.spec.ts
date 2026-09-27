import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDialogStore } from '@/stores/dialog.store'

describe('useDialogStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with dialog closed', () => {
    const store = useDialogStore()
    expect(store.isOpen).toBe(false)
    expect(store.inputValue).toBe('')
  })

  it('should open confirm dialog and resolve true on confirm', async () => {
    const store = useDialogStore()
    const promise = store.confirm({ title: 'Hapus Pesanan?', message: 'Data tidak bisa dipulihkan' })

    expect(store.isOpen).toBe(true)
    expect(store.options.title).toBe('Hapus Pesanan?')

    store.handleConfirm()
    const result = await promise
    expect(result).toBe(true)
    expect(store.isOpen).toBe(false)
  })

  it('should open confirm dialog and resolve false on cancel', async () => {
    const store = useDialogStore()
    const promise = store.confirm({ title: 'Batalkan?' })

    expect(store.isOpen).toBe(true)
    store.handleCancel()
    const result = await promise
    expect(result).toBe(false)
    expect(store.isOpen).toBe(false)
  })

  it('should handle prompt dialog and resolve with input value', async () => {
    const store = useDialogStore()
    const promise = store.prompt({ title: 'Masukkan Alasan', defaultValue: 'Stok Habis' })

    expect(store.isOpen).toBe(true)
    expect(store.inputValue).toBe('Stok Habis')

    store.inputValue = 'Bahan baku rusak'
    store.handleConfirm()

    const result = await promise
    expect(result).toBe('Bahan baku rusak')
    expect(store.isOpen).toBe(false)
  })

  it('should resolve null on prompt cancel', async () => {
    const store = useDialogStore()
    const promise = store.prompt({ title: 'Input' })

    store.handleCancel()
    const result = await promise
    expect(result).toBeNull()
  })
})
