-- Migration 014: Optimize Database Indexes

-- 1. Foreign Key Indexes for order_items
-- Speeds up JOINs and relationship lookups between orders and products
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);

-- 2. Expression Index for Date-based lookups
-- Supports current 'created_at::date' query pattern (though we will also optimize the query itself)
CREATE INDEX IF NOT EXISTS idx_orders_created_at_date ON orders ((created_at::date));

-- 3. Partial Index for Active Kitchen Orders
-- Significantly speeds up KDS and POS sidebar polling by only indexing non-finished orders
CREATE INDEX IF NOT EXISTS idx_orders_active_kitchen 
ON orders(kitchen_status) 
WHERE kitchen_status != 'done';

-- 4. Case-insensitive Index for Customer Name Search
-- Supports faster 'LOWER(customer_name)' searches
CREATE INDEX IF NOT EXISTS idx_orders_customer_name_lower ON orders (LOWER(customer_name));
