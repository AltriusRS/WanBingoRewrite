package users

import (
	"context"
	"wanshow-bingo/db"

	"github.com/gofiber/fiber/v2"
)

func GetUserPermissions(c *fiber.Ctx) error {
	userID := c.Params("id")

	ctx := context.Background()
	player, err := db.GetPlayerByID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	permissions := player.Permissions.GetAllPermissions()
	return c.JSON(fiber.Map{
		"user_id":      player.ID,
		"display_name": player.DisplayName,
		"permissions":  permissions,
	})
}
