-- ==============================================================================
-- 000005_add_product_stock_and_recipes.sql
-- Menambahkan kolom stock pada products dan balance_after / product_id pada stock_movements
-- ==============================================================================

-- 1. Tambahkan kolom stock dan min_stock pada tabel products jika belum ada
ALTER TABLE products ADD COLUMN IF NOT EXISTS stock NUMERIC(12, 2) NOT NULL DEFAULT 50;
ALTER TABLE products ADD COLUMN IF NOT EXISTS min_stock NUMERIC(12, 2) NOT NULL DEFAULT 5;

-- 2. Update stock_movements agar bisa mencatat produk jadi dan bahan baku
ALTER TABLE stock_movements ADD COLUMN IF NOT EXISTS product_id UUID REFERENCES products(id) ON DELETE SET NULL;
ALTER TABLE stock_movements ADD COLUMN IF NOT EXISTS balance_after NUMERIC(12, 4) DEFAULT 0;
ALTER TABLE stock_movements ALTER COLUMN inventory_item_id DROP NOT NULL;

-- 3. Resep lengkap untuk seluruh produk (jika belum ada)
INSERT INTO product_recipes (product_id, variant_id, inventory_item_id, quantity_required, uom, instructions)
VALUES
-- Cappuccino: Coffee Beans (0.018 kg), Milk (0.15 L), Cup (1 pcs)
('fa111111-0000-0000-0000-000000000002', NULL, 'f1111111-0000-0000-0000-000000000001', 0.0180, 'kg', '18g espresso beans'),
('fa111111-0000-0000-0000-000000000002', NULL, 'f1111111-0000-0000-0000-000000000002', 0.1500, 'liter', '150ml steamed milk'),
('fa111111-0000-0000-0000-000000000002', NULL, 'f1111111-0000-0000-0000-000000000005', 1.0000, 'pcs', '1 cup & lid'),

-- Caramel Macchiato: Coffee Beans (0.018 kg), Milk (0.18 L), Caramel Syrup (0.02 bottle), Cup (1 pcs)
('fa111111-0000-0000-0000-000000000003', NULL, 'f1111111-0000-0000-0000-000000000001', 0.0180, 'kg', '18g espresso beans'),
('fa111111-0000-0000-0000-000000000003', NULL, 'f1111111-0000-0000-0000-000000000002', 0.1800, 'liter', '180ml steamed milk'),
('fa111111-0000-0000-0000-000000000003', NULL, 'f1111111-0000-0000-0000-000000000004', 0.0200, 'bottle', '20ml caramel syrup'),
('fa111111-0000-0000-0000-000000000003', NULL, 'f1111111-0000-0000-0000-000000000005', 1.0000, 'pcs', '1 cup & lid'),

-- Matcha Green Tea Latte: Matcha Powder (0.02 pack), Milk (0.2 L), Sugar Syrup (0.02 kg), Cup (1 pcs)
('fa111111-0000-0000-0000-000000000004', NULL, 'f1111111-0000-0000-0000-000000000007', 0.0200, 'pack', '20g matcha powder'),
('fa111111-0000-0000-0000-000000000004', NULL, 'f1111111-0000-0000-0000-000000000002', 0.2000, 'liter', '200ml milk'),
('fa111111-0000-0000-0000-000000000004', NULL, 'f1111111-0000-0000-0000-000000000003', 0.0200, 'kg', '20g gula aren'),
('fa111111-0000-0000-0000-000000000004', NULL, 'f1111111-0000-0000-0000-000000000005', 1.0000, 'pcs', '1 cup & lid'),

-- Basque Burnt Cheesecake: Dairy/Cream Cheese
('fa111111-0000-0000-0000-000000000006', NULL, 'f1111111-0000-0000-0000-000000000002', 0.1000, 'liter', '100ml cream dairy')
ON CONFLICT DO NOTHING;

-- 4. Initial Stock Movements untuk Saldo Awal (jika belum ada)
INSERT INTO stock_movements (id, inventory_item_id, warehouse_id, type, quantity, balance_after, reference_type, remarks, created_at)
SELECT 
    uuid_generate_v4(),
    s.inventory_item_id,
    s.warehouse_id,
    'in_purchase',
    s.quantity,
    s.quantity,
    'initial_stock',
    'Saldo Awal Stok Bahan Baku Gudang',
    NOW() - INTERVAL '3 days'
FROM inventory_stocks s
WHERE NOT EXISTS (
    SELECT 1 FROM stock_movements WHERE inventory_item_id = s.inventory_item_id
);

