-- ======================
-- SUBSCRIPTION (PLAN/BILLING)
-- ======================

ALTER TABLE store_settings
  ADD COLUMN IF NOT EXISTS plan TEXT NOT NULL DEFAULT 'premium',
  ADD COLUMN IF NOT EXISTS paid_until TIMESTAMP NOT NULL DEFAULT (now() + interval '3650 days');

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'store_settings_plan_check') THEN
    ALTER TABLE store_settings
      ADD CONSTRAINT store_settings_plan_check CHECK (plan IN ('standard', 'premium'));
  END IF;
END $$;

