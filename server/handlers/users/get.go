package users

import (
	"context"
	"wanshow-bingo/db"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// GetByIdentifier returns partial player profile by ID or display_name
func GetByIdentifier(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	if identifier == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Identifier parameter required", 400))
	}

	player, err := db.GetPlayerByIdentifier(context.Background(), identifier)
	if err != nil {
		if err.Error() == "player not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.NewApiError("Player not found", 404))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to fetch player", 500))
	}

	// Return partial profile (exclude sensitive info like DID)
	return c.JSON(fiber.Map{
		"success": true,
		"player": fiber.Map{
			"id":           player.ID,
			"display_name": player.DisplayName,
			"avatar":       player.Avatar,
			"score":        player.Score,
			"created_at":   player.CreatedAt,
		},
	})
}

// GetAll returns all players (partial profiles)
func GetAll(c *fiber.Ctx) error {
	players, err := db.GetAllPlayers(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to fetch players", 500))
	}

	// Return partial profiles for all players
	var playerProfiles []fiber.Map
	for _, player := range players {
		playerProfiles = append(playerProfiles, fiber.Map{
			"id":           player.ID,
			"display_name": player.DisplayName,
			"avatar":       player.Avatar,
			"score":        player.Score,
			"created_at":   player.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"players": playerProfiles,
	})
}
