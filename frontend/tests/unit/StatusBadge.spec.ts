import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { NTag } from 'naive-ui'

describe('StatusBadge component', () => {
  it('should render uppercase status label and warning type for pending order', () => {
    const wrapper = mount(StatusBadge, {
      props: { status: 'pending', type: 'order' }
    })
    expect(wrapper.text()).toBe('PENDING')
    const tag = wrapper.findComponent(NTag)
    expect(tag.props('type')).toBe('warning')
  })

  it('should render correct tag type for completed status', () => {
    const wrapper = mount(StatusBadge, {
      props: { status: 'completed', type: 'order' }
    })
    expect(wrapper.text()).toBe('COMPLETED')
    const tag = wrapper.findComponent(NTag)
    expect(tag.props('type')).toBe('success')
  })

  it('should render correct tag type for cancelled / rejected status', () => {
    const wrapper = mount(StatusBadge, {
      props: { status: 'cancelled', type: 'order' }
    })
    expect(wrapper.text()).toBe('CANCELLED')
    const tag = wrapper.findComponent(NTag)
    expect(tag.props('type')).toBe('error')

    const wrapperRejected = mount(StatusBadge, {
      props: { status: 'rejected', type: 'expense' }
    })
    expect(wrapperRejected.text()).toBe('REJECTED')
    expect(wrapperRejected.findComponent(NTag).props('type')).toBe('error')
  })

  it('should fallback to default tag type for unknown status', () => {
    const wrapper = mount(StatusBadge, {
      props: { status: 'archived', type: 'po' }
    })
    expect(wrapper.text()).toBe('ARCHIVED')
    const tag = wrapper.findComponent(NTag)
    expect(tag.props('type')).toBe('default')
  })
})
