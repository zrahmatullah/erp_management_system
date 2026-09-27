import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAppStore } from '@/stores/app.store'

describe('useAppStore', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
  })

  it('should initialize with default sidebar and theme mode', () => {
    const store = useAppStore()
    expect(store.sidebarCollapsed).toBe(false)
    expect(store.themeMode).toBe('system')
    expect(['light', 'dark']).toContain(store.resolvedTheme)
  })

  it('should toggle sidebar collapsed state', () => {
    const store = useAppStore()
    store.toggleSidebar()
    expect(store.sidebarCollapsed).toBe(true)

    store.toggleSidebar()
    expect(store.sidebarCollapsed).toBe(false)
  })

  it('should set theme mode explicitly and persist to localStorage', () => {
    const store = useAppStore()
    store.setTheme('light')

    expect(store.themeMode).toBe('light')
    expect(store.resolvedTheme).toBe('light')
    expect(localStorage.getItem('cafe_erp_theme')).toBe('light')

    store.setTheme('dark')
    expect(store.themeMode).toBe('dark')
    expect(store.resolvedTheme).toBe('dark')
    expect(localStorage.getItem('cafe_erp_theme')).toBe('dark')
  })

  it('should cycle theme through modes', () => {
    const store = useAppStore()
    store.setTheme('light')
    store.cycleTheme()
    expect(store.themeMode).toBe('dark')
    store.cycleTheme()
    expect(store.themeMode).toBe('system')
    store.cycleTheme()
    expect(store.themeMode).toBe('light')
  })

  it('should set breadcrumbs and current branch', () => {
    const store = useAppStore()
    store.setBreadcrumbs([{ label: 'Home', path: '/' }, { label: 'POS' }])
    expect(store.breadcrumbs.length).toBe(2)

    store.setCurrentBranch('branch-senopati')
    expect(store.currentBranch).toBe('branch-senopati')
  })
})
