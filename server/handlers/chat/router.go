package chat

import (
	_ "wanshow-bingo/handlers/chat/stream"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/chat", ChatRouter)
}

func ChatRouter(router fiber.Router) {
	router.Post("/", middleware.RequiredAuthMiddleware, Post)
	router.Post("/s", PostSystem)
}
