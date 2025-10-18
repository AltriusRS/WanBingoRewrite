package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/db/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

// Discord OAuth configuration
var (
	discordOAuthConfig *oauth2.Config
	discordAPIBaseURL  = "https://discord.com/api/v10"
	frontendURL        string
)

// DiscordGuild represents a Discord guild (server) the user is in
type DiscordGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions string `json:"permissions"`
}

// InitDiscordOAuth initializes the Discord OAuth configuration
func InitDiscordOAuth() {
	clientID := os.Getenv("DISCORD_CLIENT_ID")
	clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	redirectURL := os.Getenv("DISCORD_REDIRECT_URL")
	frontendURL = os.Getenv("FRONTEND_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		fmt.Println("Warning: Discord OAuth not configured. Missing DISCORD_CLIENT_ID, DISCORD_CLIENT_SECRET, or DISCORD_REDIRECT_URL")
		return
	}

	if frontendURL == "" {
		fmt.Println("Warning: FRONTEND_URL not set, OAuth callback will not redirect properly")
	}

	discordOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"identify",
			"guilds",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
}

// GetDiscordAuthURL generates the Discord OAuth authorization URL
func GetDiscordAuthURL(state string) string {
	if discordOAuthConfig == nil {
		return ""
	}
	return discordOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// GetFrontendURL returns the configured frontend URL for redirects
func GetFrontendURL() string {
	return frontendURL
}

// ExchangeCodeForToken exchanges the authorization code for an access token
func ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	if discordOAuthConfig == nil {
		return nil, fmt.Errorf("Discord OAuth not configured")
	}

	ctx := context.Background()
	token, err := discordOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

// GetDiscordUser fetches the Discord user information using the access token
func GetDiscordUser(token *oauth2.Token) (*models.DiscordUser, error) {
	client := discordOAuthConfig.Client(context.Background(), token)

	resp, err := client.Get(discordAPIBaseURL + "/users/@me")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Discord user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Discord API returned status %d", resp.StatusCode)
	}

	var user models.DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode Discord user: %w", err)
	}

	return &user, nil
}

// GetDiscordGuilds fetches the Discord guilds the user is in
func GetDiscordGuilds(token *oauth2.Token) ([]DiscordGuild, error) {
	client := discordOAuthConfig.Client(context.Background(), token)

	resp, err := client.Get(discordAPIBaseURL + "/users/@me/guilds")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Discord guilds: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Discord API returned status %d", resp.StatusCode)
	}

	var guilds []DiscordGuild
	if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return nil, fmt.Errorf("failed to decode Discord guilds: %w", err)
	}

	return guilds, nil
}

// AuthMiddleware - require a valid session
func AuthMiddleware(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing session cookie",
		})
	}

	// Validate session and get player
	ctx := context.Background()
	player, err := db.ValidateSession(ctx, sessionID)
	if err != nil {
		// Session is invalid, clear the cookie
		c.ClearCookie("session_id")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired session",
		})
	}

	// Store player info in context
	c.Locals("player", player)
	return c.Next()
}

// OptionalPlayerAuthMiddleware - optional session validation
func OptionalPlayerAuthMiddleware(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.Next()
	}

	// Validate session and get player
	ctx := context.Background()
	player, err := db.ValidateSession(ctx, sessionID)
	if err == nil {
		c.Locals("player", player)
	} else {
		// Session is invalid, clear the cookie
		c.ClearCookie("session_id")
	}

	return c.Next()
}

// GetPlayerFromContext retrieves the player from the context
func GetPlayerFromContext(c *fiber.Ctx) (*models.Player, error) {
	raw := c.Locals("player")
	if raw == nil {
		return nil, fmt.Errorf("no player in context")
	}

	player, ok := raw.(*models.Player)
	if !ok {
		return nil, fmt.Errorf("invalid player type in context")
	}

	return player, nil
}

// GetDiscordUserFromContext retrieves the Discord user from the context
func GetDiscordUserFromContext(c *fiber.Ctx) (*models.DiscordUser, error) {
	raw := c.Locals("discord_user")
	if raw == nil {
		return nil, fmt.Errorf("no Discord user in context")
	}

	user, ok := raw.(*models.DiscordUser)
	if !ok {
		return nil, fmt.Errorf("invalid Discord user type in context")
	}

	return user, nil
}

// SetDiscordSessionCookie sets a secure Discord session cookie
func SetDiscordSessionCookie(c *fiber.Ctx, token *oauth2.Token) {
	cookie := &fiber.Cookie{
		Name:     "discord-token",
		Value:    token.AccessToken,
		Path:     "/",
		Domain:   "api.bingo.local",              // Match other cookies
		Expires:  time.Now().Add(24 * time.Hour), // Discord tokens typically last 24 hours
		HTTPOnly: true,
		Secure:   true,   // HTTPS enabled via Caddy proxy
		SameSite: "None", // Allow cross-site with HTTPS
	}
	c.Cookie(cookie)
}

// ClearDiscordSessionCookie clears the Discord session cookie
func ClearDiscordSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "discord-token",
		Value:    "",
		Path:     "/",
		Domain:   "api.bingo.local",
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})
}

// RequirePermissionMiddleware checks if the authenticated user has a specific permission
func RequirePermissionMiddleware(permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		player, err := GetPlayerFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		// Find the permission constant by name
		var requiredPerm models.Permission
		found := false
		for perm, name := range models.PermissionNames {
			if name == permissionName {
				requiredPerm = perm
				found = true
				break
			}
		}

		if !found {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid permission name",
			})
		}

		if !player.Permissions.HasPermission(requiredPerm) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}
