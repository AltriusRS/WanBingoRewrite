package users

import (
	"wanshow-bingo/handlers/users/me"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/users", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	// Protected routes - require authentication
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware)

	protected.Get("/me", me.Get)
	protected.Put("/me", me.Put)

	// Public routes - no authentication required
	router.Get("/", GetAll)
	router.Get("/:identifier", GetByIdentifier)
}
