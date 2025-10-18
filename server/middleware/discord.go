package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

// Discord OAuth configuration
var (
	discordOAuthConfig *oauth2.Config
	discordAPIBaseURL  = "https://discord.com/api/v10"
)

// DiscordUser represents a Discord user from the API
type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Verified      bool   `json:"verified"`
}

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

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		fmt.Println("Warning: Discord OAuth not configured. Missing DISCORD_CLIENT_ID, DISCORD_CLIENT_SECRET, or DISCORD_REDIRECT_URL")
		return
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
func GetDiscordUser(token *oauth2.Token) (*DiscordUser, error) {
	client := discordOAuthConfig.Client(context.Background(), token)

	resp, err := client.Get(discordAPIBaseURL + "/users/@me")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Discord user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Discord API returned status %d", resp.StatusCode)
	}

	var user DiscordUser
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

// DiscordAuthMiddleware - require a valid Discord session
func DiscordAuthMiddleware(c *fiber.Ctx) error {
	discordToken := c.Cookies("discord-token")
	if discordToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Discord session cookie",
		})
	}

	// Parse the token from the cookie
	token := &oauth2.Token{
		AccessToken: discordToken,
		TokenType:   "Bearer",
	}

	// Verify the token is still valid by fetching user info
	user, err := GetDiscordUser(token)
	if err != nil {
		// Token is invalid, clear the cookie
		c.ClearCookie("discord-token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired Discord session",
		})
	}

	// Store user info in context
	c.Locals("discord_user", user)
	return c.Next()
}

// OptionalDiscordAuthMiddleware - optional Discord session validation
func OptionalDiscordAuthMiddleware(c *fiber.Ctx) error {
	discordToken := c.Cookies("discord-token")
	if discordToken == "" {
		return c.Next()
	}

	// Parse the token from the cookie
	token := &oauth2.Token{
		AccessToken: discordToken,
		TokenType:   "Bearer",
	}

	// Verify the token is still valid by fetching user info
	user, err := GetDiscordUser(token)
	if err == nil {
		c.Locals("discord_user", user)
	} else {
		// Token is invalid, clear the cookie
		c.ClearCookie("discord-token")
	}

	return c.Next()
}

// GetDiscordUserFromContext retrieves the Discord user from the context
func GetDiscordUserFromContext(c *fiber.Ctx) (*DiscordUser, error) {
	raw := c.Locals("discord_user")
	if raw == nil {
		return nil, fmt.Errorf("no Discord user in context")
	}

	user, ok := raw.(*DiscordUser)
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
		Expires:  time.Now().Add(24 * time.Hour), // Discord tokens typically last 24 hours
		HTTPOnly: true,
		Secure:   os.Getenv("NODE_ENV") == "production", // Only secure in production
		SameSite: "Lax",
	}
	c.Cookie(cookie)
}

// ClearDiscordSessionCookie clears the Discord session cookie
func ClearDiscordSessionCookie(c *fiber.Ctx) {
	c.ClearCookie("discord-token")
}
