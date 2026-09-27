import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { requireAuth, requireRole } from '@/router/guards'
import { useAuthStore } from '@/stores/auth.store'

describe('router navigation guards', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
  })

  describe('requireAuth', () => {
    it('should redirect unauthenticated user to /login when meta.requiresAuth is true', () => {
      const next = vi.fn()
      const to = { meta: { requiresAuth: true } } as any
      const from = {} as any

      requireAuth(to, from, next)
      expect(next).toHaveBeenCalledWith('/login')
    })

    it('should allow navigation when route does not require auth', () => {
      const next = vi.fn()
      const to = { meta: { requiresAuth: false } } as any
      const from = {} as any

      requireAuth(to, from, next)
      expect(next).toHaveBeenCalledWith()
    })

    it('should allow navigation when user is authenticated', () => {
      const authStore = useAuthStore()
      authStore.token = 'valid-token'

      const next = vi.fn()
      const to = { meta: { requiresAuth: true } } as any
      const from = {} as any

      requireAuth(to, from, next)
      expect(next).toHaveBeenCalledWith()
    })
  })

  describe('requireRole', () => {
    it('should allow user with matching role', () => {
      const authStore = useAuthStore()
      authStore.user = { roleId: 'admin' } as any

      const next = vi.fn()
      const guard = requireRole(['admin', 'manager'])
      guard({} as any, {} as any, next)

      expect(next).toHaveBeenCalledWith()
    })

    it('should redirect user with unauthorized role to /forbidden', () => {
      const authStore = useAuthStore()
      authStore.user = { roleId: 'cashier' } as any

      const next = vi.fn()
      const guard = requireRole(['admin', 'finance_manager'])
      guard({} as any, {} as any, next)

      expect(next).toHaveBeenCalledWith('/forbidden')
    })
  })
})
