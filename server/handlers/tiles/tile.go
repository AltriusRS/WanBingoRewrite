package tilerouter

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// GetTileByID retrieves a single tile by its ID
func GetTileByID(c *fiber.Ctx) error {
	tileID := c.Params("tile_id")
	if tileID == "" {
		return utils.NewApiError("Tile ID is required", 0x0601).AsResponse(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tile, err := db.GetTileByID(ctx, tileID)
	if err != nil {
		log.Printf("failed to get tile %s: %v", tileID, err)
		return utils.NewApiError("Tile not found", 0x0602).AsResponse(c)
	}

	return c.JSON(tile)
}
