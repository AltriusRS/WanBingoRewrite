package discord

import (
	"wanshow-bingo/handlers/auth/discord/callback"
	"wanshow-bingo/handlers/auth/discord/login"
	"wanshow-bingo/handlers/auth/discord/logout"

	"github.com/gofiber/fiber/v2"
)

func Register(router fiber.Router) {
	// Discord OAuth routes
	router.Get("/login", login.Get)

	router.Get("/callback", callback.Get)

	router.Post("/discord/logout", logout.Post)
}
