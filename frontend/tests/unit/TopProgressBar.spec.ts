import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import TopProgressBar from '@/components/common/TopProgressBar.vue'
import { useLoadingStore } from '@/stores/loading.store'

describe('TopProgressBar component', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should not display bar when not loading', () => {
    const store = useLoadingStore()
    store.isLoading = false
    store.progress = 0

    const wrapper = mount(TopProgressBar)
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('should display bar with correct progress width when loading', async () => {
    const store = useLoadingStore()
    store.isLoading = true
    store.progress = 65

    const wrapper = mount(TopProgressBar)
    expect(wrapper.find('.fixed').exists()).toBe(true)
    const bar = wrapper.find('.bg-gradient-to-r')
    expect(bar.attributes('style')).toContain('width: 65%')
  })
})
