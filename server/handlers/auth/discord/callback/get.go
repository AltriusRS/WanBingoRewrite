package callback

import (
	"context"
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
	if storedState == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Missing state cookie", 400))
	}

	// Verify state parameter matches stored state
	if state != storedState {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Invalid state parameter", 400))
	}

	// Clear the state cookie
	c.ClearCookie("oauth_state")

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
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  c.Context().Time().Add(24 * 60 * 60), // 24 hours
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
	})

	// Redirect to frontend application
	frontendURL := middleware.GetFrontendURL()
	if frontendURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Frontend URL not configured", 500))
	}

	return c.Redirect(frontendURL, fiber.StatusTemporaryRedirect)
}
