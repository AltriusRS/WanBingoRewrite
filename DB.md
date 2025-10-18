# Database

---

This document counts as the documentation for the data types inside
the database and should be kept updated.

## Players

| Field          | Data Type      | Description                                                    |
|----------------|----------------|----------------------------------------------------------------|
| `id`           | `VARCHAR(10)`  | Generate unique column to identify each player with an account |
| `did`          | `VARCHAR(20)`  | The Discord profile ID of the user when they authenticated.    |
| `display_name` | `VARCHAR(100)` | The display name of the player                                 |
| `avatar`       | `TEXT`         | URL or identifier for the player's avatar                      |
| `settings`     | `JSONB`        | Player settings (chat name color, pronouns, sound, language)   |
| `score`        | `INTEGER`      | Overall player score                                           |
| `created_at`   | `TIMESTAMP`    | Account creation timestamp                                     |
| `updated_at`   | `TIMESTAMP`    | Last update timestamp (auto-updated via trigger)               |
| `deleted_at`   | `TIMESTAMP`    | Soft delete timestamp                                          |

---

## Shows

| Field               | Data Type     | Description                                     |
|---------------------|---------------|-------------------------------------------------|
| `id`                | `VARCHAR(10)` | Unique identifier for each show                 |
| `youtube_id`        | `VARCHAR(20)` | YouTube video ID                                |
| `scheduled_time`    | `TIMESTAMP`   | When the show is scheduled to start             |
| `actual_start_time` | `TIMESTAMP`   | When the show actually started                  |
| `thumbnail`         | `TEXT`        | URL to show thumbnail                           |
| `metadata`          | `JSONB`       | Additional metadata (hosts, sponsors, duration) |
| `created_at`        | `TIMESTAMP`   | Record creation timestamp                       |
| `updated_at`        | `TIMESTAMP`   | Last update timestamp (auto-updated via trigger)|
| `deleted_at`        | `TIMESTAMP`   | Soft delete timestamp                           |

---

## Tiles

| Field        | Data Type      | Description                                   |
|--------------|----------------|-----------------------------------------------|
| `id`         | `VARCHAR(10)`  | Unique identifier for each tile               |
| `title`      | `VARCHAR(200)` | The text displayed on the tile                |
| `category`   | `VARCHAR(50)`  | Category classification for the tile          |
| `last_drawn` | `TIMESTAMP`    | Last time this tile was drawn in a show       |
| `created_by` | `VARCHAR(10)`  | Player ID who created the tile                |
| `settings`   | `JSONB`        | Tile settings (needs_context, has_timer, etc) |
| `created_at` | `TIMESTAMP`    | Tile creation timestamp                       |
| `updated_at` | `TIMESTAMP`    | Last update timestamp (auto-updated via trigger)|
| `deleted_at` | `TIMESTAMP`    | Soft delete timestamp                         |

---

## Show Tiles

Tile information for each show entry, includes dynamic data such as score weighting.

| Field        | Data Type     | Description                                     |
|--------------|---------------|-------------------------------------------------|
| `show_id`    | `VARCHAR(10)` | Reference to the show                           |
| `tile_id`    | `VARCHAR(10)` | Reference to the tile                           |
| `weight`     | `INTEGER`     | Weight/rarity of tile in this show              |
| `score`      | `INTEGER`     | Points awarded for this tile in this show       |
| `created_at` | `TIMESTAMP`   | Record creation timestamp                       |
| `updated_at` | `TIMESTAMP`   | Last update timestamp (auto-updated via trigger)|
| `deleted_at` | `TIMESTAMP`   | Soft delete timestamp                           |

---

## Boards

| Field                     | Data Type     | Description                                     |
|---------------------------|---------------|-------------------------------------------------|
| `id`                      | `VARCHAR(10)` | Unique identifier for each board                |
| `player_id`               | `VARCHAR(10)` | Reference to the player                         |
| `show_id`                 | `VARCHAR(10)` | Reference to the show                           |
| `tiles`                   | `TEXT[]`      | Array of tile IDs on this board                 |
| `winner`                  | `BOOLEAN`     | Whether this board won                          |
| `total_score`             | `INTEGER`     | Total score for this board                      |
| `regeneration_diminisher` | `INTEGER`     | Diminisher value for board regeneration         |
| `created_at`              | `TIMESTAMP`   | Board creation timestamp                        |
| `updated_at`              | `TIMESTAMP`   | Last update timestamp (auto-updated via trigger)|
| `deleted_at`              | `TIMESTAMP`   | Soft delete timestamp                           |

---

## Tile Confirmations

| Field               | Data Type     | Description                                     |
|---------------------|---------------|-------------------------------------------------|
| `id`                | `VARCHAR(10)` | Unique identifier for confirmation              |
| `show_id`           | `VARCHAR(10)` | Reference to the show                           |
| `tile_id`           | `VARCHAR(10)` | Reference to the tile                           |
| `confirmed_by`      | `VARCHAR(10)` | Player or host ID who confirmed                 |
| `context`           | `TEXT`        | Additional context for the confirmation         |
| `confirmation_time` | `TIMESTAMP`   | When the tile was confirmed                     |
| `created_at`        | `TIMESTAMP`   | Record creation timestamp                       |
| `updated_at`        | `TIMESTAMP`   | Last update timestamp (auto-updated via trigger)|
| `deleted_at`        | `TIMESTAMP`   | Soft delete timestamp                           |
