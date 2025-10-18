# Database Documentation

This document provides comprehensive information about the WAN Bingo database schema, relationships, and operational procedures.

## Overview

The application uses **PostgreSQL 18 (beta)** as its primary database with the following key characteristics:

- **Connection:** pgx v5 driver with connection pooling
- **Performance:** Async I/O operations enabled for enhanced performance
- **Migrations:** Schema defined in `db/schema.sql`
- **Soft Deletes:** All tables use `deleted_at` for data preservation
- **Timestamps:** All use `TIMESTAMP WITH TIME ZONE`
- **IDs:** 10-character alphanumeric strings (except sessions)

## Schema Overview

```mermaid
graph TD
    A[players] --> B[sessions]
    A --> C[boards]
    A --> D[messages]
    A --> E[timers]

    F[shows] --> C
    F --> G[show_tiles]
    F --> H[tile_confirmations]
    F --> D
    F --> E

    I[tiles] --> G
    I --> H
    I --> J[boards.tiles]
```

## Core Tables

### Players

**Purpose:** User accounts linked to Discord OAuth

**Key Fields:**
- `did` - Discord user ID (unique)
- `display_name` - Username with discriminator
- `score` - Total bingo points
- `settings` - User preferences (JSONB)

**Relationships:**
- 1:N with sessions, boards, messages, timers

### Sessions

**Purpose:** Authentication tokens for logged-in users

**Key Fields:**
- `player_id` - Reference to authenticated user
- `expires_at` - Automatic expiration (24 hours)

**Maintenance:** Background cleanup removes expired sessions

### Shows

**Purpose:** WAN Show episodes and metadata

**Key Fields:**
- `youtube_id` - YouTube video identifier
- `scheduled_time` - Planned broadcast time
- `actual_start_time` - Actual start time
- `metadata` - Show information (JSONB)

**Relationships:**
- 1:N with boards, messages, timers, tile_confirmations

### Tiles

**Purpose:** Bingo squares with scoring and categorization

**Key Fields:**
- `title` - Display text
- `category` - Grouping (Linus, Events, Sponsors, etc.)
- `weight` - Random selection probability
- `score` - Points when confirmed
- `settings` - Tile configuration (JSONB)

**Relationships:**
- 1:N with show_tiles, tile_confirmations

## Game Tables

### Boards

**Purpose:** Player bingo cards for specific shows

**Key Fields:**
- `tiles` - Array of 25 tile IDs
- `winner` - Boolean bingo status
- `total_score` - Current points
- `potential_score` - Maximum possible points

**Constraints:**
- Unique per player per show
- Exactly 25 tiles required

### Show Tiles

**Purpose:** Links tiles to shows with show-specific overrides

**Key Fields:**
- `weight` - Show-specific probability modifier
- `score` - Show-specific point override

**Purpose:** Allows different tile pools per episode

### Tile Confirmations

**Purpose:** Records when tiles are marked during live shows

**Key Fields:**
- `confirmation_time` - When tile was confirmed
- `confirmed_by` - Which player confirmed it
- `context` - Additional confirmation details

## Communication Tables

### Messages

**Purpose:** Chat messages during shows

**Key Fields:**
- `contents` - Message text
- `system` - Boolean for system messages
- `replying` - Optional thread reference

**Features:**
- Threaded conversations
- System message support
- Player attribution (with deleted user handling)

### Timers

**Purpose:** Countdown timers for shows

**Key Fields:**
- `title` - Display name
- `duration` - Length in seconds
- `is_active` - Current running status
- `starts_at` / `expires_at` - Timer lifecycle

**Features:**
- Automatic expiration monitoring
- SSE event broadcasting
- Owner-based permissions

## Special Records

### System Accounts

```sql
-- Deleted user placeholder
INSERT INTO players (id, did, display_name, settings)
VALUES ('DELETED', '0', '[Deleted User]', '{"system_account": true}');

-- System user for automated messages
INSERT INTO players (id, did, display_name, settings)
VALUES ('SYSTEM', '1', '[System]', '{"system_account": true}');
```

### Sample Data

The schema includes sample tiles and show data for development:

```sql
-- Sample show
INSERT INTO shows (id, youtube_id, scheduled_time, metadata)
VALUES ('Y2kz75uBC8', 'YVHXYqMPyzc', '2025-10-11T00:30:00Z',
        '{"title": "Piracy Is Dangerous And Harmful", "fp_vod": "w3A5fKcfTi"}');

-- Sample tiles (70+ predefined tiles)
INSERT INTO tiles (id, title, category, weight, score) VALUES
('pYhro7iTSQ', 'Linus or Luke or Dan sighs', 'Events', 0.35, 15),
('3AXPH39mmU', 'Linus ignores Luke to change the topic', 'Linus', 0.42, 8),
-- ... many more tiles
```

## Indexes

### Performance Indexes

```sql
-- Players
CREATE INDEX idx_players_did ON players (did);

-- Sessions
CREATE INDEX idx_sessions_player_id ON sessions (player_id);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

-- Shows
CREATE INDEX idx_shows_scheduled_time ON shows (scheduled_time);
CREATE INDEX idx_shows_youtube_id ON shows (youtube_id);

-- Tiles
CREATE INDEX idx_tiles_category ON tiles (category);
CREATE INDEX idx_tiles_created_by ON tiles (created_by);
CREATE INDEX idx_tiles_last_drawn ON tiles (last_drawn);

-- Boards
CREATE INDEX idx_boards_player_id ON boards (player_id);
CREATE INDEX idx_boards_show_id ON boards (show_id);
CREATE INDEX idx_boards_winner ON boards (winner);

-- Messages
CREATE INDEX idx_messages_show_id ON messages (show_id);
CREATE INDEX idx_messages_player_id ON messages (player_id);

-- Timers
CREATE INDEX idx_timers_show_id ON timers (show_id);
CREATE INDEX idx_timers_created_by ON timers (created_by);
CREATE INDEX idx_timers_expires_at ON timers (expires_at);
CREATE INDEX idx_timers_is_active ON timers (is_active);
```

### Unique Constraints

```sql
-- One board per player per show
ALTER TABLE boards ADD CONSTRAINT unique_player_show_board
    UNIQUE (player_id, show_id);

-- One tile per show (junction table)
ALTER TABLE show_tiles ADD CONSTRAINT pk_show_tiles
    PRIMARY KEY (show_id, tile_id);
```

## Auto-updates

All tables automatically update `updated_at` timestamps:

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Applied to all tables with BEFORE UPDATE triggers
CREATE TRIGGER update_players_updated_at
    BEFORE UPDATE ON players
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## Data Types

### JSONB Fields

**Player Settings:**
```json
{
  "chat_name_color": "#ff0000",
  "sound_on_mention": true,
  "pronouns": "he/him",
  "interface_language": "en",
  "timezone": "America/New_York"
}
```

**Tile Settings:**
```json
{
  "needs_context": true,
  "has_timer": false,
  "timer_duration": 300,
  "color": "#00ff00"
}
```

**Show Metadata:**
```json
{
  "title": "Episode Title",
  "fp_vod": "floatplane_id",
  "hosts": ["Linus", "Luke"],
  "sponsors": ["dbrand", "SquareSpace"],
  "duration": 14400
}
```

**Timer Settings:**
```json
{
  "color": "#ff0000",
  "sound": true,
  "recurring": false,
  "warning_time": 30
}
```

### Arrays

**Board Tiles:**
```sql
-- Array of exactly 25 tile IDs
tiles TEXT[] NOT NULL CHECK (array_length(tiles, 1) = 25)
```

## Database Operations

### Connection Management

```go
// Initialize connection pool
func Init() {
    pool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }
}
```

### Transaction Support

All database functions support optional transactions:

```go
func PersistPlayer(ctx context.Context, player *Player, tx ...pgx.Tx) error {
    if len(tx) > 0 {
        // Use provided transaction
        return tx[0].Exec(ctx, query, args...)
    } else {
        // Use connection pool
        return pool.Exec(ctx, query, args...)
    }
}
```

### Query Patterns

**With Transactions:**
```go
tx, err := pool.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// Multiple operations
err = PersistPlayer(ctx, player, tx)
if err != nil {
    return err
}

err = CreateSession(ctx, player.ID, tx)
if err != nil {
    return err
}

return tx.Commit(ctx)
```

**Without Transactions:**
```go
return PersistPlayer(ctx, player) // Uses pool automatically
```

## Maintenance Procedures

### Cleanup Expired Sessions

```sql
-- Manual cleanup
DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;

-- Automatic via application background task
SELECT COUNT(*) FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
```

### Soft Delete Handling

```sql
-- Find deleted records
SELECT * FROM players WHERE deleted_at IS NOT NULL;

-- Permanent deletion (rare)
DELETE FROM players WHERE deleted_at < CURRENT_TIMESTAMP - INTERVAL '30 days';
```

### Index Maintenance

```sql
-- Check index usage
SELECT * FROM pg_stat_user_indexes WHERE schemaname = 'public';

-- Reindex if needed
REINDEX INDEX idx_players_did;
```

## Backup and Recovery

### Backup Strategy

```bash
# Full database backup
pg_dump -U username -h hostname database_name > backup.sql

# Schema only
pg_dump -U username -h hostname --schema-only database_name > schema.sql

# Data only
pg_dump -U username -h hostname --data-only database_name > data.sql
```

### Restore Procedure

```bash
# Create fresh database
createdb -U username new_database

# Restore schema
psql -U username -d new_database < schema.sql

# Restore data
psql -U username -d new_database < data.sql
```

## Monitoring

### Key Metrics

- **Connection Pool Status**
- **Query Performance**
- **Table Sizes**
- **Index Hit Rates**

### Health Checks

```sql
-- Connection test
SELECT 1;

-- Table counts
SELECT 'players' as table_name, COUNT(*) as count FROM players
UNION ALL
SELECT 'shows', COUNT(*) FROM shows
UNION ALL
SELECT 'tiles', COUNT(*) FROM tiles;

-- Active sessions
SELECT COUNT(*) FROM sessions WHERE expires_at > CURRENT_TIMESTAMP;
```

## Development Setup

### Local Database

```bash
# Create database
createdb wan_bingo_dev

# Run schema
psql -d wan_bingo_dev < db/schema.sql

# Verify setup
psql -d wan_bingo_dev -c "SELECT COUNT(*) FROM tiles;"
```

### Environment Variables

```bash
DATABASE_URL=postgres://user:password@localhost:5432/wan_bingo_dev
DB_MAX_CONNECTIONS=10
DB_MIN_CONNECTIONS=2
```

## Migration Strategy

### Schema Changes

1. **Create Migration Script**
   ```sql
   -- migration_001_add_new_field.sql
   ALTER TABLE players ADD COLUMN new_field TEXT;
   CREATE INDEX idx_players_new_field ON players (new_field);
   ```

2. **Backup Data**
3. **Apply Migration**
4. **Update Application Code**
5. **Test Thoroughly**

### Data Migrations

```sql
-- Example: Migrate existing data
UPDATE players SET new_field = 'default_value' WHERE new_field IS NULL;
```

## Performance Optimization

### Query Optimization

- **Use appropriate indexes**
- **Avoid SELECT *** in production code**
- **Use pagination for large result sets**
- **Batch operations when possible**

### Connection Pool Tuning

```go
config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
config.MaxConns = 10
config.MinConns = 2
config.MaxConnLifetime = time.Hour
```

### Monitoring Queries

```sql
-- Slow queries
SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;

-- Table bloat
SELECT schemaname, tablename, n_dead_tup, n_live_tup
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC;
```</content>
</xai:function_call">Database Documentation