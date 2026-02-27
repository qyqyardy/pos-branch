-- ======================
-- MIDTRANS SUPPORT
-- ======================

ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS payment_status TEXT NOT NULL DEFAULT 'completed',
  ADD COLUMN IF NOT EXISTS payment_token TEXT;

-- For existing orders, we assume they are completed (cash/qris manual)
-- For new Midtrans orders, we will set them to 'pending' initially.
