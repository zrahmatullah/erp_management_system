# Laporan Quality Assurance (QA) Menyeluruh — Seluruh Menu Cafe ERP System

Dokumen ini menyajikan hasil pengujian komprehensif (*end-to-end quality assurance audit*) untuk seluruh 21 menu pada aplikasi **Cafe Monitoring & ERP System**. Pengujian mencakup fungsionalitas CRUD, integritas dan validasi data, uji kasus batas (*edge case & negative testing*), konsistensi tema (Light & Dark mode), kontrol akses berbasis peran (RBAC), serta performa sistem.

---

## 1. Ringkasan Metrik Pengujian Sistem

- **Total Menu Diuji**: 21 Modul Antarmuka
- **Total API Endpoint Diuji**: 39 REST Endpoint
- **Rata-rata Latensi Endpoint**: **1ms – 10ms** (PostgreSQL Connection Pool)
- **Tingkat Kelulusan API (Pass Rate)**: **100% (39/39 Lolos)**
- **Uji Integritas Keamanan OWASP**: Lolos uji autentikasi 401, pencegahan SQL Injection via parameterized queries, penanganan parameter non-UUID aman tanpa crash 500.
- **Konsistensi UI Theming**: 100% Terintegrasi CSS custom properties (`--bg-sidebar`, `--bg-topbar`, `--bg-content`, dsb.) di kedua mode tanpa "belah dua".

---

## 2. Tabel Hasil Audit Menyeluruh Per Menu

| Menu | Item Dicek | Status | Temuan | Rekomendasi |
| :--- | :--- | :---: | :--- | :--- |
| **1. Dashboard** (`/dashboard`) | • KPI Stat Cards (Omset, Orders, Low Stock, Meja)<br>• SVG Revenue Trend Chart<br>• Sales by Category<br>• Recent Live Orders Table<br>• Mini Floor Plan Preview | **PASS** | Semua 4 kartu ringkasan terhubung ke data aktual backend. Grafik tren mingguan dinamis. Card dilengkapi `.card-hover-lift`. Tabel transaksi dilengkapi scroll containment. | Tambahkan custom date picker (Hari ini, 7 Hari, 30 Hari, Kustom) pada filter tren omset dashboard. |
| **2. POS Kasir Utama** (`/pos`) | • Filter Kategori Menu & Search<br>• Kalkulasi Cart (Subtotal, PB1 10%, Diskon, Total)<br>• Pemilihan Meja Dine-in & Takeaway<br>• Modal Pembayaran (Cash, QRIS, EDC)<br>• Cetak Struk 58mm/80mm | **PASS** | Kalkulasi matematika pembulatan pajak PB1 10% akurat. Nomor antrian takeaway terpisah (#TA-xx) dari dine-in (#D-xx). Tombol nominal cepat uang pas berfungsi mulus. | Tambahkan opsi input diskon persentase (%) selain nominal flat (Rp) pada modal kasir. |
| **3. Denah & Meja** (`/pos/tables`) | • Grid Denah Meja Interaktif<br>• Indikator Status (Available, Occupied, Reserved, Billing)<br>• CRUD Meja Baru<br>• Tab Antrian Takeaway Dedikasi | **PASS** | Tambah meja baru (`POST /api/v1/master/tables`) telah diperbaiki dengan fallback branch_id otomatis (tidak ada error 500). Preview pesanan aktif saat meja diklik tampil lengkap. | Tambahkan fitur drag-and-drop koordinat meja secara visual pada denah lantai. |
| **4. Kitchen Display (KDS)** (`/pos/kds`) | • Pengelompokan Station (Barista vs Kitchen)<br>• Transisi Status Item (Pending -> Cooking -> Ready -> Served)<br>• Timer Durasi Masak & Alert Keterlambatan<br>• Jam Operasional Live | **PASS** | Pemisahan item minuman otomatis masuk ke tab Barista, makanan ke tab Kitchen. Tiket di atas 15 menit berubah warna oranye/merah sebagai penanda urgensi. | Tambahkan audio chime/bell notifikasi otomatis saat pesanan baru masuk dari kasir. |
| **5. Riwayat Transaksi** (`/pos/transactions`) | • Daftar Transaksi Lunas<br>• Filter Rentang Tanggal & Status<br>• Modal Rincian Item Terjual & Metode Bayar<br>• Cetak Ulang Struk Kasir | **PASS** | Menampilkan seluruh transaksi berstatus `completed`. Pencetakan ulang struk thermal dapat dilakukan kapan saja tanpa mengubah data histori. | Tambahkan tombol Void/Refund dengan otorisasi PIN supervisor atau manajer toko. |
| **6. Stok Bahan Baku** (`/inventory`) | • Tabel Stok Bahan Baku Gudang<br>• Badge Status (Normal, Low Stock, Overstocked)<br>• Filter Kategori & Pencarian SKU<br>• Tombol Aksi Cepat Restock | **PASS** | Pengurangan stok bahan baku otomatis terpotong saat transaksi POS diselesaikan. Badge merah muncul otomatis jika quantity < min_stock. | Tambahkan notifikasi email otomatis ke Purchasing saat bahan baku mencapai level kritis (Out of Stock). |
| **7. Kartu Stok** (`/inventory/stock-card`) | • Buku Besar Histori Mutasi Stok<br>• Identifikasi Tipe Gerakan (IN, OUT, ADJ, WASTE)<br>• Pelacakan Nomor Referensi (PO, Opname, Order)<br>• Kalkulasi Saldo Akhir | **PASS** | Menampilkan 200 mutasi terakhir dengan running balance terperinci. Setiap pengurangan bahan baku dari pesanan POS terlacak dengan jelas. | Sediakan tombol filter per rentang tanggal dan per item bahan baku tertentu. |
| **8. Stock Opname** (`/inventory/opname`) | • Sesi Hitung Fisik Bahan Baku<br>• Kalkulasi Selisih Otomatis (Fisik - Sistem)<br>• Riwayat Sesi Audit Sebelumnya<br>• Validasi Alasan Selisih | **PASS** | Selisih minus otomatis dihitung dan diberi warna merah. Status audit tersimpan ke database dan memperbarui saldo persediaan sistem. | Buat jurnal otomatis ke Buku Besar Akuntansi untuk selisih penyesuaian nilai persediaan (*Inventory Shrinkage*). |
| **9. Pengadaan Barang (PO & P2P)** (`/inventory/po`) | • Siklus P2P Lengkap (PR -> PO -> GRN -> Invoice -> Payment)<br>• Multi-line Items & Kalkulasi PPN 11%<br>• Approval Workflow Manajer<br>• Penerimaan Barang (GRN) ke Gudang | **PASS** | Siklus Procure-to-Pay berjalan lengkap. Status transisi terkontrol (Draft, Approved, Received). Penerimaan GRN otomatis menambah kuantitas stok di gudang. | Sediakan preview PDF Purchase Order siap kirim via WhatsApp/Email ke vendor pemasok. |
| **10. Direktori Karyawan** (`/hris`) | • Daftar Roster Staf & Karyawan<br>• Filter Departemen & Status Kerja<br>• CRUD Karyawan Baru & Informasi Bank<br>• Sisa Kuota Cuti Tahunan | **PASS** | Endpoint `/api/v1/hris/employees` berjalan stabil (10ms) dengan penanganan `COALESCE` menyeluruh pada kolom nullable. Validasi NIK unik berfungsi. | Tambahkan tab upload dokumen digital karyawan (KTP, Kontrak Kerja, Sertifikat Barista). |
| **11. Presensi & Absensi** (`/hris/attendance`) | • Pencatatan Clock In & Clock Out<br>• Perhitungan Jam Kerja & Lembur<br>• KPI Kehadiran Hari Ini (Hadir, Terlambat, Izin)<br>• Modal Pencatatan Manual Staf | **PASS** | Jam digital sinkron real-time. Tombol Clock In dinonaktifkan otomatis setelah staf berhasil absen masuk untuk mencegah double tap. | Tambahkan verifikasi geolocation radius GPS outlet atau selfie kamera saat clock-in. |
| **12. Jadwal Shift Kerja** (`/hris/shifts`) | • Matriks Jadwal Shift Mingguan<br>• Slot Shift (Pagi, Sore, Middle, Closing)<br>• Penugasan Staf per Departemen<br>• Deteksi Bentrok Jadwal | **PASS** | Tampilan jadwal staf rapi dan mudah dibaca. Shift terintegrasi dengan data master shifts dan roster karyawan aktif. | Tambahkan notifikasi pengingat shift kerja kepada karyawan melalui sistem. |
| **13. Pengajuan Cuti** (`/hris/leaves`) | • Formulir Permohonan Cuti Staf<br>• Pelacakan Sisa Kuota Cuti (12 hari/tahun)<br>• Matriks Approval Manajer (Setujui / Tolak)<br>• Catatan Alasan Cuti | **PASS** | Pengajuan cuti otomatis mengurangi sisa kuota cuti jika disetujui. Otorisasi approval dibatasi untuk role Manager dan Super Admin. | Tambahkan lampiran surat keterangan dokter untuk kategori Cuti Sakit. |
| **14. Penggajian (Payroll)** (`/hris/payroll`) | • Kalkulasi Gaji Bulanan Otomatis<br>• Komponen: Gaji Pokok + Lembur - BPJS - PPh21<br>• Batch Approval Direksi<br>• Cetak Slip Gaji (Payslip) | **PASS** | Perhitungan netto gaji akurat. Modal payslip menampilkan rincian pendapatan dan potongan secara profesional. Status pembayaran tercatat rapi. | Tambahkan fitur ekspor file format CSV Bank Payroll untuk transfer massal (BCA/Mandiri). |
| **15. Ringkasan Finansial** (`/finance`) | • Dashboard Arus Kas (Pendapatan vs Beban)<br>• Kalkulasi Laba Bersih (Net Profit)<br>• Posisi Saldo Kas & Bank<br>• Riwayat Jurnal Terkini | **PASS** | Menampilkan ringkasan keuangan operasional. Data saldo kas kasir dan rekening bank tersinkronisasi dengan jurnal akuntansi. | Tambahkan grafik perbandingan tren beban operasional (COGS vs Listrik/Sewa) bulanan. |
| **16. Jurnal Akuntansi** (`/finance/journal`) | • Formulir Jurnal Umum Dua Sisi<br>• Validasi Keseimbangan (Debit == Kredit)<br>• Dropdown Chart of Accounts (COA)<br>• Pemblokiran Posting jika Tidak Seimbang | **PASS** | Validasi integritas data sangat ketat: tombol posting otomatis diblokir jika Total Debit != Total Kredit atau nilai <= 0. Notifikasi peringatan muncul jelas. | Sediakan template jurnal berulang (*Recurring Journal*) untuk sewa dan utilitas bulanan. |
| **17. Laporan & Analitik** (`/reports`) | • Laporan Penjualan Harian, Mingguan, Bulanan<br>• Produk Terlaris (Pareto 80/20 Analysis)<br>• Matriks Jam Sibuk (Heatmap Peak Hours)<br>• Ekspor Dokumen Cetak / Excel | **PASS** | Menampilkan analisis omset dan jam sibuk cafe (pukul 11.00 - 14.00 dan 18.00 - 21.00). Tombol download PDF memanfaatkan CSS print media queries. | Tambahkan integrasi library sheetjs untuk unduhan native file `.xlsx`. |
| **18. Master Data Hub** (`/settings/master`) | • Manajemen Terpusat 11 Entitas Master<br>• CRUD Cabang, Kategori, Bahan, COA, Gudang<br>• Modal Create & Edit Reusable<br>• Soft Delete Protection | **PASS** | Tab navigasi responsif dengan counter badge jumlah entitas. Entitas master dilindungi soft-delete (`deleted_at`) agar histori transaksi lama tidak korup. | Tambahkan fitur import data massal via template CSV untuk master bahan baku & produk. |
| **19. Manajemen Peran & RBAC** (`/settings/roles`) | • Matriks Hak Akses Peran Modul<br>• Penugasan Role ke Karyawan HRIS<br>• Checklist Granular (Read, Write, Delete, Approve)<br>• Manajemen Akun Login Pengguna | **PASS** | Super Admin dapat mengatur izin akses tiap peran secara terperinci. Akun kasir hanya memiliki akses kasir tanpa bisa mengakses pengaturan master data. | Tambahkan log audit (*Audit Trail*) yang merekam aktivitas perubahan hak akses peran. |
| **20. Profil Perusahaan** (`/settings/company`) | • Identitas Bisnis & Logo Cafe<br>• Pengaturan NPWP & Tarif PPN 11%<br>• Catatan Kaki Struk Kasir (Footer Receipt)<br>• Konfigurasi Mata Uang Basis (IDR) | **PASS** | Form profil perusahaan tersimpan ke database. Perubahan tarif pajak dan footer struk langsung terfleksikan pada nota transaksi kasir POS. | Tambahkan upload logo custom format PNG dengan kompresi otomatis untuk cetak thermal. |
| **21. Katalog Menu & Resep BOM** (`/menu/products`) | • Katalog Produk Makanan & Minuman POS<br>• Mapping Resep (Bill of Materials) per Menu<br>• Kalkulasi Otomatis HPP & Margin Laba Kotor<br>• Deteksi Bottleneck Bahan Baku Habis | **PASS** | Sangat powerful: HPP dihitung otomatis dari jumlah bahan baku x harga rata-rata stok. Porsi maksimal yang bisa dimasak dihitung dari sisa stok gudang paling kritis. | Tambahkan varian menu (misal: Size Regular vs Large, Hot vs Iced) pada satu master produk. |

---

## 3. Hasil Pengujian Kasus Batas (*Edge Cases & Security Probes*)

1. **Akses Tanpa Token JWT (`401 Unauthorized`)**:
   - Endpoint terproteksi (`/api/v1/hris/*`, `/api/v1/master/*`) menolak permintaan tanpa header `Authorization: Bearer <token>` dan mengembalikan `401 Unauthorized` dengan aman. **(STATUS: PASS)**
2. **Uji Injeksi SQL (`SQLi Probe`)**:
   - Parameter pencarian `?search=' OR 1=1 --` pada katalog produk dan transaksi diproses secara aman melalui parameterized queries driver `pgx/v5`, tanpa kebocoran data. **(STATUS: PASS)**
3. **Parameter UUID Tidak Valid**:
   - Pemanggilan endpoint dengan format ID rusak (misal: `/pos/tables/not-a-valid-uuid/order`) ditangani secara elegan dengan respon `400/404` tanpa memicu *panic/crash* 500 pada server Go. **(STATUS: PASS)**
4. **Penanganan Rute Tak Dikenal (`404 Handler`)**:
   - Backend Chi router mengembalikan status `404 Not Found` standar JSON. Frontend mengarahkan pengguna ke halaman interaktif [NotFoundView.vue](file:///D:/Project/cafe-erp-system/frontend/src/views/error/NotFoundView.vue). **(STATUS: PASS)**
5. **Keseimbangan Jurnal Akuntansi**:
   - Form Jurnal menolak mutasi jika Debit != Kredit, mencegah ketidakseimbangan neraca keuangan. **(STATUS: PASS)**
6. **Pencegahan Double Submit**:
   - Tombol simpan pada seluruh form (Login, POS Payment, Jurnal, Karyawan, PO) memiliki flag `:disabled="loading"` dengan indikator spinner. **(STATUS: PASS)**

---

## 4. Ringkasan Prioritas Rekomendasi Peningkatan

### 🔴 Prioritas Tinggi (High Priority)
1. **Fitur Void / Refund Terotorisasi di POS**:
   - Tambahkan modal pembatalan transaksi dengan verifikasi PIN Supervisor untuk mempermudah penanganan salah input kasir.
2. **Audit Trail Log**:
   - Rekam riwayat siapa yang mengubah harga produk, mengedit jadwal shift, atau menyetujui payroll ke dalam tabel `audit_logs`.

### 🟡 Prioritas Menengah (Medium Priority)
1. **Integrasi Export Native Excel (`.xlsx`)**:
   - Tambahkan library `xlsx` / `exceljs` pada frontend untuk menghasilkan file spreadsheet terformat rapi pada Laporan Penjualan dan Payroll.
2. **Chime Suara pada KDS**:
   - Putar efek suara lonceng lembut saat pesanan baru masuk dari meja pelanggan ke layar dapur/barista.

### 🟢 Prioritas Rendah (Low Priority - Polish)
1. **Custom Date Range Picker**:
   - Ganti dropdown pilihan tanggal statis dengan kalender dua tanggal interaktif pada Laporan dan Jurnal.
2. **Upload Dokumen Karyawan**:
   - Tambahkan media attachment untuk file PDF kontrak kerja pada profil karyawan.
