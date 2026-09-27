import { describe, it, expect, vi } from 'vitest'
import { generateReceiptHtml, createReceiptBlobUrl, printReceiptBlob, type ReceiptData } from '@/utils/receiptPrinter'

describe('receiptPrinter utility', () => {
  const sampleData: ReceiptData = {
    order_number: 'ORD-2026-0927-001',
    queue_number: '12',
    table_number: 'T-04',
    order_type: 'dine_in',
    customer_name: 'Dewi Lestari',
    paid_at: '27 Sep 2026, 14:30',
    cashier_name: 'Budi Santoso',
    items: [
      { name: 'Kopi Susu Gula Aren', quantity: 2, unit_price: 25000, total_price: 50000 },
      { name: 'Croissant Butter', quantity: 1, unit_price: 28000, total_price: 28000 }
    ],
    subtotal: 78000,
    discount: 7800,
    tax: 7722,
    total: 77922,
    payment_method: 'QRIS',
    amount_paid: 77922,
    change_due: 0
  }

  it('should generate well-structured HTML receipt for dine-in order', () => {
    const html = generateReceiptHtml(sampleData)
    expect(html).toContain('CAFE ERP SYSTEM')
    expect(html).toContain('ORD-2026-0927-001')
    expect(html).toContain('NO. ANTRIAN: #12')
    expect(html).toContain('Meja: T-04')
    expect(html).toContain('Dewi Lestari')
    expect(html).toContain('Kopi Susu Gula Aren')
    expect(html).toContain('Croissant Butter')
    expect(html).toContain('QRIS')
    expect(html).toContain('Rp 77.922')
    expect(html).toContain('Budi Santoso')
  })

  it('should handle takeaway order without table and queue number', () => {
    const takeawayData: ReceiptData = {
      order_number: 'ORD-2026-0927-002',
      queue_number: '-',
      order_type: 'takeaway',
      items: [
        { name: 'Americano Ice', quantity: 1, price: 20000 }
      ],
      subtotal: 20000,
      tax: 2000,
      total: 22000,
      payment_method: 'Cash',
      amount_paid: 50000,
      change_due: 28000
    }

    const html = generateReceiptHtml(takeawayData)
    expect(html).toContain('ORD-2026-0927-002')
    expect(html).toContain('Tipe: takeaway')
    expect(html).not.toContain('NO. ANTRIAN: #-')
    expect(html).toContain('Americano Ice')
    expect(html).toContain('Cash')
  })

  it('should create blob url for receipt data', () => {
    global.URL.createObjectURL = vi.fn().mockReturnValue('blob:http://localhost:5173/mock-receipt-blob')
    const blobUrl = createReceiptBlobUrl(sampleData)
    expect(blobUrl).toBe('blob:http://localhost:5173/mock-receipt-blob')
  })

  it('should trigger printReceiptBlob via hidden iframe', async () => {
    global.URL.createObjectURL = vi.fn().mockReturnValue('blob:http://localhost:5173/mock-receipt-blob')
    global.URL.revokeObjectURL = vi.fn()

    const printPromise = printReceiptBlob(sampleData)
    expect(printPromise).toBeInstanceOf(Promise)
  })
})
