# WAN Bingo API Documentation

This document provides comprehensive documentation for the WAN Bingo server API endpoints, including request/response formats, authentication requirements, and usage examples.

## Table of Contents

- [Authentication](#authentication)
- [Users](#users)
- [Shows](#shows)
- [Tiles](#tiles)
- [Timers](#timers)
- [Chat](#chat)
- [Error Handling](#error-handling)
- [Pagination](#pagination)

## Authentication

The API uses Discord OAuth2 for authentication. Most endpoints require authentication via session cookies.

### OAuth Flow

1. **GET** `/auth/discord/login` - Redirect to Discord OAuth
2. **GET** `/auth/discord/callback` - Handle OAuth callback, set session cookie
3. **POST** `/auth/discord/logout` - Clear session and redirect

### Session Management

Authenticated requests include a `session_id` cookie. The server validates sessions against the database.

---

## Users

User management and profile endpoints.

### GET /users

Get paginated list of all users.

**Authentication:** Optional (shows partial profiles)

**Query Parameters:**
- `page` (integer, default: 1) - Page number
- `limit` (integer, default: 50, max: 100) - Items per page
- `order_by` (string, default: "display_name") - Sort field: `display_name`, `created_at`, `score`
- `order_dir` (string, default: "asc") - Sort direction: `asc`, `desc`

**Response:**
```json
{
  "success": true,
  "players": [
    {
      "id": "abc123def4",
      "display_name": "LinusTech#1337",
      "avatar": "https://cdn.discordapp.com/avatars/...",
      "score": 1500,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total_count": 150,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  },
  "filters": {
    "order_by": "display_name",
    "order_dir": "asc"
  }
}
```

### GET /users/:identifier

Get user profile by ID or display name.

**Authentication:** Optional

**Path Parameters:**
- `identifier` (string) - User ID or display name

**Response:**
```json
{
  "success": true,
  "player": {
    "id": "abc123def4",
    "display_name": "LinusTech#1337",
    "avatar": "https://cdn.discordapp.com/avatars/...",
    "score": 1500,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### GET /users/me

Get authenticated user's full profile.

**Authentication:** Required

**Response:**
```json
{
  "id": "abc123def4",
  "did": "123456789012345678",
  "display_name": "LinusTech#1337",
  "avatar": "https://cdn.discordapp.com/avatars/...",
  "settings": {
    "chat_name_color": "#ff0000",
    "sound_on_mention": true
  },
  "score": 1500,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-20T14:22:00Z"
}
```

---

## Shows

WAN Show episode management.

### GET /shows/latest

Get the most recent show by scheduled time.

**Authentication:** Optional

**Response:**
```json
{
  "id": "Y2kz75uBC8",
  "youtube_id": "YVHXYqMPyzc",
  "scheduled_time": "2025-10-11T00:30:00Z",
  "actual_start_time": "2025-10-11T00:05:06Z",
  "thumbnail": "https://pbs.floatplane.com/stream_thumbnails/...",
  "metadata": {
    "title": "Piracy Is Dangerous And Harmful",
    "fp_vod": "w3A5fKcfTi",
    "duration": 14400
  },
  "created_at": "2025-10-10T23:46:40Z",
  "updated_at": "2025-10-16T19:12:58Z"
}
```

### GET /shows/:id

Get show by ID.

**Authentication:** Optional

**Path Parameters:**
- `id` (string) - Show ID

**Response:** Same as `/shows/latest`

---

## Tiles

Bingo tile management and board generation.

### GET /tiles

Get paginated list of all tiles.

**Authentication:** Optional

**Query Parameters:**
- `page` (integer, default: 1) - Page number
- `limit` (integer, default: 50, max: 100) - Items per page
- `category` (string, optional) - Filter by category
- `order_by` (string, default: "created_at") - Sort field: `created_at`, `title`, `category`, `weight`, `score`
- `order_dir` (string, default: "desc") - Sort direction: `asc`, `desc`

**Response:**
```json
{
  "tiles": [
    {
      "id": "pYhro7iTSQ",
      "title": "Linus or Luke or Dan sighs",
      "category": "Events",
      "last_drawn": "2024-01-15T20:30:00Z",
      "weight": 0.35,
      "score": 15,
      "created_by": "abc123def4",
      "settings": {},
      "created_at": "2024-01-10T10:00:00Z",
      "updated_at": "2024-01-15T20:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total_count": 150,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  },
  "filters": {
    "category": "Linus",
    "order_by": "created_at",
    "order_dir": "desc"
  }
}
```

### GET /tiles/:tile_id

Get individual tile by ID.

**Authentication:** Optional

**Path Parameters:**
- `tile_id` (string) - Tile ID

**Response:** Individual tile object (same format as in tiles array above)

### GET /tiles/show

Get tile IDs for the current show.

**Authentication:** Optional

**Response:**
```json
{
  "show_id": "Y2kz75uBC8",
  "tile_ids": [
    "pYhro7iTSQ",
    "3AXPH39mmU",
    "kJdmP6gDSi"
  ]
}
```

### GET /tiles/me

Get authenticated user's bingo board for current show.

**Authentication:** Required

**Response:**
```json
{
  "board_id": "brd_abc123",
  "show_id": "Y2kz75uBC8",
  "player_id": "usr_abc123",
  "tiles": [
    {
      "id": "pYhro7iTSQ",
      "title": "Linus or Luke or Dan sighs",
      "category": "Events",
      "weight": 0.35,
      "score": 15
    }
  ],
  "winner": false,
  "total_score": 0,
  "potential_score": 0,
  "created_at": "2024-01-15T20:30:00Z"
}
```

### GET /tiles/anonymous

Generate anonymous bingo board (for guest users).

**Authentication:** Optional

**Response:** Same as `/tiles/me` but with `"is_anonymous": true` and temporary board ID

---

## Timers

Countdown timer management for shows.

### GET /timers

Get paginated list of timers.

**Authentication:** Optional

**Query Parameters:**
- `page` (integer, default: 1) - Page number
- `limit` (integer, default: 50, max: 100) - Items per page
- `show_id` (string, optional) - Filter by show ID
- `is_active` (boolean, optional) - Filter by active status
- `order_by` (string, default: "created_at") - Sort field
- `order_dir` (string, default: "desc") - Sort direction

**Response:**
```json
{
  "timers": [
    {
      "id": "tmr_abc123",
      "title": "Commercial Break",
      "duration": 300,
      "created_by": "usr_abc123",
      "show_id": "Y2kz75uBC8",
      "starts_at": "2024-01-15T20:30:00Z",
      "expires_at": "2024-01-15T20:35:00Z",
      "is_active": true,
      "settings": {
        "color": "#ff0000",
        "sound": true
      },
      "created_at": "2024-01-15T20:25:00Z",
      "updated_at": "2024-01-15T20:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total_count": 25,
    "total_pages": 1,
    "has_next": false,
    "has_prev": false
  },
  "filters": {
    "show_id": "Y2kz75uBC8",
    "is_active": true,
    "order_by": "created_at",
    "order_dir": "desc"
  }
}
```

### GET /timers/:id

Get individual timer by ID.

**Authentication:** Optional

**Path Parameters:**
- `id` (string) - Timer ID

**Response:** Individual timer object

### POST /timers

Create a new timer.

**Authentication:** Required

**Request Body:**
```json
{
  "title": "Commercial Break",
  "duration": 300,
  "show_id": "Y2kz75uBC8",
  "settings": {
    "color": "#ff0000"
  }
}
```

**Response:** Created timer object

### PUT /timers/:id

Update an existing timer.

**Authentication:** Required (timer owner only)

**Path Parameters:**
- `id` (string) - Timer ID

**Request Body:** Same as POST, only provided fields are updated

**Response:** Updated timer object

### DELETE /timers/:id

Delete a timer.

**Authentication:** Required (timer owner only)

**Path Parameters:**
- `id` (string) - Timer ID

**Response:**
```json
{
  "success": true
}
```

### POST /timers/:id/start

Start/activate a timer.

**Authentication:** Required (timer owner only)

**Path Parameters:**
- `id` (string) - Timer ID

**Response:** Updated timer object with `starts_at` and `expires_at` set

### POST /timers/:id/stop

Stop/deactivate a timer.

**Authentication:** Required (timer owner only)

**Path Parameters:**
- `id` (string) - Timer ID

**Response:** Updated timer object with `is_active: false`

---

## Chat

Real-time chat functionality via Server-Sent Events.

### GET /chat/stream

Connect to chat SSE stream.

**Authentication:** Optional (affects permissions)

**Response:** Server-Sent Events stream

### POST /chat

Send a chat message.

**Authentication:** Required

**Request Body:**
```json
{
  "contents": "Hello everyone!"
}
```

**Response:**
```json
{
  "success": true,
  "message": {
    "id": "msg_abc123",
    "show_id": "Y2kz75uBC8",
    "player_id": "usr_abc123",
    "contents": "Hello everyone!",
    "system": false,
    "created_at": "2024-01-15T20:30:00Z"
  }
}
```

### POST /chat/s

Send a system message.

**Authentication:** Required (admin/moderator)

**Request Body:** Same as regular chat message

---

## Error Handling

All API errors follow a consistent format:

```json
{
  "message": "Human-readable error message",
  "code": 0x0101
}
```

### Common Error Codes

- `0x0101` - Database connection failed
- `0x0201` - Query execution failed
- `0x0301` - Show not found
- `0x0401` - Authentication required
- `0x0501` - Insufficient tiles for board generation
- `0x0601` - Tile not found
- `0x0701` - Timer not found
- `0x0721` - Authentication required for timer operations

---

## Pagination

Endpoints that return lists support pagination with the following parameters:

**Query Parameters:**
- `page` (integer, default: 1) - Page number (1-based)
- `limit` (integer, default: 50, max: 100) - Items per page

**Response Format:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total_count": 150,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

---

## Rate Limiting

- API endpoints are rate limited by IP address
- Authenticated users have higher limits
- SSE connections have separate limits

## Data Types

- **Timestamps:** ISO 8601 format with timezone (e.g., `2024-01-15T20:30:00Z`)
- **IDs:** 10-character alphanumeric strings
- **Booleans:** Standard JSON booleans
- **Numbers:** Integers or floats as appropriate
- **JSON Objects:** Flexible structure for settings/metadata