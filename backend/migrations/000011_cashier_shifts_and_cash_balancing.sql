-- ==============================================================================
-- 000011: CASHIER SHIFTS & BLIND CASH BALANCING (F-01)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS cashier_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    cashier_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    supervisor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    shift_name VARCHAR(50) NOT NULL DEFAULT 'Shift Pagi',
    opening_cash_float NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    actual_cash_counted NUMERIC(15, 2) DEFAULT 0.00,
    expected_cash_total NUMERIC(15, 2) DEFAULT 0.00,
    cash_difference NUMERIC(15, 2) DEFAULT 0.00,
    total_cash_sales NUMERIC(15, 2) DEFAULT 0.00,
    total_non_cash_sales NUMERIC(15, 2) DEFAULT 0.00,
    total_cash_drops NUMERIC(15, 2) DEFAULT 0.00,
    total_refunds NUMERIC(15, 2) DEFAULT 0.00,
    status VARCHAR(20) NOT NULL DEFAULT 'open', -- 'open', 'closed', 'audited'
    opened_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cashier_shift_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shift_id UUID NOT NULL REFERENCES cashier_shifts(id) ON DELETE CASCADE,
    movement_type VARCHAR(30) NOT NULL, -- 'cash_drop', 'paid_out', 'cash_in'
    amount NUMERIC(15, 2) NOT NULL,
    reason VARCHAR(255) NOT NULL,
    authorized_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Link payment to active cashier shift if applicable
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'payments' AND column_name = 'shift_id'
    ) THEN
        ALTER TABLE payments ADD COLUMN shift_id UUID REFERENCES cashier_shifts(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_cashier_shifts_branch_status ON cashier_shifts(branch_id, status);
CREATE INDEX IF NOT EXISTS idx_cashier_shifts_cashier_status ON cashier_shifts(cashier_id, status);
CREATE INDEX IF NOT EXISTS idx_cashier_shift_movements_shift ON cashier_shift_movements(shift_id);
CREATE INDEX IF NOT EXISTS idx_payments_shift_id ON payments(shift_id);
