package me

import (
	"wanshow-bingo/avatar"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func Put(c *fiber.Ctx) error {
	player, err := middleware.GetPlayerFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewApiError("Not authenticated", 401))
	}

	var body struct {
		DisplayName string         `json:"display_name"`
		Avatar      *string        `json:"avatar"`
		Settings    map[string]any `json:"settings"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewApiError("Invalid request body", 400))
	}

	// Update player
	player.DisplayName = body.DisplayName
	player.Avatar = body.Avatar
	player.Settings = &body.Settings

	if err := middleware.UpdatePlayer(c.Context(), player); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to update player", 500))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":           player.ID,
			"did":          player.DID,
			"display_name": player.DisplayName,
			"avatar":       avatar.GetAvatarURL(avatarKey(player.Avatar)),
			"settings":     player.Settings,
			"score":        player.Score,
			"created_at":   player.CreatedAt,
			"updated_at":   player.UpdatedAt,
		},
	})
}
