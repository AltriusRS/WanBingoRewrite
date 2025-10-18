package auth

import (
	"wanshow-bingo/handlers/auth/discord"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/auth", registerAuth)
}

func registerAuth(router fiber.Router) {
	router.Route("/discord", discord.Register)
}
