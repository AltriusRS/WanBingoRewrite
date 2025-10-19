package me

import (
	"wanshow-bingo/avatar"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// avatarKey safely extracts string value from *string
func avatarKey(avatar *string) string {
	if avatar == nil {
		return ""
	}
	return *avatar
}

func Get(c *fiber.Ctx) error {
	player, err := middleware.GetPlayerFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewApiError("Not authenticated", 401))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":           player.ID,
			"did":          player.DID,
			"display_name": player.DisplayName,
			"avatar":       avatar.GetAvatarURL(avatarKey(player.Avatar)),
			"settings":     player.Settings,
			"permissions":  player.Permissions,
			"score":        player.Score,
			"created_at":   player.CreatedAt,
			"updated_at":   player.UpdatedAt,
		},
	})
}
