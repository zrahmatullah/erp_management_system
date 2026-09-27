import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePOSStore } from '@/stores/pos.store'

describe('usePOSStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default state', () => {
    const store = usePOSStore()
    expect(store.cartItems).toEqual([])
    expect(store.cartItemCount).toBe(0)
    expect(store.subtotal).toBe(0)
    expect(store.grandTotal).toBe(0)
    expect(store.orderType).toBe(0)
    expect(store.selectedTable).toBeNull()
  })

  it('should add new item to cart', () => {
    const store = usePOSStore()
    store.addToCart({
      productId: 'prod-1',
      quantity: 2,
      unitPrice: 25000,
      notes: ''
    })

    expect(store.cartItems.length).toBe(1)
    expect(store.cartItems[0].productId).toBe('prod-1')
    expect(store.cartItems[0].quantity).toBe(2)
    expect(store.cartItemCount).toBe(2)
  })

  it('should increment quantity when existing item is added', () => {
    const store = usePOSStore()
    store.addToCart({ productId: 'prod-1', quantity: 1, unitPrice: 25000 })
    store.addToCart({ productId: 'prod-1', quantity: 3, unitPrice: 25000 })

    expect(store.cartItems.length).toBe(1)
    expect(store.cartItems[0].quantity).toBe(4)
    expect(store.cartItemCount).toBe(4)
  })

  it('should remove item from cart by productId', () => {
    const store = usePOSStore()
    store.addToCart({ productId: 'prod-1', quantity: 1, unitPrice: 25000 })
    store.addToCart({ productId: 'prod-2', quantity: 2, unitPrice: 30000 })

    store.removeFromCart('prod-1')
    expect(store.cartItems.length).toBe(1)
    expect(store.cartItems[0].productId).toBe('prod-2')
  })

  it('should update quantity of existing item', () => {
    const store = usePOSStore()
    store.addToCart({ productId: 'prod-1', quantity: 1, unitPrice: 25000 })
    store.updateQuantity('prod-1', 5)

    expect(store.cartItems[0].quantity).toBe(5)
    expect(store.cartItemCount).toBe(5)
  })

  it('should clear cart and reset discount rate', () => {
    const store = usePOSStore()
    store.addToCart({ productId: 'prod-1', quantity: 2, unitPrice: 25000 })
    store.applyDiscount(0.1)

    store.clearCart()
    expect(store.cartItems).toEqual([])
    expect(store.cartItemCount).toBe(0)
    expect(store.subtotal).toBe(0)
  })

  it('should calculate subtotal, discount, tax, and grandTotal correctly', () => {
    const store = usePOSStore()
    // 2 items * 10000 = 20000 subtotal
    store.addToCart({ productId: 'prod-1', quantity: 2, unitPrice: 10000 })
    expect(store.subtotal).toBe(20000)

    // 10% discount -> 2000
    store.applyDiscount(0.10)
    expect(store.discountAmount).toBe(2000)

    // Taxable = 20000 - 2000 = 18000; Tax 11% = 1980
    expect(store.taxAmount).toBe(1980)

    // Grand total = 18000 + 1980 = 19980
    expect(store.grandTotal).toBe(19980)
  })

  it('should set order type and select table', () => {
    const store = usePOSStore()
    store.setOrderType(1)
    expect(store.orderType).toBe(1)

    store.selectTable('table-12')
    expect(store.selectedTable).toBe('table-12')
  })
})
