import { describe, it, expect, beforeEach } from 'vitest'
import { useAuth } from '@/composables/useAuth'

describe('useAuth composable', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('should initialize and handle login and logout correctly', () => {
    const { user, isAuthenticated, login, logout } = useAuth()

    expect(isAuthenticated.value).toBe(false)
    expect(user.value).toBeNull()

    const mockUser = { id: 'u-1', username: 'owner_cafe' } as any
    login('secret-token-xyz', mockUser)

    expect(isAuthenticated.value).toBe(true)
    expect(user.value).toEqual(mockUser)
    expect(localStorage.getItem('token')).toBe('secret-token-xyz')

    logout()
    expect(isAuthenticated.value).toBe(false)
    expect(user.value).toBeNull()
    expect(localStorage.getItem('token')).toBeNull()
  })
})
