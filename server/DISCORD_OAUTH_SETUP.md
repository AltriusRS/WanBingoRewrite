# Discord OAuth 2.0 Setup Guide

This guide explains how to set up Discord OAuth 2.0 authentication in your Go server.

## Prerequisites

1. A Discord application registered at https://discord.com/developers/applications
2. Go 1.25.2 or later
3. The required dependencies (already added to go.mod)

## Environment Variables

Add the following environment variables to your `.env` file or system environment:

```bash
# Discord OAuth Configuration
DISCORD_CLIENT_ID=your_discord_client_id_here
DISCORD_CLIENT_SECRET=your_discord_client_secret_here
DISCORD_REDIRECT_URL=http://localhost:8080/auth/discord/callback
FRONTEND_URL=http://localhost:3000

# Optional: Set to production for secure cookies
NODE_ENV=development
```

## Discord Application Setup

1. Go to https://discord.com/developers/applications
2. Create a new application or select an existing one
3. Go to the "OAuth2" section
4. Add the following redirect URI: `http://localhost:8080/auth/discord/callback`
5. Copy the Client ID and Client Secret
6. Set the required scopes:
   - `identify` - Get user's basic information
   - `email` - Get user's email address
   - `guilds` - Get user's guilds (servers)

## API Endpoints

The following endpoints are now available:

### Authentication Endpoints

- `GET /auth/discord/login` - Initiates Discord OAuth flow
- `GET /auth/discord/callback` - Handles Discord OAuth callback and redirects to frontend
- `POST /auth/discord/logout` - Logs out the Discord user

### User Endpoints

#### Public Endpoints (no authentication required)
- `GET /users` - Get all players (partial profiles)
- `GET /users/:identifier` - Get player by ID or display_name (partial profile)

#### Protected Endpoints (require authentication)
- `GET /users/me` - Get current user information (full profile)

## Usage Examples

### Frontend Integration

```javascript
// Redirect to Discord login
window.location.href = 'http://localhost:8080/auth/discord/login';
// After authentication, user will be redirected back to FRONTEND_URL

// Check if user is authenticated (call this on page load)
fetch('http://localhost:8080/users/me', {
  credentials: 'include'
})
.then(response => response.json())
.then(data => {
  if (data.success) {
    console.log('Current User:', data.user);
  } else {
    console.log('Not authenticated');
  }
});

// Get all players
fetch('http://localhost:8080/users')
.then(response => response.json())
.then(data => {
  if (data.success) {
    console.log('All Players:', data.players);
  }
});

// Get specific player by ID or display_name
fetch('http://localhost:8080/users/somePlayerId')
.then(response => response.json())
.then(data => {
  if (data.success) {
    console.log('Player:', data.player);
  }
});

// Logout
fetch('/auth/discord/logout', {
  method: 'POST',
  credentials: 'include'
})
.then(response => response.json())
.then(data => {
  console.log(data.message);
});
```

### Backend Middleware Usage

```go
// Get current user (full profile)
app.Get("/users/me", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
    player, err := middleware.GetPlayerFromContext(c)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"success": true, "user": fiber.Map{
        "id": player.ID,
        "did": player.DID,
        "display_name": player.DisplayName,
        "avatar": player.Avatar,
        "settings": player.Settings,
        "score": player.Score,
        "created_at": player.CreatedAt,
        "updated_at": player.UpdatedAt,
    }})
})

// Get player by identifier (partial profile)
app.Get("/users/:identifier", func(c *fiber.Ctx) error {
    identifier := c.Params("identifier")
    player, err := db.GetPlayerByIdentifier(context.Background(), identifier)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Player not found"})
    }
    return c.JSON(fiber.Map{"success": true, "player": fiber.Map{
        "id": player.ID,
        "display_name": player.DisplayName,
        "avatar": player.Avatar,
        "score": player.Score,
        "created_at": player.CreatedAt,
    }})
})

// Optional authentication
app.Get("/optional", middleware.OptionalPlayerAuthMiddleware, func(c *fiber.Ctx) error {
    player, err := middleware.GetPlayerFromContext(c)
    if err != nil {
        return c.JSON(fiber.Map{"authenticated": false})
    }
    return c.JSON(fiber.Map{"authenticated": true, "player": player})
})
```

## Security Features

- **State Parameter**: CSRF protection using random state parameter
- **Secure Cookies**: HTTP-only cookies with proper security settings
- **Token Validation**: Automatic token validation on each request
- **CORS Configuration**: Proper CORS setup for Discord OAuth

## Error Handling

The API returns consistent error responses:

```json
{
  "success": false,
  "message": "Error description",
  "code": 400
}
```

Common error codes:
- `400` - Bad Request (missing parameters, invalid state)
- `401` - Unauthorized (invalid/expired token)
- `500` - Internal Server Error (Discord API errors, configuration issues)

## Testing

1. Set environment variables including `FRONTEND_URL=http://localhost:3000` (or your frontend URL)
2. Start your server: `go run main.go`
3. Start your frontend application on the configured FRONTEND_URL
4. Visit `http://localhost:8080/auth/discord/login` from your frontend
5. Complete the Discord OAuth flow
6. You should be redirected back to your frontend with authentication cookies set
7. Test protected endpoints with the session cookie

## Production Considerations

1. Set `NODE_ENV=production` for secure cookies
2. Use HTTPS for the redirect URL
3. Update CORS origins for your production domain
4. Consider implementing token refresh logic for long-lived sessions
5. Add rate limiting for OAuth endpoints
6. Implement proper logging and monitoring
