-- Add state column to shows table
ALTER TABLE shows ADD COLUMN IF NOT EXISTS state TEXT NOT NULL DEFAULT 'scheduled';