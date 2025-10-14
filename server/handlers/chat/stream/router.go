package stream

import (
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/chat/stream", StreamRouter)
}

func StreamRouter(router fiber.Router) {
	router.Get("/", middleware.OptionalAuthMiddleware, Get)
}
