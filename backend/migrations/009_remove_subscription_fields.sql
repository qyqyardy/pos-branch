-- Drop subscription related columns and constraints
ALTER TABLE store_settings DROP CONSTRAINT IF EXISTS store_settings_plan_check;
ALTER TABLE store_settings DROP COLUMN IF EXISTS plan;
ALTER TABLE store_settings DROP COLUMN IF EXISTS paid_until;
