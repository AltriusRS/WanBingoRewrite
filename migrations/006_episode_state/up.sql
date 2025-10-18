-- Add state column to shows table
ALTER TABLE shows ADD COLUMN state TEXT NOT NULL DEFAULT 'scheduled';