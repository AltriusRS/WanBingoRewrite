package tilerouter

import (
	"context"
	"log"
	"time"
	"wanshow-bingo/db"
	"wanshow-bingo/utils"

	"github.com/gofiber/fiber/v2"
)

// GetShowTiles returns tile IDs for the most recent show with pagination
func GetShowTiles(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get pagination params
	limit := c.QueryInt("limit", 0) // 0 means no limit
	offset := c.QueryInt("offset", 0)

	// Get the latest show
	latestShow, err := db.GetLatestShow(ctx)
	if err != nil {
		log.Printf("failed to get latest show: %v", err)
		return utils.NewApiError("Failed to get latest show", 0x0301).AsResponse(c)
	}

	// Ensure show has enough tiles
	showTiles, err := db.GetShowTiles(ctx, latestShow.ID)
	if err != nil {
		log.Printf("failed to get show tiles: %v", err)
		return utils.NewApiError("Failed to get show tiles", 0x0302).AsResponse(c)
	}
	if len(showTiles) < db.GetTilesPerShow() {
		err = db.PopulateShowTilesWithRandom(ctx, latestShow.ID)
		if err != nil {
			log.Printf("failed to populate show tiles: %v", err)
			return utils.NewApiError("Failed to populate show tiles", 0x0304).AsResponse(c)
		}
	}

	// Get tile IDs for this show
	tileIDs, err := db.GetShowTileIDs(ctx, latestShow.ID)
	if err != nil {
		log.Printf("failed to get show tile IDs: %v", err)
		return utils.NewApiError("Failed to get show tiles", 0x0302).AsResponse(c)
	}

	// Apply pagination
	total := len(tileIDs)
	if limit > 0 {
		end := offset + limit
		if offset > len(tileIDs) {
			offset = len(tileIDs)
		}
		if end > len(tileIDs) {
			end = len(tileIDs)
		}
		tileIDs = tileIDs[offset:end]
	}

	return c.JSON(fiber.Map{
		"show_id":  latestShow.ID,
		"tile_ids": tileIDs,
		"total":    total,
		"offset":   offset,
		"limit":    limit,
	})
}
