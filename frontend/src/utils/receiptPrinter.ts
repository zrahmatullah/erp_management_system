export interface ReceiptData {
  order_number: string
  queue_number?: string
  table_number?: string
  customer_name?: string
  order_type?: string
  date?: string
  paid_at?: string
  cashier_name?: string
  items: Array<{
    name: string
    quantity: number
    unit_price?: number
    price?: number
    total_price?: number
  }>
  subtotal: number
  tax: number
  discount?: number
  total: number
  payment_method: string
  amount_paid: number
  change_due: number
}

export function generateReceiptHtml(data: ReceiptData): string {
  const itemsHtml = (data.items || []).map(item => {
    const unitPrice = item.unit_price || item.price || 0
    const itemTotal = item.total_price || (unitPrice * item.quantity)
    return `
      <tr>
        <td style="padding: 3px 0; text-align: left; vertical-align: top;">
          <div style="font-weight: bold; color: #111;">${item.name}</div>
          <div style="font-size: 11px; color: #666;">${item.quantity} x Rp ${unitPrice.toLocaleString('id-ID')}</div>
        </td>
        <td style="padding: 3px 0; text-align: right; vertical-align: bottom; font-weight: bold; color: #111;">
          Rp ${itemTotal.toLocaleString('id-ID')}
        </td>
      </tr>
    `
  }).join('')

  const now = new Date()
  const displayDate = data.paid_at || data.date || now.toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short'
  })

  const queueBadge = data.queue_number && data.queue_number !== '-'
    ? `<div style="display: inline-block; background: #1e1b4b; color: #fff; padding: 4px 10px; border-radius: 6px; font-size: 14px; font-weight: 900; margin-top: 4px;">
        NO. ANTRIAN: #${data.queue_number}
       </div>`
    : ''

  const tableOrType = data.order_type === 'dine_in' || data.order_type === 'Dine-in'
    ? `Meja: ${data.table_number || '-'}`
    : `Tipe: ${data.order_type || 'Takeaway'}`

  return `
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>Struk Pembayaran - ${data.order_number}</title>
      <style>
        @page {
          margin: 0;
          size: 80mm auto;
        }
        body {
          font-family: 'Courier New', Courier, monospace, system-ui;
          font-size: 12px;
          line-height: 1.4;
          color: #000;
          background: #fff;
          margin: 0;
          padding: 12px;
          max-width: 320px;
          margin: 0 auto;
        }
        .text-center { text-align: center; }
        .text-right { text-align: right; }
        .font-bold { font-weight: bold; }
        .border-top { border-top: 1px dashed #444; }
        .border-bottom { border-bottom: 1px dashed #444; }
        .border-double { border-top: 2px solid #000; border-bottom: 2px solid #000; }
        table { width: 100%; border-collapse: collapse; }
        .barcode {
          letter-spacing: 4px;
          font-size: 14px;
          font-weight: bold;
          margin-top: 8px;
        }
      </style>
    </head>
    <body>
      <div class="text-center" style="margin-bottom: 12px;">
        <div style="font-size: 16px; font-weight: 900; letter-spacing: 1px;">CAFE ERP SYSTEM</div>
        <div style="font-size: 11px; color: #555;">Kopi Premium & Kitchen</div>
        <div style="font-size: 10px; color: #777;">Jl. Senopati Raya No. 45, Jakarta Selatan</div>
        <div style="font-size: 10px; color: #777;">Telp: (021) 7289-0012</div>
        ${queueBadge}
      </div>

      <div class="border-top border-bottom" style="padding: 6px 0; margin-bottom: 10px; font-size: 11px;">
        <div style="display: flex; justify-content: space-between;">
          <span>No. Order: <strong>${data.order_number}</strong></span>
          <span>${tableOrType}</span>
        </div>
        <div style="display: flex; justify-content: space-between; margin-top: 2px;">
          <span>Tamu: ${data.customer_name || 'Pelanggan'}</span>
          <span>Kasir: ${data.cashier_name || 'Jane Cashier'}</span>
        </div>
        <div style="color: #666; font-size: 10px; margin-top: 2px;">
          ${displayDate}
        </div>
      </div>

      <div style="margin-bottom: 10px;">
        <table>
          <tbody>
            ${itemsHtml}
          </tbody>
        </table>
      </div>

      <div class="border-top" style="padding-top: 6px; font-size: 11px; margin-bottom: 8px;">
        <div style="display: flex; justify-content: space-between; margin-bottom: 2px;">
          <span>Subtotal</span>
          <span>Rp ${(data.subtotal || 0).toLocaleString('id-ID')}</span>
        </div>
        ${data.discount ? `
        <div style="display: flex; justify-content: space-between; margin-bottom: 2px; color: #15803d;">
          <span>Diskon</span>
          <span>-Rp ${(data.discount).toLocaleString('id-ID')}</span>
        </div>` : ''}
        <div style="display: flex; justify-content: space-between; margin-bottom: 4px;">
          <span>PPN (10%)</span>
          <span>Rp ${(data.tax || 0).toLocaleString('id-ID')}</span>
        </div>
        <div class="border-double" style="display: flex; justify-content: space-between; padding: 6px 0; font-size: 14px; font-weight: 900;">
          <span>TOTAL</span>
          <span>Rp ${(data.total || 0).toLocaleString('id-ID')}</span>
        </div>
      </div>

      <div style="font-size: 11px; margin-bottom: 14px;">
        <div style="display: flex; justify-content: space-between; margin-bottom: 2px;">
          <span>Metode Bayar:</span>
          <strong style="text-transform: uppercase;">${data.payment_method || 'TUNAI'}</strong>
        </div>
        <div style="display: flex; justify-content: space-between; margin-bottom: 2px;">
          <span>Bayar (Tendered):</span>
          <span>Rp ${(data.amount_paid || data.total || 0).toLocaleString('id-ID')}</span>
        </div>
        <div style="display: flex; justify-content: space-between; font-weight: bold;">
          <span>Kembalian:</span>
          <span>Rp ${(data.change_due || 0).toLocaleString('id-ID')}</span>
        </div>
      </div>

      <div class="border-top text-center" style="padding-top: 10px; color: #444;">
        <div style="font-weight: bold; font-size: 11px;">TERIMA KASIH ATAS KUNJUNGAN ANDA</div>
        <div style="font-size: 10px; color: #777; margin-top: 2px;">Simpan struk ini sebagai bukti pembayaran yang sah.</div>
        <div class="barcode">||| | ||||| || |||| ||| |||</div>
        <div style="font-size: 9px; color: #888; margin-top: 4px;">WiFi: CafeERP-Guest | Pass: kopienak2026</div>
      </div>
    </body>
    </html>
  `
}

/**
 * Generates an HTML Blob URL for the receipt
 */
export function createReceiptBlobUrl(data: ReceiptData): string {
  const html = generateReceiptHtml(data)
  const blob = new Blob([html], { type: 'text/html;charset=utf-8' })
  return URL.createObjectURL(blob)
}

/**
 * Prints the receipt automatically using a hidden iframe loaded with the Blob URL
 */
export function printReceiptBlob(data: ReceiptData): Promise<void> {
  return new Promise((resolve) => {
    try {
      const blobUrl = createReceiptBlobUrl(data)
      const iframe = document.createElement('iframe')
      iframe.style.position = 'fixed'
      iframe.style.right = '0'
      iframe.style.bottom = '0'
      iframe.style.width = '0'
      iframe.style.height = '0'
      iframe.style.border = 'none'
      iframe.src = blobUrl

      iframe.onload = () => {
        setTimeout(() => {
          try {
            iframe.contentWindow?.focus()
            iframe.contentWindow?.print()
          } catch (e) {
            console.error('Iframe print error:', e)
          } finally {
            setTimeout(() => {
              document.body.removeChild(iframe)
              URL.revokeObjectURL(blobUrl)
              resolve()
            }, 1000)
          }
        }, 300)
      }

      document.body.appendChild(iframe)
    } catch (err) {
      console.error('Failed to trigger receipt blob print:', err)
      resolve()
    }
  })
}

