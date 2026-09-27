import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatCard from '@/components/common/StatCard.vue'

describe('StatCard component', () => {
  it('should render title and value correctly', () => {
    const wrapper = mount(StatCard, {
      props: {
        title: 'Total Penjualan',
        value: 'Rp 15.000.000',
        trend: 'up',
        trendValue: '+12%',
        color: 'green'
      }
    })

    expect(wrapper.text()).toContain('Total Penjualan')
    expect(wrapper.text()).toContain('Rp 15.000.000')
    expect(wrapper.text()).toContain('↑')
    expect(wrapper.text()).toContain('+12% vs kemarin')
  })

  it('should render downward trend symbol', () => {
    const wrapper = mount(StatCard, {
      props: {
        title: 'Rata-rata Order',
        value: 'Rp 45.000',
        trend: 'down',
        trendValue: '-5%',
        color: 'red'
      }
    })

    expect(wrapper.text()).toContain('↓')
    expect(wrapper.text()).toContain('-5% vs kemarin')
  })

  it('should render neutral trend symbol', () => {
    const wrapper = mount(StatCard, {
      props: {
        title: 'Meja Terisi',
        value: 8,
        trend: 'neutral',
        trendValue: '0%',
        color: 'blue'
      }
    })

    expect(wrapper.text()).toContain('-')
    expect(wrapper.text()).toContain('0% vs kemarin')
  })
})
