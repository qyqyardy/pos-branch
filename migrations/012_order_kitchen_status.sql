-- Migration 012: Add kitchen_status to orders
ALTER TABLE orders ADD COLUMN IF NOT EXISTS kitchen_status TEXT NOT NULL DEFAULT 'pending';

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'orders_kitchen_status_check') THEN
    ALTER TABLE orders
      ADD CONSTRAINT orders_kitchen_status_check CHECK (kitchen_status IN ('pending', 'preparing', 'ready', 'done'));
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_orders_kitchen_status ON orders(kitchen_status);
