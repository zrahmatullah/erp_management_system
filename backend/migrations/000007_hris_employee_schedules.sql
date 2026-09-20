-- 000007_hris_employee_schedules.sql
-- Create table for employee shift schedules to persist weekly rosters

CREATE TABLE IF NOT EXISTS employee_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    shift_id UUID REFERENCES work_shifts(id) ON DELETE SET NULL,
    date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_emp_schedule_date UNIQUE (employee_id, date)
);

CREATE INDEX IF NOT EXISTS idx_emp_schedules_date ON employee_schedules(date);
CREATE INDEX IF NOT EXISTS idx_emp_schedules_employee ON employee_schedules(employee_id);

-- Add remaining_leave column to employees if not exists
ALTER TABLE employees ADD COLUMN IF NOT EXISTS remaining_leave INT DEFAULT 12;

-- Populate default schedules for current week for existing employees
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

