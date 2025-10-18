package host

import (
	"wanshow-bingo/handlers/host/stream"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/host", HostRouter)
}

func HostRouter(router fiber.Router) {
	router.Get("/stream", middleware.AuthMiddleware, middleware.RequirePermissionMiddleware("can_host"), stream.Get)
}
