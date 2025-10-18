package me

import (
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

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
			"avatar":       player.Avatar,
			"settings":     player.Settings,
			"score":        player.Score,
			"created_at":   player.CreatedAt,
			"updated_at":   player.UpdatedAt,
		},
	})
}
