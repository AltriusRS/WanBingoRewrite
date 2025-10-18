package tilerouter

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// GetShowTiles returns all tile IDs for the most recent show
func GetShowTiles(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		log.Printf("failed to get latest show: %v", err)
		return utils.NewApiError("Failed to get latest show", 0x0301).AsResponse(c)
	}

	// Get tile IDs for this show
	tileIDs, err := db.GetShowTileIDs(ctx, latestShow.ID)
	if err != nil {
		log.Printf("failed to get show tiles: %v", err)
		return utils.NewApiError("Failed to get show tiles", 0x0302).AsResponse(c)
	}

	// If no show-specific tiles, return all tiles
	if len(tileIDs) == 0 {
		allTiles, err := db.GetAllTiles(ctx)
		if err != nil {
			log.Printf("failed to get all tiles: %v", err)
			return utils.NewApiError("Failed to get tiles", 0x0303).AsResponse(c)
		}
		tileIDs = make([]string, len(allTiles))
		for i, tile := range allTiles {
			tileIDs[i] = tile.ID
		}
	}

	return c.JSON(fiber.Map{
		"show_id":  latestShow.ID,
		"tile_ids": tileIDs,
	})
}
