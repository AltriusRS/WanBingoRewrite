-- Boards and confirmations schema

-- Boards table
CREATE TABLE boards
(
    id                      VARCHAR(10) PRIMARY KEY,
    player_id               VARCHAR(10) REFERENCES players (id) ON DELETE CASCADE,
    show_id                 VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    tiles                   TEXT[] NOT NULL,
    winner                  BOOLEAN                  DEFAULT FALSE,
    total_score             float8                   DEFAULT 0,
    potential_score         float8                   DEFAULT 0,
    regeneration_diminisher float8                   DEFAULT 1.0,
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at              TIMESTAMP WITH TIME ZONE,
    UNIQUE (player_id, show_id)
);

-- Create indexes for common queries
CREATE INDEX idx_boards_player_id ON boards (player_id);
CREATE INDEX idx_boards_show_id ON boards (show_id);
CREATE INDEX idx_boards_winner ON boards (winner);

-- Tile Confirmations table
CREATE TABLE tile_confirmations
(
    id                VARCHAR(10) PRIMARY KEY,
    show_id           VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    tile_id           VARCHAR(10) REFERENCES tiles (id) ON DELETE CASCADE,
    confirmed_by      VARCHAR(10) REFERENCES players (id) ON DELETE SET NULL,
    context           TEXT,
    confirmation_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP WITH TIME ZONE
);

-- Create indexes for common queries
CREATE INDEX idx_tile_confirmations_show_id ON tile_confirmations (show_id);
CREATE INDEX idx_tile_confirmations_tile_id ON tile_confirmations (tile_id);
CREATE INDEX idx_tile_confirmations_confirmed_by ON tile_confirmations (confirmed_by);
CREATE INDEX idx_tile_confirmations_time ON tile_confirmations (confirmation_time);

-- Create triggers
CREATE TRIGGER update_boards_updated_at
    BEFORE UPDATE
    ON boards
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tile_confirmations_updated_at
    BEFORE UPDATE
    ON tile_confirmations
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE boards IS 'Stores player bingo boards for each show';
COMMENT ON TABLE tile_confirmations IS 'Records when tiles are confirmed during a show';
COMMENT ON COLUMN boards.tiles IS 'Array of tile IDs representing the player''s bingo board layout';