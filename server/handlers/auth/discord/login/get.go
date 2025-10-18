package login

import (
	"crypto/rand"
	"encoding/hex"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	// Generate a random state parameter for security
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to generate state", 500))
	}
	state := hex.EncodeToString(stateBytes)

	// Store state in a secure cookie for validation
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		Expires:  c.Context().Time().Add(10 * 60), // 10 minutes
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
	})

	// Get Discord OAuth URL
	authURL := middleware.GetDiscordAuthURL(state)
	if authURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Discord OAuth not configured", 500))
	}

	// Redirect to Discord OAuth
	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}
