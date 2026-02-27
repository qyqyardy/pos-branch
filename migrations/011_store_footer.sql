-- Migration 011: Add footer_message to store_settings
ALTER TABLE store_settings ADD COLUMN IF NOT EXISTS footer_message TEXT DEFAULT '';
