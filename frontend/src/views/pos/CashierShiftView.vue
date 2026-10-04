<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-white p-6 rounded-2xl shadow-sm border border-stone-100">
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-2xl font-bold text-stone-900">Shift Kasir & Cash Balancing</h1>
          <span
            v-if="currentShift"
            class="px-3 py-1 text-xs font-semibold rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 flex items-center gap-1.5"
          >
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            Shift Aktif: {{ currentShift.shift_name }}
          </span>
          <span
            v-else
            class="px-3 py-1 text-xs font-semibold rounded-full bg-amber-50 text-amber-700 border border-amber-200"
          >
            Tidak Ada Shift Berjalan
          </span>
        </div>
        <p class="text-sm text-stone-500 mt-1">
          Kontrol modal awal kasir, mutasi setor brankas (cash drop), dan audit blind closing harian.
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <button
          v-if="!currentShift"
          @click="showOpenModal = true"
          class="px-4 py-2 bg-amber-600 hover:bg-amber-700 text-white rounded-xl font-medium text-sm transition-colors shadow-sm flex items-center gap-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Buka Shift Baru
        </button>

        <template v-else>
          <button
            @click="showMovementModal = true"
            class="px-4 py-2 bg-stone-100 hover:bg-stone-200 text-stone-700 rounded-xl font-medium text-sm transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/>
            </svg>
            Setor Kas / Drop
          </button>

          <button
            @click="showCloseModal = true"
            class="px-4 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-xl font-medium text-sm transition-colors shadow-sm flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            Tutup Shift (Blind Count)
          </button>
        </template>
      </div>
    </div>

    <!-- Active Shift Metrics Cards -->
    <div v-if="currentShift" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white p-5 rounded-2xl border border-stone-100 shadow-sm">
        <span class="text-xs font-medium text-stone-500 uppercase tracking-wider">Modal Awal Kasir</span>
        <div class="mt-2 text-2xl font-bold text-stone-900">{{ formatRupiah(currentShift.opening_cash_float) }}</div>
        <div class="text-xs text-stone-400 mt-1">Dibuka pada: {{ formatDateTime(currentShift.opened_at) }}</div>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-stone-100 shadow-sm">
        <span class="text-xs font-medium text-stone-500 uppercase tracking-wider">Penjualan Tunai (Cash)</span>
        <div class="mt-2 text-2xl font-bold text-emerald-600">{{ formatRupiah(currentShift.total_cash_sales) }}</div>
        <div class="text-xs text-stone-400 mt-1">Total tunai bersih dari transaksi POS</div>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-stone-100 shadow-sm">
        <span class="text-xs font-medium text-stone-500 uppercase tracking-wider">Setor Kas / Drop</span>
        <div class="mt-2 text-2xl font-bold text-amber-600">-{{ formatRupiah(currentShift.total_cash_drops) }}</div>
        <div class="text-xs text-stone-400 mt-1">Uang tunai disetor ke brankas</div>
      </div>

      <div class="bg-gradient-to-br from-amber-500 to-amber-700 p-5 rounded-2xl text-white shadow-md">
        <span class="text-xs font-medium text-amber-100 uppercase tracking-wider">Estimasi Kas di Laci</span>
        <div class="mt-2 text-2xl font-bold">{{ formatRupiah(currentShift.expected_current_cash) }}</div>
        <div class="text-xs text-amber-100 mt-1">Modal + Penjualan Cash - Drop</div>
      </div>
    </div>

    <!-- Live Movements Table (If Shift Active) -->
    <div v-if="currentShift" class="bg-white rounded-2xl border border-stone-100 shadow-sm overflow-hidden p-6">
      <h2 class="text-base font-bold text-stone-900 mb-4 flex items-center justify-between">
        <span>Log Mutasi Kas Shift Ini</span>
        <span class="text-xs text-stone-400 font-normal">Riwayat setor brankas & penarikan kas kecil</span>
      </h2>

      <div v-if="!currentShift.recent_movements || currentShift.recent_movements.length === 0" class="text-center py-8 text-stone-400 text-sm">
        Belum ada mutasi setor kas atau penarikan pada shift ini.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-stone-50 text-stone-600 text-xs uppercase">
            <tr>
              <th class="px-4 py-3">Waktu</th>
              <th class="px-4 py-3">Tipe Mutasi</th>
              <th class="px-4 py-3">Nominal</th>
              <th class="px-4 py-3">Keterangan / Alasan</th>
              <th class="px-4 py-3">Otorisasi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-stone-100">
            <tr v-for="m in currentShift.recent_movements" :key="m.id" class="hover:bg-stone-50/50">
              <td class="px-4 py-3 text-stone-500">{{ formatDateTime(m.created_at) }}</td>
              <td class="px-4 py-3">
                <span
                  class="px-2 py-0.5 text-xs font-medium rounded-full"
                  :class="m.movement_type === 'cash_drop' ? 'bg-amber-100 text-amber-800' : 'bg-stone-100 text-stone-700'"
                >
                  {{ m.movement_type === 'cash_drop' ? 'Setor Brankas' : m.movement_type }}
                </span>
              </td>
              <td class="px-4 py-3 font-semibold text-stone-900">{{ formatRupiah(m.amount) }}</td>
              <td class="px-4 py-3 text-stone-600">{{ m.reason }}</td>
              <td class="px-4 py-3 text-stone-500">{{ m.authorized_by_name || '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Shift History Table -->
    <div class="bg-white rounded-2xl border border-stone-100 shadow-sm overflow-hidden p-6">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-5">
        <div>
          <h2 class="text-lg font-bold text-stone-900">Riwayat Shift Kasir & Audit Z-Report</h2>
          <p class="text-xs text-stone-500 mt-0.5">Daftar rekonsiliasi kasir, selisih audit, dan cetak laporan Z-Report.</p>
        </div>

        <button
          @click="fetchShifts"
          class="px-3 py-1.5 text-xs text-stone-600 bg-stone-100 hover:bg-stone-200 rounded-lg transition-colors flex items-center gap-1.5"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh Data
        </button>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-stone-50 text-stone-600 text-xs uppercase font-medium">
            <tr>
              <th class="px-4 py-3">Kasir & Shift</th>
              <th class="px-4 py-3">Waktu Buka / Tutup</th>
              <th class="px-4 py-3 text-right">Modal Awal</th>
              <th class="px-4 py-3 text-right">Penjualan Cash</th>
              <th class="px-4 py-3 text-right">Uang Fisik Dihitung</th>
              <th class="px-4 py-3 text-right">Selisih Kas</th>
              <th class="px-4 py-3 text-center">Status</th>
              <th class="px-4 py-3 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-stone-100">
            <tr v-if="shifts.length === 0">
              <td colspan="8" class="text-center py-8 text-stone-400">Belum ada data shift kasir.</td>
            </tr>
            <tr v-for="s in shifts" :key="s.id" class="hover:bg-stone-50/50">
              <td class="px-4 py-3">
                <div class="font-medium text-stone-900">{{ s.cashier_name }}</div>
                <div class="text-xs text-stone-400">{{ s.shift_name }} • {{ s.branch_name }}</div>
              </td>
              <td class="px-4 py-3 text-xs text-stone-500">
                <div>Buka: {{ formatDateTime(s.opened_at) }}</div>
                <div v-if="s.closed_at">Tutup: {{ formatDateTime(s.closed_at) }}</div>
              </td>
              <td class="px-4 py-3 text-right font-medium text-stone-700">{{ formatRupiah(s.opening_cash_float) }}</td>
              <td class="px-4 py-3 text-right font-medium text-emerald-600">{{ formatRupiah(s.total_cash_sales) }}</td>
              <td class="px-4 py-3 text-right font-medium text-stone-900">
                {{ s.status === 'closed' ? formatRupiah(s.actual_cash_counted) : '-' }}
              </td>
              <td class="px-4 py-3 text-right">
                <span
                  v-if="s.status === 'closed'"
                  class="font-semibold text-xs px-2 py-0.5 rounded-full inline-block"
                  :class="{
                    'bg-emerald-100 text-emerald-800': s.cash_difference === 0,
                    'bg-rose-100 text-rose-800': s.cash_difference < 0,
                    'bg-blue-100 text-blue-800': s.cash_difference > 0
                  }"
                >
                  {{ s.cash_difference > 0 ? '+' : '' }}{{ formatRupiah(s.cash_difference) }}
                </span>
                <span v-else class="text-xs text-stone-400">-</span>
              </td>
              <td class="px-4 py-3 text-center">
                <span
                  class="px-2.5 py-1 text-xs font-semibold rounded-full"
                  :class="s.status === 'open' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-stone-100 text-stone-600'"
                >
                  {{ s.status === 'open' ? 'Berjalan' : 'Selesai' }}
                </span>
              </td>
              <td class="px-4 py-3 text-center">
                <button
                  @click="openSummaryModal(s.id)"
                  class="px-2.5 py-1 text-xs text-amber-700 hover:text-amber-800 bg-amber-50 hover:bg-amber-100 rounded-lg transition-colors font-medium"
                >
                  Z-Report
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL: Buka Shift Baru -->
    <div v-if="showOpenModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm">
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-xl border border-stone-100">
        <h3 class="text-lg font-bold text-stone-900">Buka Shift Kasir Baru</h3>
        <p class="text-xs text-stone-500 mt-1">Masukkan modal awal uang tunai (cash float) yang tersedia di laci kasir.</p>

        <form @submit.prevent="submitOpenShift" class="mt-4 space-y-4">
          <div>
            <label class="block text-xs font-semibold text-stone-700 uppercase mb-1">Nama Sesi Shift</label>
            <select v-model="openForm.shift_name" class="w-full px-3 py-2 border rounded-xl text-sm focus:ring-2 focus:ring-amber-500 outline-none">
              <option value="Shift Pagi">Shift Pagi</option>
              <option value="Shift Siang">Shift Siang</option>
              <option value="Shift Sore">Shift Sore</option>
              <option value="Shift Malam">Shift Malam</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-stone-700 uppercase mb-1">Modal Awal Tunai (Cash Float)</label>
            <div class="relative">
              <span class="absolute left-3 top-2.5 text-stone-400 text-sm">Rp</span>
              <input
                v-model.number="openForm.opening_cash_float"
                type="number"
                min="0"
                step="1000"
                required
                class="w-full pl-10 pr-3 py-2 border rounded-xl text-sm focus:ring-2 focus:ring-amber-500 outline-none font-semibold"
                placeholder="Contoh: 300000"
              />
            </div>
            <span class="text-xs text-stone-400 mt-1 block">Uang kembalian yang disiapkan sebelum buka kasir.</span>
          </div>

          <div>
            <label class="block text-xs font-semibold text-stone-700 uppercase mb-1">Catatan Tambahan (Opsional)</label>
            <textarea
              v-model="openForm.notes"
              rows="2"
              class="w-full px-3 py-2 border rounded-xl text-sm focus:ring-2 focus:ring-amber-500 outline-none"
              placeholder="Contoh: Termasuk pecahan 5.000 dan 10.000"
            ></textarea>
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button
              type="button"
              @click="showOpenModal = false"
              class="px-4 py-2 text-stone-600 hover:bg-stone-100 rounded-xl text-sm font-medium"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="loadingAction"
              class="px-5 py-2 bg-amber-600 hover:bg-amber-700 text-white rounded-xl text-sm font-semibold shadow-sm transition-colors disabled:opacity-50"
            >
              {{ loadingAction ? 'Membuka...' : 'Mulai Shift' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: Setor Kas / Drop -->
    <div v-if="showMovementModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm">
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-xl border border-stone-100">
        <h3 class="text-lg font-bold text-stone-900">Setor Kas Tengah Hari (Cash Drop)</h3>
        <p class="text-xs text-stone-500 mt-1">Gunakan ini untuk menyetor uang berlebih di laci kasir ke brankas supervisor.</p>

        <form @submit.prevent="submitMovement" class="mt-4 space-y-4">
          <div>
            <label class="block text-xs font-semibold text-stone-700 uppercase mb-1">Tipe Mutasi</label>
            <select v-model="movementForm.movement_type" class="w-full px-3 py-2 border rounded-xl text-sm focus:ring-2 focus:ring-amber-500 outline-none">
              <option value="cash_drop">Setor Kas ke Brankas (Cash Drop)</option>
              <option value="paid_out">Kas Keluar Operasional (Paid Out)</option>
              <option value="cash_in">Penambahan Modal Kas (Cash In)</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-semibold text-stone-700 uppercase mb-1">Nominal Setoran</label>
            <div class="relative">
              <span class="absolute left-3 top-2.5 text-stone-400 text-sm">Rp</span>
              <input
                v-model.number="movementForm.amount"
                type="number"
                min="1000"
                step="1000"
                required
                class="w-full pl-10 pr-3 py-2 border rounded-xl text-sm focus:ring-2 focus:ring-amber-500 outline-none font-semibold"
                placeholder="Contoh: 1000000"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-stone-700 uppercase mb-1">Alasan / Keterangan</label>
            <input
              v-model="movementForm.reason"
              type="text"
              required
              class="w-full px-3 py-2 border rounded-xl text-sm focus:ring-2 focus:ring-amber-500 outline-none"
              placeholder="Contoh: Setoran uang pecahan 100k ke brankas supervisor"
            />
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button
              type="button"
              @click="showMovementModal = false"
              class="px-4 py-2 text-stone-600 hover:bg-stone-100 rounded-xl text-sm font-medium"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="loadingAction"
              class="px-5 py-2 bg-amber-600 hover:bg-amber-700 text-white rounded-xl text-sm font-semibold shadow-sm transition-colors disabled:opacity-50"
            >
              {{ loadingAction ? 'Menyimpan...' : 'Simpan Mutasi' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: Tutup Shift (Blind Count) -->
    <div v-if="showCloseModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm">
      <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-xl border border-stone-100">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center font-bold">
            🔒
          </div>
          <div>
            <h3 class="text-lg font-bold text-stone-900">Tutup Shift Kasir (Blind Closing)</h3>
            <p class="text-xs text-stone-500">Hitung seluruh uang fisik (kertas & koin) yang ada di laci kasir saat ini.</p>
          </div>
        </div>

        <div class="mt-4 p-3 bg-amber-50 border border-amber-200 rounded-xl text-xs text-amber-800">
          <strong>Perhatian:</strong> Sistem sengaja tidak menampilkan estimasi saldo kas untuk menjaga objektivitas audit kasir. Masukkan nominal uang riil yang Anda hitung.
        </div>

        <form @submit.prevent="submitCloseShift" class="mt-4 space-y-4">
          <div>
            <label class="block text-xs font-semibold text-stone-700 uppercase mb-1">Total Uang Fisik Hasil Hitung (Actual Cash)</label>
            <div class="relative">
              <span class="absolute left-3 top-2.5 text-stone-400 text-sm">Rp</span>
              <input
                v-model.number="closeForm.actual_cash_counted"
                type="number"
                min="0"
                step="100"
                required
                class="w-full pl-10 pr-3 py-2 border rounded-xl text-base focus:ring-2 focus:ring-rose-500 outline-none font-bold text-stone-900"
                placeholder="0"
              />
            </div>
            <span class="text-xs text-stone-400 mt-1 block">Jumlahkan seluruh uang kertas dan koin di laci kasir.</span>
          </div>

          <div>
            <label class="block text-xs font-semibold text-stone-700 uppercase mb-1">Catatan Penutupan</label>
            <textarea
              v-model="closeForm.notes"
              rows="2"
              class="w-full px-3 py-2 border rounded-xl text-sm focus:ring-2 focus:ring-rose-500 outline-none"
              placeholder="Contoh: Kondisi kas tertata rapi, tidak ada komplain pelanggan"
            ></textarea>
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button
              type="button"
              @click="showCloseModal = false"
              class="px-4 py-2 text-stone-600 hover:bg-stone-100 rounded-xl text-sm font-medium"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="loadingAction"
              class="px-5 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-sm font-semibold shadow-sm transition-colors disabled:opacity-50"
            >
              {{ loadingAction ? 'Merekonsiliasi...' : 'Selesaikan Blind Closing' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: Z-Report Summary View -->
    <div v-if="selectedSummary" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm">
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-xl border border-stone-100 max-h-[90vh] overflow-y-auto">
        <div class="text-center pb-4 border-b border-dashed border-stone-200">
          <h2 class="text-lg font-black uppercase tracking-wider text-stone-900">CAFE ERP SYSTEM</h2>
          <div class="text-xs text-stone-500 mt-0.5">LAPORAN REKONSILIASI Z-REPORT KASIR</div>
          <div class="text-xs text-stone-400 mt-1">{{ selectedSummary.shift.branch_name }} • {{ selectedSummary.shift.shift_name }}</div>
        </div>

        <div class="mt-4 space-y-2 text-xs text-stone-600">
          <div class="flex justify-between">
            <span>Kasir:</span>
            <span class="font-semibold text-stone-900">{{ selectedSummary.shift.cashier_name }}</span>
          </div>
          <div class="flex justify-between">
            <span>Waktu Buka:</span>
            <span>{{ formatDateTime(selectedSummary.shift.opened_at) }}</span>
          </div>
          <div class="flex justify-between">
            <span>Waktu Tutup:</span>
            <span>{{ formatDateTime(selectedSummary.shift.closed_at) }}</span>
          </div>
        </div>

        <div class="my-4 border-t border-dashed border-stone-200 pt-3 space-y-2 text-xs">
          <div class="flex justify-between text-stone-600">
            <span>Modal Awal (Float):</span>
            <span class="font-semibold">{{ formatRupiah(selectedSummary.shift.opening_cash_float) }}</span>
          </div>
          <div class="flex justify-between text-stone-600">
            <span>Penjualan Tunai:</span>
            <span class="font-semibold text-emerald-600">{{ formatRupiah(selectedSummary.shift.total_cash_sales) }}</span>
          </div>
          <div class="flex justify-between text-stone-600">
            <span>Setoran Kas (Drop):</span>
            <span class="font-semibold text-amber-600">-{{ formatRupiah(selectedSummary.shift.total_cash_drops) }}</span>
          </div>
          <div class="flex justify-between text-stone-900 font-bold pt-1 border-t border-stone-100">
            <span>Ekspektasi Kas Sistem:</span>
            <span>{{ formatRupiah(selectedSummary.shift.expected_cash_total) }}</span>
          </div>
          <div class="flex justify-between text-stone-900 font-bold">
            <span>Uang Fisik Dihitung:</span>
            <span>{{ formatRupiah(selectedSummary.shift.actual_cash_counted) }}</span>
          </div>
          <div
            class="flex justify-between font-bold text-sm p-2 rounded-lg mt-2"
            :class="{
              'bg-emerald-50 text-emerald-700': selectedSummary.shift.cash_difference === 0,
              'bg-rose-50 text-rose-700': selectedSummary.shift.cash_difference < 0,
              'bg-blue-50 text-blue-700': selectedSummary.shift.cash_difference > 0
            }"
          >
            <span>Selisih (Discrepancy):</span>
            <span>{{ selectedSummary.shift.cash_difference > 0 ? '+' : '' }}{{ formatRupiah(selectedSummary.shift.cash_difference) }}</span>
          </div>
        </div>

        <div v-if="selectedSummary.payment_breakdown && selectedSummary.payment_breakdown.length > 0" class="my-4 border-t border-dashed border-stone-200 pt-3">
          <div class="text-xs font-bold text-stone-800 uppercase mb-2">Rincian Metode Pembayaran:</div>
          <div v-for="p in selectedSummary.payment_breakdown" :key="p.method" class="flex justify-between text-xs text-stone-600 py-0.5">
            <span>{{ p.method }} ({{ p.count }}x):</span>
            <span class="font-semibold text-stone-900">{{ formatRupiah(p.total) }}</span>
          </div>
          <div class="flex justify-between text-xs font-bold text-stone-900 pt-2 border-t border-stone-100 mt-2">
            <span>Total Omzet Kotor:</span>
            <span>{{ formatRupiah(selectedSummary.gross_total_sales) }}</span>
          </div>
        </div>

        <div class="flex gap-2 pt-4 border-t border-stone-200">
          <button
            @click="selectedSummary = null"
            class="flex-1 py-2 bg-stone-100 hover:bg-stone-200 text-stone-700 rounded-xl text-xs font-semibold"
          >
            Tutup
          </button>
          <button
            @click="printSummary"
            class="flex-1 py-2 bg-amber-600 hover:bg-amber-700 text-white rounded-xl text-xs font-semibold shadow-sm flex items-center justify-center gap-1.5"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"/>
            </svg>
            Cetak Struk
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import {
  getCurrentShift,
  openShift,
  recordCashMovement,
  closeShift,
  getShiftSummary,
  getShifts,
} from '@/services/shift.service';

const currentShift = ref<any>(null);
const shifts = ref<any[]>([]);
const selectedSummary = ref<any>(null);
const loadingAction = ref(false);

const showOpenModal = ref(false);
const showMovementModal = ref(false);
const showCloseModal = ref(false);

const openForm = ref({
  shift_name: 'Shift Pagi',
  opening_cash_float: 300000,
  notes: '',
});

const movementForm = ref({
  movement_type: 'cash_drop' as const,
  amount: 500000,
  reason: 'Setoran kas ke brankas supervisor',
});

const closeForm = ref({
  actual_cash_counted: 0,
  notes: '',
});

const fetchCurrentShift = async () => {
  try {
    const res = await getCurrentShift();
    currentShift.value = res.data?.data || null;
  } catch (err) {
    console.error('Error fetching current shift:', err);
    currentShift.value = null;
  }
};

const fetchShifts = async () => {
  try {
    const res = await getShifts({ per_page: 20 });
    shifts.value = res.data?.data || [];
  } catch (err) {
    console.error('Error fetching shifts history:', err);
  }
};

const submitOpenShift = async () => {
  loadingAction.value = true;
  try {
    await openShift(openForm.value);
    showOpenModal.value = false;
    openForm.value.opening_cash_float = 300000;
    openForm.value.notes = '';
    await fetchCurrentShift();
    await fetchShifts();
  } catch (err: any) {
    alert(err.response?.data?.message || 'Gagal membuka shift');
  } finally {
    loadingAction.value = false;
  }
};

const submitMovement = async () => {
  if (!currentShift.value) return;
  loadingAction.value = true;
  try {
    await recordCashMovement(currentShift.value.id, movementForm.value);
    showMovementModal.value = false;
    movementForm.value.amount = 500000;
    await fetchCurrentShift();
  } catch (err: any) {
    alert(err.response?.data?.message || 'Gagal mencatat mutasi kas');
  } finally {
    loadingAction.value = false;
  }
};

const submitCloseShift = async () => {
  if (!currentShift.value) return;
  loadingAction.value = true;
  try {
    const res = await closeShift(currentShift.value.id, closeForm.value);
    showCloseModal.value = false;
    closeForm.value.actual_cash_counted = 0;
    closeForm.value.notes = '';
    const closedId = currentShift.value.id;
    await fetchCurrentShift();
    await fetchShifts();
    await openSummaryModal(closedId);
  } catch (err: any) {
    alert(err.response?.data?.message || 'Gagal menutup shift');
  } finally {
    loadingAction.value = false;
  }
};

const openSummaryModal = async (shiftId: string) => {
  try {
    const res = await getShiftSummary(shiftId);
    selectedSummary.value = res.data?.data;
  } catch (err) {
    console.error('Error fetching summary:', err);
  }
};

const printSummary = () => {
  window.print();
};

const formatRupiah = (val: number | undefined | null) => {
  if (val === undefined || val === null) return 'Rp 0';
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(val);
};

const formatDateTime = (val: string | undefined | null) => {
  if (!val) return '-';
  const d = new Date(val);
  return d.toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

onMounted(() => {
  fetchCurrentShift();
  fetchShifts();
});
</script>
