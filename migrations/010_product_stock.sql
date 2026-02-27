-- Migration 010: Add stock column to products
ALTER TABLE products ADD COLUMN stock INTEGER DEFAULT 0;
