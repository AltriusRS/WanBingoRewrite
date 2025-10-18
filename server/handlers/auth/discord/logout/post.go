package logout

import (
	"wanshow-bingo/middleware"

	"github.com/gofiber/fiber/v2"
)

func Post(c *fiber.Ctx) error {
	// Clear Discord session cookie
	middleware.ClearDiscordSessionCookie(c)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully logged out from Discord",
	})
}
