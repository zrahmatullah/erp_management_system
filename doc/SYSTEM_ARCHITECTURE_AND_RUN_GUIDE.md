# Panduan Menjalankan Sistem Cafe ERP & Arsitektur

Dokumen ini menjelaskan cara menjalankan **Backend (Go Gin / REST API)** dan **Frontend (Vue 3 + Vite + Tailwind CSS)** secara lengkap melalui terminal Windows PowerShell / CMD / Bash.

---

## 1. Prasyarat Sistem (Prerequisites)

1. **Database PostgreSQL**:
   - Status: Berjalan lokal di port `5432`.
   - Nama Database: `cafe_erp`
   - Kredensial: User `postgres`, Password `samp3321` (berdasarkan konfigurasi file `.env_database`).
2. **Go (Golang)**: Versi 1.21 ke atas (`go version`).
3. **Node.js**: Versi 18 atau 20 LTS (`node -v`, `npm -v`).

---

## 2. Cara Menjalankan Backend di Terminal

### Langkah A: Pindah ke Direktori Backend
Buka terminal baru (PowerShell atau Command Prompt), lalu arahkan ke folder `backend`:
```powershell
cd d:\Project\cafe-erp-system\backend
```

### Langkah B: Migrasi Database (Jika Belum Dimigrasi)
Jika database baru dibuat atau perlu re-run migrasi dan seeder:
```powershell
go run cmd/migrate/main.go
```
*Output sukses:* `Database cafe_erp migrated and seeded successfully! 35 tables created.`

### Langkah C: Menjalankan Server Backend
Jalankan perintah berikut:
```powershell
go run cmd/server/main.go
```

**Informasi Port & Endpoint Backend:**
- **URL Base**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/health` (Mengembalikan `{"status":"ok","time":"..."}`)
- **API Master Hub**: `http://localhost:8080/api/v1/master/branches`, `categories`, `products`, dll.
- **Login Endpoint**: `POST http://localhost:8080/api/v1/auth/login`

> **Catatan Windows Service**: Backend telah dilengkapi proteksi *fail-safe* jika Redis belum diaktifkan di mesin lokal, backend akan otomatis mencatat peringatan dan tetap melayani seluruh traffic API REST via PostgreSQL secara optimal.

---

## 3. Cara Menjalankan Frontend di Terminal

### Langkah A: Buka Terminal Baru
Buka tab atau jendela terminal baru terpisah (jangan menutup terminal backend):
```powershell
cd d:\Project\cafe-erp-system\frontend
```

### Langkah B: Instalasi Dependensi (Jika Belum)
```powershell
npm install
```

### Langkah C: Menjalankan Server Development (Vite)
Jalankan server dev dengan perintah:
```powershell
npm run dev
```

**Informasi URL Frontend:**
- **Akses Aplikasi di Browser**: `http://localhost:5173` atau `http://127.0.0.1:5173`
- Mode Hot-Reload aktif otomatis saat Anda mengubah komponen Vue.

### Langkah D: Kompilasi / Build untuk Production
Untuk membuat paket distribusi production siap deploy:
```powershell
npm run build
```
File hasil kompilasi akan tersimpan di folder `frontend/dist/`.

---

## 4. Kredensial Login Sistem (Default Seed)

Gunakan akun berikut pada halaman login (`http://localhost:5173/login`):

| Peran (Role) | Username / Email | Password | Hak Akses Utama |
| :--- | :--- | :--- | :--- |
| **Super Admin** | `admin` / `admin@cafeerp.com` | `Admin@123` | Akses penuh ke seluruh menu & Super Admin Master Hub |
| **Store Manager** | `manager` / `manager@cafeerp.com` | `Manager@123` | POS, Inventory, KDS, Laporan Cabang |
| **Cashier Lead** | `cashier` / `cashier@cafeerp.com` | `Cashier@123` | POS Kasir & Transaksi Meja |
| **Barista / Kitchen** | `barista` / `barista@cafeerp.com` | `Barista@123` | Kitchen Display System (KDS) |
| **HR Specialist** | `hr` / `hr@cafeerp.com` | `HR@12345` | HRIS, Presensi, Jadwal Shift, Slip Gaji |
| **Finance Officer** | `finance` / `finance@cafeerp.com` | `Finance@123` | Financial Overview & Input Jurnal Umum |

> **Fitur Quick Demo Login**: Pada halaman login terdapat tombol *"1-Click Demo Login"* untuk langsung masuk ke masing-masing peran tanpa mengetik manual.

---

## 5. Arsitektur Komunikasi & Integrasi

```text
[ Browser / Klien (Port 5173) ]
       │
       │ HTTP / JSON API (CORS Allowed)
       ▼
[ Go Gin Web Server (Port 8080) ]
       ├── Auth Middleware (JWT HS256)
       ├── Master CRUD Handler (/api/v1/master/*)
       └── Database Driver (pgx/v5)
              │
              ▼
[ PostgreSQL Database (Port 5432 / cafe_erp) ]
  (35 Tabel Relasional, 60 Foreign Keys, Multi-Branch Schema)
```

