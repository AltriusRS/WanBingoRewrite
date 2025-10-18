package auth

import (
	"wanshow-bingo/handlers/auth/discord"

	"github.com/gofiber/fiber/v2"
)

func init() {

}

func registerAuth(router fiber.Router) {
	router.Route("/discord", discord.Register)
}
