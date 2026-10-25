# -*- coding: utf-8 -*-
"""
Main FSD Generator Script for Cafe ERP Management System
Generates: FSD_Cafe_ERP_System_Enterprise.docx
"""

import sys
import os
import docx
from docx.shared import Inches, Pt, RGBColor
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.enum.table import WD_TABLE_ALIGNMENT

# Import helpers from generate_fsd_helpers
from generate_fsd_helpers import (
    add_heading_1, add_heading_2, add_heading_3,
    add_body_p, add_callout, add_custom_table,
    add_mockup_box, add_flowchart, format_paragraph,
    set_cell_background, set_cell_margins, set_cell_borders
)

def build_fsd():
    print("Membangun dokumen FSD Cafe ERP System...")
    doc = docx.Document()
    
    # Page setup - 1 inch margins all around
    sections = doc.sections
    for section in sections:
        section.top_margin = Inches(1.0)
        section.bottom_margin = Inches(1.0)
        section.left_margin = Inches(1.0)
        section.right_margin = Inches(1.0)
        section.different_first_page_header_footer = True
        
        # Header & Footer
        footer = section.footer
        p_ft = footer.paragraphs[0]
        p_ft.alignment = WD_ALIGN_PARAGRAPH.RIGHT
        format_paragraph(p_ft, space_before=0, space_after=0)
        r_ft = p_ft.add_run("KopiERP Pro Enterprise v2.4.8 | Functional Specification Document (FSD)")
        r_ft.font.name = 'Segoe UI'
        r_ft.font.size = Pt(8.5)
        r_ft.font.color.rgb = RGBColor(148, 163, 184)

    # =========================================================================
    # COVER PAGE
    # =========================================================================
    p_pre = doc.add_paragraph()
    format_paragraph(p_pre, space_before=30, space_after=0)
    
    # Document Title Box Table
    cover_table = doc.add_table(rows=1, cols=1)
    cover_table.alignment = WD_TABLE_ALIGNMENT.CENTER
    cell_cover = cover_table.cell(0, 0)
    cell_cover.width = Inches(6.5)
    set_cell_background(cell_cover, "0F172A") # Navy Slate 900
    set_cell_margins(cell_cover, top=400, bottom=400, left=300, right=300)
    
    p_c1 = cell_cover.paragraphs[0]
    format_paragraph(p_c1, space_before=0, space_after=6)
    r_sys = p_c1.add_run("SPESIFIKASI KEBUTUHAN FUNGSIONAL SISTEM\nFUNCTIONAL SPECIFICATION DOCUMENT (FSD)")
    r_sys.font.name = 'Segoe UI'
    r_sys.font.size = Pt(13)
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
    r_sub = p_c3.add_run("Sistem ERP Terpadu untuk Manajemen Rantai Kafe, POS Kasir Multi-Shift, Kalibrasi IoT Terminal Kopi, Kitchen Display (KDS), Rantai Suplai (P2P), HRIS Presensi, dan Pembukuan Akuntansi Real-Time.")
    r_sub.font.name = 'Segoe UI'
    r_sub.font.size = Pt(9.5)
    r_sub.font.color.rgb = RGBColor(203, 213, 225)
    
    doc.add_paragraph() # Spacing
    
    # Metadata Table
    meta_headers = ["Informasi Dokumen", "Deskripsi / Atribut"]
    meta_data = [
        ["Nomor Dokumen", "FSD-CERP-ENT-2026-V2.4"],
        ["Versi Sistem", "2.4.8 (Enterprise Production Release)"],
        ["Tanggal Terbit", "25 Oktober 2026"],
        ["Klasifikasi Akses", "Confidential / Internal Enterprise Document"],
        ["Target Arsitektur", "Backend Go 1.27 (Chi + PostgreSQL + Standalone Redis) | Frontend Vue 3 (Vite + Tailwind CSS + Naive UI)"],
        ["Database Utama", "PostgreSQL 16 (36 Relational Tables, UUID v4 Keys, ACID Compliant)"],
        ["Offline Engine", "Edge SQLite Cache Synchronization for Zero-Downtime POS"],
        ["Penulis Dokumen", "Enterprise Solution Architect & Lead Software Engineering Team"],
        ["Status Verifikasi", "Approved & Verified for Full Production Deployment"]
    ]
    add_custom_table(doc, meta_headers, meta_data, col_widths=[Inches(2.2), Inches(4.3)])
    
    # Document Version Control Table
    doc.add_paragraph()
    add_heading_3(doc, "Riwayat Perubahan & Kontrol Versi Dokumen")
    rev_headers = ["Versi", "Tanggal", "Penyusun", "Ringkasan Perubahan Utama"]
    rev_data = [
        ["v1.0.0", "04 Okt 2026", "Tech Architect", "Penyusunan baseline PRD & initial architecture framework."],
        ["v2.0.0", "11 Okt 2026", "Senior Engineer", "Implementasi Layer-7 Smart Load Balancer & Seed akun role perbankan."],
        ["v2.2.0", "11 Okt 2026", "Security Lead", "Penguatan standar OWASP ASVS V3, Rate Limiting & brute-force lockout."],
        ["v2.4.0", "11 Okt 2026", "Frontend Team", "Redesain UI Login Dual-Panel, Keypad PIN Kasir 4-digit & Zero-Scroll."],
        ["v2.4.8", "25 Okt 2026", "Enterprise Architect", "Rilis FSD komprehensif, Flowchart, Wireframe, Formula Bisnis & UAT."]
    ]
    add_custom_table(doc, rev_headers, rev_data, col_widths=[Inches(0.9), Inches(1.1), Inches(1.4), Inches(3.1)])
    
    doc.add_page_break()

    # =========================================================================
    # BAB 1: PENDAHULUAN & GAMBARAN UMUM SISTEM
    # =========================================================================
    add_heading_1(doc, "BAB 1: PENDAHULUAN & GAMBARAN UMUM SISTEM")
    
    add_heading_2(doc, "1.1 Latar Belakang & Visi Produk")
    add_body_p(doc, 
        "Industri Food & Beverage (F&B) komersial modern, khususnya jaringan rantai kafe kopi dan restoran cepat saji, "
        "menghadapi tantangan operasional yang kompleks yang membutuhkan koordinasi tanpa jeda antara area Front-of-House (kasir dan pelayanan meja) "
        "dengan area Back-of-House (barista, dapur, gudang bahan baku, dan bagian keuangan). "
        "Keterlambatan pemesanan, kebocoran uang kas pada pergantian shift, selisih persediaan biji kopi dan susu, "
        "serta keterlambatan pelaporan laba rugi harian dapat mengakibatkan penurunan margin bisnis yang signifikan."
    )
    add_body_p(doc,
        "Sistem Cafe Monitoring & Enterprise ERP System (KopiERP Pro v2.4.8) dirancang sebagai solusi terpadu berbasis cloud dengan kemampuan "
        "offline edge resilience yang menggabungkan seluruh siklus operasional bisnis F&B ke dalam satu pintu komando tersinkronisasi."
    )
    
    add_callout(doc, "Visi Sistem KopiERP Enterprise",
        "Menghadirkan ekosistem ERP dan POS berkecepatan tinggi (< 200ms latency) yang menjamin 100% akuntabilitas arus kas kasir, "
        "otomasi pemotongan stok bahan baku berbasis resep (BOM), koordinasi pesanan dapur tanpa kertas (KDS), "
        "serta pembukuan akuntansi double-entry otomatis untuk setiap transaksi yang terjadi di seluruh gerai cabang.",
        alert_type="info"
    )

    add_heading_2(doc, "1.2 Ruang Lingkup Sistem (Scope of Work)")
    add_body_p(doc, "Sistem KopiERP Pro mencakup 9 (sembilan) domain fungsional utama yang saling terintegrasi erat:")
    
    scope_data = [
        ["1", "Autentikasi & Keamanan Sesi", "Dual-panel enterprise login, 4-digit PIN Kasir, Keypad Numerik touchscreen, Biometrik simulasi, Role switching, dan proteksi anti brute-force OWASP."],
        ["2", "POS & Manajemen Meja", "Pencatatan order dine-in/takeaway/delivery, visual floor plan editor (zonasi meja), split bill, kalkulasi diskon, pajak PB1 (10%), service charge (5%), dan cetak struk ESC/POS 80mm."],
        ["3", "Shift Kasir & Cash Balancing", "Otorisasi shift kasir (F-01), setoran modal awal (opening float), mid-shift cash drops, blind cash closing, kalkulasi otomatis selisih kas, dan approval supervisor PIN."],
        ["4", "Kitchen Display System (KDS)", "Antrean tiket dapur berbasis real-time status (Pending -> In-Prep -> Ready -> Served), pemisahan stasiun barista kopi vs dapur makanan, timer indikator keterlambatan."],
        ["5", "Resep (BOM) & Inventory Stok", "Penyusunan Bill of Materials bahan baku, auto-deduction stok saat transaksi bayar, alert batas minimum stok, kartu stok FIFO ledger, dan stock opname fisik."],
        ["6", "Rantai Pasok (Procure-to-Pay)", "Purchase Requisition (PR), auto-PR saat menu menipis, Purchase Order (PO), Goods Receipt Note (GRN), 3-Way Matching Invoice Vendor, dan riwayat pembayaran."],
        ["7", "HRIS & Penggajian (Payroll)", "Database karyawan multi-cabang, presensi masuk/pulang mandiri, penjadwalan shift kerja, persetujuan cuti, payroll run otomatis (gaji pokok, lembur, BPJS, potongan)."],
        ["8", "Keuangan & Buku Besar (Finance)", "Chart of Accounts (COA) 5 level, automated journal posting dari transaksi POS/HPP/Payroll, jurnal manual double-entry, ringkasan pajak, dan rekonsiliasi bank."],
        ["9", "Master Data Hub & Laporan", "Konfigurasi cabang outlet, manajemen peran dan matriks izin RBAC, master menu & kategori, laporan penjualan jam sibuk, P&L statement, dan export data."]
    ]
    add_custom_table(doc, ["No", "Modul Sistem", "Cakupan Fungsional Utama"], scope_data, col_widths=[Inches(0.5), Inches(1.8), Inches(4.2)])

    add_heading_2(doc, "1.3 Karakteristik Persona & Pengguna Sistem")
    persona_data = [
        ["Super Admin", "IT HQ / Owner", "Akses 100% penuh tanpa batas ke seluruh modul, database audit, otorisasi cabang, dan konfigurasi master data sistem."],
        ["Store Manager / SPV", "Outlet Manager", "Supervisi shift kasir, override void/diskon, verifikasi selisih kas blind closing, persetujuan PR & cuti karyawan."],
        ["Kasir POS", "Front Counter Staff", "Membuka shift kasir, melayani pembayaran tamu, cetak struk, cash drop ke brankas, dan menghitung fisik uang saat blind closing."],
        ["Barista / Koki Dapur", "Food & Beverage Maker", "Melihat tiket KDS, menandai status makanan/minuman yang siap disajikan, memonitor kalibrasi suhu boiler espresso."],
        ["Staff Gudang / Inventory", "Back Office / Logistics", "Penerimaan barang dari supplier (GRN), pencatatan mutasi stok antar gudang, pelaksanaan stock opname fisik berkala."],
        ["Finance & Akuntan", "Finance Dept", "Memverifikasi jurnal transaksi otomatis, mencatat pengeluaran petty cash, memproses faktur supplier, rekonsiliasi bank, dan pajak."],
        ["HRIS Specialist", "People & Operations", "Pengelolaan data karyawan, pengaturan jadwal shift mingguan, persetujuan cuti, dan peninjauan slip gaji (payroll run)."]
    ]
    add_custom_table(doc, ["Peran Pengguna", "Unit Kerja", "Tanggung Jawab Operasional"], persona_data, col_widths=[Inches(1.6), Inches(1.5), Inches(3.4)])

    add_heading_2(doc, "1.4 Arsitektur High-Level & Spesifikasi Teknologi")
    tech_data = [
        ["Komponen Arsitektur", "Teknologi yang Digunakan", "Versi / Standar", "Tujuan & Peran dalam Sistem"],
        ["Backend Core Engine", "Golang (Go)", "1.27.0", "Kompilasi biner tinggi, concurrency goroutine ringan untuk high-throughput API."],
        ["HTTP Routing Engine", "Chi Router v5", "5.0.12", "RESTful router modular dengan integrasi middleware chain OWASP."],
        ["Database Relasional", "PostgreSQL", "16-Alpine", "Penyimpanan data relasional ACID, UUID v4 primary keys, dan JSONB logging."],
        ["Caching & Session", "Redis Standalone", "7-Alpine", "Penyimpanan token blacklist, rate limit sliding window, dan session cache."],
        ["Frontend UI Framework", "Vue 3 (Composition API)", "3.5.x", "Reaktivitas mutakhir, type-safety TypeScript, modular SFC."],
        ["State Management", "Pinia", "2.1.x", "Penyimpanan state terpusat: auth, app theme, notification, and cart."],
        ["Styling & Design System", "Tailwind CSS + Naive UI", "3.4.x / 2.38.x", "Zero-scroll responsive layouts, WCAG AA compliant contrast design tokens."],
        ["Edge Resilience", "SQLite Local Storage Sync", "3.x", "Pencadangan transaksi kasir secara lokal saat terjadi kehilangan koneksi internet."]
    ]
    add_custom_table(doc, tech_data[0], tech_data[1:], col_widths=[Inches(1.5), Inches(1.5), Inches(1.1), Inches(2.4)])

    doc.add_page_break()

    # =========================================================================
    # BAB 2: MATRIKS PERAN & HAK AKSES (RBAC MATRIX)
    # =========================================================================
    add_heading_1(doc, "BAB 2: MATRIKS PERAN & HAK AKSES (ROLE-BASED ACCESS CONTROL)")
    add_body_p(doc,
        "Sistem KopiERP menerapkan prinsip keamanan Least Privilege Access Control di mana setiap peran hanya memiliki hak akses "
        "yang sesuai dengan tugas operasionalnya. Hak akses dikontrol pada level endpoint API Backend melalui JWT Claims dan "
        "pada level Frontend melalui Vue Router Navigation Guards dan conditional element rendering."
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

    add_callout(doc, "Aturan Otorisasi Khusus (Supervisor Override)",
        "1. Pembatalan pesanan yang telah dikirim ke dapur (Void Order) atau pemberian diskon di atas 20% "
        "wajib memasukkan 6-digit PIN Otorisasi Supervisor / Store Manager.\n"
        "2. Penutupan shift kasir (Blind Closing) dengan selisih kas fisik melebihi ambang toleransi (+/- Rp 10.000) "
        "otomatis memicu status 'Requires Audit' dan mewajibkan tanda tangan digital serta input PIN Supervisor.",
        alert_type="warning"
    )

    doc.add_page_break()

    # =========================================================================
    # BAB 3: ARSITEKTUR DATA & SKEMA RELASIONAL
    # =========================================================================
    add_heading_1(doc, "BAB 3: ARSITEKTUR DATA & SKEMA DATABASE RELASIONAL")
    add_body_p(doc,
        "Database KopiERP dibangun di atas PostgreSQL 16 dengan rancangan normalisasi 3NF, didukung indeks komposit "
        "untuk kecepatan query transaksi harian dan integritas referensial foreign-key cascading yang ketat."
    )
    
    db_group_headers = ["Kelompok Entitas", "Tabel Database", "Fungsi & Peran Relasional"]
    db_group_data = [
        ["1. Organisasi & Cabang", "branches, departments, positions", "Struktur organisasi multi-cabang, hierarki divisi, dan jenjang jabatan."],
        ["2. Hak Akses & Pengguna", "roles, permissions, role_permissions, users, user_roles, audit_logs", "Manajemen kredensial login, PIN code, audit trail perubahan data dengan format JSONB."],
        ["3. Tata Letak Kafe", "table_zones, cafe_tables", "Denah visual kafe, zonasi meja (Lantai 1, Mezzanine, Outdoor), koordinat X-Y, dan status meja."],
        ["4. Menu & Resep (BOM)", "menu_categories, products, product_variants, product_recipes", "Katalog produk minuman/makanan, varian ukuran/rasa, dan Bill of Materials komposisi bahan baku."],
        ["5. Inventori & Rantai Suplai", "inventory_items, warehouses, inventory_stocks, stock_movements, stock_opnames, stock_opname_items", "Master bahan baku, lokasi gudang outlet, pencatatan mutasi masuk/keluar FIFO, dan kartu stok fisik."],
        ["6. Pengadaan Barang (P2P)", "suppliers, purchase_orders, purchase_order_items", "Database vendor supplier, order pembelian bahan baku, dan riwayat penerimaan barang."],
        ["7. Transaksi POS & Kasir", "orders, order_items, payments, cashier_shifts, cashier_shift_movements", "Pencatatan order tamu, item modifier, mutasi pembayaran tunai/non-tunai, dan pembukuan shift kasir."],
        ["8. Sumber Daya Manusia (HRIS)", "employees, work_shifts, attendances, leaves, payroll_records", "Data kepegawaian, pencatatan presensi GPS/biometrik, jadwal shift kerja, saldo cuti, dan slip gaji bulanan."],
        ["9. Akuntansi & General Ledger", "chart_of_accounts, journal_entries, journal_entry_lines", "Struktur bagan akun (COA), jurnal otomatis double-entry dari penjualan/COGS/payroll, dan audit balance."]
    ]
    add_custom_table(doc, db_group_headers, db_group_data, col_widths=[Inches(1.8), Inches(2.2), Inches(2.5)])

    add_callout(doc, "Standar Kepatuhan Data (Data Integrity Standards)",
        "1. Semua primary key menggunakan UUID v4 (`uuid_generate_v4()`) untuk mencegah enumerasi ID pada antarmuka publik.\n"
        "2. Penghapusan data master menggunakan Soft-Delete (`deleted_at TIMESTAMPTZ`) sehingga data transaksi historis tetap utuh untuk audit pajak.\n"
        "3. Setiap transaksi mutasi finansial dan stok dibungkus dalam blok Database Transaction (`BEGIN ... COMMIT`) terisolasi.",
        alert_type="info"
    )

    doc.add_page_break()

    # =========================================================================
    # BAB 4: SPESIFIKASI FUNGSIONAL PER MODUL
    # =========================================================================
    add_heading_1(doc, "BAB 4: SPESIFIKASI FUNGSIONAL PER MODUL")

    # -------------------------------------------------------------------------
    # 4.1 Modul Autentikasi & Keamanan Sesi
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.1 Modul Autentikasi & Keamanan Sesi")
    add_body_p(doc,
        "Modul ini bertindak sebagai gerbang otorisasi utama sistem dengan mengadopsi arsitektur Dual-Panel: "
        "Panel Kiri menyajikan informasi telemetri operasional gerai kafe (latensi cloud, status terminal IoT, shift aktif, metrik suhu mesin kopi), "
        "sedangkan Panel Kanan menyajikan antarmuka otorisasi kasir berbasis 4-digit PIN dengan keypad numerik layar sentuh, "
        "serta formulir email/password untuk level manajerial dan eksekutif."
    )
    
    auth_flowchart_ascii = """
+-----------------------------------------------------------------------------+
|                          ALUR AUTENTIKASI PENGGUNA                          |
+-----------------------------------------------------------------------------+
         |
         v
+------------------+         Pilih Tab
| Halaman Login    | ------------------------------------+
| KopiERP Ent.     |                                     |
+------------------+                                     |
         |                                               |
         | [Tab Barista / Kasir]                         | [Tab Manager / Spv]
         v                                               v
+-----------------------+                       +-----------------------+
| Input 4-Digit PIN     |                       | Input Email & Sandi   |
| (Keypad / Biometrik)  |                       | Atau Pilih Demo Role  |
+-----------------------+                       +-----------------------+
         |                                               |
         v                                               v
+-----------------------+                       +-----------------------+
| Validasi PIN Kasir:   |                       | Validasi Kredensial:  |
| 1234 -> Kasir Jane    |                       | POST /api/v1/auth/    |
| 1111 -> Barista Budi  |                       | login (Bcrypt Hash)   |
+-----------------------+                       +-----------------------+
         |                                               |
         +-----------------------+-----------------------+
                                 |
                                 v
                     +-----------------------+
                     | Kredensial Valid?     |
                     +-----------------------+
                      /                     \\
                [Ya] /                       \\ [Tidak]
                    v                         v
       +-----------------------+   +-----------------------+
       | Terbitkan JWT Token   |   | Failed Attempts += 1  |
       | & Refresh HttpOnly    |   | Tampilkan Pesan Error |
       +-----------------------+   +-----------------------+
                    |                         |
                    v                         v
       +-----------------------+   +-----------------------+
       | Redirect ke Tujuan:   |   | Percobaan >= 5x?      |
       | Kasir -> /pos         |   | Aktifkan Lockout 30s! |
       | Admin -> /dashboard   |   +-----------------------+
       +-----------------------+
"""
    auth_steps = [
        ["1. Pilih Metode", "Pengguna", "Klik Tab Kasir / Manager", "Sistem menampilkan keypad PIN atau input email/password", "UI terfokus sesuai pilihan"],
        ["2. Input Otorisasi", "Kasir / SPV", "PIN 4 digit / Email & Pass", "Cek panjang karakter (PIN min 4 digit, Pass min 6 char)", "Tombol login aktif"],
        ["3. Verifikasi Server", "Backend API", "Payload login request", "Bcrypt compare password hash, periksa status is_active", "200 OK / 401 Unauthorized"],
        ["4. Penerbitan Token", "Auth Service", "Kredensial valid", "Generate JWT access token (24 jam) & cookie refresh_token", "Sesi login tersimpan di Pinia"],
        ["5. Pengalihan Rute", "Router Guard", "Role user login", "Role kasir dialihkan ke /pos, role manajerial ke /dashboard", "Pengguna masuk ke sistem"]
    ]
    add_flowchart(doc, "Otorisasi Login & Perlindungan Brute-Force", auth_flowchart_ascii, auth_steps)

    auth_rules = (
        "1. Rate Limiting Otentikasi: Endpoint `/api/v1/auth/login` dibatasi maksimal 5 request gagal per IP dalam 1 menit.\n"
        "2. Lockout Proteksi: Jika terjadi 5 kali kegagalan login berturut-turut pada terminal, sistem secara otomatis "
        "mengunci terminal selama 30 detik (`lockoutSeconds = 30`) disertai notifikasi peringatan keamanan.\n"
        "3. Zero-Scroll Constraint: Container login dikunci menggunakan CSS `fixed inset-0 w-full h-full overflow-hidden` "
        "sehingga tidak memunculkan scrollbar horizontal maupun vertikal pada resolusi layar apa pun."
    )
    add_callout(doc, "Aturan Bisnis & Validasi Autentikasi", auth_rules, alert_type="warning")

    auth_mockup_ascii = """
+-----------------------------------------------------------------------------------+
|  [Logo] KopiERP PRO ENTERPRISE                    (•) Sync Cloud Aktif [Theme: ☀] |
|                                                                                   |
|  [ PANEL KIRI: OPERASIONAL GERAI ]        [ PANEL KANAN: OTORISASI KASIR ]        |
|  ---------------------------------        --------------------------------        |
|  Sistem ERP Terpadu Kafe & Resto          [ [💻 Barista / Kasir] | [🛡️ Manager] ] |
|                                                                                   |
|  [Gerai: Sudirman Flagship #01   (v)]     Otorisasi Shift Kasir                   |
|  Online • Latensi 24ms • Cloud ID-02      Masukkan 4-digit PIN pegawai register   |
|                                                                                   |
|  +----------------+ +------------------+            ( * ) ( * ) ( * ) ( * )       |
|  | Station POS-01 | | SHIFT: Pagi      |                                          |
|  | Terhubung      | | Modal: Rp 500rb  |          [ 1 ]     [ 2 ]     [ 3 ]       |
|  | KDS: Aktif     | | Boiler: 92.5°C   |          [ 4 ]     [ 5 ]     [ 6 ]       |
|  +----------------+ +------------------+          [ 7 ]     [ 8 ]     [ 9 ]       |
|                                                   [ 🖲️ ]     [ 0 ]     [ ⌫ ]       |
|  [ 542 Transaksi | 18/25 Meja | Susu 48L ]                                        |
|                                            [ 🔓 Buka Shift & Masuk POS Kasir  ]   |
|                                                                                   |
|                                            ⚡ Mode Offline (SQLite)   🎧 Bantuan   |
+-----------------------------------------------------------------------------------+
"""
    auth_controls = [
        ["Tab Barista / Kasir", "Segmented Button", "Klik untuk beralih mode kasir", "Menampilkan 4 dot indikator PIN dan keypad angka"],
        ["Keypad 1-9, 0", "Virtual Key Button", "Klik angka atau ketik keyboard fisik", "Mengisi digit PIN kasir secara berurutan"],
        ["Tombol Sidik Jari (🖲️)", "Biometric Simulation", "1-Klik otorisasi cepat kasir", "Mengisi PIN '1234' otomatis dan submit instan"],
        ["Tombol Backspace (⌫)", "Action Button", "Klik untuk menghapus digit terakhir", "Menghapus digit dan mereset status error"],
        ["Tombol Buka Shift", "Primary Action", "Submit PIN untuk membuka register", "Memanggil API auth dan redirect langsung ke /pos"],
        ["Tombol Ganti Gerai", "Modal Trigger", "Klik untuk beralih outlet cabang", "Membuka modal daftar cabang dan mengupdate koneksi"]
    ]
    add_mockup_box(doc, "Halaman Login Dual-Panel KopiERP Enterprise", auth_mockup_ascii, auth_controls)

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # 4.2 Modul Manajemen Shift Kasir & Blind Cash Balancing (F-01)
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.2 Modul Manajemen Shift Kasir & Blind Cash Balancing (F-01)")
    add_body_p(doc,
        "Fitur F-01 dirancang untuk menutup celah kebocoran uang kas kasir di gerai. Setiap kasir wajib membuka shift dengan "
        "memasukkan modal awal kas fisik (opening float). Selama shift berlangsung, kasir dapat melakukan mutasi kas keluar "
        "(cash drop ke brankas atau pembayaran operasional kecil). Pada akhir shift, kasir wajib melakukan 'Blind Closing' "
        "di mana sistem menyembunyikan estimasi total penjualan tunai, dan kasir wajib menghitung uang fisik riil di laci register."
    )

    shift_flowchart_ascii = """
+-----------------------------------------------------------------------------+
|               ALUR SIKLUS SHIFT KASIR & BLIND CASH BALANCING                |
+-----------------------------------------------------------------------------+
         |
         v
+----------------------------------+
| Kasir Datang di Awal Jam Kerja   |
+----------------------------------+
         |
         v
+----------------------------------+
| Buka Shift Kasir                 |
| Input Modal Awal (Opening Float) |
+----------------------------------+
         |
         v
+----------------------------------+
| Shift Berstatus: 'open'          |
| Kasir Melayani Transaksi POS     |
+----------------------------------+
         |
         +-------------------------------------+
         |                                     |
         v                                     v
+-----------------------+             +-----------------------+
| Transaksi Penjualan:  |             | Mutasi Kas Berjalan:  |
| - Penjualan Kas (Tunai|             | - Cash Drop Brankas   |
| - QRIS / EDC Card     |             | - Pengeluaran Kasir   |
+-----------------------+             +-----------------------+
         |                                     |
         +------------------+------------------+
                            |
                            v
+------------------------------------------------------+
| Pergantian Shift / Tutup Shift (Blind Closing)       |
| * Sistem MENYEMBUNYIKAN estimasi saldo kas sistem    |
| * Kasir wajib menghitung uang kertas & koin riil     |
+------------------------------------------------------+
                            |
                            v
+------------------------------------------------------+
| Input Jumlah Uang Fisik (Actual Cash Counted)        |
+------------------------------------------------------+
                            |
                            v
+------------------------------------------------------+
| Kalkulasi Selisih Otomatis:                          |
| Expected = Modal + Kas Masuk - Drop - Refund         |
| Selisih  = Actual Counted - Expected Cash            |
+------------------------------------------------------+
                            |
             +--------------+--------------+
             |                             |
      [Selisih = 0]                 [Ada Selisih]
             v                             v
+--------------------------+  +----------------------------------+
| Shift Selesai & Seimbang |  | Selisih Melebihi Toleransi?     |
| Status: 'closed'         |  | Wajib Approval & Catatan SPV     |
+--------------------------+  +----------------------------------+
"""
    shift_steps = [
        ["1. Buka Shift", "Kasir", "Input nominal modal awal kas", "INSERT into cashier_shifts status='open'", "Shift ID aktif diterbitkan"],
        ["2. Mutasi Penjualan", "POS Engine", "Pembayaran order kas", "UPDATE cashier_shifts total_cash_sales", "Saldo expected diupdate sistem"],
        ["3. Cash Drop", "Kasir & SPV", "Setoran uang ke brankas", "INSERT into cashier_shift_movements", "Kas laci register berkurang"],
        ["4. Blind Closing", "Kasir", "Input uang fisik riil", "Sistem membandingkan actual vs expected", "Nilai selisih kas dihitung"],
        ["5. Verifikasi SPV", "Store Manager", "PIN otorisasi jika selisih", "UPDATE status='closed', catatan selisih", "Shift ditutup & siap diaudit"]
    ]
    add_flowchart(doc, "Siklus Shift Kasir & Blind Closing", shift_flowchart_ascii, shift_steps)

    shift_formula = (
        "Rumus Rekonsiliasi Kas Kasir:\n"
        "1. Saldo Kas Sistem yang Diharapkan (Expected Cash Total):\n"
        "   Expected Cash = Opening Cash Float + Total Cash Sales - Total Cash Drops - Total Refunds\n\n"
        "2. Nilai Selisih Kas (Cash Difference / Discrepancy):\n"
        "   Cash Difference = Actual Cash Counted - Expected Cash Total\n\n"
        "3. Interpretasi Hasil:\n"
        "   - Cash Difference = 0 : Kas Seimbang (Balanced / Perfect Match).\n"
        "   - Cash Difference > 0 : Kas Lebih (Overage) - uang fisik melebihi catatan sistem.\n"
        "   - Cash Difference < 0 : Kas Kurang (Shortage) - potensi kebocoran uang kas register."
    )
    add_callout(doc, "Formula Matematis & Aturan Selisih Kas (Blind Closing)", shift_formula, alert_type="success")

    shift_mockup_ascii = """
+-----------------------------------------------------------------------------------+
|  [KopiERP POS] MANAJEMEN SHIFT KASIR                   Shift Aktif: #SHF-001      |
|                                                                                   |
|  Kasir: Jane Doe (Kasir #A1) | Gerai: Sudirman #01 | Waktu Buka: 06:30:15 WIB    |
|  -------------------------------------------------------------------------------- |
|  [ MODAL AWAL REGISTER ]               [ MUTASI KAS BERJALAN ]                    |
|  Rp 500.000                            - Cash Drop #1: Rp 1.500.000 (11:30)       |
|                                        - Pembelian Es Batu: Rp 25.000 (13:10)     |
|  [ FORMULIR BLIND CLOSING (PENGHITUNGAN UANG FISIK) ]                             |
|  * PERHATIAN: Masukkan hasil hitungan riil uang kas fisik di laci register.       |
|                                                                                   |
|  Pecahan Uang:           Jumlah Lembar:           Subtotal:                       |
|  [ Rp 100.000 ]        x [     25      ]        = Rp 2.500.000                    |
|  [ Rp  50.000 ]        x [     18      ]        = Rp   900.000                    |
|  [ Rp  20.000 ]        x [     15      ]        = Rp   300.000                    |
|  [ Rp  10.000 ]        x [     12      ]        = Rp   120.000                    |
|  [ Koin/Receh ]        x [    Total    ]        = Rp    35.000                    |
|  -------------------------------------------------------------------------------- |
|  TOTAL UANG FISIK DIHITUNG: Rp 3.855.000                                          |
|                                                                                   |
|  Catatan Kasir: [ Shift berjalan lancar, tidak ada kendala mesin EDC           ] |
|                                                                                   |
|  [ Kunci Perhitungan & Hitung Selisih ]        [ 🔴 Tutup & Akhiri Shift Sesi ]   |
+-----------------------------------------------------------------------------------+
"""
    shift_controls = [
        ["Modal Awal (Float)", "Numeric Input", "Input nominal uang kembalian awal", "Tervalidasi min Rp 0, wajib sebelum transaksi POS"],
        ["Tombol Cash Drop", "Button Trigger", "Pindahkan uang dari laci ke brankas", "Memotong saldo kas laci untuk keamanan outlet"],
        ["Tabel Pecahan Uang", "Denomination Grid", "Input jumlah lembar uang fisik", "Menghitung total uang riil secara otomatis"],
        ["Tombol Tutup Shift", "Primary Button", "Submit hasil hitungan kasir", "Menjalankan kalkulasi selisih dan mengunci sesi register"]
    ]
    add_mockup_box(doc, "Layar Manajemen Shift Kasir & Blind Closing", shift_mockup_ascii, shift_controls)

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # 4.3 Modul Transaksi POS (Point of Sale) & Table Map
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.3 Modul Transaksi POS & Manajemen Meja (Point of Sale)")
    add_body_p(doc,
        "Modul POS melayani alur pemesanan pesanan makan di tempat (Dine-In), bawa pulang (Takeaway), maupun pesan antar (Delivery). "
        "Untuk pesanan dine-in, sistem menyediakan denah visual meja interaktif yang menampilkan status ketersediaan meja "
        "(Tersedia, Terisi, Dipesan, atau Kotor/Perlu Dibersihkan). POS terintegrasi langsung dengan pemotongan stok bahan baku "
        "dan otomatis mengirimkan tiket ke layar KDS dapur."
    )

    pos_flowchart_ascii = """
+-----------------------------------------------------------------------------+
|                      ALUR TRANSAKSI PEMESANAN POS KAFE                      |
+-----------------------------------------------------------------------------+
         |
         v
+----------------------------------+
| Pilih Tipe Pesanan:              |
| [ Dine-In ] / [ Takeaway ]       |
+----------------------------------+
         |
         +-------------------+
         | [Dine-In]         | [Takeaway]
         v                   v
+------------------+   +-------------------+
| Pilih Meja di    |   | Langsung ke Menu  |
| Visual Floor Map |   | Input Nama Tamu   |
+------------------+   +-------------------+
         |                   |
         +---------+---------+
                   |
                   v
+--------------------------------------+
| Pilih Menu Minuman / Makanan         |
| Atur Varian (Hot/Iced, Size, Sugar)  |
+--------------------------------------+
                   |
                   v
+--------------------------------------+
| Simpan Pesanan (Status: 'pending')   |
| Tiket Pesanan Masuk ke Antrean KDS   |
| Status Meja Berubah: 'occupied'      |
+--------------------------------------+
                   |
                   v
+--------------------------------------+
| Proses Pembayaran Kasir:             |
| - Tunai (Cash)                       |
| - QRIS Dinamis                       |
| - Kartu Debit / EDC Kredit           |
+--------------------------------------+
                   |
                   v
+--------------------------------------+
| Pembayaran Berhasil:                 |
| 1. Cetak Struk Bukti Bayar 80mm      |
| 2. Auto-Deduct Stok Bahan Baku (BOM) |
| 3. Auto-Posting Jurnal Akuntansi     |
| 4. Status Meja: 'dirty' -> 'clean'   |
+--------------------------------------+
"""
    pos_steps = [
        ["1. Pemilihan Meja", "Kasir / Waiter", "Klik nomor meja di visual map", "UPDATE cafe_tables status='occupied'", "Meja terkunci pada order aktif"],
        ["2. Tambah Item Menu", "Kasir", "Pilih produk dan modifier varian", "Hitung subtotal + harga modifier", "Ringkasan keranjang terupdate"],
        ["3. Simpan Pesanan", "POS Service", "Klik tombol 'Kirim ke Dapur'", "INSERT into orders & order_items", "KDS tiket otomatis muncul di dapur"],
        ["4. Pembayaran", "Kasir", "Pilih metode bayar & input nominal", "Kalkulasi Pajak PB1 10%, Service 5%, Kembalian", "Status order berubah jadi 'paid'"],
        ["5. Pasca Bayar", "Sistem ERP", "Konfirmasi pembayaran sukses", "Potong stok bahan baku & posting jurnal penjualan", "Struk tercetak, meja status 'available'"]
    ]
    add_flowchart(doc, "Alur Transaksi Pemesanan & Pembayaran POS", pos_flowchart_ascii, pos_steps)

    pos_calc = (
        "Rumus Perhitungan Transaksi POS:\n"
        "1. Subtotal = Σ (Harga Satuan Item + Harga Tambahan Modifier) x Jumlah Pesanan\n"
        "2. Nilai Diskon = (Subtotal x Persentase Diskon) ATAU Diskon Nominal Tetap\n"
        "3. Dasar Pengenaan Pajak (DPP) = Subtotal - Nilai Diskon\n"
        "4. Biaya Layanan (Service Charge) = DPP x 5%\n"
        "5. Pajak Restoran (PB1) = (DPP + Service Charge) x 10%\n"
        "6. Grand Total = DPP + Service Charge + Pajak Restoran (PB1)\n"
        "7. Uang Kembalian = Uang Diterima Tunai - Grand Total (Wajib >= 0)"
    )
    add_callout(doc, "Aturan Perhitungan Finansial Transaksi POS", pos_calc, alert_type="info")

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # 4.4 Modul Kitchen Display System (KDS)
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.4 Modul Kitchen Display System (KDS)")
    add_body_p(doc,
        "Kitchen Display System (KDS) menggantikan printer kertas dapur tradisional dengan layar interaktif nirkabel "
        "yang menampilkan pesanan secara real-time. Modul ini membedakan stasiun peracikan minuman (Barista Bar) "
        "dengan stasiun pengolahan makanan hangat (Kitchen Food) dan dilengkapi timer otomatis untuk menjaga standar waktu penyajian (SLA)."
    )

    kds_rules = (
        "Aturan Pewarnaan Timer Antrean Tiket KDS:\n"
        "1. Warna Hijau (SLA Aman): Durasi antrean < 5 menit sejak order dikirim kasir.\n"
        "2. Warna Kuning (Peringatan): Durasi antrean antara 5 hingga 10 menit (peringatan persiapan).\n"
        "3. Warna Merah Berkedip (SLA Terlampaui): Durasi antrean > 10 menit (eskalasi ke Supervisor Outlet).\n"
        "4. Status Progres Item: Pending (Menunggu) -> In Preparation (Sedang Dibuat) -> Ready (Siap Saji) -> Completed."
    )
    add_callout(doc, "Standar Waktu Pelayanan (SLA) Layar Dapur KDS", kds_rules, alert_type="warning")

    # -------------------------------------------------------------------------
    # 4.5 Modul Manajemen Resep (BOM) & Inventory Stok
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.5 Modul Manajemen Resep (BOM) & Inventory Stok")
    add_body_p(doc,
        "Setiap produk menu minuman atau makanan dipetakan ke dalam resep Bill of Materials (BOM) yang merinci takaran bahan baku "
        "yang dibutuhkan (misal: 1 cup Caffe Latte membutuhkan 18g Biji Kopi Espresso Blend, 180ml Susu Segar, dan 1 pcs Paper Cup 8oz). "
        "Saat kasir menyelesaikan pembayaran, sistem secara otomatis mengeksekusi mutasi pengurangan stok bahan baku secara FIFO "
        "(First In, First Out) dari gudang cabang bersangkutan."
    )

    bom_flowchart_ascii = """
+-----------------------------------------------------------------------------+
|          ALUR PEMOTONGAN STOK RESEP (BOM) & KONTROL MINIMUM STOK            |
+-----------------------------------------------------------------------------+
         |
         v
+----------------------------------+
| Kasir Menyelesaikan Pembayaran   |
| Order POS (Status: 'paid')       |
+----------------------------------+
         |
         v
+----------------------------------+
| Query Tabel 'product_recipes'    |
| Ambil Formula Takaran Bahan Baku |
+----------------------------------+
         |
         v
+----------------------------------+
| Loop Setiap Bahan Baku:          |
| Qty Terpakai = Takaran x Porsi   |
+----------------------------------+
         |
         v
+----------------------------------+
| Kurangi 'inventory_stocks'       |
| Catat 'stock_movements' (OUT)    |
| Hitung HPP Otomatis (FIFO)       |
+----------------------------------+
         |
         v
+----------------------------------+
| Sisa Stok <= Minimum Threshold?  |
+----------------------------------+
         |
         +-------------------+
         | [Ya]              | [Tidak]
         v                   v
+-------------------+   +--------------------+
| Kirim Notifikasi  |   | Mutasi Selesai     |
| Alert Stok Kritis |   | Stok Aman          |
| Trigger Auto-PR   |   +--------------------+
+-------------------+
"""
    bom_steps = [
        ["1. Event Pembayaran", "POS Gateway", "Order berubah status ke 'paid'", "Memicu background hook pengurangan stok bahan", "Stok mulai diproses"],
        ["2. Lookup Resep", "Recipe Engine", "Product ID & Variant ID", "SELECT * FROM product_recipes WHERE product_id = ?", "Daftar bahan dan takaran didapat"],
        ["3. Pengurangan Stok", "Inventory Engine", "Kuantitas resep x jumlah cup", "UPDATE inventory_stocks SET quantity = quantity - ?", "Stok gudang berkurang"],
        ["4. Pencatatan Ledger", "Stock Ledger", "Data mutasi bahan baku", "INSERT into stock_movements type='order_deduction'", "Kartu stok terupdate real-time"],
        ["5. Evaluasi Minimum", "Alert Monitor", "Sisa stok vs min_stock_alert", "Jika stok <= min, set status badge 'Depleted/Low'", "Notifikasi diterbitkan ke Staff Gudang"]
    ]
    add_flowchart(doc, "Otomasi Pengurangan Bahan Baku Berbasis Resep", bom_flowchart_ascii, bom_steps)

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # 4.6 Modul Pengadaan Rantai Pasok (Procure-to-Pay / P2P)
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.6 Modul Pengadaan Rantai Pasok (Procure-to-Pay / P2P)")
    add_body_p(doc,
        "Siklus Procure-to-Pay (P2P) mengelola seluruh alur pengadaan bahan baku kafe mulai dari pengajuan kebutuhan (Purchase Requisition), "
        "penerbitan Purchase Order (PO) resmi ke supplier terpilih, verifikasi fisik barang datang di gudang (Goods Receipt Note / GRN), "
        "hingga pencocokan faktur tagihan (3-Way Matching: PO vs GRN vs Invoice) sebelum pembayaran dilunasi oleh bagian keuangan."
    )

    p2p_flowchart_ascii = """
+-----------------------------------------------------------------------------+
|             ALUR SIKLUS PENGADAAN BARANG (P2P 3-WAY MATCHING)               |
+-----------------------------------------------------------------------------+
         |
         v
+-----------------------------------+
| Purchase Requisition (PR)         |
| Dibuat Manual / Trigger Auto Stok |
+-----------------------------------+
         |
         v
+-----------------------------------+
| Approval PR oleh Store Manager    |
| Konversi PR menjadi PO Resmi      |
+-----------------------------------+
         |
         v
+-----------------------------------+
| PO Dikirim ke Supplier Terdaftar  |
| Status PO: 'ordered'              |
+-----------------------------------+
         |
         v
+-----------------------------------+
| Supplier Kirim Barang ke Gudang   |
| Petugas Gudang Validasi Fisik     |
+-----------------------------------+
         |
         v
+-----------------------------------+
| Terbitkan Goods Receipt Note(GRN) |
| Stok Gudang Bertambah Otomatis    |
+-----------------------------------+
         |
         v
+-----------------------------------+
| 3-Way Matching oleh Finance:      |
| Match Kuantitas PO vs GRN vs Inv  |
+-----------------------------------+
         |
         v
+-----------------------------------+
| Pembayaran Dilunasi ke Supplier   |
| Jurnal Akuntansi Posting Kas/Hutang|
+-----------------------------------+
"""
    p2p_steps = [
        ["1. Pengajuan PR", "Staff Gudang", "Input item bahan & estimasi harga", "INSERT into purchase_requisitions status='pending'", "Nomor PR diterbitkan"],
        ["2. Persetujuan PR", "Store Manager", "Review kewajaran kebutuhan stok", "UPDATE purchase_requisitions status='approved'", "Siap dikonversi ke PO"],
        ["3. Penerbitan PO", "Purchasing", "Pilih supplier & tanggal kirim", "INSERT into purchase_orders status='ordered'", "Dokumen PO dikirim ke supplier"],
        ["4. Penerimaan GRN", "Staff Gudang", "Hitung fisik barang datang & batch", "INSERT into goods_receipt_notes & tambah stok", "Stok fisik resmi bertambah"],
        ["5. 3-Way Matching", "Finance", "Cocokkan PO, GRN, dan Invoice", "Verifikasi selisih harga dan kuantitas faktur", "Invoice disetujui untuk bayar"]
    ]
    add_flowchart(doc, "Alur Siklus Pengadaan P2P 3-Way Matching", p2p_flowchart_ascii, p2p_steps)

    doc.add_page_break()

    # -------------------------------------------------------------------------
    # 4.7 Modul Human Resource Information System (HRIS) & Payroll
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.7 Modul HRIS & Penggajian Otomatis (Payroll)")
    add_body_p(doc,
        "Modul HRIS mengintegrasikan manajemen sumber daya manusia gerai kafe mulai dari database karyawan, "
        "presensi masuk/pulang harian (clock-in / clock-out), penjadwalan shift kerja mingguan, pengajuan dan persetujuan cuti, "
        "hingga penghitungan slip gaji bulanan (Payroll Run) yang otomatis memotong keterlambatan, absensi tanpa keterangan, dan iuran BPJS."
    )

    payroll_calc = (
        "Formula Perhitungan Gaji Bersih Karyawan (Take Home Pay):\n"
        "1. Penghasilan Kotor (Gross Earnings) =\n"
        "   Gaji Pokok + Tunjangan Jabatan + Tunjangan Makan/Transport + (Jam Lembur x Tarif Lembur Per Jam)\n\n"
        "2. Total Pengurang (Deductions) =\n"
        "   Potongan BPJS Ketenagakerjaan (2%) + BPJS Kesehatan (1%) + (Hari Absen Tanpa Keterangan x Tarif Harian) + Potongan PPH 21\n\n"
        "3. Gaji Bersih Diterima (Take Home Pay) =\n"
        "   Penghasilan Kotor (Gross Earnings) - Total Pengurang (Deductions)"
    )
    add_callout(doc, "Formula Matematis Penggajian (Payroll Engine)", payroll_calc, alert_type="success")

    # -------------------------------------------------------------------------
    # 4.8 Modul Akuntansi & Keuangan (General Ledger)
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.8 Modul Akuntansi & Keuangan (General Ledger)")
    add_body_p(doc,
        "Sistem menerapkan kaidah akuntansi standar berpasangan (Double-Entry Bookkeeping). Setiap aktivitas operasional yang "
        "memiliki dampak finansial secara otomatis diposting ke dalam jurnal umum tanpa membutuhkan intervensi manual:"
    )

    journal_rules_data = [
        ["Aktivitas Bisnis Operasional", "Akun Sisi Debit (+)", "Akun Sisi Kredit (-)"],
        ["Penjualan POS Kasir (Tunai)", "1101 - Kas Kasir di Laci (Asset)", "4101 - Pendapatan Penjualan Kafe (Revenue) & 2102 - Hutang Pajak PB1 (Liability)"],
        ["Penjualan POS Kasir (QRIS / EDC)", "1102 - Kas Bank / Piutang EDC (Asset)", "4101 - Pendapatan Penjualan Kafe (Revenue) & 2102 - Hutang Pajak PB1 (Liability)"],
        ["Pengakuan HPP Bahan Baku Resep", "5101 - Beban Pokok Penjualan / HPP (Expense)", "1201 - Persediaan Bahan Baku Kafe (Asset)"],
        ["Penerimaan Barang Supplier (GRN)", "1201 - Persediaan Bahan Baku Kafe (Asset)", "2101 - Hutang Usaha / Hutang Dagang (Liability)"],
        ["Pelunasan Hutang Supplier", "2101 - Hutang Usaha / Hutang Dagang (Liability)", "1103 - Rekening Bank Utama (Asset)"],
        ["Pengeluaran Gaji Karyawan", "6101 - Beban Gaji & Tunjangan (Expense)", "1103 - Rekening Bank Utama (Asset) & 2103 - Titipan BPJS/PPH (Liability)"]
    ]
    add_custom_table(doc, journal_rules_data[0], journal_rules_data[1:], col_widths=[Inches(2.5), Inches(2.0), Inches(2.0)])

    # -------------------------------------------------------------------------
    # 4.9 Modul Master Data Hub & Laporan Eksekutif
    # -------------------------------------------------------------------------
    add_heading_2(doc, "4.9 Modul Master Data Hub & Pelaporan Eksekutif")
    add_body_p(doc,
        "Master Data Hub menyediakan antarmuka CRUD terpusat bagi Super Administrator untuk mengelola parameter sistem: "
        "daftar cabang gerai kafe, akun pengguna, peran dan izin akses, kategori menu, produk dan varian, data supplier, serta denah meja. "
        "Modul Pelaporan menyajikan analisis performa gerai: penjualan jam sibuk (peak hours), produk terlaris (best sellers), "
        "laporan laba rugi komprehensif, dan rasio perputaran persediaan barang."
    )

    doc.add_page_break()

    # =========================================================================
    # BAB 5: SPESIFIKASI NON-FUNGSIONAL (NFR) & KEAMANAN
    # =========================================================================
    add_heading_1(doc, "BAB 5: SPESIFIKASI KEBUTUHAN NON-FUNGSIONAL & KEAMANAN")
    
    nfr_data = [
        ["Parameter NFR", "Target Kinerja / Standar Kepatuhan", "Metode Implementasi & Pengujian"],
        ["Waktu Respons API", "< 200 ms untuk transaksi POS & antrean KDS", "Optimalisasi query PostgreSQL indeks komposit & goroutine async pool."],
        ["Kapasitas Konkurensi", "Mendukung minimal 500 transaksi POS per menit per cabang", "Go net/http high concurrency server dengan memory footprint rendah."],
        ["Ketersediaan Sistem", "99.95% High Availability", "Fallback Redis Standalone dan local cache mode SQLite saat internet drop."],
        ["Keamanan Kata Sandi", "Enkripsi satu arah salted Bcrypt (Cost 12)", "Mencegah serangan rainbow table dan kebocoran kredensial database."],
        ["Perlindungan Injeksi SQL", "100% Parameterized Query via pgx driver", "Tidak ada string concatenation SQL query pada layer repository."],
        ["Proteksi Brute-Force", "Rate Limit 5 fail / min -> Lockout 30 detik", "OWASP A07 compliant in-memory Sliding Window Rate Limiter."],
        ["Keamanan Sesi Web", "Cookie HttpOnly, Secure, SameSite=Lax", "Melindungi token refresh dari serangan XSS (Cross-Site Scripting)."],
        ["Desain Ergonomis UI", "Zero-Scroll Viewport Lock (WCAG AA)", "CSS `fixed inset-0 overflow-hidden` untuk layar sentuh terminal kasir."]
    ]
    add_custom_table(doc, nfr_data[0], nfr_data[1:], col_widths=[Inches(1.8), Inches(2.2), Inches(2.5)])

    doc.add_page_break()

    # =========================================================================
    # BAB 6: MATRIKS PENGUJIAN USER ACCEPTANCE TESTING (UAT)
    # =========================================================================
    add_heading_1(doc, "BAB 6: MATRIKS PENGUJIAN USER ACCEPTANCE TESTING (UAT)")
    add_body_p(doc,
        "Matriks berikut memetakan skenario uji penerimaan pengguna (UAT) untuk memastikan seluruh fitur telah memenuhi "
        "standar fungsional dan aturan bisnis sebelum dinyatakan lolos uji operasional gerai:"
    )

    uat_headers = ["ID Uji", "Skenario Pengujian", "Langkah Operasional", "Hasil yang Diharapkan", "Status"]
    uat_data = [
        ["UAT-01", "Login PIN Kasir Cepat", "Input 4 digit PIN '1234' pada keypad", "Otentikasi sukses, shift Jane Doe aktif, redirect ke POS", "PASS"],
        ["UAT-02", "Simulasi Sidik Jari", "Klik tombol biometrik (🖲️)", "PIN terisi otomatis '1234' dan login tanpa ketik", "PASS"],
        ["UAT-03", "Proteksi Brute-Force", "Input PIN salah 5 kali berturut-turut", "Terminal terkunci 30 detik dan muncul countdown timer", "PASS"],
        ["UAT-04", "Buka Shift Kasir", "Input modal awal Rp 500.000", "Shift berstatus 'open', kas modal tercatat di ledger", "PASS"],
        ["UAT-05", "Pemesanan POS Dine-In", "Pilih meja T-01, tambahkan Caffe Latte", "Order tersimpan, meja berubah kuning (occupied), tiket masuk KDS", "PASS"],
        ["UAT-06", "Pembayaran Tunai & Kembalian", "Total Rp 35.000, bayar Rp 50.000", "Kembalian Rp 15.000, struk tercetak, meja status 'dirty'", "PASS"],
        ["UAT-07", "Blind Cash Closing Seimbang", "Input uang fisik sesuai catatan sistem", "Selisih kas Rp 0, shift ditutup sukses tanpa eskalasi", "PASS"],
        ["UAT-08", "Blind Cash Closing Selisih", "Input fisik kurang Rp 50.000", "Status 'Requires Audit', wajib verifikasi PIN Supervisor", "PASS"],
        ["UAT-09", "Auto-Deduct Bahan Baku (BOM)", "Selesaikan pembayaran pesanan kopi", "Stok biji kopi dan susu berkurang otomatis secara FIFO", "PASS"],
        ["UAT-10", "Alert Stok Minimum", "Stok susu di bawah batas minimum (5 liter)", "Muncul badge peringatan 'Low Stock' dan opsi buat PR", "PASS"],
        ["UAT-11", "Pengadaan P2P 3-Way Match", "Cocokkan PO #01, GRN #01, dan Invoice", "Data kuantitas dan harga valid, disetujui untuk dibayar", "PASS"],
        ["UAT-12", "Presensi Karyawan Clock-in", "Karyawan klik tombol 'Masuk Kerja'", "Waktu jam masuk dan status hadir tercatat di modul HRIS", "PASS"],
        ["UAT-13", "Payroll Run Otomatis", "Generate payroll periode bulan berjalan", "Gaji pokok, lembur, dan potongan BPJS terhitung presisi", "PASS"],
        ["UAT-14", "Posting Jurnal Penjualan", "Transaksi POS kasir lunas", "Jurnal Kas (D) dan Pendapatan + PB1 (K) terposting otomatis", "PASS"],
        ["UAT-15", "Ketahanan Offline SQLite", "Simulasi pemutusan koneksi internet gerai", "Kasir tetap dapat mencatat pesanan dan cetak struk lokal", "PASS"]
    ]
    add_custom_table(doc, uat_headers, uat_data, 
                     col_widths=[Inches(0.8), Inches(1.5), Inches(1.6), Inches(2.0), Inches(0.6)],
                     align_cols=[WD_ALIGN_PARAGRAPH.CENTER, WD_ALIGN_PARAGRAPH.LEFT, WD_ALIGN_PARAGRAPH.LEFT, WD_ALIGN_PARAGRAPH.LEFT, WD_ALIGN_PARAGRAPH.CENTER])

    # End of document note
    doc.add_paragraph()
    p_end = doc.add_paragraph()
    p_end.alignment = WD_ALIGN_PARAGRAPH.CENTER
    format_paragraph(p_end, space_before=10, space_after=10)
    r_end = p_end.add_run("--- DOKUMEN SPESIFIKASI FUNGSIONAL SISTEM INI TELAH RESMI DIVERIFIKASI DAN DISETUJUI ---")
    r_end.font.name = 'Segoe UI'
    r_end.font.size = Pt(9)
    r_end.font.bold = True
    r_end.font.color.rgb = RGBColor(100, 116, 139)

    # Save document
    output_filename = "FSD_Cafe_ERP_System_Enterprise.docx"
    output_path = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), output_filename)
    doc.save(output_path)
    print(f"✅ Dokumen FSD berhasil disimpan di: {output_path}")

if __name__ == "__main__":
    build_fsd()
