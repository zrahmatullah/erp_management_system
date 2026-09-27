import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useNotificationStore } from '@/stores/notification.store'

describe('useNotificationStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  it('should initialize with empty toasts and notifications', () => {
    const store = useNotificationStore()
    expect(store.toasts).toEqual([])
    expect(store.notifications).toEqual([])
    expect(store.unreadCount).toBe(0)
  })

  it('should add toast with success helper', () => {
    const store = useNotificationStore()
    const id = store.success('Operasi Berhasil', 'Sukses!')

    expect(id).toBeDefined()
    expect(store.toasts.length).toBe(1)
    expect(store.toasts[0].type).toBe('success')
    expect(store.toasts[0].message).toBe('Operasi Berhasil')
    expect(store.toasts[0].title).toBe('Sukses!')
  })

  it('should add toast with error helper with 5000ms duration', () => {
    const store = useNotificationStore()
    store.error('Koneksi database gagal')

    expect(store.toasts.length).toBe(1)
    expect(store.toasts[0].type).toBe('error')
    expect(store.toasts[0].duration).toBe(5000)
  })

  it('should remove toast manually by id', () => {
    const store = useNotificationStore()
    const id = store.warning('Peringatan stok menipis')
    expect(store.toasts.length).toBe(1)

    store.removeToast(id)
    expect(store.toasts.length).toBe(0)
  })

  it('should auto-remove toast after specified duration', () => {
    const store = useNotificationStore()
    store.info('Info update sistem')
    expect(store.toasts.length).toBe(1)

    // Fast-forward timers
    vi.advanceTimersByTime(4500)
    expect(store.toasts.length).toBe(0)
  })

  it('should compute unread notifications count', () => {
    const store = useNotificationStore()
    store.notifications = [
      { id: '1', title: 'Order Baru', message: 'Meja 5', isRead: false, createdAt: '2026-09-27' },
      { id: '2', title: 'Stok Habis', message: 'Susu UHT', isRead: true, createdAt: '2026-09-27' },
      { id: '3', title: 'Shift Dimulai', message: 'Kasir 1', isRead: false, createdAt: '2026-09-27' }
    ]

    expect(store.unreadCount).toBe(2)
  })
})
