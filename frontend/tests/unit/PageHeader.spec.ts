import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PageHeader from '@/components/common/PageHeader.vue'
import { NButton } from 'naive-ui'

describe('PageHeader component', () => {
  it('should render title and subtitle', () => {
    const wrapper = mount(PageHeader, {
      props: {
        title: 'Manajemen Meja',
        subtitle: 'Daftar denah meja cafe'
      }
    })

    expect(wrapper.text()).toContain('Manajemen Meja')
    expect(wrapper.text()).toContain('Daftar denah meja cafe')
  })

  it('should render back button when showBack is true', async () => {
    let backCalled = false
    const wrapper = mount(PageHeader, {
      props: {
        title: 'Detail Transaksi',
        showBack: true
      },
      global: {
        mocks: {
          $router: {
            back: () => {
              backCalled = true
            }
          }
        }
      }
    })

    const btn = wrapper.findComponent(NButton)
    expect(btn.exists()).toBe(true)
    expect(wrapper.text()).toContain('←')
    await btn.trigger('click')
    expect(backCalled).toBe(true)
  })
})
