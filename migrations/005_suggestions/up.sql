-- Tile Suggestions schema

CREATE TABLE tile_suggestions
(
    id           VARCHAR(10) PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    tile_name    VARCHAR(50)  NOT NULL,
    reason       TEXT         NOT NULL,
    status       VARCHAR(20)  NOT NULL DEFAULT 'pending',
    reviewed_by  VARCHAR(10)  REFERENCES players (id) ON DELETE SET NULL,
    reviewed_at  TIMESTAMP WITH TIME ZONE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tile_suggestions_status ON tile_suggestions (status);
CREATE INDEX idx_tile_suggestions_created_at ON tile_suggestions (created_at);

CREATE TRIGGER update_tile_suggestions_updated_at
    BEFORE UPDATE
    ON tile_suggestions
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE tile_suggestions IS 'Stores user-submitted tile suggestions for review';