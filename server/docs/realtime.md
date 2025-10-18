# Real-time Events (SSE)

This document describes the Server-Sent Events (SSE) system used for real-time communication in the WAN Bingo application.

## Overview

The application uses Server-Sent Events for real-time features including:
- Live chat messaging
- Timer expiration notifications
- Show status updates
- User presence information

## Connection Details

### Chat Stream

**Endpoint:** `GET /chat/stream`

**Authentication:** Optional (affects permissions)

**Headers:**
```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
X-Accel-Buffering: no
```

### Host Stream

**Endpoint:** `GET /host/stream` (admin/moderator only)

**Authentication:** Required (host permissions)

**Purpose:** Receives administrative events like timer expirations

## Event Format

All SSE events follow this JSON structure:

```json
{
  "id": "unique-event-id",
  "opcode": "event.type",
  "data": {
    // Event-specific payload
  }
}
```

**Fields:**
- `id` - Unique event identifier (generated using nanoid)
- `opcode` - Event type identifier
- `data` - Event payload (varies by event type)

## Chat Events

### hub.connected

Sent when client successfully connects to chat stream.

```json
{
  "id": "abc123def4",
  "opcode": "hub.connected",
  "data": "Connected to chat hub."
}
```

### hub.authenticated

Sent after connection to indicate client permissions.

```json
{
  "id": "abc123def4",
  "opcode": "hub.authenticated",
  "data": {
    "canChat": true,
    "canHost": false,
    "canModerate": false
  }
}
```

### hub.connections.count

Sent periodically to indicate total connected users.

```json
{
  "id": "abc123def4",
  "opcode": "hub.connections.count",
  "data": 42
}
```

### chat.message

Sent when a new chat message is posted.

```json
{
  "id": "msg_abc123",
  "opcode": "chat.message",
  "data": {
    "id": "msg_abc123",
    "show_id": "Y2kz75uBC8",
    "player_id": "usr_abc123",
    "contents": "Hello everyone!",
    "system": false,
    "replying": null,
    "created_at": "2024-01-15T20:30:00Z"
  }
}
```

### chat.players

Sent on connection to provide information about chat participants.

```json
{
  "id": "plr_list_001",
  "opcode": "chat.players",
  "data": {
    "players": [
      {
        "id": "usr_abc123",
        "display_name": "LinusTech#1337",
        "avatar": "https://cdn.discordapp.com/avatars/..."
      },
      {
        "id": "usr_def456",
        "display_name": "LukeLafr#0001",
        "avatar": "https://cdn.discordapp.com/avatars/..."
      }
    ],
    "count": 2
  }
}
```

## Show Events

### whenplane.aggregate

Sent when show status updates are received from the whenplane service.

```json
{
  "id": "wp_agg_001",
  "opcode": "whenplane.aggregate",
  "data": {
    "id": "Y2kz75uBC8",
    "youtube_id": "YVHXYqMPyzc",
    "scheduled_time": "2025-10-11T00:30:00Z",
    "actual_start_time": "2025-10-11T00:05:06Z",
    "thumbnail": "https://pbs.floatplane.com/stream_thumbnails/...",
    "metadata": {
      "title": "Piracy Is Dangerous And Harmful",
      "fp_vod": "w3A5fKcfTi",
      "hosts": ["Linus", "Luke"],
      "live": true
    }
  }
}
```

## Timer Events

### timer.expired

Sent when a timer expires (to host stream only).

```json
{
  "id": "tmr_exp_001",
  "opcode": "timer.expired",
  "data": {
    "timer_id": "tmr_abc123",
    "title": "Commercial Break",
    "show_id": "Y2kz75uBC8",
    "created_by": "usr_abc123",
    "expired_at": "2024-01-15T20:35:00Z"
  }
}
```

## Event Flow Examples

### New User Connecting to Chat

1. Client connects to `GET /chat/stream`
2. Server sends `hub.connected`
3. Server sends `hub.authenticated` with permissions
4. Server sends `chat.players` with participant list
5. Server sends recent `chat.message` events (chat history)
6. Server periodically sends `hub.connections.count`

### Timer Lifecycle

1. User creates timer via `POST /timers`
2. User starts timer via `POST /timers/:id/start`
3. Background monitor checks every 2 seconds for expired timers (high accuracy for stream delays)
4. When timer expires:
   - `timer.expired` event sent to host stream
   - Timer automatically deactivated in database

### Message Broadcasting

1. User sends message via `POST /chat`
2. Message saved to database
3. `chat.message` event broadcast to all connected chat clients
4. Message appears in real-time for all users

## Client Implementation

### JavaScript SSE Client

```javascript
const eventSource = new EventSource('/chat/stream');

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  handleEvent(data.opcode, data.data);
};

function handleEvent(opcode, data) {
  switch (opcode) {
    case 'hub.connected':
      console.log('Connected:', data);
      break;
    case 'chat.message':
      displayMessage(data);
      break;
    case 'timer.expired':
      handleTimerExpired(data);
      break;
  }
}
```

### Connection Management

- **Reconnection:** Clients should implement automatic reconnection on disconnect
- **Heartbeat:** Server sends periodic connection count updates
- **Authentication:** Permission changes require reconnection
- **Rate Limiting:** SSE connections have separate rate limits from HTTP APIs

## Server Implementation Details

### Hub Architecture

- **Chat Hub:** Handles user chat and general events
- **Host Hub:** Handles administrative events (timers, moderation)
- **Broadcasting:** Events sent to appropriate hub based on target audience
- **Connection Tracking:** Real-time user count maintained

### Background Services

- **Timer Monitor:** Cron job checking for expired timers every 2 seconds
- **Event Broadcasting:** Automatic distribution to connected clients
- **Connection Cleanup:** Automatic removal of disconnected clients

### Performance Considerations

- **Connection Limits:** Maximum concurrent SSE connections per IP
- **Message Buffering:** Internal queues prevent message loss
- **Event Filtering:** Clients only receive relevant events
- **Memory Management:** Automatic cleanup of disconnected clients

## Error Handling

### Connection Errors

- **Network Issues:** Clients should implement exponential backoff reconnection
- **Server Errors:** SSE stream may terminate with error messages
- **Authentication Failures:** Stream closes with 401-equivalent behavior

### Event Processing Errors

- **Invalid JSON:** Events logged but not broadcast
- **Missing Fields:** Events validated before broadcasting
- **Database Errors:** Logged with fallback behavior

## Security Considerations

- **Origin Validation:** CORS headers restrict allowed origins
- **Authentication:** Session validation on connection
- **Rate Limiting:** Connection and message rate limits
- **Input Sanitization:** All event data validated and sanitized
- **Permission Checks:** Events filtered based on user permissions</content>
</xai:function_call">Real-time Events Documentation