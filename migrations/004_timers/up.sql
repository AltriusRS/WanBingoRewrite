-- Timers schema

CREATE TABLE IF NOT EXISTS timers
(
    id         VARCHAR(10) PRIMARY KEY,
    title      VARCHAR(200) NOT NULL,
    duration   INTEGER      NOT NULL, -- Duration in seconds
    created_by VARCHAR(10)  REFERENCES players (id) ON DELETE SET NULL,
    show_id    VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    starts_at  TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active  BOOLEAN                  DEFAULT FALSE,
    settings   JSONB                    DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_timers_show_id ON timers (show_id);
CREATE INDEX IF NOT EXISTS idx_timers_created_by ON timers (created_by);
CREATE INDEX IF NOT EXISTS idx_timers_expires_at ON timers (expires_at);
CREATE INDEX IF NOT EXISTS idx_timers_is_active ON timers (is_active);

CREATE OR REPLACE TRIGGER update_timers_updated_at
    BEFORE UPDATE
    ON timers
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();