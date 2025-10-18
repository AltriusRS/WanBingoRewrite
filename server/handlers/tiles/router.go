package tilerouter

import (
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/tiles", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	router.Get("/", Get)
	router.Get("/show", GetShowTiles)
	router.Get("/me", middleware.AuthMiddleware, GetMyBoard)
	router.Post("/me/regenerate", middleware.AuthMiddleware, RegenerateMyBoard)
	router.Get("/anonymous", GetAnonymousBoard)
	router.Post("/anonymous/regenerate", RegenerateAnonymousBoard)
	router.Get("/:tile_id", GetTileByID)
	router.Post("/confirmations", middleware.AuthMiddleware, ConfirmTile)
}
