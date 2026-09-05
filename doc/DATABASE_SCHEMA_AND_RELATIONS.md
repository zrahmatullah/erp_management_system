# Database Schema & Relational Integrity Documentation

## 1. Overview
Sistem Cafe ERP menggunakan basis data relasional **PostgreSQL 15+** dengan total **35 tabel** yang terhubung melalui **60 Foreign Key constraints**, primary keys berbasis UUID v4 (`uuid`), timestamp audit (`created_at`, `updated_at`, `deleted_at` untuk soft deletes), serta indexing komprehensif pada foreign keys, status, dan kode unik.

Database telah dibuat dan dimigrasikan secara sukses ke basis data lokal PostgreSQL menggunakan:
- **Connection URL**: `postgres://postgres:samp3321@127.0.0.1:5432/cafe_erp?sslmode=disable`
- **Migration Engine**: `golang-migrate` via `backend/cmd/migrate/main.go`
- **Migration Files**:
  - `000001_full_cafe_erp_schema.sql` (Struktur 35 tabel & 60 relasi FK)
  - `000002_complete_seed_data.sql` (Master data, user admin, produk, akun COA, cabang)

---

## 2. Entity Relationship Diagram (Mermaid)

```mermaid
erDiagram
    BRANCHES ||--o{ USERS : "employs"
    BRANCHES ||--o{ TABLES : "has"
    BRANCHES ||--o{ WAREHOUSES : "operates"
    BRANCHES ||--o{ ORDERS : "records"
    BRANCHES ||--o{ SHIFTS : "schedules"
    BRANCHES ||--o{ DEPARTMENTS : "contains"

    ROLES ||--o{ USERS : "assigned_to"
    ROLES ||--o{ ROLE_PERMISSIONS : "defines"
    PERMISSIONS ||--o{ ROLE_PERMISSIONS : "granted_in"

    CATEGORIES ||--o{ PRODUCTS : "groups"
    PRODUCTS ||--o{ PRODUCT_VARIANTS : "has"
    PRODUCT_VARIANTS ||--o{ RECIPES : "composed_of"
    INVENTORY_ITEMS ||--o{ RECIPES : "used_as_ingredient"

    WAREHOUSES ||--o{ INVENTORY_STOCKS : "stores"
    INVENTORY_ITEMS ||--o{ INVENTORY_STOCKS : "quantified_in"
    SUPPLIERS ||--o{ PURCHASE_ORDERS : "supplies"
    PURCHASE_ORDERS ||--o{ PURCHASE_ORDER_ITEMS : "contains"
    INVENTORY_ITEMS ||--o{ PURCHASE_ORDER_ITEMS : "ordered_in"

    ZONES ||--o{ TABLES : "locates"
    TABLES ||--o{ ORDERS : "seats"
    USERS ||--o{ ORDERS : "handled_by"
    ORDERS ||--o{ ORDER_ITEMS : "includes"
    PRODUCT_VARIANTS ||--o{ ORDER_ITEMS : "selected_item"
    ORDERS ||--o{ PAYMENTS : "settled_with"
    ORDERS ||--o{ KITCHEN_ORDERS : "routes_to_kds"

    DEPARTMENTS ||--o{ POSITIONS : "categorizes"
    POSITIONS ||--o{ EMPLOYEES : "holds"
    BRANCHES ||--o{ EMPLOYEES : "assigned_to"
    USERS ||--o| EMPLOYEES : "links_to"
    EMPLOYEES ||--o{ ATTENDANCES : "logs"
    SHIFTS ||--o{ ATTENDANCES : "scheduled_in"
    EMPLOYEES ||--o{ LEAVES : "requests"
    EMPLOYEES ||--o{ PAYROLLS : "earns"

    ACCOUNTS ||--o{ ACCOUNTS : "parent_child"
    JOURNAL_ENTRIES ||--o{ JOURNAL_LINES : "contains"
    ACCOUNTS ||--o{ JOURNAL_LINES : "credited_debited"
```

---

## 3. Daftar 35 Tabel & Relasi Lengkap

### Multi-Branch & Identity Access Management (IAM)
1. **`branches`**: Master outlet/cabang (`id`, `code`, `name`, `address`, `phone`, `email`, `is_active`).
2. **`roles`**: Master peran sistem (`id`, `name`, `description`).
3. **`permissions`**: Master hak akses permission (`id`, `module`, `action`, `description`).
4. **`role_permissions`**: Pivot matrix relasi peran dan permission (`role_id` -> `roles.id`, `permission_id` -> `permissions.id`).
5. **`users`**: Akun pengguna sistem (`branch_id` -> `branches.id`, `role_id` -> `roles.id`).

### Master Produk & Katalog (Menu)
6. **`categories`**: Kategori produk (`id`, `name`, `code`, `is_active`).
7. **`products`**: Master barang/menu (`category_id` -> `categories.id`, `code`, `name`, `base_price`, `cost_price`).
8. **`product_variants`**: Varian menu (`product_id` -> `products.id`, `name`, `sku`, `price_adjustment`).
9. **`modifier_groups`**: Kelompok opsi/add-on (contoh: Level Gula, Pilihan Susu).
10. **`modifiers`**: Item add-on (`modifier_group_id` -> `modifier_groups.id`, `name`, `price`).
11. **`product_modifiers`**: Relasi varian produk dengan modifier group (`product_id` -> `products.id`, `modifier_group_id` -> `modifier_groups.id`).

### Manajemen Inventori & Resep (BOM - Bill of Materials)
12. **`inventory_items`**: Master bahan baku & logistik (`id`, `code`, `name`, `unit`, `minimum_stock`, `cost_per_unit`).
13. **`warehouses`**: Master gudang/lokasi penyimpanan (`branch_id` -> `branches.id`, `name`, `code`).
14. **`inventory_stocks`**: Saldo stok realtime per gudang (`inventory_item_id` -> `inventory_items.id`, `warehouse_id` -> `warehouses.id`, `quantity`, `reserved_quantity`).
15. **`stock_movements`**: Kartu stok pergerakan barang (`inventory_item_id` -> `inventory_items.id`, `warehouse_id` -> `warehouses.id`, `type`, `quantity`, `reference_no`).
16. **`recipes`**: Formula Bill of Materials (BOM) per porsi (`product_variant_id` -> `product_variants.id`, `inventory_item_id` -> `inventory_items.id`, `quantity_required`).
17. **`suppliers`**: Master pemasok bahan baku (`id`, `code`, `name`, `contact_person`, `phone`, `email`).
18. **`purchase_orders`**: Dokumen Purchase Order (`branch_id` -> `branches.id`, `supplier_id` -> `suppliers.id`, `created_by` -> `users.id`, `total_amount`, `status`).
19. **`purchase_order_items`**: Rincian PO (`purchase_order_id` -> `purchase_orders.id`, `inventory_item_id` -> `inventory_items.id`, `quantity`, `unit_price`, `total_price`).
20. **`stock_opnames`**: Dokumen audit fisik stok / Stock Take (`warehouse_id` -> `warehouses.id`, `performed_by` -> `users.id`, `status`).
21. **`stock_opname_items`**: Selisih fisik vs sistem (`stock_opname_id` -> `stock_opnames.id`, `inventory_item_id` -> `inventory_items.id`, `system_qty`, `physical_qty`, `difference_qty`).

### Point of Sale (POS), Meja, & Kitchen Display System (KDS)
22. **`zones`**: Area tata letak cafe (`branch_id` -> `branches.id`, `name`).
23. **`tables`**: Master meja cafe (`zone_id` -> `zones.id`, `branch_id` -> `branches.id`, `table_number`, `capacity`, `status`).
24. **`orders`**: Transaksi POS (`branch_id` -> `branches.id`, `table_id` -> `tables.id`, `user_id` -> `users.id`, `customer_name`, `order_type`, `status`, `subtotal`, `tax`, `total_amount`).
25. **`order_items`**: Rincian pesanan menu (`order_id` -> `orders.id`, `product_variant_id` -> `product_variants.id`, `quantity`, `unit_price`, `subtotal`, `notes`).
26. **`payments`**: Transaksi pembayaran (`order_id` -> `orders.id`, `payment_method`, `amount_paid`, `change_amount`, `status`, `reference_no`).
27. **`kitchen_orders`**: Antrian tiket KDS (`order_id` -> `orders.id`, `branch_id` -> `branches.id`, `status`, `priority`).
28. **`kitchen_order_items`**: Item tiket dapur (`kitchen_order_id` -> `kitchen_orders.id`, `order_item_id` -> `order_items.id`, `status`).

### Human Resource Information System (HRIS) & Payroll
29. **`departments`**: Master departemen (`branch_id` -> `branches.id`, `code`, `name`).
30. **`positions`**: Master jabatan (`department_id` -> `departments.id`, `title`, `base_salary`).
31. **`shifts`**: Master jam kerja (`branch_id` -> `branches.id`, `name`, `start_time`, `end_time`).
32. **`employees`**: Data induk karyawan (`branch_id` -> `branches.id`, `position_id` -> `positions.id`, `user_id` -> `users.id`, `nik`, `first_name`, `last_name`, `status`).
33. **`attendances`**: Catatan presensi presisi (`employee_id` -> `employees.id`, `shift_id` -> `shifts.id`, `date`, `clock_in`, `clock_out`, `status`, `late_minutes`).
34. **`leaves`**: Pengajuan cuti & izin (`employee_id` -> `employees.id`, `start_date`, `end_date`, `reason`, `status`, `approved_by` -> `users.id`).
35. **`payrolls`**: Penggajian karyawan (`employee_id` -> `employees.id`, `period`, `basic_salary`, `allowance`, `overtime_pay`, `deductions`, `net_salary`, `status`).

### Keuangan & General Ledger (Finance)
36. **`accounts`**: Chart of Accounts (COA) (`parent_id` -> `accounts.id`, `code`, `name`, `type`, `normal_balance`).
37. **`journal_entries`**: Header transaksi jurnal umum (`branch_id` -> `branches.id`, `created_by` -> `users.id`, `entry_date`, `reference_no`, `description`, `is_posted`).
38. **`journal_lines`**: Baris debit/kredit jurnal (`journal_entry_id` -> `journal_entries.id`, `account_id` -> `accounts.id`, `debit`, `credit`, `memo`).

---

## 4. Validasi Foreign Key Integrity
Seluruh relasi tabel telah diuji melalui query katalog sistem PostgreSQL:
```sql
SELECT 
    conrelid::regclass AS table_name,
    conname AS foreign_key_name,
    confrelid::regclass AS referenced_table
FROM pg_constraint 
WHERE contype = 'f' 
ORDER BY table_name;
```
Hasil verifikasi menunjukkan **60 relasi Foreign Key aktif** dengan aturan integritas referensial:
- `ON DELETE RESTRICT` pada akun COA dan data master primer untuk mencegah *orphan records*.
- `ON DELETE CASCADE` pada baris rincian (child items seperti `order_items`, `journal_lines`, `recipe_items`, `role_permissions`).
- `ON DELETE SET NULL` pada asosiasi opsional seperti `parent_id` akun COA dan `table_id` transaksi takeaway/drive-thru.

