-- ==========================
-- SETTINGS + ROLES + FINANCE
-- ==========================

-- Store profile (single row, id=1)
CREATE TABLE IF NOT EXISTS store_settings (
  id INT PRIMARY KEY,
  name TEXT NOT NULL,
  tagline TEXT NOT NULL,
  address_line1 TEXT NOT NULL,
  address_line2 TEXT NOT NULL,
  phone TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO store_settings (id, name, tagline, address_line1, address_line2, phone)
VALUES (1, 'WARKOP', 'Point of Sale', 'Jl. Contoh No. 1', 'Kota Kamu', '08xx-xxxx-xxxx')
ON CONFLICT (id) DO NOTHING;

-- Users: created_at + role constraint
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT now();

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'users_role_check') THEN
    ALTER TABLE users
      ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'cashier', 'finance'));
  END IF;
END $$;

-- Orders: store metadata + payment fields for reporting
ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS order_type TEXT NOT NULL DEFAULT 'dine_in',
  ADD COLUMN IF NOT EXISTS table_no TEXT,
  ADD COLUMN IF NOT EXISTS guest_count INT,
  ADD COLUMN IF NOT EXISTS customer_name TEXT,
  ADD COLUMN IF NOT EXISTS payment_method TEXT NOT NULL DEFAULT 'cash',
  ADD COLUMN IF NOT EXISTS received BIGINT,
  ADD COLUMN IF NOT EXISTS change BIGINT;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'orders_order_type_check') THEN
    ALTER TABLE orders
      ADD CONSTRAINT orders_order_type_check CHECK (order_type IN ('dine_in', 'take_away'));
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'orders_payment_method_check') THEN
    ALTER TABLE orders
      ADD CONSTRAINT orders_payment_method_check CHECK (payment_method IN ('cash', 'qris'));
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_orders_cashier_id ON orders(cashier_id);

-- Manual bookkeeping: cash ledger (income/expense)
CREATE TABLE IF NOT EXISTS cash_ledger (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  entry_date DATE NOT NULL,
  type TEXT NOT NULL,
  amount BIGINT NOT NULL,
  description TEXT,
  created_by UUID REFERENCES users(id),
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'cash_ledger_type_check') THEN
    ALTER TABLE cash_ledger
      ADD CONSTRAINT cash_ledger_type_check CHECK (type IN ('income', 'expense'));
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'cash_ledger_amount_check') THEN
    ALTER TABLE cash_ledger
      ADD CONSTRAINT cash_ledger_amount_check CHECK (amount > 0);
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_cash_ledger_entry_date ON cash_ledger(entry_date);

