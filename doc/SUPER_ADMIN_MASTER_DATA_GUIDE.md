# Panduan Master Data Hub (Super Admin CRUD)

Sistem Cafe ERP menyediakan modul khusus terpusat bagi **Super Admin** untuk mengelola seluruh data master sistem melalui antarmuka visual terpadu pada menu:
`http://localhost:5173/settings/master`

---

## 1. Daftar 13 Entitas Master CRUD

Modul Super Admin Master Hub mencakup 13 data master fundamental bisnis:

| No | Modul Master | Deskripsi Entitas | Endpoint API Backend | Operasi Didukung |
| :--- | :--- | :--- | :--- | :--- |
| 1 | **Branches (Cabang)** | Data outlet cafe, alamat, kontak, & status operasional | `/api/v1/master/branches` | Create, Read, Update, Delete |
| 2 | **Users (Pengguna)** | Akun pengguna sistem terhubung cabang & role | `/api/v1/master/users` | Create, Read, Update, Delete |
| 3 | **Roles & Permissions** | Hak akses peran sistem (Admin, Kasir, Barista, dll.) | `/api/v1/master/roles` | Create, Read, Update, Delete |
| 4 | **Categories (Kategori)** | Klasifikasi menu (Coffee, Non-Coffee, Food, Pastry) | `/api/v1/master/categories` | Create, Read, Update, Delete |
| 5 | **Products (Produk/Menu)** | Master item penjualan beserta harga jual & HPP | `/api/v1/master/products` | Create, Read, Update, Delete |
| 6 | **Inventory Items** | Bahan baku (biji kopi, susu, sirup, kemasan) | `/api/v1/master/inventory-items` | Create, Read, Update, Delete |
| 7 | **Tables (Meja)** | Nomor meja, kapasitas, dan alokasi zona cafe | `/api/v1/master/tables` | Create, Read, Update, Delete |
| 8 | **Accounts (COA)** | Bagan Akun Standar Akuntansi Keuangan | `/api/v1/master/accounts` | Create, Read, Update, Delete |
| 9 | **Suppliers (Pemasok)** | Vendor penyedia bahan baku dan nomor kontak | `/api/v1/master/suppliers` | Create, Read, Update, Delete |
| 10 | **Warehouses (Gudang)** | Lokasi penyimpanan stok per cabang | `/api/v1/master/warehouses` | Create, Read, Update, Delete |
| 11 | **Shifts (Jam Kerja)** | Pengaturan waktu kerja (Morning, Middle, Closing) | `/api/v1/master/shifts` | Create, Read, Update, Delete |
| 12 | **Departments** | Struktur departemen (Operations, Kitchen, HR, dll.) | `/api/v1/master/departments` | Create, Read, Update, Delete |
| 13 | **Positions (Jabatan)** | Jenjang jabatan karyawan & gaji pokok acuan | `/api/v1/master/positions` | Create, Read, Update, Delete |

---

## 2. Fitur Antarmuka Super Admin

Antarmuka Master Hub (`MasterDataHub.vue`) dirancang dengan standar UX enterprise:
1. **Sidebar Tab Switcher**: Navigasi cepat berpindah antar 13 entitas master tanpa berpindah halaman.
2. **Universal Live Search**: Pencarian real-time berdasarkan kode, nama, atau parameter spesifik.
3. **Modal Dialog Tambah & Edit**: Form validasi interaktif dengan deteksi tipe kolom (angka, teks, dropdown relasi, switch aktif/nonaktif).
4. **Delete Confirmation Safe-Guard**: Konfirmasi sebelum penghapusan data dengan perlindungan relasional (*Foreign Key Restrict*).
5. **Dynamic Badge & Status**: Visual status aktif/non-aktif, badge kategori, dan format mata uang Rupiah (`Rp xxx.xxx`).
6. **Integrasi Backend Otomatis**: Secara otomatis memanggil REST API backend saat server Go aktif, dan memiliki fallback data cerdas (*smart in-memory state*) saat server offline.

---

## 3. Struktur Endpoint API Backend Master

Implementasi backend berada di `backend/internal/delivery/http/handler/master_handler.go` dengan routing:

```go
masterGroup := v1.Group("/master")
{
    // Branches
    masterGroup.GET("/branches", masterHandler.GetBranches)
    masterGroup.POST("/branches", masterHandler.CreateBranch)
    masterGroup.PUT("/branches/:id", masterHandler.UpdateBranch)
    masterGroup.DELETE("/branches/:id", masterHandler.DeleteBranch)

    // Categories & Products
    masterGroup.GET("/categories", masterHandler.GetCategories)
    masterGroup.POST("/categories", masterHandler.CreateCategory)
    masterGroup.GET("/products", masterHandler.GetProducts)
    masterGroup.POST("/products", masterHandler.CreateProduct)

    // Inventory & Warehouses
    masterGroup.GET("/inventory-items", masterHandler.GetInventoryItems)
    masterGroup.GET("/warehouses", masterHandler.GetWarehouses)
    masterGroup.GET("/suppliers", masterHandler.GetSuppliers)

    // Tables & Floor
    masterGroup.GET("/tables", masterHandler.GetTables)
    masterGroup.POST("/tables", masterHandler.CreateTable)

    // HRIS Masters
    masterGroup.GET("/departments", masterHandler.GetDepartments)
    masterGroup.GET("/positions", masterHandler.GetPositions)
    masterGroup.GET("/shifts", masterHandler.GetShifts)

    // COA Accounts
    masterGroup.GET("/accounts", masterHandler.GetAccounts)
}
```
Setiap handler terhubung langsung ke transaksi basis data PostgreSQL melalui koneksi `*pgxpool.Pool`.

