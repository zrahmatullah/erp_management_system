INSERT INTO roles (id, name, description) VALUES
('00000000-0000-0000-0000-000000000001', 'Super Admin', 'Full access'),
('00000000-0000-0000-0000-000000000002', 'Owner', 'Read full reports'),
('00000000-0000-0000-0000-000000000003', 'Manager', 'Branch management'),
('00000000-0000-0000-0000-000000000004', 'Cashier', 'POS access');

-- Password is Admin@123
INSERT INTO users (id, username, email, password_hash, full_name, is_active) VALUES
('11111111-1111-1111-1111-111111111111', 'admin', 'admin@cafe-erp.com', '$2a$12$Z0o2.J.YJ1oN2w3K6T0buexM.iP9Lp0nLp5R6m3p3mK5C4R7D0QjO', 'Super Administrator', TRUE);
