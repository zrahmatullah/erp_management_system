-- ==============================================================================
-- Migration: 000006_p2p_procurement_cycle.sql
-- Procure-to-Pay (P2P) Full Lifecycle & Accounting Integration
-- ==============================================================================

-- 1. Ensure Standard Chart of Accounts (COA) Exist for P2P Cycle
INSERT INTO chart_of_accounts (id, code, name, account_type, description, is_active)
SELECT 'a1111111-0000-0000-0000-000000000001', '11101', 'Kas Operasional', 'asset', 'Kas Kecil & Kasir Operasional', true
WHERE NOT EXISTS (SELECT 1 FROM chart_of_accounts WHERE code = '11101');

INSERT INTO chart_of_accounts (id, code, name, account_type, description, is_active)
SELECT 'a1111111-0000-0000-0000-000000000002', '11102', 'Bank BCA Operasional', 'asset', 'Rekening Bank Giro Operasional', true
WHERE NOT EXISTS (SELECT 1 FROM chart_of_accounts WHERE code = '11102');

INSERT INTO chart_of_accounts (id, code, name, account_type, description, is_active)
SELECT 'a1111111-0000-0000-0000-000000000003', '11301', 'Persediaan Bahan Baku & Barang Dagang', 'asset', 'Persediaan Fisik Gudang Cafe', true
WHERE NOT EXISTS (SELECT 1 FROM chart_of_accounts WHERE code = '11301');

INSERT INTO chart_of_accounts (id, code, name, account_type, description, is_active)
SELECT 'a1111111-0000-0000-0000-000000000004', '11401', 'PPN Masukan (Pajak Masukan)', 'asset', 'Pajak Pertambahan Nilai Masukan', true
WHERE NOT EXISTS (SELECT 1 FROM chart_of_accounts WHERE code = '11401');

INSERT INTO chart_of_accounts (id, code, name, account_type, description, is_active)
SELECT 'a1111111-0000-0000-0000-000000000005', '21101', 'Hutang Usaha (Accounts Payable)', 'liability', 'Hutang Dagang Supplier Resmi', true
WHERE NOT EXISTS (SELECT 1 FROM chart_of_accounts WHERE code = '21101');

INSERT INTO chart_of_accounts (id, code, name, account_type, description, is_active)
SELECT 'a1111111-0000-0000-0000-000000000006', '21201', 'Hutang Belum Difakturkan (Unbilled AP)', 'liability', 'Akrual Penerimaan Barang Belum Ditagih', true
WHERE NOT EXISTS (SELECT 1 FROM chart_of_accounts WHERE code = '21201');


-- 2. Table: purchase_requisitions (PR)
CREATE TABLE IF NOT EXISTS purchase_requisitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    pr_number VARCHAR(50) UNIQUE NOT NULL,
    requested_by UUID REFERENCES users(id) ON DELETE SET NULL,
    department VARCHAR(100) NOT NULL DEFAULT 'Kitchen & Bar',
    status VARCHAR(30) NOT NULL DEFAULT 'pending_approval', -- 'draft', 'pending_approval', 'approved', 'rejected', 'converted_to_po'
    required_date DATE DEFAULT CURRENT_DATE + INTERVAL '3 days',
    notes TEXT,
    total_estimated_cost NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS purchase_requisition_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    purchase_requisition_id UUID NOT NULL REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    quantity NUMERIC(12, 4) NOT NULL,
    estimated_unit_price NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    estimated_total_price NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


-- 3. Enhance purchase_orders with P2P & Tax Fields
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS pr_id UUID REFERENCES purchase_requisitions(id) ON DELETE SET NULL;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS dpp_amount NUMERIC(12, 2) DEFAULT 0.00;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS tax_type VARCHAR(20) DEFAULT 'exclude'; -- 'include', 'exclude', 'non_pkp'
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS tax_rate NUMERIC(5, 2) DEFAULT 11.00;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS tax_amount NUMERIC(12, 2) DEFAULT 0.00;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS payment_terms VARCHAR(50) DEFAULT 'net_30'; -- 'cod', 'net_7', 'net_14', 'net_30'


-- 4. Table: goods_receipt_notes (GRN)
CREATE TABLE IF NOT EXISTS goods_receipt_notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    grn_number VARCHAR(50) UNIQUE NOT NULL,
    delivery_order_number VARCHAR(100) NOT NULL, -- Nomor Surat Jalan Supplier
    received_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'received_full', -- 'received_full', 'received_partial', 'returned'
    notes TEXT,
    received_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS goods_receipt_note_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    grn_id UUID NOT NULL REFERENCES goods_receipt_notes(id) ON DELETE CASCADE,
    po_item_id UUID REFERENCES purchase_order_items(id) ON DELETE SET NULL,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    quantity_ordered NUMERIC(12, 4) NOT NULL,
    quantity_received NUMERIC(12, 4) NOT NULL,
    quantity_rejected NUMERIC(12, 4) NOT NULL DEFAULT 0.0000,
    condition_status VARCHAR(30) NOT NULL DEFAULT 'good', -- 'good', 'damaged', 'expired'
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


-- 5. Table: vendor_invoices (Tagihan Vendor & Faktur Pajak 3-Way Match)
CREATE TABLE IF NOT EXISTS vendor_invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE RESTRICT,
    grn_id UUID REFERENCES goods_receipt_notes(id) ON DELETE SET NULL,
    invoice_number VARCHAR(100) NOT NULL,
    tax_invoice_number VARCHAR(100), -- Nomor e-Faktur PPN Masukan
    invoice_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL DEFAULT CURRENT_DATE + INTERVAL '30 days',
    subtotal NUMERIC(12, 2) NOT NULL DEFAULT 0.00, -- DPP
    tax_amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00, -- PPN Masukan
    total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    paid_amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(30) NOT NULL DEFAULT 'unpaid', -- 'unpaid', 'partially_paid', 'paid', 'cancelled'
    notes TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);


-- 6. Table: vendor_payments (Pelunasan Hutang Usaha / AP Settlement)
CREATE TABLE IF NOT EXISTS vendor_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    vendor_invoice_id UUID NOT NULL REFERENCES vendor_invoices(id) ON DELETE CASCADE,
    payment_number VARCHAR(50) UNIQUE NOT NULL,
    payment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    payment_method VARCHAR(50) NOT NULL DEFAULT 'bank_transfer', -- 'bank_transfer', 'cash', 'giro'
    payment_account_id UUID REFERENCES chart_of_accounts(id) ON DELETE SET NULL,
    amount_paid NUMERIC(12, 2) NOT NULL,
    reference_number VARCHAR(100), -- No Ref Transfer Bank
    notes TEXT,
    paid_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


-- 7. Seed Sample P2P Cycle Data for Initial Testing
INSERT INTO purchase_requisitions (id, branch_id, pr_number, department, status, required_date, notes, total_estimated_cost)
VALUES 
('de111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'PR-202609-001', 'Barista & Coffee Bar', 'approved', CURRENT_DATE + INTERVAL '2 days', 'Kebutuhan biji kopi arabika & susu fresh milk mingguan', 1750000.00)
ON CONFLICT (pr_number) DO NOTHING;

INSERT INTO purchase_requisition_items (id, purchase_requisition_id, inventory_item_id, quantity, estimated_unit_price, estimated_total_price, notes)
VALUES
('de222222-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000001', 'f1111111-0000-0000-0000-000000000001', 10.0, 120000.00, 1200000.00, 'House Blend Espresso Beans'),
('de222222-0000-0000-0000-000000000002', 'de111111-0000-0000-0000-000000000001', 'f1111111-0000-0000-0000-000000000002', 25.0, 22000.00, 550000.00, 'Fresh Milk Pasteurisasi 1L')
ON CONFLICT DO NOTHING;
