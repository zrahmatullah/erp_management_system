# Panduan Autentikasi, Notifikasi Popup, dan Standarisasi Icon UI

Dokumen ini menjelaskan perbaikan autentikasi API Login, integrasi sistem notifikasi popup/toast interaktif, serta standarisasi icon antarmuka di seluruh aplikasi menggunakan library **`lucide-vue-next`**.

---

## 1. Resolusi Issue API Login 401 Unauthorized

### Akar Masalah (Root Causes)
Ketika frontend mengirimkan permintaan `POST http://localhost:5173/api/v1/auth/login` (yang diproxy ke backend `http://localhost:8080/api/v1/auth/login`), sistem mengembalikan `401 Unauthorized`. Setelah investigasi menyeluruh pada log PostgreSQL dan internal backend Go, ditemukan dua akar masalah:

1. **Hash Bcrypt Tidak Valid pada Seed Data Awal**:
   - Hash password placeholder pada database seed (`$2a$10$N9qo8uLOickgx2ZMRZoMye...`) tidak cocok dengan plaintext yang diharapkan (`Admin@123`).
2. **Error Pemindaian Kolom Nullable (NULL Scan Error)**:
   - Pada repositori `backend/internal/repository/postgres/user_repo.go`, metode `GetByID`, `GetByEmail`, `GetByUsername`, dan `List` memindai kolom `phone` dan `avatar_url` langsung ke tipe primitif `string`.
   - Di PostgreSQL, nilai kedua kolom tersebut pada seed data adalah `NULL`. Driver `pgx` menolak pemindaian `NULL` ke dalam Go `string` (`converting NULL to string is unsupported`), menyebabkan pemanggilan `Login` mengembalikan error sebelum pengecekan password selesai.

### Tindakan Perbaikan yang Diterapkan
1. **Pembaruan Query Repository (`user_repo.go`)**:
   - Menggunakan fungsi SQL `COALESCE(phone, '')` dan `COALESCE(avatar_url, '')` sehingga nilai `NULL` dikonversi menjadi string kosong (`""`) yang kompatibel dengan struct domain Go:
   ```go
   query := `
       SELECT u.id, u.branch_id, u.role_id, u.username, u.email, u.password_hash,
              u.full_name, COALESCE(u.phone, ''), COALESCE(u.avatar_url, ''), u.pin_code,
              u.is_active, u.created_at, u.updated_at,
              r.name as role_name, COALESCE(b.name, '') as branch_name
       FROM users u
       JOIN roles r ON u.role_id = r.id
       LEFT JOIN branches b ON u.branch_id = b.id
       WHERE u.email = $1
   `
   ```
2. **Pembaruan Seed Akun dengan Bcrypt Hash yang Valid**:
   - Password standar untuk seluruh akun seed telah di-hash ulang menggunakan cost 12:
     - Plaintext: **`Admin@123`**
     - Hash: `$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK`
   - Migrasi telah dieksekusi ke database `cafe_erp`.

### Daftar Akun & Kredensial untuk Pengujian

| Username | Email | Password | Role | Akses Modul |
| :--- | :--- | :--- | :--- | :--- |
| **`admin`** | `admin@cafe-erp.com` | `Admin@123` | **Super Admin** | Akses Penuh (Dashboard, POS, Kitchen, Stock, HR, Finance, Master Hub) |
| **`sarah.manager`** | `sarah@cafe-erp.com` | `Admin@123` | **Branch Manager** | Operasional Cabang, Laporan, Stock Opname, Approval PO |
| **`jane.cashier`** | `cashier1@cafe-erp.com` | `Admin@123` | **Kasir (POS)** | POS Kasir, Transaksi Penjualan, Cetak Struk, Pembayaran |
| **`chef.john`** | `chef@cafe-erp.com` | `Admin@123` | **Kitchen / Barista**| Kitchen Display System (KDS), Update Status Pesanan |
| **`ahmad.warehouse`**| `warehouse@cafe-erp.com` | `Admin@123` | **Staff Gudang** | Manajemen Stok, Penerimaan Barang, Purchase Order |

---

## 2. Sistem Popup Notifikasi (Toast Notification System)

Sistem notifikasi modern dibangun tanpa library eksternal yang memberatkan, menggunakan kombinasi **Pinia Reactive Store** dan **Vue 3 Animated Component**.

### Arsitektur Notifikasi

```
[ Aksi Pengguna / API Call ]
            │
            ▼
[ notification.store.ts ] ──► (addToast: success / error / warning / info)
            │
            ▼
[ ToastContainer.vue ] (Teleport / Global Mount di App.vue)
            │
            ▼
[ Popup Toast Tampil di Kanan Atas ] (Auto-dismiss 4.5s, Close Button, Lucide Icon)
```

### File Komponen
1. **`frontend/src/stores/notification.store.ts`**:
   - Mengelola state daftar notifikasi aktif (`toasts`).
   - Menyediakan helper: `notify.success(msg, title)`, `notify.error(msg, title)`, `notify.warning(msg, title)`, `notify.info(msg, title)`.
   - Timer otomatis menghapus notifikasi setelah 4.5 detik.
2. **`frontend/src/components/common/ToastContainer.vue`**:
   - Menampilkan popup notifikasi dengan animasi transisi (`transition-group`), indikator warna (hijau, merah, kuning, biru), dan icon Lucide (`CheckCircle2`, `AlertCircle`, `AlertTriangle`, `Info`, `X`).
   - Posisi: *Top-Right Fixed Overlay* (`z-50`).
3. **Pemasangan Global di `frontend/src/App.vue`**:
   - Diletakkan di level akar aplikasi sehingga popup muncul di halaman manapun (Login, POS, Master Hub, Settings).

### Contoh Penggunaan pada Kode:
```typescript
import { useNotificationStore } from '@/stores/notification.store'
const notify = useNotificationStore()

// Login Sukses
notify.success('Selamat datang kembali, ' + user.full_name, 'Login Berhasil')

// Login Gagal
notify.error('Username atau kata sandi tidak valid', 'Gagal Masuk')

// CRUD Master Data
notify.success('Data Kategori berhasil diperbarui', 'Master Data Disimpan')
```

---

## 3. Standarisasi Icon UI (`lucide-vue-next`)

Semua emoji non-formal (seperti ☕, 🚪, 🛒, 👤, ⏰, 🔍, 👥, 🏖️, ✏️, 📦, 📊) telah digantikan 100% dengan icon vektor profesional dari **`lucide-vue-next`**.

### Pemetaan Icon per Komponen

| Halaman / Komponen | Elemen UI | Icon Lucide yang Digunakan |
| :--- | :--- | :--- |
| **Sidebar Navigation** | Brand Logo<br>Dashboard<br>POS Kasir<br>Kitchen Display<br>Inventory<br>Karyawan (HR)<br>Keuangan<br>Hak Akses<br>Master Hub<br>Logout | `Coffee`<br>`LayoutDashboard`<br>`UtensilsCrossed`<br>`MonitorPlay`<br>`Warehouse`<br>`Users`<br>`Receipt`<br>`Shield`<br>`Database` / `Settings`<br>`LogOut` |
| **Header (AppHeader)** | Search Bar<br>Bantuan<br>Notifikasi<br>Dropdown Profil & Cabang | `Search`<br>`HelpCircle`<br>`Bell`<br>`ChevronDown`, `User`, `Database`, `Settings`, `LogOut` |
| **Login (LoginView)** | Brand Logo<br>Username / Email Field<br>Password Field<br>Tombol Masuk | `Coffee`<br>`User`<br>`Lock`, `Eye`, `EyeOff`<br>`LogIn` |
| **Dashboard** | Revenue Metric<br>Total Orders Metric<br>Stock Alert Metric<br>Table Occupancy Metric | `DollarSign`<br>`ClipboardList`<br>`Package`<br>`Armchair`, `TrendingUp`, `AlertTriangle` |
| **POS Kasir** | Pilih Kasir<br>Waktu Shift<br>Cari Menu<br>Keranjang Belanja<br>Tambah/Kurang Qty | `User`<br>`Clock`<br>`Search`<br>`ShoppingCart`<br>`Plus`, `Minus` |
| **Payment Modal (POS)** | Dine-in / Takeaway<br>Nomor Meja<br>QRIS<br>Tunai (Cash)<br>Kartu Debit/Kredit | `Armchair`, `ShoppingBag`<br>`MapPin`<br>`QrCode`<br>`Banknote`<br>`CreditCard` |
| **Kitchen Display (KDS)**| Header KDS<br>Waktu Tunggu Pesanan<br>Selesaikan Pesanan | `ChefHat`<br>`Clock`<br>`Check` |
| **Manajemen Meja** | Denah Meja & Kapasitas | `Armchair`, `Users`, `Plus`, `Edit`, `Trash2` |
| **Manajemen Stok & Opname** | Cari Item<br>Refresh Stok<br>Tambah Item<br>Opname Filter & Validasi | `Search`<br>`RotateCw`<br>`Plus`<br>`Warehouse`, `Calendar`, `User`, `Archive` |
| **Manajemen Karyawan** | Total Karyawan<br>Aktif / Cuti<br>Aksi Detail / Edit | `Users`<br>`UserCheck`, `CalendarOff`<br>`Eye`, `Edit`, `MoreVertical` |
| **Slip Gaji (Payslip Modal)**| Header Cafe<br>Cetak / Download Slip | `Coffee`<br>`Printer`, `Download`, `X` |
| **Master Hub (13 Modul)** | Tab Navigasi Master Data | `Building2`, `User`, `Shield`, `FolderTree`, `Coffee`, `Package`, `Armchair`, `Receipt`, `Warehouse`, `Truck`, `Users`, `Briefcase`, `Clock` |

---

## 4. Cara Menjalankan & Menguji Aplikasi

### Menjalankan Backend Go
```powershell
cd d:\Project\cafe-erp-system\backend
go run cmd/server/main.go
```
*Server aktif pada port `8080` (Healthcheck: `http://localhost:8080/health`).*

### Menjalankan Frontend Vue
```powershell
cd d:\Project\cafe-erp-system\frontend
npm run dev
```
*Aplikasi aktif pada port `5173` (`http://localhost:5173`).*

### Langkah Uji Coba:
1. Buka browser pada alamat `http://localhost:5173`.
2. Masukkan email `admin@cafe-erp.com` dan password `Admin@123`.
3. Klik tombol **Masuk**. Popup notifikasi hijau *"Selamat datang kembali, Super Administrator!"* akan muncul di pojok kanan atas.
4. Coba masukkan password salah (misal: `Admin@999`). Popup notifikasi merah *"Username atau kata sandi tidak valid"* akan muncul.
5. Masuk ke menu **Super Admin -> Master Data Hub** (`http://localhost:5173/settings/master`).
6. Tambah atau ubah data master (misal: Cabang atau Kategori Menu). Popup konfirmasi hijau akan muncul secara otomatis saat data berhasil disimpan.

