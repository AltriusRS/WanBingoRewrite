package logout

import (
	"context"
	"wanshow-bingo/db"

	"github.com/gofiber/fiber/v2"
)

func Post(c *fiber.Ctx) error {
	// Get session ID from cookie
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		// Delete session from database
		ctx := context.Background()
		err := db.DeleteSession(ctx, sessionID)
		if err != nil {
			// Log error but continue with cookie clearing
			c.Context().Logger().Printf("Failed to delete session: %v", err)
		}
	}

	// Clear session cookie
	c.ClearCookie("session_id")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully logged out",
	})
}
