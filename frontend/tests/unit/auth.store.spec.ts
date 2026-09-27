import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth.store'
import type { User, Permission } from '@/types/auth'

describe('useAuthStore', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
  })

  it('should initialize with null user and token when localStorage is empty', () => {
    const store = useAuthStore()
    expect(store.user).toBeNull()
    expect(store.token).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(store.fullName).toBe('')
    expect(store.currentRole).toBe('')
  })

  it('should successfully login and persist data in localStorage', () => {
    const store = useAuthStore()
    const mockUser: User = {
      id: 'usr-001',
      username: 'budi_kasir',
      email: 'budi@cafe.com',
      firstName: 'Budi',
      lastName: 'Santoso',
      roleId: 'cashier',
      branchId: 'br-01',
      isActive: true,
      createdAt: '2026-01-01',
      updatedAt: '2026-01-01'
    }
    const mockPerms: Permission[] = [
      { id: 'p1', module: 'pos', action: 'create', description: 'Create POS Order' },
      { id: 'p2', module: 'pos', action: 'read', description: 'Read POS Order' }
    ]

    store.login('jwt-access-token-123', 'refresh-token-456', mockUser, mockPerms)

    expect(store.isAuthenticated).toBe(true)
    expect(store.token).toBe('jwt-access-token-123')
    expect(store.refreshTokenValue).toBe('refresh-token-456')
    expect(store.fullName).toBe('Budi Santoso')
    expect(store.currentRole).toBe('cashier')

    // localStorage check
    expect(localStorage.getItem('token')).toBe('jwt-access-token-123')
    expect(localStorage.getItem('refreshToken')).toBe('refresh-token-456')
    expect(JSON.parse(localStorage.getItem('user')!)).toEqual(mockUser)
  })

  it('should correctly evaluate hasPermission', () => {
    const store = useAuthStore()
    const mockPerms: Permission[] = [
      { id: 'p1', module: 'finance', action: 'approve', description: 'Approve expense' }
    ]
    store.permissions = mockPerms

    expect(store.hasPermission('finance', 'approve')).toBe(true)
    expect(store.hasPermission('finance', 'delete')).toBe(false)
    expect(store.hasPermission('hris', 'read')).toBe(false)
  })

  it('should clear all auth states and localStorage upon logout', () => {
    const store = useAuthStore()
    store.token = 'existing-token'
    store.user = { id: 'usr-1', username: 'admin' } as any
    localStorage.setItem('token', 'existing-token')

    store.logout()

    expect(store.token).toBeNull()
    expect(store.user).toBeNull()
    expect(store.permissions).toEqual([])
    expect(store.isAuthenticated).toBe(false)
    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('user')).toBeNull()
  })

  it('should update user via setUser', () => {
    const store = useAuthStore()
    const updatedUser = { id: 'usr-1', firstName: 'Jane', lastName: 'Doe' } as any
    store.setUser(updatedUser)
    expect(store.user).toEqual(updatedUser)
    expect(store.fullName).toBe('Jane Doe')
  })
})
