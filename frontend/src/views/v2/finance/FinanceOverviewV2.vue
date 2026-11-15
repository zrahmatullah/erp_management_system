<template>
  <div class="finance-v2 space-y-6 animate-in fade-in duration-200">
    <!-- Page Header V2 -->
    <PageHeaderV2
      title="Manajemen Keuangan & Akuntansi"
      description="Pusat pencatatan pembukuan ganda (double-entry), pemantauan arus kas, rekonsiliasi perbankan, dan kepatuhan pajak."
      badge="General Ledger"
    >
      <template #actions>
        <button
          type="button"
          @click="refreshCurrentTab"
          class="h-9 w-9 rounded-xl border flex items-center justify-center transition-all duration-150 cursor-pointer text-slate-600 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 bg-white hover:bg-slate-50 dark:bg-slate-800 dark:hover:bg-slate-700 border-slate-200/80 dark:border-slate-700/80 shadow-xs hover:border-blue-300 dark:hover:border-blue-600"
          title="Perbarui Data"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>

        <button
          v-if="activeTab === 'journal' || activeTab === 'overview'"
          type="button"
          @click="showNewJournalModal = true"
          class="h-9 px-3.5 rounded-xl text-xs font-bold bg-[#4361EE] hover:bg-[#3752DE] text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer shadow-blue-500/20 hover:shadow-md"
        >
          <Plus class="w-4 h-4" />
          <span>Buat Jurnal Baru</span>
        </button>

        <button
          v-else-if="activeTab === 'coa'"
          type="button"
          @click="showNewAccountModal = true"
          class="h-9 px-3.5 rounded-xl text-xs font-bold bg-[#4361EE] hover:bg-[#3752DE] text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer shadow-blue-500/20 hover:shadow-md"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah Akun COA</span>
        </button>

        <button
          v-else-if="activeTab === 'expenses'"
          type="button"
          @click="showNewExpenseModal = true"
          class="h-9 px-3.5 rounded-xl text-xs font-bold bg-[#4361EE] hover:bg-[#3752DE] text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer shadow-blue-500/20 hover:shadow-md"
        >
          <Plus class="w-4 h-4" />
          <span>Catat Pengeluaran</span>
        </button>

        <button
          v-else-if="activeTab === 'budget'"
          type="button"
          @click="showNewBudgetModal = true"
          class="h-9 px-3.5 rounded-xl text-xs font-bold bg-[#4361EE] hover:bg-[#3752DE] text-white transition-all shadow-xs flex items-center gap-1.5 cursor-pointer shadow-blue-500/20 hover:shadow-md"
        >
          <Plus class="w-4 h-4" />
          <span>Atur Anggaran</span>
        </button>
      </template>

      <!-- Finance Navigation Sub-tabs -->
      <template #tabs>
        <div class="flex items-center gap-5 text-xs font-semibold overflow-x-auto pb-2 select-none">
          <button
            v-for="t in financeTabs"
            :key="t.id"
            type="button"
            @click="switchTab(t.id)"
            class="pb-2 -mb-2 cursor-pointer transition-colors whitespace-nowrap"
            :class="activeTab === t.id 
              ? 'text-blue-600 dark:text-blue-400 border-b-2 border-blue-600' 
              : 'text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200'"
          >
            {{ t.label }}
          </button>
        </div>
      </template>
    </PageHeaderV2>

    <!-- ======================================================================= -->
    <!-- TAB 1: OVERVIEW & CASHFLOW                                              -->
    <!-- ======================================================================= -->
    <div v-if="activeTab === 'overview'" class="space-y-6">
      <!-- 4 Financial Stat Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <StatCardV2
          title="Total Pendapatan"
          :value="`Rp ${formatNum(overviewData.total_revenue || 128450000)}`"
          trend="up"
          trendValue="+14.2%"
          subtitle="Akumulasi bulan berjalan"
          color="emerald"
          :icon="DollarSign"
        />
        <StatCardV2
          title="Beban Operasional"
          :value="`Rp ${formatNum(overviewData.total_expenses || 78320000)}`"
          trend="neutral"
          trendValue="Sesuai Budget"
          subtitle="HPP & Pengeluaran Kas"
          color="rose"
          :icon="Receipt"
        />
        <StatCardV2
          title="Laba Bersih Operasional"
          :value="`Rp ${formatNum(overviewData.net_profit || 50130000)}`"
          trend="up"
          trendValue="+18.5%"
          subtitle="Net Profit Margin 39.0%"
          color="blue"
          :icon="TrendingUp"
        />
        <StatCardV2
          title="Saldo Kas & Bank"
          :value="`Rp ${formatNum(overviewData.cash_on_hand || 32500000)}`"
          subtitle="Cash & Bank BCA"
          color="indigo"
          :icon="Landmark"
        />
      </div>

      <!-- Cashflow Bar Chart & Bank Accounts -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Monthly Revenue vs Expense Visualization (8 cols) -->
        <div 
          class="lg:col-span-8 p-6 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="flex items-center justify-between mb-6">
            <div>
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Tren Pendapatan vs Beban Bulanan</h3>
              <p class="text-xs text-slate-400 mt-0.5">Perbandingan performa finansial per bulan (dalam juta Rupiah)</p>
            </div>
            <div class="flex items-center gap-4 text-xs font-semibold">
              <span class="flex items-center gap-1.5 text-blue-600">
                <span class="w-3 h-3 rounded bg-blue-600"></span> Pendapatan
              </span>
              <span class="flex items-center gap-1.5 text-rose-500">
                <span class="w-3 h-3 rounded bg-rose-500"></span> Beban
              </span>
            </div>
          </div>

          <!-- Bar Chart Columns -->
          <div class="h-56 flex items-end justify-between gap-2 pt-6 pb-2 border-b" style="border-color: var(--border-color);">
            <div 
              v-for="bar in monthlyData" 
              :key="bar.m"
              class="flex-1 flex flex-col items-center gap-1.5 h-full justify-end group cursor-pointer"
            >
              <div class="w-full flex items-end justify-center gap-1 h-44">
                <div 
                  class="w-3.5 sm:w-4 bg-blue-600 hover:bg-blue-700 rounded-t-md transition-all relative"
                  :style="{ height: `${(bar.rev / 80) * 100}%` }"
                  :title="`${bar.m} Pendapatan: Rp ${bar.rev} Juta`"
                ></div>
                <div 
                  class="w-3.5 sm:w-4 bg-rose-500 hover:bg-rose-600 rounded-t-md transition-all relative"
                  :style="{ height: `${(bar.exp / 80) * 100}%` }"
                  :title="`${bar.m} Beban: Rp ${bar.exp} Juta`"
                ></div>
              </div>
              <span class="text-[10px] font-bold text-slate-400 group-hover:text-slate-900 dark:group-hover:text-slate-100 transition-colors">
                {{ bar.m }}
              </span>
            </div>
          </div>
        </div>

        <!-- Akun Kas & Rekening (4 cols) -->
        <div 
          class="lg:col-span-4 p-6 rounded-2xl border shadow-xs space-y-4 flex flex-col justify-between"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div>
            <div class="flex items-center justify-between mb-1">
              <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Likuiditas Kas & Bank</h3>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">
                Sinkron
              </span>
            </div>
            <p class="text-xs text-slate-400 mb-4">Saldo rekening operasional resmi</p>

            <div class="space-y-3">
              <div class="p-3.5 rounded-xl border flex items-center justify-between"
                   style="background-color: var(--bg-content); border-color: var(--border-color);">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-lg bg-blue-100 text-blue-700 font-bold text-xs flex items-center justify-center">
                    BCA
                  </div>
                  <div>
                    <div class="text-xs font-bold text-slate-900 dark:text-slate-100">BCA Giro Operasional</div>
                    <div class="text-[10.5px] font-mono text-slate-400">882-0192-331</div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-xs font-bold font-mono text-slate-900 dark:text-slate-100">Rp 24.850.000</div>
                  <div class="text-[10px] text-emerald-600 font-semibold">Aktif</div>
                </div>
              </div>

              <div class="p-3.5 rounded-xl border flex items-center justify-between"
                   style="background-color: var(--bg-content); border-color: var(--border-color);">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-lg bg-amber-100 text-amber-700 font-bold text-xs flex items-center justify-center">
                    CASH
                  </div>
                  <div>
                    <div class="text-xs font-bold text-slate-900 dark:text-slate-100">Kas Tunai Kasir POS</div>
                    <div class="text-[10.5px] text-slate-400">Drawer Outlet Senopati</div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-xs font-bold font-mono text-slate-900 dark:text-slate-100">Rp 7.650.000</div>
                  <div class="text-[10px] text-emerald-600 font-semibold">Aktif</div>
                </div>
              </div>
            </div>
          </div>

          <div class="pt-3 border-t space-y-2" style="border-color: var(--border-color);">
            <button 
              type="button"
              @click="showNewJournalModal = true" 
              class="w-full py-2 px-3 rounded-xl bg-blue-600 hover:bg-blue-700 text-white text-xs font-semibold flex items-center justify-center gap-1.5 shadow-xs transition-colors cursor-pointer"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Buat Jurnal Umum</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Jurnal Transaksi Terbaru -->
      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-50">Jurnal Transaksi Terbaru</h3>
            <p class="text-xs text-slate-400">Entri transaksi debit dan kredit akun akuntansi</p>
          </div>
          <button 
            type="button" 
            @click="activeTab = 'journal'"
            class="text-xs font-semibold text-blue-600 dark:text-blue-400 hover:underline cursor-pointer"
          >
            Lihat Semua Jurnal →
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">No. Bukti / Ref</th>
                <th class="py-3 px-4">Tanggal</th>
                <th class="py-3 px-4">Keterangan / Deskripsi Transaksi</th>
                <th class="py-3 px-4 text-right">Debit</th>
                <th class="py-3 px-4 text-right">Kredit</th>
                <th class="py-3 px-4 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr 
                v-for="j in recentJournals" 
                :key="j.id || j.ref"
                class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="py-3.5 px-4 font-mono font-bold text-blue-600 dark:text-blue-400">{{ j.reference_no || j.ref }}</td>
                <td class="py-3.5 px-4 text-slate-500 font-mono text-[11px]">{{ j.entry_date || j.date }}</td>
                <td class="py-3.5 px-4 font-medium text-slate-800 dark:text-slate-200">{{ j.description }}</td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(j.total_debit || j.debit) }}
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(j.total_credit || j.credit) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <StatusBadgeV2 :status="j.status === 'posted' ? 'completed' : 'pending'" :label="j.status || 'Posted'" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 2: JURNAL UMUM (JOURNAL ENTRIES)                                    -->
    <!-- ======================================================================= -->
    <div v-else-if="activeTab === 'journal'" class="space-y-4">
      <div 
        class="p-4 rounded-2xl border shadow-xs flex flex-col sm:flex-row items-center justify-between gap-3"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="relative w-full sm:w-80">
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            v-model="journalSearch"
            type="text"
            placeholder="Cari referensi atau deskripsi jurnal..."
            class="w-full pl-9 pr-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500"
            style="border-color: var(--border-color);"
          />
        </div>
        <div class="flex items-center gap-2 w-full sm:w-auto justify-end">
          <select
            v-model="journalStatusFilter"
            class="px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-700 dark:text-slate-300 outline-none"
            style="border-color: var(--border-color);"
          >
            <option value="all">Semua Status</option>
            <option value="posted">Posted</option>
            <option value="draft">Draft</option>
          </select>
        </div>
      </div>

      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">No. Referensi</th>
                <th class="py-3 px-4">Tanggal Entri</th>
                <th class="py-3 px-4">Deskripsi Pembukuan</th>
                <th class="py-3 px-4 text-right">Total Debit</th>
                <th class="py-3 px-4 text-right">Total Kredit</th>
                <th class="py-3 px-4 text-center">Status</th>
                <th class="py-3 px-4 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr 
                v-for="j in filteredJournals" 
                :key="j.id"
                class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="py-3.5 px-4 font-mono font-bold text-blue-600 dark:text-blue-400">{{ j.reference_no }}</td>
                <td class="py-3.5 px-4 text-slate-500 font-mono text-[11px]">{{ j.entry_date }}</td>
                <td class="py-3.5 px-4 font-medium text-slate-800 dark:text-slate-200">{{ j.description }}</td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(j.total_debit) }}
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(j.total_credit) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <StatusBadgeV2 :status="j.status === 'posted' ? 'completed' : 'pending'" :label="j.status" />
                </td>
                <td class="py-3.5 px-4 text-center">
                  <button
                    type="button"
                    @click="viewJournalLines(j)"
                    class="px-2.5 py-1 text-[11px] font-semibold text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-950/40 rounded-lg cursor-pointer"
                  >
                    Rincian
                  </button>
                </td>
              </tr>
              <tr v-if="filteredJournals.length === 0">
                <td colspan="7" class="py-8 text-center text-slate-400 text-xs">
                  Tidak ada entri jurnal yang ditemukan.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 3: BAGAN AKUN (CHART OF ACCOUNTS / COA)                              -->
    <!-- ======================================================================= -->
    <div v-else-if="activeTab === 'coa'" class="space-y-4">
      <div 
        class="p-4 rounded-2xl border shadow-xs flex flex-col sm:flex-row items-center justify-between gap-3"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="relative w-full sm:w-80">
          <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            v-model="coaSearch"
            type="text"
            placeholder="Cari kode atau nama akun..."
            class="w-full pl-9 pr-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500"
            style="border-color: var(--border-color);"
          />
        </div>
        <div class="flex items-center gap-2 w-full sm:w-auto justify-end">
          <select
            v-model="coaTypeFilter"
            class="px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-700 dark:text-slate-300 outline-none"
            style="border-color: var(--border-color);"
          >
            <option value="all">Semua Tipe Akun</option>
            <option value="asset">Aset (Assets)</option>
            <option value="liability">Kewajiban (Liabilities)</option>
            <option value="equity">Ekuitas (Equity)</option>
            <option value="revenue">Pendapatan (Revenue)</option>
            <option value="expense">Beban (Expenses)</option>
          </select>
        </div>
      </div>

      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">Kode Akun</th>
                <th class="py-3 px-4">Nama Akun Buku Besar</th>
                <th class="py-3 px-4">Klasifikasi Tipe</th>
                <th class="py-3 px-4">Posisi Normal</th>
                <th class="py-3 px-4 text-right">Saldo Berjalan</th>
                <th class="py-3 px-4 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr 
                v-for="acc in filteredAccounts" 
                :key="acc.id"
                class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="py-3.5 px-4 font-mono font-bold text-slate-900 dark:text-white">{{ acc.code }}</td>
                <td class="py-3.5 px-4 font-medium text-slate-800 dark:text-slate-200">
                  <div>{{ acc.name }}</div>
                  <div class="text-[11px] text-slate-400">{{ acc.description }}</div>
                </td>
                <td class="py-3.5 px-4">
                  <span class="inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-blue-50 dark:bg-blue-950/60 text-blue-600 dark:text-blue-400">
                    {{ acc.account_type }}
                  </span>
                </td>
                <td class="py-3.5 px-4 uppercase text-[11px] text-slate-500 font-mono">
                  {{ acc.normal_balance }}
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(acc.balance) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <StatusBadgeV2 :status="acc.is_active ? 'completed' : 'cancelled'" :label="acc.is_active ? 'Aktif' : 'Nonaktif'" />
                </td>
              </tr>
              <tr v-if="filteredAccounts.length === 0">
                <td colspan="6" class="py-8 text-center text-slate-400 text-xs">
                  Tidak ada bagan akun yang cocok.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 4: BEBAN & PENGELUARAN (EXPENSES & PETTY CASH)                      -->
    <!-- ======================================================================= -->
    <div v-else-if="activeTab === 'expenses'" class="space-y-4">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <StatCardV2
          title="Total Pengeluaran Bulan Ini"
          :value="`Rp ${formatNum(totalExpenseMonth)}`"
          subtitle="Petty cash & biaya operasional"
          color="rose"
          :icon="Receipt"
        />
        <StatCardV2
          title="Menunggu Persetujuan"
          :value="`${pendingExpenseCount} Pengajuan`"
          subtitle="Persetujuan Manager Finance"
          color="amber"
          :icon="Scale"
        />
        <StatCardV2
          title="Selesai Dicairkan"
          :value="`${expenses.filter(e => e.status === 'reimbursed').length} Transaksi`"
          subtitle="Sudah diposting ke jurnal"
          color="emerald"
          :icon="CheckCircle2"
        />
      </div>

      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">No. Pengeluaran</th>
                <th class="py-3 px-4">Tanggal</th>
                <th class="py-3 px-4">Kategori & Keterangan</th>
                <th class="py-3 px-4">Diajukan Oleh</th>
                <th class="py-3 px-4 text-right">Nominal</th>
                <th class="py-3 px-4 text-center">Status</th>
                <th class="py-3 px-4 text-center">Aksi Approval</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr 
                v-for="e in expenses" 
                :key="e.id"
                class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="py-3.5 px-4 font-mono font-bold text-slate-900 dark:text-white">{{ e.expense_number }}</td>
                <td class="py-3.5 px-4 text-slate-500 font-mono text-[11px]">{{ e.expense_date }}</td>
                <td class="py-3.5 px-4 font-medium text-slate-800 dark:text-slate-200">
                  <div class="font-bold text-slate-900 dark:text-slate-100">{{ e.category }}</div>
                  <div class="text-[11px] text-slate-400">{{ e.description }}</div>
                </td>
                <td class="py-3.5 px-4 text-slate-600 dark:text-slate-300">{{ e.submitted_by }}</td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(e.amount) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <StatusBadgeV2 :status="e.status === 'reimbursed' ? 'completed' : e.status === 'approved' ? 'in_progress' : 'pending'" :label="e.status" />
                </td>
                <td class="py-3.5 px-4 text-center">
                  <div v-if="e.status === 'submitted'" class="flex items-center justify-center gap-1.5">
                    <button
                      type="button"
                      @click="updateExpenseStatus(e.id, 'approved')"
                      class="px-2 py-1 rounded bg-emerald-600 text-white font-semibold text-[10.5px] cursor-pointer hover:bg-emerald-700"
                    >
                      Setujui
                    </button>
                    <button
                      type="button"
                      @click="updateExpenseStatus(e.id, 'rejected')"
                      class="px-2 py-1 rounded bg-rose-600 text-white font-semibold text-[10.5px] cursor-pointer hover:bg-rose-700"
                    >
                      Tolak
                    </button>
                  </div>
                  <div v-else-if="e.status === 'approved'">
                    <button
                      type="button"
                      @click="updateExpenseStatus(e.id, 'reimbursed')"
                      class="px-2 py-1 rounded bg-blue-600 text-white font-semibold text-[10.5px] cursor-pointer hover:bg-blue-700"
                    >
                      Cairkan Kas
                    </button>
                  </div>
                  <span v-else class="text-[11px] text-slate-400">Selesai</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 5: PAJAK RESTO & KEPATUHAN (TAX PB1)                                -->
    <!-- ======================================================================= -->
    <div v-else-if="activeTab === 'tax'" class="space-y-5">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <StatCardV2
          title="Omset Kena Pajak (PB1)"
          :value="`Rp ${formatNum(taxData.gross_sales || 2944700)}`"
          subtitle="Pajak Restoran Daerah 10%"
          color="blue"
          :icon="Percent"
        />
        <StatCardV2
          title="PB1 Terutang"
          :value="`Rp ${formatNum(taxData.pb1_tax || 294470)}`"
          subtitle="Tarif 10% Penjualan F&B"
          color="amber"
          :icon="Receipt"
        />
        <StatCardV2
          title="Total Kewajiban Pajak"
          :value="`Rp ${formatNum(taxData.total_tax_liability || 5728383)}`"
          subtitle="Termasuk PPh 21 Karyawan"
          color="rose"
          :icon="ShieldCheck"
        />
      </div>

      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
          <h3 class="font-bold text-sm text-slate-900 dark:text-white">Rekapitulasi Transaksi Pajak (Masa Pajak: {{ taxData.period }})</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">Jenis Pajak</th>
                <th class="py-3 px-4">No. Referensi Faktur/Slip</th>
                <th class="py-3 px-4 text-center">Tarif</th>
                <th class="py-3 px-4 text-right">Dasar Pengenaan Pajak (DPP)</th>
                <th class="py-3 px-4 text-right">Nominal Pajak</th>
                <th class="py-3 px-4 text-center">Jatuh Tempo</th>
                <th class="py-3 px-4 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr 
                v-for="t in taxData.tax_transactions" 
                :key="t.id"
                class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="py-3.5 px-4 font-bold text-slate-900 dark:text-white">{{ t.tax_type }}</td>
                <td class="py-3.5 px-4 font-mono text-slate-600 dark:text-slate-300">{{ t.reference }}</td>
                <td class="py-3.5 px-4 text-center font-bold text-blue-600">{{ t.rate }}%</td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(t.taxable_amount) }}
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-amber-600">
                  Rp {{ formatNum(t.tax_amount) }}
                </td>
                <td class="py-3.5 px-4 text-center text-slate-500 font-mono">{{ t.due_date }}</td>
                <td class="py-3.5 px-4 text-center">
                  <StatusBadgeV2 :status="t.status === 'terutang' ? 'pending' : 'completed'" :label="t.status" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 6: ALOKASI ANGGARAN (BUDGETING & COST CONTROL)                      -->
    <!-- ======================================================================= -->
    <div v-else-if="activeTab === 'budget'" class="space-y-5">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <StatCardV2
          title="Total Pagu Anggaran"
          :value="`Rp ${formatNum(budgetData.total_budget || 73000000)}`"
          subtitle="Pagu periode berjalan"
          color="indigo"
          :icon="Wallet"
        />
        <StatCardV2
          title="Realisasi Pengeluaran"
          :value="`Rp ${formatNum(budgetData.total_actual || 85021938)}`"
          subtitle="Akumulasi aktual terbit"
          color="rose"
          :icon="Receipt"
        />
        <StatCardV2
          title="Persentase Penyerapan"
          :value="`${(budgetData.utilization_pct || 116.4).toFixed(1)}%`"
          subtitle="Utilitas anggaran operasional"
          color="amber"
          :icon="TrendingUp"
        />
      </div>

      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
          <h3 class="font-bold text-sm text-slate-900 dark:text-white">Alokasi Anggaran vs Realisasi Pengeluaran (Periode: {{ budgetData.period }})</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">Kategori Pengeluaran</th>
                <th class="py-3 px-4 text-right">Pagu Anggaran</th>
                <th class="py-3 px-4 text-right">Realisasi Aktual</th>
                <th class="py-3 px-4 text-right">Sisa Anggaran</th>
                <th class="py-3 px-4 w-48 text-center">Penyerapan (%)</th>
                <th class="py-3 px-4 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr 
                v-for="b in budgetData.items" 
                :key="b.id"
                class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="py-3.5 px-4 font-bold text-slate-900 dark:text-white">
                  <div>{{ b.category }}</div>
                  <div class="text-[11px] text-slate-400 font-normal">{{ b.notes }}</div>
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100">
                  Rp {{ formatNum(b.budget_amount) }}
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-rose-600">
                  Rp {{ formatNum(b.actual_amount) }}
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold" :class="b.remaining < 0 ? 'text-rose-600' : 'text-emerald-600'">
                  Rp {{ formatNum(b.remaining) }}
                </td>
                <td class="py-3.5 px-4">
                  <div class="flex items-center gap-2">
                    <div class="flex-1 bg-slate-200 dark:bg-slate-700 rounded-full h-2 overflow-hidden">
                      <div 
                        class="h-full rounded-full transition-all duration-300"
                        :class="b.utilization > 100 ? 'bg-rose-500' : b.utilization > 85 ? 'bg-amber-500' : 'bg-emerald-500'"
                        :style="{ width: `${Math.min(b.utilization, 100)}%` }"
                      ></div>
                    </div>
                    <span class="font-mono font-bold text-[11px]" :class="b.utilization > 100 ? 'text-rose-600' : 'text-slate-600'">
                      {{ b.utilization.toFixed(0) }}%
                    </span>
                  </div>
                </td>
                <td class="py-3.5 px-4 text-center">
                  <StatusBadgeV2 :status="b.status === 'on_track' ? 'completed' : b.status === 'warning' ? 'in_progress' : 'cancelled'" :label="b.status === 'on_track' ? 'Aman' : b.status === 'warning' ? 'Mendekati Limit' : 'Over Budget'" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- TAB 7: REKONSILIASI BANK (BANK RECONCILIATION)                          -->
    <!-- ======================================================================= -->
    <div v-else-if="activeTab === 'reconciliation'" class="space-y-6">
      <!-- 3 Balance Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div 
          class="p-5 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <p class="text-xs font-medium text-slate-500 mb-1">Saldo Rekening Koran Bank</p>
          <p class="text-xl font-bold font-mono text-slate-900 dark:text-white">Rp {{ formatNum(reconData.statement_balance) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">{{ reconData.bank_name }}</p>
        </div>

        <div 
          class="p-5 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <p class="text-xs font-medium text-slate-500 mb-1">Saldo Buku Kas ERP (Akun 1102)</p>
          <p class="text-xl font-bold font-mono text-slate-900 dark:text-white">Rp {{ formatNum(reconData.erp_book_balance) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">General Ledger Buku Besar</p>
        </div>

        <div 
          class="p-5 rounded-2xl border shadow-xs"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <p class="text-xs font-medium text-slate-500 mb-1">Selisih Rekonsiliasi</p>
          <p class="text-xl font-bold font-mono text-blue-600 dark:text-blue-400">Rp {{ formatNum(reconData.difference) }}</p>
          <p class="text-[11px] text-slate-400 mt-1">Cocok: {{ reconData.matched_count }} | Belum: {{ reconData.unmatched_count }}</p>
        </div>
      </div>

      <!-- Action & Auto Match Toolbar -->
      <div 
        class="flex flex-col sm:flex-row items-center justify-between gap-4 p-4 rounded-2xl border shadow-xs"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="flex items-center gap-3 w-full sm:w-auto">
          <label class="text-xs font-bold text-slate-600 dark:text-slate-400 whitespace-nowrap">Pilih Akun Bank:</label>
          <select 
            v-model="selectedBank"
            @change="fetchReconciliation"
            class="px-3 py-2 border rounded-xl text-xs bg-transparent text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            style="border-color: var(--border-color);"
          >
            <option value="Bank BCA Operasional (1102)">Bank BCA Operasional (1102)</option>
            <option value="Bank Mandiri Operasional (1103)">Bank Mandiri Operasional (1103)</option>
          </select>
        </div>
        <div class="flex items-center gap-3 w-full sm:w-auto justify-end">
          <button 
            type="button"
            @click="autoMatchAll"
            class="inline-flex items-center gap-2 px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-semibold shadow-xs transition-colors cursor-pointer"
          >
            <CheckCircle2 class="w-4 h-4" />
            <span>Auto-Match Sesuai Mutasi</span>
          </button>
        </div>
      </div>

      <!-- Bank Statement Mutasi Table -->
      <div 
        class="rounded-2xl border shadow-xs overflow-hidden"
        style="background-color: var(--bg-primary); border-color: var(--border-color);"
      >
        <div class="p-4 border-b flex items-center justify-between" style="border-color: var(--border-color);">
          <h3 class="font-bold text-slate-900 dark:text-white text-sm">Mutasi Rekening Koran vs Pencatatan Kas ERP</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b text-slate-400 font-semibold tracking-wider uppercase text-[10.5px]" style="border-color: var(--border-color); background-color: var(--bg-content);">
                <th class="py-3 px-4">Tanggal Mutasi</th>
                <th class="py-3 px-4">No. Referensi Bank</th>
                <th class="py-3 px-4">Keterangan Transaksi</th>
                <th class="py-3 px-4 text-center">Jenis</th>
                <th class="py-3 px-4 text-right">Nominal (Rp)</th>
                <th class="py-3 px-4 text-center">Status Cocok</th>
                <th class="py-3 px-4 text-center">Aksi Rekonsiliasi</th>
              </tr>
            </thead>
            <tbody class="divide-y" style="border-color: var(--border-color);">
              <tr 
                v-for="s in reconData.statements" 
                :key="s.id" 
                class="hover:bg-slate-50/60 dark:hover:bg-slate-800/40 transition-colors"
              >
                <td class="py-3.5 px-4 text-slate-600 dark:text-slate-300 font-mono whitespace-nowrap">{{ s.statement_date }}</td>
                <td class="py-3.5 px-4 font-mono font-bold text-slate-900 dark:text-white whitespace-nowrap">{{ s.ref_number }}</td>
                <td class="py-3.5 px-4 text-slate-700 dark:text-slate-200 max-w-sm">{{ s.description }}</td>
                <td class="py-3.5 px-4 text-center">
                  <span 
                    class="inline-block px-2 py-0.5 rounded text-[10px] font-bold uppercase"
                    :class="s.type === 'credit' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-rose-50 text-rose-700 border border-rose-200'"
                  >
                    {{ s.type === 'credit' ? 'Kredit (Masuk)' : 'Debet (Keluar)' }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-right font-mono font-bold text-slate-900 dark:text-slate-100 whitespace-nowrap">
                  Rp {{ formatNum(s.amount) }}
                </td>
                <td class="py-3.5 px-4 text-center">
                  <StatusBadgeV2 :status="s.match_status === 'matched' ? 'completed' : 'pending'" :label="s.match_status === 'matched' ? 'Matched' : 'Unmatched'" />
                </td>
                <td class="py-3.5 px-4 text-center">
                  <button 
                    type="button"
                    @click="toggleMatchReconciliation(s)"
                    class="px-2.5 py-1 rounded-lg text-[11px] font-semibold transition-colors cursor-pointer"
                    :class="s.match_status === 'matched' 
                      ? 'bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 text-slate-700 dark:text-slate-300' 
                      : 'bg-blue-600 hover:bg-blue-700 text-white shadow-2xs'"
                  >
                    {{ s.match_status === 'matched' ? 'Batalkan Cocok' : 'Cocokkan Transaksi' }}
                  </button>
                </td>
              </tr>
              <tr v-if="reconData.statements.length === 0">
                <td colspan="7" class="py-8 text-center text-slate-400 text-xs">
                  Tidak ada data mutasi rekening koran untuk periode ini.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ======================================================================= -->
    <!-- MODALS SECTION (View Lines, Create Journal, Create COA, Create Expense) -->
    <!-- ======================================================================= -->

    <!-- Modal 1: Rincian Baris Jurnal -->
    <Teleport to="body">
      <div v-if="showLinesModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4">
        <div 
          class="rounded-2xl max-w-2xl w-full shadow-2xl border overflow-hidden animate-in zoom-in-95 duration-200"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="p-5 border-b flex items-center justify-between" style="border-color: var(--border-color);">
            <h3 class="font-bold text-slate-900 dark:text-white text-base">Rincian Baris Jurnal: {{ selectedJournalRef }}</h3>
            <button @click="showLinesModal = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 cursor-pointer">✕</button>
          </div>
          <div class="p-5 overflow-x-auto">
            <table class="w-full text-left text-xs border-collapse">
              <thead>
                <tr class="border-b text-slate-400 uppercase text-[10px]" style="border-color: var(--border-color);">
                  <th class="py-2.5 px-3">Kode Akun</th>
                  <th class="py-2.5 px-3">Nama Akun</th>
                  <th class="py-2.5 px-3">Deskripsi</th>
                  <th class="py-2.5 px-3 text-right">Debit</th>
                  <th class="py-2.5 px-3 text-right">Kredit</th>
                </tr>
              </thead>
              <tbody class="divide-y" style="border-color: var(--border-color);">
                <tr v-for="l in selectedJournalLines" :key="l.id">
                  <td class="py-2.5 px-3 font-mono font-bold text-blue-600">{{ l.account_code }}</td>
                  <td class="py-2.5 px-3 text-slate-800 dark:text-slate-200">{{ l.account_name }}</td>
                  <td class="py-2.5 px-3 text-slate-500">{{ l.description }}</td>
                  <td class="py-2.5 px-3 text-right font-mono font-bold text-slate-900 dark:text-white">Rp {{ formatNum(l.debit) }}</td>
                  <td class="py-2.5 px-3 text-right font-mono font-bold text-slate-900 dark:text-white">Rp {{ formatNum(l.credit) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="p-4 border-t flex justify-end" style="border-color: var(--border-color);">
            <button @click="showLinesModal = false" class="px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 rounded-xl text-xs font-semibold cursor-pointer">
              Tutup
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Modal 2: Buat Jurnal Baru -->
    <Teleport to="body">
      <div v-if="showNewJournalModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4">
        <div 
          class="rounded-2xl max-w-2xl w-full shadow-2xl border overflow-hidden animate-in zoom-in-95 duration-200"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="p-5 border-b flex items-center justify-between" style="border-color: var(--border-color);">
            <h3 class="font-bold text-slate-900 dark:text-white text-base">Buat Entri Jurnal Umum Baru</h3>
            <button @click="showNewJournalModal = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 cursor-pointer">✕</button>
          </div>
          <form @submit.prevent="submitNewJournal" class="p-5 space-y-4">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">No. Referensi</label>
                <input v-model="newJournalForm.reference_no" type="text" placeholder="JV-2026-..." required class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Tanggal Entri</label>
                <input v-model="newJournalForm.entry_date" type="date" required class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
              </div>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Keterangan / Memo</label>
              <input v-model="newJournalForm.description" type="text" placeholder="Deskripsi transaksi..." required class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>

            <!-- Lines -->
            <div class="space-y-2 pt-2 border-t" style="border-color: var(--border-color);">
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold text-slate-800 dark:text-slate-200">Baris Debit & Kredit</span>
                <button type="button" @click="addJournalLine" class="text-xs text-blue-600 dark:text-blue-400 font-semibold hover:underline cursor-pointer">+ Tambah Baris</button>
              </div>
              <div v-for="(line, idx) in newJournalForm.lines" :key="idx" class="grid grid-cols-12 gap-2 items-center">
                <select v-model="line.account_id" class="col-span-5 px-2.5 py-1.5 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none" style="border-color: var(--border-color);">
                  <option v-for="acc in accounts" :key="acc.id" :value="acc.id">{{ acc.code }} - {{ acc.name }}</option>
                </select>
                <input v-model.number="line.debit" type="number" placeholder="Debit" class="col-span-3 px-2.5 py-1.5 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none text-right font-mono" style="border-color: var(--border-color);" />
                <input v-model.number="line.credit" type="number" placeholder="Kredit" class="col-span-3 px-2.5 py-1.5 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none text-right font-mono" style="border-color: var(--border-color);" />
                <button type="button" @click="removeJournalLine(idx)" class="col-span-1 text-slate-400 hover:text-rose-500 text-center cursor-pointer">✕</button>
              </div>
              <div class="flex items-center justify-between text-xs font-bold pt-2">
                <span>Total Debit: Rp {{ formatNum(newJournalTotalDebit) }}</span>
                <span>Total Kredit: Rp {{ formatNum(newJournalTotalCredit) }}</span>
              </div>
              <p v-if="!isNewJournalBalanced" class="text-xs text-rose-500 font-semibold">⚠️ Total Debit dan Kredit harus seimbang!</p>
            </div>

            <div class="p-4 border-t flex justify-end gap-2" style="border-color: var(--border-color);">
              <button type="button" @click="showNewJournalModal = false" class="px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 text-slate-700 dark:text-slate-200 rounded-xl text-xs font-semibold cursor-pointer">Batal</button>
              <button type="submit" :disabled="!isNewJournalBalanced || loading" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold cursor-pointer disabled:opacity-50">Posting Jurnal</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Modal 3: Tambah Bagan Akun (COA) -->
    <Teleport to="body">
      <div v-if="showNewAccountModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4">
        <div 
          class="rounded-2xl max-w-md w-full shadow-2xl border overflow-hidden animate-in zoom-in-95 duration-200"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="p-5 border-b flex items-center justify-between" style="border-color: var(--border-color);">
            <h3 class="font-bold text-slate-900 dark:text-white text-base">Tambah Bagan Akun (COA)</h3>
            <button @click="showNewAccountModal = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 cursor-pointer">✕</button>
          </div>
          <form @submit.prevent="submitNewAccount" class="p-5 space-y-3">
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Kode Akun</label>
              <input v-model="newAccountForm.code" type="text" placeholder="Contoh: 1103" required class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Nama Akun</label>
              <input v-model="newAccountForm.name" type="text" placeholder="Nama akun buku besar..." required class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Klasifikasi Akun</label>
              <select v-model="newAccountForm.account_type" class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none" style="border-color: var(--border-color);">
                <option value="asset">Aset (Assets)</option>
                <option value="liability">Kewajiban (Liabilities)</option>
                <option value="equity">Ekuitas (Equity)</option>
                <option value="revenue">Pendapatan (Revenue)</option>
                <option value="expense">Beban (Expenses)</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Deskripsi Akun</label>
              <input v-model="newAccountForm.description" type="text" placeholder="Keterangan peruntukan akun..." class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>
            <div class="p-4 border-t flex justify-end gap-2" style="border-color: var(--border-color);">
              <button type="button" @click="showNewAccountModal = false" class="px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 text-slate-700 dark:text-slate-200 rounded-xl text-xs font-semibold cursor-pointer">Batal</button>
              <button type="submit" :disabled="loading" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold cursor-pointer">Simpan Akun</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Modal 4: Catat Pengeluaran Baru -->
    <Teleport to="body">
      <div v-if="showNewExpenseModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4">
        <div 
          class="rounded-2xl max-w-md w-full shadow-2xl border overflow-hidden animate-in zoom-in-95 duration-200"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="p-5 border-b flex items-center justify-between" style="border-color: var(--border-color);">
            <h3 class="font-bold text-slate-900 dark:text-white text-base">Catat Pengeluaran Kas Baru</h3>
            <button @click="showNewExpenseModal = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 cursor-pointer">✕</button>
          </div>
          <form @submit.prevent="submitNewExpense" class="p-5 space-y-3">
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Kategori Pengeluaran</label>
              <select v-model="newExpenseForm.category" class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none" style="border-color: var(--border-color);">
                <option value="Operasional & Kas Kecil">Operasional & Kas Kecil</option>
                <option value="Bahan Baku">Bahan Baku</option>
                <option value="Pemeliharaan Peralatan">Pemeliharaan Peralatan</option>
                <option value="Promosi & Pemasaran">Promosi & Pemasaran</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Nominal (Rp)</label>
              <input v-model.number="newExpenseForm.amount" type="number" placeholder="Nominal pengeluaran..." required min="1" class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none font-mono focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Keterangan / Keperluan</label>
              <input v-model="newExpenseForm.description" type="text" placeholder="Detail transaksi..." required class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Diajukan Oleh</label>
              <input v-model="newExpenseForm.submitted_by" type="text" placeholder="Nama staf / kasir..." required class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>
            <div class="p-4 border-t flex justify-end gap-2" style="border-color: var(--border-color);">
              <button type="button" @click="showNewExpenseModal = false" class="px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 text-slate-700 dark:text-slate-200 rounded-xl text-xs font-semibold cursor-pointer">Batal</button>
              <button type="submit" :disabled="loading" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold cursor-pointer">Ajukan Pengeluaran</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Modal 5: Atur Anggaran Baru -->
    <Teleport to="body">
      <div v-if="showNewBudgetModal" class="fixed inset-0 z-50 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4">
        <div 
          class="rounded-2xl max-w-md w-full shadow-2xl border overflow-hidden animate-in zoom-in-95 duration-200"
          style="background-color: var(--bg-primary); border-color: var(--border-color);"
        >
          <div class="p-5 border-b flex items-center justify-between" style="border-color: var(--border-color);">
            <h3 class="font-bold text-slate-900 dark:text-white text-base">Atur Alokasi Anggaran Bulanan</h3>
            <button @click="showNewBudgetModal = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 cursor-pointer">✕</button>
          </div>
          <form @submit.prevent="submitNewBudget" class="p-5 space-y-3">
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Kategori Pengeluaran</label>
              <select v-model="newBudgetForm.category" class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none" style="border-color: var(--border-color);">
                <option value="Bahan Baku & Minuman">Bahan Baku & Minuman</option>
                <option value="Sewa Tempat & Gedung">Sewa Tempat & Gedung</option>
                <option value="Utilitas Listrik, Air & Gas">Utilitas Listrik, Air & Gas</option>
                <option value="Promosi & Pemasaran">Promosi & Pemasaran</option>
                <option value="Operasional & Kas Kecil">Operasional & Kas Kecil</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Pagu Anggaran (Rp)</label>
              <input v-model.number="newBudgetForm.budget_amount" type="number" placeholder="Nominal pagu..." required min="1" class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none font-mono focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-700 dark:text-slate-300 mb-1">Catatan Alokasi</label>
              <input v-model="newBudgetForm.notes" type="text" placeholder="Peruntukan anggaran..." class="w-full px-3 py-2 text-xs rounded-xl border bg-transparent text-slate-900 dark:text-white outline-none focus:border-blue-500" style="border-color: var(--border-color);" />
            </div>
            <div class="p-4 border-t flex justify-end gap-2" style="border-color: var(--border-color);">
              <button type="button" @click="showNewBudgetModal = false" class="px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 text-slate-700 dark:text-slate-200 rounded-xl text-xs font-semibold cursor-pointer">Batal</button>
              <button type="submit" :disabled="loading" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-semibold cursor-pointer">Simpan Pagu</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import PageHeaderV2 from '@/components/v2/common/PageHeaderV2.vue'
import StatCardV2 from '@/components/v2/common/StatCardV2.vue'
import StatusBadgeV2 from '@/components/v2/common/StatusBadgeV2.vue'
import {
  RefreshCw,
  Plus,
  DollarSign,
  Receipt,
  TrendingUp,
  Landmark,
  Search,
  Scale,
  CheckCircle2,
  Percent,
  ShieldCheck,
  Wallet
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

const notifyStore = useNotificationStore()
const loading = ref(false)
const activeTab = ref('overview')

const financeTabs = [
  { id: 'overview', label: 'Ringkasan Eksekutif' },
  { id: 'journal', label: 'Buku Jurnal Umum' },
  { id: 'coa', label: 'Bagan Akun (COA)' },
  { id: 'expenses', label: 'Beban & Pengeluaran' },
  { id: 'tax', label: 'Pajak Resto (PB1)' },
  { id: 'budget', label: 'Alokasi Budget' },
  { id: 'reconciliation', label: 'Rekonsiliasi Bank' }
]

// -----------------------------------------------------------------------------
// OVERVIEW & CHART DATA
// -----------------------------------------------------------------------------
const overviewData = ref({
  total_revenue: 128450000,
  total_expenses: 78320000,
  net_profit: 50130000,
  cash_on_hand: 32500000
})

const monthlyData = [
  { m: 'Jan', rev: 54, exp: 32 },
  { m: 'Feb', rev: 48, exp: 30 },
  { m: 'Mar', rev: 55, exp: 37 },
  { m: 'Apr', rev: 61, exp: 42 },
  { m: 'May', rev: 63, exp: 43 },
  { m: 'Jun', rev: 68, exp: 41 },
  { m: 'Jul', rev: 62, exp: 39 },
  { m: 'Aug', rev: 64, exp: 37 },
  { m: 'Sep', rev: 60, exp: 36 },
  { m: 'Oct', rev: 74, exp: 44 },
  { m: 'Nov', rev: 69, exp: 45 },
  { m: 'Dec', rev: 73, exp: 47 }
]

const recentJournals = ref<any[]>([])

const fetchOverviewData = async () => {
  try {
    loading.value = true
    const res = await axios.get('/api/v1/finance/overview')
    if (res.data?.data) {
      overviewData.value = {
        total_revenue: res.data.data.total_revenue || 128450000,
        total_expenses: res.data.data.total_expenses || 78320000,
        net_profit: res.data.data.net_profit || 50130000,
        cash_on_hand: res.data.data.cash_on_hand || 32500000
      }
    }
  } catch (err) {
    console.error('Failed to load overview:', err)
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// JOURNALS
// -----------------------------------------------------------------------------
const journals = ref<any[]>([])
const journalSearch = ref('')
const journalStatusFilter = ref('all')
const showLinesModal = ref(false)
const showNewJournalModal = ref(false)
const selectedJournalRef = ref('')
const selectedJournalLines = ref<any[]>([])

const newJournalForm = ref({
  reference_no: '',
  entry_date: new Date().toISOString().split('T')[0],
  description: '',
  lines: [
    { account_id: '', description: '', debit: 0, credit: 0 },
    { account_id: '', description: '', debit: 0, credit: 0 }
  ]
})

const newJournalTotalDebit = computed(() => {
  return newJournalForm.value.lines.reduce((sum, l) => sum + (Number(l.debit) || 0), 0)
})

const newJournalTotalCredit = computed(() => {
  return newJournalForm.value.lines.reduce((sum, l) => sum + (Number(l.credit) || 0), 0)
})

const isNewJournalBalanced = computed(() => {
  return (
    newJournalTotalDebit.value > 0 &&
    newJournalTotalDebit.value === newJournalTotalCredit.value
  )
})

const filteredJournals = computed(() => {
  return journals.value.filter(j => {
    const matchSearch = 
      (j.reference_no || '').toLowerCase().includes(journalSearch.value.toLowerCase()) ||
      (j.description || '').toLowerCase().includes(journalSearch.value.toLowerCase())
    const matchStatus = journalStatusFilter.value === 'all' || j.status === journalStatusFilter.value
    return matchSearch && matchStatus
  })
})

const fetchJournals = async () => {
  try {
    const res = await axios.get('/api/v1/finance/journals')
    if (res.data?.data) {
      journals.value = res.data.data
      recentJournals.value = res.data.data.slice(0, 5)
    }
  } catch (err) {
    console.error('Failed to load journals:', err)
  }
}

const viewJournalLines = async (journal: any) => {
  selectedJournalRef.value = journal.reference_no
  try {
    const res = await axios.get(`/api/v1/finance/journals/${journal.id}/lines`)
    selectedJournalLines.value = res.data?.data || []
    showLinesModal.value = true
  } catch (err) {
    console.error('Failed to load journal lines:', err)
    notifyStore.error('Gagal memuat rincian jurnal baris')
  }
}

const addJournalLine = () => {
  newJournalForm.value.lines.push({
    account_id: accounts.value[0]?.id || '',
    description: '',
    debit: 0,
    credit: 0
  })
}

const removeJournalLine = (idx: number) => {
  if (newJournalForm.value.lines.length > 2) {
    newJournalForm.value.lines.splice(idx, 1)
  }
}

const submitNewJournal = async () => {
  if (!isNewJournalBalanced.value) {
    notifyStore.error('Jurnal harus seimbang: Total Debit harus sama dengan Total Kredit!')
    return
  }
  try {
    loading.value = true
    await axios.post('/api/v1/finance/journals', newJournalForm.value)
    notifyStore.success('Jurnal umum berhasil diposting!')
    showNewJournalModal.value = false
    newJournalForm.value = {
      reference_no: '',
      entry_date: new Date().toISOString().split('T')[0],
      description: '',
      lines: [
        { account_id: accounts.value[0]?.id || '', description: '', debit: 0, credit: 0 },
        { account_id: accounts.value[1]?.id || '', description: '', debit: 0, credit: 0 }
      ]
    }
    await fetchJournals()
    await fetchOverviewData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal memposting jurnal')
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// CHART OF ACCOUNTS (COA)
// -----------------------------------------------------------------------------
const accounts = ref<any[]>([])
const coaSearch = ref('')
const coaTypeFilter = ref('all')
const showNewAccountModal = ref(false)

const newAccountForm = ref({
  code: '',
  name: '',
  account_type: 'asset',
  description: '',
  is_active: true
})

const filteredAccounts = computed(() => {
  return accounts.value.filter(a => {
    const matchSearch = 
      (a.code || '').toLowerCase().includes(coaSearch.value.toLowerCase()) ||
      (a.name || '').toLowerCase().includes(coaSearch.value.toLowerCase())
    const matchType = coaTypeFilter.value === 'all' || a.account_type === coaTypeFilter.value
    return matchSearch && matchType
  })
})

const fetchAccounts = async () => {
  try {
    const res = await axios.get('/api/v1/finance/coa')
    if (res.data?.data) {
      accounts.value = res.data.data
      if (newJournalForm.value.lines[0].account_id === '' && accounts.value.length > 0) {
        newJournalForm.value.lines[0].account_id = accounts.value[0].id
        newJournalForm.value.lines[1].account_id = accounts.value[1]?.id || accounts.value[0].id
      }
    }
  } catch (err) {
    console.error('Failed to load COA:', err)
  }
}

const submitNewAccount = async () => {
  if (!newAccountForm.value.code || !newAccountForm.value.name) {
    notifyStore.error('Kode dan Nama Akun wajib diisi!')
    return
  }
  try {
    loading.value = true
    await axios.post('/api/v1/finance/coa', newAccountForm.value)
    notifyStore.success('Bagan Akun (COA) berhasil ditambahkan!')
    showNewAccountModal.value = false
    newAccountForm.value = {
      code: '',
      name: '',
      account_type: 'asset',
      description: '',
      is_active: true
    }
    await fetchAccounts()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menambah akun')
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// EXPENSES
// -----------------------------------------------------------------------------
const expenses = ref<any[]>([])
const showNewExpenseModal = ref(false)

const newExpenseForm = ref({
  category: 'Operasional & Kas Kecil',
  account_id: '',
  amount: 0,
  expense_date: new Date().toISOString().split('T')[0],
  submitted_by: 'Kasir / Staff',
  description: '',
  receipt_url: ''
})

const totalExpenseMonth = computed(() => {
  return expenses.value.reduce((sum, e) => sum + (Number(e.amount) || 0), 0)
})

const pendingExpenseCount = computed(() => {
  return expenses.value.filter(e => e.status === 'submitted').length
})

const fetchExpenses = async () => {
  try {
    const res = await axios.get('/api/v1/finance/expenses')
    if (res.data?.data) {
      expenses.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load expenses:', err)
  }
}

const submitNewExpense = async () => {
  if (newExpenseForm.value.amount <= 0 || !newExpenseForm.value.description) {
    notifyStore.error('Nominal harus lebih dari 0 dan Keterangan wajib diisi!')
    return
  }
  try {
    loading.value = true
    await axios.post('/api/v1/finance/expenses', newExpenseForm.value)
    notifyStore.success('Pengeluaran berhasil diajukan!')
    showNewExpenseModal.value = false
    newExpenseForm.value = {
      category: 'Operasional & Kas Kecil',
      account_id: '',
      amount: 0,
      expense_date: new Date().toISOString().split('T')[0],
      submitted_by: 'Kasir / Staff',
      description: '',
      receipt_url: ''
    }
    await fetchExpenses()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal mengajukan pengeluaran')
  } finally {
    loading.value = false
  }
}

const updateExpenseStatus = async (id: string, status: 'approved' | 'rejected' | 'reimbursed') => {
  try {
    loading.value = true
    await axios.put(`/api/v1/finance/expenses/${id}/status`, {
      status,
      approved_by: 'Finance Manager'
    })
    const actionLabel = status === 'reimbursed' ? 'dicairkan & dijurnal' : status === 'approved' ? 'disetujui' : 'ditolak'
    notifyStore.success(`Pengeluaran berhasil ${actionLabel}!`)
    await fetchExpenses()
    await fetchJournals()
    await fetchOverviewData()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal memperbarui status pengeluaran')
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// TAX MANAGEMENT
// -----------------------------------------------------------------------------
const taxData = ref({
  gross_sales: 0,
  pb1_tax: 0,
  ppn_masukan: 0,
  pph21_tax: 0,
  total_tax_liability: 0,
  period: 'September 2026',
  tax_transactions: [] as any[]
})

const fetchTaxSummary = async () => {
  try {
    const res = await axios.get('/api/v1/finance/tax/summary')
    if (res.data?.data) {
      taxData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load tax summary:', err)
  }
}

// -----------------------------------------------------------------------------
// BUDGETING
// -----------------------------------------------------------------------------
const budgetData = ref({
  period: '2026-09',
  total_budget: 0,
  total_actual: 0,
  remaining: 0,
  utilization_pct: 0,
  items: [] as any[]
})

const showNewBudgetModal = ref(false)
const newBudgetForm = ref({
  period: '2026-09',
  category: 'Bahan Baku & Minuman',
  budget_amount: 0,
  notes: ''
})

const fetchBudgets = async () => {
  try {
    const res = await axios.get('/api/v1/finance/budgets')
    if (res.data?.data) {
      budgetData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load budgets:', err)
  }
}

const submitNewBudget = async () => {
  if (newBudgetForm.value.budget_amount <= 0) {
    notifyStore.error('Nominal anggaran harus lebih dari 0!')
    return
  }
  try {
    loading.value = true
    await axios.post('/api/v1/finance/budgets', newBudgetForm.value)
    notifyStore.success('Anggaran berhasil diperbarui!')
    showNewBudgetModal.value = false
    await fetchBudgets()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan anggaran')
  } finally {
    loading.value = false
  }
}

// -----------------------------------------------------------------------------
// BANK RECONCILIATION
// -----------------------------------------------------------------------------
const reconData = ref({
  bank_name: 'Bank BCA Operasional (1102)',
  period: '2026-09',
  statement_balance: 0,
  erp_book_balance: 0,
  difference: 0,
  matched_count: 0,
  unmatched_count: 0,
  statements: [] as any[]
})

const selectedBank = ref('Bank BCA Operasional (1102)')

const fetchReconciliation = async () => {
  try {
    const res = await axios.get(`/api/v1/finance/reconciliation?bank_name=${encodeURIComponent(selectedBank.value)}`)
    if (res.data?.data) {
      reconData.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to load bank reconciliation:', err)
  }
}

const toggleMatchReconciliation = async (item: any) => {
  const newStatus = item.match_status === 'matched' ? 'unmatched' : 'matched'
  try {
    await axios.post('/api/v1/finance/reconciliation/match', {
      id: item.id,
      match_status: newStatus
    })
    item.match_status = newStatus
    if (newStatus === 'matched') {
      reconData.value.matched_count++
      reconData.value.unmatched_count--
      notifyStore.success(`Transaksi ${item.ref_number} ditandai Cocok (Matched)!`)
    } else {
      reconData.value.matched_count--
      reconData.value.unmatched_count++
      notifyStore.info(`Kesesuaian ${item.ref_number} dibatalkan.`)
    }
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal mengubah status rekonsiliasi')
  }
}

const autoMatchAll = () => {
  let matched = 0
  reconData.value.statements.forEach(async (s: any) => {
    if (s.match_status !== 'matched') {
      s.match_status = 'matched'
      matched++
      try {
        await axios.post('/api/v1/finance/reconciliation/match', {
          id: s.id,
          match_status: 'matched'
        })
      } catch (e) {
        console.error(e)
      }
    }
  })
  reconData.value.matched_count += matched
  reconData.value.unmatched_count = 0
  notifyStore.success(`Auto-Match selesai! ${matched} transaksi berhasil dicocokkan otomatis.`)
}

// -----------------------------------------------------------------------------
// TAB SWITCHING & REFRESH
// -----------------------------------------------------------------------------
const switchTab = async (tabId: string) => {
  activeTab.value = tabId
  await loadTabData(tabId)
}

const refreshCurrentTab = async () => {
  loading.value = true
  await loadTabData(activeTab.value)
  loading.value = false
  notifyStore.info('Data berhasil diperbarui', 'Sinkronisasi')
}

const loadTabData = async (tabId: string) => {
  if (tabId === 'overview') {
    await Promise.all([fetchOverviewData(), fetchJournals()])
  } else if (tabId === 'journal') {
    await fetchJournals()
  } else if (tabId === 'coa') {
    await fetchAccounts()
  } else if (tabId === 'expenses') {
    await fetchExpenses()
  } else if (tabId === 'tax') {
    await fetchTaxSummary()
  } else if (tabId === 'budget') {
    await fetchBudgets()
  } else if (tabId === 'reconciliation') {
    await fetchReconciliation()
  }
}

const formatNum = (num: number) => Number(num || 0).toLocaleString('id-ID')

onMounted(async () => {
  await Promise.all([
    fetchOverviewData(),
    fetchJournals(),
    fetchAccounts()
  ])
})
</script>

<style scoped>
.finance-v2 {
  min-height: calc(100vh - 120px);
}
</style>
