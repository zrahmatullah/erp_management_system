# Standarisasi UI & Pemetaan Mockup (Design System)

Sistem Cafe ERP telah mengadopsi seluruh rancangan antarmuka dari **20 file mockup** di folder `Mockups_UI/` dengan standarisasi visual modern berbasis **Tailwind CSS v3**, palet warna konsisten, serta komponen Vue 3 modular.

---

## 1. Design System Tokens & Standar Visual

- **Primary Color (Brand Accent)**: Royal Blue `#2563EB` (Tailwind `blue-600`) dengan variasi hover `#1D4ED8` (`blue-700`) dan soft tint `bg-blue-50`.
- **Sidebar & Shell Color**: Slate 900 `#0F172A` / `#1E293B` dengan aksen aktif `bg-blue-600` dan teks sekunder `text-slate-400`.
- **Neutral Surface & Background**: 
  - Body background: `#F8FAFC` (`slate-50`)
  - Card background: `#FFFFFF` (`white`) dengan border tipis `#E2E8F0` (`slate-200`) dan bayangan lembut (`shadow-xs` / `shadow-sm`).
- **Typography**: Inter / System Sans-serif, ukuran teks hierarkis:
  - Header: `text-2xl font-bold tracking-tight text-slate-900`
  - Subheader: `text-xl font-bold text-slate-900`
  - Body: `text-sm text-slate-700`
  - Label / Metadata: `text-xs font-semibold uppercase tracking-wider text-slate-400`
- **Border Radius**: `rounded-xl` (12px) untuk input & tombol, `rounded-2xl` (16px) untuk kartu kontainer utama.
- **Iconography**: `lucide-vue-next` dengan ketebalan stroke seragam.

---

## 2. Tabel Pemetaan 20 File Mockup ke Komponen Vue 3

| No | File Mockup Gambar di `Mockups_UI/` | File Komponen Vue 3 Implementasi | Route Path | Fitur Utama |
| :--- | :--- | :--- | :--- | :--- |
| 1 | `Auth/login_mockup_1788489124004.jpg` | `src/views/auth/LoginView.vue` | `/login` | Ilustrasi modern cafe, form kredensial, remember me, quick 1-click role switcher |
| 2 | `Dashboard/dashboard_mockup_1788488540267.jpg` | `src/views/dashboard/DashboardView.vue` | `/dashboard` | 4 KPI cards (Penjualan, Transaksi, Nilai Rata-rata, Harian), grafik penjualan, live orders |
| 3 | `POS/pos_mockup_1788488573190.jpg` | `src/views/pos/POSView.vue` | `/pos` | Grid katalog menu, filter kategori, panel keranjang kasir, perhitungan subtotal & pajak |
| 4 | `POS/payment_modal_mockup_1788508491875.jpg` | `src/views/pos/PaymentModal.vue` | Modal di `/pos` | Metode pembayaran (Cash, QRIS, EDC BCA/Mandiri), uang pas, kalkulasi kembalian |
| 5 | `POS/table_management_mockup_1788508012211.jpg` | `src/views/pos/TableManagementView.vue` | `/pos/tables` | Tab zona (Indoor AC, Outdoor Garden, Bar), kartu status meja (Available, Occupied, Reserved) |
| 6 | `POS/kitchen_display_mockup_1788507978701.jpg` | `src/views/pos/KitchenDisplayView.vue` | `/pos/kds` | Kolom tiket KDS (Queue, Cooking, Ready), timer durasi masak, checklist rincian item |
| 7 | `Inventory/inventory_mockup_1788488740443.jpg` | `src/views/inventory/StockList.vue` | `/inventory` | Tab All Items, Low Stock, Filter Kategori, status warning stok menipis, search |
| 8 | `Inventory/stock_opname_mockup_1788508726697.jpg` | `src/views/inventory/StockOpnameView.vue` | `/inventory/opname` | Form Stock Take fisik vs sistem, kalkulasi otomatis selisih unit & nilai rupiah |
| 9 | `Inventory/purchase_order_mockup_1788508351093.jpg` | `src/views/inventory/PurchaseOrderList.vue` | `/inventory/po` | Status PO (Approved, Pending, Completed), rincian supplier, nilai nominal, tombol aksi |
| 10 | `HRIS/hris_mockup_1788488905266.jpg` | `src/views/hris/EmployeeList.vue` | `/hris` | Direktori karyawan, pencarian NIK & nama, filter departemen, status aktif |
| 11 | `HRIS/employee_detail_mockup_1788508045796.jpg` | `src/views/hris/EmployeeDetailModal.vue` | Modal di `/hris` | Detail komprehensif profil staf, riwayat presensi, gaji, kontak darurat, rekening bank |
| 12 | `HRIS/attendance_mockup_1788508895427.jpg` | `src/views/hris/AttendanceView.vue` | `/hris/attendance` | Log presensi, jam clock-in/out, status On Time/Late, kalkulasi keterlambatan menit |
| 13 | `HRIS/shift_schedule_mockup_1788508534154.jpg` | `src/views/hris/ShiftScheduleView.vue` | `/hris/shifts` | Kalender jadwal shift mingguan, alokasi karyawan ke Morning/Middle/Closing shift |
| 14 | `HRIS/leave_management_mockup_1788508148211.jpg` | `src/views/hris/LeaveManagementView.vue` | `/hris/leaves` | Saldo cuti tahunan, tabel permohonan izin/cuti, tombol persetujuan Approve / Reject |
| 15 | `HRIS/payroll_mockup_1788508104353.jpg` | `src/views/hris/PayrollRun.vue` | `/hris/payroll` | Stepper 5 langkah proses payroll, tabel rekap gaji bruto, BPJS, PPh21, dan netto |
| 16 | `HRIS/payslip_mockup_1788508775586.jpg` | `src/views/hris/PayslipModal.vue` | Modal di `/hris/payroll` | Format resmi dokumen Slip Gaji, komparasi Pendapatan vs Potongan, tombol Print |
| 17 | `Finance/finance_mockup_1788488937126.jpg` | `src/views/finance/FinanceOverview.vue` | `/finance` | 4 kartu metrik keuangan, grafik batang Rev vs Exp, grafik tren kas, donut biaya |
| 18 | `Finance/journal_entry_mockup_1788508398257.jpg` | `src/views/finance/JournalEntryForm.vue` | `/finance/journal` | Form multi-baris jurnal umum, dropdown akun COA, validasi balance total Debit = Kredit |
| 19 | `Reports/reports_mockup_1788488969793.jpg` | `src/views/reports/SalesReport.vue` | `/reports` | Filter periode & cabang, kurva perbandingan periode lalu, Top 10 menu, heatmap jam sibuk |
| 20 | `Settings/role_permission_mockup_1788508435826.jpg` | `src/views/settings/RoleManagement.vue` | `/settings/roles` | Matrix hak akses permission (View, Create, Edit, Delete) per modul untuk setiap role |

---

## 3. Komponen Layout Global

- **`src/components/layout/AppLayout.vue`**: Shell induk yang membungkus seluruh halaman terproteksi.
- **`src/components/layout/Sidebar.vue`**: Sidebar navigasi kiri bergaya Slate Dark dengan header logo cafe, menu hierarki collapsible, badge notifikasi, profil user, dan akses Super Admin Master Hub.
- **`src/components/layout/AppHeader.vue`**: Bar bagian atas dengan nama outlet aktif, input pencarian cepat, lonceng notifikasi interaktif, bantuan, dan avatar dropdown profil pengguna.

