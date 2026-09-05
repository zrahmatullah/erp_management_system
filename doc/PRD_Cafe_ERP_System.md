# ☕ Product Requirements Document (PRD): Cafe ERP System

| Dokumen | Spesifikasi Kebutuhan Produk & Desain Database |
|---|---|
| **Nama Proyek** | Cafe ERP System (All-in-One Cafe Management Platform) |
| **Versi** | 1.0.0 |
| **Status** | Approved for Development |
| **Target OS & Arsitektur** | Web-based (Go Chi Backend + Vue 3 Frontend + PostgreSQL + Redis) |

---

## 1. Executive Summary & Visi Produk

### 1.1 Visi
Membangun platform Enterprise Resource Planning (ERP) terpadu khusus industri F&B/Kafe yang mengintegrasikan penjualan kasir (POS), pesanan dapur & bar (KDS), otomatisasi pengurangan bahan baku resep (BOM), logistik multi-cabang, absensi & penggajian (HRIS), serta pembukuan akuntansi otomatis dalam satu antarmuka terpusat.

### 1.2 Masalah Operasional yang Diselesaikan
1. **Kebocoran & Selisih Bahan Baku**: Ketiadaan sinkronisasi otomatis antara pesanan kasir dengan stok fisik bahan dasar (biji kopi, susu, sirup, cup).
2. **Keterlambatan Komunikasi Dapur/Bar**: Pesanan manual berbasis kertas rentan hilang atau salah urutan.
3. **Kompleksitas Shift & Lembur**: Pencatatan jadwal kerja barista/waiter yang berantakan menyulitkan perhitungan slip gaji bulanan.
4. **Pencatatan Keuangan Terpisah**: Kasir menggunakan aplikasi terpisah dari pembukuan akunting sehingga perhitungan laba-rugi (*Profit & Loss*) memakan waktu berhari-hari.

---

## 2. Persona Pengguna & Tingkat Akses (RBAC)

| Peran (Role) | Tanggung Jawab Utama | Akses Sistem |
|---|---|---|
| **Super Admin** | Konfigurasi sistem global, setup cabang, manajemen lisensi | Akses penuh (*Full Access*) |
| **Owner / Direksi** | Pemantauan KPI omset, margin laba-rugi, persetujuan PO bernilai besar | Dashboard analitik, approval PO, laporan keuangan |
| **Branch Manager** | Operasional harian cabang, approval cuti, stock opname, shift tukar | Modul cabang, persetujuan internal, supervisor POS |
| **Cashier (Kasir)** | Transaksi penjualan, buka/tutup laci kasir, cetak struk | POS, Meja, Pembayaran |
| **Barista & Kitchen** | Memproses minuman & makanan sesuai tiket KDS | Layar Kitchen Display System (KDS) |
| **Warehouse / Logistik** | Menerima pengiriman PO, transfer antar cabang, opname | Modul Inventaris, Stock Movement |
| **HR Admin** | Data staf, rekap absensi, approval izin, perhitungan gaji | Modul HRIS & Payroll |
| **Accountant** | Jurnal penyesuaian, kas kecil (petty cash), neraca, pajak | Modul Keuangan & Buku Besar |

---

## 3. Rincian Fitur & Modul Fungsional

### 🔐 Modul 1: Autentikasi, Cabang & Hak Akses
- **Multi-Branch Tenant**: Sistem mendukung banyak outlet/cabang dalam 1 akun perusahaan.
- **Role & Permission Matrix**: Hak akses granular per tombol/modul (*View, Create, Edit, Delete, Approve, Export*).
- **Fast Switch / Cashier PIN**: Pergantian kasir cepat tanpa logout akun utama menggunakan 4-6 digit PIN.
- **Audit Logging**: Mencatat timestamp, IP, User ID, dan detail setiap aksi kritikal (void transaksi, diskon manual, koreksi stok).

### ☕ Modul 2: Menu, Varian, Modifiers & Resep (BOM)
- **Kategori & Katalog Produk**: Struktur hierarki menu (e.g., *Beverages -> Coffee -> Espresso Based*).
- **Varian Produk**: Pilihan ukuran (*Regular / Large*), suhu (*Hot / Iced*), tipe biji kopi (*House Blend / Single Origin*).
- **Modifiers & Add-ons**: Tambahan opsional berbayar (*Extra Shot, Oat Milk, Vanilla Syrup*).
- **Bill of Materials (BOM / Resep)**:
  - Setiap varian terikat pada komposisi bahan baku dari inventaris.
  - *Contoh Resep Iced Caramel Latte (Large)*:
    - 20 gram Espresso Beans
    - 220 ml Fresh Milk
    - 25 ml Caramel Syrup
    - 1 pcs Cup 16oz + Straw + Lid
- **Multi-UOM (Unit of Measure)**: Konversi otomatis antara satuan beli (*kg, liter, karton*) ke satuan pakai (*gram, ml, pcs*).

### 🖥️ Modul 3: Point of Sale (POS) & Manajemen Meja
- **Tipe Pesanan**:
  - *Dine-in*: Terhubung ke nomor meja di denah.
  - *Takeaway*: Antrean nama/nomor panggil.
  - *Delivery / Online*: Integrasi nomor resi / pesanan pihak ketiga.
- **Visual Table Floor Plan**:
  - Peta visual meja cafe (Zonasi: *Indoor AC, Outdoor, VIP, Bar Barista*).
  - Status meja dinamis: *Kosong (Hijau), Terisi (Merah), Tagihan Dicetak (Kuning), Kotor (Abu-abu)*.
- **Operasional Kasir**:
  - Buka shift (*Starting Cash Drawer*), Tutup shift (*Cash Drop & Rekonsiliasi Kasir*).
  - Fitur *Split Bill* (per item atau bagi rata) dan *Merge Bill* (gabung meja).
  - Pajak Restoran (PB1 10%) dan Service Charge (5%) dinamis.
- **Multi-Payment Gateway**:
  - Tunai (hitung kembalian otomatis).
  - QRIS Dinamis/Statis.
  - Mesin EDC (Debit / Kredit).
  - Voucher Promo & Loyalty Points.

### 🍳 Modul 4: Kitchen Display System (KDS) & Bar Station
- **Routing Otomatis**: Pesanan makanan dikirim ke monitor Kitchen, pesanan kopi/minuman dikirim ke monitor Barista.
- **Status Tiket Real-time**:
  - *Pending* (Baru masuk)
  - *Preparing* (Sedang dibuat)
  - *Ready* (Siap saji)
  - *Served* (Sudah diantar ke meja)
- **Time Alerts**: Indikator warna kartu tiket berdasarkan durasi tunggu (> 10 menit kuning, > 15 menit merah berkedip).

### 📦 Modul 5: Logistik, Inventaris & Pengadaan (Supply Chain)
- **Otomasi Pemotongan Stok**: Begitu pesanan POS berstatus dibayar/diproses, stok bahan baku berkurang otomatis sesuai BOM.
- **Multi-Warehouse**: Pemisahan gudang cabang utama (*Main Storage*), pos bar (*Bar Stock*), dan dapur (*Kitchen Stock*).
- **Purchase Order (PO)**:
  - Rekomendasi pembelian otomatis saat bahan mencapai *Minimum Reorder Point*.
  - Alur persetujuan bertingkat (*Approval Workflow*).
  - Penerimaan Barang (*Goods Receipt Note*) yang memperbarui nilai rata-rata HPP (*Moving Average Cost*).
- **Stock Opname & Penyesuaian**:
  - Penghitungan fisik berkala vs stok sistem.
  - Pencatatan selisih: *Wastage, Rusak, Kedaluwarsa, atau Hilang*.

### 👥 Modul 6: HRIS, Penjadwalan Shift & Payroll
- **Data Induk Karyawan**: NIK, profil, kontak darurat, jabatan, status kontrak, info rekening bank.
- **Roster & Penjadwalan Shift**:
  - Pengaturan template shift (*Pagi 07:00-15:00, Siang 12:00-20:00, Closing 15:00-23:00*).
  - Pengajuan & persetujuan tukar shift antar-karyawan.
- **Absensi & Kehadiran**:
  - Clock-in & Clock-out berbasis PIN atau lokasi geofencing.
  - Perhitungan otomatis menit keterlambatan dan jam lembur.
- **Penggajian (Payroll)**:
  - Komponen: Gaji Pokok, Tunjangan Kehadiran, Uang Lembur, Insentif Target Omset.
  - Potongan: Keterlambatan, Kasbon, BPJS.
  - Generasi slip gaji digital (PDF).

### 💰 Modul 7: Keuangan & Akuntansi (Finance & Accounting)
- **Bagan Akun Standar (Standard Chart of Accounts)**:
  - Aset (Kas Kasir, Kas Operasional, Rekening Bank, Persediaan Bahan Baku).
  - Kewajiban (Hutang Usaha Supplier, Hutang Pajak Restoran).
  - Ekuitas (Modal Disetor, Laba Ditahan).
  - Pendapatan (Penjualan Makanan, Minuman, Merchandise, Biaya Servis).
  - HPP (Harga Pokok Penjualan Bahan Baku).
  - Beban (Gaji, Listrik, Sewa Tempat, Kerusakan Bahan).
- **Auto-Journaling Engine**:
  - Penjualan POS otomatis menjurnal: `Debit Kas/Bank, Kredit Pendapatan, Debit HPP, Kredit Persediaan`.
  - Pembelian bahan baku otomatis menjurnal: `Debit Persediaan, Kredit Hutang Usaha / Kas`.
- **Kas Kecil (Petty Cash)**: Klaim pengeluaran harian darurat oleh Supervisor Cabang.
- **Laporan Finansial**: Laporan Laba Rugi (*P&L*), Neraca (*Balance Sheet*), dan Arus Kas (*Cash Flow*).

### 📈 Modul 8: Laporan & Analitik (Dashboard BI)
- **Metrik Utama (Real-time KPI)**: Total omset hari ini, jumlah tiket pesanan, rata-rata belanja tamu (*AOV*), meja aktif.
- **Analisis Menu Engineering (Boston Matrix)**:
  - *Stars* (Margin tinggi, volume tinggi).
  - *Cash Cows* (Margin rendah, volume tinggi).
  - *Puzzles* (Margin tinggi, volume rendah).
  - *Dogs* (Margin rendah, volume rendah).
- **Laporan Jam Sibuk (Hourly Heatmap)**: Membantu penjadwalan efisien staf di jam terpadat.

---

## 4. Desain Skema Database (PostgreSQL)

```
                    ┌─────────────┐
                    │  branches   │
                    └──────┬──────┘
                           │
       ┌───────────────────┼───────────────────┐
       ▼                   ▼                   ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│    users     │    │ table_zones  │    │  warehouses  │
└──────┬───────┘    └──────┬───────┘    └──────┬───────┘
       │                   ▼                   ▼
       │            ┌──────────────┐    ┌──────────────┐
       │            │ cafe_tables  │    │invent_stocks │
       │            └──────┬───────┘    └──────┬───────┘
       ▼                   ▼                   ▲
┌──────────────┐    ┌──────────────┐           │ (deduct)
│    orders    │◄───┤ order_items  │           │
└──────┬───────┘    └──────┬───────┘           │
       │                   ▼                   │
       │            ┌──────────────┐           │
       │            │   recipes    ├───────────┘
       ▼            └──────┬───────┘
┌──────────────┐           ▼
│   payments   │    ┌──────────────┐
└──────────────┘    │invent_items  │
                    └──────────────┘
```

### 4.1 Tabel Organisasi & Pengguna
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL, -- super_admin, manager, cashier, barista, chef, hr, accountant
    pin_code VARCHAR(6),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### 4.2 Tabel Meja & Denah Cafe
```sql
CREATE TABLE table_zones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL -- 'Indoor AC', 'Outdoor Smoking', 'Lantai 2'
);

CREATE TABLE cafe_tables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    zone_id UUID REFERENCES table_zones(id) ON DELETE SET NULL,
    table_number VARCHAR(20) NOT NULL,
    capacity INT DEFAULT 4,
    status VARCHAR(20) DEFAULT 'available', -- 'available', 'occupied', 'billed', 'dirty'
    pos_x INT DEFAULT 0,
    pos_y INT DEFAULT 0
);
```

### 4.3 Tabel Menu & Resep (BOM)
```sql
CREATE TABLE menu_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES menu_categories(id),
    name VARCHAR(150) NOT NULL,
    sku VARCHAR(50) UNIQUE,
    description TEXT,
    base_price NUMERIC(12, 2) NOT NULL,
    target_station VARCHAR(20) NOT NULL DEFAULT 'kitchen', -- 'kitchen', 'barista', 'cashier'
    image_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, -- e.g. 'Hot', 'Iced', 'Large'
    price_adjustment NUMERIC(12, 2) DEFAULT 0.00
);

CREATE TABLE inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    uom VARCHAR(20) NOT NULL, -- 'gram', 'ml', 'pcs', 'pack'
    category VARCHAR(50),
    min_stock NUMERIC(12, 2) DEFAULT 0,
    average_cost NUMERIC(12, 2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id),
    quantity_required NUMERIC(12, 4) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### 4.4 Tabel Inventaris & Logistik
```sql
CREATE TABLE warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inventory_stocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    quantity NUMERIC(12, 4) DEFAULT 0,
    UNIQUE(warehouse_id, inventory_item_id)
);

CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    type VARCHAR(30) NOT NULL, -- 'in_purchase', 'out_pos_sales', 'waste', 'transfer', 'adjustment'
    quantity NUMERIC(12, 4) NOT NULL,
    cost_per_unit NUMERIC(12, 2),
    reference_id UUID,
    remarks TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### 4.5 Tabel Pesanan POS & Pembayaran
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    table_id UUID REFERENCES cafe_tables(id) ON DELETE SET NULL,
    cashier_id UUID REFERENCES users(id),
    waiter_id UUID REFERENCES users(id),
    order_type VARCHAR(20) NOT NULL DEFAULT 'dine_in', -- 'dine_in', 'takeaway', 'delivery'
    status VARCHAR(30) NOT NULL DEFAULT 'pending',     -- 'pending', 'processing', 'completed', 'void'
    customer_name VARCHAR(100),
    subtotal NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    discount_amount NUMERIC(12, 2) DEFAULT 0.00,
    tax_amount NUMERIC(12, 2) DEFAULT 0.00,
    service_charge NUMERIC(12, 2) DEFAULT 0.00,
    total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),
    quantity INT NOT NULL DEFAULT 1,
    unit_price NUMERIC(12, 2) NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL,
    kitchen_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'cooking', 'ready', 'served'
    station VARCHAR(20) NOT NULL,                 -- 'kitchen' / 'barista'
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id),
    payment_method VARCHAR(50) NOT NULL, -- 'cash', 'qris', 'card', 'transfer'
    amount_paid NUMERIC(12, 2) NOT NULL,
    change_due NUMERIC(12, 2) DEFAULT 0.00,
    reference_number VARCHAR(100),
    status VARCHAR(20) DEFAULT 'completed',
    processed_by UUID REFERENCES users(id),
    paid_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### 4.6 Tabel HRIS & Penggajian
```sql
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    branch_id UUID NOT NULL REFERENCES branches(id),
    nik VARCHAR(50) UNIQUE NOT NULL,
    job_title VARCHAR(50) NOT NULL,
    basic_salary NUMERIC(12, 2) DEFAULT 0.00,
    employment_type VARCHAR(20) DEFAULT 'full_time',
    bank_account_number VARCHAR(50),
    bank_name VARCHAR(50),
    joined_date DATE NOT NULL
);

CREATE TABLE work_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id),
    name VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL
);

CREATE TABLE attendances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    shift_id UUID REFERENCES work_shifts(id),
    clock_in TIMESTAMPTZ NOT NULL,
    clock_out TIMESTAMPTZ,
    status VARCHAR(20) DEFAULT 'present',
    notes TEXT
);

CREATE TABLE payroll_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    basic_salary NUMERIC(12, 2) NOT NULL,
    incentives NUMERIC(12, 2) DEFAULT 0.00,
    deductions NUMERIC(12, 2) DEFAULT 0.00,
    net_salary NUMERIC(12, 2) NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    paid_at TIMESTAMPTZ
);
```

### 4.7 Tabel Keuangan & Buku Besar
```sql
CREATE TABLE chart_of_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    account_type VARCHAR(30) NOT NULL -- 'asset', 'liability', 'equity', 'revenue', 'expense'
);

CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id),
    reference_number VARCHAR(100) NOT NULL,
    entry_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE journal_entry_lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    journal_entry_id UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES chart_of_accounts(id),
    debit NUMERIC(14, 2) DEFAULT 0.00,
    credit NUMERIC(14, 2) DEFAULT 0.00
);
```

---

## 5. Roadmap Prioritas Pengembangan (Phase by Phase)

```
┌────────────────────────────────────────────────────────┐
│  FASE 1: Fondasi & Minimum Viable Product (POS Kasir) │
│  • Auth & User Multi-Branch                            │
│  • Kategori, Menu & Varian                             │
│  • Kasir POS (Order Dine-in & Bayar Tunai/QRIS)        │
│  • Laporan Penjualan Harian                            │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│  FASE 2: Kitchen KDS, Resep & Auto-Stok               │
│  • Kitchen Display System (Layar Barista & Dapur)      │
│  • Recipe BOM (Otomatis potong stok saat order)        │
│  • Master Stok Bahan Baku & Stock Opname               │
│  • Pembelian ke Supplier (Purchase Order)              │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│  FASE 3: Operasional Meja & HRIS                       │
│  • Table Management & Visual Floor Map                 │
│  • Penjadwalan Shift Barista & Absensi                 │
│  • Perhitungan Gaji & Lembur (Payroll)                 │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│  FASE 4: Finance, Akuntansi & Analitik Lanjutan        │
│  • Chart of Accounts & Auto-Journaling                 │
│  • Laba Rugi Real-Time (P&L)                           │
│  • Menu Matrix Analysis (Menu Terlaris & Profit Margin)│
└───────────────────────────┬────────────────────────────┘
```

---

## 6. Persyaratan Non-Fungsional (NFR)
1. **Kecepatan Transaksi**: Waktu respon checkout kasir dan pembuatan tiket pesanan harus `< 300ms`.
2. **Ketersediaan Offline-Ready (POS)**: POS kasir dirancang dengan caching lokal (IndexedDB) agar transaksi tetap bisa dicatat saat koneksi internet terputus sementara.
3. **Integritas Data**: Menggunakan transaksi basis data ACID PostgreSQL untuk menjamin pemotongan stok dan pencatatan pembayaran tidak mengalami *race condition*.
4. **Auditability**: Seluruh tindakan pembatalan pesanan (*void*) dan perubahan data stok wajib tercatat lengkap dengan identitas pengguna.
