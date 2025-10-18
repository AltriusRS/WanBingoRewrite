# Database Models

This document describes all database models used in the WAN Bingo application, including their fields, relationships, and constraints.

## Table of Contents

- [Player](#player)
- [Session](#session)
- [Show](#show)
- [Tile](#tile)
- [ShowTile](#showtile)
- [Board](#board)
- [TileConfirmation](#tileconfirmation)
- [Message](#message)
- [Timer](#timer)

## Player

Represents a user account linked to Discord OAuth.

```sql
CREATE TABLE players (
    id           VARCHAR(10) PRIMARY KEY,
    did          VARCHAR(20) UNIQUE NOT NULL,
    display_name VARCHAR(100)       NOT NULL,
    avatar       TEXT,
    settings     JSONB              DEFAULT '{}'::jsonb,
    score        INTEGER            DEFAULT 0,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE
);
```

**Fields:**
- `id` - Primary key, 10-character unique identifier
- `did` - Discord user ID (string, unique)
- `display_name` - User's display name (e.g., "LinusTech#1337")
- `avatar` - Discord avatar URL
- `settings` - JSON object for user preferences:
  - `chat_name_color` - Hex color for chat display
  - `sound_on_mention` - Boolean for mention notifications
  - `pronouns` - User's pronouns
  - `interface_language` - UI language preference
- `score` - User's total bingo score
- `created_at` - Account creation timestamp
- `updated_at` - Last modification timestamp
- `deleted_at` - Soft delete timestamp

**Relationships:**
- One-to-many with `sessions`
- One-to-many with `boards`
- One-to-many with `messages`
- One-to-many with `timers` (created_by)

**Indexes:**
- `idx_players_did` on `did`

## Session

Manages user authentication sessions.

```sql
CREATE TABLE sessions (
    id         VARCHAR(32) PRIMARY KEY,
    player_id  VARCHAR(10) REFERENCES players (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

**Fields:**
- `id` - Primary key, 32-character session identifier
- `player_id` - Foreign key to players table
- `created_at` - Session creation timestamp
- `updated_at` - Last activity timestamp
- `expires_at` - Session expiration timestamp (24 hours from creation)
- `deleted_at` - Soft delete timestamp

**Relationships:**
- Many-to-one with `players`

**Indexes:**
- `idx_sessions_player_id` on `player_id`
- `idx_sessions_expires_at` on `expires_at`

## Show

Represents a WAN Show episode.

```sql
CREATE TABLE shows (
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
```

**Fields:**
- `id` - Primary key, 10-character unique identifier
- `youtube_id` - YouTube video ID (optional, unique)
- `scheduled_time` - Planned show start time
- `actual_start_time` - Actual show start time
- `thumbnail` - Show thumbnail image URL
- `metadata` - JSON object containing:
  - `title` - Show title
  - `fp_vod` - Floatplane VOD identifier
  - `hosts` - Array of host names
  - `sponsors` - Array of sponsor information
  - `duration` - Show duration in seconds
- `created_at` - Record creation timestamp
- `updated_at` - Last modification timestamp
- `deleted_at` - Soft delete timestamp

**Relationships:**
- One-to-many with `show_tiles`
- One-to-many with `boards`
- One-to-many with `tile_confirmations`
- One-to-many with `messages`
- One-to-many with `timers`

**Indexes:**
- `idx_shows_scheduled_time` on `scheduled_time`
- `idx_shows_youtube_id` on `youtube_id`

## Tile

Defines bingo tiles that can be used in games.

```sql
CREATE TABLE tiles (
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
```

**Fields:**
- `id` - Primary key, 10-character unique identifier
- `title` - Tile description/title
- `category` - Grouping category (e.g., "Linus", "Events", "Sponsors")
- `last_drawn` - Last time this tile was used in a game
- `weight` - Probability weight for random selection (0.30-0.66)
- `score` - Point value when tile is confirmed (5-50)
- `created_by` - Player who created the tile
- `settings` - JSON object for tile configuration:
  - `needs_context` - Boolean, requires additional context when confirmed
  - `has_timer` - Boolean, tile triggers a timer
  - `timer_duration` - Integer, timer duration in seconds
- `created_at` - Tile creation timestamp
- `updated_at` - Last modification timestamp
- `deleted_at` - Soft delete timestamp

**Relationships:**
- Many-to-one with `players` (created_by)
- One-to-many with `show_tiles`
- One-to-many with `tile_confirmations`

**Indexes:**
- `idx_tiles_category` on `category`
- `idx_tiles_created_by` on `created_by`
- `idx_tiles_last_drawn` on `last_drawn`

## ShowTile

Junction table linking tiles to specific shows with dynamic scoring.

```sql
CREATE TABLE show_tiles (
    show_id    VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    tile_id    VARCHAR(10) REFERENCES tiles (id) ON DELETE CASCADE,
    weight     float8                   DEFAULT 1,
    score      float8                   DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (show_id, tile_id)
);
```

**Fields:**
- `show_id` - Foreign key to shows table
- `tile_id` - Foreign key to tiles table
- `weight` - Show-specific weight multiplier
- `score` - Show-specific score override (0 = use tile default)
- `created_at` - Record creation timestamp
- `updated_at` - Last modification timestamp
- `deleted_at` - Soft delete timestamp

**Relationships:**
- Many-to-one with `shows`
- Many-to-one with `tiles`

**Indexes:**
- `idx_show_tiles_show_id` on `show_id`
- `idx_show_tiles_tile_id` on `tile_id`

## Board

Stores a player's bingo board for a specific show.

```sql
CREATE TABLE boards (
    id                      VARCHAR(10) PRIMARY KEY,
    player_id               VARCHAR(10) REFERENCES players (id) ON DELETE CASCADE,
    show_id                 VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    tiles                   TEXT[] NOT NULL,
    winner                  BOOLEAN                  DEFAULT FALSE,
    total_score             float8                   DEFAULT 0,
    potential_score         float8                   DEFAULT 0,
    regeneration_diminisher float8                   DEFAULT 0,
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at              TIMESTAMP WITH TIME ZONE
);
```

**Fields:**
- `id` - Primary key, 10-character unique identifier
- `player_id` - Foreign key to players table
- `show_id` - Foreign key to shows table
- `tiles` - Array of 25 tile IDs representing the board layout
- `winner` - Boolean indicating if player has bingo
- `total_score` - Current confirmed score
- `potential_score` - Maximum possible score if all tiles confirmed
- `regeneration_diminisher` - Penalty for board regeneration
- `created_at` - Board creation timestamp
- `updated_at` - Last modification timestamp
- `deleted_at` - Soft delete timestamp

**Relationships:**
- Many-to-one with `players`
- Many-to-one with `shows`

**Constraints:**
- `UNIQUE (player_id, show_id)` - One board per player per show

**Indexes:**
- `idx_boards_player_id` on `player_id`
- `idx_boards_show_id` on `show_id`
- `idx_boards_winner` on `winner`

## TileConfirmation

Records when tiles are confirmed during live shows.

```sql
CREATE TABLE tile_confirmations (
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
```

**Fields:**
- `id` - Primary key, 10-character unique identifier
- `show_id` - Foreign key to shows table
- `tile_id` - Foreign key to tiles table
- `confirmed_by` - Player who confirmed the tile
- `context` - Additional context about the confirmation
- `confirmation_time` - When the tile was confirmed
- `created_at` - Record creation timestamp
- `updated_at` - Last modification timestamp
- `deleted_at` - Soft delete timestamp

**Relationships:**
- Many-to-one with `shows`
- Many-to-one with `tiles`
- Many-to-one with `players` (confirmed_by)

**Indexes:**
- `idx_tile_confirmations_show_id` on `show_id`
- `idx_tile_confirmations_tile_id` on `tile_id`
- `idx_tile_confirmations_confirmed_by` on `confirmed_by`
- `idx_tile_confirmations_time` on `confirmation_time`

## Message

Chat messages sent during shows.

```sql
CREATE TABLE messages (
    id         VARCHAR(10) PRIMARY KEY,
    show_id    VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE       NOT NULL,
    player_id  VARCHAR(10) REFERENCES players (id) ON DELETE SET DEFAULT NOT NULL DEFAULT 'DELETED',
    contents   TEXT                                                      NOT NULL,
    system     BOOLEAN                                                   NOT NULL DEFAULT FALSE,
    replying   VARCHAR(10) references messages (id),
    created_at TIMESTAMP WITH TIME ZONE                                           DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE                                           DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

**Fields:**
- `id` - Primary key, 10-character unique identifier
- `show_id` - Foreign key to shows table
- `player_id` - Foreign key to players table (defaults to 'DELETED' for deleted users)
- `contents` - Message text content
- `system` - Boolean indicating if this is a system message
- `replying` - Optional reference to another message (for threading)
- `created_at` - Message creation timestamp
- `updated_at` - Last modification timestamp
- `deleted_at` - Soft delete timestamp

**Relationships:**
- Many-to-one with `shows`
- Many-to-one with `players`
- Self-referencing for replies

**Indexes:**
- `idx_messages_show_id` on `show_id`
- `idx_messages_player_id` on `player_id`

## Timer

Countdown timers for shows.

```sql
CREATE TABLE timers (
    id          VARCHAR(10) PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    duration    INTEGER      NOT NULL,
    created_by  VARCHAR(10)  REFERENCES players (id) ON DELETE SET NULL,
    show_id     VARCHAR(10)  REFERENCES shows (id) ON DELETE CASCADE,
    starts_at   TIMESTAMP WITH TIME ZONE,
    expires_at  TIMESTAMP WITH TIME ZONE,
    is_active   BOOLEAN                  DEFAULT FALSE,
    settings    JSONB                    DEFAULT '{}'::jsonb,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE
);
```

**Fields:**
- `id` - Primary key, 10-character unique identifier
- `title` - Timer display name
- `duration` - Duration in seconds
- `created_by` - Player who created the timer
- `show_id` - Associated show
- `starts_at` - When timer was started
- `expires_at` - When timer will expire
- `is_active` - Whether timer is currently running
- `settings` - JSON configuration:
  - `color` - Hex color for display
  - `sound` - Boolean for audio alerts
  - `recurring` - Boolean for auto-restart
- `created_at` - Timer creation timestamp
- `updated_at` - Last modification timestamp
- `deleted_at` - Soft delete timestamp

**Relationships:**
- Many-to-one with `players` (created_by)
- Many-to-one with `shows`

**Indexes:**
- `idx_timers_show_id` on `show_id`
- `idx_timers_created_by` on `created_by`
- `idx_timers_expires_at` on `expires_at`
- `idx_timers_is_active` on `is_active`

## Special Records

### Deleted User Placeholder

```sql
INSERT INTO players (id, did, display_name, avatar, settings, score, created_at, updated_at)
VALUES ('DELETED', '0', '[Deleted User]', NULL, '{"system_account": true}', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
```

### System User

```sql
INSERT INTO players (id, did, display_name, avatar, settings, score, created_at, updated_at)
VALUES ('SYSTEM', '1', '[System]', NULL, '{"system_account": true}', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
```

## Data Types and Constraints

- **IDs**: 10-character alphanumeric strings (except sessions: 32 chars)
- **Timestamps**: All use `TIMESTAMP WITH TIME ZONE`
- **Soft Deletes**: All tables use `deleted_at` for soft deletion
- **JSON Fields**: `settings` and `metadata` use JSONB for flexible storage
- **Foreign Keys**: Cascading deletes where appropriate
- **Unique Constraints**: Prevent duplicate boards per player per show
- **Default Values**: Sensible defaults for scores, weights, timestamps

## Auto-Updates

All tables have triggers that automatically update the `updated_at` timestamp on row modifications.</content>
</xai:function_call">Database Models Documentation