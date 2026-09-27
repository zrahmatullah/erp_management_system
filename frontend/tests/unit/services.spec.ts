import { describe, it, expect, vi } from 'vitest'
import * as authService from '@/services/auth.service'
import * as orderService from '@/services/order.service'
import * as roleService from '@/services/role.service'
import * as userService from '@/services/user.service'
import api from '@/plugins/axios'

vi.mock('@/plugins/axios', () => {
  return {
    default: {
      get: vi.fn().mockResolvedValue({ data: { success: true } }),
      post: vi.fn().mockResolvedValue({ data: { success: true } }),
      put: vi.fn().mockResolvedValue({ data: { success: true } }),
      delete: vi.fn().mockResolvedValue({ data: { success: true } })
    }
  }
})

describe('API Services', () => {
  describe('auth.service', () => {
    it('should call auth endpoints', async () => {
      await authService.login({ username: 'admin', password: '123' })
      expect(api.post).toHaveBeenCalledWith('/auth/login', { username: 'admin', password: '123' })

      await authService.logout()
      expect(api.post).toHaveBeenCalledWith('/auth/logout')

      await authService.refreshToken({ refresh_token: 'xyz' })
      expect(api.post).toHaveBeenCalledWith('/auth/refresh', { refresh_token: 'xyz' })

      await authService.forgotPassword({ email: 'test@cafe.com' })
      expect(api.post).toHaveBeenCalledWith('/auth/forgot-password', { email: 'test@cafe.com' })

      await authService.resetPassword({ token: 't1', new_password: 'p1' })
      expect(api.post).toHaveBeenCalledWith('/auth/reset-password', { token: 't1', new_password: 'p1' })

      await authService.getProfile()
      expect(api.get).toHaveBeenCalledWith('/auth/profile')
    })
  })

  describe('order.service', () => {
    it('should call order endpoints', async () => {
      await orderService.createOrder({ total: 50000 })
      expect(api.post).toHaveBeenCalledWith('/orders', { total: 50000 })

      await orderService.getOrders({ status: 'completed' })
      expect(api.get).toHaveBeenCalledWith('/orders', { params: { status: 'completed' } })

      await orderService.getOrderById('ord-1')
      expect(api.get).toHaveBeenCalledWith('/orders/ord-1')

      await orderService.updateOrderStatus('ord-1', 'paid')
      expect(api.put).toHaveBeenCalledWith('/orders/ord-1/status', { status: 'paid' })

      await orderService.cancelOrder('ord-1')
      expect(api.post).toHaveBeenCalledWith('/orders/ord-1/cancel')

      await orderService.voidOrder('ord-1')
      expect(api.post).toHaveBeenCalledWith('/orders/ord-1/void')

      await orderService.getKitchenQueue()
      expect(api.get).toHaveBeenCalledWith('/orders/kitchen-queue')
    })
  })

  describe('role.service', () => {
    it('should call role endpoints', async () => {
      await roleService.getRoles()
      expect(api.get).toHaveBeenCalledWith('/roles')

      await roleService.createRole({ name: 'Barista' })
      expect(api.post).toHaveBeenCalledWith('/roles', { name: 'Barista' })

      await roleService.updateRole('r1', { name: 'Lead Barista' })
      expect(api.put).toHaveBeenCalledWith('/roles/r1', { name: 'Lead Barista' })

      await roleService.deleteRole('r1')
      expect(api.delete).toHaveBeenCalledWith('/roles/r1')

      await roleService.getPermissions()
      expect(api.get).toHaveBeenCalledWith('/permissions')

      await roleService.assignPermissions('r1', { permissions: ['p1'] })
      expect(api.post).toHaveBeenCalledWith('/roles/r1/permissions', { permissions: ['p1'] })
    })
  })

  describe('user.service', () => {
    it('should call user endpoints', async () => {
      await userService.getUsers({ branch_id: 'b1' })
      expect(api.get).toHaveBeenCalledWith('/users', { params: { branch_id: 'b1' } })

      await userService.getUserById('u1')
      expect(api.get).toHaveBeenCalledWith('/users/u1')

      await userService.createUser({ username: 'cashier1' })
      expect(api.post).toHaveBeenCalledWith('/users', { username: 'cashier1' })

      await userService.updateUser('u1', { first_name: 'John' })
      expect(api.put).toHaveBeenCalledWith('/users/u1', { first_name: 'John' })

      await userService.deleteUser('u1')
      expect(api.delete).toHaveBeenCalledWith('/users/u1')

      await userService.assignRole('u1', { role_id: 'r1' })
      expect(api.post).toHaveBeenCalledWith('/users/u1/role', { role_id: 'r1' })
    })
  })
})
