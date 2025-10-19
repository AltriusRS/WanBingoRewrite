package callback

import (
	"context"
	"os"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	// Get the state parameter from query
	state := c.Query("state")
	if state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Missing state parameter", 400))
	}

	// Get the stored state from cookie
	storedState := c.Cookies("oauth_state")
	utils.Debugf("OAuth callback - state from URL: %s", state)
	utils.Debugf("OAuth callback - stored state from cookie: %s", storedState)

	if storedState == "" {
		// Debug: Log all cookies to help troubleshoot
		cookies := c.GetReqHeaders()["Cookie"]
		utils.Debugf("Available cookies: %v", cookies)
		utils.Debugf("oauth_state cookie: %s", c.Cookies("oauth_state"))

		// In development, allow bypassing state validation if explicitly requested
		if os.Getenv("SKIP_OAUTH_STATE_VALIDATION") == "true" {
			utils.Debugf("Skipping OAuth state validation (development mode)")
			storedState = state // Use the state from the URL parameter
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Missing state cookie - please try logging in again", 400))
		}
	}

	// Verify state parameter matches stored state
	if state != storedState {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Invalid state parameter", 400))
	}

	// Clear the state cookie by setting it to expire immediately
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		Domain:   "", // Use current domain
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	// Get the authorization code
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Missing authorization code", 400))
	}

	// Exchange code for token
	token, err := middleware.ExchangeCodeForToken(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to exchange code for token", 500))
	}

	// Get Discord user information
	discordUser, err := middleware.GetDiscordUser(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to fetch Discord user", 500))
	}
	utils.Debugf("AUTH: Discord user data - ID: %s, Username: %s, Avatar: %s", discordUser.ID, discordUser.Username, discordUser.Avatar)

	// Find or create player in database
	utils.Debugf("AUTH: About to call FindOrCreatePlayer for Discord user %s with avatar hash: %s", discordUser.ID, discordUser.Avatar)
	ctx := context.Background()
	player, err := db.FindOrCreatePlayer(ctx, discordUser)
	if err != nil {
		utils.Debugf("AUTH: FindOrCreatePlayer failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to create/find player", 500))
	}
	avatarStr := ""
	if player.Avatar != nil {
		avatarStr = *player.Avatar
	}
	utils.Debugf("AUTH: FindOrCreatePlayer completed, player ID: %s, avatar: %s", player.ID, avatarStr)

	// Create session for player
	sessionID, err := db.CreateSession(ctx, player.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to create session", 500))
	}

	// Set session cookie
	utils.Debugf("Setting session cookie for player %s with session ID: %s", player.ID, sessionID)
	domain := os.Getenv("COOKIE_DOMAIN")
	env := os.Getenv("ENV")
	nodeEnv := os.Getenv("NODE_ENV")
	host := c.Hostname()
	utils.Debugf("Environment variables - ENV: '%s', NODE_ENV: '%s', COOKIE_DOMAIN: '%s'", env, nodeEnv, domain)
	utils.Debugf("Request hostname: '%s'", host)

	if domain == "" {
		// For local development, don't set a domain to allow localhost cookies
		if env == "development" || nodeEnv == "development" || host == "localhost" || host == "127.0.0.1" {
			domain = "" // Empty domain for localhost
		} else {
			// Try to extract domain from hostname
			// For example, if hostname is "api.bingo-demo.totallyfake.dev", domain should be ".bingo-demo.totallyfake.dev"
			if len(host) > 4 && host[:4] == "api." {
				domain = "." + host[4:] // Remove "api." prefix
			} else {
				domain = ".bingo.local" // Fallback
			}
		}
	}

	utils.Debugf("Using cookie domain '%s' for session cookie", domain)

	// Determine if we should use secure cookies
	isSecure := os.Getenv("ENV") != "development" && os.Getenv("NODE_ENV") != "development"
	sameSite := "Lax"
	if isSecure {
		sameSite = "None" // Allow cross-site cookies with HTTPS
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
		HTTPOnly: true,
		Secure:   isSecure, // Only secure in production
		SameSite: sameSite,
	})

	// Redirect to frontend application
	frontendURL := middleware.GetFrontendURL()
	if frontendURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Frontend URL not configured", 500))
	}

	// Redirect immediately to frontend
	// The cookies should be set by the time the browser processes this redirect
	return c.Redirect(frontendURL, fiber.StatusFound)
}
