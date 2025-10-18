package users

import (
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/admin/users", UsersRouter)
}

func UsersRouter(router fiber.Router) {
	router.Get("/:id/permissions", GetUserPermissions)
	router.Put("/:id/permissions", UpdateUserPermissions)
}
