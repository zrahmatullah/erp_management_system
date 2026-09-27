import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useLoadingStore } from '@/stores/loading.store'

describe('useLoadingStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  it('should initialize with not loading and progress 0', () => {
    const store = useLoadingStore()
    expect(store.isLoading).toBe(false)
    expect(store.progress).toBe(0)
  })

  it('should start loading with progress 15 and increase on interval', () => {
    const store = useLoadingStore()
    store.start()

    expect(store.isLoading).toBe(true)
    expect(store.progress).toBe(15)

    vi.advanceTimersByTime(250)
    expect(store.progress).toBeGreaterThan(15)
  })

  it('should finish loading, set progress to 100, then reset', () => {
    const store = useLoadingStore()
    store.start()
    store.finish()

    expect(store.progress).toBe(100)

    // After 300ms + 200ms
    vi.advanceTimersByTime(550)
    expect(store.isLoading).toBe(false)
    expect(store.progress).toBe(0)
  })
})
