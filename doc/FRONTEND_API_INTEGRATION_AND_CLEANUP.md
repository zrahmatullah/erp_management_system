# Panduan Integrasi API End-to-End & Pembersihan Data Dummy Frontend

Dokumen ini mencatat pembersihan menyeluruh terhadap seluruh data dummy/hardcoded pada antarmuka frontend (seperti `const staffRoster = ref([...])`, dsb.) dan migrasinya ke panggilan REST API langsung ke basis data PostgreSQL.

---

## 1. Ringkasan Perubahan

Sesuai instruksi, seluruh penugasan array hardcoded statis di Vue ref (`ref([...])`) telah **dihapus 100%**. State komponen sekarang diinisialisasi dalam kondisi kosong (`ref<any[]>([])`) dan diisi secara dinamis melalui siklus hidup `onMounted` menggunakan library Axios.

Jika server API atau database mengembalikan data baru, antarmuka langsung merefleksikan data tersebut secara otomatis tanpa perlu mengubah kode frontend.

---

## 2. Pemetaan Modul & Endpoint API Backend

Berikut daftar komponen frontend yang telah dihubungkan penuh ke database:

| Komponen Vue | Halaman / Modul | Endpoint API Backend | Deskripsi Sumber Data Database |
| :--- | :--- | :--- | :--- |
| **`ShiftScheduleView.vue`** | HRIS → Jadwal Shift | `GET /api/v1/hris/schedules` | Mengambil data karyawan (`employees`) dan pola rotasi shift (`work_shifts`). |
| **`TableManagementView.vue`** | POS → Manajemen Meja | `GET /api/v1/pos/tables`<br>`POST /api/v1/master/tables` | Mengambil status meja live (Available, Occupied, Reserved) & menambah meja baru ke database. |
| **`AttendanceView.vue`** | HRIS → Presensi & Absensi | `GET /api/v1/hris/attendances` | Log presensi, jam masuk/keluar, status kehadiran, dan menit keterlambatan. |
| **`LeaveManagementView.vue`** | HRIS → Pengajuan Cuti | `GET /api/v1/hris/leaves` | Permohonan cuti, jenis cuti (Tahunan, Sakit, Alasan Penting), dan status persetujuan. |
| **`POSView.vue`** | POS → Kasir Penjualan | `GET /api/v1/pos/products`<br>`GET /api/v1/pos/tables` | Katalog menu, harga, kategori dinamis, dan pemilihan meja pesanan aktif. |
| **`StockOpnameView.vue`** | Inventory → Stock Opname | `GET /api/v1/inventory/stocks` | Stok sistem bahan baku, input hitungan fisik, dan kalkulasi selisih otomatis. |
| **`PurchaseOrderList.vue`** | Inventory → Purchase Orders | `GET /api/v1/inventory/purchase-orders`<br>`GET /api/v1/inventory/stocks` | Detail PO, status timeline approval, dan item pesanan supplier. |
| **`EmployeeList.vue`** | HRIS → Daftar Karyawan | `GET /api/v1/hris/employees` | Direktori staf cabang, departemen, jabatan, dan metrik dinamis. |
| **`StockList.vue`** | Inventory → Stok Bahan Baku | `GET /api/v1/inventory/stocks` | Stok aktual, batas minimum reorder, dan status level persediaan. |
| **`KitchenDisplayView.vue`** | POS → Kitchen Display (KDS) | `GET /api/v1/kds/tickets`<br>`PUT /api/v1/kds/items/:id/status` | Tiket dapur real-time dan pembaruan status (Memasak, Siap, Disajikan). |
| **`FinanceOverview.vue`** | Finance → Ikhtisar Keuangan | `GET /api/v1/finance/overview`<br>`GET /api/v1/finance/journals` | Revenue, biaya operasional, laba bersih, kas di tangan, dan daftar jurnal umum. |
| **`JournalEntryForm.vue`** | Finance → Form Jurnal | `GET /api/v1/master/accounts` | Bagan akun (COA) resmi dari master data keuangan. |
| **`SalesReport.vue`** | Reports → Laporan Penjualan | `GET /api/v1/pos/products` | Penjualan produk terlaris, kuantitas, dan kontribusi omzet. |
| **`ProductList.vue`** | Menu → Master Produk | `GET /api/v1/master/products` | Daftar menu penjualan langsung dari database. |
| **`CategoryList.vue`** | Menu → Kategori Menu | `GET /api/v1/master/categories` | Kategori menu dengan fungsionalitas penambahan kategori baru. |
| **`DashboardView.vue`** | Dashboard Utama | `GET /api/v1/dashboard/stats`<br>`GET /api/v1/pos/orders`<br>`GET /api/v1/pos/tables` | KPI penjualan, total transaksi, peringatan stok, pesanan terkini, dan denah meja mini. |

---

## 3. Implementasi Endpoint Baru: `GET /api/v1/hris/schedules`

Untuk memenuhi kebutuhan jadwal rotasi kerja mingguan staf tanpa data dummy:
* Ditambahkan method `GetShiftSchedules` pada [`backend/internal/delivery/http/handler/operational_handler.go`](file:///d:/Project/cafe-erp-system/backend/internal/delivery/http/handler/operational_handler.go).
* Route didaftarkan pada grup `/hris/schedules` di [`backend/internal/delivery/http/router.go`](file:///d:/Project/cafe-erp-system/backend/internal/delivery/http/router.go).
* Menggabungkan data dari tabel `employees` dan master `work_shifts` di database PostgreSQL.

---

## 4. Cara Menjalankan & Menguji

### 1. Jalankan Backend Go
```powershell
cd d:\Project\cafe-erp-system\backend
go run cmd/server/main.go
```
*Server mendengarkan pada `http://localhost:8080`.*

### 2. Jalankan Frontend Vite
```powershell
cd d:\Project\cafe-erp-system\frontend
npm run dev
```
*Aplikasi aktif pada `http://localhost:5173`.*

### 3. Verifikasi Data:
- Masuk ke menu **HRIS → Shift Schedule**: Seluruh karyawan (`Sarah Johnson`, `John Doe`, `Sarah Andini`, `Jane Doe`) muncul langsung dari database `cafe_erp`.
- Masuk ke menu **POS Kasir**: Produk `Cafe Latte`, `Caramel Macchiato`, dll. dimuat dari database dengan harga dasar dan stasiun pemrosesan aslinya.
- Masuk ke menu **Inventory → Stock Opname**: Stok sistem berasal dari tabel `inventory_stocks`.
- Masuk ke menu **Finance → Overview**: Total Revenue, Total Expenses, dan jurnal umum berimbang berasal dari kalkulasi order dan general ledger.

