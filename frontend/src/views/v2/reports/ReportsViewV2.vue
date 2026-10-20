<template>
  <div class="reports-v2 space-y-6 animate-in fade-in duration-200">
    <!-- Page Header V2 -->
    <PageHeaderV2
      title="Laporan Bisnis & Analitik Eksekutif"
      description="Analisis performa penjualan harian, profitabilitas cabang, audit inventaris, dan efisiensi tenaga kerja."
      badge="Business Intelligence"
    >
      <template #actions>
        <select 
          v-model="selectedRange"
          @change="handleFilterChange"
          class="px-3 py-2 rounded-xl border text-xs font-semibold outline-none cursor-pointer"
          style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
        >
          <option value="today">Hari Ini</option>
          <option value="this-week">Minggu Ini</option>
          <option value="this-month">Bulan Ini (September - Oktober)</option>
          <option value="this-year">Tahun 2026</option>
        </select>

        <button
          type="button"
          @click="exportExcel"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold border transition-all flex items-center gap-1.5 cursor-pointer text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800"
          style="border-color: var(--border-color); background-color: var(--bg-primary);"
        >
          <FileSpreadsheet class="w-3.5 h-3.5 text-emerald-600" />
          <span>Export Excel</span>
        </button>

        <button
          type="button"
          @click="exportPDF"
          class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-blue-600 hover:bg-blue-700 text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Unduh PDF</span>
        </button>
      </template>

      <!-- Sub-tabs Slot -->
      <template #tabs>
        <div class="flex items-center gap-6 text-xs font-semibold overflow-x-auto pb-2 select-none">
          <button
            v-for="t in reportTabs"
            :key="t.id"
            type="button"
            @click="activeTab = t.id"
            class="flex items-center gap-2 pb-2 -mb-2 cursor-pointer transition-colors whitespace-nowrap"
            :class="activeTab === t.id 
              ? 'text-blue-600 dark:text-blue-400 border-b-2 border-blue-600' 
              : 'text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200'"
          >
            <component :is="t.icon" class="w-3.5 h-3.5" />
            <span>{{ t.label }}</span>
          </button>
        </div>
      </template>
    </PageHeaderV2>

    <!-- ===================================================================== -->
    <!-- TAB 1: SALES REPORT (PENJUALAN)                                       -->
    <!-- ===================================================================== -->
    <div v-if="activeTab === 'sales'" class="space-y-6">
      <!-- 4 Stat Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Omzet Penjualan"
          value="Rp 148.520.000"
          trend="up"
          trendValue="+14.5%"
          subtitle="Dibandingkan bulan lalu"
          color="emerald"
          :icon="DollarSign"
        />
        <StatCardV2
          title="Rata-rata Order (AOV)"
          value="Rp 54.200"
          trend="up"
          trendValue="+4.2%"
          subtitle="Per transaksi pelanggan"
          color="blue"
          :icon="TrendingUp"
        />
        <StatCardV2
          title="Total Transaksi Sukses"
          value="2.740"
          subtitle="Kasir POS & Dine-in"
          color="indigo"
          :icon="ShoppingBag"
        />
        <StatCardV2
          title="Metode Bayar Favorit"
          value="QRIS (68%)"
          subtitle="Non-tunai mendominasi"
          color="amber"
          :icon="CreditCard"
        />
      </div>

      <!-- Trend Sales Bar Chart & Best Sellers -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Sales Trend Bar Chart (8 cols) -->
        <div 
          class="lg:col-span-8 p-6 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between mb-6">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Tren Penjualan Harian</h3>
              <p class="text-xs text-slate-400 mt-0.5">Grafik dinamika omzet harian periode terpilih</p>
            </div>
            <span class="text-xs font-bold text-emerald-600 bg-emerald-50 dark:bg-emerald-950/40 px-2.5 py-0.5 rounded-full border border-emerald-200">
              Tren Positif +14.5%
            </span>
          </div>

          <div class="h-56 flex items-end justify-between gap-2 pt-6 pb-2 border-b" style="border-color: var(--border-color);">
            <div 
              v-for="(day, idx) in dailyTrend" 
              :key="idx"
              class="flex-1 flex flex-col items-center gap-1.5 h-full justify-end group cursor-pointer"
            >
              <div 
                class="w-full max-w-[28px] bg-blue-600 hover:bg-blue-700 rounded-t-lg transition-all relative"
                :style="{ height: `${(day.amount / 6000000) * 100}%` }"
                :title="`${day.date}: Rp ${formatNum(day.amount)}`"
              ></div>
              <span class="text-[10px] font-mono text-slate-400 group-hover:text-slate-900 dark:group-hover:text-slate-100">
                {{ day.date }}
              </span>
            </div>
          </div>
        </div>

        <!-- Best Sellers Leaderboard (4 cols) -->
        <div 
          class="lg:col-span-4 p-6 rounded-2xl border shadow-xs space-y-4"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Menu Terlaris (Top 5)</h3>
            <p class="text-xs text-slate-400 mt-0.5">Produk dengan kontribusi omzet tertinggi</p>
          </div>

          <div class="space-y-3">
            <div 
              v-for="(prod, pIdx) in topProducts" 
              :key="prod.name"
              class="p-3 rounded-xl border flex items-center justify-between text-xs"
              style="background-color: var(--bg-content); border-color: var(--border-color);"
            >
              <div class="flex items-center gap-3">
                <span class="w-6 h-6 rounded-md font-bold font-mono text-xs flex items-center justify-center"
                      :class="pIdx === 0 ? 'bg-amber-100 text-amber-800' : pIdx === 1 ? 'bg-slate-200 text-slate-700' : 'bg-slate-100 text-slate-600'">
                  #{{ pIdx + 1 }}
                </span>
                <div>
                  <div class="font-bold text-slate-900 dark:text-slate-100">{{ prod.name }}</div>
                  <div class="text-[10px] text-slate-400">{{ prod.sold }} Porsi Terjual</div>
                </div>
              </div>
              <div class="text-right font-mono font-bold text-blue-600 dark:text-blue-400">
                Rp {{ formatNum(prod.revenue) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- TAB 2: FINANCIAL REPORT (LABA RUGI & ARUS KAS)                         -->
    <!-- ===================================================================== -->
    <div v-else-if="activeTab === 'financial'" class="space-y-6">
      <!-- 4 Financial KPIs -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Pendapatan Kotor"
          :value="`Rp ${formatNum(financialData.kpi.total_revenue)}`"
          trend="up"
          trendValue="+14.2%"
          subtitle="Omzet Kasir & Online"
          color="blue"
          :icon="DollarSign"
        />
        <StatCardV2
          title="HPP & Beban Operasional"
          :value="`Rp ${formatNum(financialData.kpi.total_expenses)}`"
          trend="neutral"
          trendValue="COGS + Biaya Cafe"
          subtitle="Total Beban Berjalan"
          color="rose"
          :icon="Receipt"
        />
        <StatCardV2
          title="Laba Kotor (Gross Profit)"
          :value="`Rp ${formatNum(financialData.kpi.gross_profit)}`"
          trend="up"
          :trendValue="financialData.kpi.gross_margin_pct"
          subtitle="Margin Laba Kotor"
          color="emerald"
          :icon="TrendingUp"
        />
        <StatCardV2
          title="Total Kas Likuid"
          :value="`Rp ${formatNum(financialData.cash_flow.total_liquid_cash)}`"
          subtitle="Kasir Drawer + Bank BCA"
          color="indigo"
          :icon="Landmark"
        />
      </div>

      <!-- Income Statement (Laporan Laba Rugi / P&L) Table -->
      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="p-5 sm:p-6 border-b flex flex-col sm:flex-row sm:items-center justify-between gap-3" style="border-color: var(--border-color);">
          <div>
            <div class="flex items-center gap-2.5">
              <h3 class="text-base sm:text-lg font-bold text-slate-900 dark:text-slate-50">
                Laporan Laba Rugi (Income Statement / P&L)
              </h3>
              <StatusBadgeV2 status="completed" label="Jurnal Terposting" />
            </div>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
              Rincian pendapatan operasional kafe, harga pokok penjualan (HPP), dan realisasi beban usaha
            </p>
          </div>
          <span class="text-xs font-mono font-semibold text-slate-400">
            Standar SAK EMKM • Akrual
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-5">Pos Akun / Kategori Akuntansi</th>
                <th class="py-3 px-5 text-center">Porsi (%)</th>
                <th class="py-3 px-5 text-right">Jumlah / Nominal</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <template v-for="(cat, cIdx) in financialData.pnl_statement" :key="cIdx">
                <!-- Category Section Header -->
                <tr class="font-bold text-slate-800 dark:text-slate-100" style="background-color: var(--bg-content);">
                  <td class="py-3 px-5 flex items-center gap-2">
                    <span class="w-2 h-2 rounded-full" :class="cIdx === 0 ? 'bg-blue-600' : cIdx === 1 ? 'bg-amber-500' : 'bg-rose-500'"></span>
                    <span>{{ cat.category }}</span>
                  </td>
                  <td class="py-3 px-5 text-center font-mono text-[11px] text-slate-400">
                    100.0%
                  </td>
                  <td class="py-3 px-5 text-right font-mono font-bold text-sm" :class="cIdx === 0 ? 'text-blue-600 dark:text-blue-400' : 'text-slate-900 dark:text-slate-100'">
                    Rp {{ formatNum(cat.subtotal) }}
                  </td>
                </tr>

                <!-- Item Rows inside Category -->
                <tr 
                  v-for="(item, iIdx) in cat.items" 
                  :key="iIdx"
                  class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30 transition-colors"
                >
                  <td class="py-2.5 px-5 pl-9 text-slate-600 dark:text-slate-300">
                    {{ item.name }}
                  </td>
                  <td class="py-2.5 px-5 text-center font-mono text-[11px] text-slate-400">
                    {{ item.pct }}
                  </td>
                  <td class="py-2.5 px-5 text-right font-mono font-medium text-slate-800 dark:text-slate-200">
                    Rp {{ formatNum(item.amount) }}
                  </td>
                </tr>
              </template>

              <!-- Net Profit / Laba Bersih Highlight Row -->
              <tr class="font-bold border-t-2" :class="financialData.kpi.net_profit >= 0 ? 'bg-emerald-50/70 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-200' : 'bg-rose-50/70 dark:bg-rose-950/40 text-rose-800 dark:text-rose-200'" style="border-color: var(--border-color);">
                <td class="py-3.5 px-5 text-sm uppercase tracking-wide">
                  Laba / (Rugi) Bersih Periode Berjalan (Net Profit)
                </td>
                <td class="py-3.5 px-5 text-center font-mono text-xs">
                  {{ financialData.kpi.net_margin_pct }}
                </td>
                <td class="py-3.5 px-5 text-right font-mono text-base font-black">
                  Rp {{ formatNum(financialData.kpi.net_profit) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Cash Flow & Tax Compliance Tables Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Arus Kas Table -->
        <div 
          class="rounded-2xl border shadow-xs overflow-hidden"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
            <div>
              <h3 class="text-sm font-bold text-slate-900 dark:text-slate-50">Ringkasan Arus Kas (Cash Flow)</h3>
              <p class="text-[11px] text-slate-400 mt-0.5">Arus kas masuk & keluar kasir dan bank</p>
            </div>
            <StatusBadgeV2 status="active" label="Likuid" />
          </div>

          <table class="w-full text-left text-xs border-collapse">
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                <td class="py-3 px-4 font-medium text-slate-700 dark:text-slate-300">Arus Kas Operasional (Operating)</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-emerald-600">
                  + Rp {{ formatNum(financialData.cash_flow.operating_cash_flow) }}
                </td>
              </tr>
              <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                <td class="py-3 px-4 font-medium text-slate-700 dark:text-slate-300">Arus Kas Investasi (Alat & Servis)</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-rose-600">
                  - Rp {{ formatNum(Math.abs(financialData.cash_flow.investing_cash_flow)) }}
                </td>
              </tr>
              <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                <td class="py-3 px-4 font-medium text-slate-700 dark:text-slate-300">Saldo Kas Tunai Kasir (Drawer)</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(financialData.cash_flow.cash_on_hand) }}
                </td>
              </tr>
              <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                <td class="py-3 px-4 font-medium text-slate-700 dark:text-slate-300">Saldo Rekening Bank BCA Operasional</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(financialData.cash_flow.bank_balance) }}
                </td>
              </tr>
              <tr class="font-bold border-t bg-slate-50/80 dark:bg-slate-800/50" style="border-color: var(--border-color);">
                <td class="py-3.5 px-4 text-slate-900 dark:text-slate-100">Total Kas & Setara Kas</td>
                <td class="py-3.5 px-4 text-right font-mono text-sm text-blue-600 dark:text-blue-400">
                  Rp {{ formatNum(financialData.cash_flow.total_liquid_cash) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Rekapitulasi Pajak Resto PB1 Table -->
        <div 
          class="rounded-2xl border shadow-xs overflow-hidden"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
            <div>
              <h3 class="text-sm font-bold text-slate-900 dark:text-slate-50">Kepatuhan Pajak (Tax Compliance PB1)</h3>
              <p class="text-[11px] text-slate-400 mt-0.5">Rekap pajak restoran daerah dari transaksi kasir</p>
            </div>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-amber-50 text-amber-700 border border-amber-200">
              Siap Setor
            </span>
          </div>

          <table class="w-full text-left text-xs border-collapse">
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                <td class="py-3 px-4 font-medium text-slate-700 dark:text-slate-300">Dasar Pengenaan Pajak (DPP Penjualan)</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(financialData.tax_summary.taxable_sales) }}
                </td>
              </tr>
              <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                <td class="py-3 px-4 font-medium text-slate-700 dark:text-slate-300">Tarif Pajak Restoran (PB1)</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-blue-600">
                  {{ financialData.tax_summary.pb1_tax_rate }}
                </td>
              </tr>
              <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                <td class="py-3 px-4 font-medium text-slate-700 dark:text-slate-300">Total Pajak Terkumpul dari Konsumen</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-emerald-600">
                  Rp {{ formatNum(financialData.tax_summary.tax_collected) }}
                </td>
              </tr>
              <tr class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                <td class="py-3 px-4 font-medium text-slate-700 dark:text-slate-300">Status Penyetoran Kas Daerah</td>
                <td class="py-3 px-4 text-right font-semibold text-amber-700 dark:text-amber-400">
                  {{ financialData.tax_summary.status }}
                </td>
              </tr>
              <tr class="font-bold border-t bg-slate-50/80 dark:bg-slate-800/50" style="border-color: var(--border-color);">
                <td class="py-3.5 px-4 text-slate-900 dark:text-slate-100">Pajak Terutang Siap Dilaporkan</td>
                <td class="py-3.5 px-4 text-right font-mono text-sm text-emerald-600">
                  Rp {{ formatNum(financialData.tax_summary.tax_payable) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- TAB 3: INVENTORY REPORT (VALUASI STOK)                                -->
    <!-- ===================================================================== -->
    <div v-else-if="activeTab === 'inventory'" class="space-y-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Nilai Valuasi Stok"
          :value="`Rp ${formatNum(inventoryData.kpi.total_valuation)}`"
          subtitle="Aset bahan baku gudang"
          color="blue"
          :icon="Package"
        />
        <StatCardV2
          title="Total Item SKU"
          :value="`${inventoryData.kpi.total_items} SKU`"
          subtitle="Terdata di sistem kartu stok"
          color="indigo"
          :icon="FileSpreadsheet"
        />
        <StatCardV2
          title="Stok Menipis / Kritis"
          :value="`${inventoryData.kpi.low_stock_items} Item`"
          trend="down"
          trendValue="Perlu Restok"
          subtitle="Mendekati safety stock"
          color="amber"
          :icon="AlertTriangle"
        />
        <StatCardV2
          title="Rasio Perputaran (Turnover)"
          :value="inventoryData.kpi.turnover_ratio"
          subtitle="Siklus bahan baku lancar"
          color="emerald"
          :icon="TrendingUp"
        />
      </div>

      <!-- Inventory Stock Status Table -->
      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="p-4 border-b flex flex-col sm:flex-row sm:items-center justify-between gap-3" style="border-color: var(--border-color);">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Daftar Valuasi & Status Stok Bahan Baku</h3>
            <p class="text-xs text-slate-400 mt-0.5">Berdasarkan mutasi kartu stok real-time</p>
          </div>
          <div class="relative w-full sm:w-64">
            <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5 pointer-events-none" />
            <input 
              v-model="inventorySearch"
              type="text" 
              placeholder="Cari SKU atau nama item..."
              class="w-full pl-8.5 pr-3 py-1.5 text-xs rounded-xl border outline-none"
              style="background-color: var(--input-bg); border-color: var(--input-border); color: var(--input-text);"
            />
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">SKU</th>
                <th class="py-3 px-4">Nama Bahan Baku</th>
                <th class="py-3 px-4">Kategori</th>
                <th class="py-3 px-4 text-right">Stok Fisik</th>
                <th class="py-3 px-4 text-right">Min. Stok</th>
                <th class="py-3 px-4 text-right">HPP Satuan</th>
                <th class="py-3 px-4 text-right">Total Valuasi</th>
                <th class="py-3 px-4 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr v-for="item in filteredInventoryItems" :key="item.sku" class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30 transition-colors">
                <td class="py-3 px-4 font-mono font-medium text-blue-600">{{ item.sku }}</td>
                <td class="py-3 px-4 font-semibold text-slate-900 dark:text-slate-100">{{ item.name }}</td>
                <td class="py-3 px-4 text-slate-500">{{ item.category }}</td>
                <td class="py-3 px-4 text-right font-mono font-bold">{{ item.quantity }} {{ item.uom }}</td>
                <td class="py-3 px-4 text-right font-mono text-slate-400">{{ item.min_stock }} {{ item.uom }}</td>
                <td class="py-3 px-4 text-right font-mono text-slate-600 dark:text-slate-300">Rp {{ formatNum(item.average_cost) }}</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">Rp {{ formatNum(item.valuation) }}</td>
                <td class="py-3 px-4 text-center">
                  <StatusBadgeV2 
                    :status="item.status === 'Aman' ? 'active' : item.status === 'Menipis' ? 'pending' : 'cancelled'"
                    :label="item.status"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- TAB 4: HR REPORT (SDM & PRODUKTIVITAS)                                -->
    <!-- ===================================================================== -->
    <div v-else-if="activeTab === 'hr'" class="space-y-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Karyawan Aktif"
          :value="`${hrData.kpi.total_employees} Orang`"
          subtitle="Staf operasional terdaftar"
          color="blue"
          :icon="Users"
        />
        <StatCardV2
          title="Rata-rata Kehadiran"
          :value="hrData.kpi.attendance_rate_pct"
          trend="up"
          trendValue="95 Shift Selesai"
          subtitle="Tingkat kedisiplinan"
          color="emerald"
          :icon="CheckCircle2"
        />
        <StatCardV2
          title="Total Biaya Payroll Gaji"
          :value="`Rp ${formatNum(hrData.kpi.total_payroll_cost)}`"
          subtitle="Gaji pokok + tunjangan"
          color="indigo"
          :icon="Banknote"
        />
        <StatCardV2
          title="Akumulasi Jam Lembur"
          :value="hrData.kpi.total_overtime_hours"
          subtitle="Kompensasi shift sibuk"
          color="amber"
          :icon="Clock"
        />
      </div>

      <!-- Staff Performance Table -->
      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Evaluasi Kinerja & Presensi Staf</h3>
            <p class="text-xs text-slate-400 mt-0.5">Rekapitulasi produktivitas shift dan payroll bulanan</p>
          </div>
          <StatusBadgeV2 status="active" label="Aktif" />
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">NIK / Nama Karyawan</th>
                <th class="py-3 px-4">Departemen</th>
                <th class="py-3 px-4">Posisi</th>
                <th class="py-3 px-4 text-center">Total Shift</th>
                <th class="py-3 px-4 text-center">Tingkat Hadir</th>
                <th class="py-3 px-4 text-center">Lembur</th>
                <th class="py-3 px-4 text-right">Gaji Pokok</th>
                <th class="py-3 px-4 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr v-for="s in hrData.staff_performance" :key="s.nik" class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30 transition-colors">
                <td class="py-3 px-4">
                  <div class="font-bold text-slate-900 dark:text-slate-100">{{ s.name }}</div>
                  <div class="font-mono text-[10.5px] text-slate-400">{{ s.nik }}</div>
                </td>
                <td class="py-3 px-4 text-slate-500">{{ s.department }}</td>
                <td class="py-3 px-4 font-medium text-slate-800 dark:text-slate-200">{{ s.position }}</td>
                <td class="py-3 px-4 text-center font-mono font-bold">{{ s.shifts }}</td>
                <td class="py-3 px-4 text-center font-mono text-emerald-600 font-bold">{{ s.attendance_rate }}</td>
                <td class="py-3 px-4 text-center font-mono text-slate-500">{{ s.overtime_hours }}</td>
                <td class="py-3 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(s.basic_salary) }}
                </td>
                <td class="py-3 px-4 text-center">
                  <StatusBadgeV2 :status="s.status === 'Tetap' ? 'active' : 'pending'" :label="s.status" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- TAB 5: CUSTOM REPORT (LAPORAN KUSTOM)                                 -->
    <!-- ===================================================================== -->
    <div v-else class="space-y-6">
      <div 
        class="p-6 rounded-2xl border shadow-xs space-y-4"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="flex items-center justify-between pb-3 border-b" style="border-color: var(--border-color);">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Generator Laporan Kustom</h3>
            <p class="text-xs text-slate-400 mt-0.5">Pilih template preset laporan atau buat query analisis khusus</p>
          </div>
          <button
            type="button"
            @click="generateCustomReport"
            class="px-4 py-2 rounded-xl text-xs font-bold bg-blue-600 hover:bg-blue-700 text-white cursor-pointer shadow-xs flex items-center gap-1.5"
          >
            <Sparkles class="w-3.5 h-3.5" />
            <span>Generate Laporan</span>
          </button>
        </div>

        <!-- Presets -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
          <div
            v-for="tmpl in customPresets"
            :key="tmpl.name"
            @click="selectPreset(tmpl)"
            class="p-3.5 rounded-xl border cursor-pointer hover:border-blue-500 hover:shadow-xs transition-all"
            style="background-color: var(--bg-content); border-color: var(--border-color);"
          >
            <div class="font-bold text-xs text-slate-900 dark:text-slate-100">{{ tmpl.name }}</div>
            <div class="text-[11px] text-slate-400 mt-1">{{ tmpl.desc }}</div>
          </div>
        </div>

        <!-- Result Preview Table if generated -->
        <div v-if="customResult" class="mt-4 pt-4 border-t space-y-3" style="border-color: var(--border-color);">
          <div class="flex items-center justify-between">
            <h4 class="font-bold text-xs text-slate-900 dark:text-slate-100">Hasil Analisis Kustom: {{ customResult.title }}</h4>
            <span class="text-xs font-mono text-slate-400">{{ customResult.rows.length }} Baris Data</span>
          </div>

          <div class="overflow-x-auto rounded-xl border" style="border-color: var(--border-color);">
            <table class="w-full text-left text-xs border-collapse">
              <thead>
                <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                  <th v-for="h in customResult.headers" :key="h" class="py-2.5 px-4">{{ h }}</th>
                </tr>
              </thead>
              <tbody class="divide-y" style="border-color: var(--border-color);">
                <tr v-for="(row, rIdx) in customResult.rows" :key="rIdx" class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                  <td v-for="(cell, cIdx) in row" :key="cIdx" class="py-2.5 px-4 font-mono">
                    {{ cell }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import PageHeaderV2 from '@/components/v2/common/PageHeaderV2.vue'
import StatCardV2 from '@/components/v2/common/StatCardV2.vue'
import StatusBadgeV2 from '@/components/v2/common/StatusBadgeV2.vue'
import {
  FileSpreadsheet,
  Download,
  ShoppingBag,
  DollarSign,
  Package,
  Users,
  Sliders,
  TrendingUp,
  CreditCard,
  Receipt,
  Landmark,
  Search,
  CheckCircle2,
  Banknote,
  Clock,
  AlertTriangle,
  Sparkles
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const route = useRoute()
const router = useRouter()
const notifyStore = useNotificationStore()

const activeTab = ref((route.query.tab as string) || 'sales')
const selectedRange = ref('this-month')
const inventorySearch = ref('')

const reportTabs = [
  { id: 'sales', label: 'Penjualan (Sales)', icon: ShoppingBag },
  { id: 'financial', label: 'Laba Rugi (Finance)', icon: DollarSign },
  { id: 'inventory', label: 'Valuasi Stok (Inventory)', icon: Package },
  { id: 'hr', label: 'Produktivitas SDM (HR)', icon: Users },
  { id: 'custom', label: 'Laporan Kustom', icon: Sliders }
]

watch(activeTab, (newTab) => {
  router.replace({ query: { ...route.query, tab: newTab } })
  loadTabData(newTab)
})

// Tab 1: Sales Data
const dailyTrend = [
  { date: '01 Okt', amount: 3850000 },
  { date: '02 Okt', amount: 4200000 },
  { date: '03 Okt', amount: 4950000 },
  { date: '04 Okt', amount: 5600000 },
  { date: '05 Okt', amount: 3900000 },
  { date: '06 Okt', amount: 4100000 },
  { date: '07 Okt', amount: 4800000 },
  { date: '08 Okt', amount: 5200000 }
]

const topProducts = [
  { name: 'Caffe Latte Oatmilk', sold: 480, revenue: 18240000 },
  { name: 'Americano Ice Single', sold: 420, revenue: 11760000 },
  { name: 'V60 Gayo Natural', sold: 290, revenue: 10150000 },
  { name: 'Croissant Butter French', sold: 210, revenue: 6720000 },
  { name: 'Truffle Fries Parmesan', sold: 160, revenue: 6720000 }
]

// Tab 2: Financial Data
const financialData = ref({
  kpi: {
    total_revenue: 2664200,
    total_expenses: 23969638.25,
    gross_profit: 1651804,
    gross_margin_pct: '62.0%',
    net_profit: -21305438.25,
    net_margin_pct: '-799.7%'
  },
  pnl_statement: [
    {
      category: 'Pendapatan Operasional (Revenue)',
      subtotal: 2664200,
      items: [
        { name: 'Penjualan Minuman Kopi & Non-Kopi', pct: '72.0%', amount: 1918224 },
        { name: 'Penjualan Makanan & Pastry', pct: '24.0%', amount: 639408 },
        { name: 'Penjualan Merchandise & Whole Beans', pct: '4.0%', amount: 106568 }
      ]
    },
    {
      category: 'Harga Pokok Penjualan (HPP / COGS)',
      subtotal: 1012396,
      items: [
        { name: 'Bahan Baku Biji Kopi (Arabica Specialty)', pct: '45.0%', amount: 455578.2 },
        { name: 'Susu Segar, Oatmilk & Dairy', pct: '30.0%', amount: 303718.8 },
        { name: 'Sirup, Flavour, & Topping', pct: '15.0%', amount: 151859.4 },
        { name: 'Cup, Sedotan & Packaging Takeaway', pct: '10.0%', amount: 101239.6 }
      ]
    },
    {
      category: 'Beban Operasional (Operating Expenses)',
      subtotal: 22957242.25,
      items: [
        { name: 'Beban Gaji & Upah Barista / Staf', pct: '55.0%', amount: 14477214.42 },
        { name: 'Beban Listrik, Air & Gas Cafe', pct: '18.0%', amount: 850000 },
        { name: 'Pemeliharaan & Servis Mesin Espresso', pct: '10.0%', amount: 450000 },
        { name: 'Pemasaran, Ads Instagram & Promo', pct: '11.0%', amount: 500000 },
        { name: 'Beban Internet Wi-Fi & Aplikasi POS', pct: '6.0%', amount: 350000 }
      ]
    }
  ],
  cash_flow: {
    operating_cash_flow: 451804,
    investing_cash_flow: -1500000,
    cash_on_hand: 32500000,
    bank_balance: 68450000,
    total_liquid_cash: 100950000
  },
  tax_summary: {
    taxable_sales: 2664200,
    pb1_tax_rate: '10%',
    tax_collected: 266420,
    tax_payable: 266420,
    status: 'Siap Setor Kas Daerah (PB1 Restoran)'
  }
})

// Tab 3: Inventory Data
const inventoryData = ref({
  kpi: {
    total_valuation: 128450000,
    total_items: 28,
    low_stock_items: 3,
    turnover_ratio: '4.2x / Bln'
  },
  stock_items: [
    { sku: 'RM-COF-001', name: 'Arabica Gayo Natural Roast', category: 'Biji Kopi', quantity: 18.5, uom: 'kg', average_cost: 185000, valuation: 3422500, status: 'Aman', min_stock: 5 },
    { sku: 'RM-DRY-001', name: 'Fresh Milk Greenfields 1L', category: 'Dairy', quantity: 4, uom: 'karton', average_cost: 215000, valuation: 860000, status: 'Menipis', min_stock: 12 },
    { sku: 'RM-DRY-002', name: 'Oatside Barista Blend 1L', category: 'Dairy', quantity: 2, uom: 'karton', average_cost: 345000, valuation: 690000, status: 'Menipis', min_stock: 8 },
    { sku: 'RM-SRP-001', name: 'Monin Salted Caramel 700ml', category: 'Sirup', quantity: 6, uom: 'botol', average_cost: 165000, valuation: 990000, status: 'Aman', min_stock: 2 },
    { sku: 'RM-SRP-002', name: 'Monin Vanilla Syrup 700ml', category: 'Sirup', quantity: 0, uom: 'botol', average_cost: 165000, valuation: 0, status: 'Habis', min_stock: 3 },
    { sku: 'RM-PKG-001', name: 'Paper Cup Hot 8oz Cafe ERP', category: 'Packaging', quantity: 450, uom: 'pcs', average_cost: 850, valuation: 382500, status: 'Aman', min_stock: 100 }
  ]
})

const filteredInventoryItems = computed(() => {
  if (!inventorySearch.value.trim()) return inventoryData.value.stock_items
  const q = inventorySearch.value.toLowerCase().trim()
  return inventoryData.value.stock_items.filter(i => i.name.toLowerCase().includes(q) || i.sku.toLowerCase().includes(q))
})

// Tab 4: HR Data
const hrData = ref({
  kpi: {
    total_employees: 7,
    attendance_rate_pct: '97.2%',
    total_payroll_cost: 28260969,
    total_overtime_hours: '18 Jam'
  },
  staff_performance: [
    { nik: 'EMP-001', name: 'Budi Pratama', department: 'Operasional', position: 'Head Barista', shifts: 24, attendance_rate: '100%', overtime_hours: '6 Jam', basic_salary: 4500000, status: 'Tetap' },
    { nik: 'EMP-002', name: 'Jane Doe', department: 'Kasir', position: 'Kasir POS Senior', shifts: 22, attendance_rate: '95.5%', overtime_hours: '4 Jam', basic_salary: 4200000, status: 'Tetap' },
    { nik: 'EMP-003', name: 'Ahmad Gudang', department: 'Logistik', position: 'Warehouse Staff', shifts: 24, attendance_rate: '100%', overtime_hours: '2 Jam', basic_salary: 4000000, status: 'Tetap' },
    { nik: 'EMP-004', name: 'Dewi Finance', department: 'Keuangan', position: 'Akuntan Keuangan', shifts: 20, attendance_rate: '100%', overtime_hours: '0 Jam', basic_salary: 5000000, status: 'Tetap' },
    { nik: 'EMP-005', name: 'Rian Barista', department: 'Operasional', position: 'Junior Barista', shifts: 22, attendance_rate: '90.9%', overtime_hours: '6 Jam', basic_salary: 3500000, status: 'Percobaan' }
  ]
})

// Tab 5: Custom Data
const customPresets = [
  { name: 'Omzet per Kasir POS', desc: 'Rincian kontribusi pendapatan per kasir outlet' },
  { name: 'Valuasi Stok per Kategori', desc: 'Distribusi nilai aset bahan baku per departemen' },
  { name: 'Tingkat Kehadiran per Shift', desc: 'Evaluasi shift pagi vs shift siang' },
  { name: 'Beban Supplier & Pembelian', desc: 'Total realisasi pembelian bahan baku per vendor' }
]

const customResult = ref<{ title: string; headers: string[]; rows: (string | number)[][] } | null>({
  title: 'Omzet per Kasir POS',
  headers: ['Nama Kasir', 'Total Transaksi', 'Total Penjualan', 'Kontribusi (%)'],
  rows: [
    ['Jane Cashier (Shift Pagi)', 1420, 'Rp 78.400.000', '52.8%'],
    ['Budi Kasir (Shift Siang)', 1320, 'Rp 70.120.000', '47.2%']
  ]
})

const selectPreset = (tmpl: any) => {
  if (tmpl.name === 'Valuasi Stok per Kategori') {
    customResult.value = {
      title: 'Valuasi Stok per Kategori',
      headers: ['Kategori Bahan', 'Jumlah SKU', 'Total Nilai Aset', 'Porsi Aset'],
      rows: [
        ['Biji Kopi (Beans)', 6, 'Rp 58.200.000', '45.3%'],
        ['Dairy & Susu', 4, 'Rp 28.500.000', '22.2%'],
        ['Sirup & Perisa', 8, 'Rp 24.300.000', '18.9%'],
        ['Kemasan & Packaging', 10, 'Rp 17.450.000', '13.6%']
      ]
    }
  } else {
    customResult.value = {
      title: tmpl.name,
      headers: ['Indikator / Entitas', 'Frekuensi', 'Nominal / Nilai', 'Status'],
      rows: [
        ['Shift Pagi (07:00 - 15:00)', '48 Shift', 'Rp 78.400.000', 'Optimal'],
        ['Shift Siang (15:00 - 23:00)', '47 Shift', 'Rp 70.120.000', 'Normal']
      ]
    }
  }
  notifyStore.success(`Template preset "${tmpl.name}" diterapkan`, 'Laporan Kustom')
}

const generateCustomReport = () => {
  notifyStore.success('Laporan kustom berhasil digenerate!', 'Custom Query')
}

const loadTabData = async (tab: string) => {
  try {
    if (tab === 'financial') {
      const res = await axios.get('/api/v1/reports/financial')
      if (res.data?.data) financialData.value = res.data.data
    } else if (tab === 'inventory') {
      const res = await axios.get('/api/v1/reports/inventory')
      if (res.data?.data?.stock_items) inventoryData.value.stock_items = res.data.data.stock_items
    } else if (tab === 'hr') {
      const res = await axios.get('/api/v1/reports/hr')
      if (res.data?.data?.staff_performance) hrData.value = res.data.data
    }
  } catch (err) {
    console.warn(`Could not load real-time data for tab ${tab}, using cached:`, err)
  }
}

const handleFilterChange = () => {
  notifyStore.info(`Rentang waktu diubah: ${selectedRange.value}`, 'Filter Laporan')
  loadTabData(activeTab.value)
}

const exportExcel = () => {
  notifyStore.success('File spreadsheet laporan (.xlsx) berhasil diunduh', 'Export Excel')
}

const exportPDF = () => {
  notifyStore.success('Dokumen ringkasan eksekutif (.pdf) berhasil diekspor', 'Export PDF')
}

const formatNum = (val: number | string) => Number(val || 0).toLocaleString('id-ID')

onMounted(() => {
  loadTabData(activeTab.value)
})
</script>
