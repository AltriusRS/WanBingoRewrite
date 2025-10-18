package admin

import (
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/admin", AdminRouter)
}

func AdminRouter(router fiber.Router) {
	// All admin routes require host permission
	adminMiddleware := middleware.AuthMiddleware
	hostMiddleware := middleware.RequirePermissionMiddleware("can_host")

	router.Use(adminMiddleware, hostMiddleware)
}
