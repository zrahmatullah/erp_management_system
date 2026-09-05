-- Migration: 000004_add_queue_number.sql
-- Add queue_number column to orders table for Dine-in, Takeaway, and Delivery queue tracking

ALTER TABLE orders ADD COLUMN IF NOT EXISTS queue_number VARCHAR(20);

CREATE INDEX IF NOT EXISTS idx_orders_queue_type_date 
ON orders (order_type, created_at) 
WHERE deleted_at IS NULL;

