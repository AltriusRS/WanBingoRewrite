package tilerouter

import (
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/tiles", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	router.Get("/", Get)
}
