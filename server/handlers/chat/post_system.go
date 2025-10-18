package chat

import (
	"context"
	"wanshow-bingo/db"
	"wanshow-bingo/middleware"

	"github.com/gofiber/fiber/v2"
)

func Post(ctx *fiber.Ctx) error {
	duser, err := middleware.GetDiscordUserFromContext(ctx)

	if err != nil || duser == nil {
		return fiber.ErrUnauthorized
	}

	pool := db.Pool()

	playerRow := pool.QueryRow(context.Background(), "SELECT * FROM players WHERE did = $1", duser.ID)

	if playerRow == nil {
		return fiber.ErrNotFound
	}

	return nil
}
