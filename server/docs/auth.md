# Authentication & Authorization

This document describes the authentication and authorization system used in the WAN Bingo application.

## Overview

The application uses **Discord OAuth2** for user authentication with **session-based** authorization. This provides secure, scalable user management while leveraging Discord's robust identity system.

## OAuth2 Flow

### 1. Authorization Request

**Endpoint:** `GET /auth/discord/login`

**Purpose:** Redirects user to Discord OAuth authorization page

**Query Parameters:**
- `state` (string) - CSRF protection token (auto-generated)

**Flow:**
1. User clicks "Login with Discord"
2. Server generates random `state` token
3. Server stores `state` in cookie
4. Server redirects to Discord OAuth URL

### 2. Discord Authorization

**Discord URL:** `https://discord.com/api/oauth2/authorize`

**Parameters:**
- `client_id` - Application's Discord client ID
- `redirect_uri` - Callback URL (`/auth/discord/callback`)
- `response_type` - `code`
- `scope` - `identify` (basic user info)
- `state` - CSRF token from step 1

### 3. OAuth Callback

**Endpoint:** `GET /auth/discord/callback`

**Purpose:** Handle Discord's authorization response

**Query Parameters:**
- `code` (string) - Authorization code from Discord
- `state` (string) - CSRF token for validation

**Process:**
1. Validate `state` parameter matches stored cookie
2. Exchange `code` for access token
3. Fetch user info from Discord API
4. Find or create user in database
5. Create session token
6. Set session cookie
7. Redirect to frontend application

### 4. Session Management

**Cookie:** `session_id`

**Attributes:**
- Value: 32-character session token
- Path: `/`
- Expires: 24 hours from creation
- HttpOnly: `true` (prevents JavaScript access)
- Secure: `true` in production
- SameSite: `Lax`

## User Registration

### Automatic User Creation

When a user completes OAuth for the first time:

1. **Check Existing User**
   ```sql
   SELECT * FROM players WHERE did = ? AND deleted_at IS NULL
   ```

2. **Create New User** (if not found)
   ```sql
   INSERT INTO players (id, did, display_name, avatar, settings, score)
   VALUES (?, ?, ?, ?, '{}', 0)
   ```

3. **Update Existing User** (if found)
   ```sql
   UPDATE players
   SET display_name = ?, avatar = ?, updated_at = CURRENT_TIMESTAMP
   WHERE id = ?
   ```

### User Data Mapping

| Discord API Field | Database Field | Notes |
|-------------------|----------------|-------|
| `id` | `did` | Discord user ID |
| `username` + `#` + `discriminator` | `display_name` | e.g., "LinusTech#1337" |
| `avatar` | `avatar` | Full Discord CDN URL |

## Session Validation

### Middleware Implementation

**Location:** `middleware/discord.go`

**Function:** `AuthMiddleware` and `RequiredAuthMiddleware`

**Process:**
1. Extract `session_id` from cookies
2. Query database for valid session
   ```sql
   SELECT p.* FROM players p
   JOIN sessions s ON p.id = s.player_id
   WHERE s.id = ? AND s.expires_at > CURRENT_TIMESTAMP
   AND s.deleted_at IS NULL AND p.deleted_at IS NULL
   ```
3. Attach player data to request context
4. Continue to handler or return 401

### Session Expiration

- **Duration:** 24 hours from creation
- **Automatic Cleanup:** Background process removes expired sessions
- **Sliding Expiration:** Sessions extend on activity (optional)

## Authorization Levels

### Permission System

The application uses role-based permissions for different features:

```json
{
  "canChat": true,           // Send chat messages
  "canHost": false,          // Host/admin controls
  "canModerate": false,      // Moderate chat/users
  "canSendMessages": true,   // Alias for canChat
  "canSendWhispers": true,   // Private messages
  "canDeleteOwnMessages": true,
  "canDeleteMessages": false,
  "canBanUsers": false,
  "canKickUsers": false,
  "canMuteUsers": false,
  "canUnmuteUsers": false,
  "canManageChat": false,
  "canManageChatPermissions": false,
  "canSuggestTiles": true,   // Suggest new bingo tiles
  "canReviewTiles": false,   // Review tile suggestions
  "canApproveTiles": false,  // Approve tiles for use
  "canManageTiles": false,   // Full tile management
  "canPromotePlayers": false,
  "canModifyShowData": false,
  "canManageTimers": false   // Create/manage timers
}
```

### Current Implementation

**All Authenticated Users:**
- `canChat`: `true`
- `canSendMessages`: `true`
- `canSendWhispers`: `true`
- `canDeleteOwnMessages`: `true`
- `canSuggestTiles`: `true`

**Guests (Unauthenticated):**
- All permissions: `false`

### Future Extensions

The permission system is designed to be extensible:

```sql
-- Future: Permission table
CREATE TABLE permissions (
    id VARCHAR(10) PRIMARY KEY,
    player_id VARCHAR(10) REFERENCES players(id),
    permission VARCHAR(50) NOT NULL,
    granted_by VARCHAR(10) REFERENCES players(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Security Measures

### CSRF Protection

- **State Parameter:** Random token generated per OAuth request
- **Cookie Storage:** State stored in httpOnly cookie
- **Validation:** Server compares cookie and query parameter

### Session Security

- **HttpOnly Cookies:** Prevent JavaScript access to session tokens
- **Secure Cookies:** HTTPS-only in production
- **SameSite Protection:** CSRF protection for cross-site requests
- **Expiration:** Automatic cleanup of expired sessions

### Rate Limiting

- **OAuth Requests:** Limited per IP address
- **API Calls:** Rate limited by session/IP
- **SSE Connections:** Separate limits for real-time features

### Input Validation

- **Discord Data:** All user data validated before storage
- **Session Tokens:** Format validation and database verification
- **Permissions:** Server-side validation on all actions

## API Authentication

### Required Authentication

**Endpoints requiring authentication:**
- `POST /chat` - Send messages
- `GET /users/me` - User profile
- `GET /tiles/me` - User bingo board
- `POST /timers` - Create timer
- `PUT /timers/:id` - Update timer
- `DELETE /timers/:id` - Delete timer
- `POST /timers/:id/start` - Start timer
- `POST /timers/:id/stop` - Stop timer

### Optional Authentication

**Endpoints with enhanced features when authenticated:**
- `GET /chat/stream` - Full permissions vs. read-only
- `GET /users` - Shows partial profiles for all users
- `GET /tiles` - Same data for all users

## Error Responses

### Authentication Errors

```json
{
  "message": "Authentication required",
  "code": 0x0401
}
```

### Authorization Errors

```json
{
  "message": "Access denied",
  "code": 0x0404
}
```

### OAuth Errors

```json
{
  "message": "Invalid state parameter",
  "code": 0x0A01
}
```

## Client Integration

### JavaScript Implementation

```javascript
// Check authentication status
async function checkAuth() {
  try {
    const response = await fetch('/users/me');
    if (response.ok) {
      const user = await response.json();
      setAuthenticatedUser(user);
    } else {
      setUnauthenticated();
    }
  } catch (error) {
    console.error('Auth check failed:', error);
  }
}

// Login flow
function login() {
  window.location.href = '/auth/discord/login';
}

// Logout
async function logout() {
  await fetch('/auth/discord/logout', { method: 'POST' });
  window.location.reload();
}
```

### Session Persistence

- **Automatic Checks:** Verify session validity on page load
- **Token Refresh:** Handle session expiration gracefully
- **Reconnection:** SSE streams handle authentication changes

## Database Schema

### Sessions Table

```sql
CREATE TABLE sessions (
    id VARCHAR(32) PRIMARY KEY,
    player_id VARCHAR(10) REFERENCES players(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for performance
CREATE INDEX idx_sessions_player_id ON sessions(player_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

### Players Table (Relevant Fields)

```sql
CREATE TABLE players (
    id VARCHAR(10) PRIMARY KEY,
    did VARCHAR(20) UNIQUE NOT NULL,  -- Discord ID
    display_name VARCHAR(100) NOT NULL,
    avatar TEXT,
    settings JSONB DEFAULT '{}'::jsonb,
    score INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

## Monitoring & Analytics

### Authentication Metrics

- **Login Success/Failure Rates**
- **Session Duration Statistics**
- **OAuth Conversion Rates**
- **Permission Usage Patterns**

### Security Monitoring

- **Failed Authentication Attempts**
- **Suspicious Session Activity**
- **Rate Limit Violations**
- **CSRF Attempt Detection**

## Configuration

### Environment Variables

```bash
# Discord OAuth
DISCORD_CLIENT_ID=your_client_id
DISCORD_CLIENT_SECRET=your_client_secret
DISCORD_REDIRECT_URI=https://yourapp.com/auth/discord/callback

# Frontend URL
FRONTEND_URL=https://yourapp.com

# Session Settings
SESSION_DURATION_HOURS=24
```

### OAuth Application Setup

1. **Create Discord Application:** https://discord.com/developers/applications
2. **Add Redirect URI:** `https://yourapp.com/auth/discord/callback`
3. **Copy Client ID/Secret:** To environment variables
4. **Configure Scopes:** `identify` (basic user info)

## Troubleshooting

### Common Issues

**"Invalid state parameter":**
- Check that cookies are enabled
- Verify redirect URI matches Discord app configuration
- Ensure HTTPS in production

**"Session expired":**
- Sessions automatically expire after 24 hours
- Users need to re-authenticate
- Check server time synchronization

**"Access denied":**
- Verify user permissions for requested action
- Check if session is still valid
- Confirm user account is not deleted

### Debug Mode

Enable detailed logging for authentication issues:

```go
// In middleware/discord.go
utils.Debugf("Auth check for session: %s", sessionID)
utils.Debugf("Player found: %+v", player)
```</content>
</xai:function_call">Authentication Documentation