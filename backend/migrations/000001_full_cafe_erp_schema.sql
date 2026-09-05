-- ==============================================================================
-- CAFE ERP SYSTEM: FULL RELATIONAL SCHEMA (POSTGRESQL)
-- ==============================================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ------------------------------------------------------------------------------
-- 1. CABANG / BRANCHES
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    phone VARCHAR(30),
    email VARCHAR(100),
    manager_id UUID,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- ------------------------------------------------------------------------------
-- 2. HAK AKSES, PERAN & PENGGUNA
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    UNIQUE(module, action)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    role_id UUID REFERENCES roles(id) ON DELETE SET NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(30),
    avatar_url TEXT,
    pin_code VARCHAR(6),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_name VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100),
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ------------------------------------------------------------------------------
-- 3. MEJA, ZONASI & DENAH CAFE
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS table_zones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS cafe_tables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    zone_id UUID REFERENCES table_zones(id) ON DELETE SET NULL,
    table_number VARCHAR(20) NOT NULL,
    capacity INT DEFAULT 4,
    status VARCHAR(20) DEFAULT 'available', -- 'available', 'occupied', 'reserved', 'billing', 'dirty'
    pos_x INT DEFAULT 0,
    pos_y INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE(branch_id, table_number)
);

-- ------------------------------------------------------------------------------
-- 4. KATALOG MENU & RESEP (BILL OF MATERIALS)
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS menu_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    icon VARCHAR(50),
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES menu_categories(id) ON DELETE RESTRICT,
    branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    name VARCHAR(150) NOT NULL,
    sku VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    base_price NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    target_station VARCHAR(20) NOT NULL DEFAULT 'kitchen', -- 'kitchen', 'barista', 'cashier'
    image_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, -- e.g. 'Regular', 'Large', 'Hot', 'Iced'
    sku VARCHAR(50),
    additional_price NUMERIC(12, 2) DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    category VARCHAR(50) NOT NULL, -- 'coffee_beans', 'dairy', 'syrup', 'dry_goods', 'packaging', 'cleaning'
    uom VARCHAR(20) NOT NULL,      -- 'kg', 'gram', 'liter', 'ml', 'pcs', 'pack', 'box'
    min_stock NUMERIC(12, 4) DEFAULT 0,
    max_stock NUMERIC(12, 4) DEFAULT 0,
    average_cost NUMERIC(12, 2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS product_recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    quantity_required NUMERIC(12, 4) NOT NULL,
    uom VARCHAR(20) NOT NULL,
    instructions TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- ------------------------------------------------------------------------------
-- 5. LOGISTIK, GUDANG, PENGADAAN & STOK
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) DEFAULT 'main', -- 'main', 'bar', 'kitchen', 'cold_storage'
    address TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS inventory_stocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    quantity NUMERIC(12, 4) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(warehouse_id, inventory_item_id)
);

CREATE TABLE IF NOT EXISTS stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL, -- 'in_purchase', 'out_pos_sales', 'waste', 'transfer', 'adjustment'
    quantity NUMERIC(12, 4) NOT NULL,
    cost_per_unit NUMERIC(12, 2) DEFAULT 0.00,
    reference_id UUID,
    reference_type VARCHAR(50), -- 'order', 'purchase_order', 'opname', 'transfer'
    remarks TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    contact_person VARCHAR(100),
    phone VARCHAR(30),
    email VARCHAR(100),
    address TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS purchase_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,
    po_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'draft', -- 'draft', 'submitted', 'manager_approved', 'owner_approved', 'sent', 'partially_received', 'received', 'rejected', 'cancelled'
    total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    ordered_date DATE DEFAULT CURRENT_DATE,
    expected_date DATE,
    received_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS purchase_order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    quantity NUMERIC(12, 4) NOT NULL,
    unit_price NUMERIC(12, 2) NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL,
    quantity_received NUMERIC(12, 4) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stock_opnames (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    opname_number VARCHAR(50) UNIQUE NOT NULL,
    opname_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'in_progress', -- 'in_progress', 'completed', 'approved'
    notes TEXT,
    conducted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS stock_opname_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stock_opname_id UUID NOT NULL REFERENCES stock_opnames(id) ON DELETE CASCADE,
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    system_stock NUMERIC(12, 4) NOT NULL,
    physical_stock NUMERIC(12, 4) NOT NULL,
    difference NUMERIC(12, 4) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ------------------------------------------------------------------------------
-- 6. POS TRANSAKSI & PEMBAYARAN
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    table_id UUID REFERENCES cafe_tables(id) ON DELETE SET NULL,
    customer_name VARCHAR(100),
    cashier_id UUID REFERENCES users(id) ON DELETE SET NULL,
    waiter_id UUID REFERENCES users(id) ON DELETE SET NULL,
    order_type VARCHAR(20) NOT NULL DEFAULT 'dine_in', -- 'dine_in', 'takeaway', 'delivery'
    status VARCHAR(30) NOT NULL DEFAULT 'pending',     -- 'pending', 'processing', 'completed', 'void'
    subtotal NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    discount_amount NUMERIC(12, 2) DEFAULT 0.00,
    tax_amount NUMERIC(12, 2) DEFAULT 0.00,
    service_charge NUMERIC(12, 2) DEFAULT 0.00,
    total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price NUMERIC(12, 2) NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL,
    kitchen_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'cooking', 'ready', 'served'
    station VARCHAR(20) NOT NULL DEFAULT 'barista', -- 'kitchen', 'barista'
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    payment_method VARCHAR(50) NOT NULL, -- 'cash', 'qris', 'card', 'transfer', 'ewallet'
    amount_paid NUMERIC(12, 2) NOT NULL,
    change_due NUMERIC(12, 2) DEFAULT 0.00,
    reference_number VARCHAR(100),
    status VARCHAR(20) DEFAULT 'completed',
    processed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    paid_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ------------------------------------------------------------------------------
-- 7. HRIS, STRUKTUR ORGANISASI, SHIFT, PRESENSI & PAYROLL
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    manager_id UUID REFERENCES users(id) ON DELETE SET NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS positions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    level INT DEFAULT 1,
    base_salary NUMERIC(12, 2) DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    position_id UUID REFERENCES positions(id) ON DELETE SET NULL,
    nik VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100),
    phone VARCHAR(30),
    dob DATE,
    gender VARCHAR(10),
    marital_status VARCHAR(20),
    address TEXT,
    national_id VARCHAR(50),
    tax_id VARCHAR(50),
    bank_name VARCHAR(50),
    bank_account VARCHAR(50),
    bank_account_name VARCHAR(100),
    job_title VARCHAR(100),
    basic_salary NUMERIC(12, 2) DEFAULT 0.00,
    employment_type VARCHAR(20) DEFAULT 'full_time',
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'probation', 'resigned', 'terminated'
    photo_url TEXT,
    join_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS work_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL, -- 'Pagi', 'Siang', 'Malam'
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    color VARCHAR(20) DEFAULT '#2563eb',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS attendances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    shift_id UUID REFERENCES work_shifts(id) ON DELETE SET NULL,
    clock_in TIMESTAMPTZ NOT NULL,
    clock_out TIMESTAMPTZ,
    status VARCHAR(20) DEFAULT 'present', -- 'present', 'late', 'absent', 'leave'
    late_minutes INT DEFAULT 0,
    overtime_minutes INT DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS leaves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type VARCHAR(50) NOT NULL, -- 'Cuti Tahunan', 'Cuti Sakit', 'Cuti Alasan Penting'
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INT NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS payroll_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    basic_salary NUMERIC(12, 2) NOT NULL,
    allowances NUMERIC(12, 2) DEFAULT 0.00,
    overtime_pay NUMERIC(12, 2) DEFAULT 0.00,
    gross_salary NUMERIC(12, 2) NOT NULL,
    bpjs_deduction NUMERIC(12, 2) DEFAULT 0.00,
    tax_deduction NUMERIC(12, 2) DEFAULT 0.00,
    other_deductions NUMERIC(12, 2) DEFAULT 0.00,
    total_deductions NUMERIC(12, 2) DEFAULT 0.00,
    net_salary NUMERIC(12, 2) NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- ------------------------------------------------------------------------------
-- 8. KEUANGAN, BAGAN AKUN & BUKU BESAR
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS chart_of_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    account_type VARCHAR(30) NOT NULL, -- 'asset', 'liability', 'equity', 'revenue', 'expense'
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS journal_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    reference_number VARCHAR(100) UNIQUE NOT NULL,
    entry_date DATE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT,
    status VARCHAR(20) DEFAULT 'posted', -- 'draft', 'posted'
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS journal_entry_lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    journal_entry_id UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,
    description TEXT,
    debit NUMERIC(14, 2) DEFAULT 0.00,
    credit NUMERIC(14, 2) DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ------------------------------------------------------------------------------
-- 9. INDEKS UNTUK PERFORMA QUERY & RELASI
-- ------------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_users_branch ON users(branch_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_variants_product ON product_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_recipes_product ON product_recipes(product_id);
CREATE INDEX IF NOT EXISTS idx_recipes_variant ON product_recipes(variant_id);
CREATE INDEX IF NOT EXISTS idx_stocks_warehouse ON inventory_stocks(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_stocks_item ON inventory_stocks(inventory_item_id);
CREATE INDEX IF NOT EXISTS idx_tables_branch_zone ON cafe_tables(branch_id, zone_id);
CREATE INDEX IF NOT EXISTS idx_orders_branch ON orders(branch_id);
CREATE INDEX IF NOT EXISTS idx_orders_table ON orders(table_id);
CREATE INDEX IF NOT EXISTS idx_orders_created ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_order ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_employees_branch ON employees(branch_id);
CREATE INDEX IF NOT EXISTS idx_employees_dept ON employees(department_id);
CREATE INDEX IF NOT EXISTS idx_attendances_emp ON attendances(employee_id);
CREATE INDEX IF NOT EXISTS idx_journal_lines_entry ON journal_entry_lines(journal_entry_id);

