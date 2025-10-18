-- Core schema: players, sessions, shows, tiles, show_tiles

-- Players table
CREATE TABLE players
(
    id           VARCHAR(10) PRIMARY KEY,
    did          VARCHAR(20) UNIQUE NOT NULL,
    display_name VARCHAR(100)       NOT NULL,
    avatar       TEXT,
    settings     JSONB                    DEFAULT '{}'::jsonb,
    score        INTEGER                  DEFAULT 0,
    permissions  BIGINT                   DEFAULT 0,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE
);

-- Create index on Discord ID for faster lookups
CREATE INDEX idx_players_did ON players (did);

-- Sessions table
CREATE TABLE sessions
(
    id         VARCHAR(32) PRIMARY KEY,
    player_id  VARCHAR(10) REFERENCES players (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on player_id for faster lookups
CREATE INDEX idx_sessions_player_id ON sessions (player_id);
-- Create index on expires_at for cleanup queries
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

-- Shows table
CREATE TABLE shows
(
    id                VARCHAR(10) PRIMARY KEY,
    youtube_id        VARCHAR(20) UNIQUE,
    scheduled_time    TIMESTAMP WITH TIME ZONE,
    actual_start_time TIMESTAMP WITH TIME ZONE,
    thumbnail         TEXT,
    metadata          JSONB                    DEFAULT '{}'::jsonb,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP WITH TIME ZONE
);

-- Create index on scheduled_time for querying upcoming/past shows
CREATE INDEX idx_shows_scheduled_time ON shows (scheduled_time);
CREATE INDEX idx_shows_youtube_id ON shows (youtube_id);

-- Tiles table
CREATE TABLE tiles
(
    id         VARCHAR(10) PRIMARY KEY,
    title      VARCHAR(200) NOT NULL,
    category   VARCHAR(50),
    last_drawn TIMESTAMP WITH TIME ZONE,
    weight     FLOAT8       NOT NULL    DEFAULT 0.30 + floor(random() * 36) * 0.02,
    score      FLOAT8       NOT NULL    DEFAULT floor(random() * 46 + 5),
    created_by VARCHAR(10)  REFERENCES players (id) ON DELETE SET NULL,
    settings   JSONB                    DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for common queries
CREATE INDEX idx_tiles_category ON tiles (category);
CREATE INDEX idx_tiles_created_by ON tiles (created_by);
CREATE INDEX idx_tiles_last_drawn ON tiles (last_drawn);

-- Show Tiles junction table
CREATE TABLE show_tiles
(
    show_id    VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    tile_id    VARCHAR(10) REFERENCES tiles (id) ON DELETE CASCADE,
    weight     float8                   DEFAULT 1,
    score      float8                   DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (show_id, tile_id)
);

-- Create indexes for common queries
CREATE INDEX idx_show_tiles_show_id ON show_tiles (show_id);
CREATE INDEX idx_show_tiles_tile_id ON show_tiles (tile_id);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to auto-update updated_at timestamps
CREATE TRIGGER update_players_updated_at
    BEFORE UPDATE
    ON players
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sessions_updated_at
    BEFORE UPDATE
    ON sessions
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shows_updated_at
    BEFORE UPDATE
    ON shows
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tiles_updated_at
    BEFORE UPDATE
    ON tiles
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_show_tiles_updated_at
    BEFORE UPDATE
    ON show_tiles
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Add comments to document the schema
COMMENT ON TABLE players IS 'Stores user account information';
COMMENT ON TABLE sessions IS 'Stores user session tokens for authentication';
COMMENT ON TABLE shows IS 'Stores information about WAN show episodes';
COMMENT ON TABLE tiles IS 'Stores bingo tile definitions';
COMMENT ON TABLE show_tiles IS 'Junction table linking tiles to specific shows with dynamic scoring';

COMMENT ON COLUMN players.did IS 'Discord user ID from OAuth authentication';
COMMENT ON COLUMN players.settings IS 'JSON object containing user preferences (chat_name_color, pronouns, sound_on_mention, interface_language, etc)';
COMMENT ON COLUMN tiles.settings IS 'JSON object containing tile configuration (needs_context, has_timer, timer_duration, etc)';
COMMENT ON COLUMN shows.metadata IS 'JSON object containing additional show data (hosts, sponsors, duration, etc)';
COMMENT ON COLUMN boards.tiles IS 'Array of tile IDs representing the player''s bingo board layout';