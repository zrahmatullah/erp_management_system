# Panduan Alur Operasional POS, Manajemen Meja, Nomor Antrian & KDS

Dokumen ini menjelaskan arsitektur dan alur operasional komprehensif untuk modul **Point of Sale (POS)**, **Manajemen Meja & Denah**, sistem **Nomor Antrian Terpisah**, serta **Kitchen Display System (KDS)** pada Cafe ERP System.

---

## 1. Ringkasan Alur Kerja (End-to-End Workflow)

```
[ KASIR POS ]
  │
  ├─► Pilih Menu & Tipe:
  │     ├── [Dine-In] ────► Pilih Meja (e.g. T-02) ──► "Kirim ke Meja"
  │     │                                                 │
  │     │                                                 ▼
  │     │                                       [Status Meja -> Occupied (Merah)]
  │     │                                       [Generate Antrian: D-01]
  │     │
  │     ├── [Takeaway] ───► "Pesan Takeaway" ──► [Generate Antrian: TA-01]
  │     │
  │     └── [Delivery] ───► "Pesan Delivery" ──► [Generate Antrian: DL-01]
  │
  ▼
[ KITCHEN DISPLAY SYSTEM (KDS) ]
  ├── Filter Tab: Semua Tiket | Dine-In (#D) | Takeaway (#TA) | Delivery (#DL)
  ├── 4 Kanban Kolom: Pending -> Cooking -> Ready -> Served
  └── Badge Antrian Kontras + Timer Tunggu Berwarna (Hijau <10m, Kuning 10-20m, Merah >20m)
  │
  ▼
[ MANAJEMEN MEJA & DENAH ]
  ├── Tab 1-4: Zona Meja Fisik (Lantai 1, Lantai 2, Outdoor)
  │     └── Meja Terisi (Occupied): Badge Antrian (#D-xx), Order No, Tamu, Total Tagihan
  │     └── Klik Meja -> Modal Detail -> Tombol "Bayar Meja Ini"
  │
  └── Tab 5: "Pesanan Takeaway & Delivery" (List Antrian Aktif)
        └── Board antrian kartu: No Antrian (#TA-xx, #DL-xx), Pelanggan, Status Dapur, Total
        └── Tombol langsung: "Bayar & Selesaikan"
  │
  ▼
[ PROSES PEMBAYARAN KASIR ]
  ├── Kasir Pembayaran Meja (Tab 2 di POS) / Modal Pembayaran
  ├── Masukkan Nominal / Scan QRIS Dinamis / Kartu
  └── Klik "Bayar Sekarang"
        │
        ▼
  [Status Meja -> Available (Hijau Kembali Secara Otomatis)]
  [Status Order -> Completed]
  [Status Dapur -> Served]
```

---

## 2. Pemisahan Nomor Antrian Harian (Queue Numbering System)

Sistem secara otomatis menghasilkan nomor antrian berurutan harian (`queue_number`) yang dipisahkan berdasarkan tipe pesanan:

| Tipe Pesanan | Prefix | Format Contoh | Target Operasional |
| :--- | :--- | :--- | :--- |
| **Dine-In** | `D-` | `D-01`, `D-02`, dst. | Tamu yang makan di tempat, mengikat nomor meja fisik kafe |
| **Takeaway** | `TA-` | `TA-01`, `TA-02`, dst. | Tamu bungkus, menunggu di lounge pengambilan antrian |
| **Delivery** | `DL-` | `DL-01`, `DL-02`, dst. | Pesanan kurir online (GoFood, GrabFood, ShopeeFood) |

Nomor antrian ter-reset otomatis setiap hari berdasarkan fungsi tanggal database (`created_at >= CURRENT_DATE`).

---

## 3. Otomatisasi Status Meja (Table Occupancy Automation)

1. **Ketika Pesanan Dine-In Dibuat**:
   - Endpoint: `POST /api/v1/pos/orders`
   - Parameter: `order_type: 'dine_in'`, `table_number: 'T-xx'`
   - Database:
     ```sql
     UPDATE cafe_tables SET status = 'occupied', updated_at = NOW() WHERE table_number = 'T-xx';
     ```
   - Visual di Denah: Meja seketika berubah warna menjadi **Merah (Occupied)**, memunculkan badge `#D-xx`, nomor order, nama tamu, dan total harga pesanan.

2. **Ketika Pesanan Dibayar Lunas**:
   - Endpoint: `POST /api/v1/pos/orders/{id}/pay`
   - Database Transaction:
     - Catat record pembayaran di tabel `payments`.
     - Update order menjadi `completed`: `UPDATE orders SET status = 'completed' WHERE id = $1`.
     - Lepaskan meja kembali ke kondisi siap pakai:
       ```sql
       UPDATE cafe_tables SET status = 'available', updated_at = NOW() WHERE id = table_id;
       ```
   - Visual di Denah: Meja seketika kembali berwarna **Hijau (Available)** dan siap menerima tamu baru.

---

## 4. Fitur Baru Manajemen Meja: Tab Pesanan Takeaway & Delivery

Pada halaman `/pos/tables`, tersedia tab ke-5: **"Pesanan Takeaway & Delivery"**.
Tab ini menyediakan papan pemantauan antrian aktif non-meja:
- Menampilkan nomor antrian besar (`TA-01`, `DL-01`) dengan kode warna oranye (Takeaway) dan ungu (Delivery).
- Nama pelanggan dan timer tunggu berjalan (`⏱ x menit lalu`).
- Ringkasan item menu dan status dapur (`cooking` / `ready`).
- Catatan khusus pelanggan (misal: "Plastik terpisah", "Less ice").
- Tombol **"Bayar & Selesaikan"** yang langsung memunculkan modal pembayaran kasir.

---

## 5. Fitur Baru Kitchen Display System (KDS)

Halaman `/pos/kds` dilengkapi dengan antarmuka real-time dark mode:
1. **Sub-Tab Filter Tipe Antrian**:
   - `Semua Tiket`: Menampilkan seluruh pesanan aktif.
   - `Dine-In (Meja)`: Khusus tiket meja makan di tempat.
   - `Takeaway (Bungkus)`: Khusus pesanan kemasan bungkus.
   - `Delivery (Kurir)`: Khusus pesanan pengiriman kurir.
2. **Filter Station**:
   - `Barista`: Menampilkan produk minuman & kopi.
   - `Kitchen`: Menampilkan masakan dapur & makanan utama.
3. **Indikator Timer Urgensi**:
   - **Hijau** ($< 10$ menit): Standar waktu persiapan awal.
   - **Kuning** ($10 - 20$ menit): Peringatan pesanan perlu dipercepat.
   - **Merah Berkedip** ($> 20$ menit): Prioritas utama pesanan terlambat.
4. **Alur Transisi Kanban**:
   - `Pesanan Masuk (Pending)` $\rightarrow$ Klik *"Mulai Masak"* $\rightarrow$ `Sedang Dimasak (Cooking)`
   - `Sedang Dimasak (Cooking)` $\rightarrow$ Klik *"Selesai Masak"* $\rightarrow$ `Siap Saji (Ready)`
   - `Siap Saji (Ready)` $\rightarrow$ Klik *"Sajikan / Pick-Up"* $\rightarrow$ `Disajikan (Served)`
5. **Auto Refresh**:
   - KDS melakukan sinkronisasi otomatis dengan server setiap 12 detik tanpa perlu me-refresh browser secara manual.

---

## 6. Daftar Endpoint API Terkait

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/v1/pos/tables` | Mengambil seluruh meja beserta info order aktif (`active_order`) |
| `GET` | `/api/v1/pos/tables/{id}/order` | Mengambil detail order dan rincian item pesanan aktif dari meja |
| `PUT` | `/api/v1/pos/tables/{id}/status` | Mengubah status meja secara manual (`available`, `reserved`, dll.) |
| `GET` | `/api/v1/pos/orders/active` | Mengambil seluruh tagihan aktif di kasir untuk Tab 2 |
| `POST` | `/api/v1/pos/orders` | Membuat pesanan baru, generate `queue_number`, dan update meja `occupied` |
| `POST` | `/api/v1/pos/orders/{id}/pay` | Memproses pembayaran, set order `completed`, dan kembalikan meja ke `available` |
| `GET` | `/api/v1/pos/takeaways` | Mengambil daftar antrian aktif Takeaway & Delivery untuk Tab 5 |
| `GET` | `/api/v1/kds/tickets` | Mengambil tiket KDS dapur dengan nomor antrian dan tipe pesanan |
| `PUT` | `/api/v1/kds/items/{id}/status` | Mengupdate progres masak item dapur (`cooking`, `ready`, `served`) |
