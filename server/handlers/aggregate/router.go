package aggregaterouter

import (
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/aggregate", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	router.Get("/", Get)
	router.Post("/", Post)
}
