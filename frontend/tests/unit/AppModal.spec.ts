import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AppModal from '@/components/common/AppModal.vue'

describe('AppModal component', () => {
  it('should render modal with title when show is true', () => {
    const wrapper = mount(AppModal, {
      props: {
        show: true,
        title: 'Konfirmasi Hapus',
        subtitle: 'Tindakan ini tidak bisa dibatalkan'
      },
      slots: {
        default: '<p class="modal-body-content">Apakah Anda yakin?</p>',
        footer: '<button class="confirm-btn">Ya, Hapus</button>'
      },
      global: {
        stubs: {
          Teleport: true,
          Transition: true
        }
      }
    })

    expect(wrapper.text()).toContain('Konfirmasi Hapus')
    expect(wrapper.text()).toContain('Tindakan ini tidak bisa dibatalkan')
    expect(wrapper.find('.modal-body-content').text()).toBe('Apakah Anda yakin?')
    expect(wrapper.find('.confirm-btn').text()).toBe('Ya, Hapus')
  })

  it('should emit close when close button is clicked', async () => {
    const wrapper = mount(AppModal, {
      props: {
        show: true,
        title: 'Edit Data'
      },
      global: {
        stubs: {
          Teleport: true,
          Transition: true
        }
      }
    })

    const closeBtn = wrapper.find('button')
    await closeBtn.trigger('click')
    expect(wrapper.emitted('close')).toHaveLength(1)
  })
})
