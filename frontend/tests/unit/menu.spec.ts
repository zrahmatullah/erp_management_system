import { describe, it, expect } from 'vitest'
import { MAIN_MENU } from '@/constants/menu'

describe('MAIN_MENU constant', () => {
  it('should define all required ERP navigation items', () => {
    expect(Array.isArray(MAIN_MENU)).toBe(true)
    expect(MAIN_MENU.length).toBe(8)

    const keys = MAIN_MENU.map(m => m.key)
    expect(keys).toContain('dashboard')
    expect(keys).toContain('pos')
    expect(keys).toContain('menu')
    expect(keys).toContain('inventory')
    expect(keys).toContain('hris')
    expect(keys).toContain('finance')
    expect(keys).toContain('reports')
    expect(keys).toContain('settings')
  })

  it('should have valid path starting with slash for all menu items', () => {
    MAIN_MENU.forEach(item => {
      expect(item.path.startsWith('/')).toBe(true)
      expect(item.label.length).toBeGreaterThan(0)
    })
  })
})
