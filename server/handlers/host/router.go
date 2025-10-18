package host

import (
	"wanshow-bingo/handlers/host/stream"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/host", HostRouter)
}

func HostRouter(router fiber.Router) {
	router.Get("/stream", stream.Get)
	router.Post("/test-message", PostTestMessage)
}
