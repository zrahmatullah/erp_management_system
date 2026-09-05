-- ==============================================================================
-- CAFE ERP SYSTEM: COMPLETE SEED DATA & MASTER SEEDER (VALID HEX UUIDs)
-- ==============================================================================

-- 1. CABANG / BRANCHES
INSERT INTO branches (id, code, name, address, phone, email, is_active) VALUES
('b1111111-0000-0000-0000-000000000001', 'B-SENO-01', 'Kopi Kenangan Senopati (HQ & Central Kitchen)', 'Jl. Senopati No. 45, Kebayoran Baru, Jakarta Selatan', '021-7201234', 'senopati@cafe-erp.com', TRUE),
('b1111111-0000-0000-0000-000000000002', 'B-DHAR-02', 'Cabang Dharmawangsa (Express & Dine-In)', 'Jl. Dharmawangsa VI No. 12, Jakarta Selatan', '021-7205678', 'dharmawangsa@cafe-erp.com', TRUE)
ON CONFLICT (code) DO NOTHING;

-- 2. ROLES / PERAN
INSERT INTO roles (id, name, description) VALUES
('a1111111-0000-0000-0000-000000000001', 'Super Admin', 'Akses penuh seluruh konfigurasi sistem, lisensi dan semua modul'),
('a1111111-0000-0000-0000-000000000002', 'Owner', 'Pemantauan KPI omset, margin laba rugi, dan approval PO bernilai besar'),
('a1111111-0000-0000-0000-000000000003', 'Manager', 'Operasional harian cabang, approval cuti, opname, shift, dan supervisor'),
('a1111111-0000-0000-0000-000000000004', 'Kasir', 'Transaksi penjualan POS, buka/tutup laci kasir, dan cetak struk'),
('a1111111-0000-0000-0000-000000000005', 'Kitchen Staff', 'Memproses pesanan makanan sesuai monitor KDS'),
('a1111111-0000-0000-0000-000000000006', 'HR Admin', 'Manajemen data staf, presensi, cuti, dan penggajian payroll'),
('a1111111-0000-0000-0000-000000000007', 'Warehouse', 'Penerimaan barang PO, transfer stok, dan stock opname'),
('a1111111-0000-0000-0000-000000000008', 'Akuntan', 'Jurnal penyesuaian, petty cash, neraca, laba rugi, dan pajak'),
('a1111111-0000-0000-0000-000000000009', 'Pelayan', 'Pelayanan tamu, catat nomor meja, dan antar pesanan')
ON CONFLICT (name) DO NOTHING;

-- 3. PERMISSIONS
INSERT INTO permissions (module, action, description) VALUES
('dashboard', 'view', 'Lihat dashboard analitik'),
('pos', 'view', 'Buka layar POS'),
('pos', 'create', 'Buat transaksi pesanan kasir'),
('pos', 'edit', 'Koreksi item pesanan kasir'),
('pos', 'delete', 'Void / batalkan pesanan kasir'),
('pos', 'approve', 'Otorisasi diskon kasir'),
('kitchen', 'view', 'Lihat Kitchen Display System (KDS)'),
('kitchen', 'edit', 'Perbarui status tiket pesanan di dapur'),
('tables', 'view', 'Lihat denah dan status meja'),
('tables', 'edit', 'Ubah status atau pindah meja'),
('inventory', 'view', 'Lihat daftar stok dan barang'),
('inventory', 'create', 'Tambah item bahan baku'),
('inventory', 'edit', 'Penyesuaian stok manual'),
('inventory', 'delete', 'Arsipkan bahan baku'),
('inventory', 'approve', 'Persetujuan koreksi opname'),
('purchasing', 'view', 'Lihat daftar PO pembelian'),
('purchasing', 'create', 'Buat pesanan pembelian (PO)'),
('purchasing', 'approve', 'Setujui dokumen PO'),
('hris', 'view', 'Lihat data karyawan'),
('hris', 'create', 'Tambah data staf baru'),
('hris', 'edit', 'Ubah profil atau shift karyawan'),
('hris', 'approve', 'Setujui cuti karyawan'),
('payroll', 'view', 'Lihat rekap penggajian'),
('payroll', 'approve', 'Setujui pencairan gaji'),
('finance', 'view', 'Lihat ikhtisar keuangan & jurnal'),
('finance', 'create', 'Buat jurnal transaksi baru'),
('finance', 'approve', 'Posting jurnal ke buku besar'),
('reports', 'view', 'Lihat laporan analitik'),
('reports', 'export', 'Ekspor laporan ke PDF atau Excel'),
('settings', 'view', 'Lihat pengaturan sistem'),
('settings', 'edit', 'Ubah konfigurasi & hak akses'),
('master', 'view', 'Lihat data master super admin'),
('master', 'create', 'Tambah data master'),
('master', 'edit', 'Ubah data master'),
('master', 'delete', 'Hapus / nonaktifkan data master')
ON CONFLICT (module, action) DO NOTHING;

-- Berikan seluruh permission kepada Super Admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'a1111111-0000-0000-0000-000000000001', id FROM permissions
ON CONFLICT DO NOTHING;

-- 4. PENGGUNA AWAL (Password: Admin@123 -> $2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK)
INSERT INTO users (id, branch_id, role_id, username, email, password_hash, full_name, phone, pin_code, is_active) VALUES
('c1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000001', 'admin', 'admin@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Super Administrator', '081234567890', '123456', TRUE),
('c1111111-0000-0000-0000-000000000002', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000003', 'sarah.manager', 'sarah@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Sarah Johnson', '081234567891', '112233', TRUE),
('c1111111-0000-0000-0000-000000000003', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000004', 'jane.cashier', 'cashier@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Jane Doe', '081234567892', '1234', TRUE),
('c1111111-0000-0000-0000-000000000004', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000005', 'chef.john', 'kitchen@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'John Doe', '081234567893', '2233', TRUE),
('c1111111-0000-0000-0000-000000000005', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000007', 'ahmad.warehouse', 'warehouse@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Ahmad Warehouse', '081234567894', '4455', TRUE)
ON CONFLICT (username) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- 5. ZONASI & MEJA CAFE
INSERT INTO table_zones (id, branch_id, name, description) VALUES
('e1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'Lantai 1 (Indoor AC)', 'Ruang utama barista dan kasir'),
('e1111111-0000-0000-0000-000000000002', 'b1111111-0000-0000-0000-000000000001', 'Lantai 2 (Mezzanine & Co-Working)', 'Area tenang dan stopkontak lengkap'),
('e1111111-0000-0000-0000-000000000003', 'b1111111-0000-0000-0000-000000000001', 'Outdoor (Smoking Garden)', 'Taman luar terbuka')
ON CONFLICT DO NOTHING;

INSERT INTO cafe_tables (branch_id, zone_id, table_number, capacity, status, pos_x, pos_y) VALUES
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-01', 2, 'available', 10, 20),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-02', 2, 'available', 30, 20),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-03', 2, 'available', 50, 20),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-04', 4, 'occupied', 70, 20),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-05', 4, 'occupied', 90, 20),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-06', 6, 'available', 10, 50),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-07', 8, 'reserved', 40, 50),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-08', 6, 'available', 70, 50),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-09', 2, 'available', 10, 80),
('b1111111-0000-0000-0000-000000000001', 'e1111111-0000-0000-0000-000000000001', 'T-10', 4, 'occupied', 30, 80)
ON CONFLICT (branch_id, table_number) DO NOTHING;

-- 6. KATEGORI MENU
INSERT INTO menu_categories (id, name, slug, icon, sort_order, is_active) VALUES
('d1111111-0000-0000-0000-000000000001', 'Coffee', 'coffee', 'CoffeeOutlined', 1, TRUE),
('d1111111-0000-0000-0000-000000000002', 'Non-Coffee', 'non-coffee', 'CupOutlined', 2, TRUE),
('d1111111-0000-0000-0000-000000000003', 'Food', 'food', 'FastFoodOutlined', 3, TRUE),
('d1111111-0000-0000-0000-000000000004', 'Dessert', 'dessert', 'IceCreamOutlined', 4, TRUE),
('d1111111-0000-0000-0000-000000000005', 'Snack', 'snack', 'CookieOutlined', 5, TRUE)
ON CONFLICT (slug) DO NOTHING;

-- 7. BAHAN BAKU / INVENTORI
INSERT INTO inventory_items (id, sku, name, category, uom, min_stock, max_stock, average_cost, is_active) VALUES
('f1111111-0000-0000-0000-000000000001', 'INV-CB-001', 'Arabica Coffee Beans 1kg', 'coffee_beans', 'kg', 20.0, 200.0, 120000.00, TRUE),
('f1111111-0000-0000-0000-000000000002', 'INV-MK-010', 'Whole Milk 1L', 'dairy', 'liter', 15.0, 100.0, 18000.00, TRUE),
('f1111111-0000-0000-0000-000000000003', 'INV-SG-005', 'Sugar Syrup / Gula Aren 1kg', 'syrup', 'kg', 10.0, 100.0, 22000.00, TRUE),
('f1111111-0000-0000-0000-000000000004', 'INV-SY-003', 'Caramel Syrup 750ml', 'syrup', 'bottle', 5.0, 50.0, 65000.00, TRUE),
('f1111111-0000-0000-0000-000000000005', 'INV-CU-101', 'Paper Cup 12oz & Lid', 'packaging', 'pcs', 100.0, 1000.0, 450.00, TRUE),
('f1111111-0000-0000-0000-000000000006', 'INV-CR-005', 'Butter Croissant Dough', 'dry_goods', 'pcs', 20.0, 150.0, 8500.00, TRUE),
('f1111111-0000-0000-0000-000000000007', 'INV-MT-020', 'Matcha Powder 500g', 'dry_goods', 'pack', 5.0, 30.0, 95000.00, TRUE)
ON CONFLICT (sku) DO NOTHING;

-- 8. MENU & VARIAN
INSERT INTO products (id, category_id, branch_id, name, sku, description, base_price, target_station, image_url, is_active) VALUES
('fa111111-0000-0000-0000-000000000001', 'd1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'Cafe Latte', 'PRD-LATTE', 'Espresso klasik dipadukan dengan steamed fresh milk gurih', 35000.00, 'barista', 'https://images.unsplash.com/photo-1570968915860-54d5c301fa9f?w=400', TRUE),
('fa111111-0000-0000-0000-000000000002', 'd1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'Cappuccino', 'PRD-CAPPU', 'Espresso seimbang dengan foam susu tebal dan taburan cocoa', 38000.00, 'barista', 'https://images.unsplash.com/photo-1534778101976-62847782c213?w=400', TRUE),
('fa111111-0000-0000-0000-000000000003', 'd1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'Caramel Macchiato', 'PRD-MACCH', 'Espresso dilapisi susu vanilla manis dan saus karamel leleh', 42000.00, 'barista', 'https://images.unsplash.com/photo-1485808191679-5f86510681a2?w=400', TRUE),
('fa111111-0000-0000-0000-000000000004', 'd1111111-0000-0000-0000-000000000002', 'b1111111-0000-0000-0000-000000000001', 'Matcha Green Tea Latte', 'PRD-MATCHA', 'Bubuk matcha Uji premium dengan susu segar murni', 38000.00, 'barista', 'https://images.unsplash.com/photo-1536256263959-770b48d82b0a?w=400', TRUE),
('fa111111-0000-0000-0000-000000000005', 'd1111111-0000-0000-0000-000000000003', 'b1111111-0000-0000-0000-000000000001', 'Butter Croissant', 'PRD-CROISS', 'Pastry khas Perancis renyah berlapis dengan aroma butter wangi', 35000.00, 'kitchen', 'https://images.unsplash.com/photo-1555507036-ab1f4038808a?w=400', TRUE),
('fa111111-0000-0000-0000-000000000006', 'd1111111-0000-0000-0000-000000000004', 'b1111111-0000-0000-0000-000000000001', 'Basque Burnt Cheesecake', 'PRD-CHEESE', 'Kue keju lembut meleleh dengan permukaan karamelisasi gosong', 45000.00, 'kitchen', 'https://images.unsplash.com/photo-1533134242443-d4fd215305ad?w=400', TRUE)
ON CONFLICT (sku) DO NOTHING;

INSERT INTO product_variants (id, product_id, name, sku, additional_price) VALUES
('ea111111-0000-0000-0000-000000000001', 'fa111111-0000-0000-0000-000000000001', 'Hot', 'LATTE-HOT', 0.00),
('ea111111-0000-0000-0000-000000000002', 'fa111111-0000-0000-0000-000000000001', 'Iced', 'LATTE-ICE', 3000.00),
('ea111111-0000-0000-0000-000000000003', 'fa111111-0000-0000-0000-000000000002', 'Regular', 'CAPPU-REG', 0.00),
('ea111111-0000-0000-0000-000000000004', 'fa111111-0000-0000-0000-000000000002', 'Large', 'CAPPU-LRG', 5000.00)
ON CONFLICT DO NOTHING;

-- 9. RESEP (BOM) UNTUK OTOMATISASI PENGURANGAN STOK
INSERT INTO product_recipes (product_id, variant_id, inventory_item_id, quantity_required, uom, instructions) VALUES
('fa111111-0000-0000-0000-000000000001', 'ea111111-0000-0000-0000-000000000002', 'f1111111-0000-0000-0000-000000000001', 0.0180, 'kg', '18 gram espresso grind'),
('fa111111-0000-0000-0000-000000000001', 'ea111111-0000-0000-0000-000000000002', 'f1111111-0000-0000-0000-000000000002', 0.2000, 'liter', '200 ml cold fresh milk'),
('fa111111-0000-0000-0000-000000000001', 'ea111111-0000-0000-0000-000000000002', 'f1111111-0000-0000-0000-000000000005', 1.0000, 'pcs', '1 pcs cup & lid'),
('fa111111-0000-0000-0000-000000000005', NULL, 'f1111111-0000-0000-0000-000000000006', 1.0000, 'pcs', 'Bake croissant for 12 mins at 180C')
ON CONFLICT DO NOTHING;

-- 10. GUDANG & SALDO STOK AWAL
INSERT INTO warehouses (id, branch_id, name, type, address) VALUES
('ba111111-0000-0000-0000-000000000011', 'b1111111-0000-0000-0000-000000000001', 'Main Store Senopati', 'main', 'Lantai Basement Gudang Utama'),
('ba111111-0000-0000-0000-000000000012', 'b1111111-0000-0000-0000-000000000001', 'Bar Station Floor 1', 'bar', 'Bar Counter Lantai 1'),
('ba111111-0000-0000-0000-000000000013', 'b1111111-0000-0000-0000-000000000001', 'Kitchen Prep Storage', 'kitchen', 'Dapur Utama Belakang')
ON CONFLICT DO NOTHING;

INSERT INTO inventory_stocks (warehouse_id, inventory_item_id, quantity) VALUES
('ba111111-0000-0000-0000-000000000011', 'f1111111-0000-0000-0000-000000000001', 85.0),
('ba111111-0000-0000-0000-000000000011', 'f1111111-0000-0000-0000-000000000002', 80.0),
('ba111111-0000-0000-0000-000000000011', 'f1111111-0000-0000-0000-000000000003', 20.0),
('ba111111-0000-0000-0000-000000000011', 'f1111111-0000-0000-0000-000000000004', 14.0),
('ba111111-0000-0000-0000-000000000011', 'f1111111-0000-0000-0000-000000000005', 400.0),
('ba111111-0000-0000-0000-000000000011', 'f1111111-0000-0000-0000-000000000006', 98.0),
('ba111111-0000-0000-0000-000000000011', 'f1111111-0000-0000-0000-000000000007', 15.0)
ON CONFLICT (warehouse_id, inventory_item_id) DO UPDATE SET quantity = EXCLUDED.quantity;

-- 11. SUPPLIER & VENDOR
INSERT INTO suppliers (id, code, name, contact_person, phone, email, address, is_active) VALUES
('5a111111-0000-0000-0000-000000000001', 'SUP-KOPI-01', 'PT Kopi Nusantara Jaya', 'Hendrik Wijaya', '08119876543', 'hendrik@kopinusantara.id', 'Kawasan Industri Pulogadung, Jakarta Timur', TRUE),
('5a111111-0000-0000-0000-000000000002', 'SUP-SUSU-02', 'PT Susu Segar Sejahtera', 'Dewi Anggraini', '08128765432', 'order@sususegar.com', 'Jl. Raya Bogor KM 28, Cimanggis, Depok', TRUE),
('5a111111-0000-0000-0000-000000000003', 'SUP-KEMAS-03', 'PT Sentosa Packaging Solusindo', 'Budi Pratama', '08137654321', 'sales@sentosapack.co.id', 'Cengkareng, Jakarta Barat', TRUE)
ON CONFLICT (code) DO NOTHING;

-- 12. HRIS - DEPARTEMEN & JABATAN
INSERT INTO departments (id, name, description) VALUES
('de111111-0000-0000-0000-000000000001', 'Management', 'Direksi dan Operasional Manajemen'),
('de111111-0000-0000-0000-000000000002', 'Front of House', 'Barista, Kasir dan Pelayan'),
('de111111-0000-0000-0000-000000000003', 'Kitchen', 'Chef dan Staf Dapur'),
('de111111-0000-0000-0000-000000000004', 'Warehouse & Supply', 'Staf Gudang dan Pengadaan')
ON CONFLICT DO NOTHING;

INSERT INTO positions (id, department_id, title, level, base_salary) VALUES
('ec111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000001', 'Cafe Manager', 3, 8500000.00),
('ec111111-0000-0000-0000-000000000002', 'de111111-0000-0000-0000-000000000002', 'Head Barista', 2, 6000000.00),
('ec111111-0000-0000-0000-000000000003', 'de111111-0000-0000-0000-000000000002', 'Barista', 1, 4800000.00),
('ec111111-0000-0000-0000-000000000004', 'de111111-0000-0000-0000-000000000002', 'Cashier', 1, 4800000.00),
('ec111111-0000-0000-0000-000000000005', 'de111111-0000-0000-0000-000000000003', 'Chef / Kitchen Cook', 2, 6500000.00)
ON CONFLICT DO NOTHING;

INSERT INTO employees (id, branch_id, department_id, position_id, nik, first_name, last_name, email, phone, job_title, basic_salary, status, join_date) VALUES
('eb111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000001', 'ec111111-0000-0000-0000-000000000001', 'EMP-001', 'Sarah', 'Johnson', 'sarah@cafe-erp.com', '081234567891', 'Cafe Manager', 8500000.00, 'active', '2023-01-10'),
('eb111111-0000-0000-0000-000000000002', 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000003', 'ec111111-0000-0000-0000-000000000005', 'EMP-002', 'John', 'Doe', 'kitchen@cafe-erp.com', '081234567893', 'Head Chef', 6500000.00, 'active', '2023-02-15'),
('eb111111-0000-0000-0000-000000000003', 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000002', 'ec111111-0000-0000-0000-000000000003', 'EMP-003', 'Sarah', 'Andini', 'sarah.andini@email.com', '0812-3456-7890', 'Barista', 5500000.00, 'active', '2023-03-01'),
('eb111111-0000-0000-0000-000000000004', 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000002', 'ec111111-0000-0000-0000-000000000004', 'EMP-004', 'Jane', 'Doe', 'cashier@cafe-erp.com', '081234567892', 'Cashier', 4800000.00, 'active', '2023-04-12')
ON CONFLICT (nik) DO NOTHING;

-- 13. SHIFT KERJA
INSERT INTO work_shifts (id, branch_id, name, start_time, end_time, color) VALUES
('ac111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'Pagi', '07:00:00', '15:00:00', '#2563eb'),
('ac111111-0000-0000-0000-000000000002', 'b1111111-0000-0000-0000-000000000001', 'Siang', '15:00:00', '23:00:00', '#16a34a'),
('ac111111-0000-0000-0000-000000000003', 'b1111111-0000-0000-0000-000000000001', 'Malam', '23:00:00', '07:00:00', '#9333ea')
ON CONFLICT DO NOTHING;

-- 14. BAGAN AKUN (CHART OF ACCOUNTS)
INSERT INTO chart_of_accounts (id, code, name, account_type, description) VALUES
('ab111111-0000-0000-0000-000000000001', '1101', 'Kas Kasir (Cash Drawer)', 'asset', 'Kas tunai di laci kasir outlet'),
('ab111111-0000-0000-0000-000000000002', '1102', 'Bank BCA Operasional', 'asset', 'Rekening bank operasional utama'),
('ab111111-0000-0000-0000-000000000003', '1201', 'Persediaan Bahan Baku', 'asset', 'Total nilai fisik persediaan bahan'),
('ab111111-0000-0000-0000-000000000004', '2101', 'Hutang Dagang Supplier', 'liability', 'Hutang atas pembelian bahan baku'),
('ab111111-0000-0000-0000-000000000005', '3101', 'Modal Disetor', 'equity', 'Modal pendiri kafe'),
('ab111111-0000-0000-0000-000000000006', '4101', 'Pendapatan Penjualan Minuman', 'revenue', 'Omset dari kategori kopi & non-kopi'),
('ab111111-0000-0000-0000-000000000007', '4102', 'Pendapatan Penjualan Makanan', 'revenue', 'Omset dari kategori makanan & pastry'),
('ab111111-0000-0000-0000-000000000008', '5101', 'HPP Bahan Baku (COGS)', 'expense', 'Harga pokok penjualan dari pemakaian resep'),
('ab111111-0000-0000-0000-000000000009', '6100', 'Beban Sewa Outlet', 'expense', 'Sewa gedung dan tempat usaha'),
('ab111111-0000-0000-0000-000000000010', '6200', 'Beban Listrik, Air & Internet', 'expense', 'Biaya utilitas bulanan')
ON CONFLICT (code) DO NOTHING;

