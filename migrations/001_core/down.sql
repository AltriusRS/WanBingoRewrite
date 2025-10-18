-- Undo core schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_show_tiles_updated_at ON show_tiles;
DROP TRIGGER IF EXISTS update_tiles_updated_at ON tiles;
DROP TRIGGER IF EXISTS update_shows_updated_at ON shows;
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
DROP TRIGGER IF EXISTS update_players_updated_at ON players;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS show_tiles CASCADE;
DROP TABLE IF EXISTS tiles CASCADE;
DROP TABLE IF EXISTS shows CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS players CASCADE;