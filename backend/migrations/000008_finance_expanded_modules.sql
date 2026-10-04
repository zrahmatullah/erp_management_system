-- 000008_finance_expanded_modules.sql
-- Financial modules: Budgets, Bank Reconciliations, Expenses

CREATE TABLE IF NOT EXISTS public.bank_reconciliations (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    branch_id UUID REFERENCES public.branches(id) ON DELETE SET NULL,
    bank_name VARCHAR(60) NOT NULL,
    period VARCHAR(20) NOT NULL,
    statement_date DATE NOT NULL,
    ref_number VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    type VARCHAR(10) NOT NULL,
    match_status VARCHAR(20) DEFAULT 'unmatched',
    matched_journal_id UUID REFERENCES public.journal_entries(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.budgets (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    branch_id UUID REFERENCES public.branches(id) ON DELETE SET NULL,
    period VARCHAR(20) NOT NULL,
    category VARCHAR(60) NOT NULL,
    account_id UUID REFERENCES public.chart_of_accounts(id) ON DELETE RESTRICT,
    budget_amount NUMERIC(14,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT budgets_period_category_key UNIQUE (period, category)
);

CREATE TABLE IF NOT EXISTS public.expenses (
    id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    branch_id UUID REFERENCES public.branches(id) ON DELETE SET NULL,
    category VARCHAR(60) NOT NULL,
    account_id UUID REFERENCES public.chart_of_accounts(id) ON DELETE RESTRICT,
    expense_number VARCHAR(50) NOT NULL UNIQUE,
    expense_date DATE DEFAULT CURRENT_DATE NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    submitted_by VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    receipt_url TEXT,
    status VARCHAR(30) DEFAULT 'submitted',
    approved_by VARCHAR(100),
    approval_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Seed Budgets
INSERT INTO public.budgets (id, branch_id, period, category, account_id, budget_amount, notes)
VALUES
('3860aeca-6d61-47a4-9e97-3a8192ccdcf3', NULL, '2026-09', 'Bahan Baku & Minuman', 'ab111111-0000-0000-0000-000000000008', 45000000.00, 'Alokasi bahan kopi, susu, sirup, dan pastry'),
('97dc99dc-b6c6-4cbf-b29d-1a08d74b875c', NULL, '2026-09', 'Utilitas Listrik, Air & Gas', 'ab111111-0000-0000-0000-000000000010', 6000000.00, 'PLN pascabayar, PDAM & tabung gas LPG 12kg'),
('619593bc-659b-4ebc-ade4-ef65098fef88', NULL, '2026-09', 'Sewa Tempat & Gedung', 'ab111111-0000-0000-0000-000000000009', 15000000.00, 'Biaya sewa outlet ruko bulanan'),
('ffd613ff-8537-4be7-9438-fa72e8e87413', NULL, '2026-09', 'Promosi & Pemasaran', 'ab111111-0000-0000-0000-000000000010', 4000000.00, 'Digital ads, voucher loyalty & event komunitas'),
('65a1a26a-4664-4d12-99ae-e31d9eed0f2e', NULL, '2026-09', 'Operasional & Kas Kecil', 'ab111111-0000-0000-0000-000000000010', 3000000.00, 'Kebersihan, ATK, packaging take-away')
ON CONFLICT (period, category) DO NOTHING;

-- Seed Expenses
INSERT INTO public.expenses (id, branch_id, category, account_id, expense_number, expense_date, amount, submitted_by, description, receipt_url, status)
VALUES
('ee111111-0000-0000-0000-000000000002', 'b1111111-0000-0000-0000-000000000001', 'Bahan Baku', 'ab111111-0000-0000-0000-000000000008', 'EXP-2026-09-002', '2026-09-12', 750000.00, 'Chef John (Head Chef)', 'Pembelian darurat es batu kristal & daun mint segar di pasar lokal', 'https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400', 'approved'),
('ee111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'Operasional & Kas Kecil', 'ab111111-0000-0000-0000-000000000010', 'EXP-2026-09-001', '2026-09-05', 350000.00, 'Jane Doe (Cashier SPV)', 'Pembelian ATK kasir, kertas thermal struk & spidol marker', 'https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400', 'reimbursed'),
('ee111111-0000-0000-0000-000000000003', 'b1111111-0000-0000-0000-000000000001', 'Pemeliharaan Peralatan', 'ab111111-0000-0000-0000-000000000010', 'EXP-2026-09-003', '2026-09-18', 1250000.00, 'Sarah Jenkins (Store Manager)', 'Servis rutin mesin espresso La Marzocco & penggantian gasket seal', 'https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400', 'approved'),
('ee111111-0000-0000-0000-000000000004', 'b1111111-0000-0000-0000-000000000001', 'Promosi & Pemasaran', 'ab111111-0000-0000-0000-000000000010', 'EXP-2026-09-004', '2026-09-22', 500000.00, 'Budi Santoso (Marketing)', 'Instagram & TikTok Ads promo bundling kopi pastry weekend', 'https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400', 'submitted')
ON CONFLICT (expense_number) DO NOTHING;

-- Seed Bank Reconciliations
INSERT INTO public.bank_reconciliations (id, branch_id, bank_name, period, statement_date, ref_number, description, amount, type, match_status)
VALUES
('222e44dd-5327-4b4e-bb37-566427d8b75a', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-02', 'BCA-TX-00101', 'Setoran Penjualan Kasir POS 01-Sep', 4520000.00, 'credit', 'matched'),
('02603371-3553-42c1-be83-90e6dcf07628', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-02', 'BCA-TX-00102', 'Pembayaran Sewa Outlet Bulan Sep', 15000000.00, 'debit', 'matched'),
('147629bc-29f5-42ec-9092-ac6dc779661f', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-03', 'BCA-TX-00103', 'Auto-debet Tagihan Listrik PLN', 2500000.00, 'debit', 'matched'),
('443dca43-54ea-4ce6-aecc-7c3cebf5f231', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-04', 'BCA-TX-00104', 'Transfer Supplier Biji Kopi PO-0891', 6000000.00, 'debit', 'matched'),
('9b6e825e-4a9b-4def-9fad-f4ca0c10c745', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-10', 'BCA-TX-00105', 'QRIS Settlement BCA Merchant H+1', 3840000.00, 'credit', 'matched'),
('385204d6-0654-4275-9666-11a2f9a073a6', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-15', 'BCA-TX-00106', 'Biaya Administrasi Bulanan Rekening', 25000.00, 'debit', 'unmatched'),
('a82a2a0a-4c48-47ba-b7d3-732e50d3f730', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-15', 'BCA-TX-00107', 'Pendapatan Bunga Jasa Giro', 45600.00, 'credit', 'unmatched'),
('32a6bb23-3bba-47ae-a890-fb75fe73db14', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-20', 'BCA-TX-00108', 'Transfer Payroll Gaji Karyawan', 28260969.00, 'debit', 'matched'),
('48a17046-bc0a-4c55-b92a-aa998bbdfaec', NULL, 'Bank BCA Operasional (1102)', '2026-09', '2026-09-25', 'BCA-TX-00109', 'EDC BCA Settlement Weekend Sales', 5120000.00, 'credit', 'unmatched')
ON CONFLICT (id) DO NOTHING;
