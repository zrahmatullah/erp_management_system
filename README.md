# ☕ Cafe ERP Management System

Enterprise Resource Planning (ERP) & Point of Sale (POS) modern terpadu untuk industri Food & Beverage (F&B) / Cafe & Restoran. Dibangun dengan performa tinggi menggunakan arsitektur **Clean Architecture REST API (Go)** dan antarmuka reaktif **Single Page Application (Vue 3 + Vite + Tailwind CSS)**.

---

## 📑 Daftar Isi
- [Fitur Utama](#-fitur-utama)
- [Teknologi yang Digunakan](#-teknologi-yang-digunakan)
- [Struktur Proyek](#-struktur-proyek)
- [Prasyarat Sistem](#-prasyarat-sistem)
- [Panduan Setup Lokal (Local Development)](#-panduan-setup-lokal-local-development)
  - [1. Konfigurasi Database PostgreSQL](#1-konfigurasi-database-postgresql)
  - [2. Menjalankan Backend API (Go)](#2-menjalankan-backend-api-go)
  - [3. Menjalankan Frontend Web (Vue 3)](#3-menjalankan-frontend-web-vue-3)
- [Kredensial Akun Default](#-kredensial-akun-default)
- [Panduan Deployment Produksi (Production Deployment)](#-panduan-deployment-produksi-production-deployment)
  - [Opsi A: Docker & Docker Compose (Rekomendasi)](#opsi-a-docker--docker-compose-rekomendasi)
  - [Opsi B: Linux VPS (Systemd + Nginx Reverse Proxy + SSL)](#opsi-b-linux-vps-systemd--nginx-reverse-proxy--ssl)
- [Keamanan & Pengecualian `.env`](#-keamanan--pengecualian-env)
- [Dokumentasi Lengkap Sistem](#-dokumentasi-lengkap-sistem)

---

## ✨ Fitur Utama

### 1. Point of Sale (POS) & Denah Meja Interaktif
- **Denah Meja Real-time (`/pos/tables`)**: Tata letak meja visual per lantai/zona dengan badge status real-time (*Available*, *Occupied*, *Billing*, *Reserved*).
- **Popup Detail Meja**: Klik meja untuk melihat item yang sedang dipesan, memantau status dapur (*Pending / Cooking / Served*), dan memesan menu tambahan langsung ke dapur.
- **Kasir 2-Tab (`/pos`)**:
  - *Tab 1*: Katalog menu cepat, pencarian menu, dan keranjang belanja (*Dine-in* / *Takeaway*).
  - *Tab 2*: Kasir Pembayaran Meja khusus untuk melunasi tagihan meja aktif dengan perhitungan durasi tamu, cetak bill sementara, dan integrasi modal pembayaran.
- **Pelunasan & Multi-Metode Pembayaran**: Cash (dengan kalkulator kembalian otomatis & quick cash pills), QRIS Dinamis, Transfer Bank, dan Debit/EDC. Meja otomatis dikosongkan saat pembayaran lunas.

### 2. Kitchen Display System (KDS) (`/kds`)
- Tiket pesanan dapur real-time dikelompokkan per stasiun (*Kitchen / Barista*).
- Update status per item: *Pending* ➔ *Cooking* ➔ *Served*.

### 3. Inventory & Manajemen Logistik
- **Daftar Stok Bahan Baku (`/inventory/stock`)**: Monitoring kuantitas bahan, satuan baku, dan indikator peringatan stok menipis (*Low Stock Warning*).
- **Purchase Order (PO) (`/inventory/purchase-orders`)**: Pengadaan barang ke supplier, approval hierarki, dan tracking penerimaan barang gudang.
- **Stock Opname (`/inventory/opname`)**: Audit fisik berkala vs sistem dengan pencatatan selisih & alasan adjustment.

### 4. HRIS & Sistem Penggajian (Payroll)
- **Direktori Staf (`/hris/employees`)**: Database karyawan, departemen, jabatan, dan riwayat kerja.
- **Presensi & Jam Kerja (`/hris/attendance`)**: Log absensi harian, tracking keterlambatan menit, dan jam lembur.
- **Jadwal Shift Roster (`/hris/shifts`)**: Jadwal mingguan staf terintegrasi langsung ke master shift.
- **Pengajuan Cuti (`/hris/leave`)**: Workflow permohonan dan persetujuan cuti.
- **Payroll Wizard 5-Langkah (`/hris/payroll`)**: Perhitungan gaji pokok, tunjangan, uang lembur, potongan PPh21 & BPJS, dan cetak slip gaji resmi (*Payslip*).

### 5. Keuangan & Akuntansi (Finance)
- **Chart of Accounts (COA)**: Standar bagan akun 5 digit.
- **Jurnal Umum Berimbang (`/finance/journal`)**: Input transaksi debit/kredit multi-baris otomatis seimbang.
- **Ringkasan Arus Kas (`/finance`)**: Visualisasi pendapatan vs pengeluaran, net profit, dan riwayat jurnal.

### 6. Super Admin Master Data CRUD Hub (`/settings/master`)
- Antarmuka terpusat untuk mengelola **13 entitas master data**: Cabang, User, Role & Hak Akses Matrix, Kategori Menu, Produk & Resep, Bahan Baku, Meja & Zona, Akun COA, Pemasok (Supplier), Gudang, Shift Kerja, Departemen, dan Jabatan.

---

## 🛠 Teknologi yang Digunakan

| Lapisan | Teknologi | Versi | Deskripsi |
| :--- | :--- | :--- | :--- |
| **Backend** | **Go (Golang)** | `1.22+` | Bahasa pemrograman backend cepat dan hemat memori |
| **Router** | **Chi Router** | `v5` | Router HTTP idiomatik, ringan, dan kompatibel dengan net/http |
| **Database Driver** | **pgx** | `v5` | Driver PostgreSQL native berkinerja tinggi dengan pooling |
| **Database** | **PostgreSQL** | `15 / 16+` | Basis data relasional dengan 35 tabel & 60 Foreign Keys |
| **Frontend** | **Vue.js** | `3.4+` | Framework UI progresif berbasis Composition API (`<script setup>`) |
| **Build Tool** | **Vite** | `5.4+` | Bundler kilat untuk hot module replacement (HMR) dan build produksi |
| **Styling** | **Tailwind CSS** | `3.4+` | Utility-first CSS framework |
| **Icons** | **Lucide Vue Next** | `latest` | Library icon SVG modern dan konsisten |
| **State Manager** | **Pinia** | `2.1+` | Manajemen state terpusat (Auth, Notifikasi Toast, Dialog) |

---

## 📁 Struktur Proyek

```
cafe-erp-system/
├── .env.example             # Template konfigurasi environment root
├── .gitignore               # Pengecualian ketat berkas sensitif (.env, node_modules, binaries)
├── README.md                # Dokumentasi utama proyek & panduan instalasi
├── Mockups_UI/              # 20 acuan visual desain UI sistem
├── doc/                     # Dokumentasi arsitektur dan panduan teknis mendalam
│   ├── PRD_Cafe_ERP_System.md
│   ├── SYSTEM_ARCHITECTURE_AND_RUN_GUIDE.md
│   ├── DATABASE_SCHEMA_AND_RELATIONS.md
│   ├── SUPER_ADMIN_MASTER_DATA_GUIDE.md
│   ├── AUTH_AND_UI_STANDARDIZATION_GUIDE.md
│   ├── MODAL_AND_DIALOG_STANDARDIZATION_GUIDE.md
│   ├── POS_TABLE_ORDER_AND_BILLING_GUIDE.md
│   └── DEPLOYMENT_GUIDE.md
├── backend/                 # Source code Backend Go REST API
│   ├── cmd/server/main.go   # Entrypoint aplikasi server Go
│   ├── internal/
│   │   ├── config/          # Loader konfigurasi environment
│   │   ├── delivery/http/   # Router, Handlers (Master & Operational)
│   │   ├── repository/      # Layer query PostgreSQL (pgx)
│   │   └── usecase/         # Logika bisnis sistem
│   ├── migrations/          # Berkas skrip SQL DDL & Seeder
│   │   ├── 000001_full_cafe_erp_schema.sql
│   │   ├── 000002_complete_seed_data.sql
│   │   └── 000003_operational_seed_data.sql
│   └── go.mod, go.sum
└── frontend/                # Source code Frontend Vue 3 + Vite
    ├── src/
    │   ├── components/      # Common UI (AppModal, ToastContainer, GlobalDialog)
    │   ├── layouts/         # AppLayout, Sidebar, AppHeader
    │   ├── stores/          # Pinia stores (auth, notification, dialog)
    │   ├── views/           # Halaman SPA (POS, Tables, KDS, HRIS, Finance, Settings)
    │   └── main.ts, App.vue
    ├── package.json
    └── vite.config.ts
```

---

## 📋 Prasyarat Sistem

Sebelum memulai, pastikan perangkat atau server Anda telah terpasang:
- **Git** (versi 2.30+)
- **Go** (versi 1.22 atau lebih baru)
- **Node.js** (versi 18 LTS atau 20 LTS) & **npm** (versi 9+)
- **PostgreSQL** (versi 15 atau 16)
- **Terminal Shell**: Bash / Zsh (Linux/macOS) atau PowerShell (Windows)

---

## 🚀 Panduan Setup Lokal (Local Development)

### 1. Konfigurasi Database PostgreSQL
1. Buat database baru bernama `cafe_erp` di instance PostgreSQL lokal:
   ```sql
   CREATE DATABASE cafe_erp;
   ```
2. Jalankan skrip migrasi skema tabel dan data seed berurutan:
   ```bash
   # Di Linux / macOS:
   PGPASSWORD='your_password' psql -U postgres -d cafe_erp -f backend/migrations/000001_full_cafe_erp_schema.sql
   PGPASSWORD='your_password' psql -U postgres -d cafe_erp -f backend/migrations/000002_complete_seed_data.sql
   PGPASSWORD='your_password' psql -U postgres -d cafe_erp -f backend/migrations/000003_operational_seed_data.sql

   # Di Windows (PowerShell):
   $env:PGPASSWORD='your_password'
   psql -U postgres -d cafe_erp -f backend/migrations/000001_full_cafe_erp_schema.sql
   psql -U postgres -d cafe_erp -f backend/migrations/000002_complete_seed_data.sql
   psql -U postgres -d cafe_erp -f backend/migrations/000003_operational_seed_data.sql
   ```

---

### 2. Menjalankan Backend API (Go)
1. Salin berkas konfigurasi `.env`:
   ```bash
   cd backend
   cp .env.example .env
   ```
2. Sesuaikan kredensial koneksi database pada `backend/.env`:
   ```env
   DB_HOST=127.0.0.1
   DB_PORT=5432
   DB_DATABASE=cafe_erp
   DB_USERNAME=postgres
   DB_PASSWORD=your_password
   APP_PORT=8080
   JWT_SECRET=rahasia-kunci-jwt-minimal-32-karakter
   ```
3. Unduh dependensi Go dan jalankan server:
   ```bash
   go mod download
   go run cmd/server/main.go
   ```
   *Server backend akan aktif di: `http://localhost:8080` (Health check: `http://localhost:8080/health`)*

---

### 3. Menjalankan Frontend Web (Vue 3)
1. Buka terminal baru, masuk ke direktori `frontend`:
   ```bash
   cd frontend
   cp .env.example .env
   ```
2. Pasang paket dependensi npm:
   ```bash
   npm install
   ```
3. Jalankan development server:
   ```bash
   npm run dev
   ```
   *Aplikasi web akan dapat diakses melalui browser di: `http://localhost:5173`*

---

## 🔑 Kredensial Akun Default

Sistem menyediakan akun demo untuk berbagai tingkatan hak akses (Password default untuk seluruh akun demo: **`Admin@123`**):

| Role | Username | Email | Hak Akses Utama |
| :--- | :--- | :--- | :--- |
| **Super Admin** | `admin` | `admin@cafe-erp.com` | Akses penuh ke seluruh fitur & Master Data Hub |
| **Store Manager**| `manager` | `manager@cafe-erp.com` | Manajemen operasional cabang, laporan, & inventory |
| **Cashier (Kasir)**| `kasir` | `kasir@cafe-erp.com` | POS Meja, Kasir Pembayaran (Billing), & Transaksi |
| **Barista / Kitchen**| `barista` | `barista@cafe-erp.com` | Kitchen Display System (KDS) & Order Status |
| **Warehouse Staff**| `gudang` | `gudang@cafe-erp.com` | Stok bahan baku, Purchase Order, & Opname |

> **Tips Cepat**: Pada halaman login (`/login`), Anda juga dapat mengklik salah satu tombol **"1-Click Demo Login"** untuk langsung masuk sebagai role yang dipilih tanpa perlu mengetik manual.

---

## 🌐 Panduan Deployment Produksi (Production Deployment)

### Opsi A: Docker & Docker Compose (Rekomendasi)

1. Salin berkas environment di root proyek:
   ```bash
   cp .env.example .env
   ```
2. Sesuaikan konfigurasi `DB_PASSWORD` dan `JWT_SECRET` pada `.env`.
3. Jalankan seluruh container (PostgreSQL, Backend API, Frontend Nginx):
   ```bash
   docker compose up -d --build
   ```
4. Cek status container:
   ```bash
   docker compose ps
   docker compose logs -f
   ```

---

### Opsi B: Linux VPS (Systemd + Nginx Reverse Proxy + SSL)

#### 1. Setup Backend sebagai Linux Daemon (Systemd)
1. Kompilasi binary Go untuk arsitektur Linux:
   ```bash
   cd backend
   CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/server cmd/server/main.go
   chmod +x bin/server
   ```
2. Buat file service Systemd di `/etc/systemd/system/cafe-erp.service`:
   ```ini
   [Unit]
   Description=Cafe ERP Backend Service
   After=network.target postgresql.service

   [Service]
   Type=simple
   User=www-data
   WorkingDirectory=/var/www/cafe-erp-system/backend
   ExecStart=/var/www/cafe-erp-system/backend/bin/server
   Restart=always
   RestartSec=5s
   EnvironmentFile=/var/www/cafe-erp-system/backend/.env

   [Install]
   WantedBy=multi-user.target
   ```
3. Aktifkan service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now cafe-erp
   ```

#### 2. Build Frontend untuk Produksi
```bash
cd frontend
npm ci
npm run build
# Hasil kompilasi siap disajikan berada pada folder: frontend/dist
```

#### 3. Konfigurasi Nginx & Let's Encrypt SSL
1. Buat konfigurasi server block di `/etc/nginx/sites-available/cafe-erp`:
   ```nginx
   server {
       listen 80;
       server_name cafe-erp.domainanda.com;

       root /var/www/cafe-erp-system/frontend/dist;
       index index.html;

       location / {
           try_files $uri $uri/ /index.html;
       }

       location /api/v1/ {
           proxy_pass http://127.0.0.1:8080/api/v1/;
           proxy_http_version 1.1;
           proxy_set_header Upgrade $http_upgrade;
           proxy_set_header Connection 'upgrade';
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }

       location /health {
           proxy_pass http://127.0.0.1:8080/health;
       }
   }
   ```
2. Aktifkan konfigurasi dan reload Nginx:
   ```bash
   sudo ln -s /etc/nginx/sites-available/cafe-erp /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx
   ```
3. Pasang sertifikat SSL gratis via Certbot:
   ```bash
   sudo certbot --nginx -d cafe-erp.domainanda.com
   ```

---

## 🔒 Keamanan & Pengecualian `.env`

Repository ini menerapkan aturan ketat untuk mencegah kebocoran informasi kredensial:
- Seluruh berkas `.env`, `.env.*`, dan `.env_database` telah dimasukkan ke dalam [`.gitignore`](file:///d:/Project/cafe-erp-system/.gitignore).
- Jangan pernah mengunggah kredensial database asli, secret key JWT, atau password server ke remote repository.
- Selalu gunakan berkas template [`.env.example`](file:///d:/Project/cafe-erp-system/.env.example) untuk membagikan referensi variabel lingkungan.

---

## 📚 Dokumentasi Lengkap Sistem

Dokumentasi arsitektur sistem, skema basis data, dan modul fungsional tersedia di folder [`doc/`](file:///d:/Project/cafe-erp-system/doc/):

- [`doc/OWASP_SECURITY_STANDARDIZATION.md`](file:///d:/Project/cafe-erp-system/doc/OWASP_SECURITY_STANDARDIZATION.md) : Standarisasi dan kepatuhan keamanan sistem berdasarkan OWASP Top 10 (2021) & ASVS v4.0.
- [`doc/DEBUGGING_AND_BUGFIX_REPORT.md`](file:///d:/Project/cafe-erp-system/doc/DEBUGGING_AND_BUGFIX_REPORT.md) : Laporan lengkap audit sistem, perbaikan bug inventory, PO auto-receipt, stock opname, dan jurnal umum.
- [`doc/DEPLOYMENT_GUIDE.md`](file:///d:/Project/cafe-erp-system/doc/DEPLOYMENT_GUIDE.md) : Panduan komprehensif deployment produksi (Docker & VPS Systemd/Nginx).
- [`doc/POS_TABLE_ORDER_AND_BILLING_GUIDE.md`](file:///d:/Project/cafe-erp-system/doc/POS_TABLE_ORDER_AND_BILLING_GUIDE.md) : Panduan alur POS denah meja interaktif & sistem kasir 2-tab.
- [`doc/MODAL_AND_DIALOG_STANDARDIZATION_GUIDE.md`](file:///d:/Project/cafe-erp-system/doc/MODAL_AND_DIALOG_STANDARDIZATION_GUIDE.md) : Standarisasi UI komponen AppModal dan Global Dialog reaktif.
- [`doc/DATABASE_SCHEMA_AND_RELATIONS.md`](file:///d:/Project/cafe-erp-system/doc/DATABASE_SCHEMA_AND_RELATIONS.md) : Diagram ERD Mermaid, relasi 35 tabel, dan 60 Foreign Key constraints.
- [`doc/SUPER_ADMIN_MASTER_DATA_GUIDE.md`](file:///d:/Project/cafe-erp-system/doc/SUPER_ADMIN_MASTER_DATA_GUIDE.md) : Panduan operasional CRUD untuk 13 entitas master data.
- [`doc/AUTH_AND_UI_STANDARDIZATION_GUIDE.md`](file:///d:/Project/cafe-erp-system/doc/AUTH_AND_UI_STANDARDIZATION_GUIDE.md) : Dokumentasi autentikasi JWT, sistem toast notifikasi, dan ikon Lucide.
- [`doc/UI_STANDARDIZATION_AND_MOCKUP_MAPPING.md`](file:///d:/Project/cafe-erp-system/doc/UI_STANDARDIZATION_AND_MOCKUP_MAPPING.md) : Pemetaan 20 visual mockups ke komponen antarmuka Vue 3.

---

## 📄 Lisensi
Hak Cipta © 2026 Cafe ERP Management System. Seluruh hak cipta dilindungi.
