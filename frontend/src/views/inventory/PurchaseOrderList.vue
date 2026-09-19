<template>
  <div class="space-y-6">
    <!-- Header & P2P Cycle Overview -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2.5">
          <h2 class="text-2xl font-black text-slate-900 tracking-tight">Siklus Pengadaan & Pembayaran</h2>
          <span class="px-2.5 py-0.5 rounded-full text-xs font-black bg-blue-100 text-blue-800 border border-blue-200 shadow-xs">
            ERP P2P Hub
          </span>
        </div>
        <p class="text-xs text-slate-500 mt-1">
          Alur standar pengadaan barang: Permintaan (PR) &rarr; Pemesanan (PO) &rarr; Penerimaan Stok (GRN) &rarr; Faktur & PPN 11% &rarr; Pelunasan Kas/Bank
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          @click="refreshCurrentTab"
          class="p-2.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 text-slate-600 transition-colors shadow-xs cursor-pointer"
          title="Segarkan Data"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>

        <!-- Dynamic Action Button based on Active Tab -->
        <button
          v-if="activeTab === 'pr'"
          @click="openCreatePRModal"
          class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>Buat PR Baru</span>
        </button>
        <button
          v-else-if="activeTab === 'po'"
          @click="openCreatePOModal"
          class="px-4 py-2.5 bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white text-xs font-bold rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>Buat PO Baru</span>
        </button>
        <button
          v-else-if="activeTab === 'grn'"
          @click="openCreateGRNModal"
          class="px-4 py-2.5 bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white text-xs font-bold rounded-xl shadow-md shadow-emerald-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <PackageCheck class="w-4 h-4" />
          <span>Catat Penerimaan (GRN)</span>
        </button>
        <button
          v-else-if="activeTab === 'invoice'"
          @click="openCreateInvoiceModal"
          class="px-4 py-2.5 bg-amber-600 hover:bg-amber-700 active:bg-amber-800 text-white text-xs font-bold rounded-xl shadow-md shadow-amber-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <FileText class="w-4 h-4" />
          <span>Input Faktur Vendor</span>
        </button>
        <button
          v-else-if="activeTab === 'payment'"
          @click="openCreatePaymentModal"
          class="px-4 py-2.5 bg-violet-600 hover:bg-violet-700 active:bg-violet-800 text-white text-xs font-bold rounded-xl shadow-md shadow-violet-600/30 transition-all cursor-pointer flex items-center gap-1.5"
        >
          <CreditCard class="w-4 h-4" />
          <span>Pelunasan Tagihan</span>
        </button>
      </div>
    </div>

    <!-- P2P 5-Stage Stepper Navigation -->
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-2.5 bg-slate-100/70 p-1.5 rounded-2xl border border-slate-200/80">
      <button
        v-for="(step, idx) in p2pSteps"
        :key="step.id"
        @click="activeTab = step.id"
        class="flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-bold transition-all text-left cursor-pointer"
        :class="activeTab === step.id ? 'bg-white text-blue-600 shadow-md shadow-slate-200' : 'text-slate-600 hover:text-slate-900 hover:bg-white/60'"
      >
        <span
          class="w-6 h-6 rounded-lg flex items-center justify-center text-[10px] font-black shrink-0"
          :class="activeTab === step.id ? 'bg-blue-600 text-white' : 'bg-slate-200 text-slate-600'"
        >
          {{ idx + 1 }}
        </span>
        <div class="min-w-0">
          <div class="truncate text-[11px] leading-tight">{{ step.label }}</div>
          <div class="text-[9px] font-medium text-slate-400 uppercase tracking-wider">{{ step.code }}</div>
        </div>
      </button>
    </div>

    <!-- Search & Filters -->
    <div class="flex flex-col sm:flex-row items-center justify-between gap-3 bg-white p-3 rounded-2xl border border-slate-200/80 shadow-xs">
      <div class="relative w-full sm:w-80">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari nomor dokumen, supplier, catatan..."
          class="w-full pl-9 pr-4 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 focus:ring-2 focus:ring-blue-100 outline-hidden transition-all"
        />
      </div>

      <div class="flex items-center gap-2 w-full sm:w-auto overflow-x-auto pb-1 sm:pb-0">
        <button
          v-for="status in currentFilterOptions"
          :key="status.value"
          @click="activeStatusFilter = status.value"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold whitespace-nowrap transition-colors cursor-pointer"
          :class="activeStatusFilter === status.value ? 'bg-blue-600 text-white shadow-xs' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'"
        >
          {{ status.label }}
        </button>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- TAB 1: PURCHASE REQUISITION (PR)                                         -->
    <!-- ========================================================================= -->
    <div v-if="activeTab === 'pr'" class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 border-b border-slate-200/80 text-[11px] uppercase font-bold text-slate-500 tracking-wider">
            <tr>
              <th class="py-3 px-4">Nomor PR</th>
              <th class="py-3 px-4">Departemen</th>
              <th class="py-3 px-4">Tgl Pengajuan</th>
              <th class="py-3 px-4">Dibutuhkan</th>
              <th class="py-3 px-4">Estimasi Nilai</th>
              <th class="py-3 px-4">Status Alur</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="7" class="py-8 text-center text-slate-400">
                <RotateCw class="w-5 h-5 animate-spin mx-auto mb-2 text-blue-600" />
                Memuat data PR...
              </td>
            </tr>
            <tr v-else-if="filteredPRList.length === 0">
              <td colspan="7" class="py-12 text-center text-slate-400">
                Belum ada Purchase Requisition (PR). Klik "Buat PR Baru" untuk mengajukan pengadaan bahan baku.
              </td>
            </tr>
            <tr v-for="pr in filteredPRList" :key="pr.id" class="hover:bg-slate-50/70 transition-colors">
              <td class="py-3 px-4">
                <div class="flex items-center gap-2">
                  <span class="font-mono font-bold text-blue-600">{{ pr.pr_number }}</span>
                  <span
                    v-if="pr.pr_number?.startsWith('PR-MENU-') || pr.notes?.includes('Ref Menu')"
                    class="px-2 py-0.5 rounded-md text-[9px] font-black uppercase tracking-wider bg-cyan-100 text-cyan-800 border border-cyan-200"
                  >
                    Restock Menu
                  </span>
                </div>
                <div v-if="pr.notes" class="text-[10px] text-slate-400 truncate max-w-xs mt-0.5" :title="pr.notes">
                  {{ pr.notes }}
                </div>
              </td>
              <td class="py-3 px-4 font-semibold text-slate-900">{{ pr.department }}</td>
              <td class="py-3 px-4">{{ formatDate(pr.created_at) }}</td>
              <td class="py-3 px-4">{{ formatDate(pr.required_date) }}</td>
              <td class="py-3 px-4 font-semibold text-slate-800">Rp {{ formatNum(pr.estimated_total) }}</td>
              <td class="py-3 px-4">
                <span :class="getStatusBadgeClass(pr.status)" class="px-2.5 py-0.5 rounded-full text-[10px] font-bold">
                  {{ pr.status === 'pending_approval' ? 'Pending Gudang' : pr.status }}
                </span>
              </td>
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <button
                    v-if="pr.status === 'draft' || pr.status === 'submitted' || pr.status === 'pending_approval'"
                    @click="approvePR(pr)"
                    class="p-1.5 rounded-lg text-emerald-600 hover:bg-emerald-50 cursor-pointer"
                    title="Setujui PR"
                  >
                    <CheckCircle2 class="w-4 h-4" />
                  </button>
                  <button
                    v-if="pr.status === 'approved' || pr.status === 'pending_approval'"
                    @click="openConvertPRModal(pr)"
                    class="px-2.5 py-1 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-bold text-[11px] shadow-xs cursor-pointer flex items-center gap-1"
                    title="Konversi ke Purchase Order (PO)"
                  >
                    <span>Jadikan PO</span>
                    <ArrowRight class="w-3 h-3" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- TAB 2: PURCHASE ORDER (PO)                                               -->
    <!-- ========================================================================= -->
    <div v-else-if="activeTab === 'po'" class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 border-b border-slate-200/80 text-[11px] uppercase font-bold text-slate-500 tracking-wider">
            <tr>
              <th class="py-3 px-4">Nomor PO</th>
              <th class="py-3 px-4">Ref. PR</th>
              <th class="py-3 px-4">Supplier / Vendor</th>
              <th class="py-3 px-4">Tgl Pemesanan</th>
              <th class="py-3 px-4">Estimasi Sampai</th>
              <th class="py-3 px-4">Total Biaya</th>
              <th class="py-3 px-4">Status PO</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="8" class="py-8 text-center text-slate-400">
                <RotateCw class="w-5 h-5 animate-spin mx-auto mb-2 text-blue-600" />
                Memuat data Purchase Order...
              </td>
            </tr>
            <tr v-else-if="filteredPOList.length === 0">
              <td colspan="8" class="py-12 text-center text-slate-400">
                Tidak ada Purchase Order yang cocok.
              </td>
            </tr>
            <tr v-for="po in filteredPOList" :key="po.id" class="hover:bg-slate-50/70 transition-colors">
              <td class="py-3 px-4 font-mono font-bold text-blue-600">{{ po.po_number }}</td>
              <td class="py-3 px-4 font-mono text-[11px] text-slate-500">{{ po.pr_number || '-' }}</td>
              <td class="py-3 px-4 font-semibold text-slate-900">{{ po.supplier?.name || po.supplier_name || 'Vendor Umum' }}</td>
              <td class="py-3 px-4">{{ formatDate(po.order_date || po.created_at) }}</td>
              <td class="py-3 px-4">{{ formatDate(po.expected_date) }}</td>
              <td class="py-3 px-4 font-bold text-slate-900">Rp {{ formatNum(po.total_amount) }}</td>
              <td class="py-3 px-4">
                <span :class="getStatusBadgeClass(po.status)" class="px-2.5 py-0.5 rounded-full text-[10px] font-bold">
                  {{ po.status }}
                </span>
              </td>
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <button
                    @click="viewPODetail(po)"
                    class="p-1.5 rounded-lg text-slate-500 hover:bg-slate-100 cursor-pointer"
                    title="Lihat Detail PO"
                  >
                    <Eye class="w-4 h-4" />
                  </button>
                  <button
                    v-if="po.status === 'draft' || po.status === 'submitted'"
                    @click="approvePO(po)"
                    class="p-1.5 rounded-lg text-blue-600 hover:bg-blue-50 cursor-pointer"
                    title="Setujui PO (Manager)"
                  >
                    <Check class="w-4 h-4" />
                  </button>
                  <button
                    v-if="po.status === 'approved' || po.status === 'manager_approved' || po.status === 'sent'"
                    @click="openGRNForPO(po)"
                    class="px-2.5 py-1 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-[11px] shadow-xs cursor-pointer flex items-center gap-1"
                  >
                    <span>Terima (GRN)</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- TAB 3: GOODS RECEIPT NOTE (GRN)                                          -->
    <!-- ========================================================================= -->
    <div v-else-if="activeTab === 'grn'" class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 border-b border-slate-200/80 text-[11px] uppercase font-bold text-slate-500 tracking-wider">
            <tr>
              <th class="py-3 px-4">Nomor GRN</th>
              <th class="py-3 px-4">Ref. Nomor PO</th>
              <th class="py-3 px-4">No. Surat Jalan Vendor</th>
              <th class="py-3 px-4">Tgl Penerimaan</th>
              <th class="py-3 px-4">QC Status</th>
              <th class="py-3 px-4">Jurnal Akuntansi</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="7" class="py-8 text-center text-slate-400">
                <RotateCw class="w-5 h-5 animate-spin mx-auto mb-2 text-blue-600" />
                Memuat data GRN...
              </td>
            </tr>
            <tr v-else-if="filteredGRNList.length === 0">
              <td colspan="7" class="py-12 text-center text-slate-400">
                Belum ada catatan penerimaan barang gudang (GRN).
              </td>
            </tr>
            <tr v-for="grn in filteredGRNList" :key="grn.id" class="hover:bg-slate-50/70 transition-colors">
              <td class="py-3 px-4 font-mono font-bold text-emerald-600">{{ grn.grn_number }}</td>
              <td class="py-3 px-4 font-mono text-[11px] text-blue-600 font-bold">{{ grn.po_number || '-' }}</td>
              <td class="py-3 px-4 font-mono text-slate-700">{{ grn.delivery_note_number || '-' }}</td>
              <td class="py-3 px-4">{{ formatDate(grn.received_date || grn.created_at) }}</td>
              <td class="py-3 px-4">
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                  :class="grn.qc_status === 'passed' ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'"
                >
                  {{ grn.qc_status }}
                </span>
              </td>
              <td class="py-3 px-4">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-mono font-bold bg-blue-50 text-blue-700 border border-blue-200/60">
                  Auto Accrual Posted
                </span>
              </td>
              <td class="py-3 px-4 text-right">
                <button
                  @click="openInvoiceForPO(grn)"
                  class="px-2.5 py-1 rounded-lg bg-amber-600 hover:bg-amber-700 text-white font-bold text-[11px] shadow-xs cursor-pointer inline-flex items-center gap-1"
                >
                  <span>Proses Faktur</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- TAB 4: INVOICE & FAKTUR PAJAK (PPN 11%)                                   -->
    <!-- ========================================================================= -->
    <div v-else-if="activeTab === 'invoice'" class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 border-b border-slate-200/80 text-[11px] uppercase font-bold text-slate-500 tracking-wider">
            <tr>
              <th class="py-3 px-4">No. Faktur Vendor</th>
              <th class="py-3 px-4">Ref. PO</th>
              <th class="py-3 font-mono px-4">No. Faktur Pajak</th>
              <th class="py-3 px-4">DPP Nilai</th>
              <th class="py-3 px-4">PPN 11%</th>
              <th class="py-3 px-4">Total Tagihan</th>
              <th class="py-3 px-4">Jatuh Tempo</th>
              <th class="py-3 px-4">Status</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="9" class="py-8 text-center text-slate-400">
                <RotateCw class="w-5 h-5 animate-spin mx-auto mb-2 text-blue-600" />
                Memuat data Faktur...
              </td>
            </tr>
            <tr v-else-if="filteredInvoiceList.length === 0">
              <td colspan="9" class="py-12 text-center text-slate-400">
                Belum ada tagihan faktur vendor.
              </td>
            </tr>
            <tr v-for="inv in filteredInvoiceList" :key="inv.id" class="hover:bg-slate-50/70 transition-colors">
              <td class="py-3 px-4 font-mono font-bold text-slate-900">{{ inv.invoice_number }}</td>
              <td class="py-3 px-4 font-mono text-[11px] text-blue-600 font-semibold">{{ inv.po_number || '-' }}</td>
              <td class="py-3 px-4 font-mono text-[11px] text-amber-700">{{ inv.tax_invoice_number || '-' }}</td>
              <td class="py-3 px-4">Rp {{ formatNum(inv.dpp_amount) }}</td>
              <td class="py-3 px-4 font-bold text-amber-700">Rp {{ formatNum(inv.tax_amount) }}</td>
              <td class="py-3 px-4 font-black text-slate-900">Rp {{ formatNum(inv.total_amount) }}</td>
              <td class="py-3 px-4">{{ formatDate(inv.due_date) }}</td>
              <td class="py-3 px-4">
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-bold"
                  :class="inv.status === 'paid' ? 'bg-emerald-100 text-emerald-800' : 'bg-amber-100 text-amber-800'"
                >
                  {{ inv.status }}
                </span>
              </td>
              <td class="py-3 px-4 text-right">
                <button
                  v-if="inv.status !== 'paid'"
                  @click="openPaymentForInvoice(inv)"
                  class="px-2.5 py-1 rounded-lg bg-violet-600 hover:bg-violet-700 text-white font-bold text-[11px] shadow-xs cursor-pointer inline-flex items-center gap-1"
                >
                  <CreditCard class="w-3 h-3" />
                  <span>Bayar</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- TAB 5: VENDOR PAYMENT SETTLEMENT (KAS/BANK)                              -->
    <!-- ========================================================================= -->
    <div v-else-if="activeTab === 'payment'" class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-600">
          <thead class="bg-slate-50 border-b border-slate-200/80 text-[11px] uppercase font-bold text-slate-500 tracking-wider">
            <tr>
              <th class="py-3 px-4">No. Bukti Bayar</th>
              <th class="py-3 px-4">No. Faktur Vendor</th>
              <th class="py-3 px-4">Tgl Pembayaran</th>
              <th class="py-3 px-4">Metode Bayar</th>
              <th class="py-3 px-4">Akun Sumber</th>
              <th class="py-3 px-4">Nominal Dilunasi</th>
              <th class="py-3 px-4">Jurnal Pengurang AP</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="7" class="py-8 text-center text-slate-400">
                <RotateCw class="w-5 h-5 animate-spin mx-auto mb-2 text-blue-600" />
                Memuat data Pelunasan...
              </td>
            </tr>
            <tr v-else-if="filteredPaymentList.length === 0">
              <td colspan="7" class="py-12 text-center text-slate-400">
                Belum ada catatan transaksi pelunasan kas / bank ke vendor.
              </td>
            </tr>
            <tr v-for="pay in filteredPaymentList" :key="pay.id" class="hover:bg-slate-50/70 transition-colors">
              <td class="py-3 px-4 font-mono font-bold text-violet-600">{{ pay.payment_number }}</td>
              <td class="py-3 px-4 font-mono font-bold text-slate-800">{{ pay.invoice_number || '-' }}</td>
              <td class="py-3 px-4">{{ formatDate(pay.payment_date || pay.created_at) }}</td>
              <td class="py-3 px-4 uppercase font-semibold text-slate-700">{{ pay.payment_method }}</td>
              <td class="py-3 px-4 font-mono text-[11px] text-slate-600">Kode: {{ pay.account_code || '11102' }}</td>
              <td class="py-3 px-4 font-black text-emerald-600">Rp {{ formatNum(pay.amount) }}</td>
              <td class="py-3 px-4">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-mono font-bold bg-emerald-50 text-emerald-700 border border-emerald-200/60">
                  AP Settled (Dr AP / Cr Kas)
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ========================================================================= -->
    <!-- MODAL 1: CREATE PR                                                       -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreatePRModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center p-4 bg-slate-950/45 backdrop-blur-sm">
        <div class="bg-white rounded-3xl max-w-xl w-full p-6 shadow-2xl border border-slate-100 max-h-[90vh] flex flex-col">
          <div class="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h3 class="text-lg font-black text-slate-900">Buat Purchase Requisition (PR)</h3>
              <p class="text-xs text-slate-500">Tahap 1: Pengajuan pengadaan internal cafe</p>
            </div>
            <button @click="showCreatePRModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <form @submit.prevent="submitPR" class="space-y-4 pt-4 overflow-y-auto flex-1">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Departemen Pengaju</label>
                <input
                  v-model="newPR.department"
                  type="text"
                  required
                  placeholder="e.g. Barista & Dapur Utama"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Tanggal Dibutuhkan</label>
                <input
                  v-model="newPR.required_date"
                  type="date"
                  required
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Catatan Keperluan</label>
              <textarea
                v-model="newPR.notes"
                rows="2"
                placeholder="Alasan pengadaan, estimasi event, atau kebutuhan stok tipis"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden resize-none"
              ></textarea>
            </div>

            <!-- Items -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-bold text-slate-700">Item Bahan Baku yang Diminta</span>
                <button
                  type="button"
                  @click="addPRItem"
                  class="text-[11px] font-bold text-blue-600 hover:underline flex items-center gap-1 cursor-pointer"
                >
                  <Plus class="w-3.5 h-3.5" /> Tambah Baris
                </button>
              </div>

              <div class="space-y-2 max-h-48 overflow-y-auto pr-1">
                <div v-for="(it, idx) in newPR.items" :key="idx" class="flex items-center gap-2 p-2 bg-slate-50 rounded-xl border border-slate-200/60">
                  <select
                    v-model="it.inventory_item_id"
                    required
                    class="flex-1 px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden"
                  >
                    <option value="" disabled>Pilih Bahan Baku / Stok</option>
                    <option v-for="stock in availableStocks" :key="stock.id" :value="stock.id">
                      {{ stock.name }} (Stok saat ini: {{ stock.quantity || 0 }} {{ stock.unit || 'unit' }})
                    </option>
                  </select>
                  <input
                    v-model.number="it.quantity"
                    type="number"
                    min="1"
                    required
                    placeholder="Qty"
                    class="w-20 px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden text-right"
                  />
                  <input
                    v-model.number="it.estimated_price"
                    type="number"
                    min="0"
                    placeholder="Est. Harga"
                    class="w-28 px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden text-right"
                  />
                  <button
                    type="button"
                    @click="removePRItem(idx)"
                    :disabled="newPR.items.length <= 1"
                    class="p-1 text-slate-400 hover:text-red-500 disabled:opacity-30 cursor-pointer"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
              <button
                type="button"
                @click="showCreatePRModal = false"
                class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-100 rounded-xl cursor-pointer"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer disabled:opacity-60"
              >
                {{ submitting ? 'Menyimpan...' : 'Kirim Permintaan (PR)' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
    </Teleport>

    <!-- ========================================================================= -->
    <!-- MODAL 2: CONVERT PR TO PO                                                -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showConvertPRModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center p-4 bg-slate-950/45 backdrop-blur-sm">
        <div class="bg-white rounded-3xl max-w-lg w-full p-6 shadow-2xl border border-slate-100">
          <div class="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h3 class="text-lg font-black text-slate-900">Konversi PR ke Purchase Order (PO)</h3>
              <p class="text-xs text-slate-500 font-mono">Ref PR: {{ selectedPRToConvert?.pr_number }}</p>
            </div>
            <button @click="showConvertPRModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <form @submit.prevent="submitConvertPR" class="space-y-4 pt-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Pilih Supplier Resmi Vendor</label>
              <select
                v-model="convertPOForm.supplier_id"
                required
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden"
              >
                <option value="" disabled>Pilih Supplier</option>
                <option v-for="s in supplierOptions" :key="s.id" :value="s.id">
                  {{ s.name }} ({{ s.contact_person || 'Kontak' }} - {{ s.phone || '-' }})
                </option>
              </select>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Syarat Pembayaran</label>
                <select
                  v-model="convertPOForm.payment_terms"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                >
                  <option value="cash">Tunai (Cash On Delivery)</option>
                  <option value="net_14">Net 14 Hari</option>
                  <option value="net_30">Net 30 Hari</option>
                  <option value="net_60">Net 60 Hari</option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Estimasi Pengiriman</label>
                <input
                  v-model="convertPOForm.expected_date"
                  type="date"
                  required
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Catatan PO</label>
              <textarea
                v-model="convertPOForm.notes"
                rows="2"
                placeholder="Instruksi pengiriman, lokasi bongkar muat..."
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-hidden resize-none"
              ></textarea>
            </div>

            <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
              <button
                type="button"
                @click="showConvertPRModal = false"
                class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-100 rounded-xl cursor-pointer"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer disabled:opacity-60"
              >
                {{ submitting ? 'Memproses...' : 'Terbitkan PO Resmi' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
    </Teleport>

    <!-- ========================================================================= -->
    <!-- MODAL 3: CREATE DIRECT PO                                                -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreatePOModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center p-4 bg-slate-950/45 backdrop-blur-sm">
        <div class="bg-white rounded-3xl max-w-xl w-full p-6 shadow-2xl border border-slate-100 max-h-[90vh] flex flex-col">
          <div class="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h3 class="text-lg font-black text-slate-900">Buat Purchase Order (PO) Langsung</h3>
              <p class="text-xs text-slate-500">Tahap 2: Pemesanan resmi ke supplier</p>
            </div>
            <button @click="showCreatePOModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <form @submit.prevent="submitDirectPO" class="space-y-4 pt-4 overflow-y-auto flex-1">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Supplier / Vendor</label>
                <select
                  v-model="newPO.supplier_id"
                  required
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden"
                >
                  <option value="" disabled>Pilih Supplier</option>
                  <option v-for="s in supplierOptions" :key="s.id" :value="s.id">
                    {{ s.name }}
                  </option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Estimasi Pengiriman</label>
                <input
                  v-model="newPO.expected_date"
                  type="date"
                  required
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Catatan</label>
              <textarea
                v-model="newPO.notes"
                rows="2"
                placeholder="Instruksi tambahan untuk vendor..."
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl outline-hidden resize-none"
              ></textarea>
            </div>

            <!-- Items -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-bold text-slate-700">Daftar Bahan Baku Dipesan</span>
                <button
                  type="button"
                  @click="addPOItem"
                  class="text-[11px] font-bold text-blue-600 hover:underline flex items-center gap-1 cursor-pointer"
                >
                  <Plus class="w-3.5 h-3.5" /> Tambah Baris
                </button>
              </div>

              <div class="space-y-2 max-h-44 overflow-y-auto pr-1">
                <div v-for="(it, idx) in newPO.items" :key="idx" class="flex items-center gap-2 p-2 bg-slate-50 rounded-xl border border-slate-200/60">
                  <select
                    v-model="it.inventory_item_id"
                    required
                    class="flex-1 px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden"
                  >
                    <option value="" disabled>Pilih Bahan Baku</option>
                    <option v-for="stock in availableStocks" :key="stock.id" :value="stock.id">
                      {{ stock.name }} ({{ stock.unit || 'unit' }})
                    </option>
                  </select>
                  <input
                    v-model.number="it.quantity"
                    type="number"
                    min="1"
                    required
                    placeholder="Qty"
                    class="w-20 px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden text-right"
                  />
                  <input
                    v-model.number="it.unit_price"
                    type="number"
                    min="0"
                    placeholder="Harga Satuan"
                    class="w-28 px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden text-right"
                  />
                  <button
                    type="button"
                    @click="removePOItem(idx)"
                    :disabled="newPO.items.length <= 1"
                    class="p-1 text-slate-400 hover:text-red-500 disabled:opacity-30 cursor-pointer"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
              <button
                type="button"
                @click="showCreatePOModal = false"
                class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-100 rounded-xl cursor-pointer"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold text-xs rounded-xl shadow-md shadow-blue-600/30 transition-all cursor-pointer disabled:opacity-60"
              >
                {{ submitting ? 'Menyimpan...' : 'Terbitkan PO' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
    </Teleport>

    <!-- ========================================================================= -->
    <!-- MODAL 4: CREATE GRN                                                      -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreateGRNModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center p-4 bg-slate-950/45 backdrop-blur-sm">
        <div class="bg-white rounded-3xl max-w-xl w-full p-6 shadow-2xl border border-slate-100 max-h-[90vh] flex flex-col">
          <div class="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h3 class="text-lg font-black text-slate-900">Catat Penerimaan Barang Gudang (GRN)</h3>
              <p class="text-xs text-slate-500">Tahap 3: Stok bertambah & jurnal akrual persediaan tercatat</p>
            </div>
            <button @click="showCreateGRNModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <form @submit.prevent="submitGRN" class="space-y-4 pt-4 overflow-y-auto flex-1">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Purchase Order Terkait</label>
                <select
                  v-model="newGRN.po_id"
                  required
                  @change="onGRNPODetailsChange"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden"
                >
                  <option value="" disabled>Pilih PO</option>
                  <option v-for="po in poList" :key="po.id" :value="po.id">
                    {{ po.po_number }} - {{ po.supplier?.name || po.supplier_name || 'Vendor' }}
                  </option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">No. Surat Jalan Vendor</label>
                <input
                  v-model="newGRN.delivery_note_number"
                  type="text"
                  required
                  placeholder="e.g. SJ-2026/09/001"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Status QC (Pemeriksaan Kualitas)</label>
                <select
                  v-model="newGRN.qc_status"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                >
                  <option value="passed">Lulus QC Penuh (Passed)</option>
                  <option value="partial">Sebagian Lolos (Partial)</option>
                  <option value="rejected">Ditolak / Rusak (Rejected)</option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Catatan Penerimaan</label>
                <input
                  v-model="newGRN.notes"
                  type="text"
                  placeholder="Kondisi kemasan, suhu bahan, dll."
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                />
              </div>
            </div>

            <!-- Items to Receive -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-bold text-slate-700">Rincian Fisik Barang Diterima</span>
                <button
                  type="button"
                  @click="addGRNItem"
                  class="text-[11px] font-bold text-emerald-600 hover:underline flex items-center gap-1 cursor-pointer"
                >
                  <Plus class="w-3.5 h-3.5" /> Tambah Baris
                </button>
              </div>

              <div class="space-y-2 max-h-44 overflow-y-auto pr-1">
                <div v-for="(it, idx) in newGRN.items" :key="idx" class="flex items-center gap-2 p-2 bg-slate-50 rounded-xl border border-slate-200/60">
                  <select
                    v-model="it.inventory_item_id"
                    required
                    class="flex-1 px-2.5 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden"
                  >
                    <option value="" disabled>Pilih Bahan Baku</option>
                    <option v-for="stock in availableStocks" :key="stock.id" :value="stock.id">
                      {{ stock.name }} ({{ stock.unit || 'unit' }})
                    </option>
                  </select>
                  <input
                    v-model.number="it.quantity_received"
                    type="number"
                    min="1"
                    required
                    placeholder="Qty Terima"
                    title="Jumlah Diterima"
                    class="w-24 px-2 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden text-right"
                  />
                  <input
                    v-model.number="it.quantity_rejected"
                    type="number"
                    min="0"
                    placeholder="Qty Tolak"
                    title="Jumlah Ditolak"
                    class="w-20 px-2 py-1.5 text-xs bg-white border border-slate-200 rounded-lg outline-hidden text-right"
                  />
                  <button
                    type="button"
                    @click="removeGRNItem(idx)"
                    :disabled="newGRN.items.length <= 1"
                    class="p-1 text-slate-400 hover:text-red-500 disabled:opacity-30 cursor-pointer"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
              <button
                type="button"
                @click="showCreateGRNModal = false"
                class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-100 rounded-xl cursor-pointer"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-xl shadow-md shadow-emerald-600/30 transition-all cursor-pointer disabled:opacity-60"
              >
                {{ submitting ? 'Memproses...' : 'Simpan & Tambah Stok Gudang' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
    </Teleport>

    <!-- ========================================================================= -->
    <!-- MODAL 5: CREATE INVOICE (FAKTUR & PPN 11%)                                -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreateInvoiceModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center p-4 bg-slate-950/45 backdrop-blur-sm">
        <div class="bg-white rounded-3xl max-w-lg w-full p-6 shadow-2xl border border-slate-100">
          <div class="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h3 class="text-lg font-black text-slate-900">Input Faktur Pembelian & PPN 11%</h3>
              <p class="text-xs text-slate-500">Tahap 4: Tagihan Vendor & Pengakuan Pajak Masukan</p>
            </div>
            <button @click="showCreateInvoiceModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <form @submit.prevent="submitInvoice" class="space-y-4 pt-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Purchase Order Terkait</label>
              <select
                v-model="newInvoice.po_id"
                required
                @change="onInvoicePOChange"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden"
              >
                <option value="" disabled>Pilih PO</option>
                <option v-for="po in poList" :key="po.id" :value="po.id">
                  {{ po.po_number }} - Rp {{ formatNum(po.total_amount) }}
                </option>
              </select>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">No. Invoice Tagihan Vendor</label>
                <input
                  v-model="newInvoice.invoice_number"
                  type="text"
                  required
                  placeholder="e.g. INV-SUP/2026/099"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">No. Seri Faktur Pajak (PPN 11%)</label>
                <input
                  v-model="newInvoice.tax_invoice_number"
                  type="text"
                  placeholder="e.g. 010.000-26.12345678"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Dasar Pengenaan Pajak (DPP)</label>
                <input
                  v-model.number="newInvoice.dpp_amount"
                  type="number"
                  min="0"
                  required
                  @input="recalcTax"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden text-right font-bold"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">PPN Masukan (11%)</label>
                <input
                  v-model.number="newInvoice.tax_amount"
                  type="number"
                  min="0"
                  class="w-full px-3 py-2 text-xs bg-slate-100 border border-slate-200 rounded-xl text-amber-700 outline-hidden text-right font-bold"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Tanggal Jatuh Tempo</label>
                <input
                  v-model="newInvoice.due_date"
                  type="date"
                  required
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Total Tagihan (DPP + PPN)</label>
                <div class="px-3 py-2 text-xs bg-blue-50 border border-blue-200 rounded-xl font-black text-blue-900 text-right">
                  Rp {{ formatNum(newInvoice.dpp_amount + newInvoice.tax_amount) }}
                </div>
              </div>
            </div>

            <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
              <button
                type="button"
                @click="showCreateInvoiceModal = false"
                class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-100 rounded-xl cursor-pointer"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-5 py-2 bg-amber-600 hover:bg-amber-700 text-white font-bold text-xs rounded-xl shadow-md shadow-amber-600/30 transition-all cursor-pointer disabled:opacity-60"
              >
                {{ submitting ? 'Memproses...' : 'Catat Faktur & Akun Hutang (AP)' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
    </Teleport>

    <!-- ========================================================================= -->
    <!-- MODAL 6: CREATE PAYMENT                                                  -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreatePaymentModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center p-4 bg-slate-950/45 backdrop-blur-sm">
        <div class="bg-white rounded-3xl max-w-lg w-full p-6 shadow-2xl border border-slate-100">
          <div class="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h3 class="text-lg font-black text-slate-900">Pelunasan Tagihan Kas/Bank Vendor</h3>
              <p class="text-xs text-slate-500">Tahap 5: Settlement Hutang Dagang (AP)</p>
            </div>
            <button @click="showCreatePaymentModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <form @submit.prevent="submitPayment" class="space-y-4 pt-4">
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Pilih Faktur yang Belum Dilunasi</label>
              <select
                v-model="newPayment.invoice_id"
                required
                @change="onPaymentInvoiceChange"
                class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-blue-600 outline-hidden"
              >
                <option value="" disabled>Pilih Tagihan Faktur</option>
                <option v-for="inv in unpaidInvoices" :key="inv.id" :value="inv.id">
                  {{ inv.invoice_number }} - Rp {{ formatNum(inv.total_amount) }} (Due: {{ formatDate(inv.due_date) }})
                </option>
              </select>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Metode Bayar</label>
                <select
                  v-model="newPayment.payment_method"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                >
                  <option value="bank_transfer">Transfer Bank BCA/Mandiri</option>
                  <option value="cash">Kas Operasional (Petty Cash)</option>
                  <option value="giro">Giro / Cek Bank</option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Akun Sumber (COA)</label>
                <select
                  v-model="newPayment.account_code"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden font-mono"
                >
                  <option value="11102">11102 - Bank BCA Operasional</option>
                  <option value="11101">11101 - Kas Kasir / Peti Kas</option>
                  <option value="11103">11103 - Bank Mandiri Bisnis</option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">Nominal Pelunasan (Rp)</label>
                <input
                  v-model.number="newPayment.amount"
                  type="number"
                  min="1"
                  required
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden text-right font-black text-slate-900"
                />
              </div>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">No. Referensi / Ref Bank</label>
                <input
                  v-model="newPayment.reference_number"
                  type="text"
                  placeholder="e.g. TRF-BCA-982312"
                  class="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-200 rounded-xl focus:bg-white outline-hidden"
                />
              </div>
            </div>

            <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
              <button
                type="button"
                @click="showCreatePaymentModal = false"
                class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-100 rounded-xl cursor-pointer"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-5 py-2 bg-violet-600 hover:bg-violet-700 text-white font-bold text-xs rounded-xl shadow-md shadow-violet-600/30 transition-all cursor-pointer disabled:opacity-60"
              >
                {{ submitting ? 'Memproses...' : 'Posting Pelunasan (Settled)' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
    </Teleport>

    <!-- ========================================================================= -->
    <!-- MODAL 7: VIEW PO DETAIL                                                  -->
    <!-- ========================================================================= -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showDetailModal" class="fixed md:left-64 inset-y-0 right-0 left-0 z-50 flex items-center justify-center p-4 bg-slate-950/45 backdrop-blur-sm">
        <div class="bg-white rounded-3xl max-w-lg w-full p-6 shadow-2xl border border-slate-100">
          <div class="flex items-center justify-between pb-4 border-b border-slate-100">
            <div>
              <h3 class="text-lg font-black text-slate-900">Detail Purchase Order</h3>
              <p class="text-xs text-blue-600 font-mono font-bold">{{ selectedPODetail?.po_number }}</p>
            </div>
            <button @click="showDetailModal = false" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer">
              <X class="w-5 h-5" />
            </button>
          </div>

          <div class="py-4 space-y-4">
            <div class="grid grid-cols-2 gap-3 text-xs">
              <div class="bg-slate-50 p-3 rounded-xl border border-slate-200/60">
                <span class="text-slate-400 font-semibold block mb-0.5">Supplier:</span>
                <span class="font-bold text-slate-900">{{ selectedPODetail?.supplier?.name || selectedPODetail?.supplier_name || 'Vendor Umum' }}</span>
              </div>
              <div class="bg-slate-50 p-3 rounded-xl border border-slate-200/60">
                <span class="text-slate-400 font-semibold block mb-0.5">Status:</span>
                <span :class="getStatusBadgeClass(selectedPODetail?.status)" class="px-2 py-0.5 rounded-full text-[10px] font-bold">
                  {{ selectedPODetail?.status }}
                </span>
              </div>
              <div class="bg-slate-50 p-3 rounded-xl border border-slate-200/60">
                <span class="text-slate-400 font-semibold block mb-0.5">Tgl Pesan:</span>
                <span class="font-medium text-slate-800">{{ formatDate(selectedPODetail?.order_date || selectedPODetail?.created_at) }}</span>
              </div>
              <div class="bg-slate-50 p-3 rounded-xl border border-slate-200/60">
                <span class="text-slate-400 font-semibold block mb-0.5">Tgl Ekspektasi:</span>
                <span class="font-medium text-slate-800">{{ formatDate(selectedPODetail?.expected_date) }}</span>
              </div>
            </div>

            <!-- Items Table -->
            <div>
              <h4 class="text-xs font-bold text-slate-700 mb-2">Item Barang Dipesan:</h4>
              <div class="bg-slate-50 rounded-xl border border-slate-200/60 overflow-hidden">
                <table class="w-full text-left text-xs">
                  <thead class="bg-slate-100 text-[10px] uppercase font-bold text-slate-500">
                    <tr>
                      <th class="p-2">Item</th>
                      <th class="p-2 text-right">Qty</th>
                      <th class="p-2 text-right">Harga</th>
                      <th class="p-2 text-right">Subtotal</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-200/60">
                    <tr v-for="(it, i) in (selectedPODetail?.items || [])" :key="i">
                      <td class="p-2 font-medium text-slate-800">{{ it.inventory_item_name || it.item_name || 'Bahan Baku' }}</td>
                      <td class="p-2 text-right font-mono">{{ it.quantity }}</td>
                      <td class="p-2 text-right font-mono">Rp {{ formatNum(it.unit_price) }}</td>
                      <td class="p-2 text-right font-mono font-bold text-slate-900">Rp {{ formatNum(it.quantity * it.unit_price) }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <div class="flex items-center justify-between p-3 rounded-xl bg-blue-50 border border-blue-200 text-xs">
              <span class="font-bold text-blue-900">Total Nilai Pemesanan:</span>
              <span class="text-sm font-black text-blue-900">Rp {{ formatNum(selectedPODetail?.total_amount) }}</span>
            </div>
          </div>

          <div class="pt-3 border-t border-slate-100 flex justify-end">
            <button
              @click="showDetailModal = false"
              class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold text-xs rounded-xl cursor-pointer"
            >
              Tutup
            </button>
          </div>
        </div>
      </div>
    </Transition>
    </Teleport>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import {
  RotateCw,
  Plus,
  Search,
  Eye,
  Check,
  CheckCircle2,
  PackageCheck,
  FileText,
  CreditCard,
  Trash2,
  X,
  ArrowRight
} from 'lucide-vue-next'
import { useNotificationStore } from '@/stores/notification.store'

type TabType = 'pr' | 'po' | 'grn' | 'invoice' | 'payment'

const notifyStore = useNotificationStore()

const activeTab = ref<TabType>('po')
const loading = ref(false)
const submitting = ref(false)
const searchQuery = ref('')
const activeStatusFilter = ref('all')

const p2pSteps: { id: TabType; code: string; label: string }[] = [
  { id: 'pr', code: 'Tahap 1', label: 'Permintaan (PR)' },
  { id: 'po', code: 'Tahap 2', label: 'Pemesanan (PO)' },
  { id: 'grn', code: 'Tahap 3', label: 'Penerimaan (GRN)' },
  { id: 'invoice', code: 'Tahap 4', label: 'Faktur & PPN 11%' },
  { id: 'payment', code: 'Tahap 5', label: 'Pelunasan Kas/Bank' }
]

// Data state
const prList = ref<any[]>([])
const poList = ref<any[]>([])
const grnList = ref<any[]>([])
const invoiceList = ref<any[]>([])
const paymentList = ref<any[]>([])
const supplierOptions = ref<any[]>([])
const availableStocks = ref<any[]>([])

// Modals
const showCreatePRModal = ref(false)
const showConvertPRModal = ref(false)
const showCreatePOModal = ref(false)
const showCreateGRNModal = ref(false)
const showCreateInvoiceModal = ref(false)
const showCreatePaymentModal = ref(false)
const showDetailModal = ref(false)

const selectedPRToConvert = ref<any>(null)
const selectedPODetail = ref<any>(null)

// Forms
const newPR = ref({
  department: 'Dapur & Bar',
  required_date: new Date(Date.now() + 3 * 86400000).toISOString().split('T')[0],
  notes: '',
  items: [{ inventory_item_id: '', quantity: 10, estimated_price: 25000 }]
})

const convertPOForm = ref({
  supplier_id: '',
  payment_terms: 'net_30',
  expected_date: new Date(Date.now() + 5 * 86400000).toISOString().split('T')[0],
  notes: ''
})

const newPO = ref({
  supplier_id: '',
  expected_date: new Date(Date.now() + 7 * 86400000).toISOString().split('T')[0],
  notes: '',
  items: [{ inventory_item_id: '', quantity: 10, unit_price: 25000 }]
})

const newGRN = ref({
  po_id: '',
  delivery_note_number: '',
  qc_status: 'passed',
  notes: '',
  items: [{ inventory_item_id: '', quantity_received: 10, quantity_rejected: 0, unit_price: 25000 }]
})

const newInvoice = ref({
  po_id: '',
  invoice_number: '',
  tax_invoice_number: '',
  dpp_amount: 0,
  tax_amount: 0,
  due_date: new Date(Date.now() + 30 * 86400000).toISOString().split('T')[0]
})

const newPayment = ref({
  invoice_id: '',
  account_code: '11102',
  payment_method: 'bank_transfer',
  amount: 0,
  reference_number: ''
})

// Formatters
const formatNum = (num: number) => {
  return Number(num || 0).toLocaleString('id-ID')
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Intl.DateTimeFormat('id-ID', { dateStyle: 'medium' }).format(new Date(dateStr))
  } catch {
    return dateStr
  }
}

const getStatusBadgeClass = (status: string) => {
  switch (status?.toLowerCase()) {
    case 'approved':
    case 'passed':
    case 'paid':
      return 'bg-emerald-100 text-emerald-800'
    case 'sent':
    case 'manager_approved':
      return 'bg-blue-100 text-blue-800'
    case 'draft':
    case 'submitted':
    case 'pending_approval':
    case 'unpaid':
      return 'bg-amber-100 text-amber-800'
    case 'rejected':
    case 'cancelled':
      return 'bg-rose-100 text-rose-800'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

// Fetch Data
const fetchPRList = async () => {
  try {
    const res = await axios.get('/api/v1/inventory/pr')
    prList.value = res.data?.data || []
  } catch (err) {
    console.error('Failed to load PR:', err)
  }
}

const fetchPOList = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/v1/inventory/purchase-orders')
    poList.value = Array.isArray(res.data) ? res.data : (res.data?.data || [])
  } catch (err) {
    console.error('Failed to load PO:', err)
  } finally {
    loading.value = false
  }
}

const fetchGRNList = async () => {
  try {
    const res = await axios.get('/api/v1/inventory/grn')
    grnList.value = res.data?.data || []
  } catch (err) {
    console.error('Failed to load GRN:', err)
  }
}

const fetchInvoiceList = async () => {
  try {
    const res = await axios.get('/api/v1/inventory/invoices')
    invoiceList.value = res.data?.data || []
  } catch (err) {
    console.error('Failed to load Invoices:', err)
  }
}

const fetchPaymentList = async () => {
  try {
    const res = await axios.get('/api/v1/inventory/vendor-payments')
    paymentList.value = res.data?.data || []
  } catch (err) {
    console.error('Failed to load Payments:', err)
  }
}

const fetchMasterData = async () => {
  try {
    const [supRes, stockRes] = await Promise.all([
      axios.get('/api/v1/master/suppliers'),
      axios.get('/api/v1/inventory/stocks')
    ])
    supplierOptions.value = Array.isArray(supRes.data) ? supRes.data : (supRes.data?.data || [])
    availableStocks.value = Array.isArray(stockRes.data) ? stockRes.data : (stockRes.data?.data || [])
  } catch (err) {
    console.error('Failed to load suppliers/stocks:', err)
  }
}

const refreshCurrentTab = async () => {
  loading.value = true
  await Promise.all([
    fetchPRList(),
    fetchPOList(),
    fetchGRNList(),
    fetchInvoiceList(),
    fetchPaymentList()
  ])
  loading.value = false
}

// Filters & Computeds
const currentFilterOptions = computed(() => {
  return [
    { label: 'Semua Status', value: 'all' },
    { label: 'Draft', value: 'draft' },
    { label: 'Disetujui', value: 'approved' },
    { label: 'Selesai / Lunas', value: 'paid' }
  ]
})

const unpaidInvoices = computed(() => {
  return invoiceList.value.filter(inv => inv.status !== 'paid')
})

const filteredPRList = computed(() => {
  return prList.value.filter(pr => {
    const matchStatus = activeStatusFilter.value === 'all' || pr.status === activeStatusFilter.value
    const q = searchQuery.value.toLowerCase()
    const matchQuery = !q || (pr.pr_number && pr.pr_number.toLowerCase().includes(q)) || (pr.department && pr.department.toLowerCase().includes(q))
    return matchStatus && matchQuery
  })
})

const filteredPOList = computed(() => {
  return poList.value.filter(po => {
    const matchStatus = activeStatusFilter.value === 'all' || po.status === activeStatusFilter.value
    const q = searchQuery.value.toLowerCase()
    const sup = po.supplier?.name || po.supplier_name || ''
    const matchQuery = !q || (po.po_number && po.po_number.toLowerCase().includes(q)) || sup.toLowerCase().includes(q)
    return matchStatus && matchQuery
  })
})

const filteredGRNList = computed(() => {
  return grnList.value.filter(grn => {
    const q = searchQuery.value.toLowerCase()
    return !q || (grn.grn_number && grn.grn_number.toLowerCase().includes(q)) || (grn.po_number && grn.po_number.toLowerCase().includes(q))
  })
})

const filteredInvoiceList = computed(() => {
  return invoiceList.value.filter(inv => {
    const q = searchQuery.value.toLowerCase()
    return !q || (inv.invoice_number && inv.invoice_number.toLowerCase().includes(q)) || (inv.po_number && inv.po_number.toLowerCase().includes(q))
  })
})

const filteredPaymentList = computed(() => {
  return paymentList.value.filter(pay => {
    const q = searchQuery.value.toLowerCase()
    return !q || (pay.payment_number && pay.payment_number.toLowerCase().includes(q)) || (pay.invoice_number && pay.invoice_number.toLowerCase().includes(q))
  })
})

// Modal Openers
const openCreatePRModal = () => {
  newPR.value = {
    department: 'Dapur & Bar',
    required_date: new Date(Date.now() + 3 * 86400000).toISOString().split('T')[0],
    notes: '',
    items: [{ inventory_item_id: availableStocks.value[0]?.id || '', quantity: 10, estimated_price: 25000 }]
  }
  showCreatePRModal.value = true
}

const addPRItem = () => {
  newPR.value.items.push({
    inventory_item_id: availableStocks.value[0]?.id || '',
    quantity: 10,
    estimated_price: 25000
  })
}

const removePRItem = (idx: number) => {
  if (newPR.value.items.length > 1) {
    newPR.value.items.splice(idx, 1)
  }
}

const submitPR = async () => {
  if (newPR.value.items.some(i => !i.inventory_item_id || i.quantity <= 0)) {
    notifyStore.warning('Pilih item barang dan isi kuantitas valid!', 'Validasi')
    return
  }
  submitting.value = true
  try {
    const res = await axios.post('/api/v1/inventory/pr', newPR.value)
    notifyStore.success(res.data?.message || 'PR berhasil diajukan!', 'Sukses')
    showCreatePRModal.value = false
    await fetchPRList()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal mengajukan PR', 'Error')
  } finally {
    submitting.value = false
  }
}

const approvePR = async (pr: any) => {
  try {
    const res = await axios.put(`/api/v1/inventory/pr/${pr.id}/status`, { status: 'approved' })
    notifyStore.success(res.data?.message || 'PR berhasil disetujui!', 'Sukses')
    await fetchPRList()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menyetujui PR', 'Error')
  }
}

const openConvertPRModal = (pr: any) => {
  selectedPRToConvert.value = pr
  convertPOForm.value = {
    supplier_id: supplierOptions.value[0]?.id || '',
    payment_terms: 'net_30',
    expected_date: new Date(Date.now() + 5 * 86400000).toISOString().split('T')[0],
    notes: `Dikonversi dari ${pr.pr_number}`
  }
  showConvertPRModal.value = true
}

const submitConvertPR = async () => {
  if (!selectedPRToConvert.value) return
  submitting.value = true
  try {
    const res = await axios.post(`/api/v1/inventory/pr/${selectedPRToConvert.value.id}/convert-po`, convertPOForm.value)
    notifyStore.success(res.data?.message || 'PR berhasil dikonversi ke PO!', 'Sukses')
    showConvertPRModal.value = false
    activeTab.value = 'po'
    await Promise.all([fetchPRList(), fetchPOList()])
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal konversi ke PO', 'Error')
  } finally {
    submitting.value = false
  }
}

const openCreatePOModal = () => {
  newPO.value = {
    supplier_id: supplierOptions.value[0]?.id || '',
    expected_date: new Date(Date.now() + 7 * 86400000).toISOString().split('T')[0],
    notes: '',
    items: [{ inventory_item_id: availableStocks.value[0]?.id || '', quantity: 10, unit_price: 25000 }]
  }
  showCreatePOModal.value = true
}

const addPOItem = () => {
  newPO.value.items.push({
    inventory_item_id: availableStocks.value[0]?.id || '',
    quantity: 10,
    unit_price: 25000
  })
}

const removePOItem = (idx: number) => {
  if (newPO.value.items.length > 1) {
    newPO.value.items.splice(idx, 1)
  }
}

const submitDirectPO = async () => {
  if (!newPO.value.supplier_id || newPO.value.items.some(i => !i.inventory_item_id || i.quantity <= 0)) {
    notifyStore.warning('Lengkapi supplier dan rincian item dengan benar!', 'Validasi')
    return
  }
  submitting.value = true
  try {
    const res = await axios.post('/api/v1/inventory/purchase-orders', newPO.value)
    notifyStore.success(res.data?.message || 'PO berhasil diterbitkan!', 'Sukses')
    showCreatePOModal.value = false
    await fetchPOList()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menerbitkan PO', 'Error')
  } finally {
    submitting.value = false
  }
}

const approvePO = async (po: any) => {
  try {
    const res = await axios.put(`/api/v1/inventory/purchase-orders/${po.id}/status`, { status: 'approved' })
    notifyStore.success(res.data?.message || 'PO berhasil disetujui!', 'Sukses')
    await fetchPOList()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menyetujui PO', 'Error')
  }
}

const viewPODetail = (po: any) => {
  selectedPODetail.value = po
  showDetailModal.value = true
}

const openGRNForPO = (po: any) => {
  newGRN.value.po_id = po.id
  newGRN.value.delivery_note_number = `SJ-${po.po_number || '001'}`
  if (po.items && po.items.length > 0) {
    newGRN.value.items = po.items.map((it: any) => ({
      inventory_item_id: it.inventory_item_id,
      quantity_received: it.quantity,
      quantity_rejected: 0,
      unit_price: it.unit_price
    }))
  } else {
    newGRN.value.items = [{
      inventory_item_id: availableStocks.value[0]?.id || '',
      quantity_received: 10,
      quantity_rejected: 0,
      unit_price: 25000
    }]
  }
  activeTab.value = 'grn'
  showCreateGRNModal.value = true
}

const openCreateGRNModal = () => {
  newGRN.value = {
    po_id: poList.value[0]?.id || '',
    delivery_note_number: '',
    qc_status: 'passed',
    notes: '',
    items: [{
      inventory_item_id: availableStocks.value[0]?.id || '',
      quantity_received: 10,
      quantity_rejected: 0,
      unit_price: 25000
    }]
  }
  showCreateGRNModal.value = true
}

const onGRNPODetailsChange = () => {
  const po = poList.value.find(p => p.id === newGRN.value.po_id)
  if (po && po.items && po.items.length > 0) {
    newGRN.value.items = po.items.map((it: any) => ({
      inventory_item_id: it.inventory_item_id,
      quantity_received: it.quantity,
      quantity_rejected: 0,
      unit_price: it.unit_price
    }))
  }
}

const addGRNItem = () => {
  newGRN.value.items.push({
    inventory_item_id: availableStocks.value[0]?.id || '',
    quantity_received: 10,
    quantity_rejected: 0,
    unit_price: 25000
  })
}

const removeGRNItem = (idx: number) => {
  if (newGRN.value.items.length > 1) {
    newGRN.value.items.splice(idx, 1)
  }
}

const submitGRN = async () => {
  if (!newGRN.value.po_id || newGRN.value.items.some(i => !i.inventory_item_id || i.quantity_received <= 0)) {
    notifyStore.warning('Pilih PO dan masukkan kuantitas barang diterima!', 'Validasi')
    return
  }
  submitting.value = true
  try {
    const res = await axios.post('/api/v1/inventory/grn', newGRN.value)
    notifyStore.success(res.data?.message || 'GRN berhasil dicatat & stok bertambah!', 'Sukses')
    showCreateGRNModal.value = false
    await Promise.all([fetchGRNList(), fetchPOList(), fetchMasterData()])
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan GRN', 'Error')
  } finally {
    submitting.value = false
  }
}

const openInvoiceForPO = (grn: any) => {
  newInvoice.value.po_id = grn.po_id
  const po = poList.value.find(p => p.id === grn.po_id)
  const total = Number(po?.total_amount || 1000000)
  newInvoice.value.dpp_amount = Math.round(total / 1.11)
  newInvoice.value.tax_amount = total - newInvoice.value.dpp_amount
  newInvoice.value.invoice_number = `INV-${po?.po_number || Date.now()}`
  activeTab.value = 'invoice'
  showCreateInvoiceModal.value = true
}

const openCreateInvoiceModal = () => {
  const firstPO = poList.value[0]
  const total = Number(firstPO?.total_amount || 1000000)
  const dpp = Math.round(total / 1.11)
  newInvoice.value = {
    po_id: firstPO?.id || '',
    invoice_number: `INV-${Date.now().toString().slice(-6)}`,
    tax_invoice_number: `010.000-26.${Date.now().toString().slice(-8)}`,
    dpp_amount: dpp,
    tax_amount: total - dpp,
    due_date: new Date(Date.now() + 30 * 86400000).toISOString().split('T')[0]
  }
  showCreateInvoiceModal.value = true
}

const onInvoicePOChange = () => {
  const po = poList.value.find(p => p.id === newInvoice.value.po_id)
  if (po) {
    const total = Number(po.total_amount || 0)
    newInvoice.value.dpp_amount = Math.round(total / 1.11)
    newInvoice.value.tax_amount = total - newInvoice.value.dpp_amount
    newInvoice.value.invoice_number = `INV-${po.po_number}`
  }
}

const recalcTax = () => {
  newInvoice.value.tax_amount = Math.round(newInvoice.value.dpp_amount * 0.11)
}

const submitInvoice = async () => {
  if (!newInvoice.value.po_id || !newInvoice.value.invoice_number) {
    notifyStore.warning('Lengkapi data nomor invoice dan PO terkait!', 'Validasi')
    return
  }
  submitting.value = true
  try {
    const res = await axios.post('/api/v1/inventory/invoices', newInvoice.value)
    notifyStore.success(res.data?.message || 'Faktur berhasil dicatat & jurnal hutang di-posting!', 'Sukses')
    showCreateInvoiceModal.value = false
    await fetchInvoiceList()
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal menyimpan faktur', 'Error')
  } finally {
    submitting.value = false
  }
}

const openPaymentForInvoice = (inv: any) => {
  newPayment.value.invoice_id = inv.id
  newPayment.value.amount = Number(inv.total_amount || 0)
  newPayment.value.reference_number = `TRF-BCA-${Date.now().toString().slice(-6)}`
  showCreatePaymentModal.value = true
}

const openCreatePaymentModal = () => {
  const inv = unpaidInvoices.value[0]
  newPayment.value = {
    invoice_id: inv?.id || '',
    account_code: '11102',
    payment_method: 'bank_transfer',
    amount: Number(inv?.total_amount || 0),
    reference_number: `TRF-BCA-${Date.now().toString().slice(-6)}`
  }
  showCreatePaymentModal.value = true
}

const onPaymentInvoiceChange = () => {
  const inv = invoiceList.value.find(i => i.id === newPayment.value.invoice_id)
  if (inv) {
    newPayment.value.amount = Number(inv.total_amount || 0)
  }
}

const submitPayment = async () => {
  if (!newPayment.value.invoice_id || newPayment.value.amount <= 0) {
    notifyStore.warning('Pilih invoice dan isi nominal pelunasan valid!', 'Validasi')
    return
  }
  submitting.value = true
  try {
    const res = await axios.post('/api/v1/inventory/vendor-payments', newPayment.value)
    notifyStore.success(res.data?.message || 'Pelunasan berhasil diposting & hutang diselesaikan!', 'Sukses')
    showCreatePaymentModal.value = false
    await Promise.all([fetchInvoiceList(), fetchPaymentList()])
  } catch (err: any) {
    notifyStore.error(err.response?.data?.message || 'Gagal posting pelunasan', 'Error')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  refreshCurrentTab()
  fetchMasterData()
})
</script>
