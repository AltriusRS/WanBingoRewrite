-- WAN Show Community Bingo — PostgreSQL base schema
-- This schema is designed to match the MVP described in README and allow for
-- straightforward evolution into the stretch goals (history, leaderboards, etc.).
--
-- Notes
-- - Handlers already query table `tile_router (text, category, weight, last_drawn_show)`.
--   This table definition is backward compatible with the existing code.
-- - Timestamps use timestamptz for clarity across time zones.
-- - Use with PostgreSQL 13+ is recommended.

-- Enable useful extensions (optional)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1) Core content: Tiles available for boards
CREATE TABLE IF NOT EXISTS tiles
(
    id              VARCHAR(10) PRIMARY KEY,
    text            TEXT             NOT NULL UNIQUE,      -- Displayed tile text
    category        TEXT             NOT NULL,             -- e.g., Linus, Luke, Dan, Sponsor
    weight          DOUBLE PRECISION NOT NULL DEFAULT 1.0, -- used for weighted random selection
    last_drawn_show VARCHAR(10),                           -- loosely used marker; app may treat as show index
    is_active       BOOLEAN          NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ      NOT NULL DEFAULT now()
);

-- Keep updated_at fresh
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tiles_set_updated_at ON tiles;
CREATE TRIGGER tiles_set_updated_at
    BEFORE UPDATE
    ON tiles
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX IF NOT EXISTS tiles_category_idx ON tiles (category);
CREATE INDEX IF NOT EXISTS tiles_active_idx ON tiles (is_active);

-- 2) Shows (episodes) pulled from WhenPlane (or similar)
CREATE TABLE IF NOT EXISTS shows
(
    id           VARCHAR(10) PRIMARY KEY,
    yt_id        TEXT UNIQUE, -- YouTube video ID
    fp_id        TEXT UNIQUE, -- FloatPlane video ID
    starts_at    TIMESTAMPTZ, -- scheduled start
    went_live_at TIMESTAMPTZ, -- actual live timestamp if known
    title        TEXT,
    is_live      BOOLEAN     NOT NULL DEFAULT FALSE,
    metadata     JSONB,       -- raw payload from WhenPlane/socket if desired
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS shows_live_idx ON shows (is_live);
CREATE INDEX IF NOT EXISTS shows_starts_idx ON shows (starts_at);
CREATE INDEX IF NOT EXISTS shows_created_idx ON shows (created_at);



CREATE TABLE IF NOT EXISTS players
(
    id           VARCHAR(10) PRIMARY KEY,
    wos_id       TEXT        NOT NULL UNIQUE,
    display_name VARCHAR(30) NOT NULL,
    permissions  int8        NOT NULL DEFAULT 0,
    settings     JSONB       NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ          DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS players_deleted_idx ON players (deleted_at);

CREATE INDEX IF NOT EXISTS players_wos_id_idx ON players (wos_id);


CREATE TABLE IF NOT EXISTS leaderboard
(
    player_id VARCHAR(10) NOT NULL REFERENCES players (id) ON DELETE CASCADE,
    score     int4        NOT NULL,
    PRIMARY KEY (player_id)
);

CREATE INDEX IF NOT EXISTS leaderboard_score_idx ON leaderboard (score desc);

-- 3) Tile confirmations during a show (drives chat/system events)
CREATE TABLE IF NOT EXISTS tile_confirmations
(
    id           VARCHAR(10) PRIMARY KEY,
    show_id      VARCHAR(10) NOT NULL REFERENCES shows (id) ON DELETE CASCADE,
    tile_id      VARCHAR(10) NOT NULL REFERENCES tiles (id) ON DELETE CASCADE,
    context      TEXT,                                                   -- optional context entered by host (e.g., "during sponsor seg")
    confirmed_by VARCHAR(10) REFERENCES players (id) ON DELETE SET NULL, -- free-form identifier for host/operator
    confirmed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Avoid duplicate confirmations of the same tile for the same show
CREATE UNIQUE INDEX IF NOT EXISTS tile_confirmations_unique
    ON tile_confirmations (show_id, tile_id);

CREATE INDEX IF NOT EXISTS tile_conf_show_time_idx ON tile_confirmations (show_id, confirmed_at);

-- 4) Host locks (short-lived locks to prevent duplicate confirmations)
-- These entries are intended to be transient; the app should periodically
-- delete expired locks.
CREATE TABLE IF NOT EXISTS host_locks
(
    tile_id    VARCHAR(10) NOT NULL REFERENCES tiles (id) ON DELETE CASCADE,
    show_id    VARCHAR(10) NOT NULL REFERENCES shows (id) ON DELETE CASCADE,
    locked_by  TEXT        NOT NULL, -- identifier for host (name or session id)
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tile_id, show_id)
);

CREATE INDEX IF NOT EXISTS host_locks_expires_idx ON host_locks (expires_at);

-- 5) Chat messages archive (optional; SSE can remain ephemeral if you prefer)
CREATE TABLE IF NOT EXISTS chat_messages
(
    id         VARCHAR(10) PRIMARY KEY,
    show_id    VARCHAR(10) REFERENCES shows (id) ON DELETE SET NULL,
    type       TEXT        NOT NULL CHECK (type IN ('user', 'system')),
    username   TEXT, -- for 'user' type; optional for 'system'
    message    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS chat_messages_show_time_idx ON chat_messages (show_id, created_at);

-- 6) Simple view to assist leaderboards/future reporting (optional example)
CREATE OR REPLACE VIEW v_show_confirmations AS
SELECT s.id         AS show_id,
       s.title,
       s.starts_at,
       COUNT(tc.id) AS confirmations
FROM shows s
         LEFT JOIN tile_confirmations tc ON tc.show_id = s.id
GROUP BY s.id;

CREATE TABLE IF NOT EXISTS player_boards
(
    id         VARCHAR(10) PRIMARY KEY,
    show_id    VARCHAR(10)   NOT NULL REFERENCES shows (id) ON DELETE CASCADE,
    tiles      VARCHAR(10)[] NOT NULL,
    player_id  VARCHAR(10)   NOT NULL REFERENCES players (id) ON DELETE CASCADE,
    won        BOOLEAN       NOT NULL DEFAULT FALSE,
    locked     BOOLEAN       NOT NULL DEFAULT FALSE,
    refreshes  int2          NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now()
);