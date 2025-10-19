package timers

import (
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/timers", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	router.Get("/", GetTimers)
	router.Get("/:id", GetTimerByID)

	// Protected routes that require authentication
	auth := router.Group("", middleware.AuthMiddleware)
	auth.Post("/", CreateTimer)
	auth.Put("/:id", UpdateTimer)
	auth.Delete("/:id", DeleteTimer)
	auth.Post("/:id/start", StartTimer)
	auth.Post("/:id/stop", StopTimer)
}
