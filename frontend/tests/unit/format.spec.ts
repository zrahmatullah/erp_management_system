import { describe, it, expect } from 'vitest'
import { formatCurrency, formatDate } from '@/utils/format'

describe('format utilities', () => {
  describe('formatCurrency', () => {
    it('should format positive integer to IDR currency', () => {
      const result = formatCurrency(50000)
      expect(result).toMatch(/Rp\s*50\.000/)
    })

    it('should format zero correctly', () => {
      const result = formatCurrency(0)
      expect(result).toMatch(/Rp\s*0/)
    })

    it('should format negative amount correctly', () => {
      const result = formatCurrency(-25000)
      expect(result).toContain('25.000')
    })

    it('should format large millions and billions amount', () => {
      const result = formatCurrency(1500000000)
      expect(result).toMatch(/1\.500\.000\.000/)
    })
  })

  describe('formatDate', () => {
    it('should format ISO date string to localized date', () => {
      const formatted = formatDate('2026-09-27')
      expect(formatted).toBeDefined()
      expect(typeof formatted).toBe('string')
      expect(formatted).toContain('2026')
    })

    it('should format datetime string properly', () => {
      const formatted = formatDate('2026-01-15T10:30:00Z')
      expect(formatted).toBeDefined()
      expect(formatted).toContain('2026')
    })
  })
})
