import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { CartItem, OrderType, Order } from '@/types/order';

export const usePOSStore = defineStore('pos', () => {
  const cartItems = ref<CartItem[]>([]);
  const currentOrder = ref<Order | null>(null);
  const orderType = ref<OrderType>(0); // 0 = DINE_IN
  const selectedTable = ref<string | null>(null);
  const activeOrders = ref<Order[]>([]);
  const discountRate = ref(0);

  const cartItemCount = computed(() => cartItems.value.reduce((total, item) => total + item.quantity, 0));
  const subtotal = computed(() => cartItems.value.reduce((total, item) => total + (10000 * item.quantity), 0)); // Placeholder price
  const discountAmount = computed(() => subtotal.value * discountRate.value);
  const taxAmount = computed(() => (subtotal.value - discountAmount.value) * 0.11);
  const grandTotal = computed(() => subtotal.value - discountAmount.value + taxAmount.value);

  const addToCart = (item: CartItem) => {
    const existing = cartItems.value.find(i => i.productId === item.productId);
    if (existing) {
      existing.quantity += item.quantity;
    } else {
      cartItems.value.push(item);
    }
  };

  const removeFromCart = (productId: string) => {
    cartItems.value = cartItems.value.filter(i => i.productId !== productId);
  };

  const updateQuantity = (productId: string, quantity: number) => {
    const item = cartItems.value.find(i => i.productId === productId);
    if (item) item.quantity = quantity;
  };

  const clearCart = () => {
    cartItems.value = [];
    discountRate.value = 0;
  };

  const setOrderType = (type: OrderType) => {
    orderType.value = type;
  };

  const selectTable = (tableId: string) => {
    selectedTable.value = tableId;
  };

  const applyDiscount = (rate: number) => {
    discountRate.value = rate;
  };

  return {
    cartItems, currentOrder, orderType, selectedTable, activeOrders,
    cartItemCount, subtotal, discountAmount, taxAmount, grandTotal,
    addToCart, removeFromCart, updateQuantity, clearCart, setOrderType, selectTable, applyDiscount
  };
});
