package timers

import (
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/timers", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	router.Get("/", GetTimers)
	router.Get("/:id", GetTimerByID)
	router.Post("/", CreateTimer)
	router.Put("/:id", UpdateTimer)
	router.Delete("/:id", DeleteTimer)
	router.Post("/:id/start", StartTimer)
	router.Post("/:id/stop", StopTimer)
}
