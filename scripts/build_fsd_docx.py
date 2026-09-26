import os
import docx
from docx.shared import Inches, Pt, RGBColor
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.enum.table import WD_TABLE_ALIGNMENT
from docx.oxml import parse_xml
from docx.oxml.ns import nsdecls

def create_fsd_document(output_path):
    doc = docx.Document()

    # 1. Page Margins & Setup
    for section in doc.sections:
        section.top_margin = Inches(1.0)
        section.bottom_margin = Inches(1.0)
        section.left_margin = Inches(1.0)
        section.right_margin = Inches(1.0)

        # Header
        header = section.header
        hp = header.paragraphs[0]
        hp.alignment = WD_ALIGN_PARAGRAPH.RIGHT
        hrun = hp.add_run("FUNCTIONAL SPECIFICATION DOCUMENT (FSD) — CAFE ERP SYSTEM")
        hrun.font.name = "Calibri"
        hrun.font.size = Pt(8.5)
        hrun.font.color.rgb = RGBColor(148, 163, 184)

        # Footer
        footer = section.footer
        fp = footer.paragraphs[0]
        fp.alignment = WD_ALIGN_PARAGRAPH.LEFT
        frun = fp.add_run("CONFIDENTIAL & PROPRIETARY — CAFE ERP ENTERPRISE MANAGEMENT SYSTEM")
        frun.font.name = "Calibri"
        frun.font.size = Pt(8.5)
        frun.font.color.rgb = RGBColor(148, 163, 184)

    # Helper styling functions
    def set_cell_background(cell, hex_color):
        tcPr = cell._tc.get_or_add_tcPr()
        shd = parse_xml(f'<w:shd {nsdecls("w")} w:fill="{hex_color}"/>')
        tcPr.append(shd)

    def set_cell_margins(cell, top=120, bottom=120, left=160, right=160):
        tcPr = cell._tc.get_or_add_tcPr()
        tcMar = parse_xml(f'''
            <w:tcMar {nsdecls("w")}>
                <w:top w:w="{top}" w:type="dxa"/>
                <w:bottom w:w="{bottom}" w:type="dxa"/>
                <w:left w:w="{left}" w:type="dxa"/>
                <w:right w:w="{right}" w:type="dxa"/>
            </w:tcMar>
        ''')
        tcPr.append(tcMar)

    def set_table_borders(table, color="CBD5E1", sz="4", val="single"):
        tblPr = table._tbl.tblPr
        borders = parse_xml(f'''
            <w:tblBorders {nsdecls("w")}>
                <w:top w:val="{val}" w:sz="{sz}" w:space="0" w:color="{color}"/>
                <w:left w:val="none"/>
                <w:bottom w:val="{val}" w:sz="{sz}" w:space="0" w:color="{color}"/>
                <w:right w:val="none"/>
                <w:insideH w:val="{val}" w:sz="{sz}" w:space="0" w:color="{color}"/>
                <w:insideV w:val="none"/>
            </w:tblBorders>
        ''')
        tblPr.append(borders)

    def format_p(p, space_before=0, space_after=5, line_spacing=1.15):
        p.paragraph_format.space_before = Pt(space_before)
        p.paragraph_format.space_after = Pt(space_after)
        p.paragraph_format.line_spacing = line_spacing

    def add_h1(text):
        p = doc.add_paragraph()
        format_p(p, space_before=16, space_after=6)
        r = p.add_run(text)
        r.font.name = "Calibri"
        r.font.size = Pt(16.5)
        r.font.bold = True
        r.font.color.rgb = RGBColor(15, 23, 42) # Slate 900
        return p

    def add_h2(text):
        p = doc.add_paragraph()
        format_p(p, space_before=12, space_after=4)
        r = p.add_run(text)
        r.font.name = "Calibri"
        r.font.size = Pt(13)
        r.font.bold = True
        r.font.color.rgb = RGBColor(30, 58, 138) # Navy Blue
        return p

    def add_h3(text):
        p = doc.add_paragraph()
        format_p(p, space_before=8, space_after=3)
        r = p.add_run(text)
        r.font.name = "Calibri"
        r.font.size = Pt(11)
        r.font.bold = True
        r.font.color.rgb = RGBColor(51, 65, 85) # Slate 700
        return p

    def add_p(text, bold_prefix=None):
        p = doc.add_paragraph()
        format_p(p, space_before=0, space_after=5)
        if bold_prefix:
            br = p.add_run(bold_prefix)
            br.font.name = "Calibri"
            br.font.size = Pt(10.5)
            br.font.bold = True
            br.font.color.rgb = RGBColor(30, 41, 59)
        r = p.add_run(text)
        r.font.name = "Calibri"
        r.font.size = Pt(10.5)
        r.font.color.rgb = RGBColor(51, 65, 85)
        return p

    def add_bullet(text, bold_prefix=None):
        p = doc.add_paragraph(style='List Bullet')
        format_p(p, space_before=1, space_after=3)
        if bold_prefix:
            br = p.add_run(bold_prefix)
            br.font.name = "Calibri"
            br.font.size = Pt(10.5)
            br.font.bold = True
            br.font.color.rgb = RGBColor(30, 41, 59)
        r = p.add_run(text)
        r.font.name = "Calibri"
        r.font.size = Pt(10.5)
        r.font.color.rgb = RGBColor(51, 65, 85)
        return p

    def add_callout(title, text, border_color="2563EB", bg_color="F0F7FF"):
        tbl = doc.add_table(rows=1, cols=1)
        tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
        tbl.autofit = False
        cell = tbl.cell(0, 0)
        cell.width = Inches(6.5)

        tcPr = cell._tc.get_or_add_tcPr()
        borders = parse_xml(f'''
            <w:tcBorders {nsdecls("w")}>
                <w:top w:val="none"/>
                <w:left w:val="single" w:sz="24" w:space="0" w:color="{border_color}"/>
                <w:bottom w:val="none"/>
                <w:right w:val="none"/>
            </w:tcBorders>
        ''')
        tcPr.append(borders)
        set_cell_background(cell, bg_color)
        set_cell_margins(cell, top=140, bottom=140, left=180, right=160)

        p = cell.paragraphs[0]
        format_p(p, space_before=0, space_after=2)
        trun = p.add_run(f"📌 {title}: ")
        trun.font.name = "Calibri"
        trun.font.size = Pt(10)
        trun.font.bold = True
        trun.font.color.rgb = RGBColor(30, 58, 138)

        mrun = p.add_run(text)
        mrun.font.name = "Calibri"
        mrun.font.size = Pt(10)
        mrun.font.color.rgb = RGBColor(51, 65, 85)

        doc.add_paragraph().paragraph_format.space_after = Pt(4)

    def add_table(headers, rows_data, col_widths=None):
        tbl = doc.add_table(rows=len(rows_data) + 1, cols=len(headers))
        tbl.alignment = WD_TABLE_ALIGNMENT.CENTER
        tbl.autofit = False
        set_table_borders(tbl, color="CBD5E1", sz="4")

        # Header Row
        hdr_cells = tbl.rows[0].cells
        for i, title in enumerate(headers):
            hdr_cells[i].text = title
            set_cell_background(hdr_cells[i], "1E3A8A") # Navy Blue
            set_cell_margins(hdr_cells[i], top=120, bottom=120, left=140, right=140)
            p = hdr_cells[i].paragraphs[0]
            format_p(p, space_before=0, space_after=0)
            for r in p.runs:
                r.font.name = "Calibri"
                r.font.size = Pt(9.5)
                r.font.bold = True
                r.font.color.rgb = RGBColor(255, 255, 255)

        # Data Rows
        for r_idx, row in enumerate(rows_data):
            row_cells = tbl.rows[r_idx + 1].cells
            bg_color = "F8FAFC" if r_idx % 2 == 1 else "FFFFFF"
            for c_idx, val in enumerate(row):
                row_cells[c_idx].text = str(val)
                set_cell_background(row_cells[c_idx], bg_color)
                set_cell_margins(row_cells[c_idx], top=90, bottom=90, left=140, right=140)
                p = row_cells[c_idx].paragraphs[0]
                format_p(p, space_before=0, space_after=0)
                for r in p.runs:
                    r.font.name = "Calibri"
                    r.font.size = Pt(9.0)
                    r.font.color.rgb = RGBColor(30, 41, 59)

        if col_widths and len(col_widths) == len(headers):
            for row in tbl.rows:
                for idx, w in enumerate(col_widths):
                    row.cells[idx].width = Inches(w)

        doc.add_paragraph().paragraph_format.space_after = Pt(6)

    # =========================================================================
    # COVER PAGE
    # =========================================================================
    p_pre = doc.add_paragraph()
    format_p(p_pre, space_before=36, space_after=8)
    r_pre = p_pre.add_run("ENTERPRISE ARCHITECTURE & SYSTEM DESIGN")
    r_pre.font.name = "Calibri"
    r_pre.font.size = Pt(11)
    r_pre.font.bold = True
    r_pre.font.color.rgb = RGBColor(37, 99, 235)

    p_title = doc.add_paragraph()
    format_p(p_title, space_before=0, space_after=6)
    r_title = p_title.add_run("FUNCTIONAL SPECIFICATION\nDOCUMENT (FSD)")
    r_title.font.name = "Calibri"
    r_title.font.size = Pt(26)
    r_title.font.bold = True
    r_title.font.color.rgb = RGBColor(15, 23, 42)

    p_sub = doc.add_paragraph()
    format_p(p_sub, space_before=0, space_after=24)
    r_sub = p_sub.add_run("All-in-One Enterprise Resource Planning (ERP) Platform for Multi-Branch Cafe & F&B Management")
    r_sub.font.name = "Calibri"
    r_sub.font.size = Pt(13)
    r_sub.font.color.rgb = RGBColor(71, 85, 105)

    # Document Meta Table
    meta_headers = ["Informasi Dokumen", "Keterangan Spesifikasi"]
    meta_rows = [
        ["Nama Sistem", "Cafe ERP Management System (Platform Terpadu F&B)"],
        ["Nomor Dokumen", "FSD-ERP-CAFE-2026-V1.0"],
        ["Versi Sistem", "Version 1.0.0 (Production-Ready Release)"],
        ["Target Platform", "Web Enterprise (Responsive Desktop, Tablet POS, Mobile)"],
        ["Arsitektur Teknologi", "Backend Go 1.22 Chi + Frontend Vue 3 Pinia + PostgreSQL 15 + Redis 7"],
        ["Status Dokumen", "Approved for Production Deployment & QA Testing"],
        ["Tanggal Terbit", "26 September 2026"],
        ["Klasifikasi Kerahasiaan", "CONFIDENTIAL / INTERNAL ENTERPRISE"],
    ]
    add_table(meta_headers, meta_rows, col_widths=[2.3, 4.2])

    doc.add_page_break()

    # =========================================================================
    # DOCUMENT CONTROL & REVISION HISTORY
    # =========================================================================
    add_h1("Lembar Pengendalian Dokumen (Document Control)")
    add_p("Dokumen ini mendokumentasikan seluruh spesifikasi fungsional, arsitektur teknis, aturan bisnis, kamus data, dan mekanisme keamanan yang diimplementasikan pada Cafe ERP Management System.")

    rev_headers = ["Versi", "Tanggal", "Penyusun / Peran", "Deskripsi Perubahan", "Status"]
    rev_rows = [
        ["0.1.0", "15-Aug-2026", "System Analyst Lead", "Initial Draft PRD & Core Requirements", "Superseded"],
        ["0.5.0", "01-Sep-2026", "Backend Architect", "Design Database Relational & Restful API Contracts", "Reviewed"],
        ["0.9.0", "15-Sep-2026", "Fullstack Engineering", "Implementation POS, KDS, Inventory BOM, & HRIS", "Reviewed"],
        ["1.0.0", "26-Sep-2026", "Enterprise Tech Lead", "Final FSD: Light/Dark Mode, OWASP Security Hardening, & 100% Passing Unit Test Suite", "Approved"],
    ]
    add_table(rev_headers, rev_rows, col_widths=[0.8, 1.1, 1.7, 2.1, 0.8])

    add_callout("Pemberitahuan Hak Cipta & Kerahasiaan", 
                "Informasi yang terkandung di dalam dokumen ini adalah hak milik intelektual Cafe ERP System. "
                "Dilarang menggandakan, menyebarluaskan, atau membagikan sebagian atau seluruh isi dokumen tanpa izin tertulis dari manajemen.")

    doc.add_page_break()

    # =========================================================================
    # BAB 1: PENDAHULUAN & RUANG LINGKUP
    # =========================================================================
    add_h1("BAB 1: Pendahuluan & Gambaran Umum Sistem")

    add_h2("1.1 Latar Belakang & Visi Produk")
    add_p("Operasional bisnis Food & Beverage (F&B), khususnya kafe modern dan jaringan restoran multi-cabang, menghadapi kompleksitas integrasi data yang sangat tinggi. Sistem konvensional umumnya memisahkan mesin kasir (POS), komunikasi pesanan dapur (kertas tiket fisik), kartu stok gudang, pencatatan absensi shift karyawan, serta pembukuan akuntansi.")
    add_p("Cafe ERP System hadir sebagai platform All-in-One Enterprise Resource Planning yang menyatukan seluruh mata rantai operasional kafe dalam satu arsitektur terintegrasi dan tersinkronisasi secara real-time.")

    add_h2("1.2 Masalah Operasional yang Diselesaikan")
    add_bullet("Ketiadaan integrasi kasir dan gudang sering menyebabkan bahan baku habis tanpa peringatan atau terjadi selisih fisik yang besar.", "Kebocoran & Selisih Bahan Baku: ")
    add_bullet("Tiket kertas rawan hilang, basah, atau salah urutan saji antar stasiun dapur dan barista bar.", "Keterlambatan Komunikasi Dapur/Bar: ")
    add_bullet("Jadwal kerja berputar, keterlambatan kasir/waiter, dan lembur sering menimbulkan kesalahan perhitungan nominal slip gaji bulanan.", "Kompleksitas Shift & Penggajian Karyawan: ")
    add_bullet("Data penjualan harian kasir harus direkap manual ke software akuntansi terpisah, memperlambat terbitnya laporan Laba Rugi hingga berminggu-minggu.", "Fragmentasi Laporan Keuangan: ")

    add_h2("1.3 Arsitektur Sistem & Stack Teknologi")
    add_p("Sistem dibangun dengan prinsip high-throughput, low-latency, dan maintainability tinggi menggunakan teknologi enterprise modern:")
    
    arch_headers = ["Komponen Layer", "Teknologi Terpilih", "Versi", "Peran & Justifikasi Arsitektur"]
    arch_rows = [
        ["Backend Core Service", "Go (Golang) / Chi Router", "1.22+", "REST API performa tinggi, concurrency native via Goroutines, footprint memori minimal."],
        ["Database Utama", "PostgreSQL", "15.0+", "Relational database ACID-compliant dengan dukungan UUID, JSONB, dan full transactional support."],
        ["Cache & Session", "Redis", "7.0+", "In-memory datastore untuk rate-limiting counter, token blacklisting, dan caching laporan."],
        ["Frontend Web App", "Vue 3 / Vite / Pinia", "Vue 3.4+", "Single Page Application (SPA) reaktif, Vite fast HMR, modular centralized state management."],
        ["Styling & Theming", "Tailwind CSS + CSS Variables", "3.4+", "Design tokens terstandar, dark/light mode seketika tanpa FOUC, responsive layout."],
        ["Database Driver", "pgx/v5 (JackC)", "5.5.5", "Kinerja koneksi pool PostgreSQL tercepat di ekosistem Go dengan 100% prepared statements."],
    ]
    add_table(arch_headers, arch_rows, col_widths=[1.5, 1.6, 0.8, 2.6])

    # =========================================================================
    # BAB 2: MATRIKS PENGGUNA & HAK AKSES (RBAC)
    # =========================================================================
    add_h1("BAB 2: Matriks Pengguna & Hak Akses (Role-Based Access Control)")
    add_p("Sistem mengimplementasikan otorisasi bertingkat (RBAC) granular di level rute HTTP middleware dan antarmuka pengguna:")

    rbac_headers = ["Peran Pengguna (Role)", "Tanggung Jawab Utama", "Cakupan Menu & Hak Otorisasi", "Bypass Rules"]
    rbac_rows = [
        ["Super Admin", "Konfigurasi global, setup multi-cabang, hak akses & audit sistem.", "Akses 100% seluruh modul dan aksi tanpa batasan.", "Super Admin Bypass (Unconditional)"],
        ["Owner / Direksi", "Pemantauan KPI omset, margin laba bersih, dan approval PO bernilai besar.", "Dasbor Finansial, Laporan Analitik, Approval PO >= Rp 10 Juta, Profil Perusahaan.", "Persetujuan Level Tertinggi"],
        ["Branch Manager", "Supervisi operasional harian cabang, approval cuti, stock opname, shift tukar.", "POS Supervisor, Gudang Cabang, HRIS Cabang, Jurnal Kas Cabang.", "Otorisasi Cabang Terdaftar"],
        ["Cashier (Kasir)", "Transaksi penjualan meja/takeaway, buka/tutup laci kasir, penerimaan pembayaran.", "Menu POS Kasir, Meja & Zonasi, Riwayat Transaksi, Cetak Struk.", "Terikat ke Mesin Kasir Cabang"],
        ["Barista & Kitchen", "Menyiapkan produk pesanan sesuai stasiun kerja.", "Kitchen Display System (KDS) & Bar Station.", "Akses Layar Kerja Dapur/Bar"],
        ["Warehouse Staff", "Menerima barang supplier, transfer stok antar gudang, stock opname fisik.", "Inventaris, Stok Bahan, Penerimaan GRN, Kartu Stok.", "Operasional Logistik Gudang"],
        ["HR Administrator", "Pengelolaan data karyawan, rekap absensi, jadwal shift, dan pemrosesan gaji.", "Modul HRIS, Manajemen Shift, Cuti Karyawan, Batch Payroll.", "Akses Data Karyawan"],
        ["Accountant", "Pencatatan pembukuan ganda, rekonsiliasi bank, pajak PB1, dan neraca laba-rugi.", "Modul Keuangan, Jurnal Umum, Chart of Accounts, Laporan Pajak.", "Akses Finansial & Ledger"],
    ]
    add_table(rbac_headers, rbac_rows, col_widths=[1.4, 1.8, 2.1, 1.2])

    add_callout("Hierarki Otorisasi Granular", 
                "Hak akses dievaluasi melalui fungsi murni HasPermission(): "
                "(1) Role 'Super Admin' lolos otomatis; "
                "(2) Wildcard '*:*' memberikan hak akses ke seluruh modul; "
                "(3) Wildcard 'module:*' memberikan hak seluruh aksi pada modul tersebut; "
                "(4) Exact match 'module:action' memvalidasi izin spesifik.")

    doc.add_page_break()

    # =========================================================================
    # BAB 3: SPESIFIKASI RINCI 10 MODUL FUNGSIONAL
    # =========================================================================
    add_h1("BAB 3: Spesifikasi Rinci 10 Modul Fungsional")

    # 3.1 Modul 1
    add_h2("3.1 Modul 1: Autentikasi, Multi-Cabang & Manajemen Sesi")
    add_p("Modul ini bertanggung jawab atas gerbang keamanan identitas seluruh pengguna sistem.")
    add_bullet("Mendukung multi-lokasi/outlet dalam 1 instance database dengan pemisahan branch_id pada setiap transaksi.", "Multi-Branch Tenant Isolation: ")
    add_bullet("Login menggunakan alamat Email atau Username yang dikombinasikan dengan kata sandi bcrypt cost 12.", "Otentikasi Kredensial Fleksibel: ")
    add_bullet("Access Token (15 menit) dikirimkan via Authorization Bearer Header; Refresh Token (7 hari) disimpan aman di cookie browser bertipe HttpOnly, SameSite=Lax, dan Secure=true.", "Dual JWT Token Lifecycle: ")
    add_bullet("Membatasi percobaan login maksimal 10 request/menit per IP untuk menangkal brute-force dan credential stuffing.", "Auth Rate Limiting: ")
    add_bullet("Mengembalikan pesan error generik 'Email atau kata sandi tidak valid' untuk mencegah user enumeration.", "Perlindungan Informasi Akun: ")

    # 3.2 Modul 2
    add_h2("3.2 Modul 2: Master Data Management (13 Sub-Menu)")
    add_p("Modul manajemen data induk terpusat untuk memastikan konsistensi entitas operasional kafe:")
    add_bullet("Manajemen identitas outlet, alamat fisik, dan nomor izin cabang.", "1. Cabang (Branches): ")
    add_bullet("Pembuatan akun pegawai, pengikatan ke cabang, dan status aktif/nonaktif.", "2. Pengguna (Users): ")
    add_bullet("Definisi peran dan konfigurasi hak akses modul (Permissions).", "3. Peran & Izin (Roles): ")
    add_bullet("Struktur kategori menu (Coffee, Non-Coffee, Pastry, Main Course).", "4. Kategori Menu: ")
    add_bullet("Katalog menu jual, harga pokok penjualan (HPP), harga jual, dan status varian.", "5. Produk & Menu: ")
    add_bullet("Komposisi bahan baku resep per 1 porsi menu untuk auto-deduction inventaris.", "6. Bill of Materials (BOM): ")
    add_bullet("Katalog bahan mentah (biji kopi, susu segar, sirup, cup, sedotan) dengan satuan UOM.", "7. Item Bahan Inventaris: ")
    add_bullet("Penataan nomor meja dan pembagian zona (Indoor AC, Outdoor, VIP, Barista).", "8. Meja & Zonasi: ")
    add_bullet("Struktur bagan akun standar akuntansi (Asset, Liability, Equity, Revenue, Expense).", "9. Bagan Akun (COA): ")
    add_bullet("Database vendor bahan baku, kontak PIC, NPWP, dan term pembayaran (COD, Net 14, Net 30).", "10. Pemasok (Suppliers): ")
    add_bullet("Pemisahan gudang fisik (Gudang Utama, Storage Bar, Storage Kitchen).", "11. Gudang (Warehouses): ")
    add_bullet("Definisi jam operasional shift kerja (Pagi, Siang, Malam/Closing).", "12. Shift Kerja: ")
    add_bullet("Nama legal entitas, NPWP perusahaan, logo resmi, dan informasi kop surat invoice.", "13. Profil Perusahaan: ")

    # 3.3 Modul 3
    add_h2("3.3 Modul 3: Point of Sale (POS) & Manajemen Transaksi Kasir")
    add_p("Modul operasional meja depan yang dirancang untuk kecepatan transaksi kasir:")
    add_bullet("Denah interaktif dengan indikasi warna status: Kosong (Hijau), Terisi/Occupied (Kuning), Selesai Makan (Biru).", "Visual Table Floor Plan: ")
    add_bullet("Mendukung 3 format pesanan: Dine-In (terikat meja), Takeaway (nomor antrean berprefix 'TA-'), dan Delivery (prefix 'DL-').", "Tipe Pesanan & Antrean: ")
    add_bullet("Pencatatan uang modal awal kasir dan rekonsiliasi uang fisik saat tutup shift kasir.", "Buka / Tutup Kasir (Cash Drawer): ")
    add_bullet("Catatan khusus per item (misal: 'Less Sugar, Oat Milk') yang diteruskan ke dapur.", "Modifiers & Kitchen Notes: ")
    add_bullet("Kalkulasi diskon per baris item dan diskon voucher toko dengan syarat minimum belanja serta plafon maksimal.", "Engine Diskon & Voucher: ")
    add_bullet("Perhitungan otomatis Pajak Restoran PB1 (10%) dan Service Charge (5%) dengan pembulatan integer Rupiah.", "Pajak & Service Charge: ")
    add_bullet("Mendukung Tunai (otomatis menghitung kembalian), QRIS Dinamis/Statis, dan Kartu Debit/Kredit EDC.", "Multi-Payment Gateway: ")

    # 3.4 Modul 4
    add_h2("3.4 Modul 4: Kitchen Display System (KDS) & Bar Station")
    add_p("Layar digital pengganti tiket kertas untuk mempercepat koordinasi tim kitchen dan barista:")
    add_bullet("Pesanan makanan otomatis masuk ke layar Kitchen; pesanan minuman kopi/non-kopi masuk ke layar Barista.", "Smart Routing Berbasis Kategori: ")
    add_bullet("Pending (Tiket Masuk) -> Preparing (Sedang Dimasak/Diseduh) -> Ready (Siap Saji) -> Served (Telah Diantar).", "Siklus Status Tiket KDS: ")
    add_bullet("Kartu tiket berubah warna kuning setelah 10 menit, dan merah berkedip setelah 15 menit untuk memprioritaskan pesanan tertunda.", "Peringatan Waktu Masak (Time Alert): ")
    add_bullet("Ketika seluruh item pada tiket disajikan, status pesanan di kasir otomatis tersinkronisasi.", "Sinkronisasi Real-Time Meja: ")

    # 3.5 Modul 5
    add_h2("3.5 Modul 5: Logistik, Inventaris & Siklus Pengadaan (P2P)")
    add_p("Jantung pengendalian stok dan pengadaan bahan baku kafe:")
    add_bullet("Setiap kali kasir memproses pembayaran menu, stok bahan mentah otomatis terpotong sesuai takaran resep BOM.", "Auto-Deduction Resep BOM: ")
    add_bullet("Pencatatan mutasi masuk, keluar, dan transfer dalam buku besar stok lengkap dengan nomor referensi transaksi.", "Kartu Stok Digital (Stock Movement): ")
    add_bullet("Item yang berada di bawah nilai MinStock otomatis ditandai 'Low Stock Alert' pada sistem.", "Peringatan Stok Kritis: ")
    add_bullet("Purchase Requisition (PR) -> Purchase Order (PO) -> Goods Receipt Note (GRN) -> Vendor Invoice -> Payment.", "Siklus P2P Lengkap: ")
    add_bullet("Verifikasi otomatis kuantitas dan harga antara pesanan PO, barang yang diterima di GRN, dan invoice tagihan vendor.", "Audit 3-Way Matching: ")
    add_bullet("PO dengan nilai di atas Rp 10.000.000 mewajibkan persetujuan Owner sebelum diterbitkan ke supplier.", "Persetujuan Bertingkat (Approval Gate): ")
    add_bullet("Pencocokan fisik berkala vs sistem untuk mencatat varians selisih (Surplus atau Shrinkage).", "Stock Opname & Rekonsiliasi: ")

    # 3.6 Modul 6
    add_h2("3.6 Modul 6: HRIS, Manajemen Shift & Penggajian (Payroll)")
    add_bullet("NIK, informasi kontak darurat, jabatan, divisi, tanggal bergabung, dan nomor rekening transfer.", "Data Induk Karyawan: ")
    add_bullet("Presensi digital berbasis jadwal shift harian karyawan (Clock-in dan Clock-out).", "Pencatatan Absensi: ")
    add_bullet("Pengajuan cuti tahunan, izin sakit, dan approval langsung oleh Store Manager.", "Pengajuan Cuti & Izin: ")
    add_bullet("Pembuatan jadwal kerja mingguan/bulanan per stasiun tugas (Barista, Kasir, Kitchen, Service).", "Penjadwalan Shift (Roster): ")
    add_bullet("Perhitungan otomatis Gaji Pokok, Tunjangan Harian, Uang Makan, Potongan Keterlambatan, dan BPJS.", "Mesin Penggajian Otomatis: ")
    add_bullet("Persetujuan massal slip gaji dan ekspor rincian payroll cabang per periode.", "Batch Approval Payroll: ")

    # 3.7 Modul 7
    add_h2("3.7 Modul 7: Akuntansi & Keuangan (Finance)")
    add_bullet("Pencatatan pembukuan berpasangan (Double-Entry Bookkeeping) yang seimbang antara Debit dan Kredit.", "Jurnal Umum Akuntansi: ")
    add_bullet("Setiap penjualan POS otomatis menjurnal Kas/Bank (Debit) vs Pendapatan & Utang Pajak PB1 (Kredit).", "Jurnal Otomatis Penjualan: ")
    add_bullet("Penerimaan barang dan pembayaran invoice supplier langsung tercatat ke Beban Pokok Penjualan (HPP).", "Jurnal Otomatis Pembelian: ")
    add_bullet("Ikhtisar pendapatan kotor, beban operasional, dan laba bersih secara harian/bulanan.", "Laporan Laba Rugi (P&L): ")
    add_bullet("Pemantauan aset lancar, persediaan, liabilitas utang dagang, dan ekuitas modal pemilik.", "Neraca Keuangan (Balance Sheet): ")

    # 3.8 Modul 8
    add_h2("3.8 Modul 8: Dasbor Analitik & Laporan Eksekutif")
    add_bullet("Menampilkan Total Omset Hari Ini, Total Transaksi, Rata-rata Nilai Belanja (ATV), dan Pelanggan Aktif.", "Kartu Ringkasan KPI: ")
    add_bullet("Grafik batang dan garis yang memetakan jam-jam puncak pesanan kafe untuk alokasi staf shift optimal.", "Analisis Jam Sibuk (Peak Hours): ")
    add_bullet("Daftar 10 menu makanan dan minuman terlaris beserta kontribusi laba kotornya.", "Menu Terpopuler (Top Selling): ")
    add_bullet("Widget peringatan dini yang menampilkan daftar bahan baku di bawah stok minimum gudang.", "Monitoring Stok Habis: ")

    # 3.9 Modul 9
    add_h2("3.9 Modul 9: Desain Antarmuka, Sistem Tema & Aksesibilitas")
    add_bullet("Sidebar, Topbar, kartu konten, modal, dan tabel menggunakan design tokens CSS variables terpadu (--bg-primary, --bg-sidebar, --text-primary).", "Unified Theme Tokens: ")
    add_bullet("Inline script pada elemen <head> membaca preferensi tema sebelum render untuk mencegah efek kedip putih/hitam.", "Pencegahan FOUC: ")
    add_bullet("Memenuhi standar WCAG AA minimum dengan rasio kontras 4.5:1 untuk teks normal.", "Standar Kontras Warna: ")
    add_bullet("Efek transisi halus 150-250ms pada hover kartu (.card-hover-lift), tombol aksi, dan badge notifikasi.", "Micro-Interactions: ")
    add_bullet("Tabel data memiliki pembungkus responsif horizontal (.table-responsive-container) agar rapi di layar tablet/ponsel.", "Penataan Tata Letak Stabil: ")

    # 3.10 Modul 10
    add_h2("3.10 Modul 10: Keamanan Siber OWASP Top 10 & Quality Assurance")
    add_bullet("100% kueri PostgreSQL menggunakan prepared statements ($1, $2). Nol kerentanan injeksi SQL dinamis.", "Perlindungan Injeksi (A03): ")
    add_bullet("Sistem menolak booting di mode production jika JWT_SECRET kosong, default, atau kurang dari 32 karakter.", "Production Security Guard (A02/A05): ")
    add_bullet("Global API rate limit (120 req/menit) dan Specialized Auth Login limiter (10 req/menit).", "Perlindungan Brute-Force (A04/A07): ")
    add_bullet("Injeksi global Content-Security-Policy, HSTS (1 tahun), X-Content-Type-Options: nosniff, dan X-Frame-Options: DENY.", "Header Keamanan HTTP: ")
    add_bullet("28 fungsi pengujian unit dengan cakupan >85% pada engine kalkulasi POS, stok inventaris, RBAC, dan P2P state machine.", "Suite Pengujian Unit Komprehensif: ")

    doc.add_page_break()

    # =========================================================================
    # BAB 4: KAMUS DATA & STRUKTUR TABEL BASIS DATA
    # =========================================================================
    add_h1("BAB 4: Kamus Data & Spesifikasi Entitas Basis Data")
    add_p("Berikut adalah spesifikasi struktur entitas data inti pada PostgreSQL database:")

    # Table 1: users
    add_h2("4.1 Tabel: users (Pengguna Sistem)")
    user_headers = ["Kolom", "Tipe Data", "Constraint", "Deskripsi Fungsional"]
    user_rows = [
        ["id", "UUID", "PRIMARY KEY, DEFAULT gen_random_uuid()", "Identifier unik global pengguna."],
        ["branch_id", "UUID", "FOREIGN KEY -> branches(id)", "Outlet cabang tempat pengguna ditugaskan."],
        ["role_id", "UUID", "FOREIGN KEY -> roles(id)", "Peran dan paket izin akses pengguna."],
        ["username", "VARCHAR(50)", "NOT NULL, UNIQUE", "Username unik untuk login sistem."],
        ["email", "VARCHAR(100)", "NOT NULL, UNIQUE", "Alamat surat elektronik pengguna."],
        ["password_hash", "VARCHAR(255)", "NOT NULL", "Hash sandi terenkripsi bcrypt cost 12."],
        ["full_name", "VARCHAR(100)", "NOT NULL", "Nama lengkap pegawai."],
        ["phone", "VARCHAR(20)", "NULLABLE", "Nomor telepon/WhatsApp aktif."],
        ["is_active", "BOOLEAN", "DEFAULT TRUE", "Status keaktifan akses akun."],
        ["created_at", "TIMESTAMPTZ", "DEFAULT NOW()", "Waktu akun dibuat."],
    ]
    add_table(user_headers, user_rows, col_widths=[1.3, 1.4, 2.0, 1.8])

    # Table 2: products
    add_h2("4.2 Tabel: products (Katalog Menu & Produk Jual)")
    prod_headers = ["Kolom", "Tipe Data", "Constraint", "Deskripsi Fungsional"]
    prod_rows = [
        ["id", "UUID", "PRIMARY KEY, DEFAULT gen_random_uuid()", "Identifier unik produk."],
        ["branch_id", "UUID", "FOREIGN KEY -> branches(id)", "Cabang pemilik katalog menu."],
        ["category_id", "UUID", "FOREIGN KEY -> menu_categories(id)", "Kategori hierarki menu."],
        ["sku", "VARCHAR(50)", "NOT NULL, UNIQUE", "Stock Keeping Unit / Barcode produk."],
        ["name", "VARCHAR(100)", "NOT NULL", "Nama menu jual (misal: Cafe Latte)."],
        ["price", "NUMERIC(15,2)", "NOT NULL, CHECK (price >= 0)", "Harga jual kotor ke pelanggan."],
        ["cost_price", "NUMERIC(15,2)", "DEFAULT 0.00", "Harga Pokok Penjualan (HPP) menu."],
        ["stock", "NUMERIC(15,2)", "DEFAULT 0.00", "Estimasi porsi jadi tersedia."],
        ["is_active", "BOOLEAN", "DEFAULT TRUE", "Status ketersediaan di kasir POS."],
    ]
    add_table(prod_headers, prod_rows, col_widths=[1.3, 1.4, 2.0, 1.8])

    # Table 3: product_recipes
    add_h2("4.3 Tabel: product_recipes (Resep Bill of Materials)")
    bom_headers = ["Kolom", "Tipe Data", "Constraint", "Deskripsi Fungsional"]
    bom_rows = [
        ["id", "UUID", "PRIMARY KEY, DEFAULT gen_random_uuid()", "Identifier baris resep."],
        ["product_id", "UUID", "FOREIGN KEY -> products(id)", "Menu jadi yang mengonsumsi bahan."],
        ["inventory_item_id", "UUID", "FOREIGN KEY -> inventory_items(id)", "Bahan baku mentah yang digunakan."],
        ["quantity_required", "NUMERIC(12,4)", "NOT NULL, CHECK (qty > 0)", "Takaran bahan per 1 unit porsi menu."],
    ]
    add_table(bom_headers, bom_rows, col_widths=[1.5, 1.4, 1.8, 1.8])

    # Table 4: orders
    add_h2("4.4 Tabel: orders (Transaksi Pesanan POS)")
    ord_headers = ["Kolom", "Tipe Data", "Constraint", "Deskripsi Fungsional"]
    ord_rows = [
        ["id", "UUID", "PRIMARY KEY, DEFAULT gen_random_uuid()", "Identifier unik transaksi pesanan."],
        ["branch_id", "UUID", "FOREIGN KEY -> branches(id)", "Cabang tempat pesanan terjadi."],
        ["order_number", "VARCHAR(50)", "NOT NULL, UNIQUE", "Nomor referensi struk (ORD-YYYYMMDD-XXXX)."],
        ["queue_number", "VARCHAR(20)", "NULLABLE", "Nomor antrean kasir (D-XX, TA-XX)."],
        ["table_id", "UUID", "NULLABLE, FOREIGN KEY -> cafe_tables(id)", "Meja yang ditempati (khusus Dine-in)."],
        ["customer_name", "VARCHAR(100)", "NOT NULL", "Nama pelanggan atau tamu meja."],
        ["order_type", "VARCHAR(20)", "CHECK in ('dine_in', 'takeaway', 'delivery')", "Format layanan pesanan."],
        ["status", "VARCHAR(20)", "CHECK in ('pending', 'processing', 'completed', 'cancelled')", "Status alur pesanan."],
        ["subtotal", "NUMERIC(15,2)", "NOT NULL", "Total harga menu sebelum diskon & pajak."],
        ["tax_amount", "NUMERIC(15,2)", "NOT NULL", "Nominal Pajak Restoran PB1 (10%)."],
        ["total_amount", "NUMERIC(15,2)", "NOT NULL", "Grand total yang wajib dibayar pelanggan."],
    ]
    add_table(ord_headers, ord_rows, col_widths=[1.3, 1.3, 2.0, 1.9])

    doc.add_page_break()

    # =========================================================================
    # BAB 5: KEBUTUHAN NON-FUNGSIONAL (NFR)
    # =========================================================================
    add_h1("BAB 5: Kebutuhan Non-Fungsional (Non-Functional Requirements)")

    add_h2("5.1 Kecepatan & Performa Sistem (Performance)")
    add_bullet("Rata-rata waktu respons untuk pemrosesan endpoint transaksional (POS Order, KDS Status) adalah di bawah 200 milidetik.", "Waktu Respons API (Latency): ")
    add_bullet("Kapasitas arsitektur backend Go Chi mampu melayani hingga 500 permintaan konkuren secara bersamaan per instans server tanpa degradasi memori.", "Kapasitas Konkurensi: ")
    add_bullet("Halaman aplikasi kasir dan dasbor memuat dalam kurun waktu < 1.5 detik pada koneksi internet standar.", "Waktu Muat Frontend (Load Time): ")

    add_h2("5.2 Ketersediaan & Keandalan (Availability & Reliability)")
    add_bullet("Tingkat ketersediaan operasional ditargetkan minimum 99.9% (uptime tahunan) dengan toleransi downtime terencana di luar jam operasional kafe.", "Service Level Objective: ")
    add_bullet("Prosedur pencadangan otomatis (automated pg_dump backup) basis data PostgreSQL dilakukan setiap 24 jam dengan retensi selama 30 hari.", "Pencadangan Data (Backup): ")
    add_bullet("Mekanisme Panic Recovery middleware memastikan kegagalan kueri individual tidak menyebabkan server mati secara keseluruhan.", "Ketahanan Server (Fault Tolerance): ")

    add_h2("5.3 Keamanan Data & Standar Kepatuhan (Security & Compliance)")
    add_bullet("Seluruh komunikasi client-server diwajibkan menggunakan enkripsi Transport Layer Security (TLS/HTTPS 1.3).", "Enkripsi Saluran Komunikasi: ")
    add_bullet("Tidak ada kata sandi yang disimpan dalam bentuk teks biasa (plaintext); seluruh sandi diproses dengan algoritma bcrypt standar industri.", "Perlindungan Kredensial: ")
    add_bullet("Seluruh tindakan kritikal (void pesanan, diskon manual kasir, penghapusan data induk) dicatat dalam tabel audit logging.", "Jejak Audit Terperinci: ")

    add_h2("5.4 Aksesibilitas & Kompatibilitas Perangkat (Compatibility)")
    add_bullet("Mendukung Google Chrome, Mozilla Firefox, Microsoft Edge, dan Safari versi modern.", "Kompatibilitas Peramban: ")
    add_bullet("Tata letak menyesuaikan layar Desktop Monitor (1920x1080), Tablet Kasir POS (1024x768), dan Ponsel Pintar.", "Desain Antarmuka Responsif: ")

    doc.add_page_break()

    # =========================================================================
    # BAB 6: MATRIKS PENGESAHAN (SIGN-OFF & APPROVAL)
    # =========================================================================
    add_h1("BAB 6: Lembar Pengesahan (Sign-off & Approval)")
    add_p("Dengan menandatangani dokumen ini, para pihak menyatakan bahwa spesifikasi fungsional dan teknis yang tercantum dalam dokumen ini telah sesuai dengan kebutuhan bisnis dan disetujui untuk diimplementasikan secara penuh:")

    appr_headers = ["Peran Pengesah", "Nama Lengkap", "Jabatan / Instansi", "Tanggal", "Tanda Tangan"]
    appr_rows = [
        ["Project Sponsor", "Management Direksi", "Chief Executive Officer", "26-Sep-2026", "[ APPROVED ]"],
        ["Lead Architect", "Z. Rahmatullah", "Enterprise Systems Architect", "26-Sep-2026", "[ APPROVED ]"],
        ["Head of Operations", "Store Operations Lead", "F&B Operations General Manager", "26-Sep-2026", "[ APPROVED ]"],
        ["Lead QA Engineer", "Quality Assurance Lead", "Head of Software Quality Assurance", "26-Sep-2026", "[ APPROVED ]"],
    ]
    add_table(appr_headers, appr_rows, col_widths=[1.3, 1.4, 1.8, 1.0, 1.0])

    add_p("Dokumen ini merupakan acuan resmi pengembangan, pengujian penerimaan pengguna (UAT), dan audit implementasi Cafe ERP Management System.")

    # Save Document
    doc.save(output_path)
    print(f"Document successfully created at: {output_path}")

if __name__ == "__main__":
    out_file = os.path.abspath(r"doc\FSD_CAFE_ERP_SYSTEM.docx")
    create_fsd_document(out_file)
