import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUiVersionStore } from '@/stores/uiVersion.store'
import { UI_VERSION_STORAGE_KEY } from '@/config/uiVersion'

describe('uiVersion.store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    document.documentElement.removeAttribute('data-ui-version')
    document.documentElement.classList.remove('ui-v1', 'ui-v2')
  })

  it('initializes with default version and applies DOM attribute', () => {
    const store = useUiVersionStore()
    expect(store.currentVersion).toBeDefined()
    expect(['v1', 'v2']).toContain(store.currentVersion)
    expect(document.documentElement.getAttribute('data-ui-version')).toBe(store.currentVersion)
  })

  it('sets version explicitly to v1 and updates state and localStorage', () => {
    const store = useUiVersionStore()
    store.setVersion('v1')
    expect(store.currentVersion).toBe('v1')
    expect(store.isV1).toBe(true)
    expect(store.isV2).toBe(false)
    expect(localStorage.getItem(UI_VERSION_STORAGE_KEY)).toBe('v1')
    expect(document.documentElement.getAttribute('data-ui-version')).toBe('v1')
    expect(document.documentElement.classList.contains('ui-v1')).toBe(true)
    expect(document.documentElement.classList.contains('ui-v2')).toBe(false)
  })

  it('sets version explicitly to v2 and updates state and localStorage', () => {
    const store = useUiVersionStore()
    store.setVersion('v2')
    expect(store.currentVersion).toBe('v2')
    expect(store.isV1).toBe(false)
    expect(store.isV2).toBe(true)
    expect(localStorage.getItem(UI_VERSION_STORAGE_KEY)).toBe('v2')
    expect(document.documentElement.getAttribute('data-ui-version')).toBe('v2')
    expect(document.documentElement.classList.contains('ui-v2')).toBe(true)
    expect(document.documentElement.classList.contains('ui-v1')).toBe(false)
  })

  it('toggles version between v1 and v2', () => {
    const store = useUiVersionStore()
    store.setVersion('v2')
    expect(store.isV2).toBe(true)

    store.toggleVersion()
    expect(store.isV1).toBe(true)
    expect(store.currentVersion).toBe('v1')

    store.toggleVersion()
    expect(store.isV2).toBe(true)
    expect(store.currentVersion).toBe('v2')
  })

  it('resets to env default and removes localStorage item', () => {
    const store = useUiVersionStore()
    store.setVersion('v1')
    expect(localStorage.getItem(UI_VERSION_STORAGE_KEY)).toBe('v1')

    store.resetToEnvDefault()
    expect(localStorage.getItem(UI_VERSION_STORAGE_KEY)).toBeNull()
  })
})

