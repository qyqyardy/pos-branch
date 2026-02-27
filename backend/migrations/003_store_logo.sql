-- =========
-- STORE LOGO
-- =========

ALTER TABLE store_settings
  ADD COLUMN IF NOT EXISTS logo_data_url TEXT NOT NULL DEFAULT '';

