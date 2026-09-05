# POS Table Management & Kasir Pembayaran (2 Tab) Guide

Dokumen ini menjelaskan implementasi fitur denah meja interaktif (`/pos/tables`) dan 2-Tab Kasir POS (`/pos`) yang terhubung langsung ke basis data PostgreSQL dan REST API Go.

---

## 1. Denah Meja Interaktif (`/pos/tables`)

Pada halaman Manajemen Meja (`TableManagementView.vue`), setiap kartu meja dapat diklik langsung untuk memunculkan modal detail (`TableDetailModal.vue`).

### Fitur Berdasarkan Status Meja:
1. **Meja Terisi (`occupied` / `billing`):**
   - **Informasi Pesanan Aktif:** Nomor pesanan (misal: `ORD-2026-0001`), nama tamu, status pesanan, dan durasi waktu tamu berada di meja.
   - **Rincian Menu yang Dipesan:** Daftar item lengkap dengan thumbnail, kuantitas, harga satuan, subtotal, catatan khusus (*notes*), dan badge status dapur (*Pending*, *Cooking*, *Served*).
   - **Rincian Billing:** Subtotal pesanan, Pajak PPN (10%), dan Total tagihan terkini.
   - **Pesan Menu Tambahan Langsung:** Terdapat menu dropdown/accordion untuk memilih menu tambahan dari katalog cepat, mengatur jumlah & catatan, lalu menekan tombol **"Kirim Menu Tambahan ke Dapur"** (`POST /api/v1/pos/orders/{id}/items`). Item baru akan langsung ditambahkan ke pesanan meja dan status dapur tercatat sebagai *Pending*.
   - **Aksi Cepat:**
     - Tombol **"Bayar Sekarang"**: Membuka modal pembayaran (`PaymentModal.vue`) secara langsung.
     - Tombol **"Ke Kasir Pembayaran"**: Mengalihkan kasir langsung ke POS Tab Pembayaran (`/pos?tab=billing&table={table_no}`).

2. **Meja Kosong (`available`):**
   - **Buka Pesanan Baru (Dine-in):** Memungkinkan pelayan/kasir memasukkan nama pelanggan, memilih item menu awal, dan langsung menekan **"Buka & Simpan Pesanan"** (`POST /api/v1/pos/orders`). Status meja otomatis berubah menjadi `occupied`.
   - **Set Meja ke Reservasi:** Tombol ubah status meja langsung ke `reserved`.

---

## 2. Kasir POS dengan 2 Tab (`/pos`)

Halaman Kasir POS (`POSView.vue`) dilengkapi sistem navigasi 2 Tab di bagian atas layar:

### Tab 1: Katalog & Pesanan Baru (`activePosTab = 'order'`)
- Digunakan untuk menerima pesanan baru secara walk-in maupun dine-in.
- Berisi navigasi kategori menu, input pencarian real-time, keranjang belanja (cart), pemilihan meja aktif, dan opsi submit:
  - **"Kirim ke Meja"** (Membuat pesanan aktif untuk meja terpilih).
  - **"Bayar Langsung"** (Untuk pelanggan takeaway/langsung bayar di tempat).

### Tab 2: Kasir Pembayaran Meja (`activePosTab = 'billing'`)
- Khusus digunakan untuk mengelola tagihan meja yang sedang aktif dan belum lunas.
- Menampilkan badge jumlah meja aktif (contoh: `2 Tagihan Aktif`).
- **Filter Zona Meja:** Filter cepat berdasarkan zona meja (`Semua Meja`, `Lantai 1`, `Lantai 2`, `Outdoor`, `Takeaway`).
- **Pencarian Tagihan:** Kasir dapat mencari berdasarkan nama tamu atau nomor meja.
- **Kartu Tagihan Meja:**
  - Nomor meja & badge durasi (e.g. `45 Menit lalu`).
  - Nomor pesanan & nama tamu.
  - Ringkasan item pesanan beserta kuantitasnya.
  - Total tagihan (termasuk pajak).
  - Tombol **"Cetak Bill Sementara"** untuk memberikan struk awal kepada pelanggan.
  - Tombol **"Bayar Meja Ini"**: Membuka modal pembayaran standar (`PaymentModal.vue`).

---

## 3. Alur Pembayaran & Pembebasan Meja

Saat kasir menekan tombol bayar pada Tab Pembayaran atau pada popup detail meja:
1. Modal `PaymentModal.vue` terbuka dengan total tagihan yang sudah terisi otomatis.
2. Kasir memilih metode pembayaran (Tunai, QRIS, Kartu Debit, Kartu Kredit).
3. Jika tunai, kasir dapat memilih pecahan cepat (`Uang Pas`, `50.000`, `100.000`, dll) dan kembalian dihitung secara otomatis.
4. Saat tombol **"Konfirmasi & Selesaikan Pembayaran"** ditekan:
   - Backend memproses endpoint `POST /api/v1/pos/orders/{id}/pay`.
   - Record transaksi baru dicatat di tabel `payments`.
   - Status pesanan di tabel `orders` diperbarui menjadi `'completed'`.
   - Meja terkait di tabel `cafe_tables` otomatis dikembalikan menjadi `'available'`.
   - Kasir diberikan opsi cetak struk resmi dan data daftar tagihan aktif otomatis diperbarui tanpa perlu refresh halaman.

---

## 4. REST API Endpoints Backend Go

Seluruh fitur ini didukung oleh endpoint REST API di backend:

| Endpoint | Method | Deskripsi |
| :--- | :--- | :--- |
| `/api/v1/pos/tables` | `GET` | Mengambil seluruh denah meja lengkap dengan objek `active_order` terkini (LATERAL JOIN). |
| `/api/v1/pos/tables/{id}/order` | `GET` | Mengambil detail spesifik pesanan aktif pada meja beserta item dan status dapur. |
| `/api/v1/pos/tables/{id}/status` | `PUT` | Memperbarui status meja (`available`, `reserved`, `occupied`, `billing`). |
| `/api/v1/pos/orders/active` | `GET` | Mengambil semua pesanan yang belum dibayar (`status NOT IN ('completed', 'cancelled')`). |
| `/api/v1/pos/orders/{id}/items` | `POST` | Menambahkan item menu pesanan baru ke pesanan meja yang sedang berjalan. |
| `/api/v1/pos/orders/{id}/pay` | `POST` | Memproses pelunasan pesanan, mencatat pembayaran, dan mengosongkan meja. |

