-- ==============================================================================
-- CAFE ERP SYSTEM: DUMMY ACCOUNTS PER ROLE & EMPLOYEE PORTAL DATA
-- Password for all accounts: Admin@123
-- Hash: $2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK
-- ==============================================================================

-- 1. Ensure Roles Exist
INSERT INTO roles (id, name, description) VALUES
('a1111111-0000-0000-0000-000000000001', 'Super Admin', 'Akses penuh seluruh konfigurasi sistem, lisensi dan semua modul'),
('a1111111-0000-0000-0000-000000000002', 'Owner', 'Pemantauan KPI omset, margin laba rugi, dan approval PO bernilai besar'),
('a1111111-0000-0000-0000-000000000003', 'Manager', 'Operasional harian cabang, approval cuti, opname, shift, dan supervisor'),
('a1111111-0000-0000-0000-000000000004', 'Kasir', 'Transaksi penjualan POS, buka/tutup laci kasir, dan cetak struk'),
('a1111111-0000-0000-0000-000000000005', 'Kitchen Staff', 'Memproses pesanan makanan sesuai monitor KDS'),
('a1111111-0000-0000-0000-000000000006', 'HR Admin', 'Manajemen data staf, presensi, cuti, dan penggajian payroll'),
('a1111111-0000-0000-0000-000000000007', 'Warehouse', 'Penerimaan barang PO, transfer stok, dan stock opname'),
('a1111111-0000-0000-0000-000000000008', 'Akuntan', 'Jurnal penyesuaian, petty cash, neraca, laba rugi, dan pajak'),
('a1111111-0000-0000-0000-000000000009', 'Pelayan', 'Pelayanan tamu, catat nomor meja, dan antar pesanan'),
('a1111111-0000-0000-0000-000000000010', 'Barista', 'Peracik kopi & minuman di bar counter')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

-- Grant permissions for new Barista role
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'a1111111-0000-0000-0000-000000000010', id 
FROM permissions 
WHERE module IN ('dashboard', 'pos', 'kitchen', 'inventory')
ON CONFLICT DO NOTHING;

-- Grant permissions for Akuntan / Finance
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'a1111111-0000-0000-0000-000000000008', id 
FROM permissions 
WHERE module IN ('dashboard', 'finance', 'reports')
ON CONFLICT DO NOTHING;

-- Grant permissions for HR Admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'a1111111-0000-0000-0000-000000000006', id 
FROM permissions 
WHERE module IN ('dashboard', 'hris', 'payroll')
ON CONFLICT DO NOTHING;

-- 2. Insert / Update Dummy Users for Every Role
-- Password: Admin@123
INSERT INTO users (id, branch_id, role_id, username, email, password_hash, full_name, phone, pin_code, is_active) VALUES
-- Super Admin
('c1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000001', 'admin', 'admin@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Super Administrator', '081234567890', '123456', TRUE),
-- Owner
('c1111111-0000-0000-0000-000000000020', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000002', 'owner', 'owner@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Bapak Hendra (Owner)', '081234567800', '999888', TRUE),
-- Manager
('c1111111-0000-0000-0000-000000000026', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000003', 'manager', 'manager@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Sarah Johnson (Manager)', '081234567891', '112233', TRUE),
-- Barista
('c1111111-0000-0000-0000-000000000021', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000010', 'barista', 'barista@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Budi Pratama (Barista)', '081234567821', '1111', TRUE),
-- Kasir
('c1111111-0000-0000-0000-000000000024', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000004', 'kasir', 'kasir@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Jane Doe (Kasir)', '081234567892', '1234', TRUE),
-- Inventory / Warehouse
('c1111111-0000-0000-0000-000000000025', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000007', 'inventory', 'inventory@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Ahmad Warehouse', '081234567894', '4455', TRUE),
-- Finance / Akuntan
('c1111111-0000-0000-0000-000000000022', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000008', 'finance', 'finance@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Faisal Finance (Akuntan)', '081234567822', '7788', TRUE),
-- HRIS / HR Admin
('c1111111-0000-0000-0000-000000000023', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000006', 'hris', 'hris@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Hani HRIS (HR Specialist)', '081234567823', '5566', TRUE)
ON CONFLICT (username) DO UPDATE SET 
    email = EXCLUDED.email,
    password_hash = EXCLUDED.password_hash,
    full_name = EXCLUDED.full_name,
    role_id = EXCLUDED.role_id,
    is_active = TRUE;

-- Add user aliases for convenience
INSERT INTO users (id, branch_id, role_id, username, email, password_hash, full_name, phone, pin_code, is_active) VALUES
('c1111111-0000-0000-0000-000000000031', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000004', 'jane.cashier', 'jane.cashier@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Jane Doe', '081234567892', '1234', TRUE),
('c1111111-0000-0000-0000-000000000032', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000007', 'ahmad.warehouse', 'ahmad.warehouse@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Ahmad Warehouse', '081234567894', '4455', TRUE),
('c1111111-0000-0000-0000-000000000033', 'b1111111-0000-0000-0000-000000000001', 'a1111111-0000-0000-0000-000000000008', 'akuntan', 'akuntan@cafe-erp.com', '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK', 'Faisal Finance', '081234567822', '7788', TRUE)
ON CONFLICT (username) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- 3. Ensure Positions exist
INSERT INTO positions (id, department_id, title, level, base_salary) VALUES
('ec111111-0000-0000-0000-000000000006', 'de111111-0000-0000-0000-000000000004', 'Warehouse Officer', 1, 5000000.00),
('ec111111-0000-0000-0000-000000000007', 'de111111-0000-0000-0000-000000000001', 'Finance & Accounting Specialist', 2, 7000000.00),
('ec111111-0000-0000-0000-000000000008', 'de111111-0000-0000-0000-000000000001', 'HR & People Operations', 2, 7500000.00)
ON CONFLICT DO NOTHING;

-- 4. Update / Insert Employees and Link with User Accounts
-- Barista (Budi Pratama)
INSERT INTO employees (id, user_id, branch_id, department_id, position_id, nik, first_name, last_name, email, phone, job_title, basic_salary, status, join_date, remaining_leave) VALUES
('eb111111-0000-0000-0000-000000000003', (SELECT id FROM users WHERE username = 'barista'), 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000002', 'ec111111-0000-0000-0000-000000000003', 'EMP-003', 'Budi', 'Pratama', 'barista@cafe-erp.com', '081234567821', 'Head Barista', 5200000.00, 'active', '2023-03-01', 10)
ON CONFLICT (nik) DO UPDATE SET 
    user_id = EXCLUDED.user_id,
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    email = EXCLUDED.email,
    job_title = EXCLUDED.job_title,
    basic_salary = EXCLUDED.basic_salary,
    remaining_leave = EXCLUDED.remaining_leave;

-- Kasir (Jane Doe)
INSERT INTO employees (id, user_id, branch_id, department_id, position_id, nik, first_name, last_name, email, phone, job_title, basic_salary, status, join_date, remaining_leave) VALUES
('eb111111-0000-0000-0000-000000000004', (SELECT id FROM users WHERE username = 'kasir'), 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000002', 'ec111111-0000-0000-0000-000000000004', 'EMP-004', 'Jane', 'Doe', 'kasir@cafe-erp.com', '081234567892', 'Kasir POS Senior', 4800000.00, 'active', '2023-04-12', 11)
ON CONFLICT (nik) DO UPDATE SET 
    user_id = EXCLUDED.user_id,
    email = EXCLUDED.email,
    job_title = EXCLUDED.job_title,
    remaining_leave = EXCLUDED.remaining_leave;

-- Warehouse / Inventory (Ahmad Warehouse)
INSERT INTO employees (id, user_id, branch_id, department_id, position_id, nik, first_name, last_name, email, phone, job_title, basic_salary, status, join_date, remaining_leave) VALUES
('eb111111-0000-0000-0000-000000000005', (SELECT id FROM users WHERE username = 'inventory'), 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000004', 'ec111111-0000-0000-0000-000000000006', 'EMP-005', 'Ahmad', 'Warehouse', 'inventory@cafe-erp.com', '081234567894', 'Staff Gudang & Logistik', 5000000.00, 'active', '2023-05-01', 12)
ON CONFLICT (nik) DO UPDATE SET 
    user_id = EXCLUDED.user_id,
    email = EXCLUDED.email,
    remaining_leave = EXCLUDED.remaining_leave;

-- Finance / Akuntan (Faisal Finance)
INSERT INTO employees (id, user_id, branch_id, department_id, position_id, nik, first_name, last_name, email, phone, job_title, basic_salary, status, join_date, remaining_leave) VALUES
('eb111111-0000-0000-0000-000000000006', (SELECT id FROM users WHERE username = 'finance'), 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000001', 'ec111111-0000-0000-0000-000000000007', 'EMP-006', 'Faisal', 'Akuntan', 'finance@cafe-erp.com', '081234567822', 'Senior Accountant', 7000000.00, 'active', '2023-02-01', 12)
ON CONFLICT (nik) DO UPDATE SET 
    user_id = EXCLUDED.user_id,
    email = EXCLUDED.email,
    remaining_leave = EXCLUDED.remaining_leave;

-- HRIS (Hani HRIS)
INSERT INTO employees (id, user_id, branch_id, department_id, position_id, nik, first_name, last_name, email, phone, job_title, basic_salary, status, join_date, remaining_leave) VALUES
('eb111111-0000-0000-0000-000000000007', (SELECT id FROM users WHERE username = 'hris'), 'b1111111-0000-0000-0000-000000000001', 'de111111-0000-0000-0000-000000000001', 'ec111111-0000-0000-0000-000000000008', 'EMP-007', 'Hani', 'Pertiwi', 'hris@cafe-erp.com', '081234567823', 'HR & People Operations', 7500000.00, 'active', '2023-01-15', 12)
ON CONFLICT (nik) DO UPDATE SET 
    user_id = EXCLUDED.user_id,
    email = EXCLUDED.email,
    remaining_leave = EXCLUDED.remaining_leave;

-- Manager (Sarah Johnson)
UPDATE employees SET 
    user_id = (SELECT id FROM users WHERE username = 'manager')
WHERE nik = 'EMP-001';

-- 5. Seed Attendance Records for Today
-- Barista (Clocked in early today, on time)
INSERT INTO attendances (id, employee_id, shift_id, clock_in, clock_out, status, late_minutes, overtime_minutes, notes) VALUES
('6a111111-0000-0000-0000-000000000021', 'eb111111-0000-0000-0000-000000000003', 'ac111111-0000-0000-0000-000000000001', CURRENT_DATE + TIME '06:52:00', NULL, 'present', 0, 0, 'Station Bar Siap')
ON CONFLICT DO NOTHING;

-- Kasir (Clocked in today)
INSERT INTO attendances (id, employee_id, shift_id, clock_in, clock_out, status, late_minutes, overtime_minutes, notes) VALUES
('6a111111-0000-0000-0000-000000000022', 'eb111111-0000-0000-0000-000000000004', 'ac111111-0000-0000-0000-000000000001', CURRENT_DATE + TIME '06:58:00', NULL, 'present', 0, 0, 'Laci kasir siap Rp 500rb')
ON CONFLICT DO NOTHING;

-- Warehouse (Clocked in today)
INSERT INTO attendances (id, employee_id, shift_id, clock_in, clock_out, status, late_minutes, overtime_minutes, notes) VALUES
('6a111111-0000-0000-0000-000000000023', 'eb111111-0000-0000-0000-000000000005', 'ac111111-0000-0000-0000-000000000001', CURRENT_DATE + TIME '07:05:00', NULL, 'present', 5, 0, 'Gudang operasional')
ON CONFLICT DO NOTHING;

-- 6. Seed Payroll for EMP-005, EMP-006, EMP-007
INSERT INTO payroll_records (id, employee_id, period_start, period_end, basic_salary, allowances, overtime_pay, gross_salary, bpjs_deduction, tax_deduction, total_deductions, net_salary, is_paid, paid_at) VALUES
('9a111111-0000-0000-0000-000000000005', 'eb111111-0000-0000-0000-000000000005', '2026-08-01', '2026-08-31', 5000000.00, 800000.00, 300000.00, 6100000.00, 200000.00, 150000.00, 350000.00, 5750000.00, TRUE, '2026-08-28 10:00:00+07'),
('9a111111-0000-0000-0000-000000000006', 'eb111111-0000-0000-0000-000000000006', '2026-08-01', '2026-08-31', 7000000.00, 1200000.00, 450000.00, 8650000.00, 280000.00, 320000.00, 600000.00, 8050000.00, TRUE, '2026-08-28 10:00:00+07'),
('9a111111-0000-0000-0000-000000000007', 'eb111111-0000-0000-0000-000000000007', '2026-08-01', '2026-08-31', 7500000.00, 1500000.00, 200000.00, 9200000.00, 300000.00, 380000.00, 680000.00, 8520000.00, TRUE, '2026-08-28 10:00:00+07')
ON CONFLICT DO NOTHING;

-- 7. Seed Weekly Shift Schedule for the Current Week
DO $$
DECLARE
    emp RECORD;
    shift_pagi UUID;
    shift_siang UUID;
    shift_malam UUID;
    curr_monday DATE;
    d INT;
    s_id UUID;
BEGIN
    SELECT id INTO shift_pagi FROM work_shifts WHERE name = 'Pagi' LIMIT 1;
    SELECT id INTO shift_siang FROM work_shifts WHERE name = 'Siang' LIMIT 1;
    SELECT id INTO shift_malam FROM work_shifts WHERE name = 'Malam' LIMIT 1;

    curr_monday := date_trunc('week', CURRENT_DATE)::DATE;

    FOR emp IN SELECT id FROM employees WHERE deleted_at IS NULL LOOP
        FOR d IN 0..6 LOOP
            IF d = 6 THEN
                s_id := NULL; -- OFF on Sunday
            ELSIF d % 2 = 0 THEN
                s_id := COALESCE(shift_pagi, shift_siang);
            ELSE
                s_id := COALESCE(shift_siang, shift_malam);
            END IF;

            INSERT INTO employee_schedules (employee_id, shift_id, date, notes)
            VALUES (emp.id, s_id, curr_monday + d, 'Jadwal shift reguler')
            ON CONFLICT (employee_id, date) DO UPDATE 
            SET shift_id = EXCLUDED.shift_id, updated_at = NOW();
        END LOOP;
    END LOOP;
END $$;
