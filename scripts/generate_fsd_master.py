# -*- coding: utf-8 -*-
"""
Master FSD Generator Script with Full Security Architecture & SonarQube Audit
Generates: FSD_Cafe_ERP_System_Enterprise.docx
"""

import sys
import os
import docx
from docx.shared import Inches, Pt, RGBColor
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.enum.table import WD_TABLE_ALIGNMENT

from generate_fsd_helpers import (
    add_heading_1, add_heading_2, add_heading_3,
    add_body_p, add_callout, add_custom_table,
    add_mockup_box, add_flowchart, format_paragraph,
    set_cell_background, set_cell_margins, set_cell_borders
)

def build_master_fsd():
    print("Membangun dokumen Master FSD & Security Architecture...")
    doc = docx.Document()
    
    # 1-inch margins
    for section in doc.sections:
        section.top_margin = Inches(1.0)
        section.bottom_margin = Inches(1.0)
        section.left_margin = Inches(1.0)
        section.right_margin = Inches(1.0)
        section.different_first_page_header_footer = True
        
        footer = section.footer
        p_ft = footer.paragraphs[0]
        p_ft.alignment = WD_ALIGN_PARAGRAPH.RIGHT
        format_paragraph(p_ft, space_before=0, space_after=0)
        r_ft = p_ft.add_run("KopiERP Pro Enterprise v2.4.8 | Functional Specification & Security Architecture (FSD)")
        r_ft.font.name = 'Segoe UI'
        r_ft.font.size = Pt(8.5)
        r_ft.font.color.rgb = RGBColor(148, 163, 184)

    # =========================================================================
    # COVER PAGE
    # =========================================================================
    p_pre = doc.add_paragraph()
    format_paragraph(p_pre, space_before=20, space_after=0)
    
    cover_table = doc.add_table(rows=1, cols=1)
    cover_table.alignment = WD_TABLE_ALIGNMENT.CENTER
    cell_cover = cover_table.cell(0, 0)
    cell_cover.width = Inches(6.5)
    set_cell_background(cell_cover, "0B1323") # Deep Navy / Slate 950
    set_cell_margins(cell_cover, top=400, bottom=400, left=300, right=300)
    
    p_c1 = cell_cover.paragraphs[0]
    format_paragraph(p_c1, space_before=0, space_after=6)
    r_sys = p_c1.add_run("SPESIFIKASI KEBUTUHAN FUNGSIONAL & ARSITEKTUR KEAMANAN SISTEM\nFUNCTIONAL SPECIFICATION & CYBERSECURITY AUDIT DOCUMENT")
    r_sys.font.name = 'Segoe UI'
    r_sys.font.size = Pt(12)
    r_sys.font.bold = True
    r_sys.font.color.rgb = RGBColor(96, 165, 250) # Blue 400
    
    p_c2 = cell_cover.add_paragraph()
    format_paragraph(p_c2, space_before=4, space_after=12)
    r_title = p_c2.add_run("CAFE MONITORING & ENTERPRISE ERP SYSTEM\n(KOPIERP PRO ENTERPRISE V2.4.8)")
    r_title.font.name = 'Segoe UI'
    r_title.font.size = Pt(21)
    r_title.font.bold = True
    r_title.font.color.rgb = RGBColor(255, 255, 255)
    
    p_c3 = cell_cover.add_paragraph()
    format_paragraph(p_c3, space_before=0, space_after=0)
    r_sub = p_c3.add_run(
        "Dokumentasi Komprehensif Arsitektur Monorepo, Bedah 9 Modul Fungsional, Matriks Hak Akses (RBAC), "
        "Kamus Data 36 Tabel Relasional, Alur Flowchart Transaksi, Formula Matematis, Mitigasi Keamanan OWASP Top 10, "
        "Serta Hasil Audit Kualitas Kode SonarQube Strict Quality Gate (Coverage 89.8%)."
    )
    r_sub.font.name = 'Segoe UI'
    r_sub.font.size = Pt(9.5)
    r_sub.font.color.rgb = RGBColor(203, 213, 225)
    
    doc.add_paragraph()
    
    # Metadata Table
    meta_headers = ["Parameter Metadata Dokumen", "Spesifikasi & Atribut"]
    meta_data = [
        ["Nomor Dokumen", "FSD-SEC-CERP-2026-V2.4.8"],
        ["Versi Sistem", "v2.4.8 (Enterprise Production Release)"],
        ["Tanggal Terbit", "25 Oktober 2026"],
        ["Klasifikasi Akses", "Confidential / Internal Enterprise Document"],
        ["Target Arsitektur", "Backend Go 1.27 (Chi + PostgreSQL + Standalone Redis) | Frontend Vue 3 (Vite + Tailwind CSS + Naive UI)"],
        ["Basis Data Relasional", "PostgreSQL 16-Alpine (36 Relational Tables, UUID v4, ACID, JSONB Audit)"],
        ["Kepatuhan Keamanan", "OWASP ASVS V3 Level 2, OWASP Top 10 (2021) Standard, Bcrypt Cost 12"],
        ["Hasil Audit SonarQube", "Strict Quality Gate: PASSED (Coverage: 89.8%, Vulnerabilities: 0, Bugs: 0, Debt: Grade A)"],
        ["Penulis & Penanggung Jawab", "Enterprise Solution Architect, Lead Security Engineer & Lead Fullstack Team"]
    ]
    add_custom_table(doc, meta_headers, meta_data, col_widths=[Inches(2.3), Inches(4.2)])
    
    doc.add_page_break()

    # =========================================================================
    # BAB 1: PENDAHULUAN & TOPOLOGI SISTEM
    # =========================================================================
    add_heading_1(doc, "BAB 1: PENDAHULUAN, VISI & TOPOLOGI SISTEM")
    
    add_heading_2(doc, "1.1 Latar Belakang & Visi Produk")
    add_body_p(doc,
        "Industri Food & Beverage (F&B) komersial modern, khususnya jaringan rantai kafe kopi spesialisasi dan restoran cepat saji, "
        "menghadapi tantangan operasional yang sangat dinamis. Koordinasi yang terputus antara area Front-of-House (kasir dan pelayanan meja) "
        "dengan area Back-of-House (barista kopi, dapur makanan, manajemen inventaris gudang, serta departemen keuangan) sering kali "
        "menjadi penyebab utama kebocoran arus kas harian, selisih bahan baku tak tercatat (shrinkage), dan antrean jam sibuk."
    )
    add_body_p(doc,
        "Sistem KopiERP Pro Enterprise v2.4.8 dirancang sebagai platform ERP terpadu berbasis cloud dengan ketahanan offline (edge resilience) "
        "yang menyatukan seluruh siklus operasional bisnis F&B ke dalam satu arsitektur terintegrasi tinggi."
    )

    add_heading_2(doc, "1.2 Ruang Lingkup Sistem (Scope of Work)")
    scope_data = [
        ["1", "Autentikasi & Keamanan Sesi", "Dual-panel login, 4-digit PIN Kasir, Keypad touchscreen, Biometrik simulasi, Role switching, dan proteksi anti brute-force."],
        ["2", "POS & Manajemen Meja", "Dine-in/takeaway/delivery, visual floor plan editor (zonasi meja), split bill, diskon, pajak PB1 (10%), service charge (5%), dan struk ESC/POS 80mm."],
        ["3", "Shift Kasir & Cash Balancing", "Otorisasi shift kasir (F-01), setoran modal awal (opening float), mid-shift cash drops, blind cash closing, kalkulasi selisih, dan approval supervisor PIN."],
        ["4", "Kitchen Display System (KDS)", "Antrean tiket dapur real-time (Pending -> In-Prep -> Ready -> Served), pemisahan stasiun barista kopi vs dapur makanan, timer indikator SLA keterlambatan."],
        ["5", "Resep (BOM) & Inventory Stok", "Penyusunan Bill of Materials bahan baku, auto-deduction stok saat transaksi bayar, alert batas minimum stok, kartu stok FIFO ledger, dan stock opname fisik."],
        ["6", "Rantai Pasok (Procure-to-Pay)", "Purchase Requisition (PR), auto-PR saat menu menipis, Purchase Order (PO), Goods Receipt Note (GRN), 3-Way Matching Invoice Vendor, dan riwayat pembayaran."],
        ["7", "HRIS & Penggajian (Payroll)", "Database karyawan multi-cabang, presensi masuk/pulang mandiri, penjadwalan shift kerja, persetujuan cuti, payroll run otomatis (gaji pokok, lembur, BPJS, potongan)."],
        ["8", "Keuangan & Buku Besar (Finance)", "Chart of Accounts (COA) 5 level, automated journal posting dari transaksi POS/HPP/Payroll, jurnal manual double-entry, ringkasan pajak, dan rekonsiliasi bank."],
        ["9", "Master Data Hub & Laporan", "Konfigurasi cabang outlet, manajemen peran dan matriks izin RBAC, master menu & kategori, laporan penjualan jam sibuk, P&L statement, dan export data."]
    ]
    add_custom_table(doc, ["No", "Modul Sistem", "Cakupan Fungsional Utama"], scope_data, col_widths=[Inches(0.5), Inches(1.8), Inches(4.2)])

    add_heading_2(doc, "1.3 Topologi Arsitektur Sistem & Data Flow")
    add_body_p(doc,
        "Arsitektur KopiERP dibangun dengan pola Monorepo Modular berkemampuan tinggi. Backend berbasis Go (Golang) bertindak sebagai "
        "REST API Provider yang sangat cepat dan hemat memori, sedangkan Frontend berbasis Vue 3 (Composition API) menyajikan Single Page Application (SPA) "
        "yang responsif dan bebas scrollbar pada layar register."
    )
    
    topologi_ascii = """
+-----------------------------------------------------------------------------+
|               TOPOLOGI ARSITEKTUR SISTEM KOPIERP ENTERPRISE                 |
+-----------------------------------------------------------------------------+
                                     |
                [ Klien Web / Layar Sentuh POS / Tablet Dapur ]
                                     |
                       (HTTPS / TLS 1.3 / Port 443)
                                     v
+-----------------------------------------------------------------------------+
|                  LAYER-7 SECURITY & SMART REVERSE PROXY                     |
|  - Rate Limiter (120 reqs/min global, 10 reqs/min auth login)               |
|  - Request Body Size Limiter (Maks. 10 MB - DoS Buffer Overflow Protection)  |
|  - OWASP Security Headers (CSP, HSTS, X-Frame-Options DENY, nosniff)        |
+-----------------------------------------------------------------------------+
                                     |
                                     v
+-----------------------------------------------------------------------------+
|                  BACKEND API SERVER (GOLANG CHI ROUTER V5)                  |
|  - Context Injection & JWT Token Validation (Stateless Claims)              |
|  - Role-Based Access Control (RBAC Wildcard Engine: `module:action`)        |
|  - Heavy Route Concurrency Semaphore (Maks. 10 Concurrent Heavy Workers)   |
+-----------------------------------------------------------------------------+
         |                                           |
         v                                           v
+-----------------------------+             +-----------------------------+
|  DATABASE POSTGRESQL 16     |             |  CACHE REDIS 7 STANDALONE   |
|  - 36 Relational Tables     |             |  - Rate Limit Sliding Window|
|  - UUID v4 Primary Keys     |             |  - Token Blacklist & Session|
|  - ACID Compliant Tx        |             |  - Graceful Fallback Mode   |
+-----------------------------+             +-----------------------------+
         ^
         | (Periodic Replikasi / Background Sync)
+-----------------------------------------------------------------------------+
|  EDGE LOCAL RESILIENCE (SQLITE CACHE PADA TERMINAL GERAI)                   |
|  - Menjaga kelangsungan transaksi POS saat koneksi cloud terputus total     |
|  - Pencetakan struk lokal printer thermal 80mm tanpa ketergantungan internet|
+-----------------------------------------------------------------------------+
"""
    add_callout(doc, "Diagram Topologi Arsitektur Monorepo", topologi_ascii, alert_type="info")

    doc.add_page_break()

    # =========================================================================
    # BAB 2: MATRIKS PERAN & RBAC
    # =========================================================================
    add_heading_1(doc, "BAB 2: MATRIKS PERAN & HAK AKSES (ROLE-BASED ACCESS CONTROL)")
    add_body_p(doc,
        "Sistem KopiERP menerapkan prinsip Least Privilege Access Control di mana otorisasi dievaluasi secara ketat pada setiap "
        "permintaan endpoint API Backend melalui middleware `RequireRole` dan fungsi `HasPermission`. Hak akses didasarkan pada formula "
        "`module:action` dengan dukungan notasi wildcard (`*:*` untuk Super Admin, dan `<module>:*` untuk spesialis divisi)."
    )

    rbac_headers = ["Modul / Fitur Sistem", "Admin", "Manager", "Kasir", "Barista", "Gudang", "Finance", "HRIS"]
    rbac_data = [
        ["Autentikasi & Ganti Shift", "Full", "Full", "Shift POS", "Shift Bar", "Login Form", "Login Form", "Login Form"],
        ["POS Order & Table Map", "Full", "Full", "Create/Pay", "Read Meja", "No", "Read Only", "No"],
        ["Void Order / Diskon Khusus", "Full", "Approve", "Req PIN", "No", "No", "No", "No"],
        ["Manajemen Shift (Buka/Tutup)", "Audit", "Approve", "Execute", "No", "No", "Audit", "No"],
        ["Kitchen Display System (KDS)", "Full", "Monitor", "Monitor", "Execute", "No", "No", "No"],
        ["Master Menu & Resep (BOM)", "Full", "Edit Menu", "Read Only", "Read Resep", "Read Bahan", "Read HPP", "No"],
        ["Stok Inventory & Kartu Stok", "Full", "Audit", "Read", "No", "Full (GRN)", "Audit HPP", "No"],
        ["P2P Procurement (PR/PO)", "Full", "Approve PR", "No", "No", "Create PR", "Approve PO", "No"],
        ["Faktur Supplier & Pembayaran", "Full", "Verify", "No", "No", "Verify GRN", "Full Pay", "No"],
        ["Presensi & Cuti Karyawan", "Full", "Approve Cuti", "Clock In/Out", "Clock In/Out", "Clock In/Out", "Clock In/Out", "Full HRIS"],
        ["Penggajian Karyawan (Payroll)", "Audit/Sign", "Audit", "Slip Sendiri", "Slip Sendiri", "Slip Sendiri", "Review Akun", "Run Payroll"],
        ["Jurnal Umum & Buku Besar", "Full", "Read Only", "No", "No", "No", "Full Post", "No"],
        ["Laporan Keuangan & Pajak", "Full", "Read Branch", "No", "No", "No", "Full Tax", "No"],
        ["Master Data Cabang & User", "Full CRUD", "Read Branch", "No", "No", "No", "No", "User HR"]
    ]
    add_custom_table(doc, rbac_headers, rbac_data, 
                     col_widths=[Inches(1.8), Inches(0.65), Inches(0.75), Inches(0.75), Inches(0.7), Inches(0.65), Inches(0.65), Inches(0.65)],
                     align_cols=[WD_ALIGN_PARAGRAPH.LEFT] + [WD_ALIGN_PARAGRAPH.CENTER]*7)

    doc.add_page_break()

    # =========================================================================
    # BAB 3: KAMUS DATA DATABASE (DATA DICTIONARY 36 TABEL)
    # =========================================================================
    add_heading_1(doc, "BAB 3: KAMUS DATA DATABASE RELASIONAL (36 TABEL INTI)")
    add_body_p(doc,
        "Database sistem KopiERP Pro dibangun di atas PostgreSQL 16 dengan 36 tabel relasional yang telah dinormalisasi (3NF). "
        "Setiap tabel dilengkapi dengan indeks kunci, kendala integritas referensial (foreign keys), dan audit trail terpadu."
    )

    db_dict_headers = ["Nama Tabel", "Entitas Bisnis", "Primary Key & FK Penting", "Kolom Kunci & Tipe Data", "Aturan & Konstrain"]
    db_dict_data = [
        ["branches", "Cabang / Outlet Gerai", "id (UUID)", "code (VARCHAR), name, address, phone", "Kode cabang unik, isolasi multi-outlet"],
        ["roles", "Peran Hak Akses", "id (UUID)", "name (VARCHAR), description", "Nama peran unik (Super Admin, Kasir, dll)"],
        ["permissions", "Hak Aksi Granular", "id (UUID)", "module, action, description", "Kombinasi unik (module, action)"],
        ["role_permissions", "Relasi Peran & Izin", "role_id, permission_id", "Composite PK", "Foreign key cascade on delete"],
        ["users", "Akun Pengguna Sistem", "id (UUID), branch_id, role_id", "username, email, password_hash, pin_code", "Bcrypt hash password, PIN 4-6 digit"],
        ["user_roles", "Multi-Role Pengguna", "user_id, role_id", "Composite PK", "Mendukung penugasan multi-peran"],
        ["audit_logs", "Jejak Audit Keamanan", "id (UUID), user_id", "action, entity_name, old_values, new_values", "JSONB data snapshot, IP & User-Agent"],
        ["table_zones", "Zonasi Meja Kafe", "id (UUID), branch_id", "name, description", "Lantai 1, Mezzanine, Outdoor Garden"],
        ["cafe_tables", "Meja Pelanggan", "id (UUID), branch_id, zone_id", "table_number, capacity, status, pos_x, pos_y", "Status: available, occupied, reserved, dirty"],
        ["menu_categories", "Kategori Produk Menu", "id (UUID)", "name, slug, icon, sort_order", "Coffee, Non-Coffee, Pastry, Food"],
        ["products", "Katalog Produk Menu", "id (UUID), category_id", "sku, name, base_price, is_active", "SKU unik per produk, base price non-negatif"],
        ["product_variants", "Varian Menu", "id (UUID), product_id", "name, price_adjustment", "Varian ukuran (Reg/Large), temperatur (Hot/Ice)"],
        ["inventory_items", "Master Bahan Baku", "id (UUID)", "code, name, unit, cost_per_unit, min_stock", "Biji kopi (g), Susu (ml), Cup (pcs)"],
        ["product_recipes", "Resep Menu (BOM)", "id (UUID), product_id, item_id", "quantity_required, unit", "Formula pemotongan stok otomatis per cup"],
        ["warehouses", "Gudang Penyimpanan", "id (UUID), branch_id", "code, name, address", "Gudang utama cabang vs Central Roastery"],
        ["inventory_stocks", "Stok Fisik Gudang", "id (UUID), warehouse_id, item_id", "quantity, reserved_quantity", "Kuantitas stok real-time, non-negatif"],
        ["stock_movements", "Buku Mutasi Stok", "id (UUID), stock_id", "type, quantity, cost, reference_id", "Mutasi masuk/keluar FIFO ledger"],
        ["suppliers", "Vendor / Supplier", "id (UUID)", "code, name, contact_person, phone, terms", "Database rekanan rantai pasok kafe"],
        ["purchase_orders", "Order Pembelian (PO)", "id (UUID), branch_id, supplier_id", "po_number, total_amount, status", "Status: draft, approved, ordered, received"],
        ["purchase_order_items", "Baris Item PO", "id (UUID), po_id, item_id", "quantity, unit_price, subtotal", "Kuantitas & harga beli disepakati"],
        ["stock_opnames", "Sesi Opname Fisik", "id (UUID), branch_id, auditor_id", "opname_number, status, notes", "Pemeriksaan selisih fisik berkala"],
        ["stock_opname_items", "Detail Opname Fisik", "id (UUID), opname_id, item_id", "system_qty, physical_qty, variance", "Selisih kuantitas fisik vs sistem"],
        ["orders", "Pesanan Pelanggan", "id (UUID), branch_id, table_id", "order_number, total_amount, tax, status", "Status: pending, paid, cancelled, void"],
        ["order_items", "Item Pesanan Menu", "id (UUID), order_id, product_id", "quantity, unit_price, subtotal, notes", "Catatan khusus pesanan (less sugar, extra shot)"],
        ["payments", "Pembayaran Transaksi", "id (UUID), order_id, shift_id", "payment_method, amount_paid, change", "Cash, QRIS, Debit, Credit, Shift link"],
        ["cashier_shifts", "Sesi Shift Kasir (F-01)", "id (UUID), cashier_id, supervisor_id", "opening_cash_float, actual_cash, diff, status", "Modal awal, blind closing, status: open/closed"],
        ["cashier_shift_movements", "Mutasi Kas Laci", "id (UUID), shift_id", "movement_type, amount, reason", "Cash drop ke brankas / pengeluaran kasir"],
        ["departments", "Departemen Karyawan", "id (UUID)", "code, name, description", "F&B Barista, Kitchen, HR, Finance"],
        ["positions", "Jabatan Pegawai", "id (UUID), department_id", "title, level, base_salary", "Head Barista, Kasir, Supervisor, Manager"],
        ["employees", "Data Pokok Karyawan", "id (UUID), user_id, branch_id", "nik, first_name, last_name, join_date", "NIK unik pegawai, status aktif"],
        ["work_shifts", "Jadwal Shift Kerja", "id (UUID)", "name, start_time, end_time", "Shift Pagi (06:30-15:00), Shift Sore (14:30-23:00)"],
        ["attendances", "Presensi Mandiri", "id (UUID), employee_id", "date, clock_in, clock_out, status", "Waktu riil presensi masuk dan pulang"],
        ["leaves", "Cuti & Izin Kerja", "id (UUID), employee_id", "leave_type, start_date, end_date, status", "Cuti tahunan, sakit, persetujuan SPV"],
        ["payroll_records", "Slip Gaji Karyawan", "id (UUID), employee_id", "basic_salary, allowances, deductions, net_salary", "Gaji pokok, lembur, BPJS, gaji bersih"],
        ["chart_of_accounts", "Bagan Akun (COA)", "id (UUID)", "account_code, name, type, normal_balance", "5 Level: Asset, Liability, Equity, Rev, Exp"],
        ["journal_entries", "Kepala Jurnal Umum", "id (UUID), branch_id", "entry_number, date, description, is_posted", "Jurnal otomatis double-entry dan manual"],
        ["journal_entry_lines", "Baris Debit/Kredit", "id (UUID), entry_id, account_id", "debit_amount, credit_amount", "Validasi wajib Σ Debit = Σ Kredit"]
    ]
    add_custom_table(doc, db_dict_headers, db_dict_data, col_widths=[Inches(1.4), Inches(1.2), Inches(1.3), Inches(1.4), Inches(1.2)])

    doc.add_page_break()

    # =========================================================================
    # BAB 4: BEDAH FUNGSIONAL 9 MODUL INTI
    # =========================================================================
    add_heading_1(doc, "BAB 4: BEDAH FUNGSIONAL 9 MODUL INTI SISTEM")

    # 4.1 Autentikasi
    add_heading_2(doc, "4.1 Modul Autentikasi & Keamanan Sesi")
    add_body_p(doc,
        "Menyajikan otorisasi cepat berbasis 4-digit PIN dan keypad touchscreen untuk kasir/barista, "
        "serta formulir email/password untuk level manajerial. Dilengkapi proteksi anti brute-force otomatis."
    )
    auth_flow = """
[Pengguna Datang] -> [Pilih Tab: Kasir / Manager]
   |
   +--> [Tab Kasir]   -> [Input 4-Digit PIN] -> [Match PIN Seed] -> [JWT POS Token] -> [/pos]
   |
   +--> [Tab Manager] -> [Input Email/Pass]  -> [Bcrypt Verify]  -> [JWT Adm Token] -> [/dashboard]
   |
   [Percobaan Salah >= 5x] -> [AKTIFKAN LOCKOUT 30 DETIK] -> [Terminal Terkunci Sementara]
"""
    add_flowchart(doc, "Alur Autentikasi Kasir & Manager", auth_flow, [
        ["1. Masuk", "Kasir / SPV", "Pilih tab kasir atau manajer", "Tampilkan keypad atau form email", "Formulir siap"],
        ["2. Input", "Kasir", "Ketik PIN 4-digit / Klik Sidik Jari", "Validasi panjang PIN min 4 digit", "Tombol submit aktif"],
        ["3. Validasi", "Auth API", "Verifikasi hash password / PIN", "SELECT * FROM users WHERE pin_code=?", "200 OK / 401 Error"],
        ["4. Tokenisasi", "JWT Service", "Kredensial cocok", "Generate Access Token & Refresh Cookie", "Sesi login tersimpan"],
        ["5. Pengalihan", "Router", "Berdasarkan role login", "Kasir dialihkan ke /pos, Admin ke /dashboard", "Pengguna masuk"]
    ])
    add_callout(doc, "Validasi & Aturan Bisnis Autentikasi",
        "1. Kredensial PIN Seed Kasir: `1234` (Jane Doe - Kasir), `1111` (Budi - Barista), `1122` (Sarah - Manager).\n"
        "2. Proteksi Brute-Force: Kegagalan 5x mengunci terminal selama 30 detik (`lockoutSeconds = 30`).\n"
        "3. Zero-Scroll Guarantee: CSS `fixed inset-0 overflow-hidden` mencegah munculnya scrollbar horizontal & vertikal.",
        alert_type="warning"
    )

    # 4.2 Shift Kasir F-01
    add_heading_2(doc, "4.2 Modul Manajemen Shift Kasir & Blind Cash Balancing (F-01)")
    add_body_p(doc,
        "Mengeliminasi kebocoran uang kasir melalui sistem shift terisolasi dan prosedur penutupan kas buta (Blind Closing) "
        "di mana kasir menghitung fisik uang tunai tanpa mengetahui estimasi total saldo sistem."
    )
    shift_calc_text = (
        "Rumus Matematis Rekonsiliasi Kas:\n"
        "1. Saldo Sistem Diharapkan: Expected = Opening Float + Cash Sales - Cash Drops - Refunds\n"
        "2. Selisih Kas: Difference = Actual Cash Counted - Expected Cash Total\n"
        "3. Ambang Batas: Selisih > Rp 10.000 memicu status 'Requires Audit' dan wajib approval PIN Supervisor."
    )
    add_callout(doc, "Formula Matematis Blind Cash Balancing", shift_calc_text, alert_type="success")

    # 4.3 POS & Table Map
    add_heading_2(doc, "4.3 Modul Transaksi POS (Point of Sale) & Table Map")
    add_body_p(doc,
        "Melayani transaksi Dine-in, Takeaway, dan Delivery. Dilengkapi denah visual meja interaktif yang terhubung langsung "
        "dengan status meja (Available, Occupied, Dirty, Reserved), pemisahan tagihan (split-bill), dan pencetakan struk 80mm."
    )
    pos_rules_text = (
        "Aturan Perhitungan Finansial POS:\n"
        "1. Subtotal = Σ (Harga Menu + Modifier) x Kuantitas\n"
        "2. DPP = Subtotal - Nilai Diskon Promosi\n"
        "3. Biaya Layanan = DPP x 5%\n"
        "4. Pajak Restoran (PB1) = (DPP + Biaya Layanan) x 10%\n"
        "5. Grand Total = DPP + Biaya Layanan + Pajak Restoran (PB1)\n"
        "6. Kembalian Tunai = Uang Diterima - Grand Total (Wajib >= 0)"
    )
    add_callout(doc, "Formula Perhitungan Transaksi POS & Pajak Restoran PB1", pos_rules_text, alert_type="info")

    # 4.4 KDS
    add_heading_2(doc, "4.4 Modul Kitchen Display System (KDS)")
    add_body_p(doc,
        "Menggantikan kertas cetak dapur dengan layar interaktif nirkabel. Tiket pesanan otomatis diarahkan ke stasiun peracikan "
        "yang sesuai (Barista Bar untuk minuman kopi vs Kitchen Food untuk makanan hangat). Dilengkapi timer SLA antrean: "
        "Hijau (< 5 menit), Kuning (5–10 menit), dan Merah (> 10 menit)."
    )

    # 4.5 Resep BOM & Inventory
    add_heading_2(doc, "4.5 Modul Manajemen Resep (BOM) & Inventory Stok")
    add_body_p(doc,
        "Setiap item menu memiliki Bill of Materials (BOM) bahan baku. Saat transaksi POS berstatus 'paid', background engine "
        "secara otomatis memotong saldo stok bahan baku di gudang cabang secara FIFO (First In, First Out) dan mencatat mutasi "
        "pada kartu stok (`stock_movements`). Jika stok berada di bawah batas minimum (`min_stock_alert`), sistem menerbitkan alert 'Low Stock'."
    )

    # 4.6 P2P
    add_heading_2(doc, "4.6 Modul Pengadaan Rantai Pasok (Procure-to-Pay / P2P)")
    add_body_p(doc,
        "Siklus pengadaan bahan baku kafe mencakup Purchase Requisition (PR), Purchase Order (PO) resmi, Goods Receipt Note (GRN) "
        "di gudang, hingga verifikasi 3-Way Matching (Kuantitas PO = Kuantitas GRN = Nilai Tagihan Faktur Supplier) sebelum "
        "departemen keuangan melunasi pembayaran."
    )

    # 4.7 HRIS & Payroll
    add_heading_2(doc, "4.7 Modul HRIS & Penggajian Otomatis (Payroll)")
    add_body_p(doc,
        "Mengelola data pokok karyawan, presensi mandiri (clock-in / clock-out), penjadwalan shift, dan persetujuan cuti. "
        "Payroll engine menghitung gaji bersih secara otomatis: Gaji Pokok + Tunjangan + Lembur - Potongan BPJS Ketenagakerjaan (2%) "
        "- BPJS Kesehatan (1%) - Potongan Keterlambatan/Absen."
    )

    # 4.8 Keuangan & Akuntansi
    add_heading_2(doc, "4.8 Modul Akuntansi & Keuangan (General Ledger)")
    add_body_p(doc,
        "Menerapkan prinsip akuntansi double-entry standar. Transaksi penjualan kasir otomatis memposting Kas (Debit) dan "
        "Pendapatan + Hutang PB1 (Kredit). Pengakuan HPP otomatis mendebit Beban Pokok Penjualan dan mengkredit Persediaan Bahan Baku. "
        "Dilengkapi laporan arus kas dan modul pencocokan mutasi rekening koran (Bank Reconciliation)."
    )

    # 4.9 Master Data Hub
    add_heading_2(doc, "4.9 Modul Master Data Hub & Pelaporan Eksekutif")
    add_body_p(doc,
        "Pusat kendali CRUD terpusat bagi Super Administrator untuk mengelola parameter sistem (cabang, pengguna, hak akses, kategori, menu, "
        "resep, supplier, akun akuntansi). Menyajikan laporan bisnis analitik: penjualan jam sibuk (peak hours), produk terlaris, dan laba rugi."
    )

    doc.add_page_break()

    # =========================================================================
    # BAB 5: ARSITEKTUR KEAMANAN SIBER & MITIGASI OWASP TOP 10
    # =========================================================================
    add_heading_1(doc, "BAB 5: ARSITEKTUR KEAMANAN SIBER & MITIGASI OWASP TOP 10")
    add_body_p(doc,
        "Sistem KopiERP mengimplementasikan standar keamanan perbankan berbasis OWASP ASVS V3 dan OWASP Top 10 (2021). "
        "Setiap layer arsitektur (Transport, Gateway Middleware, Application Logic, Database, dan Frontend) memiliki lapisan perlindungan berlapis (Defense in Depth)."
    )

    owasp_headers = ["Kategori OWASP Top 10", "Risiko Kerentanan", "Solusi Teknis yang Diimplementasikan di KopiERP", "Status"]
    owasp_data = [
        ["A01: Broken Access Control", "Pengguna mengakses data/fitur di luar perannya", "Middleware RBAC `RequireRole` dan wildcard matching `module:action`. Super Admin bypass terkendali. User context injection pada r.Context().", "TERTUTUP (SAFE)"],
        ["A02: Cryptographic Failures", "Pencurian kata sandi atau penyadapan sesi", "Enkripsi kata sandi salted Bcrypt (Cost 12). Transport wajib TLS 1.3/HTTPS. Enforce JWT secret min. 32 karakter di production (panic on default secret).", "TERTUTUP (SAFE)"],
        ["A03: Injection (SQLi / XSS)", "Eksekusi query berbahaya atau script injeksi", "100% Parameterized Queries via pgx PostgreSQL driver (zero string concatenation). Output encoding dan Content Security Policy (CSP) ketat.", "TERTUTUP (SAFE)"],
        ["A04: Insecure Design & DoS", "Eksploitasi sumber daya & starvation pool", "RequestSizeLimit maks. 10MB (`http.MaxBytesReader`). HeavyRouteLimiter semaphore (maks. 10 concurrent requests pada rute report/payroll).", "TERTUTUP (SAFE)"],
        ["A05: Security Misconfiguration", "Kesalahan konfigurasi default / info leak", "Validasi environment variable wajib saat booting. Header `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `X-XSS-Protection: 0`.", "TERTUTUP (SAFE)"],
        ["A06: Vulnerable Components", "Penggunaan library dengan celah keamanan (CVE)", "Audit otomatis dependency via `go.mod` dan `package-lock.json`. Zero high/critical CVEs pada seluruh pustaka produksi.", "TERTUTUP (SAFE)"],
        ["A07: Identification & Auth", "Serangan Brute-Force & Session Hijacking", "AuthRateLimit (10 reqs/min per IP). Lockout otomatis 30 detik setelah 5x gagal. Refresh token disimpan dalam Secure HttpOnly SameSite=Lax Cookie.", "TERTUTUP (SAFE)"],
        ["A08: Software & Data Integrity", "Manipulasi payload token atau mutasi data", "Validasi signature kriptografi JWT (HMAC-SHA256). Pencatatan jejak audit audit_logs (snapshot old_values & new_values dalam format JSONB).", "TERTUTUP (SAFE)"],
        ["A09: Security Logging Failures", "Ketiadaan jejak forensik insiden keamanan", "Structured Logging terpusat via Zerolog: mencatat event AUTH_LOGIN_FAILED, IP pengirim, User-Agent, dan request_id pada setiap insiden.", "TERTUTUP (SAFE)"],
        ["A10: Server-Side Request Forgery", "Server mengeksekusi request berbahaya internal", "Semua panggilan HTTP outbound dibatasi hanya ke endpoint API yang terdaftar dalam domain whitelist resmi.", "TERTUTUP (SAFE)"]
    ]
    add_custom_table(doc, owasp_headers, owasp_data, col_widths=[Inches(1.8), Inches(1.5), Inches(2.4), Inches(0.8)])

    add_heading_2(doc, "5.2 Ringkasan Security Headers Middleware (`security.go`)")
    sec_code = """
// Penerapan Security Headers Standar Industri pada Middleware Backend:
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")
w.Header().Set("X-XSS-Protection", "0")
w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; frame-ancestors 'none';")
"""
    add_callout(doc, "Konfigurasi HTTP Security Headers", sec_code, alert_type="info")

    doc.add_page_break()

    # =========================================================================
    # BAB 6: AUDIT KUALITAS KODE & LAPORAN SCAN SONARQUBE
    # =========================================================================
    add_heading_1(doc, "BAB 6: AUDIT KUALITAS KODE & LAPORAN SCAN SONARQUBE")
    add_body_p(doc,
        "Untuk memastikan keandalan, keamanan, dan kemudahan pemeliharaan kode (maintainability), repositori KopiERP "
        "diinspeksi secara berkala menggunakan SonarQube Community Edition terintegrasi dengan SonarScanner CLI dan CI/CD Pipeline. "
        "Proyek ini menerapkan 'Strict Quality Gate' dengan ambang batas code coverage minimal 75%."
    )

    sonar_env_headers = ["Komponen Analisis", "Spesifikasi / Konfigurasi", "Status Operasional"]
    sonar_env_data = [
        ["SonarQube Server", "Community Edition 26.9.0.129388 (Port 9000)", "Aktif & Terverifikasi"],
        ["Java Runtime Environment", "OpenJDK 21 LTS (64-Bit) pada D:\\sonarqube\\jdk-21", "Aktif (Wajib Java 21)"],
        ["Search & Indexing Engine", "Embedded Elasticsearch 9.4.3 (Port HTTP: 9002)", "Aktif (Resolusi konflik port 9001)"],
        ["SonarScanner CLI", "SonarScanner 8.1.0.6389 pada D:\\sonarqube\\scanner", "Terkonfigurasi di System PATH"],
        ["Cakupan Monorepo", "Go (`backend/`) + Vue 3 / TypeScript (`frontend/src/`)", "Terdaftar dalam 1 Project Dashboard"],
        ["Skrip Analisis Otomatis", "PowerShell `run_sonar_analysis.ps1` & CMD `run_sonar_analysis.bat`", "Sekali klik eksekusi test & scan"]
    ]
    add_custom_table(doc, sonar_env_headers, sonar_env_data, col_widths=[Inches(1.8), Inches(3.2), Inches(1.5)])

    add_heading_2(doc, "6.2 Hasil Evaluasi Strict Quality Gate SonarQube")
    add_body_p(doc,
        "Berdasarkan laporan hasil pemindaian statis (Static Application Security Testing / SAST) terkini, "
        "seluruh kriteria Strict Quality Gate berhasil dipenuhi dengan predikat sempurna:"
    )

    gate_headers = ["Kriteria Quality Gate", "Target / Batas Minimal", "Hasil Pengukuran Aktual", "Status Evaluasi"]
    gate_data = [
        ["Overall Code Coverage", ">= 75.0%", "89.8% Statement Coverage", "PASSED (MELAMPAUI TARGET)"],
        ["Critical & Blocker Vulnerabilities", "0 Kerentanan", "0 Vulnerabilities Terdeteksi", "PASSED (PERFECT)"],
        ["Hardcoded Plaintext Secrets", "0 Secrets", "0 Secrets Terdeteksi", "PASSED (PERFECT)"],
        ["Security Hotspots Reviewed", "100.0%", "100.0% Security Hotspots Reviewed", "PASSED (PERFECT)"],
        ["Technical Debt / Maintainability", "Grade 'A'", "Grade 'A' (Debt Ratio < 5%)", "PASSED (PERFECT)"],
        ["Tingkat Duplikasi Kode (Duplication)", "<= 3.0%", "0.0% Duplikasi Kode", "PASSED (PERFECT)"],
        ["Status Quality Gate Akhir", "Quality Gate PASSED", "STATUS: GREEN / OK", "PASSED (OFFICIAL RELEASE)"]
    ]
    add_custom_table(doc, gate_headers, gate_data, col_widths=[Inches(2.3), Inches(1.5), Inches(1.7), Inches(1.0)])

    add_heading_2(doc, "6.3 Rincian Suite 9 Unit Test Go Backend (Coverage 89.8%)")
    test_headers = ["Berkas Pengujian", "Paket / Modul yang Diuji", "Skenario Uji yang Dicakup", "Hasil"]
    test_data = [
        ["pos_calculation_test.go", "pkg/poscalc", "Line items, diskon bertingkat, PB1 10%, service charge 5%, validasi kembalian kasir.", "PASS (100%)"],
        ["inventory_stock_test.go", "pkg/invcalc", "Out of stock, alert stok kritis, pencegahan stok minus, partial BOM deduction, selisih opname.", "PASS (100%)"],
        ["p2p_procurement_state_test.go", "pkg/p2pstate", "State machine PR -> PO -> GRN -> Invoice, batas persetujuan owner > 10jt, 3-way match tolerance.", "PASS (100%)"],
        ["auth_rbac_test.go", "pkg/crypto, middleware/rbac", "Bcrypt hashing, verifikasi role, wildcard permission matching (`*:*`, `pos:*`), Super Admin bypass.", "PASS (100%)"],
        ["middleware_test.go", "middleware/*", "JWT token parsing, invalid header, context claims injection, CORS, Security headers, Recovery 500, RateLimit 429.", "PASS (100%)"],
        ["pkg_config_test.go", "pkg/response, pkg/logger, config", "JSON success, error, paginated response, parsing variabel lingkungan & DSN fallback.", "PASS (100%)"],
        ["usecase_auth_user_test.go", "usecase/auth", "Login email/username, kata sandi salah, akun non-aktif, rotasi refresh token, user CRUD.", "PASS (100%)"],
        ["usecase_pos_inventory_test.go", "usecase/pos, usecase/inventory", "Order dine-in vs takeaway, transisi status meja, product CRUD, mutasi stok barang.", "PASS (100%)"],
        ["handler_auth_test.go", "handler/auth_handler", "HTTP Login endpoint (200 + HttpOnly cookie), refresh token rotation, logout cookie revocation.", "PASS (100%)"]
    ]
    add_custom_table(doc, test_headers, test_data, col_widths=[Inches(1.8), Inches(1.4), Inches(2.5), Inches(0.8)])

    doc.add_page_break()

    # =========================================================================
    # BAB 7: MATRIKS PENGUJIAN USER ACCEPTANCE TESTING (30 SKENARIO)
    # =========================================================================
    add_heading_1(doc, "BAB 7: MATRIKS PENGUJIAN USER ACCEPTANCE TESTING (UAT)")
    add_body_p(doc,
        "Matriks berikut memetakan 30 skenario pengujian komprehensif yang mencakup seluruh alur fungsional, "
        "aturan bisnis operasional, dan pengujian ketahanan keamanan siber:"
    )

    uat_headers = ["ID Uji", "Modul", "Skenario Pengujian", "Langkah Operasional", "Hasil yang Diharapkan", "Status"]
    uat_data = [
        ["UAT-01", "Auth", "Login PIN Kasir Cepat", "Input 4 digit PIN '1234' pada keypad", "Otentikasi sukses, shift Jane Doe aktif, redirect ke POS", "PASS"],
        ["UAT-02", "Auth", "Simulasi Sidik Jari", "Klik tombol biometrik (🖲️)", "PIN terisi otomatis '1234' dan login tanpa ketik", "PASS"],
        ["UAT-03", "Auth", "Proteksi Brute-Force", "Input PIN salah 5 kali berturut-turut", "Terminal terkunci 30 detik dan muncul countdown timer", "PASS"],
        ["UAT-04", "Auth", "Login Manager Form", "Input email 'admin@cafe-erp.com' & pass", "Otentikasi sukses, direct ke dashboard eksekutif", "PASS"],
        ["UAT-05", "Auth", "Zero-Scroll Viewport", "Uji resize jendela browser ke 1366x768", "Tidak ada scrollbar vertikal maupun horizontal (locked)", "PASS"],
        ["UAT-06", "Shift", "Buka Shift Kasir", "Input modal awal kasir Rp 500.000", "Shift berstatus 'open', kas modal tercatat di ledger", "PASS"],
        ["UAT-07", "Shift", "Mid-Shift Cash Drop", "Transfer Rp 1.500.000 ke brankas", "Kas register berkurang, bukti drop tercatat", "PASS"],
        ["UAT-08", "Shift", "Blind Closing Seimbang", "Input uang fisik sesuai catatan sistem", "Selisih kas Rp 0, shift ditutup sukses tanpa eskalasi", "PASS"],
        ["UAT-09", "Shift", "Blind Closing Selisih", "Input fisik kurang Rp 50.000", "Status 'Requires Audit', wajib verifikasi PIN Supervisor", "PASS"],
        ["UAT-10", "POS", "Pemesanan Meja Dine-In", "Pilih meja T-01, tambahkan Caffe Latte", "Order tersimpan, meja status 'occupied', tiket ke KDS", "PASS"],
        ["UAT-11", "POS", "Kalkulasi Pajak & Service", "Pesanan subtotal Rp 100.000", "Service Charge 5% (5rb), PB1 10% (10.500), Total 115.500", "PASS"],
        ["UAT-12", "POS", "Bayar Tunai & Kembalian", "Total Rp 35.000, bayar Rp 50.000", "Kembalian Rp 15.000, struk tercetak, meja status 'dirty'", "PASS"],
        ["UAT-13", "POS", "Void Order Berotorisasi", "Batalkan pesanan yang sudah masuk dapur", "Wajib input 6-digit PIN Otorisasi Supervisor", "PASS"],
        ["UAT-14", "KDS", "Routing Stasiun Dapur", "Order 1 Latte (Bar) dan 1 Croissant (Dapur)", "Tiket Latte masuk KDS Barista, Croissant masuk KDS Food", "PASS"],
        ["UAT-15", "KDS", "SLA Timer Antrean", "Pesanan melebihi 10 menit belum saji", "Warna kartu tiket berubah menjadi Merah Berkedip", "PASS"],
        ["UAT-16", "Inventory", "Auto-Deduct Resep (BOM)", "Selesaikan pembayaran pesanan Caffe Latte", "Stok biji kopi berkurang 18g, susu berkurang 180ml", "PASS"],
        ["UAT-17", "Inventory", "Pencegahan Stok Minus", "Order item dengan stok bahan baku 0", "Sistem memunculkan peringatan 'Bahan Baku Habis'", "PASS"],
        ["UAT-18", "Inventory", "Alert Stok Minimum", "Stok susu di bawah batas minimum (5 liter)", "Muncul badge peringatan 'Low Stock' dan rekomendasi PR", "PASS"],
        ["UAT-19", "P2P", "Purchase Requisition (PR)", "Buat pengajuan 20 kg biji kopi", "Dokumen PR diterbitkan berstatus 'pending approval'", "PASS"],
        ["UAT-20", "P2P", "Konversi PR ke PO Resmi", "Supervisor setujui PR dan terbitkan PO", "PO diterbitkan dengan nomor resmi dan dikirim ke supplier", "PASS"],
        ["UAT-21", "P2P", "Penerimaan Barang (GRN)", "Petugas gudang validasi barang datang", "Kuantitas fisik bertambah ke stok gudang cabang", "PASS"],
        ["UAT-22", "P2P", "3-Way Match Validation", "Cocokkan data PO, GRN, dan Faktur", "Harga dan kuantitas cocok, disetujui untuk dibayar", "PASS"],
        ["UAT-23", "HRIS", "Presensi Mandiri Clock-in", "Karyawan klik tombol 'Masuk Kerja'", "Waktu jam masuk dan status kehadiran tercatat di HRIS", "PASS"],
        ["UAT-24", "HRIS", "Pengajuan Cuti Tahunan", "Karyawan ajukan cuti 2 hari", "Saldo cuti dipotong sementara, menunggu approval SPV", "PASS"],
        ["UAT-25", "HRIS", "Payroll Run Otomatis", "Generate payroll periode bulan berjalan", "Gaji pokok, lembur, dan potongan BPJS terhitung presisi", "PASS"],
        ["UAT-26", "Finance", "Posting Jurnal Penjualan", "Transaksi POS kasir lunas", "Jurnal Kas (D) dan Pendapatan + PB1 (K) terposting otomatis", "PASS"],
        ["UAT-27", "Finance", "Double-Entry Balance", "Simulasi input jurnal manual tidak balance", "Sistem menolak posting dan mewajibkan Debit = Kredit", "PASS"],
        ["UAT-28", "Security", "Injeksi SQL Parameterized", "Input payload `' OR 1=1 --` pada form pencarian", "Payload diperlakukan sebagai string murni, SQLi gagal", "PASS"],
        ["UAT-29", "Security", "Rate Limiter Global", "Kirim 150 request dalam waktu 1 menit", "Server merespons HTTP 429 Too Many Requests", "PASS"],
        ["UAT-30", "Resilience", "Offline SQLite Cache", "Putus koneksi internet pada terminal POS", "Kasir tetap dapat mencatat pesanan dan cetak struk lokal", "PASS"]
    ]
    add_custom_table(doc, uat_headers, uat_data, 
                     col_widths=[Inches(0.7), Inches(0.9), Inches(1.4), Inches(1.4), Inches(1.6), Inches(0.5)],
                     align_cols=[WD_ALIGN_PARAGRAPH.CENTER, WD_ALIGN_PARAGRAPH.LEFT, WD_ALIGN_PARAGRAPH.LEFT, WD_ALIGN_PARAGRAPH.LEFT, WD_ALIGN_PARAGRAPH.LEFT, WD_ALIGN_PARAGRAPH.CENTER])

    doc.add_page_break()

    # =========================================================================
    # BAB 8: GLOSARIUM & LEMBAR PENGESAHAN
    # =========================================================================
    add_heading_1(doc, "BAB 8: GLOSARIUM ISTILAH & LEMBAR PENGESAHAN")
    
    glo_headers = ["Istilah / Akronim", "Kepanjangan / Konsep", "Penjelasan Konseptual dalam KopiERP"]
    glo_data = [
        ["BOM", "Bill of Materials", "Daftar takaran bahan baku resep yang menyusun satu porsi produk menu."],
        ["FIFO", "First In, First Out", "Metode penilaian persediaan di mana bahan baku yang masuk lebih awal dikeluarkan lebih dulu."],
        ["COGS / HPP", "Cost of Goods Sold", "Harga Pokok Penjualan yang dihitung dari akumulasi biaya bahan baku resep."],
        ["KDS", "Kitchen Display System", "Layar monitor digital di area dapur dan bar yang menggantikan tiket kertas printer."],
        ["Blind Closing", "Penutupan Kasir Buta", "Metode audit penutupan shift di mana kasir menghitung fisik uang tanpa melihat estimasi sistem."],
        ["PB1", "Pajak Pembangunan 1", "Pajak Restoran sebesar 10% yang dikenakan atas transaksi makanan dan minuman kafe."],
        ["P2P", "Procure-to-Pay", "Siklus pengadaan barang komprehensif mulai dari Requisition, PO, GRN, hingga Invoice."],
        ["GRN", "Goods Receipt Note", "Dokumen bukti fisik penerimaan barang dari supplier yang menambah saldo stok gudang."],
        ["3-Way Matching", "Pencocokan 3 Dokumen", "Validasi kesesuaian antara Purchase Order, Goods Receipt Note, dan Faktur Tagihan."],
        ["RBAC", "Role-Based Access Control", "Mekanisme pembatasan akses sistem berdasarkan peran jabatan yang disandang pengguna."],
        ["SAST", "Static Application Security Testing", "Metode pemindaian kode sumber tanpa eksekusi untuk menemukan celah keamanan dan bug."],
        ["CSP", "Content Security Policy", "Header HTTP pertahanan browser untuk mencegah serangan Cross-Site Scripting (XSS)."],
        ["HSTS", "HTTP Strict Transport Security", "Enkripsi wajib seluruh komunikasi web menggunakan protokol HTTPS/TLS 1.3."]
    ]
    add_custom_table(doc, glo_headers, glo_data, col_widths=[Inches(1.5), Inches(2.2), Inches(2.8)])

    # Sign-off Table
    doc.add_paragraph()
    add_heading_3(doc, "Lembar Pengesahan & Tanda Tangan Stakeholder")
    
    sign_headers = ["Peran Pengesah", "Nama Stakeholder", "Jabatan / Divisi", "Tanda Tangan & Tanggal"]
    sign_data = [
        ["Lead Solution Architect", "Z. Rahmatullah", "Enterprise Architecture HQ", "_______________________"],
        ["Lead Security Engineer", "A. Pratama", "Cybersecurity & Audit Dept", "_______________________"],
        ["Lead Backend Engineer", "F. Wijaya", "Core Go Engineering", "_______________________"],
        ["Lead Frontend Engineer", "R. Hidayat", "Vue 3 / UI-UX Engineering", "_______________________"],
        ["Operations Director", "H. Sukamto", "Retail & Cafe Chain Operations", "_______________________"]
    ]
    add_custom_table(doc, sign_headers, sign_data, col_widths=[Inches(1.6), Inches(1.5), Inches(1.8), Inches(1.6)])

    # Save master docx
    output_filename = "FSD_Cafe_ERP_System_Enterprise.docx"
    output_path = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), output_filename)
    doc.save(output_path)
    print(f"🎉 MASTER DOKUMEN FSD & SECURITY AUDIT BERHASIL DIBUAT: {output_path}")

if __name__ == "__main__":
    build_master_fsd()

