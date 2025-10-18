package chat

import (
	"wanshow-bingo/middleware"

	"github.com/gofiber/fiber/v2"
)

func Post(ctx *fiber.Ctx) error {
	player, err := middleware.GetPlayerFromContext(ctx)

	if err != nil || player == nil {
		return fiber.ErrUnauthorized
	}

	// Player is already authenticated and available
	// The message saving logic should be implemented here
	// For now, just return success
	return ctx.JSON(fiber.Map{"success": true})
}
