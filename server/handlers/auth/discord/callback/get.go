package callback

import (
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
	user, err := middleware.GetDiscordUser(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to fetch Discord user", 500))
	}

	// Set Discord session cookie
	middleware.SetDiscordSessionCookie(c, token)

	// Return success response with user info
	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":            user.ID,
			"username":      user.Username,
			"discriminator": user.Discriminator,
			"email":         user.Email,
			"avatar":        user.Avatar,
			"verified":      user.Verified,
		},
		"message": "Successfully authenticated with Discord",
	})
}
