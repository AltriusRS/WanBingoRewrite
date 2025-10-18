package callback

import (
	"context"
	"fmt"
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

	// Find or create player in database
	ctx := context.Background()
	player, err := db.FindOrCreatePlayer(ctx, discordUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to create/find player", 500))
	}

	// Create session for player
	sessionID, err := db.CreateSession(ctx, player.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to create session", 500))
	}

	// Set session cookie
	utils.Debugf("Setting session cookie for player %s with session ID: %s", player.ID, sessionID)
	domain := os.Getenv("COOKIE_DOMAIN")
	if domain == "" {
		domain = ".bingo.local" // Default to parent domain for production
	}

	utils.Debugf("Using cookie domain '%s' for session cookie", domain)

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
		HTTPOnly: true,
		Secure:   true,   // HTTPS enabled via Caddy proxy
		SameSite: "None", // Allow cross-site cookies with HTTPS
	})

	// Redirect to frontend application
	frontendURL := middleware.GetFrontendURL()
	if frontendURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Frontend URL not configured", 500))
	}

	// Return HTML page that redirects after cookies are set
	// This gives the browser time to process the cookies before redirecting
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Authentication Successful</title>
</head>
<body>
    <p>Authentication successful! Redirecting...</p>
    <script>
        // Redirect after a delay to ensure cookies are fully set
        setTimeout(function() {
            window.location.href = '%s';
        }, 500);
    </script>
</body>
</html>`, frontendURL)

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}
