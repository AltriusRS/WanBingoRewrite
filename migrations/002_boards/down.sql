-- Undo boards and confirmations

DROP TRIGGER IF EXISTS update_tile_confirmations_updated_at ON tile_confirmations;
DROP TRIGGER IF EXISTS update_boards_updated_at ON boards;

DROP TABLE IF EXISTS tile_confirmations CASCADE;
DROP TABLE IF EXISTS boards CASCADE;