package tilerouter

import (
	"context"
	"log"
	"wanshow-bingo/db"
	"wanshow-bingo/middleware"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

func init() {
	utils.RegisterRouter("/tiles", BuildRouter)
}

func BuildRouter(router fiber.Router) {
	log.Printf("Registering tiles routes")
	router.Get("/", Get)
	router.Get("/show", GetShowTiles)
	router.Get("/me", middleware.AuthMiddleware, GetMyBoard)
	router.Post("/me/regenerate", middleware.AuthMiddleware, RegenerateMyBoard)
	router.Get("/anonymous", GetAnonymousBoard)
	router.Post("/anonymous/regenerate", RegenerateAnonymousBoard)
	router.Get("/confirmed", GetConfirmedTiles)
	router.Post("/confirmations", middleware.AuthMiddleware, ConfirmTile)
	router.Post("/win", middleware.AuthMiddleware, RecordWin)
	router.Get("/:tile_id", GetTileByID)
	log.Printf("Tiles routes registered")
}

func GetConfirmedTiles(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		log.Printf("GetConfirmedTiles: Failed to get latest show: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get latest show", 500))
	}

	log.Printf("GetConfirmedTiles: Got latest show %s", latestShow.ID)

	// Get confirmed tiles for this show
	confirmed, err := db.GetTileConfirmationsForShow(ctx, latestShow.ID)
	if err != nil {
		log.Printf("GetConfirmedTiles: Failed to get confirmations for show %s: %v", latestShow.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewApiError("Failed to get confirmed tiles", 500))
	}

	log.Printf("GetConfirmedTiles: Got %d confirmations for show %s", len(confirmed), latestShow.ID)

	// Return just the tile IDs
	tileIDs := make([]string, len(confirmed))
	for i, conf := range confirmed {
		tileIDs[i] = conf.TileID
	}

	log.Printf("GetConfirmedTiles: Returning %d confirmed tiles for show %s", len(tileIDs), latestShow.ID)
	return c.JSON(tileIDs)
}
