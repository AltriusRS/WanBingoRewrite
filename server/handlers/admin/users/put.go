package users

import (
	"context"
	"wanshow-bingo/db"

	"github.com/gofiber/fiber/v2"
)

type UpdatePermissionsRequest struct {
	Permissions map[string]bool `json:"permissions"`
}

func UpdateUserPermissions(c *fiber.Ctx) error {
	userID := c.Params("id")

	var req UpdatePermissionsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.Background()
	player, err := db.GetPlayerByID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update permissions from the request
	player.Permissions.SetPermissionsFromMap(req.Permissions)

	// Save the updated player
	err = db.PersistPlayer(ctx, player)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user permissions",
		})
	}

	permissions := player.Permissions.GetAllPermissions()
	return c.JSON(fiber.Map{
		"user_id":      player.ID,
		"display_name": player.DisplayName,
		"permissions":  permissions,
	})
}
