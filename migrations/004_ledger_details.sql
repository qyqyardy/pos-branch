-- ==================
-- LEDGER ENHANCEMENT
-- ==================

ALTER TABLE cash_ledger
  ADD COLUMN IF NOT EXISTS payment_method TEXT NOT NULL DEFAULT 'cash',
  ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'general';

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'cash_ledger_payment_method_check') THEN
    ALTER TABLE cash_ledger
      ADD CONSTRAINT cash_ledger_payment_method_check CHECK (payment_method IN ('cash', 'bank'));
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_cash_ledger_payment_method ON cash_ledger(payment_method);
CREATE INDEX IF NOT EXISTS idx_cash_ledger_category ON cash_ledger(category);

